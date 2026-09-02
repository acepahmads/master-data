<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between bg-slate-50 p-4 rounded-xl border border-slate-200">
      <div>
        <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">Source Header to Target Field Mapping</h4>
        <p class="text-2xs text-slate-500 mt-0.5">Map each detected column in your Excel workbooks to canonical Master Data fields.</p>
      </div>
      <button @click="autoDetect" class="btn btn-secondary text-xs">
        <AppIcon name="sparkles" size="xs" class="mr-1.5 text-purple-600" />
        <span>Auto-Map Columns</span>
      </button>
    </div>

    <div class="overflow-x-auto border border-slate-200 rounded-xl">
      <table class="dense-table w-full">
        <thead>
          <tr class="bg-slate-100/80 text-slate-600 uppercase text-3xs font-bold">
            <th class="w-12 text-center">#</th>
            <th class="min-w-[200px]">Detected Excel Header</th>
            <th class="min-w-[220px]">Target Master Data Field</th>
            <th class="min-w-[150px]">Mapping Status</th>
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
                    class="px-1.5 py-0.5 rounded bg-slate-100 border border-slate-200/80 font-mono text-slate-700 max-w-[200px] truncate"
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
