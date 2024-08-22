package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UEditorActionReq struct {
	g.Meta `path:"/ueditor/action" tags:"UEditor" method:"get,post" summary:"UEditor action操作"`
	UEditorReq
}

type UEditorReq struct {
	Action   string            `p:"action"`
	Callback string            `p:"callback"`
	File     *ghttp.UploadFile `p:"upfile" type:"file"`
	Base64   string            `p:"base64" `
	Source   []string          `p:"source"`
	callback string            `p:"callback"`
	Start    int               `p:"start"`
	Size     int               `p:"size"`
}

type UEditorRes struct {
	g.Meta `mime:"application/json"`
	g.Map  `json:"data"`
}

// IUEditor UEditor 接口
type IUEditor interface {
	Action(ctx context.Context, req *UEditorActionReq) (res *UEditorRes, err error)
}
