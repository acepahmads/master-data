<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <router-link to="/data-import" class="text-xs text-brand-600 hover:text-brand-800 font-semibold flex items-center gap-1">
            <AppIcon name="chevron-left" size="xs" /> Back to Batches
          </router-link>
          <span class="text-slate-400">•</span>
          <span class="font-mono text-xs text-slate-500 font-bold">{{ batch?.batchCode || 'Loading...' }}</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">Import & Migration Wizard</h1>
      </div>

      <div v-if="batch" class="flex items-center gap-2">
        <router-link :to="`/data-import/review/${batch.id}`" class="btn btn-secondary text-xs">
          <AppIcon name="eye" size="xs" class="mr-1.5" />
          <span>Open Review Workspace</span>
        </router-link>
      </div>
    </div>

    <!-- Wizard Stepper Indicator -->
    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-2xs">
      <div class="flex items-center justify-between overflow-x-auto gap-2">
        <div v-for="(step, idx) in steps" :key="step.key" class="flex items-center gap-2 shrink-0">
          <div
            class="w-7 h-7 rounded-full flex items-center justify-center font-mono font-bold text-xs transition-colors"
            :class="getStepCircleClass(idx)"
          >
            {{ idx + 1 }}
          </div>
          <div class="text-left">
            <div class="text-2xs font-bold" :class="getStepTextClass(idx)">{{ step.title }}</div>
            <div class="text-3xs text-slate-400">{{ step.subtitle }}</div>
          </div>
          <AppIcon v-if="idx < steps.length - 1" name="chevron-right" size="xs" class="text-slate-300 ml-2" />
        </div>
      </div>
    </div>

    <!-- Wizard Step Content Container -->
    <div class="app-card p-6">
      <!-- Step 1: Upload Files -->
      <div v-if="activeStep === 0" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-slate-900">Step 1: Upload Source Excel Workbooks</h3>
          <p class="text-xs text-slate-500 mt-0.5">Select or drag & drop your 2026 pricelists, vendor catalogs, or component sheets.</p>
        </div>

        <div
          class="border-2 border-dashed border-slate-300 hover:border-brand-500 rounded-2xl p-8 text-center transition-colors bg-slate-50/50 cursor-pointer"
          @click="$refs.fileInput.click()"
          @dragover.prevent
          @drop.prevent="handleFileDrop"
        >
          <input
            type="file"
            ref="fileInput"
            multiple
            accept=".xlsx,.xlsm"
            class="hidden"
            @change="handleFileSelect"
          />
          <AppIcon name="file-spreadsheet" size="lg" class="mx-auto text-brand-600 mb-2" />
          <div class="text-xs font-bold text-slate-800">Click to upload or drag & drop files here</div>
          <div class="text-3xs text-slate-400 mt-1">Supports XLSX, XLSM workbooks (up to 50MB per file)</div>
        </div>

        <!-- Staged Files List -->
        <div v-if="batch && batch.files && batch.files.length > 0" class="space-y-2">
          <h4 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Uploaded Workbooks in Batch</h4>
          <div class="divide-y divide-slate-100 border border-slate-200 rounded-xl overflow-hidden bg-white">
            <div v-for="f in batch.files" :key="f.id" class="p-3.5 flex items-center justify-between text-xs">
              <div class="flex items-center gap-3">
                <AppIcon name="file-spreadsheet" size="sm" class="text-emerald-600" />
                <div>
                  <div class="font-bold text-slate-800">{{ f.originalFileName }}</div>
                  <div class="text-3xs text-slate-400 font-mono">
                    {{ (f.fileSizeBytes / 1024).toFixed(1) }} KB • {{ f.sheetCount }} sheet(s) detected • SHA: {{ f.fileSha256 ? f.fileSha256.substring(0, 8) : '-' }}...
                  </div>
                </div>
              </div>
              <span class="px-2 py-0.5 rounded text-3xs font-mono font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                Uploaded
              </span>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end pt-4 border-t border-slate-100">
          <button
            @click="goToStep2"
            class="btn btn-primary text-xs"
            :disabled="!batch || !batch.files || batch.files.length === 0"
          >
            <span>Proceed to Sheet Selection</span>
            <AppIcon name="chevron-right" size="xs" class="ml-1" />
          </button>
        </div>
      </div>

      <!-- Step 2: Sheet Inspection & Parsing -->
      <div v-if="activeStep === 1" class="space-y-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div>
            <h3 class="text-sm font-bold text-slate-900">Step 2: Workbook Sheet Inspection & Selection</h3>
            <p class="text-xs text-slate-500 mt-0.5">Pilih sheet yang ingin di-import dan klik Preview untuk melihat sampel data sebelum ekstraksi.</p>
          </div>
          <div class="flex items-center gap-2">
            <button @click="selectAllSheets(true)" class="btn btn-secondary text-3xs py-1 px-2">Select All</button>
            <button @click="selectAllSheets(false)" class="btn btn-secondary text-3xs py-1 px-2">Deselect All</button>
          </div>
        </div>

        <div v-for="f in batch?.files" :key="f.id" class="border border-slate-200 rounded-xl p-4 bg-slate-50 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <AppIcon name="file-spreadsheet" size="xs" class="text-emerald-600" />
              <strong class="text-xs font-bold text-slate-800">{{ f.originalFileName }}</strong>
            </div>
            <span class="text-3xs text-slate-500 font-mono">{{ f.sheetCount }} Sheets</span>
          </div>

          <div class="grid grid-cols-1 gap-2.5">
            <div
              v-for="s in parseManifest(f.sheetManifestJson)"
              :key="s.name"
              class="p-3.5 bg-white rounded-xl border transition-all flex flex-col space-y-2"
              :class="isSheetSelected(f.id, s.name) ? 'border-brand-300 shadow-2xs bg-brand-50/10' : 'border-slate-200 opacity-70'"
            >
              <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <div class="flex items-center gap-3">
                  <input
                    type="checkbox"
                    :id="`sheet-${f.id}-${s.name}`"
                    :checked="isSheetSelected(f.id, s.name)"
                    @change="toggleSheet(f.id, s.name)"
                    class="w-4 h-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                  />
                  <div>
                    <label :for="`sheet-${f.id}-${s.name}`" class="font-bold text-xs text-slate-900 cursor-pointer flex flex-wrap items-center gap-2">
                      <span>{{ s.name }}</span>
                      <span
                        v-if="s.rowCount > 1 && s.sampleRows && s.sampleRows.length > 0"
                        class="px-2 py-0.2 rounded text-3xs font-mono font-bold bg-purple-50 text-purple-700 border border-purple-200"
                      >
                        ✨ AI Recommended
                      </span>
                      <span
                        class="px-2 py-0.2 rounded text-3xs font-mono font-bold"
                        :class="isSheetSelected(f.id, s.name) ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-500'"
                      >
                        {{ isSheetSelected(f.id, s.name) ? 'Selected for Import' : 'Ignored' }}
                      </span>
                    </label>
                    <div class="text-3xs text-slate-500 font-mono mt-0.5">
                      <strong class="text-slate-700">{{ s.rowCount }}</strong> baris • <strong class="text-slate-700">{{ s.columnCount }}</strong> kolom terdeteksi
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-2 shrink-0">
                  <button
                    type="button"
                    @click="openSheetPreview(f, s)"
                    class="btn btn-secondary text-2xs py-1 px-2.5 flex items-center gap-1 hover:bg-slate-100"
                  >
                    <AppIcon name="eye" size="xs" class="text-slate-500" />
                    <span>Preview Data ({{ s.sampleRows ? s.sampleRows.length : 0 }})</span>
                  </button>
                </div>
              </div>

              <!-- Header Row Selector per Sheet -->
              <div v-if="isSheetSelected(f.id, s.name)" class="flex flex-wrap items-center gap-2.5 pt-2 border-t border-slate-100 text-3xs text-slate-600 bg-slate-50/60 p-2 rounded-lg">
                <div class="flex items-center gap-1 font-semibold text-slate-700">
                  <AppIcon name="table" size="xs" class="text-brand-600" />
                  <span>Header Row:</span>
                </div>
                <select
                  v-model.number="headerRowOverrides[s.name]"
                  class="form-select text-3xs py-1 px-2 font-mono w-56 rounded-md bg-white border-slate-300 focus:border-brand-500"
                >
                  <option :value="s.headerRowIndex || 1">Auto-Detected (Baris {{ s.headerRowIndex || 1 }})</option>
                  <option v-for="r in getHeaderRowOptions(s)" :key="r" :value="r">Gunakan Baris ke-{{ r }} sebagai Header</option>
                </select>
                <span class="text-4xs text-slate-400">Pilih baris judul kolom jika format spreadsheet memiliki judul/kop surat di atas</span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-4 border-t border-slate-100">
          <button @click="activeStep = 0" class="btn btn-secondary text-xs">Back</button>
          <button
            @click="handleExecuteParse"
            class="btn btn-primary text-xs flex items-center gap-1.5"
            :disabled="parsingLoading || totalDetectedRows === 0"
          >
            <AppIcon v-if="!parsingLoading" name="sparkles" size="xs" class="text-amber-300" />
            <span v-if="parsingLoading">Extracting Staged Rows...</span>
            <span v-else>Parse & Extract {{ totalDetectedRows }} Staged Rows</span>
          </button>
        </div>
      </div>

      <!-- Step 3: Column Mapping -->
      <div v-if="activeStep === 2" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-slate-900">Step 3: Column Mapping & Canonical Normalization</h3>
          <p class="text-xs text-slate-500 mt-0.5">Confirm mapping between your Excel headers and the Master Data fields. Live sample values are shown below each column.</p>
        </div>

        <ColumnMappingGrid
          :headers="detectedHeadersList"
          :initialMapping="currentMapping"
          :columnSamples="columnSamplesMap"
          @change="updateMapping"
          @auto-detect="autoDetectMapping"
        />

        <div class="flex items-center justify-between pt-4 border-t border-slate-100">
          <button @click="activeStep = 1" class="btn btn-secondary text-xs">Back</button>
          <button @click="handleApplyMapping" class="btn btn-primary text-xs" :disabled="mappingLoading">
            <span v-if="mappingLoading">Applying Mapping...</span>
            <span v-else>Apply Mapping & Validate Rules</span>
          </button>
        </div>
      </div>

      <!-- Step 4: Rule Validation & AI Co-Pilot -->
      <div v-if="activeStep === 3" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-slate-900">Step 4: Rule Validation & AI Co-Pilot Classification</h3>
          <p class="text-xs text-slate-500 mt-0.5">Review validation results and run AI suggestions on ambiguous rows.</p>
        </div>

        <!-- Validation Summary Cards -->
        <div class="grid grid-cols-3 gap-3">
          <div class="p-4 rounded-xl bg-emerald-50 border border-emerald-200">
            <div class="text-3xs font-bold text-emerald-700 uppercase">Valid Rows</div>
            <div class="text-2xl font-black text-emerald-900 font-mono mt-1">{{ batch?.validRows || 0 }}</div>
            <div class="text-3xs text-emerald-600 mt-0.5">Ready for review & import</div>
          </div>

          <div class="p-4 rounded-xl bg-amber-50 border border-amber-200">
            <div class="text-3xs font-bold text-amber-700 uppercase">Warnings</div>
            <div class="text-2xl font-black text-amber-900 font-mono mt-1">{{ batch?.warningRows || 0 }}</div>
            <div class="text-3xs text-amber-600 mt-0.5">Missing optional details</div>
          </div>

          <div class="p-4 rounded-xl bg-rose-50 border border-rose-200">
            <div class="text-3xs font-bold text-rose-700 uppercase">Validation Errors</div>
            <div class="text-2xl font-black text-rose-900 font-mono mt-1">{{ batch?.errorRows || 0 }}</div>
            <div class="text-3xs text-rose-600 mt-0.5">Requires correction before import</div>
          </div>
        </div>

        <!-- AI Co-Pilot Section -->
        <div class="p-5 rounded-2xl bg-purple-50/60 border border-purple-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-1.5 font-bold text-purple-900 text-xs">
              <AppIcon name="sparkles" size="xs" class="text-purple-600" />
              <span>AI Classification & Cleansing Co-Pilot</span>
            </div>
            <p class="text-2xs text-purple-700 mt-1 max-w-lg">
              Automatically suggest Item Classification (Product vs Component vs Service) and Category with confidence scoring.
            </p>
          </div>
          <button @click="handleRunAI" class="btn bg-purple-600 hover:bg-purple-700 text-white text-xs shrink-0" :disabled="aiLoading">
            <AppIcon name="sparkles" size="xs" class="mr-1.5" />
            <span v-if="aiLoading">Analyzing with AI...</span>
            <span v-else>Run AI Suggestions</span>
          </button>
        </div>

        <div class="flex items-center justify-between pt-4 border-t border-slate-100">
          <button @click="activeStep = 2" class="btn btn-secondary text-xs">Back</button>
          <router-link :to="`/data-import/review/${batch.id}`" class="btn btn-primary text-xs">
            <span>Proceed to Triage Workspace</span>
            <AppIcon name="chevron-right" size="xs" class="ml-1" />
          </router-link>
        </div>
      </div>
    </div>

    <!-- Sheet Data Preview Modal -->
    <div v-if="previewModal.open" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-4xl max-h-[85vh] flex flex-col overflow-hidden">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
          <div>
            <div class="flex items-center gap-2">
              <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                Sheet Preview
              </span>
              <span class="text-xs font-semibold text-slate-500 font-mono">{{ previewModal.fileName }}</span>
            </div>
            <h3 class="text-sm font-bold text-slate-900 mt-0.5">
              Sheet: <span class="text-brand-600">{{ previewModal.sheetName }}</span> (Menampilkan Seluruh {{ previewModal.sampleRows ? previewModal.sampleRows.length : 0 }} Baris)
            </h3>
          </div>
          <button @click="previewModal.open = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200 transition-colors">
            <AppIcon name="x" size="sm" />
          </button>
        </div>

        <!-- Modal Table Body -->
        <div class="flex-1 overflow-auto p-4">
          <div v-if="!previewModal.headers || previewModal.headers.length === 0" class="p-8 text-center text-xs text-slate-400">
            Tidak ada kolom header terdeteksi pada sheet ini.
          </div>
          <table v-else class="dense-table w-full border border-slate-200 rounded-lg text-xs">
            <thead>
              <tr class="bg-slate-100 text-slate-700 text-3xs font-bold uppercase sticky top-0 shadow-xs">
                <th class="w-10 text-center">#</th>
                <th v-for="(h, idx) in previewModal.headers" :key="idx" class="whitespace-nowrap px-3 py-2 border-r border-slate-200">
                  {{ h }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 bg-white">
              <tr v-for="(row, rIdx) in previewModal.sampleRows" :key="rIdx" class="hover:bg-slate-50/80">
                <td class="text-center font-mono text-3xs text-slate-400 bg-slate-50/50 border-r border-slate-100">
                  {{ rIdx + 1 }}
                </td>
                <td v-for="(cell, cIdx) in row" :key="cIdx" class="px-3 py-2 border-r border-slate-100 truncate max-w-[220px]" :title="cell">
                  {{ cell || '-' }}
                </td>
              </tr>
              <tr v-if="previewModal.sampleRows.length === 0">
                <td :colspan="previewModal.headers.length + 1" class="p-6 text-center text-slate-400 text-xs italic">
                  Tidak ada baris data sampel.
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-3 border-t border-slate-200 bg-slate-50 flex items-center justify-between text-2xs text-slate-500">
          <span>Menampilkan seluruh baris data dari sheet untuk memverifikasi isi tabel sebelum diekstrak ke staging.</span>
          <button @click="previewModal.open = false" class="btn btn-primary text-xs py-1 px-4">
            Tutup Preview
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useImportStore } from '../../stores/importStore';
import AppIcon from '../../components/common/AppIcon.vue';
import ColumnMappingGrid from '../../components/import/ColumnMappingGrid.vue';

