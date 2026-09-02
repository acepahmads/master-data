<template>
  <div v-if="eco" class="space-y-6">
    <!-- Header / Status Banner -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div class="flex items-center gap-3">
            <span class="font-mono text-base font-extrabold text-cyan-300 bg-cyan-950/80 px-3 py-1 rounded-xl border border-cyan-800/60">
              {{ eco.code }}
            </span>
            <ChangeStatusBadge :status="eco.status" />
          </div>
          <h1 class="text-xl font-bold text-slate-100 mt-2">
            {{ eco.title }}
          </h1>
          <p class="text-xs text-slate-400 mt-1">
            Source ECR:
            <router-link v-if="eco.ecr" :to="`/ecr/${eco.ecr.id}`" class="text-cyan-400 hover:underline font-mono font-bold">
              {{ eco.ecr.code }} ({{ eco.ecr.title }})
            </router-link>
            <span v-else>—</span>
            &bull; Implementation Owner: <strong class="text-slate-300">{{ eco.implementationOwner }}</strong>
          </p>
        </div>

        <!-- Lifecycle Actions -->
        <div class="flex flex-wrap items-center gap-2.5">
          <!-- Step 1: Approve ECO -->
          <button
            v-if="eco.status === 'DRAFT'"
            @click="approveECO"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold shadow-lg shadow-emerald-500/20 transition-all disabled:opacity-50"
          >
            Approve ECO for Execution
          </button>

          <!-- Step 2: Implement ECO (Atomic Baseline Generation) -->
          <button
            v-if="eco.status === 'APPROVED'"
            @click="implementECO"
            :disabled="actionLoading"
            class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-cyan-500 to-emerald-500 hover:from-cyan-400 hover:to-emerald-400 text-slate-950 text-xs font-extrabold shadow-lg shadow-cyan-500/30 transition-all flex items-center gap-2 disabled:opacity-50"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            Implement ECO & Spawn Revision
          </button>

          <!-- Step 3: Verify Regression Results -->
          <button
            v-if="eco.status === 'IMPLEMENTED'"
            @click="verifyECO"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-purple-500 hover:bg-purple-400 text-slate-950 text-xs font-bold shadow-lg shadow-purple-500/20 transition-all flex items-center gap-2 disabled:opacity-50"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Verify Regression Results
          </button>

          <!-- Step 4: Close ECO -->
          <button
            v-if="eco.status === 'VERIFIED'"
            @click="closeECO"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-slate-700 hover:bg-slate-600 text-slate-200 text-xs font-bold transition-all disabled:opacity-50"
          >
            Close ECO
          </button>
        </div>
      </div>

      <!-- Execution Progress Stepper -->
      <div class="grid grid-cols-2 md:grid-cols-5 gap-2 pt-3 border-t border-slate-800 text-[11px]">
        <div class="p-2.5 rounded-xl border" :class="getStepClass('DRAFT', ['DRAFT', 'APPROVED', 'IMPLEMENTED', 'VERIFIED', 'CLOSED'])">
          <div class="font-bold">1. Order Created</div>
          <div class="text-[10px] text-slate-400">DRAFT</div>
        </div>
        <div class="p-2.5 rounded-xl border" :class="getStepClass('APPROVED', ['APPROVED', 'IMPLEMENTED', 'VERIFIED', 'CLOSED'])">
          <div class="font-bold">2. CCB Release</div>
          <div class="text-[10px] text-slate-400">APPROVED</div>
        </div>
        <div class="p-2.5 rounded-xl border" :class="getStepClass('IMPLEMENTED', ['IMPLEMENTED', 'VERIFIED', 'CLOSED'])">
          <div class="font-bold">3. Baseline Cloned</div>
          <div class="text-[10px] text-slate-400">IMPLEMENTED</div>
        </div>
        <div class="p-2.5 rounded-xl border" :class="getStepClass('VERIFIED', ['VERIFIED', 'CLOSED'])">
          <div class="font-bold">4. Regressions Passed</div>
          <div class="text-[10px] text-slate-400">VERIFIED</div>
        </div>
        <div class="p-2.5 rounded-xl border" :class="getStepClass('CLOSED', ['CLOSED'])">
          <div class="font-bold">5. Closed-Loop Done</div>
          <div class="text-[10px] text-slate-400">CLOSED</div>
        </div>
      </div>
    </div>

    <!-- Generated Baseline Summary Card (Shown when IMPLEMENTED or beyond) -->
    <div
      v-if="eco.targetHardwareRevision"
      class="p-6 rounded-2xl bg-emerald-950/20 border border-emerald-800/40 space-y-3"
    >
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="p-2 rounded-xl bg-emerald-500/20 text-emerald-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </span>
          <div>
            <h3 class="text-base font-bold text-slate-100">
              New Engineering Baseline Generated: <span class="font-mono text-emerald-400">{{ eco.targetHardwareRevision.code }}</span>
            </h3>
            <p class="text-xs text-slate-400">
              Historical baseline ({{ eco.sourceHardwareRevision ? eco.sourceHardwareRevision.code : 'REV-A' }}) remains immutable and untouched.
            </p>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-3 pt-2 text-xs">
        <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800">
          <span class="text-slate-500 uppercase text-[10px] font-semibold">Revision Name</span>
          <div class="font-semibold text-slate-200 mt-0.5">{{ eco.targetHardwareRevision.name }}</div>
        </div>
        <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800">
          <span class="text-slate-500 uppercase text-[10px] font-semibold">Target BOM</span>
          <div class="font-mono font-bold text-cyan-300 mt-0.5">{{ eco.targetBOM ? eco.targetBOM.bomCode : 'BOM-' + eco.targetHardwareRevision.code }}</div>
        </div>
        <div class="p-3 rounded-xl bg-slate-900/80 border border-slate-800">
          <span class="text-slate-500 uppercase text-[10px] font-semibold">Downstream Action</span>
          <div class="font-semibold text-purple-400 mt-0.5">Prototype Build & Test Regression Mandated</div>
        </div>
      </div>
    </div>

    <!-- Section: BOM Delta & Preview Engine -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-cyan-400"></span>
            Bill of Materials (BOM) Delta Preview Engine
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">
            Automated line-by-line diff of the proposed target BOM baseline against the source BOM.
          </p>
        </div>
      </div>

      <BOMChangeDiffTable :summary="previewSummary" />
    </div>
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';
import ChangeStatusBadge from '@/components/changes/ChangeStatusBadge.vue';
import BOMChangeDiffTable from '@/components/changes/BOMChangeDiffTable.vue';

