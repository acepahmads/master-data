<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-primary-400 uppercase tracking-wider mb-1">
          <span class="inline-block w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          Phase 5 — Testing & Validation Control
        </div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <i class="feather-check-circle text-emerald-400"></i>
          Testing Command Center
        </h1>
        <p class="text-sm text-slate-400 mt-0.5">
          Laboratory verification, multi-run parametric measurements, defect tracking, and validation governance.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <router-link
          to="/test-plans"
          class="px-3.5 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 rounded-lg text-sm font-medium flex items-center gap-2 transition"
        >
          <i class="feather-layers text-primary-400"></i>
          Test Plan Registry
        </router-link>
        <router-link
          to="/findings"
          class="px-3.5 py-2 bg-amber-500/10 hover:bg-amber-500/20 text-amber-300 border border-amber-500/30 rounded-lg text-sm font-medium flex items-center gap-2 transition"
        >
          <i class="feather-alert-triangle text-amber-400"></i>
          Findings Hub ({{ stats.openFindingsCount }})
        </router-link>
      </div>
    </div>

    <!-- Top KPI Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Pass Rate Card -->
      <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 relative overflow-hidden backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400 uppercase tracking-wider">Overall Pass Rate</span>
          <span class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400">
            <i class="feather-award text-lg"></i>
          </span>
        </div>
        <div class="mt-3 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-slate-100">{{ stats.overallPassRate.toFixed(1) }}%</span>
          <span class="text-xs text-slate-400">({{ stats.totalPassed }}/{{ stats.totalExecutions }} passed)</span>
        </div>
        <div class="mt-3 w-full bg-slate-700/50 rounded-full h-1.5 overflow-hidden">
          <div
            class="bg-gradient-to-r from-emerald-500 to-teal-400 h-1.5 rounded-full transition-all duration-500"
            :style="{ width: `${stats.overallPassRate}%` }"
          ></div>
        </div>
      </div>

      <!-- Active Sessions Card -->
      <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 relative overflow-hidden backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400 uppercase tracking-wider">Active Test Sessions</span>
          <span class="p-2 rounded-lg bg-indigo-500/10 text-indigo-400">
            <i class="feather-activity text-lg"></i>
          </span>
        </div>
        <div class="mt-3 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-slate-100">{{ stats.activeTestSessions }}</span>
          <span class="text-xs text-slate-400">sessions in progress</span>
        </div>
        <div class="mt-3 flex items-center gap-2 text-xs text-slate-400">
          <i class="feather-file-text text-indigo-400"></i>
          <span>{{ stats.totalTestPlans }} Master Test Protocols</span>
        </div>
      </div>

      <!-- Open Findings / Defects -->
      <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 relative overflow-hidden backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400 uppercase tracking-wider">Open Defects / Findings</span>
          <span class="p-2 rounded-lg bg-amber-500/10 text-amber-400">
            <i class="feather-alert-octagon text-lg"></i>
          </span>
        </div>
        <div class="mt-3 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-amber-400">{{ stats.openFindingsCount }}</span>
          <span class="text-xs text-slate-400">active anomalies</span>
        </div>
        <div class="mt-3 flex items-center gap-2">
          <span
            v-if="stats.criticalFindingsCount > 0"
            class="px-2 py-0.5 rounded text-[11px] font-semibold bg-rose-500/20 text-rose-300 border border-rose-500/30"
          >
            {{ stats.criticalFindingsCount }} Critical
          </span>
          <span
            v-if="stats.highFindingsCount > 0"
            class="px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/20 text-amber-300 border border-amber-500/30"
          >
            {{ stats.highFindingsCount }} High
          </span>
          <span
            v-if="stats.criticalFindingsCount === 0 && stats.highFindingsCount === 0"
            class="text-xs text-emerald-400 flex items-center gap-1"
          >
            <i class="feather-check"></i> Zero critical blockers
          </span>
        </div>
      </div>

      <!-- Change Candidates for Future Phase 6 -->
      <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 relative overflow-hidden backdrop-blur-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400 uppercase tracking-wider">Change Candidates</span>
          <span class="p-2 rounded-lg bg-purple-500/10 text-purple-400">
            <i class="feather-git-pull-request text-lg"></i>
          </span>
        </div>
        <div class="mt-3 flex items-baseline gap-2">
          <span class="text-3xl font-extrabold text-purple-300">{{ stats.changeCandidatesCount }}</span>
          <span class="text-xs text-slate-400">tagged for ECO review</span>
        </div>
        <div class="mt-3 text-xs text-slate-400 flex items-center gap-1.5">
          <i class="feather-info text-purple-400"></i>
          <span>Phase 6 ECR Readiness</span>
        </div>
      </div>
    </div>

    <!-- Main Content Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left 2 Cols: Master Protocols & Active Sessions -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Master Protocols List -->
        <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 backdrop-blur-sm">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-base font-semibold text-slate-100 flex items-center gap-2">
              <i class="feather-layers text-primary-400"></i>
              Master Test Plan Protocols
            </h2>
            <router-link to="/test-plans" class="text-xs font-medium text-primary-400 hover:text-primary-300 flex items-center gap-1">
              View All Protocols <i class="feather-arrow-right text-xs"></i>
            </router-link>
          </div>

          <div v-if="testPlans.length === 0" class="p-8 text-center text-slate-400 text-sm">
            No test plans found. Create your first master test protocol.
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="tp in testPlans.slice(0, 4)"
              :key="tp.id"
              class="p-4 rounded-lg bg-slate-900/60 border border-slate-700/40 hover:border-slate-600 transition flex items-center justify-between"
            >
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <span class="text-xs font-mono px-2 py-0.5 rounded bg-primary-500/10 text-primary-300 border border-primary-500/20 font-semibold">
                    {{ tp.code }}
                  </span>
                  <span class="text-sm font-semibold text-slate-200">{{ tp.name }}</span>
                </div>
                <p class="text-xs text-slate-400 line-clamp-1">{{ tp.description }}</p>
                <div class="flex items-center gap-3 text-[11px] text-slate-400 pt-1">
                  <span><i class="feather-tag text-slate-500"></i> {{ tp.category || 'GENERAL' }}</span>
                  <span><i class="feather-git-commit text-slate-500"></i> {{ tp.versions ? tp.versions.length : 1 }} Version(s)</span>
                  <span><i class="feather-user text-slate-500"></i> {{ tp.createdBy }}</span>
                </div>
              </div>

              <div class="flex items-center gap-2">
                <router-link
                  :to="`/test-plans/${tp.id}`"
                  class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 rounded text-xs font-medium transition"
                >
                  Workspace
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Findings & Defects -->
        <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 backdrop-blur-sm">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-base font-semibold text-slate-100 flex items-center gap-2">
              <i class="feather-alert-triangle text-amber-400"></i>
              Active Engineering Findings & Defect Candidates
            </h2>
            <router-link to="/findings" class="text-xs font-medium text-primary-400 hover:text-primary-300 flex items-center gap-1">
              View Defect Registry <i class="feather-arrow-right text-xs"></i>
            </router-link>
          </div>

          <div v-if="findings.length === 0" class="p-6 text-center text-slate-400 text-sm">
            Zero active engineering findings logged.
          </div>

          <div v-else class="space-y-2.5">
            <div
              v-for="f in findings.slice(0, 4)"
              :key="f.id"
              class="p-3.5 rounded-lg bg-slate-900/60 border border-slate-700/40 hover:border-slate-600 transition flex items-center justify-between"
            >
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <span
                    class="text-[10px] font-bold px-2 py-0.5 rounded uppercase"
                    :class="getSeverityBadge(f.severity)"
                  >
                    {{ f.severity }}
                  </span>
                  <span class="text-xs font-mono text-slate-400 font-semibold">{{ f.code }}</span>
                  <span class="text-sm font-medium text-slate-200">{{ f.title }}</span>
                </div>
                <p class="text-xs text-slate-400 line-clamp-1">{{ f.description }}</p>
              </div>

              <div class="flex items-center gap-2 shrink-0">
                <span
                  v-if="f.changeCandidateStatus === 'CANDIDATE'"
                  class="px-2 py-0.5 rounded text-[10px] font-semibold bg-purple-500/10 text-purple-300 border border-purple-500/20"
                >
                  ECO Candidate
                </span>
                <span
                  class="px-2 py-0.5 rounded text-[10px] font-medium bg-slate-800 text-slate-300 border border-slate-700"
                >
                  {{ formatDisposition(f.findingDisposition) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right 1 Col: Quick Tools & Protocol Governance -->
      <div class="space-y-6">
        <!-- Testing Governance Summary -->
        <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 backdrop-blur-sm space-y-4">
          <h3 class="text-sm font-semibold text-slate-200 flex items-center gap-2">
            <i class="feather-shield text-emerald-400"></i>
            Validation Governance Status
          </h3>

          <div class="space-y-2.5 text-xs">
            <div class="flex justify-between items-center p-2.5 rounded-lg bg-slate-900/40 border border-slate-700/30">
              <span class="text-slate-400">Validated Sessions</span>
              <span class="font-semibold text-emerald-400">{{ stats.totalValidatedSessions }}</span>
            </div>
            <div class="flex justify-between items-center p-2.5 rounded-lg bg-slate-900/40 border border-slate-700/30">
              <span class="text-slate-400">Pending Review</span>
              <span class="font-semibold text-amber-400">{{ stats.pendingValidations }}</span>
            </div>
            <div class="flex justify-between items-center p-2.5 rounded-lg bg-slate-900/40 border border-slate-700/30">
              <span class="text-slate-400">Rejected Sessions</span>
              <span class="font-semibold text-rose-400">{{ stats.totalRejectedSessions }}</span>
            </div>
          </div>

          <div class="pt-2 border-t border-slate-700/60">
            <div class="text-[11px] text-slate-400 leading-relaxed">
              <span class="font-semibold text-slate-300">Validation Blocking Rule:</span>
              Unresolved <span class="text-rose-300 font-medium">CRITICAL</span> or <span class="text-amber-300 font-medium">HIGH</span> findings strictly block formal prototype validation sign-off.
            </div>
          </div>
        </div>

        <!-- Quick Links -->
        <div class="bg-gradient-to-br from-indigo-900/40 to-slate-900 border border-indigo-500/20 rounded-xl p-5 space-y-3">
          <h3 class="text-sm font-semibold text-indigo-200 flex items-center gap-2">
            <i class="feather-cpu text-indigo-400"></i>
            Prototype Testing Workspace
          </h3>
          <p class="text-xs text-slate-300 leading-relaxed">
            Assign released test protocols to physical prototype build units and log multi-run measurements directly inside Prototype Details.
          </p>
          <router-link
            to="/prototypes"
            class="block text-center w-full py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-semibold transition"
          >
            Go to Prototype Builds
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia';
import { useTestingStore } from '@/stores/testingStore';

export default {
  name: 'TestingDashboardView',
  computed: {
    ...mapState(useTestingStore, {
      stats: 'dashboardStats',
      testPlans: 'testPlans',
      findings: 'findings',
      loading: 'loading'
    })
  },
  async mounted() {
    await Promise.all([
      this.fetchDashboard(),
      this.fetchTestPlans({ limit: 5 }),
      this.fetchFindings({ limit: 5 })
    ]);
  },
  methods: {
    ...mapActions(useTestingStore, ['fetchDashboard', 'fetchTestPlans', 'fetchFindings']),
    getSeverityBadge(sev) {
      switch (sev) {
        case 'CRITICAL':
          return 'bg-rose-500/20 text-rose-300 border border-rose-500/30';
        case 'HIGH':
          return 'bg-amber-500/20 text-amber-300 border border-amber-500/30';
        case 'MEDIUM':
          return 'bg-yellow-500/20 text-yellow-300 border border-yellow-500/30';
        case 'LOW':
          return 'bg-blue-500/20 text-blue-300 border border-blue-500/30';
        default:
          return 'bg-slate-700 text-slate-300';
      }
    },
    formatDisposition(disp) {
      if (!disp) return 'OPEN';
      return disp.replace(/_/g, ' ');
    }
  }
};
</script>
