const { Kafka } = require('kafkajs');
const { v4: uuidv4 } = require('uuid');
const winston = require('winston');

const logger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.printf(({ timestamp, level, message, service }) => {
      return `${timestamp} [${service}] ${level.toUpperCase()}: ${message}`;
    })
  ),
  transports: [
    new winston.transports.Console(),
  ],
});

const kafka = new Kafka({
  clientId: 'ecommerce-app',
  brokers: ['localhost:9093'],
  retry: {
    initialRetryTime: 300,
    retries: 10
  }
});

// Topics
const TOPICS = {
  ORDER_CREATED: 'order-created',
  ORDER_VALIDATED: 'order-validated',
  PAYMENT_PROCESSED: 'payment-processed',
  INVENTORY_UPDATED: 'inventory-updated',
  SHIPPING_PREPARED: 'shipping-prepared',
  NOTIFICATION_SENT: 'notification-sent',
  ORDER_COMPLETED: 'order-completed',
  ORDER_FAILED: 'order-failed'
};

const createProducer = async (serviceName) => {
  const producer = kafka.producer();
  await producer.connect();
  logger.info(`Producer connected`, { service: serviceName });
  return producer;
};

const createConsumer = async (groupId, serviceName) => {
  const consumer = kafka.consumer({ groupId });
  await consumer.connect();
  logger.info(`Consumer connected`, { service: serviceName });
  return consumer;
};

const subscribe = async (consumer, topics, serviceName) => {
  for (const topic of topics) {
    await consumer.subscribe({ topic, fromBeginning: false });
    logger.info(`Subscribed to ${topic}`, { service: serviceName });
  }
};

const sendMessage = async (producer, topic, message, serviceName) => {
  try {
    const messageId = uuidv4();
    await producer.send({
      topic,
      messages: [
        { 
          key: messageId, 
          value: JSON.stringify({
            ...message,
            id: messageId,
            timestamp: new Date().toISOString()
          })
        },
      ],
    });
    logger.info(`Message sent to ${topic}`, { service: serviceName });
    return messageId;
  } catch (error) {
    logger.error(`Error sending message to ${topic}: ${error.message}`, { service: serviceName });
    throw error;
  }
};

module.exports = {
  kafka,
  TOPICS,
  createProducer,
  createConsumer,
  subscribe,
  sendMessage,
  logger
};

