import { defineStore } from 'pinia';
import versionsApi from '@/api/versions';
import revisionsApi from '@/api/revisions';

export const useVersionStore = defineStore('versionStore', {
  state: () => ({
    // Reactive Product Versions from REST API
    versions: [],
    // Reactive Hardware Revisions from REST API
    revisions: [],
    // Lifecycle tree data per product
    productLifecycle: {},
    // Active comparison state
    comparisonResult: null,
    // Version & Revision Audit History
    versionHistory: [],
    // Loading and error states
    isLoading: false,
    error: null
  }),

  getters: {
    totalVersionsCount: (state) => state.versions.length,
    activeVersionsCount: (state) => state.versions.filter(v => v.status === 'Active').length,
    draftVersionsCount: (state) => state.versions.filter(v => v.status === 'Draft' || v.status === 'In Review').length,
    deprecatedVersionsCount: (state) => state.versions.filter(v => v.status === 'Deprecated' || v.status === 'Archived').length,
    totalRevisionsCount: (state) => state.revisions.length,
    releasedRevisionsCount: (state) => state.revisions.filter(r => r.status === 'Released').length,

    // Grouped by Product for the Product Tree Hierarchy View
    versionsByProduct: (state) => {
      const groups = {};
      state.versions.forEach(v => {
        const pKey = v.productCode || v.productId || 'UNKNOWN';
        if (!groups[pKey]) {
          groups[pKey] = {
            productCode: v.productCode || pKey,
            productName: v.productName || pKey,
            versions: []
          };
        }
        groups[pKey].versions.push(v);
      });
      return Object.values(groups);
    }
  },

  actions: {
    // Normalization helper for documents and changesList
    normalizeVersion(v) {
      if (!v) return null;
      let docs = v.documents;
      if (typeof docs === 'string' && docs.trim() !== '') {
        try {
          docs = JSON.parse(docs);
        } catch (e) {
          docs = [];
        }
      }
      return {
        ...v,
        documents: Array.isArray(docs) ? docs : []
      };
    },

    normalizeRevision(r) {
      if (!r) return null;
      let changes = r.changesList;
      if (typeof changes === 'string' && changes.trim() !== '') {
        try {
          changes = JSON.parse(changes);
        } catch (e) {
          changes = [r.changesList];
        }
      }
      let docs = r.documents;
      if (typeof docs === 'string' && docs.trim() !== '') {
        try {
          docs = JSON.parse(docs);
        } catch (e) {
          docs = [];
        }
      }
      return {
        ...r,
        versionId: r.productVersionId || r.versionId,
        changesList: Array.isArray(changes) ? changes : [],
        documents: Array.isArray(docs) ? docs : []
      };
    },

    // Fetch all versions with optional filtering
    async fetchVersions(params = { limit: 100 }) {
      this.isLoading = true;
      this.error = null;
      try {
        const res = await versionsApi.getAll(params);
        if (res.data && res.data.data) {
          const list = Array.isArray(res.data.data) ? res.data.data : [];
          this.versions = list.map(v => this.normalizeVersion(v));
        }
        return this.versions;
      } catch (err) {
        this.error = err.message || 'Failed to fetch product versions';
        console.error('[versionStore] fetchVersions error:', err);
        return this.versions;
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch all hardware revisions
    async fetchRevisions(params = { limit: 100 }) {
      this.isLoading = true;
      this.error = null;
      try {
        const res = await revisionsApi.getAll(params);
        if (res.data && res.data.data) {
          const list = Array.isArray(res.data.data) ? res.data.data : [];
          this.revisions = list.map(r => this.normalizeRevision(r));
        }
        return this.revisions;
      } catch (err) {
        this.error = err.message || 'Failed to fetch hardware revisions';
        console.error('[versionStore] fetchRevisions error:', err);
        return this.revisions;
      } finally {
        this.isLoading = false;
      }
    },

    // Fetch Product Lifecycle tree
    async fetchProductLifecycle(productId) {
      this.isLoading = true;
      try {
        const res = await versionsApi.getLifecycle(productId);
        if (res.data && res.data.data) {
          this.productLifecycle[productId] = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        console.error('[versionStore] fetchProductLifecycle error:', err);
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    // Compare two versions
    async compareVersions(baseId, targetId) {
      this.isLoading = true;
      try {
        const res = await versionsApi.compare(baseId, targetId);
        if (res.data && res.data.data) {
          this.comparisonResult = res.data.data;
          return res.data.data;
        }
      } catch (err) {
        console.error('[versionStore] compareVersions error:', err);
        throw err;
      } finally {
        this.isLoading = false;
      }
      return null;
    },

    // Synchronous lookups from state (refreshed automatically on fetch)
    getVersionById(versionId) {
      return this.versions.find(v => v.id === versionId || v.versionNumber === versionId);
    },

    getVersionsForProduct(productCodeOrId) {
      return this.versions.filter(v => v.productCode === productCodeOrId || v.productId === productCodeOrId);
    },

    getRevisionById(revisionId) {
      return this.revisions.find(r => r.id === revisionId);
    },

    getRevisionsForVersion(versionId) {
      const ver = this.getVersionById(versionId);
      const vId = ver ? ver.id : versionId;
      const vNum = ver ? ver.versionNumber : versionId;
      return this.revisions.filter(r => (r.versionId === vId || r.productVersionId === vId || r.versionNumber === vNum));
    },

    // Save (Create or Update) Product Version
    async saveVersion(versionData) {
      this.isLoading = true;
      try {
        const payload = {
          ...versionData,
          documents: typeof versionData.documents === 'object' ? JSON.stringify(versionData.documents) : versionData.documents
        };

        let result;
        if (versionData.id && this.versions.some(v => v.id === versionData.id)) {
          // Update existing
          const res = await versionsApi.update(versionData.id, payload);
          result = this.normalizeVersion(res.data.data);
          const index = this.versions.findIndex(v => v.id === versionData.id);
          if (index !== -1) {
            this.versions.splice(index, 1, result);
          } else {
            this.versions.unshift(result);
          }
          this.logHistory({
            actor: 'Current User',
            action: 'Version Updated',
            targetType: 'Product Version',
            targetId: result.id,
            targetName: `${result.productName} v${result.versionNumber}`,
            version: `v${result.versionNumber}`,
            revision: result.currentRevision || 'REV-A',
            badgeColor: 'blue',
            details: `Updated version parameters for ${result.productCode}.`
          });
        } else {
          // Create new
          const res = await versionsApi.create(payload);
          result = this.normalizeVersion(res.data.data);
          this.versions.unshift(result);
          // Refetch revisions as backend automatically creates initial REV-A
          await this.fetchRevisions();

          this.logHistory({
            actor: 'Current User',
            action: 'Version Created',
            targetType: 'Product Version',
            targetId: result.id,
            targetName: `${result.productName} v${result.versionNumber}`,
            version: `v${result.versionNumber}`,
            revision: 'REV-A',
            badgeColor: 'emerald',
            details: `Created new version ${result.versionNumber} for ${result.productCode}.`
          });
        }
        return result;
      } catch (err) {
        console.error('[versionStore] saveVersion error:', err);
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    // Delete Product Version
    async deleteVersion(versionId) {
      this.isLoading = true;
      try {
        await versionsApi.delete(versionId);
        const idx = this.versions.findIndex(v => v.id === versionId);
        if (idx !== -1) {
          const deleted = this.versions.splice(idx, 1)[0];
          this.logHistory({
            actor: 'Current User',
            action: 'Version Deleted',
            targetType: 'Product Version',
            targetId: versionId,
            targetName: `${deleted.productName} v${deleted.versionNumber}`,
            version: `v${deleted.versionNumber}`,
            revision: deleted.currentRevision,
            badgeColor: 'rose',
            details: `Deleted version ${deleted.versionNumber} from ${deleted.productCode}.`
          });
        }
        return true;
      } catch (err) {
        console.error('[versionStore] deleteVersion error:', err);
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    // Save (Create or Update) Hardware Revision
    async saveRevision(revisionData) {
      this.isLoading = true;
      try {
        const payload = {
          ...revisionData,
          productVersionId: revisionData.productVersionId || revisionData.versionId,
          changesList: typeof revisionData.changesList === 'object' ? JSON.stringify(revisionData.changesList) : revisionData.changesList,
          documents: typeof revisionData.documents === 'object' ? JSON.stringify(revisionData.documents) : revisionData.documents
        };

        let result;
        if (revisionData.id && this.revisions.some(r => r.id === revisionData.id)) {
          // Update existing
          const res = await revisionsApi.update(revisionData.id, payload);
          result = this.normalizeRevision(res.data.data);
          const index = this.revisions.findIndex(r => r.id === revisionData.id);
          if (index !== -1) {
            this.revisions.splice(index, 1, result);
          } else {
            this.revisions.unshift(result);
          }
          this.logHistory({
            actor: 'Current User',
            action: 'Revision Updated',
            targetType: 'Hardware Revision',
            targetId: result.id,
            targetName: `${result.productCode} ${result.code}`,
            version: `v${result.versionNumber}`,
            revision: result.code,
            badgeColor: 'blue',
            details: `Updated revision details for ${result.code}.`
          });
        } else {
          // Create new
          const res = await revisionsApi.create(payload);
          result = this.normalizeRevision(res.data.data);
          this.revisions.unshift(result);

          // Update parent version locally
          const ver = this.getVersionById(result.versionId);
          if (ver) {
            ver.currentRevision = result.code;
            ver.revisionsCount = (ver.revisionsCount || 0) + 1;
          }

          this.logHistory({
            actor: 'Current User',
            action: 'Revision Created',
            targetType: 'Hardware Revision',
            targetId: result.id,
            targetName: `${result.productCode} ${result.code}`,
            version: `v${result.versionNumber}`,
            revision: result.code,
            badgeColor: 'emerald',
            details: `Created new hardware revision ${result.code} for version ${result.versionNumber}.`
          });
        }
        return result;
      } catch (err) {
        console.error('[versionStore] saveRevision error:', err);
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    // Delete Hardware Revision
    async deleteRevision(revisionId) {
      this.isLoading = true;
      try {
        await revisionsApi.delete(revisionId);
        const idx = this.revisions.findIndex(r => r.id === revisionId);
        if (idx !== -1) {
          const deleted = this.revisions.splice(idx, 1)[0];
          this.logHistory({
            actor: 'Current User',
            action: 'Revision Deleted',
            targetType: 'Hardware Revision',
            targetId: revisionId,
            targetName: `${deleted.productCode} ${deleted.code}`,
            version: `v${deleted.versionNumber}`,
            revision: deleted.code,
            badgeColor: 'rose',
            details: `Deleted hardware revision ${deleted.code}.`
          });
        }
        return true;
      } catch (err) {
        console.error('[versionStore] deleteRevision error:', err);
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    logHistory(item) {
      this.versionHistory.unshift({
        id: `VH-00${this.versionHistory.length + 1}`,
        timestamp: 'Just now',
        date: new Date().toISOString().slice(0, 10),
        ...item
      });
    }
  }
});
