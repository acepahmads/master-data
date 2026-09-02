<template>
  <header class="h-14 bg-white/95 backdrop-blur-md border-b border-slate-200/80 px-4 md:px-6 flex items-center justify-between shrink-0 sticky top-0 z-20 shadow-xs">
    <!-- LEFT: HAMBURGER (MOBILE) + BREADCRUMBS -->
    <div class="flex items-center gap-2.5 min-w-0">
      <!-- Mobile Sidebar Toggle Button -->
      <button
        @click="store.toggleMobileSidebar"
        class="lg:hidden p-1.5 -ml-1 rounded-lg text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-colors shrink-0"
        title="Toggle Navigation Menu"
      >
        <AppIcon name="menu" size="md" />
      </button>

      <nav class="flex items-center text-xs text-slate-500 font-medium truncate">
        <router-link to="/" class="hover:text-brand-600 transition-colors flex items-center gap-1 shrink-0">
          <AppIcon name="dashboard" size="xs" class="text-slate-400" />
          <span class="hidden sm:inline">Product & Master</span>
        </router-link>

        <template v-for="(crumb, idx) in breadcrumbs">
          <AppIcon name="chevron-right" size="xs" class="mx-1.5 text-slate-300 shrink-0" :key="`sep-${idx}`" />
          <router-link
            v-if="crumb.to"
            :to="crumb.to"
            :key="`link-${idx}`"
            class="hover:text-brand-600 transition-colors truncate"
          >
            {{ crumb.text }}
          </router-link>
          <span v-else :key="`text-${idx}`" class="text-slate-900 font-bold truncate">
            {{ crumb.text }}
          </span>
        </template>
      </nav>
    </div>

    <!-- RIGHT: GLOBAL SEARCH, NOTIFICATIONS & PROFILE -->
    <div class="flex items-center gap-3 shrink-0">
      <!-- Global Quick Search Trigger (Ctrl+K) -->
      <button
        @click="openGlobalSearch"
        class="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-50 hover:bg-slate-100/90 text-slate-400 hover:text-slate-700 border border-slate-200/80 hover:border-slate-300 transition-all duration-150 text-xs shadow-xs"
        title="Quick Search (Ctrl + K)"
      >
        <AppIcon name="search" size="xs" class="text-slate-400" />
        <span class="text-xs text-slate-500 font-medium">Search components, products...</span>
        <kbd class="hidden md:inline-block px-1.5 py-0.5 text-3xs font-mono font-bold bg-white text-slate-500 border border-slate-200 rounded shadow-2xs">
          Ctrl K
        </kbd>
      </button>

      <!-- Mobile Search Icon -->
      <button
        @click="openGlobalSearch"
        class="sm:hidden btn-icon"
        title="Search"
      >
        <AppIcon name="search" size="sm" />
      </button>

      <!-- Notifications Dropdown -->
      <div class="relative">
        <button
          @click="notificationsOpen = !notificationsOpen"
          class="relative btn-icon"
          title="Notifications"
        >
          <AppIcon name="bell" size="sm" />
          <span
            v-if="unreadCount > 0"
            class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-brand-500 ring-2 ring-white animate-pulse"
          ></span>
        </button>

        <!-- Notification Popover -->
        <div
          v-if="notificationsOpen"
          v-click-outside="closeNotifications"
          class="absolute right-0 mt-2 w-80 rounded-xl bg-white border border-slate-200 shadow-dropdown py-2 z-50 animate-in fade-in duration-150"
        >
          <div class="px-3.5 py-2 border-b border-slate-100 flex items-center justify-between">
            <span class="text-xs font-bold text-slate-800 uppercase tracking-wider">Activity Notifications</span>
            <button
              v-if="unreadCount > 0"
              @click="markAllRead"
              class="text-3xs text-brand-600 hover:text-brand-800 font-semibold"
            >
              Mark all read
            </button>
          </div>
          <div class="max-h-64 overflow-y-auto divide-y divide-slate-100">
            <div
              v-for="item in store.notifications"
              :key="item.id"
              :class="['p-3 flex items-start gap-2.5 hover:bg-slate-50 transition-colors', item.read ? 'opacity-70' : 'bg-brand-50/20']"
            >
              <div
                :class="[
                  'w-2 h-2 rounded-full mt-1.5 shrink-0',
                  item.type === 'success' ? 'bg-emerald-500' :
                  item.type === 'warning' ? 'bg-amber-500' : 'bg-brand-500'
                ]"
              ></div>
              <div class="min-w-0 flex-1">
                <div class="text-xs font-semibold text-slate-800">{{ item.title }}</div>
                <div class="text-2xs text-slate-500 leading-tight mt-0.5">{{ item.message }}</div>
                <div class="text-3xs text-slate-400 mt-1 font-mono">{{ item.time }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Role Switcher Dropdown (RBAC Simulator) -->
      <div v-if="store.currentUser" class="relative">
        <button
          @click="userMenuOpen = !userMenuOpen"
          class="flex items-center gap-2 p-1 pr-2.5 rounded-lg hover:bg-slate-100/80 transition-all border border-slate-200/60 hover:border-slate-300 focus:outline-none shadow-xs"
        >
          <div class="relative shrink-0">
            <img
              :src="store.currentUser.avatar || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
              alt="User"
              class="w-7 h-7 rounded-full object-cover ring-1 ring-slate-200"
            />
            <span class="absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full bg-emerald-500 ring-2 ring-white"></span>
          </div>
          <div class="hidden md:flex flex-col text-left">
            <span class="text-xs font-bold text-slate-800 leading-tight">{{ store.currentUser.name }}</span>
            <span class="text-3xs font-semibold text-brand-600">{{ userRoleName }}</span>
          </div>
          <AppIcon name="chevron-down" size="xs" class="text-slate-400 ml-0.5" />
        </button>

        <!-- User & Role Switch Popover -->
        <div
          v-if="userMenuOpen"
          v-click-outside="closeUserMenu"
          class="absolute right-0 mt-2 w-64 rounded-xl bg-white border border-slate-200 shadow-dropdown py-2 z-50"
        >
          <div class="px-3.5 py-2.5 border-b border-slate-100 bg-slate-50/50">
            <div class="text-xs font-bold text-slate-900">{{ store.currentUser.name }}</div>
            <div class="text-3xs text-slate-500 font-mono mt-0.5">{{ store.currentUser.email }}</div>
            <div class="mt-2 flex items-center gap-1.5">
              <span class="text-3xs font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                {{ userRoleName }}
              </span>
              <span class="text-3xs text-slate-500 font-medium">{{ store.currentUser.department }}</span>
            </div>
          </div>

          <!-- Switch Role Simulator -->
          <div class="px-3.5 py-1.5 text-3xs font-bold uppercase tracking-wider text-slate-400">
            Switch Role (RBAC Simulator)
          </div>
          <div class="max-h-48 overflow-y-auto">
            <button
              v-for="u in store.users"
              :key="u.id"
              @click="selectUser(u.id)"
              :class="[
                'w-full flex items-center justify-between px-3.5 py-2 text-xs text-left hover:bg-slate-50 transition-colors',
                store.currentUser.id === u.id ? 'bg-brand-50/70 font-semibold text-brand-700' : 'text-slate-700'
              ]"
            >
              <div class="flex items-center gap-2 min-w-0">
                <img :src="u.avatar" class="w-5 h-5 rounded-full object-cover" />
                <div class="truncate">
                  <span class="block leading-none font-medium">{{ u.name }}</span>
                  <span class="text-3xs text-slate-400 leading-tight">{{ getRoleName(u.role) }}</span>
                </div>
              </div>
              <AppIcon v-if="store.currentUser.id === u.id" name="check" size="xs" class="text-brand-600" />
            </button>
          </div>

          <div class="border-t border-slate-100 mt-1 pt-1 px-2 space-y-0.5">
            <button
              @click="openSettings"
              class="w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs font-medium text-slate-700 hover:bg-slate-100 transition-colors"
            >
              <AppIcon name="settings" size="xs" />
              <span>System Settings</span>
            </button>
            <button
              @click="handleLogout"
              class="w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs font-medium text-rose-600 hover:bg-rose-50 transition-colors"
            >
              <AppIcon name="logout" size="xs" class="text-rose-500" />
              <span>Sign Out</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import { useAuthStore } from '@/stores/authStore';
