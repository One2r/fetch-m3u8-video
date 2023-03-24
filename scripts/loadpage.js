var page = require('webpage').create(),
  system = require('system'),
  address;

if (system.args.length === 1) {
  console.log('Usage: loadpage.js <some URL>');
  phantom.exit();
}

address = system.args[1];
page.settings.userAgent = 'Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html）)'
page.open(address, function(status) {
  if (status !== 'success') {
    console.log('FAIL to load the address');
  } else {
    m3u8 = page.evaluate(function() {
      var flashvars =  document.getElementsByTagName("embed")[0].getAttribute("flashvars");
      return unescape(flashvars.substr(5))
    });
    console.log("url:"+m3u8);
  }
  phantom.exit();
});
