const { 
createConsumer, 
TOPICS, 
subscribe,
logger
} = require('../../utils/kafka-utils');

const SERVICE_NAME = 'analytics-service';

// Analytics metrics
const metrics = {
totalOrders: 0,
successfulOrders: 0,
failedOrders: 0,
totalRevenue: 0,
averageOrderValue: 0,
paymentSuccessRate: 0,
productsByPopularity: new Map(),
notificationStats: {
    sent: 0,
    failed: 0
}
};

const updateMetrics = (topic, data) => {
switch (topic) {
    case TOPICS.ORDER_CREATED:
    metrics.totalOrders++;
    break;
    
    case TOPICS.ORDER_COMPLETED:
    metrics.successfulOrders++;
    break;
    
    case TOPICS.ORDER_FAILED:
    metrics.failedOrders++;
    break;
    
    case TOPICS.PAYMENT_PROCESSED:
    if (data.status === 'SUCCESS') {
        metrics.totalRevenue += data.amount;
        
        // Update average order value
        if (metrics.successfulOrders > 0) {
        metrics.averageOrderValue = metrics.totalRevenue / metrics.successfulOrders;
        }
        
        // Update payment success rate
        const totalPayments = metrics.successfulOrders + metrics.failedOrders;
        if (totalPayments > 0) {
        metrics.paymentSuccessRate = (metrics.successfulOrders / totalPayments) * 100;
        }
        
        // Update product popularity
        if (data.products) {
        for (const product of data.products) {
            const count = metrics.productsByPopularity.get(product.id) || 0;
            metrics.productsByPopularity.set(product.id, count + 1);
        }
        }
    }
    break;
    
    case TOPICS.NOTIFICATION_SENT:
    if (data.status === 'SENT') {
        metrics.notificationStats.sent++;
    } else {
        metrics.notificationStats.failed++;
    }
    break;
}
};

const start = async () => {
try {
    // Create consumer
    const consumer = await createConsumer(`${SERVICE_NAME}-group`, SERVICE_NAME);
    
    // Subscribe to all topics for analytics
    const allTopics = Object.values(TOPICS);
    await subscribe(consumer, allTopics, SERVICE_NAME);
    
    // Process messages
    await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
        try {
        const data = JSON.parse(message.value.toString());
        
        // Update metrics
        updateMetrics(topic, data);
        
        // Print metrics every 20 messages
        if ((metrics.totalOrders + metrics.notificationStats.sent) % 20 === 0) {
            logMetrics();
        }
        } catch (error) {
        logger.error(`Error processing message: ${error.message}`, { service: SERVICE_NAME });
        }
    },
    });
    
    // Log metrics every 30 seconds
    setInterval(logMetrics, 30000);
    
    logger.info('Analytics service started', { service: SERVICE_NAME });
} catch (error) {
    logger.error(`Startup error: ${error.message}`, { service: SERVICE_NAME });
    process.exit(1);
}
};

const logMetrics = () => {
logger.info('Current Analytics Metrics:', { service: SERVICE_NAME });
logger.info(`Total Orders: ${metrics.totalOrders}`, { service: SERVICE_NAME });
logger.info(`Successful Orders: ${metrics.successfulOrders}`, { service: SERVICE_NAME });
logger.info(`Failed Orders: ${metrics.failedOrders}`, { service: SERVICE_NAME });
logger.info(`Total Revenue: $${metrics.totalRevenue.toFixed(2)}`, { service: SERVICE_NAME });
logger.info(`Average Order Value: $${metrics.averageOrderValue.toFixed(2)}`, { service: SERVICE_NAME });
logger.info(`Payment Success Rate: ${metrics.paymentSuccessRate.toFixed(2)}%`, { service: SERVICE_NAME });
logger.info(`Notifications Sent: ${metrics.notificationStats.sent}`, { service: SERVICE_NAME });
logger.info(`Notifications Failed: ${metrics.notificationStats.failed}`, { service: SERVICE_NAME });

// Log top 5 most popular products
const sortedProducts = [...metrics.productsByPopularity.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 5);
    
if (sortedProducts.length > 0) {
    logger.info('Top 5 Most Popular Products:', { service: SERVICE_NAME });
    sortedProducts.forEach(([productId, count], index) => {
    logger.info(`${index + 1}. Product ${productId}: ${count} orders`, { service: SERVICE_NAME });
    });
}
};

start();

