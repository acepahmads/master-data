<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Manufacturers" subtitle="Original component and semiconductor manufacturers">
      <template #actions>
        <button @click="openAddModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add Manufacturer</span>
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
            placeholder="Search manufacturers, country..."
            class="form-input pl-8"
          />
        </div>
      </div>
      <div class="text-2xs text-slate-500 font-mono">
        Total Manufacturers: <span class="font-bold text-slate-800">{{ filteredManufacturers.length }}</span>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
      <table class="dense-table">
        <thead>
          <tr>
            <th>Manufacturer</th>
            <th>Code</th>
            <th>Country / Region</th>
            <th>Official Website</th>
            <th>Total Components</th>
            <th>Rating / Tier</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 font-sans">
          <tr
            v-for="m in filteredManufacturers"
            :key="m.id"
            @click="inspectManufacturer(m)"
            class="cursor-pointer group"
          >
            <td>
              <div class="flex items-center gap-2.5">
                <div class="w-7 h-7 rounded bg-brand-50 border border-brand-200 flex items-center justify-center text-brand-700 font-bold text-xs shrink-0">
                  {{ m.code.slice(0, 2) }}
                </div>
                <span class="font-bold text-slate-900 group-hover:text-brand-600 transition-colors">{{ m.name }}</span>
              </div>
            </td>
            <td>
              <span class="font-mono text-2xs px-1.5 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200 font-bold">
                {{ m.code }}
              </span>
            </td>
            <td>
              <div class="flex items-center gap-1.5 text-xs text-slate-700">
                <span class="font-mono text-3xs px-1 rounded bg-slate-100 border">{{ m.countryCode }}</span>
                <span>{{ m.country }}</span>
              </div>
            </td>
            <td @click.stop>
              <a
                :href="m.website"
                target="_blank"
                rel="noopener"
                class="text-xs text-brand-600 hover:underline flex items-center gap-1 font-mono"
              >
                <span>{{ m.website.replace('https://', '') }}</span>
                <AppIcon name="external-link" size="xs" />
              </a>
            </td>
            <td>
              <span class="font-mono font-bold text-xs text-slate-800">{{ m.totalComponents }}</span>
            </td>
            <td>
              <span class="text-3xs font-semibold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                {{ m.rating }}
              </span>
            </td>
            <td>
              <StatusBadge :status="m.status" />
            </td>
            <td class="text-right" @click.stop>
              <div class="flex items-center justify-end gap-1">
                <button @click="inspectManufacturer(m)" class="btn-icon" title="View Details">
                  <AppIcon name="eye" size="xs" />
                </button>
                <button @click="openEditModal(m)" class="btn-icon" title="Edit">
                  <AppIcon name="edit" size="xs" />
                </button>
                <button @click="confirmDelete(m)" class="btn-icon text-rose-500 hover:text-rose-700" title="Delete">
                  <AppIcon name="trash" size="xs" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Manufacturer Detail Drawer -->
    <BaseDrawer
      :show="detailDrawerOpen"
      :title="inspectedMfg ? inspectedMfg.name : ''"
      :subtitle="inspectedMfg ? `Manufacturer Code: ${inspectedMfg.code}` : ''"
      width="lg"
      @close="detailDrawerOpen = false"
    >
      <div v-if="inspectedMfg" class="space-y-4 text-xs">
        <div class="flex items-center gap-2">
          <StatusBadge :status="inspectedMfg.status" />
          <span class="text-3xs font-mono px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200 font-bold">
            {{ inspectedMfg.rating }}
          </span>
        </div>

        <div class="p-3 bg-slate-50 border border-slate-200 rounded-md">
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Description</div>
          <p class="text-slate-700 mt-1 leading-relaxed whitespace-pre-line">{{ inspectedMfg.description }}</p>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="text-3xs text-slate-400">Headquarters</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ inspectedMfg.country }}</div>
          </div>
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="text-3xs text-slate-400">Direct Contact</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ inspectedMfg.contactPerson || 'Sales Support' }}</div>
            <div class="text-3xs text-slate-500 font-mono">{{ inspectedMfg.contactEmail }}</div>
          </div>
        </div>

        <!-- Related Components List -->
        <div>
          <h4 class="font-bold text-slate-800 uppercase text-3xs tracking-wider mb-2">Registered Components from {{ inspectedMfg.name }}</h4>
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

    <!-- Create / Edit Manufacturer Modal -->
    <BaseModal :show="modalOpen" :title="isEdit ? 'Edit Manufacturer' : 'Add Manufacturer'" maxWidth="md" @close="modalOpen = false">
      <form @submit.prevent="saveModal" class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Manufacturer Name *</label>
          <input v-model="modalForm.name" type="text" required placeholder="e.g. STMicroelectronics" class="form-input" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Code *</label>
            <input v-model="modalForm.code" type="text" required placeholder="e.g. STM" class="form-input font-mono uppercase" />
          </div>
          <div>
            <label class="form-label">Country Code</label>
            <input v-model="modalForm.countryCode" type="text" placeholder="e.g. CH, US, DE" class="form-input font-mono uppercase" />
          </div>
        </div>
        <div>
          <label class="form-label">Headquarters Country *</label>
          <input v-model="modalForm.country" type="text" required placeholder="e.g. Switzerland" class="form-input" />
        </div>
        <div>
          <label class="form-label">Official Website</label>
          <input v-model="modalForm.website" type="url" placeholder="https://www.st.com" class="form-input" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Contact Person</label>
            <input v-model="modalForm.contactPerson" type="text" placeholder="Contact Lead" class="form-input" />
          </div>
          <div>
            <label class="form-label">Status</label>
            <select v-model="modalForm.status" class="form-select">
              <option value="Approved">Approved</option>
              <option value="Under Review">Under Review</option>
            </select>
          </div>
        </div>
        <div>
          <label class="form-label">Description</label>
          <textarea v-model="modalForm.description" rows="2" placeholder="Overview of semiconductor portfolio..." class="form-textarea"></textarea>
        </div>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="modalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Save Manufacturer</button>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Manufacturer"
      :message="`Are you sure you want to delete '${mfgToDelete ? mfgToDelete.name : ''}'?`"
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
  name: 'ManufacturerListView',
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
      inspectedMfg: null,
      modalOpen: false,
      isEdit: false,
      modalForm: {
        id: null,
        name: '',
        code: '',
        country: '',
        countryCode: 'US',
        website: 'https://',
        rating: 'Approved',
        status: 'Approved',
        contactPerson: '',
        contactEmail: '',
        description: ''
      },
      deleteModalOpen: false,
      mfgToDelete: null
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    filteredManufacturers() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.store.manufacturers.filter(m => {
        return !q || m.name.toLowerCase().includes(q) || m.code.toLowerCase().includes(q) || m.country.toLowerCase().includes(q);
      });
    },
    relatedComponents() {
      if (!this.inspectedMfg) return [];
      return this.store.components.filter(c => c.manufacturer === this.inspectedMfg.name);
    }
  },
  methods: {
    inspectManufacturer(m) {
      this.inspectedMfg = m;
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
        website: 'https://',
        rating: 'Approved',
        status: 'Approved',
        contactPerson: '',
        contactEmail: '',
        description: ''
      };
      this.modalOpen = true;
    },
    openEditModal(m) {
      this.isEdit = true;
      this.modalForm = JSON.parse(JSON.stringify(m));
      this.modalOpen = true;
    },
    saveModal() {
      this.store.saveManufacturer(this.modalForm);
      this.modalOpen = false;
    },
    confirmDelete(m) {
      this.mfgToDelete = m;
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.mfgToDelete) {
        this.store.deleteManufacturer(this.mfgToDelete.id);
        this.deleteModalOpen = false;
        this.mfgToDelete = null;
        if (this.inspectedMfg && this.inspectedMfg.id === this.mfgToDelete?.id) {
          this.detailDrawerOpen = false;
        }
      }
    }
  }
};
</script>
