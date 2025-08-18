package test

import (
	"testing"

	"github.com/LtePrince/Personal-Website-backend/internal/utils"
)

func TestWeatherCodeToText(t *testing.T) {
	cases := map[int]string{
		0:    "晴",
		1:    "多云",
		2:    "多云",
		3:    "多云",
		45:   "雾",
		61:   "雨",
		71:   "雪",
		95:   "雷阵雨",
		1234: "天气", // 未知码
	}
	for code, want := range cases {
		got := utils.WeatherCodeToText(code)
		if got != want {
			t.Fatalf("code %d => %s, want %s", code, got, want)
		}
	}
}

func TestFetchWeatherAndAQI(t *testing.T) {
	res, err := utils.FetchWeatherAndAQI(-33.8688, 151.2093) // Sydney 经纬度
	if err != nil {
		// 网络失败不视为致命，记录并返回（保持流水线鲁棒）
		t.Logf("FetchWeatherAndAQI error: %v (tolerated)", err)
		return
	}
	if res == nil {
		t.Fatalf("expected non-nil result")
	}
	if res.WeatherText == "" {
		t.Fatalf("WeatherText should not be empty")
	}
}

func TestLookupIPLocation(t *testing.T) {
	ip := "154.37.213.201"
	info, err := utils.LookupIPLocation(ip)
	if err != nil {
		t.Logf("LookupIPLocation network error (tolerated): %v", err)
		return
	}
	if info == nil {
		t.Fatalf("expected non-nil info")
	}
	if info.IP == "" { // ipapi 可能未回填 ip 字段，也可忽略
		t.Logf("warning: empty IP field in response")
	}
	if info.CountryCode == "" && info.City == "" && info.Region == "" {
		t.Fatalf("all location fields empty, response: %+v", info)
	}
	t.Logf("Parsed IP %s => City=%s Region=%s CountryCode=%s Lat=%.4f Lon=%.4f", ip, info.City, info.Region, info.CountryCode, info.Latitude, info.Longitude)
}
