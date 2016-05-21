package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"config"
)

//解析m3u8文件
func ParseM3u8(m3u8Content string) {
	config.VList = make(map[string]string)
	sReader := strings.NewReader(m3u8Content)
	bReader := bufio.NewScanner(sReader)
	bReader.Split(bufio.ScanLines)

	for bReader.Scan() {
		line := bReader.Text()
		if !strings.HasPrefix(line, "#") {
			writIntoInputs(line)
			config.VList[line] = config.M3u8Url.Scheme + "://" + config.M3u8Url.Host + "/" + line
		}
	}
}

//写入ffmpeg合并输入文件
func writIntoInputs(filename string) {
	tsFp, openErr := os.OpenFile(config.Ffmpeginputs, os.O_WRONLY|os.O_APPEND, 0775)
	if openErr != nil {
		fmt.Println("open local file", config.Ffmpeginputs, "failed,", openErr)
		return
	}
	defer tsFp.Close()
	tsFp.WriteString("file '" + filename + "'\r\n")
}

//清理垃圾文件
func DoClean() {
	fmt.Println("清理垃圾文件...")

	err := os.RemoveAll(config.Tmp)
	if err != nil {
		fmt.Printf("删除垃圾文件失败，%s", err)
	}

	fmt.Println("清理完成...")
}
