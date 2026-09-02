import client from './client';

export const importApi = {
  getBatches(params = {}) {
    return client.get('/data-import/batches', { params });
  },

  getBatchById(id) {
    return client.get(`/data-import/batches/${id}`);
  },

  createBatch(data) {
    return client.post('/data-import/batches', data);
  },

  uploadBatchFiles(batchId, formData) {
    return client.post(`/data-import/batches/${batchId}/files`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
  },

  parseBatch(batchId, selectedSheetsByFile = {}, headerRowOverrides = {}) {
    return client.post(`/data-import/batches/${batchId}/parse`, { selectedSheetsByFile, headerRowOverrides });
  },

  applyMapping(batchId, mapping) {
    return client.post(`/data-import/batches/${batchId}/mapping`, { mapping });
  },

  runAIClassification(batchId) {
    return client.post(`/data-import/batches/${batchId}/ai-classify`);
  },

  generateMissingCodes(batchId, payload = {}) {
    return client.post(`/data-import/batches/${batchId}/generate-codes`, payload);
  },

  getStagedRows(batchId, params = {}) {
    return client.get(`/data-import/batches/${batchId}/rows`, { params });
  },

  triageRow(rowId, payload) {
    return client.put(`/data-import/rows/${rowId}`, payload);
  },

  approveBatch(batchId) {
    return client.post(`/data-import/batches/${batchId}/approve`);
  },

  commitBatch(batchId, payload = {}) {
    return client.post(`/data-import/batches/${batchId}/commit`, payload);
  },

  deleteBatch(batchId) {
    return client.delete(`/data-import/batches/${batchId}`);
  }
};

export default importApi;
