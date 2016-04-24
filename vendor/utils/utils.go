package utils

import (
	"bufio"
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
			config.VList[line] = config.M3u8Url.Scheme + "://" + config.M3u8Url.Host + "/" + line
		}
	}
}
