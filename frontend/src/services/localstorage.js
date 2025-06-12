// src/services/api.js
import { API_BASE_URL } from '../config/constants';

const getLocalStorage = async (endpoint, options = {}) => {
    const data = JSON.parse(localStorage.getItem(endpoint)) || [];
    return data;
};

const saveLocalStorage = async (endpoint, data) => {
    const existingData = JSON.parse(localStorage.getItem(endpoint)) || [];
    localStorage.setItem(endpoint, JSON.stringify([...existingData, data]));
}

// Local Storage API endpoints
export const localStorageAPI = {

  // Brokers
  createBroker: (data) => saveLocalStorage('brokers', data),
  getBrokers: () => getLocalStorage('brokers'),
  getSelectedBroker: () => getLocalStorage('selectedBroker'),
};

export default localStorageAPI;
