package main

import (
	"bytes"
	"encoding/json"
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

//SaveAndView 保存
func SaveAndView(data ReportData) {
	go func() {
		// encode json
		data, _ := json.Marshal(data)
		res, err := http.Post("/report", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
	}()
}
