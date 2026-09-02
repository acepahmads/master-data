<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity" @click="handleBackdropClick"></div>

    <div class="flex min-h-full items-center justify-center p-4 text-center">
      <div
        :class="[
          'relative transform overflow-hidden rounded-lg bg-white text-left shadow-modal transition-all border border-slate-200 w-full',
          maxWidthClass
        ]"
      >
        <!-- Modal Header -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-slate-200 bg-slate-50/70">
          <div class="flex items-center gap-2 min-w-0">
            <slot name="icon" />
            <h3 class="text-sm font-bold text-slate-900 truncate">{{ title }}</h3>
          </div>
          <button
            @click="close"
            class="p-1 rounded text-slate-400 hover:text-slate-700 hover:bg-slate-200/60 focus:outline-none transition-colors"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-4 max-h-[calc(85vh-110px)] overflow-y-auto">
          <slot />
        </div>

        <!-- Modal Footer -->
        <div v-if="$slots.footer" class="flex items-center justify-end gap-2 px-4 py-3 border-t border-slate-200 bg-slate-50/50">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AppIcon from './AppIcon.vue';

export default {
  name: 'BaseModal',
  components: { AppIcon },
  props: {
    show: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    maxWidth: {
      type: String,
      default: 'md' // 'sm', 'md', 'lg', 'xl', '2xl', '3xl', '4xl'
    },
    closeOnBackdrop: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    maxWidthClass() {
      switch (this.maxWidth) {
        case 'sm': return 'max-w-sm';
        case 'md': return 'max-w-md';
        case 'lg': return 'max-w-lg';
        case 'xl': return 'max-w-xl';
        case '2xl': return 'max-w-2xl';
        case '3xl': return 'max-w-3xl';
        case '4xl': return 'max-w-4xl';
        case '5xl': return 'max-w-5xl';
        default: return 'max-w-md';
      }
    }
  },
  methods: {
    close() {
      this.$emit('close');
    },
    handleBackdropClick() {
      if (this.closeOnBackdrop) {
        this.close();
      }
    }
  }
};
</script>
