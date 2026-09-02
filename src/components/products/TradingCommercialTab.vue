<template>
  <div class="space-y-6">
    <!-- Header Banner -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl bg-purple-50 border border-purple-200">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-purple-600 text-white flex items-center justify-center shadow-xs">
          <AppIcon name="shopping-bag" size="sm" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-purple-950">Trading Commercial & Sourcing Master</h2>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-purple-200/70 text-purple-800">
              TRADING PRODUCT
            </span>
          </div>
          <p class="text-xs text-purple-700 mt-0.5">
            Procurement pricing, profit margin, cross-currency exchange conversion, and distributor lead times.
          </p>
        </div>
      </div>
      <button @click="openEditModal" class="btn btn-primary">
        <AppIcon name="edit" size="xs" class="mr-1.5 text-white" />
        <span>Edit Commercial Terms</span>
      </button>
    </div>

    <!-- Commercial Metrics Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- 1. Purchase Cost -->
      <div class="app-card p-4 bg-white border-l-4 border-l-purple-500 shadow-2xs space-y-1.5">
        <div class="flex items-center justify-between">
          <span class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Purchase Cost</span>
          <span class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-purple-50 text-purple-700 border border-purple-200">
            {{ detail.purchaseCurrency || detail.currency || 'USD' }}
          </span>
        </div>
        <div class="font-mono font-bold text-lg text-slate-900">
          {{ formatCurrency(detail.purchasePrice, detail.purchaseCurrency || detail.currency) }}
        </div>
        <div class="text-3xs font-mono text-slate-500">
          ≈ {{ formatCurrency(detail.purchasePriceInIDR || (detail.purchasePrice * (detail.effectiveExchangeRate || 16500)), 'IDR') }}
        </div>
      </div>

      <!-- 2. Selling Price -->
      <div class="app-card p-4 bg-white border-l-4 border-l-emerald-500 shadow-2xs space-y-1.5">
        <div class="flex items-center justify-between">
          <span class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Selling Price (MSRP)</span>
          <span class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-emerald-50 text-emerald-700 border border-emerald-200">
            {{ detail.sellingCurrency || detail.purchaseCurrency || detail.currency || 'USD' }}
          </span>
        </div>
        <div class="font-mono font-bold text-lg text-emerald-600">
          {{ formatCurrency(detail.sellingPrice, detail.sellingCurrency || detail.purchaseCurrency || detail.currency) }}
        </div>
        <div class="text-3xs font-mono text-emerald-700">
          ≈ {{ formatCurrency(detail.sellingPriceInIDR || (detail.sellingPrice * (detail.effectiveExchangeRate || 16500)), 'IDR') }}
        </div>
      </div>

      <!-- 3. Profit Margin -->
      <div class="app-card p-4 bg-white border-l-4 border-l-brand-500 shadow-2xs space-y-1.5">
        <div class="flex items-center justify-between">
          <span class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Profit Margin</span>
          <span class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-brand-50 text-brand-700 border border-brand-200">
            MARKUP
          </span>
        </div>
        <div class="font-mono font-bold text-lg text-brand-600">
          +{{ detail.profitMarginPercentage || marginPercentage }}%
        </div>
        <div class="text-3xs text-slate-400">
          Cost × (1 + {{ detail.profitMarginPercentage || marginPercentage }}%)
        </div>
      </div>

      <!-- 4. Exchange Rate Reference -->
      <div class="app-card p-4 bg-white border-l-4 border-l-indigo-500 shadow-2xs space-y-1.5">
        <div class="flex items-center justify-between">
          <span class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Exchange Rate</span>
          <span class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded" :class="detail.exchangeRateMode === 'MANUAL' ? 'bg-amber-50 text-amber-700 border border-amber-200' : 'bg-indigo-50 text-indigo-700 border border-indigo-200'">
            {{ detail.exchangeRateMode || 'AUTO' }}
          </span>
        </div>
        <div class="font-mono font-bold text-base text-indigo-950">
          1 USD = Rp {{ Number(detail.effectiveExchangeRate || 16500).toLocaleString('id-ID') }}
        </div>
        <div class="text-3xs text-slate-400 truncate">
          {{ detail.exchangeRateProvider || 'Frankfurter API' }}
        </div>
      </div>
    </div>

    <!-- Sourcing & Warranty Details (2 Columns) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Manufacturer & Sourcing Identifiers -->
      <div class="app-card p-5 bg-white space-y-4">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-100">
          <AppIcon name="factory" size="xs" class="text-purple-600" />
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Manufacturer & Part Numbers</h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
          <div>
            <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Manufacturer</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ detail.manufacturerName || 'Not specified' }}</div>
          </div>
          <div>
            <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Manufacturer Part Number (MPN)</div>
            <div class="font-mono font-bold text-purple-700 bg-purple-50 px-2 py-0.5 rounded border border-purple-200 mt-0.5 inline-block text-xs">
              {{ detail.manufacturerPartNumber || 'N/A' }}
            </div>
          </div>
          <div>
            <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Authorized Supplier / Distributor</div>
            <div class="font-bold text-slate-800 mt-0.5">{{ detail.supplierName || 'Not specified' }}</div>
          </div>
          <div>
            <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Supplier SKU / Order Code</div>
            <div class="font-mono text-slate-700 mt-0.5">{{ detail.supplierPartNumber || 'N/A' }}</div>
          </div>
        </div>
      </div>

      <!-- Warranty & Commercial Notes -->
      <div class="app-card p-5 bg-white space-y-4">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-100">
          <AppIcon name="shield-check" size="xs" class="text-emerald-600" />
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Warranty & Procurement Terms</h3>
        </div>

        <div class="space-y-3 text-xs">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Lead Time</div>
              <div class="text-slate-800 font-mono font-bold mt-0.5">{{ detail.leadTimeDays || 0 }} Days</div>
            </div>
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Warranty Coverage</div>
              <div class="text-slate-800 font-medium mt-0.5 flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                <span>{{ detail.warrantyInformation || 'Standard 12 Months' }}</span>
              </div>
            </div>
          </div>
          <div>
            <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Commercial Notes</div>
            <div class="text-slate-600 bg-slate-50 p-3 rounded-lg border border-slate-200 text-xs mt-1 leading-relaxed">
              {{ detail.notes || 'No commercial or procurement restrictions noted for this trading item.' }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Commercial Modal with Live Multi-Currency Engine -->
    <BaseModal
      :isOpen="modalOpen"
      title="Edit Commercial & Multi-Currency Pricing"
      @close="modalOpen = false"
    >
      <div class="space-y-5 text-xs">
        <!-- Sourcing Fields -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="form-label">Manufacturer Name</label>
            <input v-model="editForm.manufacturerName" type="text" class="form-input" placeholder="e.g. Vaisala Oyj" />
          </div>
          <div>
            <label class="form-label">Manufacturer Part Number (MPN)</label>
            <input v-model="editForm.manufacturerPartNumber" type="text" class="form-input font-mono" placeholder="e.g. WXT536-A" />
          </div>
          <div>
            <label class="form-label">Supplier / Distributor</label>
            <input v-model="editForm.supplierName" type="text" class="form-input" placeholder="e.g. DigiKey Electronics" />
          </div>
          <div>
            <label class="form-label">Supplier SKU / Part Number</label>
            <input v-model="editForm.supplierPartNumber" type="text" class="form-input font-mono" placeholder="e.g. DK-WXT536-ND" />
          </div>
        </div>

        <!-- PRICING & PROFIT ENGINE WORKSPACE -->
        <div class="p-4 rounded-xl bg-purple-50/50 border border-purple-200 space-y-4">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-purple-950 uppercase tracking-wider">Multi-Currency Pricing Calculator</span>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-purple-100 text-purple-800">LIVE ENGINE</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <!-- 1. Purchase Price -->
            <div class="p-3 bg-white rounded-lg border border-purple-200 space-y-2">
              <span class="text-3xs font-bold uppercase text-slate-500 block">Purchase Cost</span>
              <div class="grid grid-cols-5 gap-1.5">
                <input
                  v-model.number="editForm.purchasePrice"
                  @input="recalculateModalPricing"
                  type="number"
                  step="0.01"
                  min="0"
                  class="form-input col-span-3 font-mono font-bold text-xs"
                />
                <select
                  v-model="editForm.purchaseCurrency"
                  @change="onModalPurchaseCurrencyChange"
                  class="form-select col-span-2 font-mono text-3xs"
                >
                  <option value="USD">USD</option>
                  <option value="IDR">IDR</option>
                  <option value="EUR">EUR</option>
                  <option value="SGD">SGD</option>
                </select>
              </div>
            </div>

            <!-- 2. Profit Margin -->
            <div class="p-3 bg-white rounded-lg border border-purple-200 space-y-2">
              <span class="text-3xs font-bold uppercase text-slate-500 block">Profit Margin (%)</span>
              <div class="relative">
                <input
                  v-model.number="editForm.profitMarginPercentage"
                  @input="recalculateModalPricing"
                  type="number"
                  step="1"
                  min="0"
                  max="500"
                  class="form-input font-mono font-bold text-xs pr-6"
                />
                <span class="absolute right-2.5 top-2.5 text-xs text-slate-400 font-bold">%</span>
              </div>
            </div>

            <!-- 3. Selling Price -->
            <div class="p-3 bg-purple-100/60 rounded-lg border border-purple-300 space-y-2">
              <span class="text-3xs font-bold uppercase text-purple-900 block">Derived Selling Price</span>
              <div class="grid grid-cols-5 gap-1.5">
                <input
                  :value="editForm.sellingPrice"
                  readonly
                  type="number"
                  class="form-input col-span-3 font-mono font-bold text-xs bg-white text-purple-900 border-purple-300 cursor-not-allowed"
                />
                <select
                  v-model="editForm.sellingCurrency"
                  @change="recalculateModalPricing"
                  class="form-select col-span-2 font-mono text-3xs border-purple-300 bg-white"
                >
                  <option value="USD">USD</option>
                  <option value="IDR">IDR</option>
                  <option value="EUR">EUR</option>
                  <option value="SGD">SGD</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Exchange Rate & IDR Previews -->
          <div class="p-3.5 bg-white rounded-lg border border-purple-200 space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-3xs font-bold uppercase text-slate-700">Exchange Rate (Kurs 1 USD = Rp)</span>
              <span class="text-3xs font-mono font-bold text-indigo-700 truncate max-w-[120px]">
                {{ editForm.exchangeRateProvider || 'ExchangeRate-API' }}
              </span>
            </div>

            <div class="relative">
              <span class="absolute left-3 top-2.5 text-xs font-mono font-bold text-slate-400">Rp</span>
              <input
                v-model.number="editForm.manualExchangeRate"
                @input="onModalManualRateInput"
                type="number"
                step="1"
                min="1"
                class="form-input pl-9 font-mono font-bold text-xs bg-white text-slate-900 border-indigo-300 focus:border-brand-500 focus:ring-brand-500 shadow-inner"
                placeholder="17749"
              />
            </div>

            <div class="flex items-center justify-between text-3xs pt-1 border-t border-indigo-100">
              <button
                type="button"
                @click="useModalLiveRate"
                class="text-brand-700 hover:text-brand-900 font-bold flex items-center gap-1"
                title="Click to apply live market rate"
              >
                <AppIcon name="check" size="2xs" class="text-emerald-600" />
                <span>Sync Live: Rp {{ Number(editForm.liveExchangeRate || 17749).toLocaleString('id-ID') }}</span>
              </button>
              <span
                class="px-1.5 py-0.2 rounded font-mono text-3xs font-bold"
                :class="editForm.exchangeRateMode === 'MANUAL' ? 'bg-amber-50 text-amber-700 border border-amber-200' : 'bg-emerald-50 text-emerald-700 border border-emerald-200'"
              >
                {{ editForm.exchangeRateMode === 'MANUAL' ? 'Manual Input' : 'Auto Synced' }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 pt-2 border-t border-slate-100 text-3xs font-mono">
              <div>
                <span class="text-slate-400 block">IDR Purchase Cost:</span>
                <span class="font-bold text-slate-800">Rp {{ Number(editForm.purchasePriceInIDR || 0).toLocaleString('id-ID') }}</span>
              </div>
              <div>
                <span class="text-slate-400 block">IDR Selling Price:</span>
                <span class="font-bold text-emerald-700">Rp {{ Number(editForm.sellingPriceInIDR || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="form-label">Procurement Lead Time (Days)</label>
            <input v-model.number="editForm.leadTimeDays" type="number" class="form-input font-mono" placeholder="14" />
          </div>
          <div>
            <label class="form-label">Warranty Terms</label>
            <input v-model="editForm.warrantyInformation" type="text" class="form-input" placeholder="e.g. 24 Months Factory Limited Warranty" />
          </div>
        </div>

        <div>
          <label class="form-label">Procurement & Sourcing Notes</label>
          <textarea v-model="editForm.notes" rows="3" class="form-textarea" placeholder="Commercial notes, calibration requirements, minimum order quantities..."></textarea>
        </div>
      </div>

      <template #footer>
        <button @click="modalOpen = false" class="btn btn-secondary">Cancel</button>
        <button @click="saveDetails" class="btn btn-primary">Save Commercial Terms</button>
      </template>
    </BaseModal>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import { exchangeRateApi } from '@/api/exchangeRate';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'TradingCommercialTab',
  components: {
    BaseModal,
    AppIcon
  },
  props: {
    productId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      modalOpen: false,
      editForm: {
        manufacturerName: '',
        manufacturerPartNumber: '',
        supplierName: '',
        supplierPartNumber: '',
        purchasePrice: 0,
        purchaseCurrency: 'USD',
        profitMarginPercentage: 25,
        sellingPrice: 0,
        sellingCurrency: 'USD',
        exchangeRateMode: 'AUTO',
        liveExchangeRate: 16500,
        manualExchangeRate: 16500,
        effectiveExchangeRate: 16500,
        exchangeRateProvider: 'Frankfurter API',
        purchasePriceInIDR: 0,
        sellingPriceInIDR: 0,
        leadTimeDays: 0,
        warrantyInformation: '',
        notes: ''
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    detail() {
      return this.store.tradingDetails[this.productId] || {};
    },
    marginPercentage() {
      const p = Number(this.detail.purchasePrice || 0);
      const s = Number(this.detail.sellingPrice || 0);
      if (s <= 0) return 0;
      return Math.round(((s - p) / p) * 100);
    }
  },
  async mounted() {
    await this.store.fetchTradingDetail(this.productId);
  },
  methods: {
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
    openEditModal() {
      this.editForm = {
        manufacturerName: this.detail.manufacturerName || '',
        manufacturerPartNumber: this.detail.manufacturerPartNumber || '',
        supplierName: this.detail.supplierName || '',
        supplierPartNumber: this.detail.supplierPartNumber || '',
        purchasePrice: this.detail.purchasePrice || 0,
        purchaseCurrency: this.detail.purchaseCurrency || this.detail.currency || 'USD',
        profitMarginPercentage: this.detail.profitMarginPercentage || 25,
        sellingPrice: this.detail.sellingPrice || 0,
        sellingCurrency: this.detail.sellingCurrency || this.detail.purchaseCurrency || this.detail.currency || 'USD',
        exchangeRateMode: this.detail.exchangeRateMode || 'AUTO',
        liveExchangeRate: this.detail.liveExchangeRate || 16500,
        manualExchangeRate: this.detail.manualExchangeRate || this.detail.effectiveExchangeRate || 16500,
        effectiveExchangeRate: this.detail.effectiveExchangeRate || 16500,
        exchangeRateProvider: this.detail.exchangeRateProvider || 'Frankfurter API',
        purchasePriceInIDR: this.detail.purchasePriceInIDR || 0,
        sellingPriceInIDR: this.detail.sellingPriceInIDR || 0,
        leadTimeDays: this.detail.leadTimeDays || 0,
        warrantyInformation: this.detail.warrantyInformation || '',
        notes: this.detail.notes || ''
      };
      this.modalOpen = true;
      this.recalculateModalPricing();
    },
    onModalPurchaseCurrencyChange() {
      this.editForm.sellingCurrency = this.editForm.purchaseCurrency;
      this.recalculateModalPricing();
    },
    onModalManualRateInput() {
      this.editForm.exchangeRateMode = 'MANUAL';
      const rate = Number(this.editForm.manualExchangeRate) || 0;
      this.editForm.effectiveExchangeRate = rate > 0 ? rate : (this.editForm.liveExchangeRate || 17749);
      this.recalculateModalPricing();
    },
    useModalLiveRate() {
      this.editForm.exchangeRateMode = 'AUTO';
      const live = this.editForm.liveExchangeRate || 17749;
      this.editForm.manualExchangeRate = live;
      this.editForm.effectiveExchangeRate = live;
      this.recalculateModalPricing();
    },
    async recalculateModalPricing() {
      try {
        const res = await exchangeRateApi.calculatePricing({
          purchasePrice: Number(this.editForm.purchasePrice) || 0,
          purchaseCurrency: this.editForm.purchaseCurrency || 'USD',
          profitMarginPercentage: Number(this.editForm.profitMarginPercentage) || 0,
          sellingCurrency: this.editForm.sellingCurrency || 'USD',
          exchangeRateMode: this.editForm.exchangeRateMode || 'AUTO',
          manualExchangeRate: Number(this.editForm.manualExchangeRate) || 0
        });
        if (res && res.data) {
          this.editForm.sellingPrice = res.data.sellingPrice;
          this.editForm.liveExchangeRate = res.data.liveExchangeRate;
          this.editForm.effectiveExchangeRate = res.data.effectiveExchangeRate;
          this.editForm.exchangeRateProvider = res.data.exchangeRateProvider;
          this.editForm.purchasePriceInIDR = res.data.purchasePriceInIDR;
          this.editForm.sellingPriceInIDR = res.data.sellingPriceInIDR;
        }
      } catch (err) {
        const cost = Number(this.editForm.purchasePrice) || 0;
        const margin = Number(this.editForm.profitMarginPercentage) || 0;
        const baseSelling = cost * (1 + margin / 100);
        const effRate = this.editForm.exchangeRateMode === 'MANUAL' && this.editForm.manualExchangeRate > 0
          ? this.editForm.manualExchangeRate
          : (this.editForm.liveExchangeRate || 16500);

        this.editForm.effectiveExchangeRate = effRate;

        if (this.editForm.purchaseCurrency === this.editForm.sellingCurrency) {
          this.editForm.sellingPrice = Math.round(baseSelling * 100) / 100;
        } else if (this.editForm.purchaseCurrency === 'USD' && this.editForm.sellingCurrency === 'IDR') {
          this.editForm.sellingPrice = Math.round(baseSelling * effRate);
        } else if (this.editForm.purchaseCurrency === 'IDR' && this.editForm.sellingCurrency === 'USD') {
          this.editForm.sellingPrice = Math.round((baseSelling / effRate) * 100) / 100;
        }

        this.editForm.purchasePriceInIDR = this.editForm.purchaseCurrency === 'IDR' ? cost : Math.round(cost * effRate);
        this.editForm.sellingPriceInIDR = this.editForm.sellingCurrency === 'IDR' ? this.editForm.sellingPrice : Math.round(this.editForm.sellingPrice * effRate);
      }
    },
    async saveDetails() {
      try {
        await this.store.saveTradingDetail(this.productId, this.editForm);
        this.modalOpen = false;
        await this.store.fetchTradingDetail(this.productId);
        this.store.showToast('Commercial Terms Saved', 'Product pricing and multi-currency terms updated successfully.', 'success');
      } catch (err) {
        alert('Failed to save trading details: ' + (err.message || err));
      }
    }
  }
};
</script>
