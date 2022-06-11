package main

import (
	"fmt"
	"os"
)

func test() {
	curr_wd, err := os.Getwd()

	if err != nil {

		fmt.Println(err)

		os.Exit(1)
	}

	// Print the current working directory
	fmt.Println("--d")
	fmt.Println(curr_wd)
	fmt.Println("--d")
}
