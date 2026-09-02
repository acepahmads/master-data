<template>
  <div class="space-y-5">
    <!-- Page Header -->
    <PageHeader
      title="Product Versions"
      subtitle="Manage hardware iterations, engineering version history, and product lifecycle state baselines"
    >
      <template #badge>
        <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-brand-50 text-brand-700 border border-brand-200">
          {{ versionStore.totalVersionsCount }} Versions
        </span>
      </template>
      <template #actions>
        <router-link to="/versions/compare" class="btn btn-secondary">
          <AppIcon name="git-branch" size="xs" class="mr-1.5 text-slate-500" />
          <span>Compare Versions</span>
        </router-link>
        <router-link to="/versions/create" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
          <span>Create Version</span>
        </router-link>
      </template>
    </PageHeader>

    <!-- Version KPI Summary (5 Compact Cards) -->
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3.5">
      <KpiCard
        title="Total Products"
        :value="uniqueProductsCount"
        icon="product"
        iconColor="blue"
        trend="Active"
        trendType="up"
        subtitle="Versioned IoT Hardware"
        :sparkVariant="1"
      />
      <KpiCard
        title="Total Versions"
        :value="versionStore.totalVersionsCount"
        icon="layers"
        iconColor="purple"
        trend="Baselines"
        trendType="up"
        subtitle="Engineering Releases"
        :sparkVariant="2"
      />
      <KpiCard
        title="Active Versions"
        :value="versionStore.activeVersionsCount"
        icon="check-circle"
        iconColor="emerald"
        trend="Production"
        trendType="up"
        subtitle="Current Active Deployments"
        :sparkVariant="3"
      />
      <KpiCard
        title="Draft Versions"
        :value="versionStore.draftVersionsCount"
        icon="edit"
        iconColor="amber"
        trend="In Progress"
        trendType="neutral"
        subtitle="Under Engineering Review"
        :sparkVariant="4"
      />
      <KpiCard
        title="Deprecated"
        :value="versionStore.deprecatedVersionsCount"
        icon="alert-circle"
        iconColor="cyan"
        trend="Archived"
        trendType="down"
        subtitle="Historical Predecessors"
        :sparkVariant="5"
      />
    </div>

    <!-- Filter & View Mode Bar -->
    <div class="app-card p-3.5 flex flex-col md:flex-row md:items-center justify-between gap-3 bg-white">
      <div class="flex flex-wrap items-center gap-2.5 flex-1">
        <!-- Search Input -->
        <div class="relative min-w-[240px] max-w-sm flex-1">
          <AppIcon name="search" size="xs" class="absolute left-3 top-2.5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search version (e.g. 2.1.0), product, or engineer..."
            class="form-input pl-8.5 text-xs font-mono"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-2.5 top-2 text-slate-400 hover:text-slate-600"
          >
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <!-- Product Filter -->
        <select v-model="selectedProduct" class="form-select w-auto min-w-[160px] text-xs">
          <option value="">All Products</option>
          <option v-for="p in uniqueProducts" :key="p.code" :value="p.code">
            {{ p.code }} — {{ p.name }}
          </option>
        </select>

        <!-- Status Filter -->
        <select v-model="selectedStatus" class="form-select w-auto min-w-[130px] text-xs">
          <option value="">All Statuses</option>
          <option value="Active">Active</option>
          <option value="Draft">Draft</option>
          <option value="Deprecated">Deprecated</option>
          <option value="Archived">Archived</option>
        </select>

        <!-- Reset Button -->
        <button
          v-if="isFiltered"
          @click="resetFilters"
          class="btn btn-ghost text-slate-500 hover:text-slate-800"
        >
          <AppIcon name="refresh" size="xs" class="mr-1" />
          <span>Reset</span>
        </button>
      </div>

      <!-- View Switcher (Tree Hierarchy vs Flat Table) -->
      <div class="flex items-center gap-2 shrink-0">
        <div class="flex items-center p-0.5 rounded-lg bg-slate-100 border border-slate-200">
          <button
            @click="viewMode = 'tree'"
            :class="[
              'px-2.5 py-1 rounded-md text-xs font-semibold flex items-center gap-1.5 transition-all',
              viewMode === 'tree' ? 'bg-white text-brand-600 shadow-2xs' : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <AppIcon name="tree" size="xs" />
            <span>Product Tree</span>
          </button>
          <button
            @click="viewMode = 'table'"
            :class="[
              'px-2.5 py-1 rounded-md text-xs font-semibold flex items-center gap-1.5 transition-all',
              viewMode === 'table' ? 'bg-white text-brand-600 shadow-2xs' : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <AppIcon name="layers" size="xs" />
            <span>Flat Table</span>
          </button>
        </div>
      </div>
    </div>

    <!-- VIEW MODE 1: PRODUCT HIERARCHY TREE VIEW -->
    <div v-if="viewMode === 'tree'" class="space-y-4">
      <div v-if="filteredVersions.length === 0" class="app-card p-8">
        <EmptyState
          title="No versions found"
          description="No product versions match your current search filter criteria."
          icon="layers"
        >
          <template #action>
            <button @click="resetFilters" class="btn btn-secondary mr-2">Reset Filters</button>
            <router-link to="/versions/create" class="btn btn-primary">Create Version</router-link>
          </template>
        </EmptyState>
      </div>
      <LifecycleTree v-else :filterProductCode="selectedProduct" />
    </div>

    <!-- VIEW MODE 2: FLAT TABLE VIEW -->
    <div v-else class="app-card overflow-hidden">
      <div v-if="filteredVersions.length === 0" class="p-8">
        <EmptyState
          title="No versions found"
          description="No product versions match your active search filters."
          icon="layers"
        >
          <template #action>
            <button @click="resetFilters" class="btn btn-secondary mr-2">Reset Filters</button>
            <router-link to="/versions/create" class="btn btn-primary">Create Version</router-link>
          </template>
        </EmptyState>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="dense-table">
          <thead>
            <tr>
              <th>Product Master</th>
              <th>Version</th>
              <th>Version Name</th>
              <th>Status</th>
              <th>Hardware Revision</th>
              <th>Lead Engineer</th>
              <th>Last Updated</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="ver in filteredVersions"
              :key="ver.id"
              @click="$router.push(`/versions/${ver.id}`)"
              class="cursor-pointer group"
            >
              <td>
                <div>
                  <div class="font-bold text-slate-900 text-xs truncate max-w-xs group-hover:text-brand-600 transition-colors">
                    {{ ver.productName }}
                  </div>
                  <span class="font-mono font-semibold text-3xs text-brand-700 bg-brand-50 px-1.5 py-0.2 rounded border border-brand-200">
                    {{ ver.productCode }}
                  </span>
                </div>
              </td>
              <td>
                <span class="font-mono font-bold text-sm text-slate-900 bg-slate-100 px-2 py-0.5 rounded border border-slate-200 shadow-2xs">
                  v{{ ver.versionNumber }}
                </span>
              </td>
              <td>
                <div class="text-slate-700 truncate max-w-xs font-medium text-xs">
                  {{ ver.versionName }}
                </div>
                <div v-if="ver.baseVersionId" class="text-3xs text-slate-400 font-mono">
                  Base: {{ ver.baseVersionId }}
                </div>
              </td>
              <td>
                <VersionBadge :status="ver.status" />
              </td>
              <td>
                <span class="font-mono font-bold text-xs text-slate-800 bg-slate-50 px-2 py-0.5 rounded border border-slate-200">
                  {{ ver.currentRevision }}
                </span>
                <span class="text-3xs text-slate-400 ml-1">({{ ver.revisionsCount }} revs)</span>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <img :src="ver.ownerAvatar" class="w-5 h-5 rounded-full object-cover border border-slate-200" />
                  <span class="text-xs font-medium text-slate-800">{{ ver.owner }}</span>
                </div>
              </td>
              <td class="font-mono text-3xs text-slate-500">
                {{ ver.updatedAt || ver.releaseDate }}
              </td>
              <td class="text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <router-link
                    :to="`/versions/${ver.id}`"
                    class="btn-icon"
                    title="View Version Details"
                  >
                    <AppIcon name="eye" size="xs" />
                  </router-link>
                  <router-link
                    :to="`/versions/${ver.id}/edit`"
                    class="btn-icon"
                    title="Edit Version"
                  >
                    <AppIcon name="edit" size="xs" />
                  </router-link>
                  <router-link
                    :to="`/versions/compare?vA=${ver.id}`"
                    class="btn-icon text-slate-400 hover:text-brand-600"
                    title="Compare Version"
                  >
                    <AppIcon name="git-branch" size="xs" />
                  </router-link>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import PageHeader from '@/components/common/PageHeader.vue';
