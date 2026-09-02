<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <span class="p-2 rounded-xl bg-purple-500/10 text-purple-400">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </span>
          Engineering Change Traceability Matrix
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Closed-loop Directed Acyclic Graph (DAG) demonstrating complete lineage from Defect Finding to Re-test Validation.
        </p>
      </div>

      <button
        @click="refreshData"
        class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors flex items-center gap-2 self-start md:self-auto"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh Traceability
      </button>
    </div>

    <!-- Traceability Chains Visualization Component -->
    <TraceabilityGraph :chains="traceabilityChains" />
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';
import TraceabilityGraph from '@/components/changes/TraceabilityGraph.vue';

export default {
  name: 'TraceabilityMatrixView',
  components: {
    TraceabilityGraph
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    traceabilityChains() {
      return this.changeStore.traceabilityChains;
    }
  },
  async mounted() {
    await this.refreshData();
  },
  methods: {
    async refreshData() {
      await this.changeStore.fetchTraceability();
    }
  }
};
</script>
