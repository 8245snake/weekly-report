package main

import (
	"fmt"
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

func handleIndex(w http.ResponseWriter, r *http.Request) {

	//今日開いたなら今週の月曜が起点のはず
	var start time.Time
	if n := int(time.Now().Weekday()); n == 0 {
		//日曜の場合は特殊な計算
		start = time.Now().AddDate(0, 0, 6)
	} else {
		start = time.Now().AddDate(0, 0, -n+1)
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

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form
	types := [...]ContentType{ContentTypeJisseki, ContentTypeYotei}
	suffixes := [...]string{"date", "chk", "subtxt", "txt"}

	report := NewReportData()
	for _, contentstype := range types {
		for _, week := range Weeks {
			for _, suffix := range suffixes {
				key := fmt.Sprintf("%s-%s-%s", contentstype, week, suffix)
				value := form.Get(key)
				report.SetParam(contentstype, week, suffix, value)
			}
		}
	}
	report.Title = form.Get("title")
	report.Tasks = form.Get("tasks")
	report.Schedule = form.Get("schedule")
	SaveJSON(&report)

	if err := tplReport.Execute(w, report); err != nil {
		log.Printf("failed to execute template: %v", err)
	}

}

func main() {
	port := "3000"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/submit", handleSubmit)
	log.Printf("Server listening on port %s", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
