import axios from 'axios';
// We use a fallback if process.env.API_BASE_URL is not set by a bundler
// 10.0.2.2 is the localhost alias for Android Emulators
const BASE_URL = 'http://10.0.2.2:8000';

import AsyncStorage from '@react-native-async-storage/async-storage';
import { Alert } from 'react-native';

const api = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
});

// Interceptor to attach JWT token
api.interceptors.request.use(
  async (config) => {
    try {
      const token = await AsyncStorage.getItem('userToken');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    } catch (error) {
      console.error('Error fetching token from AsyncStorage:', error);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor to handle errors globally
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const message = error.response.data?.error || 'Terjadi kesalahan pada server';
      Alert.alert('Error', message);
    } else {
      Alert.alert('Error', 'Tidak dapat terhubung ke server');
    }
    return Promise.reject(error);
  }
);

export default api;
