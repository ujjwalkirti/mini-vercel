const { Kafka } = require('kafkajs')

class KafkaConsumerService {
    /**
     * Constructs a new KafkaConsumerService instance.
     * @param {Kafka} kafka - An instance of the Kafka class from the kafkajs library.
     * @param {string} groupId - The group ID to subscribe to the topic with.
     */
    constructor(kafka, groupId) {
        this.kafka = kafka;
        this.consumer = this.kafka.consumer({ groupId: groupId });
    }

    async connect() {
        await this.consumer.connect()
    }

    async listenForMessagesInBatch(topic, callback) {
        await this.connect();
        await this.consumer.subscribe({ topic, fromBeginning: true });
        await this.consumer.run({
            eachBatch: async ({ batch, heartbeat, commitOffsetsIfNecessary, resolveOffset }) => {
                for (const message of batch.messages) {
                    await callback(message);
                    resolveOffset(message.offset)
                    await commitOffsetsIfNecessary(message.offset)
                    await heartbeat()
                }
            }
        });
    }

    async disconnect() {
        await this.consumer.disconnect()
    }
}

module.exports = KafkaConsumerService
