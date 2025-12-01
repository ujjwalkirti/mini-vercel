import { Kafka } from "kafkajs";

class KafkaProducerService {
    /**
     * Constructs a new KafkaService instance.
     * @param {Kafka} kafka - An instance of the Kafka class from the kafkajs library.
     */
    constructor(kafka) {
        this.kafka = kafka;
        this.producer = this.kafka.producer()
    }

    async connect() {
        await this.producer.connect()
    }

    /**
     * Sends a single message to a Kafka topic.
     * @param {string} topic - The name of the Kafka topic to send the message to.
     * @param {string} key - The key of the message to send.
     * @param {string} message - The value of the message to send.
     * @returns {Promise<void>} A promise that resolves when the message has been sent.
     */
    async generateMessage(topic, key, message) {
        await this.producer.send({
            topic: topic,
            messages: [
                { key: key, value: JSON.stringify(message) }
            ]
        })
    }

    async generateContinuousMessages(topic, message, interval) {
        await this.producer.connect()
        setInterval(() => {
            this.producer.send({
                topic: topic,
                messages: [
                    { value: message }
                ]
            })
        }, interval);
    }

    async disconnect() {
        await this.producer.disconnect()
    }
}

export default KafkaProducerService
