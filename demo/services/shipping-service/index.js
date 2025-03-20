const { faker } = require('@faker-js/faker');
const { 
  createProducer, 
  createConsumer, 
  TOPICS, 
  subscribe, 
  sendMessage,
  logger
} = require('../../utils/kafka-utils');

const SERVICE_NAME = 'shipping-service';

const prepareShipment = (orderData) => {
  const carriers = ['AnPost', 'UPS', 'DPD', 'Fastway'];
  const selectedCarrier = carriers[Math.floor(Math.random() * carriers.length)];
  
  return new Promise(resolve => {
    setTimeout(() => {
      const success = Math.random() < 0.95;
      
      if (success) {
        resolve({
          status: 'SUCCESS',
          trackingNumber: faker.string.alphanumeric(12).toUpperCase(),
          carrier: selectedCarrier,
          estimatedDelivery: new Date(Date.now() + (Math.floor(Math.random() * 5) + 2) * 24 * 60 * 60 * 1000).toISOString(),
          address: orderData.customer.address
        });
      } else {
        resolve({
          status: 'FAILED',
          reason: 'Unable to prepare shipment'
        });
      }
    }, 700 + Math.random() * 1200);
  });
};

const start = async () => {
  try {
    const producer = await createProducer(SERVICE_NAME);
    const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
    
    await subscribe(consumer, [TOPICS.INVENTORY_UPDATED], SERVICE_NAME);
    
    await consumer.run({
      eachMessage: async ({ topic, partition, message }) => {
        try {
          const messageData = JSON.parse(message.value.toString());
          logger.info(`Received message from ${topic}`, { service: SERVICE_NAME });
          
          if (messageData.status === 'SUCCESS') {
            logger.info(`Preparing shipment for order ${messageData.orderId}`, { service: SERVICE_NAME });
            
            // Prepare shipment
            
            await sendMessage(producer, TOPICS.SHIPPING_PREPARED, {
                orderId: messageData.orderId,
                status: shipmentResult.status,
                ...shipmentResult,
                customer: messageData.customer  
              }, SERVICE_NAME);
            
            logger.info(`Shipment ${shipmentResult.status} for order ${messageData.orderId}`, { service: SERVICE_NAME });
          }
        } catch (error) {
          logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
        }
      },
    });
    
    logger.info('Shipping service started', { service: SERVICE_NAME });
  } catch (error) {
    logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
    process.exit(1);
  }
};

start();
