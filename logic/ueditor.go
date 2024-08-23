package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/topascend/go-gf-ueditor/api"
	"github.com/topascend/go-gf-ueditor/lib"
	"github.com/topascend/go-gf-ueditor/service"
	"math/rand"
	"regexp"

	"strconv"
	"strings"
	"time"
)

func init() {
	service.RegisterUeditor(New())
}

func New() *sUEditor {
	return &sUEditor{}
}

type sUEditor struct {
}

func (s *sUEditor) Action(ctx context.Context, req *api.UEditorReq) (res g.Map, err error) {
	switch req.Action {
	//获取配置
	case "config":
		err = s.uEditorConfig(ctx, req.Callback) //特殊处理返回结果

	//上传图片
	case Config(ctx, "imageActionName"):
		res, err = s.uEditorUpload(ctx, req)

	//上传视频  上传文件
	case Config(ctx, "fileActionName"), Config(ctx, "videoActionName"):
		res, err = s.uEditorUpload(ctx, req)

	//上传涂鸦
	case Config(ctx, "scrawlActionName"):
		res, err = s.uEditorScrawl(ctx, req)

	//列出图片  列出文件
	case Config(ctx, "imageManagerActionName"), Config(ctx, "fileManagerActionName"):
		res, err = s.uEditorList(ctx, req)

	//抓取远端图片
	case Config(ctx, "catcherActionName"):
		res, err = s.uEditorCatchImage(ctx, req)

	case Config(ctx, "snapscreenActionName"):
		err = gerror.New("暂不支持截图功能")

	default:
		err = gerror.New("发生错误")
	}

	return
}

// uEditorConfig 获取uEditor配置
func (s *sUEditor) uEditorConfig(ctx context.Context, callback string) (err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Write(fmt.Sprintf("%s(%s)", callback, gconv.String(config)))

	return
}

// uEditorUpload 上传图片
func (s *sUEditor) uEditorUpload(ctx context.Context, req *api.UEditorReq) (res g.Map, err error) {
	var (
		r      = g.RequestFromCtx(ctx)
		upFile = req.File
		host   = r.GetSchema() + "://" + r.Host

		pathFormat string
		allowFiles []string
		maxSize    int64
	)

	switch req.Action {
	case Config(ctx, "imageActionName"):
		pathFormat = gconv.String(Config(ctx, "imagePathFormat"))
		allowFiles = gconv.SliceStr(Config(ctx, "imageAllowFiles"))
		maxSize = gconv.Int64(Config(ctx, "imageMaxSize"))

	case Config(ctx, "videoActionName"):
		pathFormat = gconv.String(Config(ctx, "videoPathFormat"))
		allowFiles = gconv.SliceStr(Config(ctx, "videoAllowFiles"))
		maxSize = gconv.Int64(Config(ctx, "videoMaxSize"))

	case Config(ctx, "fileActionName"):
		pathFormat = gconv.String(Config(ctx, "filePathFormat"))
		allowFiles = gconv.SliceStr(Config(ctx, "fileAllowFiles"))
		maxSize = gconv.Int64(Config(ctx, "fileMaxSize"))
	}

	// 检查文件扩展名
	oldName := upFile.Filename
	err = s.CheckExt(gfile.Ext(oldName), allowFiles)
	if err != nil {
		return
	}

	// 检查文件大小
	err = s.CheckSize(upFile.Size, maxSize)
	if err != nil {
		return
	}

	// 根据配置规则创建文件名
	newName := s.CreateUEditorFileName(pathFormat, oldName)
	absDir := lib.PublicAbsPath(ctx) + gfile.Dir(newName)

	upFile.Filename = gfile.Basename(newName) // 重命名文件名
	_, err = upFile.Save(absDir, false)       // 保存文件
	url := host + newName

	if err != nil {
		return
	}

	res = g.Map{
		"state":    "SUCCESS",
		"url":      url,
		"title":    oldName,
		"original": oldName,
	}

	return
}

// uEditorScrawl 上传涂鸦
func (s *sUEditor) uEditorScrawl(ctx context.Context, req *api.UEditorReq) (res g.Map, err error) {
	var (
		r      = g.RequestFromCtx(ctx)
		absDir = lib.PublicAbsPath(ctx)
		host   = r.GetSchema() + "://" + r.Host
	)

	fileName := s.CreateUEditorFileName(gconv.String(Config(ctx, "scrawlPathFormat")), "temp.png")
	str, err := gbase64.Decode([]byte(gconv.String(req.Base64)))
	if err != nil {
		return
	}

	err = gfile.PutContents(absDir+fileName, string(str))
	if err != nil {
		return
	}

	url := host + fileName
	res = g.Map{
		"state":    "SUCCESS",
		"url":      url,
		"title":    fileName,
		"original": fileName,
	}
	return
}

