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

// WeatherAQI 描述天气与空气质量（部分字段可为空）。
// 数值字段保持与外部 API 一致的单位：
//
//	TempC: 摄氏度
//	WindSpeedKmh: 公里/小时
//	Humidity: 相对湿度百分比
//	AQIUS: 美国 AQI 指标
//
// WeatherText: 经过 WeatherCodeToText 映射后的中文描述
type WeatherAQI struct {
	TempC        *float64
	WindSpeedKmh *float64
	Humidity     *float64
	AQIUS        *float64
	WeatherText  string
}

// WeatherCodeToText 将 Open‑Meteo weather_code 映射为中文描述。
// 参考 https://open-meteo.com/en/docs 里的 weather code 列表，按语义合并常见类别。
// 未知或缺省返回 "天气" 占位，保持与此前逻辑兼容。
func WeatherCodeToText(code int) string {
	switch code { // 分组归类
	case 0:
		return "晴"
	case 1, 2, 3:
		return "多云"
	case 45, 48:
		return "雾"
	// 毛毛雨 / 雨
	case 51, 53, 55, 56, 57, 61, 63, 65, 80, 81, 82:
		return "雨"
	// 雪 & 冻降水
	case 71, 73, 75, 77, 85, 86:
		return "雪"
	// 雷暴
	case 95, 96, 99:
		return "雷阵雨"
	default:
		return "天气"
	}
}

// FetchWeatherAndAQI 使用 Open-Meteo 获取当前天气与 US AQI
func FetchWeatherAndAQI(lat, lon float64) (*WeatherAQI, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	// 构造请求 URL：仅拉取当前需要用到的字段，减小响应体
	wURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=auto",
		lat, lon,
	)
	aqiURL := fmt.Sprintf(
		"https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%f&longitude=%f&current=us_aqi",
		lat, lon,
	)

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

	var (
		wData wResp
		aData aqiResp
	)

	// helper：执行 GET 并在 200 时解 JSON，不抛错（容忍失败）
	fetchJSON := func(url string, v any) {
		if resp, err := client.Get(url); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				if b, err := io.ReadAll(resp.Body); err == nil {
					_ = json.Unmarshal(b, v)
				}
			}
		}
	}

	fetchJSON(wURL, &wData)
	fetchJSON(aqiURL, &aData)

	weatherText := "天气"
	if wData.Current.WeatherCode != nil {
		weatherText = WeatherCodeToText(*wData.Current.WeatherCode)
	}

	return &WeatherAQI{
		TempC:        wData.Current.Temperature2M,
		WindSpeedKmh: wData.Current.WindSpeed10M,
		Humidity:     wData.Current.RelativeHumidity2M,
		AQIUS:        aData.Current.USAqi,
		WeatherText:  weatherText,
	}, nil
}
