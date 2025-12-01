import type { Server } from "socket.io";

export default class WebSocketService {
    private io: Server;

    /**
     * @param io  An instance of socket.io Server
     */
    constructor(io: Server) {
        this.io = io;
    }
}
