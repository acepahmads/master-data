import apiClient from './client';

export const usersApi = {
  getAll(params = {}) {
    return apiClient.get('/users', { params });
  },
  getById(id) {
    return apiClient.get(`/users/${id}`);
  },
  create(data) {
    return apiClient.post('/users', data);
  },
  update(id, data) {
    return apiClient.put(`/users/${id}`, data);
  },
  toggleStatus(id) {
    return apiClient.patch(`/users/${id}/status`);
  },
  resetPassword(id, password) {
    return apiClient.post(`/users/${id}/reset-password`, { password });
  },
  delete(id) {
    return apiClient.delete(`/users/${id}`);
  },
};

export default usersApi;
