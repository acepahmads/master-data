<template>
  <div class="space-y-4">
    <div class="grid grid-cols-1 md:grid-cols-5 gap-3">
      <div
        v-for="domain in domains"
        :key="domain.key"
        :class="[
          'p-4 rounded-xl border transition-all',
          getDomainCardClass(domain.key)
        ]"
      >
        <div class="flex items-center justify-between mb-2">
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">
            {{ domain.label }}
          </span>
          <span
            :class="[
              'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
              getLevelBadgeClass(getImpactLevel(domain.key))
            ]"
          >
            {{ getImpactLevel(domain.key) }}
          </span>
        </div>

        <p class="text-xs text-slate-300 line-clamp-3 min-h-[36px]">
          {{ getImpactDescription(domain.key) }}
        </p>

        <div class="mt-3 pt-2 border-t border-slate-700/50 flex items-center justify-between text-[11px] text-slate-400">
          <span v-if="domain.key === 'COST'">
            Delta: <strong :class="getCostDeltaClass">{{ formatCost(getImpactData(domain.key)?.costDelta) }}</strong>
          </span>
          <span v-else-if="domain.key === 'SCHEDULE'">
            Delay: <strong class="text-amber-300">{{ getImpactData(domain.key)?.scheduleImpactWeeks || 0 }} wks</strong>
          </span>
          <span v-else-if="domain.key === 'INVENTORY'">
            Disp: <strong class="text-slate-200">{{ getImpactData(domain.key)?.scrapDisposition || 'RUN_OUT' }}</strong>
          </span>
          <span v-else>
            Reviewer: <strong class="text-slate-300">{{ getImpactData(domain.key)?.reviewerName || 'Unassigned' }}</strong>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChangeImpactMatrix',
  props: {
    impacts: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      domains: [
        { key: 'COST', label: 'Unit & Tooling Cost' },
        { key: 'SCHEDULE', label: 'Schedule & Launch' },
        { key: 'INVENTORY', label: 'Inventory & Scrap' },
        { key: 'TESTING', label: 'Testing & Validation' },
        { key: 'COMPLIANCE', label: 'Cert & Compliance' }
      ]
    };
  },
  computed: {
    getCostDeltaClass() {
      const costItem = this.getImpactData('COST');
      if (!costItem || costItem.costDelta === 0) return 'text-slate-300';
      return costItem.costDelta > 0 ? 'text-amber-400' : 'text-emerald-400';
    }
  },
  methods: {
    getImpactData(domainKey) {
      return (this.impacts || []).find((i) => i.domain === domainKey);
    },
    getImpactLevel(domainKey) {
      const data = this.getImpactData(domainKey);
      return data?.impactLevel || 'NONE';
    },
    getImpactDescription(domainKey) {
      const data = this.getImpactData(domainKey);
      return data?.description || 'No impact analysis recorded yet for this domain.';
    },
    getLevelBadgeClass(level) {
      switch ((level || '').toUpperCase()) {
        case 'CRITICAL':
          return 'bg-rose-950/80 text-rose-300 border border-rose-800';
        case 'HIGH':
          return 'bg-amber-950/80 text-amber-300 border border-amber-800';
        case 'MEDIUM':
          return 'bg-blue-950/80 text-blue-300 border border-blue-800';
        case 'LOW':
          return 'bg-slate-800 text-slate-400 border border-slate-700';
        default:
          return 'bg-slate-800/40 text-slate-500 border border-slate-800';
      }
    },
    getDomainCardClass(domainKey) {
      const level = this.getImpactLevel(domainKey);
      switch ((level || '').toUpperCase()) {
        case 'CRITICAL':
          return 'bg-rose-950/20 border-rose-800/40 shadow-sm shadow-rose-950/30';
        case 'HIGH':
          return 'bg-amber-950/20 border-amber-800/40 shadow-sm shadow-amber-950/30';
        case 'MEDIUM':
          return 'bg-blue-950/15 border-blue-800/30';
        case 'LOW':
          return 'bg-slate-900/60 border-slate-800';
        default:
          return 'bg-slate-900/40 border-slate-800/60';
      }
    },
    formatCost(val) {
      if (val === undefined || val === null) return '$0.00';
      const prefix = val > 0 ? '+$' : '$';
      return `${prefix}${Number(val).toFixed(2)}`;
    }
  }
};
</script>
