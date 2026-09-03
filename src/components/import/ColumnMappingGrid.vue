<template>
  <div class="space-y-6">
    <!-- Action Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between bg-slate-50 p-4 rounded-xl border border-slate-200 gap-3">
      <div>
        <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">Source Header to Target Field Mapping</h4>
        <p class="text-2xs text-slate-500 mt-0.5">
          <span v-if="fileGroups && fileGroups.length > 1">
            Daftar kolom telah dikelompokkan per file workbook. Seluruh kolom utama telah otomatis dipetakan ke field Master Data.
          </span>
          <span v-else>
            Map each detected column in your Excel workbooks to canonical Master Data fields.
          </span>
        </p>
      </div>
      <button @click="autoDetect" class="btn btn-secondary text-xs shrink-0 flex items-center gap-1.5">
        <AppIcon name="sparkles" size="xs" class="text-purple-600" />
        <span>Auto-Map Columns</span>
      </button>
    </div>

    <!-- Multi-File View: Grouped Cards per File -->
    <div v-if="fileGroups && fileGroups.length > 1" class="space-y-6">
      <div
        v-for="(group, gIdx) in fileGroups"
        :key="group.fileId || gIdx"
        class="border border-slate-200 rounded-xl overflow-hidden bg-white shadow-2xs"
      >
        <!-- File Card Header -->
        <div class="px-4 py-3 bg-slate-50/90 border-b border-slate-200 flex flex-wrap items-center justify-between gap-2">
          <div class="flex items-center gap-2.5">
            <span class="w-6 h-6 rounded-full bg-brand-50 text-brand-700 font-bold text-3xs flex items-center justify-center border border-brand-200 font-mono">
              {{ gIdx + 1 }}
            </span>
            <AppIcon name="file-spreadsheet" size="xs" class="text-emerald-600" />
            <strong class="text-xs font-bold text-slate-800">{{ group.fileName }}</strong>
            <span v-if="group.sheets" class="text-3xs font-mono font-semibold text-slate-500 bg-white px-2 py-0.5 rounded border border-slate-200">
              Sheet: {{ group.sheets }}
            </span>
          </div>
          <span class="text-3xs text-slate-500 font-mono">{{ group.headers.length }} Kolom Terdeteksi</span>
        </div>

        <!-- Table for this File -->
        <div class="overflow-x-auto">
          <table class="dense-table w-full">
            <thead>
              <tr class="bg-slate-100/60 text-slate-600 uppercase text-3xs font-bold">
                <th class="w-12 text-center">#</th>
                <th class="min-w-[200px]">Detected Excel Header</th>
                <th class="min-w-[220px]">Target Master Data Field</th>
                <th class="min-w-[160px]">Mapping Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 bg-white">
              <tr v-for="(header, idx) in group.headers" :key="header" class="hover:bg-slate-50/70 transition-colors">
                <td class="text-center font-mono text-3xs text-slate-400">{{ idx + 1 }}</td>
                <td class="py-2.5">
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <span class="font-bold text-slate-900 text-xs font-mono">{{ header }}</span>
                    </div>
                    <!-- Sample data preview chips -->
                    <div
                      v-if="getSamples(group, header).length"
                      class="text-3xs text-slate-500 flex items-center gap-1.5 flex-wrap"
                    >
                      <span class="text-slate-400 font-semibold uppercase text-4xs tracking-wider">Samples:</span>
                      <span
                        v-for="(sample, sIdx) in getSamples(group, header).slice(0, 3)"
                        :key="sIdx"
                        class="px-1.5 py-0.5 rounded bg-slate-100 border border-slate-200/80 font-mono text-slate-700 max-w-[220px] truncate"
                        :title="sample"
                      >
                        "{{ sample }}"
                      </span>
                    </div>
                  </div>
                </td>
                <td>
                  <select
                    v-model="localMapping[header]"
                    @change="emitMapping"
                    class="form-select text-xs font-medium py-1.5"
                  >
                    <option value="IGNORE" class="text-slate-400 font-normal">-- Ignore / Skip Column --</option>
                    <optgroup label="Core Master Identifiers">
                      <option value="name">Product / Item Name (Required)</option>
                      <option value="code">Product Code / MPN / SKU</option>
                      <option value="category">Category Classification</option>
                      <option value="description">Technical Description</option>
                    </optgroup>
                    <optgroup label="Sourcing & Commercial">
                      <option value="manufacturer">Manufacturer / Brand</option>
                      <option value="supplier">Supplier / Distributor</option>
                      <option value="unit_cost">Purchase Cost (Buying Price)</option>
                      <option value="selling_price">Selling Price (MSRP)</option>
                      <option value="unit">Unit of Measure (UoM)</option>
                      <option value="target_market">Target Application / Market</option>
                    </optgroup>
                  </select>
                </td>
                <td>
                  <span
                    v-if="localMapping[header] && localMapping[header] !== 'IGNORE'"
                    class="px-2 py-0.5 rounded text-3xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 inline-flex items-center gap-1"
                  >
                    <span class="text-emerald-500 font-bold">✓</span>
                    <span>Mapped ➔ {{ formatField(localMapping[header]) }}</span>
                  </span>
                  <span v-else class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-500">
                    Ignored
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Single File Fallback (Standard Table) -->
    <div v-else class="overflow-x-auto border border-slate-200 rounded-xl">
      <table class="dense-table w-full">
        <thead>
          <tr class="bg-slate-100/80 text-slate-600 uppercase text-3xs font-bold">
            <th class="w-12 text-center">#</th>
            <th class="min-w-[200px]">Detected Excel Header</th>
            <th class="min-w-[220px]">Target Master Data Field</th>
            <th class="min-w-[160px]">Mapping Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="(header, idx) in headers" :key="header" class="hover:bg-slate-50/70 transition-colors">
            <td class="text-center font-mono text-3xs text-slate-400">{{ idx + 1 }}</td>
            <td class="py-2.5">
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-slate-900 text-xs font-mono">{{ header }}</span>
                </div>
                <!-- Sample data preview chips -->
                <div v-if="columnSamples && columnSamples[header] && columnSamples[header].length" class="text-3xs text-slate-500 flex items-center gap-1.5 flex-wrap">
                  <span class="text-slate-400 font-semibold uppercase text-4xs tracking-wider">Samples:</span>
                  <span
                    v-for="(sample, sIdx) in columnSamples[header].slice(0, 3)"
                    :key="sIdx"
                    class="px-1.5 py-0.5 rounded bg-slate-100 border border-slate-200/80 font-mono text-slate-700 max-w-[220px] truncate"
                    :title="sample"
                  >
                    "{{ sample }}"
                  </span>
                </div>
              </div>
            </td>
            <td>
              <select
                v-model="localMapping[header]"
                @change="emitMapping"
                class="form-select text-xs font-medium py-1.5"
              >
                <option value="IGNORE" class="text-slate-400 font-normal">-- Ignore / Skip Column --</option>
                <optgroup label="Core Master Identifiers">
                  <option value="name">Product / Item Name (Required)</option>
                  <option value="code">Product Code / MPN / SKU</option>
                  <option value="category">Category Classification</option>
                  <option value="description">Technical Description</option>
                </optgroup>
                <optgroup label="Sourcing & Commercial">
                  <option value="manufacturer">Manufacturer / Brand</option>
                  <option value="supplier">Supplier / Distributor</option>
                  <option value="unit_cost">Purchase Cost (Buying Price)</option>
                  <option value="selling_price">Selling Price (MSRP)</option>
                  <option value="unit">Unit of Measure (UoM)</option>
                  <option value="target_market">Target Application / Market</option>
                </optgroup>
              </select>
            </td>
            <td>
              <span
                v-if="localMapping[header] && localMapping[header] !== 'IGNORE'"
                class="px-2 py-0.5 rounded text-3xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 inline-flex items-center gap-1"
              >
                <span class="text-emerald-500 font-bold">✓</span>
                <span>Mapped ➔ {{ formatField(localMapping[header]) }}</span>
              </span>
              <span v-else class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-500">
                Ignored
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import AppIcon from '../common/AppIcon.vue';

