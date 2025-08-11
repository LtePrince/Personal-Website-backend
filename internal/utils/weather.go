package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// IPInfo 描述 IP 定位信息
type IPInfo struct {
	IP          string  `json:"ip"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// IsPrivateOrLoopbackIP 粗略判断私有或回环地址
func IsPrivateOrLoopbackIP(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	if parsed.IsLoopback() {
		return true
	}
	// 私有网段
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	for _, cidr := range privateCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil && network.Contains(parsed) {
			return true
		}
	}
	return false
}

// LookupIPLocation 使用 ipapi.co 解析 IP
func LookupIPLocation(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/json/", strings.TrimSpace(ip))
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("ip geolocation failed")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// WeatherAQI 描述天气与空气质量
type WeatherAQI struct {
	TempC        *float64
	WindSpeedKmh *float64
	Humidity     *float64
	AQIUS        *float64
	WeatherText  string
}

// FetchWeatherAndAQI 使用 Open-Meteo 获取当前天气与 US AQI
func FetchWeatherAndAQI(lat, lon float64) (*WeatherAQI, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	wURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=auto", lat, lon)
	aqiURL := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%f&longitude=%f&current=us_aqi", lat, lon)

	type wResp struct {
		Current struct {
			Temperature2M      *float64 `json:"temperature_2m"`
			RelativeHumidity2M *float64 `json:"relative_humidity_2m"`
			WindSpeed10M       *float64 `json:"wind_speed_10m"`
			WeatherCode        *int     `json:"weather_code"`
		} `json:"current"`
	}
	type aqiResp struct {
		Current struct {
			USAqi *float64 `json:"us_aqi"`
		} `json:"current"`
	}

	// 并行请求（简化为顺序 + 独立错误容忍）
	var wData wResp
	var aData aqiResp

	// 天气
	if resp, err := client.Get(wURL); err == nil {
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()
		if resp.StatusCode == 200 {
			if b, err := io.ReadAll(resp.Body); err == nil {
				_ = json.Unmarshal(b, &wData)
			}
		}
	}
	// AQI
	if resp, err := client.Get(aqiURL); err == nil {
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()
		if resp.StatusCode == 200 {
			if b, err := io.ReadAll(resp.Body); err == nil {
				_ = json.Unmarshal(b, &aData)
			}
		}
	}

	// 文本映射（简化）
	weatherText := "天气"
	if code := wData.Current.WeatherCode; code != nil {
		switch *code {
		case 0:
			weatherText = "晴"
		case 1, 2, 3:
			weatherText = "多云"
		case 45, 48:
			weatherText = "雾"
		case 61, 63, 65, 80, 81, 82:
			weatherText = "雨"
		case 71, 73, 75, 85, 86:
			weatherText = "雪"
		case 95, 96, 99:
			weatherText = "雷阵雨"
		default:
			weatherText = "天气"
		}
	}

	return &WeatherAQI{
		TempC:        wData.Current.Temperature2M,
		WindSpeedKmh: wData.Current.WindSpeed10M,
		Humidity:     wData.Current.RelativeHumidity2M,
		AQIUS:        aData.Current.USAqi,
		WeatherText:  weatherText,
	}, nil
}
