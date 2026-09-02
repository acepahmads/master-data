<template>
  <div class="space-y-4">
    <!-- Diff KPI Summary Banner -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3 p-4 rounded-xl bg-slate-900/80 border border-slate-800">
      <div>
        <div class="text-[11px] text-slate-400 uppercase tracking-wider">Source BOM Cost</div>
        <div class="text-base font-bold text-slate-200 font-mono mt-0.5">
          ${{ Number(summary?.sourceBomCost || 0).toFixed(4) }}
        </div>
        <div class="text-[10px] text-slate-500">{{ summary?.sourceRevisionCode || 'Source' }}</div>
      </div>
      <div>
        <div class="text-[11px] text-slate-400 uppercase tracking-wider">Target BOM Cost</div>
        <div class="text-base font-bold text-cyan-300 font-mono mt-0.5">
          ${{ Number(summary?.targetBomCost || 0).toFixed(4) }}
        </div>
        <div class="text-[10px] text-slate-500">{{ summary?.targetRevisionCode || 'Target' }}</div>
      </div>
      <div>
        <div class="text-[11px] text-slate-400 uppercase tracking-wider">Net Unit Delta</div>
        <div
          :class="[
            'text-base font-bold font-mono mt-0.5',
            summary?.netCostDelta > 0 ? 'text-amber-400' : 'text-emerald-400'
          ]"
        >
          {{ formatDelta(summary?.netCostDelta) }}
        </div>
        <div class="text-[10px] text-slate-500">
          {{ summary?.percentCostDelta ? summary.percentCostDelta.toFixed(2) + '%' : '0.00%' }}
        </div>
      </div>
      <div>
        <div class="text-[11px] text-slate-400 uppercase tracking-wider">Delta Counts</div>
        <div class="text-xs font-semibold text-slate-300 mt-1 flex items-center gap-2">
          <span class="px-1.5 py-0.5 rounded bg-amber-950/60 text-amber-300 border border-amber-800/40">
            {{ summary?.totalModified || 0 }} Modified
          </span>
          <span class="px-1.5 py-0.5 rounded bg-slate-800 text-slate-400">
            {{ summary?.totalUnchanged || 0 }} Unchanged
          </span>
        </div>
      </div>
    </div>

    <!-- Line-Item Diff Table -->
    <div class="overflow-x-auto rounded-xl border border-slate-800 bg-slate-900/40">
      <table class="w-full text-left text-xs">
        <thead class="bg-slate-950 text-slate-400 uppercase tracking-wider font-semibold border-b border-slate-800">
          <tr>
            <th class="py-3 px-4">Ref Des</th>
            <th class="py-3 px-4">Component Name & MPN</th>
            <th class="py-3 px-3 text-center">Status</th>
            <th class="py-3 px-3 text-right">Qty</th>
            <th class="py-3 px-3 text-right">Source Cost</th>
            <th class="py-3 px-3 text-right">Target Cost</th>
            <th class="py-3 px-3 text-right">Delta</th>
            <th class="py-3 px-4">Change Reason / Delta Note</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 font-mono">
          <tr
            v-for="(item, idx) in summary?.items || []"
            :key="idx"
            :class="[
              'hover:bg-slate-800/40 transition-colors',
              getRowHighlightClass(item.changeStatus)
            ]"
          >
            <td class="py-3 px-4 font-bold text-slate-200">
              {{ item.referenceDesignator }}
            </td>
            <td class="py-3 px-4 font-sans font-medium text-slate-300">
              <div>{{ item.componentName }}</div>
              <div v-if="item.mpn" class="text-[11px] text-slate-500 font-mono">{{ item.mpn }} ({{ item.package }})</div>
            </td>
            <td class="py-3 px-3 text-center">
              <span
                :class="[
                  'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                  getStatusBadgeClass(item.changeStatus)
                ]"
              >
                {{ item.changeStatus }}
              </span>
            </td>
            <td class="py-3 px-3 text-right text-slate-300">
              {{ item.targetQuantity }}
            </td>
            <td class="py-3 px-3 text-right text-slate-400">
              ${{ Number(item.sourceTotalCost || 0).toFixed(4) }}
            </td>
            <td class="py-3 px-3 text-right text-cyan-300 font-semibold">
              ${{ Number(item.targetTotalCost || 0).toFixed(4) }}
            </td>
            <td
              :class="[
                'py-3 px-3 text-right font-bold',
                item.costDelta > 0 ? 'text-amber-400' : item.costDelta < 0 ? 'text-emerald-400' : 'text-slate-500'
              ]"
            >
              {{ formatDelta(item.costDelta) }}
            </td>
            <td class="py-3 px-4 font-sans text-xs text-slate-400 truncate max-w-xs">
              {{ item.changeReason || '—' }}
            </td>
          </tr>
          <tr v-if="(!summary?.items || summary.items.length === 0)">
            <td colspan="8" class="text-center py-8 text-slate-500 font-sans">
              No BOM items available for comparison preview.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BOMChangeDiffTable',
  props: {
    summary: {
      type: Object,
      default: () => ({})
    }
  },
  methods: {
    formatDelta(val) {
      if (val === undefined || val === null || val === 0) return '$0.0000';
      const prefix = val > 0 ? '+$' : '-$';
      return `${prefix}${Math.abs(Number(val)).toFixed(4)}`;
    },
    getRowHighlightClass(status) {
      switch ((status || '').toUpperCase()) {
        case 'MODIFIED':
          return 'bg-amber-950/10';
        case 'ADDED':
          return 'bg-emerald-950/15';
        case 'REMOVED':
          return 'bg-rose-950/15';
        default:
          return '';
      }
    },
    getStatusBadgeClass(status) {
      switch ((status || '').toUpperCase()) {
        case 'MODIFIED':
          return 'bg-amber-950/80 text-amber-300 border border-amber-800';
        case 'ADDED':
          return 'bg-emerald-950/80 text-emerald-300 border border-emerald-800';
        case 'REMOVED':
          return 'bg-rose-950/80 text-rose-300 border border-rose-800';
        case 'UNCHANGED':
          return 'bg-slate-800 text-slate-400 border border-slate-700';
        default:
          return 'bg-slate-800 text-slate-400';
      }
    }
  }
};
</script>
