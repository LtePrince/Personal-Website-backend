package test

import (
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func get(t *testing.T, path string) []byte {
	resp, err := http.Get(baseURL + path)
	if err != nil {
		t.Fatalf("GET %s error: %v", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status=%d", path, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read %s body error: %v", path, err)
	}
	return b
}

func TestBlogList(t *testing.T) {
	b := get(t, "/api/Blog")
	if len(b) == 0 {
		t.Fatalf("empty response")
	}
}

func TestBlogDetail(t *testing.T) {
	b := get(t, "/api/BlogDetail?id=1")
	if len(b) == 0 {
		t.Fatalf("empty response")
	}
}

func TestLatestBlog(t *testing.T) {
	b := get(t, "/api/LatestBlog")
	if len(b) == 0 {
		t.Fatalf("empty response")
	}
}

func TestStaticImage(t *testing.T) {
	b := get(t, "/static/Image_1730389545752.jpg")
	if len(b) < 10 { // 原逻辑打印前10字节
		t.Fatalf("image too small, got %d bytes", len(b))
	}
}
