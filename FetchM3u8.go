package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/codegangsta/cli"
)

var outputFile *os.File

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
			Usage: "默认为url。\r\nm3u8:  m3u8文件url地址;\r\nurl:  包含m3u8文件地址的普通url地址",
		},

	}
	app.Action = func(c *cli.Context) {
		if c.String("url") == "" {
			fmt.Println("请输入url地址!")
			os.Exit(1)
		}
		Url := c.String("url")

		tsFp, openErr := os.OpenFile("./Downloads/"+c.String("o"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0775)
		if openErr != nil {
			fmt.Println("open local file", "./Downloads/"+c.String("o"), "failed,", openErr)
			os.Exit(1)
		}
		outputFile = tsFp
		defer outputFile.Close()
		defer tsFp.Close()

		if c.String("urltype") == "m3u8" {
			fetchMovie(Url)
		}else{
			fetchM3u8(Url)
		}
		fmt.Println("Result:", c.String("o"))
	}
	fmt.Printf("\n========%s========\n", app.Name)
	app.Run(os.Args)
}

func fetchM3u8(Url string) {
	req := httplib.Get(Url).SetTimeout(60*time.Second*30, 60*time.Second*30)
	resp, err := req.Response()
	if err != nil {
		fmt.Println("fetch video url error,", err)
		return
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read video url content error,", err)
		return
	}
	pattern := `flashvars="list=.*"`
	regx := regexp.MustCompile(pattern)
	flashVars := regx.FindString(string(respData))
	if flashVars != "" {
		size := len(flashVars)
		Subm3u8Url, err := url.QueryUnescape(flashVars[16 : size-1])
		if err != nil {
			fmt.Println(err)
		} else {
			fetchMovie(Subm3u8Url)
		}
	} else {
		fmt.Println("m3u8 playlist not found")
		return
	}
}

func fetchMovie(m3u8Url string) {
	req := httplib.Get(m3u8Url).SetTimeout(60*time.Second*30, 60*time.Second*30)
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
			tsFileUrl := fmt.Sprintf("%s://%s/%s", m3u8Uri.Scheme, m3u8Uri.Host, tsFileName)
			downloadTS(tsFileUrl)
		}
	}
}

func downloadTS(tsFileUrl string) {
	fmt.Println("downloading", tsFileUrl)
	req := httplib.Get(tsFileUrl).SetTimeout(60*time.Second*30, 60*time.Second*30)
	resp, respErr := req.Response()
	if respErr != nil {
		fmt.Println("download ts ", tsFileUrl, "failed,", respErr)
		return
	}
	defer resp.Body.Close()
	_, copyErr := io.Copy(outputFile, resp.Body)
	if copyErr != nil {
		fmt.Println("download ts", tsFileUrl, " failed,", copyErr)
	}
}
