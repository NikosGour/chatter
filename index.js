const util = require('util');
// Connect to Go WebSocket server
const ws = new WebSocket('ws://localhost:8080/');


// Event: connection opened
ws.onopen = function () {

	console.log('Connected to server');
};

// Event: message received
ws.onmessage = function (data) {
	console.log("data: ", Object.values(data));
	console.log('Received from server:', util.inspect(data, { depth: null, colors: true, showHidden: true }));
};

// Event: connection closed
ws.onclose = () => {
	console.log('Connection closed');
};
