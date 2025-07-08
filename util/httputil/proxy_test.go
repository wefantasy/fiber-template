package httputil

import (
	"app/conf"
	"app/log"
	"app/util"
	"testing"
)

func TestGetHttpTransportWithProxy(t *testing.T) {
	conf.Initialize()
	_, err := GetProxyTransportFromApi(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckProxyAvailability(t *testing.T) {
	conf.Initialize()
	log.Initialize()
	tests := []struct {
		name    string
		proxy   string
		testURL string
		want    string
	}{
		{
			name:    "Valid HTTPS proxy",
			proxy:   "http://72.10.160.90:9733",
			testURL: "https://httpbin.org/get",
			want:    "https",
		},
		{
			name:    "Valid HTTP proxy only",
			proxy:   "socks5://198.177.253.13:4145",
			testURL: "https://httpbin.org/get",
			want:    "http",
		},
		{
			name:    "Invalid proxy address",
			proxy:   ":invalid",
			testURL: "https://httpbin.org/get",
			want:    "",
		},
		{
			name:    "Unavailable proxy",
			proxy:   "http://134.209.29.120:80",
			testURL: "https://httpbin.org/get",
			want:    "",
		},
		{
			name:    "test socks4 proxy",
			proxy:   "socks4://198.177.253.13:4145",
			testURL: "https://httpbin.org/get",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例设置mock行为
			if got := CheckProxyAvailabilityWithTestUrl(util.NewRootContext(), tt.proxy, tt.testURL); got != tt.want {
				t.Errorf("CheckProxyAvailability() = %v, want %v", got, tt.want)
			}
		})
	}
}
