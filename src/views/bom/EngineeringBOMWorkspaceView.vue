<template>
  <div class="space-y-4">
    <!-- TOP HEADER / REVISION CONTEXT -->
    <div class="app-card p-5 md:p-6 bg-white space-y-4">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-4 border-b border-slate-100">
        <div class="flex items-start gap-4 min-w-0">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-tr from-brand-900 via-slate-900 to-brand-700 text-white flex items-center justify-center font-mono font-bold text-xl shadow-md shrink-0">
            BOM
          </div>
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono font-bold text-xs text-slate-900 bg-slate-100 px-2.5 py-0.5 rounded border border-slate-200">
                {{ revision ? revision.code : 'REV' }}
              </span>
              <span v-if="revision" class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                v{{ revision.versionNumber }}
              </span>
              <span v-if="revision" class="text-3xs font-mono text-slate-500 font-semibold">
                {{ revision.productCode }}
              </span>
              <span :class="getStatusBadgeClass(bom ? bom.status : 'Draft')">
                {{ bom ? bom.status : 'Draft' }}
              </span>
            </div>
            <h1 class="text-xl md:text-2xl font-bold text-slate-900 tracking-tight truncate">
              {{ bom ? bom.name : 'Engineering Bill of Materials' }}
            </h1>
            <p class="text-xs text-slate-500 mt-1 font-mono">
              BOM Code: <span class="font-bold text-slate-700">{{ bom ? bom.bomCode : 'PENDING' }}</span>
              <span v-if="bom && bom.createdBy" class="ml-3 font-sans text-slate-400">Created by: {{ bom.createdBy }}</span>
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <button
            @click="goBackToRevision"
            class="btn btn-secondary"
          >
            <AppIcon name="chevron-left" size="xs" class="mr-1.5 text-slate-500" />
            <span>Revision Detail</span>
          </button>

          <button
            v-if="bom"
            @click="openSettingsModal"
            class="btn btn-secondary"
            title="Configure Target Margin & BOM Metadata"
          >
            <AppIcon name="settings" size="xs" class="mr-1.5 text-slate-600" />
            <span>BOM Settings</span>
          </button>

          <button
            v-if="bom"
            @click="openAddItemModal(null)"
            class="btn btn-primary"
          >
            <AppIcon name="plus" size="xs" class="mr-1.5 text-white" />
            <span>Add Item / Assembly</span>
          </button>
        </div>
      </div>

      <!-- KEY ENGINEERING COST METRICS RIBBON -->
      <div v-if="bom" class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-1">
        <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200/80">
          <div class="flex items-center justify-between text-3xs font-semibold uppercase tracking-wider text-slate-400">
            <span>Estimated BOM Cost</span>
            <AppIcon name="dollar-sign" size="xs" class="text-brand-500" />
          </div>
          <div class="mt-1 flex items-baseline gap-1.5">
            <span class="text-xl font-bold text-slate-900 font-mono">${{ formatNumber(bom.totalCost) }}</span>
            <span class="text-3xs font-mono font-semibold text-slate-500">{{ bom.currency }}</span>
          </div>
          <span class="text-3xs text-slate-400 mt-0.5 block">Rollup of all leaf components</span>
        </div>

        <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200/80">
          <div class="flex items-center justify-between text-3xs font-semibold uppercase tracking-wider text-slate-400">
            <span>Target Gross Margin</span>
            <AppIcon name="trending-up" size="xs" class="text-emerald-500" />
          </div>
          <div class="mt-1 flex items-baseline gap-1">
            <span class="text-xl font-bold text-emerald-700 font-mono">{{ formatNumber(bom.targetMargin, 1) }}%</span>
          </div>
          <span class="text-3xs text-slate-400 mt-0.5 block">Configurable target profitability</span>
        </div>

        <div class="p-3.5 rounded-xl bg-emerald-50/70 border border-emerald-200/80">
          <div class="flex items-center justify-between text-3xs font-semibold uppercase tracking-wider text-emerald-700">
            <span>Estimated Selling Price</span>
            <AppIcon name="check-circle" size="xs" class="text-emerald-600" />
          </div>
          <div class="mt-1 flex items-baseline gap-1.5">
            <span class="text-xl font-bold text-emerald-900 font-mono">${{ formatNumber(bom.estimatedSellingPrice) }}</span>
            <span class="text-3xs font-mono font-semibold text-emerald-700">{{ bom.currency }}</span>
          </div>
          <span class="text-3xs text-emerald-600/80 mt-0.5 block">Cost / (1 - Margin %)</span>
        </div>

        <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200/80">
          <div class="flex items-center justify-between text-3xs font-semibold uppercase tracking-wider text-slate-400">
            <span>Structure Breakdown</span>
            <AppIcon name="layers" size="xs" class="text-indigo-500" />
          </div>
          <div class="mt-1 flex items-baseline gap-2 font-mono">
            <span class="text-lg font-bold text-slate-900">{{ costSummary ? costSummary.totalComponentsCount : 0 }}</span>
            <span class="text-3xs text-slate-400">Items /</span>
            <span class="text-lg font-bold text-indigo-600">{{ costSummary ? costSummary.totalAssembliesCount : 0 }}</span>
            <span class="text-3xs text-slate-400">Assemblies</span>
          </div>
          <span class="text-3xs text-slate-400 mt-0.5 block">Multi-level hierarchy</span>
        </div>
      </div>
    </div>

    <!-- INITIALIZE BOM PROMPT IF NO BOM EXISTS -->
    <div v-if="!loading && !bom" class="app-card p-12 text-center bg-white space-y-4">
      <div class="w-16 h-16 rounded-2xl bg-brand-50 text-brand-600 flex items-center justify-center mx-auto text-2xl font-mono font-bold shadow-inner">
        BOM
      </div>
      <div class="max-w-md mx-auto">
        <h3 class="text-lg font-bold text-slate-900">No Active Engineering BOM for this Revision</h3>
        <p class="text-xs text-slate-500 mt-1">
          Hardware Revision <strong class="font-mono text-slate-800">{{ revision ? revision.code : '' }}</strong> does not have an active engineering Bill of Materials yet. Initialize a new multi-level BOM baseline to manage components, assemblies, and cost rollups.
        </p>
      </div>
      <div>
        <button
          @click="initBOM"
          class="btn btn-primary px-6"
        >
          <AppIcon name="plus" size="xs" class="mr-2 text-white" />
          <span>Initialize Engineering BOM</span>
        </button>
      </div>
    </div>

    <!-- MAIN 3-PANE WORKSPACE -->
    <div v-if="bom" class="grid grid-cols-1 lg:grid-cols-12 gap-4">
      <!-- LEFT PANE: BOM HIERARCHY TREE EXPLORER (3 Cols) -->
      <div class="lg:col-span-3 space-y-3">
        <div class="app-card p-4 bg-white space-y-3">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div class="flex items-center gap-1.5 text-xs font-bold text-slate-900">
              <AppIcon name="folder-tree" size="xs" class="text-brand-600" />
              <span>BOM Hierarchy Tree</span>
            </div>
            <span class="text-3xs font-mono font-semibold text-slate-500 bg-slate-100 px-2 py-0.5 rounded">
              {{ flatItems.length }} nodes
            </span>
          </div>

          <!-- Tree Search Filter -->
          <div>
            <input
              v-model="treeFilter"
              type="text"
              placeholder="Filter tree items..."
              class="w-full text-xs px-2.5 py-1.5 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-sans"
            />
          </div>

          <!-- Tree List -->
          <div class="space-y-1 max-h-[560px] overflow-y-auto pr-1">
            <div
              v-for="item in filteredFlatItems"
              :key="item.id"
              @click="selectItem(item.id)"
              :style="{ paddingLeft: `${item.depth * 14 + 8}px` }"
              :class="[
                'p-2 rounded-lg text-xs cursor-pointer transition-all flex items-center justify-between group border',
                selectedItemId === item.id
                  ? 'bg-brand-50/80 border-brand-200 text-brand-900 font-semibold shadow-xs'
                  : 'border-transparent text-slate-700 hover:bg-slate-50 hover:border-slate-200'
              ]"
            >
              <div class="flex items-center gap-2 min-w-0">
                <span
                  :class="[
                    'w-5 h-5 rounded flex items-center justify-center text-3xs font-mono font-bold shrink-0',
                    item.itemType === 'SUB_ASSEMBLY'
                      ? 'bg-indigo-100 text-indigo-700'
                      : 'bg-emerald-100 text-emerald-700'
                  ]"
                >
                  {{ item.itemType === 'SUB_ASSEMBLY' ? 'A' : 'C' }}
                </span>
                <div class="min-w-0">
                  <div class="truncate text-xs font-medium">{{ item.name }}</div>
                  <div v-if="item.componentPartNumber" class="text-3xs font-mono text-slate-400 truncate">
                    {{ item.componentPartNumber }}
                  </div>
                </div>
              </div>

              <div class="text-right shrink-0 font-mono text-3xs font-semibold text-slate-500 pl-2">
                ${{ formatNumber(item.lineCost) }}
              </div>
            </div>
          </div>
        </div>

        <!-- CATEGORY COST DISTRIBUTION CARD -->
        <div v-if="costSummary && costSummary.costByCategory" class="app-card p-4 bg-white space-y-2.5">
          <div class="text-xs font-bold text-slate-900 flex items-center justify-between">
            <span>Cost by Category</span>
            <span class="text-3xs text-slate-400 font-mono font-normal">Rollup</span>
          </div>

          <div class="space-y-2 pt-1">
            <div
              v-for="(cost, catName) in costSummary.costByCategory"
              :key="catName"
              class="space-y-1"
            >
              <div class="flex items-center justify-between text-3xs font-semibold">
                <span class="text-slate-700 truncate max-w-[120px]">{{ catName }}</span>
                <span class="font-mono text-slate-900">${{ formatNumber(cost) }}</span>
              </div>
              <div class="w-full bg-slate-100 rounded-full h-1.5 overflow-hidden">
                <div
                  class="bg-brand-600 h-1.5 rounded-full transition-all duration-500"
                  :style="{ width: `${bom.totalCost > 0 ? (cost / bom.totalCost) * 100 : 0}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- CENTER PANE: ENGINEERING MULTI-LEVEL BOM TABLE (6 Cols) -->
      <div class="lg:col-span-6 space-y-3">
        <div class="app-card overflow-hidden bg-white">
          <div class="p-4 border-b border-slate-200 bg-slate-50/60 flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span class="font-bold text-xs text-slate-900">Engineering Items Table</span>
              <span class="text-3xs font-mono px-2 py-0.5 rounded-full bg-slate-200/80 text-slate-700 font-semibold">
                {{ flatItems.length }} total lines
              </span>
            </div>

            <div class="flex items-center gap-2">
              <button
                @click="openAddItemModal(null)"
                class="btn btn-secondary text-xs py-1"
              >
                <AppIcon name="plus" size="xs" class="mr-1 text-slate-600" />
                <span>Add Root Item</span>
              </button>
            </div>
          </div>

          <!-- TABLE -->
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="bg-slate-50/90 text-slate-400 uppercase text-3xs font-bold border-b border-slate-200 tracking-wider font-mono">
                  <th class="py-3 px-3">Item / Description</th>
                  <th class="py-3 px-2">MPN</th>
                  <th class="py-3 px-2 text-right">Qty</th>
                  <th class="py-3 px-2">Unit</th>
                  <th class="py-3 px-2 text-right">Unit Cost</th>
                  <th class="py-3 px-2 text-right">Line Cost</th>
                  <th class="py-3 px-2">Ref Des</th>
                  <th class="py-3 px-2 text-center">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="item in flatItems"
                  :key="item.id"
                  @click="selectItem(item.id)"
                  :class="[
                    'hover:bg-slate-50/80 transition-colors cursor-pointer',
                    selectedItemId === item.id ? 'bg-brand-50/50 font-medium' : ''
                  ]"
                >
                  <!-- Name & Depth -->
                  <td class="py-3 px-3">
                    <div class="flex items-center gap-1.5" :style="{ paddingLeft: `${item.depth * 16}px` }">
                      <span
                        :class="[
                          'text-3xs font-mono font-bold px-1.5 py-0.2 rounded shrink-0',
                          item.itemType === 'SUB_ASSEMBLY'
                            ? 'bg-indigo-100 text-indigo-800'
                            : 'bg-emerald-100 text-emerald-800'
                        ]"
                      >
                        {{ item.itemType === 'SUB_ASSEMBLY' ? 'ASM' : 'CMP' }}
                      </span>
                      <span class="font-medium text-slate-900 truncate max-w-[180px]">{{ item.name }}</span>
                    </div>
                  </td>

                  <!-- MPN -->
                  <td class="py-3 px-2 font-mono text-3xs text-slate-600 truncate max-w-[110px]">
                    {{ item.componentPartNumber || '—' }}
                  </td>

                  <!-- QTY -->
                  <td class="py-3 px-2 text-right font-mono font-bold text-slate-900">
                    {{ formatNumber(item.quantity, 2) }}
                  </td>

                  <!-- UNIT -->
                  <td class="py-3 px-2 font-mono text-3xs text-slate-500">
                    {{ item.unitCode || 'PCS' }}
                  </td>

                  <!-- UNIT COST -->
                  <td class="py-3 px-2 text-right font-mono text-slate-600 text-3xs">
                    ${{ formatNumber(item.unitCost) }}
                  </td>

                  <!-- LINE COST -->
                  <td class="py-3 px-2 text-right font-mono font-bold text-slate-900">
                    ${{ formatNumber(item.lineCost) }}
                  </td>

                  <!-- REFERENCE DESIGNATOR (JETBRAINS MONO) -->
                  <td class="py-3 px-2 font-mono text-3xs text-brand-700 max-w-[90px] truncate">
                    <span v-if="item.referenceDesignators" class="bg-brand-50 border border-brand-200/80 px-1.5 py-0.5 rounded">
                      {{ item.referenceDesignators }}
                    </span>
                    <span v-else class="text-slate-300">—</span>
                  </td>

                  <!-- ACTIONS -->
                  <td class="py-3 px-2 text-center">
                    <div class="flex items-center justify-center gap-1">
                      <button
                        v-if="item.itemType === 'SUB_ASSEMBLY'"
                        @click.stop="openAddItemModal(item.id)"
                        class="p-1 text-slate-400 hover:text-brand-600 rounded hover:bg-slate-100"
                        title="Add Child Item inside this Sub-Assembly"
                      >
                        <AppIcon name="plus" size="xs" />
                      </button>
                      <button
                        @click.stop="openEditItemModal(item)"
                        class="p-1 text-slate-400 hover:text-slate-700 rounded hover:bg-slate-100"
                        title="Edit Item"
                      >
                        <AppIcon name="edit" size="xs" />
                      </button>
                      <button
                        @click.stop="confirmDeleteItem(item)"
                        class="p-1 text-slate-400 hover:text-rose-600 rounded hover:bg-slate-100"
                        title="Delete Item"
                      >
                        <AppIcon name="trash" size="xs" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- RIGHT PANE: SELECTED ITEM INSPECTOR (3 Cols) -->
      <div class="lg:col-span-3 space-y-3">
        <div v-if="selectedItem" class="app-card p-4 bg-white space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <span
                :class="[
                  'text-3xs font-mono font-bold px-2 py-0.5 rounded',
                  selectedItem.itemType === 'SUB_ASSEMBLY'
                    ? 'bg-indigo-100 text-indigo-800'
                    : 'bg-emerald-100 text-emerald-800'
                ]"
              >
                {{ selectedItem.itemType === 'SUB_ASSEMBLY' ? 'SUB-ASSEMBLY' : 'COMPONENT' }}
              </span>
              <span class="font-bold text-xs text-slate-900">Item Inspector</span>
            </div>

            <button
              @click="openEditItemModal(selectedItem)"
              class="text-3xs font-semibold text-brand-600 hover:underline flex items-center gap-1"
            >
              <AppIcon name="edit" size="xs" />
              <span>Edit</span>
            </button>
          </div>

          <!-- Item Title & MPN -->
          <div>
            <h3 class="font-bold text-sm text-slate-900 leading-snug">
              {{ selectedItem.name }}
            </h3>
            <p v-if="selectedItem.componentPartNumber" class="text-xs font-mono text-brand-700 font-semibold mt-1">
              MPN: {{ selectedItem.componentPartNumber }}
            </p>
          </div>

          <!-- Cost Calculation Card -->
          <div class="p-3 rounded-xl bg-slate-50 border border-slate-200/80 space-y-2">
            <span class="text-3xs font-bold text-slate-400 uppercase tracking-wider block">Line Cost Breakdown</span>
            <div class="flex items-center justify-between text-xs">
              <span class="text-slate-500">Unit Cost:</span>
              <span class="font-mono font-semibold text-slate-900">${{ formatNumber(selectedItem.unitCost) }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-slate-500">Quantity:</span>
              <span class="font-mono font-bold text-slate-900">{{ selectedItem.quantity }} {{ selectedItem.unitCode }}</span>
            </div>
            <div class="pt-2 border-t border-slate-200 flex items-center justify-between text-xs">
              <span class="font-bold text-slate-800">Total Line Cost:</span>
              <span class="font-mono font-bold text-emerald-700 text-sm">${{ formatNumber(selectedItem.lineCost) }}</span>
            </div>
          </div>

          <!-- Engineering Reference Designators -->
          <div v-if="selectedItem.referenceDesignators" class="space-y-1.5">
            <span class="text-3xs font-bold text-slate-400 uppercase tracking-wider block">Reference Designators</span>
            <div class="p-2.5 rounded-lg bg-slate-900 text-emerald-400 font-mono text-xs break-all">
              {{ selectedItem.referenceDesignators }}
            </div>
          </div>

          <!-- Notes -->
          <div v-if="selectedItem.notes" class="space-y-1">
            <span class="text-3xs font-bold text-slate-400 uppercase tracking-wider block">Engineering Notes</span>
            <p class="text-xs text-slate-600 bg-slate-50 p-2.5 rounded-lg border border-slate-100 leading-relaxed font-sans">
              {{ selectedItem.notes }}
            </p>
          </div>

          <!-- Quick Navigation to Master Component -->
          <div v-if="selectedItem.componentId" class="pt-2 border-t border-slate-100">
            <router-link
              :to="`/components/${selectedItem.componentId}`"
              class="btn btn-secondary w-full justify-center text-xs py-2"
            >
              <AppIcon name="external-link" size="xs" class="mr-1.5 text-slate-500" />
              <span>View in Component Master</span>
            </router-link>
          </div>
        </div>

        <div v-else class="app-card p-8 text-center bg-white text-slate-400 text-xs">
          <AppIcon name="mouse-pointer" size="md" class="mx-auto text-slate-300 mb-2" />
          <p>Select an item from the tree or table to inspect engineering specifications.</p>
        </div>
      </div>
    </div>

    <!-- MODAL: ADD / EDIT BOM ITEM -->
    <div
      v-if="showItemModal"
      class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto"
    >
      <div class="bg-white rounded-2xl max-w-xl w-full p-6 shadow-2xl border border-slate-200 space-y-5 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-lg bg-brand-50 text-brand-600 flex items-center justify-center">
              <AppIcon name="plus" size="xs" />
            </div>
            <h3 class="font-bold text-base text-slate-900">
              {{ editingItemId ? 'Edit BOM Item' : 'Add Item to Engineering BOM' }}
            </h3>
          </div>
          <button @click="closeItemModal" class="text-slate-400 hover:text-slate-600">
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <form @submit.prevent="saveItem" class="space-y-4">
          <!-- Item Type Selector -->
          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Item Classification</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                @click="itemForm.itemType = 'COMPONENT'"
                :class="[
                  'p-3 rounded-xl border text-left transition-all flex items-center gap-2.5',
                  itemForm.itemType === 'COMPONENT'
                    ? 'border-brand-500 bg-brand-50/60 ring-2 ring-brand-500/20'
                    : 'border-slate-200 hover:bg-slate-50'
                ]"
              >
                <div class="w-7 h-7 rounded-lg bg-emerald-100 text-emerald-700 flex items-center justify-center font-bold font-mono text-xs">C</div>
                <div>
                  <div class="text-xs font-bold text-slate-900">Component</div>
                  <div class="text-3xs text-slate-400">Electronic/mechanical part</div>
                </div>
              </button>

              <button
                type="button"
                @click="itemForm.itemType = 'SUB_ASSEMBLY'"
                :class="[
                  'p-3 rounded-xl border text-left transition-all flex items-center gap-2.5',
                  itemForm.itemType === 'SUB_ASSEMBLY'
                    ? 'border-indigo-500 bg-indigo-50/60 ring-2 ring-indigo-500/20'
                    : 'border-slate-200 hover:bg-slate-50'
                ]"
              >
                <div class="w-7 h-7 rounded-lg bg-indigo-100 text-indigo-700 flex items-center justify-center font-bold font-mono text-xs">A</div>
                <div>
                  <div class="text-xs font-bold text-slate-900">Sub-Assembly</div>
                  <div class="text-3xs text-slate-400">Nested module / board</div>
                </div>
              </button>
            </div>
          </div>

          <!-- Component Master Picker (Only if COMPONENT) -->
          <div v-if="itemForm.itemType === 'COMPONENT'">
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Select Master Component *</label>
            <select
              v-model="itemForm.componentId"
              @change="onComponentSelected"
              required
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 bg-white"
            >
              <option value="" disabled>-- Select component from master library --</option>
              <option
                v-for="comp in componentsList"
                :key="comp.id"
                :value="comp.id"
              >
                {{ comp.partNumber }} — {{ comp.name }} (${{ formatNumber(comp.estimatedUnitCost) }})
              </option>
            </select>
          </div>

          <!-- Item Name -->
          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Item / Assembly Name *</label>
            <input
              v-model="itemForm.name"
              type="text"
              required
              placeholder="e.g. Primary DC-DC Power Stage or STM32 Microcontroller"
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
            />
          </div>

          <!-- Parent Sub-Assembly -->
          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Parent Sub-Assembly (Hierarchy)</label>
            <select
              v-model="itemForm.parentItemId"
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 bg-white"
            >
              <option :value="null">-- Root Item (Top-Level) --</option>
              <option
                v-for="asm in subAssemblies"
                :key="asm.id"
                :value="asm.id"
                :disabled="editingItemId === asm.id"
              >
                {{ asm.name }}
              </option>
            </select>
          </div>

          <!-- Quantity, Unit & Unit Cost -->
          <div class="grid grid-cols-3 gap-3">
            <div>
              <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Quantity *</label>
              <input
                v-model.number="itemForm.quantity"
                type="number"
                step="0.01"
                min="0.01"
                required
                class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono"
              />
            </div>
            <div>
              <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Unit</label>
              <input
                v-model="itemForm.unitCode"
                type="text"
                placeholder="PCS"
                class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono"
              />
            </div>
            <div>
              <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Unit Cost ($)</label>
              <input
                v-model.number="itemForm.unitCost"
                type="number"
                step="0.0001"
                min="0"
                :disabled="itemForm.itemType === 'SUB_ASSEMBLY'"
                class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono disabled:bg-slate-100 disabled:text-slate-400"
              />
            </div>
          </div>

          <!-- Reference Designators -->
          <div v-if="itemForm.itemType === 'COMPONENT'">
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">
              Reference Designators (JetBrains Mono)
            </label>
            <input
              v-model="itemForm.referenceDesignators"
              type="text"
              placeholder="e.g. U1, U2 or C1, C2, C3, C4 or R101-R104"
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono text-brand-700"
            />
            <span class="text-3xs text-slate-400 mt-1 block">Comma-separated PCB component identifiers</span>
          </div>

          <!-- Notes -->
          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Engineering Notes</label>
            <textarea
              v-model="itemForm.notes"
              rows="2"
              placeholder="Assembly instructions, tolerance requirements, or thermal considerations..."
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
            ></textarea>
          </div>

          <!-- Actions -->
          <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
            <button
              type="button"
              @click="closeItemModal"
              class="btn btn-secondary"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="btn btn-primary"
            >
              <AppIcon name="check" size="xs" class="mr-1.5 text-white" />
              <span>{{ editingItemId ? 'Save Changes' : 'Add to BOM' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: BOM SETTINGS & TARGET MARGIN -->
    <div
      v-if="showSettingsModal"
      class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto"
    >
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-slate-200 space-y-5 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-lg bg-slate-100 text-slate-700 flex items-center justify-center">
              <AppIcon name="settings" size="xs" />
            </div>
            <h3 class="font-bold text-base text-slate-900">BOM Configuration & Target Margin</h3>
          </div>
          <button @click="showSettingsModal = false" class="text-slate-400 hover:text-slate-600">
            <AppIcon name="x" size="xs" />
          </button>
        </div>

        <form @submit.prevent="saveBOMSettings" class="space-y-4">
          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">BOM Name</label>
            <input
              v-model="settingsForm.name"
              type="text"
              required
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
            />
          </div>

          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Target Gross Margin (%)</label>
            <div class="relative">
              <input
                v-model.number="settingsForm.targetMargin"
                type="number"
                step="0.5"
                min="1"
                max="99"
                required
                class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono pr-8"
              />
              <span class="absolute right-3 top-2 text-xs font-bold text-slate-400">%</span>
            </div>
            <span class="text-3xs text-slate-400 mt-1 block">Formula: Estimated Price = BOM Cost / (1 - Margin%)</span>
          </div>

          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">BOM Lifecycle Status</label>
            <select
              v-model="settingsForm.status"
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500 bg-white"
            >
              <option value="Draft">Draft (In Engineering)</option>
              <option value="Active">Active (Released Baseline)</option>
              <option value="Obsolete">Obsolete</option>
              <option value="Archived">Archived</option>
            </select>
          </div>

          <div>
            <label class="block text-3xs uppercase font-bold text-slate-500 mb-1">Description / Release Notes</label>
            <textarea
              v-model="settingsForm.description"
              rows="3"
              class="w-full text-xs px-3 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
            ></textarea>
          </div>

          <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
            <button
              type="button"
              @click="showSettingsModal = false"
              class="btn btn-secondary"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="btn btn-primary"
            >
              Save Configuration
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia';
import { useBOMStore } from '@/stores/bomStore';
import { useMasterStore } from '@/stores/masterStore';
import { useVersionStore } from '@/stores/versionStore';

export default {
  name: 'EngineeringBOMWorkspaceView',
  data() {
    return {
      revisionId: this.$route.params.id,
      treeFilter: '',
      showItemModal: false,
      editingItemId: null,
      itemForm: {
        itemType: 'COMPONENT',
        componentId: '',
        name: '',
        parentItemId: null,
        quantity: 1.0,
        unitCode: 'PCS',
        unitCost: 0,
        referenceDesignators: '',
        notes: ''
      },
      showSettingsModal: false,
      settingsForm: {
        name: '',
        targetMargin: 35.0,
        status: 'Active',
        description: ''
      }
    };
  },
  computed: {
    ...mapState(useBOMStore, ['currentBOM', 'costSummary', 'selectedItemId', 'selectedItem', 'flatItems', 'subAssemblies', 'loading']),
    ...mapState(useMasterStore, ['components']),
    ...mapState(useVersionStore, ['revisions']),

    bom() {
      return this.currentBOM;
    },

    revision() {
      if (this.revisions && this.revisions.length > 0) {
        return this.revisions.find(r => r.id === this.revisionId);
      }
      return null;
    },

    componentsList() {
      return this.components || [];
    },

    filteredFlatItems() {
      if (!this.treeFilter) return this.flatItems;
      const q = this.treeFilter.toLowerCase();
      return this.flatItems.filter(item =>
        item.name.toLowerCase().includes(q) ||
        (item.componentPartNumber && item.componentPartNumber.toLowerCase().includes(q))
      );
    }
  },
  async created() {
    await this.fetchRevisionContext();
    await this.fetchBOMByRevision(this.revisionId);
    if (!this.components || this.components.length === 0) {
      await this.fetchComponents();
    }
  },
  methods: {
    ...mapActions(useBOMStore, ['fetchBOMByRevision', 'createBOM', 'updateBOM', 'addItem', 'updateItem', 'deleteItem', 'selectItem']),
    ...mapActions(useMasterStore, ['fetchComponents']),
    ...mapActions(useVersionStore, ['fetchRevisions']),

    async fetchRevisionContext() {
      if (!this.revisions || this.revisions.length === 0) {
        await this.fetchRevisions();
      }
    },

    formatNumber(num, decimals = 2) {
      if (num === null || num === undefined) return '0.00';
      return Number(num).toLocaleString('en-US', {
        minimumFractionDigits: decimals,
        maximumFractionDigits: decimals
      });
    },

    getStatusBadgeClass(status) {
      switch (status) {
        case 'Active':
          return 'text-3xs font-mono font-bold px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 border border-emerald-200';
        case 'Draft':
          return 'text-3xs font-mono font-bold px-2 py-0.5 rounded bg-amber-50 text-amber-700 border border-amber-200';
        case 'Obsolete':
          return 'text-3xs font-mono font-bold px-2 py-0.5 rounded bg-rose-50 text-rose-700 border border-rose-200';
        default:
          return 'text-3xs font-mono font-bold px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200';
      }
    },

    goBackToRevision() {
      this.$router.push(`/revisions/${this.revisionId}`);
    },

    async initBOM() {
      try {
        await this.createBOM(this.revisionId, {
          targetMargin: 35.0,
          status: 'Active'
        });
      } catch (err) {
        alert(err.response?.data?.message || 'Failed to initialize BOM');
      }
    },

    openAddItemModal(parentItemId = null) {
      this.editingItemId = null;
      this.itemForm = {
        itemType: 'COMPONENT',
        componentId: '',
        name: '',
        parentItemId: parentItemId,
        quantity: 1.0,
        unitCode: 'PCS',
        unitCost: 0,
        referenceDesignators: '',
        notes: ''
      };
      this.showItemModal = true;
    },

    openEditItemModal(item) {
      this.editingItemId = item.id;
      this.itemForm = {
        itemType: item.itemType,
        componentId: item.componentId || '',
        name: item.name,
        parentItemId: item.parentItemId || null,
        quantity: item.quantity,
        unitCode: item.unitCode || 'PCS',
        unitCost: item.unitCost,
        referenceDesignators: item.referenceDesignators || '',
        notes: item.notes || ''
      };
      this.showItemModal = true;
    },

    closeItemModal() {
      this.showItemModal = false;
      this.editingItemId = null;
    },

    onComponentSelected() {
      if (!this.itemForm.componentId) return;
      const comp = this.componentsList.find(c => c.id === this.itemForm.componentId);
      if (comp) {
        if (!this.itemForm.name || this.itemForm.name === '') {
          this.itemForm.name = comp.name;
        }
        this.itemForm.unitCost = comp.estimatedUnitCost || 0;
        this.itemForm.unitCode = comp.unitCode || 'PCS';
      }
    },

    async saveItem() {
      try {
        if (this.editingItemId) {
          await this.updateItem(this.editingItemId, this.itemForm);
        } else {
          await this.addItem(this.bom.id, this.itemForm);
        }
        this.closeItemModal();
      } catch (err) {
        alert(err.response?.data?.message || 'Failed to save BOM item');
      }
    },

    async confirmDeleteItem(item) {
      if (confirm(`Are you sure you want to remove '${item.name}' from the BOM?`)) {
        try {
          await this.deleteItem(item.id);
        } catch (err) {
          alert(err.response?.data?.message || 'Failed to delete item');
        }
      }
    },

    openSettingsModal() {
      if (!this.bom) return;
      this.settingsForm = {
        name: this.bom.name,
        targetMargin: this.bom.targetMargin || 35.0,
        status: this.bom.status || 'Active',
        description: this.bom.description || ''
      };
      this.showSettingsModal = true;
    },

    async saveBOMSettings() {
      try {
        await this.updateBOM(this.bom.id, this.settingsForm);
        this.showSettingsModal = false;
      } catch (err) {
        alert(err.response?.data?.message || 'Failed to update settings');
      }
    }
  }
};
</script>