export default {
  name: 'ECODetailView',
  components: {
    ChangeStatusBadge,
    BOMChangeDiffTable
  },
  data() {
    return {
      actionLoading: false
    };
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    eco() {
      return this.changeStore.currentECO;
    },
    previewSummary() {
      return this.changeStore.bomChangePreview || {};
    }
  },
  async mounted() {
    await this.loadData();
  },
  methods: {
    async loadData() {
      const id = this.$route.params.id;
      await Promise.all([
        this.changeStore.fetchECOById(id),
        this.changeStore.fetchChangePreview(id)
      ]);
    },
    getStepClass(stepKey, activeStates) {
      if (!this.eco) return 'bg-slate-900/40 border-slate-800 text-slate-500';
      if (activeStates.includes(this.eco.status)) {
        if (this.eco.status === stepKey) {
          return 'bg-cyan-950/40 border-cyan-800/80 text-cyan-300 shadow-sm shadow-cyan-950/30';
        }
        return 'bg-emerald-950/30 border-emerald-800/40 text-emerald-300';
      }
      return 'bg-slate-900/40 border-slate-800 text-slate-500';
    },
    async approveECO() {
      if (confirm(`Approve Change Order ${this.eco.code} for execution?`)) {
        this.actionLoading = true;
        try {
          await this.changeStore.approveECO(this.eco.id);
          await this.loadData();
        } catch (err) {
          alert(err.response?.data?.message || err.message || 'Failed to approve ECO');
        } finally {
          this.actionLoading = false;
        }
      }
    },
    async implementECO() {
      // Guardrail #2: Atomic transaction creating new revision & cloning BOM
      if (confirm(`Execute implementation of ${this.eco.code}? This will spawn a new target Hardware Revision and clone the BOM.`)) {
        this.actionLoading = true;
        try {
          await this.changeStore.implementECO(this.eco.id);
          await this.loadData();
        } catch (err) {
          alert(err.response?.data?.message || err.message || 'Failed to implement ECO');
        } finally {
          this.actionLoading = false;
        }
      }
    },
    async verifyECO() {
      const notes = prompt('Enter verification notes / test run summary reference:');
      if (notes !== null) {
        this.actionLoading = true;
        try {
          await this.changeStore.verifyECO(this.eco.id, { notes });
          await this.loadData();
        } catch (err) {
          alert(err.response?.data?.message || err.message || 'Failed to verify ECO');
        } finally {
          this.actionLoading = false;
        }
      }
    },
    async closeECO() {
      if (confirm(`Formally close out ${this.eco.code}? All closed-loop verification requirements have been satisfied.`)) {
        this.actionLoading = true;
        try {
          await this.changeStore.closeECO(this.eco.id);
          await this.loadData();
        } catch (err) {
          alert(err.response?.data?.message || err.message || 'Failed to close ECO');
        } finally {
          this.actionLoading = false;
        }
      }
    }
  }
};
</script>
