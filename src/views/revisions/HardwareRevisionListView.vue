<template>
  <div class="space-y-5">
    <!-- Page Header -->
    <PageHeader
      title="Hardware Revisions"
      subtitle="Manage PCB layout iterations, schematic revisions, gerber packages, and hardware engineering changes"
    >
      <template #badge>
        <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
          {{ versionStore.totalRevisionsCount }} Revisions
        </span>
      </template>
      <template #actions>
        <button @click="openCreateModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
          <span>New Revision</span>
        </button>
      </template>
    </PageHeader>

    <!-- Filters & Search Bar -->
    <div class="app-card p-3.5 flex flex-col md:flex-row md:items-center justify-between gap-3 bg-white">
      <div class="flex flex-wrap items-center gap-2.5 flex-1">
        <!-- Search Input -->
        <div class="relative min-w-[240px] max-w-sm flex-1">
          <AppIcon name="search" size="xs" class="absolute left-3 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search revision (e.g. REV-B), PCB rev, or summary..."
            class="form-input pl-8.5 text-xs font-mono"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-2.5 top-2 text-slate-400 hover:text-slate-600"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Product Filter -->
        <select v-model="selectedProduct" class="form-select w-auto min-w-[150px] text-xs">
          <option value="">All Products</option>
          <option v-for="p in uniqueProducts" :key="p" :value="p">{{ p }}</option>
        </select>

        <!-- Status Filter -->
        <select v-model="selectedStatus" class="form-select w-auto min-w-[130px] text-xs">
          <option value="">All Statuses</option>
          <option value="Released">Released</option>
          <option value="Approved">Approved</option>
          <option value="In Review">In Review</option>
          <option value="Draft">Draft</option>
          <option value="Superseded">Superseded</option>
          <option value="Obsolete">Obsolete</option>
        </select>

        <!-- Reset Button -->
        <button
          v-if="isFiltered"
          @click="resetFilters"
          class="btn btn-ghost text-slate-500 hover:text-slate-800"
        >
          <AppIcon name="refresh" size="xs" class="mr-1" />
          <span>Reset</span>
        </button>
      </div>

      <!-- Count Summary -->
      <div class="text-3xs font-mono text-slate-500 shrink-0">
        Showing <span class="font-bold text-slate-900">{{ filteredRevisions.length }}</span> of {{ versionStore.totalRevisionsCount }} revisions
      </div>
    </div>

    <!-- Revisions Table Card -->
    <div class="app-card overflow-hidden">
      <div v-if="filteredRevisions.length === 0" class="p-8">
        <EmptyState
          title="No hardware revisions found"
          description="No engineering hardware revisions match your search filters."
          icon="cpu"
        >
          <template #action>
            <button v-if="isFiltered" @click="resetFilters" class="btn btn-secondary mr-2">Reset Filters</button>
            <button @click="openCreateModal" class="btn btn-primary">Create Revision</button>
          </template>
        </EmptyState>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="dense-table">
          <thead>
            <tr>
              <th>Revision</th>
              <th>Product & Version</th>
              <th>Revision Name & Changes</th>
              <th>PCB / Schematic Ref</th>
              <th>Status</th>
              <th>Lead Engineer</th>
              <th>Release Date</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="rev in filteredRevisions"
              :key="rev.id"
              @click="$router.push(`/revisions/${rev.id}`)"
              class="cursor-pointer group"
            >
              <td>
                <span class="font-mono font-bold text-xs text-slate-900 bg-slate-100 px-2 py-0.5 rounded border border-slate-200 shadow-2xs group-hover:border-brand-300">
                  {{ rev.code }}
                </span>
              </td>
              <td>
                <div>
                  <div class="font-bold text-slate-900 text-xs truncate max-w-xs group-hover:text-brand-600 transition-colors">
                    {{ rev.productName }}
                  </div>
                  <div class="text-3xs font-mono text-slate-500 mt-0.5 flex items-center gap-1.5">
                    <span class="font-semibold text-brand-700">{{ rev.productCode }}</span>
                    <span>•</span>
                    <span class="font-bold text-slate-800">v{{ rev.versionNumber }}</span>
                  </div>
                </div>
              </td>
              <td class="max-w-md">
                <div class="font-semibold text-slate-800 text-xs truncate">{{ rev.name }}</div>
                <p class="text-3xs text-slate-500 truncate mt-0.5 font-medium">{{ rev.changeSummary }}</p>
              </td>
              <td>
                <div class="text-3xs font-mono space-y-0.5">
                  <div class="text-slate-700 font-semibold">{{ rev.pcbRevision }}</div>
                  <div class="text-slate-400">{{ rev.schematicRevision }}</div>
                </div>
              </td>
              <td>
                <RevisionBadge :status="rev.status" />
              </td>
              <td>
                <span class="text-xs text-slate-800 font-medium">{{ rev.engineer }}</span>
              </td>
              <td class="font-mono text-3xs text-slate-500">
                {{ rev.releaseDate || rev.createdAt }}
              </td>
              <td class="text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <router-link
                    :to="`/revisions/${rev.id}`"
                    class="btn-icon"
                    title="View Revision Profile"
                  >
                    <AppIcon name="eye" size="xs" />
                  </router-link>
                  <button
                    @click="confirmDelete(rev)"
                    class="btn-icon text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                    title="Delete Revision"
                  >
                    <AppIcon name="trash" size="xs" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Quick Create Revision Modal -->
    <BaseModal
      :show="createModalOpen"
      title="Create Hardware Revision"
      maxWidth="md"
      @close="createModalOpen = false"
    >
      <div class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Target Product Version *</label>
          <select v-model="createForm.versionId" @change="onCreateVersionChange" class="form-select font-mono">
            <option value="">Select version...</option>
            <option
              v-for="v in versionStore.versions"
              :key="v.id"
              :value="v.id"
            >
              {{ v.productCode }} — v{{ v.versionNumber }} ({{ v.versionName }})
            </option>
          </select>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Revision Code *</label>
            <input
              v-model="createForm.code"
              type="text"
              placeholder="e.g. REV-C"
              class="form-input font-mono uppercase"
            />
          </div>
          <div>
            <label class="form-label">Status *</label>
            <select v-model="createForm.status" class="form-select font-mono">
              <option value="Draft">Draft</option>
              <option value="In Review">In Review</option>
              <option value="Approved">Approved</option>
              <option value="Released">Released</option>
            </select>
          </div>
        </div>
        <div>
          <label class="form-label">Revision Title *</label>
          <input
            v-model="createForm.name"
            type="text"
            placeholder="e.g. Power Plane & Filter Optimization"
            class="form-input"
          />
        </div>
        <div>
          <label class="form-label">Change Summary</label>
          <textarea
            v-model="createForm.changeSummary"
            rows="3"
            placeholder="Key schematic, layout, or decoupling modifications..."
            class="form-textarea"
          ></textarea>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="createModalOpen = false">Cancel</button>
        <button class="btn btn-primary" @click="saveCreatedRevision">Create Revision</button>
      </template>
    </BaseModal>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import PageHeader from '@/components/common/PageHeader.vue';
