package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
var privateCIDRs = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

// IsPrivateOrLoopbackIP 粗略判断是否为私有或回环 IP
func IsPrivateOrLoopbackIP(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	if parsed.IsLoopback() {
		return true
	}
	for _, cidr := range privateCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil && network.Contains(parsed) {
			return true
		}
	}
	return false
}

// ----- IP 地理位置解析 -----
// 设计目标：
// 1. 私有 / 回环地址快速返回占位，不发外网请求
// 2. 主提供商 ipapi.co，有限次重试与日志
// 3. 失败后使用 ipwho.is 作为后备
// 4. 代码结构清晰：小函数 + 明确错误语义

// providerRet 统一返回值，用于内部 provider 调用
type providerRet struct {
	info      *IPInfo
	retriable bool // 是否建议上层重试
	err       error
	name      string
}

// fetchFromIPApi 调用 ipapi.co
func fetchFromIPApi(ip string, client *http.Client) providerRet {
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return providerRet{err: err, retriable: false, name: "ipapi"}
	}
	req.Header.Set("User-Agent", "PersonalSite-Geolocate/1.0")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return providerRet{err: err, retriable: true, name: "ipapi"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if resp.StatusCode != http.StatusOK {
		// 429/5xx 视为可重试，其它直接失败
		retriable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return providerRet{err: fmt.Errorf("ipapi status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body))), retriable: retriable, name: "ipapi"}
	}
	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return providerRet{err: err, retriable: false, name: "ipapi"}
	}
	log.Printf("[IPGeo] ipapi success ip=%s city=%s region=%s countryCode=%s lat=%.4f lon=%.4f", ip, info.City, info.Region, info.CountryCode, info.Latitude, info.Longitude)
	// 至少 city 或 countryCode 有值才认为有效
	if info.City == "" && info.CountryCode == "" {
		return providerRet{err: errors.New("ipapi empty essential fields"), retriable: false, name: "ipapi"}
	}
	return providerRet{info: &info, name: "ipapi"}
}

// fetchFromIPWhoIs 调用 ipwho.is
func fetchFromIPWhoIs(ip string, client *http.Client) providerRet {
	url := fmt.Sprintf("https://ipwho.is/%s", ip)
	resp, err := client.Get(url)
	if err != nil {
		return providerRet{err: err, name: "ipwhois"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if resp.StatusCode != http.StatusOK {
		return providerRet{err: fmt.Errorf("ipwho.is status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body))), name: "ipwhois"}
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return providerRet{err: err, name: "ipwhois"}
	}
	info := &IPInfo{IP: ip}
	if v, ok := raw["city"].(string); ok {
		info.City = v
	}
	if v, ok := raw["region"].(string); ok {
		info.Region = v
	}
	if v, ok := raw["country_code"].(string); ok {
		info.CountryCode = v
	}
	if v, ok := raw["country"].(string); ok {
		info.Country = v
	}
	if v, ok := raw["latitude"].(float64); ok {
		info.Latitude = v
	}
	if v, ok := raw["longitude"].(float64); ok {
		info.Longitude = v
	}
	log.Printf("[IPGeo] ipwho.is success ip=%s city=%s region=%s countryCode=%s lat=%.4f lon=%.4f", ip, info.City, info.Region, info.CountryCode, info.Latitude, info.Longitude)
	if info.City == "" && info.CountryCode == "" {
		return providerRet{err: errors.New("ipwho.is empty essential fields"), name: "ipwhois"}
	}
	return providerRet{info: info, name: "ipwhois"}
}

// LookupIPLocation 统一对外调用：
//   - 私有/回环: 立即返回占位
//   - 主提供商 ipapi 重试 2 次
//   - 失败后 fallback ipwho.is
//   - 返回 IPInfo 或错误
func LookupIPLocation(ip string) (*IPInfo, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil, errors.New("empty ip")
	}
	if IsPrivateOrLoopbackIP(ip) {
		// 私有地址：不请求外部服务
		info := &IPInfo{IP: ip, City: "Local Network", Region: "", Country: "", CountryCode: "", Latitude: 0, Longitude: 0}
		return info, nil
	}
	client := &http.Client{Timeout: 5 * time.Second}

	// 主提供商重试
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		ret := fetchFromIPApi(ip, client)
		if ret.info != nil {
			return ret.info, nil
		}
		lastErr = ret.err
		log.Printf("[IPGeo] ipapi attempt %d error: %v (retriable=%v)", attempt+1, ret.err, ret.retriable)
		if !ret.retriable {
			break
		}
		// 退避
		time.Sleep(time.Duration(200*(attempt+1)) * time.Millisecond)
	}

	// 后备
	fb := fetchFromIPWhoIs(ip, client)
	if fb.info != nil {
		return fb.info, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fb.err
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
