
# EPUB to HTML Cleaner

> A blazingly fast, concurrent EPUB to HTML converter written in Go.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

**EPUB to HTML Cleaner** is a high-performance utility designed to extract and sanitize EPUB files into clean, readable HTML documents. Built with Go's powerful concurrency model, it processes large books in seconds and generates a fully responsive, mobile-friendly reading interface with Dark Mode support.
# Features
-   **Blazingly Fast:** Utilizes Go Goroutines and Worker Pools to process chapters concurrently.
    
-   **Responsive Design:** Generated HTML is optimized for all devices (Mobile, Tablet, Desktop).
    
-   **Theme Support:** Built-in Light/Dark mode toggle with persistence (saves your preference).
    
-   **Smart Cleaning:** Removes unnecessary tags while preserving the core structure of the book.
-    **Auto-Indexing:** Automatically generates a Table of Contents (`index.html`) and navigation links (Next/Prev) for every chapter.
-   **Privacy Focused:** All processing is done locally on your machine. No data is sent to any server.
# Installation & Usage
### Option 1: Run from Source

If you have Go installed on your machine:
1. Clone the repository:
```
git clone https://github.com/Deware-PK/epub2html.git
cd epub2html
```
2. Install dependencies:
```
go mod tidy
```
3. Run the program:
```
go run main.go
```

### Option 2: Build Binary

You can build a standalone executable file:
```
# For Windows 
go build -o epub2html.exe main.go 

# For Mac/Linux 
go build -o epub2html main.go
```

# Disclaimer

**Please Read Carefully:**

This tool is developed for **personal use and educational purposes only**. It is designed to format DRM-free EPUB files to improve the reading experience on web browsers.

-   **No DRM Removal:** This tool **does not** and **cannot** bypass or remove Digital Rights Management (DRM) protection. It only works on DRM-free files.
    
-   **Copyright Respect:** The author (**Deware**) does not endorse piracy or the unauthorized distribution of copyrighted material. Please ensure you own the rights to any content you process with this tool.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/Deware-PK/epub2html/blob/main/LICENSE) file for details.
