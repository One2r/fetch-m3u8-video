package fetch_test

import (
	"fmt"
	"testing"

	"fetch"
)

func Test_GetM3u8(t *testing.T) {
	rs := fetch.GetM3u8("http://us.sinaimg.cn/0012PjkYjx070lb4wobZ0504010000220k01.m3u8?KID=unistore,video&Expires=1460888111&ssig=QM6fyONFrr&fid=1034:6df61e808a8cb9e7ccf119b087da7b59&uid=&monitor=&vf=vshow")
	fmt.Println(rs)
}
