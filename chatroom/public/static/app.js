window.onload = function () {
    let conn;
    let msg = document.getElementById("msg");
    let log = document.getElementById("log");

    function appendLog(item) {
        let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        const params = window.location.href.split("?");
        const room = params[params.length-1];
        const nick = params[params.length-2];
        conn = new WebSocket("ws://" + document.location.host + "/ws?"+nick+"&"+room);
        conn.onclose = function (evt) {
            let item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onmessage = function (evt) {
            let messages = evt.data.split('\n');
            for (let i = 0; i < messages.length; i++) {
                let item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        let item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};