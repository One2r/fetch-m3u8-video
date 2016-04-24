package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"

	"config"
	"fetch"
	"utils"
)

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
		config.Tmp = filepath.Dir(path) + "/Downloads/tmp/"

		Url := c.String("url")
		m3u8Ct := ""

		if c.String("urltype") == "m3u8" {
			m3u8Ct = fetch.GetM3u8Content(Url)
		} else {
			m3u8Url := fetch.GetM3u8Url(Url)
			m3u8Ct = fetch.GetM3u8Content(m3u8Url)
		}

		if m3u8Ct != "" {
			utils.ParseM3u8(m3u8Ct)
			for k, v := range config.VList {
				fetch.Download(v, config.Tmp+k)
			}
		} else {
			fmt.Println("获取m3u8文件索引失败！")
			os.Exit(1)
		}

		fmt.Println("Result:", c.String("o"))
	}
	app.Run(os.Args)
}