export default {
  name: 'ColumnMappingGrid',
  components: { AppIcon },
  props: {
    headers: {
      type: Array,
      default: () => []
    },
    fileGroups: {
      type: Array,
      default: () => []
    },
    initialMapping: {
      type: Object,
      default: () => ({})
    },
    columnSamples: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      localMapping: {}
    };
  },
  watch: {
    initialMapping: {
      immediate: true,
      handler(val) {
        this.localMapping = { ...val };
      }
    }
  },
  methods: {
    getSamples(group, header) {
      if (group && group.samples && group.samples[header] && group.samples[header].length) {
        return group.samples[header];
      }
      if (this.columnSamples && this.columnSamples[header]) {
        return this.columnSamples[header];
      }
      return [];
    },
    formatField(field) {
      const labels = {
        name: 'Item Name',
        code: 'Product Code / MPN',
        category: 'Category',
        manufacturer: 'Manufacturer',
        supplier: 'Supplier',
        unit_cost: 'Purchase Cost',
        selling_price: 'Selling Price',
        unit: 'UoM',
        target_market: 'Target Market',
        description: 'Description'
      };
      return labels[field] || field;
    },
    emitMapping() {
      this.$emit('change', { ...this.localMapping });
    },
    autoDetect() {
      this.$emit('auto-detect');
    }
  }
};
</script>
