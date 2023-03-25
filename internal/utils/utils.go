package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

    "fetch-m3u8-video/internal/vars"
)

// 检测是否安装 ffmpeg
func CheckFfmpeg() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

// 检测是否安装 nodejs
func CheckNodejs() {
	cmd := exec.Command("node", "-v")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

// 解析m3u8文件
func ParseM3u8(m3u8Content string) {
	vars.VList = make(map[string]string)
	sReader := strings.NewReader(m3u8Content)
	bReader := bufio.NewScanner(sReader)
	bReader.Split(bufio.ScanLines)

	for bReader.Scan() {
		line := bReader.Text()
		if !strings.HasPrefix(line, "#") {
			writIntoInputs(line)
			vars.VList[line] = vars.M3u8Url.Scheme + "://" + vars.M3u8Url.Host + "/" + line
		}
	}
}

// 写入ffmpeg合并输入文件
func writIntoInputs(filename string) {
	tsFp, openErr := os.OpenFile(vars.FfmpeginputsFile, os.O_WRONLY|os.O_APPEND, 0775)
	if openErr != nil {
		fmt.Println("open local file", vars.FfmpeginputsFile, "failed,", openErr)
		return
	}
	defer tsFp.Close()
	tsFp.WriteString("file '" + filename + "'\r\n")
}

// 合并视频
func Concat() {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-i", vars.FfmpeginputsFile, "-c", "copy", vars.OutputDir)
	err := cmd.Run()
	if err != nil {
		fmt.Println("合并视频失败，", err)
		return
	} else {
		fmt.Println("合并视完成..")
	}
}

// 清理垃圾文件
func DoClean() {
	fmt.Println("清理垃圾文件...")

	err := os.RemoveAll(vars.TmpDir)
	if err != nil {
		fmt.Printf("删除垃圾文件失败，%s", err)
	}

	fmt.Println("清理完成...")
}
