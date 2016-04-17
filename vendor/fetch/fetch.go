package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
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
		url, err := url.QueryUnescape(urlStr[4 : size-1])
		if err != nil {
			fmt.Println(nil)
			return
		} else {
			m3u8Url = url
		}
	}
	return m3u8Url
}

//获取m3u8地址获取视频列表
func GetM3u8(m3u8Url string) (m3u8 string) {
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
