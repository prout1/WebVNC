var mouseState = {
    type: "pe",
    buttons: {
        left: false,
        right: false,
        mid: false,
        up: false,
        down: false
    },
    x: 0,
    y: 0,
}
var buttonMap = {
    "0": "left",
    "1": "mid",
    "2" : "right"
}
var scrollUp = "up",
    scrollDown = "down"

var host = "192.168.43.25";
var port = ":8080";
function sendMouseUpdate(ws) {
    ws.send(JSON.stringify(mouseState));
}
function handleMouse(canvas, ws) {
    var mouseEvent = function (e, pressFlag) {
        mouseState.buttons[buttonMap[e.button.toString()]] = pressFlag;
        mouseState.x = Math.round((e.x / canvas.width ) * 65536);
        mouseState.y = Math.round((e.y / canvas.height) * 65536);
        sendMouseUpdate(ws);
        console.log(mouseState.x + " " + mouseState.y);
    }
    canvas.onmousedown = function (e) {
        e.preventDefault();
        mouseEvent(e, true);
    }

    document.oncontextmenu = function (e) {
        e.preventDefault();
        e.stopPropagation();
    }
    canvas.addEventListener("wheel", function (e) {
        mouseState.buttons[scrollUp] = (e.deltaY < 0);
        mouseState.buttons[scrollDown] = !mouseState.buttons[scrollUp];
        sendMouseUpdate(ws);
        // restoring the scroll states to false, because we dont want to scroll continuously
        mouseState.buttons[scrollUp] = mouseState.buttons[scrollDown] = false;
    });

    canvas.onmouseup = function(e) {
        mouseEvent(e, false);
    }
}

function handleKeyboard(canvas, ws) {
    var keyBdEvent = function (e, pressFlag) {
        ws.send(JSON.stringify({
            type: "ke",
            keyCode: e.keyCode,
            press: pressFlag
        }));
    }

    document.addEventListener("keyup", function (e) {
        keyBdEvent(e, false);
    });

    document.addEventListener("keydown", function (e) {
        keyBdEvent(e, true);
    });
}

function handleUpdates(canvas, ws) {
    setInterval(function () {
        ws.send(JSON.stringify({ type: "re" })); // frame update request
        
    }, 80);
}

window.onload = function () {
    document.body = document.createElement("body");
    var canvas = document.createElement("canvas");
    var width = (1366*11)/17,
        height = (768*11)/17
    canvas.setAttribute("width",width.toString());
    canvas.setAttribute("height", height.toString());
    var ctx = canvas.getContext("2d");

    ctx.fillRect(0, 0, width, height);
    var ws = new WebSocket("ws://" + host + port);

    handleKeyboard(canvas,ws);
    handleMouse(canvas,ws);
    handleUpdates(canvas, ws);

    ws.onmessage = function (e) {
        var ctx = canvas.getContext("2d");
        var image = new Image();
        image.src = 'data:image/  jpeg;base64,' + e.data;
        image.onload = function () {
            ctx.drawImage(image, 0, 0, width, height);
        };
    }
    document.body.appendChild(canvas);
}