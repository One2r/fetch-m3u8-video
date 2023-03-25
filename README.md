# fetch-m3u8-video
fetch-m3u8-video 获取 m3u8 格式视频（Go语言版）

# 编译
```
go mod tidy
make 
make build-win
```

# 第三方软件依赖： 
- FFmpeg
## 可选依赖
- upx
- NodeJS
- Puppeteer   

# 用法  
## 指令
`-out`  
输出文件,默认在Downloads目录下，默认"YourVideo.avi"  

`--url`  
m3u8 文件 url 地址,或者是有 m3u8 地址的页面 url 地址,通过 urltype 指定 url 地址类型。  

`--urltype`  
url 类型，默认为 url。
- m3u8: 指定 url 为 m3u8 地址; 
- url: 指定 url 为有 m3u8 的页面 url 地址；  

`--loadpage`  
Puppeteer 脚本，默认"loadpage.js"

`--help, -h`  
show help  

`--version, -v`  
print the version  

## 示例  
`./fetch-m3u8-video --url "https://www.cxtvlive.com/live-tv/canal-once"

# 许可
Licensed under the MIT license
