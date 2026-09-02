<template>
  <BaseModal :show="isOpen" title="Documentation & Backend API Specification" maxWidth="3xl" @close="close">
    <template #icon>
      <div class="w-6 h-6 rounded bg-brand-50 text-brand-600 flex items-center justify-center">
        <AppIcon name="help" size="xs" />
      </div>
    </template>

    <div class="space-y-4 text-xs text-slate-700">
      <!-- Section 1: Overview -->
      <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg">
        <h4 class="font-bold text-slate-900 text-sm mb-1">Product & Master Data Control Center — Platform Architecture</h4>
        <p class="leading-relaxed text-slate-600">
          Centralized Product and Master Data Management platform governing Trading Products, Project Solution Compositions, and In-House R&D Engineering Lifecycles (Phases 1–6) over a unified Master Data Foundation.
        </p>
      </div>

      <!-- Section 2: Technical Architecture -->
      <div>
        <h5 class="font-bold text-slate-900 uppercase tracking-wider text-2xs mb-2">Target Stack & Compatibility</h5>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-2 font-mono text-2xs">
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="font-semibold text-brand-700 mb-1">Frontend Layer</div>
            <ul class="space-y-0.5 text-slate-600">
              <li>• Vue.js 2.7 (Composition API + Options API)</li>
              <li>• Tailwind CSS (Enterprise Design System)</li>
              <li>• Pinia State Management</li>
              <li>• Vue Router SPA Navigation</li>
            </ul>
          </div>
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded">
            <div class="font-semibold text-emerald-700 mb-1">Target Backend & Database</div>
            <ul class="space-y-0.5 text-slate-600">
              <li>• Golang 1.22+ (Gin Web Framework)</li>
              <li>• GORM (Object Relational Mapping)</li>
              <li>• MySQL 8.0 (Relational Master Database)</li>
              <li>• RESTful OpenAPI 3.0 Standard</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Section 3: REST API Endpoint Mapping -->
      <div>
        <h5 class="font-bold text-slate-900 uppercase tracking-wider text-2xs mb-2">Phase 1 REST API Endpoints (Gin / GORM)</h5>
        <div class="border border-slate-200 rounded-lg overflow-hidden font-mono text-2xs">
          <table class="min-w-full divide-y divide-slate-200">
            <thead class="bg-slate-100 font-semibold text-slate-700">
              <tr>
                <th class="px-3 py-2 text-left">Method</th>
                <th class="px-3 py-2 text-left">Endpoint</th>
                <th class="px-3 py-2 text-left">Description</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 bg-white">
              <tr>
                <td class="px-3 py-1.5 text-emerald-600 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/products</td>
                <td class="px-3 py-1.5 text-slate-600">List and create IoT Product Masters</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-blue-600 font-bold">GET/PUT/DEL</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/products/:id</td>
                <td class="px-3 py-1.5 text-slate-600">Get, update, or delete single product</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-emerald-600 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/components</td>
                <td class="px-3 py-1.5 text-slate-600">List and create Engineering Components</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-blue-600 font-bold">GET/PUT/DEL</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/components/:id</td>
                <td class="px-3 py-1.5 text-slate-600">Get, update, or delete component specs</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-purple-600 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/categories/tree</td>
                <td class="px-3 py-1.5 text-slate-600">Hierarchical category tree structure</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-slate-700 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/manufacturers</td>
                <td class="px-3 py-1.5 text-slate-600">Manufacturer directory and rating</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-slate-700 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/suppliers</td>
                <td class="px-3 py-1.5 text-slate-600">Supplier master data & procurement terms</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-slate-700 font-bold">GET/POST</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/units</td>
                <td class="px-3 py-1.5 text-slate-600">Units of Measurement (UoM)</td>
              </tr>
              <tr>
                <td class="px-3 py-1.5 text-amber-600 font-bold">GET/PUT</td>
                <td class="px-3 py-1.5 text-slate-800">/api/v1/roles/:id/permissions</td>
                <td class="px-3 py-1.5 text-slate-600">RBAC permission matrix management</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <template #footer>
      <button type="button" class="btn btn-secondary" @click="close">
        Close
      </button>
    </template>
  </BaseModal>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import BaseModal from './BaseModal.vue';
import AppIcon from './AppIcon.vue';

export default {
  name: 'HelpModal',
  components: { BaseModal, AppIcon },
  computed: {
    store() {
      return useMasterStore();
    },
    isOpen() {
      return this.store.helpModalOpen;
    }
  },
  methods: {
    close() {
      this.store.helpModalOpen = false;
    }
  }
};
</script>
