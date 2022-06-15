package main

import (
	"bufio"
	"fmt"
	hook "github.com/robotn/gohook"
	"kdy-maker/kdy"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	go stdInputHookBgListenExit()
	fmt.Println("LOOK!! 可达鸭表情包生成！选择输入并回车, 按Esc键退出O(∩_∩)O...")
	fmt.Println()
	for true {
		kdy.RunMaker(readCustom())
		fmt.Println(">>> done\n>>> 完成，位于Windows桌面。")
		fmt.Println()
	}
}

func readCustom() kdy.Custom {
	var customSel kdy.Custom
	var c, f int
	in := bufio.NewReader(os.Stdin)
	fmt.Println("输入可达鸭左手边和右手边文字[用中(英)文逗号隔开，两边各不超过四个字数]：")

	for true {
		text, e := in.ReadString('\n')
		// ReadString读\n结束并接收\n，此处去除最后的\n,windows是\r\n
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")
		// 最多8+1空格=9字符
		if e != nil || utf8.RuneCountInString(text) > 9 {
			fmt.Println("输入有误哦，重新再输一次吧：")
			continue
		}
		// 按中英文逗号','、'，'分割
		texts := strings.FieldsFunc(text, func(r rune) bool {
			return r == '，' || r == ','
		})
		if len(texts) != 2 || utf8.RuneCountInString(texts[0]) > 4 || utf8.RuneCountInString(texts[1]) > 4 {
			fmt.Println("输入有误哦，重新再输一次吧：")
			continue
		}
		customSel.LeftText = texts[0]
		customSel.RightText = texts[1]
		break
	}
	fmt.Println("颜色选择，0纯黑色(默认) 1亮绿色 2暗红色")
	_, e := fmt.Scanf("%d\n", &c)
	if e != nil {
		println("无效输入,已选纯黑色(默认)")
		c = 0
	}

	fmt.Println("字体选择 0宋体(默认) 1楷体 2隶书 3黑体 4幼圆")
	_, e = fmt.Scanf("%d\n", &f)
	if e != nil {
		println("无效输入,已选宋体(默认)")
		f = 0
	}
	customSel.RGB = kdy.RgbMap[c]
	customSel.FontPath = kdy.FontPathMap[f]
	return customSel
}

func stdInputHookBgListenExit() {
	evChan := hook.Start()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown && ev.Keychar == 27 {
			os.Exit(0)
		}
	}
}
