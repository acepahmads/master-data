<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-100 flex items-center gap-3">
          <span class="p-2 rounded-xl bg-blue-500/10 text-blue-400">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </span>
          Engineering Change Requests (ECR)
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Formal requests proposing modifications to hardware architecture, components, or BOM items.
        </p>
      </div>

      <router-link
        to="/ecr/create"
        class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-bold text-sm shadow-lg shadow-cyan-500/20 transition-all self-start md:self-auto"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New ECR Wizard
      </router-link>
    </div>

    <!-- Filter Bar -->
    <div class="p-4 rounded-2xl bg-slate-900/60 border border-slate-800 flex flex-wrap items-center gap-3">
      <div class="flex-1 min-w-[240px]">
        <input
          v-model="filters.search"
          @input="applyFilters"
          type="text"
          placeholder="Search by code, title, or keyword..."
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
        <option value="UNDER_REVIEW">Under Review</option>
        <option value="PENDING_APPROVAL">Pending Approval</option>
        <option value="APPROVED">Approved</option>
        <option value="REJECTED">Rejected</option>
        <option value="CONVERTED_TO_ECO">Converted to ECO</option>
      </select>

      <select
        v-model="filters.priority"
        @change="applyFilters"
        class="px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
      >
        <option value="">All Priorities</option>
        <option value="CRITICAL">Critical</option>
        <option value="HIGH">High</option>
        <option value="MEDIUM">Medium</option>
        <option value="LOW">Low</option>
      </select>

      <select
        v-model="filters.origin"
        @change="applyFilters"
        class="px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:border-cyan-500 focus:outline-none"
      >
        <option value="">All Origins</option>
        <option value="TEST_FINDING">Test Finding</option>
        <option value="CUSTOMER_REQUIREMENT">Customer Requirement</option>
        <option value="SUPPLIER_CHANGE">Supplier Change</option>
        <option value="OBSOLESCENCE">Obsolescence</option>
        <option value="COST_REDUCTION">Cost Reduction</option>
        <option value="RELIABILITY">Reliability</option>
        <option value="REGULATORY">Regulatory</option>
      </select>
    </div>

    <!-- ECR Table -->
    <div class="overflow-x-auto rounded-2xl border border-slate-800 bg-slate-900/60 shadow-xl">
      <table class="w-full text-left text-xs">
        <thead class="bg-slate-950 text-slate-400 uppercase tracking-wider font-semibold border-b border-slate-800">
          <tr>
            <th class="py-3.5 px-4">ECR Code</th>
            <th class="py-3.5 px-4">Title & Scope</th>
            <th class="py-3.5 px-3">Status</th>
            <th class="py-3.5 px-3">Priority</th>
            <th class="py-3.5 px-3">Origin</th>
            <th class="py-3.5 px-3">Source Revision</th>
            <th class="py-3.5 px-3">Lead</th>
            <th class="py-3.5 px-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60">
          <tr
            v-for="ecr in ecrs"
            :key="ecr.id"
            class="hover:bg-slate-800/30 transition-colors cursor-pointer"
            @click="$router.push(`/ecr/${ecr.id}`)"
          >
            <td class="py-3.5 px-4 font-mono font-bold text-slate-200">
              {{ ecr.code }}
            </td>
            <td class="py-3.5 px-4">
              <div class="font-semibold text-slate-200 max-w-sm truncate">{{ ecr.title }}</div>
              <div class="text-[11px] text-slate-400 truncate max-w-sm mt-0.5">{{ ecr.businessReason || ecr.description }}</div>
            </td>
            <td class="py-3.5 px-3">
              <ChangeStatusBadge :status="ecr.status" />
            </td>
            <td class="py-3.5 px-3">
              <ChangePriorityBadge :priority="ecr.priority" />
            </td>
            <td class="py-3.5 px-3 font-mono text-[11px] text-slate-300">
              <span v-if="ecr.sourceFinding" class="text-rose-300 bg-rose-950/40 px-1.5 py-0.5 rounded border border-rose-800/40">
                {{ ecr.sourceFinding.code }}
              </span>
              <span v-else class="text-slate-400">
                {{ ecr.origin }}
              </span>
            </td>
            <td class="py-3.5 px-3 font-mono text-slate-300">
              {{ ecr.sourceHardwareRevision ? ecr.sourceHardwareRevision.code : '—' }}
            </td>
            <td class="py-3.5 px-3 text-slate-300">
              {{ ecr.assignedLead || ecr.requestedBy || '—' }}
            </td>
            <td class="py-3.5 px-4 text-right" @click.stop>
              <router-link
                :to="`/ecr/${ecr.id}`"
                class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors"
              >
                Inspect
              </router-link>
            </td>
          </tr>

          <tr v-if="ecrs.length === 0">
            <td colspan="8" class="text-center py-12 text-slate-500">
              No change requests found matching current criteria.
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
import ChangePriorityBadge from '@/components/changes/ChangePriorityBadge.vue';

export default {
  name: 'ECRRegistryView',
  components: {
    ChangeStatusBadge,
    ChangePriorityBadge
  },
  data() {
    return {
      filters: {
        search: '',
        status: '',
        priority: '',
        origin: ''
      }
    };
  },
  computed: {
    changeStore() {
      return useChangeStore();
    },
    ecrs() {
      return this.changeStore.ecrs;
    }
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
      await this.changeStore.fetchECRs(this.filters);
    },
    applyFilters() {
      this.fetchData();
    }
  }
};
</script>
