<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-primary-400 uppercase tracking-wider mb-1">
          <router-link to="/testing" class="hover:underline flex items-center gap-1">
            <i class="feather-arrow-left text-xs"></i> Testing Command Center
          </router-link>
        </div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <i class="feather-alert-triangle text-amber-400"></i>
          Engineering Findings & Defect Registry
        </h1>
        <p class="text-sm text-slate-400 mt-0.5">
          Cross-build defect tracking, failure root cause analyses, disposition management, and future change candidates.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-amber-600 hover:bg-amber-500 text-white rounded-lg text-sm font-semibold flex items-center gap-2 transition shadow-lg shadow-amber-900/30"
        >
          <i class="feather-plus"></i>
          Log Engineering Finding
        </button>
      </div>
    </div>

    <!-- Filters & Search Bar -->
    <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 flex flex-col md:flex-row gap-4 items-center justify-between backdrop-blur-sm">
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto flex-1">
        <div class="relative w-full max-w-xs">
          <i class="feather-search absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-sm"></i>
          <input
            v-model="searchQuery"
            @input="loadFindings"
            type="text"
            placeholder="Search findings by code, title, root cause..."
            class="w-full pl-9 pr-4 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-primary-500 transition"
          />
        </div>

        <select
          v-model="filterSeverity"
          @change="loadFindings"
          class="px-3 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-xs text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">All Severities</option>
          <option value="CRITICAL">CRITICAL (Blocker)</option>
          <option value="HIGH">HIGH</option>
          <option value="MEDIUM">MEDIUM</option>
          <option value="LOW">LOW</option>
          <option value="INFO">INFO</option>
        </select>

        <select
          v-model="filterDisposition"
          @change="loadFindings"
          class="px-3 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-xs text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">All Dispositions</option>
          <option value="OPEN">OPEN (Under Investigation)</option>
          <option value="BENCH_REWORK">BENCH REWORK</option>
          <option value="DESIGN_CHANGE_RECOMMENDED">DESIGN CHANGE RECOMMENDED</option>
          <option value="ACCEPTED_AS_IS">ACCEPTED AS-IS</option>
          <option value="CLOSED">CLOSED</option>
        </select>

        <select
          v-model="filterChangeStatus"
          @change="loadFindings"
          class="px-3 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-xs text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">All Change Statuses</option>
          <option value="CANDIDATE">ECO Candidate Only</option>
          <option value="DEFERRED">Deferred Changes</option>
          <option value="ACCEPTED_RISK">Accepted Risk</option>
          <option value="NONE">None</option>
        </select>
      </div>

      <div class="text-xs text-slate-400">
        Found <span class="font-semibold text-slate-200">{{ findings.length }}</span> items
      </div>
    </div>

    <!-- Findings Table -->
    <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl overflow-hidden backdrop-blur-sm">
      <div v-if="loading && findings.length === 0" class="p-12 text-center text-slate-400 text-sm">
        <i class="feather-loader animate-spin text-2xl text-amber-400 mb-2"></i>
        <p>Loading defect registry...</p>
      </div>

      <div v-else-if="findings.length === 0" class="p-12 text-center text-slate-400 text-sm">
        <i class="feather-check-circle text-4xl text-emerald-400 mb-3 inline-block"></i>
        <h3 class="text-base font-semibold text-slate-200">No findings match the current filter</h3>
        <p class="text-xs text-slate-500 mt-1">Zero active anomalies or defects detected under selected criteria.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-300">
          <thead class="bg-slate-900/60 uppercase font-semibold text-slate-400 border-b border-slate-700/60">
            <tr>
              <th class="px-4 py-3">Code & Severity</th>
              <th class="px-4 py-3">Finding Title & Description</th>
              <th class="px-4 py-3">Category</th>
              <th class="px-4 py-3">Disposition</th>
              <th class="px-4 py-3">Change Scope</th>
              <th class="px-4 py-3">Assignee</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-700/40">
            <tr v-for="f in findings" :key="f.id" class="hover:bg-slate-750 transition">
              <td class="px-4 py-3.5 space-y-1">
                <span class="font-mono font-bold text-xs text-slate-200 block">{{ f.code }}</span>
                <span
                  class="inline-block text-[10px] font-extrabold px-2 py-0.5 rounded uppercase"
                  :class="getSeverityBadge(f.severity)"
                >
                  {{ f.severity }}
                </span>
              </td>
              <td class="px-4 py-3.5 max-w-sm">
                <div class="font-bold text-slate-100 text-sm">{{ f.title }}</div>
                <p class="text-xs text-slate-400 line-clamp-2 mt-0.5 leading-relaxed">{{ f.description }}</p>
                <div v-if="f.rootCause" class="text-[11px] text-amber-300/80 italic mt-1">
                  Root Cause: {{ f.rootCause }}
                </div>
              </td>
              <td class="px-4 py-3.5">
                <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-slate-900 text-slate-300 border border-slate-700">
                  {{ f.category || 'GENERAL' }}
                </span>
              </td>
              <td class="px-4 py-3.5">
                <span class="px-2.5 py-1 rounded-md text-[11px] font-bold block w-fit" :class="getDispositionClass(f.findingDisposition)">
                  {{ formatDisp(f.findingDisposition) }}
                </span>
              </td>
              <td class="px-4 py-3.5">
                <div v-if="f.changeCandidateStatus === 'CANDIDATE'" class="space-y-1">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-purple-500/20 text-purple-300 border border-purple-500/30">
                    ECO CANDIDATE
                  </span>
                  <span class="text-[10px] text-slate-400 block font-mono">{{ formatScope(f.recommendedChangeScope) }}</span>
                </div>
                <div v-else class="text-[11px] text-slate-500">
                  None
                </div>
              </td>
              <td class="px-4 py-3.5 text-slate-300 text-xs">
                <div>{{ f.assignedTo || 'Unassigned' }}</div>
                <span class="text-[10px] text-slate-500">By {{ f.reportedBy }}</span>
              </td>
              <td class="px-4 py-3.5 text-right space-x-1.5 whitespace-nowrap">
                <button
                  v-if="f.changeCandidateStatus === 'CANDIDATE'"
                  @click="createECRFromFinding(f)"
                  class="px-2.5 py-1.5 bg-cyan-500 hover:bg-cyan-400 text-slate-950 rounded text-xs font-bold shadow-sm transition"
                  title="Create ECR from Candidate Finding"
                >
                  Create ECR
                </button>
                <button
                  @click="openEditModal(f)"
                  class="px-2.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 rounded text-xs font-semibold transition"
                >
                  Disposition
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Finding Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl max-h-[90vh] overflow-y-auto">
        <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-alert-triangle text-amber-400"></i>
          Log Engineering Finding
        </h2>

        <form @submit.prevent="submitCreateFinding" class="space-y-4 text-sm">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Target Prototype Build ID *</label>
            <input
              v-model="createForm.prototypeBuildId"
              type="text"
              required
              placeholder="e.g. PROTO-2024-001"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono text-xs"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Finding Title *</label>
            <input
              v-model="createForm.title"
              type="text"
              required
              placeholder="e.g. Buck Inductor L1 Thermal Overheat"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Severity *</label>
              <select
                v-model="createForm.severity"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
              >
                <option value="CRITICAL">CRITICAL (Validation Blocker)</option>
                <option value="HIGH">HIGH (Major Interface)</option>
                <option value="MEDIUM">MEDIUM (Tolerance Drift)</option>
                <option value="LOW">LOW (Cosmetic / Bench)</option>
                <option value="INFO">INFO</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Category</label>
              <select
                v-model="createForm.category"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
              >
                <option value="THERMAL">Thermal</option>
                <option value="ELECTRICAL">Electrical</option>
                <option value="RF_WIRELESS">RF & Wireless</option>
                <option value="FUNCTIONAL">Functional</option>
                <option value="COMMUNICATIONS">Communications</option>
                <option value="MECHANICAL">Mechanical</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Technical Observation & Description *</label>
            <textarea
              v-model="createForm.description"
              required
              rows="3"
              placeholder="Detailed description of the observed defect or parametric violation..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Change Candidate</label>
              <select
                v-model="createForm.changeCandidateStatus"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
              >
                <option value="CANDIDATE">CANDIDATE (Mark for ECO)</option>
                <option value="NONE">NONE (Local Rework Only)</option>
                <option value="DEFERRED">DEFERRED</option>
                <option value="ACCEPTED_RISK">ACCEPTED RISK</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Recommended Scope</label>
              <select
                v-model="createForm.recommendedChangeScope"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
              >
                <option value="NONE">None</option>
                <option value="PCB_LAYOUT">PCB Layout & Copper</option>
                <option value="SCHEMATIC_CIRCUIT">Schematic Circuit</option>
                <option value="BOM_COMPONENT">BOM Component Swap</option>
                <option value="MECHANICAL_ENCLOSURE">Mechanical Enclosure</option>
                <option value="FIRMWARE_LOGIC">Firmware Logic</option>
              </select>
            </div>
          </div>

          <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-amber-600 hover:bg-amber-500 text-white rounded-lg text-xs font-bold"
            >
              Log Finding
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit / Disposition Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
        <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-edit text-primary-400"></i>
          Update Finding Disposition & Root Cause
        </h2>

        <form @submit.prevent="submitUpdateFinding" class="space-y-4 text-sm">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Disposition *</label>
              <select
                v-model="editForm.findingDisposition"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-bold"
              >
                <option value="OPEN">OPEN (Under Investigation)</option>
                <option value="BENCH_REWORK">BENCH REWORK (Fixed on Board)</option>
                <option value="DESIGN_CHANGE_RECOMMENDED">DESIGN CHANGE RECOMMENDED</option>
                <option value="ACCEPTED_AS_IS">ACCEPTED AS-IS</option>
                <option value="CLOSED">CLOSED (Resolved)</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Change Candidate</label>
              <select
                v-model="editForm.changeCandidateStatus"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
              >
                <option value="CANDIDATE">CANDIDATE</option>
                <option value="NONE">NONE</option>
                <option value="DEFERRED">DEFERRED</option>
                <option value="ACCEPTED_RISK">ACCEPTED RISK</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Root Cause Analysis</label>
            <textarea
              v-model="editForm.rootCause"
              rows="2"
              placeholder="Root cause identified by engineering team..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Containment / Resolution Notes</label>
            <textarea
              v-model="editForm.resolutionNotes"
              rows="2"
              placeholder="Rework actions taken or proposed layout fix..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
            <button
              type="button"
              @click="showEditModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-xs font-bold"
            >
              Save Disposition
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia';
import { useTestingStore } from '@/stores/testingStore';

