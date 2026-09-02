<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full font-medium tracking-tight whitespace-nowrap text-2xs font-mono',
      badgeClass
    ]"
  >
    <span :class="['w-1.5 h-1.5 rounded-full shrink-0', dotClass]"></span>
    <span>{{ status }}</span>
  </span>
</template>

<script>
export default {
  name: 'VersionBadge',
  props: {
    status: {
      type: String,
      default: 'Active'
    }
  },
  computed: {
    normalizedStatus() {
      return (this.status || '').toLowerCase().trim();
    },
    badgeClass() {
      const s = this.normalizedStatus;
      if (s === 'active' || s === 'released' || s === 'production') {
        return 'bg-emerald-50 text-emerald-700 border border-emerald-200/80';
      }
      if (s === 'in review' || s === 'under review') {
        return 'bg-amber-50 text-amber-700 border border-amber-200/80';
      }
      if (s === 'draft' || s === 'prototype' || s === 'new') {
        return 'bg-slate-100 text-slate-700 border border-slate-200';
      }
      if (s === 'deprecated') {
        return 'bg-amber-50/80 text-amber-800 border border-amber-300/80';
      }
      if (s === 'archived' || s === 'obsolete' || s === 'eol') {
        return 'bg-rose-50 text-rose-700 border border-rose-200/80';
      }
      return 'bg-slate-100 text-slate-600 border border-slate-200';
    },
    dotClass() {
      const s = this.normalizedStatus;
      if (s === 'active' || s === 'released' || s === 'production') {
        return 'bg-emerald-500';
      }
      if (s === 'in review' || s === 'under review') {
        return 'bg-amber-500';
      }
      if (s === 'draft' || s === 'prototype' || s === 'new') {
        return 'bg-slate-400';
      }
      if (s === 'deprecated') {
        return 'bg-amber-600';
      }
      if (s === 'archived' || s === 'obsolete' || s === 'eol') {
        return 'bg-rose-500';
      }
      return 'bg-slate-400';
    }
  }
};
</script>
