<template>
  <div class="max-w-4xl mx-auto space-y-5 pb-16">
    <!-- Header -->
    <PageHeader
      :title="isEdit ? `Edit Component: ${form.partNumber}` : 'Register New Engineering Component'"
      :subtitle="isEdit ? 'Update parametric specifications, datasheets, and master attributes' : 'Add a new electronic or mechanical component to the central library'"
    >
      <template #actions>
        <button type="button" @click="cancel" class="btn btn-secondary">
          Cancel
        </button>
        <button type="button" @click="save('Draft')" class="btn btn-secondary">
          Save Draft
        </button>
        <button type="button" @click="save(form.status || 'Active')" class="btn btn-primary">
          <AppIcon name="check" size="xs" class="mr-1.5" />
          <span>{{ isEdit ? 'Save Changes' : 'Save Component' }}</span>
        </button>
      </template>
    </PageHeader>

    <form @submit.prevent="save(form.status || 'Active')" class="space-y-5">
      <!-- SECTION 01: BASIC INFORMATION -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="component" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Section 01 — Basic Information</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Part Number (MPN) -->
          <div>
            <label class="form-label">Manufacturer Part Number (MPN) *</label>
            <input
              v-model="form.partNumber"
              type="text"
              required
              placeholder="e.g. STM32H753BIT6"
              class="form-input font-mono font-bold"
            />
          </div>

          <!-- Internal Code (IPN) -->
          <div>
            <label class="form-label">Internal Part Number (IPN) *</label>
            <input
              v-model="form.internalCode"
              type="text"
              required
              placeholder="e.g. IPN-MCU-001"
              class="form-input font-mono"
            />
          </div>

          <!-- Component Name -->
          <div>
            <label class="form-label">Component Display Name *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="e.g. Arm Cortex-M7 480MHz High-Performance MCU"
              class="form-input font-medium"
            />
          </div>

          <!-- Category -->
          <div>
            <label class="form-label">Component Category *</label>
            <select v-model="form.category" required class="form-select">
              <option value="" disabled>Select Category</option>
              <option v-for="cat in store.flatCategories" :key="cat.id" :value="cat.rawName">
                {{ cat.name }}
              </option>
            </select>
          </div>

          <!-- Manufacturer -->
          <div>
            <label class="form-label">Manufacturer *</label>
            <select v-model="form.manufacturer" required class="form-select">
              <option value="" disabled>Select Manufacturer</option>
              <option v-for="m in store.manufacturers" :key="m.id" :value="m.name">
                {{ m.name }} ({{ m.country }})
              </option>
            </select>
          </div>

          <!-- Supplier -->
          <div>
            <label class="form-label">Primary Supplier *</label>
            <select v-model="form.supplier" required class="form-select">
              <option value="" disabled>Select Supplier</option>
              <option v-for="s in store.suppliers" :key="s.id" :value="s.name">
                {{ s.name }}
              </option>
            </select>
          </div>

          <!-- Unit of Measure -->
          <div>
            <label class="form-label">Standard Unit of Measure (UoM) *</label>
            <select v-model="form.unit" required class="form-select">
              <option v-for="u in store.units" :key="u.id" :value="u.code">
                {{ u.code }} — {{ u.name }} ({{ u.category }})
              </option>
            </select>
          </div>

          <!-- Lifecycle Status -->
          <div>
            <label class="form-label">Lifecycle Status *</label>
            <select v-model="form.status" required class="form-select">
              <option value="Active">Active</option>
              <option value="In Review">In Review</option>
              <option value="Draft">Draft</option>
              <option value="EOL">EOL (End of Life)</option>
            </select>
          </div>

          <!-- Package / Case -->
          <div>
            <label class="form-label">Package / Footprint</label>
            <input
              v-model="form.packageType"
              type="text"
              placeholder="e.g. LQFP-208, QFN-32, SOT-23"
              class="form-input font-mono"
            />
          </div>

          <!-- Mounting Type -->
          <div>
            <label class="form-label">Mounting Type</label>
            <select v-model="form.mountingType" class="form-select">
              <option value="Surface Mount (SMD)">Surface Mount (SMD)</option>
              <option value="Through Hole (THT)">Through Hole (THT)</option>
              <option value="Panel Mount">Panel Mount</option>
              <option value="Chassis Mount">Chassis Mount</option>
            </select>
          </div>

          <!-- Purchase Price (Harga Beli) -->
          <div>
            <label class="form-label flex items-center justify-between">
              <span>Purchase Price (Harga Beli)</span>
              <span class="text-3xs text-slate-400 font-mono">Per Unit / UoM</span>
            </label>
            <div class="grid grid-cols-5 gap-2">
              <input
                v-model.number="form.estimatedUnitCost"
                type="number"
                step="any"
                min="0"
                placeholder="0.00"
                class="form-input col-span-3 font-mono font-bold text-xs"
              />
              <select v-model="form.currency" class="form-select col-span-2 font-mono text-xs">
                <option value="IDR">IDR (Rp)</option>
                <option value="USD">USD ($)</option>
                <option value="EUR">EUR (€)</option>
                <option value="SGD">SGD ($)</option>
                <option value="CNY">CNY (¥)</option>
              </select>
            </div>
          </div>

          <!-- Purchase Link (Link Pembelian) -->
          <div>
            <label class="form-label flex items-center justify-between">
              <span>Purchase Link (Link Pembelian)</span>
              <a
                v-if="form.purchaseLink"
                :href="form.purchaseLink"
                target="_blank"
                rel="noopener noreferrer"
                class="text-3xs text-brand-600 hover:text-brand-800 font-bold flex items-center gap-0.5"
                title="Test open purchase link"
              >
                <span>Test Link</span>
                <AppIcon name="external-link" size="2xs" />
              </a>
            </label>
            <div class="relative">
              <input
                v-model="form.purchaseLink"
                type="url"
                placeholder="e.g. https://www.digikey.com/... or Tokopedia / Shopee link"
                class="form-input text-xs pr-8"
              />
              <span class="absolute right-2.5 top-2.5 text-slate-400">
                <AppIcon name="shopping-bag" size="xs" />
              </span>
            </div>
          </div>
        </div>

        <!-- Compliance Toggles -->
        <div class="mt-4 pt-3 border-t border-slate-100 flex flex-wrap gap-4 text-xs">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" v-model="form.rohsCompliant" class="rounded text-brand-600 focus:ring-brand-500" />
            <span class="font-medium text-slate-700">EU RoHS 3 Compliant</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" v-model="form.reachCompliant" class="rounded text-brand-600 focus:ring-brand-500" />
            <span class="font-medium text-slate-700">REACH SVHC Free</span>
          </label>
        </div>
      </div>

      <!-- SECTION 02: DESCRIPTION -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="file-text" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Section 02 — Description</h2>
        </div>

        <div class="space-y-4">
          <div>
            <label class="form-label">Short Summary</label>
            <input
              v-model="form.shortDescription"
              type="text"
              placeholder="Concise 1-line engineering summary"
              class="form-input"
            />
          </div>
          <div>
            <label class="form-label">Detailed Functional Description</label>
            <textarea
              v-model="form.detailedDescription"
              rows="3"
              placeholder="Provide complete electrical architecture details, clock features, bus configurations..."
              class="form-textarea"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- SECTION 03: TECHNICAL SPECIFICATIONS (DYNAMIC KEY-VALUE ROWS) -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-200 mb-4 gap-2">
          <div class="flex items-center gap-2">
            <AppIcon name="zap" size="sm" class="text-amber-500" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Section 03 — Technical Specifications</h2>
          </div>
          <div class="flex items-center gap-2">
            <!-- Preset Quick-Fill Templates -->
            <button type="button" @click="loadTemplate('mcu')" class="btn btn-secondary text-3xs py-1">
              + MCU Template
            </button>
            <button type="button" @click="loadTemplate('sensor')" class="btn btn-secondary text-3xs py-1">
              + Sensor Template
            </button>
            <button type="button" @click="addSpecRow" class="btn btn-primary text-3xs py-1">
              <AppIcon name="plus" size="xs" class="mr-1" />
              <span>Add Row</span>
            </button>
          </div>
        </div>

        <div v-if="form.specs.length === 0" class="py-6 text-center text-slate-400 text-xs border border-dashed rounded-lg">
          No parametric specifications yet. Click "Add Row" or load a template above.
        </div>

        <div v-else class="space-y-2">
          <div class="grid grid-cols-12 gap-2 text-2xs font-bold text-slate-500 uppercase px-1">
            <div class="col-span-5">Specification Parameter</div>
            <div class="col-span-4">Value</div>
            <div class="col-span-2">Unit</div>
            <div class="col-span-1 text-center">Action</div>
          </div>

          <div
            v-for="(spec, idx) in form.specs"
            :key="idx"
            class="grid grid-cols-12 gap-2 items-center"
          >
            <div class="col-span-5">
              <input
                v-model="spec.name"
                type="text"
                placeholder="e.g. Supply Voltage, Clock Speed"
                class="form-input text-xs"
                required
              />
            </div>
            <div class="col-span-4">
              <input
                v-model="spec.value"
                type="text"
                placeholder="e.g. 3.3, 480"
                class="form-input font-mono text-xs font-bold"
                required
              />
            </div>
            <div class="col-span-2">
              <input
                v-model="spec.unit"
                type="text"
                placeholder="e.g. V, MHz, mA"
                class="form-input font-mono text-xs"
              />
            </div>
            <div class="col-span-1 text-center">
              <button
                type="button"
                @click="removeSpecRow(idx)"
                class="p-1 rounded text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                title="Remove Specification"
              >
                <AppIcon name="trash" size="xs" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 04: DOCUMENTS -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="file-text" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Section 04 — Engineering Documents</h2>
        </div>

        <div class="space-y-4">
          <!-- Hidden File Input for Real Document Upload -->
          <input
            ref="docFileInput"
            type="file"
            accept=".pdf,.step,.stp,.cad,.zip,.docx,.png,.jpg"
            class="hidden"
            @change="handleDocumentUpload"
          />

          <div
            @click="$refs.docFileInput.click()"
            class="border-2 border-dashed border-slate-300 hover:border-brand-500 rounded-lg p-5 text-center cursor-pointer bg-slate-50/50 hover:bg-brand-50/20 transition-all"
          >
            <AppIcon name="upload" size="md" class="mx-auto text-brand-600 mb-1" />
            <div class="text-xs font-bold text-slate-800">Click to Upload Datasheet, 3D STEP Model, or Certificate</div>
            <div class="text-3xs text-slate-500 mt-0.5">PDF, STEP, CAD up to 25MB (Uploaded directly to Backend Storage)</div>
          </div>

          <!-- Document List -->
          <div v-if="form.documents && form.documents.length > 0" class="space-y-2">
            <div
              v-for="(doc, idx) in form.documents"
              :key="doc.id || idx"
              class="flex items-center justify-between p-2.5 rounded-lg border border-slate-200 bg-slate-50/60"
            >
              <div class="flex items-center gap-2.5 min-w-0">
                <AppIcon name="file-text" size="sm" class="text-brand-600 shrink-0" />
                <div class="min-w-0">
                  <div class="text-xs font-bold text-slate-800 truncate">{{ doc.name || doc.fileName }}</div>
                  <div class="text-3xs text-slate-500 font-mono">{{ doc.type || doc.documentType }} • {{ doc.size || doc.fileSize }}</div>
                </div>
              </div>
              <button
                type="button"
                @click="removeDocument(idx)"
                class="text-slate-400 hover:text-rose-600 p-1 rounded"
              >
                <AppIcon name="trash" size="xs" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- STICKY FOOTER ACTION BAR -->
      <div class="fixed bottom-0 left-0 right-0 z-40 bg-white/95 backdrop-blur-md border-t border-slate-200 py-3 px-6 shadow-modal">
        <div class="max-w-4xl mx-auto flex items-center justify-between">
          <button type="button" @click="cancel" class="btn btn-secondary">
            Cancel
          </button>
          <div class="flex items-center gap-2">
            <button type="button" @click="save('Draft')" class="btn btn-secondary">
              Save Draft
            </button>
            <button type="submit" class="btn btn-primary">
              <AppIcon name="check" size="xs" class="mr-1.5" />
              <span>{{ isEdit ? 'Save Changes' : 'Save Component' }}</span>
            </button>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ComponentFormView',
  components: {
    PageHeader,
    AppIcon
  },
  data() {
    return {
      form: {
        id: null,
        partNumber: '',
        internalCode: '',
        name: '',
        category: '',
        categoryId: '',
        manufacturer: '',
        manufacturerId: '',
        supplier: '',
        supplierId: '',
        unit: 'PCS',
        unitId: 'UOM-001',
        status: 'Active',
        packageType: '',
        mountingType: 'Surface Mount (SMD)',
        estimatedUnitCost: 0,
        currency: 'IDR',
        purchaseLink: '',
        rohsCompliant: true,
        reachCompliant: true,
        shortDescription: '',
        detailedDescription: '',
        specs: [
          { name: 'Supply Voltage', value: '3.3', unit: 'V' },
          { name: 'Operating Temp', value: '-40 to +85', unit: '°C' }
        ],
        documents: []
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    componentId() {
      return this.$route.params.id;
    },
    isEdit() {
      return Boolean(this.componentId);
    }
  },
  created() {
    if (this.isEdit) {
      const existing = this.store.components.find(c => c.id === this.componentId || c.partNumber === this.componentId);
      if (existing) {
        this.form = JSON.parse(JSON.stringify(existing));
      } else {
        this.$router.push('/components');
      }
    } else {
      this.form.internalCode = `IPN-ENG-${String(this.store.components.length + 1).padStart(3, '0')}`;
    }
  },
  methods: {
    addSpecRow() {
      this.form.specs.push({ name: '', value: '', unit: '' });
    },
    removeSpecRow(idx) {
      this.form.specs.splice(idx, 1);
    },
    loadTemplate(type) {
      if (type === 'mcu') {
        this.form.specs = [
          { name: 'Core Architecture', value: 'Arm Cortex-M7', unit: '' },
          { name: 'Clock Speed', value: '480', unit: 'MHz' },
          { name: 'Flash Memory', value: '2048', unit: 'KB' },
          { name: 'RAM Size', value: '1024', unit: 'KB' },
          { name: 'Supply Voltage', value: '1.8 - 3.6', unit: 'V' },
          { name: 'Operating Temperature', value: '-40 to +85', unit: '°C' }
        ];
        this.store.showToast('Template Loaded', 'MCU specification template applied.', 'info');
      } else if (type === 'sensor') {
        this.form.specs = [
          { name: 'Supply Voltage (VDD)', value: '1.71 - 3.6', unit: 'V' },
          { name: 'Interface Protocol', value: 'I2C / SPI', unit: '' },
          { name: 'Measurement Range', value: '300 to 1100', unit: 'hPa' },
          { name: 'Active Current Draw', value: '3.7', unit: 'µA' },
          { name: 'Operating Temperature', value: '-40 to +85', unit: '°C' }
        ];
        this.store.showToast('Template Loaded', 'Sensor specification template applied.', 'info');
      }
    },
    async handleDocumentUpload(event) {
      const file = event.target.files[0];
      if (!file) return;
      try {
        const uploadRes = await this.store.uploadFile(file);
        if (uploadRes && uploadRes.data) {
          const docType = file.name.endsWith('.step') || file.name.endsWith('.stp') ? '3D CAD Model' : 'Datasheet';
          this.form.documents.push({
            id: `DOC-${Date.now()}`,
            name: uploadRes.data.fileName,
            type: docType,
            path: uploadRes.data.url,
            size: uploadRes.data.fileSize,
            mimeType: uploadRes.data.mimeType,
            uploadedBy: this.store.currentUser.name
          });
          this.store.showToast('Document Uploaded', `${file.name} saved to server.`, 'success');
        }
      } catch (err) {
        console.error('Document upload error:', err);
      }
    },
    removeDocument(idx) {
      this.form.documents.splice(idx, 1);
    },
    cancel() {
      if (this.isEdit) {
        this.$router.push(`/components/${this.componentId}`);
      } else {
        this.$router.push('/components');
      }
    },
    async save(statusOverride) {
      if (!this.form.partNumber || !this.form.name || !this.form.category || !this.form.manufacturer) {
        this.store.showToast('Validation Error', 'Please complete all required fields (*).', 'error');
        return;
      }
      if (statusOverride) {
        this.form.status = statusOverride;
      }

      // Resolve relational IDs
      const mfgObj = this.store.manufacturers.find(m => m.name === this.form.manufacturer);
      if (mfgObj) this.form.manufacturerId = mfgObj.id;

      const supObj = this.store.suppliers.find(s => s.name === this.form.supplier);
      if (supObj) this.form.supplierId = supObj.id;

      const catObj = this.store.flatCategories.find(c => c.rawName === this.form.category);
      if (catObj) this.form.categoryId = catObj.id;

      const unitObj = this.store.units.find(u => u.code === this.form.unit);
      if (unitObj) this.form.unitId = unitObj.id;

      try {
        const saved = await this.store.saveComponent(this.form);
        this.$router.push(`/components/${saved.id || saved.partNumber}`);
      } catch (err) {
        console.error('Save component error:', err);
      }
    }
  }
};
</script>
