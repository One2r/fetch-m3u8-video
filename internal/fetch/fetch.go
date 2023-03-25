package fetch

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"

	config "fetch-m3u8-video/internal/configs"
)

// 从普通页面获取m3u8地址
func GetM3u8Url(Url string) (m3u8Url string) {
	cmd := exec.Command("node", config.LoadPageJS, Url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		reg := regexp.MustCompile(`(url:).*`)
		urlStr := reg.FindString(out.String())
		size := len(urlStr)
        if (size <= 4){
            return
        }
		rawUrl, err := url.QueryUnescape(urlStr[4 : size-1])
		if err != nil {
			fmt.Println(nil)
			return
		} else {
			config.M3u8Url, _ = url.Parse(rawUrl)
			m3u8Url = rawUrl
		}
	}
	return m3u8Url
}

// 获取m3u8地址获取视频列表
func GetM3u8Content(m3u8Url string) (m3u8 string) {
	client := resty.New()
	client.SetTimeout(15 * time.Second)
	resp, respErr := client.R().Get(m3u8Url)
	if respErr != nil {
		fmt.Println("fetch m3u8 playlist error,", respErr)
		return
	}
	m3u8 = resp.String()
	return m3u8
}

// 视频下载
func Download(url string, fileName string, ch chan int) {
	client := resty.New()
	client.SetTimeout(60 * time.Second * 30)
	_, respErr := client.R().SetOutput(fileName).Get(url)
	if respErr != nil {
		fmt.Println("download ", url, "failed,", respErr)
		return
	}
	ch <- 1
	return
}
