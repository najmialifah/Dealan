import api from '../config/api';

export const processPayment = async (payload) => {
  const response = await api.post('/payments', payload);
  return response.data;
};
