package main

import (
	"bytes"
	"fmt"
	"text/template"
)

const tamplateReport = `{{.Title}}

■先週の実績
{{.Jisseki.Mon.DispDate}} {{.Jisseki.Mon.SubText}}
{{.Jisseki.Mon.Text}}

{{.Jisseki.Tue.DispDate}} {{.Jisseki.Tue.SubText}}
{{.Jisseki.Tue.Text}}

{{.Jisseki.Wed.DispDate}} {{.Jisseki.Wed.SubText}}
{{.Jisseki.Wed.Text}}

{{.Jisseki.Thu.DispDate}} {{.Jisseki.Thu.SubText}}
{{.Jisseki.Thu.Text}}

{{.Jisseki.Fri.DispDate}} {{.Jisseki.Fri.SubText}}
{{.Jisseki.Fri.Text}}

{{.Jisseki.Sat.DispDate}} {{.Jisseki.Sat.SubText}}
{{.Jisseki.Sat.Text}}

{{.Jisseki.Sun.DispDate}} {{.Jisseki.Sun.SubText}}
{{.Jisseki.Sun.Text}}

■今週の予定
{{.Yotei.Mon.DispDate}} {{.Yotei.Mon.SubText}}
{{.Yotei.Mon.Text}}

{{.Yotei.Tue.DispDate}} {{.Yotei.Tue.SubText}}
{{.Yotei.Tue.Text}}

{{.Yotei.Wed.DispDate}} {{.Yotei.Wed.SubText}}
{{.Yotei.Wed.Text}}

{{.Yotei.Thu.DispDate}} {{.Yotei.Thu.SubText}}
{{.Yotei.Thu.Text}}

{{.Yotei.Fri.DispDate}} {{.Yotei.Fri.SubText}}
{{.Yotei.Fri.Text}}

{{.Yotei.Sat.DispDate}} {{.Yotei.Sat.SubText}}
{{.Yotei.Sat.Text}}

{{.Yotei.Sun.DispDate}} {{.Yotei.Sun.SubText}}
{{.Yotei.Sun.Text}}

■適宜やること
{{.Tasks}}

■近日中の予定
{{.Schedule}}
`

//ApplyTemplate テンプレートを適用して文字列を返す
func ApplyTemplate(tpl string, data interface{}) string {
	// template.New(<テンプレート名>).Parse(<文字列>)
	t, err := template.New("sample").Parse(tpl)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Execute(io.Writer(出力先), <データ>)
	var st string
	bt := bytes.NewBufferString(st)
	if err = t.Execute(bt, data); err != nil {
		fmt.Printf("%v\n", err)
	}

	return bt.String()
}
