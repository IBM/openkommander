const { 
    createProducer, 
    createConsumer, 
    TOPICS, 
    subscribe, 
    sendMessage,
    logger
  } = require('../../utils/kafka-utils');
  
  const SERVICE_NAME = 'inventory-service';
  
  const inventory = new Map();
  
  // initialize with some random products
  for (let i = 0; i < 50; i++) {
    const productId = `prod-${i}`;
    inventory.set(productId, {
      id: productId,
      stock: Math.floor(Math.random() * 100) + 10
    });
  }
  
  const checkInventory = (products) => {
    const unavailableProducts = [];
  
    for (const product of products) {
      const inventoryItem = inventory.get(product.id);
      
      if (!inventoryItem || inventoryItem.stock < product.quantity) {
        unavailableProducts.push(product.id);
      }
    }
  
    return {
      success: unavailableProducts.length === 0,
      unavailableProducts
    };
  };
  
  const updateInventory = (products) => {
    for (const product of products) {
      const inventoryItem = inventory.get(product.id);
      if (inventoryItem) {
        inventoryItem.stock -= product.quantity;
        inventory.set(product.id, inventoryItem);
      } else {
        inventory.set(product.id, {
          id: product.id,
          stock: Math.floor(Math.random() * 100) - product.quantity
        });
      }
    }
  };
  
  const start = async () => {
    try {
      const producer = await createProducer(SERVICE_NAME);
      const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
      
      await subscribe(consumer, [
        TOPICS.ORDER_VALIDATED,
        TOPICS.PAYMENT_PROCESSED
      ], SERVICE_NAME);
      
      await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
          try {
            const messageData = JSON.parse(message.value.toString());
            logger.info(`Received message from ${topic}`, { service: SERVICE_NAME });
            
            switch (topic) {
              case TOPICS.PAYMENT_PROCESSED:
                if (messageData.status === 'SUCCESS') {
                  logger.info(`Processing inventory for order ${messageData.orderId}`, { service: SERVICE_NAME });
                  
                  // Simulate processing time
                  await new Promise(resolve => setTimeout(resolve, 500 + Math.random() * 1000));
                  
                  // Update inventory (80% success rate)
                  const success = Math.random() < 0.8;
                  if (success) {
                    updateInventory(messageData.products);
                    await sendMessage(producer, TOPICS.INVENTORY_UPDATED, {
                      orderId: messageData.orderId,
                      status: 'SUCCESS',
                      products: messageData.products,
                      customer: messageData.customer  
                    }, SERVICE_NAME);
                  } else {
                    await sendMessage(producer, TOPICS.INVENTORY_UPDATED, {
                      orderId: messageData.orderId,
                      status: 'FAILED',
                      reason: 'Inventory update failed',
                      customer: messageData.customer  
                    }, SERVICE_NAME);
                  }
                }
                break;
            }
          } catch (error) {
            logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
          }
        },
      });
      
      logger.info('Inventory service started', { service: SERVICE_NAME });
    } catch (error) {
      logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
      process.exit(1);
    }
  };
  
  start();