import { defineStore } from 'pinia';
import importApi from '../api/importApi';

const extractData = (res) => {
  if (res && res.data !== undefined) return res.data;
  return res;
};

export const useImportStore = defineStore('import', {
  state: () => ({
    batches: [],
    currentBatch: null,
    stagedRows: [],
    totalRows: 0,
    currentPage: 1,
    pageSize: 50,
    loading: false,
    tableLoading: false,
    error: null,
    commitResult: null,
    selectedYear: '2026', // 2026 active vs archive
    statusFilter: 'ALL',
    validationFilter: 'ALL',
    decisionFilter: 'ALL',
    classificationFilter: 'ALL',
    searchQuery: '',
    selectedRow: null,
    isInspectorOpen: false
  }),

  getters: {
    active2026Batches: (state) => (state.batches || []).filter(b => b && b.sourceYear === '2026'),
    archiveBatches: (state) => (state.batches || []).filter(b => b && b.sourceYear !== '2026'),
    totalPages: (state) => Math.ceil(state.totalRows / state.pageSize) || 1,
    batchSummary: (state) => {
      if (!state.currentBatch) return null;
      return {
        total: state.currentBatch.rowsCount || 0,
        valid: state.currentBatch.validRows || 0,
        warning: state.currentBatch.warningRows || 0,
        error: state.currentBatch.errorRows || 0,
        approved: state.currentBatch.approvedRows || 0,
        rejected: state.currentBatch.rejectedRows || 0,
        imported: state.currentBatch.importedRows || 0
      };
    }
  },

  actions: {
    async fetchBatches() {
      this.loading = true;
      this.error = null;
      try {
        const res = await importApi.getBatches({
          year: this.selectedYear,
          status: this.statusFilter
        });
        const data = extractData(res);
        this.batches = Array.isArray(data) ? data : (data?.data || []);
      } catch (err) {
        this.error = err.message || 'Failed to fetch batches';
      } finally {
        this.loading = false;
      }
    },

    async fetchBatchById(id) {
      this.loading = true;
      this.error = null;
      try {
        const res = await importApi.getBatchById(id);
        this.currentBatch = extractData(res);
        return this.currentBatch;
      } catch (err) {
        this.error = err.message || 'Failed to fetch batch details';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createBatch(payload) {
      this.loading = true;
      try {
        const res = await importApi.createBatch(payload);
        const batch = extractData(res);
        if (batch) {
          this.batches.unshift(batch);
          this.currentBatch = batch;
        }
        return batch;
      } catch (err) {
        this.error = err.message || 'Failed to create batch';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async uploadFiles(batchId, formData) {
      this.loading = true;
      try {
        const res = await importApi.uploadBatchFiles(batchId, formData);
        await this.fetchBatchById(batchId);
        return extractData(res);
      } catch (err) {
        this.error = err.message || 'Failed to upload files';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async parseBatch(batchId, selectedSheetsByFile, headerRowOverrides = {}) {
      this.loading = true;
      try {
        const res = await importApi.parseBatch(batchId, selectedSheetsByFile, headerRowOverrides);
        this.currentBatch = extractData(res);
        return this.currentBatch;
      } catch (err) {
        this.error = err.message || 'Failed to parse batch';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async applyMapping(batchId, mapping) {
      this.loading = true;
      try {
        const res = await importApi.applyMapping(batchId, mapping);
        this.currentBatch = extractData(res);
        return this.currentBatch;
      } catch (err) {
        this.error = err.message || 'Failed to apply mapping';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async runAIClassification(batchId) {
      this.loading = true;
      try {
        const res = await importApi.runAIClassification(batchId);
        this.currentBatch = extractData(res);
        await this.fetchStagedRows(batchId);
        return this.currentBatch;
      } catch (err) {
        this.error = err.message || 'Failed to run AI classification';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async generateMissingCodes(batchId = this.currentBatch?.id, payload = {}) {
      if (!batchId) return;
      this.loading = true;
      this.error = null;
      try {
        const res = await importApi.generateMissingCodes(batchId, payload);
        const data = extractData(res);
        if (data) {
          this.currentBatch = data;
        }
        await this.fetchStagedRows(batchId);
        return data;
      } catch (err) {
        this.error = err.message || 'Failed to generate missing codes';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchStagedRows(batchId = this.currentBatch?.id) {
      if (!batchId) return;
      this.tableLoading = true;
      try {
        const res = await importApi.getStagedRows(batchId, {
          status: this.validationFilter,
          decision: this.decisionFilter,
          classification: this.classificationFilter,
          search: this.searchQuery,
          page: this.currentPage,
          limit: this.pageSize
        });
        const data = extractData(res);
        this.stagedRows = Array.isArray(data) ? data : (data?.data || []);
        this.totalRows = res.total || (Array.isArray(data) ? data.length : data?.total) || 0;
      } catch (err) {
        this.error = err.message || 'Failed to fetch staged rows';
      } finally {
        this.tableLoading = false;
      }
    },

    async triageRow(rowId, payload) {
      try {
        const res = await importApi.triageRow(rowId, payload);
        const updatedRow = extractData(res);
        const idx = this.stagedRows.findIndex(r => r.id === rowId);
        if (idx !== -1 && updatedRow) {
          this.stagedRows.splice(idx, 1, updatedRow);
        }
        if (this.selectedRow && this.selectedRow.id === rowId && updatedRow) {
          this.selectedRow = updatedRow;
        }
        // Update batch counts
        if (this.currentBatch) {
          await this.fetchBatchById(this.currentBatch.id);
        }
        return updatedRow;
      } catch (err) {
        this.error = err.message || 'Failed to triage row';
        throw err;
      }
    },

    async approveBatch(batchId) {
      this.loading = true;
      try {
        const res = await importApi.approveBatch(batchId);
        this.currentBatch = extractData(res);
        await this.fetchStagedRows(batchId);
        return this.currentBatch;
      } catch (err) {
        this.error = err.message || 'Failed to approve batch';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async commitBatch(batchId, payload = {}) {
      this.loading = true;
      try {
        const res = await importApi.commitBatch(batchId, payload);
        this.commitResult = extractData(res);
        await this.fetchBatchById(batchId);
        await this.fetchStagedRows(batchId);
        return this.commitResult;
      } catch (err) {
        this.error = err.message || 'Failed to commit batch';
        throw err;
      } finally {
        this.loading = false;
      }
    },

    openInspector(row) {
      this.selectedRow = row;
      this.isInspectorOpen = true;
    },

    closeInspector() {
      this.isInspectorOpen = false;
      this.selectedRow = null;
    },

    async deleteBatch(id) {
      try {
        await importApi.deleteBatch(id);
      } catch (e) {
        console.warn(`[ImportStore] Delete batch ${id} API error/fallback:`, e);
      }
      this.batches = this.batches.filter(b => b.id !== id && b.batchCode !== id);
      if (this.currentBatch && (this.currentBatch.id === id || this.currentBatch.batchCode === id)) {
        this.currentBatch = null;
      }
    },

    async deleteMultipleBatches(ids) {
      if (!ids || ids.length === 0) return;
      this.loading = true;
      try {
        for (const id of ids) {
          try {
            await importApi.deleteBatch(id);
          } catch (e) {
            console.warn(`Failed to delete batch ${id}:`, e);
          }
        }
        const idSet = new Set(ids);
        this.batches = this.batches.filter(b => !idSet.has(b.id) && !idSet.has(b.batchCode));
      } finally {
        this.loading = false;
      }
    },

    async deleteAllBatches() {
      this.loading = true;
      try {
        const allIds = this.batches.map(b => b.id || b.batchCode);
        for (const id of allIds) {
          try {
            await importApi.deleteBatch(id);
          } catch (e) {
            console.warn(`Failed to delete batch ${id}:`, e);
          }
        }
        this.batches = [];
        this.currentBatch = null;
        this.stagedRows = [];
        this.totalRows = 0;
      } finally {
        this.loading = false;
      }
    }
  }
});

export default useImportStore;
