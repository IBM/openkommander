// src/pages/TopicsPage.jsx
import React from 'react';
import { TextInput, NumberInput } from '@carbon/react';
import ResourcePage from '../components/common/ResourcePage';
import StatusTag from '../components/common/StatusTag';
import { TABLE_CONFIGS } from '../config/constants';

// Updated form config to match API
const TOPIC_FORM_CONFIG = {
  initialData: {
    name: '',
    partitions: 1,
    replication_factor: 1, // snake_case to match API
    // API doesn't support retentionMs directly
  }
};

// Transform API data to match UI expectations
const transformRow = (apiTopic) => {
  return {
    ...apiTopic,
    // Ensure UI properties exist with right names
    replicationFactor: apiTopic.replicationFactor || apiTopic.replication_factor,
    status: determineTopicStatus(apiTopic),
  };
};

// Helper function to determine topic status from API data
const determineTopicStatus = (topic) => {
  if (topic.replicas === topic.in_sync_replicas && topic.replicas > 0) {
    return 'active';
  } else if (topic.replicas === 0) {
    return 'inactive';
  } else {
    return 'warning';
  }
};

const TopicsPage = () => {
  const renderFormContent = ({ formData, setFormData, isEditing }) => (
    <>
      <TextInput
        id="name"
        labelText="Topic Name"
        value={formData.name}
        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        disabled={isEditing}
      />
      <NumberInput
        id="partitions"
        label="Partitions"
        value={formData.partitions}
        onChange={(e) => setFormData({ ...formData, partitions: parseInt(e.target.value) })}
        min={1}
        max={100}
      />
      <NumberInput
        id="replication_factor"
        label="Replication Factor"
        value={formData.replication_factor}
        onChange={(e) => setFormData({ 
          ...formData, 
          replication_factor: parseInt(e.target.value),
          replicationFactor: parseInt(e.target.value) // Keep both properties for compatibility
        })}
        min={1}
        max={3}
      />
    </>
  );

  const renderCustomCell = (cell, row, index) => {
    if (cell.info.header === 'actions') {
      return (
        <div className="flex gap-2">
          {TABLE_CONFIGS.topics.actions.map((action, actionIndex) => (
            <Button
              key={actionIndex}
              kind="ghost"
              size="sm"
              renderIcon={action.icon}
              iconDescription={action.description}
              hasIconOnly
              onClick={() => action.onClick(row)} // Ensure action triggers correctly
            />
          ))}
        </div>
      );
    }
    if (cell.info.header === 'status') {
      return <StatusTag 
        value={cell.value} 
        config={{
          'active': 'green',
          'inactive': 'red',
          'warning': 'yellow'
        }} 
      />;
    }
    return cell.value;
  };

  return (
    <ResourcePage
      endpoint="topics"
      tableConfig={TABLE_CONFIGS.topics}
      formConfig={TOPIC_FORM_CONFIG}
      renderFormContent={renderFormContent}
      renderCustomCell={renderCustomCell}
      transformRow={transformRow}
    />
  );
};

export default TopicsPage;
