<template>
  <div class="space-y-5">
    <!-- Top Action Bar -->
    <div class="app-card p-4 bg-white flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-cyan-50 border border-cyan-200 text-cyan-700 flex items-center justify-center shadow-2xs">
          <AppIcon name="layers" size="sm" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-slate-900">Project Composition & Architecture Workspace</h2>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-cyan-50 text-cyan-700 border border-cyan-200">
              PROJECT PRODUCT
            </span>
          </div>
          <p class="text-3xs text-slate-500 mt-0.5">Solution package hierarchy referencing R&D products, commercial hardware, and components</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button @click="openAddModal(null, 'SUB_ASSEMBLY')" class="btn btn-secondary">
          <AppIcon name="folder-plus" size="xs" class="mr-1.5 text-slate-500" />
          <span>+ Add Sub-Assembly</span>
        </button>
        <button @click="openAddModal(null, 'PRODUCT')" class="btn btn-primary">
          <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
          <span>+ Add Item</span>
        </button>
      </div>
    </div>

    <!-- Main Workspace (2-Column Architecture Explorer) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
      <!-- LEFT PANEL: Interactive Tree Explorer (7 Cols) -->
      <div class="lg:col-span-7 space-y-3">
        <div class="app-card overflow-hidden bg-white">
          <!-- Tree Header with Filter -->
          <div class="px-4 py-3 bg-slate-50 border-b border-slate-200 flex items-center justify-between gap-3">
            <div class="relative flex-1 max-w-xs">
              <AppIcon name="search" size="xs" class="absolute left-2.5 top-2.5 text-slate-400" />
              <input
                v-model="treeSearch"
                type="text"
                placeholder="Filter structure items..."
                class="form-input pl-7.5 py-1 text-xs"
              />
            </div>
            <div class="flex items-center gap-2">
              <span class="text-3xs font-mono text-slate-500">
                {{ flatItems.length }} Total Elements
              </span>
              <button
                @click="toggleAll"
                class="text-3xs font-semibold text-brand-600 hover:text-brand-800 px-2 py-1 rounded hover:bg-slate-200/60"
              >
                {{ allExpanded ? 'Collapse All' : 'Expand All' }}
              </button>
            </div>
          </div>

          <!-- Financial Summary Strip -->
          <div v-if="flatItems.length > 0" class="p-3 bg-gradient-to-r from-slate-50 via-purple-50/20 to-blue-50/30 border-b border-slate-200 grid grid-cols-3 gap-2 text-xs">
            <div class="p-2 rounded-lg bg-white border border-slate-200/80 shadow-2xs">
              <div class="text-4xs font-bold uppercase tracking-wider text-slate-400">Total Harga Jual (Revenue)</div>
              <div class="text-xs font-black font-mono text-emerald-700 mt-0.5">{{ formatCurrency(totalProjectRevenue) }}</div>
            </div>
            <div class="p-2 rounded-lg bg-white border border-slate-200/80 shadow-2xs">
              <div class="text-4xs font-bold uppercase tracking-wider text-slate-400">Total Biaya/Modal (Cost)</div>
              <div class="text-xs font-black font-mono text-amber-700 mt-0.5">{{ formatCurrency(totalProjectCost) }}</div>
            </div>
            <div class="p-2 rounded-lg bg-white border border-slate-200/80 shadow-2xs">
              <div class="text-4xs font-bold uppercase tracking-wider text-slate-400">Est. Gross Margin</div>
              <div class="text-xs font-black font-mono text-indigo-700 mt-0.5">{{ projectMarginPercent }}% ({{ formatCurrency(projectGrossProfit) }})</div>
            </div>
          </div>

          <!-- Empty Tree State -->
          <div v-if="flatItems.length === 0" class="p-8 text-center">
            <div class="w-12 h-12 rounded-2xl bg-cyan-50 text-cyan-600 border border-cyan-200 flex items-center justify-center mx-auto mb-3">
              <AppIcon name="layers" size="md" />
            </div>
            <div class="text-xs font-bold text-slate-800">No Structure Items Configured</div>
            <p class="text-3xs text-slate-500 max-w-sm mx-auto mt-1">
              Start building this solution package by adding Sub-Assemblies, referencing internal R&D Products, or attaching Trading Hardware.
            </p>
            <div class="flex items-center justify-center gap-2 mt-4">
              <button @click="openAddModal(null, 'SUB_ASSEMBLY')" class="btn btn-secondary btn-sm">
                + Add Sub-Assembly
              </button>
              <button @click="openAddModal(null, 'PRODUCT')" class="btn btn-primary btn-sm">
                + Add Product Reference
              </button>
            </div>
          </div>

          <!-- Hierarchical Tree Render -->
          <div v-else class="p-3 space-y-2">
            <div
              v-for="rootNode in rootTreeItems"
              :key="rootNode.id"
              class="space-y-1.5"
            >
              <!-- Tree Node Element -->
              <div
                @click="selectItem(rootNode)"
                :class="[
                  'p-2.5 rounded-lg border transition-all flex items-center justify-between cursor-pointer group select-none text-xs gap-2',
                  selectedItemId === rootNode.id
                    ? 'border-brand-500 bg-brand-50/40 ring-1 ring-brand-500/20'
                    : 'border-slate-200/90 bg-white hover:border-slate-300 hover:bg-slate-50/70'
                ]"
              >
                <div class="flex items-center gap-2.5 min-w-0 flex-1">
                  <button
                    v-if="hasChildren(rootNode.id)"
                    @click.stop="toggleExpand(rootNode.id)"
                    class="p-0.5 rounded hover:bg-slate-200 text-slate-400 hover:text-slate-700"
                  >
                    <AppIcon :name="isExpanded(rootNode.id) ? 'chevron-down' : 'chevron-right'" size="xs" />
                  </button>
                  <span v-else class="w-4"></span>

                  <!-- Item Type Icon -->
                  <div
                    :class="[
                      'w-7 h-7 rounded-lg flex items-center justify-center shrink-0 text-3xs font-bold shadow-2xs border',
                      rootNode.itemType === 'SUB_ASSEMBLY' ? 'bg-amber-50 text-amber-700 border-amber-200' :
                      rootNode.itemType === 'PRODUCT' ? (rootNode.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-brand-50 text-brand-700 border-brand-200') :
                      'bg-emerald-50 text-emerald-700 border-emerald-200'
                    ]"
                  >
                    <AppIcon
                      :name="rootNode.itemType === 'SUB_ASSEMBLY' ? 'folder' : rootNode.itemType === 'PRODUCT' ? 'product' : 'component'"
                      size="xs"
                    />
                  </div>

                  <div class="min-w-0 flex-1">
                    <div class="flex items-center gap-1.5">
                      <span class="font-bold text-slate-900 truncate">{{ rootNode.name }}</span>
                      <span
                        v-if="rootNode.productType"
                        :class="[
                          'text-3xs font-mono font-bold px-1.5 py-0.2 rounded border',
                          rootNode.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-emerald-50 text-emerald-700 border-emerald-200'
                        ]"
                      >
                        {{ rootNode.productType }}
                      </span>
                    </div>
                    <div class="text-3xs text-slate-400 font-mono mt-0.5 flex items-center gap-2">
                      <span v-if="rootNode.productCode">{{ rootNode.productCode }}</span>
                      <span v-else-if="rootNode.componentMpn">{{ rootNode.componentMpn }}</span>
                      <span>•</span>
                      <span>Qty: {{ rootNode.quantity }} {{ rootNode.unitName }}</span>
                    </div>
                  </div>
                </div>

                <!-- Price Information Column -->
                <div class="text-right font-mono text-3xs shrink-0 pl-2">
                  <div class="font-bold text-slate-900">
                    {{ formatCurrency((rootNode.sellingPrice || 0) * (rootNode.quantity || 1), rootNode.currency) }}
                  </div>
                  <div class="text-4xs text-slate-400 flex items-center justify-end gap-1 mt-0.5">
                    <span v-if="rootNode.costPrice > 0" class="text-amber-700">Beli: {{ formatCurrency(rootNode.costPrice, rootNode.currency) }}</span>
                    <span v-else class="text-slate-400">Beli: Rp 0</span>
                    <span>|</span>
                    <span class="text-emerald-700 font-medium">Jual: {{ formatCurrency(rootNode.sellingPrice, rootNode.currency) }}</span>
                  </div>
                </div>

                <div class="flex items-center gap-1 opacity-80 group-hover:opacity-100 shrink-0">
                  <button
                    v-if="rootNode.itemType === 'SUB_ASSEMBLY'"
                    @click.stop="openAddModal(rootNode.id, 'PRODUCT')"
                    class="p-1 rounded text-slate-400 hover:text-brand-600 hover:bg-white"
                    title="Add Item Under Sub-Assembly"
                  >
                    <AppIcon name="plus" size="xs" />
                  </button>
                  <button
                    @click.stop="confirmDelete(rootNode)"
                    class="p-1 rounded text-slate-400 hover:text-rose-600 hover:bg-white"
                    title="Remove Item"
                  >
                    <AppIcon name="trash" size="xs" />
                  </button>
                </div>
              </div>

              <!-- Nested Children (Level 2) -->
              <div
                v-if="hasChildren(rootNode.id) && isExpanded(rootNode.id)"
                class="pl-6 space-y-1.5 border-l-2 border-slate-200 ml-4"
              >
                <div
                  v-for="childNode in getChildren(rootNode.id)"
                  :key="childNode.id"
                  @click="selectItem(childNode)"
                  :class="[
                    'p-2 rounded-lg border transition-all flex items-center justify-between cursor-pointer group text-xs gap-2',
                    selectedItemId === childNode.id
                      ? 'border-brand-500 bg-brand-50/40 ring-1 ring-brand-500/20'
                      : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50/60'
                  ]"
                >
                  <div class="flex items-center gap-2 min-w-0 flex-1">
                    <div
                      :class="[
                        'w-6 h-6 rounded-md flex items-center justify-center shrink-0 text-3xs font-bold border',
                        childNode.itemType === 'PRODUCT' ? (childNode.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-brand-50 text-brand-700 border-brand-200') :
                        'bg-emerald-50 text-emerald-700 border-emerald-200'
                      ]"
                    >
                      <AppIcon :name="childNode.itemType === 'PRODUCT' ? 'product' : 'component'" size="xs" />
                    </div>

                    <div class="min-w-0 flex-1">
                      <div class="flex items-center gap-1.5">
                        <span class="font-semibold text-slate-900 truncate">{{ childNode.name }}</span>
                        <span
                          v-if="childNode.productType"
                          :class="[
                            'text-3xs font-mono font-bold px-1.5 py-0.2 rounded border',
                            childNode.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-emerald-50 text-emerald-700 border-emerald-200'
                          ]"
                        >
                          {{ childNode.productType }}
                        </span>
                      </div>
                      <div class="text-3xs text-slate-400 font-mono flex items-center gap-2">
                        <span v-if="childNode.productCode">{{ childNode.productCode }}</span>
                        <span v-else-if="childNode.componentMpn">{{ childNode.componentMpn }}</span>
                        <span>•</span>
                        <span>Qty: {{ childNode.quantity }} {{ childNode.unitName }}</span>
                      </div>
                    </div>
                  </div>

                  <!-- Child Price Information Column -->
                  <div class="text-right font-mono text-3xs shrink-0 pl-2">
                    <div class="font-bold text-slate-900">
                      {{ formatCurrency((childNode.sellingPrice || 0) * (childNode.quantity || 1), childNode.currency) }}
                    </div>
                    <div class="text-4xs text-slate-400 flex items-center justify-end gap-1 mt-0.5">
                      <span v-if="childNode.costPrice > 0" class="text-amber-700">Beli: {{ formatCurrency(childNode.costPrice, childNode.currency) }}</span>
                      <span v-else class="text-slate-400">Beli: Rp 0</span>
                      <span>|</span>
                      <span class="text-emerald-700 font-medium">Jual: {{ formatCurrency(childNode.sellingPrice, childNode.currency) }}</span>
                    </div>
                  </div>

                  <div class="flex items-center gap-1 opacity-80 group-hover:opacity-100 shrink-0">
                    <button
                      @click.stop="confirmDelete(childNode)"
                      class="p-1 rounded text-slate-400 hover:text-rose-600 hover:bg-white"
                      title="Remove Item"
                    >
                      <AppIcon name="trash" size="xs" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT PANEL: Selected Item Inspector (5 Cols) -->
      <div class="lg:col-span-5 space-y-4">
        <div v-if="selectedItem" class="app-card p-5 bg-white space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="info" size="xs" class="text-brand-600" />
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Element Inspector</h3>
            </div>
            <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200">
              {{ selectedItem.itemType }}
            </span>
          </div>

          <div class="space-y-3 text-xs">
            <div>
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Element Title</div>
              <div class="font-bold text-slate-900 text-sm mt-0.5">{{ selectedItem.name }}</div>
            </div>

            <!-- Referenced Product Details -->
            <div v-if="selectedItem.itemType === 'PRODUCT' && selectedItem.productCode" class="p-3 bg-slate-50 rounded-lg border border-slate-200 space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-3xs font-bold uppercase tracking-wider text-slate-500">Referenced Product Master</span>
                <span
                  :class="[
                    'text-3xs font-mono font-bold px-1.5 py-0.2 rounded border',
                    selectedItem.productType === 'TRADING' ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-emerald-50 text-emerald-700 border-emerald-200'
                  ]"
                >
                  {{ selectedItem.productType }} PRODUCT
                </span>
              </div>
              <div class="font-mono font-bold text-slate-800 text-xs">{{ selectedItem.productCode }}</div>
              <p class="text-3xs text-slate-500">
                This project element points to the master product record. It maintains its own independent lifecycle and specifications.
              </p>
              <router-link
                :to="`/products/${selectedItem.productId || selectedItem.productCode}`"
                class="text-3xs font-semibold text-brand-600 hover:text-brand-800 flex items-center gap-1 pt-1"
              >
                <span>Open Product Master ➔</span>
              </router-link>
            </div>

            <!-- Referenced Component Details -->
            <div v-else-if="selectedItem.itemType === 'COMPONENT'" class="p-3 bg-slate-50 rounded-lg border border-slate-200 space-y-2">
              <div class="text-3xs font-bold uppercase tracking-wider text-slate-500">Component Master Link</div>
              <div class="font-mono font-bold text-slate-800 text-xs">{{ selectedItem.componentMpn || 'Generic Hardware' }}</div>
              <router-link
                v-if="selectedItem.componentId"
                :to="`/components/${selectedItem.componentId}`"
                class="text-3xs font-semibold text-brand-600 hover:text-brand-800 flex items-center gap-1 pt-1"
              >
                <span>Open Component Master ➔</span>
              </router-link>
            </div>

            <div class="grid grid-cols-2 gap-3 pt-1">
              <div>
                <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Installed Quantity</div>
                <div class="font-mono font-bold text-slate-900 mt-0.5">{{ selectedItem.quantity }} {{ selectedItem.unitName }}</div>
              </div>
              <div>
                <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Lifecycle Status</div>
                <div class="font-medium text-emerald-600 mt-0.5">{{ selectedItem.status || 'Active' }}</div>
              </div>
            </div>

            <!-- Commercial & Pricing Terms Card -->
            <div class="p-3.5 bg-gradient-to-br from-slate-50 to-purple-50/30 rounded-xl border border-slate-200 space-y-2.5">
              <div class="flex items-center justify-between pb-2 border-b border-slate-200/80">
                <div class="flex items-center gap-1.5 text-3xs font-bold uppercase tracking-wider text-slate-700">
                  <AppIcon name="dollar-sign" size="2xs" class="text-emerald-600" />
                  <span>Commercial & Pricing Terms</span>
                </div>
                <span class="text-3xs font-mono font-bold px-1.5 py-0.2 rounded bg-white text-slate-700 border border-slate-200">
                  {{ selectedItem.currency || 'IDR' }}
                </span>
              </div>

              <div class="grid grid-cols-2 gap-3 text-xs">
                <div>
                  <div class="text-4xs font-semibold uppercase tracking-wider text-slate-400">Harga Beli (Cost)</div>
                  <div class="font-mono font-bold text-amber-800 text-xs mt-0.5">
                    {{ selectedItem.costPrice > 0 ? formatCurrency(selectedItem.costPrice, selectedItem.currency) : 'Rp 0' }}
                  </div>
                </div>
                <div>
                  <div class="text-4xs font-semibold uppercase tracking-wider text-slate-400">Harga Jual (Sell)</div>
                  <div class="font-mono font-bold text-emerald-800 text-xs mt-0.5">
                    {{ selectedItem.sellingPrice > 0 ? formatCurrency(selectedItem.sellingPrice, selectedItem.currency) : 'Rp 0' }}
                  </div>
                </div>
                <div>
                  <div class="text-4xs font-semibold uppercase tracking-wider text-slate-400">Kuantitas Terpasang</div>
                  <div class="font-mono font-bold text-slate-900 text-xs mt-0.5">
                    {{ selectedItem.quantity }} {{ selectedItem.unitName }}
                  </div>
                </div>
                <div>
                  <div class="text-4xs font-semibold uppercase tracking-wider text-slate-400">Subtotal Nilai Jual</div>
                  <div class="font-mono font-black text-brand-700 text-xs mt-0.5">
                    {{ formatCurrency((selectedItem.sellingPrice || 0) * (selectedItem.quantity || 1), selectedItem.currency) }}
                  </div>
                </div>
              </div>
            </div>

            <div v-if="selectedItem.notes">
              <div class="text-3xs font-semibold uppercase tracking-wider text-slate-400">Architecture Notes</div>
              <div class="text-3xs text-slate-700 bg-slate-50 p-2.5 rounded border border-slate-200 mt-1 leading-relaxed whitespace-pre-line">
                {{ selectedItem.notes }}
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center gap-2">
              <button @click="openEditModal(selectedItem)" class="btn btn-secondary btn-sm flex-1">
                <AppIcon name="edit" size="xs" class="mr-1" />
                <span>Edit Element</span>
              </button>
              <button @click="confirmDelete(selectedItem)" class="btn btn-danger btn-sm">
                <AppIcon name="trash" size="xs" />
              </button>
            </div>
          </div>
        </div>

        <div v-else class="app-card p-8 bg-white text-center text-slate-400">
          <AppIcon name="mouse-pointer" size="md" class="mx-auto mb-2 text-slate-300" />
          <div class="text-xs font-semibold text-slate-600">No Item Selected</div>
          <p class="text-3xs text-slate-400 mt-0.5">Click an item in the System Architecture Tree to inspect its master reference and specifications.</p>
        </div>
      </div>
    </div>

    <!-- Add/Edit Item Modal -->
    <BaseModal
      :isOpen="modalOpen"
      :title="isEditModal ? 'Edit Structure Item' : 'Add Project Structure Element'"
      @close="modalOpen = false"
    >
      <div class="space-y-4 text-xs">
        <!-- Item Type Selector (3 options) -->
        <div>
          <label class="form-label">Element Classification Type *</label>
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              @click="form.itemType = 'SUB_ASSEMBLY'"
              :class="[
                'p-2.5 rounded-lg border text-center transition-all',
                form.itemType === 'SUB_ASSEMBLY'
                  ? 'border-amber-500 bg-amber-50/50 text-amber-900 font-bold'
                  : 'border-slate-200 hover:border-slate-300 text-slate-600'
              ]"
            >
              <div class="text-xs">Sub-Assembly</div>
              <div class="text-3xs text-slate-400 mt-0.5">Grouping Module</div>
            </button>

            <button
              type="button"
              @click="form.itemType = 'PRODUCT'"
              :class="[
                'p-2.5 rounded-lg border text-center transition-all',
                form.itemType === 'PRODUCT'
                  ? 'border-brand-500 bg-brand-50/50 text-brand-900 font-bold'
                  : 'border-slate-200 hover:border-slate-300 text-slate-600'
              ]"
            >
              <div class="text-xs">Product Link</div>
              <div class="text-3xs text-slate-400 mt-0.5">Trading or R&D</div>
            </button>

            <button
              type="button"
              @click="form.itemType = 'COMPONENT'"
              :class="[
                'p-2.5 rounded-lg border text-center transition-all',
                form.itemType === 'COMPONENT'
                  ? 'border-emerald-500 bg-emerald-50/50 text-emerald-900 font-bold'
                  : 'border-slate-200 hover:border-slate-300 text-slate-600'
              ]"
            >
              <div class="text-xs">Component</div>
              <div class="text-3xs text-slate-400 mt-0.5">Physical Part</div>
            </button>
          </div>
        </div>

        <!-- Parent Sub-Assembly selector -->
        <div>
          <label class="form-label">Parent Structure Node (Optional)</label>
          <select v-model="form.parentItemId" class="form-select">
            <option :value="null">Root (Top Level Package)</option>
            <option
              v-for="sub in subAssemblies"
              :key="sub.id"
              :value="sub.id"
            >
              {{ sub.name }}
            </option>
          </select>
        </div>

        <!-- Product Reference Selector (When PRODUCT selected) -->
        <div v-if="form.itemType === 'PRODUCT'">
          <label class="form-label">Select Master Product *</label>
          <select v-model="form.productId" @change="onProductSelect" class="form-select font-medium" required>
            <option value="" disabled>Choose Product Master</option>
            <option
              v-for="p in availableProducts"
              :key="p.id"
              :value="p.id"
            >
              [{{ p.productType }}] {{ p.code }} — {{ p.name }}
            </option>
          </select>
        </div>

        <!-- Component Selector (When COMPONENT selected) -->
        <div v-else-if="form.itemType === 'COMPONENT'">
          <label class="form-label">Select Component Master (Optional link)</label>
          <select v-model="form.componentId" @change="onComponentSelect" class="form-select font-mono">
            <option value="">-- Generic Installation Part --</option>
            <option
              v-for="c in store.components"
              :key="c.id"
              :value="c.id"
            >
              {{ c.partNumber }} — {{ c.name }}
            </option>
          </select>
        </div>

        <!-- Item Name -->
        <div>
          <label class="form-label">Element Display Name *</label>
          <input v-model="form.name" type="text" class="form-input" placeholder="e.g. Core Telemetry Module" required />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="form-label">Quantity *</label>
            <input v-model.number="form.quantity" type="number" step="0.1" min="0.1" class="form-input font-mono" required />
          </div>
          <div>
            <label class="form-label">Unit of Measure</label>
            <input v-model="form.unitName" type="text" class="form-input" placeholder="pcs, kit, set, m" />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="form-label">Harga Beli / Cost (IDR / USD)</label>
            <input v-model.number="form.costPrice" type="number" step="0.01" min="0" class="form-input font-mono" placeholder="0" />
          </div>
          <div>
            <label class="form-label">Harga Jual / Selling Price (IDR / USD)</label>
            <input v-model.number="form.sellingPrice" type="number" step="0.01" min="0" class="form-input font-mono" placeholder="0" />
          </div>
        </div>

        <div>
          <label class="form-label">Architecture / Wiring Notes</label>
          <textarea v-model="form.notes" rows="2" class="form-textarea" placeholder="Connection interfaces, cabling requirements, mechanical placement..."></textarea>
        </div>
      </div>

      <template #footer>
        <button @click="modalOpen = false" class="btn btn-secondary">Cancel</button>
        <button @click="saveForm" class="btn btn-primary">
          {{ isEditModal ? 'Save Changes' : 'Add to Structure' }}
        </button>
      </template>
    </BaseModal>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ProjectStructureWorkspace',
  components: {
    BaseModal,
    AppIcon
  },
  props: {
    productId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      treeSearch: '',
      allExpanded: true,
      expandedNodes: {},
      selectedItemId: null,
      modalOpen: false,
      isEditModal: false,
      form: {
        id: null,
        itemType: 'PRODUCT',
        parentItemId: null,
        productId: '',
        productCode: '',
        productType: '',
        componentId: '',
        componentMpn: '',
        name: '',
        quantity: 1.0,
        costPrice: 0,
        sellingPrice: 0,
        currency: 'IDR',
        unitName: 'pcs',
        notes: ''
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    flatItems() {
      return this.store.projectItems[this.productId] || [];
    },
    filteredItems() {
      if (!this.treeSearch) return this.flatItems;
      const q = this.treeSearch.toLowerCase().trim();
      return this.flatItems.filter(i =>
        i.name.toLowerCase().includes(q) ||
        (i.productCode && i.productCode.toLowerCase().includes(q)) ||
        (i.componentMpn && i.componentMpn.toLowerCase().includes(q))
      );
    },
    rootTreeItems() {
      return this.filteredItems.filter(i => !i.parentItemId);
    },
    subAssemblies() {
      return this.flatItems.filter(i => i.itemType === 'SUB_ASSEMBLY');
    },
    selectedItem() {
      if (!this.selectedItemId) return null;
      return this.flatItems.find(i => i.id === this.selectedItemId);
    },
    availableProducts() {
      // Exclude self project to prevent cycle
      return this.store.products.filter(p => p.id !== this.productId && p.code !== this.productId);
    },
    totalProjectRevenue() {
      return this.flatItems.reduce((sum, item) => sum + ((item.sellingPrice || 0) * (item.quantity || 1)), 0);
    },
    totalProjectCost() {
      return this.flatItems.reduce((sum, item) => sum + ((item.costPrice || 0) * (item.quantity || 1)), 0);
    },
    projectGrossProfit() {
      return this.totalProjectRevenue - this.totalProjectCost;
    },
    projectMarginPercent() {
      if (this.totalProjectRevenue <= 0) return 0;
      return ((this.projectGrossProfit / this.totalProjectRevenue) * 100).toFixed(1);
    }
  },
  watch: {
    flatItems: {
      immediate: true,
      handler(items) {
        if (items && items.length > 0) {
          items.forEach(i => {
            if (this.expandedNodes[i.id] === undefined) {
              this.$set(this.expandedNodes, i.id, true);
            }
          });
          if (!this.selectedItemId && items[0]) {
            this.selectedItemId = items[0].id;
          }
        }
      }
    }
  },
  async mounted() {
    await Promise.all([
      this.store.fetchProjectItems(this.productId),
      this.store.products.length === 0 ? this.store.fetchProducts() : Promise.resolve(),
      this.store.components.length === 0 ? this.store.fetchComponents() : Promise.resolve()
    ]);
  },
  methods: {
    formatCurrency(val, currency = 'IDR') {
      if (val === null || val === undefined) return 'Rp 0';
      const num = Number(val);
      if (isNaN(num)) return 'Rp 0';
      const formatted = num.toLocaleString('id-ID', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
      });
      return `${currency || 'Rp'} ${formatted}`;
    },
    selectItem(item) {
      this.selectedItemId = item.id;
    },
    hasChildren(itemId) {
      return this.flatItems.some(i => i.parentItemId === itemId);
    },
    getChildren(itemId) {
      return this.filteredItems.filter(i => i.parentItemId === itemId);
    },
    isExpanded(nodeId) {
      return this.expandedNodes[nodeId] !== false;
    },
    toggleExpand(nodeId) {
      this.$set(this.expandedNodes, nodeId, !this.isExpanded(nodeId));
    },
    toggleAll() {
      this.allExpanded = !this.allExpanded;
      this.flatItems.forEach(i => {
        this.$set(this.expandedNodes, i.id, this.allExpanded);
      });
    },
    openAddModal(parentId = null, defaultType = 'PRODUCT') {
      this.isEditModal = false;
      this.form = {
        id: null,
        itemType: defaultType,
        parentItemId: parentId,
        productId: '',
        productCode: '',
        productType: '',
        componentId: '',
        componentMpn: '',
        name: '',
        quantity: 1.0,
        costPrice: 0,
        sellingPrice: 0,
        currency: 'IDR',
        unitName: 'pcs',
        notes: ''
      };
      this.modalOpen = true;
    },
    openEditModal(item) {
      this.isEditModal = true;
      this.form = { ...item };
      this.modalOpen = true;
    },
    async onProductSelect() {
      const p = this.store.products.find(item => item.id === this.form.productId || item.code === this.form.productId);
      if (p) {
        this.form.name = p.name;
        this.form.productCode = p.code;
        this.form.productType = p.productType;
        if (p.productType === 'TRADING' && p.id) {
          const detail = await this.store.fetchTradingDetail(p.id);
          if (detail) {
            this.form.costPrice = detail.purchasePrice || 0;
            this.form.sellingPrice = detail.sellingPrice || 0;
            this.form.currency = detail.currency || 'IDR';
          }
        }
      }
    },
    onComponentSelect() {
      const c = this.store.components.find(item => item.id === this.form.componentId);
      if (c) {
        this.form.name = c.name;
        this.form.componentMpn = c.partNumber;
        this.form.costPrice = c.estimatedUnitCost || 0;
        this.form.currency = c.currency || 'USD';
      }
    },
    async saveForm() {
      if (!this.form.name) {
        alert('Please provide an element name.');
        return;
      }

      try {
        if (this.isEditModal) {
          await this.store.updateProjectItem(this.productId, this.form.id, this.form);
        } else {
          await this.store.createProjectItem(this.productId, this.form);
        }
        this.modalOpen = false;
      } catch (err) {
        alert('Failed to save project item: ' + (err.message || err));
      }
    },
    async confirmDelete(item) {
      if (confirm(`Remove '${item.name}' from this project package?`)) {
        try {
          await this.store.deleteProjectItem(this.productId, item.id);
          if (this.selectedItemId === item.id) {
            this.selectedItemId = null;
          }
        } catch (err) {
          alert('Failed to delete item: ' + (err.message || err));
        }
      }
    }
  }
};
</script>
