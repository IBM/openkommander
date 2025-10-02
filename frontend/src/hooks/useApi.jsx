import { useState, useCallback } from 'react';
import { API_BASE_URL } from '../config/constants';

export const useApi = (endpoint) => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`);
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const result = await response.json();
      setData(result);
      setError(null);
    } catch (err) {
      setError(`Failed to fetch ${endpoint}: ${err.message}`);
    } finally {
      setLoading(false);
    }
  }, [endpoint]);

  const addItem = useCallback(async (itemData) => {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(itemData),
      });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const result = await response.json();
      await fetchData();
      return result;
    } catch (err) {
      setError(`Failed to add ${endpoint}: ${err.message}`);
      return null;
    }
  }, [endpoint, fetchData]);

  const updateItem = useCallback(async (id, itemData) => {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(itemData),
      });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const result = await response.json();
      await fetchData();
      return result;
    } catch (err) {
      setError(`Failed to update ${endpoint}: ${err.message}`);
      return null;
    }
  }, [endpoint, fetchData]);

  const deleteItem = useCallback(async (id) => {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}/${id}`, {
        method: 'DELETE',
      });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      await fetchData();
      return true;
    } catch (err) {
      setError(`Failed to delete ${endpoint}: ${err.message}`);
      return false;
    }
  }, [endpoint, fetchData]);

  return {
    data,
    loading,
    error,
    setError,
    fetchData,
    addItem,
    updateItem,
    deleteItem
  };
};
