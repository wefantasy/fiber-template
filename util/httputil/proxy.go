package httputil

import (
	"app/conf"
	"app/log"
	"app/util"
	"context"
	"crypto/tls"
	"fmt"
	"h12.io/socks"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ExecRequestWithProxy 用于执行带有代理的 HTTP 请求。
// method: HTTP 请求方法（GET、POST 等）
// url: 请求的 URL
// body: 请求体，可以是字符串或结构体
// header: 请求头
// retry: 重试次数
//
// 返回值:
//
//	string: 响应体
//	error: 如果请求失败，则返回错误
func ExecRequestWithProxy(method, url string, body any, header map[string]string, retry int) (string, error) {
	for i := 0; i < retry; i++ {
		log.Infof("Request with proxy, retry %d", i)
		transport, err := GetProxyTransportFromApi()
		if err != nil {
			continue
		}
		resp, err := RequestBase(method, url, body, header, transport, 15*time.Second)
		if err != nil {
			if IsStatusFailed(err) {
				return "", err
			}
			continue
		}
		return resp, nil
	}
	return "", fmt.Errorf("request with proxy failed after %d retries", retry)
}

// GetProxyTransportFromApi 用于从 API 获取一个代理的 Transport。
//
// 返回值:
//
//	*http.Transport: 解析成功则返回配置好的 Transport 实例，失败则返回默认的空 Transport
//	error: 如果获取代理失败，则返回错误
func GetProxyTransportFromApi() (*http.Transport, error) {
	header := make(map[string]string)
	header["X-API-Secret"] = conf.Proxy.Secret

	if conf.Proxy.BaseUrl == "" {
		return nil, fmt.Errorf("proxy server is empty")
	}
	respBody, err := Request(http.MethodGet, conf.Proxy.BaseUrl, "", header, nil)
	if err != nil {
		return nil, err
	}

	protocol, err := util.JsonIndex(respBody,
		"data.protocol")
	if err != nil {
		return nil, err
	}
	ip, _ := util.JsonIndex(respBody,
		"data.ip")
	port, _ := util.JsonIndex(respBody,
		"data.port")
	proxy := fmt.Sprintf("%s://%s:%s", protocol, ip, port)
	return GetTransportWithUrl(proxy)
}

// GetTransportWithUrl 用于获取一个 HTTP 或 SOCKS 代理的 Transport。
//
// 参数:
//
// proxy: 代理服务器的 URL (例如 "http://127.0.0.1:8080", "socks5://127.0.0.1:1080")
//
// 返回值:
//
//	*http.Transport: 解析成功则返回配置好的 Transport 实例，失败则返回默认的空 Transport
//	error: 如果解析代理地址失败，则返回错误
func GetTransportWithUrl(proxy string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Debug("proxy format failed: ", err)
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	if proxyURL.Scheme == "socks5" || proxyURL.Scheme == "socks4" || proxyURL.Scheme == "socks4a" {
		transport.Proxy = nil
		dialer := socks.Dial(proxyURL.String())
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Infof("Dialing to %s through SOCKS proxy %s", addr, proxyURL.Host)
			return dialer(network, addr)
		}
	} else {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			log.Infof("Dialing through HTTP proxy %s", proxyURL.Host)
			return proxyURL, nil
		}
	}

	return transport, nil
}

// CheckProxyAvailability 用于测试一个代理地址的可用性。
// 它首先尝试直接通过代理请求测试地址，如果失败，则会禁用证书校验后重试。
//
// 参数:
//
//	proxy: 代理服务器的地址 (例如 "http://127.0.0.1:8080", "socks5://127.0.0.1:1080")
//	testURL:   用于测试的URL (例如 "https://httpbin.org/get", "https://ip-api.io/json")
//
// 返回值:
//
//	string: "http"/"https"/""。
func CheckProxyAvailability(proxy string) string {
	return CheckProxyAvailabilityWithTestUrl(proxy, "https://httpbin.org/get")
}

// CheckProxyAvailabilityWithTestUrl 用于测试一个代理地址的可用性。
// 通过proxy访问testUrl来测试可用性，它首先尝试直接通过代理请求测试地址，如果失败，则会禁用证书校验后重试。
func CheckProxyAvailabilityWithTestUrl(proxy, testUrl string) string {
	log.Info("开始验证代理: ", proxy)

	if checkProxy(proxy, testUrl, false) {
		log.Infof("✅ 代理 %s 可用于 HTTPS", proxy)
		return "https"
	}

	if checkProxy(proxy, testUrl, true) {
		log.Infof("✅ 代理 %s 可用于 HTTP", proxy)
		return "http"
	} else {
		return ""
	}
}

// checkProxy 是一个内部辅助函数，用于执行实际的HTTP请求。
func checkProxy(proxyAddr string, testURL string, insecureSkipVerify bool) bool {
	transport, err := GetTransportWithUrl(proxyAddr)
	if err != nil {
		return false
	}

	if insecureSkipVerify {
		// 配置 TLS 客户端以决定是否跳过证书验证
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureSkipVerify}
	}

	_, err = RequestBase(http.MethodGet, testURL, "", map[string]string{}, transport, 10*time.Second)
	if err != nil {
		return false
	}
	return true
}
