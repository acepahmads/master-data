import apiClient from './client';

export const unitsApi = {
  getAll(params = {}) {
    return apiClient.get('/units', { params });
  },
  getById(id) {
    return apiClient.get(`/units/${id}`);
  },
  create(data) {
    return apiClient.post('/units', data);
  },
  update(id, data) {
    return apiClient.put(`/units/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/units/${id}`);
  },
};

export default unitsApi;
