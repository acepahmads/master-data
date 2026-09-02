<template>
  <BaseModal :show="isOpen" title="System & Environment Settings" maxWidth="lg" @close="close">
    <template #icon>
      <div class="w-6 h-6 rounded bg-slate-100 text-slate-700 flex items-center justify-center">
        <AppIcon name="settings" size="xs" />
      </div>
    </template>

    <div class="space-y-4 text-xs">
      <div>
        <label class="form-label">Active Workspace</label>
        <input type="text" class="form-input bg-slate-50" readonly value="Product & Master Data Control Center (Production Environment)" />
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="form-label">System Platform</label>
          <input type="text" class="form-input bg-slate-50 font-mono" readonly value="Phases 1–6 Complete (100%)" />
        </div>
        <div>
          <label class="form-label">Release Version</label>
          <input type="text" class="form-input bg-slate-50 font-mono" readonly value="v1.0.0" />
        </div>
      </div>

      <div class="border-t border-slate-200 pt-3">
        <div class="font-bold text-slate-800 text-xs mb-2">Display & Density Preferences</div>
        <div class="space-y-2">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" checked class="rounded text-brand-600 focus:ring-brand-500" />
            <span class="text-xs text-slate-700">Compact High-Density Table Layout (Optimized for 1440px/1280px)</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" checked class="rounded text-brand-600 focus:ring-brand-500" />
            <span class="text-xs text-slate-700">Show Quick Action Dropdowns on Hover</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" checked class="rounded text-brand-600 focus:ring-brand-500" />
            <span class="text-xs text-slate-700">Real-time Engineering Audit Trail Logging</span>
          </label>
        </div>
      </div>
    </div>

    <template #footer>
      <button type="button" class="btn btn-primary" @click="save">
        Save Preferences
      </button>
    </template>
  </BaseModal>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import BaseModal from './BaseModal.vue';
import AppIcon from './AppIcon.vue';

export default {
  name: 'SettingsModal',
  components: { BaseModal, AppIcon },
  computed: {
    store() {
      return useMasterStore();
    },
    isOpen() {
      return this.store.settingsModalOpen;
    }
  },
  methods: {
    close() {
      this.store.settingsModalOpen = false;
    },
    save() {
      this.store.showToast('Settings Updated', 'System preferences have been saved.', 'success');
      this.close();
    }
  }
};
</script>
