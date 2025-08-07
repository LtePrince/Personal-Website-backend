package main

import (
	"fmt"
	"io"
	"net/http"
)

func test1() {
	url := "http://154.37.213.201:8080/pages/Blog"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Response from server: %s\n", body)
}

func test2() {
	url := "http://154.37.213.201:8080/pages/BlogDetail?id=1"
	// url := "http://101.132.86.173:8080/pages/BlogDetail?id=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Response from server: %s\n", body)
}

func test3() (body []byte) {
	url := "http://154.37.213.201:8080/static/Image_1730389545752.jpg"
	// url := "http://101.132.86.173:8080/static/Image_1730389545752.jpg"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return nil
	}
	fmt.Printf("Response Success. %s\n", body[:10])
	return body
}

func test4() {
	url := "http://154.37.213.201:8080/pages/LatestBlog"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Response from /pages/LatestBlog: %s\n", body)
}

func main() {
	fmt.Println("Starting tests...")
	fmt.Println("Test1")
	test1()
	fmt.Println("Test2")
	test2()
	fmt.Println("Test3")
	test3()
	fmt.Println("Test4")
	test4()
	fmt.Println("Tests completed.")
}
