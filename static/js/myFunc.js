

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
    if (tell_page_move && !is_submit ) {
        event = event || window.event;
        event.returnValue = 'ページから移動しますか？';
    }
}

function saveAndStayPage() {
    savePage();
    $('#toast1').toast('show');
}


