<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <PageHeader title="Component Categories" subtitle="Hierarchical engineering master taxonomy">
      <template #actions>
        <button @click="openAddRootCategoryModal" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5" />
          <span>Add Root Category</span>
        </button>
      </template>
    </PageHeader>

    <!-- Main Tree + Detail 2-Column Split Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
      <!-- LEFT PANEL: CATEGORY HIERARCHY TREE (5 Cols) -->
      <div class="lg:col-span-5 bg-white rounded-lg border border-slate-200 shadow-card overflow-hidden">
        <div class="p-3 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <AppIcon name="tree" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Category Hierarchy</h2>
          </div>
          <span class="text-3xs font-mono text-slate-500">Click to inspect</span>
        </div>

        <!-- Search categories in tree -->
        <div class="p-2.5 border-b border-slate-100 bg-slate-50/50">
          <div class="relative">
            <AppIcon name="search" size="xs" class="absolute left-2.5 top-2.5 text-slate-400" />
            <input
              v-model="treeSearch"
              type="text"
              placeholder="Filter taxonomy tree..."
              class="form-input pl-8 py-1 text-2xs"
            />
          </div>
        </div>

        <!-- Tree Nodes List -->
        <div class="p-2 space-y-1 max-h-[600px] overflow-y-auto">
          <div
            v-for="rootNode in filteredCategories"
            :key="rootNode.id"
            class="space-y-1"
          >
            <!-- Root Level Item -->
            <div
              @click="selectCategory(rootNode)"
              :class="[
                'flex items-center justify-between px-2.5 py-1.5 rounded-md cursor-pointer text-xs transition-colors group',
                selectedCat && selectedCat.id === rootNode.id
                  ? 'bg-brand-50 text-brand-900 font-bold border border-brand-200'
                  : 'hover:bg-slate-100 text-slate-800'
              ]"
            >
              <div class="flex items-center gap-2 min-w-0">
                <button
                  @click.stop="toggleNodeExpand(rootNode.id)"
                  class="p-0.5 rounded text-slate-400 hover:text-slate-700"
                >
                  <AppIcon
                    :name="expandedNodeIds.includes(rootNode.id) ? 'chevron-down' : 'chevron-right'"
                    size="xs"
                  />
                </button>
                <AppIcon name="folder" size="xs" class="text-brand-600 shrink-0" />
                <span class="truncate">{{ rootNode.name }}</span>
              </div>
              <span class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-100 text-slate-500 group-hover:bg-white">
                {{ rootNode.totalComponents }}
              </span>
            </div>

            <!-- Nested Subcategories (Level 2) -->
            <div
              v-if="expandedNodeIds.includes(rootNode.id)"
              class="pl-6 space-y-0.5 border-l-2 border-slate-100 ml-3.5 mt-0.5"
            >
              <div
                v-for="child in rootNode.children"
                :key="child.id"
                @click="selectCategory(child, rootNode)"
                :class="[
                  'flex items-center justify-between px-2 py-1 rounded cursor-pointer text-xs transition-colors group',
                  selectedCat && selectedCat.id === child.id
                    ? 'bg-brand-50 text-brand-900 font-bold border border-brand-200'
                    : 'hover:bg-slate-100 text-slate-700'
                ]"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <span class="w-1.5 h-1.5 rounded-full bg-slate-300 group-hover:bg-brand-500 shrink-0"></span>
                  <span class="truncate">{{ child.name }}</span>
                </div>
                <span class="text-3xs font-mono px-1.5 py-0.2 rounded bg-slate-100 text-slate-500">
                  {{ child.totalComponents }}
                </span>
              </div>

              <!-- Quick inline Add subcategory button -->
              <button
                @click="openAddSubcategoryModal(rootNode)"
                class="w-full text-left px-2 py-1 text-3xs font-semibold text-brand-600 hover:text-brand-800 hover:bg-brand-50/50 rounded flex items-center gap-1 mt-1 transition-colors"
              >
                <AppIcon name="plus" size="xs" />
                <span>Add Subcategory to {{ rootNode.code }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT PANEL: CATEGORY DETAIL PANEL (7 Cols) -->
      <div class="lg:col-span-7 space-y-4">
        <div v-if="selectedCat" class="bg-white rounded-lg border border-slate-200 shadow-card p-5 space-y-5">
          <!-- Detail Header -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-200 gap-3">
            <div>
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                  {{ selectedCat.code }}
                </span>
                <StatusBadge :status="selectedCat.status" />
              </div>
              <h2 class="text-lg font-bold text-slate-900 tracking-tight mt-1">{{ selectedCat.name }}</h2>
              <div v-if="selectedParent" class="text-2xs text-slate-500 mt-0.5">
                Parent: <strong class="text-slate-700">{{ selectedParent.name }}</strong> ({{ selectedParent.code }})
              </div>
            </div>

            <div class="flex items-center gap-2">
              <button @click="openAddSubcategoryModal(selectedCat)" class="btn btn-secondary text-2xs">
                <AppIcon name="plus" size="xs" class="mr-1" />
                <span>Add Child</span>
              </button>
              <button @click="openEditModal(selectedCat)" class="btn btn-secondary text-2xs">
                <AppIcon name="edit" size="xs" class="mr-1" />
                <span>Edit</span>
              </button>
              <button @click="confirmDelete(selectedCat)" class="btn btn-danger text-2xs">
                <AppIcon name="trash" size="xs" />
              </button>
            </div>
          </div>

          <!-- Structured Information Fields -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
            <div class="p-3 bg-slate-50 border border-slate-200 rounded-md">
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Total Registered Components</div>
              <div class="text-xl font-bold font-mono text-brand-700 mt-0.5">{{ selectedCat.totalComponents }}</div>
            </div>
            <div class="p-3 bg-slate-50 border border-slate-200 rounded-md">
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Taxonomy Tier</div>
              <div class="text-sm font-bold text-slate-800 mt-1">{{ selectedCat.parentId ? 'Level 2 (Subcategory)' : 'Level 1 (Root Category)' }}</div>
            </div>
          </div>

          <!-- Description -->
          <div>
            <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Description & Scope</h3>
            <p class="text-xs text-slate-600 bg-slate-50 p-3 rounded-md border border-slate-200 leading-relaxed">
              {{ selectedCat.description || 'No description provided for this category branch.' }}
            </p>
          </div>

          <!-- Subcategories Table (if root or has children) -->
          <div v-if="selectedCat.children && selectedCat.children.length > 0">
            <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-2">Direct Child Subcategories ({{ selectedCat.children.length }})</h3>
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <table class="dense-table">
                <thead>
                  <tr>
                    <th>Subcategory Code</th>
                    <th>Name</th>
                    <th>Components</th>
                    <th class="text-right">Action</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="child in selectedCat.children" :key="child.id">
                    <td class="font-mono font-bold text-slate-800">{{ child.code }}</td>
                    <td class="font-semibold text-brand-700 cursor-pointer" @click="selectCategory(child, selectedCat)">
                      {{ child.name }}
                    </td>
                    <td class="font-mono text-2xs">{{ child.totalComponents }}</td>
                    <td class="text-right">
                      <button @click="selectCategory(child, selectedCat)" class="btn-icon" title="View details">
                        <AppIcon name="chevron-right" size="xs" />
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Sample Mapped Components in Category -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">Sample Mapped Components</h3>
              <router-link to="/components" class="text-3xs text-brand-600 hover:underline">
                View all in Central Library →
              </router-link>
            </div>
            <div class="border border-slate-200 rounded-lg overflow-hidden">
              <table class="dense-table">
                <thead>
                  <tr>
                    <th>MPN</th>
                    <th>Name</th>
                    <th>Manufacturer</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="c in sampleCategoryComponents"
                    :key="c.id"
                    @click="$router.push(`/components/${c.id}`)"
                    class="cursor-pointer hover:bg-slate-50"
                  >
                    <td class="font-mono font-bold text-brand-700">{{ c.partNumber }}</td>
                    <td class="text-slate-800 font-medium truncate max-w-xs">{{ c.name }}</td>
                    <td class="text-slate-600">{{ c.manufacturer }}</td>
                    <td><StatusBadge :status="c.status" /></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div v-else class="bg-white rounded-lg border border-slate-200 shadow-card p-10 text-center text-slate-400">
          <AppIcon name="tree" size="lg" class="mx-auto mb-2 text-slate-300" />
          <h3 class="text-sm font-bold text-slate-800">Select a Category</h3>
          <p class="text-xs text-slate-500 mt-1">Select any root or subcategory node from the left tree to view properties and child relations.</p>
        </div>
      </div>
    </div>

    <!-- Create / Edit Category Modal -->
    <BaseModal :show="categoryModalOpen" :title="isEditModal ? 'Edit Category' : 'Create Category'" maxWidth="md" @close="categoryModalOpen = false">
      <form @submit.prevent="saveCategoryModal" class="space-y-4 text-xs">
        <div>
          <label class="form-label">Parent Category</label>
          <select v-model="modalForm.parentId" class="form-select">
            <option :value="null">-- Root Level Category --</option>
            <option v-for="cat in store.categories" :key="cat.id" :value="cat.id">
              {{ cat.name }} ({{ cat.code }})
            </option>
          </select>
        </div>

        <div>
          <label class="form-label">Category Name *</label>
          <input v-model="modalForm.name" type="text" required placeholder="e.g. Communication Modules" class="form-input" />
        </div>

        <div>
          <label class="form-label">Category Code *</label>
          <input v-model="modalForm.code" type="text" required placeholder="e.g. ELEC-COMM" class="form-input font-mono uppercase" />
        </div>

        <div>
          <label class="form-label">Description</label>
          <textarea v-model="modalForm.description" rows="3" placeholder="Category scope and engineering definitions..." class="form-textarea"></textarea>
        </div>

        <div>
          <label class="form-label">Status</label>
          <select v-model="modalForm.status" class="form-select">
            <option value="Active">Active</option>
            <option value="Draft">Draft</option>
          </select>
        </div>

        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
          <button type="button" class="btn btn-secondary" @click="categoryModalOpen = false">Cancel</button>
          <button type="submit" class="btn btn-primary">Save Category</button>
        </div>
      </form>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <ConfirmDialog
      :show="deleteModalOpen"
      title="Delete Category Node"
      :message="`Are you sure you want to delete category '${categoryToDelete ? categoryToDelete.name : ''}'? Any subcategories will also be affected.`"
      @confirm="executeDelete"
      @cancel="deleteModalOpen = false"
    />
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import PageHeader from '@/components/common/PageHeader.vue';
import StatusBadge from '@/components/common/StatusBadge.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'CategoryListView',
  components: {
    PageHeader,
    StatusBadge,
    ConfirmDialog,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      treeSearch: '',
      expandedNodeIds: ['CAT-ELEC', 'CAT-MECH', 'CAT-OPT'],
      selectedCat: null,
      selectedParent: null,
      categoryModalOpen: false,
      isEditModal: false,
      modalForm: {
        id: null,
        name: '',
        code: '',
        description: '',
        parentId: null,
        status: 'Active'
      },
      deleteModalOpen: false,
      categoryToDelete: null
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    filteredCategories() {
      if (!this.treeSearch) return this.store.categories;
      const q = this.treeSearch.toLowerCase().trim();
      return this.store.categories.filter(c => {
        return c.name.toLowerCase().includes(q) || c.code.toLowerCase().includes(q) ||
          (c.children && c.children.some(child => child.name.toLowerCase().includes(q) || child.code.toLowerCase().includes(q)));
      });
    },
    sampleCategoryComponents() {
      if (!this.selectedCat) return [];
      return this.store.components.filter(c => {
        return c.category === this.selectedCat.name || (this.selectedCat.children && this.selectedCat.children.some(ch => ch.name === c.category));
      }).slice(0, 5);
    }
  },
  created() {
    if (this.store.categories.length > 0) {
      this.selectedCat = this.store.categories[0];
    }
  },
  methods: {
    toggleNodeExpand(nodeId) {
      const idx = this.expandedNodeIds.indexOf(nodeId);
      if (idx !== -1) {
        this.expandedNodeIds.splice(idx, 1);
      } else {
        this.expandedNodeIds.push(nodeId);
      }
    },
    selectCategory(node, parent = null) {
      this.selectedCat = node;
      this.selectedParent = parent;
    },
    openAddRootCategoryModal() {
      this.isEditModal = false;
      this.modalForm = {
        id: null,
        name: '',
        code: '',
        description: '',
        parentId: null,
        status: 'Active'
      };
      this.categoryModalOpen = true;
    },
    openAddSubcategoryModal(parentNode) {
      this.isEditModal = false;
      this.modalForm = {
        id: null,
        name: '',
        code: `${parentNode.code}-`,
        description: '',
        parentId: parentNode.id,
        status: 'Active'
      };
      this.categoryModalOpen = true;
    },
    openEditModal(node) {
      this.isEditModal = true;
      this.modalForm = {
        id: node.id,
        name: node.name,
        code: node.code,
        description: node.description,
        parentId: node.parentId,
        status: node.status
      };
      this.categoryModalOpen = true;
    },
    saveCategoryModal() {
      this.store.saveCategory(this.modalForm);
      this.categoryModalOpen = false;
      if (!this.expandedNodeIds.includes(this.modalForm.parentId)) {
        this.expandedNodeIds.push(this.modalForm.parentId);
      }
    },
    confirmDelete(node) {
      this.categoryToDelete = node;
      this.deleteModalOpen = true;
    },
    executeDelete() {
      if (this.categoryToDelete) {
        this.store.deleteCategory(this.categoryToDelete.id);
        this.deleteModalOpen = false;
        this.categoryToDelete = null;
        this.selectedCat = this.store.categories[0] || null;
      }
    }
  }
};
</script>
