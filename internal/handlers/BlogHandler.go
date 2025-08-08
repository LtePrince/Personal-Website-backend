package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/LtePrince/Personal-Website-backend/internal/utils"
)

// Handler 解析请求并调用相应的处理函数
func Handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/Blog":
		log.Printf("\033[32m[Log]\033[0mBlogHandler")
		blogHandler(w, r)
	case r.URL.Path == "/api/LatestBlog":
		log.Printf("\033[32m[Log]\033[0mLatestBlogHandler")
		LatestBlogHandler(w, r)
	case len(r.URL.Path) >= len("/api/BlogDetail") && r.URL.Path[:len("/api/BlogDetail")] == "/api/BlogDetail":
		log.Printf("\033[32m[Log]\033[0mBlogContentHandler")
		BlogContentHandler(w, r)
	default:
		log.Printf("\033[31m[Log]\033[0mNot Found")
		http.NotFound(w, r)
	}
}

// BlogHandler 处理 /pages/Blog 请求
func blogHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求方法
	method := r.Method

	// 获取请求的 URL 路径
	path := r.URL.Path

	// 获取请求头部信息
	userAgent := r.Header.Get("User-Agent")

	// 获取请求体内容（假设请求体是文本）
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	// 输出请求信息
	log.Printf("\033[32m[Log]\033[0m------Method: %s\n", method)
	log.Printf("\033[32m[Log]\033[0m------Path: %s\n", path)
	log.Printf("\033[32m[Log]\033[0m------User-Agent: %s\n", userAgent)
	log.Printf("\033[32m[Log]\033[0m------Body: %s\n", string(body))

	// 获取博客标题和摘要
	blogs, err := utils.GetBlogInfo()
	if err != nil {
		http.Error(w, "Error fetching blog titles and summaries", http.StatusInternalServerError)
		log.Printf("Error fetching blog titles and summaries: %v", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func LatestBlogHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求方法
	method := r.Method

	// 获取请求的 URL 路径
	path := r.URL.Path

	// 获取请求头部信息
	userAgent := r.Header.Get("User-Agent")

	// 获取请求体内容（假设请求体是文本）
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	// 输出请求信息
	log.Printf("\033[32m[Log]\033[0m------Method: %s\n", method)
	log.Printf("\033[32m[Log]\033[0m------Path: %s\n", path)
	log.Printf("\033[32m[Log]\033[0m------User-Agent: %s\n", userAgent)
	log.Printf("\033[32m[Log]\033[0m------Body: %s\n", string(body))

	// 获取最新博客内容
	latestBlog, err := utils.GetLatestBlog()
	if err != nil {
		http.Error(w, "Error fetching latest blog", http.StatusInternalServerError)
		log.Printf("Error fetching latest blog: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latestBlog)
}

func BlogContentHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求方法
	method := r.Method

	// 获取请求的 URL 路径
	path := r.URL.Path

	// 获取请求头部信息
	userAgent := r.Header.Get("User-Agent")

	// 获取请求体内容（假设请求体是文本）
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	// 输出请求信息
	log.Printf("\033[32m[Log]\033[0m------Method: %s\n", method)
	log.Printf("\033[32m[Log]\033[0m------Path: %s\n", path)
	log.Printf("\033[32m[Log]\033[0m------User-Agent: %s\n", userAgent)
	log.Printf("\033[32m[Log]\033[0m------Body: %s\n", string(body))

	// 从path中获取博客ID
	// 假设路径格式为 /pages/BlogDetail?id=xxx
	// 解析查询参数
	query := r.URL.Query()
	id := query.Get("id")
	if id == "" {
		http.Error(w, "Missing blog ID", http.StatusBadRequest)
		log.Println("Missing blog ID")
		return
	}

	log.Printf("\033[32m[Log]\033[0m------ID: %s\n", id)

	// 字符串转int
	blogID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		log.Println("Invalid blog ID:", id)
		return
	}

	// 获取博客内容
	blogContent, err := utils.GetBlogContentByID(blogID)
	if err != nil {
		http.Error(w, "Error fetching blog content", http.StatusInternalServerError)
		log.Printf("Error fetching blog content: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogContent)
}
