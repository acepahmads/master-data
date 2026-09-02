import client from './client';

export const changesApi = {
  // Dashboard & Traceability
  getDashboard() {
    return client.get('/changes/dashboard');
  },

  getTraceability(productId = '') {
    const params = productId ? { productId } : {};
    return client.get('/changes/traceability', { params });
  },

  // ECR CRUD & Lifecycle
  getECRs(params = {}) {
    return client.get('/ecr', { params });
  },

  getECRById(id) {
    return client.get(`/ecr/${id}`);
  },

  createECR(data) {
    return client.post('/ecr', data);
  },

  updateECR(id, data) {
    return client.put(`/ecr/${id}`, data);
  },

  deleteECR(id) {
    return client.delete(`/ecr/${id}`);
  },

  submitECR(id) {
    return client.post(`/ecr/${id}/submit`);
  },

  // ECR Change Items
  addChangeItem(ecrId, data) {
    return client.post(`/ecr/${ecrId}/items`, data);
  },

  deleteChangeItem(itemId) {
    return client.delete(`/ecr/items/${itemId}`);
  },

  // ECR Impact Analysis
  saveImpact(ecrId, data) {
    return client.post(`/ecr/${ecrId}/impact`, data);
  },

  // ECR Technical Reviews
  submitReview(ecrId, data) {
    return client.post(`/ecr/${ecrId}/reviews`, data);
  },

  // ECR CCB Approvals
  submitApproval(ecrId, data) {
    return client.post(`/ecr/${ecrId}/approvals`, data);
  },

  // ECO Creation & Lifecycle
  createECO(data) {
    return client.post('/eco', data);
  },

  getECOs(params = {}) {
    return client.get('/eco', { params });
  },

  getECOById(id) {
    return client.get(`/eco/${id}`);
  },

  approveECO(id) {
    return client.post(`/eco/${id}/approve`);
  },

  getChangePreview(id) {
    return client.get(`/eco/${id}/change-preview`);
  },

  implementECO(id) {
    return client.post(`/eco/${id}/implement`);
  },

  verifyECO(id, data = {}) {
    return client.post(`/eco/${id}/verify`, data);
  },

  closeECO(id) {
    return client.post(`/eco/${id}/close`);
  },
};
