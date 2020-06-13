package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

var tplEditer *template.Template
var tplReport *template.Template

func init() {
	//テンプレート読み込み
	if t, err := template.ParseFiles("pages/editer.html", "pages/_daily-item.html"); err == nil {
		tplEditer = t
	} else {
		log.Fatalf("template error: %v", err)
	}
	if t, err := template.ParseFiles("pages/report.html"); err == nil {
		tplReport = t
	} else {
		log.Fatalf("template error: %v", err)
	}
}

//handleEdit 編集画面
func handleEdit(w http.ResponseWriter, r *http.Request) {

	//パラメータ解析
	r.ParseForm()
	form := r.Form

	//レポートの開始日を決定
	date := form.Get("start")
	var start time.Time
	if date != "" {
		var err error
		start, err = time.Parse("2006-01-02", date)
		if err != nil {
			return
		}
	} else {
		//今日開いたなら今週の月曜が起点のはず
		if n := int(time.Now().Weekday()); n == 0 {
			//日曜の場合は特殊な計算
			start = time.Now().AddDate(0, 0, -6)
		} else {
			start = time.Now().AddDate(0, 0, -n+1)
		}
	}

	//構造体を作る
	var data ReportData
	jsonpath := SearchJSON(start)
	if jsonpath != "" {
		//保存済みのJSONがあれば復元
		if d, err := RestoreJSON(jsonpath); err == nil {
			data = d
			data.CompleteOmitedParam()
		}
	} else {
		data = NewReportDataToday(start)
	}

	if err := tplEditer.Execute(w, data); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

//handleSave 保存処理
func handleSave(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var report ReportData
	if err := json.Unmarshal(body, &report); err != nil {
		log.Fatal(err)
	}

	SaveJSON(&report)
}

//handleViewReport レポート表示
func handleViewReport(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	form := r.Form

	start, err := time.Parse("2006-01-02", form.Get("start"))
	if err != nil {
		return
	}
	//構造体を作る
	var report ReportData
	jsonpath := SearchJSON(start)
	if jsonpath != "" {
		if d, err := RestoreJSON(jsonpath); err == nil {
			report = d
			report.CompleteOmitedParam()
		} else {
			return
		}
	} else {
		return
	}

	if err := tplReport.Execute(w, report); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func main() {
	port := "3000"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/report", handleViewReport)
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
