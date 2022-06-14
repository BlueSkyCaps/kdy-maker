package kdy

import (
	"encoding/json"
	"github.com/BlueSkyCaps/GoGif/gof/img_op"
	"github.com/BlueSkyCaps/commGon"
	"github.com/fogleman/gg"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var orgFrameBasePath = path.Join("static", "kdy_frames")
var tmpFrameBasePath = path.Join("static", "tmp")
var finalOutBasePath = path.Join("static", "out")

var pointPath = path.Join("config", "point.json")
var orgFramePathsList []string
var tmpFramePathsList []string
var gifFrameNames []string
var customSel Custom
var point Point

func init() {
	// 事先static目录下没有tmp和out文件夹，创建
	info, err := os.Stat(tmpFrameBasePath)
	if info == nil && err != nil {
		commGon.CreateFolder(tmpFrameBasePath)
	}
	info, err = os.Stat(finalOutBasePath)
	if info == nil && err != nil {
		commGon.CreateFolder(finalOutBasePath)
	}

	orgFramesDir, _ := os.ReadDir(orgFrameBasePath)
	var tmpNames []string
	for _, ele := range orgFramesDir {
		tmpNames = append(tmpNames, ele.Name())
	}
	commGon.SortStringSlice(tmpNames, false)
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
func RunMaker(custom Custom) {
	// 获取输入了的自定义选项
	customSel = custom
	// 开始读取图片对象
	readGifHandler()
}
func readGifHandler() {
	// 读出配置文件
	jsonBytes, err := ioutil.ReadFile(pointPath)
	if err != nil {
		commGon.DebugPrint(err)
	}
	//反序列化json为结构体
	err = json.Unmarshal(jsonBytes, &point)
	if err != nil {
		commGon.DebugPrint(err)
	}

	for i, mPath := range orgFramePathsList {
		// 加载当前迭代的帧配置数据 帧图一一对应json配置序列元素
		currentPoint := currentPointIndexLoad(i, point)
		// 加载当前帧原始图
		currentImage, err := gg.LoadPNG(mPath)
		if err != nil {
			commGon.DebugPrint(err)
		}
		gifPtr := gg.NewContextForImage(currentImage)
		// 设置rgb颜色 值为0到1，原始数值/255得到相对的rgb值
		gifPtr.SetRGB(customSel.RGB[0], customSel.RGB[1], customSel.RGB[2])
		if currentPoint.LeftFlag {
			/* 开始渲染左边*/
			if err := gifPtr.LoadFontFace(customSel.FontPath, currentPoint.LeftSize); err != nil {
				commGon.DebugPrint(err)
			}
			gifPtr.DrawString(customSel.LeftText, currentPoint.LeftX, currentPoint.LeftY)
		}
		if currentPoint.RightFlag {
			/* 开始渲染右边*/
			if err := gifPtr.LoadFontFace(customSel.FontPath, currentPoint.RightSize); err != nil {
				commGon.DebugPrint(err)
			}
			gifPtr.DrawString(customSel.RightText, currentPoint.RightX, currentPoint.RightY)
		}
		// 保存当前帧效果图
		sErr := gifPtr.SavePNG(path.Join(tmpFrameBasePath, strconv.Itoa(i+1)+".png"))
		if sErr != nil {
			commGon.DebugPrint(err)
		}
	}
	// 开始转换成gif歌格式
	img_op.ConvertToGif(tmpFramePathsList, finalOutBasePath)
	// 开始生成最终gif动图
	var outSize = img_op.Size{X: point.ClassicGif.X, Y: point.ClassicGif.Y}
	outName := img_op.OpGifFileToGifDone(finalOutBasePath, gifFrameNames, finalOutBasePath, outSize, point.ClassicGif.Interval)
	// 将动图移到Windows桌面
	dirClear(outName)
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

func dirClear(outName string) {
	homeDir, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		from, _ := syscall.UTF16PtrFromString(path.Join(finalOutBasePath, outName))
		to, _ := syscall.UTF16PtrFromString(path.Join(homeDir, "desktop", outName))
		// 调用winApi移动文件 os.Rename无法跨卷移动（或改用创建写入方式）
		err := syscall.MoveFile(from, to)
		if err != nil {
			commGon.DebugPrint(err)
		}
	} else {
		err := os.Rename(path.Join(finalOutBasePath, outName), path.Join(homeDir, "desktop", outName))
		if err != nil {
			commGon.DebugPrint(err)
		}
	}
	commGon.RemoveFolderChildren(tmpFrameBasePath)
	commGon.RemoveFolderChildren(finalOutBasePath)

}
