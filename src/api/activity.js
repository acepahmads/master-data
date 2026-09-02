import apiClient from './client';

export const activityApi = {
  getAll(params = {}) {
    return apiClient.get('/activity', { params });
  },
};

export default activityApi;
