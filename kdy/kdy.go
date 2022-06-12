package kdy

import (
	"encoding/json"
	"github.com/BlueSkyCaps/GoGif/gof/img_op"
	"github.com/fogleman/gg"
	"io/ioutil"
	"kdy-maker/common"
	"os"
	"path"
	"strconv"
	"strings"
)

var orgFrameBasePath = path.Join("static", "kdy_frames")
var tmpFrameBasePath = path.Join("static", "tmp")
var finalOutBasePath = path.Join("static", "out")

var pointPath = path.Join("config", "point.json")
var orgFramePathsList []string
var tmpFramePathsList []string
var gifFrameNames []string
var point Point

func init() {
	orgFramesDir, _ := os.ReadDir(orgFrameBasePath)
	var tmpNames []string
	for _, ele := range orgFramesDir {
		tmpNames = append(tmpNames, ele.Name())
	}
	common.SortStringSlice(tmpNames, false)
	// 原始帧图在 "static/kdy_frames/" 目录下
	for _, n := range tmpNames {
		orgFramePathsList = append(orgFramePathsList, path.Join(orgFrameBasePath, n))
	}
	// 生成后的帧图在 "static\tmp\" 目录下 需要将"/"替换成windows分隔符"\" 否则GoGif包会寻找字符得到错误文件名
	for _, n := range tmpNames {
		replacedSpPath := strings.Replace(path.Join(tmpFrameBasePath, n), "/", "\\", -1)
		tmpFramePathsList = append(tmpFramePathsList, replacedSpPath)
	}
	// "static\tmp\"
	tmpFrameBasePath = strings.Replace(path.Join(tmpFrameBasePath), "/", "\\", -1)
	// "static\out\"
	finalOutBasePath = strings.Replace(path.Join(finalOutBasePath), "/", "\\", -1)
	// 文件名x.png改为x.gif 用于生成gif动图 gifFrameNames存储的是已经形成了的每一张gif格式的图片
	for _, n := range tmpNames {
		gifFrameNames = append(gifFrameNames, strings.Replace(n, ".png", ".gif", 1))
	}

}
func RunMaker() {
	// 开始读取图片对象
	readGifHandler()
}
func readGifHandler() {
	// 读出配置文件
	jsonBytes, err := ioutil.ReadFile(pointPath)
	if err != nil {
		common.DebugPrint(err)
	}
	//反序列化json为结构体
	err = json.Unmarshal(jsonBytes, &point)
	if err != nil {
		common.DebugPrint(err)
	}

	for i, mPath := range orgFramePathsList {
		// 加载当前迭代的帧配置数据 帧图一一对应json配置序列元素
		currentPoint := currentPointIndexLoad(i, point)
		// 加载当前帧原始图
		currentImage, err := gg.LoadPNG(mPath)
		if err != nil {
			common.DebugPrint(err)
		}
		gifPtr := gg.NewContextForImage(currentImage)
		// 设置rgb颜色 值为0到1，原始数值/255得到相对的rgb值
		gifPtr.SetRGB(0, 0.8, 0.4) //默认亮绿色
		if currentPoint.LeftFlag {
			/* 开始渲染左边*/
			if err := gifPtr.LoadFontFace("C:/Windows/Fonts/simsun.ttc", currentPoint.LeftSize); err != nil {
				common.DebugPrint(err)
			}
			gifPtr.DrawString("拒绝", currentPoint.LeftX, currentPoint.LeftY)
		}
		if currentPoint.RightFlag {
			/* 开始渲染右边*/
			if err := gifPtr.LoadFontFace("C:/Windows/Fonts/simsun.ttc", currentPoint.RightSize); err != nil {
				common.DebugPrint(err)
			}
			gifPtr.DrawString("加班", currentPoint.RightX, currentPoint.RightY)
		}
		// 保存当前帧效果图
		sErr := gifPtr.SavePNG(path.Join(tmpFrameBasePath, strconv.Itoa(i+1)+".png"))
		if sErr != nil {
			common.DebugPrint(err)
		}
	}
	// 开始转换成gif歌格式
	img_op.ConvertToGif(tmpFramePathsList, finalOutBasePath)
	// 开始生成最终gif动图

	var outSize = img_op.Size{X: 240, Y: 238}
	img_op.OpGifFileToGifDone(finalOutBasePath, gifFrameNames, finalOutBasePath, outSize, 0.2)
}

// 加载当前帧对应的清单数据到CurrentPoint结构体中
func currentPointIndexLoad(index int, point Point) CurrentPoint {
	var currentPoint CurrentPoint
	// 分割，第一个元素是当前帧左边纸数据 第二个元素是当前帧右边纸数据
	p := strings.Split(point.Classic[index], ",")
	if p[0] != "n" {
		// 分割，分别是当前帧左边纸横坐标 纵坐标 字体大小
		xys := strings.Split(p[0], "-")
		currentPoint.LeftFlag = true
		currentPoint.LeftX, _ = strconv.ParseFloat(xys[0], 64)
		currentPoint.LeftY, _ = strconv.ParseFloat(xys[1], 64)
		currentPoint.LeftSize, _ = strconv.ParseFloat(xys[2], 64)
	} else {
		currentPoint.LeftFlag = false
	}
	if p[1] != "n" {
		// 分割，分别是当前帧右边纸横坐标 纵坐标 字体大小
		xys := strings.Split(p[1], "-")
		currentPoint.RightFlag = true
		currentPoint.RightX, _ = strconv.ParseFloat(xys[0], 64)
		currentPoint.RightY, _ = strconv.ParseFloat(xys[1], 64)
		currentPoint.RightSize, _ = strconv.ParseFloat(xys[2], 64)
	} else {
		currentPoint.RightFlag = false
	}
	return currentPoint
}
