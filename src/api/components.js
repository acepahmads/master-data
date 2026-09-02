import apiClient from './client';

export const componentsApi = {
  getAll(params = {}) {
    return apiClient.get('/components', { params });
  },
  getById(id) {
    return apiClient.get(`/components/${id}`);
  },
  create(data) {
    return apiClient.post('/components', data);
  },
  update(id, data) {
    return apiClient.put(`/components/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/components/${id}`);
  },
  addDocument(componentId, data) {
    return apiClient.post(`/components/${componentId}/documents`, data);
  },
  deleteDocument(componentId, docId) {
    return apiClient.delete(`/components/${componentId}/documents/${docId}`);
  },
};

export default componentsApi;
