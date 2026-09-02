import apiClient from './client';

export const productsApi = {
  getAll(params = {}) {
    return apiClient.get('/products', { params });
  },
  getById(id) {
    return apiClient.get(`/products/${id}`);
  },
  create(data) {
    return apiClient.post('/products', data);
  },
  update(id, data) {
    return apiClient.put(`/products/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/products/${id}`);
  },

  // Trading Detail (Phase 2C)
  getTradingDetail(id) {
    return apiClient.get(`/products/${id}/trading-detail`);
  },
  updateTradingDetail(id, data) {
    return apiClient.put(`/products/${id}/trading-detail`, data);
  },

  // Project Structure (Phase 2C)
  getProjectItems(id, params = {}) {
    return apiClient.get(`/products/${id}/project-items`, { params });
  },
  createProjectItem(projectId, data) {
    return apiClient.post(`/products/${projectId}/project-items`, data);
  },
  updateProjectItem(itemId, data) {
    return apiClient.put(`/project-items/${itemId}`, data);
  },
  deleteProjectItem(itemId) {
    return apiClient.delete(`/project-items/${itemId}`);
  },
  reorderProjectItems(projectId, itemIds) {
    return apiClient.post(`/products/${projectId}/project-items/reorder`, { itemIds });
  },
};

export default productsApi;
