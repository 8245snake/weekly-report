package main

//ContentType 予定or実績
type ContentType string

const (
	//ContentTypeJisseki 実績
	ContentTypeJisseki ContentType = "jisseki"
	//ContentTypeYotei 予定
	ContentTypeYotei ContentType = "yotei"
)

//ReportData 週報データ全て
type ReportData struct {
	Jisseki  WeeklyContents `json:"jisseki"`
	Yotei    WeeklyContents `json:"yotei"`
	Title    string         `json:"title"`
	Tasks    string         `json:"task"`
	Schedule string         `json:"schedule"`
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
	Type        ContentType `json:"type"`
	IsHolyday   bool        `json:"is_holyday"`
	Prefix      string      `json:"prefix,omitempty"`
	Style       string      `json:"style,omitempty"`
	DateValue   string      `json:"date"`
	DispDate    string      `json:"disp_date,omitempty"`
	Placeholder string      `json:"placeholder,omitempty"`
	Text        string      `json:"text"`
	SubText     string      `json:"subtext"`
}

//Weeks 週の名前
var Weeks = [7]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

//WeeksJP 週の名前
var WeeksJP = [7]string{"(日)", "(月)", "(火)", "(水)", "(木)", "(金)", "(土)"}

//NewReportData 初期化
func NewReportData() ReportData {
	return ReportData{
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
