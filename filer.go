package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

//EnumJSON JSONデータの列挙
func EnumJSON() []string {
	pattern := "./data/*-*-*_*-*-*.json"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return []string{}
	}
	return files
}

//SearchJSON 保存済みデータの検索
func SearchJSON(startDate time.Time) string {
	pattern := "./data/" + startDate.Format("2006-01-02") + "_*.json"
	if files, err := filepath.Glob(pattern); err != nil {
		fmt.Println("EnumTempFiles_Error :", err.Error())
	} else {
		if len(files) > 0 {
			return files[0]
		}
	}
	return ""
}

//SaveJSON JSONファイルに保存する
func SaveJSON(jsonStruct *ReportData) error {
	filePath := fmt.Sprintf("./data/%s_%s.json", jsonStruct.Jisseki.Mon.DateValue, jsonStruct.Yotei.Sun.DateValue)
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer fp.Close()

	e := json.NewEncoder(fp)
	e.SetIndent("", "  ")
	if err := e.Encode(jsonStruct); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//RestoreJSON JSONから生成
func RestoreJSON(jsonpath string) (ReportData, error) {
	file, err := os.Open(jsonpath)
	if err != nil {
		msg := fmt.Sprintf("%s Open error : %v", jsonpath, err)
		fmt.Println(msg)
		return ReportData{}, err
	}
	defer file.Close()
	d := json.NewDecoder(file)
	d.DisallowUnknownFields() // エラーの場合 json: unknown field "JSONのフィールド名"
	var jsonstruct ReportData
	if err := d.Decode(&jsonstruct); err != nil && err != io.EOF {
		msg := fmt.Sprintf("%s Decode error : %v", jsonpath, err)
		fmt.Println(msg)
		return ReportData{}, err
	}
	return jsonstruct, nil
}

//getReportData 日付から判断してReportDataを作成する
func getReportData(startDate string) (ReportData, error) {
	var data ReportData
	var start time.Time
	if startDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return data, fmt.Errorf("日付の形式が不正です。yyyy-mm-ddで月曜日を指定してください")
		}
		if start.Weekday() != 1 {
			return data, fmt.Errorf("startには必ず月曜日を指定してください")
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
	return data, nil
}
