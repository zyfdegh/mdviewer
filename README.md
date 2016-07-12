# mdviewer
![Mdviewer logo](https://raw.githubusercontent.com/zyfdegh/mdviewer/master/raw/golang-markdown-viewer.png)

Mdviewer(mdv) is a markdown server, it displays markdown files in your broswer.

# Run

```sh
git clone https://github.com/zyfdegh/mdviewer.git
go get github.com/russross/blackfriday
cd mdviewer && ./build.sh
cd bin && ./mdv README.md
```

Open your favourite web browser and type:
```sh
127.0.0.1:8080
```

**bin/README.md** should be well styled in web page.

# Thanks

This project was mainly inspired by [ryanuber/readme-server][1]

# Related

**Call-web-broswer**: Call system default web browser to open a url [zyfdegh/call-web-broswer][2]

[1]:https://github.com/ryanuber/readme-server
[2]:https://github.com/zyfdegh/call-web-broswer