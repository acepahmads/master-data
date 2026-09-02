import apiClient from './client';

export const versionsApi = {
  getAll(params = {}) {
    return apiClient.get('/product-versions', { params });
  },
  getById(id) {
    return apiClient.get(`/product-versions/${id}`);
  },
  getLifecycle(productId) {
    return apiClient.get(`/products/${productId}/lifecycle`);
  },
  compare(baseId, targetId) {
    return apiClient.get('/product-versions/compare', {
      params: { base_id: baseId, target_id: targetId }
    });
  },
  create(data) {
    return apiClient.post('/product-versions', data);
  },
  update(id, data) {
    return apiClient.put(`/product-versions/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/product-versions/${id}`);
  }
};

export default versionsApi;
