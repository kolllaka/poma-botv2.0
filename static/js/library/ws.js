let socket
function connectWS(WEBSOCKET, handler) {

	socket = new WebSocket(WEBSOCKET);

	socket.onopen = function (e) {
		console.log(`ws conn open to  ${WEBSOCKET}`);
	};

	socket.onmessage = function (e) {
		msgStruct = JSON.parse(e.data)
		console.log(`[message] ${msgStruct}`);
		// обработчик
		if (!handler) {
			console.log("handler not allowed");

			return
		}

		handler()
	};

	socket.onclose = function (event) {
		if (event.wasClean) {
			console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
		} else {
			// например, сервер убил процесс или сеть недоступна
			// обычно в этом случае event.code 1006
			console.log('[close ws] Соединение прервано');
		}

		setTimeout(function () {
			connectWS();
		}, 1000);
	};

	socket.onerror = function (error) {
		console.log(`[error] ${error}`);
		socket.close();
	};
}