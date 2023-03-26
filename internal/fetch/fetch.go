package fetch

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net/url"
    "os/exec"
    "strings"
    "time"

    resty "github.com/go-resty/resty/v2"

    "fetch-m3u8-video/internal/vars"
)

// GetM3u8Url 从普通页面获取m3u8地址
func GetM3u8Url(Url string) (m3u8Urls []string) {
	cmd := exec.Command("node", vars.LoadPageJS, Url)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	} else {
        err:=json.Unmarshal([]byte(out.String()),&m3u8Urls)
        if err != nil {
            fmt.Println("解析json数据错误：",err)
            return
        }
	}
	return m3u8Urls
}

// GetM3u8Content 获取m3u8地址获取视频列表
func GetM3u8Content(url *url.URL) (vList []string) {
    fmt.Println(url.String())
	client := resty.New()
	client.SetTimeout(15 * time.Second)
	resp, respErr := client.R().Get(url.String())
	if respErr != nil {
		fmt.Println("fetch m3u8 playlist error,", respErr)
		return
	}
    content := resp.String()
    sReader := strings.NewReader(content)
    bReader := bufio.NewScanner(sReader)
    bReader.Split(bufio.ScanLines)
    for bReader.Scan() {
        line := bReader.Text()
        if !strings.HasPrefix(line, "#") {
            var tsUrl string
            if strings.HasPrefix("http",line) {
                tsUrl = line
            }else{
                tsUrl = url.Scheme + "://" + url.Host + url.Path[0:strings.LastIndex(url.Path,"/")+1]+line
            }
            if strings.Contains(line,".m3u8"){
                newM3u8Url, _ := url.Parse(tsUrl)
                vList = append(vList,GetM3u8Content(newM3u8Url)...)
            }else{
                vList = append(vList,tsUrl)
            }
        }
    }
	return vList
}

// Download 视频下载
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
