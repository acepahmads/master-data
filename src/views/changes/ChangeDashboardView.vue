<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <svg class="w-7 h-7 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Engineering Change Control Center
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Closed-loop ECR/ECO governance, CCB multi-tier sign-off, automated baseline cloning, and traceability.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <router-link
          to="/ecr/create"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-semibold text-sm shadow-lg shadow-cyan-500/20 transition-all"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          New ECR
        </router-link>
        <router-link
          to="/changes/traceability"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-sm font-semibold border border-slate-700 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Traceability Matrix
        </router-link>
      </div>
    </div>

    <!-- KPI Metric Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Total ECRs</span>
          <span class="p-2 rounded-xl bg-blue-500/10 text-blue-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </span>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-slate-100">{{ stats.totalEcrs }}</span>
          <span class="text-xs text-amber-400 font-medium">{{ stats.underReviewEcrs }} Under Review</span>
        </div>
        <div class="mt-2 text-xs text-slate-500">
          {{ stats.approvedEcrs }} approved &bull; {{ stats.convertedEcrs }} converted to ECO
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Total ECOs</span>
          <span class="p-2 rounded-xl bg-cyan-500/10 text-cyan-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2" />
            </svg>
          </span>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-cyan-300">{{ stats.totalEcos }}</span>
          <span class="text-xs text-emerald-400 font-medium">{{ stats.implementedEcos }} Implemented</span>
        </div>
        <div class="mt-2 text-xs text-slate-500">
          {{ stats.draftEcos }} draft &bull; {{ stats.approvedEcos }} ready for baseline execution
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Closed-Loop Verified</span>
          <span class="p-2 rounded-xl bg-purple-500/10 text-purple-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </span>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-purple-300">{{ stats.verifiedEcos }}</span>
          <span class="text-xs text-slate-400 font-medium">Post-regression OK</span>
        </div>
        <div class="mt-2 text-xs text-slate-500">
          {{ stats.closedEcos }} fully closed out
        </div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Critical / High Changes</span>
          <span class="p-2 rounded-xl bg-rose-500/10 text-rose-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </span>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-rose-400">{{ stats.criticalHighEcrs }}</span>
          <span class="text-xs text-slate-400 font-medium">High Attention</span>
        </div>
        <div class="mt-2 text-xs text-slate-500">
          Avg cycle: {{ stats.avgImplementationCycleDays || 0 }} days
        </div>
      </div>
    </div>

    <!-- Quick Navigation Sub-Panels -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent ECRs List -->
      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-cyan-400"></span>
            Active Change Requests (ECR)
          </h3>
          <router-link to="/ecr" class="text-xs text-cyan-400 hover:text-cyan-300 font-medium">
            View All &rarr;
          </router-link>
        </div>

        <div class="divide-y divide-slate-800/80">
          <div
            v-for="ecr in recentEcrs"
            :key="ecr.id"
            class="py-3 flex items-center justify-between hover:bg-slate-800/20 px-2 rounded-lg transition-colors cursor-pointer"
            @click="$router.push(`/ecr/${ecr.id}`)"
          >
            <div>
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs font-bold text-slate-200">{{ ecr.code }}</span>
                <ChangeStatusBadge :status="ecr.status" />
                <ChangePriorityBadge :priority="ecr.priority" />
              </div>
              <div class="text-xs text-slate-400 mt-1 line-clamp-1 font-medium">
                {{ ecr.title }}
              </div>
            </div>
            <span class="text-xs text-slate-500 font-mono">{{ ecr.sourceHardwareRevision ? ecr.sourceHardwareRevision.code : '—' }}</span>
          </div>

          <div v-if="recentEcrs.length === 0" class="text-center py-6 text-xs text-slate-500">
            No active ECRs found.
          </div>
        </div>
      </div>

      <!-- Recent ECOs List -->
      <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
            Implementation Change Orders (ECO)
          </h3>
          <router-link to="/eco" class="text-xs text-cyan-400 hover:text-cyan-300 font-medium">
            View All &rarr;
          </router-link>
        </div>

        <div class="divide-y divide-slate-800/80">
          <div
            v-for="eco in recentEcos"
            :key="eco.id"
            class="py-3 flex items-center justify-between hover:bg-slate-800/20 px-2 rounded-lg transition-colors cursor-pointer"
            @click="$router.push(`/eco/${eco.id}`)"
          >
            <div>
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs font-bold text-cyan-300">{{ eco.code }}</span>
                <ChangeStatusBadge :status="eco.status" />
              </div>
              <div class="text-xs text-slate-400 mt-1 line-clamp-1 font-medium">
                {{ eco.title }}
              </div>
            </div>
            <div class="text-right">
              <div class="text-xs text-slate-300 font-medium">{{ eco.implementationOwner || '—' }}</div>
              <div class="text-[10px] text-slate-500">{{ eco.targetHardwareRevision ? eco.targetHardwareRevision.code : 'Baseline Pending' }}</div>
            </div>
          </div>

          <div v-if="recentEcos.length === 0" class="text-center py-6 text-xs text-slate-500">
            No active ECOs found.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';
import ChangeStatusBadge from '@/components/changes/ChangeStatusBadge.vue';
import ChangePriorityBadge from '@/components/changes/ChangePriorityBadge.vue';

export default {
  name: 'ChangeDashboardView',
  components: {
    ChangeStatusBadge,
    ChangePriorityBadge
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    stats() {
      return this.changeStore.dashboardStats;
    },
    recentEcrs() {
      return (this.changeStore.ecrs || []).slice(0, 5);
    },
    recentEcos() {
      return (this.changeStore.ecos || []).slice(0, 5);
    }
  },
  async mounted() {
    await Promise.all([
      this.changeStore.fetchDashboard(),
      this.changeStore.fetchECRs({ limit: 5 }),
      this.changeStore.fetchECOs({ limit: 5 })
    ]);
  }
};
</script>
