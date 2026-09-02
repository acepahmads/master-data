<template>
  <div class="space-y-6">
    <!-- PAGE HEADER -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <span class="px-2.5 py-0.5 rounded-md text-3xs font-mono font-bold bg-brand-50 text-brand-700 border border-brand-200">
            DEVELOPER & INTEGRATIONS
          </span>
          <span class="text-xs text-slate-400">REST API v1</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 tracking-tight mt-1 flex items-center gap-2">
          <AppIcon name="key" size="sm" class="text-brand-600" />
          <span>API Access Tokens (Static Keys)</span>
        </h1>
        <p class="text-xs text-slate-500 mt-1 max-w-2xl">
          Generate static Bearer tokens for external applications (ERP, Mobile Apps, Scripts) to access the product catalog directly without a login flow.
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <button @click="openCreateModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Generate New API Token</span>
        </button>
      </div>
    </div>

    <!-- QUICK STATS BANNER -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="app-card p-4 bg-white flex items-center gap-4">
        <div class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shrink-0 border border-brand-100">
          <AppIcon name="key" size="sm" />
        </div>
        <div>
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Active Tokens</div>
          <div class="text-lg font-mono font-bold text-slate-900 mt-0.5">{{ activeTokensCount }}</div>
        </div>
      </div>

      <div class="app-card p-4 bg-white flex items-center gap-4">
        <div class="w-10 h-10 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center shrink-0 border border-purple-100">
          <AppIcon name="infinity" size="sm" />
        </div>
        <div>
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Unlimited Tokens</div>
          <div class="text-lg font-mono font-bold text-slate-900 mt-0.5">{{ unlimitedTokensCount }}</div>
        </div>
      </div>

      <div class="app-card p-4 bg-white flex items-center gap-4">
        <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0 border border-emerald-100">
          <AppIcon name="shield-check" size="sm" />
        </div>
        <div>
          <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Auth Header Format</div>
          <div class="text-xs font-mono font-semibold text-slate-800 mt-0.5 truncate">Bearer iot_live_...</div>
        </div>
      </div>
    </div>

    <!-- TOKENS LIST TABLE -->
    <div class="app-card bg-white overflow-hidden shadow-2xs border border-slate-200">
      <div class="p-4 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-slate-50/50">
        <div class="flex items-center gap-2">
          <AppIcon name="list" size="xs" class="text-slate-400" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Generated API Keys</h2>
          <span class="text-3xs font-mono px-2 py-0.5 rounded-full bg-slate-200 text-slate-700 font-bold">
            {{ tokens.length }}
          </span>
        </div>

        <button @click="loadTokens" class="btn btn-secondary btn-xs" :disabled="loading">
          <AppIcon name="refresh-cw" size="xs" :class="{ 'animate-spin': loading }" class="mr-1" />
          <span>Refresh</span>
        </button>
      </div>

      <div v-if="loading" class="p-12 text-center text-slate-400 text-xs">
        <AppIcon name="refresh-cw" size="md" class="animate-spin mx-auto mb-2 text-brand-600" />
        <span>Loading API tokens...</span>
      </div>

      <div v-else-if="tokens.length === 0" class="p-12 text-center">
        <div class="w-12 h-12 rounded-full bg-slate-100 text-slate-400 flex items-center justify-center mx-auto mb-3">
          <AppIcon name="key" size="md" />
        </div>
        <h3 class="text-sm font-bold text-slate-800">No API Tokens Generated Yet</h3>
        <p class="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
          Create static API tokens with configurable expiration (or Unlimited) to allow third-party applications to query products directly.
        </p>
        <button @click="openCreateModal" class="btn btn-primary mt-4">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Generate First API Token</span>
        </button>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b border-slate-200 bg-slate-50/90 text-3xs font-bold uppercase tracking-wider text-slate-500">
              <th class="py-3 px-4">Token Name & Label</th>
              <th class="py-3 px-4">Token Key Preview</th>
              <th class="py-3 px-4">Permission Scope</th>
              <th class="py-3 px-4">Expiration (Masa Berlaku)</th>
              <th class="py-3 px-4">Last Used</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="t in tokens" :key="t.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3 px-4">
                <div class="font-bold text-slate-900">{{ t.name }}</div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5">Created by {{ t.createdBy }} • {{ formatDate(t.createdAt) }}</div>
              </td>
              <td class="py-3 px-4 font-mono text-3xs">
                <span class="px-2 py-1 rounded bg-slate-100 text-slate-700 border border-slate-200 font-bold">
                  {{ t.tokenPreview }}
                </span>
              </td>
              <td class="py-3 px-4">
                <span
                  class="px-2 py-0.5 rounded-full text-3xs font-mono font-bold"
                  :class="t.scope === 'FULL_ACCESS' ? 'bg-purple-50 text-purple-700 border border-purple-200' : 'bg-blue-50 text-blue-700 border border-blue-200'"
                >
                  {{ t.scope === 'FULL_ACCESS' ? 'FULL ACCESS (R/W)' : 'READ ONLY (Catalog)' }}
                </span>
              </td>
              <td class="py-3 px-4">
                <div v-if="!t.expiresAt" class="flex items-center gap-1.5 text-emerald-600 font-semibold text-3xs">
                  <AppIcon name="infinity" size="xs" />
                  <span>Unlimited (Never Expire)</span>
                </div>
                <div v-else class="text-3xs font-mono" :class="isExpired(t.expiresAt) ? 'text-rose-600 font-bold' : 'text-slate-700'">
                  <span>{{ formatDate(t.expiresAt) }}</span>
                  <span v-if="isExpired(t.expiresAt)" class="ml-1 px-1.5 py-0.2 rounded bg-rose-50 border border-rose-200 text-rose-700 font-bold">EXPIRED</span>
                </div>
              </td>
              <td class="py-3 px-4 text-3xs text-slate-500 font-mono">
                {{ t.lastUsedAt ? formatDate(t.lastUsedAt) : 'Never used yet' }}
              </td>
              <td class="py-3 px-4">
                <span
                  class="px-2 py-0.5 rounded text-3xs font-mono font-bold"
                  :class="t.status === 'Active' && !isExpired(t.expiresAt) ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
                >
                  {{ t.status === 'Active' && !isExpired(t.expiresAt) ? 'ACTIVE' : 'REVOKED' }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <button
                  v-if="t.status === 'Active'"
                  @click="confirmRevoke(t)"
                  class="btn btn-danger btn-xs"
                  title="Revoke Token Access"
                >
                  <AppIcon name="trash" size="xs" class="mr-1" />
                  <span>Revoke</span>
                </button>
                <span v-else class="text-3xs text-slate-400 italic">Revoked</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- QUICK USAGE GUIDE CARD -->
    <div class="app-card p-5 bg-slate-900 text-white border border-slate-800 space-y-3">
      <div class="flex items-center gap-2">
        <AppIcon name="code" size="sm" class="text-brand-400" />
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-200">Quick Integration Example</h3>
      </div>
      <p class="text-xs text-slate-400">
        Directly query endpoints without calling <code>/api/v1/auth/login</code> by passing the token in the <code>Authorization</code> header:
      </p>
      <div class="p-3 bg-slate-950 rounded-lg font-mono text-3xs text-slate-300 overflow-x-auto border border-slate-800">
        <div class="text-slate-500"># cURL Example:</div>
        <div class="text-brand-300">curl -X GET "http://localhost:8080/api/v1/products?limit=50" \</div>
        <div class="text-brand-300 pl-4">-H "Authorization: Bearer iot_live_YOUR_TOKEN_HERE" \</div>
        <div class="text-brand-300 pl-4">-H "Content-Type: application/json"</div>
      </div>
    </div>

    <!-- MODAL 1: GENERATE TOKEN MODAL -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-md p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <AppIcon name="key" size="sm" class="text-brand-600" />
            <h3 class="text-sm font-bold text-slate-900">Generate New API Access Token</h3>
          </div>
          <button @click="showCreateModal = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <div class="space-y-4 text-xs">
          <!-- Token Name -->
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Token Name / Application Label *</label>
            <input
              v-model="form.name"
              type="text"
              class="input w-full"
              placeholder="e.g. Integrasi ERP SAP, Mobile App Sales, Script Python"
            />
            <p class="text-3xs text-slate-400 mt-1">A descriptive name to remember what application uses this token.</p>
          </div>

          <!-- Expiration -->
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Expiration Period (Masa Berlaku) *</label>
            <select v-model="form.expiresInDays" class="select w-full">
              <option :value="0">♾️ Unlimited (Never Expire / Tanpa Batas Waktu)</option>
              <option :value="7">⏱️ 7 Hari</option>
              <option :value="30">⏱️ 30 Hari (1 Bulan)</option>
              <option :value="90">⏱️ 90 Hari (3 Bulan)</option>
              <option :value="365">⏱️ 365 Hari (1 Tahun)</option>
            </select>
            <p class="text-3xs text-slate-400 mt-1">
              Select <strong>Unlimited</strong> if this token is used for continuous background server integration.
            </p>
          </div>

          <!-- Scope -->
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Permission Scope *</label>
            <select v-model="form.scope" class="select w-full">
              <option value="ALL_READ">🟢 Read-Only (Hanya Baca: Products, BOM, Components, Revisions)</option>
              <option value="FULL_ACCESS">🟡 Full Access (Read & Write)</option>
            </select>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
          <button @click="showCreateModal = false" class="btn btn-secondary text-xs">Cancel</button>
          <button @click="generateToken" class="btn btn-primary text-xs" :disabled="submitting || !form.name.trim()">
            <span v-if="submitting">Generating...</span>
            <span v-else>Generate Token</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 2: CREATED TOKEN SUCCESS MODAL (COPY ONCE) -->
    <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/70 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-lg p-6 space-y-4">
        <div class="flex items-center gap-3 text-emerald-600 border-b border-slate-100 pb-3">
          <div class="w-10 h-10 rounded-full bg-emerald-50 border border-emerald-200 flex items-center justify-center shrink-0">
            <AppIcon name="check" size="sm" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-slate-900">API Token Generated Successfully</h3>
            <p class="text-3xs text-slate-500 font-normal">Copy your token now. For security reasons, it will not be shown again.</p>
          </div>
        </div>

        <div class="space-y-3">
          <div>
            <div class="text-3xs font-semibold uppercase text-slate-400 mb-1">Your Static API Token</div>
            <div class="flex items-center gap-2">
              <input
                type="text"
                readonly
                :value="createdRawToken"
                class="input font-mono text-xs bg-slate-50 text-brand-700 font-bold select-all flex-1"
              />
              <button @click="copyToken" class="btn btn-primary shrink-0 text-xs">
                <AppIcon name="copy" size="xs" class="mr-1" />
                <span>{{ copied ? 'Copied!' : 'Copy' }}</span>
              </button>
            </div>
          </div>

          <div class="p-3 bg-amber-50 rounded-xl border border-amber-200 text-3xs text-amber-800 space-y-1">
            <div class="font-bold flex items-center gap-1">
              <AppIcon name="alert-triangle" size="xs" />
              <span>Important Security Notice</span>
            </div>
            <div>
              Store this token in a safe place (e.g. environment variable <code>API_BEARER_TOKEN</code> in your external app). If you lose it, you will need to generate a new token.
            </div>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-end">
          <button @click="showSuccessModal = false" class="btn btn-primary text-xs">
            <span>Done / I have copied it</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { apiClient } from '@/api/client';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ApiTokenView',
  components: { AppIcon },
  data() {
    return {
      tokens: [],
      loading: false,
      showCreateModal: false,
      showSuccessModal: false,
      submitting: false,
      copied: false,
      createdRawToken: '',
      form: {
        name: '',
        expiresInDays: 0,
        scope: 'ALL_READ'
      }
    };
  },
  computed: {
    activeTokensCount() {
      return this.tokens.filter(t => t.status === 'Active' && !this.isExpired(t.expiresAt)).length;
    },
    unlimitedTokensCount() {
      return this.tokens.filter(t => t.status === 'Active' && !t.expiresAt).length;
    }
  },
  mounted() {
    this.loadTokens();
  },
  methods: {
    async loadTokens() {
      this.loading = true;
      try {
        const res = await apiClient.get('/api-tokens');
        if (res.success) {
          this.tokens = res.data || [];
        }
      } catch (err) {
        console.error('Failed to load API tokens:', err);
      } finally {
        this.loading = false;
      }
    },
    openCreateModal() {
      this.form = {
        name: '',
        expiresInDays: 0,
        scope: 'ALL_READ'
      };
      this.showCreateModal = true;
    },
    async generateToken() {
      if (!this.form.name.trim()) return;
      this.submitting = true;
      try {
        const res = await apiClient.post('/api-tokens', this.form);
        if (res.success) {
          this.createdRawToken = res.data.rawToken;
          this.showCreateModal = false;
          this.showSuccessModal = true;
          this.copied = false;
          await this.loadTokens();
        }
      } catch (err) {
        alert(err.message || 'Failed to generate token');
      } finally {
        this.submitting = false;
      }
    },
    copyToken() {
      navigator.clipboard.writeText(this.createdRawToken);
      this.copied = true;
      setTimeout(() => {
        this.copied = false;
      }, 2500);
    },
    async confirmRevoke(tokenItem) {
      if (!confirm(`Are you sure you want to revoke API Token "${tokenItem.name}"? External apps using this token will lose access immediately.`)) {
        return;
      }
      try {
        await apiClient.delete(`/api-tokens/${tokenItem.id}`);
        await this.loadTokens();
      } catch (err) {
        alert(err.message || 'Failed to revoke token');
      }
    },
    isExpired(dateStr) {
      if (!dateStr) return false;
      return new Date(dateStr) < new Date();
    },
    formatDate(dateStr) {
      if (!dateStr) return '-';
      const d = new Date(dateStr);
      return d.toLocaleDateString('id-ID', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  }
};
</script>
