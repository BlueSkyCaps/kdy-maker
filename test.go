package main

import (
	"fmt"
	"strconv"
)

func test() {
	fa, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", 1.223), 32)
	fa = fa * float64(100)
	ia := int(fa)
	print(ia)
}
