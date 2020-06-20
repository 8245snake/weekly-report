package main

import (
	"fmt"
	"time"
)

//ContentType 予定or実績
type ContentType string

const (
	//DateValueFormat キー日付のフォーマット
	DateValueFormat string = "2006-01-02"
	//ContentTypeJisseki 実績
	ContentTypeJisseki ContentType = "jisseki"
	//ContentTypeYotei 予定
	ContentTypeYotei ContentType = "yotei"
	//PlaceholderJisseki プレースホルダ（実績）
	PlaceholderJisseki string = "実績を入力する"
	//PlaceholderYotei プレースホルダ（予定）
	PlaceholderYotei string = "予定を入力する"
)

//ControlData ページ作成用データ
type ControlData struct {
	PlaceholderJisseki string
	PlaceholderYotei   string
}

//ReportData 週報データ全て
type ReportData struct {
	Jisseki     WeeklyContents `json:"jisseki"`
	Yotei       WeeklyContents `json:"yotei"`
	Title       string         `json:"title"`
	Tasks       string         `json:"task"`
	Schedule    string         `json:"schedule"`
	LastUpdate  string         `json:"last_update"`
	ControlData ControlData    `json:"-"`
}

//WeeklyContents 週毎のデータ
type WeeklyContents struct {
	Mon *DailyItem `json:"mon"`
	Tue *DailyItem `json:"tue"`
	Wed *DailyItem `json:"wed"`
	Thu *DailyItem `json:"thu"`
	Fri *DailyItem `json:"fri"`
	Sat *DailyItem `json:"sat"`
	Sun *DailyItem `json:"sun"`
}

//DailyItem 日毎のデータ
type DailyItem struct {
	Type      ContentType `json:"type"`
	IsHolyday bool        `json:"is_holyday"`
	Prefix    string      `json:"prefix,omitempty"`
	DateValue string      `json:"date"`
	DispDate  string      `json:"disp_date,omitempty"`
	Text      string      `json:"text"`
	SubText   string      `json:"subtext"`
}

//Weeks 週の名前
var Weeks = [7]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

//WeeksJP 週の名前（日本語）
var WeeksJP = [7]string{"(日)", "(月)", "(火)", "(水)", "(木)", "(金)", "(土)"}

//NewReportData 初期化
func NewReportData() ReportData {
	data := ReportData{
		Jisseki: WeeklyContents{
			Mon: &DailyItem{Type: ContentTypeJisseki},
			Tue: &DailyItem{Type: ContentTypeJisseki},
			Wed: &DailyItem{Type: ContentTypeJisseki},
			Thu: &DailyItem{Type: ContentTypeJisseki},
			Fri: &DailyItem{Type: ContentTypeJisseki},
			Sat: &DailyItem{Type: ContentTypeJisseki},
			Sun: &DailyItem{Type: ContentTypeJisseki},
		},
		Yotei: WeeklyContents{
			Mon: &DailyItem{Type: ContentTypeYotei},
			Tue: &DailyItem{Type: ContentTypeYotei},
			Wed: &DailyItem{Type: ContentTypeYotei},
			Thu: &DailyItem{Type: ContentTypeYotei},
			Fri: &DailyItem{Type: ContentTypeYotei},
			Sat: &DailyItem{Type: ContentTypeYotei},
			Sun: &DailyItem{Type: ContentTypeYotei},
		},
	}
	data.ControlData.PlaceholderJisseki = PlaceholderJisseki
	data.ControlData.PlaceholderYotei = PlaceholderYotei
	return data
}

//NewReportDataToday startDateを起点として2週間分の週報データを作成する
func NewReportDataToday(startDate time.Time) ReportData {
	jisseki := WeeklyContents{
		Mon: NewDailyItem(ContentTypeJisseki, false, startDate),
		Tue: NewDailyItem(ContentTypeJisseki, false, startDate.AddDate(0, 0, 1)),
		Wed: NewDailyItem(ContentTypeJisseki, false, startDate.AddDate(0, 0, 2)),
		Thu: NewDailyItem(ContentTypeJisseki, false, startDate.AddDate(0, 0, 3)),
		Fri: NewDailyItem(ContentTypeJisseki, false, startDate.AddDate(0, 0, 4)),
		Sat: NewDailyItem(ContentTypeJisseki, true, startDate.AddDate(0, 0, 5)),
		Sun: NewDailyItem(ContentTypeJisseki, true, startDate.AddDate(0, 0, 6)),
	}
	yotei := WeeklyContents{
		Mon: NewDailyItem(ContentTypeYotei, false, startDate.AddDate(0, 0, 7)),
		Tue: NewDailyItem(ContentTypeYotei, false, startDate.AddDate(0, 0, 8)),
		Wed: NewDailyItem(ContentTypeYotei, false, startDate.AddDate(0, 0, 9)),
		Thu: NewDailyItem(ContentTypeYotei, false, startDate.AddDate(0, 0, 10)),
		Fri: NewDailyItem(ContentTypeYotei, false, startDate.AddDate(0, 0, 11)),
		Sat: NewDailyItem(ContentTypeYotei, true, startDate.AddDate(0, 0, 12)),
		Sun: NewDailyItem(ContentTypeYotei, true, startDate.AddDate(0, 0, 13)),
	}
	layout := "2006/01/02"
	title := fmt.Sprintf("【】%s～%s", startDate.Format(layout), startDate.AddDate(0, 0, 6).Format(layout))
	data := NewReportData()
	data.Jisseki = jisseki
	data.Yotei = yotei
	data.Title = title
	return data
}

