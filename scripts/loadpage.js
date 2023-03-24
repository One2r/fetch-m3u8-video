const puppeteer = require('puppeteer');
address = system.args[1];
(async () => {
  const browser = await puppeteer.launch( {args: ["--no-sandbox"]});
  const page = await browser.newPage();
  await page.goto(address, {waitUntil: 'networkidle2'});
  
  const m3u8 = await page.evaluate(() => {
    var flashvars =  document.getElementsByTagName("embed")[0].getAttribute("flashvars");
    return unescape(flashvars.substr(5))
  });
  console.log("url:"+m3u8);
  await browser.close();
})();
