# Copy-files

###multithreaded copy file 
### 介绍
这个程序是在一次工作中对大量的视频文件做拷贝操作时写出来的，支持多线程，比单线程快5倍以上.其中数据源为一个后缀为xlsx格式的表格(吐槽下用的第三方模块不支持xls,各位童靴在使用的时候请自行装换为xlsx)
### BUILD
首先得将整个项目拷贝到本机的`$GOPATH/src`下

`go install zxywork`

### exec

`./zxywork 1.xlsx /temp/test`


