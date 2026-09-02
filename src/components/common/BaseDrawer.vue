<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-hidden">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs transition-opacity" @click="close"></div>

    <div class="fixed inset-y-0 right-0 max-w-full flex pl-10">
      <div
        :class="[
          'w-screen bg-white shadow-2xl border-l border-slate-200 flex flex-col transition-all duration-300',
          widthClass
        ]"
      >
        <!-- Header -->
        <div class="px-4 py-3 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
          <div class="flex items-center gap-2 min-w-0">
            <slot name="icon" />
            <div>
              <h3 class="text-sm font-bold text-slate-900 truncate">{{ title }}</h3>
              <p v-if="subtitle" class="text-2xs text-slate-500">{{ subtitle }}</p>
            </div>
          </div>
          <button
            @click="close"
            class="p-1 rounded text-slate-400 hover:text-slate-700 hover:bg-slate-200/60 focus:outline-none transition-colors"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Body -->
        <div class="flex-1 overflow-y-auto p-4">
          <slot />
        </div>

        <!-- Footer -->
        <div v-if="$slots.footer" class="px-4 py-3 border-t border-slate-200 bg-slate-50/70 flex items-center justify-end gap-2">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AppIcon from './AppIcon.vue';

export default {
  name: 'BaseDrawer',
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
    subtitle: {
      type: String,
      default: ''
    },
    width: {
      type: String,
      default: 'md' // 'sm', 'md', 'lg', 'xl', '2xl'
    }
  },
  computed: {
    widthClass() {
      switch (this.width) {
        case 'sm': return 'max-w-sm';
        case 'md': return 'max-w-md';
        case 'lg': return 'max-w-lg';
        case 'xl': return 'max-w-xl';
        case '2xl': return 'max-w-2xl';
        default: return 'max-w-md';
      }
    }
  },
  methods: {
    close() {
      this.$emit('close');
    }
  }
};
</script>
