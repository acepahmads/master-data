<template>
  <div class="space-y-5">
    <!-- Page Header -->
    <PageHeader title="IoT Products Master" subtitle="Manage electronics product definitions, target applications, and product architecture classifications">
      <template #badge>
        <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-brand-50 text-brand-700 border border-brand-200">
          {{ store.products.length }} Products
        </span>
      </template>
      <template #actions>
        <button @click="openImportModal" class="btn btn-secondary">
          <AppIcon name="upload" size="xs" class="mr-1.5 text-slate-500" />
          <span>Import CSV</span>
        </button>
        <button @click="exportProducts" class="btn btn-secondary">
          <AppIcon name="download" size="xs" class="mr-1.5 text-slate-500" />
          <span>Export JSON</span>
        </button>
        <router-link to="/products/create" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
          <span>Add Product</span>
        </router-link>
      </template>
    </PageHeader>

    <!-- Product Type Classification Filter Tabs -->
    <div class="flex items-center gap-2 border-b border-slate-200 pb-2">
      <button
        @click="selectedProductType = ''"
        :class="[
          'px-3 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center gap-2',
          selectedProductType === ''
            ? 'bg-slate-900 text-white shadow-2xs'
            : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'
        ]"
      >
        <span>All Products</span>
        <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full" :class="selectedProductType === '' ? 'bg-slate-700 text-slate-200' : 'bg-slate-100 text-slate-600'">
          {{ store.products.length }}
        </span>
      </button>

      <button
        @click="selectedProductType = 'RND'"
        :class="[
          'px-3 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center gap-2',
          selectedProductType === 'RND'
            ? 'bg-emerald-600 text-white shadow-2xs'
            : 'bg-white text-slate-600 hover:bg-emerald-50 hover:text-emerald-700 border border-slate-200'
        ]"
      >
        <AppIcon name="cpu" size="xs" />
        <span>R&D Products</span>
        <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full" :class="selectedProductType === 'RND' ? 'bg-emerald-700 text-emerald-100' : 'bg-emerald-50 text-emerald-700'">
          {{ rndCount }}
        </span>
      </button>

      <button
        @click="selectedProductType = 'TRADING'"
        :class="[
          'px-3 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center gap-2',
          selectedProductType === 'TRADING'
            ? 'bg-purple-600 text-white shadow-2xs'
            : 'bg-white text-slate-600 hover:bg-purple-50 hover:text-purple-700 border border-slate-200'
        ]"
      >
        <AppIcon name="shopping-bag" size="xs" />
        <span>Trading Products</span>
        <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full" :class="selectedProductType === 'TRADING' ? 'bg-purple-700 text-purple-100' : 'bg-purple-50 text-purple-700'">
          {{ trdCount }}
        </span>
      </button>

      <button
        @click="selectedProductType = 'PROJECT'"
        :class="[
          'px-3 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center gap-2',
          selectedProductType === 'PROJECT'
            ? 'bg-cyan-600 text-white shadow-2xs'
            : 'bg-white text-slate-600 hover:bg-cyan-50 hover:text-cyan-700 border border-slate-200'
        ]"
      >
        <AppIcon name="layers" size="xs" />
        <span>Project Solutions</span>
        <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full" :class="selectedProductType === 'PROJECT' ? 'bg-cyan-700 text-cyan-100' : 'bg-cyan-50 text-cyan-700'">
          {{ prjCount }}
        </span>
      </button>
    </div>

    <!-- Filter & Search Workspace -->
    <div class="app-card p-3.5 flex flex-col md:flex-row md:items-center justify-between gap-3 bg-white">
      <div class="flex flex-wrap items-center gap-2.5 flex-1">
        <!-- Search Input -->
        <div class="relative min-w-[240px] max-w-sm flex-1">
          <AppIcon name="search" size="xs" class="absolute left-3 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by code, product name, or lead..."
            class="form-input pl-8.5 text-xs"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-2.5 top-2 text-slate-400 hover:text-slate-600"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Category Filter -->
        <select v-model="selectedCategory" class="form-select w-auto min-w-[150px] text-xs">
          <option value="">All Categories</option>
          <option v-for="cat in uniqueCategories" :key="cat" :value="cat">{{ cat }}</option>
        </select>

        <!-- Status Filter -->
        <select v-model="selectedStatus" class="form-select w-auto min-w-[130px] text-xs">
          <option value="">All Statuses</option>
          <option value="Active">Active</option>
          <option value="Draft">Draft</option>
          <option value="In Review">In Review</option>
          <option value="Obsolete">Obsolete</option>
        </select>

        <!-- Reset Filter Button -->
        <button
          v-if="isFiltered"
          @click="resetFilters"
          class="btn btn-ghost text-slate-500 hover:text-slate-800"
          title="Reset all filters"
        >
          <AppIcon name="refresh" size="xs" class="mr-1" />
          <span>Reset</span>
        </button>
      </div>

      <!-- Count summary -->
      <div class="text-3xs text-slate-500 font-mono shrink-0">
        Showing <span class="font-bold text-slate-900">{{ filteredProducts.length }}</span> of {{ store.products.length }} products
      </div>
    </div>

    <!-- Product Data Table Card -->
    <div class="app-card overflow-hidden">
      <div v-if="filteredProducts.length === 0" class="p-8">
        <EmptyState
          title="No products found"
          :description="searchQuery || selectedCategory || selectedStatus || selectedProductType ? 'No products match your active search filters.' : 'Get started by creating your first IoT Product Master.'"
          icon="product"
        >
          <template #action>
            <button
              v-if="searchQuery || selectedCategory || selectedStatus || selectedProductType"
              @click="resetFilters"
              class="btn btn-secondary mr-2"
            >
              Reset Filters
            </button>
            <router-link to="/products/create" class="btn btn-primary">
              <AppIcon name="plus" size="xs" class="mr-1.5" />
              <span>Create Product</span>
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
              <th class="py-2.5 px-2.5 min-w-[180px] lg:min-w-[200px]">Product Master</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap min-w-[130px]">Product Code</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-20">Type</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Category</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-28">Cost / Purchase</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-32">Selling Price</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-24">Status</th>
              <th class="py-2.5 px-2.5 whitespace-nowrap w-20 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr
              v-for="p in paginatedProducts"
              :key="p.id"
              @click="$router.push(`/products/${p.id || p.code}`)"
              class="cursor-pointer hover:bg-slate-50/70 transition-colors"
            >
              <td class="py-2 px-2 text-center" @click.stop>
                <input
                  type="checkbox"
                  :value="p.id"
                  v-model="selectedIds"
                  class="rounded text-brand-600 focus:ring-brand-500 cursor-pointer"
                />
              </td>
              <td class="py-2 px-2.5">
                <div class="flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-lg bg-slate-100 border border-slate-200/80 overflow-hidden shrink-0 group-hover:border-brand-300 transition-colors shadow-2xs flex items-center justify-center">
                    <img
                      v-if="p.image && !failedImgs[p.id || p.code]"
                      :src="getImageUrl(p.image)"
                      alt="Product Thumb"
                      class="w-full h-full object-cover"
                      @error="onImgError(p.id || p.code)"
                    />
                    <div v-else class="w-full h-full flex items-center justify-center text-slate-400 bg-slate-50">
                      <AppIcon name="product" size="xs" />
                    </div>
                  </div>
                  <div class="min-w-0">
                    <div
                      class="font-bold text-slate-900 hover:text-brand-600 transition-colors truncate max-w-[170px] lg:max-w-xs text-xs"
                      :title="p.name + (p.targetMarket ? ' — ' + p.targetMarket : '')"
                    >
                      {{ p.name }}
                    </div>
                    <div class="text-3xs text-slate-400 font-medium truncate mt-0.5">
                      Lead: <span class="text-slate-600 font-semibold">{{ p.projectLead || 'Unassigned' }}</span>
                    </div>
                  </div>
                </div>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span class="font-mono font-bold text-brand-700 text-xs bg-brand-50/80 px-2 py-0.5 rounded border border-brand-200/60 shadow-2xs whitespace-nowrap inline-block">
                  {{ p.code }}
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span
                  :class="[
                    'text-3xs font-mono font-bold px-2 py-0.5 rounded-full border shadow-2xs inline-flex items-center gap-1 whitespace-nowrap',
                    p.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' :
                    p.productType === 'PROJECT' ? 'bg-cyan-50 text-cyan-700 border-cyan-200' :
                    'bg-emerald-50 text-emerald-700 border-emerald-200'
                  ]"
                >
                  <AppIcon :name="p.productType === 'TRADING' ? 'shopping-bag' : p.productType === 'PROJECT' ? 'layers' : 'cpu'" size="xs" />
                  <span>{{ p.productType || 'RND' }}</span>
                </span>
              </td>
              <td class="py-2 px-2.5 whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 border border-slate-200/80 truncate block max-w-[120px]">
                  {{ p.category }}
                </span>
              </td>

              <!-- COST / PURCHASE COLUMN -->
              <td class="py-2 px-2.5">
                <div
                  v-if="p.costSummary && p.costSummary.available"
                  class="space-y-0.5"
                  :title="p.costSummary.currency !== 'IDR' && p.costSummary.amountInIDR ? '≈ Rp ' + Number(p.costSummary.amountInIDR).toLocaleString('id-ID') : ''"
                >
                  <div class="font-mono font-bold text-xs text-slate-900">
                    {{ formatCurrency(p.costSummary.amount, p.costSummary.currency) }}
                  </div>
                  <div
                    :class="[
                      'text-3xs font-medium truncate',
                      p.productType === 'TRADING' ? 'text-purple-600' :
                      p.productType === 'PROJECT' ? 'text-cyan-600' :
                      'text-emerald-600'
                    ]"
                  >
                    {{ p.costSummary.label }}
                  </div>
                </div>
                <div v-else class="space-y-0.5">
                  <div class="font-mono text-xs text-slate-300">—</div>
                  <div class="text-3xs text-slate-400 font-medium">
                    {{ p.costSummary ? p.costSummary.label : 'No Cost' }}
                  </div>
                </div>
              </td>

              <!-- SELLING PRICE COLUMN -->
              <td class="py-2 px-2.5">
                <div
                  v-if="p.sellingPriceSummary && p.sellingPriceSummary.available"
                  class="space-y-0.5"
                  :title="p.sellingPriceSummary.currency !== 'IDR' && p.sellingPriceSummary.amountInIDR ? '≈ Rp ' + Number(p.sellingPriceSummary.amountInIDR).toLocaleString('id-ID') : ''"
                >
                  <div class="flex items-center gap-1.5">
                    <span
                      :class="[
                        'font-mono font-bold text-xs',
                        p.productType === 'TRADING' ? 'text-purple-700' :
                        p.productType === 'PROJECT' ? 'text-cyan-700' :
                        'text-emerald-700'
                      ]"
                    >
                      {{ formatCurrency(p.sellingPriceSummary.amount, p.sellingPriceSummary.currency) }}
                    </span>
                    <span
                      v-if="p.sellingPriceSummary.margin"
                      class="px-1 py-0.2 rounded font-mono font-bold text-3xs text-brand-700 bg-brand-50 border border-brand-200"
                    >
                      +{{ p.sellingPriceSummary.margin }}%
                    </span>
                  </div>
                  <div class="text-3xs text-slate-500 font-medium truncate">
                    {{ p.sellingPriceSummary.label }}
                  </div>
                </div>
                <div v-else class="space-y-0.5">
                  <div class="font-mono text-xs text-slate-300">—</div>
                  <div class="text-3xs text-slate-400 font-medium">
                    {{ p.sellingPriceSummary ? p.sellingPriceSummary.label : 'Not Calculated' }}
                  </div>
                </div>
              </td>

              <td class="py-2 px-2.5">
                <StatusBadge :status="p.status" />
              </td>

              <td class="py-2 px-2.5 text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <router-link
                    :to="`/products/${p.id || p.code}`"
                    class="btn-icon"
                    title="View Product Details"
                  >
                    <AppIcon name="eye" size="xs" />
                  </router-link>
                  <router-link
                    :to="`/products/${p.id || p.code}/edit`"
                    class="btn-icon"
                    title="Edit Product"
                  >
                    <AppIcon name="edit" size="xs" />
                  </router-link>
                  <button
                    @click="confirmDelete(p)"
                    class="btn-icon text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                    title="Delete Product"
                  >
                    <AppIcon name="trash" size="xs" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Footer -->
      <div v-if="filteredProducts.length > 0" class="px-4 py-3 bg-slate-50/80 border-t border-slate-200 flex items-center justify-between text-xs">
        <div class="text-3xs text-slate-500 font-mono">
          Page {{ currentPage }} of {{ totalPages || 1 }} ({{ filteredProducts.length }} results)
        </div>
        <div class="flex items-center gap-1">
          <button
            @click="currentPage--"
            :disabled="currentPage <= 1"
            class="btn btn-secondary btn-sm"
          >
            Previous
          </button>
          <button
            @click="currentPage++"
            :disabled="currentPage >= totalPages"
            class="btn btn-secondary btn-sm"
          >
            Next
          </button>
        </div>
      </div>
    </div>

      <!-- Confirm Delete Dialog -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Product Master"
      :message="`Are you sure you want to delete '${productToDelete ? productToDelete.name : ''}' (${productToDelete ? productToDelete.code : ''})? This will archive the master definition.`"
      confirmText="Delete Product"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />

    <!-- Confirm Delete Selected Dialog -->
    <ConfirmDialog
      :show="deleteSelectedModalOpen"
      title="Delete Selected Products"
      :message="`Are you sure you want to delete the ${selectedIds.length} selected products? This action cannot be undone.`"
      confirmText="Delete Selected"
      @confirm="executeDeleteSelected"
      @cancel="deleteSelectedModalOpen = false"
    />

    <!-- Danger Purge / Delete All Modal (Super Admin) -->
    <BaseModal :show="deleteAllModalOpen" title="⚠️ Danger Zone: Delete All Products" maxWidth="md" @close="closeDeleteAllModal">
      <div class="space-y-4 text-xs">
        <div class="p-3.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-800 space-y-2">
          <div class="flex items-center gap-2 font-bold text-rose-900 text-sm">
            <AppIcon name="alert-triangle" size="sm" class="text-rose-600 shrink-0" />
            <span>Permanent Purge Warning</span>
          </div>
          <p class="leading-relaxed">
            You are about to delete <strong class="underline font-mono">{{ store.products.length }} products</strong> from the system. This will permanently remove product records, associated bill of materials linkages, and historical definitions.
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
          <span>{{ isDeletingAll ? 'Purging Records...' : `Delete All (${store.products.length}) Products` }}</span>
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
          <div class="w-6 h-6 rounded-full bg-brand-500/20 text-brand-400 flex items-center justify-center font-mono font-bold text-xs">
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
            title="Super Admin Only: Purge entire products database"
          >
            <AppIcon name="alert-triangle" size="xs" class="text-rose-400 group-hover:scale-110 transition-transform" />
            <span>Delete All ({{ store.products.length }})</span>
            <span class="text-3xs font-mono px-1.5 py-0.2 rounded bg-rose-900/80 border border-rose-600/60 text-rose-200 uppercase tracking-widest font-bold">
              Super Admin
            </span>
          </button>
        </div>
      </div>
    </transition>

    <!-- Import Modal Simulation -->
    <BaseModal :show="importModalOpen" title="Import Products from CSV / Excel" maxWidth="lg" @close="importModalOpen = false">
      <div class="space-y-3">
        <div class="border-2 border-dashed border-slate-300 rounded-xl p-6 text-center hover:bg-slate-50 transition-colors cursor-pointer">
          <AppIcon name="upload" size="lg" class="mx-auto text-brand-600 mb-2" />
          <p class="text-xs font-bold text-slate-800">Click or drag CSV file to upload</p>
          <p class="text-3xs text-slate-500 mt-1">Supported formats: .CSV, .XLSX (Max 5MB)</p>
        </div>
        <div class="text-3xs text-slate-500">
          Template headers required: <code class="font-mono bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">code, name, productType, category, status, description</code>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="importModalOpen = false">Cancel</button>
        <button class="btn btn-primary" @click="executeImport">Upload & Map Columns</button>
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
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';
import { getImageUrl } from '@/utils/image';