// uEditorList 图片 文件列表
func (s *sUEditor) uEditorList(ctx context.Context, req *api.UEditorReq) (res g.Map, err error) {
	var (
		r     = g.RequestFromCtx(ctx)
		start = req.Start
		size  = req.Size

		absDir = lib.PublicAbsPath(ctx)
		host   = r.GetSchema() + "://" + r.Host
		end    = start + size

		listPath   = absDir
		allowFiles []string
	)

	// 图片列表  文件列表
	if req.Action == Config(ctx, "imageManagerActionName") {
		listPath = gconv.String(Config(ctx, "imageManagerListPath"))
		allowFiles = gconv.SliceStr(Config(ctx, "imageManagerAllowFiles"))
	} else if req.Action == Config(ctx, "fileManagerActionName") {
		listPath = gconv.String(Config(ctx, "fileManagerListPath"))
		allowFiles = gconv.SliceStr(Config(ctx, "fileManagerAllowFiles"))
	}

	allowExt := ""
	for _, ext := range allowFiles {
		allowExt += "*" + ext + ","
	}

	fileList, err := gfile.ScanDir(absDir+listPath, allowExt, true)
	if err != nil {
		return
	}

	absDir = strings.Replace(absDir, "\\", "/", -1)
	length := len(fileList)
	list := make([]map[string]string, 0)
	for i := min(length, end) - 1; i < length && i >= 0 && i >= start; i-- {
		mtime := gfile.MTimestamp(fileList[i])
		temp := strings.Replace(fileList[i], "\\", "/", -1)
		temp = strings.Replace(temp, absDir, host, 1)
		list = append(list, map[string]string{"url": temp, "mtime": strconv.FormatInt(mtime, 10)})
	}

	res = g.Map{
		"state": "SUCCESS",
		"list":  list,
		"start": start,
		"total": length,
	}

	return
}

// uEditorCatchImage 上传图片
func (s *sUEditor) uEditorCatchImage(ctx context.Context, req *api.UEditorReq) (res g.Map, allErr error) {
	var (
		r      = g.RequestFromCtx(ctx)
		source = req.Source

		absDir = lib.PublicAbsPath(ctx)
		host   = r.GetSchema() + "://" + r.Host
	)

	list := make([]g.MapStrStr, len(source), len(source))

	var err error
	for key, img := range source {
		ext := gfile.Ext(img)
		err = s.CheckExt(ext, gconv.SliceStr(Config(ctx, "catcherAllowFiles")))
		if err != nil {
			allErr = errors.Join(allErr, gerror.New(img+" "+err.Error()))
			continue
		}

		fileName := s.CreateUEditorFileName(gconv.String(Config(ctx, "catcherPathFormat")), img)
		absFile := absDir + fileName
		size, err := s.saveRemoteImage(ctx, img, absFile)
		if err != nil {
			allErr = errors.Join(allErr, gerror.New(img+" "+err.Error()))
			continue
		}

		url := host + fileName
		list[key] = g.MapStrStr{
			"state":    "SUCCESS",
			"url":      url,
			"size":     strconv.FormatInt(size, 10),
			"title":    fileName,
			"original": fileName,
			"source":   fileName,
		}
	}

	if allErr != nil {
		return
	}

	res = g.Map{
		"state": "SUCCESS",
		"list":  list,
	}

	return
}

// CheckExt 检查文件类型
func (s *sUEditor) CheckExt(ext string, AllowExt []string) (err error) {
	if !lib.Contain(AllowExt, strings.ToLower(ext)) {
		err = gerror.New("不支持的文件类型")
	}

	return
}

// CheckSize 检查文件大小
func (s *sUEditor) CheckSize(size int64, maxSize int64) (err error) {
	if size > maxSize {
		err = gerror.New("文件大小超出限制")
	}

	return
}

// saveRemoteImage 保存远程图片
func (s *sUEditor) saveRemoteImage(ctx context.Context, url, filename string) (size int64, err error) {
	var r *gclient.Response
	if r, err = g.Client().Get(ctx, url); err != nil {
		return
	}
	size = r.ContentLength

	defer r.Close()
	err = gfile.PutBytes(filename, r.ReadAll())
	if err != nil {
		return 0, err
	}

	return

	/*// 发送HTTP GET请求以获取图片
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	size = resp.ContentLength
	err = s.CheckSize(size, gconv.Int64(Config("catcherMaxSize")))
	if err != nil {
		return
	}

	absDir := gfile.Dir(filename)
	if !gfile.IsDir(absDir) {
		err = gfile.Mkdir(absDir)
	}

	// 创建一个文件用于写入
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()

	// 将响应体复制到文件
	_, err = io.Copy(out, resp.Body)

	return*/
}

// CreateUEditorFileName 解析并替换模板中的占位符
func (s *sUEditor) CreateUEditorFileName(template string, fileName string) string {
	ext := gfile.Ext(fileName)
	name := gfile.Name(fileName)

	now := time.Now()
	year := now.Format("2006")                     // 四位年份
	yearShort := now.Format("06")                  // 两位年份
	month := now.Format("01")                      // 两位月份
	day := now.Format("02")                        // 两位日期
	hour := now.Format("15")                       // 两位小时
	minute := now.Format("04")                     // 两位分钟
	second := now.Format("05")                     // 两位秒
	timestamp := strconv.FormatInt(now.Unix(), 10) // 时间戳

	rand2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	re := regexp.MustCompile(`\{rand:(\d+)\}`)
	matches := re.FindStringSubmatch(template)

	randomStr := ""
	if len(matches) >= 2 {
		num := gconv.Int(matches[1])
		for i := 0; i < num; i++ {
			randomStr += strconv.Itoa(rand2.Intn(10)) // 生成6位随机数
		}

		template = strings.ReplaceAll(template, matches[0], randomStr)
	}

	// 替换占位符
	template = strings.ReplaceAll(template, "{filename}", name)
	template = strings.ReplaceAll(template, "{yyyy}", year)
	template = strings.ReplaceAll(template, "{yy}", yearShort)
	template = strings.ReplaceAll(template, "{mm}", month)
	template = strings.ReplaceAll(template, "{dd}", day)
	template = strings.ReplaceAll(template, "{hh}", hour)
	template = strings.ReplaceAll(template, "{ii}", minute)
	template = strings.ReplaceAll(template, "{ss}", second)
	template = strings.ReplaceAll(template, "{time}", timestamp)
	//template = strings.ReplaceAll(template, matches[0], randomStr)

	return template + ext
}
