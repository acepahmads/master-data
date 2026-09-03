<template>
  <div class="space-y-6">
    <!-- 1. DASHBOARD HERO & COMMAND OVERVIEW -->
    <div class="bg-gradient-to-r from-slate-900 via-navy-850 to-slate-900 rounded-2xl p-5 md:p-6 text-white shadow-lg border border-slate-800 relative overflow-hidden">
      <!-- Background Ambient Glow -->
      <div class="absolute -right-16 -top-16 w-64 h-64 bg-brand-500/10 rounded-full blur-3xl pointer-events-none"></div>
      <div class="absolute right-1/3 -bottom-16 w-48 h-48 bg-cyan-500/10 rounded-full blur-2xl pointer-events-none"></div>

      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div class="flex items-center gap-2 mb-1.5">
            <span class="w-2 h-2 rounded-full bg-emerald-400 ring-4 ring-emerald-400/20 animate-pulse"></span>
            <span class="text-3xs font-mono font-bold tracking-widest uppercase text-emerald-400">Control Center Online</span>
            <span class="text-slate-600">•</span>
            <span class="text-3xs font-mono text-slate-400">REST API v1.0.0 Online</span>
          </div>
          <h1 class="text-xl md:text-2xl font-bold tracking-tight text-white flex items-center gap-2">
            <span>Good day, {{ currentUserName }}</span>
            <span class="text-lg">⚡</span>
          </h1>
          <p class="text-xs text-slate-300 max-w-2xl mt-1 leading-relaxed">
            Centralized platform governing Product Master (Trading, Project, Manufacture & Service), central component library, parametric specifications, and multi-tier engineering lifecycles.
          </p>
        </div>

        <!-- Quick Action Buttons -->
        <div class="flex items-center gap-2.5 shrink-0">
          <button
            @click="$router.push('/products/create')"
            class="px-3.5 py-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-white text-xs font-semibold border border-slate-700 hover:border-slate-600 transition-all duration-150 flex items-center gap-1.5 shadow-sm active:scale-95"
          >
            <AppIcon name="plus" size="xs" class="text-slate-300" />
            <span>New Product</span>
          </button>
          <button
            @click="$router.push('/components/create')"
            class="px-3.5 py-2 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-md shadow-brand-600/30 hover:shadow-glow-brand transition-all duration-150 flex items-center gap-1.5 active:scale-95"
          >
            <AppIcon name="plus" size="xs" class="text-white" />
            <span>Register Component</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 2. KPI CARDS SECTION (5 ENRICHED STATS) -->
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3.5">
      <KpiCard
        title="Total Products"
        :value="store.totalProductsCount"
        icon="product"
        iconColor="blue"
        trend="+14% MoM"
        trendType="up"
        subtitle="IoT Active Masters"
        :sparkVariant="1"
      />
      <KpiCard
        title="Central Library"
        :value="store.totalComponentsCount"
        icon="component"
        iconColor="emerald"
        trend="+24 items"
        trendType="up"
        subtitle="Verified Hardware Parts"
        :sparkVariant="2"
      />
      <KpiCard
        title="Categories"
        :value="store.totalCategoriesCount"
        icon="tree"
        iconColor="purple"
        trend="3 Tiers"
        trendType="neutral"
        subtitle="Hierarchical Nodes"
        :sparkVariant="3"
      />
      <KpiCard
        title="Manufacturers"
        :value="store.totalManufacturersCount"
        icon="manufacturer"
        iconColor="amber"
        trend="100% Valid"
        trendType="up"
        subtitle="Semiconductor Partners"
        :sparkVariant="4"
      />
      <KpiCard
        title="Suppliers"
        :value="store.totalSuppliersCount"
        icon="supplier"
        iconColor="cyan"
        trend="Tier 1"
        trendType="up"
        subtitle="Authorized Distributors"
        :sparkVariant="5"
      />
    </div>

    <!-- 3. MAIN 2-COLUMN OPERATIONAL GRID -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- LEFT 8 COLS: Recent Products Showcase & Engineering Components Library -->
      <div class="lg:col-span-8 space-y-6">
        <!-- A. RECENT PRODUCTS SHOWCASE -->
        <div class="app-card overflow-hidden">
          <div class="px-5 py-3.5 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div class="flex items-center gap-2.5">
              <div class="w-6 h-6 rounded-md bg-brand-50 text-brand-600 flex items-center justify-center border border-brand-200/60">
                <AppIcon name="product" size="xs" />
              </div>
              <div>
                <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Recent IoT Products</h2>
                <p class="text-3xs text-slate-400">Master product records and hardware configurations</p>
              </div>
            </div>
            <router-link
              to="/products"
              class="text-xs font-semibold text-brand-600 hover:text-brand-700 flex items-center gap-1 group"
            >
              <span>View All ({{ store.totalProductsCount }})</span>
              <AppIcon name="chevron-right" size="xs" class="group-hover:translate-x-0.5 transition-transform" />
            </router-link>
          </div>

          <div class="overflow-x-auto">
            <table class="dense-table">
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Category</th>
                  <th>Status</th>
                  <th>Lead Engineer</th>
                  <th>Last Updated</th>
                  <th class="text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="p in recentProducts"
                  :key="p.id"
                  @click="$router.push(`/products/${p.id || p.code}`)"
                  class="cursor-pointer group"
                >
                  <td>
                    <div class="flex items-center gap-3">
                      <div class="w-9 h-9 rounded-lg bg-slate-100 border border-slate-200 overflow-hidden shrink-0 group-hover:border-brand-300 transition-colors">
                        <img
                          v-if="p.image"
                          :src="p.image"
                          alt="Product"
                          class="w-full h-full object-cover"
                        />
                        <div v-else class="w-full h-full flex items-center justify-center text-slate-400">
                          <AppIcon name="product" size="xs" />
                        </div>
                      </div>
                      <div class="min-w-0">
                        <div class="font-bold text-slate-900 truncate max-w-xs group-hover:text-brand-600 transition-colors">
                          {{ p.name }}
                        </div>
                        <div class="text-3xs font-mono font-semibold text-brand-700 mt-0.5">
                          {{ p.code }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 border border-slate-200/80">
                      {{ p.category }}
                    </span>
                  </td>
                  <td>
                    <StatusBadge :status="p.status" />
                  </td>
                  <td>
                    <div class="text-xs font-medium text-slate-700 truncate max-w-[120px]">
                      {{ p.projectLead || 'Unassigned' }}
                    </div>
                  </td>
                  <td class="font-mono text-3xs text-slate-500">
                    {{ formatTimestamp(p.updatedAt || p.createdAt) }}
                  </td>
                  <td class="text-right" @click.stop>
                    <router-link
                      :to="`/products/${p.id || p.code}`"
                      class="btn-icon inline-flex"
                      title="View Product Detail"
                    >
                      <AppIcon name="eye" size="xs" />
                    </router-link>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- B. RECENT ENGINEERING COMPONENTS LIBRARY -->
        <div class="app-card overflow-hidden">
          <div class="px-5 py-3.5 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div class="flex items-center gap-2.5">
              <div class="w-6 h-6 rounded-md bg-emerald-50 text-emerald-600 flex items-center justify-center border border-emerald-200/60">
                <AppIcon name="component" size="xs" />
              </div>
              <div>
                <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Hardware Components Library</h2>
                <p class="text-3xs text-slate-400">Electronic parts, semiconductors, sensors, and passives</p>
              </div>
            </div>
            <router-link
              to="/components"
              class="text-xs font-semibold text-brand-600 hover:text-brand-700 flex items-center gap-1 group"
            >
              <span>View Library ({{ store.totalComponentsCount }})</span>
              <AppIcon name="chevron-right" size="xs" class="group-hover:translate-x-0.5 transition-transform" />
            </router-link>
          </div>

          <div class="overflow-x-auto">
            <table class="dense-table">
              <thead>
                <tr>
                  <th>Part Number (MPN)</th>
                  <th>Component Name</th>
                  <th>Category</th>
                  <th>Manufacturer</th>
                  <th>Compliance</th>
                  <th>Status</th>
                  <th class="text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="c in recentComponents"
                  :key="c.id"
                  @click="$router.push(`/components/${c.id || c.partNumber}`)"
                  class="cursor-pointer group"
                >
                  <td>
                    <span class="font-mono font-bold text-slate-900 group-hover:text-brand-600 transition-colors">
                      {{ c.partNumber }}
                    </span>
                  </td>
                  <td>
                    <div class="text-slate-700 truncate max-w-xs font-medium">
                      {{ c.name }}
                    </div>
                  </td>
                  <td>
                    <span class="px-2 py-0.5 rounded text-3xs font-medium bg-slate-100 text-slate-700 border border-slate-200/80">
                      {{ c.category }}
                    </span>
                  </td>
                  <td>
                    <span class="text-slate-800 font-semibold text-xs">
                      {{ c.manufacturer }}
                    </span>
                  </td>
                  <td>
                    <span
                      v-if="c.rohsCompliant || c.roHsCompliant"
                      class="px-1.5 py-0.5 rounded text-3xs font-mono font-bold bg-emerald-50 text-emerald-700 border border-emerald-200"
                      title="RoHS Compliant"
                    >
                      RoHS
                    </span>
                    <span v-else class="text-3xs text-slate-400 font-mono">N/A</span>
                  </td>
                  <td>
                    <StatusBadge :status="c.status" />
                  </td>
                  <td class="text-right" @click.stop>
                    <router-link
                      :to="`/components/${c.id || c.partNumber}`"
                      class="btn-icon inline-flex"
                      title="View Component Detail"
                    >
                      <AppIcon name="eye" size="xs" />
                    </router-link>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- RIGHT 4 COLS: Activity Timeline, Quick Launchpad & Data Health -->
      <div class="lg:col-span-4 space-y-4">
        <!-- A. QUICK ACTIONS LAUNCHPAD -->
        <div class="bg-white rounded-xl border border-slate-200/90 p-4 shadow-2xs">
          <div class="flex items-center justify-between pb-2.5 border-b border-slate-100 mb-3">
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
              <AppIcon name="layers" size="xs" class="text-brand-600" />
              <span>Quick Actions</span>
            </h2>
            <span class="text-3xs font-mono text-slate-400">Shortcuts</span>
          </div>

          <div class="grid grid-cols-2 gap-2.5">
            <router-link
              to="/products/create"
              class="p-2.5 rounded-lg border border-slate-200/80 hover:border-brand-300 bg-slate-50/60 hover:bg-brand-50/40 transition-all duration-150 group flex items-start gap-2.5"
            >
              <div class="w-7 h-7 rounded-lg bg-brand-50 text-brand-600 flex items-center justify-center shrink-0 border border-brand-200/60 group-hover:scale-105 transition-transform">
                <AppIcon name="plus" size="xs" />
              </div>
              <div class="min-w-0">
                <div class="text-xs font-bold text-slate-800 group-hover:text-brand-600 truncate">Add Product</div>
                <div class="text-3xs text-slate-400 truncate">New Master</div>
              </div>
            </router-link>

            <router-link
              to="/components/create"
              class="p-2.5 rounded-lg border border-slate-200/80 hover:border-emerald-300 bg-slate-50/60 hover:bg-emerald-50/40 transition-all duration-150 group flex items-start gap-2.5"
            >
              <div class="w-7 h-7 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0 border border-emerald-200/60 group-hover:scale-105 transition-transform">
                <AppIcon name="component" size="xs" />
              </div>
              <div class="min-w-0">
                <div class="text-xs font-bold text-slate-800 group-hover:text-emerald-700 truncate">Add Component</div>
                <div class="text-3xs text-slate-400 truncate">Specs & Docs</div>
              </div>
            </router-link>

            <router-link
              to="/manufacturers"
              class="p-2.5 rounded-lg border border-slate-200/80 hover:border-amber-300 bg-slate-50/60 hover:bg-amber-50/40 transition-all duration-150 group flex items-start gap-2.5"
            >
              <div class="w-7 h-7 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 border border-amber-200/60 group-hover:scale-105 transition-transform">
                <AppIcon name="manufacturer" size="xs" />
              </div>
              <div class="min-w-0">
                <div class="text-xs font-bold text-slate-800 group-hover:text-amber-700 truncate">Manufacturers</div>
                <div class="text-3xs text-slate-400 truncate">Fab Directory</div>
              </div>
            </router-link>

            <router-link
              to="/suppliers"
              class="p-2.5 rounded-lg border border-slate-200/80 hover:border-cyan-300 bg-slate-50/60 hover:bg-cyan-50/40 transition-all duration-150 group flex items-start gap-2.5"
            >
              <div class="w-7 h-7 rounded-lg bg-cyan-50 text-cyan-600 flex items-center justify-center shrink-0 border border-cyan-200/60 group-hover:scale-105 transition-transform">
                <AppIcon name="supplier" size="xs" />
              </div>
              <div class="min-w-0">
                <div class="text-xs font-bold text-slate-800 group-hover:text-cyan-700 truncate">Suppliers</div>
                <div class="text-3xs text-slate-400 truncate">Distributors</div>
              </div>
            </router-link>
          </div>
        </div>

        <!-- B. RECENT ENGINEERING ACTIVITY AUDIT -->
        <div class="bg-white rounded-xl border border-slate-200/90 p-4 shadow-2xs">
          <div class="flex items-center justify-between pb-2.5 border-b border-slate-100 mb-3">
            <div class="flex items-center gap-2">
              <div class="w-5 h-5 rounded bg-brand-50 text-brand-600 flex items-center justify-center">
                <AppIcon name="activity" size="xs" />
              </div>
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Live Engineering Audit</h2>
            </div>
            <span class="text-3xs font-mono font-semibold text-emerald-600 bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-200 flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
              <span>Realtime</span>
            </span>
          </div>

          <div class="space-y-2 max-h-[340px] overflow-y-auto pr-0.5">
            <div
              v-for="act in store.activities.slice(0, 6)"
              :key="act.id"
              class="p-2.5 rounded-lg bg-slate-50/80 border border-slate-150/70 hover:bg-slate-100/70 transition-all flex items-start gap-2.5 group"
            >
              <!-- Avatar or Initial -->
              <div class="relative shrink-0 mt-0.5">
                <img
                  v-if="act.userAvatar"
                  :src="act.userAvatar"
                  alt="Avatar"
                  class="w-6 h-6 rounded-full object-cover ring-1 ring-slate-300"
                />
                <div
                  v-else
                  :class="[
                    'w-6 h-6 rounded-full flex items-center justify-center text-3xs font-bold text-white',
                    act.badgeColor === 'emerald' ? 'bg-emerald-600' :
                    act.badgeColor === 'amber' ? 'bg-amber-600' :
                    act.badgeColor === 'rose' ? 'bg-rose-600' :
                    act.badgeColor === 'purple' ? 'bg-purple-600' : 'bg-brand-600'
                  ]"
                >
                  {{ (act.user || 'S').charAt(0) }}
                </div>
              </div>

              <!-- Activity Content -->
              <div class="min-w-0 flex-1">
                <div class="flex items-center justify-between gap-1.5 mb-0.5">
                  <span class="text-xs font-bold text-slate-900 truncate">{{ act.user || 'System' }}</span>
                  <span class="text-3xs font-mono text-slate-400 shrink-0">{{ formatTime(act.timestamp) }}</span>
                </div>

                <div class="text-xs text-slate-700 leading-snug">
                  <span class="text-slate-600">{{ act.action }}: </span>
                  <span class="font-mono font-semibold text-brand-600">{{ act.targetName || act.entityName || act.targetId }}</span>
                </div>

                <div class="flex items-center gap-1.5 mt-1 text-3xs font-mono text-slate-400">
                  <span class="px-1.5 py-0.2 rounded bg-slate-200/70 text-slate-600 font-semibold">{{ act.targetType || 'Log' }}</span>
                  <span v-if="act.targetId" class="truncate">ID: {{ act.targetId }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- C. MASTER DATA & PLATFORM INTEGRITY -->
        <div class="bg-gradient-to-br from-slate-900 via-navy-850 to-slate-900 rounded-xl p-4 text-white shadow-md border border-slate-800 space-y-2.5 relative overflow-hidden">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <AppIcon name="shield" size="xs" class="text-emerald-400" />
              <h3 class="text-xs font-bold tracking-wide">Platform & Master Data Integrity</h3>
            </div>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-emerald-950 text-emerald-300 border border-emerald-800">
              100% AUDITED
            </span>
          </div>

          <p class="text-2xs text-slate-300 leading-relaxed">
            All registered products (Trading, Project, Manufacture & Service), components, suppliers, and change orders comply with enterprise PLM standards across Phases 1–6.
          </p>

          <div class="grid grid-cols-2 gap-2 pt-2 font-mono text-3xs text-slate-400 border-t border-slate-800">
            <div>
              <span class="block text-slate-400">RoHS Compliance</span>
              <span class="text-emerald-400 font-bold text-xs">100% Verified</span>
            </div>
            <div>
              <span class="block text-slate-400">Security & RBAC</span>
              <span class="text-brand-400 font-bold text-xs">Enforced</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import KpiCard from '@/components/common/KpiCard.vue';
import StatusBadge from '@/components/common/StatusBadge.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'DashboardView',
  components: {
    KpiCard,
    StatusBadge,
    AppIcon
  },
  computed: {
    store() {
      return useMasterStore();
    },
    currentUserName() {
      return this.store.currentUser ? this.store.currentUser.name : 'Engineer';
    },
    recentProducts() {
      return this.store.products.slice(0, 5);
    },
    recentComponents() {
      return this.store.components.slice(0, 5);
    }
  },
  methods: {
    formatTimestamp(ts) {
      if (!ts) return 'Recent';
      if (ts.length > 10) {
        return ts.substring(0, 10);
      }
      return ts;
    },
    formatTime(ts) {
      if (!ts) return 'Just now';
      if (ts.includes('hours ago')) {
        const num = ts.split(' ')[0];
        return `${num}h ago`;
      }
      if (ts.includes('hour ago')) {
        return '1h ago';
      }
      if (ts.includes('mins ago') || ts.includes('min ago')) {
        const num = ts.split(' ')[0];
        return `${num}m ago`;
      }
      if (ts.includes('days ago') || ts.includes('day ago')) {
        const num = ts.split(' ')[0];
        return `${num}d ago`;
      }
      return ts;
    }
  }
};
</script>
