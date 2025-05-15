package main

import (
	"fmt"
	"net/http"

	"github.com/LtePrince/Personal-Website-backend/internal/handlers"
)

func main() {
	// 静态资源服务，访问 /static/xxx.jpg 实际读取 static 目录下的文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/home/adolph/workspace/Personal-website/blogs/static"))))

	http.HandleFunc("/", handlers.Handler)
	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
