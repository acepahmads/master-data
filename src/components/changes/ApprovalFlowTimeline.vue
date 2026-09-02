<template>
  <div class="space-y-4">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div
        v-for="stage in ccbStages"
        :key="stage.key"
        :class="[
          'p-4 rounded-xl border transition-all',
          getStageContainerClass(stage.key)
        ]"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <span :class="['w-2 h-2 rounded-full', getStatusDotClass(getApproval(stage.key)?.decision)]"></span>
            <span class="text-xs font-bold uppercase tracking-wider text-slate-300">
              {{ stage.title }}
            </span>
          </div>
          <span
            :class="[
              'px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider',
              getDecisionBadgeClass(getApproval(stage.key)?.decision)
            ]"
          >
            {{ getApproval(stage.key)?.decision || 'PENDING' }}
          </span>
        </div>

        <div class="text-xs text-slate-400 mb-2">
          Required Role: <span class="text-slate-300 font-medium">{{ stage.role }}</span>
        </div>

        <div v-if="getApproval(stage.key)" class="space-y-2 mt-3 pt-2 border-t border-slate-800 text-xs">
          <div class="flex items-center justify-between text-slate-400">
            <span>Sign-off:</span>
            <span class="text-slate-200 font-semibold">{{ getApproval(stage.key).approverName }}</span>
          </div>
          <div v-if="getApproval(stage.key).decidedAt" class="text-[11px] text-slate-500">
            {{ formatDate(getApproval(stage.key).decidedAt) }}
          </div>
          <div v-if="getApproval(stage.key).justification" class="p-2 rounded bg-slate-950/60 border border-slate-800 text-[11px] text-slate-300 italic">
            "{{ getApproval(stage.key).justification }}"
          </div>
        </div>

        <div v-else class="mt-4 pt-2 border-t border-slate-800/60 flex items-center justify-between">
          <span class="text-[11px] text-slate-500 italic">Awaiting CCB review...</span>
          <button
            v-if="canApprove"
            @click="$emit('open-approval-dialog', stage.key)"
            class="px-2.5 py-1 rounded-lg bg-emerald-600/20 hover:bg-emerald-600/30 text-emerald-300 border border-emerald-500/40 text-xs font-semibold transition-colors"
          >
            Sign Off
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ApprovalFlowTimeline',
  props: {
    approvals: {
      type: Array,
      default: () => []
    },
    canApprove: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      ccbStages: [
        { key: 'TIER_1_TECH_LEAD', title: 'Tier 1: Tech Lead', role: 'Hardware / System Lead' },
        { key: 'TIER_2_ENG_MANAGER', title: 'Tier 2: Eng Manager', role: 'Engineering Manager' },
        { key: 'TIER_3_QUALITY_DIRECTOR', title: 'Tier 3: Quality Director', role: 'Quality & Reliability Director' }
      ]
    };
  },
  methods: {
    getApproval(stageKey) {
      return (this.approvals || []).find((a) => a.stage === stageKey);
    },
    getStatusDotClass(decision) {
      switch ((decision || '').toUpperCase()) {
        case 'APPROVED':
          return 'bg-emerald-400';
        case 'REJECTED':
          return 'bg-rose-400';
        case 'REQUEST_CHANGES':
          return 'bg-amber-400';
        default:
          return 'bg-slate-500 animate-pulse';
      }
    },
    getDecisionBadgeClass(decision) {
      switch ((decision || '').toUpperCase()) {
        case 'APPROVED':
          return 'bg-emerald-950/80 text-emerald-300 border border-emerald-800';
        case 'REJECTED':
          return 'bg-rose-950/80 text-rose-300 border border-rose-800';
        case 'REQUEST_CHANGES':
          return 'bg-amber-950/80 text-amber-300 border border-amber-800';
        default:
          return 'bg-slate-800 text-slate-400 border border-slate-700';
      }
    },
    getStageContainerClass(stageKey) {
      const appr = this.getApproval(stageKey);
      if (!appr) return 'bg-slate-900/40 border-slate-800/80';
      switch (appr.decision) {
        case 'APPROVED':
          return 'bg-emerald-950/15 border-emerald-800/40 shadow-sm shadow-emerald-950/20';
        case 'REJECTED':
          return 'bg-rose-950/20 border-rose-800/50 shadow-sm shadow-rose-950/20';
        case 'REQUEST_CHANGES':
          return 'bg-amber-950/15 border-amber-800/40';
        default:
          return 'bg-slate-900/60 border-slate-800';
      }
    },
    formatDate(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      return d.toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  }
};
</script>
