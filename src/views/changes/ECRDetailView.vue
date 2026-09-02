<template>
  <div v-if="ecr" class="space-y-6">
    <!-- Header / Status Banner -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div class="flex items-center gap-3">
            <span class="font-mono text-base font-extrabold text-cyan-400 bg-cyan-950/80 px-3 py-1 rounded-xl border border-cyan-800/60">
              {{ ecr.code }}
            </span>
            <ChangeStatusBadge :status="ecr.status" />
            <ChangePriorityBadge :priority="ecr.priority" />
          </div>
          <h1 class="text-xl font-bold text-slate-100 mt-2">
            {{ ecr.title }}
          </h1>
          <p class="text-xs text-slate-400 mt-1">
            Source Revision: <strong class="text-slate-200 font-mono">{{ ecr.sourceHardwareRevision ? ecr.sourceHardwareRevision.code : 'REV-001' }}</strong>
            &bull; Requested By: <strong class="text-slate-300">{{ ecr.requestedBy }}</strong>
            &bull; Assigned Lead: <strong class="text-slate-300">{{ ecr.assignedLead || 'Unassigned' }}</strong>
          </p>
        </div>

        <!-- Action Control Buttons -->
        <div class="flex flex-wrap items-center gap-2.5">
          <button
            v-if="ecr.status === 'DRAFT'"
            @click="submitForReview"
            class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 text-xs font-bold shadow-lg shadow-amber-500/20 transition-all"
          >
            Submit for Technical Review
          </button>

          <button
            v-if="ecr.status === 'APPROVED'"
            @click="openConvertToECOModal"
            class="px-4 py-2 rounded-xl bg-cyan-500 hover:bg-cyan-400 text-slate-950 text-xs font-bold shadow-lg shadow-cyan-500/20 transition-all flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
            Create ECO
          </button>

          <router-link
            v-if="ecr.eco"
            :to="`/eco/${ecr.eco.id}`"
            class="px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold shadow-lg shadow-indigo-600/20 transition-all flex items-center gap-2"
          >
            <span>View ECO ({{ ecr.eco.code }})</span>
            &rarr;
          </router-link>
        </div>
      </div>

      <!-- Description & Business Justification -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-3 border-t border-slate-800 text-xs">
        <div>
          <span class="text-slate-500 font-semibold uppercase tracking-wider text-[10px]">Description</span>
          <p class="text-slate-300 mt-1 leading-relaxed">{{ ecr.description }}</p>
        </div>
        <div>
          <span class="text-slate-500 font-semibold uppercase tracking-wider text-[10px]">Business Reason & Technical Justification</span>
          <p class="text-slate-300 mt-1 leading-relaxed">{{ ecr.businessReason || '—' }}</p>
          <p v-if="ecr.technicalJustification" class="text-slate-400 italic mt-1">{{ ecr.technicalJustification }}</p>
        </div>
      </div>
    </div>

    <!-- Section 1: Change Items Breakdown -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-cyan-400"></span>
            Proposed Change Items ({{ (ecr.changeItems || []).length }})
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">Line-item hardware, component, or schematic modifications.</p>
        </div>
        <button
          v-if="ecr.status === 'DRAFT' || ecr.status === 'UNDER_REVIEW'"
          @click="showAddItemModal = true"
          class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors"
        >
          + Add Change Item
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <ChangeItemCard
          v-for="item in ecr.changeItems || []"
          :key="item.id"
          :item="item"
          :can-delete="ecr.status === 'DRAFT' || ecr.status === 'UNDER_REVIEW'"
          @delete="deleteItem"
        />
        <div v-if="(!ecr.changeItems || ecr.changeItems.length === 0)" class="col-span-2 text-center py-6 text-xs text-slate-500">
          No change items added yet.
        </div>
      </div>
    </div>

    <!-- Section 2: 5-Domain Impact Analysis -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-amber-400"></span>
            Cross-Functional Impact Analysis
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">Evaluation across Cost, Schedule, Inventory Scrap, Testing, and Compliance.</p>
        </div>
        <button
          @click="showImpactModal = true"
          class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors"
        >
          Edit Impact Data
        </button>
      </div>

      <ChangeImpactMatrix :impacts="ecr.impactAnalyses || []" />
    </div>

    <!-- Section 3: CCB Multi-Tier Approval Workflow -->
    <div class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-slate-200 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
            Change Control Board (CCB) Sign-Off
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">Multi-tier formal approval required prior to Engineering Change Order generation.</p>
        </div>
      </div>

      <ApprovalFlowTimeline
        :approvals="ecr.approvals || []"
        :can-approve="ecr.status === 'UNDER_REVIEW' || ecr.status === 'PENDING_APPROVAL'"
        @open-approval-dialog="openApprovalModal"
      />
    </div>

    <!-- Add Change Item Modal -->
    <div v-if="showAddItemModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 space-y-4">
        <h3 class="text-base font-bold text-slate-100">Add Proposed Change Item</h3>
        <div class="space-y-3 text-xs">
          <div>
            <label class="block text-slate-400 mb-1">Change Type</label>
            <select v-model="newItem.changeType" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200">
              <option value="BOM_COMPONENT">BOM Component Substitution</option>
              <option value="COMPONENT_VALUE">Component Value / Tolerance</option>
              <option value="PCB_FOOTPRINT">PCB Footprint / Land Pattern</option>
              <option value="SCHEMATIC_NET">Schematic Netlist / Circuit</option>
              <option value="FIRMWARE_CONFIG">Firmware Pin / Register Config</option>
            </select>
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Affected Reference Designator (e.g. U3, R12)</label>
            <input v-model="newItem.affectedReference" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono" />
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Change Title *</label>
            <input v-model="newItem.title" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200" />
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Description / Reason</label>
            <textarea v-model="newItem.description" rows="2" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200"></textarea>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showAddItemModal = false" class="px-4 py-2 rounded-xl bg-slate-800 text-slate-300 text-xs">Cancel</button>
          <button @click="saveNewItem" class="px-4 py-2 rounded-xl bg-cyan-500 text-slate-950 text-xs font-bold">Add Item</button>
        </div>
      </div>
    </div>

    <!-- Impact Modal -->
    <div v-if="showImpactModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4">
        <h3 class="text-base font-bold text-slate-100">Record Domain Impact</h3>
        <div class="space-y-3 text-xs">
          <div>
            <label class="block text-slate-400 mb-1">Domain</label>
            <select v-model="impactForm.domain" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200">
              <option value="COST">Cost (Unit & Tooling)</option>
              <option value="SCHEDULE">Schedule & Milestones</option>
              <option value="INVENTORY">Inventory & Scrap</option>
              <option value="TESTING">Testing & Regression</option>
              <option value="COMPLIANCE">Compliance & Regulatory</option>
            </select>
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Impact Level</label>
            <select v-model="impactForm.impactLevel" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200">
              <option value="NONE">None</option>
              <option value="LOW">Low</option>
              <option value="MEDIUM">Medium</option>
              <option value="HIGH">High</option>
              <option value="CRITICAL">Critical</option>
            </select>
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Cost Delta ($)</label>
            <input v-model.number="impactForm.costDelta" type="number" step="0.01" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono" />
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Schedule Delay (Weeks)</label>
            <input v-model.number="impactForm.scheduleImpactWeeks" type="number" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 font-mono" />
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Impact Description</label>
            <textarea v-model="impactForm.description" rows="2" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200"></textarea>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showImpactModal = false" class="px-4 py-2 rounded-xl bg-slate-800 text-slate-300 text-xs">Cancel</button>
          <button @click="saveImpactData" class="px-4 py-2 rounded-xl bg-cyan-500 text-slate-950 text-xs font-bold">Save Impact</button>
        </div>
      </div>
    </div>

    <!-- CCB Approval Modal -->
    <div v-if="showApprovalModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 space-y-4">
        <h3 class="text-base font-bold text-slate-100">CCB Sign-off Decision</h3>
        <div class="space-y-3 text-xs">
          <div>
            <label class="block text-slate-400 mb-1">Stage</label>
            <input :value="approvalForm.stage" disabled class="w-full px-3 py-2 rounded-xl bg-slate-950/50 border border-slate-800 text-slate-400 font-mono" />
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Decision *</label>
            <select v-model="approvalForm.decision" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200">
              <option value="APPROVED">Approved</option>
              <option value="REQUEST_CHANGES">Request Changes</option>
              <option value="REJECTED">Rejected</option>
            </select>
          </div>
          <div>
            <label class="block text-slate-400 mb-1">Justification / Review Comments</label>
            <textarea v-model="approvalForm.justification" rows="3" placeholder="State reasons for approval or change requests..." class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200"></textarea>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showApprovalModal = false" class="px-4 py-2 rounded-xl bg-slate-800 text-slate-300 text-xs">Cancel</button>
          <button @click="submitApprovalDecision" class="px-4 py-2 rounded-xl bg-emerald-500 text-slate-950 text-xs font-bold">Submit Decision</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';
