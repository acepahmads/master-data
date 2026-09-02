<template>
  <div class="space-y-4">
    <div
      v-for="(chain, idx) in chains"
      :key="idx"
      class="p-5 rounded-2xl border border-slate-800 bg-slate-900/60 space-y-4"
    >
      <!-- Chain Header -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <div class="flex items-center gap-3">
          <span class="font-mono text-xs font-bold text-cyan-400 bg-cyan-950/60 px-2.5 py-1 rounded-lg border border-cyan-800/40">
            {{ chain.ecrCode }}
          </span>
          <span class="text-sm font-semibold text-slate-200">{{ chain.ecrTitle }}</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-xs text-slate-400">Status:</span>
          <span class="px-2 py-0.5 rounded text-[11px] font-bold uppercase bg-slate-800 text-slate-300 border border-slate-700">
            {{ chain.ecrStatus }}
          </span>
        </div>
      </div>

      <!-- Closed-Loop Step Nodes Flow -->
      <div class="grid grid-cols-1 md:grid-cols-6 gap-3 pt-2">
        <!-- Node 1: Origin / Finding -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">1. Origin</div>
            <div v-if="chain.sourceFindingCode" class="font-mono text-xs text-rose-400 font-semibold truncate">
              {{ chain.sourceFindingCode }}
            </div>
            <div v-else class="text-xs text-slate-300 font-medium">Direct Engineering</div>
          </div>
          <div class="mt-2 text-[10px] text-slate-500 truncate">
            {{ chain.sourceFindingCode ? 'Defect Finding' : 'Non-test origin' }}
          </div>
        </div>

        <!-- Node 2: Source Revision -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">2. Source Rev</div>
            <div class="font-mono text-xs text-amber-300 font-bold">
              {{ chain.sourceRevisionCode }}
            </div>
          </div>
          <div class="mt-2 text-[10px] text-slate-500">
            Source Baseline
          </div>
        </div>

        <!-- Node 3: CCB Approvals -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">3. CCB Review</div>
            <div class="text-xs font-semibold" :class="chain.ccbApprovalCount >= 3 ? 'text-emerald-400' : 'text-amber-400'">
              {{ chain.ccbApprovalCount || 0 }}/3 Approved
            </div>
          </div>
          <div class="mt-2 text-[10px] text-slate-500">
            Multi-tier Sign-off
          </div>
        </div>

        <!-- Node 4: ECO -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">4. Change Order</div>
            <div v-if="chain.ecoCode" class="font-mono text-xs text-cyan-300 font-bold truncate">
              {{ chain.ecoCode }}
            </div>
            <div v-else class="text-xs text-slate-500 italic">Pending ECO</div>
          </div>
          <div class="mt-2 text-[10px]" :class="chain.ecoStatus === 'IMPLEMENTED' ? 'text-emerald-400' : 'text-slate-500'">
            {{ chain.ecoStatus || 'Not Created' }}
          </div>
        </div>

        <!-- Node 5: Target Revision & BOM -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">5. Target Rev</div>
            <div v-if="chain.targetRevisionCode" class="font-mono text-xs text-emerald-400 font-bold">
              {{ chain.targetRevisionCode }}
            </div>
            <div v-else class="text-xs text-slate-500 italic">Not Spawned</div>
          </div>
          <div class="mt-2 text-[10px] text-slate-500 truncate">
            {{ chain.targetBOMCode || 'New BOM Pending' }}
          </div>
        </div>

        <!-- Node 6: Verification & Regression -->
        <div class="p-3 rounded-xl bg-slate-950/80 border border-slate-800 flex flex-col justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-slate-500 tracking-wider mb-1">6. Verification</div>
            <div v-if="chain.regressionBuildRequired" class="text-xs font-semibold text-purple-400">
              Mandatory
            </div>
            <div v-else class="text-xs text-slate-500">Not required</div>
          </div>
          <div class="mt-2 text-[10px]" :class="chain.closedLoopVerified ? 'text-emerald-400 font-bold' : 'text-amber-400'">
            {{ chain.closedLoopVerified ? 'CLOSED LOOP OK' : 'Verification Open' }}
          </div>
        </div>
      </div>
    </div>

    <div v-if="(!chains || chains.length === 0)" class="text-center py-12 rounded-2xl border border-slate-800 bg-slate-900/40 text-slate-500 text-sm">
      No engineering change traceability chains available.
    </div>
  </div>
</template>

<script>
export default {
  name: 'TraceabilityGraph',
  props: {
    chains: {
      type: Array,
      default: () => []
    }
  }
};
</script>
