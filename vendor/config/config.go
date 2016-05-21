package config

import (
	"net/url"
)

//视频输出文件
var Output string = ""

//phantomjs可执行文件地址
var Phantomjs string = ""

//ffmpeg可执行文件地址
var Ffmpeg string = ""

//phantomjs执行js脚本，用于加载普通页面获取m3u8地址
var LoadPageJS string = ""

//m3u8 URL对象
var M3u8Url *url.URL

//解析m3u后的视频列表
var VList map[string]string

//临时视频目录
var Tmp string = ""

//ffmpeg合并输入文件地址
var Ffmpeginputs = ""
