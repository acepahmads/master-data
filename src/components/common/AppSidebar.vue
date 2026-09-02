<template>
  <div>
    <!-- Mobile/Tablet Backdrop Overlay -->
    <transition
      enter-active-class="transition-opacity duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="store.mobileSidebarOpen"
        @click="store.closeMobileSidebar"
        class="fixed inset-0 bg-slate-950/70 backdrop-blur-xs z-40 lg:hidden"
      ></div>
    </transition>

    <!-- Sidebar Element -->
    <aside
      :class="[
        'bg-[#0b1120] text-slate-300 flex flex-col shrink-0 border-r border-slate-800/80 transition-all duration-300 select-none shadow-xl h-full',
        // Desktop positioning & width
        'lg:static lg:translate-x-0',
        collapsed ? 'lg:w-16' : 'lg:w-64',
        // Mobile / Tablet off-canvas drawer
        'fixed inset-y-0 left-0 z-50 w-72 max-w-[85vw]',
        store.mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'
      ]"
    >
      <!-- Top Brand Header -->
      <div class="h-14 flex items-center justify-between px-3.5 border-b border-slate-800/60 bg-[#070b15]/60 backdrop-blur-sm">
        <router-link to="/" class="flex items-center gap-2.5 overflow-hidden group">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-brand-600 to-cyan-500 flex items-center justify-center text-white shadow-md shadow-brand-500/20 group-hover:shadow-glow-brand transition-all duration-200 shrink-0">
            <AppIcon name="component" size="sm" class="text-white" />
          </div>
          <div v-if="!collapsed || store.mobileSidebarOpen" class="flex flex-col truncate">
            <div class="flex items-center gap-1.5">
              <span class="font-bold text-xs text-white tracking-tight">Product & Master</span>
              <span class="text-3xs font-bold px-1.5 py-0.2 rounded bg-brand-950 text-brand-300 border border-brand-700/60 font-mono">v1.0.0</span>
            </div>
            <span class="text-3xs text-slate-400 font-medium tracking-wider uppercase">Control Center</span>
          </div>
        </router-link>

        <!-- Desktop Collapse / Expand Toggle Button -->
        <button
          @click="toggleCollapse"
          class="hidden lg:flex p-1 rounded-md text-slate-400 hover:text-white hover:bg-slate-800/70 focus:outline-none transition-all duration-150"
          :title="collapsed ? 'Expand Sidebar' : 'Collapse Sidebar'"
        >
          <AppIcon :name="collapsed ? 'chevron-right' : 'chevron-left'" size="xs" />
        </button>

        <!-- Mobile Close Drawer Button -->
        <button
          @click="store.closeMobileSidebar"
          class="lg:hidden p-1.5 rounded-md text-slate-400 hover:text-white hover:bg-slate-800/70 focus:outline-none transition-all"
          title="Close Navigation"
        >
          <AppIcon name="x" size="sm" />
        </button>
      </div>

    <!-- Navigation Groups -->
    <div class="flex-1 overflow-y-auto overflow-x-hidden py-3.5 px-2.5 space-y-4">
      <!-- CORE SECTION -->
      <div>
        <div v-if="!collapsed" class="px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300">
          Core Overview
        </div>
        <router-link
          to="/"
          exact
          active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
          class="flex items-center gap-2.5 px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group relative"
          :title="collapsed ? 'Dashboard' : ''"
        >
          <AppIcon name="dashboard" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
          <span v-if="!collapsed" class="truncate">Dashboard</span>
        </router-link>

        <router-link
          to="/progress"
          active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
          class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group mt-1"
          :title="collapsed ? 'Architecture & Progress' : ''"
        >
          <div class="flex items-center gap-2.5 min-w-0">
            <AppIcon name="roadmap" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
            <span v-if="!collapsed" class="truncate">Platform Architecture</span>
          </div>
          <span v-if="!collapsed" class="text-3xs font-mono font-bold px-1.5 py-0.5 rounded bg-emerald-950 text-emerald-300 border border-emerald-800/80">
            100%
          </span>
        </router-link>
      </div>

      <!-- PRODUCT & MASTER DATA GROUP (COLLAPSIBLE) -->
      <div>
        <div
          v-if="!collapsed"
          @click="masterDataExpanded = !masterDataExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Product & Master Data</span>
          <AppIcon :name="masterDataExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="masterDataExpanded || collapsed" class="space-y-1">
          <!-- Products -->
          <router-link
            to="/products"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Products' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="product" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Products</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalProductsCount }}
            </span>
          </router-link>

          <!-- Components -->
          <router-link
            to="/components"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Components Library' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="component" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Components</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalComponentsCount }}
            </span>
          </router-link>

          <!-- Component Categories -->
          <router-link
            to="/categories"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Component Categories' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="tree" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Categories</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalCategoriesCount }}
            </span>
          </router-link>

          <!-- Manufacturers -->
          <router-link
            to="/manufacturers"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Manufacturers' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="manufacturer" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Manufacturers</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalManufacturersCount }}
            </span>
          </router-link>

          <!-- Suppliers -->
          <router-link
            to="/suppliers"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Suppliers' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="supplier" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Suppliers</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalSuppliersCount }}
            </span>
          </router-link>

          <!-- Units of Measure -->
          <router-link
            to="/units"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Units of Measurement' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="unit" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Units (UoM)</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalUnitsCount }}
            </span>
          </router-link>

          <!-- Data Import & Migration (Cross-Cutting Platform Foundation) -->
          <router-link
            to="/data-import"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Data Import & Migration' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="import" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Data Import</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-purple-950 text-purple-300 border border-purple-800/60">
              2026
            </span>
          </router-link>
        </div>
      </div>

      <!-- PRODUCT LIFECYCLE GROUP -->
      <div>
        <div
          v-if="!collapsed"
          @click="lifecycleExpanded = !lifecycleExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Product Versions</span>
          <AppIcon :name="lifecycleExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="lifecycleExpanded || collapsed" class="space-y-1">
          <!-- Product Versions -->
          <router-link
            to="/versions"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Product Versions' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="layers" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Product Versions</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ versionStore.totalVersionsCount }}
            </span>
          </router-link>

          <!-- Hardware Revisions -->
          <router-link
            to="/revisions"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Hardware Revisions' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="cpu" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Hardware Revisions</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ versionStore.totalRevisionsCount }}
            </span>
          </router-link>
        </div>
      </div>

      <!-- PROTOTYPE BUILDS GROUP -->
      <div>
        <div
          v-if="!collapsed"
          @click="prototypeExpanded = !prototypeExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Prototype Builds</span>
          <AppIcon :name="prototypeExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="prototypeExpanded || collapsed" class="space-y-1">
          <!-- Command Center -->
          <router-link
            to="/prototypes"
            exact
            active-class="bg-indigo-500/20 text-indigo-300 font-semibold border-l-2 border-indigo-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Prototype Command Center' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="cpu" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-indigo-500\/20]:text-indigo-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Build Dashboard</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-indigo-950 text-indigo-300 border border-indigo-700/60 font-bold">
              {{ protoStore.totalPrototypesCount }}
            </span>
          </router-link>

          <!-- All Builds Registry -->
          <router-link
            to="/prototypes/list"
            active-class="bg-indigo-500/20 text-indigo-300 font-semibold border-l-2 border-indigo-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'All Prototype Builds' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="list" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-indigo-500\/20]:text-indigo-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Builds Registry</span>
            </div>
          </router-link>
        </div>
      </div>

      <!-- TESTING & VALIDATION GROUP -->
      <div>
        <div
          v-if="!collapsed"
          @click="testingExpanded = !testingExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Testing & Validation</span>
          <AppIcon :name="testingExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="testingExpanded || collapsed" class="space-y-1">
          <!-- Testing Dashboard -->
          <router-link
            to="/testing"
            exact
            active-class="bg-emerald-500/20 text-emerald-300 font-semibold border-l-2 border-emerald-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Testing Command Center' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="shield" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-emerald-500\/20]:text-emerald-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Testing Center</span>
            </div>
          </router-link>

          <!-- Master Test Plans -->
          <router-link
            to="/test-plans"
            active-class="bg-emerald-500/20 text-emerald-300 font-semibold border-l-2 border-emerald-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Master Test Plans' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="layers" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-emerald-500\/20]:text-emerald-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Test Protocols</span>
            </div>
          </router-link>

          <!-- Engineering Findings Hub -->
          <router-link
            to="/findings"
            active-class="bg-emerald-500/20 text-emerald-300 font-semibold border-l-2 border-emerald-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Defects & Findings' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="alert-triangle" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-emerald-500\/20]:text-emerald-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Defects & Findings</span>
            </div>
          </router-link>
        </div>
      </div>

      <!-- ENGINEERING CHANGES GROUP -->
      <div>
        <div
          v-if="!collapsed"
          @click="changesExpanded = !changesExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Engineering Changes</span>
          <AppIcon :name="changesExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="changesExpanded || collapsed" class="space-y-1">
          <!-- Changes Dashboard -->
          <router-link
            to="/changes"
            exact
            active-class="bg-cyan-500/20 text-cyan-300 font-semibold border-l-2 border-cyan-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Change Control Center' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="dashboard" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-cyan-500\/20]:text-cyan-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Change Control</span>
            </div>
          </router-link>

          <!-- ECR Registry -->
          <router-link
            to="/ecr"
            active-class="bg-cyan-500/20 text-cyan-300 font-semibold border-l-2 border-cyan-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'ECR Requests' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="file" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-cyan-500\/20]:text-cyan-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">ECR Requests</span>
            </div>
          </router-link>

          <!-- ECO Change Orders -->
          <router-link
            to="/eco"
            active-class="bg-cyan-500/20 text-cyan-300 font-semibold border-l-2 border-cyan-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'ECO Change Orders' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="layers" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-cyan-500\/20]:text-cyan-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">ECO Orders</span>
            </div>
          </router-link>

          <!-- Traceability Matrix -->
          <router-link
            to="/changes/traceability"
            active-class="bg-cyan-500/20 text-cyan-300 font-semibold border-l-2 border-cyan-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Traceability Matrix' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="roadmap" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-cyan-500\/20]:text-cyan-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Traceability Matrix</span>
            </div>
          </router-link>
        </div>
      </div>

      <!-- ADMINISTRATION GROUP -->
      <div>
        <div
          v-if="!collapsed"
          @click="adminExpanded = !adminExpanded"
          class="flex items-center justify-between px-2.5 mb-1.5 text-3xs font-bold uppercase tracking-widest text-slate-300 cursor-pointer hover:text-white transition-colors"
        >
          <span>Administration</span>
          <AppIcon :name="adminExpanded ? 'chevron-down' : 'chevron-right'" size="xs" class="text-slate-300" />
        </div>

        <div v-show="adminExpanded || collapsed" class="space-y-1">
          <!-- Users -->
          <router-link
            to="/users"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'User Directory' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="users" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">User Directory</span>
            </div>
            <span v-if="!collapsed" class="text-3xs font-mono px-1.5 py-0.5 rounded bg-slate-800/90 text-slate-400 border border-slate-700/50 group-[.bg-brand-500\/15]:bg-brand-950 group-[.bg-brand-500\/15]:text-brand-300 group-[.bg-brand-500\/15]:border-brand-700/60">
              {{ store.totalUsersCount }}
            </span>
          </router-link>

          <!-- Roles & Permissions -->
          <router-link
            to="/roles"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'Roles & RBAC Matrix' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="shield" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">Roles & RBAC</span>
            </div>
          </router-link>

          <!-- API Access Tokens -->
          <router-link
            to="/api-tokens"
            active-class="bg-brand-500/15 text-brand-300 font-semibold border-l-2 border-brand-500 shadow-sm"
            class="flex items-center justify-between px-2.5 py-2 rounded-md text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all duration-150 group"
            :title="collapsed ? 'API Access Tokens' : ''"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <AppIcon name="key" size="sm" class="text-slate-400 group-hover:text-white group-[.bg-brand-500\/15]:text-brand-400 transition-colors" />
              <span v-if="!collapsed" class="truncate">API Access Tokens</span>
            </div>
            <span v-if="!collapsed" class="text-4xs font-mono font-bold px-1.5 py-0.2 rounded bg-purple-900/60 text-purple-300 border border-purple-700/50">
              API
            </span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Bottom Utility & Profile Area -->
    <div class="p-2.5 border-t border-slate-800/70 bg-[#070b15]/70 backdrop-blur-sm space-y-1">
      <!-- Settings Modal Trigger -->
      <button
        @click="openSettings"
        class="w-full flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-xs text-slate-400 hover:text-white hover:bg-slate-800/50 transition-colors"
        :title="collapsed ? 'Settings' : ''"
      >
        <AppIcon name="settings" size="sm" />
        <span v-if="!collapsed" class="truncate">System Settings</span>
      </button>

      <!-- Help / Guide Trigger -->
      <button
        @click="openHelp"
        class="w-full flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-xs text-slate-400 hover:text-white hover:bg-slate-800/50 transition-colors"
        :title="collapsed ? 'Help & Documentation' : ''"
      >
        <AppIcon name="help" size="sm" />
        <span v-if="!collapsed" class="truncate">Documentation & API</span>
      </button>

      <!-- User Profile & Quick Switcher Footer -->
      <div v-if="(!collapsed || store.mobileSidebarOpen) && store.currentUser" class="pt-2 mt-1 border-t border-slate-800/60">
        <div class="flex items-center justify-between gap-2 px-2 py-1.5 rounded-lg bg-slate-850/60 border border-slate-800/80 hover:border-slate-700/80 transition-colors group">
          <div class="flex items-center gap-2.5 min-w-0 flex-1">
            <div class="relative shrink-0">
              <img
                :src="store.currentUser.avatar || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
                alt="User Avatar"
                class="w-7 h-7 rounded-full object-cover ring-1 ring-slate-700"
              />
              <span class="absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full bg-emerald-500 ring-2 ring-[#0b1120]"></span>
            </div>
            <div class="min-w-0 flex-1">
              <div class="text-xs font-semibold text-white truncate leading-tight">{{ store.currentUser.name }}</div>
              <div class="text-3xs text-brand-400 font-medium truncate">{{ userRoleName }}</div>
            </div>
          </div>
          <button
            @click="handleLogout"
            class="p-1 rounded-md text-slate-400 hover:text-rose-400 hover:bg-slate-800/80 transition-colors shrink-0"
            title="Sign Out of Workstation"
          >
            <AppIcon name="logout" size="xs" />
          </button>
        </div>
      </div>
    </div>
  </aside>
