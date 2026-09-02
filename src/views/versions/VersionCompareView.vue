<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <PageHeader
      title="Version Comparison"
      subtitle="Side-by-side engineering version and hardware revision diff inspection"
    >
      <template #actions>
        <router-link to="/versions" class="btn btn-secondary">
          <AppIcon name="chevron-left" size="xs" class="mr-1.5 text-slate-500" />
          <span>Back to Versions</span>
        </router-link>
      </template>
    </PageHeader>

    <!-- Version Selectors Card -->
    <div class="app-card p-5 bg-white space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 items-center">
        <!-- Version A (Base) -->
        <div class="p-4 rounded-xl border border-slate-200/90 bg-slate-50/50 space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-3xs font-mono font-bold uppercase tracking-wider text-slate-500">Base Version (A)</span>
            <span class="text-3xs font-mono px-2 py-0.5 rounded bg-slate-200 text-slate-700 font-bold">PREDECESSOR</span>
          </div>
          <select v-model="selectedVersionAId" class="form-select font-mono text-xs">
            <option
              v-for="v in versionStore.versions"
              :key="v.id"
              :value="v.id"
            >
              {{ v.productCode }} — v{{ v.versionNumber }} ({{ v.versionName }})
            </option>
          </select>
        </div>

        <!-- Version B (Target) -->
        <div class="p-4 rounded-xl border border-brand-200/90 bg-brand-50/30 space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-3xs font-mono font-bold uppercase tracking-wider text-brand-700">Target Version (B)</span>
            <span class="text-3xs font-mono px-2 py-0.5 rounded bg-brand-600 text-white font-bold">COMPARISON TARGET</span>
          </div>
          <select v-model="selectedVersionBId" class="form-select font-mono text-xs">
            <option
              v-for="v in versionStore.versions"
              :key="v.id"
              :value="v.id"
            >
              {{ v.productCode }} — v{{ v.versionNumber }} ({{ v.versionName }})
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Comparison Matrix Card -->
    <div v-if="versionA && versionB" class="app-card overflow-hidden">
      <div class="px-5 py-3.5 bg-slate-50/80 border-b border-slate-200 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <AppIcon name="git-branch" size="xs" class="text-brand-600" />
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">
            Version Delta Comparison Matrix
          </h3>
        </div>
        <div class="flex items-center gap-3 text-3xs font-mono">
          <span class="flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-amber-500"></span>
            <span class="text-slate-600">Modified Parameter</span>
          </span>
          <span class="flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-slate-300"></span>
            <span class="text-slate-500">Identical</span>
          </span>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100 text-left text-xs">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-4 py-3 font-bold text-slate-500 uppercase tracking-wider text-3xs w-1/4">Property</th>
              <th class="px-4 py-3 font-bold text-slate-900 uppercase tracking-wider text-3xs w-3/8 bg-slate-100/50">
                Base: v{{ versionA.versionNumber }}
              </th>
              <th class="px-4 py-3 font-bold text-brand-700 uppercase tracking-wider text-3xs w-3/8 bg-brand-50/30">
                Target: v{{ versionB.versionNumber }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 font-sans">
            <!-- Product Code -->
            <tr :class="isDifferent(versionA.productCode, versionB.productCode) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Product Code</td>
              <td class="px-4 py-3 font-mono font-bold text-slate-800">{{ versionA.productCode }}</td>
              <td class="px-4 py-3 font-mono font-bold text-slate-800">{{ versionB.productCode }}</td>
            </tr>

            <!-- Version Number -->
            <tr :class="isDifferent(versionA.versionNumber, versionB.versionNumber) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Version Number</td>
              <td class="px-4 py-3 font-mono font-bold text-slate-900">v{{ versionA.versionNumber }}</td>
              <td class="px-4 py-3 font-mono font-bold text-brand-700">v{{ versionB.versionNumber }}</td>
            </tr>

            <!-- Version Name -->
            <tr :class="isDifferent(versionA.versionName, versionB.versionName) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Release Title</td>
              <td class="px-4 py-3 text-slate-800 font-medium">{{ versionA.versionName }}</td>
              <td class="px-4 py-3 text-slate-800 font-medium">{{ versionB.versionName }}</td>
            </tr>

            <!-- Lifecycle Status -->
            <tr :class="isDifferent(versionA.status, versionB.status) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Lifecycle Status</td>
              <td class="px-4 py-3">
                <VersionBadge :status="versionA.status" />
              </td>
              <td class="px-4 py-3">
                <VersionBadge :status="versionB.status" />
              </td>
            </tr>

            <!-- Current Revision -->
            <tr :class="isDifferent(versionA.currentRevision, versionB.currentRevision) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Current Revision</td>
              <td class="px-4 py-3 font-mono font-bold text-slate-800">{{ versionA.currentRevision }}</td>
              <td class="px-4 py-3 font-mono font-bold text-slate-800">{{ versionB.currentRevision }}</td>
            </tr>

            <!-- Lead Engineer -->
            <tr :class="isDifferent(versionA.owner, versionB.owner) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Engineering Owner</td>
              <td class="px-4 py-3 text-slate-800">{{ versionA.owner }}</td>
              <td class="px-4 py-3 text-slate-800">{{ versionB.owner }}</td>
            </tr>

            <!-- Architecture Scope -->
            <tr :class="isDifferent(versionA.description, versionB.description) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Description / Scope</td>
              <td class="px-4 py-3 text-slate-600 leading-relaxed text-xs">{{ versionA.description }}</td>
              <td class="px-4 py-3 text-slate-600 leading-relaxed text-xs">{{ versionB.description }}</td>
            </tr>

            <!-- Release Notes -->
            <tr :class="isDifferent(versionA.releaseNotes, versionB.releaseNotes) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Release Notes</td>
              <td class="px-4 py-3 text-slate-700 font-mono text-3xs whitespace-pre-line bg-slate-50/50 rounded">{{ versionA.releaseNotes }}</td>
              <td class="px-4 py-3 text-slate-700 font-mono text-3xs whitespace-pre-line bg-brand-50/20 rounded">{{ versionB.releaseNotes }}</td>
            </tr>

            <!-- Artifacts Count -->
            <tr :class="isDifferent(versionA.documents ? versionA.documents.length : 0, versionB.documents ? versionB.documents.length : 0) ? 'bg-amber-50/30' : ''">
              <td class="px-4 py-3 font-semibold text-slate-500">Attached Documents</td>
              <td class="px-4 py-3 font-mono text-xs">{{ versionA.documents ? versionA.documents.length : 0 }} files</td>
              <td class="px-4 py-3 font-mono text-xs">{{ versionB.documents ? versionB.documents.length : 0 }} files</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import PageHeader from '@/components/common/PageHeader.vue';
import VersionBadge from '@/components/versions/VersionBadge.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'VersionCompareView',
  components: {
    PageHeader,
    VersionBadge,
    AppIcon
  },
  data() {
    return {
      selectedVersionAId: '',
      selectedVersionBId: ''
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    versionA() {
      return this.versionStore.getVersionById(this.selectedVersionAId);
    },
    versionB() {
      return this.versionStore.getVersionById(this.selectedVersionBId);
    }
  },
  async mounted() {
    await this.versionStore.fetchVersions();
    const vA = this.$route.query.vA;
    const vB = this.$route.query.vB;

    if (vA && this.versionStore.getVersionById(vA)) {
      this.selectedVersionAId = vA;
    } else if (this.versionStore.versions.length > 1) {
      this.selectedVersionAId = this.versionStore.versions[1].id;
    } else if (this.versionStore.versions.length > 0) {
      this.selectedVersionAId = this.versionStore.versions[0].id;
    }

    if (vB && this.versionStore.getVersionById(vB)) {
      this.selectedVersionBId = vB;
    } else if (this.versionStore.versions.length > 0) {
      this.selectedVersionBId = this.versionStore.versions[0].id;
    }
  },
  methods: {
    isDifferent(valA, valB) {
      return String(valA || '').trim() !== String(valB || '').trim();
    }
  }
};
</script>
