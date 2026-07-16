package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/**/*
var staticFS embed.FS

//go:embed static/robots.txt
var robotsTxt embed.FS

func main() {
	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8999"
	}
	log.Printf("🏥 仁和医院门户网站已启动: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// parseTemplates 解析所有页面模板
func parseTemplates() (*template.Template, error) {
	return template.New("").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFS(templatesFS, "templates/*.html")
}

// setupRouter 配置路由
func setupRouter(tmpl *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	robotsSub, _ := fs.Sub(robotsTxt, "static")
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, robotsSub, "robots.txt")
	})

	pages := []string{"index", "departments", "doctors", "about", "contact"}
	for _, page := range pages {
		page := page
		mux.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			renderPage(w, tmpl, page)
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if len(r.URL.Path) > 7 && r.URL.Path[:8] == "/backup/" {
				serveBackup(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}
		renderPage(w, tmpl, "index")
	})

	return mux
}

func renderPage(w http.ResponseWriter, tmpl *template.Template, page string) {
	data := map[string]interface{}{
		"Title": getPageTitle(page),
		"Page":  page,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("渲染页面 %s 失败: %v", page, err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
	}
}

func serveBackup(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(".", "backup", filepath.Base(r.URL.Path))

	// 防止路径遍历
	if filepath.Base(r.URL.Path) != r.URL.Path[len("/backup/"):] {
		http.NotFound(w, r)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(r.URL.Path))
	http.ServeFile(w, r, filePath)
}

func getPageTitle(page string) string {
	titles := map[string]string{
		"index":       "首页",
		"departments": "科室介绍",
		"doctors":     "专家团队",
		"about":       "医院简介",
		"contact":     "联系我们",
	}
	if t, ok := titles[page]; ok {
		return t + " - 仁和医院"
	}
	return "仁和医院"
}
