<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <span class="p-2 rounded-xl bg-cyan-500/10 text-cyan-400">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2" />
            </svg>
          </span>
          Engineering Change Orders (ECO)
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Authorized change execution baselines: automated BOM cloning, revision incrementation, and regression mandate.
        </p>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="p-4 rounded-2xl bg-slate-900/60 border border-slate-800 flex flex-wrap items-center gap-3">
      <div class="flex-1 min-w-[240px]">
        <input
          v-model="filters.search"
          @input="applyFilters"
          type="text"
          placeholder="Search by ECO code, title, or owner..."
          class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 placeholder-slate-500 text-xs focus:border-cyan-500 focus:outline-none"
        />
      </div>

      <select
        v-model="filters.status"
        @change="applyFilters"
        class="px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
      >
        <option value="">All Statuses</option>
        <option value="DRAFT">Draft</option>
        <option value="APPROVED">Approved (Ready to Implement)</option>
        <option value="IMPLEMENTED">Implemented (Baseline Generated)</option>
        <option value="VERIFIED">Verified (Regression OK)</option>
        <option value="CLOSED">Closed</option>
      </select>
    </div>

    <!-- ECO Table -->
    <div class="overflow-x-auto rounded-2xl border border-slate-800 bg-slate-900/60 shadow-xl">
      <table class="w-full text-left text-xs">
        <thead class="bg-slate-950 text-slate-400 uppercase tracking-wider font-semibold border-b border-slate-800">
          <tr>
            <th class="py-3.5 px-4">ECO Code</th>
            <th class="py-3.5 px-4">Title & Description</th>
            <th class="py-3.5 px-3">Status</th>
            <th class="py-3.5 px-3">Source ECR</th>
            <th class="py-3.5 px-3">Source Rev</th>
            <th class="py-3.5 px-3">Target Rev Baseline</th>
            <th class="py-3.5 px-3">Owner</th>
            <th class="py-3.5 px-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60">
          <tr
            v-for="eco in ecos"
            :key="eco.id"
            class="hover:bg-slate-800/30 transition-colors cursor-pointer"
            @click="$router.push(`/eco/${eco.id}`)"
          >
            <td class="py-3.5 px-4 font-mono font-bold text-cyan-300">
              {{ eco.code }}
            </td>
            <td class="py-3.5 px-4">
              <div class="font-semibold text-slate-200 max-w-sm truncate">{{ eco.title }}</div>
              <div class="text-[11px] text-slate-400 truncate max-w-sm mt-0.5">{{ eco.implementationNotes || '—' }}</div>
            </td>
            <td class="py-3.5 px-3">
              <ChangeStatusBadge :status="eco.status" />
            </td>
            <td class="py-3.5 px-3 font-mono text-[11px] text-slate-300">
              <span v-if="eco.ecr" class="text-blue-300 bg-blue-950/40 px-1.5 py-0.5 rounded border border-blue-800/40">
                {{ eco.ecr.code }}
              </span>
              <span v-else>—</span>
            </td>
            <td class="py-3.5 px-3 font-mono text-slate-300">
              {{ eco.sourceHardwareRevision ? eco.sourceHardwareRevision.code : '—' }}
            </td>
            <td class="py-3.5 px-3 font-mono font-bold text-emerald-400">
              {{ eco.targetHardwareRevision ? eco.targetHardwareRevision.code : 'Pending Baseline' }}
            </td>
            <td class="py-3.5 px-3 text-slate-300">
              {{ eco.implementationOwner || '—' }}
            </td>
            <td class="py-3.5 px-4 text-right" @click.stop>
              <router-link
                :to="`/eco/${eco.id}`"
                class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors"
              >
                Inspect
              </router-link>
            </td>
          </tr>

          <tr v-if="ecos.length === 0">
            <td colspan="8" class="text-center py-12 text-slate-500">
              No engineering change orders found matching current criteria.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { useChangeStore } from '@/stores/changeStore';
import ChangeStatusBadge from '@/components/changes/ChangeStatusBadge.vue';

export default {
  name: 'ECORegistryView',
  components: {
    ChangeStatusBadge
  },
  data() {
    return {
      filters: {
        search: '',
        status: ''
      }
    };
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    ecos() {
      return this.changeStore.ecos;
    }
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
      await this.changeStore.fetchECOs(this.filters);
    },
    applyFilters() {
      this.fetchData();
    }
  }
};
</script>