export default {
  name: 'EngineeringFindingsView',
  data() {
    return {
      searchQuery: '',
      filterSeverity: '',
      filterDisposition: '',
      filterChangeStatus: '',
      showCreateModal: false,
      showEditModal: false,
      createForm: {
        prototypeBuildId: 'PROTO-2024-001',
        title: '',
        description: '',
        severity: 'HIGH',
        category: 'THERMAL',
        changeCandidateStatus: 'CANDIDATE',
        recommendedChangeScope: 'PCB_LAYOUT'
      },
      editForm: {
        id: '',
        findingDisposition: 'OPEN',
        changeCandidateStatus: 'CANDIDATE',
        rootCause: '',
        resolutionNotes: ''
      }
    };
  },
  computed: {
    ...mapState(useTestingStore, ['findings', 'loading'])
  },
  mounted() {
    this.loadFindings();
  },
  methods: {
    ...mapActions(useTestingStore, ['fetchFindings', 'createFinding', 'updateFinding']),
    async loadFindings() {
      await this.fetchFindings({
        q: this.searchQuery,
        severity: this.filterSeverity,
        disposition: this.filterDisposition,
        changeStatus: this.filterChangeStatus
      });
    },
    openCreateModal() {
      this.createForm = {
        prototypeBuildId: 'PROTO-2024-001',
        title: '',
        description: '',
        severity: 'HIGH',
        category: 'THERMAL',
        changeCandidateStatus: 'CANDIDATE',
        recommendedChangeScope: 'PCB_LAYOUT'
      };
      this.showCreateModal = true;
    },
    async submitCreateFinding() {
      try {
        await this.createFinding(this.createForm);
        this.showCreateModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
    openEditModal(f) {
      this.editForm = {
        id: f.id,
        findingDisposition: f.findingDisposition || 'OPEN',
        changeCandidateStatus: f.changeCandidateStatus || 'NONE',
        rootCause: f.rootCause || '',
        resolutionNotes: f.resolutionNotes || ''
      };
      this.showEditModal = true;
    },
    async submitUpdateFinding() {
      try {
        await this.updateFinding(this.editForm.id, this.editForm);
        this.showEditModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
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
    getDispositionClass(disp) {
      switch (disp) {
        case 'CLOSED':
          return 'bg-emerald-500/10 text-emerald-300 border border-emerald-500/20';
        case 'ACCEPTED_AS_IS':
          return 'bg-blue-500/10 text-blue-300 border border-blue-500/20';
        case 'DESIGN_CHANGE_RECOMMENDED':
          return 'bg-purple-500/10 text-purple-300 border border-purple-500/20';
        case 'BENCH_REWORK':
          return 'bg-indigo-500/10 text-indigo-300 border border-indigo-500/20';
        default:
          return 'bg-amber-500/10 text-amber-300 border border-amber-500/20';
      }
    },
    formatDisp(d) {
      if (!d) return 'OPEN';
      return d.replace(/_/g, ' ');
    },
    formatScope(s) {
      if (!s || s === 'NONE') return 'None';
      return s.replace(/_/g, ' ');
    },
    createECRFromFinding(f) {
      this.$router.push({
        path: '/ecr/create',
        query: {
          findingId: f.id,
          title: `Remediate ${f.code}: ${f.title}`,
          productId: f.productId || ''
        }
      });
    }
  }
};
</script>
