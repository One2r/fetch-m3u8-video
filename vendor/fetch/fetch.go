package fetch

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/astaxie/beego/httplib"

	"config"
)

//从普通页面获取m3u8地址
func GetM3u8Url(Url string) (m3u8Url string) {
	cmd := exec.Command(config.Phantomjs, config.LoadPageJS, Url)
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

//获取m3u8地址获取视频列表
func GetM3u8Content(m3u8Url string) (m3u8 string) {
	req := httplib.Get(m3u8Url).SetTimeout(15*time.Second, 15*time.Second)
	resp, respErr := req.Response()
	if respErr != nil {
		fmt.Println("fetch m3u8 playlist error,", respErr)
		return
	}
	defer resp.Body.Close()
	respData, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println("read m3u8 playlist content error,", readErr)
		return
	} else {
		m3u8 = string(respData)
	}
	return m3u8
}

//视频下载
func Download(url string, fileName string, ch chan int) {
	req := httplib.Get(url).SetTimeout(60*time.Second*30, 60*time.Second*30)
	resp, respErr := req.Response()
	if respErr != nil {
		fmt.Println("download ", url, "failed,", respErr)
		return
	}
	defer resp.Body.Close()
	tsFp, openErr := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0775)
	if openErr != nil {
		fmt.Println("open local file", fileName, "failed,", openErr)
		return
	}
	defer tsFp.Close()
	_, copyErr := io.Copy(tsFp, resp.Body)
	if copyErr != nil {
		fmt.Println("download ", url, " failed,", copyErr)
	}
	ch <- 1
	return
}
