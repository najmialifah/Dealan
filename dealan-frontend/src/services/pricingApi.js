import api from '../config/api';

export const negotiatePrice = async (payload) => {
  // endpoint: POST /pricing/negotiate based on pricing-service routes
  const response = await api.post('/pricing/negotiate', payload);
  return response.data;
};
