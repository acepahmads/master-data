import Vue from 'vue';
import { defineStore } from 'pinia';
import authApi from '@/api/auth';
import productsApi from '@/api/products';
import componentsApi from '@/api/components';
import categoriesApi from '@/api/categories';
import manufacturersApi from '@/api/manufacturers';
import suppliersApi from '@/api/suppliers';
import customersApi from '@/api/customers';
import unitsApi from '@/api/units';
import usersApi from '@/api/users';
import rolesApi from '@/api/roles';
import activityApi from '@/api/activity';
import uploadApi from '@/api/upload';

// Helper to safely normalize user objects and role display names
export function normalizeUser(user) {
  if (!user) return null;
  let roleObj = null;
  let roleName = 'Unknown Role';

  if (typeof user.role === 'object' && user.role !== null) {
    roleObj = user.role;
    roleName = user.role.name || 'Unknown Role';
  } else if (typeof user.role === 'string' && user.role.trim() !== '') {
    roleName = user.role;
    roleObj = { id: user.roleId || '', name: roleName, permissions: [] };
  } else if (user.roleName) {
    roleName = user.roleName;
  }

  return {
    ...user,
    role: roleObj || user.role,
    roleName: roleName
  };
}

export const useMasterStore = defineStore('master', {
  state: () => ({
    token: localStorage.getItem('iot_token') || '',
    currentUser: null,
    permissions: [],
    isLoading: false,

    // Notifications
    notifications: [
      { id: 'notif-1', title: 'Component Status Changed', message: 'ESP32-S3-WROOM-1 approved for production', time: '10 min ago', read: false, type: 'success' },
      { id: 'notif-2', title: 'Manufacturer Audit Due', message: 'Bosch Sensortec annual audit verification pending', time: '2 hours ago', read: false, type: 'warning' },
      { id: 'notif-3', title: 'New Product Registered', message: 'Smart Edge Gateway (GW-400X) created by Sarah J.', time: '1 day ago', read: true, type: 'info' }
    ],

    // UI state
    globalSearchOpen: false,
    settingsModalOpen: false,
    helpModalOpen: false,
    sidebarCollapsed: false,
    mobileSidebarOpen: false,
    toast: {
      show: false,
      message: '',
      type: 'success',
      title: ''
    },

    // Master Data collections from REST API
    products: [],
    components: [],
    categories: [],
    flatCategoryList: [],
    manufacturers: [],
    suppliers: [],
    customers: [],
    units: [],
    users: [],
    roles: [],
    allPermissions: [],
    activities: [],
    tradingDetails: {},
    projectItems: {}
  }),

  getters: {
    totalProductsCount: (state) => state.products.length,
    totalComponentsCount: (state) => state.components.length,
    totalCategoriesCount: (state) => {
      let count = 0;
      const countNodes = (nodes) => {
        if (!nodes) return;
        nodes.forEach(node => {
          count++;
          if (node.children && node.children.length > 0) {
            countNodes(node.children);
          }
        });
      };
      countNodes(state.categories);
      return count;
    },
    totalManufacturersCount: (state) => state.manufacturers.length,
    totalSuppliersCount: (state) => state.suppliers.length,
    totalCustomersCount: (state) => state.customers.length,
    totalUnitsCount: (state) => state.units.length,
    totalUsersCount: (state) => state.users.length,

    // Clean human-readable role name for UI display
    currentUserRoleName: (state) => {
      if (!state.currentUser) return 'Viewer';
      if (typeof state.currentUser.role === 'object' && state.currentUser.role !== null) {
        return state.currentUser.role.name || state.currentUser.roleName || 'Viewer';
      }
      return state.currentUser.role || state.currentUser.roleName || 'Viewer';
    },

    // Check if current user is Super Admin
    isSuperAdmin: (state) => {
      if (!state.currentUser) return false;
      const role = state.currentUser.role;
      if (typeof role === 'object' && role !== null) {
        return role.name === 'Super Admin' || role.id === 'ROLE-ADMIN';
      }
      return role === 'Super Admin' || role === 'ROLE-ADMIN' || state.currentUser.roleName === 'Super Admin';
    },

    // Flat category list for select dropdowns
    flatCategories: (state) => {
      const list = [];
      const traverse = (nodes, depth = 0) => {
        if (!nodes) return;
        nodes.forEach(node => {
          list.push({
            id: node.id,
            name: `${'— '.repeat(depth)}${node.name}`,
            rawName: node.name,
            code: node.code,
            depth: depth,
            parentId: node.parentId
          });
          if (node.children && node.children.length > 0) {
            traverse(node.children, depth + 1);
          }
        });
      };
      traverse(state.categories);
      return list;
    }
  },

  actions: {
    // Show Toast
    showToast(title, message, type = 'success') {
      this.toast = {
        show: true,
        title,
        message,
        type
      };
      setTimeout(() => {
        this.toast.show = false;
      }, 4000);
    },

    // INITIALIZE AUTH & MASTER DATA
    async initAuth() {
      try {
        if (this.token) {
          const res = await authApi.getMe();
          if (res && res.data && res.data.user) {
            this.currentUser = normalizeUser(res.data.user);
            this.permissions = res.data.permissions || [];
            return this.currentUser;
          }
        }
        this.currentUser = null;
        this.permissions = [];
        this.token = '';
        localStorage.removeItem('iot_token');
        return null;
      } catch (err) {
        console.warn('[MasterStore] Auth token validation error:', err.message);
        this.currentUser = null;
        this.permissions = [];
        this.token = '';
        localStorage.removeItem('iot_token');
        return null;
      }
    },

    async fetchInitialData() {
      this.isLoading = true;
      const user = await this.initAuth();
      if (!user) {
        this.isLoading = false;
        return;
      }
      try {
        await Promise.allSettled([
          this.fetchProducts(),
          this.fetchComponents(),
          this.fetchCategories(),
          this.fetchManufacturers(),
          this.fetchSuppliers(),
          this.fetchUnits(),
          this.fetchUsers(),
          this.fetchRoles(),
          this.fetchActivities()
        ]);
      } catch (err) {
        console.error('[MasterStore] Initial data load error:', err);
      } finally {
        this.isLoading = false;
      }
    },

    async login(email, password) {
      try {
        const res = await authApi.login(email, password);
        if (res && res.data && res.data.token) {
          this.token = res.data.token;
          localStorage.setItem('iot_token', this.token);
          this.currentUser = normalizeUser(res.data.user);
          this.permissions = res.data.permissions || [];
          this.showToast('Login Successful', `Welcome back, ${this.currentUser.name}!`);
          await this.fetchInitialData();
          return res.data;
        } else {
          throw new Error('Invalid response structure');
        }
      } catch (err) {
        this.showToast('Login Failed', err.message, 'error');
        throw err;
      }
    },

    logout() {
      this.token = '';
      this.currentUser = null;
      this.permissions = [];
      localStorage.removeItem('iot_token');
      sessionStorage.removeItem('iot_token');
    },

    async switchUser(userId) {
      const user = this.users.find(u => u.id === userId);
      if (user) {
        this.currentUser = normalizeUser(user);
        const roleLabel = this.currentUser.roleName || (typeof user.role === 'object' ? user.role.name : user.role) || user.department;
        this.showToast('User Switched', `Active session preview set to ${user.name} (${roleLabel})`, 'info');
      }
    },

    // PRODUCTS CRUD
    async fetchProducts(params = {}) {
      try {
        const res = await productsApi.getAll({ limit: 100, ...params });
        if (res && res.data) {
          this.products = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchProducts error:', err);
      }
    },

    async saveProduct(productData) {
      try {
        const isEdit = Boolean(productData.id && this.products.some(p => p.id === productData.id || p.code === productData.code));
        let res;
        if (isEdit) {
          res = await productsApi.update(productData.id || productData.code, productData);
          this.showToast('Product Updated', `${productData.name} was updated successfully.`);
        } else {
          res = await productsApi.create(productData);
          this.showToast('Product Created', `${productData.name} was created successfully.`);
        }
        await this.fetchProducts();
        await this.fetchActivities();
        return res && res.data ? res.data : productData;
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteProduct(id) {
      try {
        const p = this.products.find(item => item.id === id || item.code === id);
        const name = p ? p.name : id;
        await productsApi.delete(id);
        this.products = this.products.filter(item => item.id !== id && item.code !== id);
        this.showToast('Product Deleted', `${name} has been removed.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteMultipleProducts(ids) {
      if (!ids || ids.length === 0) return;
      this.isLoading = true;
      let count = 0;
      try {
        for (const id of ids) {
          try {
            await productsApi.delete(id);
            count++;
          } catch (e) {
            console.warn(`Failed to delete product ${id}:`, e);
          }
        }
        const idSet = new Set(ids);
        this.products = this.products.filter(item => !idSet.has(item.id) && !idSet.has(item.code));
        this.showToast('Batch Delete', `Successfully deleted ${count} products.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      } finally {
        this.isLoading = false;
      }
    },

    async deleteAllProducts() {
      this.isLoading = true;
      let count = 0;
      try {
        const allIds = this.products.map(p => p.id || p.code);
        for (const id of allIds) {
          try {
            await productsApi.delete(id);
            count++;
          } catch (e) {
            console.warn(`Failed to delete product ${id}:`, e);
          }
        }
        this.products = [];
        this.showToast('Purge Completed', `All ${count} products have been deleted by Super Admin.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete All Failed', err.message, 'error');
      } finally {
        this.isLoading = false;
      }
    },

    // TRADING PRODUCT DETAILS (Phase 2C)
    async fetchTradingDetail(productId) {
      try {
        const res = await productsApi.getTradingDetail(productId);
        if (res && res.data) {
          Vue.set(this.tradingDetails, productId, res.data);
          return res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchTradingDetail error:', err);
      }
      return null;
    },

    async saveTradingDetail(productId, detailData) {
      try {
        const res = await productsApi.updateTradingDetail(productId, detailData);
        if (res && res.data) {
          Vue.set(this.tradingDetails, productId, res.data);
          this.showToast('Trading Detail Saved', 'Commercial sourcing data updated.');
          await this.fetchActivities();
          return res.data;
        }
      } catch (err) {
        this.showToast('Save Failed', err.response?.data?.message || err.message, 'error');
        throw err;
      }
    },

    // PROJECT PRODUCT ITEMS & STRUCTURE (Phase 2C)
    async fetchProjectItems(projectId) {
      try {
        const res = await productsApi.getProjectItems(projectId);
        if (res && res.data) {
          Vue.set(this.projectItems, projectId, res.data);
          return res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchProjectItems error:', err);
      }
      return [];
    },

    async createProjectItem(projectId, itemData) {
      try {
        const res = await productsApi.createProjectItem(projectId, itemData);
        await this.fetchProjectItems(projectId);
        await this.fetchActivities();
        this.showToast('Item Added', `${itemData.name} added to project structure.`);
        return res && res.data ? res.data : null;
      } catch (err) {
        this.showToast('Add Failed', err.response?.data?.message || err.message, 'error');
        throw err;
      }
    },

    async updateProjectItem(projectId, itemId, itemData) {
      try {
        const res = await productsApi.updateProjectItem(itemId, itemData);
        await this.fetchProjectItems(projectId);
        await this.fetchActivities();
        this.showToast('Item Updated', 'Project structure item updated.');
        return res && res.data ? res.data : null;
      } catch (err) {
        this.showToast('Update Failed', err.response?.data?.message || err.message, 'error');
        throw err;
      }
    },

    async deleteProjectItem(projectId, itemId) {
      try {
        await productsApi.deleteProjectItem(itemId);
        await this.fetchProjectItems(projectId);
        await this.fetchActivities();
        this.showToast('Item Removed', 'Item removed from project structure.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.response?.data?.message || err.message, 'error');
        throw err;
      }
    },

    async reorderProjectItems(projectId, itemIds) {
      try {
        await productsApi.reorderProjectItems(projectId, itemIds);
        await this.fetchProjectItems(projectId);
      } catch (err) {
        console.error('[MasterStore] reorderProjectItems error:', err);
      }
    },

    // COMPONENTS CRUD
    async fetchComponents(params = {}) {
      try {
        const res = await componentsApi.getAll({ limit: 100, ...params });
        if (res && res.data) {
          this.components = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchComponents error:', err);
      }
    },

    async saveComponent(componentData) {
      try {
        const isEdit = Boolean(componentData.id && this.components.some(c => c.id === componentData.id || c.partNumber === componentData.partNumber));
        let res;
        if (isEdit) {
          res = await componentsApi.update(componentData.id || componentData.partNumber, componentData);
          this.showToast('Component Updated', `${componentData.partNumber} was updated successfully.`);
        } else {
          res = await componentsApi.create(componentData);
          this.showToast('Component Created', `${componentData.partNumber} was registered successfully.`);
        }
        await this.fetchComponents();
        await this.fetchActivities();
        return res && res.data ? res.data : componentData;
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteComponent(id) {
      try {
        const c = this.components.find(item => item.id === id || item.partNumber === id);
        const partNumber = c ? c.partNumber : id;
        await componentsApi.delete(id);
        this.components = this.components.filter(item => item.id !== id && item.partNumber !== id);
        this.showToast('Component Deleted', `${partNumber} has been removed.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteMultipleComponents(ids) {
      if (!ids || ids.length === 0) return;
      this.isLoading = true;
      let count = 0;
      try {
        for (const id of ids) {
          try {
            await componentsApi.delete(id);
            count++;
          } catch (e) {
            console.warn(`Failed to delete component ${id}:`, e);
          }
        }
        const idSet = new Set(ids);
        this.components = this.components.filter(item => !idSet.has(item.id) && !idSet.has(item.partNumber));
        this.showToast('Batch Delete', `Successfully deleted ${count} components.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      } finally {
        this.isLoading = false;
      }
    },

    async deleteAllComponents() {
      this.isLoading = true;
      let count = 0;
      try {
        const allIds = this.components.map(c => c.id || c.partNumber);
        for (const id of allIds) {
          try {
            await componentsApi.delete(id);
            count++;
          } catch (e) {
            console.warn(`Failed to delete component ${id}:`, e);
          }
        }
        this.components = [];
        this.showToast('Purge Completed', `All ${count} components have been deleted by Super Admin.`, 'warning');
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Delete All Failed', err.message, 'error');
      } finally {
        this.isLoading = false;
      }
    },

    async addComponentDocument(componentId, docData) {
      try {
        const res = await componentsApi.addDocument(componentId, docData);
        await this.fetchComponents();
        await this.fetchActivities();
        this.showToast('Document Attached', `${docData.name} has been attached.`, 'info');
        return res && res.data ? res.data : docData;
      } catch (err) {
        this.showToast('Upload Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteComponentDocument(componentId, docId) {
      try {
        await componentsApi.deleteDocument(componentId, docId);
        await this.fetchComponents();
        this.showToast('Document Removed', 'Document deleted.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // CATEGORIES CRUD
    async fetchCategories() {
      try {
        const [treeRes, flatRes] = await Promise.all([
          categoriesApi.getTree(),
          categoriesApi.getFlat()
        ]);
        if (treeRes && treeRes.data) {
          this.categories = treeRes.data;
        }
        if (flatRes && flatRes.data) {
          this.flatCategoryList = flatRes.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchCategories error:', err);
      }
    },

    async saveCategory(categoryData) {
      try {
        const isEdit = Boolean(categoryData.id && this.flatCategoryList.some(c => c.id === categoryData.id));
        if (isEdit) {
          await categoriesApi.update(categoryData.id, categoryData);
          this.showToast('Category Updated', `${categoryData.name} updated.`);
        } else {
          await categoriesApi.create(categoryData);
          this.showToast('Category Created', `${categoryData.name} category created.`);
        }
        await this.fetchCategories();
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteCategory(id) {
      try {
        await categoriesApi.delete(id);
        this.showToast('Category Removed', 'Category was deleted.', 'warning');
        await this.fetchCategories();
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
        throw err;
      }
    },

    // MANUFACTURERS CRUD
    async fetchManufacturers(params = {}) {
      try {
        const res = await manufacturersApi.getAll({ limit: 100, ...params });
        if (res && res.data) {
          this.manufacturers = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchManufacturers error:', err);
      }
    },

    async saveManufacturer(mfgData) {
      try {
        const isEdit = Boolean(mfgData.id && this.manufacturers.some(m => m.id === mfgData.id));
        if (isEdit) {
          await manufacturersApi.update(mfgData.id, mfgData);
          this.showToast('Manufacturer Updated', `${mfgData.name} updated.`);
        } else {
          await manufacturersApi.create(mfgData);
          this.showToast('Manufacturer Added', `${mfgData.name} created.`);
        }
        await this.fetchManufacturers();
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteManufacturer(id) {
      try {
        await manufacturersApi.delete(id);
        this.manufacturers = this.manufacturers.filter(item => item.id !== id);
        this.showToast('Manufacturer Deleted', 'Manufacturer removed.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // SUPPLIERS CRUD
    async fetchSuppliers(params = {}) {
      try {
        const res = await suppliersApi.getAll({ limit: 100, ...params });
        if (res && res.data) {
          this.suppliers = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchSuppliers error:', err);
      }
    },

    async saveSupplier(supplierData) {
      try {
        const isEdit = Boolean(supplierData.id && this.suppliers.some(s => s.id === supplierData.id));
        if (isEdit) {
          await suppliersApi.update(supplierData.id, supplierData);
          this.showToast('Supplier Updated', `${supplierData.name} updated.`);
        } else {
          await suppliersApi.create(supplierData);
          this.showToast('Supplier Added', `${supplierData.name} created.`);
        }
        await this.fetchSuppliers();
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteSupplier(id) {
      try {
        await suppliersApi.delete(id);
        this.suppliers = this.suppliers.filter(item => item.id !== id);
        this.showToast('Supplier Deleted', 'Supplier removed.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // CUSTOMERS CRUD
    async fetchCustomers(params = {}) {
      try {
        const res = await customersApi.getAll({ limit: 100, ...params });
        if (res && res.data) {
          this.customers = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchCustomers error:', err);
      }
    },

    async saveCustomer(customerData) {
      try {
        const isEdit = Boolean(customerData.id && this.customers.some(c => c.id === customerData.id));
        if (isEdit) {
          await customersApi.update(customerData.id, customerData);
          this.showToast('Customer Updated', `${customerData.name || customerData.company} updated.`);
        } else {
          await customersApi.create(customerData);
          this.showToast('Customer Added', `${customerData.name || customerData.company} created.`);
        }
        await this.fetchCustomers();
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteCustomer(id) {
      try {
        await customersApi.delete(id);
        this.customers = this.customers.filter(item => item.id !== id);
        this.showToast('Customer Deleted', 'Customer removed.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // UNITS OF MEASURE CRUD
    async fetchUnits() {
      try {
        const res = await unitsApi.getAll({ limit: 100 });
        if (res && res.data) {
          this.units = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchUnits error:', err);
      }
    },

    async saveUnit(unitData) {
      try {
        const isEdit = Boolean(unitData.id && this.units.some(u => u.id === unitData.id));
        if (isEdit) {
          await unitsApi.update(unitData.id, unitData);
          this.showToast('Unit Updated', `${unitData.name} updated.`);
        } else {
          await unitsApi.create(unitData);
          this.showToast('Unit Created', `${unitData.name} (${unitData.symbol}) created.`);
        }
        await this.fetchUnits();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async deleteUnit(id) {
      try {
        await unitsApi.delete(id);
        this.units = this.units.filter(item => item.id !== id);
        this.showToast('Unit Deleted', 'Unit removed.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // USERS CRUD
    async fetchUsers() {
      try {
        const res = await usersApi.getAll({ limit: 100 });
        if (res && res.data) {
          this.users = res.data.map(normalizeUser);
        }
      } catch (err) {
        console.error('[MasterStore] fetchUsers error:', err);
      }
    },

    async saveUser(userData) {
      try {
        const isEdit = Boolean(userData.id && this.users.some(u => u.id === userData.id));
        if (isEdit) {
          await usersApi.update(userData.id, userData);
          this.showToast('User Updated', `User ${userData.name} updated.`);
        } else {
          await usersApi.create(userData);
          this.showToast('User Created', `User ${userData.name} created.`);
        }
        await this.fetchUsers();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async toggleUserStatus(id) {
      try {
        const res = await usersApi.toggleStatus(id);
        await this.fetchUsers();
        const updated = res && res.data ? res.data : null;
        if (updated) {
          this.showToast('User Status', `${updated.name} is now ${updated.status}.`, 'info');
        }
      } catch (err) {
        this.showToast('Update Failed', err.message, 'error');
      }
    },

    async resetUserPassword(id, newPassword) {
      try {
        await usersApi.resetPassword(id, newPassword);
        this.showToast('Password Reset', 'Password updated successfully.', 'success');
      } catch (err) {
        this.showToast('Reset Failed', err.message, 'error');
      }
    },

    async deleteUser(id) {
      try {
        await usersApi.delete(id);
        this.users = this.users.filter(item => item.id !== id);
        this.showToast('User Deleted', 'User account removed.', 'warning');
      } catch (err) {
        this.showToast('Delete Failed', err.message, 'error');
      }
    },

    // ROLES & PERMISSIONS
    async fetchRoles() {
      try {
        const [rolesRes, permsRes] = await Promise.all([
          rolesApi.getAll(),
          rolesApi.getAllPermissions()
        ]);
        if (rolesRes && rolesRes.data) {
          this.roles = rolesRes.data;
        }
        if (permsRes && permsRes.data) {
          this.allPermissions = permsRes.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchRoles error:', err);
      }
    },

    async saveRole(roleData) {
      try {
        const isEdit = Boolean(roleData.id && this.roles.some(r => r.id === roleData.id));
        if (isEdit) {
          await rolesApi.update(roleData.id, roleData);
          this.showToast('Role Updated', `${roleData.name} updated.`);
        } else {
          await rolesApi.create(roleData);
          this.showToast('Role Created', `${roleData.name} created.`);
        }
        await this.fetchRoles();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    async saveRolePermissions(roleId, permissionCodes) {
      try {
        // Flatten matrix object if passed as nested structure
        let codes = [];
        if (Array.isArray(permissionCodes)) {
          codes = permissionCodes;
        } else if (typeof permissionCodes === 'object' && permissionCodes !== null) {
          Object.keys(permissionCodes).forEach(mod => {
            const actions = permissionCodes[mod];
            Object.keys(actions).forEach(act => {
              if (actions[act]) {
                codes.push(`${mod}.${act}`);
              }
            });
          });
        }

        await rolesApi.updatePermissions(roleId, codes);
        this.showToast('Permissions Saved', 'Role permissions updated successfully.');
        await this.fetchRoles();
        await this.fetchActivities();
      } catch (err) {
        this.showToast('Save Failed', err.message, 'error');
        throw err;
      }
    },

    // ACTIVITIES
    async fetchActivities(params = {}) {
      try {
        const res = await activityApi.getAll({ limit: 20, ...params });
        if (res && res.data) {
          this.activities = res.data;
        }
      } catch (err) {
        console.error('[MasterStore] fetchActivities error:', err);
      }
    },

    logActivity(action, targetType, targetName, targetId, badgeColor = 'blue') {
      const now = new Date();
      const dateStr = now.toISOString().slice(0, 16).replace('T', ' ');
      this.activities.unshift({
        id: `ACT-${Date.now()}`,
        user: this.currentUser.name,
        userAvatar: this.currentUser.avatar,
        action,
        targetType,
        targetName,
        targetId,
        timestamp: 'Just now',
        date: dateStr,
        badgeColor
      });
    },

    // MOBILE SIDEBAR ACTIONS
    toggleMobileSidebar() {
      this.mobileSidebarOpen = !this.mobileSidebarOpen;
    },
    closeMobileSidebar() {
      this.mobileSidebarOpen = false;
    },
    openMobileSidebar() {
      this.mobileSidebarOpen = true;
    },

    // FILE UPLOAD HELPER
    async uploadFile(file) {
      try {
        const res = await uploadApi.uploadFile(file);
        if (res && res.data) {
          const uploadBase = import.meta.env.VITE_API_UPLOAD_URL || 'http://localhost:8080';
          return {
            ...res.data,
            fullUrl: `${uploadBase}${res.data.url}`
          };
        }
        return res;
      } catch (err) {
        this.showToast('Upload Error', err.message, 'error');
        throw err;
      }
    }
  }
});
