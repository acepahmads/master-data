<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-primary-400 uppercase tracking-wider mb-1">
          <router-link to="/test-plans" class="hover:underline flex items-center gap-1">
            <i class="feather-arrow-left text-xs"></i> All Protocols
          </router-link>
          <span>/</span>
          <span class="text-slate-400 font-mono">{{ plan ? plan.code : '' }}</span>
        </div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <i class="feather-layers text-primary-400"></i>
          {{ plan ? plan.name : 'Test Plan Protocol' }}
        </h1>
        <p class="text-sm text-slate-400 mt-0.5">
          {{ plan ? plan.description : 'Loading test plan protocol...' }}
        </p>
      </div>

      <div class="flex items-center gap-3" v-if="currentVersion">
        <button
          v-if="currentVersion.status === 'DRAFT'"
          @click="openAddCaseModal"
          class="px-3.5 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-xs font-semibold flex items-center gap-1.5 transition shadow"
        >
          <i class="feather-plus"></i> Add Test Case
        </button>

        <button
          v-if="currentVersion.status === 'DRAFT'"
          @click="openReleaseModal"
          class="px-3.5 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-semibold flex items-center gap-1.5 transition shadow"
        >
          <i class="feather-check"></i> Formally Release Version
        </button>

        <button
          v-if="currentVersion.status === 'RELEASED'"
          @click="openCloneModal"
          class="px-3.5 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-semibold flex items-center gap-1.5 transition shadow"
        >
          <i class="feather-copy"></i> Clone to New Version
        </button>
      </div>
    </div>

    <!-- Version Switcher Strip -->
    <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 flex flex-col md:flex-row items-center justify-between gap-4 backdrop-blur-sm">
      <div class="flex items-center gap-2 overflow-x-auto w-full md:w-auto">
        <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider mr-2 shrink-0">Versions:</span>
        <button
          v-for="ver in versions"
          :key="ver.id"
          @click="selectVersion(ver)"
          class="px-3 py-1.5 rounded-lg text-xs font-bold font-mono transition flex items-center gap-2 shrink-0"
          :class="currentVersion && currentVersion.id === ver.id ? 'bg-primary-600 text-white shadow' : 'bg-slate-900/80 text-slate-300 hover:bg-slate-700'"
        >
          <span>{{ ver.versionNumber }}</span>
          <span
            class="text-[10px] px-1.5 py-0.2 rounded"
            :class="ver.status === 'RELEASED' ? 'bg-emerald-500/20 text-emerald-300' : 'bg-amber-500/20 text-amber-300'"
          >
            {{ ver.status }}
          </span>
        </button>

        <button
          @click="openNewVersionModal"
          class="px-2.5 py-1.5 rounded-lg text-xs font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-700/50 transition flex items-center gap-1 shrink-0"
        >
          <i class="feather-plus text-xs"></i> New Draft
        </button>
      </div>

      <!-- Immutability Status Badge -->
      <div v-if="currentVersion">
        <div
          v-if="currentVersion.status === 'RELEASED'"
          class="px-3 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-semibold flex items-center gap-1.5"
        >
          <i class="feather-lock text-xs"></i>
          IMMUTABLE RELEASED PROTOCOL VERSION
        </div>
        <div
          v-else
          class="px-3 py-1 rounded-lg bg-amber-500/10 border border-amber-500/30 text-amber-400 text-xs font-semibold flex items-center gap-1.5"
        >
          <i class="feather-edit text-xs"></i>
          DRAFT PROTOCOL VERSION (EDITABLE)
        </div>
      </div>
    </div>

    <!-- Test Cases Table -->
    <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl overflow-hidden backdrop-blur-sm">
      <div class="p-4 border-b border-slate-700/60 flex items-center justify-between">
        <h2 class="text-sm font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-list text-primary-400"></i>
          Parametric Test Cases ({{ testCases.length }})
        </h2>
        <span class="text-xs text-slate-400">
          Version {{ currentVersion ? currentVersion.versionNumber : '' }} Specification
        </span>
      </div>

      <div v-if="testCases.length === 0" class="p-12 text-center text-slate-400 text-sm">
        <i class="feather-clipboard text-3xl text-slate-500 mb-2 inline-block"></i>
        <p>No test cases defined for this version.</p>
        <button
          v-if="currentVersion && currentVersion.status === 'DRAFT'"
          @click="openAddCaseModal"
          class="mt-3 px-3 py-1.5 bg-primary-600 hover:bg-primary-500 text-white text-xs font-medium rounded-lg transition"
        >
          Add First Test Case
        </button>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-300">
          <thead class="bg-slate-900/60 uppercase font-semibold text-slate-400 border-b border-slate-700/60">
            <tr>
              <th class="px-4 py-3 w-12 text-center">Seq</th>
              <th class="px-4 py-3">Code & Name</th>
              <th class="px-4 py-3">Category</th>
              <th class="px-4 py-3">Type</th>
              <th class="px-4 py-3">Parametric Limits / Acceptance</th>
              <th class="px-4 py-3">Instructions</th>
              <th v-if="currentVersion && currentVersion.status === 'DRAFT'" class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-700/40 font-mono">
            <tr v-for="tc in testCases" :key="tc.id" class="hover:bg-slate-750 transition font-sans">
              <td class="px-4 py-3.5 text-center font-mono font-bold text-slate-400">{{ tc.sequence }}</td>
              <td class="px-4 py-3.5">
                <span class="font-mono font-bold text-primary-300 text-xs block">{{ tc.code }}</span>
                <span class="font-semibold text-slate-200 text-sm block">{{ tc.name }}</span>
                <span v-if="tc.description" class="text-xs text-slate-400 block mt-0.5 line-clamp-1">{{ tc.description }}</span>
              </td>
              <td class="px-4 py-3.5">
                <span class="px-2 py-0.5 rounded text-[11px] font-semibold bg-slate-900 text-slate-300 border border-slate-700">
                  {{ tc.category || 'GENERAL' }}
                </span>
              </td>
              <td class="px-4 py-3.5">
                <span class="text-xs font-mono font-medium text-indigo-300">{{ formatType(tc.testType) }}</span>
              </td>
              <td class="px-4 py-3.5 font-mono">
                <div v-if="tc.testType === 'NUMERIC_RANGE'" class="text-slate-200 text-xs">
                  <span class="text-emerald-400 font-bold">{{ tc.minimumValue }}</span>
                  <span class="text-slate-400 mx-1">to</span>
                  <span class="text-emerald-400 font-bold">{{ tc.maximumValue }}</span>
                  <span class="text-primary-300 font-bold ml-1">{{ tc.unit }}</span>
                  <span v-if="tc.targetValue" class="text-slate-400 block text-[10px]">Target: {{ tc.targetValue }} {{ tc.unit }}</span>
                </div>
                <div v-else-if="tc.testType === 'NUMERIC_MIN'" class="text-slate-200 text-xs">
                  <span class="text-slate-400">&ge; </span>
                  <span class="text-emerald-400 font-bold">{{ tc.minimumValue }}</span>
                  <span class="text-primary-300 font-bold ml-1">{{ tc.unit }}</span>
                </div>
                <div v-else-if="tc.testType === 'NUMERIC_MAX'" class="text-slate-200 text-xs">
                  <span class="text-slate-400">&le; </span>
                  <span class="text-emerald-400 font-bold">{{ tc.maximumValue }}</span>
                  <span class="text-primary-300 font-bold ml-1">{{ tc.unit }}</span>
                </div>
                <div v-else-if="tc.testType === 'BOOLEAN_PASS_FAIL'" class="text-slate-200 text-xs">
                  Expected: <span class="text-emerald-400 font-bold">{{ tc.expectedBoolean ? 'TRUE' : 'FALSE' }}</span>
                </div>
                <div v-else class="text-slate-400 text-xs">
                  {{ tc.acceptanceNotes || 'Visual inspection standard' }}
                </div>
              </td>
              <td class="px-4 py-3.5 text-xs text-slate-400 max-w-xs line-clamp-2">
                {{ tc.instructions || 'Follow standard bench procedure.' }}
              </td>
              <td v-if="currentVersion && currentVersion.status === 'DRAFT'" class="px-4 py-3.5 text-right space-x-2">
                <button
                  @click="deleteCase(tc.id)"
                  class="p-1.5 rounded hover:bg-rose-500/20 text-rose-400 transition"
                  title="Delete Test Case"
                >
                  <i class="feather-trash text-sm"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Test Case Modal -->
    <div
      v-if="showCaseModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-xl w-full p-6 space-y-4 shadow-2xl max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-700/60 pb-3">
          <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
            <i class="feather-plus-circle text-primary-400"></i>
            Add Parametric Test Case
          </h2>
          <button @click="showCaseModal = false" class="text-slate-400 hover:text-slate-200">
            <i class="feather-x text-lg"></i>
          </button>
        </div>

        <form @submit.prevent="submitAddCase" class="space-y-4 text-sm">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Test Code *</label>
              <input
                v-model="caseForm.code"
                type="text"
                required
                placeholder="e.g. TC-PWR-01"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500 font-mono"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Sequence #</label>
              <input
                v-model.number="caseForm.sequence"
                type="number"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Test Case Title *</label>
            <input
              v-model="caseForm.name"
              type="text"
              required
              placeholder="e.g. 3.3V System Rail Regulation Test"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Category</label>
              <select
                v-model="caseForm.category"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
              >
                <option value="ELECTRICAL">Electrical & Power</option>
                <option value="FUNCTIONAL">Functional & Logic</option>
                <option value="RF_WIRELESS">RF & Wireless</option>
                <option value="POWER_CONSUMPTION">Power Consumption</option>
                <option value="THERMAL">Thermal Stress</option>
                <option value="COMMUNICATIONS">Communications</option>
                <option value="VISUAL_INSPECTION">Visual Inspection</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Evaluation Type</label>
              <select
                v-model="caseForm.testType"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
              >
                <option value="NUMERIC_RANGE">Numeric Range (Min - Max)</option>
                <option value="NUMERIC_MIN">Numeric Minimum (&ge; Min)</option>
                <option value="NUMERIC_MAX">Numeric Maximum (&le; Max)</option>
                <option value="BOOLEAN_PASS_FAIL">Boolean (True / False)</option>
                <option value="VISUAL_INSPECTION">Visual Inspection</option>
              </select>
            </div>
          </div>

          <!-- Parametric Fields if Numeric -->
          <div
            v-if="['NUMERIC_RANGE', 'NUMERIC_MIN', 'NUMERIC_MAX'].includes(caseForm.testType)"
            class="grid grid-cols-3 gap-3 p-3 bg-slate-900/60 rounded-lg border border-slate-700/50"
          >
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Min Value</label>
              <input
                v-model.number="caseForm.minimumValue"
                type="number"
                step="any"
                placeholder="3.15"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Max Value</label>
              <input
                v-model.number="caseForm.maximumValue"
                type="number"
                step="any"
                placeholder="3.45"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Unit</label>
              <input
                v-model="caseForm.unit"
                type="text"
                placeholder="V, uA, degC"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Step-by-Step Procedure Instructions</label>
            <textarea
              v-model="caseForm.instructions"
              rows="3"
              placeholder="Instructions for lab technician to execute this test..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
            <button
              type="button"
              @click="showCaseModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium transition"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-xs font-semibold transition"
            >
              Add Test Case
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Clone Version Modal -->
    <div
      v-if="showCloneModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-md w-full p-6 space-y-4 shadow-2xl">
        <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-copy text-indigo-400"></i>
          Clone to New Draft Version
        </h2>
        <p class="text-xs text-slate-400">
          This creates an editable new draft version by copying all test cases and parameters from {{ currentVersion ? currentVersion.versionNumber : '' }}.
        </p>

        <div class="space-y-3">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">New Version Number *</label>
            <input
              v-model="cloneForm.newVersionNumber"
              type="text"
              placeholder="e.g. v1.1"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Change Description</label>
            <textarea
              v-model="cloneForm.description"
              rows="2"
              placeholder="Notes on what is changing in this version..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            ></textarea>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
          <button
            @click="showCloneModal = false"
            class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium transition"
          >
            Cancel
          </button>
          <button
            @click="submitClone"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-semibold transition"
          >
            Clone Version
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia';
import { useTestingStore } from '@/stores/testingStore';

