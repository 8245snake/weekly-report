// +build wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

//checkHolyday 休日チェックを入れたときの処理
func checkHolyday(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	checked := document.Call("getElementById", args[0].String()+"chk").Get("checked").Bool()
	if checked {
		document.Call("getElementById", args[0].String()+"container").Call("setAttribute", "class", "col-lg-2")
		document.Call("getElementById", args[0].String()+"date-disp").Call("setAttribute", "class", "h6")
	} else {
		document.Call("getElementById", args[0].String()+"container").Call("setAttribute", "class", "col-lg-4")
		document.Call("getElementById", args[0].String()+"date-disp").Call("setAttribute", "class", "h3")
	}

	return nil
}

func savePage(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	types := [...]ContentType{ContentTypeJisseki, ContentTypeYotei}
	suffixes := [...]string{"chk", "date", "subtxt", "txt"}

	report := NewReportData()
	for _, contentstype := range types {
		for _, week := range Weeks {
			for _, suffix := range suffixes {
				key := fmt.Sprintf("%s-%s-%s", contentstype, week, suffix)
				var value string = ""
				if suffix == "chk" {
					if document.Call("getElementById", key).Get("checked").Bool() {
						value = "checked"
					}
				} else {
					value = document.Call("getElementById", key).Get("value").String()
				}
				report.SetParam(contentstype, week, suffix, value)
			}
		}
	}
	report.Title = document.Call("getElementById", "title").Get("value").String()
	report.Tasks = document.Call("getElementById", "tasks").Get("value").String()
	report.Schedule = document.Call("getElementById", "schedule").Get("value").String()
	report.LastUpdate = time.Now().Format("2006/01/02 15:04:05")
	Save(report)
	return nil
}

func previewReport(this js.Value, args []js.Value) interface{} {
	//デッドロック回避のため無名関数のgorutineで実行する
	go func() {
		start := args[0].String()
		report := GetReportData(start)
		document := js.Global().Get("document")
		//プレビューを入れる
		text := ApplyTemplate(tamplateReport, report)
		document.Call("getElementById", "preview").Set("value", text)

		//隠しコントロールにキーを格納
		document.Call("getElementById", "start").Set("value", start)
	}()
	return nil
}

//registerCallbacks 関数をjavascript側で使えるように登録する
func registerCallbacks() {
	js.Global().Set("checkHolyday", js.FuncOf(checkHolyday))
	js.Global().Set("savePage", js.FuncOf(savePage))
	js.Global().Set("previewReport", js.FuncOf(previewReport))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
