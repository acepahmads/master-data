<template>
  <div class="space-y-5">
    <!-- PAGE HEADER -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-slate-200">
      <div>
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center border border-indigo-200">
            <AppIcon name="list" size="sm" />
          </div>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <span>Prototype Builds Registry</span>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 border border-slate-200">
              {{ filteredPrototypes.length }} of {{ protoStore.prototypes.length }} Builds
            </span>
          </h1>
        </div>
        <p class="text-xs text-slate-500 mt-1 font-medium">
          Detailed engineering directory of all physical prototype units, frozen BOM baselines, and execution records.
        </p>
      </div>

      <div class="flex items-center gap-2.5 shrink-0">
        <router-link to="/prototypes" class="btn btn-secondary text-xs flex items-center gap-1.5">
          <AppIcon name="dashboard" size="xs" />
          <span>Dashboard</span>
        </router-link>
        <router-link to="/prototypes/create" class="btn btn-primary text-xs flex items-center gap-1.5 bg-gradient-to-r from-indigo-600 to-brand-600 border-none shadow-sm">
          <AppIcon name="plus" size="xs" />
          <span>+ New Prototype Build</span>
        </router-link>
      </div>
    </div>

    <!-- FILTER BAR -->
    <div class="app-card p-4 bg-white space-y-3 shadow-card">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-3 text-xs">
        <!-- Search -->
        <div class="md:col-span-2 relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search code, name, purpose, or owner..."
            class="form-input text-xs pl-8 w-full"
          />
          <AppIcon name="search" size="xs" class="absolute left-2.5 top-2.5 text-slate-400" />
        </div>

        <!-- Product Filter -->
        <div>
          <select v-model="selectedProduct" class="form-select text-xs w-full">
            <option value="">All R&D Products</option>
            <option v-for="prod in uniqueProducts" :key="prod.code" :value="prod.code">
              {{ prod.code }} — {{ prod.name }}
            </option>
          </select>
        </div>

        <!-- Status Filter -->
        <div>
          <select v-model="selectedStatus" class="form-select text-xs w-full">
            <option value="">All Build Statuses</option>
            <option value="DRAFT">DRAFT</option>
            <option value="PLANNING">PLANNING</option>
            <option value="COMPONENT_PREPARATION">COMPONENT PREPARATION</option>
            <option value="IN_ASSEMBLY">IN ASSEMBLY</option>
            <option value="ASSEMBLY_COMPLETE">ASSEMBLY COMPLETE</option>
            <option value="READY_FOR_TESTING">READY FOR TESTING</option>
            <option value="COMPLETED">COMPLETED</option>
            <option value="CANCELLED">CANCELLED</option>
          </select>
        </div>

        <!-- Owner Filter -->
        <div>
          <select v-model="selectedOwner" class="form-select text-xs w-full">
            <option value="">All Engineers</option>
            <option v-for="owner in uniqueOwners" :key="owner" :value="owner">
              {{ owner }}
            </option>
          </select>
        </div>
      </div>

      <!-- Filter Active Summary & Reset -->
      <div v-if="isFiltered" class="flex items-center justify-between pt-2 border-t border-slate-100 text-3xs">
        <span class="text-slate-500 font-medium">
          Showing {{ filteredPrototypes.length }} filtered build records.
        </span>
        <button
          @click="resetFilters"
          class="text-brand-600 hover:text-brand-800 font-bold flex items-center gap-1"
        >
          <AppIcon name="refresh" size="2xs" />
          <span>Reset Filters</span>
        </button>
      </div>
    </div>

    <!-- PROTOTYPE BUILDS TABLE -->
    <div class="app-card overflow-hidden bg-white shadow-card">
      <div class="overflow-x-auto">
        <table class="app-table w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-700 text-3xs uppercase font-mono tracking-wider border-b border-slate-200">
              <th class="py-3 px-4">Prototype Code</th>
              <th class="py-3 px-4">Product & Build Name</th>
              <th class="py-3 px-4">Version & Revision</th>
              <th class="py-3 px-4">Build # / Qty</th>
              <th class="py-3 px-4">Purpose</th>
              <th class="py-3 px-4">Build Status</th>
              <th class="py-3 px-4">Assembly Progress</th>
              <th class="py-3 px-4">Engineer Owner</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-xs">
            <tr
              v-for="p in filteredPrototypes"
              :key="p.id"
              class="hover:bg-slate-50/80 transition-colors group"
            >
              <!-- Code -->
              <td class="py-3 px-4 font-mono font-bold text-indigo-700">
                <router-link :to="`/prototypes/${p.id}`" class="hover:underline">
                  {{ p.code }}
                </router-link>
              </td>

              <!-- Product & Name -->
              <td class="py-3 px-4">
                <div class="font-bold text-slate-900 group-hover:text-indigo-600 transition-colors">
                  {{ p.name }}
                </div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5">
                  {{ p.productName }} ({{ p.productCode }})
                </div>
              </td>

              <!-- Version & Revision -->
              <td class="py-3 px-4 font-mono text-3xs">
                <div class="flex items-center gap-1.5">
                  <span class="px-1.5 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200 font-bold">
                    {{ p.versionCode }}
                  </span>
                  <span class="text-slate-400">➔</span>
                  <span class="px-1.5 py-0.5 rounded bg-emerald-50 text-emerald-700 border border-emerald-200 font-bold">
                    {{ p.hardwareRevisionCode }}
                  </span>
                </div>
              </td>

              <!-- Build # / Qty -->
              <td class="py-3 px-4 font-mono text-3xs">
                <span class="font-bold text-slate-800">Build #{{ p.buildNumber }}</span>
                <span class="text-slate-400 block">{{ p.quantity }} Units</span>
              </td>

              <!-- Purpose -->
              <td class="py-3 px-4 text-3xs font-medium text-slate-700">
                <span class="px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 border border-slate-200 inline-block">
                  {{ p.purpose }}
                </span>
              </td>

              <!-- Build Status -->
              <td class="py-3 px-4">
                <span
                  class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full border shadow-2xs inline-block"
                  :class="getStatusBadgeClass(p.status)"
                >
                  {{ formatStatus(p.status) }}
                </span>
              </td>

              <!-- Assembly Progress -->
              <td class="py-3 px-4 w-36">
                <div class="flex items-center justify-between text-3xs font-mono mb-1">
                  <span class="text-slate-400">{{ getCompletedStagesCount(p) }}/8 Stages</span>
                  <span class="font-bold text-indigo-600">{{ getAssemblyProgress(p) }}%</span>
                </div>
                <div class="w-full h-1.5 rounded-full bg-slate-100 overflow-hidden border border-slate-200/60">
                  <div
                    class="h-full rounded-full bg-gradient-to-r from-indigo-500 to-brand-500 transition-all duration-300"
                    :style="{ width: `${getAssemblyProgress(p)}%` }"
                  ></div>
                </div>
              </td>

              <!-- Engineer -->
              <td class="py-3 px-4 text-xs font-medium text-slate-700">
                <div class="flex items-center gap-1.5">
                  <div class="w-5 h-5 rounded-full bg-indigo-100 text-indigo-700 font-bold text-3xs flex items-center justify-center">
                    {{ p.engineeringOwner ? p.engineeringOwner.charAt(0) : 'E' }}
                  </div>
                  <span>{{ p.engineeringOwner }}</span>
                </div>
                <div class="text-3xs font-mono text-slate-400 mt-0.5">Due: {{ p.targetBuildDate }}</div>
              </td>

              <!-- Actions -->
              <td class="py-3 px-4 text-right">
                <router-link
                  :to="`/prototypes/${p.id}`"
                  class="btn btn-secondary py-1 px-2.5 text-xs inline-flex items-center gap-1 group-hover:bg-indigo-50 group-hover:text-indigo-700 group-hover:border-indigo-200"
                >
                  <span>Workspace</span>
                  <AppIcon name="chevron-right" size="2xs" />
                </router-link>
              </td>
            </tr>

            <tr v-if="filteredPrototypes.length === 0">
              <td colspan="9" class="py-12 text-center text-slate-400">
                <div class="flex flex-col items-center justify-center gap-2">
                  <AppIcon name="cpu" size="lg" class="text-slate-300" />
                  <span class="text-xs font-bold text-slate-700">No prototype builds found</span>
                  <span class="text-3xs text-slate-400">Try adjusting your filters or create a new prototype build.</span>
                  <button @click="resetFilters" class="btn btn-secondary btn-sm text-xs mt-2">
                    Clear Filters
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { usePrototypeStore } from '@/stores/prototypeStore';
import { useMasterStore } from '@/stores/masterStore';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'PrototypeListView',
  components: {
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedProduct: '',
      selectedStatus: '',
      selectedOwner: ''
    };
  },
  computed: {
    protoStore() {
      return usePrototypeStore();
    },
    masterStore() {
      return useMasterStore();
    },
    uniqueProducts() {
      const prods = {};
      this.protoStore.prototypes.forEach(p => {
        if (!prods[p.productCode]) {
          prods[p.productCode] = { code: p.productCode, name: p.productName };
        }
      });
      return Object.values(prods);
    },
    uniqueOwners() {
      const set = new Set(this.protoStore.prototypes.map(p => p.engineeringOwner).filter(Boolean));
      return Array.from(set);
    },
    isFiltered() {
      return Boolean(this.searchQuery || this.selectedProduct || this.selectedStatus || this.selectedOwner);
    },
    filteredPrototypes() {
      return this.protoStore.prototypes.filter(p => {
        if (this.searchQuery) {
          const q = this.searchQuery.toLowerCase();
          const matchCode = p.code.toLowerCase().includes(q);
          const matchName = p.name.toLowerCase().includes(q);
          const matchPurpose = (p.purpose || '').toLowerCase().includes(q);
          const matchOwner = (p.engineeringOwner || '').toLowerCase().includes(q);
          if (!matchCode && !matchName && !matchPurpose && !matchOwner) return false;
        }

        if (this.selectedProduct && p.productCode !== this.selectedProduct && p.productId !== this.selectedProduct) {
          return false;
        }

        if (this.selectedStatus && p.status !== this.selectedStatus) {
          return false;
        }

        if (this.selectedOwner && p.engineeringOwner !== this.selectedOwner) {
          return false;
        }

        return true;
      });
    }
  },
  methods: {
    formatStatus(status) {
      if (!status) return 'DRAFT';
      return status.replace(/_/g, ' ');
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'COMPLETED':
          return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'READY_FOR_TESTING':
          return 'bg-purple-50 text-purple-700 border-purple-200';
        case 'IN_ASSEMBLY':
        case 'ASSEMBLY_COMPLETE':
          return 'bg-blue-50 text-blue-700 border-blue-200';
        case 'COMPONENT_PREPARATION':
          return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'CANCELLED':
          return 'bg-red-50 text-red-700 border-red-200';
        default:
          return 'bg-slate-100 text-slate-600 border-slate-200';
      }
    },
    getCompletedStagesCount(p) {
      if (!p || !p.assemblyStages) return 0;
      return p.assemblyStages.filter(s => s.status === 'COMPLETED').length;
    },
    getAssemblyProgress(p) {
      if (!p || !p.assemblyStages || p.assemblyStages.length === 0) return 0;
      const count = p.assemblyStages.filter(s => s.status === 'COMPLETED').length;
      return Math.round((count / p.assemblyStages.length) * 100);
    },
    resetFilters() {
      this.searchQuery = '';
      this.selectedProduct = '';
      this.selectedStatus = '';
      this.selectedOwner = '';
    }
  },
  async mounted() {
    try {
      await this.protoStore.fetchPrototypes();
    } catch (err) {
      console.error('[PrototypeListView] mounted error:', err);
    }
  }
};
</script>
