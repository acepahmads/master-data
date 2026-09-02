<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-primary-400 uppercase tracking-wider mb-1">
          <router-link to="/testing" class="hover:underline flex items-center gap-1">
            <i class="feather-arrow-left text-xs"></i> Testing Command Center
          </router-link>
        </div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <i class="feather-layers text-primary-400"></i>
          Master Test Plan Protocols
        </h1>
        <p class="text-sm text-slate-400 mt-0.5">
          Reusable engineering test suites, parametric verification criteria, and versioned protocol specifications.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-sm font-medium flex items-center gap-2 transition shadow-lg shadow-primary-900/30"
        >
          <i class="feather-plus"></i>
          Create Test Plan
        </button>
      </div>
    </div>

    <!-- Filters & Search Bar -->
    <div class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-4 flex flex-col md:flex-row gap-4 items-center justify-between backdrop-blur-sm">
      <div class="flex items-center gap-3 w-full md:w-auto flex-1">
        <div class="relative w-full max-w-md">
          <i class="feather-search absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-sm"></i>
          <input
            v-model="searchQuery"
            @input="handleSearch"
            type="text"
            placeholder="Search test plans by code, name, or description..."
            class="w-full pl-9 pr-4 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-primary-500 transition"
          />
        </div>

        <select
          v-model="selectedCategory"
          @change="handleFilter"
          class="px-3 py-2 bg-slate-900/80 border border-slate-700/60 rounded-lg text-sm text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">All Categories</option>
          <option value="ELECTRICAL">Electrical & Power</option>
          <option value="FUNCTIONAL">Functional & Logic</option>
          <option value="RF_WIRELESS">RF & Wireless</option>
          <option value="POWER_CONSUMPTION">Power Consumption</option>
          <option value="THERMAL">Thermal Stress</option>
          <option value="COMMUNICATIONS">Communications</option>
          <option value="VISUAL_INSPECTION">Visual Inspection</option>
        </select>
      </div>

      <div class="text-xs text-slate-400">
        Showing <span class="font-semibold text-slate-200">{{ testPlans.length }}</span> protocols
      </div>
    </div>

    <!-- Protocols Grid -->
    <div v-if="loading && testPlans.length === 0" class="p-12 text-center text-slate-400 text-sm">
      <i class="feather-loader animate-spin text-2xl text-primary-400 mb-2"></i>
      <p>Loading master test plans...</p>
    </div>

    <div v-else-if="testPlans.length === 0" class="p-12 text-center bg-slate-800/40 border border-slate-700/40 rounded-xl">
      <i class="feather-layers text-4xl text-slate-500 mb-3 inline-block"></i>
      <h3 class="text-base font-semibold text-slate-200">No test plans found</h3>
      <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
        Create a new test plan protocol to define parametric limits, verification steps, and versioned test cases.
      </p>
      <button
        @click="openCreateModal"
        class="mt-4 px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white text-xs font-medium rounded-lg transition"
      >
        Create Protocol
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="plan in testPlans"
        :key="plan.id"
        class="bg-slate-800/80 border border-slate-700/60 rounded-xl p-5 hover:border-slate-600 transition flex flex-col justify-between backdrop-blur-sm group"
      >
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-mono px-2.5 py-1 rounded bg-primary-500/10 text-primary-300 border border-primary-500/20 font-bold">
              {{ plan.code }}
            </span>
            <span class="px-2 py-0.5 rounded text-[11px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              {{ plan.status }}
            </span>
          </div>

          <div>
            <h3 class="text-base font-bold text-slate-100 group-hover:text-primary-300 transition line-clamp-1">
              {{ plan.name }}
            </h3>
            <p class="text-xs text-slate-400 line-clamp-2 mt-1 leading-relaxed">
              {{ plan.description || 'No description provided.' }}
            </p>
          </div>

          <div class="pt-2 border-t border-slate-700/40 grid grid-cols-2 gap-2 text-xs text-slate-400">
            <div>
              <span class="text-slate-500 block text-[10px] uppercase">Category</span>
              <span class="font-medium text-slate-300">{{ plan.category || 'GENERAL' }}</span>
            </div>
            <div>
              <span class="text-slate-500 block text-[10px] uppercase">Versions</span>
              <span class="font-medium text-slate-300">{{ plan.versions ? plan.versions.length : 1 }} Release(s)</span>
            </div>
          </div>
        </div>

        <div class="mt-5 pt-3 border-t border-slate-700/40 flex items-center justify-between">
          <span class="text-[11px] text-slate-500">By {{ plan.createdBy }}</span>
          <router-link
            :to="`/test-plans/${plan.id}`"
            class="px-3.5 py-1.5 bg-slate-900 hover:bg-slate-700 text-slate-200 border border-slate-700 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition"
          >
            Workspace <i class="feather-arrow-right text-xs"></i>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
    >
      <div class="bg-slate-800 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
        <div class="flex items-center justify-between border-b border-slate-700/60 pb-3">
          <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
            <i class="feather-layers text-primary-400"></i>
            Create Master Test Plan
          </h2>
          <button @click="showModal = false" class="text-slate-400 hover:text-slate-200">
            <i class="feather-x text-lg"></i>
          </button>
        </div>

        <form @submit.prevent="submitCreatePlan" class="space-y-4 text-sm">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Protocol Code</label>
            <input
              v-model="form.code"
              type="text"
              placeholder="e.g. TP-GW400X-EVT"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500 font-mono"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Protocol Name *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="e.g. Smart Gateway EVT Qualification Suite"
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Category</label>
              <select
                v-model="form.category"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
              >
                <option value="ELECTRICAL">Electrical & Power</option>
                <option value="FUNCTIONAL">Functional & Logic</option>
                <option value="RF_WIRELESS">RF & Wireless</option>
                <option value="POWER_CONSUMPTION">Power Consumption</option>
                <option value="THERMAL">Thermal Stress</option>
                <option value="COMMUNICATIONS">Communications</option>
                <option value="VISUAL_INSPECTION">Visual Inspection</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Initial Version</label>
              <input
                v-model="form.initialVersion"
                type="text"
                placeholder="v1.0"
                class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Description & Scope</label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Detailed description of the test protocol, target hardware, and scope..."
              class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg text-slate-200 focus:outline-none focus:border-primary-500"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-700/60 flex items-center justify-end gap-3">
            <button
              type="button"
              @click="showModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg text-xs font-medium transition"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-xs font-semibold transition flex items-center gap-1.5"
            >
              <i v-if="submitting" class="feather-loader animate-spin text-xs"></i>
              Create Protocol
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia';
import { useTestingStore } from '@/stores/testingStore';