import KpiCard from '@/components/common/KpiCard.vue';
import VersionBadge from '@/components/versions/VersionBadge.vue';
import LifecycleTree from '@/components/versions/LifecycleTree.vue';
import EmptyState from '@/components/common/EmptyState.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ProductVersionListView',
  components: {
    PageHeader,
    KpiCard,
    VersionBadge,
    LifecycleTree,
    EmptyState,
    AppIcon
  },
  data() {
    return {
      searchQuery: '',
      selectedProduct: '',
      selectedStatus: '',
      viewMode: 'tree' // 'tree' or 'table'
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    uniqueProducts() {
      const map = new Map();
      this.versionStore.versions.forEach(v => {
        if (!map.has(v.productCode)) {
          map.set(v.productCode, { code: v.productCode, name: v.productName });
        }
      });
      return Array.from(map.values());
    },
    uniqueProductsCount() {
      return this.uniqueProducts.length;
    },
    isFiltered() {
      return Boolean(this.searchQuery || this.selectedProduct || this.selectedStatus);
    },
    filteredVersions() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.versionStore.versions.filter(v => {
        const matchesSearch = !q || v.versionNumber.toLowerCase().includes(q) || v.versionName.toLowerCase().includes(q) || v.productName.toLowerCase().includes(q) || v.productCode.toLowerCase().includes(q) || (v.owner && v.owner.toLowerCase().includes(q));
        const matchesProduct = !this.selectedProduct || v.productCode === this.selectedProduct;
        const matchesStatus = !this.selectedStatus || v.status === this.selectedStatus;
        return matchesSearch && matchesProduct && matchesStatus;
      });
    }
  },
  methods: {
    resetFilters() {
      this.searchQuery = '';
      this.selectedProduct = '';
      this.selectedStatus = '';
    }
  },
  async mounted() {
    await Promise.all([
      this.versionStore.fetchVersions(),
      this.versionStore.fetchRevisions()
    ]);
  }
};
</script>
