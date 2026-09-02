<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Roles & Permissions" subtitle="Role-Based Access Control (RBAC) and granular security policies">
      <template #actions>
        <button @click="openAddRoleModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add New Role</span>
        </button>
      </template>
    </PageHeader>

    <!-- 2-Column Split: Role List on Left, Matrix on Right -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
      <!-- LEFT 4 COLS: ROLE LIST -->
      <div class="lg:col-span-4 bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
        <div class="p-3.5 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <AppIcon name="shield" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Assigned Roles</h2>
          </div>
          <span class="text-3xs font-mono text-slate-500">{{ store.roles.length }} Roles</span>
        </div>

        <div class="p-2 space-y-1">
          <div
            v-for="role in store.roles"
            :key="role.id"
            @click="selectRole(role)"
            :class="[
              'p-3 rounded-lg cursor-pointer transition-all border flex items-center justify-between group',
              selectedRole && selectedRole.id === role.id
                ? 'bg-brand-50/70 border-brand-300 ring-1 ring-brand-200'
                : 'border-transparent hover:bg-slate-50 hover:border-slate-200'
            ]"
          >
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span :class="['font-bold text-xs', selectedRole && selectedRole.id === role.id ? 'text-brand-900' : 'text-slate-800']">
                  {{ role.name }}
                </span>
                <span v-if="role.isSystem" class="text-3xs font-mono px-1 rounded bg-slate-100 text-slate-500 border border-slate-200">
                  SYSTEM
                </span>
              </div>
              <p class="text-3xs text-slate-500 mt-0.5 truncate">{{ role.description }}</p>
            </div>
            <AppIcon
              name="chevron-right"
              size="xs"
              :class="[
                'shrink-0 transition-transform',
                selectedRole && selectedRole.id === role.id ? 'text-brand-600 translate-x-0.5' : 'text-slate-300 group-hover:text-slate-500'
              ]"
            />
          </div>
        </div>
      </div>

      <!-- RIGHT 8 COLS: ROLE DETAILS & PERMISSION MATRIX -->
      <div class="lg:col-span-8 space-y-4">
        <div v-if="selectedRole" class="bg-white rounded-lg border border-slate-200 shadow-card p-5 space-y-5">
          <!-- Role Header -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-200 gap-3">
            <div>
              <div class="flex items-center gap-2">
                <h2 class="text-lg font-bold text-slate-900 tracking-tight">{{ selectedRole.name }}</h2>
                <span class="font-mono text-3xs px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200 font-bold">
                  {{ selectedRole.code }}
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-1">{{ selectedRole.description }}</p>
            </div>

            <div class="flex items-center gap-2">
              <button @click="resetMatrix" class="btn btn-secondary text-2xs">
                Reset
              </button>
              <button @click="savePermissions" class="btn btn-primary text-2xs">
                <AppIcon name="check" size="xs" class="mr-1" />
                <span>Save Permissions</span>
              </button>
            </div>
          </div>

          <!-- PERMISSION MATRIX -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
                <AppIcon name="lock" size="xs" class="text-slate-400" />
                <span>Module Permission Matrix</span>
              </h3>
              <div class="flex items-center gap-2 text-2xs">
                <button @click="globalSelectAll" class="text-brand-600 hover:underline font-semibold">Select All</button>
                <span>•</span>
                <button @click="globalClearAll" class="text-slate-500 hover:underline">Clear All</button>
              </div>
            </div>

            <!-- MODULE 1: PRODUCTS -->
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <div class="px-3.5 py-2 bg-slate-50/80 border-b border-slate-200 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <AppIcon name="product" size="xs" class="text-brand-600" />
                  <span class="font-bold text-xs text-slate-800">PRODUCTS</span>
                </div>
                <div class="flex items-center gap-2 text-3xs">
                  <button @click="toggleModuleAll('products', true)" class="text-brand-600 hover:underline font-medium">Select All</button>
                  <button @click="toggleModuleAll('products', false)" class="text-slate-400 hover:underline">Clear</button>
                </div>
              </div>
              <div class="p-3 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3 text-xs">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.view" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>View</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.create" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Create</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.edit" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Edit</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.delete" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Delete</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.export" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Export</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.products.approve" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Approve</span>
                </label>
              </div>
            </div>

            <!-- MODULE 2: COMPONENTS -->
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <div class="px-3.5 py-2 bg-slate-50/80 border-b border-slate-200 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <AppIcon name="component" size="xs" class="text-emerald-600" />
                  <span class="font-bold text-xs text-slate-800">COMPONENTS & SPECIFICATIONS</span>
                </div>
                <div class="flex items-center gap-2 text-3xs">
                  <button @click="toggleModuleAll('components', true)" class="text-brand-600 hover:underline font-medium">Select All</button>
                  <button @click="toggleModuleAll('components', false)" class="text-slate-400 hover:underline">Clear</button>
                </div>
              </div>
              <div class="p-3 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3 text-xs">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.view" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>View</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.create" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Create</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.edit" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Edit</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.delete" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Delete</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.export" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Export</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.components.approve" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Approve</span>
                </label>
              </div>
            </div>

            <!-- MODULE 3: MASTER DATA (CATEGORIES, MFGS, SUPPLIERS, UNITS) -->
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <div class="px-3.5 py-2 bg-slate-50/80 border-b border-slate-200 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <AppIcon name="tree" size="xs" class="text-purple-600" />
                  <span class="font-bold text-xs text-slate-800">MASTER DATA ENTITIES (Categories, Manufacturers, Suppliers, UoM)</span>
                </div>
              </div>
              <div class="p-3 divide-y divide-slate-100 text-xs">
                <div class="py-2 grid grid-cols-5 gap-3 items-center">
                  <span class="font-semibold text-slate-700">Categories</span>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.categories.view" class="rounded text-brand-600" /><span>View</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.categories.create" class="rounded text-brand-600" /><span>Create</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.categories.edit" class="rounded text-brand-600" /><span>Edit</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.categories.delete" class="rounded text-brand-600" /><span>Delete</span></label>
                </div>
                <div class="py-2 grid grid-cols-5 gap-3 items-center">
                  <span class="font-semibold text-slate-700">Manufacturers</span>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.manufacturers.view" class="rounded text-brand-600" /><span>View</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.manufacturers.create" class="rounded text-brand-600" /><span>Create</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.manufacturers.edit" class="rounded text-brand-600" /><span>Edit</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.manufacturers.delete" class="rounded text-brand-600" /><span>Delete</span></label>
                </div>
                <div class="py-2 grid grid-cols-5 gap-3 items-center">
                  <span class="font-semibold text-slate-700">Suppliers</span>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.suppliers.view" class="rounded text-brand-600" /><span>View</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.suppliers.create" class="rounded text-brand-600" /><span>Create</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.suppliers.edit" class="rounded text-brand-600" /><span>Edit</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.suppliers.delete" class="rounded text-brand-600" /><span>Delete</span></label>
                </div>
                <div class="py-2 grid grid-cols-5 gap-3 items-center">
                  <span class="font-semibold text-slate-700">Units (UoM)</span>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.units.view" class="rounded text-brand-600" /><span>View</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.units.create" class="rounded text-brand-600" /><span>Create</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.units.edit" class="rounded text-brand-600" /><span>Edit</span></label>
                  <label class="flex items-center gap-1.5 cursor-pointer"><input type="checkbox" v-model="matrix.units.delete" class="rounded text-brand-600" /><span>Delete</span></label>
                </div>
              </div>
            </div>

            <!-- MODULE 4: ADMINISTRATION -->
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <div class="px-3.5 py-2 bg-slate-50/80 border-b border-slate-200 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <AppIcon name="settings" size="xs" class="text-slate-700" />
                  <span class="font-bold text-xs text-slate-800">ADMINISTRATION & SECURITY POLICIES</span>
                </div>
              </div>
              <div class="p-3 grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.administration.manageUsers" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Manage Users</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.administration.manageRoles" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>Manage Roles</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.administration.auditLogs" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>View Audit Logs</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="checkbox" v-model="matrix.administration.systemSettings" class="rounded text-brand-600 focus:ring-brand-500" />
                  <span>System Settings</span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Role Modal -->
    <BaseModal :show="addRoleModalOpen" title="Create Custom Engineering Role" maxWidth="md" @close="addRoleModalOpen = false">
      <form @submit.prevent="saveNewRole" class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Role Name *</label>
          <input v-model="newRoleForm.name" type="text" required placeholder="e.g. Compliance Lead" class="form-input" />
        </div>
        <div>
          <label class="form-label">Role Code *</label>
          <input v-model="newRoleForm.code" type="text" required placeholder="e.g. COMPLIANCE_LEAD" class="form-input font-mono uppercase" />
        </div>
        <div>
          <label class="form-label">Description</label>
          <textarea v-model="newRoleForm.description" rows="3" placeholder="Access scope and responsibilities..." class="form-textarea"></textarea>
        </div>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="addRoleModalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Create Role</button>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'RoleListView',
  components: {
    PageHeader,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      selectedRole: null,
      matrix: {
        products: {},
        components: {},
        categories: {},
        manufacturers: {},
        suppliers: {},
        units: {},
        administration: {}
      },
      addRoleModalOpen: false,
      newRoleForm: {
        name: '',
        code: '',
        description: ''
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    }
  },
  created() {
    if (this.store.roles.length > 0) {
      this.selectRole(this.store.roles[0]);
    }
  },
  methods: {
    selectRole(role) {
      this.selectedRole = role;
      this.matrix = JSON.parse(JSON.stringify(role.permissions));
    },
    resetMatrix() {
      if (this.selectedRole) {
        this.matrix = JSON.parse(JSON.stringify(this.selectedRole.permissions));
      }
    },
    toggleModuleAll(moduleKey, value) {
      if (this.matrix[moduleKey]) {
        Object.keys(this.matrix[moduleKey]).forEach(k => {
          this.$set(this.matrix[moduleKey], k, value);
        });
      }
    },
    globalSelectAll() {
      Object.keys(this.matrix).forEach(mod => {
        this.toggleModuleAll(mod, true);
      });
    },
    globalClearAll() {
      Object.keys(this.matrix).forEach(mod => {
        this.toggleModuleAll(mod, false);
      });
    },
    savePermissions() {
      if (this.selectedRole) {
        this.store.saveRolePermissions(this.selectedRole.id, this.matrix);
      }
    },
    openAddRoleModal() {
      this.newRoleForm = {
        name: '',
        code: '',
        description: ''
      };
      this.addRoleModalOpen = true;
    },
    saveNewRole() {
      this.store.saveRole(this.newRoleForm);
      this.addRoleModalOpen = false;
      this.selectRole(this.store.roles[this.store.roles.length - 1]);
    }
  }
};
</script>
