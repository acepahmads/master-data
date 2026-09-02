import client from './client';

export const prototypesApi = {
  // Get dashboard KPIs and summaries
  getDashboard() {
    return client.get('/prototypes/dashboard');
  },

  // Get all prototype builds with filters and pagination
  getAll(params = {}) {
    return client.get('/prototypes', { params });
  },

  // Get single prototype by ID or Code
  getById(id) {
    return client.get(`/prototypes/${id}`);
  },

  // Get prototype builds for a specific Hardware Revision
  getByRevision(revisionId) {
    return client.get(`/revisions/${revisionId}/prototypes`);
  },

  // Create new prototype build with frozen BOM snapshot
  create(data) {
    return client.post('/prototypes', data);
  },

  // Update prototype metadata & status
  update(id, data) {
    return client.put(`/prototypes/${id}`, data);
  },

  // Delete prototype build
  delete(id) {
    return client.delete(`/prototypes/${id}`);
  },

  // Get Frozen BOM Snapshot
  getBOMSnapshot(id) {
    return client.get(`/prototypes/${id}/bom-snapshot`);
  },

  // Component Preparation & Readiness tracking
  getComponents(id) {
    return client.get(`/prototypes/${id}/components`);
  },

  updateComponent(id, componentId, data) {
    return client.put(`/prototypes/${id}/components/${componentId}`, data);
  },

  // Assembly Stages workflow
  getAssembly(id) {
    return client.get(`/prototypes/${id}/assembly`);
  },

  updateAssemblyStage(id, stageId, data) {
    return client.put(`/prototypes/${id}/assembly/${stageId}`, data);
  },

  // Engineering Notes
  getNotes(id) {
    return client.get(`/prototypes/${id}/notes`);
  },

  createNote(id, data) {
    return client.post(`/prototypes/${id}/notes`, data);
  },

  updateNote(id, noteId, data) {
    return client.put(`/prototypes/${id}/notes/${noteId}`, data);
  },

  deleteNote(id, noteId) {
    return client.delete(`/prototypes/${id}/notes/${noteId}`);
  }
};

export default prototypesApi;
