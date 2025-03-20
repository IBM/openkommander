const { 
createProducer, 
createConsumer, 
TOPICS, 
subscribe, 
sendMessage,
logger
} = require('../../utils/kafka-utils');

const SERVICE_NAME = 'notification-service';

const sendNotification = (type, data) => {
    return new Promise(resolve => {
      setTimeout(() => {
        const success = Math.random() < 0.98;
        
        const recipientEmail = data.customer && data.customer.email ? data.customer.email : 'unknown@example.com';
        
        resolve({
          type,
          recipient: recipientEmail,
          channel: Math.random() > 0.5 ? 'EMAIL' : 'SMS',
          status: success ? 'SENT' : 'FAILED',
          sentAt: new Date().toISOString()
        });
      }, 300 + Math.random() * 700);
    });
};

const start = async () => {
try {
    const producer = await createProducer(SERVICE_NAME);
    const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
    
    await subscribe(consumer, [
    TOPICS.PAYMENT_PROCESSED,
    TOPICS.SHIPPING_PREPARED,
    TOPICS.ORDER_COMPLETED,
    TOPICS.ORDER_FAILED
    ], SERVICE_NAME);
    
    await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
        try {
        const messageData = JSON.parse(message.value.toString());
        logger.info(`Received message from ${topic}`, { service: SERVICE_NAME });
        
        switch (topic) {
            case TOPICS.PAYMENT_PROCESSED:
            if (messageData.status === 'SUCCESS') {
                logger.info(`Sending payment confirmation for order ${messageData.orderId}`, { service: SERVICE_NAME });
                const notificationResult = await sendNotification('PAYMENT_CONFIRMATION', messageData);
                
                await sendMessage(producer, TOPICS.NOTIFICATION_SENT, {
                orderId: messageData.orderId,
                notificationType: 'PAYMENT_CONFIRMATION',
                ...notificationResult
                }, SERVICE_NAME);
            }
            break;
            
            case TOPICS.SHIPPING_PREPARED:
            if (messageData.status === 'SUCCESS') {
                logger.info(`Sending shipping notification for order ${messageData.orderId}`, { service: SERVICE_NAME });
                const notificationResult = await sendNotification('SHIPPING_CONFIRMATION', messageData);
                
                await sendMessage(producer, TOPICS.NOTIFICATION_SENT, {
                orderId: messageData.orderId,
                notificationType: 'SHIPPING_CONFIRMATION',
                trackingNumber: messageData.trackingNumber,
                carrier: messageData.carrier,
                ...notificationResult
                }, SERVICE_NAME);
            }
            break;
            
            case TOPICS.ORDER_COMPLETED:
            logger.info(`Sending order completion notification for order ${messageData.orderId}`, { service: SERVICE_NAME });
            const completionResult = await sendNotification('ORDER_COMPLETED', messageData);
            
            await sendMessage(producer, TOPICS.NOTIFICATION_SENT, {
                orderId: messageData.orderId,
                notificationType: 'ORDER_COMPLETED',
                ...completionResult
            }, SERVICE_NAME);
            break;
            
            case TOPICS.ORDER_FAILED:
            logger.info(`Sending order failure notification for order ${messageData.orderId}`, { service: SERVICE_NAME });
            const failureResult = await sendNotification('ORDER_FAILED', messageData);
            
            await sendMessage(producer, TOPICS.NOTIFICATION_SENT, {
                orderId: messageData.orderId,
                notificationType: 'ORDER_FAILED',
                reason: messageData.reason,
                ...failureResult
            }, SERVICE_NAME);
            break;
        }
        } catch (error) {
        logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
        }
    },
    });
    
    logger.info('Notification service started', { service: SERVICE_NAME });
} catch (error) {
    logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
    process.exit(1);
}
};

start();
