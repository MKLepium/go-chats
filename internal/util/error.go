package util

import "fmt"

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}
}
