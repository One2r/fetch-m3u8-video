# FetchM3u8
Fetch M3u8 Video 获取m3u8文件索引视频（Go语言版）

Go版本：go1.5.3  

Go包依赖：  
github.com/codegangsta/cli  
github.com/astaxie/beego/httplib

第三方软件依赖：  
PhantomJS  
FFmpeg  

#用法  
##指令
`-o`  
输出文件,默认在Downloads目录下，默认"YourVideo.avi"  

`--url`  
m3u8文件url地址,或者是包含m3u8文件地址的普通url地址,通过urltype指定url地址类型  

`--urltype`  
url类型，默认为url。m3u8:  m3u8文件url地址;url:  包含m3u8文件地址的普通url地址  

`--phantomjs`  
phantomjs可执行文件地址，默认"/usr/bin/phantomjs"  

`--ffmpeg`  
ffmpeg可执行文件地址，默认"/usr/bin/ffmpeg"  

`--help, -h`  
show help  

`--version, -v`  
print the version  

##示例  
`./FetchM3u8 --url "http://weibo.com/p/134132432132413241" --phantomjs /home/svr/phantomjs/bin/phantomjs --ffmpeg /home/svr/ffmpeg/bin/ffmpeg`  

#鸣谢
感谢@jemygraw在Go友团分享的代码，地址：[http://golanghome.com/post/645](http://golanghome.com/post/645)

#许可
Licensed under the MIT license
