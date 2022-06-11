package kdy

// Point 配置文件point.json对应的结构体类型
type Point struct {
	Classic []string `json:"classic"`
}

// CurrentPoint 当前某个Point的具体数据，是Point序列中的一个
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
