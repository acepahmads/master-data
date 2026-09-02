import apiClient from './client';

export const manufacturersApi = {
  getAll(params = {}) {
    return apiClient.get('/manufacturers', { params });
  },
  getById(id) {
    return apiClient.get(`/manufacturers/${id}`);
  },
  create(data) {
    return apiClient.post('/manufacturers', data);
  },
  update(id, data) {
    return apiClient.put(`/manufacturers/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/manufacturers/${id}`);
  },
};

export default manufacturersApi;
