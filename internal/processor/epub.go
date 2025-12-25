package processor

import (
	"archive/zip"
	"epub2html/internal/cleaner"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type EpubProcessor struct {
	Cleaner *cleaner.HTMLCleaner
}

type ChapterInfo struct {
	Index    int
	FileName string
	HTMLPath string
}

func (p *EpubProcessor) Process(epubPath string) error {
	// Create output directory if not exists
	if err := os.MkdirAll(p.Cleaner.Cfg.OutputDir, os.ModePerm); err != nil {
		return err
	}

	// Open the EPUB file (which is a ZIP archive)
	reader, err := zip.OpenReader(epubPath)
	if err != nil {
		return fmt.Errorf("couldn't open epub: %v", err)
	}
	defer reader.Close()

	// Filter only HTML files and sort by file name (very important!)
	// Because zip does not guarantee order, if not sorted, chapter 10 may come before chapter 1
	var targetFiles []*zip.File
	for _, f := range reader.File {
		if isHTML(f.Name) {
			targetFiles = append(targetFiles, f)
		}
	}

	// Sort by file name A-Z
	sort.Slice(targetFiles, func(i, j int) bool {
		return targetFiles[i].Name < targetFiles[j].Name
	})

	totalFiles := len(targetFiles)
	fmt.Printf("Found %d HTML files. Processing...\n", totalFiles)

	// Setup Worker Pool
	jobs := make(chan struct {
		f *zip.File
		i int
	}, totalFiles)
	
	results := make(chan ChapterInfo, totalFiles)
	
	var wg sync.WaitGroup

	workers := p.Cleaner.Cfg.Workers
	if workers <= 0 { workers = 1 }

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// Read
				content, _ := readFileContent(job.f)
				
				// Create
				cleanContent := p.Cleaner.Clean(content)

				// Beautify with WrapHTML
				title := fmt.Sprintf("Chapter %d", job.i)
				finalHTML := p.Cleaner.WrapHTML(title, cleanContent, job.i, totalFiles)

				// Rename output file
				outName := fmt.Sprintf("chapter_%03d.html", job.i)
				outPath := filepath.Join(p.Cleaner.Cfg.OutputDir, outName)

				// Save to output directory
				err := os.WriteFile(outPath, []byte(finalHTML), 0644)
				if err == nil {
					fmt.Printf("[Worker %d] Saved: %s\n", workerID, outName)
					results <- ChapterInfo{Index: job.i, FileName: job.f.Name, HTMLPath: outName}
				}
			}
		}(w)
	}

	for i, f := range targetFiles {
		jobs <- struct {
			f *zip.File
			i int
		}{f, i}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var chapters []ChapterInfo
	for res := range results {
		chapters = append(chapters, res)
	}

	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Index < chapters[j].Index
	})

	return p.generateIndex(chapters)
}

func (p *EpubProcessor) generateIndex(chapters []ChapterInfo) error {
	var listItems string
	for _, ch := range chapters {
		listItems += fmt.Sprintf("<li><a href='%s'>Chapter %d : %s</a></li>", ch.HTMLPath, ch.Index, ch.FileName)
	}
	
	content := fmt.Sprintf("<h1>Table of Contents</h1><ul>%s</ul>", listItems)
	indexHTML := p.Cleaner.WrapHTML("Table of Contents", content, -1, -1) // -1 คือไม่มีปุ่ม prev/next

	return os.WriteFile(filepath.Join(p.Cleaner.Cfg.OutputDir, "index.html"), []byte(indexHTML), 0644)
}

func isHTML(name string) bool {
	return len(name) > 5 && (name[len(name)-5:] == ".html" || name[len(name)-6:] == ".xhtml")
}

func readFileContent(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil { return "", err }
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil { return "", err }
	return string(content), nil
}