import { defineStore } from 'pinia';
import { bomApi } from '@/api/bom';

export const useBOMStore = defineStore('bom', {
  state: () => ({
    currentBOM: null,
    costSummary: null,
    selectedItemId: null,
    loading: false,
    error: null,
    productCostSummary: {},
    projectCostSummary: {}
  }),

  getters: {
    selectedItem: (state) => {
      if (!state.currentBOM || !state.selectedItemId) return null;
      
      const findItem = (items) => {
        if (!items) return null;
        for (const item of items) {
          if (item.id === state.selectedItemId) return item;
          if (item.subItems && item.subItems.length > 0) {
            const found = findItem(item.subItems);
            if (found) return found;
          }
        }
        return null;
      };

      return findItem(state.currentBOM.items);
    },

    flatItems: (state) => {
      if (!state.currentBOM || !state.currentBOM.items) return [];
      
      const flat = [];
      const traverse = (items, depth = 0, parentName = '') => {
        for (const item of items) {
          flat.push({
            ...item,
            depth,
            parentName
          });
          if (item.subItems && item.subItems.length > 0) {
            traverse(item.subItems, depth + 1, item.name);
          }
        }
      };
      
      traverse(state.currentBOM.items);
      return flat;
    },

    subAssemblies: (state) => {
      if (!state.currentBOM || !state.currentBOM.items) return [];
      const list = [];
      const traverse = (items) => {
        for (const item of items) {
          if (item.itemType === 'SUB_ASSEMBLY') {
            list.push(item);
          }
          if (item.subItems && item.subItems.length > 0) {
            traverse(item.subItems);
          }
        }
      };
      traverse(state.currentBOM.items);
      return list;
    }
  },

  actions: {
    async fetchBOMByRevision(revisionId) {
      this.loading = true;
      this.error = null;
      try {
        const response = await bomApi.getByRevision(revisionId);
        if (response.data && response.data.success) {
          this.currentBOM = response.data.data;
          // Auto select first item if none selected
          if (this.currentBOM.items && this.currentBOM.items.length > 0 && !this.selectedItemId) {
            this.selectedItemId = this.currentBOM.items[0].id;
          }
          // Fetch cost summary in parallel
          this.fetchCostSummary(this.currentBOM.id);
        } else {
          this.currentBOM = null;
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to fetch BOM';
        this.currentBOM = null;
      } finally {
        this.loading = false;
      }
    },

    async createBOM(revisionId, data) {
      this.loading = true;
      try {
        const response = await bomApi.create(revisionId, data);
        if (response.data && response.data.success) {
          this.currentBOM = response.data.data;
          await this.fetchCostSummary(this.currentBOM.id);
          return response.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to initialize BOM';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateBOM(id, data) {
      try {
        const response = await bomApi.update(id, data);
        if (response.data && response.data.success) {
          if (this.currentBOM && this.currentBOM.id === id) {
            this.currentBOM = { ...this.currentBOM, ...response.data.data };
          }
          await this.fetchCostSummary(id);
          return response.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to update BOM';
        throw err;
      }
    },

    async fetchCostSummary(bomId) {
      try {
        const response = await bomApi.getCostSummary(bomId);
        if (response.data && response.data.success) {
          this.costSummary = response.data.data;
        }
      } catch (err) {
        console.error('Error fetching cost summary:', err);
      }
    },

    async addItem(bomId, itemData) {
      this.loading = true;
      try {
        const response = await bomApi.addItem(bomId, itemData);
        if (response.data && response.data.success) {
          // Re-fetch entire BOM to get fresh hierarchy and recalculations
          if (this.currentBOM) {
            const revId = this.currentBOM.hardwareRevisionId;
            await this.fetchBOMByRevision(revId);
          }
          return response.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to add BOM item';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateItem(itemId, itemData) {
      this.loading = true;
      try {
        const response = await bomApi.updateItem(itemId, itemData);
        if (response.data && response.data.success) {
          if (this.currentBOM) {
            const revId = this.currentBOM.hardwareRevisionId;
            await this.fetchBOMByRevision(revId);
          }
          return response.data.data;
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to update BOM item';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteItem(itemId) {
      this.loading = true;
      try {
        const response = await bomApi.deleteItem(itemId);
        if (response.data && response.data.success) {
          if (this.selectedItemId === itemId) {
            this.selectedItemId = null;
          }
          if (this.currentBOM) {
            const revId = this.currentBOM.hardwareRevisionId;
            await this.fetchBOMByRevision(revId);
          }
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to delete BOM item';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async reorderItems(bomId, itemIds) {
      try {
        await bomApi.reorderItems(bomId, itemIds);
        if (this.currentBOM) {
          const revId = this.currentBOM.hardwareRevisionId;
          await this.fetchBOMByRevision(revId);
        }
      } catch (err) {
        console.error('Failed to reorder items:', err);
      }
    },

    async fetchProductCostSummary(productId) {
      try {
        const response = await bomApi.getProductCostSummary(productId);
        if (response.data && response.data.success) {
          this.productCostSummary = {
            ...this.productCostSummary,
            [productId]: response.data.data
          };
          return response.data.data;
        }
      } catch (err) {
        console.error('Failed to fetch product cost summary:', err);
      }
    },

    async fetchProjectCostSummary(projectId) {
      try {
        const response = await bomApi.getProjectCostSummary(projectId);
        if (response.data && response.data.success) {
          this.projectCostSummary = {
            ...this.projectCostSummary,
            [projectId]: response.data.data
          };
          return response.data.data;
        }
      } catch (err) {
        console.error('Failed to fetch project cost summary:', err);
      }
    },

    selectItem(itemId) {
      this.selectedItemId = itemId;
    }
  }
});
