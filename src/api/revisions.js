import apiClient from './client';

export const revisionsApi = {
  getAll(params = {}) {
    return apiClient.get('/hardware-revisions', { params });
  },
  getById(id) {
    return apiClient.get(`/hardware-revisions/${id}`);
  },
  create(data) {
    return apiClient.post('/hardware-revisions', data);
  },
  update(id, data) {
    return apiClient.put(`/hardware-revisions/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/hardware-revisions/${id}`);
  }
};

export default revisionsApi;
