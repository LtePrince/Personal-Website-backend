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

// ---- 通用辅助 ----
// getClientIP 提取客户端真实 IP（支持常见代理头），失败时回退 RemoteAddr
func getClientIP(r *http.Request) string {
	// 优先 Cloudflare
	if ip := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); ip != "" {
		return ip
	}
	// 再看 X-Forwarded-For 第一段
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			p := strings.TrimSpace(parts[0])
			if p != "" {
				return p
			}
		}
	}
	// 直接 RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

// beaufortLevel 近似计算蒲福风级
func beaufortLevel(kmh int) int {
	switch {
	case kmh < 1:
		return 0
	case kmh <= 5:
		return 1
	case kmh <= 11:
		return 2
	case kmh <= 19:
		return 3
	case kmh <= 28:
		return 4
	case kmh <= 38:
		return 5
	case kmh <= 49:
		return 6
	case kmh <= 61:
		return 7
	case kmh <= 74:
		return 8
	case kmh <= 88:
		return 9
	case kmh <= 102:
		return 10
	case kmh <= 117:
		return 11
	default:
		return 12
	}
}

// buildLocation 拼接展示字符串（city · region · countryCode）
func buildLocation(city, region, countryCode string) string {
	var segs []string
	if city != "" {
		segs = append(segs, city)
	}
	if region != "" {
		segs = append(segs, region)
	}
	if countryCode != "" {
		segs = append(segs, countryCode)
	}
	return strings.Join(segs, " · ")
}

// writeJSONHeaders 统一写基础 JSON 响应头
func writeJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

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
		log.Printf("\033[31m[Log]\033[0mNot Found: %s", r.URL.Path)
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
	// 日志
	log.Printf("\033[32m[Log]\033[0m------Method: %s\n", r.Method)
	log.Printf("\033[32m[Log]\033[0m------Path: %s\n", r.URL.Path)
	log.Printf("\033[32m[Log]\033[0m------User-Agent: %s\n", r.Header.Get("User-Agent"))

	writeJSONHeaders(w)

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

	ip := getClientIP(r)
	if utils.IsPrivateOrLoopbackIP(ip) {
		json.NewEncoder(w).Encode(resp{
			City:        "localhost",
			Location:    "localhost",
			WeatherText: "N/A",
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// 公网 IP 定位
	var city, region, countryCode string
	var lat, lon float64
	var haveCoord bool
	if info, err := utils.LookupIPLocation(ip); err == nil && info != nil {
		city, region, countryCode = info.City, info.Region, info.CountryCode
		if info.Latitude != 0 || info.Longitude != 0 {
			lat, lon, haveCoord = info.Latitude, info.Longitude, true
		}
	}

	// 请求天气
	var (
		weatherText                                  = "天气"
		tempPtr, windPtr, windLvlPtr, humPtr, aqiPtr *int
	)
	if haveCoord {
		if wdata, err := utils.FetchWeatherAndAQI(lat, lon); err == nil && wdata != nil {
			if wdata.TempC != nil {
				v := int(*wdata.TempC)
				tempPtr = &v
			}
			if wdata.WindSpeedKmh != nil {
				v := int(*wdata.WindSpeedKmh)
				windPtr = &v
				lvl := beaufortLevel(v)
				windLvlPtr = &lvl
			}
			if wdata.Humidity != nil {
				v := int(*wdata.Humidity)
				humPtr = &v
			}
			if wdata.AQIUS != nil {
				v := int(*wdata.AQIUS)
				aqiPtr = &v
			}
			weatherText = wdata.WeatherText
		}
	}

	json.NewEncoder(w).Encode(resp{
		City:         city,
		Region:       region,
		CountryCode:  countryCode,
		Location:     buildLocation(city, region, countryCode),
		TemperatureC: tempPtr,
		WindSpeedKmh: windPtr,
		WindLevel:    windLvlPtr,
		Humidity:     humPtr,
		AQIUS:        aqiPtr,
		WeatherText:  weatherText,
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
}
