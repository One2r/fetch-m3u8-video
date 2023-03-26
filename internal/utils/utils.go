package utils

import (
    "fmt"
    "log"
    "os"
    "os/exec"
)

// CheckFfmpeg 检测是否安装 ffmpeg
func CheckFfmpeg() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

// CheckNodejs 检测是否安装 nodejs
func CheckNodejs() {
	cmd := exec.Command("node", "-v")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

// Concat 合并视频
func Concat(input string, output string) {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", input, "-c", "copy", output)
	err := cmd.Run()
	if err != nil {
		fmt.Println("合并视频失败，", err)
		return
	} else {
		fmt.Println("合并视完成..")
	}
}

// DoClean 清理垃圾文件
func DoClean(tmpDir string) {
	fmt.Println("清理垃圾文件...")

	err := os.RemoveAll(tmpDir)
	if err != nil {
		fmt.Printf("删除垃圾文件失败，%s", err)
	}

	fmt.Println("清理完成...")
}
