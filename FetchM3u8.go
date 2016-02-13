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

var outputFile *File

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
			Usage: "m3u8文件url地址",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("url") == "" {
			fmt.Println("请输入m3u8文件的url地址!")
			os.Exit(1)
		}
		Url := c.String("url")

		tsFp, openErr := os.OpenFile("./Downloads/"+c.String("o"), os.O_CREATE|os.O_WRONLY, 0775)
		if openErr != nil {
			fmt.Println("open local file", "./Downloads/"+c.String("o"), "failed,", openErr)
			os.Exit(1)
		}
		outputFile = tsFp
		fmt.Println(tsFp)

		fetchM3u8(Url)
	}
	fmt.Printf("\n========%s========\n", app.Name)
	app.Run(os.Args)
}

func fetchM3u8(Url string) {
	req := httplib.Get(Url)
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
	req := httplib.Get(m3u8Url)
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
	videoFolder := fmt.Sprintf("video_%d", time.Now().Unix())
	mkdErr := os.Mkdir(videoFolder, 0775)
	if mkdErr != nil {
		fmt.Println("mkdir for m3u8 playlist failed,", mkdErr)
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
			tsLocalFileName := fmt.Sprintf("%s/%s", videoFolder, line)
			tsFileUrl := fmt.Sprintf("%s://%s/%s", m3u8Uri.Scheme, m3u8Uri.Host, tsFileName)
			downloadTS(tsFileUrl, tsLocalFileName)
		}
	}
	fmt.Println("Result:", videoFolder)
}

func downloadTS(tsFileUrl string, tsLocalFileName string) {
	fmt.Println("downloading", tsFileUrl)
	req := httplib.Get(tsFileUrl)
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
}
