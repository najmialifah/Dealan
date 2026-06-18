import api from '../config/api';

export const createOrder = async (payload) => {
  const response = await api.post('/api/v1/orders/', payload);
  return response.data;
};