export default {
  name: 'ImportWizardView',
  components: { AppIcon, ColumnMappingGrid },
  data() {
    return {
      store: useImportStore(),
      activeStep: 0,
      parsingLoading: false,
      mappingLoading: false,
      aiLoading: false,
      currentMapping: {},
      selectedSheetsMap: {}, // { fileId: { sheetName: true/false } }
      headerRowOverrides: {}, // { sheetName: 1-based row index }
      previewModal: {
        open: false,
        fileName: '',
        sheetName: '',
        headers: [],
        sampleRows: []
      },
      steps: [
        { key: 'upload', title: '1. Upload Files', subtitle: 'Multiple XLSX/XLSM' },
        { key: 'sheets', title: '2. Sheet Manifest', subtitle: 'Workbook scan' },
        { key: 'mapping', title: '3. Column Mapping', subtitle: 'Canonical fields' },
        { key: 'validate', title: '4. Validate & AI', subtitle: 'Rules & suggestions' }
      ]
    };
  },
  computed: {
    batch() {
      return this.store.currentBatch;
    },
    totalDetectedRows() {
      let sum = 0;
      for (const f of this.batch?.files || []) {
        const manifest = this.parseManifest(f.sheetManifestJson);
        for (const s of manifest) {
          if (this.isSheetSelected(f.id, s.name)) {
            sum += s.rowCount || 0;
          }
        }
      }
      return sum;
    },
    detectedHeadersList() {
      const set = new Set();
      for (const f of this.batch?.files || []) {
        const manifest = this.parseManifest(f.sheetManifestJson);
        for (const m of manifest) {
          if (this.isSheetSelected(f.id, m.name)) {
            for (const h of m.headers || []) {
              if (h && h.trim()) {
                set.add(h.trim());
              }
            }
          }
        }
      }
      return Array.from(set);
    },
    columnSamplesMap() {
      const map = {};
      for (const f of this.batch?.files || []) {
        const manifest = this.parseManifest(f.sheetManifestJson);
        for (const m of manifest) {
          if (this.isSheetSelected(f.id, m.name) && m.columnSamples) {
            for (const [col, samples] of Object.entries(m.columnSamples)) {
              if (!map[col]) map[col] = [];
              for (const s of samples) {
                if (s && !map[col].includes(s)) {
                  map[col].push(s);
                }
              }
            }
          }
        }
      }
      return map;
    }
  },
  async mounted() {
    const id = this.$route.params.id;
    if (id) {
      await this.store.fetchBatchById(id);
      this.initSelectedSheets();
      if (this.batch && this.batch.columnMappingJson) {
        try {
          this.currentMapping = JSON.parse(this.batch.columnMappingJson);
        } catch (e) {}
      }
      if (this.batch && (this.batch.status === 'PARSED' || this.batch.status === 'VALIDATING' || this.batch.status === 'READY_FOR_REVIEW')) {
        this.activeStep = 2;
      }
    }
  },
  methods: {
    getHeaderRowOptions(sheet) {
      if (!sheet) return [1];
      const maxRows = Math.min(sheet.rowCount || 50, 50);
      const limit = Math.max(sheet.headerRowIndex || 1, maxRows);
      const arr = [];
      for (let i = 1; i <= limit; i++) {
        arr.push(i);
      }
      return arr;
    },
    parseManifest(jsonStr) {
      try {
        return JSON.parse(jsonStr || '[]');
      } catch (e) {
        return [];
      }
    },
    initSelectedSheets() {
      const map = {};
      const headerMap = {};
      for (const f of this.batch?.files || []) {
        map[f.id] = {};
        const manifest = this.parseManifest(f.sheetManifestJson);
        for (const s of manifest) {
          // AI Smart Recommendation: select sheets with real data rows, leave blank sheets unchecked
          const isCandidate = (s.rowCount > 1 && s.sampleRows && s.sampleRows.length > 0);
          map[f.id][s.name] = isCandidate;
          if (this.headerRowOverrides[s.name] === undefined) {
            headerMap[s.name] = s.headerRowIndex || 1;
          }
        }
      }
      this.selectedSheetsMap = map;
      this.headerRowOverrides = { ...this.headerRowOverrides, ...headerMap };
    },
    isSheetSelected(fileId, sheetName) {
      if (!this.selectedSheetsMap[fileId]) return true;
      return this.selectedSheetsMap[fileId][sheetName] !== false;
    },
    toggleSheet(fileId, sheetName) {
      if (!this.selectedSheetsMap[fileId]) {
        this.$set(this.selectedSheetsMap, fileId, {});
      }
      const current = this.isSheetSelected(fileId, sheetName);
      this.$set(this.selectedSheetsMap[fileId], sheetName, !current);
    },
    selectAllSheets(val) {
      for (const f of this.batch?.files || []) {
        if (!this.selectedSheetsMap[f.id]) {
          this.$set(this.selectedSheetsMap, f.id, {});
        }
        const manifest = this.parseManifest(f.sheetManifestJson);
        for (const s of manifest) {
          this.$set(this.selectedSheetsMap[f.id], s.name, val);
        }
      }
    },
    openSheetPreview(file, sheet) {
      this.previewModal = {
        open: true,
        fileName: file.originalFileName,
        sheetName: sheet.name,
        headers: sheet.headers || [],
        sampleRows: sheet.sampleRows || []
      };
    },
    goToStep2() {
      this.initSelectedSheets();
      this.activeStep = 1;
    },
    getStepCircleClass(idx) {
      if (idx === this.activeStep) return 'bg-brand-600 text-white ring-4 ring-brand-100';
      if (idx < this.activeStep) return 'bg-emerald-500 text-white';
      return 'bg-slate-100 text-slate-400';
    },
    getStepTextClass(idx) {
      if (idx === this.activeStep) return 'text-brand-700';
      if (idx < this.activeStep) return 'text-slate-800';
      return 'text-slate-400';
    },
    async handleFileSelect(e) {
      const files = e.target.files;
      if (!files || files.length === 0) return;
      await this.uploadFiles(files);
    },
    async handleFileDrop(e) {
      const files = e.dataTransfer.files;
      if (!files || files.length === 0) return;
      await this.uploadFiles(files);
    },
    async uploadFiles(files) {
      const formData = new FormData();
      for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i]);
      }
      try {
        await this.store.uploadFiles(this.batch.id, formData);
        this.initSelectedSheets();
      } catch (err) {
        alert('Upload failed: ' + err.message);
      }
    },
    async handleExecuteParse() {
      this.parsingLoading = true;
      try {
        // Collect selected sheet names per file ID
        const selectedByFile = {};
        for (const f of this.batch?.files || []) {
          selectedByFile[f.id] = [];
          const manifest = this.parseManifest(f.sheetManifestJson);
          for (const s of manifest) {
            if (this.isSheetSelected(f.id, s.name)) {
              selectedByFile[f.id].push(s.name);
            }
          }
        }

        await this.store.parseBatch(this.batch.id, selectedByFile, this.headerRowOverrides);
        await this.store.fetchBatchById(this.batch.id);
        if (this.batch && this.batch.columnMappingJson) {
          try {
            this.currentMapping = JSON.parse(this.batch.columnMappingJson);
          } catch (e) {}
        }
        this.activeStep = 2;
      } catch (err) {
        alert('Parsing failed: ' + err.message);
      } finally {
        this.parsingLoading = false;
      }
    },
    updateMapping(newMapping) {
      this.currentMapping = newMapping;
    },
    autoDetectMapping() {
      this.handleExecuteParse();
    },
    async handleApplyMapping() {
      this.mappingLoading = true;
      try {
        await this.store.applyMapping(this.batch.id, this.currentMapping);
        await this.store.fetchBatchById(this.batch.id);
        this.activeStep = 3;
      } catch (err) {
        alert('Mapping application failed: ' + err.message);
      } finally {
        this.mappingLoading = false;
      }
    },
    async handleRunAI() {
      this.aiLoading = true;
      try {
        await this.store.runAIClassification(this.batch.id);
        await this.store.fetchBatchById(this.batch.id);
        alert('AI classification suggestions successfully generated for staged rows!');
      } catch (err) {
        alert('AI classification notice: ' + err.message);
      } finally {
        this.aiLoading = false;
      }
    }
  }
};
</script>
