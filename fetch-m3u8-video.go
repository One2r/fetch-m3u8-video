package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"

	config "fetch-m3u8-video/internal/configs"
	"fetch-m3u8-video/internal/fetch"
	"fetch-m3u8-video/internal/utils"
)

func main() {
	app := cli.NewApp()
	app.Name = "fetch-m3u8-video"
	app.Usage = "获取 m3u8 格式视频"
	app.Version = "2.0.0"
	app.Author = "One2r"
	app.Email = "601941036@qq.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Value: "",
			Usage: "m3u8 文件 url 地址,或者是包含 m3u8 地址的普通 url 地址,通过 urltype 指定 url 地址类型",
		},
		cli.StringFlag{
			Name:  "urltype",
			Value: "url",
			Usage: "url 类型，默认为 url。m3u8:  m3u8文件url地址; url:  包含m3u8文件地址的普通 url 地址",
		},
        cli.StringFlag{
			Name:  "out",
			Value: "YourVideo.avi",
			Usage: "输出文件,默认在 downloads 目录下",
		},
	}

	app.Action = func(c *cli.Context) {
		utils.CheckFfmpeg()

		if c.String("url") == "" {
			fmt.Println("请输入url地址!")
			os.Exit(1)
		}

        if c.String("urltype") == "url" {
            utils.CheckNodejs()
        }

		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)

		config.Output = filepath.Dir(path) + "/downloads/" + c.String("o")
		config.LoadPageJS = filepath.Dir(path) + "/scripts/loadpage.js"
		config.Tmp = filepath.Dir(path) + "/downloads/tmp/"

		err := os.Mkdir(filepath.Dir(path)+"/downloads/tmp/", 0755)
		if err != nil {
			fmt.Println("创建临时文件目录失败，", err)
			os.Exit(1)
		}

		config.Ffmpeginputs = filepath.Dir(path) + "/downloads/tmp/inputs.txt"

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
