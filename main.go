package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

var tplIndex *template.Template
var tplEditer *template.Template
var tplReport *template.Template

//TemplateParts テンプレート部品たち共通化
var TemplateParts = []string{
	"pages/_head.html",
	"pages/_navber.html",
	"pages/_daily-item.html",
	"pages/_weekly-item.html",
	"pages/_toast.html",
	"pages/_report-text.html",
	"pages/_modal-confirm.html",
}

func init() {
	//ホーム画面テンプレート
	if t, err := template.ParseFiles(append([]string{"pages/index.html"}, TemplateParts...)...); err == nil {
		tplIndex = t
	} else {
		log.Fatalf("template error: %v", err)
	}
	//編集画面テンプレート
	if t, err := template.ParseFiles(append([]string{"pages/editer.html"}, TemplateParts...)...); err == nil {
		tplEditer = t
	} else {
		log.Fatalf("template error: %v", err)
	}
	//レポート画面テンプレート
	if t, err := template.ParseFiles(append([]string{"pages/report.html"}, TemplateParts...)...); err == nil {
		tplReport = t
	} else {
		log.Fatalf("template error: %v", err)
	}
}

//handleIndex トップ画面
func handleIndex(w http.ResponseWriter, r *http.Request) {
	summary := LoadSummary()
	if err := tplIndex.Execute(w, summary); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

//handleEdit 編集画面
func handleEdit(w http.ResponseWriter, r *http.Request) {
	//パラメータ解析
	r.ParseForm()
	form := r.Form
	//レポートの開始日を決定
	startdate := form.Get("start")
	data, err := getReportData(startdate)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//レンダリング
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
	//パラメータ解析
	r.ParseForm()
	form := r.Form
	//レポートの開始日を決定
	startdate := form.Get("start")
	report, err := getReportData(startdate)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err := tplReport.Execute(w, report); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

//responseReportData 開始日をキーにしてReportDataを返すAPI
func responseReportData(w http.ResponseWriter, r *http.Request) {
	//パラメータ解析
	r.ParseForm()
	form := r.Form
	//レポートの開始日を決定
	startdate := form.Get("start")
	report, err := getReportData(startdate)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	jsindata, _ := json.Marshal(report)
	w.Write(jsindata)
}

func main() {
	port := "3000"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/report", handleViewReport)
	http.HandleFunc("/api/report", responseReportData)
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
