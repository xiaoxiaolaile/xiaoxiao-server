### 学习傻妞的功能
[参考地址](https://github.com/cdle/sillyGirl)

### 前端界面地址
[https://github.com/xiaoxiaolaile/xiaoxiao-web](https://github.com/xiaoxiaolaile/xiaoxiao-web)

### 使用
先编译前端， 下载代码，安装node依赖 npm run build

生成的dist目录下的文件拷贝到 internal/core/static的目录下

然后可以编译go的代码

执行 go build (编译到不同环境可以网上搜索)

### 使用和傻妞一样的数据库

可以找到傻妞的sillyGirl.cache文件，放到编译好的二进制文件，执行后打开web界面。默认是localhost:8080，之后的操作和傻妞类似。

### 遇到的问题
傻妞的订阅功能不知道怎么实现，还有定时（自我感觉没有用，没有实现）



