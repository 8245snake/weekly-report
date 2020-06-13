// +build !wasm

package main

import (
	"fmt"
)

func main() {
	fmt.Println("aeouo")
	data := NewReportData()
	Save(data)
}