import AppIcon from './AppIcon.vue';

export default {
  name: 'AppHeader',
  components: { AppIcon },
  directives: {
    'click-outside': {
      bind(el, binding) {
        el.clickOutsideEvent = (event) => {
          if (!(el === event.target || el.contains(event.target))) {
            binding.value(event);
          }
        };
        document.body.addEventListener('click', el.clickOutsideEvent);
      },
      unbind(el) {
        document.body.removeEventListener('click', el.clickOutsideEvent);
      }
    }
  },
  data() {
    return {
      notificationsOpen: false,
      userMenuOpen: false
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    unreadCount() {
      return this.store.notifications.filter(n => !n.read).length;
    },
    breadcrumbs() {
      const path = this.$route.path;
      if (path === '/') return [{ text: 'Dashboard' }];
      
      const parts = path.split('/').filter(Boolean);
      const crumbs = [];

      if (parts[0] === 'products') {
        crumbs.push({ text: 'Products', to: '/products' });
        if (parts[1] === 'create') {
          crumbs.push({ text: 'New Product Master' });
        } else if (parts[2] === 'edit') {
          crumbs.push({ text: `Edit ${parts[1]}` });
        } else if (parts[1]) {
          crumbs.push({ text: parts[1] });
        }
      } else if (parts[0] === 'components') {
        crumbs.push({ text: 'Components', to: '/components' });
        if (parts[1] === 'create') {
          crumbs.push({ text: 'New Component' });
        } else if (parts[2] === 'edit') {
          crumbs.push({ text: `Edit ${parts[1]}` });
        } else if (parts[1]) {
          crumbs.push({ text: parts[1] });
        }
      } else if (parts[0] === 'categories') {
        crumbs.push({ text: 'Component Categories' });
      } else if (parts[0] === 'manufacturers') {
        crumbs.push({ text: 'Manufacturers' });
      } else if (parts[0] === 'suppliers') {
        crumbs.push({ text: 'Suppliers' });
      } else if (parts[0] === 'units') {
        crumbs.push({ text: 'Units of Measure (UoM)' });
      } else if (parts[0] === 'users') {
        crumbs.push({ text: 'User Management' });
      } else if (parts[0] === 'roles') {
        crumbs.push({ text: 'Roles & RBAC' });
      } else if (parts[0] === 'progress' || parts[0] === 'project-progress') {
        crumbs.push({ text: 'Platform Architecture', to: '/progress' });
        crumbs.push({ text: 'Roadmap & Progress' });
      } else if (parts[0] === 'versions') {
        crumbs.push({ text: 'Product Versions', to: '/versions' });
        if (parts[1] === 'create') {
          crumbs.push({ text: 'New Version Baseline' });
        } else if (parts[1] === 'compare') {
          crumbs.push({ text: 'Version Diff Comparison' });
        } else if (parts[2] === 'edit') {
          crumbs.push({ text: `Edit ${parts[1]}` });
        } else if (parts[1]) {
          crumbs.push({ text: parts[1] });
        }
      } else if (parts[0] === 'revisions') {
        crumbs.push({ text: 'Hardware Revisions', to: '/revisions' });
        if (parts[1]) {
          crumbs.push({ text: parts[1] });
        }
      }

      return crumbs;
    },
    userRoleName() {
      if (!this.store.currentUser) return 'Viewer';
      if (typeof this.store.currentUser.role === 'object' && this.store.currentUser.role !== null) {
        return this.store.currentUser.role.name || this.store.currentUser.roleName || 'Viewer';
      }
      return this.store.currentUser.role || this.store.currentUser.roleName || 'Viewer';
    }
  },
  methods: {
    getRoleName(role) {
      if (!role) return 'Viewer';
      if (typeof role === 'object' && role !== null) {
        return role.name || 'Viewer';
      }
      return role;
    },
    openGlobalSearch() {
      this.store.globalSearchOpen = true;
    },
    closeNotifications() {
      this.notificationsOpen = false;
    },
    closeUserMenu() {
      this.userMenuOpen = false;
    },
    markAllRead() {
      this.store.notifications.forEach(n => { n.read = true; });
    },
    selectUser(userId) {
      this.store.switchUser(userId);
      this.userMenuOpen = false;
    },
    openSettings() {
      this.store.settingsModalOpen = true;
      this.userMenuOpen = false;
    },
    handleLogout() {
      const authStore = useAuthStore();
      authStore.logout();
      this.store.logout();
      this.userMenuOpen = false;
      if (this.$route.path !== '/login') {
        this.$router.replace('/login');
      }
    }
  }
};
</script>

<style scoped>
.text-3xs {
  font-size: 0.625rem;
}
.shadow-2xs {
  box-shadow: 0 1px 1px 0 rgba(0, 0, 0, 0.05);
}
</style>
