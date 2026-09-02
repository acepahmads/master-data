import client from './client';

export const bomApi = {
  // Get BOM by Hardware Revision ID
  getByRevision(revisionId) {
    return client.get(`/hardware-revisions/${revisionId}/bom`);
  },

  // Create or Initialize BOM for Revision
  create(revisionId, data) {
    return client.post(`/hardware-revisions/${revisionId}/bom`, data);
  },

  // Get BOM by ID
  getById(id) {
    return client.get(`/boms/${id}`);
  },

  // Update BOM metadata / margin
  update(id, data) {
    return client.put(`/boms/${id}`, data);
  },

  // Get BOM Cost Summary (top cost drivers, category breakdown)
  getCostSummary(id) {
    return client.get(`/boms/${id}/cost-summary`);
  },

  // Add Item to BOM
  addItem(bomId, data) {
    return client.post(`/boms/${bomId}/items`, data);
  },

  // Update BOM Item
  updateItem(itemId, data) {
    return client.put(`/bom-items/${itemId}`, data);
  },

  // Delete BOM Item
  deleteItem(itemId) {
    return client.delete(`/bom-items/${itemId}`);
  },

  // Reorder BOM Items
  reorderItems(bomId, itemIds) {
    return client.post(`/boms/${bomId}/items/reorder`, { itemIds });
  },

  // Get R&D Product Cost Summary
  getProductCostSummary(productId) {
    return client.get(`/products/${productId}/cost-summary`);
  },

  // Get Project Solution Cost Summary
  getProjectCostSummary(projectId) {
    return client.get(`/products/${projectId}/project-cost-summary`);
  }
};

export default bomApi;
