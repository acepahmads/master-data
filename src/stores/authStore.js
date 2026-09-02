import { defineStore } from 'pinia';
import authApi from '@/api/auth';

/**
 * Safely normalizes user objects and extracts a clean, displayable role name string.
 * Guarantees that raw JSON or permission objects are NEVER displayed in the UI.
 */
export function normalizeUser(user) {
  if (!user) return null;

  let roleObj = null;
  let roleName = 'Viewer';

  if (typeof user.role === 'object' && user.role !== null) {
    roleObj = user.role;
    roleName = user.role.name || 'Viewer';
  } else if (typeof user.role === 'string' && user.role.trim() !== '') {
    roleName = user.role;
    roleObj = { id: user.roleId || '', name: roleName, permissions: [] };
  } else if (user.roleName) {
    roleName = user.roleName;
  }

  // Ensure avatar has a fallback
  const avatar = user.avatar || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80';

  return {
    ...user,
    role: roleObj || user.role,
    roleName: roleName,
    avatar: avatar
  };
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('iot_token') || '',
    currentUser: null,
    permissions: [],
    isLoading: false,
    authInitialized: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.currentUser,
    
    userRoleName: (state) => {
      if (!state.currentUser) return 'Unknown Role';
      if (typeof state.currentUser.role === 'object' && state.currentUser.role !== null) {
        return state.currentUser.role.name || state.currentUser.roleName || 'Viewer';
      }
      return state.currentUser.role || state.currentUser.roleName || 'Viewer';
    },

    isAdmin: (state) => {
      if (!state.currentUser) return false;
      const role = state.currentUser.role;
      if (typeof role === 'object' && role !== null) {
        return role.name === 'Super Admin' || role.id === 'ROLE-ADMIN';
      }
      return role === 'Super Admin' || role === 'ROLE-ADMIN';
    },

    hasPermission: (state) => (permissionCode) => {
      if (!state.currentUser) return false;
      // Admin has blanket permissions
      if (state.currentUser.role && (state.currentUser.role.name === 'Super Admin' || state.currentUser.role === 'ROLE-ADMIN' || state.currentUser.role.id === 'ROLE-ADMIN')) {
        return true;
      }
      return state.permissions.includes(permissionCode);
    }
  },

  actions: {
    /**
     * Authenticate user with Email and Password
     * Calls POST /api/v1/auth/login
     */
    async login(email, password, remember = true) {
      this.isLoading = true;
      try {
        const res = await authApi.login(email, password);
        if (res && res.data && res.data.token) {
          this.token = res.data.token;
          if (remember) {
            localStorage.setItem('iot_token', this.token);
          } else {
            sessionStorage.setItem('iot_token', this.token);
            localStorage.setItem('iot_token', this.token);
          }

          this.currentUser = normalizeUser(res.data.user);
          this.permissions = res.data.permissions || [];
          this.authInitialized = true;
          return res.data;
        } else {
          throw new Error('Invalid response structure from authentication server');
        }
      } catch (err) {
        this.logout();
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    /**
     * Fetch current logged-in user details from JWT token
     * Calls GET /api/v1/auth/me
     */
    async loadCurrentUser() {
      if (!this.token) {
        this.currentUser = null;
        this.permissions = [];
        this.authInitialized = true;
        return null;
      }

      this.isLoading = true;
      try {
        const res = await authApi.getMe();
        if (res && res.data && res.data.user) {
          this.currentUser = normalizeUser(res.data.user);
          this.permissions = res.data.permissions || [];
          this.authInitialized = true;
          return this.currentUser;
        } else {
          this.logout();
          return null;
        }
      } catch (err) {
        console.warn('[AuthStore] Failed to load authenticated user profile:', err.message);
        this.logout();
        return null;
      } finally {
        this.isLoading = false;
        this.authInitialized = true;
      }
    },

    /**
     * Terminate user session
     */
    logout() {
      this.token = '';
      this.currentUser = null;
      this.permissions = [];
      this.authInitialized = true;
      localStorage.removeItem('iot_token');
      sessionStorage.removeItem('iot_token');
    }
  }
});
