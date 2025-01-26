package api

// this file defines the Restful api of Blog
type BlogParams struct {
	// BlogId is the id of the blog
	BlogId string `json:"blogId"`
}

type BlogResponse struct {
	// Title is the title of the blog
	Title string `json:"title"`
	// Content is the content of the blog
	Summary string `json:"summary"`
}

type Error struct {
	// Code is the error code
	Code int `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}
