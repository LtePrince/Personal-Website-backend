package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	case r.URL.Path == "/api/Weather":
		log.Printf("\033[32m[Log]\033[0mWeatherHandler")
		WeatherHandler(w, r)
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

// WeatherHandler 根据客户端 IP 返回所在城市天气与空气质量
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// 基本日志
	method := r.Method
	path := r.URL.Path
	userAgent := r.Header.Get("User-Agent")
	log.Printf("\033[32m[Log]\033[0m------Method: %s\n", method)
	log.Printf("\033[32m[Log]\033[0m------Path: %s\n", path)
	log.Printf("\033[32m[Log]\033[0m------User-Agent: %s\n", userAgent)

	// 提取客户端 IP（支持代理）
	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			// 取第一个 IP
			if comma := strings.Index(ip, ","); comma > 0 {
				ip = strings.TrimSpace(ip[:comma])
			}
		}
	}
	if ip == "" {
		// RemoteAddr 可能包含端口
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}

	// 允许本地开发时通过查询参数覆盖经纬度
	q := r.URL.Query()
	latParam := q.Get("lat")
	lonParam := q.Get("lon")

	// 默认经纬度：悉尼
	lat := -33.8688
	lon := 151.2093
	city := "Sydney"
	region := "NSW"
	countryCode := "AU"

	// 如果提供了经纬度参数，优先使用
	if latParam != "" && lonParam != "" {
		if v, err := strconv.ParseFloat(latParam, 64); err == nil {
			lat = v
		}
		if v, err := strconv.ParseFloat(lonParam, 64); err == nil {
			lon = v
		}
	} else {
		// 否则尝试 IP 定位
		if !utils.IsPrivateOrLoopbackIP(ip) {
			if info, err := utils.LookupIPLocation(ip); err == nil {
				if info.Latitude != 0 || info.Longitude != 0 {
					lat = info.Latitude
					lon = info.Longitude
				}
				if info.City != "" {
					city = info.City
				}
				if info.Region != "" {
					region = info.Region
				}
				if info.CountryCode != "" {
					countryCode = info.CountryCode
				}
			}
		}
	}

	// 拉取天气与空气质量
	weather, err := utils.FetchWeatherAndAQI(lat, lon)
	if err != nil {
		log.Printf("Error FetchWeatherAndAQI: %v", err)
		// 不中断，返回占位
	}

	// 输出 JSON
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	type resp struct {
		City         string `json:"city"`
		Region       string `json:"region"`
		CountryCode  string `json:"countryCode"`
		Location     string `json:"location"`
		TemperatureC *int   `json:"temperatureC"`
		WindSpeedKmh *int   `json:"windSpeedKmh"`
		WindLevel    *int   `json:"windLevel"`
		Humidity     *int   `json:"humidity"`
		AQIUS        *int   `json:"aqiUS"`
		WeatherText  string `json:"weatherText"`
		UpdatedAt    string `json:"updatedAt"`
	}

	var tPtr, wPtr, wLvlPtr, hPtr, aqiPtr *int
	if weather.TempC != nil {
		t := int(*weather.TempC)
		tPtr = &t
	}
	if weather.WindSpeedKmh != nil {
		wv := int(*weather.WindSpeedKmh)
		wPtr = &wv
	}
	if weather.Humidity != nil {
		hv := int(*weather.Humidity)
		hPtr = &hv
	}
	if weather.AQIUS != nil {
		av := int(*weather.AQIUS)
		aqiPtr = &av
	}

	// 计算风级（蒲福风级近似）
	if wPtr != nil {
		v := *wPtr
		lvl := 0
		switch {
		case v < 1:
			lvl = 0
		case v <= 5:
			lvl = 1
		case v <= 11:
			lvl = 2
		case v <= 19:
			lvl = 3
		case v <= 28:
			lvl = 4
		case v <= 38:
			lvl = 5
		case v <= 49:
			lvl = 6
		case v <= 61:
			lvl = 7
		case v <= 74:
			lvl = 8
		case v <= 88:
			lvl = 9
		case v <= 102:
			lvl = 10
		case v <= 117:
			lvl = 11
		default:
			lvl = 12
		}
		wLvlPtr = &lvl
	}

	payload := resp{
		City:        city,
		Region:      region,
		CountryCode: countryCode,
		Location: city + func() string {
			if region != "" {
				return " · " + region
			}
			return ""
		}() + func() string {
			if countryCode != "" {
				return " · " + countryCode
			}
			return ""
		}(),
		TemperatureC: tPtr,
		WindSpeedKmh: wPtr,
		WindLevel:    wLvlPtr,
		Humidity:     hPtr,
		AQIUS:        aqiPtr,
		WeatherText:  weather.WeatherText,
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(payload)
}