//NewDailyItem 初期化
func NewDailyItem(contentType ContentType, isHolyday bool, date time.Time) *DailyItem {
	data := &DailyItem{}
	data.Type = contentType
	data.DateValue = date.Format(DateValueFormat)
	data.Prefix = fmt.Sprintf("%s-%s-", contentType, Weeks[date.Weekday()])
	data.IsHolyday = isHolyday
	data.DispDate = date.Format("2006/01/02") + WeeksJP[date.Weekday()]
	return data
}

//SetParam データ代入
//suffixes := [...]string{"date", "chk", "subtxt", "txt"}
func (item *DailyItem) SetParam(suffix string, value string) {
	switch suffix {
	case "date":
		item.DateValue = value
	case "chk":
		item.IsHolyday = (value != "")
	case "subtxt":
		item.SubText = value
		if value == "" {
			if item.IsHolyday {
				item.SubText = "休日"
			} else {
				item.SubText = "社内"
			}
		} else {
			item.SubText = value
		}
	case "txt":
		item.Text = value
	}
}

//CompleteOmitedParam 省略されたパラメータを復活
func (item *DailyItem) CompleteOmitedParam() {

	date, err := time.Parse(DateValueFormat, item.DateValue)
	if err != nil {
		return
	}
	var dailyItem *DailyItem
	if date.Weekday() == 0 || date.Weekday() == 6 {
		// 休日
		dailyItem = NewDailyItem(item.Type, true, date)
	} else {
		// 平日
		dailyItem = NewDailyItem(item.Type, false, date)
	}
	item.Prefix = dailyItem.Prefix
	item.DispDate = dailyItem.DispDate
}

//CompleteOmitedParam 省略されたパラメータを復活
func (data *ReportData) CompleteOmitedParam() {
	data.Jisseki.Sun.CompleteOmitedParam()
	data.Jisseki.Mon.CompleteOmitedParam()
	data.Jisseki.Tue.CompleteOmitedParam()
	data.Jisseki.Wed.CompleteOmitedParam()
	data.Jisseki.Thu.CompleteOmitedParam()
	data.Jisseki.Fri.CompleteOmitedParam()
	data.Jisseki.Sat.CompleteOmitedParam()
	data.Yotei.Sun.CompleteOmitedParam()
	data.Yotei.Mon.CompleteOmitedParam()
	data.Yotei.Tue.CompleteOmitedParam()
	data.Yotei.Wed.CompleteOmitedParam()
	data.Yotei.Thu.CompleteOmitedParam()
	data.Yotei.Fri.CompleteOmitedParam()
	data.Yotei.Sat.CompleteOmitedParam()
	data.ControlData.PlaceholderJisseki = PlaceholderJisseki
	data.ControlData.PlaceholderYotei = PlaceholderYotei
}

//SetParam データ代入
func (data *ReportData) SetParam(tp ContentType, week string, suffix string, value string) {
	switch tp {
	case ContentTypeJisseki:
		switch week {
		case Weeks[0]:
			data.Jisseki.Sun.SetParam(suffix, value)
		case Weeks[1]:
			data.Jisseki.Mon.SetParam(suffix, value)
		case Weeks[2]:
			data.Jisseki.Tue.SetParam(suffix, value)
		case Weeks[3]:
			data.Jisseki.Wed.SetParam(suffix, value)
		case Weeks[4]:
			data.Jisseki.Thu.SetParam(suffix, value)
		case Weeks[5]:
			data.Jisseki.Fri.SetParam(suffix, value)
		case Weeks[6]:
			data.Jisseki.Sat.SetParam(suffix, value)
		}
	case ContentTypeYotei:
		switch week {
		case Weeks[0]:
			data.Yotei.Sun.SetParam(suffix, value)
		case Weeks[1]:
			data.Yotei.Mon.SetParam(suffix, value)
		case Weeks[2]:
			data.Yotei.Tue.SetParam(suffix, value)
		case Weeks[3]:
			data.Yotei.Wed.SetParam(suffix, value)
		case Weeks[4]:
			data.Yotei.Thu.SetParam(suffix, value)
		case Weeks[5]:
			data.Yotei.Fri.SetParam(suffix, value)
		case Weeks[6]:
			data.Yotei.Sat.SetParam(suffix, value)
		}
	}
}
