import { Kafka } from "kafkajs";
import type { Consumer, EachBatchPayload, KafkaMessage, Offsets } from "kafkajs";

export type KafkaBatchMessageCallback = (message: KafkaMessage) => Promise<void> | void;

export default class KafkaConsumerService {
    private kafka: Kafka;
    private consumer: Consumer;

    /**
     * @param kafka  An instance of Kafka from kafkajs
     * @param groupId  Consumer group ID
     */
    constructor(kafka: Kafka, groupId: string) {
        this.kafka = kafka;
        this.consumer = this.kafka.consumer({ groupId });
    }

    async connect(): Promise<void> {
        await this.consumer.connect();
    }

    async listenForMessagesInBatch(
        topic: string,
        callback: KafkaBatchMessageCallback
    ): Promise<void> {
        await this.connect();
        await this.consumer.subscribe({ topic, fromBeginning: true });

        await this.consumer.run({
            eachBatch: async ({ batch, heartbeat, commitOffsetsIfNecessary, resolveOffset }: EachBatchPayload) => {
                for (const message of batch.messages) {
                    await callback(message);

                    resolveOffset(message.offset);
                    await commitOffsetsIfNecessary(message.offset as unknown as Offsets);
                    await heartbeat();
                }
            }
        });
    }

    async disconnect(): Promise<void> {
        await this.consumer.disconnect();
    }
}
