<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider',
      badgeClass
    ]"
  >
    <span :class="['w-1.5 h-1.5 rounded-full', dotClass]"></span>
    {{ formatStatus(status) }}
  </span>
</template>

<script>
export default {
  name: 'ChangeStatusBadge',
  props: {
    status: {
      type: String,
      required: true
    }
  },
  computed: {
    normalizedStatus() {
      return (this.status || '').toUpperCase();
    },
    badgeClass() {
      switch (this.normalizedStatus) {
        case 'DRAFT':
          return 'bg-slate-800/80 text-slate-300 border border-slate-700/60';
        case 'UNDER_REVIEW':
        case 'IN_REVIEW':
          return 'bg-amber-950/40 text-amber-300 border border-amber-800/50';
        case 'PENDING_APPROVAL':
          return 'bg-blue-950/40 text-blue-300 border border-blue-800/50';
        case 'APPROVED':
          return 'bg-emerald-950/40 text-emerald-300 border border-emerald-800/50';
        case 'REJECTED':
          return 'bg-rose-950/40 text-rose-300 border border-rose-800/50';
        case 'CONVERTED_TO_ECO':
          return 'bg-indigo-950/40 text-indigo-300 border border-indigo-800/50';
        case 'IN_IMPLEMENTATION':
          return 'bg-cyan-950/40 text-cyan-300 border border-cyan-800/50';
        case 'IMPLEMENTED':
          return 'bg-emerald-950/40 text-emerald-300 border border-emerald-800/50';
        case 'VERIFIED':
          return 'bg-purple-950/40 text-purple-300 border border-purple-800/50';
        case 'CLOSED':
          return 'bg-slate-900/60 text-slate-400 border border-slate-800';
        default:
          return 'bg-slate-800 text-slate-400 border border-slate-700';
      }
    },
    dotClass() {
      switch (this.normalizedStatus) {
        case 'DRAFT':
          return 'bg-slate-400';
        case 'UNDER_REVIEW':
        case 'IN_REVIEW':
          return 'bg-amber-400 animate-pulse';
        case 'PENDING_APPROVAL':
          return 'bg-blue-400 animate-pulse';
        case 'APPROVED':
        case 'IMPLEMENTED':
          return 'bg-emerald-400';
        case 'REJECTED':
          return 'bg-rose-400';
        case 'CONVERTED_TO_ECO':
          return 'bg-indigo-400';
        case 'IN_IMPLEMENTATION':
          return 'bg-cyan-400 animate-pulse';
        case 'VERIFIED':
          return 'bg-purple-400';
        case 'CLOSED':
          return 'bg-slate-500';
        default:
          return 'bg-slate-400';
      }
    }
  },
  methods: {
    formatStatus(status) {
      if (!status) return 'UNKNOWN';
      return status.replace(/_/g, ' ');
    }
  }
};
</script>
