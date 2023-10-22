const aug = new Augury("aug");
const WEBSOCKET = 'ws://127.0.0.1:8080/aug/ws'
const handler = () => {
	aug.change(msgStruct)
}

connectWS(WEBSOCKET, handler);
