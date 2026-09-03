<template>
  <div class="space-y-5">
    <!-- Page Header -->
    <PageHeader title="Central Components Library" subtitle="Normalized electronic components, semiconductors, parametric specifications, and compliance datasheets">
      <template #badge>
        <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
          {{ store.components.length }} Active Parts
        </span>
      </template>
      <template #actions>
        <button @click="openImportModal" class="btn btn-secondary">
          <AppIcon name="upload" size="xs" class="mr-1.5 text-slate-500" />
          <span>Import CSV</span>
        </button>
        <button @click="exportComponents" class="btn btn-secondary">
          <AppIcon name="download" size="xs" class="mr-1.5 text-slate-500" />
          <span>Export JSON</span>
        </button>
        <router-link to="/components/create" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
          <span>Add Component</span>
        </router-link>
      </template>
    </PageHeader>

    <!-- Advanced Filter & Search Workspace -->
    <div class="app-card p-3.5 space-y-3 bg-white">
      <!-- Top row: Search & Core Dropdowns -->
      <div class="flex flex-wrap items-center gap-2.5">
        <!-- Search Input -->
        <div class="relative min-w-[260px] max-w-sm flex-1">
          <AppIcon name="search" size="xs" class="absolute left-3 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search MPN, part name, package..."
            class="form-input pl-8.5 text-xs font-mono"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-2.5 top-2 text-slate-400 hover:text-slate-600"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Category Dropdown -->
        <select v-model="selectedCategory" class="form-select w-auto min-w-[160px] text-xs">
          <option value="">All Categories</option>
          <option v-for="cat in store.flatCategories" :key="cat.id" :value="cat.rawName">
            {{ cat.name }}
          </option>
        </select>

        <!-- Manufacturer Dropdown -->
        <select v-model="selectedManufacturer" class="form-select w-auto min-w-[150px] text-xs">
          <option value="">All Manufacturers</option>
          <option v-for="m in store.manufacturers" :key="m.id" :value="m.name">
            {{ m.name }}
          </option>
        </select>

        <!-- Supplier Dropdown -->
        <select v-model="selectedSupplier" class="form-select w-auto min-w-[140px] text-xs">
          <option value="">All Suppliers</option>
          <option v-for="s in store.suppliers" :key="s.id" :value="s.name">
            {{ s.name }}
          </option>
        </select>

        <!-- Status Dropdown -->
        <select v-model="selectedStatus" class="form-select w-auto min-w-[120px] text-xs">
          <option value="">All Statuses</option>
          <option value="Active">Active</option>
          <option value="In Review">In Review</option>
          <option value="Draft">Draft</option>
          <option value="EOL">EOL</option>
        </select>

        <!-- Advanced Filter Drawer Toggle -->
        <button
          @click="filterDrawerOpen = true"
          class="btn btn-secondary"
          title="More Filters"
        >
          <AppIcon name="filter" size="xs" class="mr-1 text-slate-500" />
          <span>More Filters</span>
        </button>

        <!-- Reset Button -->
        <button
          v-if="isFiltered"
          @click="resetFilters"
          class="btn btn-ghost text-slate-500 hover:text-slate-800"
        >
          <AppIcon name="refresh" size="xs" class="mr-1" />
          <span>Reset</span>
        </button>
      </div>

      <!-- Quick Category Filter Chips -->
      <div class="flex items-center gap-1.5 pt-2 border-t border-slate-100 overflow-x-auto text-2xs">
        <button
          v-for="chip in categoryChips"
          :key="chip.key"
          @click="setCategoryChip(chip.key)"
          :class="[
            'px-2.5 py-1 rounded-full font-medium transition-all whitespace-nowrap border shadow-2xs',
            activeCategoryChip === chip.key
              ? 'bg-brand-50 border-brand-300 text-brand-700 font-bold'
              : 'bg-slate-50/80 border-slate-200/80 text-slate-600 hover:bg-slate-100 hover:border-slate-300'
          ]"
        >
          {{ chip.label }}
        </button>
      </div>
    </div>

    <!-- Component Data Table Card -->
    <div class="app-card overflow-hidden">
      <div v-if="filteredComponents.length === 0" class="p-8">
        <EmptyState
          title="No components found"
          :description="isFiltered ? 'No engineering components match your search filters.' : 'Get started by creating your first component record.'"
          icon="component"
        >
          <template #action>
            <button v-if="isFiltered" @click="resetFilters" class="btn btn-secondary mr-2">
              Reset Filters
            </button>
            <router-link to="/components/create" class="btn btn-primary">
              <AppIcon name="plus" size="xs" class="mr-1.5" />
              <span>Add Component</span>
            </router-link>
          </template>
        </EmptyState>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs divide-y divide-slate-100">
          <thead>
            <tr class="bg-slate-50/90 text-slate-500 uppercase tracking-wider text-3xs font-bold select-none border-b border-slate-200/80">
              <th class="w-8 py-2.5 px-2 text-center">
                <input
                  type="checkbox"
                  v-model="selectAll"
                  class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer"
                />
              </th>
              <th class="py-2.5 px-2.5 whitespace-nowrap min-w-[140px]">Part Number (MPN)</th>
              <th class="py-2.5 px-2.5 min-w-[180px] lg:min-w-[210px]">Component Details</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Category</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-36">Manufacturer / Sourcing</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Purchase Cost</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-16 text-center">UoM</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24">Status</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr
              v-for="c in paginatedComponents"
              :key="c.id"
              @click="$router.push(`/components/${c.id || c.partNumber}`)"
              class="cursor-pointer hover:bg-slate-50/70 transition-colors"
            >
              <td class="py-2 px-2 text-center" @click.stop>
                <input
                  type="checkbox"
                  :value="c.id"
                  v-model="selectedIds"
                  class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer"
                />
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <div class="flex items-center gap-1.5">
                  <span class="font-mono font-bold text-slate-900 group-hover:text-brand-600 transition-colors text-xs whitespace-nowrap">
                    {{ c.partNumber }}
                  </span>
                  <span
                    v-if="c.rohsCompliant || c.roHsCompliant"
                    class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-emerald-50 text-emerald-700 border border-emerald-200 shrink-0"
                    title="RoHS Compliant"
                  >
                    RoHS
                  </span>
                </div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5">{{ c.internalCode || 'IPN-STD' }}</div>
              </td>
              <td class="py-2 px-2.5">
                <div class="flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-lg bg-slate-100 border border-slate-200 overflow-hidden shrink-0 flex items-center justify-center text-slate-400 shadow-2xs">
                    <img
                      v-if="c.image"
                      :src="c.image"
                      alt="Component"
                      class="w-full h-full object-cover"
                    />
                    <AppIcon v-else name="component" size="xs" />
                  </div>
                  <div class="min-w-0">
                    <div class="text-slate-900 font-bold truncate max-w-[170px] lg:max-w-xs text-xs" :title="c.name">{{ c.name }}</div>
                    <div class="text-3xs text-slate-400 font-mono truncate mt-0.5">Pkg: <span class="font-medium text-slate-600">{{ c.packageType || 'Standard' }}</span></div>
                  </div>
                </div>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 border border-slate-200/80 truncate block max-w-[120px]">
                  {{ c.category }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <div class="font-semibold text-slate-800 text-xs truncate max-w-[140px]" :title="c.manufacturer">
                  {{ c.manufacturer }}
                </div>
                <div class="text-3xs text-slate-400 font-medium truncate max-w-[140px]" :title="'Supplier: ' + (c.supplier || 'Direct')">
                  Dist: <span class="text-slate-600">{{ c.supplier || 'Direct' }}</span>
                </div>
              </td>
              <!-- Purchase Cost & Link -->
              <td class="py-2 px-2.5 whitespace-nowrap">
                <div v-if="c.estimatedUnitCost > 0" class="space-y-0.5">
                  <div class="font-mono font-bold text-xs text-slate-900">
                    {{ formatCurrency(c.estimatedUnitCost, c.currency) }}
                  </div>
                  <div class="text-3xs text-slate-400 font-mono">
                    Per {{ c.unit || 'pcs' }}
                  </div>
                </div>
                <div v-else class="text-slate-300 font-mono text-xs">—</div>
              </td>

              <td class="py-2 px-2.5 whitespace-nowrap text-center">
                <span class="font-mono font-bold text-3xs text-slate-700 bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">
                  {{ c.unit || 'pcs' }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <StatusBadge :status="c.status" />
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <a
                    v-if="c.purchaseLink"
                    :href="c.purchaseLink"
                    target="_blank"
                    rel="noopener noreferrer"
                    @click.stop
                    class="btn-icon text-brand-600 hover:text-brand-800 hover:bg-brand-50"
                    title="Open Purchase Link"
                  >
                    <AppIcon name="shopping-bag" size="xs" />
                  </a>
                  <router-link
                    :to="`/components/${c.id || c.partNumber}`"
                    class="btn-icon"
                    title="View Profile"
                  >
                    <AppIcon name="eye" size="xs" />
                  </router-link>
                  <router-link
                    :to="`/components/${c.id || c.partNumber}/edit`"
                    class="btn-icon"
                    title="Edit Specifications"
                  >
                    <AppIcon name="edit" size="xs" />
                  </router-link>
                  <button
                    @click="confirmDelete(c)"
                    class="btn-icon text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                    title="Delete Component"
                  >
                    <AppIcon name="trash" size="xs" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Table Pagination Bar -->
      <div v-if="filteredComponents.length > 0" class="px-4 py-3 border-t border-slate-100 bg-slate-50/70 flex items-center justify-between text-xs text-slate-600">
        <div class="flex items-center gap-2">
          <span class="text-3xs text-slate-500 font-semibold uppercase tracking-wider">Per page:</span>
          <select v-model="pageSize" class="form-select w-auto py-1 text-2xs font-mono">
            <option :value="5">5</option>
            <option :value="10">10</option>
            <option :value="20">20</option>
          </select>
        </div>
        <div class="flex items-center gap-3 font-mono text-xs">
          <span class="text-3xs text-slate-500">Page <strong class="text-slate-800">{{ currentPage }}</strong> of {{ totalPages || 1 }}</span>
          <div class="flex items-center gap-1">
            <button
              :disabled="currentPage === 1"
              @click="currentPage--"
              class="p-1.5 rounded-lg border border-slate-200 bg-white disabled:opacity-40 hover:bg-slate-100 transition-colors shadow-2xs"
            >
              <AppIcon name="chevron-left" size="xs" />
            </button>
            <button
              :disabled="currentPage >= totalPages"
              @click="currentPage++"
              class="p-1.5 rounded-lg border border-slate-200 bg-white disabled:opacity-40 hover:bg-slate-100 transition-colors shadow-2xs"
            >
              <AppIcon name="chevron-right" size="xs" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Advanced Filter Drawer -->
    <BaseDrawer
      :show="filterDrawerOpen"
      title="Advanced Component Filters"
      subtitle="Refine electronic & mechanical component criteria"
      width="sm"
      @close="filterDrawerOpen = false"
    >
      <div class="space-y-4 text-xs">
        <div>
          <label class="form-label">Mounting Type</label>
          <select v-model="filterMounting" class="form-select">
            <option value="">All Mounting Types</option>
            <option value="Surface Mount">Surface Mount (SMD)</option>
            <option value="Through Hole">Through Hole (THT)</option>
            <option value="Panel Mount">Panel Mount</option>
          </select>
        </div>

        <div>
          <label class="form-label">Compliance Filter</label>
          <label class="flex items-center gap-2 cursor-pointer mt-1">
            <input type="checkbox" v-model="filterRohsOnly" class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer" />
            <span class="text-slate-700 font-medium">RoHS 3 Compliant Only</span>
          </label>
        </div>

        <div>
          <label class="form-label">Packaging Standard</label>
          <select v-model="filterUnit" class="form-select">
            <option value="">All Units</option>
            <option v-for="u in store.units" :key="u.id" :value="u.code">{{ u.name }} ({{ u.code }})</option>
          </select>
        </div>
      </div>

      <template #footer>
        <button class="btn btn-secondary" @click="resetFilters">Reset</button>
        <button class="btn btn-primary" @click="filterDrawerOpen = false">Apply Filters</button>
      </template>
    </BaseDrawer>

    <!-- Confirm Delete Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Engineering Component"
      :message="`Are you sure you want to delete component '${componentToDelete ? componentToDelete.partNumber : ''}'?`"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />

    <!-- Confirm Delete Selected Dialog -->
    <ConfirmDialog
      :show="deleteSelectedModalOpen"
      title="Delete Selected Components"
      :message="`Are you sure you want to delete the ${selectedIds.length} selected components? This action cannot be undone.`"
      confirmText="Delete Selected"
      @confirm="executeDeleteSelected"
      @cancel="deleteSelectedModalOpen = false"
    />

    <!-- Danger Purge / Delete All Modal (Super Admin) -->
    <BaseModal :show="deleteAllModalOpen" title="⚠️ Danger Zone: Delete All Components" maxWidth="md" @close="closeDeleteAllModal">
      <div class="space-y-4 text-xs">
        <div class="p-3.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-800 space-y-2">
          <div class="flex items-center gap-2 font-bold text-rose-900 text-sm">
            <AppIcon name="alert-triangle" size="sm" class="text-rose-600 shrink-0" />
            <span>Permanent Purge Warning</span>
          </div>
          <p class="leading-relaxed">
            You are about to delete <strong class="underline font-mono">{{ store.components.length }} engineering components</strong> from the central library. This will permanently erase component datasheets, pin configurations, and catalog mappings.
          </p>
          <div class="text-3xs font-semibold uppercase tracking-wider text-rose-700 bg-rose-100/80 px-2 py-1 rounded border border-rose-300">
            Authorization: Super Admin Privilege Only
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">
            To proceed, type <span class="font-mono font-bold text-rose-600 select-all bg-rose-50 px-1 py-0.5 rounded border border-rose-200">DELETE ALL</span> below:
          </label>
          <input
            v-model="confirmDeleteAllText"
            type="text"
            placeholder="Type DELETE ALL to confirm"
            class="form-input text-xs font-mono text-center tracking-wider font-bold border-rose-300 focus:border-rose-500 focus:ring-rose-500"
            autocomplete="off"
            @keyup.enter="confirmDeleteAllText.trim() === 'DELETE ALL' && !isDeletingAll && executeDeleteAll()"
          />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeDeleteAllModal" :disabled="isDeletingAll">Cancel</button>
        <button
          class="btn bg-rose-600 hover:bg-rose-700 text-white font-bold disabled:opacity-40 disabled:cursor-not-allowed shadow-md shadow-rose-600/20 flex items-center gap-1.5"
          :disabled="confirmDeleteAllText.trim() !== 'DELETE ALL' || isDeletingAll"
          @click="executeDeleteAll"
        >
          <AppIcon v-if="!isDeletingAll" name="trash" size="xs" />
          <span>{{ isDeletingAll ? 'Purging Library...' : `Delete All (${store.components.length}) Components` }}</span>
        </button>
      </template>
    </BaseModal>

    <!-- Floating Bulk Action Bar (Option B) -->
    <transition
      enter-active-class="transition transform duration-200 ease-out"
      enter-from-class="translate-y-12 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition transform duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-12 opacity-0"
    >
      <div
        v-if="selectedIds.length > 0"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 z-40 bg-slate-900/95 text-white border border-slate-700/80 shadow-2xl rounded-2xl px-5 py-3 flex items-center gap-4 backdrop-blur-md max-w-[90vw]"
      >
        <div class="flex items-center gap-2 pr-3 border-r border-slate-700/80">
          <div class="w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center font-mono font-bold text-xs">
            {{ selectedIds.length }}
          </div>
          <span class="text-xs font-semibold text-slate-200">Selected</span>
        </div>

        <button
          @click="selectedIds = []"
          class="text-xs text-slate-400 hover:text-white transition-colors underline-offset-2 hover:underline"
        >
          Deselect
        </button>

        <div class="flex items-center gap-2">
          <!-- Delete Selected Button -->
          <button
            @click="confirmDeleteSelected"
            class="btn btn-sm bg-rose-600/90 hover:bg-rose-600 text-white font-semibold text-xs border border-rose-500/40 shadow-sm flex items-center gap-1.5"
          >
            <AppIcon name="trash" size="xs" />
            <span>Delete Selected ({{ selectedIds.length }})</span>
          </button>

          <!-- Delete All Records Button (Super Admin Only) -->
          <button
            v-if="isSuperAdmin"
            @click="openDeleteAllModal"
            class="btn btn-sm bg-rose-950 hover:bg-rose-900 text-rose-300 hover:text-rose-100 font-semibold text-xs border border-rose-700/80 shadow-sm flex items-center gap-1.5 group"
            title="Super Admin Only: Purge entire components database"
          >
            <AppIcon name="alert-triangle" size="xs" class="text-rose-400 group-hover:scale-110 transition-transform" />
            <span>Delete All ({{ store.components.length }})</span>
            <span class="text-3xs font-mono px-1.5 py-0.2 rounded bg-rose-900/80 border border-rose-600/60 text-rose-200 uppercase tracking-widest font-bold">
              Super Admin
            </span>
          </button>
        </div>
      </div>
    </transition>

    <!-- Import Modal -->
    <BaseModal :show="importModalOpen" title="Import Components from CSV" maxWidth="lg" @close="importModalOpen = false">
      <div class="space-y-3">
        <div class="border-2 border-dashed border-slate-300 rounded-xl p-6 text-center hover:bg-slate-50 transition-colors cursor-pointer">
          <AppIcon name="upload" size="lg" class="mx-auto text-brand-600 mb-2" />
          <p class="text-xs font-bold text-slate-800">Upload Component Master CSV</p>
          <p class="text-3xs text-slate-500 mt-1">Headers: partNumber, name, category, manufacturer, supplier, unit, packageType</p>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="importModalOpen = false">Cancel</button>
        <button class="btn btn-primary" @click="executeImport">Upload & Process</button>
      </template>
    </BaseModal>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import StatusBadge from '@/components/common/StatusBadge.vue';
import EmptyState from '@/components/common/EmptyState.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import BaseDrawer from '@/components/common/BaseDrawer.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ComponentListView',
  components: {
    PageHeader,
    StatusBadge,
    EmptyState,
    ConfirmDialog,
    BaseDrawer,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedCategory: '',
      selectedManufacturer: '',
      selectedSupplier: '',
      selectedStatus: '',
      activeCategoryChip: 'all',
      filterMounting: '',
      filterRohsOnly: false,
      filterUnit: '',
      filterDrawerOpen: false,
      selectedIds: [],
      currentPage: 1,
      pageSize: 10,
      deleteModalOpen: false,
      componentToDelete: null,
      deleteSelectedModalOpen: false,
      deleteAllModalOpen: false,
      confirmDeleteAllText: '',
      isDeletingAll: false,
      importModalOpen: false
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    isSuperAdmin() {
      return this.store.isSuperAdmin;
    },
    categoryChips() {
      return [
        { key: 'all', label: 'All Components' },
        { key: 'mcu', label: 'Microcontrollers' },
        { key: 'sensors', label: 'Sensors' },
        { key: 'rf', label: 'Wireless & RF' },
        { key: 'power', label: 'Power & PMIC' },
        { key: 'mech', label: 'Mechanical & Connectors' },
        { key: 'passives', label: 'Passives' }
      ];
    },
    isFiltered() {
      return Boolean(
        this.searchQuery ||
        this.selectedCategory ||
        this.selectedManufacturer ||
        this.selectedSupplier ||
        this.selectedStatus ||
        this.activeCategoryChip !== 'all' ||
        this.filterMounting ||
        this.filterRohsOnly ||
        this.filterUnit
      );
    },
    filteredComponents() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.store.components.filter(c => {
        const matchesSearch = !q || (c.partNumber && c.partNumber.toLowerCase().includes(q)) || (c.name && c.name.toLowerCase().includes(q)) || (c.packageType && c.packageType.toLowerCase().includes(q));
        const matchesCategory = !this.selectedCategory || c.category === this.selectedCategory;
        const matchesMfg = !this.selectedManufacturer || c.manufacturer === this.selectedManufacturer;
        const matchesSupplier = !this.selectedSupplier || c.supplier === this.selectedSupplier;
        const matchesStatus = !this.selectedStatus || c.status === this.selectedStatus;
        const matchesUnit = !this.filterUnit || c.unit === this.filterUnit;
        const matchesRohs = !this.filterRohsOnly || Boolean(c.rohsCompliant || c.roHsCompliant);
        const matchesMounting = !this.filterMounting || (c.mountingType && c.mountingType.includes(this.filterMounting));

        // Chip logic
        let matchesChip = true;
        if (this.activeCategoryChip === 'mcu') matchesChip = c.category && c.category.includes('Microcontroller');
        else if (this.activeCategoryChip === 'sensors') matchesChip = c.category && c.category.includes('Sensor');
        else if (this.activeCategoryChip === 'rf') matchesChip = c.category && (c.category.includes('Wireless') || c.category.includes('RF'));
        else if (this.activeCategoryChip === 'power') matchesChip = c.category && c.category.includes('Power');
        else if (this.activeCategoryChip === 'mech') matchesChip = c.category && (c.category.includes('Enclosure') || c.category.includes('Connector'));
        else if (this.activeCategoryChip === 'passives') matchesChip = c.category && c.category.includes('Passive');

        return matchesSearch && matchesCategory && matchesMfg && matchesSupplier && matchesStatus && matchesUnit && matchesRohs && matchesMounting && matchesChip;
      });
    },
    totalPages() {
      return Math.ceil(this.filteredComponents.length / this.pageSize);
    },
    paginatedComponents() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.filteredComponents.slice(start, start + this.pageSize);
    },
    selectAll: {
      get() {
        return this.filteredComponents.length > 0 && this.selectedIds.length === this.filteredComponents.length;
      },
      set(val) {
        if (val) {
          this.selectedIds = this.filteredComponents.map(c => c.id);
        } else {
          this.selectedIds = [];
        }
      }
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
    formatTimestamp(ts) {
      if (!ts) return 'Recent';
      if (ts.length > 10) return ts.substring(0, 10);
      return ts;
    },
    setCategoryChip(key) {
      this.activeCategoryChip = key;
      this.currentPage = 1;
    },
    resetFilters() {
      this.searchQuery = '';
      this.selectedCategory = '';
      this.selectedManufacturer = '';
      this.selectedSupplier = '';
      this.selectedStatus = '';
      this.activeCategoryChip = 'all';
      this.filterMounting = '';
      this.filterRohsOnly = false;
      this.filterUnit = '';
      this.filterDrawerOpen = false;
      this.currentPage = 1;
    },
    confirmDelete(c) {
      this.componentToDelete = c;
      this.deleteModalOpen = true;
    },
    async executeDelete() {
      if (this.componentToDelete) {
        await this.store.deleteComponent(this.componentToDelete.id || this.componentToDelete.partNumber);
        this.deleteModalOpen = false;
        this.componentToDelete = null;
      }
    },
    confirmDeleteSelected() {
      if (this.selectedIds.length === 0) return;
      this.deleteSelectedModalOpen = true;
    },
    async executeDeleteSelected() {
      await this.store.deleteMultipleComponents(this.selectedIds);
      this.selectedIds = [];
      this.deleteSelectedModalOpen = false;
    },
    openDeleteAllModal() {
      this.confirmDeleteAllText = '';
      this.deleteAllModalOpen = true;
    },
    closeDeleteAllModal() {
      this.deleteAllModalOpen = false;
      this.confirmDeleteAllText = '';
      this.isDeletingAll = false;
    },
    async executeDeleteAll() {
      if (this.confirmDeleteAllText.trim() !== 'DELETE ALL') return;
      this.isDeletingAll = true;
      try {
        await this.store.deleteAllComponents();
        this.selectedIds = [];
        this.closeDeleteAllModal();
      } finally {
        this.isDeletingAll = false;
      }
    },
    openImportModal() {
      this.importModalOpen = true;
    },
    executeImport() {
      this.importModalOpen = false;
      this.store.showToast('Components Imported', 'Successfully imported components from CSV.', 'success');
    },
    exportComponents() {
      const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(this.store.components, null, 2));
      const downloadAnchor = document.createElement('a');
      downloadAnchor.setAttribute("href", dataStr);
      downloadAnchor.setAttribute("download", "iot_components_engineering_library.json");
      document.body.appendChild(downloadAnchor);
      downloadAnchor.click();
      downloadAnchor.remove();
      this.store.showToast('Export Completed', 'Exported components engineering database.', 'success');
    }
  }
};
</script>
