import client from './client';

export const testingApi = {
  // Testing Dashboard
  getDashboard() {
    return client.get('/testing/dashboard');
  },

  // Master Test Plans
  getTestPlans(params = {}) {
    return client.get('/test-plans', { params });
  },

  getTestPlanById(id) {
    return client.get(`/test-plans/${id}`);
  },

  createTestPlan(data) {
    return client.post('/test-plans', data);
  },

  updateTestPlan(id, data) {
    return client.put(`/test-plans/${id}`, data);
  },

  deleteTestPlan(id) {
    return client.delete(`/test-plans/${id}`);
  },

  // Test Plan Versions
  createVersion(testPlanId, data) {
    return client.post(`/test-plans/${testPlanId}/versions`, data);
  },

  getVersionById(id) {
    return client.get(`/test-plan-versions/${id}`);
  },

  updateVersion(id, data) {
    return client.put(`/test-plan-versions/${id}`, data);
  },

  releaseVersion(id, data = {}) {
    return client.post(`/test-plan-versions/${id}/release`, data);
  },

  cloneVersion(id, data) {
    return client.post(`/test-plan-versions/${id}/clone`, data);
  },

  // Test Cases
  createTestCase(versionId, data) {
    return client.post(`/test-plan-versions/${versionId}/test-cases`, data);
  },

  updateTestCase(id, data) {
    return client.put(`/test-cases/${id}`, data);
  },

  deleteTestCase(id) {
    return client.delete(`/test-cases/${id}`);
  },

  // Prototype Test Sessions
  getPrototypeSessions(prototypeId) {
    return client.get(`/prototypes/${prototypeId}/test-sessions`);
  },

  createPrototypeSession(prototypeId, data) {
    return client.post(`/prototypes/${prototypeId}/test-sessions`, data);
  },

  getSessionById(id) {
    return client.get(`/test-sessions/${id}`);
  },

  updateSession(id, data) {
    return client.put(`/test-sessions/${id}`, data);
  },

  // Test Executions & Measurements
  getSessionExecutions(sessionId) {
    return client.get(`/test-sessions/${sessionId}/executions`);
  },

  recordExecution(sessionId, data) {
    return client.post(`/test-sessions/${sessionId}/executions`, data);
  },

  createNextRun(sessionId, caseId) {
    return client.post(`/test-sessions/${sessionId}/executions/${caseId}/next-run`);
  },

  addEvidence(executionId, data) {
    return client.post(`/test-executions/${executionId}/evidence`, data);
  },

  // Engineering Findings
  getFindings(params = {}) {
    return client.get('/findings', { params });
  },

  getFindingById(id) {
    return client.get(`/findings/${id}`);
  },

  createFinding(data) {
    return client.post('/findings', data);
  },

  updateFinding(id, data) {
    return client.put(`/findings/${id}`, data);
  },

  deleteFinding(id) {
    return client.delete(`/findings/${id}`);
  },

  // Validation Decisions
  getValidationDecisions(sessionId) {
    return client.get(`/test-sessions/${sessionId}/validation-decisions`);
  },

  submitValidationDecision(sessionId, data) {
    return client.post(`/test-sessions/${sessionId}/validation-decisions`, data);
  }
};
