// src/pages/BrokersPage.jsx
import React, { useState } from 'react';
import ResourcePage from '../components/common/ResourcePage';
import StatusTag from '../components/common/StatusTag';
import { TABLE_CONFIGS } from '../config/constants';
import api from '../services/api';

const BrokersPage = () => {
  const [brokerConfig, setBrokerConfig] = useState(null);

  const handleViewConfig = async (broker) => {
    try {
      // Use API to get broker config (not in OpenAPI but we'll assume it exists)
      const config = await api.getBrokerDetails(broker.id);
      setBrokerConfig(config);
      return config;
    } catch (err) {
      console.error('Failed to fetch broker configuration');
      return null;
    }
  };

  // Transform broker data from API to UI format
  const transformRow = (apiBroker) => {
    return {
      ...apiBroker,
      // Calculate status based on in_sync_partitions vs partitions
      status: apiBroker.in_sync_partitions === apiBroker.partitions ? 'alive' : 'dead',
      // Calculate partition count for UI
      partitionCount: apiBroker.partitions_leader || apiBroker.partitions || 0
    };
  };

  const renderCustomCell = (cell, row, index) => {
    if (cell.info.header === 'status') {
      return <StatusTag 
        value={cell.value}
        config={{
          'alive': 'green',
          'dead': 'red'
        }}
      />;
    }
    return cell.value;
  };

  const customActions = [{
    key: 'onViewConfig',
    handler: handleViewConfig
  }];

  return (
    <ResourcePage
      endpoint="brokers"
      tableConfig={TABLE_CONFIGS.brokers}
      renderCustomCell={renderCustomCell}
      customActions={customActions}
      transformRow={transformRow}
    />
  );
};

export default BrokersPage;

