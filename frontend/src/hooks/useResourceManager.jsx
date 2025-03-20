// src/hooks/useResourceManager.jsx
import { useState, useCallback } from 'react';
import api from '../services/api';

const resourceEndpoints = {
  topics: {
    getAll: api.getTopics,
    getOne: api.getTopic,
    create: api.createTopic,
    update: null, // Not available in API
    delete: api.deleteTopic
  },
  brokers: {
    getAll: api.getBrokers,
    getOne: api.getBrokerDetails
  },
  consumers: {
    getAll: api.getConsumerGroups,
    getOne: api.getConsumerGroup,
    create: api.createConsumer,
    delete: api.stopConsumer
  },
  schemas: {
    // These would need to be implemented in the API
    getAll: () => Promise.resolve([]),
    create: (data) => Promise.resolve(data),
    update: (id, data) => Promise.resolve(data),
    delete: (id) => Promise.resolve(true)
  },
  acls: {
    // These would need to be implemented in the API
    getAll: () => Promise.resolve([]),
    create: (data) => Promise.resolve(data),
    update: (id, data) => Promise.resolve(data),
    delete: (id) => Promise.resolve(true)
  }
};

export const useResourceManager = (resourceType, config) => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState(config?.initialData || {});

  const endpoints = resourceEndpoints[resourceType] || {};

  const fetchData = useCallback(async () => {
    if (!endpoints.getAll) {
      setError(`No API method to fetch ${resourceType}`);
      return;
    }

    setLoading(true);
    try {
      const result = await endpoints.getAll();
      const formattedData = Array.isArray(result) ? result : [];
      
      // For consumer groups, we need to transform the groupId
      if (resourceType === 'consumers') {
        setData(formattedData.map(item => ({
          ...item,
          groupId: item.groupId || item.group_id,
          id: item.groupId || item.group_id
        })));
      } else {
        // For other resources, ensure they have an ID
        setData(formattedData.map(item => ({
          ...item,
          id: item.id || item.name || item.subject
        })));
      }
      
      setError(null);
    } catch (err) {
      setError(`Failed to fetch ${resourceType}: ${err.message}`);
    } finally {
      setLoading(false);
    }
  }, [resourceType, endpoints]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (isEditing && !endpoints.update) {
      setError(`API doesn't support updating ${resourceType}`);
      return false;
    }
    
    if (!isEditing && !endpoints.create) {
      setError(`API doesn't support creating ${resourceType}`);
      return false;
    }

    try {
      if (isEditing) {
        await endpoints.update(formData.id, formData);
      } else {
        await endpoints.create(formData);
      }
      
      await fetchData();
      setFormData(config?.initialData || {});
      setIsModalOpen(false);
      return true;
    } catch (err) {
      setError(`Failed to ${isEditing ? 'update' : 'create'} ${resourceType}: ${err.message}`);
      return false;
    }
  };

  const handleAdd = () => {
    setIsEditing(false);
    setFormData(config?.initialData || {});
    setIsModalOpen(true);
  };

  const handleEdit = (item) => {
    setIsEditing(true);
    setFormData(item);
    setIsModalOpen(true);
  };

  const handleDelete = async (id) => {
    if (!endpoints.delete) {
      setError(`API doesn't support deleting ${resourceType}`);
      return false;
    }

    if (window.confirm(`Are you sure you want to delete this ${resourceType.slice(0, -1)}?`)) {
      try {
        await endpoints.delete(id);
        await fetchData();
        return true;
      } catch (err) {
        setError(`Failed to delete ${resourceType}: ${err.message}`);
        return false;
      }
    }
    return false;
  };

  return {
    data,
    loading,
    error,
    setError,
    fetchData,
    isModalOpen,
    setIsModalOpen,
    isEditing,
    formData,
    setFormData,
    handleSubmit,
    handleAdd,
    handleEdit,
    handleDelete
  };
};
