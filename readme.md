# go-ueditor 插件使用
### 基于goframe框架的ueditor插件，支持图片上传、文件上传等功能。
### goframe 2.0     https://goframe.org/display/gf
### ueditor plus 3.9    https://open-doc.modstart.com/ueditor-plus/


## 1. 复制plugins/ueditor到项目目录下
## 2. 在cmd.go中引入插件
```
go get github.com/topascend/go-gf-ueditor
```


```
package cmd

import (
	// 百度编辑器插件
	ueditorController "UEditorGoFrame/internal/plugins/ueditor/controller"
	
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)
var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 绑定ueditor插件
				group.Bind(
					ueditorController.NewUEditor(),
				)
			})
			s.Run()
			return nil
		},
	}
)
```
## 3. 配置前端配置 ueditor.config.js
```
    // 为编辑器实例添加一个路径，这个不能被注释
    UEDITOR_HOME_URL: URL,
    // 需要能跨域的静态资源请求，主要用户弹窗页面等静态资源
    UEDITOR_CORS_URL: CORS_URL,

    // 是否开启Debug模式
    debug: true,

    // 服务器统一请求接口路径
    //serverUrl: "/ueditor-plus/_demo_server/handle.php",
    serverUrl: "/ueditor/action",
    // 服务器统一请求头信息，会在所有请求中带上该信息
    
    // 服务器返回参数统一转换方法，可以在这里统一处理返回参数
    serverResponsePrepare: function( res ){
        console.log('serverResponsePrepare', res);
        return res;
    },
```
### 服务器统一请求接口路径为 /ueditor/action
### 测试功能地址为  http://127.0.0.1:8000/ueditor