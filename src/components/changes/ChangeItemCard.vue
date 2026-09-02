<template>
  <div class="p-4 rounded-xl border border-slate-800 bg-slate-900/60 hover:border-slate-700 transition-colors">
    <div class="flex items-start justify-between gap-3 mb-2">
      <div class="flex items-center gap-2">
        <span class="px-2 py-0.5 rounded text-[11px] font-bold uppercase tracking-wider bg-slate-800 text-cyan-400 border border-cyan-900/50">
          {{ item.changeType }}
        </span>
        <span v-if="item.affectedReference" class="font-mono text-xs font-semibold text-amber-300 bg-amber-950/40 px-2 py-0.5 rounded border border-amber-800/40">
          Ref: {{ item.affectedReference }}
        </span>
      </div>
      <button
        v-if="canDelete"
        @click="$emit('delete', item.id)"
        class="text-slate-500 hover:text-rose-400 p-1 transition-colors"
        title="Remove Change Item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>

    <h4 class="text-sm font-semibold text-slate-200 mb-1">
      {{ item.title }}
    </h4>

    <p v-if="item.description" class="text-xs text-slate-400 mb-3 leading-relaxed">
      {{ item.description }}
    </p>

    <!-- Component Comparison Sub-Block -->
    <div
      v-if="item.affectedComponent || item.currentValue || item.proposedValue"
      class="grid grid-cols-2 gap-2 p-2.5 rounded-lg bg-slate-950/60 border border-slate-800/80 text-xs font-mono"
    >
      <div>
        <div class="text-[10px] text-slate-500 uppercase">Current State</div>
        <div class="text-slate-300 truncate">{{ item.currentValue || 'Original BOM part' }}</div>
      </div>
      <div>
        <div class="text-[10px] text-cyan-400 uppercase">Proposed Delta</div>
        <div class="text-cyan-300 truncate font-semibold">
          {{ item.affectedComponent ? item.affectedComponent.name : item.proposedValue || 'New Specification' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChangeItemCard',
  props: {
    item: {
      type: Object,
      required: true
    },
    canDelete: {
      type: Boolean,
      default: false
    }
  }
};
</script>