export default {
  name: 'TestPlanDetailView',
  data() {
    return {
      showCaseModal: false,
      showCloneModal: false,
      caseForm: {
        code: '',
        sequence: 1,
        name: '',
        category: 'ELECTRICAL',
        testType: 'NUMERIC_RANGE',
        unit: 'V',
        minimumValue: null,
        maximumValue: null,
        targetValue: null,
        instructions: ''
      },
      cloneForm: {
        newVersionNumber: 'v1.1',
        description: ''
      }
    };
  },
  computed: {
    ...mapState(useTestingStore, {
      plan: 'currentTestPlan',
      currentVersion: 'currentVersion',
      loading: 'loading'
    }),
    planId() {
      return this.$route.params.id;
    },
    versions() {
      return this.plan?.versions || [];
    },
    testCases() {
      return this.currentVersion?.testCases || [];
    }
  },
  async mounted() {
    await this.fetchTestPlanById(this.planId);
  },
  methods: {
    ...mapActions(useTestingStore, [
      'fetchTestPlanById',
      'createTestCase',
      'deleteTestCase',
      'releaseVersion',
      'cloneVersion'
    ]),
    selectVersion(ver) {
      const store = useTestingStore();
      store.currentVersion = ver;
    },
    openAddCaseModal() {
      this.caseForm = {
        code: `TC-${this.plan?.category ? this.plan.category.slice(0, 3) : 'TST'}-${String(this.testCases.length + 1).padStart(2, '0')}`,
        sequence: this.testCases.length + 1,
        name: '',
        category: this.plan?.category || 'ELECTRICAL',
        testType: 'NUMERIC_RANGE',
        unit: 'V',
        minimumValue: null,
        maximumValue: null,
        targetValue: null,
        instructions: ''
      };
      this.showCaseModal = true;
    },
    async submitAddCase() {
      try {
        await this.createTestCase(this.currentVersion.id, this.caseForm);
        this.showCaseModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
    async deleteCase(caseId) {
      if (confirm('Delete this test case?')) {
        try {
          await this.deleteTestCase(caseId);
        } catch (err) {
          alert(err.response?.data?.message || err.message);
        }
      }
    },
    async openReleaseModal() {
      if (confirm('Formally RELEASE this test plan version? It will become strictly IMMUTABLE.')) {
        try {
          await this.releaseVersion(this.currentVersion.id);
        } catch (err) {
          alert(err.response?.data?.message || err.message);
        }
      }
    },
    openCloneModal() {
      const nextVer = `v1.${this.versions.length}`;
      this.cloneForm = {
        newVersionNumber: nextVer,
        description: `Cloned from ${this.currentVersion?.versionNumber}`
      };
      this.showCloneModal = true;
    },
    async submitClone() {
      try {
        await this.cloneVersion(this.currentVersion.id, this.cloneForm);
        this.showCloneModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
    formatType(t) {
      if (!t) return '';
      return t.replace(/_/g, ' ');
    }
  }
};
</script>
