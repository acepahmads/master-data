<template>
  <div class="space-y-6">
    <!-- PAGE HEADER -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-slate-200">
      <div>
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center border border-indigo-200">
            <AppIcon name="cpu" size="sm" />
          </div>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <span>Prototype Build Command Center</span>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
              PHASE 4
            </span>
          </h1>
        </div>
        <p class="text-xs text-slate-500 mt-1 font-medium">
          Physical prototype unit tracking, frozen BOM snapshots, component readiness, and multi-stage assembly progress.
        </p>
      </div>

      <div class="flex items-center gap-2.5 shrink-0">
        <router-link to="/prototypes/list" class="btn btn-secondary text-xs flex items-center gap-1.5">
          <AppIcon name="list" size="xs" />
          <span>All Builds</span>
        </router-link>
        <router-link to="/prototypes/create" class="btn btn-primary text-xs flex items-center gap-1.5 bg-gradient-to-r from-indigo-600 to-brand-600 border-none shadow-sm">
          <AppIcon name="plus" size="xs" />
          <span>+ New Prototype Build</span>
        </router-link>
      </div>
    </div>

    <!-- 5 KPI CARDS -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3.5">
      <KpiCard
        title="Total Prototype Builds"
        :value="`${protoStore.totalPrototypesCount}`"
        icon="cpu"
        iconColor="indigo"
        trend="Active Builds"
        trendType="up"
        subtitle="Physical build records"
        :sparkVariant="1"
      />
      <KpiCard
        title="Component Kitting"
        :value="`${protoStore.componentPrepCount}`"
        icon="layers"
        iconColor="amber"
        trend="In Readiness"
        trendType="up"
        subtitle="Parts allocation check"
        :sparkVariant="2"
      />
      <KpiCard
        title="In Assembly"
        :value="`${protoStore.inAssemblyCount}`"
        icon="tool"
        iconColor="blue"
        trend="Active Bench"
        trendType="up"
        subtitle="SMT & wiring stage"
        :sparkVariant="3"
      />
      <KpiCard
        title="Ready for Testing"
        :value="`${protoStore.readyForTestingCount}`"
        icon="shield-check"
        iconColor="purple"
        trend="Phase 5 Gate"
        trendType="up"
        subtitle="Verification handoff"
        :sparkVariant="4"
      />
      <KpiCard
        title="Builds Completed"
        :value="`${protoStore.completedCount}`"
        icon="check-circle"
        iconColor="emerald"
        trend="Verified"
        trendType="up"
        subtitle="Fully assembled units"
        :sparkVariant="5"
      />
    </div>

    <!-- MAIN DASHBOARD CONTENT (2-COLUMN) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- LEFT 8 COLS: ACTIVE PROTOTYPE BUILDS -->
      <div class="lg:col-span-8 space-y-4">
        <div class="app-card p-5 space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="activity" size="xs" class="text-indigo-600" />
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Active Prototype Builds</h2>
            </div>
            <span class="text-3xs font-mono text-slate-400 font-semibold">{{ protoStore.prototypes.length }} Total Builds</span>
          </div>

          <!-- PROTOTYPE CARDS LIST -->
          <div class="space-y-3">
            <div
              v-for="p in protoStore.prototypes"
              :key="p.id"
              class="p-4 rounded-xl border border-slate-200/90 bg-white hover:border-indigo-300 hover:shadow-card transition-all group"
            >
              <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
                <div class="min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-mono text-xs font-bold text-indigo-700 bg-indigo-50 px-2 py-0.5 rounded border border-indigo-200">
                      {{ p.code }}
                    </span>
                    <span
                      class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full border shadow-2xs"
                      :class="getStatusBadgeClass(p.status)"
                    >
                      {{ formatStatus(p.status) }}
                    </span>
                    <span class="text-3xs font-mono text-slate-400">Build #{{ p.buildNumber }} ({{ p.quantity }} Units)</span>
                  </div>
                  <router-link
                    :to="`/prototypes/${p.id}`"
                    class="font-bold text-sm text-slate-900 group-hover:text-indigo-600 transition-colors"
                  >
                    {{ p.name }}
                  </router-link>
                  <div class="text-3xs text-slate-500 mt-0.5 flex flex-wrap items-center gap-2">
                    <span>Product: <strong class="text-slate-700 font-mono">{{ p.productCode }}</strong></span>
                    <span>•</span>
                    <span>Revision: <strong class="text-slate-700 font-mono">{{ p.hardwareRevisionCode }}</strong></span>
                    <span>•</span>
                    <span>Owner: <strong class="text-slate-700">{{ p.engineeringOwner }}</strong></span>
                  </div>
                </div>

                <div class="flex sm:flex-col items-end justify-between sm:justify-center shrink-0">
                  <div class="text-right">
                    <span class="text-3xs text-slate-400 block font-mono">Target Date</span>
                    <span class="text-xs font-bold font-mono text-slate-800">{{ p.targetBuildDate }}</span>
                  </div>
                  <router-link
                    :to="`/prototypes/${p.id}`"
                    class="btn btn-secondary py-1 px-2.5 text-3xs font-bold mt-1.5 flex items-center gap-1 group-hover:bg-indigo-50 group-hover:text-indigo-700 group-hover:border-indigo-200"
                  >
                    <span>Open Workspace</span>
                    <AppIcon name="chevron-right" size="2xs" />
                  </router-link>
                </div>
              </div>

              <!-- PROGRESS BAR & ASSEMBLY STAGES METRICS -->
              <div class="pt-3 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <div class="flex-1">
                  <div class="flex items-center justify-between text-3xs font-mono font-bold mb-1">
                    <span class="text-slate-400">Assembly Checklist</span>
                    <span class="text-indigo-600">{{ getCompletedStagesCount(p) }} / 8 Stages Done ({{ getAssemblyProgress(p) }}%)</span>
                  </div>
                  <div class="w-full h-1.5 rounded-full bg-slate-100 overflow-hidden border border-slate-200/60">
                    <div
                      class="h-full rounded-full bg-gradient-to-r from-indigo-500 to-brand-500 transition-all duration-500"
                      :style="{ width: `${getAssemblyProgress(p)}%` }"
                    ></div>
                  </div>
                </div>

                <!-- Frozen BOM Snapshot Summary Pill -->
                <div class="flex items-center gap-2 bg-slate-50 px-2.5 py-1.5 rounded-lg border border-slate-200/60 text-3xs font-mono shrink-0">
                  <span class="text-slate-400">BOM Snapshot:</span>
                  <span class="font-bold text-slate-700">{{ p.bomSnapshot.totalItems }} Parts</span>
                  <span class="text-slate-300">|</span>
                  <span class="font-bold text-emerald-700">${{ Number(p.bomSnapshot.totalEstimatedCost || 0).toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT 4 COLS: STATUS DISTRIBUTION & ACTIVITY STREAM -->
      <div class="lg:col-span-4 space-y-6">
        <!-- STATUS DISTRIBUTION -->
        <div class="app-card p-4.5 space-y-3.5">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="pie-chart" size="xs" class="text-indigo-600" />
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Build Status Breakdown</h2>
            </div>
            <span class="text-3xs font-mono text-slate-400">Workflow Pipeline</span>
          </div>

          <div class="space-y-2.5 text-xs">
            <div class="flex items-center justify-between text-3xs font-mono">
              <span class="text-slate-600 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-amber-400"></span>
                <span>Component Prep</span>
              </span>
              <span class="font-bold text-slate-800">{{ protoStore.componentPrepCount }} Builds</span>
            </div>
            <div class="flex items-center justify-between text-3xs font-mono">
              <span class="text-slate-600 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                <span>In Assembly</span>
              </span>
              <span class="font-bold text-slate-800">{{ protoStore.inAssemblyCount }} Builds</span>
            </div>
            <div class="flex items-center justify-between text-3xs font-mono">
              <span class="text-slate-600 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-purple-500"></span>
                <span>Ready for Testing</span>
              </span>
              <span class="font-bold text-slate-800">{{ protoStore.readyForTestingCount }} Builds</span>
            </div>
            <div class="flex items-center justify-between text-3xs font-mono">
              <span class="text-slate-600 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                <span>Completed</span>
              </span>
              <span class="font-bold text-slate-800">{{ protoStore.completedCount }} Builds</span>
            </div>
            <div class="flex items-center justify-between text-3xs font-mono">
              <span class="text-slate-600 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-slate-300"></span>
                <span>Draft / Planning</span>
              </span>
              <span class="font-bold text-slate-800">{{ protoStore.draftCount }} Builds</span>
            </div>
          </div>

          <!-- Quick Create CTA Card -->
          <div class="p-3 bg-indigo-50/60 rounded-xl border border-indigo-200/80 space-y-2 mt-3">
            <div class="font-bold text-xs text-indigo-900">Start New Prototype Run</div>
            <p class="text-3xs text-indigo-700 leading-relaxed">
              Capture a frozen snapshot of an approved Hardware Revision BOM and track assembly stages.
            </p>
            <router-link to="/prototypes/create" class="btn btn-primary btn-sm text-xs w-full bg-indigo-600 hover:bg-indigo-700 text-center block">
              + Launch Build Wizard
            </router-link>
          </div>
        </div>

        <!-- RECENT BUILD ACTIVITY STREAM -->
        <div class="app-card p-4.5 space-y-3.5">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="activity" size="xs" class="text-indigo-600" />
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Recent Build Events</h2>
            </div>
            <span class="text-3xs font-mono text-slate-400">Live Log</span>
          </div>

          <div class="relative pl-3 space-y-3 before:absolute before:left-1 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
            <div
              v-for="act in recentActivities"
              :key="act.id"
              class="relative text-xs pl-2"
            >
              <span class="absolute -left-3 top-1 w-2 h-2 rounded-full bg-indigo-500 ring-2 ring-white shadow-2xs"></span>
              <div class="font-bold text-slate-800 text-3xs leading-snug">{{ act.action }}</div>
              <p class="text-3xs text-slate-500 mt-0.5 leading-snug">{{ act.description }}</p>
              <div class="flex items-center justify-between text-3xs font-mono text-slate-400 mt-1">
                <span>{{ act.author }}</span>
                <span>{{ act.timestamp }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { usePrototypeStore } from '@/stores/prototypeStore';
import KpiCard from '@/components/common/KpiCard.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'PrototypeDashboardView',
  components: {
    KpiCard,
    AppIcon
  },
  computed: {
    protoStore() {
      return usePrototypeStore();
    },
    recentActivities() {
      const all = [];
      this.protoStore.prototypes.forEach(p => {
        if (p.activityHistory) {
          p.activityHistory.forEach(a => {
            all.push({ ...a, protoCode: p.code });
          });
        }
      });
      return all.slice(0, 6);
    }
  },
  methods: {
    formatStatus(status) {
      if (!status) return 'DRAFT';
      return status.replace(/_/g, ' ');
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'COMPLETED':
          return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'READY_FOR_TESTING':
          return 'bg-purple-50 text-purple-700 border-purple-200';
        case 'IN_ASSEMBLY':
        case 'ASSEMBLY_COMPLETE':
          return 'bg-blue-50 text-blue-700 border-blue-200';
        case 'COMPONENT_PREPARATION':
          return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'CANCELLED':
          return 'bg-red-50 text-red-700 border-red-200';
        default:
          return 'bg-slate-100 text-slate-600 border-slate-200';
      }
    },
    getCompletedStagesCount(p) {
      if (!p || !p.assemblyStages) return 0;
      return p.assemblyStages.filter(s => s.status === 'COMPLETED').length;
    },
    getAssemblyProgress(p) {
      if (!p || !p.assemblyStages || p.assemblyStages.length === 0) return 0;
      const count = p.assemblyStages.filter(s => s.status === 'COMPLETED').length;
      return Math.round((count / p.assemblyStages.length) * 100);
    }
  },
  async mounted() {
    try {
      await Promise.all([
        this.protoStore.fetchDashboard(),
        this.protoStore.fetchPrototypes()
      ]);
    } catch (err) {
      console.error('[PrototypeDashboardView] mounted error:', err);
    }
  }
};
</script>
