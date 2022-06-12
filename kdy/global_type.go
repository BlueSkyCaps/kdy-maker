package kdy

// Point 配置文件point.json对应的结构体类型
type Point struct {
	Classic    []string `json:"classic"`
	ClassicGif struct {
		X        int     `json:"x"`
		Y        int     `json:"y"`
		Interval float32 `json:"interval"`
	} `json:"classic-gif"`
}

// CurrentPoint 当前某个Point.Classic的具体数据，是Point.Classic序列中的一个
type CurrentPoint struct {
	LeftFlag bool    //左边纸有数据，否则为false
	LeftX    float64 //左边纸有数据横坐标，否则置为0
	LeftY    float64 //左边纸有数据纵坐标，否则置为0
	LeftSize float64 //左边纸有数据字体大小，否则置为0

	RightFlag bool    //右边纸有数据，否则为false
	RightX    float64 //右边纸有数据横坐标，否则置为0
	RightY    float64 //右边纸有数据纵坐标，否则置为0
	RightSize float64 //右边纸有数据字体大小，否则置为0
}

type Custom struct {
	RGB                 []float64
	FontPath            string
	LeftText, RightText string
}

// FontPathMap 字体存储路径
var FontPathMap = map[int]string{
	0: "C:/Windows/Fonts/simsun.ttc",
	1: "C:/Windows/Fonts/simkai.ttf",
	2: "C:/Windows/Fonts/SIMLI.TTF",
	3: "C:/Windows/Fonts/simhei.ttf",
	4: "C:/Windows/Fonts/SIMYOU.TTF"}

// RgbMap 预置的rgb颜色
var RgbMap = map[int][]float64{
	0: {0, 0, 0},       //纯黑色
	1: {0, 0.8, 0.4},   //亮绿色
	2: {0.8, 0.1, 0.1}} //暗红色
