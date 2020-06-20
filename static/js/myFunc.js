window.onload = function (event) {
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("static/wasm/main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
}

//保存ボタンがあるフォームから移動しようとしたら警告する
//ただしsubmitの場合は除く
window.onbeforeunload = function (event) {
    var tell_page_move = (this.document.getElementsByName("save").length > 0);
    var is_submit = (event.srcElement.activeElement.type == "submit")
    if (tell_page_move && !is_submit) {
        event = event || window.event;
        event.returnValue = 'ページから移動しますか？';
    }
}

//保存して週報作成
function saveAndStayPage() {
    savePage();
    $('#toast1').toast('show');
}

// Date型を文字列（yyyy-mm-dd）に変換
function getStringFromDate(date) {
    var year_str = date.getFullYear();
    var month_str = date.getMonth();
    var day_str = date.getDate();

    month_str = ('0' + month_str).slice(-2);
    day_str = ('0' + day_str).slice(-2);

    format_str = 'YYYY-MM-DD';
    format_str = format_str.replace(/YYYY/g, year_str);
    format_str = format_str.replace(/MM/g, month_str);
    format_str = format_str.replace(/DD/g, day_str);

    return format_str;
};

// 先週のデータをコピー
function copyFromLastWeek(id) {
    var start = this.document.getElementsByName("start")[0].value;
    var dateArr = start.split('-');
    var date = new Date(dateArr[0], dateArr[1], dateArr[2]);

    date.setDate(date.getDate() - 7);
    qry = getStringFromDate(date);

    var mode = id.split('-')[1];

    fetch('/api/report?start=' + qry)
        .then(response => response.json())
        .then(data => {
            var text = '';
            switch (mode) {
                case 'tasks':
                    text = data.task;
                    break;
                case 'schedule':
                    text = data.schedule;
                    break;
            }
            this.document.getElementById(mode).value = text;
        });
}
