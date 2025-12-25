package main

import (
	"bufio"
	"epub2html/internal/cleaner"
	"epub2html/internal/config"
	"epub2html/internal/processor"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	Version = "v1.0.0"
)

func main() {

	printBanner()

	reader := bufio.NewReader(os.Stdin)
	cfg := config.NewDefaultConfig()

	fmt.Print("Enter EPUB file path (e.g. C:\\books\\test.epub): ")
	inputPath, _ := reader.ReadString('\n')
	cfg.EpubPath = cleanInput(inputPath)

	if cfg.EpubPath == "" {
		log.Fatal("Error: EPUB path cannot be empty.")
	}

	fmt.Print("Enter Output directory (e.g. C:\\books\\output): ")
	inputOut, _ := reader.ReadString('\n')
	cfg.OutputDir = cleanInput(inputOut)

	if cfg.OutputDir == "" {
		cfg.OutputDir = "clean_output"
		fmt.Println("No output folder specified. Using default: 'clean_output'")
	}

	cfg.Workers = runtime.NumCPU()

	c := cleaner.NewHTMLCleaner(cfg)
	proc := &processor.EpubProcessor{Cleaner: c}

	fmt.Printf("\nStarting Process...\nFile: %s\nOutput: %s\nWorkers: %d\n\n", cfg.EpubPath, cfg.OutputDir, cfg.Workers)

	err := proc.Process(cfg.EpubPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone! Check your output folder.")

	fmt.Println("Press 'Enter' to exit...")
	reader.ReadString('\n')
}

func cleanInput(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\"'")
	return text
}

func printBanner() {
	art := `
 _______  _______  __   __  _______  _______  __   __  _______  __   __  ___      
|       ||       ||  | |  ||   _   ||       ||  | |  ||       ||  |_|  ||   |    
|    ___||    _  ||  | |  ||  |_|  ||____   ||  |_|  ||_     _||       ||   |    
|   |___ |   |_| ||  |_|  ||       | ____|  ||       |  |   |  |       ||   |    
|    ___||    ___||       ||  _   | | ______||       |  |   |  |       ||   |___ 
|   |___ |   |    |       || |_|   || |_____ |   _   |  |   |  | ||_|| ||       |
|_______||___|    |_______||_______||_______||__| |__|  |___|  |_|   |_||_______|
`
	fmt.Println("====================================================================================")
	fmt.Println(art)
	fmt.Printf("Version: %s\n", Version)
	fmt.Println("====================================================================================")
	fmt.Println()
}
