package main

import (
	"sort"
	"strings"
)

//Summary 本体
type Summary struct {
	List SummaryItems
}

//SummaryItem サマリ
type SummaryItem struct {
	Start        string
	Title        string
	LastUpdate   string
	JissekiRange string
	YoteiRange   string
}

//SummaryItems ソート用
type SummaryItems []SummaryItem

func (p SummaryItems) Len() int {
	return len(p)
}

func (p SummaryItems) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SummaryItems) Less(i, j int) bool {
	return p[i].Start < p[j].Start
}

//Add 追加
// func (p *SummaryItems) Add(elem SummaryItem) {
// 	p = append(p, elem)
// }

//LoadSummary JSONファイルを読み込んでサマリを作成する
func LoadSummary() Summary {
	var summary Summary
	pathes := EnumJSON()
	for _, path := range pathes {
		var item SummaryItem
		data, err := RestoreJSON(path)
		if err != nil {
			continue
		}
		item.Title = data.Title
		item.Start = data.Jisseki.Mon.DateValue
		item.LastUpdate = data.LastUpdate
		item.JissekiRange = strings.Replace(data.Jisseki.Mon.DateValue+"～"+data.Jisseki.Sun.DateValue, "-", "/", -1)
		item.YoteiRange = strings.Replace(data.Yotei.Mon.DateValue+"～"+data.Yotei.Sun.DateValue, "-", "/", -1)
		summary.List = append(summary.List, item)
	}
	sort.Sort(summary.List)
	sort.Sort(sort.Reverse(summary.List))
	return summary
}
