<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <router-link to="/data-import" class="text-xs text-brand-600 hover:text-brand-800 font-semibold flex items-center gap-1">
            <AppIcon name="chevron-left" size="xs" /> Batches
          </router-link>
          <span class="text-slate-400">•</span>
          <router-link :to="`/data-import/wizard/${batch?.id}`" class="text-xs text-slate-500 hover:text-slate-700 font-semibold">
            Wizard
          </router-link>
          <span class="text-slate-400">•</span>
          <span class="font-mono text-xs text-brand-700 font-bold">{{ batch?.batchCode || 'Loading...' }}</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">Staging Review & Triage Workspace</h1>
      </div>

      <div class="flex items-center gap-2">
        <span v-if="batch?.status === 'IMPORTED'" class="px-3 py-1.5 rounded-lg text-xs font-mono font-bold bg-emerald-100 text-emerald-800 border border-emerald-300 flex items-center gap-1.5 mr-1">
          <AppIcon name="check-circle" size="xs" />
          <span>Import Committed</span>
        </span>

        <button
          @click="handleGenerateCodes"
          class="btn bg-indigo-50 hover:bg-indigo-100 text-indigo-700 border border-indigo-200 text-xs"
          :disabled="store.loading"
          title="Auto-Generate Standardized MPN / Item Codes and Live Sync with Master Data"
        >
          <AppIcon name="sparkles" size="xs" class="mr-1 text-indigo-600" />
          <span>Auto-Generate Codes</span>
        </button>

        <button
          @click="handleBulkApprove"
          class="btn btn-secondary text-xs"
          :disabled="store.loading"
        >
          <AppIcon name="check" size="xs" class="mr-1 text-emerald-600" />
          <span>Approve All Valid Rows</span>
        </button>

        <button
          @click="showCommitModal = true"
          class="btn btn-primary text-xs"
          :disabled="store.loading || (batch?.approvedRows === 0)"
        >
          <AppIcon name="database" size="xs" class="mr-1.5" />
          <span>{{ batch?.status === 'IMPORTED' ? 'Re-Sync to Master (' + (batch?.approvedRows || 0) + ')' : 'Commit to Master (' + (batch?.approvedRows || 0) + ')' }}</span>
        </button>
      </div>
    </div>

    <!-- Batch Summary Strip -->
    <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-6 gap-3">
      <div class="metric-card">
        <div class="text-3xs font-bold text-slate-500 uppercase">Staged Rows</div>
        <div class="text-lg font-black text-slate-900 font-mono mt-1">{{ batch?.rowsCount || 0 }}</div>
      </div>
      <div class="metric-card">
        <div class="text-3xs font-bold text-emerald-600 uppercase">Approved</div>
        <div class="text-lg font-black text-emerald-700 font-mono mt-1">{{ batch?.approvedRows || 0 }}</div>
      </div>
      <div class="metric-card">
        <div class="text-3xs font-bold text-amber-600 uppercase">Warnings</div>
        <div class="text-lg font-black text-amber-700 font-mono mt-1">{{ batch?.warningRows || 0 }}</div>
      </div>
      <div class="metric-card">
        <div class="text-3xs font-bold text-rose-600 uppercase">Errors</div>
        <div class="text-lg font-black text-rose-700 font-mono mt-1">{{ batch?.errorRows || 0 }}</div>
      </div>
      <div class="metric-card">
        <div class="text-3xs font-bold text-slate-500 uppercase">Rejected</div>
        <div class="text-lg font-black text-slate-700 font-mono mt-1">{{ batch?.rejectedRows || 0 }}</div>
      </div>
      <div class="metric-card border-brand-200 bg-brand-50/20">
        <div class="text-3xs font-bold text-brand-700 uppercase">Imported</div>
        <div class="text-lg font-black text-brand-900 font-mono mt-1">{{ batch?.importedRows || 0 }}</div>
      </div>
    </div>

    <!-- Master Parent Project Definition Banner -->
    <div class="bg-gradient-to-r from-blue-50/80 via-indigo-50/40 to-slate-50 border border-blue-200/80 rounded-xl p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-2xs">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs">
          <AppIcon name="briefcase" size="sm" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-3xs font-bold uppercase tracking-wider text-blue-700 font-mono">Master Parent Project Definition</span>
            <span class="px-2 py-0.2 rounded text-4xs font-bold bg-blue-100 text-blue-800 border border-blue-200">PROJECT Mode</span>
          </div>
          <p class="text-2xs text-slate-500 mt-0.5">Semua item sensor, komponen & jasa akan masuk ke katalog Master dan otomatis dirangkai ke dalam Struktur Proyek Induk ini saat di-commit.</p>
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0 sm:max-w-md w-full sm:w-auto">
        <label class="text-3xs font-semibold text-slate-600 shrink-0">Nama Proyek Induk:</label>
        <input
          type="text"
          v-model="parentProjectName"
          class="form-input text-xs font-bold py-1.5 px-3 bg-white border-blue-300 focus:border-blue-500 rounded-lg w-full sm:w-80 shadow-2xs"
          placeholder="e.g. RAB Maintenance SPARING - Lipo R1"
        />
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-3 bg-white p-3.5 rounded-xl border border-slate-200 shadow-2xs">
      <div class="flex flex-wrap items-center gap-2">
        <input
          type="text"
          v-model="store.searchQuery"
          @input="handleSearch"
          placeholder="Search staged names, codes, categories..."
          class="form-input text-xs w-64"
        />

        <select v-model="store.validationFilter" @change="fetchRows" class="form-select text-xs py-1.5 w-auto">
          <option value="ALL">All Validations</option>
          <option value="VALID">Valid Only</option>
          <option value="WARNING">Warnings</option>
          <option value="ERROR">Errors Only</option>
        </select>

        <select v-model="store.decisionFilter" @change="fetchRows" class="form-select text-xs py-1.5 w-auto">
          <option value="ALL">All Decisions</option>
          <option value="APPROVED">Approved</option>
          <option value="PENDING">Pending</option>
          <option value="REJECTED">Rejected</option>
          <option value="EDITED">Edited</option>
        </select>

        <select v-model="store.classificationFilter" @change="fetchRows" class="form-select text-xs py-1.5 w-auto">
          <option value="ALL">All Types</option>
          <option value="PRODUCT">PRODUCT</option>
          <option value="COMPONENT">COMPONENT</option>
          <option value="SERVICE">SERVICE</option>
          <option value="ACCESSORY">ACCESSORY</option>
        </select>
      </div>

      <div class="text-2xs text-slate-500 font-mono">
        Showing <strong class="text-slate-800">{{ store.stagedRows.length }}</strong> of {{ store.totalRows }} records
      </div>
    </div>

    <!-- High-Density Triage Table -->
    <div class="app-card overflow-hidden">
      <div v-if="store.tableLoading" class="p-12 text-center text-slate-400 text-xs">
        <AppIcon name="database" size="md" class="mx-auto text-slate-300 animate-pulse mb-2" />
        Loading staged records...
      </div>

      <div v-else-if="store.stagedRows.length === 0" class="p-12 text-center text-slate-400 text-xs">
        No staged records match the active filter criteria.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs divide-y divide-slate-100">
          <thead>
            <tr class="bg-slate-50/90 text-slate-500 uppercase tracking-wider text-3xs font-bold select-none border-b border-slate-200/80">
              <th class="py-2.5 px-2.5 w-12 text-center">Row</th>
              <th class="py-2.5 px-2.5 min-w-[160px]">Normalized Name</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap min-w-[110px]">Code / MPN</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24">Item Class</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-20">Type</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Category</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28 text-right">Cost / Price</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-20 text-center">Validation</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24">AI Co-Pilot</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-20 text-center">Decision</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr
              v-for="r in store.stagedRows"
              :key="r.id"
              @click="store.openInspector(r)"
              class="cursor-pointer hover:bg-slate-50/80 transition-colors"
              :class="{ 'bg-brand-50/40': store.selectedRow?.id === r.id }"
            >
              <td class="py-2 px-2.5 text-center font-mono text-3xs text-slate-400">
                #{{ r.rowNumber }}
              </td>
              <td class="py-2 px-2.5">
                <div class="font-bold text-slate-900 truncate max-w-[180px] text-xs" :title="r.normalizedName">
                  {{ r.normalizedName || 'Untitled' }}
                </div>
                <div class="text-3xs text-slate-400 font-mono truncate max-w-[180px]">
                  Src: {{ r.fileName }}
                </div>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap font-mono font-bold text-xs text-slate-800">
                {{ r.normalizedCode || '-' }}
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span class="px-1.5 py-0.2 rounded text-3xs font-mono font-bold" :class="getClassificationBadge(r.itemClassification)">
                  {{ r.itemClassification }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap" @click.stop>
                <select
                  v-if="r.itemClassification === 'PRODUCT'"
                  v-model="r.productType"
                  @change="quickUpdateType(r, r.productType)"
                  class="form-select text-3xs font-bold py-0.5 px-1.5 rounded bg-slate-50 border border-slate-200 text-slate-800 focus:ring-brand-500 cursor-pointer"
                >
                  <option value="TRADING">TRADING</option>
                  <option value="PROJECT">PROJECT</option>
                  <option value="MANUFACTURE">MANUFACTURE</option>
                  <option value="SERVICE">SERVICE</option>
                </select>
                <span v-else class="text-3xs text-slate-300">-</span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 truncate block max-w-[110px]" :title="r.normalizedCategory">
                  {{ r.normalizedCategory || 'General' }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap text-right font-mono text-xs">
                <div class="font-bold text-slate-900">
                  {{ formatCurrency(r.sellingPrice > 0 ? r.sellingPrice : r.unitCost, r.currency) }}
                </div>
                <div v-if="r.sellingPrice > 0 && r.unitCost > 0" class="text-3xs text-slate-400 font-normal">
                  Cost: {{ formatCurrency(r.unitCost, r.currency) }}
                </div>
                <div v-else-if="r.sellingPrice > 0" class="text-3xs text-slate-400 font-normal">
                  Harga Jual
                </div>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap text-center">
                <span
                  class="px-1.5 py-0.5 rounded text-3xs font-mono font-bold"
                  :class="getValidationBadge(r.validationStatus)"
                >
                  {{ r.validationStatus }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <AiSuggestionBadge :confidence="r.aiConfidence" :label="r.aiSuggestedType" />
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap text-center">
                <span
                  class="px-1.5 py-0.5 rounded text-3xs font-mono font-bold uppercase"
                  :class="getDecisionBadge(r.reviewDecision)"
                >
                  {{ r.reviewDecision }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="quickDecision(r, 'APPROVED')"
                    class="p-1 rounded hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 transition-colors"
                    title="Quick Approve"
                  >
                    <AppIcon name="check" size="xs" />
                  </button>
                  <button
                    @click="quickDecision(r, 'REJECTED')"
                    class="p-1 rounded hover:bg-rose-50 text-slate-400 hover:text-rose-600 transition-colors"
                    title="Quick Reject"
                  >
                    <AppIcon name="x" size="xs" />
                  </button>
                  <button
                    @click="store.openInspector(r)"
                    class="btn-icon"
                    title="Open Inspector"
                  >
                    <AppIcon name="eye" size="xs" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="store.totalRows > 0" class="px-4 py-3 border-t border-slate-100 bg-slate-50/70 flex items-center justify-between text-xs text-slate-600">
        <span class="text-3xs text-slate-500 font-mono">
          Page <strong class="text-slate-800">{{ store.currentPage }}</strong> of {{ store.totalPages }}
        </span>
        <div class="flex items-center gap-1">
          <button
            :disabled="store.currentPage === 1"
            @click="prevPage"
            class="btn btn-secondary text-2xs py-1 px-2.5"
          >
            Previous
          </button>
          <button
            :disabled="store.currentPage >= store.totalPages"
            @click="nextPage"
            class="btn btn-secondary text-2xs py-1 px-2.5"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <!-- Row Inspector Drawer -->
    <RowInspectorDrawer
      :isOpen="store.isInspectorOpen"
      :row="store.selectedRow"
      @close="store.closeInspector"
      @triage="handleInspectorTriage"
    />

    <!-- Formal Commit Modal -->
    <div v-if="showCommitModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-lg p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">Authorize & Commit Master Data Migration</h3>
          <button @click="showCommitModal = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <div class="space-y-3 text-xs text-slate-600">
          <p>
            You are about to commit <strong class="text-slate-900 font-bold font-mono">{{ batch?.approvedRows }} approved rows</strong> from batch
            <strong class="font-mono text-brand-700 font-bold">{{ batch?.batchCode }}</strong> into production Master tables.
          </p>

          <div class="p-3 bg-slate-50 rounded-xl border border-slate-200 text-3xs font-mono space-y-1">
            <div>• Database Transactions: <strong>Enabled (GORM Atomic)</strong></div>
            <div>• Idempotency Check: <strong>Enforced</strong></div>
            <div>• Lineage & Provenance Envelope: <strong>Preserved</strong></div>
            <div>• Unapproved Error Rows ({{ batch?.errorRows }}): <strong>Skipped & Staged</strong></div>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
          <button @click="showCommitModal = false" class="btn btn-secondary text-xs">Cancel</button>
          <button @click="executeCommit" class="btn btn-primary text-xs" :disabled="committing">
            <span v-if="committing">Committing Master Data...</span>
            <span v-else>Authorize & Commit</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useImportStore } from '../../stores/importStore';
import AppIcon from '../../components/common/AppIcon.vue';
import AiSuggestionBadge from '../../components/import/AiSuggestionBadge.vue';
import RowInspectorDrawer from '../../components/import/RowInspectorDrawer.vue';

export default {
  name: 'ImportReviewWorkspace',
  components: { AppIcon, AiSuggestionBadge, RowInspectorDrawer },
  data() {
    return {
      store: useImportStore(),
      showCommitModal: false,
      committing: false,
      parentProjectName: ''
    };
  },
  computed: {
    batch() {
      return this.store.currentBatch;
    }
  },
  watch: {
    batch: {
      immediate: true,
      handler(val) {
        if (val && !this.parentProjectName) {
          this.initParentProjectName();
        }
      }
    }
  },
  async mounted() {
    const id = this.$route.params.id;
    if (id) {
      await this.store.fetchBatchById(id);
      await this.store.fetchStagedRows(id);
      this.initParentProjectName();
    }
  },
  methods: {
    initParentProjectName() {
      const raw = this.batch?.files?.[0]?.originalFileName || this.batch?.batchName || '';
      if (!raw) return;
      let cleaned = raw.replace(/\.[^/.]+$/, '');
      cleaned = cleaned.replace(/^[A-Z0-9_-]+\s*[-:]\s*/i, '');
      this.parentProjectName = cleaned || 'RAB Maintenance SPARING - Lipo R1';
    },
    handleSearch() {
      this.store.currentPage = 1;
      this.store.fetchStagedRows();
    },
    fetchRows() {
      this.store.currentPage = 1;
      this.store.fetchStagedRows();
    },
    prevPage() {
      if (this.store.currentPage > 1) {
        this.store.currentPage--;
        this.store.fetchStagedRows();
      }
    },
    nextPage() {
      if (this.store.currentPage < this.store.totalPages) {
        this.store.currentPage++;
        this.store.fetchStagedRows();
      }
    },
    getClassificationBadge(cls) {
      switch (cls) {
        case 'PRODUCT': return 'bg-brand-50 text-brand-700 border border-brand-200';
        case 'COMPONENT': return 'bg-indigo-50 text-indigo-700 border border-indigo-200';
        case 'SERVICE': return 'bg-amber-50 text-amber-700 border border-amber-200';
        default: return 'bg-slate-100 text-slate-700';
      }
    },
    getValidationBadge(v) {
      switch (v) {
        case 'VALID': return 'bg-emerald-50 text-emerald-700 border border-emerald-200';
        case 'WARNING': return 'bg-amber-50 text-amber-700 border border-amber-200';
        case 'ERROR': return 'bg-rose-50 text-rose-700 border border-rose-200';
        default: return 'bg-slate-100 text-slate-500';
      }
    },
    getDecisionBadge(d) {
      switch (d) {
        case 'APPROVED': return 'bg-emerald-100 text-emerald-800';
        case 'REJECTED': return 'bg-rose-100 text-rose-800';
        case 'EDITED': return 'bg-blue-100 text-blue-800';
        default: return 'bg-slate-100 text-slate-600';
      }
    },
    formatCurrency(val, curr = 'IDR') {
      if (!val) return '0 ' + curr;
      return new Intl.NumberFormat('en-US').format(val) + ' ' + curr;
    },
    async quickDecision(row, decision) {
      try {
        await this.store.triageRow(row.id, {
          itemClassification: row.itemClassification,
          productType: row.productType,
          normalizedName: row.normalizedName,
          normalizedCode: row.normalizedCode,
          normalizedCategory: row.normalizedCategory,
          description: row.description,
          unitCost: row.unitCost,
          sellingPrice: row.sellingPrice,
          currency: row.currency,
          decision: decision,
          reviewNotes: row.reviewNotes
        });
      } catch (err) {
        alert('Triage error: ' + err.message);
      }
    },
    async quickUpdateType(row, type) {
      try {
        await this.store.triageRow(row.id, {
          itemClassification: row.itemClassification,
          productType: type,
          normalizedName: row.normalizedName,
          normalizedCode: row.normalizedCode,
          normalizedCategory: row.normalizedCategory,
          description: row.description,
          unitCost: row.unitCost,
          sellingPrice: row.sellingPrice,
          currency: row.currency,
          decision: row.reviewDecision,
          reviewNotes: row.reviewNotes
        });
      } catch (err) {
        alert('Failed to update product type: ' + err.message);
      }
    },
    async handleInspectorTriage(payload) {
      if (!this.store.selectedRow) return;
      try {
        await this.store.triageRow(this.store.selectedRow.id, payload);
        this.store.closeInspector();
      } catch (err) {
        alert('Failed to save triage changes: ' + err.message);
      }
    },
    async handleGenerateCodes() {
      try {
        await this.store.generateMissingCodes(this.batch.id, {
          parentProjectName: this.parentProjectName
        });
      } catch (err) {
        alert('Failed to generate codes: ' + err.message);
      }
    },
    async handleBulkApprove() {
      if (!confirm('Approve all valid and warned rows for Master Data migration?')) return;
      try {
        await this.store.approveBatch(this.batch.id);
      } catch (err) {
        alert('Bulk approve failed: ' + err.message);
      }
    },
    async executeCommit() {
      this.committing = true;
      try {
        const res = await this.store.commitBatch(this.batch.id, {
          parentProjectName: this.parentProjectName
        });
        this.showCommitModal = false;
        alert(`Migration Successful!\nCommitted: ${res.totalCommitted} items (Products: ${res.createdProducts}, Components: ${res.createdComponents}, Services: ${res.createdServices})${res.parentProjectName ? '\nParent Project: ' + res.parentProjectName : ''}`);
      } catch (err) {
        alert('Commit failed: ' + err.message);
      } finally {
        this.committing = false;
      }
    }
  }
};
</script>
