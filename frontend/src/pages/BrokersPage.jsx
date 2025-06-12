// src/pages/BrokersPage.jsx
import React, { useState } from 'react';
import ResourcePage from '../components/common/ResourcePage';
import StatusTag from '../components/common/StatusTag';
import { TABLE_CONFIGS, FORM_CONFIGS } from '../config/constants';
import api from '../services/api';
import { TextInput } from '@carbon/react';

const BrokersPage = () => {
  const BROKER_FORM_CONFIG = {
    initialData: {
      host: '',
      port: 9092,
    }
  };

  const [brokerConfig, setBrokerConfig] = useState(null);

  const selectedBroker = localStorage.getItem('selectedBroker') || '';
  console.log('Selected broker from local storage:', selectedBroker);

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
    console.log('Transforming broker data:', apiBroker);
    // const [host, port] = apiBroker.addr.split(':');
    const address = apiBroker.host + ':' + (apiBroker.port || 9092) || apiBroker.addr || '';

    return {
      ...apiBroker,
      host: apiBroker.host || apiBroker.addr.split(':')[0],
      port: apiBroker.port || parseInt(apiBroker.addr.split(':')[1], 10),
      // Calculate status based on in_sync_partitions vs partitions
      status: address == selectedBroker ? 'selected' : 'inactive',
      // Calculate partition count for UI
      partitionCount: apiBroker.partitions_leader || apiBroker.partitions || 0
    };
  };

  const renderCustomCell = (cell, row, index) => {
    if (cell.info.header === 'status') {
      return <StatusTag 
        value={cell.value}
        config={{
          'selected': 'green',
          'inactive': 'grey'
        }}
      />;
    }
    return cell.value;
  };

  const renderFormContent = ({ formData, setFormData, isEditing }) => (
    <>
      <TextInput
        id="host"
        labelText="Host"
        value={formData.host}
        onChange={(e) => setFormData({ ...formData, host: e.target.value })}
        disabled={isEditing}
      />
      <TextInput
        id="port"
        labelText="Port"
        value={formData.port}
        onChange={(e) => setFormData({ ...formData, port: parseInt(e.target.value) || 0})}
        disabled={isEditing}
      />
    </>
  );

  const customActions = [{
    key: 'onViewConfig',
    handler: handleViewConfig
  }];

  return (
    <ResourcePage
      endpoint="brokers"
      tableConfig={TABLE_CONFIGS.brokers}
      formConfig={BROKER_FORM_CONFIG}
      renderFormContent={renderFormContent}
      renderCustomCell={renderCustomCell}
      customActions={customActions}
      transformRow={transformRow}
    />
  );
};

export default BrokersPage;

