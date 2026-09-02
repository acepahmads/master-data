<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
    <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs transition-opacity" @click="close"></div>

    <div class="flex min-h-full items-start justify-center p-4 pt-16 text-center">
      <div class="relative transform overflow-hidden rounded-xl bg-white text-left shadow-2xl transition-all border border-slate-200 w-full max-w-2xl">
        <!-- Search Input Bar -->
        <div class="flex items-center px-4 py-3 border-b border-slate-200 gap-3 bg-slate-50/50">
          <AppIcon name="search" size="md" class="text-slate-400 shrink-0" />
          <input
            ref="searchInput"
            v-model="query"
            type="text"
            placeholder="Quick search products, components (MPN), manufacturers, categories..."
            class="w-full bg-transparent border-none text-sm text-slate-800 focus:outline-none placeholder-slate-400 font-medium"
            @keydown.esc="close"
          />
          <kbd class="font-mono text-3xs px-2 py-1 rounded bg-slate-200/80 text-slate-600 border border-slate-300 font-semibold shrink-0">
            ESC
          </kbd>
        </div>

        <!-- Filter category tabs -->
        <div class="flex items-center gap-1.5 px-4 py-2 border-b border-slate-100 bg-slate-50/30 text-2xs overflow-x-auto">
          <button
            v-for="t in filterTabs"
            :key="t.key"
            @click="activeTab = t.key"
            :class="[
              'px-2.5 py-1 rounded-md font-medium transition-colors whitespace-nowrap',
              activeTab === t.key ? 'bg-brand-600 text-white' : 'text-slate-600 hover:bg-slate-200/60'
            ]"
          >
            {{ t.label }} ({{ t.count }})
          </button>
        </div>

        <!-- Search Results List -->
        <div class="max-h-96 overflow-y-auto p-2 divide-y divide-slate-100">
          <div v-if="filteredResults.length === 0" class="py-10 text-center text-slate-400">
            <AppIcon name="search" size="lg" class="mx-auto text-slate-300 mb-2" />
            <p class="text-xs">No matching master data found for "{{ query }}"</p>
          </div>

          <div
            v-for="item in filteredResults"
            :key="item.type + '-' + item.id"
            @click="selectItem(item)"
            class="p-2.5 hover:bg-slate-50 rounded-lg cursor-pointer transition-colors flex items-center justify-between group"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div :class="['w-8 h-8 rounded-md flex items-center justify-center shrink-0 text-white shadow-2xs', item.iconBg]">
                <AppIcon :name="item.icon" size="sm" />
              </div>
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="text-xs font-bold text-slate-900 truncate group-hover:text-brand-600">{{ item.title }}</span>
                  <span class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-100 text-slate-600 border border-slate-200">{{ item.code }}</span>
                </div>
                <div class="text-2xs text-slate-500 truncate mt-0.5">{{ item.subtitle }}</div>
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <span class="text-3xs font-medium uppercase tracking-wider text-slate-400 bg-slate-100 px-2 py-0.5 rounded">{{ item.type }}</span>
              <AppIcon name="chevron-right" size="xs" class="text-slate-300 group-hover:text-brand-600 transition-transform group-hover:translate-x-0.5" />
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-4 py-2 border-t border-slate-200 bg-slate-50 flex items-center justify-between text-3xs text-slate-500">
          <div class="flex items-center gap-3">
            <span>Navigate with click</span>
            <span>ESC to close</span>
          </div>
          <span class="font-mono">Total Results: {{ filteredResults.length }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import AppIcon from './AppIcon.vue';

export default {
  name: 'GlobalSearchModal',
  components: { AppIcon },
  data() {
    return {
      query: '',
      activeTab: 'all'
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    isOpen() {
      return this.store.globalSearchOpen;
    },
    allResults() {
      const q = this.query.toLowerCase().trim();
      const results = [];

      // Products
      this.store.products.forEach(p => {
        if (!q || p.name.toLowerCase().includes(q) || p.code.toLowerCase().includes(q) || p.category.toLowerCase().includes(q)) {
          results.push({
            type: 'product',
            id: p.id,
            title: p.name,
            code: p.code,
            subtitle: `${p.category} • Status: ${p.status}`,
            icon: 'product',
            iconBg: 'bg-blue-600',
            route: `/products/${p.id}`
          });
        }
      });

      // Components
      this.store.components.forEach(c => {
        if (!q || c.partNumber.toLowerCase().includes(q) || c.name.toLowerCase().includes(q) || c.manufacturer.toLowerCase().includes(q) || c.internalCode.toLowerCase().includes(q)) {
          results.push({
            type: 'component',
            id: c.id,
            title: c.partNumber,
            code: c.internalCode,
            subtitle: `${c.name} • Mfg: ${c.manufacturer}`,
            icon: 'component',
            iconBg: 'bg-emerald-600',
            route: `/components/${c.id}`
          });
        }
      });

      // Manufacturers
      this.store.manufacturers.forEach(m => {
        if (!q || m.name.toLowerCase().includes(q) || m.code.toLowerCase().includes(q) || m.country.toLowerCase().includes(q)) {
          results.push({
            type: 'manufacturer',
            id: m.id,
            title: m.name,
            code: m.code,
            subtitle: `${m.country} • ${m.totalComponents} Components`,
            icon: 'manufacturer',
            iconBg: 'bg-purple-600',
            route: `/manufacturers`
          });
        }
      });

      // Categories
      this.store.flatCategories.forEach(c => {
        if (!q || c.rawName.toLowerCase().includes(q) || c.code.toLowerCase().includes(q)) {
          results.push({
            type: 'category',
            id: c.id,
            title: c.rawName,
            code: c.code,
            subtitle: `Category Branch`,
            icon: 'tree',
            iconBg: 'bg-amber-600',
            route: `/categories`
          });
        }
      });

      return results;
    },
    filterTabs() {
      const q = this.query.toLowerCase().trim();
      const all = this.allResults;
      return [
        { key: 'all', label: 'All Results', count: all.length },
        { key: 'product', label: 'Products', count: all.filter(r => r.type === 'product').length },
        { key: 'component', label: 'Components', count: all.filter(r => r.type === 'component').length },
        { key: 'manufacturer', label: 'Manufacturers', count: all.filter(r => r.type === 'manufacturer').length },
        { key: 'category', label: 'Categories', count: all.filter(r => r.type === 'category').length }
      ];
    },
    filteredResults() {
      if (this.activeTab === 'all') return this.allResults;
      return this.allResults.filter(r => r.type === this.activeTab);
    }
  },
  watch: {
    isOpen(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          if (this.$refs.searchInput) {
            this.$refs.searchInput.focus();
          }
        });
      }
    }
  },
  mounted() {
    window.addEventListener('keydown', this.handleKeyDown);
  },
  beforeDestroy() {
    window.removeEventListener('keydown', this.handleKeyDown);
  },
  methods: {
    close() {
      this.store.globalSearchOpen = false;
      this.query = '';
    },
    handleKeyDown(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        this.store.globalSearchOpen = !this.store.globalSearchOpen;
      }
    },
    selectItem(item) {
      this.close();
      if (this.$route.path !== item.route) {
        this.$router.push(item.route);
      }
    }
  }
};
</script>
