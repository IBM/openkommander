// src/pages/OverviewPage.jsx
import React, { useEffect, useState } from 'react';
import { Stack } from '@carbon/react';
import ClusterMetrics from '../components/overview/ClusterMetrics';
import ErrorNotification from '../components/common/ErrorNotification';
import api from '../services/api';

const OverviewPage = () => {
  const [metrics, setMetrics] = useState({
    brokers: 0,
    topics: 0,
    health: 'Unknown',
    messagesPerSecond: 0
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchMetrics = async () => {
      setLoading(true);
      try {
        // Get health status
        const health = await api.getHealth();
        
        // Get brokers count
        const brokers = await api.getBrokers();
        
        // Get topics count
        const topics = await api.getTopics();
        
        // Create metrics object
        const metricsData = {
          brokers: brokers.length,
          topics: topics.length,
          health: health.status || 'Unknown',
          // This is mock data since messages/sec isn't in the API
          messagesPerSecond: Math.floor(Math.random() * 1000)
        };
        
        setMetrics(metricsData);
        setError(null);
      } catch (err) {
        setError(`Failed to fetch metrics: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchMetrics();
  }, []);

  return (
    <>
      <ErrorNotification error={error} onClose={() => setError(null)} />
      <Stack gap={7}>
        <ClusterMetrics metrics={metrics} loading={loading} />
      </Stack>
    </>
  );
};

export default OverviewPage;

