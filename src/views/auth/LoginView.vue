<template>
  <div class="min-h-screen w-full flex items-center justify-center bg-slate-950 text-slate-100 relative overflow-hidden px-4 py-8">
    <!-- Ambient Background Lighting & Circuit Grids -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(14,165,233,0.15),rgba(255,255,255,0))] pointer-events-none"></div>
    <div class="absolute -top-40 -right-40 w-96 h-96 bg-brand-500/10 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Main Login Card Container -->
    <div class="w-full max-w-md relative z-10">
      <!-- Portal Brand Header -->
      <div class="text-center mb-6">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-cyan-950/80 border border-cyan-500/30 text-cyan-400 text-xs font-mono font-semibold mb-3 shadow-inner">
          <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
          <span>ENTERPRISE WORKSTATION ACCESS</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight flex flex-wrap items-center justify-center gap-2">
          <span>Product & Master Data</span>
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-brand-400">Control Center</span>
        </h1>
        <p class="text-xs text-slate-400 mt-1 font-medium">
          Centralized Product, Master Data & Engineering Lifecycle Platform
        </p>
      </div>

      <!-- Elevated White Card -->
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200/80 overflow-hidden text-slate-800">
        <div class="p-6 sm:p-8">
          <!-- Form Header -->
          <div class="mb-6">
            <h2 class="text-lg font-bold text-slate-900">Sign In to Workstation</h2>
            <p class="text-xs text-slate-500 mt-0.5">Enter your engineering credentials to access projects.</p>
          </div>

          <!-- Error Alert Banner -->
          <div
            v-if="errorMessage"
            class="mb-5 p-3.5 rounded-xl bg-rose-50 border border-rose-200/90 text-rose-800 text-xs flex items-start gap-2.5 animate-shake"
          >
            <AppIcon name="alert-circle" size="sm" class="text-rose-600 shrink-0 mt-0.5" />
            <div class="flex-1 min-w-0">
              <div class="font-bold text-rose-900">Authentication Failed</div>
              <div class="text-rose-700 mt-0.5">{{ errorMessage }}</div>
            </div>
            <button @click="errorMessage = ''" class="text-rose-400 hover:text-rose-600">
              <AppIcon name="x" size="xs" />
            </button>
          </div>

          <!-- Login Form -->
          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Email Field -->
            <div>
              <label for="email-input" class="block text-xs font-bold text-slate-700 mb-1.5">
                Workstation Email <span class="text-rose-500">*</span>
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                  <AppIcon name="user" size="xs" />
                </div>
                <input
                  id="email-input"
                  v-model.trim="email"
                  type="email"
                  required
                  placeholder="engineer@iotcontrol.io"
                  autocomplete="username"
                  class="w-full pl-9 pr-3.5 py-2.5 text-xs rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 text-slate-900 font-medium placeholder-slate-400 transition-all outline-none"
                  :disabled="isLoading"
                />
              </div>
            </div>

            <!-- Password Field -->
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <label for="password-input" class="block text-xs font-bold text-slate-700">
                  Password <span class="text-rose-500">*</span>
                </label>
              </div>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                  <AppIcon name="lock" size="xs" />
                </div>
                <input
                  id="password-input"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  placeholder="••••••••"
                  autocomplete="current-password"
                  class="w-full pl-9 pr-10 py-2.5 text-xs rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 text-slate-900 font-medium placeholder-slate-400 transition-all outline-none"
                  :disabled="isLoading"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none"
                  tabindex="-1"
                  title="Toggle password visibility"
                >
                  <AppIcon :name="showPassword ? 'eye' : 'lock'" size="xs" />
                </button>
              </div>
            </div>

            <!-- Remember Me Checkbox -->
            <div class="flex items-center justify-between pt-1">
              <label class="flex items-center gap-2 cursor-pointer select-none">
                <input
                  type="checkbox"
                  v-model="rememberMe"
                  class="w-3.5 h-3.5 rounded border-slate-300 text-brand-600 focus:ring-brand-500 focus:ring-offset-0"
                />
                <span class="text-xs text-slate-600 font-medium">Remember workstation session</span>
              </label>
            </div>

            <!-- Submit Button -->
            <div class="pt-2">
              <button
                type="submit"
                :disabled="isLoading || !email || !password"
                class="w-full py-2.5 px-4 rounded-xl text-xs font-bold text-white bg-gradient-to-r from-brand-600 to-cyan-600 hover:from-brand-500 hover:to-cyan-500 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md shadow-brand-500/20 flex items-center justify-center gap-2"
              >
                <svg
                  v-if="isLoading"
                  class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <AppIcon v-else name="login" size="xs" />
                <span>{{ isLoading ? 'Authenticating...' : 'Sign In to Workstation' }}</span>
              </button>
            </div>
          </form>

          <!-- Demo Quick Credentials Preset (Development / Auditor Helpers) -->
          <div class="mt-6 pt-5 border-t border-slate-100">
            <div class="text-3xs font-bold uppercase tracking-wider text-slate-400 mb-2.5 flex items-center justify-between">
              <span>Quick Demo Credentials</span>
              <span class="text-slate-400 font-mono font-normal">Password: admin123</span>
            </div>
            <div class="grid grid-cols-2 gap-1.5">
              <button
                v-for="acc in demoAccounts"
                :key="acc.email"
                type="button"
                @click="prefill(acc.email, 'admin123')"
                class="p-2 rounded-lg bg-slate-50 hover:bg-brand-50/80 border border-slate-200/80 hover:border-brand-200 text-left transition-all group"
              >
                <div class="text-xs font-bold text-slate-800 group-hover:text-brand-700 truncate leading-tight">{{ acc.name }}</div>
                <div class="text-3xs text-slate-500 group-hover:text-brand-600 truncate mt-0.5 flex items-center gap-1">
                  <span class="font-semibold">{{ acc.role }}</span>
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Card Footer -->
        <div class="px-6 py-3 bg-slate-50 border-t border-slate-100 flex items-center justify-between text-3xs text-slate-500">
          <div class="flex items-center gap-1.5 font-medium">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
            <span>REST API v1.0.0 Online</span>
          </div>
          <span class="font-mono">AES-256 JWT RBAC</span>
        </div>
      </div>

      <!-- Security Notice -->
      <div class="text-center mt-6 text-3xs text-slate-500 font-medium">
        Internal R&D Control Center • Authorized Engineering Personnel Only
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '@/stores/authStore';
import { useMasterStore } from '@/stores/masterStore';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'LoginView',
  components: {
    AppIcon
  },
  data() {
    return {
      email: 'alex.chen@iotcontrol.io',
      password: 'admin123',
      showPassword: false,
      rememberMe: true,
      isLoading: false,
      errorMessage: '',
      demoAccounts: [
        { name: 'Alex Chen', email: 'alex.chen@iotcontrol.io', role: 'Super Admin' },
        { name: 'Marcus Vance', email: 'marcus.vance@iotcontrol.io', role: 'Hardware Lead' },
        { name: 'Sarah Jenkins', email: 'sarah.jenkins@iotcontrol.io', role: 'R&D Manager' },
        { name: 'Maya Lin', email: 'maya.lin@iotcontrol.io', role: 'Viewer' }
      ]
    };
  },
  methods: {
    prefill(email, password) {
      this.email = email;
      this.password = password;
      this.errorMessage = '';
    },

    async handleLogin() {
      if (!this.email || !this.password) {
        this.errorMessage = 'Please enter both email and password.';
        return;
      }

      this.isLoading = true;
      this.errorMessage = '';

      try {
        const authStore = useAuthStore();
        const masterStore = useMasterStore();

        // Perform real backend authentication
        const loginData = await authStore.login(this.email, this.password, this.rememberMe);

        // Synchronize masterStore session
        masterStore.token = loginData.token;
        masterStore.currentUser = authStore.currentUser;
        masterStore.permissions = authStore.permissions;

        // Fetch master collections in background
        masterStore.fetchInitialData();

        // Redirect to requested path or dashboard
        const redirectPath = this.$route.query.redirect || '/';
        this.$router.replace(redirectPath);
      } catch (err) {
        console.error('[LoginView] Login error:', err);
        this.errorMessage = err.message || 'Invalid email or password. Please verify your credentials.';
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}
.animate-shake {
  animation: shake 0.3s ease-in-out;
}
</style>
