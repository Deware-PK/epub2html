package cleaner

import (
	"bytes"
	"epub2html/internal/config"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type HTMLCleaner struct {
	Cfg *config.Config
}

func NewHTMLCleaner(cfg *config.Config) *HTMLCleaner {
	return &HTMLCleaner{Cfg: cfg}
}

func (c *HTMLCleaner) Clean(input string) string {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return input
	}

	var buf bytes.Buffer
	c.traverse(doc, &buf)

	return buf.String()
}

func (c *HTMLCleaner) traverse(n *html.Node, buf *bytes.Buffer) {
	switch n.Type {
	case html.ElementNode:
		if c.Cfg.AllowedTags[n.Data] {
			buf.WriteString("<" + n.Data + ">")
		}
	case html.TextNode:
		buf.WriteString(n.Data)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.traverse(child, buf)
	}

	if n.Type == html.ElementNode && c.Cfg.AllowedTags[n.Data] {
		buf.WriteString("</" + n.Data + ">")
	}
}

func (c *HTMLCleaner) WrapHTML(title string, content string, i int, total int) string {
	prevLink := "<span></span>"
	if i > 0 {
		prevLink = fmt.Sprintf("<a href='chapter_%03d.html'>&laquo; Prev</a>", i-1)
	}

	nextLink := "<span></span>"
	if i >= 0 && i < total-1 {
		nextLink = fmt.Sprintf("<a href='chapter_%03d.html'>Next &raquo;</a>", i+1)
	}

	safeHead := `
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        :root { --bg-color: #f4f4f9; --container-bg: #ffffff; --text-color: #333333; --link-color: #007bff; --border-color: #dddddd; --quote-color: #666666; --shadow: 0 2px 5px rgba(0,0,0,0.05); }
        [data-theme="dark"] { --bg-color: #121212; --container-bg: #1e1e1e; --text-color: #e0e0e0; --link-color: #66b3ff; --border-color: #333333; --quote-color: #aaaaaa; --shadow: 0 2px 5px rgba(0,0,0,0.5); }
        body { font-family: 'Sarabun', sans-serif; line-height: 1.8; color: var(--text-color); background-color: var(--bg-color); margin: 0; padding: 20px; transition: background-color 0.3s, color 0.3s; }
        .container { max-width: 800px; margin: 0 auto; background-color: var(--container-bg); padding: 40px; border-radius: 8px; box-shadow: var(--shadow); }
        p { margin-bottom: 1.5em; text-align: justify; }
        .nav-links { margin-top: 40px; padding-top: 20px; border-top: 1px solid var(--border-color); display: flex; justify-content: space-between; flex-wrap: wrap; gap: 10px; }
        a { text-decoration: none; color: var(--link-color); font-weight: bold; }
        .theme-toggle { position: fixed; top: 20px; right: 20px; background: var(--container-bg); border: 1px solid var(--border-color); color: var(--text-color); padding: 8px 12px; border-radius: 20px; cursor: pointer; z-index: 1000; }
    </style>
    <script>
        function toggleTheme() {
            const currentTheme = document.documentElement.getAttribute("data-theme");
            const newTheme = currentTheme === "dark" ? "light" : "dark";
            document.documentElement.setAttribute("data-theme", newTheme);
            localStorage.setItem("theme", newTheme);
        }
        (function() {
            const savedTheme = localStorage.getItem("theme");
            if (savedTheme) document.documentElement.setAttribute("data-theme", savedTheme);
        })();
    </script>`

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    %s
</head>
<body>
    <button class="theme-toggle" onclick="toggleTheme()">🌓 Theme</button>
    <div class="container">
        %s
        <div class='nav-links'>
            %s
            <a href='index.html'>Table of Contents</a>
            %s
        </div>
    </div>
</body>
</html>`, title, safeHead, content, prevLink, nextLink)
}
