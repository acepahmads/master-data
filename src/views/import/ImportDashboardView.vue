<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <span class="px-2.5 py-0.5 rounded text-3xs font-mono font-bold uppercase tracking-wider bg-brand-50 text-brand-700 border border-brand-200">
            Cross-Cutting Platform Capability
          </span>
          <span class="text-slate-400">•</span>
          <span class="text-xs font-semibold text-slate-500 font-mono">Rest API v1.0.0</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">Product Data Import & Migration</h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Multi-file Excel ingestion, heuristic normalization, AI classification, and master data migration with complete audit provenance.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button @click="showNewBatchModal = true" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>New Import Batch</span>
        </button>
      </div>
    </div>

    <!-- Executive KPI Strip -->
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-3">
      <div class="metric-card">
        <div class="text-3xs font-bold text-slate-500 uppercase tracking-wider">Total Batches</div>
        <div class="text-xl font-black text-slate-900 font-mono mt-1">{{ totalBatchesCount }}</div>
        <div class="text-3xs text-slate-400 mt-0.5">2016–2026 Ingestion Sessions</div>
      </div>

      <div class="metric-card border-brand-200 bg-brand-50/30">
        <div class="text-3xs font-bold text-brand-700 uppercase tracking-wider">Active 2026 Scope</div>
        <div class="text-xl font-black text-brand-900 font-mono mt-1">{{ active2026RowsCount }} Rows</div>
        <div class="text-3xs text-brand-600 mt-0.5">{{ store.active2026Batches.length }} Active Batches</div>
      </div>

      <div class="metric-card">
        <div class="text-3xs font-bold text-emerald-600 uppercase tracking-wider">Validated Quality</div>
        <div class="text-xl font-black text-emerald-700 font-mono mt-1">{{ validatedRate }}%</div>
        <div class="text-3xs text-slate-400 mt-0.5">Rule Compliance Average</div>
      </div>

      <div class="metric-card">
        <div class="text-3xs font-bold text-amber-600 uppercase tracking-wider">Triage Queue</div>
        <div class="text-xl font-black text-amber-700 font-mono mt-1">{{ totalIssuesCount }}</div>
        <div class="text-3xs text-slate-400 mt-0.5">Warnings & Error Rows</div>
      </div>

      <div class="metric-card border-purple-200 bg-purple-50/20">
        <div class="text-3xs font-bold text-purple-700 uppercase tracking-wider">AI Co-Pilot</div>
        <div class="text-xl font-black text-purple-900 font-mono mt-1">Advisory</div>
        <div class="text-3xs text-purple-600 mt-0.5">Human Approval Mandatory</div>
      </div>
    </div>

    <!-- Year Scope & Filter Bar -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-3.5 rounded-xl border border-slate-200 shadow-2xs">
      <div class="flex items-center gap-1.5 p-1 bg-slate-100 rounded-lg">
        <button
          @click="setYear('2026')"
          class="px-3 py-1 rounded text-xs font-bold transition-all"
          :class="store.selectedYear === '2026' ? 'bg-white text-brand-700 shadow-2xs' : 'text-slate-600 hover:text-slate-900'"
        >
          ● 2026 Active Scope ({{ store.active2026Batches.length }})
        </button>
        <button
          @click="setYear('ALL')"
          class="px-3 py-1 rounded text-xs font-bold transition-all"
          :class="store.selectedYear === 'ALL' ? 'bg-white text-slate-900 shadow-2xs' : 'text-slate-600 hover:text-slate-900'"
        >
          All Years (2016–2026)
        </button>
        <button
          @click="setYear('ARCHIVE')"
          class="px-3 py-1 rounded text-xs font-bold transition-all"
          :class="store.selectedYear === 'ARCHIVE' ? 'bg-white text-slate-900 shadow-2xs' : 'text-slate-600 hover:text-slate-900'"
        >
          Historical Archive ({{ store.archiveBatches.length }})
        </button>
      </div>

      <div class="flex items-center gap-3">
        <select v-model="store.statusFilter" @change="store.fetchBatches" class="form-select text-xs py-1.5 w-auto">
          <option value="ALL">All Statuses</option>
          <option value="UPLOADED">Uploaded</option>
          <option value="PARSED">Parsed</option>
          <option value="VALIDATING">Validating</option>
          <option value="READY_FOR_REVIEW">Ready for Review</option>
          <option value="APPROVED">Approved</option>
          <option value="IMPORTED">Imported</option>
        </select>
        <button @click="store.fetchBatches" class="btn btn-secondary text-xs" :disabled="store.loading">
          <AppIcon name="trending-up" size="xs" :class="{ 'animate-spin': store.loading }" />
          <span>Refresh</span>
        </button>
      </div>
    </div>

    <!-- Batches Registry Table -->
    <div class="app-card overflow-hidden">
      <div v-if="store.loading && store.batches.length === 0" class="p-12 text-center text-slate-400 text-xs">
        <AppIcon name="database" size="md" class="mx-auto text-slate-300 animate-pulse mb-2" />
        Loading import batches...
      </div>

      <div v-else-if="store.batches.length === 0" class="p-12 text-center">
        <EmptyState
          title="No Import Batches Found"
          description="Create your first import batch to ingest and normalize historical 2026 Excel datasets."
          icon="import"
        >
          <template #action>
            <button @click="showNewBatchModal = true" class="btn btn-primary">
              <AppIcon name="plus" size="xs" class="mr-1.5" />
              <span>Create 2026 Import Batch</span>
            </button>
          </template>
        </EmptyState>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs divide-y divide-slate-100">
          <thead>
            <tr class="bg-slate-50/90 text-slate-500 uppercase tracking-wider text-3xs font-bold select-none border-b border-slate-200/80">
              <th class="w-8 py-2.5 px-2 text-center">
                <input
                  type="checkbox"
                  v-model="selectAll"
                  class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer"
                />
              </th>
              <th class="py-2.5 px-3 whitespace-nowrap min-w-[130px]">Batch Code</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24">Year Scope</th>
              <th class="py-2.5 px-2.5 min-w-[180px]">Source Workbook(s)</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28 text-center">Staged Rows</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28 text-center">Validation</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Status</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Uploader</th>
              <th class="py-2.5 px-3 whitespace-nowrap w-36 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr v-for="b in store.batches" :key="b.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-2.5 px-2 text-center" @click.stop>
                <input
                  type="checkbox"
                  :value="b.id"
                  v-model="selectedBatchIds"
                  class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer"
                />
              </td>
              <td class="py-2.5 px-3 whitespace-nowrap">
                <router-link :to="`/data-import/wizard/${b.id}`" class="font-mono font-bold text-brand-600 hover:text-brand-800 text-xs">
                  {{ b.batchCode }}
                </router-link>
                <div class="text-3xs text-slate-400 font-mono mt-0.5">{{ formatTimestamp(b.createdAt) }}</div>
              </td>
              <td class="py-2.5 px-2.5 whitespace-nowrap">
                <span
                  class="px-2 py-0.5 rounded text-3xs font-mono font-bold"
                  :class="b.sourceYear === '2026' ? 'bg-brand-50 text-brand-700 border border-brand-200' : 'bg-slate-100 text-slate-600 border border-slate-200'"
                >
                  {{ b.sourceYear }} Scope
                </span>
              </td>
              <td class="py-2.5 px-2.5">
                <div class="flex items-center gap-1.5 text-xs font-medium text-slate-800">
                  <AppIcon name="file-spreadsheet" size="xs" class="text-emerald-600 shrink-0" />
                  <span class="truncate max-w-[200px]" :title="b.sourceReference || 'Workbook'">
                    {{ b.sourceReference || (b.files && b.files[0] ? b.files[0].originalFileName : 'Excel Ingestion') }}
                  </span>
                </div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5">
                  {{ b.filesCount || (b.files ? b.files.length : 0) }} file(s) attached
                </div>
              </td>
              <td class="py-2.5 px-2.5 whitespace-nowrap text-center">
                <span class="font-mono font-bold text-xs text-slate-800">{{ b.rowsCount }}</span>
                <div class="text-3xs text-slate-400">total rows</div>
              </td>
              <td class="py-2.5 px-2.5 whitespace-nowrap text-center">
                <div class="flex items-center justify-center gap-1.5 font-mono text-3xs font-bold">
                  <span class="text-emerald-600" title="Valid Rows">✓ {{ b.validRows }}</span>
                  <span class="text-slate-300">|</span>
                  <span class="text-amber-600" title="Warnings">▲ {{ b.warningRows }}</span>
                  <span class="text-slate-300">|</span>
                  <span class="text-rose-600" title="Errors">✕ {{ b.errorRows }}</span>
                </div>
              </td>
              <td class="py-2.5 px-2.5 whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-3xs font-mono font-bold uppercase" :class="getStatusBadgeClass(b.status)">
                  {{ b.status }}
                </span>
              </td>
              <td class="py-2.5 px-2.5 whitespace-nowrap">
                <div class="text-slate-700 text-xs font-medium truncate max-w-[120px]">{{ b.uploadedBy }}</div>
              </td>
              <td class="py-2.5 px-3 whitespace-nowrap text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <router-link :to="`/data-import/wizard/${b.id}`" class="btn btn-secondary text-2xs py-1 px-2">
                    <span>Wizard</span>
                  </router-link>
                  <router-link :to="`/data-import/review/${b.id}`" class="btn btn-primary text-2xs py-1 px-2">
                    <span>Review</span>
                  </router-link>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Confirm Delete Selected Dialog -->
    <ConfirmDialog
      :show="deleteSelectedModalOpen"
      title="Delete Selected Import Batches"
      :message="`Are you sure you want to delete ${selectedBatchIds.length} selected import session(s)? Staged rows and uploaded workbooks will be purged.`"
      confirmText="Delete Batches"
      @confirm="executeDeleteSelected"
      @cancel="deleteSelectedModalOpen = false"
    />

    <!-- Danger Purge / Delete All Modal (Super Admin) -->
    <BaseModal :show="deleteAllModalOpen" title="⚠️ Danger Zone: Delete All Import Batches" maxWidth="md" @close="closeDeleteAllModal">
      <div class="space-y-4 text-xs">
        <div class="p-3.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-800 space-y-2">
          <div class="flex items-center gap-2 font-bold text-rose-900 text-sm">
            <AppIcon name="alert-triangle" size="sm" class="text-rose-600 shrink-0" />
            <span>Permanent Staging Purge Warning</span>
          </div>
          <p class="leading-relaxed">
            You are about to delete <strong class="underline font-mono">{{ totalBatchesCount }} import batches</strong> ({{ active2026RowsCount }} staged rows). This will permanently purge all uploaded Excel files, staging tables, and AI classification triage queues.
          </p>
          <div class="text-3xs font-semibold uppercase tracking-wider text-rose-700 bg-rose-100/80 px-2 py-1 rounded border border-rose-300">
            Authorization: Super Admin Privilege Only
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">
            To proceed, type <span class="font-mono font-bold text-rose-600 select-all bg-rose-50 px-1 py-0.5 rounded border border-rose-200">DELETE ALL</span> below:
          </label>
          <input
            v-model="confirmDeleteAllText"
            type="text"
            placeholder="Type DELETE ALL to confirm"
            class="form-input text-xs font-mono text-center tracking-wider font-bold border-rose-300 focus:border-rose-500 focus:ring-rose-500"
            autocomplete="off"
            @keyup.enter="confirmDeleteAllText.trim() === 'DELETE ALL' && !isDeletingAll && executeDeleteAll()"
          />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeDeleteAllModal" :disabled="isDeletingAll">Cancel</button>
        <button
          class="btn bg-rose-600 hover:bg-rose-700 text-white font-bold disabled:opacity-40 disabled:cursor-not-allowed shadow-md shadow-rose-600/20 flex items-center gap-1.5"
          :disabled="confirmDeleteAllText.trim() !== 'DELETE ALL' || isDeletingAll"
          @click="executeDeleteAll"
        >
          <AppIcon v-if="!isDeletingAll" name="trash" size="xs" />
          <span>{{ isDeletingAll ? 'Purging Batches...' : `Delete All (${totalBatchesCount}) Batches` }}</span>
        </button>
      </template>
    </BaseModal>

    <!-- Floating Bulk Action Bar (Option B) -->
    <transition
      enter-active-class="transition transform duration-200 ease-out"
      enter-from-class="translate-y-12 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition transform duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-12 opacity-0"
    >
      <div
        v-if="selectedBatchIds.length > 0"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 z-40 bg-slate-900/95 text-white border border-slate-700/80 shadow-2xl rounded-2xl px-5 py-3 flex items-center gap-4 backdrop-blur-md max-w-[90vw]"
      >
        <div class="flex items-center gap-2 pr-3 border-r border-slate-700/80">
          <div class="w-6 h-6 rounded-full bg-brand-500/20 text-brand-400 flex items-center justify-center font-mono font-bold text-xs">
            {{ selectedBatchIds.length }}
          </div>
          <span class="text-xs font-semibold text-slate-200">Batches Selected</span>
        </div>

        <button
          @click="selectedBatchIds = []"
          class="text-xs text-slate-400 hover:text-white transition-colors underline-offset-2 hover:underline"
        >
          Deselect
        </button>

        <div class="flex items-center gap-2">
          <!-- Delete Selected Button -->
          <button
            @click="confirmDeleteSelected"
            class="btn btn-sm bg-rose-600/90 hover:bg-rose-600 text-white font-semibold text-xs border border-rose-500/40 shadow-sm flex items-center gap-1.5"
          >
            <AppIcon name="trash" size="xs" />
            <span>Delete Selected ({{ selectedBatchIds.length }})</span>
          </button>

          <!-- Delete All Batches Button (Super Admin Only) -->
          <button
            v-if="isSuperAdmin"
            @click="openDeleteAllModal"
            class="btn btn-sm bg-rose-950 hover:bg-rose-900 text-rose-300 hover:text-rose-100 font-semibold text-xs border border-rose-700/80 shadow-sm flex items-center gap-1.5 group"
            title="Super Admin Only: Purge all import sessions and staging data"
          >
            <AppIcon name="alert-triangle" size="xs" class="text-rose-400 group-hover:scale-110 transition-transform" />
            <span>Delete All ({{ totalBatchesCount }})</span>
            <span class="text-3xs font-mono px-1.5 py-0.2 rounded bg-rose-900/80 border border-rose-600/60 text-rose-200 uppercase tracking-widest font-bold">
              Super Admin
            </span>
          </button>
        </div>
      </div>
    </transition>

    <!-- Create Batch Modal -->
    <div v-if="showNewBatchModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-md p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="text-sm font-bold text-slate-900">Create New Import Batch</h3>
          <button @click="showNewBatchModal = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <form @submit.prevent="handleCreateBatch" class="space-y-3">
          <div>
            <label class="form-label text-xs">Import Year Scope</label>
            <select v-model="newBatchForm.sourceYear" class="form-select text-xs font-mono">
              <option value="2026">2026 (Active Ingestion Scope)</option>
              <option value="2025">2025 (Historical Archive)</option>
              <option value="2024">2024 (Historical Archive)</option>
              <option value="2023">2023 (Historical Archive)</option>
              <option value="2022">2022 (Historical Archive)</option>
            </select>
          </div>

          <div>
            <label class="form-label text-xs">Source Reference / Catalog Name</label>
            <input
              type="text"
              v-model="newBatchForm.sourceReference"
              placeholder="e.g. 2026 Q1 Commercial Hardware Pricelist"
              required
              class="form-input text-xs"
            />
          </div>

          <div>
            <label class="form-label text-xs">Auditor Notes</label>
            <textarea
              v-model="newBatchForm.sourceNotes"
              rows="3"
              placeholder="Add migration notes or source description..."
              class="form-input text-xs"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
            <button type="button" @click="showNewBatchModal = false" class="btn btn-secondary text-xs">Cancel</button>
            <button type="submit" class="btn btn-primary text-xs" :disabled="creatingBatch">
              <span v-if="creatingBatch">Creating...</span>
              <span v-else>Initialize Batch</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { useImportStore } from '../../stores/importStore';
