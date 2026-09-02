<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- PAGE HEADER -->
    <div class="flex items-center justify-between pb-4 border-b border-slate-200">
      <div>
        <div class="flex items-center gap-2">
          <router-link to="/prototypes" class="text-slate-400 hover:text-slate-700">
            <AppIcon name="chevron-left" size="sm" />
          </router-link>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <span>New Prototype Build Wizard</span>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
              PHASE 4
            </span>
          </h1>
        </div>
        <p class="text-xs text-slate-500 mt-1 font-medium">
          Create an immutable prototype build run and freeze the Engineering BOM snapshot for assembly tracking.
        </p>
      </div>

      <router-link to="/prototypes/list" class="btn btn-secondary text-xs">
        Cancel
      </router-link>
    </div>

    <!-- 4-STEP WIZARD PROGRESS BAR -->
    <div class="app-card p-4 bg-white shadow-card">
      <div class="grid grid-cols-4 gap-2">
        <div
          v-for="s in steps"
          :key="s.number"
          :class="[
            'flex items-center gap-2.5 p-2 rounded-lg transition-all',
            currentStep === s.number ? 'bg-indigo-50/80 border border-indigo-200 text-indigo-900' :
            currentStep > s.number ? 'text-emerald-700' : 'text-slate-400'
          ]"
        >
          <div
            :class="[
              'w-7 h-7 rounded-full flex items-center justify-center font-mono font-bold text-xs shrink-0',
              currentStep === s.number ? 'bg-indigo-600 text-white' :
              currentStep > s.number ? 'bg-emerald-500 text-white' : 'bg-slate-100 text-slate-400'
            ]"
          >
            <AppIcon v-if="currentStep > s.number" name="check" size="xs" />
            <span v-else>{{ s.number }}</span>
          </div>
          <div class="min-w-0 hidden sm:block">
            <div class="text-3xs font-mono uppercase font-bold text-slate-400">Step 0{{ s.number }}</div>
            <div class="text-xs font-bold truncate leading-tight">{{ s.name }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- WIZARD STEP CONTENT -->
    <div class="app-card p-6 bg-white space-y-6 shadow-card">
      <!-- STEP 1: SELECT R&D PRODUCT, VERSION, HARDWARE REVISION -->
      <div v-if="currentStep === 1" class="space-y-4">
        <div class="pb-3 border-b border-slate-100">
          <h2 class="text-sm font-bold text-slate-900">Step 1: Select R&D Engineering Baseline</h2>
          <p class="text-xs text-slate-500 mt-0.5">
            Select the R&D Product, Product Version, and approved Hardware Revision to instantiate this prototype build.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <!-- Product Selector -->
          <div>
            <label class="form-label">1. R&D Product <span class="text-red-500">*</span></label>
            <select v-model="form.productId" @change="onProductChange" class="form-select w-full">
              <option value="">Select R&D Product</option>
              <option v-for="prod in rndProducts" :key="prod.id" :value="prod.id">
                {{ prod.code }} — {{ prod.name }}
              </option>
            </select>
          </div>

          <!-- Version Selector -->
          <div>
            <label class="form-label">2. Product Version <span class="text-red-500">*</span></label>
            <select v-model="form.versionId" @change="onVersionChange" class="form-select w-full" :disabled="!form.productId">
              <option value="">Select Product Version</option>
              <option v-for="ver in availableVersions" :key="ver.id" :value="ver.id">
                {{ ver.versionNumber }} ({{ ver.name }})
              </option>
            </select>
          </div>

          <!-- Revision Selector -->
          <div>
            <label class="form-label">3. Hardware Revision <span class="text-red-500">*</span></label>
            <select v-model="form.hardwareRevisionId" @change="onRevisionChange" class="form-select w-full" :disabled="!form.versionId">
              <option value="">Select Hardware Revision</option>
              <option v-for="rev in availableRevisions" :key="rev.id" :value="rev.id">
                {{ rev.code }} — {{ rev.name }} ({{ rev.status }})
              </option>
            </select>
          </div>
        </div>

        <!-- Selected Revision Context Card -->
        <div v-if="selectedRevisionObj" class="p-4 rounded-xl bg-slate-50 border border-slate-200 space-y-2 mt-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="font-mono font-bold text-xs px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 border border-emerald-200">
                {{ selectedRevisionObj.code }}
              </span>
              <span class="font-bold text-xs text-slate-800">{{ selectedRevisionObj.name }}</span>
            </div>
            <span class="text-3xs font-mono font-bold text-slate-500">
              PCB: {{ selectedRevisionObj.pcbRevision || 'PCB-v1.0' }} | Engineer: {{ selectedRevisionObj.engineer || 'Marcus Vance' }}
            </span>
          </div>
          <p class="text-3xs text-slate-600">
            {{ selectedRevisionObj.changeSummary || 'Approved hardware baseline ready for physical prototype instantiation.' }}
          </p>
        </div>
      </div>

      <!-- STEP 2: PROTOTYPE METADATA -->
      <div v-if="currentStep === 2" class="space-y-4">
        <div class="pb-3 border-b border-slate-100">
          <h2 class="text-sm font-bold text-slate-900">Step 2: Prototype Build Information</h2>
          <p class="text-xs text-slate-500 mt-0.5">
            Configure unit identification, build sequence number, planned quantity, and engineering ownership.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="form-label">Prototype Code <span class="text-red-500">*</span></label>
            <input v-model="form.code" type="text" class="form-input font-mono font-bold text-indigo-700 w-full" />
          </div>
          <div>
            <label class="form-label">Prototype Build Name <span class="text-red-500">*</span></label>
            <input v-model="form.name" type="text" placeholder="e.g. Functional Test Rig #1" class="form-input w-full" />
          </div>
          <div>
            <label class="form-label">Build Number <span class="text-red-500">*</span></label>
            <input v-model.number="form.buildNumber" type="number" min="1" max="99" class="form-input font-mono w-full" />
          </div>
          <div>
            <label class="form-label">Build Quantity (Units) <span class="text-red-500">*</span></label>
            <input v-model.number="form.quantity" type="number" min="1" max="500" class="form-input font-mono font-bold w-full" />
          </div>
          <div>
            <label class="form-label">Engineering Owner <span class="text-red-500">*</span></label>
            <select v-model="form.engineeringOwner" class="form-select w-full">
              <option value="Marcus Vance">Marcus Vance (Lead Hardware Architect)</option>
              <option value="Elena Rostova">Elena Rostova (Senior RF/EMC Engineer)</option>
              <option value="Alex Chen">Alex Chen (R&D Director)</option>
              <option value="Sarah Jenkins">Sarah Jenkins (Senior Firmware Lead)</option>
            </select>
          </div>
          <div>
            <label class="form-label">Target Build Date <span class="text-red-500">*</span></label>
            <input v-model="form.targetBuildDate" type="date" class="form-input font-mono w-full" />
          </div>
        </div>
      </div>

      <!-- STEP 3: FROZEN BOM SNAPSHOT PREVIEW -->
      <div v-if="currentStep === 3" class="space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div>
            <h2 class="text-sm font-bold text-slate-900">Step 3: BOM Snapshot Preview</h2>
            <p class="text-xs text-slate-500 mt-0.5">
              Review the immutable BOM snapshot that will be frozen and attached to this prototype build.
            </p>
          </div>
          <span class="text-3xs font-mono font-bold px-2.5 py-1 rounded-full bg-amber-50 text-amber-800 border border-amber-300 flex items-center gap-1 shadow-2xs">
            <AppIcon name="lock" size="2xs" />
            <span>BOM SNAPSHOT — READ ONLY</span>
          </span>
        </div>

        <!-- BOM Stats Banner -->
        <div class="grid grid-cols-3 gap-3 p-3 bg-slate-900 text-white rounded-xl text-xs font-mono">
          <div>
            <span class="text-3xs text-slate-400 block uppercase">Source Revision:</span>
            <span class="font-bold text-emerald-400">{{ form.hardwareRevisionCode || 'REV-A' }}</span>
          </div>
          <div>
            <span class="text-3xs text-slate-400 block uppercase">Snapshot Components:</span>
            <span class="font-bold text-white">{{ previewBOMItems.length }} Line Items</span>
          </div>
          <div>
            <span class="text-3xs text-slate-400 block uppercase">Estimated Unit BOM:</span>
            <span class="font-bold text-emerald-400">${{ previewBOMCost.toFixed(2) }} USD</span>
          </div>
        </div>

        <!-- Table of Frozen Line Items -->
        <div class="overflow-x-auto border border-slate-200 rounded-xl">
          <table class="app-table w-full text-left text-xs">
            <thead class="bg-slate-50 text-3xs font-mono text-slate-600 uppercase border-b border-slate-200">
              <tr>
                <th class="py-2.5 px-3">Component / MPN</th>
                <th class="py-2.5 px-3">Category</th>
                <th class="py-2.5 px-3 text-center">Qty / Board</th>
                <th class="py-2.5 px-3">Ref Designators</th>
                <th class="py-2.5 px-3 text-right">Unit Cost</th>
                <th class="py-2.5 px-3 text-right">Line Total</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 font-mono text-3xs">
              <tr v-for="itm in previewBOMItems" :key="itm.id" class="hover:bg-slate-50">
                <td class="py-2 px-3">
                  <div class="font-bold text-slate-900 font-sans text-xs">{{ itm.componentName }}</div>
                  <div class="text-slate-400">{{ itm.mpn }}</div>
                </td>
                <td class="py-2 px-3 font-sans">{{ itm.category }}</td>
                <td class="py-2 px-3 text-center font-bold">{{ itm.quantity }} {{ itm.uom }}</td>
                <td class="py-2 px-3 text-indigo-700 font-bold">{{ itm.refDesignators }}</td>
                <td class="py-2 px-3 text-right text-slate-700">${{ itm.unitCost.toFixed(2) }}</td>
                <td class="py-2 px-3 text-right font-bold text-slate-900">${{ itm.lineCost.toFixed(2) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <p class="text-3xs text-slate-400 font-mono italic">
          * Notice: Once this prototype build is created, this BOM snapshot will be permanently frozen. Any future changes to the active Engineering BOM will not alter this prototype baseline.
        </p>
      </div>

      <!-- STEP 4: BUILD CONFIGURATION & OBJECTIVES -->
      <div v-if="currentStep === 4" class="space-y-4">
        <div class="pb-3 border-b border-slate-100">
          <h2 class="text-sm font-bold text-slate-900">Step 4: Build Purpose & Engineering Objectives</h2>
          <p class="text-xs text-slate-500 mt-0.5">
            Specify the validation purpose, experimental objectives, and initial bench notes.
          </p>
        </div>

        <div class="space-y-4 text-xs">
          <div>
            <label class="form-label">Build Purpose <span class="text-red-500">*</span></label>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-2 pt-1">
              <button
                v-for="p in purposeOptions"
                :key="p"
                type="button"
                @click="form.purpose = p"
                :class="[
                  'p-3 rounded-lg border text-left text-xs transition-all font-semibold',
                  form.purpose === p ? 'bg-indigo-50 border-indigo-600 text-indigo-900 ring-2 ring-indigo-500/20' : 'bg-slate-50 border-slate-200 text-slate-700 hover:bg-white'
                ]"
              >
                {{ p }}
              </button>
            </div>
          </div>

          <div>
            <label class="form-label">Engineering Objective & Scope</label>
            <textarea
              v-model="form.objective"
              rows="3"
              placeholder="e.g. Verify RF power output across sub-GHz bands and validate sleep current on new 3.3V rail..."
              class="form-textarea w-full"
            ></textarea>
          </div>

          <div>
            <label class="form-label">Initial Assembly / Bench Notes</label>
            <textarea
              v-model="initialNote"
              rows="2"
              placeholder="e.g. Special hand-soldering required for SMA connector ground plane..."
              class="form-textarea w-full"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- WIZARD FOOTER NAVIGATION -->
      <div class="flex items-center justify-between pt-4 border-t border-slate-100">
        <button
          v-if="currentStep > 1"
          type="button"
          @click="currentStep--"
          class="btn btn-secondary text-xs"
        >
          <AppIcon name="chevron-left" size="xs" class="mr-1" />
          <span>Back</span>
        </button>
        <div v-else></div>

        <div class="flex items-center gap-2">
          <button
            v-if="currentStep < 4"
            type="button"
            @click="nextStep"
            class="btn btn-primary text-xs"
          >
            <span>Continue</span>
            <AppIcon name="chevron-right" size="xs" class="ml-1" />
          </button>

          <button
            v-else
            type="button"
            @click="submitPrototype"
            class="btn btn-primary text-xs bg-emerald-600 hover:bg-emerald-700 border-none shadow-sm flex items-center gap-1.5"
          >
            <AppIcon name="check" size="xs" />
            <span>Create Prototype Build</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { usePrototypeStore } from '@/stores/prototypeStore';
import { useMasterStore } from '@/stores/masterStore';
import { useVersionStore } from '@/stores/versionStore';
import { bomApi } from '@/api/bom';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'PrototypeCreateWizardView',
  components: {
    AppIcon
  },
  data() {
    return {
      currentStep: 1,
      steps: [
        { number: 1, name: 'Select Revision' },
        { number: 2, name: 'Build Info' },
        { number: 3, name: 'BOM Snapshot' },
        { number: 4, name: 'Objectives' }
      ],
      form: {
        productId: '',
        productCode: '',
        productName: '',
        versionId: '',
        versionCode: '',
        versionName: '',
        hardwareRevisionId: '',
        hardwareRevisionCode: '',
        hardwareRevisionName: '',
        code: '',
        name: '',
        buildNumber: 1,
        quantity: 5,
        status: 'COMPONENT_PREPARATION',
        purpose: 'Functional Prototype',
        objective: '',
        engineeringOwner: 'Marcus Vance',
        targetBuildDate: new Date(Date.now() + 14 * 86400000).toISOString().substring(0, 10)
      },
      initialNote: '',
      realBOMItems: [],
      loadingBOM: false,
      submitting: false,
      purposeOptions: [
        'Functional Prototype',
        'Engineering Validation (EVT)',
        'Design Validation (DVT)',
        'PCB / RF Validation',
        'Mechanical Validation',
        'Pre-Production Prototype'
      ]
    };
  },
  computed: {
    protoStore() {
      return usePrototypeStore();
    },
    masterStore() {
      return useMasterStore();
    },
    versionStore() {
      return useVersionStore();
    },
    rndProducts() {
      return this.masterStore.products.filter(p => p.productType === 'RND' || !p.productType);
    },
    availableVersions() {
      if (!this.form.productId) return [];
      const prod = this.rndProducts.find(p => p.id === this.form.productId || p.code === this.form.productId);
      if (!prod) return [];
      return this.versionStore.getVersionsForProduct(prod.code);
    },
    availableRevisions() {
      if (!this.form.versionId) return [];
      return this.versionStore.getRevisionsForVersion(this.form.versionId);
    },
    selectedRevisionObj() {
      if (!this.form.hardwareRevisionId) return null;
      return this.versionStore.revisions.find(r => r.id === this.form.hardwareRevisionId || r.code === this.form.hardwareRevisionId) || null;
    },
    previewBOMItems() {
      if (this.realBOMItems && this.realBOMItems.length > 0) {
        return this.realBOMItems;
      }
      return [
        { id: 'SNP-ITM-01', componentId: 'CMP-001', componentName: 'ESP32-S3-WROOM-1-N16R8', mpn: 'ESP32-S3-WROOM-1', category: 'Microcontrollers', quantity: 1, uom: 'PCS', unitCost: 3.45, lineCost: 3.45, refDesignators: 'U1' },
        { id: 'SNP-ITM-02', componentId: 'CMP-002', componentName: 'SX1262 LoRa Transceiver IC', mpn: 'SX1262IMLTRT', category: 'Wireless & RF', quantity: 1, uom: 'PCS', unitCost: 4.80, lineCost: 4.80, refDesignators: 'U2' },
        { id: 'SNP-ITM-03', componentId: 'CMP-004', componentName: 'TPS54302 Synchronous Step-Down Converter', mpn: 'TPS54302DDCR', category: 'Power Management', quantity: 1, uom: 'PCS', unitCost: 1.15, lineCost: 1.15, refDesignators: 'U3' },
        { id: 'SNP-ITM-04', componentId: 'CMP-006', componentName: 'Multi-layer Ceramic Capacitor 10uF 50V', mpn: 'CL21A106KOQNNNE', category: 'Passive Components', quantity: 10, uom: 'PCS', unitCost: 0.08, lineCost: 0.80, refDesignators: 'C1-C10' },
        { id: 'SNP-ITM-05', componentId: 'CMP-007', componentName: 'Thick Film Chip Resistor 10k 1%', mpn: 'RC0603FR-0710KL', category: 'Passive Components', quantity: 15, uom: 'PCS', unitCost: 0.02, lineCost: 0.30, refDesignators: 'R1-R15' }
      ];
    },
    previewBOMCost() {
      return this.previewBOMItems.reduce((acc, itm) => acc + (itm.lineCost || (itm.quantity * itm.unitCost) || 0), 0);
    }
  },
  async mounted() {
    if (this.masterStore.products.length === 0) {
      await this.masterStore.fetchProducts();
    }
    if (this.versionStore.versions.length === 0) {
      await this.versionStore.fetchVersions();
      await this.versionStore.fetchRevisions();
    }

    // Check query params e.g. /prototypes/create?revision=REV-001
    const revQuery = this.$route.query.revision;
    if (revQuery) {
      const rev = this.versionStore.revisions.find(r => r.id === revQuery || r.code === revQuery);
      if (rev) {
        const ver = this.versionStore.versions.find(v => v.id === rev.versionId || v.versionNumber === rev.versionId);
        if (ver) {
          const prod = this.rndProducts.find(p => p.code === ver.productCode || p.id === ver.productCode);
          if (prod) {
            this.form.productId = prod.id;
            this.form.productCode = prod.code;
            this.form.productName = prod.name;
          }
          this.form.versionId = ver.id;
          this.form.versionCode = ver.versionNumber;
          this.form.versionName = ver.name;
        }
        this.form.hardwareRevisionId = rev.id;
        this.form.hardwareRevisionCode = rev.code;
        this.form.hardwareRevisionName = rev.name;
        this.generateCode();
      }
    } else if (this.rndProducts.length > 0) {
      this.form.productId = this.rndProducts[0].id;
      this.onProductChange();
    }
  },
  methods: {
    onProductChange() {
      const prod = this.rndProducts.find(p => p.id === this.form.productId);
      if (prod) {
        this.form.productCode = prod.code;
        this.form.productName = prod.name;
      }
      if (this.availableVersions.length > 0) {
        this.form.versionId = this.availableVersions[0].id;
        this.onVersionChange();
      } else {
        this.form.versionId = '';
        this.form.hardwareRevisionId = '';
      }
    },
    onVersionChange() {
      const ver = this.availableVersions.find(v => v.id === this.form.versionId);
      if (ver) {
        this.form.versionCode = ver.versionNumber;
        this.form.versionName = ver.name;
      }
      if (this.availableRevisions.length > 0) {
        this.form.hardwareRevisionId = this.availableRevisions[0].id;
        this.onRevisionChange();
      } else {
        this.form.hardwareRevisionId = '';
      }
    },
    async onRevisionChange() {
      const rev = this.availableRevisions.find(r => r.id === this.form.hardwareRevisionId);
      if (rev) {
        this.form.hardwareRevisionCode = rev.code;
        this.form.hardwareRevisionName = rev.name;
      }
      this.generateCode();
      await this.loadRevisionBOM();
    },
    async loadRevisionBOM() {
      if (!this.form.hardwareRevisionId) return;
      this.loadingBOM = true;
      try {
        const res = await bomApi.getByRevision(this.form.hardwareRevisionId);
        if (res && res.data && res.data.items) {
          this.realBOMItems = res.data.items.map((itm, idx) => ({
            id: itm.id || `SNP-ITM-${idx+1}`,
            componentId: itm.componentId,
            componentName: itm.name || itm.component?.name || 'Component',
            mpn: itm.componentPartNumber || itm.component?.partNumber || '',
            category: itm.componentCategory || itm.component?.category || '',
            quantity: itm.quantity || 1,
            uom: itm.unitCode || 'PCS',
            unitCost: itm.unitCost || 0,
            lineCost: itm.lineCost || (itm.quantity * itm.unitCost) || 0,
            refDesignators: itm.referenceDesignators || ''
          }));
        }
      } catch (e) {
        console.warn('[PrototypeCreateWizardView] Could not load live revision BOM:', e);
      } finally {
        this.loadingBOM = false;
      }
    },
    generateCode() {
      const count = this.protoStore.prototypes.length + 1;
      this.form.code = `PROTO-2024-${String(count).padStart(3, '0')}`;
      if (!this.form.name) {
        this.form.name = `${this.form.productCode || 'PRD'} ${this.form.hardwareRevisionCode || 'REV-A'} Build #${this.form.buildNumber}`;
      }
    },
    nextStep() {
      if (this.currentStep === 1) {
        if (!this.form.productId || !this.form.versionId || !this.form.hardwareRevisionId) {
          alert('Please select an R&D Product, Version, and Hardware Revision.');
          return;
        }
        if (!this.form.code) this.generateCode();
      } else if (this.currentStep === 2) {
        if (!this.form.code || !this.form.name || !this.form.quantity) {
          alert('Please fill all required prototype build information.');
          return;
        }
      }
      this.currentStep++;
    },
    async submitPrototype() {
      this.submitting = true;
      try {
        const created = await this.protoStore.createPrototype({
          code: this.form.code,
          name: this.form.name,
          hardwareRevisionId: this.form.hardwareRevisionId,
          buildNumber: this.form.buildNumber,
          quantity: this.form.quantity,
          purpose: this.form.purpose,
          objective: this.form.objective,
          status: 'PLANNED',
          engineeringOwner: this.form.engineeringOwner,
          targetCompletionAt: this.form.targetBuildDate,
          initialNote: this.initialNote
        });

        if (created) {
          this.masterStore.showToast('Prototype Build Created', `Prototype ${created.code} successfully created with frozen BOM snapshot.`, 'success');
          this.$router.push(`/prototypes/${created.id}`);
        }
      } catch (err) {
        alert('Failed to create prototype build: ' + (err.response?.data?.message || err.message));
      } finally {
        this.submitting = false;
      }
    }
  }
};
</script>
