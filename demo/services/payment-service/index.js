const { 
    createProducer, 
    createConsumer, 
    TOPICS, 
    subscribe, 
    sendMessage,
    logger
  } = require('../../utils/kafka-utils');
  
const SERVICE_NAME = 'payment-service';

const processPayment = (order) => {
const paymentProviders = ['Stripe', 'PayPal', 'Visa', 'MasterCard'];
const selectedProvider = paymentProviders[Math.floor(Math.random() * paymentProviders.length)];

return new Promise(resolve => {
    setTimeout(() => {
    const success = Math.random() < 0.9;
    
    resolve({
        provider: selectedProvider,
        amount: order.totalAmount,
        currency: 'EUR',
        status: success ? 'SUCCESS' : 'FAILED',
        transactionId: `tx-${Math.random().toString(36).substring(2, 15)}`,
        timestamp: new Date().toISOString()
    });
    }, 500 + Math.random() * 1500);
});
};

const start = async () => {
try {
    const producer = await createProducer(SERVICE_NAME);
    const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
    
    await subscribe(consumer, [TOPICS.ORDER_CREATED], SERVICE_NAME);
    
    await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
        try {
        const order = JSON.parse(message.value.toString());
        logger.info(`Received order: ${order.id}`, { service: SERVICE_NAME });
        
        logger.info(`Processing payment for order ${order.id}`, { service: SERVICE_NAME });
        const paymentResult = await processPayment(order);
        
        await sendMessage(producer, TOPICS.PAYMENT_PROCESSED, {
            orderId: order.id,
            status: paymentResult.status,
            provider: paymentResult.provider,
            transactionId: paymentResult.transactionId,
            amount: paymentResult.amount,
            products: order.products,
            customer: order.customer  
          }, SERVICE_NAME);
        
        logger.info(`Payment ${paymentResult.status} for order ${order.id}`, { service: SERVICE_NAME });
        } catch (error) {
        logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
        }
    },
    });
    
    logger.info('Payment service started', { service: SERVICE_NAME });
} catch (error) {
    logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
    process.exit(1);
}
};

start();
