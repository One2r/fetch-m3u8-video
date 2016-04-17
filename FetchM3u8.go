package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/codegangsta/cli"

	"config"
	"fetch"
)

var outputFile string = ""

func main() {
	app := cli.NewApp()
	app.Name = "Fetch M3u8 Video"
	app.Usage = "获取m3u8文件索引视频"
	app.Version = "1.0.0"
	app.Author = "One2r"
	app.Email = "601941036@qq.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "o",
			Value: "YourVideo",
			Usage: "输出文件,默认在Downloads目录下",
		},
		cli.StringFlag{
			Name:  "url",
			Value: "",
			Usage: "m3u8文件url地址,或者是包含m3u8文件地址的普通url地址,通过urltype指定url地址类型",
		},
		cli.StringFlag{
			Name:  "urltype",
			Value: "url",
			Usage: "url类型，默认为url。m3u8:  m3u8文件url地址;url:  包含m3u8文件地址的普通url地址",
		},
		cli.StringFlag{
			Name:  "phantomjs",
			Value: "/usr/bin/phantomjs",
			Usage: "phantomjs可执行文件目录",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.String("url") == "" {
			fmt.Println("请输入url地址!")
			os.Exit(1)
		}
		if c.String("phantomjs") == "" {
			fmt.Println("请输入phantomjs可执行文件目录")
			os.Exit(1)
		} else {
			config.Phantomjs = c.String("phantomjs")
		}

		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		config.LoadPageJS = filepath.Dir(path) + "/loadpage.js"

		Url := c.String("url")
		outputFile = "./Downloads/" + c.String("o")
		if c.String("urltype") == "m3u8" {
			fetchMovie(Url)
		} else {
			m3u8Url := fetch.GetM3u8Url(Url)
			fetchMovie(m3u8Url)
		}
		fmt.Println("Result:", c.String("o"))
	}
	app.Run(os.Args)
}

func fetchMovie(m3u8Url string) {
	req := httplib.Get(m3u8Url).SetTimeout(30*time.Second, 30*time.Second)
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
	}
	m3u8Uri, _ := url.Parse(m3u8Url)

	sReader := strings.NewReader(string(respData))
	bReader := bufio.NewScanner(sReader)
	bReader.Split(bufio.ScanLines)

	for bReader.Scan() {
		line := bReader.Text()
		if !strings.HasPrefix(line, "#") {
			tsFileName := line
			tsLocalFileName := "./Downloads/tmp/" + line
			tsFileUrl := fmt.Sprintf("%s://%s/%s", m3u8Uri.Scheme, m3u8Uri.Host, tsFileName)
			downloadTS(tsFileUrl, tsLocalFileName)
		}
	}
}

func downloadTS(tsFileUrl string, tsLocalFileName string) {
	fmt.Println("downloading", tsFileUrl)
	req := httplib.Get(tsFileUrl).SetTimeout(60*time.Second*30, 60*time.Second*30)
	resp, respErr := req.Response()
	if respErr != nil {
		fmt.Println("download ts ", tsFileUrl, "failed,", respErr)
		return
	}
	defer resp.Body.Close()
	tsFp, openErr := os.OpenFile(tsLocalFileName, os.O_CREATE|os.O_WRONLY, 0775)
	if openErr != nil {
		fmt.Println("open local file", tsLocalFileName, "failed,", openErr)
		return
	}
	defer tsFp.Close()
	_, copyErr := io.Copy(tsFp, resp.Body)
	if copyErr != nil {
		fmt.Println("download ts", tsFileUrl, " failed,", copyErr)
	}
	return
}
