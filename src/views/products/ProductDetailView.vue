<template>
  <div v-if="product" class="space-y-5">
    <!-- PRODUCT HEADER -->
    <div class="app-card p-5 bg-white">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div class="flex items-start gap-4">
          <div class="w-16 h-16 rounded-xl object-cover border border-slate-200 shadow-2xs bg-slate-100 shrink-0 overflow-hidden flex items-center justify-center">
            <img
              v-if="product.image && !headerImgFailed"
              :src="getImageUrl(product.image)"
              alt="Product Master"
              class="w-full h-full object-cover"
              @error="headerImgFailed = true"
            />
            <div v-else class="text-slate-400">
              <AppIcon name="product" size="sm" />
            </div>
          </div>
          <div>
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                {{ product.code }}
              </span>
              <!-- PRODUCT TYPE BADGE -->
              <span
                :class="[
                  'text-3xs font-mono font-bold px-2 py-0.5 rounded-full border shadow-2xs inline-flex items-center gap-1',
                  product.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' :
                  product.productType === 'PROJECT' ? 'bg-cyan-50 text-cyan-700 border-cyan-200' :
                  'bg-emerald-50 text-emerald-700 border-emerald-200'
                ]"
              >
                <AppIcon :name="product.productType === 'TRADING' ? 'shopping-bag' : product.productType === 'PROJECT' ? 'layers' : 'cpu'" size="xs" />
                <span>{{ product.productType || 'RND' }} PRODUCT</span>
              </span>
              <span class="px-2 py-0.5 rounded text-2xs font-semibold bg-slate-100 text-slate-700 border border-slate-200">
                {{ product.category }}
              </span>
              <StatusBadge :status="product.status" />
            </div>
            <h1 class="text-xl font-bold text-slate-900 tracking-tight">{{ product.name }}</h1>
            <p class="text-xs text-slate-500 mt-1 max-w-2xl">{{ product.description }}</p>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <router-link :to="`/products/${product.id}/edit`" class="btn btn-primary">
            <AppIcon name="edit" size="xs" class="mr-1.5 text-white" />
            <span>Edit Master</span>
          </router-link>

          <!-- NEW VERSION BUTTON: STRICTLY FOR R&D PRODUCTS ONLY -->
          <router-link
            v-if="isRND"
            :to="`/versions/create?product=${product.code}`"
            class="btn btn-secondary"
          >
            <AppIcon name="plus" size="xs" class="mr-1.5 text-slate-500" />
            <span>New Version</span>
          </router-link>

          <button @click="duplicateProduct" class="btn btn-secondary" title="Duplicate Master">
            <AppIcon name="copy" size="xs" class="mr-1.5 text-slate-500" />
            <span>Duplicate</span>
          </button>
          <button @click="confirmDelete" class="btn btn-danger" title="Delete Product">
            <AppIcon name="trash" size="xs" />
          </button>
        </div>
      </div>

      <!-- Detail Dynamic Tabs Header -->
      <div class="flex items-center gap-6 border-t border-slate-100 mt-5 pt-3 text-xs">
        <button
          @click="activeTab = 'master'"
          :class="[
            'pb-2 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'master'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="file-text" size="xs" />
          <span>Master Profile</span>
        </button>

        <!-- Dynamic Tab for TRADING PRODUCTS -->
        <button
          v-if="isTrading"
          @click="activeTab = 'commercial'"
          :class="[
            'pb-2 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'commercial'
              ? 'border-purple-600 text-purple-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="shopping-bag" size="xs" />
          <span>Commercial & Sourcing</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-purple-50 text-purple-700 border border-purple-200">
            Terms
          </span>
        </button>

        <!-- Dynamic Tab for PROJECT PRODUCTS -->
        <button
          v-if="isProject"
          @click="activeTab = 'structure'"
          :class="[
            'pb-2 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'structure'
              ? 'border-cyan-600 text-cyan-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="layers" size="xs" />
          <span>Project Structure</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-cyan-50 text-cyan-700 border border-cyan-200">
            {{ projectItemCount }}
          </span>
        </button>

        <!-- Dynamic Tab for R&D PRODUCTS -->
        <button
          v-if="isRND"
          @click="activeTab = 'lifecycle'"
          :class="[
            'pb-2 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'lifecycle'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="roadmap" size="xs" />
          <span>Lifecycle & Versions</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-brand-50 text-brand-700 border border-brand-200">
            {{ productVersions.length }}
          </span>
        </button>
      </div>
    </div>

    <!-- TAB 1: MASTER DATA PROFILE (2-COLUMN) -->
    <div v-if="activeTab === 'master'" class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Left 2 Cols: Master Data Information & Detailed Specifications -->
      <div class="lg:col-span-2 space-y-5">
        <!-- PRODUCT INFORMATION -->
        <div class="app-card p-5 bg-white">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100 mb-4">
            <AppIcon name="product" size="xs" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Master Data Information</h2>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Product Code</div>
              <div class="font-mono font-bold text-slate-800 text-sm mt-0.5">{{ product.code }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Product Classification</div>
              <div class="mt-0.5">
                <span
                  :class="[
                    'text-3xs font-mono font-bold px-2 py-0.5 rounded-full border shadow-2xs inline-flex items-center gap-1',
                    product.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' :
                    product.productType === 'PROJECT' ? 'bg-cyan-50 text-cyan-700 border-cyan-200' :
                    'bg-emerald-50 text-emerald-700 border-emerald-200'
                  ]"
                >
                  {{ product.productType || 'RND' }} PRODUCT
                </span>
              </div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Product Name</div>
              <div class="font-bold text-slate-800 mt-0.5">{{ product.name }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Category</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ product.category }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Lifecycle Status</div>
              <div class="mt-1">
                <StatusBadge :status="product.status" />
              </div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Lead Architect / Owner</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ product.projectLead || 'Unassigned' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Target Market / Application</div>
              <div class="text-slate-800 font-medium mt-0.5">{{ product.targetMarket || 'General IoT Industry' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Created Date</div>
              <div class="font-mono text-slate-600 mt-0.5">{{ product.createdAt }}</div>
            </div>
          </div>
        </div>

        <!-- DESCRIPTION SECTION -->
        <div class="app-card p-5 bg-white">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100 mb-3">
            <AppIcon name="file-text" size="xs" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Scope Description</h2>
          </div>
          <div class="text-xs text-slate-700 leading-relaxed font-mono bg-slate-50 p-4 rounded-lg border border-slate-200">
            <p class="whitespace-pre-line m-0">{{ product.longDescription || product.description }}</p>
          </div>
        </div>

        <!-- HARDWARE & COMPLIANCE (FOR R&D) -->
        <div v-if="isRND" class="app-card p-5 bg-white">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100 mb-4">
            <AppIcon name="shield-check" size="xs" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Hardware Characteristics & Compliance</h2>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Power Rating</div>
              <div class="font-mono text-slate-800 mt-0.5">{{ product.powerRating || '9-36V DC, 15W' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Ingress Protection</div>
              <div class="font-mono text-slate-800 mt-0.5">{{ product.ingressProtection || 'IP67' }}</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Operating Temperature</div>
              <div class="font-mono text-slate-800 mt-0.5">{{ product.operatingTemp || '-40°C to +85°C' }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right 1 Col: Quick Sourcing / Lifecycle Context Card -->
      <div class="space-y-5">
        <!-- PRODUCT PHOTOS & GALLERY VIEWER -->
        <div class="app-card p-5 bg-white space-y-3.5">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="camera" size="xs" class="text-brand-600" />
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Product Photos</h3>
            </div>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-700">
              {{ galleryImages.length }} Photos
            </span>
          </div>

          <!-- Active Large Image Preview -->
          <div class="relative rounded-xl overflow-hidden bg-slate-100 border border-slate-200 aspect-video shadow-inner flex items-center justify-center">
            <img
              v-if="activePhoto && !detailImgFailed"
              :src="getImageUrl(activePhoto)"
              alt="Product photo"
              class="w-full h-full object-cover"
              @error="onDetailImgError"
            />
            <div v-else class="text-slate-400 flex flex-col items-center gap-1.5 py-6">
              <AppIcon name="product" size="md" />
              <span class="text-3xs font-medium">No product photo uploaded</span>
            </div>

            <!-- Active Photo Indicator Badge -->
            <div
              v-if="galleryImages.length > 0"
              class="absolute bottom-2 right-2 bg-slate-900/75 text-white text-3xs font-mono font-bold px-2 py-0.5 rounded-full backdrop-blur-xs shadow-xs"
            >
              {{ activePhotoIndex + 1 }} / {{ galleryImages.length }}
            </div>

            <div
              v-if="activePhotoIndex === 0 && galleryImages.length > 0"
              class="absolute top-2 left-2 bg-brand-600 text-white text-3xs font-bold px-2 py-0.5 rounded-full shadow-xs flex items-center gap-1"
            >
              <AppIcon name="check" size="2xs" />
              <span>Cover Photo</span>
            </div>
          </div>

          <!-- Thumbnail Switcher (if multiple images) -->
          <div v-if="galleryImages.length > 1" class="grid grid-cols-5 gap-2 pt-1">
            <button
              v-for="(img, idx) in galleryImages"
              :key="idx"
              type="button"
              @click="activePhotoIndex = idx; detailImgFailed = false"
              class="relative rounded-lg overflow-hidden border aspect-square bg-slate-100 transition-all focus:outline-none"
              :class="activePhotoIndex === idx ? 'border-brand-600 ring-2 ring-brand-500/40 shadow-xs' : 'border-slate-200 opacity-70 hover:opacity-100'"
            >
              <img :src="getImageUrl(img)" alt="thumb" class="w-full h-full object-cover" />
            </button>
          </div>

          <!-- Edit Photos Shortcut -->
          <div class="pt-1">
            <router-link
              :to="`/products/${product.id}/edit`"
              class="text-3xs text-brand-600 hover:text-brand-800 font-semibold flex items-center justify-center gap-1 py-1 rounded bg-brand-50/50 hover:bg-brand-50"
            >
              <AppIcon name="edit" size="2xs" />
              <span>Manage / Upload Photos</span>
            </router-link>
          </div>
        </div>

        <div class="app-card p-5 bg-white space-y-4">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100">
            <AppIcon name="info" size="xs" class="text-brand-600" />
            <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Classification Summary</h3>
          </div>

          <!-- CLASSIFICATION SPECIFIC DETAILS & COST ROLLUPS (PHASE 3) -->
          <div v-if="isTrading" class="space-y-3 text-xs">
            <div class="p-3 bg-purple-50 rounded-lg border border-purple-200">
              <div class="font-bold text-purple-900">Commercial Trading Device</div>
              <p class="text-3xs text-purple-700 mt-1 leading-relaxed">
                This item is an externally sourced commercial product with managed supplier pricing and delivery lead times.
              </p>
            </div>
            <button @click="activeTab = 'commercial'" class="btn btn-secondary w-full text-xs">
              View Commercial Terms
            </button>
          </div>

          <!-- PROJECT COST ROLLUP CARD -->
          <div v-else-if="isProject" class="space-y-3 text-xs">
            <div class="p-3 bg-cyan-50 rounded-lg border border-cyan-200">
              <div class="font-bold text-cyan-900">Project Solution Package</div>
              <p class="text-3xs text-cyan-700 mt-1 leading-relaxed">
                Turnkey system combining internal R&D products, commercial hardware, and installation materials.
              </p>
            </div>

            <!-- Project Aggregated Cost Breakdown -->
            <div v-if="currentProjectCost" class="p-3.5 rounded-xl bg-slate-50 border border-slate-200 space-y-2">
              <div class="flex items-center justify-between text-3xs uppercase font-bold text-slate-400">
                <span>Estimated Solution Cost</span>
                <span class="font-mono text-emerald-700 text-xs">${{ formatNumber(currentProjectCost.totalEstimatedCost) }}</span>
              </div>

              <div class="space-y-1 text-3xs font-mono pt-1">
                <div class="flex justify-between text-slate-500">
                  <span>Trading Hardware:</span>
                  <span class="text-slate-800">${{ formatNumber(currentProjectCost.tradingProductsCost) }}</span>
                </div>
                <div class="flex justify-between text-slate-500">
                  <span>R&D Subsystems:</span>
                  <span class="text-slate-800">${{ formatNumber(currentProjectCost.rndProductsCost) }}</span>
                </div>
                <div class="flex justify-between text-slate-500">
                  <span>Direct Components:</span>
                  <span class="text-slate-800">${{ formatNumber(currentProjectCost.directComponentsCost) }}</span>
                </div>
              </div>

              <div class="pt-2 border-t border-slate-200 flex justify-between items-center text-xs">
                <span class="font-bold text-slate-700">Est. Selling Price:</span>
                <span class="font-mono font-bold text-emerald-700">${{ formatNumber(currentProjectCost.estimatedSellingPrice) }}</span>
              </div>
            </div>

            <button @click="activeTab = 'structure'" class="btn btn-secondary w-full text-xs">
              Open System Architecture Tree
            </button>
          </div>

          <!-- R&D ENGINEERING BOM COST CARD -->
          <div v-else class="space-y-3 text-xs">
            <div class="p-3 bg-emerald-50 rounded-lg border border-emerald-200">
              <div class="font-bold text-emerald-900">Internal R&D Engineering</div>
              <p class="text-3xs text-emerald-700 mt-1 leading-relaxed">
                Hardware revisions track multi-level engineering BOMs with live component line cost rollups.
              </p>
            </div>

            <!-- Active Revision BOM Cost Rollup -->
            <div v-if="currentRNDCost && currentRNDCost.estimatedBomCost > 0" class="p-3.5 rounded-xl bg-slate-50 border border-slate-200 space-y-2">
              <div class="flex items-center justify-between text-3xs font-bold text-slate-500">
                <span>REVISION BASELINE:</span>
                <span class="font-mono px-1.5 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                  {{ currentRNDCost.currentRevision || 'ACTIVE' }} (v{{ currentRNDCost.currentVersion }})
                </span>
              </div>

              <div class="flex items-baseline justify-between pt-1">
                <span class="text-3xs uppercase font-bold text-slate-400">Est. BOM Cost</span>
                <span class="font-mono font-bold text-base text-slate-900">${{ formatNumber(currentRNDCost.estimatedBomCost) }} USD</span>
              </div>

              <div class="flex items-center justify-between text-3xs text-slate-500 font-mono">
                <span>Target Margin: {{ formatNumber(currentRNDCost.targetMargin, 1) }}%</span>
                <span class="text-emerald-700 font-bold">Price: ${{ formatNumber(currentRNDCost.estimatedSellingPrice) }}</span>
              </div>

              <div v-if="currentRNDCost.revisionId" class="pt-2 border-t border-slate-200">
                <router-link
                  :to="`/revisions/${currentRNDCost.revisionId}/bom`"
                  class="btn btn-primary w-full justify-center text-xs py-1.5"
                >
                  <AppIcon name="layers" size="xs" class="mr-1.5 text-white" />
                  <span>Open Revision BOM</span>
                </router-link>
              </div>
            </div>

            <button @click="activeTab = 'lifecycle'" class="btn btn-secondary w-full text-xs">
              View Version Lifecycle
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: TRADING COMMERCIAL SOURCING TAB (FOR TRADING) -->
    <div v-else-if="activeTab === 'commercial' && isTrading">
      <TradingCommercialTab :productId="product.id || product.code" />
    </div>

    <!-- TAB 3: PROJECT STRUCTURE WORKSPACE (FOR PROJECT) -->
    <div v-else-if="activeTab === 'structure' && isProject">
      <ProjectStructureWorkspace :productId="product.id || product.code" />
    </div>

    <!-- TAB 4: R&D LIFECYCLE & VERSIONS (FOR R&D) -->
    <div v-else-if="activeTab === 'lifecycle' && isRND" class="space-y-5">
      <LifecycleTree :productCode="product.code" />
    </div>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Product Master"
      :message="`Are you sure you want to delete '${product.name}' (${product.code})?`"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />
  </div>

  <div v-else class="p-8 text-center bg-white rounded-xl border border-slate-200">
    <AppIcon name="alert-circle" size="lg" class="mx-auto text-amber-500 mb-2" />
    <h3 class="text-sm font-bold text-slate-800">Product Not Found</h3>
    <p class="text-xs text-slate-500 mt-1">The requested product master could not be found.</p>
    <router-link to="/products" class="btn btn-secondary mt-4">
      Back to Products List
    </router-link>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import { useVersionStore } from '@/stores/versionStore';
import { useBOMStore } from '@/stores/bomStore';
import StatusBadge from '@/components/common/StatusBadge.vue';
import LifecycleTree from '@/components/versions/LifecycleTree.vue';
import TradingCommercialTab from '@/components/products/TradingCommercialTab.vue';
import ProjectStructureWorkspace from '@/components/products/ProjectStructureWorkspace.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import AppIcon from '@/components/common/AppIcon.vue';
import { getImageUrl } from '@/utils/image';

export default {
  name: 'ProductDetailView',
  components: {
    StatusBadge,
    LifecycleTree,
    TradingCommercialTab,
    ProjectStructureWorkspace,
    ConfirmDialog,
    AppIcon
  },
  data() {
    return {
      activeTab: 'master',
      deleteModalOpen: false,
      activePhotoIndex: 0,
      headerImgFailed: false,
      detailImgFailed: false
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    versionStore() {
      return useVersionStore();
    },
    bomStore() {
      return useBOMStore();
    },
    productId() {
      return this.$route.params.id;
    },
    product() {
      return this.store.products.find(p => p.id === this.productId || p.code === this.productId);
    },
    galleryImages() {
      if (!this.product) return [];
      if (this.product.images) {
        try {
          const parsed = typeof this.product.images === 'string' ? JSON.parse(this.product.images) : this.product.images;
          if (Array.isArray(parsed) && parsed.length > 0) {
            return parsed.filter(Boolean);
          }
        } catch (e) {
          console.warn('Could not parse product images JSON:', e);
        }
      }
      if (this.product.image) {
        return [this.product.image];
      }
      return [];
    },
    activePhoto() {
      if (this.galleryImages.length === 0) return null;
      const idx = this.activePhotoIndex < this.galleryImages.length ? this.activePhotoIndex : 0;
      return this.galleryImages[idx];
    },
    productType() {
      return this.product?.productType || 'RND';
    },
    isTrading() {
      return this.productType === 'TRADING';
    },
    isProject() {
      return this.productType === 'PROJECT';
    },
    isRND() {
      return this.productType === 'RND';
    },
    productVersions() {
      if (!this.product) return [];
      return this.versionStore.getVersionsForProduct(this.product.code);
    },
    projectItemCount() {
      if (!this.product) return 0;
      const items = this.store.projectItems[this.product.id || this.product.code] || [];
      return items.length;
    },
    currentRNDCost() {
      if (!this.product) return null;
      return this.bomStore.productCostSummary[this.product.id] || null;
    },
    currentProjectCost() {
      if (!this.product) return null;
      return this.bomStore.projectCostSummary[this.product.id] || null;
    }
  },
  methods: {
    getImageUrl(url) {
      return getImageUrl(url);
    },
    onDetailImgError() {
      this.detailImgFailed = true;
    },
    formatNumber(num, decimals = 2) {
      if (num === null || num === undefined) return '0.00';
      return Number(num).toLocaleString('en-US', {
        minimumFractionDigits: decimals,
        maximumFractionDigits: decimals
      });
    },
    confirmDelete() {
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.product) {
        this.store.deleteProduct(this.product.id);
        this.$router.push('/products');
      }
    },
    async duplicateProduct() {
      if (!this.product) return;
      const copy = {
        ...this.product,
        id: null,
        code: `${this.product.code}-COPY`,
        name: `${this.product.name} (Copy)`,
        status: 'Draft'
      };
      const created = await this.store.saveProduct(copy);
      this.$router.push(`/products/${created.id}`);
    },
    async loadProductData() {
      if (this.store.products.length === 0) {
        await this.store.fetchProducts();
      }
      // Ensure activeTab is valid for current product type
      if (this.activeTab === 'structure' && !this.isProject) {
        this.activeTab = 'master';
      } else if (this.activeTab === 'commercial' && !this.isTrading) {
        this.activeTab = 'master';
      } else if (this.activeTab === 'lifecycle' && !this.isRND) {
        this.activeTab = 'master';
      }

      if (this.product) {
        if (this.isTrading) {
          await this.store.fetchTradingDetail(this.product.id);
        } else if (this.isProject) {
          await Promise.all([
            this.store.fetchProjectItems(this.product.id),
            this.bomStore.fetchProjectCostSummary(this.product.id)
          ]);
        } else {
          await Promise.all([
            this.versionStore.fetchVersions(),
            this.versionStore.fetchRevisions(),
            this.bomStore.fetchProductCostSummary(this.product.id)
          ]);
        }
      }
    }
  },
  watch: {
    '$route.params.id': {
      immediate: false,
      handler(newId, oldId) {
        if (newId && newId !== oldId) {
          this.activeTab = 'master';
          this.activePhotoIndex = 0;
          this.headerImgFailed = false;
          this.detailImgFailed = false;
          this.loadProductData();
        }
      }
    },
    productType() {
      if (this.activeTab === 'structure' && !this.isProject) {
        this.activeTab = 'master';
      } else if (this.activeTab === 'commercial' && !this.isTrading) {
        this.activeTab = 'master';
      } else if (this.activeTab === 'lifecycle' && !this.isRND) {
        this.activeTab = 'master';
      }
    }
  },
  async mounted() {
    await this.loadProductData();
  }
};
</script>
