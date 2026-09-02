<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider',
      priorityClass
    ]"
  >
    <svg
      v-if="normalizedPriority === 'CRITICAL' || normalizedPriority === 'HIGH'"
      class="w-3 h-3"
      fill="currentColor"
      viewBox="0 0 20 20"
    >
      <path
        fill-rule="evenodd"
        d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
        clip-rule="evenodd"
      />
    </svg>
    {{ priority || 'NORMAL' }}
  </span>
</template>

<script>
export default {
  name: 'ChangePriorityBadge',
  props: {
    priority: {
      type: String,
      default: 'MEDIUM'
    }
  },
  computed: {
    normalizedPriority() {
      return (this.priority || '').toUpperCase();
    },
    priorityClass() {
      switch (this.normalizedPriority) {
        case 'CRITICAL':
          return 'bg-rose-950/60 text-rose-300 border border-rose-800/80 animate-pulse';
        case 'HIGH':
          return 'bg-amber-950/50 text-amber-300 border border-amber-800/70';
        case 'MEDIUM':
          return 'bg-blue-950/40 text-blue-300 border border-blue-800/60';
        case 'LOW':
          return 'bg-slate-800/80 text-slate-400 border border-slate-700/50';
        default:
          return 'bg-slate-800 text-slate-400 border border-slate-700';
      }
    }
  }
};
</script>
