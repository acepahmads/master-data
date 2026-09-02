<template>
  <div class="bg-white rounded-xl p-4 border border-slate-200/80 shadow-subtle hover:shadow-card-hover hover:border-slate-300 transition-all duration-200 flex flex-col justify-between group relative overflow-hidden">
    <!-- Top Row: Title & Icon -->
    <div class="flex items-center justify-between gap-2 mb-2.5">
      <span class="text-3xs font-bold uppercase tracking-wider text-slate-500 truncate">
        {{ title }}
      </span>
      <div :class="['w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-transform duration-200 group-hover:scale-105 shadow-xs', iconBgClass]">
        <AppIcon :name="icon" size="sm" :class="iconColorClass" />
      </div>
    </div>

    <!-- Middle: Metric Value & Trend -->
    <div class="flex items-baseline justify-between gap-2 my-1">
      <div class="text-2xl font-bold text-slate-900 tracking-tight font-mono">
        {{ value }}
      </div>
      <span
        v-if="trend"
        :class="[
          'text-3xs font-semibold px-2 py-0.5 rounded-full flex items-center gap-0.5 border shadow-2xs',
          trendType === 'up' ? 'bg-emerald-50 text-emerald-700 border-emerald-200/80' : 
          trendType === 'down' ? 'bg-rose-50 text-rose-700 border-rose-200/80' : 'bg-slate-50 text-slate-600 border-slate-200'
        ]"
      >
        <AppIcon v-if="trendType === 'up'" name="arrow-up-right" size="xs" />
        <AppIcon v-else-if="trendType === 'down'" name="arrow-down-right" size="xs" />
        {{ trend }}
      </span>
    </div>

    <!-- Bottom Row: Subtitle & Micro Visualization -->
    <div class="mt-2 pt-2 border-t border-slate-100/80 flex items-center justify-between text-3xs text-slate-400">
      <span class="truncate font-medium text-slate-500">{{ subtitle }}</span>

      <!-- Subtle SVG Sparkline Visualization -->
      <svg class="w-14 h-4 text-slate-300 group-hover:text-brand-500 transition-colors duration-200 shrink-0" viewBox="0 0 60 16" fill="none" stroke="currentColor">
        <path
          :d="sparklinePath"
          stroke-width="1.75"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </div>

    <!-- Bottom hover accent line -->
    <div :class="['absolute bottom-0 left-0 right-0 h-0.5 opacity-0 group-hover:opacity-100 transition-opacity duration-200', accentLineClass]"></div>
  </div>
</template>

<script>
import AppIcon from './AppIcon.vue';

export default {
  name: 'KpiCard',
  components: { AppIcon },
  props: {
    title: {
      type: String,
      required: true
    },
    value: {
      type: [Number, String],
      required: true
    },
    icon: {
      type: String,
      default: 'box'
    },
    iconColor: {
      type: String,
      default: 'blue' // blue, emerald, purple, amber, cyan
    },
    trend: {
      type: String,
      default: ''
    },
    trendType: {
      type: String,
      default: 'up'
    },
    subtitle: {
      type: String,
      default: ''
    },
    sparkVariant: {
      type: Number,
      default: 1
    }
  },
  computed: {
    iconBgClass() {
      switch (this.iconColor) {
        case 'blue': return 'bg-brand-50 text-brand-600 border border-brand-200/80';
        case 'emerald': return 'bg-emerald-50 text-emerald-600 border border-emerald-200/80';
        case 'purple': return 'bg-purple-50 text-purple-600 border border-purple-200/80';
        case 'amber': return 'bg-amber-50 text-amber-600 border border-amber-200/80';
        case 'cyan': return 'bg-cyan-50 text-cyan-600 border border-cyan-200/80';
        default: return 'bg-slate-100 text-slate-600 border border-slate-200';
      }
    },
    iconColorClass() {
      switch (this.iconColor) {
        case 'blue': return 'text-brand-600';
        case 'emerald': return 'text-emerald-600';
        case 'purple': return 'text-purple-600';
        case 'amber': return 'text-amber-600';
        case 'cyan': return 'text-cyan-600';
        default: return 'text-slate-600';
      }
    },
    accentLineClass() {
      switch (this.iconColor) {
        case 'blue': return 'bg-brand-500';
        case 'emerald': return 'bg-emerald-500';
        case 'purple': return 'bg-purple-500';
        case 'amber': return 'bg-amber-500';
        case 'cyan': return 'bg-cyan-500';
        default: return 'bg-brand-500';
      }
    },
    sparklinePath() {
      const paths = [
        'M 2 12 Q 15 4, 30 9 T 58 3',
        'M 2 14 L 12 10 L 25 12 L 38 6 L 48 8 L 58 2',
        'M 2 10 C 18 16, 35 2, 58 6',
        'M 2 13 L 15 13 L 28 8 L 42 8 L 58 3',
        'M 2 8 Q 20 14, 40 4 T 58 2'
      ];
      return paths[(this.sparkVariant - 1) % paths.length] || paths[0];
    }
  }
};
</script>

<style scoped>
.text-3xs {
  font-size: 0.625rem;
}
.shadow-2xs {
  box-shadow: 0 1px 1px 0 rgba(0, 0, 0, 0.04);
}
</style>
