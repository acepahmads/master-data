<template>
  <div class="py-12 px-6 text-center border-2 border-dashed border-slate-200/90 rounded-xl bg-slate-50/50 my-2">
    <div class="w-12 h-12 rounded-xl bg-white border border-slate-200 shadow-xs flex items-center justify-center mx-auto text-slate-400 mb-3 group-hover:scale-105 transition-transform">
      <AppIcon :name="icon" size="md" class="text-slate-500" />
    </div>
    <h3 class="text-sm font-bold text-slate-800 tracking-tight">{{ title }}</h3>
    <p class="text-xs text-slate-500 mt-1 max-w-md mx-auto leading-relaxed font-medium">{{ description }}</p>
    <div v-if="$slots.action" class="mt-5 flex items-center justify-center gap-2">
      <slot name="action" />
    </div>
  </div>
</template>

<script>
import AppIcon from './AppIcon.vue';

export default {
  name: 'EmptyState',
  components: { AppIcon },
  props: {
    title: {
      type: String,
      default: 'No records found'
    },
    description: {
      type: String,
      default: 'No items match your current filter criteria or none have been registered yet.'
    },
    icon: {
      type: String,
      default: 'search'
    }
  }
};
</script>
