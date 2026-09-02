<template>
  <div class="w-full overflow-x-auto py-3">
    <div class="flex items-center min-w-max gap-0 relative px-4">
      <!-- Background Connecting Rail -->
      <div class="absolute left-8 right-8 top-1/2 -translate-y-1/2 h-0.5 bg-slate-200 z-0"></div>

      <div
        v-for="(ver, index) in sortedVersions"
        :key="ver.id"
        class="flex items-center relative z-10"
      >
        <!-- Version Node Card -->
        <div
          @click="$emit('select-version', ver)"
          :class="[
            'p-3 rounded-xl border transition-all duration-200 cursor-pointer flex flex-col items-center text-center group min-w-[140px] max-w-[170px]',
            activeVersionId === ver.id
              ? 'bg-white border-brand-500 shadow-glow-brand ring-2 ring-brand-500/20 scale-105'
              : 'bg-white border-slate-200 hover:border-slate-300 hover:shadow-subtle'
          ]"
        >
          <!-- Version Node Dot / Icon -->
          <div
            :class="[
              'w-7 h-7 rounded-full flex items-center justify-center font-mono font-bold text-3xs mb-2 transition-transform group-hover:scale-110 shadow-xs',
              activeVersionId === ver.id ? 'bg-brand-600 text-white' :
              ver.status === 'Active' ? 'bg-emerald-100 text-emerald-700 border border-emerald-300' :
              ver.status === 'Draft' ? 'bg-amber-100 text-amber-700 border border-amber-300' :
              'bg-slate-100 text-slate-500 border border-slate-200'
            ]"
          >
            v{{ ver.versionNumber.split('.')[0] }}
          </div>

          <span class="font-mono font-bold text-xs text-slate-900 group-hover:text-brand-600 transition-colors">
            v{{ ver.versionNumber }}
          </span>

          <span class="text-3xs text-slate-500 truncate max-w-[130px] mt-0.5 font-medium">
            {{ ver.versionName }}
          </span>

          <div class="mt-2 flex items-center gap-1">
            <VersionBadge :status="ver.status" />
          </div>

          <div class="text-3xs font-mono text-slate-400 mt-1.5 flex items-center gap-1">
            <span class="font-semibold text-slate-600">{{ ver.currentRevision }}</span>
            <span>•</span>
            <span>{{ formatDate(ver.releaseDate || ver.createdAt) }}</span>
          </div>
        </div>

        <!-- Rail Spacer / Connector (except last) -->
        <div v-if="index < sortedVersions.length - 1" class="w-8 md:w-12 h-0.5"></div>
      </div>
    </div>
  </div>
</template>

<script>
import VersionBadge from './VersionBadge.vue';

export default {
  name: 'VersionTimeline',
  components: { VersionBadge },
  props: {
    versions: {
      type: Array,
      required: true
    },
    activeVersionId: {
      type: String,
      default: ''
    }
  },
  computed: {
    sortedVersions() {
      // Sort chronologically ascending
      return [...this.versions].sort((a, b) => {
        return (a.versionNumber || '').localeCompare(b.versionNumber || '', undefined, { numeric: true });
      });
    }
  },
  methods: {
    formatDate(dateStr) {
      if (!dateStr) return 'N/A';
      return dateStr.length > 7 ? dateStr.substring(0, 7) : dateStr;
    }
  }
};
</script>
