

class RedisService {
    constructor(redisClient) {
        this.redisClient = redisClient;
    }

    publishLog(log) {
        this.redisClient.publish('log', JSON.stringify({ log }));
    }


}

export default RedisService;

