<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm pointer-events-none">
    <transition
      enter-active-class="transform ease-out duration-200 transition"
      enter-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-150"
      leave-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="store.toast.show"
        :class="[
          'pointer-events-auto w-full rounded-lg p-3.5 shadow-dropdown border flex items-start gap-3 bg-white',
          borderClass
        ]"
      >
        <div :class="['w-6 h-6 rounded-full flex items-center justify-center shrink-0 mt-0.5', iconBgClass]">
          <AppIcon :name="iconName" size="xs" :class="iconColorClass" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-xs font-bold text-slate-900">{{ store.toast.title || defaultTitle }}</div>
          <div class="text-2xs text-slate-600 mt-0.5 leading-snug">{{ store.toast.message }}</div>
        </div>
        <button
          @click="store.toast.show = false"
          class="text-slate-400 hover:text-slate-700 p-0.5 rounded transition-colors"
        >
          <AppIcon name="x" size="xs" />
        </button>
      </div>
    </transition>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import AppIcon from './AppIcon.vue';

export default {
  name: 'ToastContainer',
  components: { AppIcon },
  computed: {
    store() {
      return useMasterStore();
    },
    defaultTitle() {
      switch (this.store.toast.type) {
        case 'success': return 'Action Completed';
        case 'error': return 'Error Occurred';
        case 'warning': return 'Warning';
        case 'info': return 'Notification';
        default: return 'Notification';
      }
    },
    borderClass() {
      switch (this.store.toast.type) {
        case 'success': return 'border-emerald-200 ring-1 ring-emerald-100';
        case 'error': return 'border-rose-200 ring-1 ring-rose-100';
        case 'warning': return 'border-amber-200 ring-1 ring-amber-100';
        case 'info': return 'border-brand-200 ring-1 ring-brand-100';
        default: return 'border-slate-200';
      }
    },
    iconBgClass() {
      switch (this.store.toast.type) {
        case 'success': return 'bg-emerald-50 text-emerald-600';
        case 'error': return 'bg-rose-50 text-rose-600';
        case 'warning': return 'bg-amber-50 text-amber-600';
        case 'info': return 'bg-brand-50 text-brand-600';
        default: return 'bg-slate-100 text-slate-600';
      }
    },
    iconColorClass() {
      switch (this.store.toast.type) {
        case 'success': return 'text-emerald-600';
        case 'error': return 'text-rose-600';
        case 'warning': return 'text-amber-600';
        case 'info': return 'text-brand-600';
        default: return 'text-slate-600';
      }
    },
    iconName() {
      switch (this.store.toast.type) {
        case 'success': return 'check-circle';
        case 'error': return 'x-circle';
        case 'warning': return 'alert-circle';
        case 'info': return 'info';
        default: return 'info';
      }
    }
  }
};
</script>
