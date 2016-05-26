'use strict';

var conn = new rattle.newConnection("ws://127.0.0.1:8080/ws", true); //addr, debug

console.log(conn)
conn.event("onConnect", function (evt) {
	document.getElementById("toggle").checked = true;
	conn.send("Main.Timer");
	conn.send("Main.Index")
})

//for this case, rattle correct fill field `From` 
var test = {
	Send: function () {
		var data = {};
		data.text = document.getElementById("text").value;

		var url = document.getElementById("json").checked ? "Main.JSON" : "Main.RAW"

		conn.send(url, data);
	},

	RecieveJSON: function (data) {
		document.getElementById("msgs").innerHTML = JSON.stringify(data);
	},

	RecieveRAW: function (data) {
		document.getElementById("msgs").innerHTML = data;
	}
}

function toggle() {
	console.log(conn)

	if (conn.connected) {
		conn.disconnect()
	} else {
		conn.connect()
	}
}