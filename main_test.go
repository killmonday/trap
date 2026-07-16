package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllPagesReturn200(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	pages := []string{"/", "/index", "/departments", "/doctors", "/about", "/contact"}
	for _, page := range pages {
		resp, err := http.Get(server.URL + page)
		if err != nil {
			t.Errorf("请求 %s 失败: %v", page, err)
			continue
		}
		if resp.StatusCode != 200 {
			t.Errorf("页面 %s 返回 %d，期望 200", page, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestRobotsTxt(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/robots.txt")
	if err != nil {
		t.Fatalf("请求 robots.txt 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("robots.txt 返回 %d，期望 200", resp.StatusCode)
	}
}

func TestStaticFiles(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	staticFiles := []string{
		"/static/css/style.css",
		"/static/js/main.js",
	}
	for _, f := range staticFiles {
		resp, err := http.Get(server.URL + f)
		if err != nil {
			t.Errorf("请求 %s 失败: %v", f, err)
			continue
		}
		if resp.StatusCode != 200 {
			t.Errorf("静态文件 %s 返回 %d，期望 200", f, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestHtmlContainsBackupHint(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("请求首页失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	// 验证 HTML 中包含备份路径的 hint
	if !strings.Contains(body, "backup") || !strings.Contains(body, "wwwroot.zip") {
		t.Error("首页 HTML 中未找到备份路径 hint（data-backup-path）")
	}
}

func TestJsContainsBackupHint(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/static/js/main.js")
	if err != nil {
		t.Fatalf("请求 main.js 失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	if !strings.Contains(body, "backup") || !strings.Contains(body, "wwwroot.zip") {
		t.Error("main.js 中未找到备份路径 hint")
	}
}

func TestNonExistentRouteReturns404(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/nonexistent-page")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Errorf("不存在页面返回 %d，期望 404", resp.StatusCode)
	}
}

func TestPathTraversalBlocked(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("解析模板失败: %v", err)
	}

	mux := setupRouter(tmpl)
	server := httptest.NewServer(mux)
	defer server.Close()

	// 路径遍历攻击应被阻止
	resp, err := http.Get(server.URL + "/backup/../../../etc/passwd")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Errorf("路径遍历攻击返回 %d，期望 404", resp.StatusCode)
	}
}
