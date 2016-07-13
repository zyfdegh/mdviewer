# mdviewer

![BuildStatus](https://travis-ci.org/zyfdegh/mdviewer.svg)

![Mdviewer logo](https://raw.githubusercontent.com/zyfdegh/mdviewer/master/raw/golang-markdown-viewer.png)

Mdviewer(mdv) is a markdown server, it displays markdown files in your broswer.

# Run

```sh
git clone https://github.com/zyfdegh/mdviewer.git
go get github.com/russross/blackfriday
cd mdviewer && ./build.sh
cd bin && ./mdv ../README.md
```

Open your favourite web browser and type:
```sh
127.0.0.1:8080
```

**README.md** should be styled and displayed in the web page.

# Add to PATH
If you feel good and want to use 'mdv' command in system wide, copy bin/* to $PATH.

For Unix/Linux/MacOS:
```sh
cp bin/* /usr/local/bin
```

For Windows:
```cygwin
copy bin\* C:\Windows\System32\
```

# For better experience
Install broswer plugin like [Auto-Refresh][3] to refresh web page automatically if you're using Google Chrome.

# Thanks

This project was mainly inspired by [ryanuber/readme-server][1]

# Related

**Call-web-broswer**: Call system default web browser to open a url [zyfdegh/call-web-broswer][2].

[1]:https://github.com/ryanuber/readme-server
[2]:https://github.com/zyfdegh/call-web-broswer
[3]:https://chrome.google.com/webstore/detail/auto-refresh/ifooldnmmcmlbdennkpdnlnbgbmfalko?utm_source=chrome-app-launcher-info-dialog
