import api from '../config/api';

export const submitRating = async (payload) => {
  const response = await api.post('/rating', payload);
  return response.data;
};
