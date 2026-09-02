import apiClient from './client';

export const rolesApi = {
  getAll() {
    return apiClient.get('/roles');
  },
  getById(id) {
    return apiClient.get(`/roles/${id}`);
  },
  create(data) {
    return apiClient.post('/roles', data);
  },
  update(id, data) {
    return apiClient.put(`/roles/${id}`, data);
  },
  updatePermissions(id, permissions) {
    return apiClient.put(`/roles/${id}/permissions`, { permissions });
  },
  getAllPermissions() {
    return apiClient.get('/permissions');
  },
};

export default rolesApi;