import RevisionBadge from '@/components/versions/RevisionBadge.vue';
import EmptyState from '@/components/common/EmptyState.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'HardwareRevisionListView',
  components: {
    PageHeader,
    RevisionBadge,
    EmptyState,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedProduct: '',
      selectedStatus: '',
      createModalOpen: false,
      createForm: {
        versionId: '',
        versionNumber: '',
        productCode: '',
        productName: '',
        code: 'REV-C',
        name: '',
        status: 'Draft',
        changeSummary: '',
        engineer: 'Marcus Vance'
      }
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    uniqueProducts() {
      const set = new Set(this.versionStore.revisions.map(r => r.productCode).filter(Boolean));
      return Array.from(set);
    },
    isFiltered() {
      return Boolean(this.searchQuery || this.selectedProduct || this.selectedStatus);
    },
    filteredRevisions() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.versionStore.revisions.filter(r => {
        const matchesSearch = !q || r.code.toLowerCase().includes(q) || r.name.toLowerCase().includes(q) || r.productName.toLowerCase().includes(q) || r.productCode.toLowerCase().includes(q) || (r.pcbRevision && r.pcbRevision.toLowerCase().includes(q)) || (r.changeSummary && r.changeSummary.toLowerCase().includes(q));
        const matchesProduct = !this.selectedProduct || r.productCode === this.selectedProduct;
        const matchesStatus = !this.selectedStatus || r.status === this.selectedStatus;
        return matchesSearch && matchesProduct && matchesStatus;
      });
    }
  },
  methods: {
    resetFilters() {
      this.searchQuery = '';
      this.selectedProduct = '';
      this.selectedStatus = '';
    },
    openCreateModal() {
      const firstVer = this.versionStore.versions[0];
      this.createForm = {
        versionId: firstVer ? firstVer.id : '',
        versionNumber: firstVer ? firstVer.versionNumber : '1.0.0',
        productCode: firstVer ? firstVer.productCode : '',
        productName: firstVer ? firstVer.productName : '',
        code: 'REV-C',
        name: '',
        status: 'Draft',
        changeSummary: '',
        engineer: 'Marcus Vance'
      };
      this.createModalOpen = true;
    },
    onCreateVersionChange() {
      const v = this.versionStore.getVersionById(this.createForm.versionId);
      if (v) {
        this.createForm.versionNumber = v.versionNumber;
        this.createForm.productCode = v.productCode;
        this.createForm.productName = v.productName;
      }
    },
    async saveCreatedRevision() {
      if (!this.createForm.code || !this.createForm.name || !this.createForm.versionId) {
        alert('Please fill all required fields.');
        return;
      }
      try {
        await this.versionStore.saveRevision({
          ...this.createForm,
          pcbRevision: `PCB-v${this.createForm.versionNumber}-${this.createForm.code.replace('REV-', '')}`,
          schematicRevision: `SCH-v${this.createForm.versionNumber}-01`,
          changesList: [this.createForm.changeSummary]
        });
        this.createModalOpen = false;
      } catch (err) {
        alert('Failed to save revision: ' + (err.response?.data?.message || err.message));
      }
    },
    async confirmDelete(rev) {
      if (confirm(`Are you sure you want to delete ${rev.code} (${rev.name})?`)) {
        try {
          await this.versionStore.deleteRevision(rev.id);
        } catch (err) {
          alert('Failed to delete revision: ' + (err.response?.data?.message || err.message));
        }
      }
    }
  },
  async mounted() {
    await Promise.all([
      this.versionStore.fetchVersions(),
      this.versionStore.fetchRevisions()
    ]);
  }
};
</script>