import ChangeStatusBadge from '@/components/changes/ChangeStatusBadge.vue';
import ChangePriorityBadge from '@/components/changes/ChangePriorityBadge.vue';
import ChangeImpactMatrix from '@/components/changes/ChangeImpactMatrix.vue';
import ChangeItemCard from '@/components/changes/ChangeItemCard.vue';
import ApprovalFlowTimeline from '@/components/changes/ApprovalFlowTimeline.vue';

export default {
  name: 'ECRDetailView',
  components: {
    ChangeStatusBadge,
    ChangePriorityBadge,
    ChangeImpactMatrix,
    ChangeItemCard,
    ApprovalFlowTimeline
  },
  data() {
    return {
      showAddItemModal: false,
      showImpactModal: false,
      showApprovalModal: false,
      newItem: {
        changeType: 'BOM_COMPONENT',
        affectedReference: '',
        title: '',
        description: ''
      },
      impactForm: {
        domain: 'COST',
        impactLevel: 'LOW',
        costDelta: 0.0,
        scheduleImpactWeeks: 0,
        description: ''
      },
      approvalForm: {
        stage: 'TIER_1_TECH_LEAD',
        decision: 'APPROVED',
        justification: ''
      }
    };
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    ecr() {
      return this.changeStore.currentECR;
    }
  },
  async mounted() {
    await this.loadData();
  },
  methods: {
    async loadData() {
      const id = this.$route.params.id;
      await this.changeStore.fetchECRById(id);
    },
    async submitForReview() {
      if (confirm('Submit this change request for technical review & impact analysis?')) {
        await this.changeStore.submitECR(this.ecr.id);
        await this.loadData();
      }
    },
    async deleteItem(itemId) {
      if (confirm('Remove this change item from the ECR?')) {
        await this.changeStore.deleteChangeItem(this.ecr.id, itemId);
      }
    },
    async saveNewItem() {
      if (!this.newItem.title) return alert('Title is required');
      await this.changeStore.addChangeItem(this.ecr.id, this.newItem);
      this.showAddItemModal = false;
      this.newItem = { changeType: 'BOM_COMPONENT', affectedReference: '', title: '', description: '' };
    },
    async saveImpactData() {
      await this.changeStore.saveImpact(this.ecr.id, this.impactForm);
      this.showImpactModal = false;
    },
    openApprovalModal(stageKey) {
      this.approvalForm.stage = stageKey;
      this.approvalForm.decision = 'APPROVED';
      this.approvalForm.justification = '';
      this.showApprovalModal = true;
    },
    async submitApprovalDecision() {
      await this.changeStore.submitApproval(this.ecr.id, this.approvalForm);
      this.showApprovalModal = false;
      await this.loadData();
    },
    async openConvertToECOModal() {
      // Guardrail #1 explicit conversion
      if (confirm(`Convert approved ${this.ecr.code} into an Engineering Change Order (ECO)?`)) {
        const eco = await this.changeStore.createECO({
          ecrId: this.ecr.id,
          title: `Implement ${this.ecr.code}: ${this.ecr.title}`
        });
        this.$router.push(`/eco/${eco.id}`);
      }
    }
  }
};
</script>
