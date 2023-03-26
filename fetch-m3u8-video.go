package main

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

    "github.com/urfave/cli"

    "fetch-m3u8-video/internal/fetch"
    "fetch-m3u8-video/internal/utils"
    "fetch-m3u8-video/internal/vars"
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
			Name:  "loadpage",
			Value: "loadpage.js",
			Usage: "Puppeteer 脚本",
		},
	}

	app.Action = func(c *cli.Context) {
		utils.CheckFfmpeg()

		if c.String("url") == "" {
			fmt.Println("请输入url地址!")
			return
		}

		if c.String("urltype") == "url" {
			utils.CheckNodejs()
		}

		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)

		vars.LoadPageJS = filepath.Dir(path) + "/scripts/" + c.String("loadpage")

		var m3u8List []string

		if c.String("urltype") == "m3u8" {
			m3u8List = append(m3u8List, c.String("url"))
		} else {
			res := fetch.GetM3u8Url(c.String("url"))
			if res != nil {
				m3u8List = append(m3u8List, res...)
			}
		}

		if m3u8List == nil {
			fmt.Println("未获取到 m3u8 地址")
			return
		}

		var downloadDirPrefix = filepath.Dir(path) + "/downloads/"
		for _, m3u8 := range m3u8List {
			m3u8Url, _ := url.Parse(m3u8)
			fmt.Println("下载：", m3u8Url.String(), " 地址视频：")

			downloadDir := downloadDirPrefix + strings.ReplaceAll(m3u8Url.Hostname(), ".", "_")
			_, err := os.Stat(downloadDir)
			if os.IsNotExist(err) {
				err := os.Mkdir(downloadDir, 0755)
				if err != nil {
					fmt.Println("创建文件目录失败，", err)
					continue
				}
			}

			data := []byte(m3u8Url.String())
			md5Str := md5.Sum(data)
			downloadDir += "/" + fmt.Sprintf("%x", md5Str)
			_, err = os.Stat(downloadDir)
			if os.IsNotExist(err) {
				err := os.Mkdir(downloadDir, 0755)
				if err != nil {
					fmt.Println("创建文件目录失败，", err)
					continue
				}
			}

			tmpDir := downloadDir + "/tmp"
			_, err = os.Stat(tmpDir)
			if os.IsNotExist(err) {
				err := os.Mkdir(tmpDir, 0755)
				if err != nil {
					fmt.Println("创建临时文件目录失败，", err)
					continue
				}
			}

			vList := fetch.GetM3u8Content(m3u8Url)
			if vList != nil {
				ffmpeginputsFile := downloadDir + "/inputs.txt"
				ffmpegInputs, err := os.OpenFile(ffmpeginputsFile, os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					fmt.Println("创建 ffmpeg 输入失败,", err)
					continue
				}
				chs := make([]chan int, len(vList))
				i := 0
				for _, v := range vList {
					downloadUrl, _ := url.Parse(v)
					filename := downloadUrl.Path[strings.LastIndex(downloadUrl.Path, "/")+1:]
					saveFile := tmpDir + "/" + filename
					ffmpegInputs.WriteString("file '" + saveFile + "'\r\n")
					chs[i] = make(chan int)
					go fetch.Download(downloadUrl.String(), saveFile, chs[i])
					i++
				}
				for _, ch := range chs {
					<-ch
				}
				ffmpegInputs.Close()
				outputFile := downloadDir + "/video.mp4"
				utils.Concat(ffmpeginputsFile, outputFile)
			} else {
				fmt.Println("获取m3u8文件索引失败！")
				return
			}
		}
	}
	app.Run(os.Args)
}
