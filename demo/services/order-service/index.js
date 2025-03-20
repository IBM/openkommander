const { faker } = require('@faker-js/faker');
const { v4: uuidv4 } = require('uuid');
const { 
  createProducer, 
  createConsumer, 
  TOPICS, 
  subscribe, 
  sendMessage,
  logger
} = require('../../utils/kafka-utils');

const SERVICE_NAME = 'order-service';

const generateProduct = () => {
  const productId = `prod-${Math.floor(Math.random() * 50)}`;
  
  return {
    id: productId,
    name: faker.commerce.productName(),
    price: parseFloat(faker.commerce.price()),
    quantity: Math.floor(Math.random() * 5) + 1
  };
};



const generateOrder = () => {
  const numProducts = Math.floor(Math.random() * 5) + 1;
  const products = Array.from({ length: numProducts }, generateProduct);
  const totalAmount = products.reduce((sum, product) => sum + (product.price * product.quantity), 0);
  
  return {
    id: uuidv4(),
    customer: {
      id: uuidv4(),
      name: faker.person.fullName(),
      email: faker.internet.email(),
      phone: faker.phone.number(),
      address: {
        street: faker.location.streetAddress(),
        city: faker.location.city(),
        state: faker.location.state(),
        zipCode: faker.location.zipCode(),
        country: faker.location.country()
      }
    },
    products,
    totalAmount,
    status: 'CREATED',
    createdAt: new Date().toISOString()
  };
};

const start = async () => {
  try {
    const producer = await createProducer(SERVICE_NAME);
    const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
    
    await subscribe(consumer, [
      TOPICS.PAYMENT_PROCESSED,
      TOPICS.INVENTORY_UPDATED,
      TOPICS.SHIPPING_PREPARED,
      TOPICS.ORDER_FAILED
    ], SERVICE_NAME);
    
    await consumer.run({
      eachMessage: async ({ topic, partition, message }) => {
        try {
          const messageData = JSON.parse(message.value.toString());
          logger.info(`Received message from ${topic}: ${message.value}`, { service: SERVICE_NAME });
          
          switch (topic) {
            case TOPICS.PAYMENT_PROCESSED:
              if (messageData.status === 'SUCCESS') {
                logger.info(`Payment successful for order ${messageData.orderId}`, { service: SERVICE_NAME });
              } else {
                await sendMessage(producer, TOPICS.ORDER_FAILED, {
                  orderId: messageData.orderId,
                  reason: 'Payment failed',
                  status: 'FAILED'
                }, SERVICE_NAME);
              }
              break;
              
            case TOPICS.INVENTORY_UPDATED:
              if (messageData.status === 'SUCCESS') {
                logger.info(`Inventory updated for order ${messageData.orderId}`, { service: SERVICE_NAME });
              } else {
                await sendMessage(producer, TOPICS.ORDER_FAILED, {
                  orderId: messageData.orderId,
                  reason: 'Inventory update failed',
                  status: 'FAILED'
                }, SERVICE_NAME);
              }
              break;
              
            case TOPICS.SHIPPING_PREPARED:
              if (messageData.status === 'SUCCESS') {
                logger.info(`Shipping prepared for order ${messageData.orderId}`, { service: SERVICE_NAME });
                await sendMessage(producer, TOPICS.ORDER_COMPLETED, {
                  orderId: messageData.orderId,
                  status: 'COMPLETED',
                  completedAt: new Date().toISOString()
                }, SERVICE_NAME);
              } else {
                await sendMessage(producer, TOPICS.ORDER_FAILED, {
                  orderId: messageData.orderId,
                  reason: 'Shipping preparation failed',
                  status: 'FAILED'
                }, SERVICE_NAME);
              }
              break;
              
            case TOPICS.ORDER_FAILED:
              logger.info(`Order ${messageData.orderId} failed: ${messageData.reason}`, { service: SERVICE_NAME });
              break;
          }
        } catch (error) {
          logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
        }
      },
    });
    
    setInterval(async () => {
      try {
        const order = generateOrder();
        logger.info(`Creating new order: ${order.id}`, { service: SERVICE_NAME });
        await sendMessage(producer, TOPICS.ORDER_CREATED, order, SERVICE_NAME);
      } catch (error) {
        logger.error(`Error generating order: ${error.message}`, { service: SERVICE_NAME });
      }
    }, 5000);
    
    logger.info('Order service started', { service: SERVICE_NAME });
  } catch (error) {
    logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
    process.exit(1);
  }
};

start();
