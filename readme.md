# go-ueditor 插件使用
### 基于goframe框架的ueditor插件，支持图片上传、文件上传等功能。
### goframe 2.0     https://goframe.org/display/gf
### ueditor plus 3.9    https://open-doc.modstart.com/ueditor-plus/


## 1. 下载
```
go get github.com/topascend/go-gf-ueditor
```
## 2. 在cmd.go中引入插件
```
package cmd

import (
	// 百度编辑器插件
	ueditorController "github.com/topascend/go-gf-ueditor/controller"
	
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

## 3. 自定义ueditor配置文件  config/ueditor.yaml
```
# 上传图片配置项
imageActionName: "uploadImage"
# 执行上传图片的action名称 
imageFieldName: "upfile"
# 提交的图片表单名称 
imageMaxSize: 2048000
# 上传大小限制，单位B 
imageAllowFiles: [ ".png", ".jpg", ".jpeg", ".gif", ".bmp" ]
# 上传图片格式显示 
imageCompressEnable: false # 是否压缩图片,默认是true 
imageCompressBorder: 1600 # 图片压缩最长边限制 
imageInsertAlign: "none" # 插入的图片浮动方式 
imageUrlPrefix: "" # 图片访问路径前缀 
imagePathFormat: "/ueditor/upload/images/{yyyy}{mm}{dd}/uploadImage-{yyyy}-{mm}-{dd}-{hh}-{ii}-{ss}-{rand:6}"  # 上传保存路径,可以自定义保存路径和文件名格式 
# {filename} 会替换成原文件名,配置这项需要注意中文乱码问题 
# {rand:6} 会替换成随机数,后面的数字是随机数的位数 
# {time} 会替换成时间戳 
# {yyyy} 会替换成四位年份 
# {yy} 会替换成两位年份 
# {mm} 会替换成两位月份 
# {dd} 会替换成两位日期 
# {hh} 会替换成两位小时 
# {ii} 会替换成两位分钟 
# {ss} 会替换成两位秒 
# 非法字符 \ : * ? " < > | 
# 具请体看线上文档: fex.baidu.com/ueditor/#use-format_upload_filename 

# 涂鸦图片上传配置项 
scrawlActionName: "uploadScrawl" # 执行上传涂鸦的action名称 
scrawlFieldName: "base64" # 提交的图片表单名称 
scrawlPathFormat: "/ueditor/upload/images/{yyyy}{mm}{dd}/uploadScrawl-{yyyy}-{mm}-{dd}-{hh}-{ii}-{ss}-{rand:6}" # 上传保存路径,可以自定义保存路径和文件名格式 
scrawlMaxSize: 2048000  #20*1024  # 上传大小限制，单位B 
scrawlUrlPrefix: ""  # 图片访问路径前缀 
scrawlInsertAlign: "none"

# 截图工具上传 
snapscreenActionName: "uploadImage" # 执行上传截图的action名称 
snapscreenPathFormat: "/ueditor/upload/images/{yyyy}{mm}{dd}/{yyyy}-{mm}-{dd}-{hh}-{ii}-{ss}-{rand:6}"  # 上传保存路径,可以自定义保存路径和文件名格式 
snapscreenUrlPrefix: ""  # 图片访问路径前缀 
snapscreenInsertAlign: "none"  # 插入的图片浮动方式 

# 抓取远程图片配置 
catcherLocalDomain: [ "127.0.0.1", "localhost", "img.baidu.com" ]
catcherActionName: "catchImage"  # 执行抓取远程图片的action名称 
catcherFieldName: "source" # 提交的图片列表表单名称 
catcherPathFormat: "/ueditor/upload/images/{yyyy}{mm}{dd}/catchImage-{yyyy}-{mm}-{dd}-{hh}-{ii}-{ss}-{rand:6}"  # 上传保存路径,可以自定义保存路径和文件名格式 
catcherUrlPrefix: ""  # 图片访问路径前缀 
catcherMaxSize: 2048000   # 上传大小限制，单位B 
catcherAllowFiles: [ ".png", ".jpg", ".jpeg", ".gif", ".bmp" ] # 抓取图片格式显示 

# 上传视频配置 
videoActionName: "uploadvideo" # 执行上传视频的action名称 
videoFieldName: "upfile"# 提交的视频表单名称
videoPathFormat: "/ueditor/upload/videos/{yyyy}{mm}{dd}/{yyyy}-{mm}-{dd}-{hh}-{ii}-{ss}-{rand:6}" # 上传保存路径,可以自定义保存路径和文件名格式
videoUrlPrefix: "" # 视频访问路径前缀
videoMaxSize: 1024000000 # 上传大小限制，单位B，默认1000MB 
videoAllowFiles: [ ".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid" ]  # 上传视频格式显示 

# 上传文件配置 
fileActionName: "uploadfile" # controller里,执行上传视频的action名称 
fileFieldName: "upfile" # 提交的文件表单名称 
filePathFormat: "/ueditor/upload/files/{yyyy}{mm}{dd}/{yyyy}-{mm}-{dd}-{ii}-{ss}-{rand:6}" # 上传保存路径,可以自定义保存路径和文件名格式 
fileUrlPrefix: ""  # 文件访问路径前缀 
fileMaxSize: 51200000  # 上传大小限制，单位B，默认50MB 
fileAllowFiles: [ ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid", ".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml", ]  # 上传文件格式显示

# 列出指定目录下的图片
imageManagerActionName: "listImage" # 执行图片管理的action名称
imageManagerListPath: "/ueditor/upload/images/" # 指定要列出图片的目录 
imageManagerListSize: 20  # 每次列出文件数量 
imageManagerUrlPrefix: ""  # 图片访问路径前缀 
imageManagerInsertAlign: "none"  # 插入的图片浮动方式 
imageManagerAllowFiles: [ ".png", ".jpg", ".jpeg", ".gif", ".bmp" ]  # 列出的文件类型 

# 列出指定目录下的文件 
fileManagerActionName: "listFile"               # 执行文件管理的action名称 
fileManagerListPath: "/ueditor/upload/files/" # 指定要列出文件的目录 
fileManagerUrlPrefix: ""  # 文件访问路径前缀 
fileManagerListSize: 20  # 每次列出文件数量 
fileManagerAllowFiles: [ ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid", ".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml", ] # 列出的文件类型 

```

## 4. 配置前端配置 ueditor.config.js
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