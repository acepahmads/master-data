import { defineStore } from 'pinia';
import { testingApi } from '@/api/testing';

export const useTestingStore = defineStore('testing', {
  state: () => ({
    dashboardStats: {
      totalTestPlans: 0,
      activeTestSessions: 0,
      totalExecutions: 0,
      totalPassed: 0,
      totalFailed: 0,
      totalBlocked: 0,
      overallPassRate: 0.0,
      openFindingsCount: 0,
      criticalFindingsCount: 0,
      highFindingsCount: 0,
      changeCandidatesCount: 0,
      pendingValidations: 0,
      totalValidatedSessions: 0,
      totalRejectedSessions: 0
    },
    testPlans: [],
    testPlansPagination: { page: 1, limit: 20, total: 0 },
    currentTestPlan: null,
    currentVersion: null,
    prototypeSessions: {}, // key: prototypeId -> array of sessions
    currentSession: null,
    findings: [],
    findingsPagination: { page: 1, limit: 20, total: 0 },
    currentFinding: null,
    loading: false,
    error: null
  }),

  getters: {
    getPlanById: (state) => (id) => {
      return state.testPlans.find(p => p.id === id || p.code === id) || state.currentTestPlan;
    },
    getSessionsForPrototype: (state) => (protoId) => {
      return state.prototypeSessions[protoId] || [];
    },
    openFindings: (state) => {
      return state.findings.filter(f => f.status !== 'CLOSED');
    },
    changeCandidateFindings: (state) => {
      return state.findings.filter(f => f.changeCandidateStatus === 'CANDIDATE');
    }
  },

  actions: {
    async fetchDashboard() {
      this.loading = true;
      try {
        const res = await testingApi.getDashboard();
        if (res.data?.success && res.data.data) {
          this.dashboardStats = res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async fetchTestPlans(params = {}) {
      this.loading = true;
      try {
        const res = await testingApi.getTestPlans(params);
        if (res.data?.success && res.data.data) {
          this.testPlans = res.data.data;
          if (res.data.meta) {
            this.testPlansPagination = res.data.meta;
          }
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async fetchTestPlanById(id) {
      this.loading = true;
      try {
        const res = await testingApi.getTestPlanById(id);
        if (res.data?.success && res.data.data) {
          this.currentTestPlan = res.data.data;
          if (this.currentTestPlan.versions && this.currentTestPlan.versions.length > 0) {
            this.currentVersion = this.currentTestPlan.versions[0];
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createTestPlan(data) {
      this.loading = true;
      try {
        const res = await testingApi.createTestPlan(data);
        if (res.data?.success && res.data.data) {
          this.testPlans.unshift(res.data.data);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateTestPlan(id, data) {
      this.loading = true;
      try {
        const res = await testingApi.updateTestPlan(id, data);
        if (res.data?.success && res.data.data) {
          const idx = this.testPlans.findIndex(p => p.id === id);
          if (idx !== -1) {
            this.testPlans.splice(idx, 1, res.data.data);
          }
          if (this.currentTestPlan && this.currentTestPlan.id === id) {
            this.currentTestPlan = { ...this.currentTestPlan, ...res.data.data };
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteTestPlan(id) {
      this.loading = true;
      try {
        await testingApi.deleteTestPlan(id);
        this.testPlans = this.testPlans.filter(p => p.id !== id);
        if (this.currentTestPlan && this.currentTestPlan.id === id) {
          this.currentTestPlan = null;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createVersion(testPlanId, data) {
      this.loading = true;
      try {
        const res = await testingApi.createVersion(testPlanId, data);
        if (res.data?.success && res.data.data) {
          if (this.currentTestPlan) {
            if (!this.currentTestPlan.versions) this.currentTestPlan.versions = [];
            this.currentTestPlan.versions.unshift(res.data.data);
            this.currentVersion = res.data.data;
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async releaseVersion(versionId, data = {}) {
      this.loading = true;
      try {
        const res = await testingApi.releaseVersion(versionId, data);
        if (res.data?.success && res.data.data) {
          this.currentVersion = res.data.data;
          if (this.currentTestPlan && this.currentTestPlan.versions) {
            const idx = this.currentTestPlan.versions.findIndex(v => v.id === versionId);
            if (idx !== -1) {
              this.currentTestPlan.versions.splice(idx, 1, res.data.data);
            }
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async cloneVersion(versionId, data) {
      this.loading = true;
      try {
        const res = await testingApi.cloneVersion(versionId, data);
        if (res.data?.success && res.data.data) {
          if (this.currentTestPlan) {
            if (!this.currentTestPlan.versions) this.currentTestPlan.versions = [];
            this.currentTestPlan.versions.unshift(res.data.data);
            this.currentVersion = res.data.data;
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createTestCase(versionId, data) {
      this.loading = true;
      try {
        const res = await testingApi.createTestCase(versionId, data);
        if (res.data?.success && res.data.data) {
          if (this.currentVersion && this.currentVersion.id === versionId) {
            if (!this.currentVersion.testCases) this.currentVersion.testCases = [];
            this.currentVersion.testCases.push(res.data.data);
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateTestCase(id, data) {
      this.loading = true;
      try {
        const res = await testingApi.updateTestCase(id, data);
        if (res.data?.success && res.data.data) {
          if (this.currentVersion && this.currentVersion.testCases) {
            const idx = this.currentVersion.testCases.findIndex(tc => tc.id === id);
            if (idx !== -1) {
              this.currentVersion.testCases.splice(idx, 1, res.data.data);
            }
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteTestCase(id) {
      this.loading = true;
      try {
        await testingApi.deleteTestCase(id);
        if (this.currentVersion && this.currentVersion.testCases) {
          this.currentVersion.testCases = this.currentVersion.testCases.filter(tc => tc.id !== id);
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Prototype Test Sessions
    async fetchPrototypeSessions(prototypeId) {
      this.loading = true;
      try {
        const res = await testingApi.getPrototypeSessions(prototypeId);
        if (res.data?.success && res.data.data) {
          this.$patch({
            prototypeSessions: {
              ...this.prototypeSessions,
              [prototypeId]: res.data.data
            }
          });
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async createPrototypeSession(prototypeId, data) {
      this.loading = true;
      try {
        const res = await testingApi.createPrototypeSession(prototypeId, data);
        if (res.data?.success && res.data.data) {
          const list = this.prototypeSessions[prototypeId] || [];
          list.unshift(res.data.data);
          this.$patch({
            prototypeSessions: {
              ...this.prototypeSessions,
              [prototypeId]: list
            }
          });
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchSessionById(id) {
      this.loading = true;
      try {
        const res = await testingApi.getSessionById(id);
        if (res.data?.success && res.data.data) {
          this.currentSession = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async recordExecution(sessionId, data) {
      try {
        const res = await testingApi.recordExecution(sessionId, data);
        if (res.data?.success && res.data.data) {
          await this.fetchSessionById(sessionId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async createNextRun(sessionId, caseId) {
      try {
        const res = await testingApi.createNextRun(sessionId, caseId);
        if (res.data?.success && res.data.data) {
          await this.fetchSessionById(sessionId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async addEvidence(executionId, data) {
      try {
        const res = await testingApi.addEvidence(executionId, data);
        if (res.data?.success && res.data.data && this.currentSession) {
          await this.fetchSessionById(this.currentSession.id);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    // Findings
    async fetchFindings(params = {}) {
      this.loading = true;
      try {
        const res = await testingApi.getFindings(params);
        if (res.data?.success && res.data.data) {
          this.findings = res.data.data;
          if (res.data.meta) {
            this.findingsPagination = res.data.meta;
          }
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async createFinding(data) {
      this.loading = true;
      try {
        const res = await testingApi.createFinding(data);
        if (res.data?.success && res.data.data) {
          this.findings.unshift(res.data.data);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateFinding(id, data) {
      this.loading = true;
      try {
        const res = await testingApi.updateFinding(id, data);
        if (res.data?.success && res.data.data) {
          const idx = this.findings.findIndex(f => f.id === id);
          if (idx !== -1) {
            this.findings.splice(idx, 1, res.data.data);
          }
          if (this.currentFinding && this.currentFinding.id === id) {
            this.currentFinding = res.data.data;
          }
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Validation Decisions
    async submitValidationDecision(sessionId, data) {
      this.loading = true;
      try {
        const res = await testingApi.submitValidationDecision(sessionId, data);
        if (res.data?.success && res.data.data) {
          await this.fetchSessionById(sessionId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
