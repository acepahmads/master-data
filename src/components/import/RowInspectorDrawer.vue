<template>
  <div v-if="isOpen && row" class="fixed inset-y-0 right-0 z-50 w-full max-w-xl bg-white shadow-2xl border-l border-slate-200 flex flex-col transform transition-transform duration-300">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
      <div>
        <div class="flex items-center gap-2">
          <span class="font-mono font-bold text-xs px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
            Row #{{ row.rowNumber }}
          </span>
          <span class="text-xs font-semibold text-slate-500 font-mono">
            {{ row.fileName }} / {{ row.sheetName }}
          </span>
        </div>
        <h3 class="text-sm font-bold text-slate-900 mt-1 truncate max-w-md">
          {{ row.normalizedName || 'Unclassified Row' }}
        </h3>
      </div>
      <button @click="$emit('close')" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200 transition-colors">
        <AppIcon name="x" size="sm" />
      </button>
    </div>

    <!-- Scrollable Body -->
    <div class="flex-1 overflow-y-auto p-6 space-y-6 text-xs">
      <!-- Validation & Diagnostics -->
      <div>
        <h4 class="font-bold text-slate-700 text-3xs uppercase tracking-wider mb-2">Validation Status & Diagnostics</h4>
        <div class="p-3 rounded-xl border" :class="validationClass">
          <div class="flex items-center gap-2 font-bold mb-1">
            <span class="uppercase tracking-wider text-3xs font-mono">{{ row.validationStatus }}</span>
          </div>
          <ul v-if="errors && errors.length > 0" class="space-y-1 text-rose-700 text-2xs font-medium">
            <li v-for="(err, i) in errors" :key="'e-'+i" class="flex items-center gap-1.5">
              <span>✕</span> <span>{{ err }}</span>
            </li>
          </ul>
          <ul v-if="warnings && warnings.length > 0" class="space-y-1 text-amber-700 text-2xs font-medium mt-1">
            <li v-for="(w, i) in warnings" :key="'w-'+i" class="flex items-center gap-1.5">
              <span>▲</span> <span>{{ w }}</span>
            </li>
          </ul>
          <div v-if="(!errors || errors.length === 0) && (!warnings || warnings.length === 0)" class="text-emerald-700 text-2xs font-medium">
            ✓ All mandatory fields, numeric types, and currencies validated successfully.
          </div>
        </div>
      </div>

      <!-- AI Co-Pilot Suggestion -->
      <div v-if="row.aiConfidence > 0" class="p-4 rounded-xl bg-purple-50/70 border border-purple-200/80">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-1.5">
            <AppIcon name="sparkles" size="xs" class="text-purple-600" />
            <span class="font-bold text-purple-900 text-3xs uppercase tracking-wider">AI Classification Suggestion</span>
          </div>
          <span class="px-2 py-0.5 rounded-full text-3xs font-mono font-bold bg-purple-100 text-purple-800 border border-purple-300">
            {{ row.aiConfidence }}% Confidence
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2 text-2xs text-purple-950 font-medium mb-2">
          <div>Item Type: <strong class="font-bold">{{ row.aiSuggestedType }}</strong></div>
          <div>Product Type: <strong class="font-bold">{{ row.aiSuggestedProdType || 'N/A' }}</strong></div>
          <div class="col-span-2">Suggested Category: <strong class="font-bold">{{ row.aiSuggestedCategory || 'General' }}</strong></div>
        </div>
        <div class="text-3xs text-purple-700 italic border-t border-purple-200/60 pt-2">
          Advisory suggestion only. Human approval is strictly mandatory before master commit.
        </div>
      </div>

      <!-- Duplicate Detection -->
      <div>
        <h4 class="font-bold text-slate-700 text-3xs uppercase tracking-wider mb-2">Duplicate & Master Data Matching</h4>
        <div class="p-3 rounded-xl border" :class="duplicateClass">
          <div class="flex items-center justify-between font-mono font-bold text-2xs mb-1">
            <span>Status: {{ row.duplicateStatus }}</span>
            <span v-if="row.matchedMasterId" class="text-slate-600">ID: {{ row.matchedMasterId }}</span>
          </div>
          <div v-if="duplicates && duplicates.length > 0" class="space-y-1.5 mt-2">
            <div v-for="(cand, idx) in duplicates" :key="idx" class="p-2 bg-white rounded-lg border border-slate-200 text-2xs">
              <div class="font-bold text-slate-800">{{ cand.name }} ({{ cand.code }})</div>
              <div class="text-3xs text-slate-500 font-mono">Matched: {{ cand.matchType }} in {{ cand.masterType }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Classification & Editable Form -->
      <div class="space-y-3 bg-slate-50 p-4 rounded-xl border border-slate-200">
        <h4 class="font-bold text-slate-700 text-3xs uppercase tracking-wider">Classification & Master Destination</h4>
        
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label text-3xs">Item Classification</label>
            <select v-model="editForm.itemClassification" class="form-select text-xs">
              <option value="PRODUCT">PRODUCT (Finished Device)</option>
              <option value="COMPONENT">COMPONENT (Electronic Part)</option>
              <option value="SPARE_PART">SPARE_PART</option>
              <option value="ACCESSORY">ACCESSORY</option>
              <option value="SERVICE">SERVICE (Installation/Labor)</option>
              <option value="CONSUMABLE">CONSUMABLE</option>
              <option value="OTHER">OTHER</option>
            </select>
          </div>

          <div v-if="editForm.itemClassification === 'PRODUCT'">
            <label class="form-label text-3xs">Product Type</label>
            <select v-model="editForm.productType" class="form-select text-xs">
              <option value="TRADING">TRADING (Commercial Sourced)</option>
              <option value="PROJECT">PROJECT (Turnkey Solution)</option>
              <option value="MANUFACTURE">MANUFACTURE (In-House Manufactured)</option>
            </select>
          </div>
        </div>

        <div>
          <label class="form-label text-3xs">Normalized Name</label>
          <input type="text" v-model="editForm.normalizedName" class="form-input text-xs" />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="form-label text-3xs">Product Code / MPN</label>
            <input type="text" v-model="editForm.normalizedCode" class="form-input text-xs font-mono" />
          </div>
          <div>
            <label class="form-label text-3xs">Category</label>
            <input type="text" v-model="editForm.normalizedCategory" class="form-input text-xs" />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-2">
          <div>
            <label class="form-label text-3xs">Purchase Cost</label>
            <input type="number" step="0.01" v-model.number="editForm.unitCost" class="form-input text-xs font-mono" />
          </div>
          <div>
            <label class="form-label text-3xs">Selling Price</label>
            <input type="number" step="0.01" v-model.number="editForm.sellingPrice" class="form-input text-xs font-mono" />
          </div>
          <div>
            <label class="form-label text-3xs">Currency</label>
            <select v-model="editForm.currency" class="form-select text-xs font-mono">
              <option value="IDR">IDR</option>
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="SGD">SGD</option>
              <option value="CNY">CNY</option>
            </select>
          </div>
        </div>

        <div>
          <label class="form-label text-3xs">Scope Description & Sub-Points</label>
          <textarea
            v-model="editForm.description"
            rows="5"
            class="form-input text-xs font-mono whitespace-pre-line leading-relaxed"
            placeholder="Detailed scope description, checklists, and sub-points..."
          ></textarea>
        </div>

        <div>
          <label class="form-label text-3xs">Review Notes / Rationale</label>
          <textarea v-model="editForm.reviewNotes" rows="2" class="form-input text-xs" placeholder="Add auditor review notes..."></textarea>
        </div>
      </div>

      <!-- Raw Excel Cell Snapshot (Immutable) -->
      <div>
        <div class="flex items-center justify-between mb-2">
          <h4 class="font-bold text-slate-700 text-3xs uppercase tracking-wider">Immutable Source Cell Snapshot</h4>
          <span class="text-3xs text-slate-400 font-mono">Preserved for Audit</span>
        </div>
        <div class="bg-slate-900 text-slate-100 p-3 rounded-xl font-mono text-3xs overflow-x-auto">
          <pre>{{ formattedRawData }}</pre>
        </div>
      </div>

      <!-- Provenance Envelope -->
      <div class="bg-slate-50 p-3 rounded-xl border border-slate-200 text-3xs font-mono text-slate-600 space-y-1">
        <div><strong>Batch ID:</strong> {{ row.batchId }}</div>
        <div><strong>Source File:</strong> {{ row.fileName }}</div>
        <div><strong>Sheet / Row:</strong> {{ row.sheetName }} : Row #{{ row.rowNumber }}</div>
        <div><strong>Review Status:</strong> {{ row.reviewDecision }} ({{ row.reviewedBy || 'Unassigned' }})</div>
      </div>
    </div>

    <!-- Footer Actions -->
    <div class="p-4 border-t border-slate-200 bg-slate-50 flex items-center justify-between">
      <button @click="handleReject" class="btn text-xs bg-rose-50 text-rose-700 hover:bg-rose-100 border border-rose-200">
        Reject Row
      </button>

      <div class="flex items-center gap-2">
        <button @click="handleSaveEdit" class="btn btn-secondary text-xs">
          Save Edits
        </button>
        <button @click="handleApprove" class="btn btn-primary text-xs">
          <AppIcon name="check" size="xs" class="mr-1" />
          <span>Approve Row</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import AppIcon from '../common/AppIcon.vue';

export default {
  name: 'RowInspectorDrawer',
  components: { AppIcon },
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    row: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      editForm: {
        itemClassification: 'PRODUCT',
        productType: 'TRADING',
        normalizedName: '',
        normalizedCode: '',
        normalizedCategory: '',
        description: '',
        unitCost: 0,
        sellingPrice: 0,
        currency: 'IDR',
        reviewNotes: ''
      }
    };
  },
  computed: {
    errors() {
      try {
        const val = this.row?.validationErrorsJson;
        if (!val || val === 'null' || val === 'undefined') return [];
        const res = JSON.parse(val);
        return Array.isArray(res) ? res : [];
      } catch (e) {
        return [];
      }
    },
    warnings() {
      try {
        const val = this.row?.validationWarningsJson;
        if (!val || val === 'null' || val === 'undefined') return [];
        const res = JSON.parse(val);
        return Array.isArray(res) ? res : [];
      } catch (e) {
        return [];
      }
    },
    duplicates() {
      try {
        const val = this.row?.duplicateCandidatesJson;
        if (!val || val === 'null' || val === 'undefined') return [];
        const res = JSON.parse(val);
        return Array.isArray(res) ? res : [];
      } catch (e) {
        return [];
      }
    },
    formattedRawData() {
      try {
        const val = this.row?.rawDataJson;
        if (!val || val === 'null') return '{}';
        const obj = JSON.parse(val);
        return JSON.stringify(obj, null, 2);
      } catch (e) {
        return this.row?.rawDataJson || '{}';
      }
    },
    validationClass() {
      if (this.row?.validationStatus === 'ERROR') return 'bg-rose-50 border-rose-200 text-rose-800';
      if (this.row?.validationStatus === 'WARNING') return 'bg-amber-50 border-amber-200 text-amber-800';
      return 'bg-emerald-50 border-emerald-200 text-emerald-800';
    },
    duplicateClass() {
      if (this.row?.duplicateStatus === 'MATCHED_EXISTING') return 'bg-blue-50 border-blue-200 text-blue-900';
      if (this.row?.duplicateStatus === 'POSSIBLE_DUPLICATE') return 'bg-amber-50 border-amber-200 text-amber-900';
      return 'bg-slate-50 border-slate-200 text-slate-800';
    }
  },
  watch: {
    row: {
      immediate: true,
      handler(val) {
        if (val) {
          this.editForm = {
            itemClassification: val.itemClassification || 'PRODUCT',
            productType: val.productType || 'TRADING',
            normalizedName: val.normalizedName || '',
            normalizedCode: val.normalizedCode || '',
            normalizedCategory: val.normalizedCategory || '',
            description: val.description || '',
            unitCost: val.unitCost || 0,
            sellingPrice: val.sellingPrice || 0,
            currency: val.currency || 'IDR',
            reviewNotes: val.reviewNotes || ''
          };
        }
      }
    }
  },
  methods: {
    handleSaveEdit() {
      this.$emit('triage', {
        ...this.editForm,
        decision: 'EDITED'
      });
    },
    handleApprove() {
      this.$emit('triage', {
        ...this.editForm,
        decision: 'APPROVED'
      });
    },
    handleReject() {
      this.$emit('triage', {
        ...this.editForm,
        decision: 'REJECTED'
      });
    }
  }
};
</script>