export default {
  name: 'TestPlanRegistryView',
  data() {
    return {
      searchQuery: '',
      selectedCategory: '',
      showModal: false,
      submitting: false,
      form: {
        code: '',
        name: '',
        category: 'ELECTRICAL',
        initialVersion: 'v1.0',
        description: ''
      }
    };
  },
  computed: {
    ...mapState(useTestingStore, ['testPlans', 'loading'])
  },
  mounted() {
    this.loadPlans();
  },
  methods: {
    ...mapActions(useTestingStore, ['fetchTestPlans', 'createTestPlan']),
    async loadPlans() {
      await this.fetchTestPlans({
        q: this.searchQuery,
        category: this.selectedCategory
      });
    },
    handleSearch() {
      this.loadPlans();
    },
    handleFilter() {
      this.loadPlans();
    },
    openCreateModal() {
      this.form = {
        code: `TP-PROT-${Math.floor(1000 + Math.random() * 9000)}`,
        name: '',
        category: 'ELECTRICAL',
        initialVersion: 'v1.0',
        description: ''
      };
      this.showModal = true;
    },
    async submitCreatePlan() {
      this.submitting = true;
      try {
        const created = await this.createTestPlan(this.form);
        this.showModal = false;
        this.$router.push(`/test-plans/${created.id}`);
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      } finally {
        this.submitting = false;
      }
    }
  }
};
</script>
