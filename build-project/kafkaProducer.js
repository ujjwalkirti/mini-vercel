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
     * Generates a message to be sent to Kafka.
     * @param {string} topic - The topic to send the message to.
     * @param {object} keys - The keys to be included in the message.
     * @param {string} message - The message to be included in the message.
     * @returns {Promise<void>} A promise resolving when the message has been sent.
     */
    async generateMessage(topic, keys, message) {
        await this.producer.send({
            topic: topic,
            messages: [
                { key: 'log', value: JSON.stringify({ ...keys, log: message }) }
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
