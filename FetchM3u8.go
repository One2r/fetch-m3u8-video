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
			Value: "YourVideo.avi",
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
			Usage: "phantomjs可执行文件地址",
		},
		cli.StringFlag{
			Name:  "ffmpeg",
			Value: "/usr/bin/ffmpeg",
			Usage: "ffmpeg可执行文件地址",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.String("url") == "" {
			fmt.Println("请输入url地址!")
			os.Exit(1)
		}
		if c.String("phantomjs") == "" {
			fmt.Println("请输入phantomjs可执行文件地址")
			os.Exit(1)
		} else {
			config.Phantomjs = c.String("phantomjs")
		}
		if c.String("ffmpeg") == "" {
			fmt.Println("请输入ffmpeg可执行文件地址")
			os.Exit(1)
		} else {
			config.Ffmpeg = c.String("ffmpeg")
		}

		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)

		config.Output = filepath.Dir(path) + "/Downloads/" + c.String("o")
		config.LoadPageJS = filepath.Dir(path) + "/loadpage.js"
		config.Tmp = filepath.Dir(path) + "/Downloads/tmp/"

		err := os.Mkdir(filepath.Dir(path)+"/Downloads/tmp/", 0755)
		if err != nil {
			fmt.Println("创建临时文件目录失败，", err)
			os.Exit(1)
		}

		config.Ffmpeginputs = filepath.Dir(path) + "/Downloads/tmp/inputs.txt"

		ffmpegInputs, err := os.Create(config.Ffmpeginputs)
		defer ffmpegInputs.Close()
		if err != nil {
			fmt.Println("创建文件失败,", err)
			return
		}

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
			chs := make([]chan int, len(config.VList))
			i := 0
			for k, v := range config.VList {
				chs[i] = make(chan int)
				go fetch.Download(v, config.Tmp+k, chs[i])
				i++
			}
			for _, ch := range chs {
				<-ch
			}
		} else {
			fmt.Println("获取m3u8文件索引失败！")
			os.Exit(1)
		}

		utils.Concat()
		utils.DoClean()

		fmt.Println("Result:", c.String("o"))
	}
	app.Run(os.Args)
}
