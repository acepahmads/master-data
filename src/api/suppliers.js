import apiClient from './client';

export const suppliersApi = {
  getAll(params = {}) {
    return apiClient.get('/suppliers', { params });
  },
  getById(id) {
    return apiClient.get(`/suppliers/${id}`);
  },
  create(data) {
    return apiClient.post('/suppliers', data);
  },
  update(id, data) {
    return apiClient.put(`/suppliers/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/suppliers/${id}`);
  },
};

export default suppliersApi;
