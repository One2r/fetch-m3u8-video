package config

import (
	"net/url"
)

//phantomjs可执行文件目录
var Phantomjs string = ""

//phantomjs执行js脚本，用于加载普通页面获取m3u8地址
var LoadPageJS string = ""

//m3u8 URL对象
var M3u8Url *url.URL

//解析m3u后的视频列表
var VList map[string]string

//临时视频目录
var Tmp string = ""
