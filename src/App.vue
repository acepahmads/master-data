<template>
  <div id="app" class="flex h-screen w-screen overflow-hidden bg-slate-100 font-sans text-slate-800 antialiased">
    <!-- Standalone View for Public / Login Routes -->
    <template v-if="isPublicRoute">
      <div class="flex-1 h-full w-full overflow-y-auto">
        <transition name="fade" mode="out-in">
          <router-view />
        </transition>
      </div>
    </template>

    <!-- Authenticated Engineering Workstation Shell -->
    <template v-else>
      <!-- LEFT SIDEBAR -->
      <AppSidebar />

      <!-- RIGHT MAIN WRAPPER -->
      <div class="flex-1 flex flex-col min-w-0 h-full overflow-hidden">
        <!-- TOP COMPACT HEADER -->
        <AppHeader />

        <!-- PAGE CONTENT AREA (Desktop 1440px Canvas) -->
        <main class="flex-1 overflow-y-auto p-4 md:p-6 lg:p-7 bg-slate-100/90">
          <div class="max-w-[1440px] mx-auto min-w-0">
            <transition name="fade" mode="out-in">
              <router-view />
            </transition>
          </div>
        </main>
      </div>

      <!-- GLOBAL MODALS & UTILITIES -->
      <GlobalSearchModal />
      <SettingsModal />
      <HelpModal />
    </template>

    <!-- Global Toast Alerts (Always Available) -->
    <ToastContainer />
  </div>
</template>

<script>
import AppSidebar from '@/components/common/AppSidebar.vue';
import AppHeader from '@/components/common/AppHeader.vue';
import GlobalSearchModal from '@/components/common/GlobalSearchModal.vue';
import SettingsModal from '@/components/common/SettingsModal.vue';
import HelpModal from '@/components/common/HelpModal.vue';
import ToastContainer from '@/components/common/ToastContainer.vue';

import { useMasterStore } from '@/stores/masterStore';
import { useAuthStore } from '@/stores/authStore';

export default {
  name: 'App',
  components: {
    AppSidebar,
    AppHeader,
    GlobalSearchModal,
    SettingsModal,
    HelpModal,
    ToastContainer
  },
  computed: {
    isPublicRoute() {
      return this.$route.matched.some(record => record.meta && record.meta.public) || this.$route.path === '/login';
    }
  },
  created() {
    const authStore = useAuthStore();
    const store = useMasterStore();

    // If a token is stored, validate session and load user
    if (authStore.token) {
      authStore.loadCurrentUser().then(user => {
        if (user) {
          store.currentUser = user;
          store.permissions = authStore.permissions;
          store.fetchInitialData();
        }
      });
    }

    window.addEventListener('keydown', this.handleGlobalKeydown);
  },
  beforeDestroy() {
    window.removeEventListener('keydown', this.handleGlobalKeydown);
  },
  methods: {
    handleGlobalKeydown(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const store = useMasterStore();
        store.globalSearchOpen = !store.globalSearchOpen;
      }
    }
  }
};
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease-in-out, transform 0.15s ease-in-out;
}
.fade-enter {
  opacity: 0;
  transform: translateY(4px);
}
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
