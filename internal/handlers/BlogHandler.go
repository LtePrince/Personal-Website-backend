package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LtePrince/Personal-Website-backend/internal/utils"
)

// Handler 解析请求并调用相应的处理函数
func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/pages/Blog":
		//fmt.Println("BlogHandler")
		blogHandler(w, r)
	default:
		//fmt.Println("Not Found")
		http.NotFound(w, r)
	}
}

// BlogHandler 处理 /pages/Blog 请求
func blogHandler(w http.ResponseWriter, r *http.Request) {

	blogs, err := utils.GetBlogTitlesAndSummaries()
	if err != nil {
		log.Fatalf("Error fetching blog titles and summaries: %v", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}
