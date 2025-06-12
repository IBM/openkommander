// src/pages/OverviewPage.jsx
import React, { useEffect, useState, useRef } from 'react';
import { Stack, Toggle, Tile } from '@carbon/react';
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
  const [lastRefreshed, setLastRefreshed] = useState(new Date());
  const [loading, setLoading] = useState(true);
  const [metricloading, setMetricLoading] = useState(false);
  const [error, setError] = useState(null);
  const [pollingEnabled, setPollingEnabled] = useState(true);
  const intervalRef = useRef(null);

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

        const messagesPerMinute = await api.getMessagePerMinute();

        const messagesPerMinuteTotal = messagesPerMinute.filter(msg => msg.topic === 'total')[0];
        const messagesPerMinuteDuarte = messagesPerMinute.filter(msg => msg.topic === 'duarte')[0];
        const totalMessages = messagesPerMinuteTotal['produced_count'] + messagesPerMinuteTotal['consumed_count'];
        console.log('Fetched metrics:', {
          messagesPerMinute,
          messagesPerMinuteTotal,
          messagesPerMinuteDuarte,
          totalMessages

        });

        // Create metrics object
        const metricsData = {
          brokers: brokers.length,
          topics: topics.length,
          health: health.status || 'Unknown',
          messagesPerMinute: totalMessages || 0,
        };

        setMetrics(metricsData);
        setError(null);
      } catch (err) {
        setError(`Failed to fetch metrics: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    // Initial fetch
    fetchMetrics();
  }, []);

  useEffect(() => {
    if (pollingEnabled) {
      intervalRef.current = setInterval(async () => {
        if (!metricloading) {
          setMetricLoading(true);
          try {
            const health = await api.getHealth();
            const brokers = await api.getBrokers();
            const topics = await api.getTopics();
            const messagesPerMinute = await api.getMessagePerMinute();
            const messagesPerMinuteTotal = messagesPerMinute.filter(msg => msg.topic === 'total')[0];
            const totalMessages = messagesPerMinuteTotal['produced_count'] + messagesPerMinuteTotal['consumed_count'];

            const metricsData = {
              brokers: brokers.length,
              topics: topics.length,
              health: health.status || 'Unknown',
              messagesPerMinute: totalMessages || 0,
            };

            setMetrics(metricsData);
            setLastRefreshed(new Date());
          } catch (err) {
            console.error('Error fetching metrics during polling:', err);
          } finally {
            setMetricLoading(false);
          }
        }
      }, 3000);
    }

    // Cleanup interval on component unmount
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [pollingEnabled, metricloading]);

  return (
    <>
      <ErrorNotification error={error} onClose={() => setError(null)} />
      <Stack gap={7}>
        <ClusterMetrics metrics={metrics} loading={loading} lastRefreshed={lastRefreshed} />
        <Tile>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h4 style={{ margin: 0 }}>Polling Settings</h4>
            <Toggle
              id="polling-toggle"
              labelText="Auto-refresh metrics"
              hideLabel={false}
              toggled={pollingEnabled}
              onToggle={(checked) => setPollingEnabled(checked)}
            />
          </div>
          <p style={{ margin: '8px 0 0 0', fontSize: '14px', color: '#6f6f6f' }}>
            {pollingEnabled
              ? 'Metrics will refresh automatically every 3 seconds'
              : 'Metrics will only refresh when you reload the page'}
          </p>
        </Tile>
      </Stack>
    </>
  );
};

export default OverviewPage;

