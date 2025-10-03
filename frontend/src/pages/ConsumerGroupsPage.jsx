// src/pages/ConsumerGroupsPage.jsx
import React from 'react';
import ResourcePage from '../components/common/ResourcePage';
import StatusTag from '../components/common/StatusTag';
import { TABLE_CONFIGS } from '../config/constants';
import api from '../services/api';

const ConsumerGroupsPage = () => {
  // Hide add button since consumer groups are created by clients
  const tableConfigWithoutAdd = {
    ...TABLE_CONFIGS.consumerGroups,
    hideAdd: true
  };

  const handleViewAssignments = async (group) => {
    try {
      // Use the group_id or groupId, depending on what's available
      const groupId = group.groupId || group.group_id;
      const data = await api.getConsumerGroupAssignments(groupId);
      return data.assignments;
    } catch (err) {
      console.error('Failed to fetch consumer group assignments');
      return null;
    }
  };

  const handleResetOffsets = async (groupId) => {
    if (window.confirm('Are you sure you want to reset offsets for this consumer group?')) {
      try {
        await api.resetConsumerOffsets(groupId);
        return true;
      } catch (err) {
        console.error('Failed to reset consumer group offsets');
        return false;
      }
    }
  };

  // Transform API data to match UI expectations
  const transformRow = (apiConsumerGroup) => {
    return {
      ...apiConsumerGroup,
      // Ensure UI properties exist with right names
      groupId: apiConsumerGroup.groupId || apiConsumerGroup.group_id,
      status: apiConsumerGroup.state || 'unknown',
      // Convert topics from count to array if necessary
      topics: Array.isArray(apiConsumerGroup.topics) 
        ? apiConsumerGroup.topics 
        : Array(apiConsumerGroup.topics || 0).fill('').map((_, i) => `Topic ${i+1}`),
      // Preserve lag
      lag: apiConsumerGroup.lag || 0,
    };
  };

  const renderCustomCell = (cell, row, index) => {
    if (cell.info.header === 'status') {
      // Map API statuses to UI status types
      const statusMap = {
        'Stable': 'stable',
        'Rebalancing': 'rebalancing',
        'Dead': 'dead',
        // Add any other status mappings
      };
      
      const statusValue = statusMap[cell.value] || cell.value.toLowerCase();
      
      return <StatusTag 
        value={statusValue}
        config={{
          'stable': 'green',
          'rebalancing': 'yellow',
          'dead': 'red'
        }}
      />;
    }
    
    if (cell.info.header === 'topics') {
      // If topics is an array, join with commas
      if (Array.isArray(cell.value)) {
        return cell.value.join(', ');
      }
      // If it's a number (from API), show as count
      if (typeof cell.value === 'number') {
        return `${cell.value} topic${cell.value !== 1 ? 's' : ''}`;
      }
      return cell.value;
    }
    
    return cell.value;
  };

  const customActions = [
    {
      key: 'onViewAssignments',
      handler: handleViewAssignments
    },
    {
      key: 'onResetOffsets',
      handler: handleResetOffsets
    }
  ];

  return (
    <ResourcePage
      endpoint="consumers"
      tableConfig={tableConfigWithoutAdd}
      renderCustomCell={renderCustomCell}
      customActions={customActions}
      transformRow={transformRow}
      // No formConfig since consumer groups can't be created manually
    />
  );
};

export default ConsumerGroupsPage;
