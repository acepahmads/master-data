import apiClient from './client';

export const categoriesApi = {
  getTree() {
    return apiClient.get('/categories/tree');
  },
  getFlat() {
    return apiClient.get('/categories');
  },
  getById(id) {
    return apiClient.get(`/categories/${id}`);
  },
  create(data) {
    return apiClient.post('/categories', data);
  },
  update(id, data) {
    return apiClient.put(`/categories/${id}`, data);
  },
  delete(id) {
    return apiClient.delete(`/categories/${id}`);
  },
};

export default categoriesApi;
