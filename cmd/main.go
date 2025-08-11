package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LtePrince/Personal-Website-backend/internal/handlers"
)

func main() {
	// 环境变量读取
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		// 默认线上静态资源路径（保持与现状一致）
		staticDir = "/www/wwwroot/Personal-Blog-db/static"
	}

	// 静态资源服务，访问 /static/xxx.jpg 实际读取 static 目录下的文件
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/home/adolph/workspace/Personal-website/blogs/static"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/", handlers.Handler)
	fmt.Printf("Server is listening on port %s (static: %s) ...\n", port, staticDir)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
