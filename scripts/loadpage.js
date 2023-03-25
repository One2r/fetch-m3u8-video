/**
 * demo loadpage.js for https://www.cxtvlive.com/live-tv/canal-once
 * 返回 ["http://xxxx.com/xxx/xxx.m3u8","http://xxxx.com/yyy/yyy.m3u8"] 格式数据
 */
const puppeteer = require('puppeteer');
var url = process.argv[2];

(async () => {
  const browser = await puppeteer.launch( {args: ["--no-sandbox"]});
  const page = await browser.newPage();
  await page.goto(url, {waitUntil: 'networkidle2'});

  const m3u8_list = await page.evaluate(() => {
    var m3u8_list =[];
    var videos = document.getElementsByTagName("video")
    if (videos){
        for (i = 0; i < videos.length; i++) { 
            var source = videos[i].getElementsByTagName("source")
            if (source){
                var url = source[0].getAttribute("src")
                m3u8_list.push(url)
            }
         }
    }
    return m3u8_list;
  });
  console.log(m3u8_list);
  await browser.close();
})();
