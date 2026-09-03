<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Customers" subtitle="Client accounts, buyer organizations, and procurement contacts from import quotes & direct records">
      <template #actions>
        <button @click="openAddModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add Customer</span>
        </button>
      </template>
    </PageHeader>

    <!-- KPI Metric Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <!-- Total Customers -->
      <div class="app-card p-3.5 bg-white border border-slate-200 flex items-center justify-between">
        <div>
          <span class="text-3xs uppercase tracking-wider font-semibold text-slate-400">Total Customers</span>
          <div class="font-mono font-bold text-xl text-slate-900 mt-0.5">{{ store.customers.length }}</div>
        </div>
        <div class="w-9 h-9 rounded-lg bg-brand-50 border border-brand-200 flex items-center justify-center text-brand-600">
          <AppIcon name="users" size="sm" />
        </div>
      </div>

      <!-- Active Clients -->
      <div class="app-card p-3.5 bg-white border border-slate-200 flex items-center justify-between">
        <div>
          <span class="text-3xs uppercase tracking-wider font-semibold text-slate-400">Active Clients</span>
          <div class="font-mono font-bold text-xl text-emerald-600 mt-0.5">{{ activeCustomersCount }}</div>
        </div>
        <div class="w-9 h-9 rounded-lg bg-emerald-50 border border-emerald-200 flex items-center justify-center text-emerald-600">
          <AppIcon name="check-circle" size="sm" />
        </div>
      </div>

      <!-- Corporate Companies -->
      <div class="app-card p-3.5 bg-white border border-slate-200 flex items-center justify-between">
        <div>
          <span class="text-3xs uppercase tracking-wider font-semibold text-slate-400">Companies / Orgs</span>
          <div class="font-mono font-bold text-xl text-indigo-600 mt-0.5">{{ uniqueCompaniesCount }}</div>
        </div>
        <div class="w-9 h-9 rounded-lg bg-indigo-50 border border-indigo-200 flex items-center justify-center text-indigo-600">
          <AppIcon name="briefcase" size="sm" />
        </div>
      </div>

      <!-- Sourced from Import XLS -->
      <div class="app-card p-3.5 bg-white border border-slate-200 flex items-center justify-between">
        <div>
          <span class="text-3xs uppercase tracking-wider font-semibold text-slate-400">From Import XLS</span>
          <div class="font-mono font-bold text-xl text-cyan-600 mt-0.5">{{ importSourcedCount }}</div>
        </div>
        <div class="w-9 h-9 rounded-lg bg-cyan-50 border border-cyan-200 flex items-center justify-center text-cyan-600">
          <AppIcon name="file-text" size="sm" />
        </div>
      </div>
    </div>

    <!-- Filters & Search Workspace -->
    <div class="bg-white rounded-lg border border-slate-200 p-3.5 shadow-subtle flex flex-wrap items-center justify-between gap-3">
      <div class="flex flex-wrap items-center gap-2.5 flex-1 min-w-[280px]">
        <!-- Search Input -->
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <AppIcon name="search" size="xs" class="absolute left-2.5 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search contact name, company, phone, email, notes..."
            class="form-input pl-8"
          />
        </div>

        <!-- Status Filter -->
        <select v-model="selectedStatus" class="form-select text-xs w-36">
          <option value="All">All Status</option>
          <option value="Active">Active</option>
          <option value="Inactive">Inactive</option>
          <option value="Lead">Lead</option>
        </select>

        <!-- Source Filter -->
        <select v-model="selectedSource" class="form-select text-xs w-40">
          <option value="All">All Sources</option>
          <option value="IMPORT_XLS">From Import XLS</option>
          <option value="MANUAL">Manual Entry</option>
        </select>

        <!-- Reset Button -->
        <button
          v-if="isFiltered"
          @click="resetFilters"
          class="btn btn-secondary text-xs text-slate-500 hover:text-slate-700"
        >
          <AppIcon name="x" size="xs" class="mr-1" />
          <span>Reset</span>
        </button>
      </div>

      <div class="text-2xs text-slate-500 font-mono">
        Showing <span class="font-bold text-slate-800">{{ filteredCustomers.length }}</span> of {{ store.customers.length }} customers
      </div>
    </div>

    <!-- Customers Table -->
    <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
      <div v-if="filteredCustomers.length === 0" class="p-12 text-center text-slate-400 space-y-3">
        <div class="w-12 h-12 rounded-full bg-slate-100 flex items-center justify-center mx-auto text-slate-400">
          <AppIcon name="users" size="md" />
        </div>
        <div class="text-sm font-semibold text-slate-700">No customers found</div>
        <p class="text-xs text-slate-500 max-w-sm mx-auto">
          No customer records match your filter criteria. Add a new customer or import quotations via Excel.
        </p>
        <button @click="openAddModal" class="btn btn-primary text-xs">
          <AppIcon name="plus" size="xs" class="mr-1" />
          <span>Add First Customer</span>
        </button>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="dense-table w-full">
          <thead>
            <tr>
              <th class="whitespace-nowrap">Customer Code</th>
              <th class="whitespace-nowrap">Contact Person (Name)</th>
              <th class="whitespace-nowrap">Company / Organization</th>
              <th class="whitespace-nowrap">Phone</th>
              <th class="whitespace-nowrap">Email</th>
              <th class="whitespace-nowrap">Description / Remarks</th>
              <th class="whitespace-nowrap text-center">Source</th>
              <th class="whitespace-nowrap text-center">Status</th>
              <th class="text-right whitespace-nowrap">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 font-sans">
            <tr
              v-for="c in filteredCustomers"
              :key="c.id"
              class="hover:bg-slate-50/70 transition-colors group"
            >
              <!-- Code -->
              <td class="whitespace-nowrap font-mono text-2xs font-bold text-slate-600">
                {{ c.code || c.id }}
              </td>

              <!-- Name with Initial Avatar -->
              <td class="whitespace-nowrap">
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-full bg-brand-50 border border-brand-200 flex items-center justify-center text-brand-700 font-bold text-2xs shrink-0 uppercase">
                    {{ getInitials(c.name || c.company) }}
                  </div>
                  <div class="font-bold text-slate-900 group-hover:text-brand-600 transition-colors">
                    {{ c.name || '—' }}
                  </div>
                </div>
              </td>

              <!-- Company -->
              <td class="whitespace-nowrap">
                <div v-if="c.company" class="flex items-center gap-1.5 font-medium text-slate-700">
                  <AppIcon name="briefcase" size="2xs" class="text-slate-400" />
                  <span>{{ c.company }}</span>
                </div>
                <span v-else class="text-slate-400 text-3xs italic">No Company</span>
              </td>

              <!-- Phone -->
              <td class="whitespace-nowrap font-mono text-2xs">
                <a
                  v-if="c.phone && c.phone !== '-'"
                  :href="`tel:${c.phone}`"
                  class="text-brand-600 hover:text-brand-800 hover:underline flex items-center gap-1"
                >
                  <AppIcon name="phone" size="2xs" class="text-slate-400" />
                  <span>{{ c.phone }}</span>
                </a>
                <span v-else class="text-slate-400 text-3xs">—</span>
              </td>

              <!-- Email -->
              <td class="whitespace-nowrap font-mono text-2xs">
                <a
                  v-if="c.email && c.email !== '-'"
                  :href="`mailto:${c.email}`"
                  class="text-brand-600 hover:text-brand-800 hover:underline flex items-center gap-1"
                >
                  <AppIcon name="mail" size="2xs" class="text-slate-400" />
                  <span>{{ c.email }}</span>
                </a>
                <span v-else class="text-slate-400 text-3xs">—</span>
              </td>

              <!-- Description / Remarks -->
              <td class="max-w-[220px] truncate text-slate-500 text-2xs" :title="c.descr">
                {{ c.descr || '—' }}
              </td>

              <!-- Source Badge -->
              <td class="text-center whitespace-nowrap">
                <span
                  v-if="c.source === 'IMPORT_XLS'"
                  class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-cyan-50 text-cyan-700 border border-cyan-200 inline-flex items-center gap-1"
                >
                  <AppIcon name="file-text" size="2xs" />
                  <span>IMPORT XLS</span>
                </span>
                <span
                  v-else
                  class="text-3xs font-mono font-medium px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 border border-slate-200 inline-flex items-center gap-1"
                >
                  <span>MANUAL</span>
                </span>
              </td>

              <!-- Status Badge -->
              <td class="text-center whitespace-nowrap">
                <span
                  class="text-3xs font-semibold px-2 py-0.5 rounded-full border"
                  :class="[
                    c.status === 'Active' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' :
                    c.status === 'Lead' ? 'bg-amber-50 text-amber-700 border-amber-200' :
                    'bg-slate-100 text-slate-600 border-slate-200'
                  ]"
                >
                  {{ c.status || 'Active' }}
                </span>
              </td>

              <!-- Actions -->
              <td class="text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click.stop="openEditModal(c)"
                    class="p-1 rounded hover:bg-slate-100 text-slate-500 hover:text-brand-600 transition-colors"
                    title="Edit Customer"
                  >
                    <AppIcon name="edit" size="xs" />
                  </button>
                  <button
                    @click.stop="confirmDelete(c)"
                    class="p-1 rounded hover:bg-rose-50 text-slate-400 hover:text-rose-600 transition-colors"
                    title="Delete Customer"
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

    <!-- Modal Form: Add / Edit Customer -->
    <BaseModal
      :isOpen="isModalOpen"
      :title="isEditMode ? 'Edit Customer' : 'Add New Customer'"
      @close="closeModal"
    >
      <form @submit.prevent="saveCustomer" class="space-y-4 text-xs">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Contact Name -->
          <div>
            <label class="form-label">Contact Person (Name) <span class="text-rose-500">*</span></label>
            <input
              v-model="form.name"
              type="text"
              required
              class="form-input"
              placeholder="e.g. Bpk. Hendra Gunawan"
            />
          </div>

          <!-- Company Name -->
          <div>
            <label class="form-label">Company / Client Organization</label>
            <input
              v-model="form.company"
              type="text"
              class="form-input"
              placeholder="e.g. PT. Indah Kiat Pulp & Paper"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Phone -->
          <div>
            <label class="form-label">Phone / WhatsApp</label>
            <input
              v-model="form.phone"
              type="text"
              class="form-input font-mono"
              placeholder="e.g. 0812-3456-7890"
            />
          </div>

          <!-- Email -->
          <div>
            <label class="form-label">Email Address</label>
            <input
              v-model="form.email"
              type="email"
              class="form-input font-mono"
              placeholder="e.g. hendra@client.co.id"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Status -->
          <div>
            <label class="form-label">Customer Status</label>
            <select v-model="form.status" class="form-select">
              <option value="Active">Active</option>
              <option value="Lead">Lead</option>
              <option value="Inactive">Inactive</option>
            </select>
          </div>

          <!-- Code (Optional / Auto) -->
          <div>
            <label class="form-label">Customer Code</label>
            <input
              v-model="form.code"
              type="text"
              class="form-input font-mono"
              placeholder="Leave blank for auto CUST-2024-xxx"
            />
          </div>
        </div>

        <!-- Description / Address / Remarks -->
        <div>
          <label class="form-label">Description / Address / Remarks</label>
          <textarea
            v-model="form.descr"
            rows="3"
            class="form-input"
            placeholder="Office address, project requirements, or client notes..."
          ></textarea>
        </div>

        <!-- Modal Footer Buttons -->
        <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
          <button type="button" @click="closeModal" class="btn btn-secondary">
            Cancel
          </button>
          <button type="submit" :disabled="isSaving" class="btn btn-primary">
            <AppIcon v-if="!isSaving" name="check" size="xs" class="mr-1.5" />
            <span>{{ isSaving ? 'Saving...' : (isEditMode ? 'Update Customer' : 'Save Customer') }}</span>
          </button>
        </div>
      </form>
    </BaseModal>

    <!-- Confirm Delete Modal -->
    <BaseModal
      :isOpen="isDeleteModalOpen"
      title="Delete Customer"
      @close="isDeleteModalOpen = false"
    >
      <div class="space-y-4 text-xs">
        <p class="text-slate-600">
          Are you sure you want to delete customer
          <strong class="text-slate-900">{{ customerToDelete?.name || customerToDelete?.company }}</strong>?
          This action will mark the customer as removed.
        </p>
        <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
          <button type="button" @click="isDeleteModalOpen = false" class="btn btn-secondary">
            Cancel
          </button>
          <button type="button" @click="handleDelete" class="btn btn-danger">
            Delete Customer
          </button>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'CustomerListView',
  components: {
    PageHeader,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedStatus: 'All',
      selectedSource: 'All',
      isModalOpen: false,
      isEditMode: false,
      isSaving: false,
      isDeleteModalOpen: false,
      customerToDelete: null,
      form: {
        id: '',
        code: '',
        name: '',
        company: '',
        phone: '',
        email: '',
        descr: '',
        status: 'Active',
        source: 'MANUAL'
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    activeCustomersCount() {
      return this.store.customers.filter(c => c.status === 'Active').length;
    },
    uniqueCompaniesCount() {
      const set = new Set(this.store.customers.map(c => c.company).filter(Boolean));
      return set.size;
    },
    importSourcedCount() {
      return this.store.customers.filter(c => c.source === 'IMPORT_XLS').length;
    },
    isFiltered() {
      return this.searchQuery !== '' || this.selectedStatus !== 'All' || this.selectedSource !== 'All';
    },
    filteredCustomers() {
      return this.store.customers.filter(c => {
        const matchesStatus = this.selectedStatus === 'All' || c.status === this.selectedStatus;
        const matchesSource = this.selectedSource === 'All' || c.source === this.selectedSource;
        
        if (!matchesStatus || !matchesSource) return false;

        if (!this.searchQuery) return true;
        const q = this.searchQuery.toLowerCase();
        return (
          (c.name && c.name.toLowerCase().includes(q)) ||
          (c.company && c.company.toLowerCase().includes(q)) ||
          (c.code && c.code.toLowerCase().includes(q)) ||
          (c.phone && c.phone.toLowerCase().includes(q)) ||
          (c.email && c.email.toLowerCase().includes(q)) ||
          (c.descr && c.descr.toLowerCase().includes(q))
        );
      });
    }
  },
  async mounted() {
    await this.store.fetchCustomers();
  },
  methods: {
    getInitials(name) {
      if (!name) return 'CU';
      const parts = name.trim().split(/\s+/);
      if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
      return (parts[0][0] + parts[1][0]).toUpperCase();
    },
    resetFilters() {
      this.searchQuery = '';
      this.selectedStatus = 'All';
      this.selectedSource = 'All';
    },
    openAddModal() {
      this.isEditMode = false;
      this.form = {
        id: '',
        code: '',
        name: '',
        company: '',
        phone: '',
        email: '',
        descr: '',
        status: 'Active',
        source: 'MANUAL'
      };
      this.isModalOpen = true;
    },
    openEditModal(customer) {
      this.isEditMode = true;
      this.form = { ...customer };
      this.isModalOpen = true;
    },
    closeModal() {
      this.isModalOpen = false;
    },
    async saveCustomer() {
      this.isSaving = true;
      try {
        await this.store.saveCustomer(this.form);
        this.closeModal();
      } catch (err) {
        console.error('Failed to save customer:', err);
      } finally {
        this.isSaving = false;
      }
    },
    confirmDelete(customer) {
      this.customerToDelete = customer;
      this.isDeleteModalOpen = true;
    },
    async handleDelete() {
      if (!this.customerToDelete) return;
      try {
        await this.store.deleteCustomer(this.customerToDelete.id);
        this.isDeleteModalOpen = false;
        this.customerToDelete = null;
      } catch (err) {
        console.error('Failed to delete customer:', err);
      }
    }
  }
};
</script>
