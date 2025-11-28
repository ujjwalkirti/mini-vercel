import Redis from "ioredis";


class RedisService {
    /**
     * Creates a new RedisService instance.
     * @param {Redis} redisClient - An instance of ioredis.
     */
    constructor(redisClient) {
        this.redisClient = redisClient;
    }

    async publishLog(channel, log) {
        this.redisClient.publish(channel, log);
    }


}

export default RedisService;

