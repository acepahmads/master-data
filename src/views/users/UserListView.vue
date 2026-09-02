<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="User Management" subtitle="Engineering personnel, security credentials, and department assignments">
      <template #actions>
        <button @click="openAddModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add User</span>
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
            placeholder="Search users, email, department..."
            class="form-input pl-8"
          />
        </div>
      </div>
      <div class="text-2xs text-slate-500 font-mono">
        Active Accounts: <span class="font-bold text-slate-800">{{ filteredUsers.length }}</span>
      </div>
    </div>

    <!-- User Data Table -->
    <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
      <table class="dense-table">
        <thead>
          <tr>
            <th>User</th>
            <th>Email</th>
            <th>Department</th>
            <th>Role</th>
            <th>2FA Security</th>
            <th>Last Login</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 font-sans">
          <tr v-for="u in filteredUsers" :key="u.id" class="hover:bg-slate-50">
            <td>
              <div class="flex items-center gap-2.5">
                <img :src="u.avatar" class="w-8 h-8 rounded-full object-cover border border-slate-200 shrink-0" />
                <div>
                  <div class="font-bold text-slate-900">{{ u.name }}</div>
                  <div class="text-3xs text-slate-400 font-mono">@{{ u.username }}</div>
                </div>
              </div>
            </td>
            <td>
              <span class="font-mono text-xs text-slate-700">{{ u.email }}</span>
            </td>
            <td>
              <span class="text-slate-700 text-xs font-medium">{{ u.department }}</span>
            </td>
            <td>
              <span class="px-2 py-0.5 rounded text-3xs font-semibold bg-brand-50 text-brand-700 border border-brand-200">
                {{ getRoleName(u.role) }}
              </span>
            </td>
            <td>
              <span v-if="u.twoFactor" class="flex items-center gap-1 text-3xs font-semibold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200 w-fit">
                <AppIcon name="lock" size="xs" />
                <span>Enabled</span>
              </span>
              <span v-else class="text-3xs text-slate-400 font-mono">Disabled</span>
            </td>
            <td>
              <span class="font-mono text-2xs text-slate-500">{{ u.lastLogin }}</span>
            </td>
            <td>
              <StatusBadge :status="u.status" />
            </td>
            <td class="text-right">
              <div class="flex items-center justify-end gap-1">
                <button @click="openEditModal(u)" class="btn-icon" title="Edit User">
                  <AppIcon name="edit" size="xs" />
                </button>
                <button
                  @click="store.toggleUserStatus(u.id)"
                  class="btn-icon"
                  :title="u.status === 'Active' ? 'Suspend Account' : 'Activate Account'"
                >
                  <AppIcon :name="u.status === 'Active' ? 'x' : 'check'" size="xs" />
                </button>
                <button @click="resetPassword(u)" class="btn-icon text-amber-600" title="Reset Password">
                  <AppIcon name="refresh" size="xs" />
                </button>
                <button @click="confirmDelete(u)" class="btn-icon text-rose-500 hover:text-rose-700" title="Delete User">
                  <AppIcon name="trash" size="xs" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit User Modal -->
    <BaseModal :show="modalOpen" :title="isEdit ? 'Edit User Account' : 'Create User Account'" maxWidth="md" @close="modalOpen = false">
      <form @submit.prevent="saveModal" class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Full Name *</label>
          <input v-model="modalForm.name" type="text" required placeholder="e.g. Alex Chen" class="form-input" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Username *</label>
            <input v-model="modalForm.username" type="text" required placeholder="e.g. achen" class="form-input font-mono" />
          </div>
          <div>
            <label class="form-label">Work Email *</label>
            <input v-model="modalForm.email" type="email" required placeholder="alex.chen@iotcontrol.io" class="form-input font-mono" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label">Department *</label>
            <select v-model="modalForm.department" required class="form-select">
              <option value="Hardware R&D">Hardware R&D</option>
              <option value="Engineering Leadership">Engineering Leadership</option>
              <option value="Embedded Systems">Embedded Systems</option>
              <option value="Supply Chain & Sourcing">Supply Chain & Sourcing</option>
              <option value="Quality & Compliance">Quality & Compliance</option>
            </select>
          </div>
          <div>
            <label class="form-label">RBAC Role *</label>
            <select v-model="modalForm.role" required class="form-select">
              <option v-for="r in store.roles" :key="r.id" :value="r.name">
                {{ r.name }}
              </option>
            </select>
          </div>
        </div>
        <div class="flex items-center gap-2 pt-1">
          <input type="checkbox" v-model="modalForm.twoFactor" id="twoFa" class="rounded text-brand-600 focus:ring-brand-500" />
          <label for="twoFa" class="cursor-pointer text-slate-700">Enforce Two-Factor Authentication (2FA)</label>
        </div>
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="modalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Save User</button>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete User Account"
      :message="`Are you sure you want to delete user account for '${userToDelete ? userToDelete.name : ''}'?`"
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
  name: 'UserListView',
  components: {
    PageHeader,
    StatusBadge,
    ConfirmDialog,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      modalOpen: false,
      isEdit: false,
      modalForm: {
        id: null,
        name: '',
        username: '',
        email: '',
        department: 'Hardware R&D',
        role: 'Hardware Engineer',
        status: 'Active',
        twoFactor: true
      },
      deleteModalOpen: false,
      userToDelete: null
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    filteredUsers() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.store.users.filter(u => {
        const roleStr = this.getRoleName(u.role);
        return !q || u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q) || u.department.toLowerCase().includes(q) || roleStr.toLowerCase().includes(q);
      });
    }
  },
  methods: {
    getRoleName(role) {
      if (!role) return 'Unknown Role';
      if (typeof role === 'object' && role !== null) {
        return role.name || 'Unknown Role';
      }
      return role;
    },
    openAddModal() {
      this.isEdit = false;
      this.modalForm = {
        id: null,
        name: '',
        username: '',
        email: '',
        department: 'Hardware R&D',
        role: 'Hardware Engineer',
        status: 'Active',
        twoFactor: true
      };
      this.modalOpen = true;
    },
    openEditModal(u) {
      this.isEdit = true;
      this.modalForm = JSON.parse(JSON.stringify(u));
      this.modalOpen = true;
    },
    saveModal() {
      this.store.saveUser(this.modalForm);
      this.modalOpen = false;
    },
    resetPassword(u) {
      this.store.showToast('Password Reset Email', `Temporary credentials dispatched to ${u.email}`, 'info');
    },
    confirmDelete(u) {
      this.userToDelete = u;
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.userToDelete) {
        this.store.deleteUser(this.userToDelete.id);
        this.deleteModalOpen = false;
        this.userToDelete = null;
      }
    }
  }
};
</script>
