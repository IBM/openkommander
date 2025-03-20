// src/services/api.js
import { API_BASE_URL } from '../config/constants';

const apiRequest = async (endpoint, options = {}) => {
  try {
    const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    // For DELETE requests that don't return JSON
    if (response.status === 204 || options.method === 'DELETE') {
      return { success: true };
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error(`API Error (${endpoint}):`, error);
    throw error;
  }
};

// API endpoints
export const api = {
  // Health check
  getHealth: () => apiRequest('health'),

  // Brokers
  getBrokers: () => apiRequest('brokers'),
  getBrokerDetails: (id) => apiRequest(`brokers/${id}`),

  // Topics
  getTopics: () => apiRequest('topics'),
  getTopic: (name) => apiRequest(`topics/${name}`),
  createTopic: (data) => apiRequest('topics', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  deleteTopic: (name) => apiRequest(`topics/${name}`, {
    method: 'DELETE',
  }),

  // Consumer Groups
  getConsumerGroups: () => apiRequest('consumers'),
  getConsumerGroup: (group) => apiRequest(`consumers/${group}`),
  createConsumer: (data) => apiRequest('consumers', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  stopConsumer: (id) => apiRequest(`consumers/${id}`, {
    method: 'DELETE',
  }),

  // Messages
  produceMessage: (topic, message, key = null) => {
    const url = key ? `messages/${topic}?key=${key}` : `messages/${topic}`;
    return apiRequest(url, {
      method: 'POST',
      body: JSON.stringify(message),
    });
  },

  // Clusters
  getClusters: () => apiRequest('clusters'),
  getCluster: (name) => apiRequest(`clusters/${name}`),
  getClusterBrokers: (name) => apiRequest(`clusters/${name}/brokers`),
  getClusterTopics: (name) => apiRequest(`clusters/${name}/topics`),

  // Custom methods to adapt to your existing UI
  resetConsumerOffsets: (groupId) => apiRequest(`consumers/${groupId}/reset`, {
    method: 'POST',
  }),
  getConsumerGroupAssignments: (groupId) => apiRequest(`consumers/${groupId}/assignments`),
};

export default api;