export default {
  name: 'ProductListView',
  components: {
    PageHeader,
    StatusBadge,
    EmptyState,
    ConfirmDialog,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedCategory: '',
      selectedStatus: '',
      selectedProductType: '',
      selectedIds: [],
      currentPage: 1,
      pageSize: 10,
      deleteModalOpen: false,
      productToDelete: null,
      deleteSelectedModalOpen: false,
      deleteAllModalOpen: false,
      confirmDeleteAllText: '',
      isDeletingAll: false,
      importModalOpen: false,
      failedImgs: {}
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    isSuperAdmin() {
      return this.store.isSuperAdmin;
    },
    rndCount() {
      return this.store.products.filter(p => p.productType === 'RND' || !p.productType).length;
    },
    trdCount() {
      return this.store.products.filter(p => p.productType === 'TRADING').length;
    },
    prjCount() {
      return this.store.products.filter(p => p.productType === 'PROJECT').length;
    },
    uniqueCategories() {
      const set = new Set(this.store.products.map(p => p.category).filter(Boolean));
      return Array.from(set);
    },
    isFiltered() {
      return Boolean(this.searchQuery || this.selectedCategory || this.selectedStatus || this.selectedProductType);
    },
    filteredProducts() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.store.products.filter(p => {
        const matchesSearch = !q || (p.name && p.name.toLowerCase().includes(q)) || (p.code && p.code.toLowerCase().includes(q)) || (p.projectLead && p.projectLead.toLowerCase().includes(q));
        const matchesCategory = !this.selectedCategory || p.category === this.selectedCategory;
        const matchesStatus = !this.selectedStatus || p.status === this.selectedStatus;
        const pType = p.productType || 'RND';
        const matchesType = !this.selectedProductType || pType === this.selectedProductType;
        return matchesSearch && matchesCategory && matchesStatus && matchesType;
      });
    },
    totalPages() {
      return Math.ceil(this.filteredProducts.length / this.pageSize);
    },
    paginatedProducts() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.filteredProducts.slice(start, start + this.pageSize);
    },
    selectAll: {
      get() {
        return this.filteredProducts.length > 0 && this.selectedIds.length === this.filteredProducts.length;
      },
      set(val) {
        if (val) {
          this.selectedIds = this.filteredProducts.map(p => p.id);
        } else {
          this.selectedIds = [];
        }
      }
    }
  },
  async mounted() {
    if (this.store.products.length === 0) {
      await this.store.fetchProducts();
    }
  },
  methods: {
    getImageUrl(url) {
      return getImageUrl(url);
    },
    formatCurrency(amount, currency) {
      if (amount === null || amount === undefined) return '—';
      const curr = (currency || 'USD').toUpperCase();
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
    onImgError(id) {
      if (id) {
        this.failedImgs = { ...this.failedImgs, [id]: true };
      }
    },
    formatTimestamp(ts) {
      if (!ts) return 'Recent';
      if (ts.length > 10) return ts.substring(0, 10);
      return ts;
    },
    resetFilters() {
      this.searchQuery = '';
      this.selectedCategory = '';
      this.selectedStatus = '';
      this.selectedProductType = '';
      this.currentPage = 1;
    },
    confirmDelete(p) {
      this.productToDelete = p;
      this.deleteModalOpen = true;
    },
    async executeDelete() {
      if (this.productToDelete) {
        await this.store.deleteProduct(this.productToDelete.id || this.productToDelete.code);
        this.deleteModalOpen = false;
        this.productToDelete = null;
      }
    },
    confirmDeleteSelected() {
      if (this.selectedIds.length === 0) return;
      this.deleteSelectedModalOpen = true;
    },
    async executeDeleteSelected() {
      await this.store.deleteMultipleProducts(this.selectedIds);
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
        await this.store.deleteAllProducts();
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
      this.store.showToast('Import Simulation', 'Successfully parsed sample products from CSV.', 'success');
    },
    exportProducts() {
      const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(this.store.products, null, 2));
      const downloadAnchor = document.createElement('a');
      downloadAnchor.setAttribute("href", dataStr);
      downloadAnchor.setAttribute("download", "iot_products_master_data.json");
      document.body.appendChild(downloadAnchor);
      downloadAnchor.click();
      downloadAnchor.remove();
      this.store.showToast('Export Completed', 'Exported products master dataset.', 'success');
    }
  }
};
</script>
