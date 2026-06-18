import api from '../config/api';

export const findDriver = async (payload) => {
  const response = await api.post('/api/v1/match', payload);
  return response.data;
};
