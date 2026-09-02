<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Suppliers" subtitle="Authorized distributors and procurement partners">
      <template #actions>
        <button @click="openAddModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add Supplier</span>
        </button>
      </template>
    </PageHeader>

    <!-- Filters & Search -->
    <div class="bg-white rounded-lg border border-slate-200 p-3 shadow-subtle flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-2.5 flex-1 min-w-[240px] max-w-sm">
        <div class="relative w-full">
          <AppIcon name="search" size="xs" class="absolute left-2.5 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search suppliers, contact, country..."
            class="form-input pl-8"
          />
        </div>
      </div>
      <div class="text-2xs text-slate-500 font-mono">
        Active Suppliers: <span class="font-bold text-slate-800">{{ filteredSuppliers.length }}</span>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
      <table class="dense-table">
        <thead>
          <tr>
            <th>Supplier</th>
            <th>Supplier Code</th>
            <th>Country</th>
            <th>Key Contact</th>
            <th>Tier / Agreement</th>
            <th>Total Components</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 font-sans">
          <tr
            v-for="s in filteredSuppliers"
            :key="s.id"
            @click="inspectSupplier(s)"
            class="cursor-pointer group"
          >
            <td>
              <div class="flex items-center gap-2.5">
                <div class="w-7 h-7 rounded bg-brand-50 border border-brand-200 flex items-center justify-center text-brand-700 font-bold text-xs shrink-0">
                  {{ s.name.slice(0, 2) }}
                </div>
                <span class="font-bold text-slate-900 group-hover:text-brand-600 transition-colors">{{ s.name }}</span>
              </div>
            </td>
            <td>
              <span class="font-mono text-2xs px-1.5 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200 font-bold">
                {{ s.code }}
              </span>
            </td>
            <td>
              <div class="flex items-center gap-1.5 text-xs text-slate-700">
                <span class="font-mono text-3xs px-1 rounded bg-slate-100 border">{{ s.countryCode }}</span>
                <span>{{ s.country }}</span>
              </div>
            </td>
            <td>
              <div class="text-xs font-semibold text-slate-800">{{ s.contact }}</div>
              <div class="text-3xs text-slate-400 font-mono">{{ s.email }}</div>
            </td>
            <td>
              <span class="text-3xs font-semibold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                {{ s.tier }}
              </span>
            </td>
            <td>
              <span class="font-mono font-bold text-xs text-slate-800">{{ s.totalComponents }}</span>
            </td>
            <td>
              <StatusBadge :status="s.status" />
            </td>
            <td class="text-right" @click.stop>
              <div class="flex items-center justify-end gap-1">
                <button @click="inspectSupplier(s)" class="btn-icon" title="View Details">
                  <AppIcon name="eye" size="xs" />
                </button>
                <button @click="openEditModal(s)" class="btn-icon" title="Edit">
                  <AppIcon name="edit" size="xs" />
                </button>
                <button @click="confirmDelete(s)" class="btn-icon text-rose-500 hover:text-rose-700" title="Delete">
                  <AppIcon name="trash" size="xs" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Supplier Detail Drawer -->
    <BaseDrawer
      :show="detailDrawerOpen"
      :title="inspectedSup ? inspectedSup.name : ''"
      :subtitle="inspectedSup ? `Supplier Code: ${inspectedSup.code}` : ''"
      width="lg"
      @close="detailDrawerOpen = false"
    >
      <div v-if="inspectedSup" class="space-y-4 text-xs">
        <div class="flex items-center gap-2">
          <StatusBadge :status="inspectedSup.status" />
          <span class="text-3xs font-mono px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200 font-bold">
            {{ inspectedSup.tier }}
          </span>
        </div>

        <div class="p-3 bg-slate-50 border border-slate-200 rounded-md">
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Company Overview</div>
          <p class="text-slate-700 mt-1 leading-relaxed whitespace-pre-line">{{ inspectedSup.description }}</p>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="text-3xs text-slate-400">Payment & Terms</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ inspectedSup.terms || 'Net 30 Days' }}</div>
            <div class="text-3xs text-slate-500 mt-0.5">Avg Lead Time: {{ inspectedSup.leadTimeAvg || '2-5 Days' }}</div>
          </div>
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="text-3xs text-slate-400">Procurement Contact</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ inspectedSup.contact }}</div>
            <div class="text-3xs text-slate-500 font-mono">{{ inspectedSup.email }}</div>
            <div class="text-3xs text-slate-500 font-mono">{{ inspectedSup.phone }}</div>
          </div>
        </div>

        <div>
          <h4 class="font-bold text-slate-800 uppercase text-3xs tracking-wider mb-2">Components Sourced from {{ inspectedSup.name }}</h4>
          <div class="border border-slate-200 rounded-md overflow-hidden">
            <table class="dense-table">
              <thead>
                <tr>
                  <th>MPN</th>
                  <th>Name</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="c in relatedComponents"
                  :key="c.id"
                  @click="$router.push(`/components/${c.id}`)"
                  class="cursor-pointer hover:bg-slate-50"
                >
                  <td class="font-mono font-bold text-brand-700">{{ c.partNumber }}</td>
                  <td class="text-slate-800 font-medium truncate max-w-xs">{{ c.name }}</td>
                  <td><StatusBadge :status="c.status" /></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </BaseDrawer>

    <!-- Create / Edit Supplier Modal -->
    <BaseModal :show="modalOpen" :title="isEdit ? 'Edit Supplier' : 'Add Supplier'" maxWidth="md" @close="modalOpen = false">
      <form @submit.prevent="saveModal" class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Supplier Name *</label>
          <input v-model="modalForm.name" type="text" required placeholder="e.g. DigiKey Electronics" class="form-input" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Supplier Code *</label>
            <input v-model="modalForm.code" type="text" required placeholder="e.g. DIGIKEY-US" class="form-input font-mono uppercase" />
          </div>
          <div>
            <label class="form-label">Country Code</label>
            <input v-model="modalForm.countryCode" type="text" placeholder="e.g. US, CN, DE" class="form-input font-mono uppercase" />
          </div>
        </div>
        <div>
          <label class="form-label">Country *</label>
          <input v-model="modalForm.country" type="text" required placeholder="e.g. United States" class="form-input" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Contact Person</label>
            <input v-model="modalForm.contact" type="text" placeholder="e.g. Dave Larson" class="form-input" />
          </div>
          <div>
            <label class="form-label">Contact Email</label>
            <input v-model="modalForm.email" type="email" placeholder="sales@digikey.com" class="form-input font-mono" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Tier / Status</label>
            <select v-model="modalForm.tier" class="form-select">
              <option value="Tier 1 Authorized Distributor">Tier 1 Authorized Distributor</option>
              <option value="Global Volume Partner">Global Volume Partner</option>
              <option value="Fast-Turn Regional Partner">Fast-Turn Regional Partner</option>
            </select>
          </div>
          <div>
            <label class="form-label">Status</label>
            <select v-model="modalForm.status" class="form-select">
              <option value="Active">Active</option>
              <option value="Under Review">Under Review</option>
            </select>
          </div>
        </div>
        <div>
          <label class="form-label">Description</label>
          <textarea v-model="modalForm.description" rows="2" placeholder="Procurement terms and logistics capabilities..." class="form-textarea"></textarea>
        </div>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="modalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Save Supplier</button>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Supplier"
      :message="`Are you sure you want to delete supplier '${supToDelete ? supToDelete.name : ''}'?`"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import StatusBadge from '@/components/common/StatusBadge.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import BaseDrawer from '@/components/common/BaseDrawer.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'SupplierListView',
  components: {
    PageHeader,
    StatusBadge,
    ConfirmDialog,
    BaseDrawer,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      detailDrawerOpen: false,
      inspectedSup: null,
      modalOpen: false,
      isEdit: false,
      modalForm: {
        id: null,
        name: '',
        code: '',
        country: '',
        countryCode: 'US',
        contact: '',
        email: '',
        phone: '',
        tier: 'Tier 1 Authorized Distributor',
        status: 'Active',
        description: ''
      },
      deleteModalOpen: false,
      supToDelete: null
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    filteredSuppliers() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.store.suppliers.filter(s => {
        return !q || s.name.toLowerCase().includes(q) || s.code.toLowerCase().includes(q) || s.contact.toLowerCase().includes(q);
      });
    },
    relatedComponents() {
      if (!this.inspectedSup) return [];
      return this.store.components.filter(c => c.supplier === this.inspectedSup.name);
    }
  },
  methods: {
    inspectSupplier(s) {
      this.inspectedSup = s;
      this.detailDrawerOpen = true;
    },
    openAddModal() {
      this.isEdit = false;
      this.modalForm = {
        id: null,
        name: '',
        code: '',
        country: '',
        countryCode: 'US',
        contact: '',
        email: '',
        phone: '',
        tier: 'Tier 1 Authorized Distributor',
        status: 'Active',
        description: ''
      };
      this.modalOpen = true;
    },
    openEditModal(s) {
      this.isEdit = true;
      this.modalForm = JSON.parse(JSON.stringify(s));
      this.modalOpen = true;
    },
    saveModal() {
      this.store.saveSupplier(this.modalForm);
      this.modalOpen = false;
    },
    confirmDelete(s) {
      this.supToDelete = s;
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.supToDelete) {
        this.store.deleteSupplier(this.supToDelete.id);
        this.deleteModalOpen = false;
        this.supToDelete = null;
        if (this.inspectedSup && this.inspectedSup.id === this.supToDelete?.id) {
          this.detailDrawerOpen = false;
        }
      }
    }
  }
};
</script>
