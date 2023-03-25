const puppeteer = require('puppeteer');
var address = process.argv[2];

(async () => {
  const browser = await puppeteer.launch( {args: ["--no-sandbox"]});
  const page = await browser.newPage();
  await page.goto(address, {waitUntil: 'networkidle2'});
  
  const m3u8 = await page.evaluate(() => {
    var flash =  document.getElementsByTagName("embed")[0]
    if (flash){
        var flashvars = flash.getAttribute("flashvars");
        return unescape(flashvars.substr(5));
    }else{
        return "";
    }
  });
  console.log("url:"+m3u8);
  await browser.close();
})();
