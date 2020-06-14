package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//Save 保存
func Save(data ReportData) {
	go func() {
		// encode json
		data, _ := json.Marshal(data)
		res, err := http.Post("/save", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
	}()
}

//GetReportData データリクエスト
func GetReportData(start string) ReportData {
	r, err := http.Get("/api/report?start=" + start)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var report ReportData
	if err := json.Unmarshal(body, &report); err != nil {
		log.Fatal(err)
	}
	return report
}
