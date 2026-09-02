<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Page Header -->
    <PageHeader
      :title="isEdit ? `Edit Product Version v${form.versionNumber}` : 'Create Product Version'"
      :subtitle="isEdit ? `Update version parameters and release notes for ${form.productName}` : 'Register a new engineering version baseline and hardware evolution branch'"
    >
      <template #actions>
        <button @click="$router.back()" class="btn btn-secondary">
          Cancel
        </button>
        <button @click="submitForm" class="btn btn-primary">
          <AppIcon name="check" size="xs" class="mr-1.5 text-white" />
          <span>{{ isEdit ? 'Save Changes' : 'Create Version' }}</span>
        </button>
      </template>
    </PageHeader>

    <!-- Version Creation Form Card -->
    <div class="app-card p-6 bg-white space-y-6">
      <!-- Section 1: Product Selection & Creation Mode -->
      <div class="space-y-4 pb-5 border-b border-slate-100">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
          <AppIcon name="product" size="xs" class="text-brand-600" />
          <span>Product Master & Versioning Mode</span>
        </h3>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="form-label">Target IoT Product Master *</label>
            <select
              v-model="form.productCode"
              @change="onProductChange"
              :disabled="isEdit"
              class="form-select"
            >
              <option value="">Select Product...</option>
              <option
                v-for="p in masterStore.products"
                :key="p.code"
                :value="p.code"
              >
                {{ p.code }} — {{ p.name }}
              </option>
            </select>
          </div>

          <div v-if="!isEdit">
            <label class="form-label">Version Creation Mode</label>
            <div class="grid grid-cols-2 gap-2 mt-1">
              <button
                type="button"
                @click="form.creationMode = 'new'"
                :class="[
                  'p-2 rounded-lg border text-left text-xs transition-all',
                  form.creationMode === 'new'
                    ? 'border-brand-500 bg-brand-50/60 font-semibold text-brand-900 shadow-2xs'
                    : 'border-slate-200 text-slate-600 hover:bg-slate-50'
                ]"
              >
                <div class="font-bold">New Clean Baseline</div>
                <div class="text-3xs text-slate-500 font-normal">Start from fresh version</div>
              </button>

              <button
                type="button"
                @click="form.creationMode = 'inherit'"
                :class="[
                  'p-2 rounded-lg border text-left text-xs transition-all',
                  form.creationMode === 'inherit'
                    ? 'border-brand-500 bg-brand-50/60 font-semibold text-brand-900 shadow-2xs'
                    : 'border-slate-200 text-slate-600 hover:bg-slate-50'
                ]"
              >
                <div class="font-bold">Inherit Predecessor</div>
                <div class="text-3xs text-slate-500 font-normal">Branch from existing version</div>
              </button>
            </div>
          </div>
        </div>

        <!-- Inherit Source Version Selector -->
        <div v-if="!isEdit && form.creationMode === 'inherit'" class="p-3.5 rounded-xl bg-brand-50/40 border border-brand-200/80 text-xs space-y-2">
          <div class="flex items-center gap-1.5 text-brand-800 font-bold">
            <AppIcon name="git-branch" size="xs" />
            <span>Select Predecessor Baseline to Branch From</span>
          </div>
          <div class="max-w-md">
            <select v-model="form.baseVersionId" class="form-select font-mono">
              <option value="">Select source version...</option>
              <option
                v-for="v in availablePredecessors"
                :key="v.id"
                :value="v.id"
              >
                v{{ v.versionNumber }} — {{ v.versionName }} ({{ v.status }})
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Section 2: Version Identification & SemVer -->
      <div class="space-y-4 pb-5 border-b border-slate-100">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
          <AppIcon name="tag" size="xs" class="text-brand-600" />
          <span>Semantic Version Identification</span>
        </h3>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div>
            <div class="flex items-center justify-between">
              <label class="form-label">Version Number *</label>
              <span class="text-3xs font-mono text-slate-400">MAJOR.MINOR.PATCH</span>
            </div>
            <input
              v-model="form.versionNumber"
              type="text"
              placeholder="e.g. 2.2.0"
              class="form-input font-mono font-bold"
            />
          </div>

          <div class="sm:col-span-2">
            <label class="form-label">Version Name / Release Title *</label>
            <input
              v-model="form.versionName"
              type="text"
              placeholder="e.g. LTE-M / NB-IoT Cellular Edge Hub"
              class="form-input"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div>
            <label class="form-label">Lifecycle Status *</label>
            <select v-model="form.status" class="form-select font-mono font-semibold">
              <option value="Draft">Draft (In Planning)</option>
              <option value="In Review">In Review (Validation)</option>
              <option value="Active">Active (Production Baseline)</option>
              <option value="Deprecated">Deprecated</option>
              <option value="Archived">Archived</option>
            </select>
          </div>

          <div>
            <label class="form-label">Engineering Owner *</label>
            <select v-model="form.owner" class="form-select">
              <option v-for="u in masterStore.users" :key="u.id" :value="u.name">
                {{ u.name }} ({{ u.department }})
              </option>
            </select>
          </div>

          <div>
            <label class="form-label">Planned Release Date</label>
            <input
              v-model="form.releaseDate"
              type="date"
              class="form-input font-mono"
            />
          </div>
        </div>
      </div>

      <!-- Section 3: Engineering Description & Release Notes -->
      <div class="space-y-4">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
          <AppIcon name="file-text" size="xs" class="text-brand-600" />
          <span>Engineering Description & Release Notes</span>
        </h3>

        <div>
          <label class="form-label">Version Architecture Scope</label>
          <textarea
            v-model="form.description"
            rows="3"
            placeholder="Detailed engineering scope, MPU/MCU changes, interface additions, power architecture..."
            class="form-textarea"
          ></textarea>
        </div>

        <div>
          <label class="form-label">Release Notes / Changelog Points</label>
          <textarea
            v-model="form.releaseNotes"
            rows="4"
            placeholder="1. Upgraded main MCU to STM32MP157...&#10;2. Enhanced ESD protection arrays...&#10;3. Certified for CE and FCC Part 15..."
            class="form-textarea font-mono"
          ></textarea>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ProductVersionFormView',
  components: {
    PageHeader,
    AppIcon
  },
  data() {
    return {
      form: {
        id: '',
        productCode: '',
        productName: '',
        creationMode: 'new',
        baseVersionId: '',
        versionNumber: '',
        versionName: '',
        status: 'Draft',
        owner: 'Alex Chen',
        releaseDate: new Date().toISOString().slice(0, 10),
        description: '',
        releaseNotes: ''
      }
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    masterStore() {
      return useMasterStore();
    },
    versionId() {
      return this.$route.params.id;
    },
    isEdit() {
      return Boolean(this.versionId && this.versionId !== 'create');
    },
    availablePredecessors() {
      if (!this.form.productCode) return [];
      return this.versionStore.getVersionsForProduct(this.form.productCode);
    }
  },
  async mounted() {
    const productQuery = this.$route.query.product;
    if (this.masterStore.products.length === 0) {
      await this.masterStore.fetchProducts();
    }
    if (this.masterStore.users.length === 0) {
      await this.masterStore.fetchUsers();
    }
    await this.versionStore.fetchVersions();

    if (productQuery) {
      this.form.productCode = productQuery;
      this.onProductChange();
    }

    if (this.isEdit) {
      const existing = this.versionStore.getVersionById(this.versionId);
      if (existing) {
        this.form = { ...existing };
      }
    } else if (this.masterStore.products.length > 0 && !this.form.productCode) {
      this.form.productCode = this.masterStore.products[0].code;
      this.onProductChange();
    }
  },
  methods: {
    onProductChange() {
      const p = this.masterStore.products.find(item => item.code === this.form.productCode || item.id === this.form.productCode);
      if (p) {
        this.form.productName = p.name;
        this.form.productId = p.id;
        this.form.productCode = p.code;
      }
    },
    async submitForm() {
      if (!this.form.productCode || !this.form.versionNumber || !this.form.versionName) {
        alert('Please fill all required fields (Product, Version Number, Version Name).');
        return;
      }

      try {
        const saved = await this.versionStore.saveVersion({
          ...this.form,
          ownerAvatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'
        });

        this.$router.push(`/versions/${saved.id}`);
      } catch (err) {
        alert('Failed to save version: ' + (err.response?.data?.message || err.message));
      }
    }
  }
};
</script>
