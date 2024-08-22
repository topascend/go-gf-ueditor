package controller

import (
	_ "github.com/topascend/go-gf-ueditor/logic"

	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/topascend/go-gf-ueditor/api"
	"github.com/topascend/go-gf-ueditor/service"
)

/*var (
	UEditor = uEditorController{}
)
type uEditorController struct {
}*/

type UEditorController struct{}

func NewUEditor() api.IUEditor {
	return &UEditorController{}
}

// Action go 实现的 action 操作方法
func (c *UEditorController) Action(ctx context.Context, req *api.UEditorActionReq) (res *api.UEditorRes, err error) {
	r := g.RequestFromCtx(ctx)
	data, err := service.Ueditor().Action(ctx, &req.UEditorReq)

	// 错误处理
	if err != nil {
		r.Response.WriteJson(g.Map{
			"state": err.Error(),
		})
		return
	}

	// jsonp callback 处理,主要用于前后端分离测试开发使用
	if req.Callback != "" && data != nil {
		r.Response.Write(fmt.Sprintf("%s(%s)", req.Callback, gjson.New(data).String()))
		return
	}

	// 正常返回
	if data != nil {
		r.Response.WriteJson(data)
	}

	return
}

// Ueditor html页面 测试使用
/*func (c *UEditorController) Ueditor(r *ghttp.Request) {
	_ = r.Response.WriteTpl("ueditor/index.html", g.Map{})
}*/
