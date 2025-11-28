const { default: Redis } = require("ioredis");

class RedisService {
    /**
     * Constructs a new instance of the RedisService class.
     * @param {Redis} redisClient - An instance of ioredis.
     */
    constructor(redisClient) {
        this.redisClient = redisClient;
    }

    subscribeLog(callback) {
        this.redisClient.psubscribe('logs:*');

        this.redisClient.on('pmessage', (pattern, channel, message) => {
            callback(pattern, channel, message);
        });
    }
}

module.exports = RedisService;
