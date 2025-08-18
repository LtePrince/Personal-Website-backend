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
