import apiClient from './client';

export const uploadApi = {
  uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);

    return apiClient.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};

export default uploadApi;
