import React from 'react';
import { Tag } from '@carbon/react';

const StatusTag = ({ value, config = {} }) => {
  const getType = () => {
    // Check if value is undefined or null
    if (value === undefined || value === null) {
      return 'gray'; // Default color for undefined/null values
    }
    
    // If config is provided and has a mapping for this value, use it
    if (config && config[value]) {
      return config[value];
    }
    
    // Default status mapping
    const positiveStatuses = ['active', 'alive', 'stable', 'healthy'];
    const warningStatuses = ['rebalancing', 'pending'];
    
    const valueLower = String(value).toLowerCase();
    
    if (positiveStatuses.some(status => valueLower.includes(status))) {
      return 'green';
    }
    if (warningStatuses.some(status => valueLower.includes(status))) {
      return 'yellow';
    }
    return 'red';
  };

  if (value === undefined || value === null) {
    return <span>-</span>;
  }

  return <Tag type={getType()}>{value}</Tag>;
};

export default StatusTag;