</div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import { useAuthStore } from '@/stores/authStore';
import { useVersionStore } from '@/stores/versionStore';
import { usePrototypeStore } from '@/stores/prototypeStore';
import AppIcon from './AppIcon.vue';

export default {
  name: 'AppSidebar',
  components: { AppIcon },
  data() {
    return {
      masterDataExpanded: false,
      lifecycleExpanded: false,
      prototypeExpanded: false,
      testingExpanded: false,
      changesExpanded: false,
      adminExpanded: false
    };
  },
  created() {
    this.syncExpandedGroups(this.$route ? this.$route.path : '/');
  },
  watch: {
    $route(to) {
      if (this.store.mobileSidebarOpen) {
        this.store.closeMobileSidebar();
      }
      this.syncExpandedGroups(to.path);
    }
  },
  computed: {
    store() {
      return useMasterStore();
    },
    versionStore() {
      return useVersionStore();
    },
    protoStore() {
      return usePrototypeStore();
    },
    collapsed() {
      return this.store.sidebarCollapsed;
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
    syncExpandedGroups(path) {
      if (!path) path = this.$route ? this.$route.path : '/';

      this.masterDataExpanded = /^\/(products|components|categories|manufacturers|suppliers|units|data-import)/.test(path);
      this.lifecycleExpanded = /^\/(versions|revisions)/.test(path);
      this.prototypeExpanded = /^\/prototypes/.test(path);
      this.testingExpanded = /^\/(testing|test-plans|test-sessions|findings)/.test(path);
      this.changesExpanded = /^\/(changes|ecr|eco)/.test(path);
      this.adminExpanded = /^\/(users|roles)/.test(path);
    },
    toggleCollapse() {
      this.store.sidebarCollapsed = !this.store.sidebarCollapsed;
    },
    openSettings() {
      this.store.settingsModalOpen = true;
    },
    openHelp() {
      this.store.helpModalOpen = true;
    },
    handleLogout() {
      const authStore = useAuthStore();
      authStore.logout();
      this.store.logout();
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
</style>
