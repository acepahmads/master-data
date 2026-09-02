<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Units of Measurement (UoM)" subtitle="Standardized engineering units, packaging types, and physical quantities">
      <template #actions>
        <button @click="openAddModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add Unit</span>
        </button>
      </template>
    </PageHeader>

    <!-- Filter & Categories Summary -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <div
        v-for="cat in uomCategories"
        :key="cat.name"
        @click="selectedCategory = selectedCategory === cat.name ? '' : cat.name"
        :class="[
          'p-3 rounded-lg border transition-all cursor-pointer bg-white flex items-center justify-between',
          selectedCategory === cat.name ? 'border-brand-500 ring-1 ring-brand-400 shadow-sm' : 'border-slate-200 hover:border-slate-300'
        ]"
      >
        <div>
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">{{ cat.name }}</div>
          <div class="font-bold text-slate-900 text-sm mt-0.5">{{ cat.count }} Units</div>
        </div>
        <div class="w-7 h-7 rounded bg-slate-100 flex items-center justify-center text-slate-600">
          <AppIcon :name="cat.icon" size="xs" />
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
      <table class="dense-table">
        <thead>
          <tr>
            <th>Unit Code</th>
            <th>Unit Name</th>
            <th>Symbol</th>
            <th>Quantity Category</th>
            <th>Description</th>
            <th>Default Flag</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 font-sans">
          <tr v-for="u in filteredUnits" :key="u.id" class="hover:bg-slate-50">
            <td>
              <span class="font-mono font-bold text-xs text-brand-700 bg-brand-50 px-2 py-0.5 rounded border border-brand-200">
                {{ u.code }}
              </span>
            </td>
            <td>
              <span class="font-bold text-slate-800">{{ u.name }}</span>
            </td>
            <td>
              <span class="font-mono text-slate-600 bg-slate-100 px-1.5 py-0.5 rounded text-2xs">{{ u.symbol }}</span>
            </td>
            <td>
              <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 border border-slate-200">
                {{ u.category }}
              </span>
            </td>
            <td class="text-slate-500 max-w-xs truncate">{{ u.description }}</td>
            <td>
              <span v-if="u.isDefault" class="text-3xs font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200">
                PRIMARY DEFAULT
              </span>
              <span v-else class="text-3xs text-slate-400 font-mono">—</span>
            </td>
            <td>
              <StatusBadge :status="u.status" />
            </td>
            <td class="text-right">
              <div class="flex items-center justify-end gap-1">
                <button @click="openEditModal(u)" class="btn-icon" title="Edit Unit">
                  <AppIcon name="edit" size="xs" />
                </button>
                <button @click="confirmDelete(u)" class="btn-icon text-rose-500 hover:text-rose-700" title="Delete Unit">
                  <AppIcon name="trash" size="xs" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Unit Modal -->
    <BaseModal :show="modalOpen" :title="isEdit ? 'Edit Unit of Measurement' : 'Add Unit of Measurement'" maxWidth="md" @close="modalOpen = false">
      <form @submit.prevent="saveModal" class="space-y-3.5 text-xs">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Unit Code *</label>
            <input v-model="modalForm.code" type="text" required placeholder="e.g. REEL" class="form-input font-mono uppercase" />
          </div>
          <div>
            <label class="form-label">Symbol *</label>
            <input v-model="modalForm.symbol" type="text" required placeholder="e.g. rl" class="form-input font-mono" />
          </div>
        </div>
        <div>
          <label class="form-label">Unit Name *</label>
          <input v-model="modalForm.name" type="text" required placeholder="e.g. SMD Tape Reel" class="form-input" />
        </div>
        <div>
          <label class="form-label">Category *</label>
          <select v-model="modalForm.category" required class="form-select">
            <option value="Count">Count / Discrete</option>
            <option value="Packaging">Packaging Container</option>
            <option value="Length">Length / Dimension</option>
            <option value="Weight">Weight / Mass</option>
          </select>
        </div>
        <div>
          <label class="form-label">Description</label>
          <textarea v-model="modalForm.description" rows="2" placeholder="Standard engineering definition and container sizing..." class="form-textarea"></textarea>
        </div>
        <div class="flex items-center gap-2 pt-1">
          <input type="checkbox" v-model="modalForm.isDefault" id="isDef" class="rounded text-brand-600 focus:ring-brand-500" />
          <label for="isDef" class="cursor-pointer text-slate-700">Set as default unit for new components</label>
        </div>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="modalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Save Unit</button>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Unit"
      :message="`Are you sure you want to delete unit '${unitToDelete ? unitToDelete.name : ''}' (${unitToDelete ? unitToDelete.code : ''})?`"
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
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'UnitListView',
  components: {
    PageHeader,
    StatusBadge,
    ConfirmDialog,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      selectedCategory: '',
      modalOpen: false,
      isEdit: false,
      modalForm: {
        id: null,
        code: '',
        name: '',
        symbol: '',
        category: 'Count',
        description: '',
        status: 'Active',
        isDefault: false
      },
      deleteModalOpen: false,
      unitToDelete: null
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    uomCategories() {
      const units = this.store.units;
      return [
        { name: 'Count', count: units.filter(u => u.category === 'Count').length, icon: 'box' },
        { name: 'Packaging', count: units.filter(u => u.category === 'Packaging').length, icon: 'layers' },
        { name: 'Length', count: units.filter(u => u.category === 'Length').length, icon: 'unit' },
        { name: 'Weight', count: units.filter(u => u.category === 'Weight').length, icon: 'unit' }
      ];
    },
    filteredUnits() {
      if (!this.selectedCategory) return this.store.units;
      return this.store.units.filter(u => u.category === this.selectedCategory);
    }
  },
  methods: {
    openAddModal() {
      this.isEdit = false;
      this.modalForm = {
        id: null,
        code: '',
        name: '',
        symbol: '',
        category: 'Count',
        description: '',
        status: 'Active',
        isDefault: false
      };
      this.modalOpen = true;
    },
    openEditModal(u) {
      this.isEdit = true;
      this.modalForm = JSON.parse(JSON.stringify(u));
      this.modalOpen = true;
    },
    saveModal() {
      this.store.saveUnit(this.modalForm);
      this.modalOpen = false;
    },
    confirmDelete(u) {
      this.unitToDelete = u;
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.unitToDelete) {
        this.store.deleteUnit(this.unitToDelete.id);
        this.deleteModalOpen = false;
        this.unitToDelete = null;
      }
    }
  }
};
</script>
