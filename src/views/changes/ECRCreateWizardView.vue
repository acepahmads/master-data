<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <div class="flex items-center gap-2 text-xs text-slate-400 mb-1">
          <router-link to="/ecr" class="hover:text-slate-200">ECRs</router-link>
          <span>/</span>
          <span class="text-slate-200">New Request</span>
        </div>
        <h1 class="text-2xl font-bold text-slate-100">Create Engineering Change Request (ECR)</h1>
      </div>
    </div>

    <!-- Wizard Form Card -->
    <form @submit.prevent="submitForm" class="p-6 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-6 shadow-xl">
      <!-- Section 1: Origin & Scope -->
      <div class="space-y-4">
        <h3 class="text-sm font-bold uppercase tracking-wider text-cyan-400 border-b border-slate-800 pb-2">
          1. Origin & Target Revision
        </h3>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Change Origin *</label>
            <select
              v-model="form.origin"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
              required
            >
              <option value="TEST_FINDING">Test Finding / Defect</option>
              <option value="CUSTOMER_REQUIREMENT">Customer Requirement</option>
              <option value="SUPPLIER_CHANGE">Supplier / PCN Notice</option>
              <option value="OBSOLESCENCE">Component Obsolescence</option>
              <option value="COST_REDUCTION">Cost Reduction Initiative</option>
              <option value="RELIABILITY">Reliability Enhancement</option>
              <option value="REGULATORY">Regulatory & Certification</option>
              <option value="OTHER">Other Engineering Need</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Source Hardware Revision *</label>
            <input
              v-model="form.sourceHardwareRevisionId"
              type="text"
              placeholder="e.g. REV-001 or ID"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none font-mono"
              required
            />
          </div>
        </div>

        <div v-if="form.origin === 'TEST_FINDING'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Source Test Finding ID</label>
            <input
              v-model="form.sourceFindingId"
              type="text"
              placeholder="e.g. FND-2024-001 or ID"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none font-mono"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Product ID</label>
            <input
              v-model="form.productId"
              type="text"
              placeholder="Associated Product ID"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none font-mono"
            />
          </div>
        </div>
      </div>

      <!-- Section 2: Change Classification -->
      <div class="space-y-4">
        <h3 class="text-sm font-bold uppercase tracking-wider text-cyan-400 border-b border-slate-800 pb-2">
          2. Classification & Priority
        </h3>

        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1.5">ECR Title *</label>
          <input
            v-model="form.title"
            type="text"
            placeholder="e.g. Replace LDO with Buck Converter for Thermal Optimization"
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            required
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Change Type *</label>
            <select
              v-model="form.changeType"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            >
              <option value="HARDWARE_DESIGN">Hardware Design / Schematic</option>
              <option value="BOM_COMPONENT">BOM Component Substitution</option>
              <option value="PCB_LAYOUT">PCB Layout / Routing</option>
              <option value="MECHANICAL">Mechanical / Enclosure</option>
              <option value="FIRMWARE">Firmware / Embedded Software</option>
              <option value="PACKAGING">Packaging / Labelling</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Priority *</label>
            <select
              v-model="form.priority"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            >
              <option value="CRITICAL">Critical (Line Stoppage / Safety)</option>
              <option value="HIGH">High (Design Block / Regression)</option>
              <option value="MEDIUM">Medium (Planned Optimization)</option>
              <option value="LOW">Low (Cosmetic / Minor Enhancement)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Assigned Engineering Lead</label>
            <input
              v-model="form.assignedLead"
              type="text"
              placeholder="e.g. Lead Hardware Engineer"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            />
          </div>
        </div>
      </div>

      <!-- Section 3: Technical Details -->
      <div class="space-y-4">
        <h3 class="text-sm font-bold uppercase tracking-wider text-cyan-400 border-b border-slate-800 pb-2">
          3. Technical Description & Justification
        </h3>

        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1.5">Change Description *</label>
          <textarea
            v-model="form.description"
            rows="3"
            placeholder="Detailed description of the proposed engineering change..."
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            required
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Business Reason</label>
            <textarea
              v-model="form.businessReason"
              rows="2"
              placeholder="Why is this change necessary from a business / product perspective?"
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Technical Justification</label>
            <textarea
              v-model="form.technicalJustification"
              rows="2"
              placeholder="Technical analysis, test data evidence, or root cause justification..."
              class="w-full px-3.5 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-800">
        <router-link
          to="/ecr"
          class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition-colors"
        >
          Cancel
        </router-link>
        <button
          type="submit"
          :disabled="loading"
          class="px-5 py-2 rounded-xl bg-cyan-500 hover:bg-cyan-400 text-slate-950 text-xs font-bold shadow-lg shadow-cyan-500/20 transition-all disabled:opacity-50"
        >
          {{ loading ? 'Creating...' : 'Create ECR Draft' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';

export default {
  name: 'ECRCreateWizardView',
  data() {
    return {
      loading: false,
      form: {
        title: '',
        origin: 'TEST_FINDING',
        sourceFindingId: '',
        productId: '',
        sourceHardwareRevisionId: 'REV-001',
        changeType: 'HARDWARE_DESIGN',
        priority: 'HIGH',
        description: '',
        businessReason: '',
        technicalJustification: '',
        assignedLead: ''
      }
    };
  },
  computed: {
    changeStore() {
      return useChangeStore();
    }
  },
  mounted() {
    const q = this.$route.query;
    if (q.findingId) {
      this.form.sourceFindingId = q.findingId;
      this.form.origin = 'TEST_FINDING';
    }
    if (q.productId) {
      this.form.productId = q.productId;
    }
    if (q.revisionId) {
      this.form.sourceHardwareRevisionId = q.revisionId;
    }
    if (q.title) {
      this.form.title = q.title;
    }
  },
  methods: {
    async submitForm() {
      this.loading = true;
      try {
        const ecr = await this.changeStore.createECR(this.form);
        this.$router.push(`/ecr/${ecr.id}`);
      } catch (err) {
        alert(err.response?.data?.message || err.message || 'Failed to create ECR');
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>
