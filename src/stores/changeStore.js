import { defineStore } from 'pinia';
import { changesApi } from '@/api/changes';

export const useChangeStore = defineStore('change', {
  state: () => ({
    dashboardStats: {
      totalEcrs: 0,
      draftEcrs: 0,
      underReviewEcrs: 0,
      approvedEcrs: 0,
      rejectedEcrs: 0,
      convertedEcrs: 0,
      totalEcos: 0,
      draftEcos: 0,
      approvedEcos: 0,
      implementedEcos: 0,
      verifiedEcos: 0,
      closedEcos: 0,
      criticalHighEcrs: 0,
      avgImplementationCycleDays: 0,
    },
    ecrs: [],
    ecrPagination: { page: 1, limit: 20, total: 0 },
    currentECR: null,
    ecos: [],
    ecoPagination: { page: 1, limit: 20, total: 0 },
    currentECO: null,
    bomChangePreview: null,
    traceabilityChains: [],
    loading: false,
    error: null,
  }),

  getters: {
    getECRById: (state) => (id) => {
      return state.ecrs.find((e) => e.id === id || e.code === id) || state.currentECR;
    },
    getECOById: (state) => (id) => {
      return state.ecos.find((e) => e.id === id || e.code === id) || state.currentECO;
    },
    approvedECRs: (state) => {
      return state.ecrs.filter((e) => e.status === 'APPROVED');
    },
    pendingImplementationECOs: (state) => {
      return state.ecos.filter((e) => e.status === 'APPROVED');
    },
    activeECOs: (state) => {
      return state.ecos.filter((e) => e.status !== 'CLOSED');
    },
  },

  actions: {
    async fetchDashboard() {
      this.loading = true;
      try {
        const res = await changesApi.getDashboard();
        if (res.data?.success && res.data.data) {
          this.dashboardStats = res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async fetchECRs(params = {}) {
      this.loading = true;
      try {
        const res = await changesApi.getECRs(params);
        if (res.data?.success && res.data.data) {
          this.ecrs = res.data.data.items || [];
          this.ecrPagination = {
            page: res.data.data.page || 1,
            limit: res.data.data.limit || 20,
            total: res.data.data.total || 0,
          };
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async fetchECRById(id) {
      this.loading = true;
      try {
        const res = await changesApi.getECRById(id);
        if (res.data?.success && res.data.data) {
          this.currentECR = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createECR(data) {
      this.loading = true;
      try {
        const res = await changesApi.createECR(data);
        if (res.data?.success && res.data.data) {
          this.currentECR = res.data.data;
          await this.fetchECRs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateECR(id, data) {
      this.loading = true;
      try {
        const res = await changesApi.updateECR(id, data);
        if (res.data?.success && res.data.data) {
          this.currentECR = res.data.data;
          await this.fetchECRs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteECR(id) {
      this.loading = true;
      try {
        await changesApi.deleteECR(id);
        await this.fetchECRs();
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async submitECR(id) {
      this.loading = true;
      try {
        const res = await changesApi.submitECR(id);
        if (res.data?.success && res.data.data) {
          this.currentECR = res.data.data;
          await this.fetchECRs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async addChangeItem(ecrId, data) {
      try {
        const res = await changesApi.addChangeItem(ecrId, data);
        if (res.data?.success) {
          await this.fetchECRById(ecrId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async deleteChangeItem(ecrId, itemId) {
      try {
        await changesApi.deleteChangeItem(itemId);
        await this.fetchECRById(ecrId);
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async saveImpact(ecrId, data) {
      try {
        const res = await changesApi.saveImpact(ecrId, data);
        if (res.data?.success) {
          await this.fetchECRById(ecrId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async submitReview(ecrId, data) {
      try {
        const res = await changesApi.submitReview(ecrId, data);
        if (res.data?.success) {
          await this.fetchECRById(ecrId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async submitApproval(ecrId, data) {
      try {
        const res = await changesApi.submitApproval(ecrId, data);
        if (res.data?.success) {
          await this.fetchECRById(ecrId);
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      }
    },

    async createECO(data) {
      this.loading = true;
      try {
        const res = await changesApi.createECO(data);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          await this.fetchECOs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchECOs(params = {}) {
      this.loading = true;
      try {
        const res = await changesApi.getECOs(params);
        if (res.data?.success && res.data.data) {
          this.ecos = res.data.data.items || [];
          this.ecoPagination = {
            page: res.data.data.page || 1,
            limit: res.data.data.limit || 20,
            total: res.data.data.total || 0,
          };
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },

    async fetchECOById(id) {
      this.loading = true;
      try {
        const res = await changesApi.getECOById(id);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async approveECO(id) {
      this.loading = true;
      try {
        const res = await changesApi.approveECO(id);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          await this.fetchECOs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchChangePreview(ecoId) {
      this.loading = true;
      try {
        const res = await changesApi.getChangePreview(ecoId);
        if (res.data?.success && res.data.data) {
          this.bomChangePreview = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async implementECO(id) {
      this.loading = true;
      try {
        const res = await changesApi.implementECO(id);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          await this.fetchECOs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async verifyECO(id, data = {}) {
      this.loading = true;
      try {
        const res = await changesApi.verifyECO(id, data);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          await this.fetchECOs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async closeECO(id) {
      this.loading = true;
      try {
        const res = await changesApi.closeECO(id);
        if (res.data?.success && res.data.data) {
          this.currentECO = res.data.data;
          await this.fetchECOs();
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchTraceability(productId = '') {
      this.loading = true;
      try {
        const res = await changesApi.getTraceability(productId);
        if (res.data?.success && res.data.data) {
          this.traceabilityChains = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    },
  },
});
