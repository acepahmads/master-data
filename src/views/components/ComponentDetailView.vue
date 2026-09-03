<template>
  <div v-if="component" class="space-y-5">
    <!-- COMPONENT HEADER -->
    <div class="bg-white rounded-lg border border-slate-200 p-5 shadow-card">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div class="flex items-start gap-4">
          <img
            :src="component.image"
            alt="Component"
            class="w-16 h-16 rounded-lg object-cover border border-slate-200 shadow-xs bg-slate-50 shrink-0"
          />
          <div>
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-slate-900 text-white shadow-xs">
                {{ component.partNumber }}
              </span>
              <span class="font-mono text-3xs text-slate-500 bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">
                {{ component.internalCode }}
              </span>
              <span class="px-2 py-0.5 rounded text-2xs font-semibold bg-slate-100 text-slate-700 border border-slate-200">
                {{ component.category }}
              </span>
              <StatusBadge :status="component.status" />
            </div>
            <h1 class="text-lg font-bold text-slate-900 tracking-tight">{{ component.name }}</h1>
            <div class="flex items-center gap-3 text-2xs text-slate-500 mt-1 font-medium">
              <span>Mfg: <strong class="text-slate-700">{{ component.manufacturer }}</strong></span>
              <span>•</span>
              <span>Supplier: <strong class="text-slate-700">{{ component.supplier }}</strong></span>
              <span>•</span>
              <span>Package: <strong class="font-mono text-slate-700">{{ component.packageType || 'N/A' }}</strong></span>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <router-link :to="`/components/${component.id}/edit`" class="btn btn-primary">
            <AppIcon name="edit" size="xs" class="mr-1.5" />
            <span>Edit Component</span>
          </router-link>
          <button @click="downloadAllDocs" class="btn btn-secondary" title="Export Datasheet Package">
            <AppIcon name="download" size="xs" class="mr-1.5" />
            <span>Datasheets</span>
          </button>
          <button @click="confirmDelete" class="btn btn-danger" title="Delete Component">
            <AppIcon name="trash" size="xs" />
          </button>
        </div>
      </div>

      <!-- TAB NAVIGATION -->
      <div class="flex items-center gap-6 mt-6 border-b border-slate-200 text-xs font-semibold">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          :class="[
            'pb-2.5 flex items-center gap-2 border-b-2 transition-colors',
            activeTab === tab.key
              ? 'border-brand-600 text-brand-600'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon :name="tab.icon" size="xs" />
          <span>{{ tab.label }}</span>
          <span
            v-if="tab.badge !== undefined"
            :class="[
              'text-3xs font-mono px-1.5 py-0.2 rounded',
              activeTab === tab.key ? 'bg-brand-100 text-brand-700' : 'bg-slate-100 text-slate-500'
            ]"
          >
            {{ tab.badge }}
          </span>
        </button>
      </div>
    </div>

    <!-- TAB 1: OVERVIEW -->
    <div v-show="activeTab === 'overview'" class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <div class="lg:col-span-2 space-y-5">
        <!-- Basic Information -->
        <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
            <AppIcon name="component" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Basic Engineering Information</h2>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Manufacturer Part Number (MPN)</div>
              <div class="font-mono font-bold text-slate-900 text-sm mt-0.5">{{ component.partNumber }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Internal Part Number (IPN)</div>
              <div class="font-mono font-bold text-slate-700 mt-0.5">{{ component.internalCode }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Component Name</div>
              <div class="font-semibold text-slate-800 mt-0.5">{{ component.name }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Category</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ component.category }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Manufacturer</div>
              <div class="text-brand-700 font-bold mt-0.5">{{ component.manufacturer }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Primary Supplier</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ component.supplier }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Standard Unit of Measure</div>
              <div class="font-mono font-bold text-slate-800 mt-0.5">{{ component.unit }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Estimated Lead Time</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ component.leadTime || 'Standard 4-8 Weeks' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Package / Case</div>
              <div class="font-mono text-slate-800 mt-0.5">{{ component.packageType || 'Standard' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Mounting Type</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ component.mountingType || 'Surface Mount' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Purchase Cost (Harga Beli)</div>
              <div class="font-mono font-bold text-slate-900 mt-0.5 text-sm">
                {{ formatCurrency(component.estimatedUnitCost, component.currency) }}
              </div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Purchase Link (Link Pembelian)</div>
              <div class="mt-0.5">
                <a
                  v-if="component.purchaseLink"
                  :href="component.purchaseLink"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center gap-1.5 text-brand-600 hover:text-brand-800 font-bold hover:underline"
                >
                  <AppIcon name="external-link" size="2xs" />
                  <span class="truncate max-w-[200px]">Open Purchase Link</span>
                </a>
                <span v-else class="text-slate-400 text-3xs italic">No Link Provided</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-3">
            <AppIcon name="file-text" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Engineering Description</h2>
          </div>
          <p class="text-xs text-slate-700 leading-relaxed whitespace-pre-line">{{ component.detailedDescription || component.shortDescription }}</p>
        </div>
      </div>

      <!-- Right 1 Col: Compliance & Quick Specs -->
      <div class="space-y-5">
        <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-3">
            <AppIcon name="shield" size="sm" class="text-emerald-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Environmental Compliance</h2>
          </div>

          <div class="space-y-2.5 text-xs">
            <div class="flex items-center justify-between p-2 rounded bg-emerald-50/60 border border-emerald-200">
              <span class="font-semibold text-emerald-900">EU RoHS 3 (2015/863/EU)</span>
              <span class="text-3xs font-mono font-bold text-emerald-700 bg-white px-2 py-0.5 rounded border border-emerald-300">
                COMPLIANT
              </span>
            </div>
            <div class="flex items-center justify-between p-2 rounded bg-emerald-50/60 border border-emerald-200">
              <span class="font-semibold text-emerald-900">REACH SVHC</span>
              <span class="text-3xs font-mono font-bold text-emerald-700 bg-white px-2 py-0.5 rounded border border-emerald-300">
                AFFECTED 0
              </span>
            </div>
            <div class="flex items-center justify-between p-2 rounded bg-slate-50 border border-slate-200">
              <span class="font-semibold text-slate-800">Halogen Free</span>
              <span class="text-3xs font-mono font-bold text-slate-600 bg-white px-2 py-0.5 rounded border border-slate-300">
                IEC 61249-2
              </span>
            </div>
          </div>
        </div>

        <div class="bg-slate-50 rounded-lg p-4 border border-slate-200 text-xs space-y-2">
          <div class="font-bold text-slate-900 flex items-center gap-1.5">
            <AppIcon name="zap" size="xs" class="text-amber-500" />
            <span>Key Parametric Highlights</span>
          </div>
          <div class="space-y-1.5 font-mono text-2xs">
            <div v-for="(s, idx) in (component.specs || []).slice(0, 4)" :key="idx" class="flex justify-between border-b border-slate-200 pb-1">
              <span class="text-slate-500">{{ s.name }}:</span>
              <span class="font-bold text-slate-800">{{ s.value }} {{ s.unit }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: TECHNICAL SPECIFICATIONS (STRUCTURED ENGINEERING TABLE) -->
    <div v-show="activeTab === 'specs'" class="space-y-4">
      <div class="bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
        <div class="px-5 py-3.5 border-b border-slate-200 bg-slate-50/70 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <AppIcon name="zap" size="sm" class="text-amber-500" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Parametric Technical Specifications</h2>
          </div>
          <router-link :to="`/components/${component.id}/edit`" class="btn btn-secondary text-2xs">
            <AppIcon name="plus" size="xs" class="mr-1" />
            <span>Add / Edit Parameter</span>
          </router-link>
        </div>

        <div v-if="!component.specs || component.specs.length === 0" class="p-6">
          <EmptyState
            title="No specifications registered"
            description="Technical engineering parameters have not yet been assigned to this component."
            icon="zap"
          />
        </div>

        <table v-else class="dense-table">
          <thead>
            <tr>
              <th class="w-1/3">Specification Parameter</th>
              <th class="w-1/3">Nominal / Rated Value</th>
              <th class="w-1/4">Engineering Unit</th>
              <th class="text-right">Type</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="(spec, idx) in component.specs" :key="idx" class="hover:bg-slate-50/80">
              <td class="font-semibold text-slate-800">{{ spec.name }}</td>
              <td class="font-mono font-bold text-brand-700">{{ spec.value }}</td>
              <td class="font-mono text-slate-600">{{ spec.unit || '—' }}</td>
              <td class="text-right">
                <span class="px-2 py-0.5 rounded text-3xs font-mono bg-slate-100 text-slate-600 border border-slate-200">
                  Parametric
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 3: DOCUMENTS & DATASHEETS -->
    <div v-show="activeTab === 'documents'" class="space-y-4">
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center justify-between pb-3 border-b border-slate-200 mb-4">
          <div class="flex items-center gap-2">
            <AppIcon name="file-text" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Attached Engineering Documents</h2>
          </div>
          <button @click="triggerUploadDoc" class="btn btn-primary text-2xs">
            <AppIcon name="upload" size="xs" class="mr-1" />
            <span>Upload Document</span>
          </button>
        </div>

        <div v-if="!component.documents || component.documents.length === 0" class="p-6">
          <EmptyState
            title="No documents uploaded"
            description="Attach manufacturer datasheets, 3D STEP models, or compliance test reports."
            icon="file-text"
          />
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5">
          <div
            v-for="doc in component.documents"
            :key="doc.id"
            class="p-3.5 rounded-lg border border-slate-200 hover:border-brand-300 hover:shadow-subtle transition-all bg-white flex flex-col justify-between group"
          >
            <div>
              <div class="flex items-start justify-between gap-2">
                <div class="w-8 h-8 rounded bg-brand-50 border border-brand-200 flex items-center justify-center text-brand-600 shrink-0">
                  <AppIcon name="file-text" size="sm" />
                </div>
                <span class="text-3xs font-mono font-semibold px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200">
                  {{ doc.type }}
                </span>
              </div>
              <div class="font-bold text-xs text-slate-800 mt-2 truncate group-hover:text-brand-600" :title="doc.name">
                {{ doc.name }}
              </div>
              <div class="text-3xs text-slate-400 font-mono mt-1">Size: {{ doc.size }} • Uploaded: {{ doc.date }}</div>
            </div>

            <div class="flex items-center gap-2 mt-4 pt-2 border-t border-slate-100">
              <button @click="previewDoc(doc)" class="btn btn-secondary text-2xs flex-1">
                <AppIcon name="eye" size="xs" class="mr-1" />
                <span>Preview</span>
              </button>
              <button @click="downloadDoc(doc)" class="btn btn-secondary text-2xs flex-1">
                <AppIcon name="download" size="xs" class="mr-1" />
                <span>Download</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 4: AUDIT TRAIL / ACTIVITY -->
    <div v-show="activeTab === 'activity'" class="space-y-4">
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center justify-between pb-3 border-b border-slate-200 mb-4">
          <div class="flex items-center gap-2">
            <AppIcon name="activity" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Component Audit Trail</h2>
          </div>
          <span class="text-3xs font-mono text-slate-400">Strict Revision History</span>
        </div>

        <div class="relative pl-4 space-y-4 before:absolute before:left-1.5 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
          <div
            v-for="act in relatedActivities"
            :key="act.id"
            class="relative"
          >
            <span class="absolute -left-4 top-1 w-2.5 h-2.5 rounded-full bg-brand-500 ring-4 ring-white"></span>
            <div class="text-xs">
              <div class="font-semibold text-slate-900">{{ act.action }}</div>
              <div class="text-slate-500 text-3xs mt-0.5">By {{ act.user }} • {{ act.timestamp }}</div>
              <div class="text-3xs font-mono text-slate-400 mt-0.5">{{ act.date }}</div>
            </div>
          </div>

          <div class="relative">
            <span class="absolute -left-4 top-1 w-2.5 h-2.5 rounded-full bg-emerald-500 ring-4 ring-white"></span>
            <div class="text-xs">
              <div class="font-semibold text-slate-900">Component Master Registered</div>
              <div class="text-slate-500 text-3xs mt-0.5">By Engineering Sourcing Team</div>
              <div class="text-3xs font-mono text-slate-400 mt-0.5">{{ component.updatedAt }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirm Delete Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Component Master"
      :message="`Are you sure you want to delete '${component.partNumber}'?`"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />
  </div>

  <div v-else class="p-8 text-center bg-white rounded-lg border border-slate-200">
    <AppIcon name="alert-circle" size="lg" class="mx-auto text-amber-500 mb-2" />
    <h3 class="text-sm font-bold text-slate-800">Component Not Found</h3>
    <p class="text-xs text-slate-500 mt-1">The requested component could not be found in the master library.</p>
    <router-link to="/components" class="btn btn-secondary mt-4">
      Back to Components List
    </router-link>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import StatusBadge from '@/components/common/StatusBadge.vue';
import EmptyState from '@/components/common/EmptyState.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ComponentDetailView',
  components: {
    StatusBadge,
    EmptyState,
    ConfirmDialog,
    AppIcon
  },
  data() {
    return {
      activeTab: 'overview',
      deleteModalOpen: false
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    componentId() {
      return this.$route.params.id;
    },
    component() {
      return this.store.components.find(c => c.id === this.componentId || c.partNumber === this.componentId);
    },
    tabs() {
      const c = this.component || {};
      return [
        { key: 'overview', label: 'Overview', icon: 'component' },
        { key: 'specs', label: 'Technical Specifications', icon: 'zap', badge: (c.specs || []).length },
        { key: 'documents', label: 'Documents', icon: 'file-text', badge: (c.documents || []).length },
        { key: 'activity', label: 'Activity', icon: 'activity' }
      ];
    },
    relatedActivities() {
      return this.store.activities.filter(a => a.targetId === this.componentId || (this.component && a.targetName.includes(this.component.partNumber)));
    }
  },
  methods: {
    formatCurrency(amount, currency) {
      if (amount === null || amount === undefined || amount === 0) return '—';
      const curr = (currency || 'IDR').toUpperCase();
      if (curr === 'IDR') {
        return new Intl.NumberFormat('id-ID', {
          style: 'currency',
          currency: 'IDR',
          minimumFractionDigits: 0,
          maximumFractionDigits: 0
        }).format(amount);
      }
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: curr,
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(amount);
    },
    confirmDelete() {
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.component) {
        this.store.deleteComponent(this.component.id);
        this.$router.push('/components');
      }
    },
    previewDoc(doc) {
      this.store.showToast('Document Preview', `Opening ${doc.name}...`, 'info');
    },
    downloadDoc(doc) {
      this.store.showToast('Download Started', `Downloading ${doc.name} (${doc.size})`, 'success');
    },
    downloadAllDocs() {
      this.store.showToast('Datasheet Package', 'Preparing complete documentation archive for ' + this.component.partNumber, 'info');
    },
    triggerUploadDoc() {
      this.$router.push(`/components/${this.component.id}/edit`);
    }
  }
};
</script>
