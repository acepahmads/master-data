import { defineStore } from 'pinia';
import prototypesApi from '@/api/prototypes';

export const usePrototypeStore = defineStore('prototype', {
  state: () => ({
    prototypes: [],
    currentPrototype: null,
    dashboardStats: {
      totalBuilds: 0,
      plannedCount: 0,
      componentPrepCount: 0,
      inAssemblyCount: 0,
      powerOnCount: 0,
      readyForTestingCount: 0,
      completedCount: 0,
      onHoldCount: 0
    },
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      limit: 50,
      totalPages: 1
    }
  }),

  getters: {
    totalPrototypesCount: (state) => state.dashboardStats?.totalBuilds || state.prototypes.length,
    draftCount: (state) => state.dashboardStats?.plannedCount || state.prototypes.filter(p => p.status === 'PLANNED' || p.status === 'DRAFT').length,
    componentPrepCount: (state) => state.dashboardStats?.componentPrepCount || state.prototypes.filter(p => p.status === 'PREPARING' || p.status === 'COMPONENT_PREPARATION').length,
    inAssemblyCount: (state) => state.dashboardStats?.inAssemblyCount || state.prototypes.filter(p => p.status === 'IN_PROGRESS' || p.status === 'ASSEMBLY' || p.status === 'IN_ASSEMBLY').length,
    readyForTestingCount: (state) => state.dashboardStats?.readyForTestingCount || state.prototypes.filter(p => p.status === 'READY_FOR_TESTING' || p.status === 'POWER_ON').length,
    completedCount: (state) => state.dashboardStats?.completedCount || state.prototypes.filter(p => p.status === 'COMPLETED').length,

    getPrototypesByRevision: (state) => (revisionId) => {
      return state.prototypes.filter(p => p.hardwareRevisionId === revisionId || p.hardwareRevisionCode === revisionId);
    },
    getPrototypesByProduct: (state) => (productId) => {
      return state.prototypes.filter(p => p.productId === productId || p.productCode === productId);
    }
  },

  actions: {
    async fetchDashboard() {
      this.loading = true;
      try {
        const res = await prototypesApi.getDashboard();
        if (res && res.data) {
          this.dashboardStats = res.data;
        }
        return this.dashboardStats;
      } catch (err) {
        console.error('[PrototypeStore] fetchDashboard failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchPrototypes(params = {}) {
      this.loading = true;
      try {
        const query = { limit: 100, ...params };
        const res = await prototypesApi.getAll(query);
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.items || []);
          this.prototypes = list.map(item => this.normalizePrototype(item));
          if (res.meta) {
            this.pagination = res.meta;
          }
        }
        return this.prototypes;
      } catch (err) {
        console.error('[PrototypeStore] fetchPrototypes failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchPrototype(id) {
      this.loading = true;
      try {
        const res = await prototypesApi.getById(id);
        if (res && res.data) {
          this.currentPrototype = this.normalizePrototype(res.data);
          // Update in prototypes list as well
          const idx = this.prototypes.findIndex(p => p.id === this.currentPrototype.id);
          if (idx !== -1) {
            this.prototypes.splice(idx, 1, this.currentPrototype);
          }
        }
        return this.currentPrototype;
      } catch (err) {
        console.error('[PrototypeStore] fetchPrototype failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createPrototype(payload) {
      this.loading = true;
      try {
        const res = await prototypesApi.create(payload);
        if (res && res.data) {
          const created = this.normalizePrototype(res.data);
          this.prototypes.unshift(created);
          if (this.dashboardStats) {
            this.dashboardStats.totalBuilds++;
          }
          return created;
        }
        return null;
      } catch (err) {
        console.error('[PrototypeStore] createPrototype failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updatePrototype(id, payload) {
      this.loading = true;
      try {
        const res = await prototypesApi.update(id, payload);
        if (res && res.data) {
          const updated = this.normalizePrototype(res.data);
          if (this.currentPrototype && this.currentPrototype.id === id) {
            this.currentPrototype = updated;
          }
          const idx = this.prototypes.findIndex(p => p.id === id);
          if (idx !== -1) {
            this.prototypes.splice(idx, 1, updated);
          }
          return updated;
        }
        return null;
      } catch (err) {
        console.error('[PrototypeStore] updatePrototype failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updatePrototypeStatus(id, status) {
      return this.updatePrototype(id, { status });
    },

    async deletePrototype(id) {
      this.loading = true;
      try {
        await prototypesApi.delete(id);
        this.prototypes = this.prototypes.filter(p => p.id !== id && p.code !== id);
        if (this.currentPrototype && (this.currentPrototype.id === id || this.currentPrototype.code === id)) {
          this.currentPrototype = null;
        }
        if (this.dashboardStats && this.dashboardStats.totalBuilds > 0) {
          this.dashboardStats.totalBuilds--;
        }
      } catch (err) {
        console.error('[PrototypeStore] deletePrototype failed:', err);
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchBOMSnapshot(id) {
      try {
        const res = await prototypesApi.getBOMSnapshot(id);
        return res?.data || null;
      } catch (err) {
        console.error('[PrototypeStore] fetchBOMSnapshot failed:', err);
        throw err;
      }
    },

    async fetchComponentPreparation(id) {
      try {
        const res = await prototypesApi.getComponents(id);
        if (res && res.data) {
          if (this.currentPrototype) {
            this.currentPrototype.componentsReadiness = res.data.components || [];
            this.currentPrototype.readinessPercentage = res.data.readinessPercentage || 0;
          }
          return res.data;
        }
        return null;
      } catch (err) {
        console.error('[PrototypeStore] fetchComponentPreparation failed:', err);
        throw err;
      }
    },

    async updateComponentPreparation(prototypeId, compPrepId, payload) {
      try {
        const res = await prototypesApi.updateComponent(prototypeId, compPrepId, payload);
        if (res && res.data) {
          if (this.currentPrototype) {
            this.currentPrototype.readinessPercentage = res.data.readinessPercentage || this.currentPrototype.readinessPercentage;
            if (this.currentPrototype.componentsReadiness) {
              const idx = this.currentPrototype.componentsReadiness.findIndex(c => c.id === compPrepId || c.componentId === compPrepId);
              if (idx !== -1) {
                this.currentPrototype.componentsReadiness.splice(idx, 1, res.data.component);
              }
            }
          }
          return res.data;
        }
        return null;
      } catch (err) {
        console.error('[PrototypeStore] updateComponentPreparation failed:', err);
        throw err;
      }
    },

    async fetchAssemblyProgress(id) {
      try {
        const res = await prototypesApi.getAssembly(id);
        if (res && res.data && this.currentPrototype) {
          this.currentPrototype.assemblyStages = res.data || [];
        }
        return res?.data || [];
      } catch (err) {
        console.error('[PrototypeStore] fetchAssemblyProgress failed:', err);
        throw err;
      }
    },

    async updateAssemblyStage(prototypeId, stageId, payload) {
      try {
        const res = await prototypesApi.updateAssemblyStage(prototypeId, stageId, payload);
        if (res && res.data && this.currentPrototype && this.currentPrototype.assemblyStages) {
          const idx = this.currentPrototype.assemblyStages.findIndex(s => s.id === stageId || s.stage === stageId);
          if (idx !== -1) {
            this.currentPrototype.assemblyStages.splice(idx, 1, res.data);
          }
        }
        return res?.data || null;
      } catch (err) {
        console.error('[PrototypeStore] updateAssemblyStage failed:', err);
        throw err;
      }
    },

    async fetchEngineeringNotes(id) {
      try {
        const res = await prototypesApi.getNotes(id);
        if (res && res.data && this.currentPrototype) {
          this.currentPrototype.engineeringNotes = res.data || [];
        }
        return res?.data || [];
      } catch (err) {
        console.error('[PrototypeStore] fetchEngineeringNotes failed:', err);
        throw err;
      }
    },

    async addEngineeringNote(prototypeId, payload) {
      try {
        const res = await prototypesApi.createNote(prototypeId, payload);
        if (res && res.data && this.currentPrototype) {
          if (!this.currentPrototype.engineeringNotes) {
            this.currentPrototype.engineeringNotes = [];
          }
          this.currentPrototype.engineeringNotes.unshift(res.data);
        }
        return res?.data || null;
      } catch (err) {
        console.error('[PrototypeStore] addEngineeringNote failed:', err);
        throw err;
      }
    },

    async deleteEngineeringNote(prototypeId, noteId) {
      try {
        await prototypesApi.deleteNote(prototypeId, noteId);
        if (this.currentPrototype && this.currentPrototype.engineeringNotes) {
          this.currentPrototype.engineeringNotes = this.currentPrototype.engineeringNotes.filter(n => n.id !== noteId);
        }
      } catch (err) {
        console.error('[PrototypeStore] deleteEngineeringNote failed:', err);
        throw err;
      }
    },

    getPrototypeById(id) {
      if (this.currentPrototype && (this.currentPrototype.id === id || this.currentPrototype.code === id)) {
        return this.currentPrototype;
      }
      return this.prototypes.find(p => p.id === id || p.code === id) || null;
    },

    normalizePrototype(item) {
      if (!item) return null;
      let snapshot = {
        snapshotId: item.bomSnapshotId || `SNP-${item.code}`,
        frozenAt: item.createdAt,
        totalItems: 0,
        totalQuantity: 0,
        totalEstimatedCost: item.bomSnapshotCost || 0,
        currency: item.bomSnapshotCurrency || 'USD',
        items: []
      };

      if (item.bomSnapshotJson) {
        try {
          snapshot = typeof item.bomSnapshotJson === 'string' ? JSON.parse(item.bomSnapshotJson) : item.bomSnapshotJson;
        } catch (e) {
          console.warn('[PrototypeStore] Failed to parse bomSnapshotJson:', e);
        }
      } else if (item.bomSnapshot) {
        snapshot = item.bomSnapshot;
      }

      return {
        ...item,
        bomSnapshot: snapshot,
        componentsReadiness: item.componentPreparations || item.componentsReadiness || [],
        assemblyStages: item.assemblyStages || [],
        engineeringNotes: item.engineeringNotes || [],
        activityHistory: item.activityHistory || []
      };
    }
  }
});

export default usePrototypeStore;
