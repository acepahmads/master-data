<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-primary-400 uppercase tracking-wider mb-1">
          <router-link :to="session ? `/prototypes/${session.prototypeBuildId}` : '/testing'" class="hover:underline flex items-center gap-1">
            <i class="feather-arrow-left text-xs"></i> Prototype Workspace
          </router-link>
          <span>/</span>
          <span class="text-slate-400 font-mono">{{ session ? session.sessionCode : '' }}</span>
        </div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <i class="feather-cpu text-emerald-400"></i>
          {{ session ? session.name : 'Laboratory Test Execution Workspace' }}
        </h1>
        <p class="text-sm text-slate-400 mt-0.5">
          {{ session ? session.description : 'Interactive test execution, measurements, evidence, and validation governance.' }}
        </p>
      </div>

      <!-- Quick Session KPI & Action -->
      <div class="flex items-center gap-3" v-if="session">
        <div class="px-3.5 py-1.5 rounded-lg bg-slate-800 border border-slate-700 text-xs font-semibold text-slate-300 flex items-center gap-2">
          <span>Pass Rate:</span>
          <span class="text-emerald-400 font-bold text-sm">{{ (session.passRatePercentage || 0).toFixed(0) }}%</span>
          <span class="text-slate-500">({{ session.passedCases }}/{{ session.totalCases }})</span>
        </div>

        <button
          @click="showValidationModal = true"
          class="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-bold transition shadow-lg flex items-center gap-1.5"
        >
          <i class="feather-shield"></i> Validation Sign-Off
        </button>
      </div>
    </div>

    <!-- Active Validation Decision Banner if present -->
    <div
      v-if="latestDecision"
      class="p-4 rounded-xl border flex items-center justify-between backdrop-blur-sm"
      :class="getValidationBannerClass(latestDecision.decision)"
    >
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-black/20 text-lg">
          <i :class="getValidationIcon(latestDecision.decision)"></i>
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold uppercase tracking-wider">Validation Verdict:</span>
            <span class="text-sm font-extrabold">{{ latestDecision.decision }}</span>
            <span class="text-xs opacity-75">by {{ latestDecision.decidedBy }}</span>
          </div>
          <p class="text-xs mt-0.5 opacity-90">{{ latestDecision.decisionSummary || latestDecision.justification }}</p>
        </div>
      </div>
      <button
        @click="showHistoryModal = true"
        class="px-3 py-1 bg-black/20 hover:bg-black/40 rounded text-xs font-semibold transition"
      >
        View Decision History ({{ session.decisions ? session.decisions.length : 1 }})
      </button>
    </div>

    <!-- 3-Pane Laboratory Execution Workspace -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[650px]">
      <!-- Left Pane: Test Case Navigator (3 Cols) -->
      <div class="lg:col-span-3 bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 flex flex-col backdrop-blur-sm">
        <div class="flex items-center justify-between pb-3 border-b border-slate-700/60 mb-3">
          <h2 class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
            <i class="feather-list text-primary-400"></i> Test Cases ({{ snapshotCases.length }})
          </h2>
          <span class="text-[11px] font-mono text-emerald-400 font-semibold">
            {{ session ? session.passedCases : 0 }} / {{ snapshotCases.length }}
          </span>
        </div>

        <div class="space-y-2 overflow-y-auto max-h-[580px] pr-1 flex-1">
          <div
            v-for="tc in snapshotCases"
            :key="tc.id"
            @click="selectTestCase(tc)"
            class="p-3 rounded-lg border cursor-pointer transition text-xs space-y-1.5"
            :class="selectedCase && selectedCase.id === tc.id ? 'bg-primary-950/50 border-primary-500/60 ring-1 ring-primary-500/40' : 'bg-slate-900/60 border-slate-700/40 hover:border-slate-600'"
          >
            <div class="flex items-center justify-between">
              <span class="font-mono font-bold text-[11px] text-slate-300">{{ tc.code }}</span>
              <span
                class="px-1.5 py-0.2 rounded text-[10px] font-bold"
                :class="getCaseStatusClass(getLatestExecution(tc.id)?.status)"
              >
                {{ getLatestExecution(tc.id)?.status || 'PENDING' }}
              </span>
            </div>
            <div class="font-medium text-slate-200 line-clamp-1 text-xs">{{ tc.name }}</div>
            <div class="flex items-center justify-between text-[10px] text-slate-400 font-mono">
              <span>{{ tc.category }}</span>
              <span v-if="getLatestExecution(tc.id)?.runNumber > 1" class="text-indigo-300 font-semibold">
                Run #{{ getLatestExecution(tc.id).runNumber }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Center Pane: Measurement Bench & Evaluation (6 Cols) -->
      <div class="lg:col-span-6 bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 flex flex-col justify-between backdrop-blur-sm">
        <div v-if="selectedCase" class="space-y-5">
          <!-- Case Header -->
          <div class="border-b border-slate-700/60 pb-4">
            <div class="flex items-center justify-between">
              <span class="text-xs font-mono px-2.5 py-0.5 rounded bg-primary-500/10 text-primary-300 border border-primary-500/20 font-bold">
                {{ selectedCase.code }} (Seq #{{ selectedCase.sequence }})
              </span>
              <span class="text-xs font-mono text-slate-400 font-semibold">
                Active Iteration: Run #{{ currentRunNumber }}
              </span>
            </div>
            <h2 class="text-lg font-bold text-slate-100 mt-2">{{ selectedCase.name }}</h2>
            <p class="text-xs text-slate-400 mt-1 leading-relaxed">
              {{ selectedCase.instructions || 'Follow standardized laboratory bench protocol.' }}
            </p>
          </div>

          <!-- Expected Tolerance Criteria -->
          <div class="p-4 bg-slate-900/80 rounded-xl border border-slate-700/60 space-y-2">
            <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
              <i class="feather-target text-emerald-400"></i> Parametric Specification
            </h3>
            <div class="grid grid-cols-3 gap-3 text-xs font-mono pt-1">
              <div v-if="selectedCase.testType === 'NUMERIC_RANGE'">
                <span class="text-slate-500 text-[10px] block uppercase">Min Limit</span>
                <span class="text-emerald-400 font-bold text-sm">{{ selectedCase.minimumValue }} {{ selectedCase.unit }}</span>
              </div>
              <div v-if="selectedCase.testType === 'NUMERIC_RANGE' || selectedCase.testType === 'NUMERIC_MAX'">
                <span class="text-slate-500 text-[10px] block uppercase">Max Limit</span>
                <span class="text-emerald-400 font-bold text-sm">{{ selectedCase.maximumValue }} {{ selectedCase.unit }}</span>
              </div>
              <div v-if="selectedCase.targetValue">
                <span class="text-slate-500 text-[10px] block uppercase">Target Nominal</span>
                <span class="text-primary-300 font-bold text-sm">{{ selectedCase.targetValue }} {{ selectedCase.unit }}</span>
              </div>
              <div v-if="selectedCase.testType === 'BOOLEAN_PASS_FAIL'">
                <span class="text-slate-500 text-[10px] block uppercase">Expected Logic</span>
                <span class="text-emerald-400 font-bold text-sm">{{ selectedCase.expectedBoolean ? 'TRUE' : 'FALSE' }}</span>
              </div>
            </div>
          </div>

          <!-- Measurement Entry Form -->
          <div class="p-4 bg-slate-900/40 rounded-xl border border-slate-700/40 space-y-4">
            <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
              <i class="feather-activity text-primary-400"></i> Record Measured Value
            </h3>

            <!-- Numeric Value Input -->
            <div v-if="['NUMERIC_RANGE', 'NUMERIC_MIN', 'NUMERIC_MAX'].includes(selectedCase.testType)" class="space-y-2">
              <label class="block text-xs font-medium text-slate-300">Measured Reading ({{ selectedCase.unit }})</label>
              <div class="flex items-center gap-3">
                <input
                  v-model.number="measForm.measuredValue"
                  type="number"
                  step="any"
                  placeholder="Enter numeric reading..."
                  class="flex-1 px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-lg text-slate-100 font-mono text-base font-bold focus:outline-none focus:border-primary-500"
                />
                <span class="px-3 py-2 bg-slate-800 rounded-lg font-mono text-slate-300 text-sm font-semibold border border-slate-700">
                  {{ selectedCase.unit }}
                </span>
              </div>
            </div>

            <!-- Boolean Value Toggle -->
            <div v-if="selectedCase.testType === 'BOOLEAN_PASS_FAIL'" class="space-y-2">
              <label class="block text-xs font-medium text-slate-300">Measured State</label>
              <div class="flex items-center gap-3">
                <button
                  type="button"
                  @click="measForm.measuredBoolean = true"
                  class="px-4 py-2 rounded-lg text-xs font-bold transition flex items-center gap-2"
                  :class="measForm.measuredBoolean === true ? 'bg-emerald-600 text-white' : 'bg-slate-800 text-slate-300 border border-slate-700'"
                >
                  <i class="feather-check"></i> TRUE (PASS)
                </button>
                <button
                  type="button"
                  @click="measForm.measuredBoolean = false"
                  class="px-4 py-2 rounded-lg text-xs font-bold transition flex items-center gap-2"
                  :class="measForm.measuredBoolean === false ? 'bg-rose-600 text-white' : 'bg-slate-800 text-slate-300 border border-slate-700'"
                >
                  <i class="feather-x"></i> FALSE (FAIL)
                </button>
              </div>
            </div>

            <!-- Visual / Manual Result Override -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Status Override (Optional)</label>
                <select
                  v-model="measForm.status"
                  class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-xs text-slate-200"
                >
                  <option value="">Auto-Evaluate from Value</option>
                  <option value="PASSED">Force PASSED</option>
                  <option value="FAILED">Force FAILED</option>
                  <option value="SKIPPED">Mark SKIPPED</option>
                  <option value="BLOCKED">Mark BLOCKED</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Technician Notes</label>
                <input
                  v-model="measForm.notes"
                  type="text"
                  placeholder="e.g. Clean waveform, minimal ripple"
                  class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-xs text-slate-200"
                />
              </div>
            </div>

            <!-- Optional Create Finding on Fail -->
            <div class="p-3 bg-amber-500/10 rounded-lg border border-amber-500/20 space-y-2">
              <label class="flex items-center gap-2 cursor-pointer text-xs font-semibold text-amber-300">
                <input v-model="measForm.createFinding" type="checkbox" class="rounded bg-slate-900 border-slate-700 text-amber-500" />
                <span>Log Engineering Defect Finding if Test Fails</span>
              </label>
              <div v-if="measForm.createFinding" class="grid grid-cols-2 gap-2 pt-1">
                <input
                  v-model="measForm.findingTitle"
                  type="text"
                  placeholder="Defect Title (e.g. Over-voltage on rail)"
                  class="px-2.5 py-1.5 bg-slate-900 border border-slate-700 rounded text-xs text-slate-200"
                />
                <select
                  v-model="measForm.findingSeverity"
                  class="px-2.5 py-1.5 bg-slate-900 border border-slate-700 rounded text-xs text-slate-200"
                >
                  <option value="CRITICAL">CRITICAL (Blocker)</option>
                  <option value="HIGH">HIGH (Requires Disposition)</option>
                  <option value="MEDIUM">MEDIUM</option>
                  <option value="LOW">LOW</option>
                </select>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="p-12 text-center text-slate-400 text-sm">
          Select a test case from the left navigator to record measurements.
        </div>

        <!-- Action Bar -->
        <div v-if="selectedCase" class="pt-4 border-t border-slate-700/60 flex items-center justify-between">
          <button
            @click="triggerNextRun"
            class="px-3 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-medium transition flex items-center gap-1.5"
          >
            <i class="feather-refresh-cw text-xs"></i> Create Run #{{ currentRunNumber + 1 }}
          </button>

          <button
            @click="submitMeasurement"
            :disabled="submitting"
            class="px-5 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-bold transition shadow flex items-center gap-2"
          >
            <i v-if="submitting" class="feather-loader animate-spin text-xs"></i>
            <i v-else class="feather-check-circle"></i>
            Commit Measurement & Evaluate
          </button>
        </div>
      </div>

      <!-- Right Pane: History, Evidence & Findings (3 Cols) -->
      <div class="lg:col-span-3 space-y-5">
        <!-- Run History -->
        <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 backdrop-blur-sm space-y-3">
          <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
            <i class="feather-clock text-indigo-400"></i> Execution Iterations
          </h3>

          <div v-if="selectedCaseExecutions.length === 0" class="text-xs text-slate-500 text-center py-4">
            No execution history yet.
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="ex in selectedCaseExecutions"
              :key="ex.id"
              class="p-2.5 rounded-lg bg-slate-900/60 border border-slate-700/40 text-xs space-y-1"
            >
              <div class="flex items-center justify-between font-mono">
                <span class="font-bold text-slate-200">Run #{{ ex.runNumber }}</span>
                <span class="px-1.5 py-0.2 rounded text-[10px] font-bold" :class="getCaseStatusClass(ex.status)">
                  {{ ex.status }}
                </span>
              </div>
              <div v-if="ex.measuredValue !== undefined && ex.measuredValue !== null" class="text-slate-300 font-mono text-[11px]">
                Reading: <span class="font-bold text-emerald-400">{{ ex.measuredValue }} {{ ex.unit }}</span>
              </div>
              <div v-if="ex.notes" class="text-[11px] text-slate-400 italic">"{{ ex.notes }}"</div>
              <div class="text-[10px] text-slate-500">By {{ ex.executedBy }}</div>
            </div>
          </div>
        </div>

        <!-- Evidence Attachment -->
        <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 backdrop-blur-sm space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
              <i class="feather-paperclip text-slate-400"></i> Test Evidence
            </h3>
            <button
              v-if="latestActiveExecution"
              @click="showEvidenceModal = true"
              class="text-[11px] text-primary-400 hover:underline flex items-center gap-1"
            >
              <i class="feather-plus text-xs"></i> Attach
            </button>
          </div>

          <div v-if="activeEvidences.length === 0" class="text-xs text-slate-500 text-center py-3">
            No waveforms or logs attached.
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="ev in activeEvidences"
              :key="ev.id"
              class="p-2 rounded bg-slate-900/60 border border-slate-700/40 text-xs flex items-center justify-between"
            >
              <div class="flex items-center gap-2 overflow-hidden">
                <i class="feather-file text-slate-400"></i>
                <span class="text-slate-200 truncate text-[11px]">{{ ev.fileName }}</span>
              </div>
              <span class="text-[10px] text-slate-500 font-mono">{{ ev.fileType }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Validation Sign-Off Modal -->
    <div
      v-if="showValidationModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
        <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-shield text-emerald-400"></i>
          Submit Prototype Validation Decision
        </h2>

        <form @submit.prevent="submitValidation" class="space-y-4 text-sm">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Validation Verdict *</label>
            <select
              v-model="valForm.decision"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-bold"
            >
              <option value="VALIDATED">VALIDATED (All Criteria Met)</option>
              <option value="CONDITIONALLY_VALIDATED">CONDITIONALLY VALIDATED (With Exceptions)</option>
              <option value="REJECTED">REJECTED (Failed Requirements)</option>
              <option value="PENDING_REVIEW">PENDING REVIEW</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Decision Summary *</label>
            <textarea
              v-model="valForm.decisionSummary"
              required
              rows="2"
              placeholder="Summary of test session results, pass rate, and engineering conclusion..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200"
            ></textarea>
          </div>

          <div v-if="valForm.decision === 'CONDITIONALLY_VALIDATED'">
            <label class="block text-xs font-semibold text-amber-300 mb-1">Mandatory Justification & Constraints *</label>
            <textarea
              v-model="valForm.justification"
              required
              rows="2"
              placeholder="Explain why build is accepted despite open high/medium defects (e.g. Validated for bench only)..."
              class="w-full px-3 py-2 bg-slate-900 border border-amber-500/50 rounded-lg text-slate-200"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
            <button
              type="button"
              @click="showValidationModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-bold"
            >
              Commit Decision
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Evidence Modal -->
    <div
      v-if="showEvidenceModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-md w-full p-6 space-y-4 shadow-2xl">
        <h2 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <i class="feather-paperclip text-primary-400"></i> Attach Test Evidence
        </h2>

        <div class="space-y-3 text-sm">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">File Name</label>
            <input
              v-model="evForm.fileName"
              type="text"
              placeholder="e.g. Scope_Capture_3V3_Ripple.png"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono text-xs"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">File Path / URL</label>
            <input
              v-model="evForm.filePath"
              type="text"
              placeholder="/uploads/scope_001.png"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 font-mono text-xs"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Description / Caption</label>
            <input
              v-model="evForm.description"
              type="text"
              placeholder="Oscilloscope capture at 12V inrush"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 text-xs"
            />
          </div>
        </div>

        <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
          <button
            @click="showEvidenceModal = false"
            class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium"
          >
            Cancel
          </button>
          <button
            @click="submitEvidence"
            class="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-xs font-bold"
          >
            Attach Evidence
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
  name: 'TestExecutionWorkspaceView',
  data() {
    return {
      selectedCase: null,
      showValidationModal: false,
      showHistoryModal: false,
      showEvidenceModal: false,
      submitting: false,
      measForm: {
        measuredValue: null,
        measuredBoolean: null,
        observedText: '',
        status: '',
        notes: '',
        createFinding: false,
        findingTitle: '',
        findingSeverity: 'HIGH'
      },
      valForm: {
        decision: 'VALIDATED',
        decisionSummary: '',
        justification: ''
      },
      evForm: {
        fileName: '',
        filePath: '',
        fileType: 'IMAGE',
        fileSize: 102400,
        description: ''
      }
    };
  },
  computed: {
    ...mapState(useTestingStore, {
      session: 'currentSession',
      loading: 'loading'
    }),
    sessionId() {
      return this.$route.params.id;
    },
    snapshotCases() {
      if (!this.session || !this.session.testPlanSnapshotJson) return [];
      try {
        const snap = JSON.parse(this.session.testPlanSnapshotJson);
        return snap.testCases || [];
      } catch (e) {
        return [];
      }
    },
    executions() {
      return this.session?.executions || [];
    },
    selectedCaseExecutions() {
      if (!this.selectedCase) return [];
      return this.executions
        .filter(e => e.testCaseSnapshotId === this.selectedCase.id)
        .sort((a, b) => b.runNumber - a.runNumber);
    },
    latestActiveExecution() {
      if (this.selectedCaseExecutions.length === 0) return null;
      return this.selectedCaseExecutions[0];
    },
    currentRunNumber() {
      return this.latestActiveExecution?.runNumber || 1;
    },
    activeEvidences() {
      return this.latestActiveExecution?.evidences || [];
    },
    latestDecision() {
      if (!this.session?.decisions || this.session.decisions.length === 0) return null;
      return this.session.decisions[0];
    }
  },
  async mounted() {
    await this.fetchSessionById(this.sessionId);
    if (this.snapshotCases.length > 0) {
      this.selectTestCase(this.snapshotCases[0]);
    }
  },
  methods: {
    ...mapActions(useTestingStore, [
      'fetchSessionById',
      'recordExecution',
      'createNextRun',
      'addEvidence',
      'submitValidationDecision'
    ]),
    selectTestCase(tc) {
      this.selectedCase = tc;
      const latest = this.getLatestExecution(tc.id);
      this.measForm = {
        measuredValue: latest?.measuredValue !== undefined ? latest.measuredValue : null,
        measuredBoolean: latest?.measuredBoolean !== undefined ? latest.measuredBoolean : null,
        observedText: latest?.observedText || '',
        status: '',
        notes: latest?.notes || '',
        createFinding: false,
        findingTitle: `Defect in ${tc.name}`,
        findingSeverity: 'HIGH'
      };
    },
    getLatestExecution(caseId) {
      const execs = this.executions.filter(e => e.testCaseSnapshotId === caseId);
      if (execs.length === 0) return null;
      return execs.reduce((prev, curr) => (curr.runNumber > prev.runNumber ? curr : prev));
    },
    getCaseStatusClass(status) {
      switch (status) {
        case 'PASSED':
          return 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30';
        case 'FAILED':
          return 'bg-rose-500/20 text-rose-300 border border-rose-500/30';
        case 'BLOCKED':
          return 'bg-amber-500/20 text-amber-300 border border-amber-500/30';
        case 'SKIPPED':
          return 'bg-slate-700 text-slate-400';
        default:
          return 'bg-slate-800 text-slate-400 border border-slate-700';
      }
    },
    getValidationBannerClass(dec) {
      switch (dec) {
        case 'VALIDATED':
          return 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300';
        case 'CONDITIONALLY_VALIDATED':
          return 'bg-amber-500/10 border-amber-500/30 text-amber-300';
        case 'REJECTED':
          return 'bg-rose-500/10 border-rose-500/30 text-rose-300';
        default:
          return 'bg-slate-800 border-slate-700 text-slate-300';
      }
    },
    getValidationIcon(dec) {
      switch (dec) {
        case 'VALIDATED':
          return 'feather-check-circle text-emerald-400';
        case 'CONDITIONALLY_VALIDATED':
          return 'feather-alert-triangle text-amber-400';
        case 'REJECTED':
          return 'feather-x-circle text-rose-400';
        default:
          return 'feather-clock text-slate-400';
      }
    },
    async submitMeasurement() {
      this.submitting = true;
      try {
        await this.recordExecution(this.sessionId, {
          testCaseSnapshotId: this.selectedCase.id,
          ...this.measForm
        });
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      } finally {
        this.submitting = false;
      }
    },
    async triggerNextRun() {
      try {
        await this.createNextRun(this.sessionId, this.selectedCase.id);
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
    async submitEvidence() {
      if (!this.latestActiveExecution) return;
      try {
        await this.addEvidence(this.latestActiveExecution.id, this.evForm);
        this.showEvidenceModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    },
    async submitValidation() {
      try {
        await this.submitValidationDecision(this.sessionId, this.valForm);
        this.showValidationModal = false;
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    }
  }
};
</script>
