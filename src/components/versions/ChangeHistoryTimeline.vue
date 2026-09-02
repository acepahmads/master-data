<template>
  <div class="space-y-4">
    <div class="relative pl-6 space-y-5 before:absolute before:left-2 before:top-2.5 before:bottom-2.5 before:w-0.5 before:bg-slate-200">
      <div
        v-for="item in items"
        :key="item.id"
        class="relative group text-xs"
      >
        <!-- Timeline status node -->
        <span
          :class="[
            'absolute -left-6 top-1 w-4 h-4 rounded-full ring-4 ring-white shadow-2xs flex items-center justify-center',
            item.badgeColor === 'emerald' ? 'bg-emerald-500 text-white' :
            item.badgeColor === 'blue' ? 'bg-brand-500 text-white' :
            item.badgeColor === 'amber' ? 'bg-amber-500 text-white' :
            item.badgeColor === 'rose' ? 'bg-rose-500 text-white' :
            'bg-slate-400 text-white'
          ]"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-white"></span>
        </span>

        <div class="app-card p-3.5 hover:border-slate-300 transition-all">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1 pb-1.5 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <span class="font-bold text-slate-900">{{ item.actor || 'Alex Chen' }}</span>
              <span class="text-3xs font-mono px-2 py-0.2 rounded bg-slate-100 text-slate-700 font-semibold border border-slate-200">
                {{ item.action }}
              </span>
            </div>
            <span class="text-3xs font-mono text-slate-400">{{ item.timestamp || item.date }}</span>
          </div>

          <div class="mt-2 flex items-center gap-2 text-xs">
            <span class="font-semibold text-brand-700 font-mono">{{ item.targetName }}</span>
            <span v-if="item.version" class="text-3xs font-mono px-1.5 py-0.2 rounded bg-brand-50 text-brand-700 border border-brand-200">
              {{ item.version }}
            </span>
            <span v-if="item.revision" class="text-3xs font-mono px-1.5 py-0.2 rounded bg-slate-100 text-slate-700 border border-slate-200">
              {{ item.revision }}
            </span>
          </div>

          <p class="text-xs text-slate-600 mt-1 leading-relaxed font-medium">
            {{ item.details }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChangeHistoryTimeline',
  props: {
    items: {
      type: Array,
      required: true
    }
  }
};
</script>
