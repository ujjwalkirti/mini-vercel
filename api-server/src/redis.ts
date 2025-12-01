import Redis from "ioredis";

export type RedisLogCallback = (pattern: string, channel: string, message: string) => void;

export default class RedisService {
    private redisClient: Redis;

    /**
     * @param redisClient  An instance of ioredis
     */
    constructor(redisClient: Redis) {
        this.redisClient = redisClient;
    }

    subscribeLog(callback: RedisLogCallback): void {
        this.redisClient.psubscribe("logs:*");

        this.redisClient.on("pmessage", (_pattern: string, channel: string, message: string) => {
            callback(_pattern, channel, message);
        });
    }
}
