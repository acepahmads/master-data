<template>
  <div v-if="prototype" class="space-y-5">
    <!-- PROTOTYPE HEADER & COMMAND WORKSPACE -->
    <div class="app-card p-5 bg-white shadow-card space-y-4">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
        <div class="flex items-start gap-3.5">
          <div class="w-12 h-12 rounded-xl bg-indigo-600 text-white flex items-center justify-center shadow-md shadow-indigo-500/20 shrink-0 font-mono font-bold text-base">
            #{{ prototype.buildNumber }}
          </div>
          <div>
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-indigo-50 text-indigo-700 border border-indigo-200">
                {{ prototype.code }}
              </span>
              <span
                class="text-3xs font-mono font-bold px-2.5 py-0.5 rounded-full border shadow-2xs"
                :class="getStatusBadgeClass(prototype.status)"
              >
                {{ formatStatus(prototype.status) }}
              </span>
              <span class="text-3xs font-mono px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200 font-bold">
                {{ prototype.quantity }} Units
              </span>
              <span class="text-3xs font-medium px-2 py-0.5 rounded-full bg-purple-50 text-purple-700 border border-purple-200">
                {{ prototype.purpose }}
              </span>
            </div>
            <h1 class="text-xl font-bold text-slate-900 tracking-tight">{{ prototype.name }}</h1>
            
            <div class="flex flex-wrap items-center gap-3 text-xs text-slate-500 mt-1 font-mono">
              <router-link :to="`/products/${prototype.productId}`" class="text-brand-600 hover:underline font-bold">
                {{ prototype.productCode }} ({{ prototype.productName }})
              </router-link>
              <span>➔</span>
              <router-link :to="`/versions/${prototype.versionId}`" class="text-brand-600 hover:underline">
                {{ prototype.versionCode }}
              </router-link>
              <span>➔</span>
              <router-link :to="`/revisions/${prototype.hardwareRevisionId}`" class="text-emerald-700 font-bold hover:underline">
                {{ prototype.hardwareRevisionCode }} ({{ prototype.hardwareRevisionName }})
              </router-link>
            </div>
          </div>
        </div>

        <!-- Header Actions -->
        <div class="flex flex-wrap items-center gap-2 shrink-0">
          <router-link to="/prototypes/list" class="btn btn-secondary text-xs">
            <AppIcon name="chevron-left" size="xs" class="mr-1" />
            <span>All Builds</span>
          </router-link>

          <!-- Status Advance Button -->
          <button
            @click="statusModalOpen = true"
            class="btn btn-primary text-xs bg-indigo-600 hover:bg-indigo-700 border-none shadow-sm flex items-center gap-1.5"
          >
            <AppIcon name="check" size="xs" />
            <span>Advance Status</span>
          </button>
        </div>
      </div>

      <!-- Quick Metadata Grid -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs font-mono pt-2 border-t border-slate-100">
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Engineering Owner</span>
          <span class="font-bold text-slate-800 text-xs mt-0.5 block font-sans">{{ prototype.engineeringOwner }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Target Build Date</span>
          <span class="font-bold text-slate-900 text-xs mt-0.5 block">{{ prototype.targetBuildDate }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Frozen BOM Cost</span>
          <span class="font-bold text-emerald-700 text-xs mt-0.5 block">${{ Number(prototype.bomSnapshot.totalEstimatedCost || 0).toFixed(2) }} USD</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Assembly Progress</span>
          <span class="font-bold text-indigo-600 text-xs mt-0.5 block">{{ completedStagesCount }} / 8 Stages ({{ assemblyProgressPercent }}%)</span>
        </div>
      </div>

      <!-- 6 WORKSPACE TABS HEADER -->
      <div class="flex items-center gap-4 sm:gap-6 border-t border-slate-100 pt-3 text-xs overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'pb-2 font-semibold border-b-2 transition-all flex items-center gap-1.5 whitespace-nowrap',
            activeTab === tab.id
              ? 'border-indigo-600 text-indigo-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon :name="tab.icon" size="xs" />
          <span>{{ tab.name }}</span>
          <span
            v-if="tab.badge"
            class="text-3xs font-mono px-1.5 py-0.2 rounded-full"
            :class="activeTab === tab.id ? 'bg-indigo-50 text-indigo-700 font-bold' : 'bg-slate-100 text-slate-500'"
          >
            {{ tab.badge }}
          </span>
        </button>
      </div>
    </div>

    <!-- TAB 1: OVERVIEW -->
    <div v-if="activeTab === 'overview'" class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Left 2 Cols: Scope & Summary -->
      <div class="lg:col-span-2 space-y-5">
        <!-- OBJECTIVES CARD -->
        <div class="app-card p-5 bg-white space-y-3 shadow-card">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100">
            <AppIcon name="target" size="xs" class="text-indigo-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Build Objective & Scope</h2>
          </div>
          <p class="text-xs text-slate-700 leading-relaxed">
            {{ prototype.objective || 'No specific objective provided for this prototype build run.' }}
          </p>
        </div>

        <!-- REVISION & BOM LINKAGE -->
        <div class="app-card p-5 bg-white space-y-3 shadow-card">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AppIcon name="layers" size="xs" class="text-emerald-600" />
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Source Revision Baseline</h2>
            </div>
            <router-link :to="`/revisions/${prototype.hardwareRevisionId}/bom`" class="text-3xs text-brand-600 hover:underline font-bold">
              View Active Engineering BOM ➔
            </router-link>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs font-mono">
            <div class="p-3 rounded-lg bg-slate-50 border border-slate-200">
              <span class="text-3xs text-slate-400 block uppercase">Revision Designator:</span>
              <span class="font-bold text-slate-800 text-sm mt-0.5 block">{{ prototype.hardwareRevisionCode }}</span>
              <span class="text-3xs text-slate-500 font-sans block mt-1">{{ prototype.hardwareRevisionName }}</span>
            </div>
            <div class="p-3 rounded-lg bg-slate-50 border border-slate-200">
              <span class="text-3xs text-slate-400 block uppercase">Frozen BOM Snapshot:</span>
              <span class="font-bold text-indigo-700 text-sm mt-0.5 block">{{ prototype.bomSnapshot.snapshotId }}</span>
              <span class="text-3xs text-slate-500 block mt-1">Frozen: {{ formatDate(prototype.bomSnapshot.frozenAt) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right 1 Col: Key Build Metrics -->
      <div class="space-y-5">
        <div class="app-card p-5 bg-white space-y-4 shadow-card">
          <div class="flex items-center gap-2 pb-3 border-b border-slate-100">
            <AppIcon name="cpu" size="xs" class="text-indigo-600" />
            <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Build Execution Metrics</h3>
          </div>

          <div class="space-y-3 text-xs font-mono">
            <div class="flex items-center justify-between pb-2 border-b border-slate-100">
              <span class="text-slate-500">Planned Build Units:</span>
              <span class="font-bold text-slate-900 text-sm">{{ prototype.quantity }} Units</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-slate-100">
              <span class="text-slate-500">Component Readiness:</span>
              <span class="font-bold text-emerald-600">{{ readinessPercent }}% Ready</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-slate-100">
              <span class="text-slate-500">Assembly Stages Done:</span>
              <span class="font-bold text-indigo-600">{{ completedStagesCount }} / 8 Stages</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-slate-500">Unit Snapshot BOM:</span>
              <span class="font-bold text-slate-900">${{ Number(prototype.bomSnapshot.totalEstimatedCost || 0).toFixed(2) }}</span>
            </div>
          </div>

          <button @click="activeTab = 'assembly'" class="btn btn-secondary w-full text-xs font-bold">
            Go to Assembly Checklist ➔
          </button>
        </div>
      </div>
    </div>

    <!-- TAB 2: BOM SNAPSHOT VIEW (READ ONLY) -->
    <div v-if="activeTab === 'bom'" class="app-card p-5 bg-white space-y-4 shadow-card">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
        <div>
          <div class="flex items-center gap-2">
            <AppIcon name="lock" size="xs" class="text-amber-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Immutable BOM Snapshot</h2>
          </div>
          <p class="text-3xs text-slate-400">Frozen component and cost baseline captured at build creation (Read-Only).</p>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-3xs font-mono font-bold px-2.5 py-1 rounded-full bg-amber-50 text-amber-800 border border-amber-300 flex items-center gap-1 shadow-2xs">
            <AppIcon name="lock" size="2xs" />
            <span>BOM SNAPSHOT — READ ONLY</span>
          </span>
        </div>
      </div>

      <!-- Snapshot Info Box -->
      <div class="p-3 bg-amber-50/50 rounded-xl border border-amber-200/70 flex flex-wrap items-center justify-between gap-3 text-3xs font-mono text-amber-900">
        <div>
          <strong class="font-bold">Snapshot ID:</strong> {{ prototype.bomSnapshot.snapshotId }}
        </div>
        <div>
          <strong class="font-bold">Frozen At:</strong> {{ formatDate(prototype.bomSnapshot.frozenAt) }}
        </div>
        <div>
          <strong class="font-bold">Hardware Revision:</strong> {{ prototype.hardwareRevisionCode }}
        </div>
        <div>
          <strong class="font-bold">Total Estimated Unit BOM:</strong> ${{ Number(prototype.bomSnapshot.totalEstimatedCost || 0).toFixed(2) }} USD
        </div>
      </div>

      <!-- Table of Frozen Line Items -->
      <div class="overflow-x-auto border border-slate-200 rounded-xl">
        <table class="app-table w-full text-left text-xs">
          <thead class="bg-slate-50 text-3xs font-mono text-slate-600 uppercase border-b border-slate-200">
            <tr>
              <th class="py-3 px-3.5">Component / MPN</th>
              <th class="py-3 px-3.5">Category</th>
              <th class="py-3 px-3.5 text-center">Qty / Board</th>
              <th class="py-3 px-3.5">Ref Designators</th>
              <th class="py-3 px-3.5 text-right">Unit Cost</th>
              <th class="py-3 px-3.5 text-right">Line Total</th>
              <th class="py-3 px-3.5 text-center">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 font-mono text-3xs">
            <tr v-for="itm in prototype.bomSnapshot.items" :key="itm.id" class="hover:bg-slate-50/80">
              <td class="py-2.5 px-3.5">
                <div class="font-bold text-slate-900 font-sans text-xs">{{ itm.componentName }}</div>
                <div class="text-slate-400">{{ itm.mpn }}</div>
              </td>
              <td class="py-2.5 px-3.5 font-sans">{{ itm.category }}</td>
              <td class="py-2.5 px-3.5 text-center font-bold text-slate-800">{{ itm.quantity }} {{ itm.uom }}</td>
              <td class="py-2.5 px-3.5 text-indigo-700 font-bold">{{ itm.refDesignators }}</td>
              <td class="py-2.5 px-3.5 text-right text-slate-700">${{ Number(itm.unitCost || 0).toFixed(2) }}</td>
              <td class="py-2.5 px-3.5 text-right font-bold text-slate-900">${{ Number(itm.lineCost || 0).toFixed(2) }}</td>
              <td class="py-2.5 px-3.5 text-center">
                <span class="px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 border border-slate-200 text-3xs font-bold">
                  FROZEN
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 3: COMPONENT PREPARATION & READINESS -->
    <div v-if="activeTab === 'preparation'" class="app-card p-5 bg-white space-y-4 shadow-card">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
        <div>
          <div class="flex items-center gap-2">
            <AppIcon name="component" size="xs" class="text-indigo-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Component Preparation & Lab Readiness</h2>
          </div>
          <p class="text-3xs text-slate-400">Track part availability and physical allocation for {{ prototype.quantity }} build units.</p>
        </div>
        <div class="flex items-center gap-3 font-mono text-xs">
          <span class="text-slate-500">Overall Readiness:</span>
          <span class="font-bold text-emerald-600">{{ readinessPercent }}% Complete</span>
        </div>
      </div>

      <!-- Readiness Progress Bar -->
      <div class="w-full h-2 rounded-full bg-slate-100 overflow-hidden border border-slate-200">
        <div
          class="h-full rounded-full bg-emerald-500 transition-all duration-500"
          :style="{ width: `${readinessPercent}%` }"
        ></div>
      </div>

      <!-- Component Preparation Table -->
      <div class="overflow-x-auto border border-slate-200 rounded-xl">
        <table class="app-table w-full text-left text-xs">
          <thead class="bg-slate-50 text-3xs font-mono text-slate-600 uppercase border-b border-slate-200">
            <tr>
              <th class="py-3 px-3.5">Component & MPN</th>
              <th class="py-3 px-3.5 text-center">Required Qty</th>
              <th class="py-3 px-3.5 text-center">Available Qty</th>
              <th class="py-3 px-3.5 text-center">Missing</th>
              <th class="py-3 px-3.5 text-center">Readiness Status</th>
              <th class="py-3 px-3.5">Lab Bench Notes</th>
              <th class="py-3 px-3.5 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-3xs font-mono">
            <tr v-for="c in prototype.componentsReadiness" :key="c.componentId" class="hover:bg-slate-50">
              <td class="py-2.5 px-3.5">
                <div class="font-bold text-slate-900 font-sans text-xs">{{ c.componentName }}</div>
                <div class="text-slate-400">{{ c.mpn }}</div>
              </td>
              <td class="py-2.5 px-3.5 text-center font-bold text-slate-800">{{ c.requiredQty }}</td>
              <td class="py-2.5 px-3.5 text-center font-bold text-emerald-700">{{ c.availableQty }}</td>
              <td class="py-2.5 px-3.5 text-center font-bold" :class="c.missingQty > 0 ? 'text-red-600' : 'text-slate-400'">
                {{ c.missingQty || 0 }}
              </td>
              <td class="py-2.5 px-3.5 text-center">
                <span
                  class="px-2 py-0.5 rounded-full border text-3xs font-bold"
                  :class="getReadinessBadgeClass(c.readiness)"
                >
                  {{ c.readiness }}
                </span>
              </td>
              <td class="py-2.5 px-3.5 font-sans text-slate-600">{{ c.notes || '—' }}</td>
              <td class="py-2.5 px-3.5 text-right">
                <button
                  type="button"
                  @click="toggleComponentReadiness(c)"
                  class="btn btn-secondary py-0.5 px-2 text-3xs font-bold"
                >
                  Toggle Ready
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 4: ASSEMBLY PROGRESS (CHECKLIST) -->
    <div v-if="activeTab === 'assembly'" class="app-card p-5 bg-white space-y-4 shadow-card">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
        <div>
          <div class="flex items-center gap-2">
            <AppIcon name="tool" size="xs" class="text-indigo-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Multi-Stage Assembly Progress</h2>
          </div>
          <p class="text-3xs text-slate-400">Sequential engineering workflow from PCB prep to power-on validation.</p>
        </div>
        <span class="text-3xs font-mono font-bold px-2.5 py-1 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
          {{ completedStagesCount }} / 8 STAGES COMPLETED ({{ assemblyProgressPercent }}%)
        </span>
      </div>

      <!-- Assembly Stages Vertical Flow -->
      <div class="space-y-3 relative pl-4 before:absolute before:left-2 before:top-3 before:bottom-3 before:w-0.5 before:bg-slate-200">
        <div
          v-for="stage in prototype.assemblyStages"
          :key="stage.id"
          class="p-3.5 rounded-xl border transition-all relative flex flex-col sm:flex-row sm:items-center justify-between gap-3"
          :class="[
            stage.status === 'COMPLETED' ? 'bg-emerald-50/40 border-emerald-200' :
            stage.status === 'IN_PROGRESS' ? 'bg-indigo-50/50 border-indigo-300 ring-2 ring-indigo-500/20' :
            'bg-slate-50/50 border-slate-200'
          ]"
        >
          <!-- Stage Node Dot -->
          <span
            class="absolute -left-4 top-4.5 w-3 h-3 rounded-full ring-4 ring-white shadow-2xs"
            :class="[
              stage.status === 'COMPLETED' ? 'bg-emerald-500' :
              stage.status === 'IN_PROGRESS' ? 'bg-indigo-600 animate-pulse' :
              'bg-slate-300'
            ]"
          ></span>

          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 mb-0.5">
              <span class="font-mono text-3xs font-bold px-1.5 py-0.2 rounded bg-white text-slate-700 border border-slate-200">
                STAGE 0{{ stage.stepNumber }}
              </span>
              <h3 class="font-bold text-xs text-slate-900">{{ stage.name }}</h3>
              <span
                class="text-3xs font-mono font-bold px-2 py-0.2 rounded-full border"
                :class="stage.status === 'COMPLETED' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : stage.status === 'IN_PROGRESS' ? 'bg-indigo-50 text-indigo-700 border-indigo-200' : 'bg-slate-100 text-slate-500 border-slate-200'"
              >
                {{ stage.status }}
              </span>
            </div>
            <p class="text-3xs text-slate-500 leading-snug">{{ stage.description }}</p>
            <div v-if="stage.notes" class="text-3xs font-mono text-indigo-700 mt-1">
              Note: {{ stage.notes }}
            </div>
          </div>

          <div class="flex items-center gap-3 shrink-0 text-3xs font-mono">
            <div class="text-right">
              <span class="text-slate-400 block">Assigned: {{ stage.engineer }}</span>
              <span v-if="stage.completedDate" class="text-emerald-700 font-bold block">Done: {{ stage.completedDate }}</span>
            </div>
            <button
              type="button"
              @click="toggleStage(stage.id)"
              class="btn btn-secondary py-1 px-2.5 text-3xs font-bold"
            >
              {{ stage.status === 'COMPLETED' ? 'Reopen' : stage.status === 'IN_PROGRESS' ? 'Mark Completed' : 'Start Stage' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 5: ENGINEERING NOTES -->
    <div v-if="activeTab === 'notes'" class="app-card p-5 bg-white space-y-4 shadow-card">
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div>
          <div class="flex items-center gap-2">
            <AppIcon name="file-text" size="xs" class="text-indigo-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Prototype Engineering Notes & Logs</h2>
          </div>
          <p class="text-3xs text-slate-400">Bench observations, design notes, rework observations, and prototype recommendations.</p>
        </div>
        <button
          type="button"
          @click="noteModalOpen = true"
          class="btn btn-primary text-xs flex items-center gap-1 bg-indigo-600 hover:bg-indigo-700"
        >
          <AppIcon name="plus" size="xs" />
          <span>Add Note</span>
        </button>
      </div>

      <!-- Notes Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div
          v-for="n in prototype.engineeringNotes"
          :key="n.id"
          class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-2 hover:bg-white hover:border-indigo-300 transition-all"
        >
          <div class="flex items-center justify-between">
            <span
              class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full border"
              :class="getNoteCategoryBadgeClass(n.category)"
            >
              {{ n.category }}
            </span>
            <span class="text-3xs font-mono text-slate-400">{{ formatDate(n.timestamp) }}</span>
          </div>
          <h4 class="font-bold text-xs text-slate-900">{{ n.title }}</h4>
          <p class="text-xs text-slate-600 leading-relaxed font-sans">{{ n.content }}</p>
          <div class="text-3xs font-mono text-slate-400 pt-2 border-t border-slate-100">
            Author: <strong class="text-slate-700">{{ n.author }}</strong>
          </div>
        </div>

        <div v-if="prototype.engineeringNotes.length === 0" class="col-span-2 py-8 text-center text-slate-400 text-xs">
          No engineering notes logged yet for this prototype build.
        </div>
      </div>
    </div>

    <!-- TAB 6: ACTIVITY HISTORY -->
    <div v-if="activeTab === 'history'" class="app-card p-5 bg-white space-y-4 shadow-card">
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div class="flex items-center gap-2">
          <AppIcon name="activity" size="xs" class="text-indigo-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Prototype Lifecycle Activity Log</h2>
        </div>
        <span class="text-3xs font-mono text-slate-400">{{ prototype.activityHistory.length }} Events</span>
      </div>

      <div class="relative pl-4 space-y-4 before:absolute before:left-1.5 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
        <div
          v-for="act in prototype.activityHistory"
          :key="act.id"
          class="relative text-xs pl-2"
        >
          <span class="absolute -left-4 top-1 w-2.5 h-2.5 rounded-full bg-indigo-600 ring-4 ring-white shadow-2xs"></span>
          <div class="flex items-center justify-between gap-2">
            <span class="font-bold text-slate-900">{{ act.action }}</span>
            <span class="text-3xs font-mono text-slate-400">{{ act.timestamp }}</span>
          </div>
          <p class="text-3xs text-slate-500 mt-0.5 leading-snug">{{ act.description }}</p>
          <div class="text-3xs font-mono text-slate-400 mt-1">Author: {{ act.author }}</div>
        </div>
      </div>
    </div>

    <!-- TAB 7: TESTING & VALIDATION (PHASE 5) -->
    <div v-if="activeTab === 'testing'" class="space-y-4">
      <div class="app-card p-5 bg-white space-y-4 border border-slate-200">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div>
            <h3 class="text-sm font-bold text-slate-900 flex items-center gap-2">
              <AppIcon name="shield" size="sm" class="text-emerald-600" />
              Prototype Test Sessions & Validation Governance
            </h3>
            <p class="text-3xs text-slate-500 mt-0.5">
              Parametric test execution runs against frozen immutable test plan snapshots.
            </p>
          </div>

          <button
            @click="openNewSessionModal"
            class="btn btn-primary text-xs bg-emerald-600 hover:bg-emerald-700 flex items-center gap-1.5"
          >
            <AppIcon name="plus" size="xs" />
            Initiate Test Session
          </button>
        </div>

        <div v-if="sessionsLoading" class="py-8 text-center text-xs text-slate-400">
          Loading test sessions...
        </div>

        <div v-else-if="prototypeSessions.length === 0" class="py-8 text-center space-y-2">
          <AppIcon name="shield" size="lg" class="text-slate-300 mx-auto" />
          <div class="text-xs font-semibold text-slate-700">No Test Sessions Initiated</div>
          <p class="text-3xs text-slate-500 max-w-sm mx-auto">
            Assign a released master test plan to freeze an immutable test snapshot and begin laboratory verification.
          </p>
          <button
            @click="openNewSessionModal"
            class="btn btn-secondary text-xs mt-2"
          >
            Initiate First Test Session
          </button>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="data-table w-full text-left text-xs">
            <thead>
              <tr class="bg-slate-50 text-slate-500 font-semibold border-b border-slate-200">
                <th class="py-2.5 px-3">Session Code</th>
                <th class="py-2.5 px-3">Session Name & Plan</th>
                <th class="py-2.5 px-3 text-center">Cases</th>
                <th class="py-2.5 px-3 text-center">Pass Rate</th>
                <th class="py-2.5 px-3 text-center">Status</th>
                <th class="py-2.5 px-3 text-right">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="sess in prototypeSessions" :key="sess.id" class="hover:bg-slate-50/80 transition">
                <td class="py-3 px-3 font-mono font-bold text-slate-800 text-xs">
                  {{ sess.sessionCode }}
                </td>
                <td class="py-3 px-3">
                  <div class="font-bold text-slate-900">{{ sess.name }}</div>
                  <div class="text-3xs text-slate-500">Plan ID: {{ sess.testPlanId }}</div>
                </td>
                <td class="py-3 px-3 text-center font-mono font-semibold">
                  <span class="text-emerald-600">{{ sess.passedCases }}</span> / {{ sess.totalCases }}
                </td>
                <td class="py-3 px-3 text-center font-mono">
                  <span
                    class="badge text-3xs font-bold"
                    :class="(sess.passRatePercentage || 0) >= 100 ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                  >
                    {{ (sess.passRatePercentage || 0).toFixed(0) }}%
                  </span>
                </td>
                <td class="py-3 px-3 text-center">
                  <span
                    class="badge text-3xs font-bold"
                    :class="sess.status === 'COMPLETED' ? 'bg-emerald-100 text-emerald-800' : 'bg-indigo-100 text-indigo-800'"
                  >
                    {{ sess.status }}
                  </span>
                </td>
                <td class="py-3 px-3 text-right">
                  <router-link
                    :to="`/test-sessions/${sess.id}`"
                    class="btn btn-secondary text-3xs font-semibold px-2.5 py-1"
                  >
                    Open Workspace &rarr;
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- INITIATE TEST SESSION MODAL -->
    <div v-if="newSessionModalOpen" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50">
      <div class="app-card p-5 bg-white max-w-lg w-full space-y-4 shadow-xl">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <h3 class="text-sm font-bold text-slate-900">Initiate Laboratory Test Session</h3>
          <button @click="newSessionModalOpen = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <form @submit.prevent="submitCreateSession" class="space-y-3 text-xs">
          <div>
            <label class="form-label">Select Master Test Plan *</label>
            <select v-model="sessionForm.testPlanId" @change="onPlanSelect" required class="form-select w-full">
              <option value="" disabled>Choose a Master Test Plan...</option>
              <option v-for="p in availablePlans" :key="p.id" :value="p.id">
                {{ p.code }} — {{ p.name }}
              </option>
            </select>
          </div>

          <div>
            <label class="form-label">Select Released Plan Version *</label>
            <select v-model="sessionForm.testPlanVersionId" required class="form-select w-full font-mono">
              <option value="" disabled>Choose a version...</option>
              <option v-for="v in availableVersions" :key="v.id" :value="v.id">
                {{ v.versionNumber }} ({{ v.status }}) — {{ v.testCases ? v.testCases.length : 0 }} Test Cases
              </option>
            </select>
          </div>

          <div>
            <label class="form-label">Session Name *</label>
            <input
              v-model="sessionForm.name"
              type="text"
              required
              placeholder="e.g. EVT Initial Power & RF Verification"
              class="form-input w-full"
            />
          </div>

          <div>
            <label class="form-label">Session Scope & Notes</label>
            <textarea
              v-model="sessionForm.description"
              rows="2"
              placeholder="Test session objectives and conditions..."
              class="form-textarea w-full"
            ></textarea>
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-100">
            <button type="button" @click="newSessionModalOpen = false" class="btn btn-secondary text-xs">Cancel</button>
            <button type="submit" class="btn btn-primary text-xs bg-emerald-600 hover:bg-emerald-700">
              Freeze Snapshot & Start
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- STATUS ADVANCE MODAL -->
    <div v-if="statusModalOpen" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50">
      <div class="app-card p-5 bg-white max-w-md w-full space-y-4 shadow-xl">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <h3 class="text-sm font-bold text-slate-900">Advance Prototype Build Status</h3>
          <button @click="statusModalOpen = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <div class="space-y-3 text-xs">
          <label class="form-label">Select Target Build Status:</label>
          <div class="space-y-1.5">
            <button
              v-for="st in statusOptions"
              :key="st"
              type="button"
              @click="advanceStatus(st)"
              class="w-full text-left p-2.5 rounded-lg border font-mono text-xs transition-all flex items-center justify-between"
              :class="prototype.status === st ? 'bg-indigo-50 border-indigo-500 font-bold text-indigo-900' : 'bg-slate-50 border-slate-200 hover:bg-white'"
            >
              <span>{{ formatStatus(st) }}</span>
              <AppIcon v-if="prototype.status === st" name="check" size="xs" class="text-indigo-600" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ADD NOTE MODAL -->
    <div v-if="noteModalOpen" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50">
      <div class="app-card p-5 bg-white max-w-lg w-full space-y-4 shadow-xl">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <h3 class="text-sm font-bold text-slate-900">Add Engineering Observation Note</h3>
          <button @click="noteModalOpen = false" class="text-slate-400 hover:text-slate-700">✕</button>
        </div>

        <div class="space-y-3 text-xs">
          <div>
            <label class="form-label">Category</label>
            <select v-model="newNote.category" class="form-select w-full">
              <option value="GENERAL">GENERAL NOTE</option>
              <option value="ASSEMBLY_ISSUE">ASSEMBLY ISSUE</option>
              <option value="DESIGN_OBSERVATION">DESIGN OBSERVATION</option>
              <option value="KNOWN_ISSUE">KNOWN ISSUE</option>
              <option value="RECOMMENDATION">RECOMMENDATION</option>
            </select>
          </div>
          <div>
            <label class="form-label">Title</label>
            <input v-model="newNote.title" type="text" placeholder="e.g. Solder Bridging Risk on QFN-32" class="form-input w-full" />
          </div>
          <div>
            <label class="form-label">Observation Details</label>
            <textarea v-model="newNote.content" rows="3" placeholder="Describe the engineering observation, bench reading, or issue..." class="form-textarea w-full"></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-100">
          <button @click="noteModalOpen = false" class="btn btn-secondary text-xs">Cancel</button>
          <button @click="submitNote" class="btn btn-primary text-xs bg-indigo-600 hover:bg-indigo-700">Save Note</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { usePrototypeStore } from '@/stores/prototypeStore';
import { useMasterStore } from '@/stores/masterStore';
import { useTestingStore } from '@/stores/testingStore';
import { testingApi } from '@/api/testing';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'PrototypeDetailView',
  components: {
    AppIcon
  },
  data() {
    return {
      activeTab: 'overview',
      statusModalOpen: false,
      noteModalOpen: false,
      newSessionModalOpen: false,
      sessionsLoading: false,
      prototypeSessions: [],
      availablePlans: [],
      availableVersions: [],
      sessionForm: {
        testPlanId: '',
        testPlanVersionId: '',
        name: '',
        description: ''
      },
      newNote: {
        category: 'GENERAL',
        title: '',
        content: ''
      },
      statusOptions: [
        'DRAFT',
        'PLANNING',
        'COMPONENT_PREPARATION',
        'IN_ASSEMBLY',
        'ASSEMBLY_COMPLETE',
        'READY_FOR_TESTING',
        'COMPLETED',
        'CANCELLED'
      ]
    };
  },
  computed: {
    protoStore() {
      return usePrototypeStore();
    },
    masterStore() {
      return useMasterStore();
    },
    testingStore() {
      return useTestingStore();
    },
    prototypeId() {
      return this.$route.params.id;
    },
    prototype() {
      return this.protoStore.getPrototypeById(this.prototypeId);
    },
    tabs() {
      if (!this.prototype) return [];
      return [
        { id: 'overview', name: 'Overview', icon: 'file-text' },
        { id: 'bom', name: 'BOM Snapshot', icon: 'lock', badge: `${this.prototype.bomSnapshot.totalItems} Parts` },
        { id: 'preparation', name: 'Component Prep', icon: 'component', badge: `${this.readinessPercent}%` },
        { id: 'assembly', name: 'Assembly Progress', icon: 'tool', badge: `${this.completedStagesCount}/8` },
        { id: 'notes', name: 'Engineering Notes', icon: 'edit', badge: `${this.prototype.engineeringNotes.length}` },
        { id: 'testing', name: 'Testing & Validation', icon: 'shield', badge: `${this.prototypeSessions.length} Sessions` },
        { id: 'history', name: 'Activity History', icon: 'activity' }
      ];
    },
    completedStagesCount() {
      if (!this.prototype || !this.prototype.assemblyStages) return 0;
      return this.prototype.assemblyStages.filter(s => s.status === 'COMPLETED').length;
    },
    assemblyProgressPercent() {
      if (!this.prototype || !this.prototype.assemblyStages || this.prototype.assemblyStages.length === 0) return 0;
      return Math.round((this.completedStagesCount / this.prototype.assemblyStages.length) * 100);
    },
    readinessPercent() {
      if (!this.prototype || !this.prototype.componentsReadiness || this.prototype.componentsReadiness.length === 0) return 100;
      const ready = this.prototype.componentsReadiness.filter(c => c.readiness === 'AVAILABLE').length;
      return Math.round((ready / this.prototype.componentsReadiness.length) * 100);
    }
  },
  methods: {
    formatStatus(status) {
      if (!status) return 'DRAFT';
      return status.replace(/_/g, ' ');
    },
    formatDate(dateStr) {
      if (!dateStr) return '—';
      return dateStr.substring(0, 10);
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'COMPLETED':
          return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'READY_FOR_TESTING':
          return 'bg-purple-50 text-purple-700 border-purple-200';
        case 'IN_ASSEMBLY':
        case 'ASSEMBLY_COMPLETE':
          return 'bg-blue-50 text-blue-700 border-blue-200';
        case 'COMPONENT_PREPARATION':
          return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'CANCELLED':
          return 'bg-red-50 text-red-700 border-red-200';
        default:
          return 'bg-slate-100 text-slate-600 border-slate-200';
      }
    },
    getReadinessBadgeClass(readiness) {
      switch (readiness) {
        case 'AVAILABLE':
          return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'PARTIAL':
          return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'MISSING':
          return 'bg-red-50 text-red-700 border-red-200';
        case 'SUBSTITUTE_REQUIRED':
          return 'bg-purple-50 text-purple-700 border-purple-200';
        default:
          return 'bg-slate-100 text-slate-600 border-slate-200';
      }
    },
    getNoteCategoryBadgeClass(cat) {
      switch (cat) {
        case 'ASSEMBLY_ISSUE':
          return 'bg-red-50 text-red-700 border-red-200';
        case 'DESIGN_OBSERVATION':
          return 'bg-blue-50 text-blue-700 border-blue-200';
        case 'RECOMMENDATION':
          return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'KNOWN_ISSUE':
          return 'bg-amber-50 text-amber-700 border-amber-200';
        default:
          return 'bg-slate-100 text-slate-700 border-slate-200';
      }
    },
    async advanceStatus(st) {
      try {
        await this.protoStore.updatePrototypeStatus(this.prototype.id, st);
        this.statusModalOpen = false;
        this.masterStore.showToast('Status Updated', `Prototype status changed to ${st.replace(/_/g, ' ')}`, 'success');
      } catch (err) {
        alert('Failed to update status: ' + (err.response?.data?.message || err.message));
      }
    },
    async toggleStage(stageId) {
      const stage = this.prototype.assemblyStages.find(s => s.id === stageId || s.stage === stageId);
      if (!stage) return;
      const nextStatus = stage.status === 'COMPLETED' ? 'NOT_STARTED' : stage.status === 'IN_PROGRESS' ? 'COMPLETED' : 'IN_PROGRESS';
      try {
        await this.protoStore.updateAssemblyStage(this.prototype.id, stage.id, {
          status: nextStatus,
          performedBy: this.prototype.engineeringOwner || 'Engineer',
          notes: stage.notes || ''
        });
      } catch (err) {
        alert('Assembly Stage update failed: ' + (err.response?.data?.message || err.message));
      }
    },
    async toggleComponentReadiness(comp) {
      const targetAvail = comp.availableQuantity >= comp.requiredQuantity ? 0 : comp.requiredQuantity;
      try {
        await this.protoStore.updateComponentPreparation(this.prototype.id, comp.id, {
          availableQuantity: targetAvail,
          notes: comp.notes || ''
        });
      } catch (err) {
        alert('Component preparation update failed: ' + (err.response?.data?.message || err.message));
      }
    },
    async submitNote() {
      if (!this.newNote.title || !this.newNote.content) {
        alert('Please fill note title and content.');
        return;
      }
      try {
        await this.protoStore.addEngineeringNote(this.prototype.id, {
          category: this.newNote.category,
          title: this.newNote.title,
          content: this.newNote.content
        });
        this.newNote = { category: 'GENERAL', title: '', content: '' };
        this.noteModalOpen = false;
        this.masterStore.showToast('Note Added', 'Engineering observation saved.', 'success');
      } catch (err) {
        alert('Failed to add note: ' + (err.response?.data?.message || err.message));
      }
    },
    async loadPrototypeSessions() {
      this.sessionsLoading = true;
      try {
        const res = await testingApi.getPrototypeSessions(this.prototypeId);
        if (res.data?.success && res.data.data) {
          this.prototypeSessions = res.data.data;
        }
      } catch (err) {
        console.error('Failed to load prototype test sessions:', err);
      } finally {
        this.sessionsLoading = false;
      }
    },
    async openNewSessionModal() {
      try {
        const res = await testingApi.getTestPlans({ limit: 50 });
        if (res.data?.success && res.data.data) {
          this.availablePlans = res.data.data;
        }
      } catch (err) {
        console.error(err);
      }
      this.sessionForm = {
        testPlanId: '',
        testPlanVersionId: '',
        name: `${this.prototype?.code || 'PROTO'} Verification Session`,
        description: ''
      };
      this.availableVersions = [];
      this.newSessionModalOpen = true;
    },
    onPlanSelect() {
      const selected = this.availablePlans.find(p => p.id === this.sessionForm.testPlanId);
      if (selected && selected.versions) {
        this.availableVersions = selected.versions.filter(v => v.status === 'RELEASED');
        if (this.availableVersions.length > 0) {
          this.sessionForm.testPlanVersionId = this.availableVersions[0].id;
        } else if (selected.versions.length > 0) {
          this.availableVersions = selected.versions;
          this.sessionForm.testPlanVersionId = selected.versions[0].id;
        }
      } else {
        this.availableVersions = [];
      }
    },
    async submitCreateSession() {
      try {
        const res = await testingApi.createPrototypeSession(this.prototypeId, this.sessionForm);
        this.newSessionModalOpen = false;
        this.masterStore.showToast('Session Initiated', 'Frozen test plan snapshot created.', 'success');
        await this.loadPrototypeSessions();
        if (res.data?.data?.id) {
          this.$router.push(`/test-sessions/${res.data.data.id}`);
        }
      } catch (err) {
        alert(err.response?.data?.message || err.message);
      }
    }
  },
  async mounted() {
    try {
      await Promise.all([
        this.protoStore.fetchPrototype(this.prototypeId),
        this.loadPrototypeSessions()
      ]);
    } catch (err) {
      console.error('[PrototypeDetailView] mounted error:', err);
    }
  }
};
</script>
