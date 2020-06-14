// +build !wasm

package main

import "fmt"

func main() {

	report := ReportData{}
	fmt.Printf("%v\n", report)

	dom := ApplyTemplate(tmplTest, "world")
	fmt.Printf("%s\n", dom)
}
