const { Server } = require("socket.io");

class WebSocketService {
    /**
     * Initializes a new instance of the WebSocketService class.
     * @param {Server} io - The Socket.IO instance to use for communication.
     */
    constructor(io) {
        this.io = io;
    }
}

module.exports = WebSocketService;