import { useMasterStore } from '../../stores/masterStore';
import AppIcon from '../../components/common/AppIcon.vue';
import EmptyState from '../../components/common/EmptyState.vue';
import ConfirmDialog from '../../components/common/ConfirmDialog.vue';
import BaseModal from '../../components/common/BaseModal.vue';

export default {
  name: 'ImportDashboardView',
  components: { AppIcon, EmptyState, ConfirmDialog, BaseModal },
  data() {
    return {
      store: useImportStore(),
      masterStore: useMasterStore(),
      selectedBatchIds: [],
      deleteSelectedModalOpen: false,
      deleteAllModalOpen: false,
      confirmDeleteAllText: '',
      isDeletingAll: false,
      showNewBatchModal: false,
      creatingBatch: false,
      newBatchForm: {
        sourceYear: '2026',
        sourceReference: '',
        sourceNotes: ''
      }
    };
  },
  computed: {
    isSuperAdmin() {
      return this.masterStore.isSuperAdmin;
    },
    totalBatchesCount() {
      return (this.store.batches || []).length;
    },
    active2026RowsCount() {
      return (this.store.active2026Batches || []).reduce((acc, b) => acc + (b?.rowsCount || 0), 0);
    },
    validatedRate() {
      const total = (this.store.batches || []).reduce((acc, b) => acc + (b?.rowsCount || 0), 0);
      const valid = (this.store.batches || []).reduce((acc, b) => acc + (b?.validRows || 0), 0);
      if (total === 0) return 100;
      return Math.round((valid / total) * 100);
    },
    totalIssuesCount() {
      return (this.store.batches || []).reduce((acc, b) => acc + ((b?.warningRows || 0) + (b?.errorRows || 0)), 0);
    },
    selectAll: {
      get() {
        return (this.store.batches || []).length > 0 && this.selectedBatchIds.length === this.store.batches.length;
      },
      set(val) {
        if (val) {
          this.selectedBatchIds = (this.store.batches || []).map(b => b.id);
        } else {
          this.selectedBatchIds = [];
        }
      }
    }
  },
  mounted() {
    this.store.fetchBatches();
  },
  methods: {
    confirmDeleteSelected() {
      if (this.selectedBatchIds.length === 0) return;
      this.deleteSelectedModalOpen = true;
    },
    async executeDeleteSelected() {
      await this.store.deleteMultipleBatches(this.selectedBatchIds);
      this.selectedBatchIds = [];
      this.deleteSelectedModalOpen = false;
      this.masterStore.showToast('Batches Deleted', 'Selected import batches have been removed.', 'warning');
    },
    openDeleteAllModal() {
      this.confirmDeleteAllText = '';
      this.deleteAllModalOpen = true;
    },
    closeDeleteAllModal() {
      this.deleteAllModalOpen = false;
      this.confirmDeleteAllText = '';
      this.isDeletingAll = false;
    },
    async executeDeleteAll() {
      if (this.confirmDeleteAllText.trim() !== 'DELETE ALL') return;
      this.isDeletingAll = true;
      try {
        await this.store.deleteAllBatches();
        this.selectedBatchIds = [];
        this.closeDeleteAllModal();
        this.masterStore.showToast('Purge Completed', 'All import batches and staging records have been deleted.', 'warning');
      } finally {
        this.isDeletingAll = false;
      }
    },
    setYear(year) {
      this.store.selectedYear = year;
      this.store.fetchBatches();
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'IMPORTED': return 'bg-emerald-100 text-emerald-800 border border-emerald-300';
        case 'APPROVED': return 'bg-blue-100 text-blue-800 border border-blue-300';
        case 'READY_FOR_REVIEW': return 'bg-purple-100 text-purple-800 border border-purple-300';
        case 'VALIDATING': return 'bg-amber-100 text-amber-800 border border-amber-300';
        case 'PARSED': return 'bg-sky-100 text-sky-800 border border-sky-300';
        case 'FAILED': return 'bg-rose-100 text-rose-800 border border-rose-300';
        default: return 'bg-slate-100 text-slate-700 border border-slate-200';
      }
    },
    formatTimestamp(dateStr) {
      if (!dateStr) return '-';
      const d = new Date(dateStr);
      return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    },
    async handleCreateBatch() {
      this.creatingBatch = true;
      try {
        const batch = await this.store.createBatch(this.newBatchForm);
        this.showNewBatchModal = false;
        this.$router.push(`/data-import/wizard/${batch.id}`);
      } catch (err) {
        alert('Failed to initialize batch: ' + err.message);
      } finally {
        this.creatingBatch = false;
      }
    }
  }
};
</script>
