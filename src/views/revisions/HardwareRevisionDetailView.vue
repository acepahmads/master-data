<template>
  <div v-if="revision" class="space-y-5">
    <!-- REVISION HERO -->
    <div class="app-card p-5 md:p-6 bg-white space-y-4">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-4 border-b border-slate-100">
        <div class="flex items-start gap-4 min-w-0">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-tr from-slate-900 to-slate-700 text-white flex items-center justify-center font-mono font-bold text-lg shadow-md shrink-0">
            {{ revision.code.replace('REV-', '') }}
          </div>
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono font-bold text-sm text-slate-900 bg-slate-100 px-2.5 py-0.5 rounded border border-slate-200">
                {{ revision.code }}
              </span>
              <RevisionBadge :status="revision.status" />
              <router-link
                :to="`/versions/${revision.versionId}`"
                class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200 hover:underline"
              >
                v{{ revision.versionNumber }}
              </router-link>
              <span class="text-3xs font-mono text-slate-500 font-semibold">{{ revision.productCode }}</span>
            </div>
            <h1 class="text-xl md:text-2xl font-bold text-slate-900 tracking-tight truncate">
              {{ revision.productName }} — {{ revision.name }}
            </h1>
            <p class="text-xs text-slate-500 mt-1 font-medium max-w-3xl">
              {{ revision.changeSummary }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <router-link
            :to="`/versions/${revision.versionId}`"
            class="btn btn-secondary"
          >
            <AppIcon name="chevron-left" size="xs" class="mr-1.5 text-slate-500" />
            <span>Back to Version</span>
          </router-link>
          <router-link
            :to="`/revisions/${revision.id}/bom`"
            class="btn btn-primary bg-gradient-to-r from-brand-600 to-indigo-600 border-none shadow-sm"
          >
            <AppIcon name="layers" size="xs" class="mr-1.5 text-white" />
            <span>Engineering BOM</span>
          </router-link>
          <router-link
            :to="`/prototypes/create?revision=${revision.id}`"
            class="btn btn-secondary text-xs flex items-center gap-1.5 text-indigo-700 bg-indigo-50 border-indigo-200 hover:bg-indigo-100 font-bold"
          >
            <AppIcon name="cpu" size="xs" />
            <span>+ Prototype Build</span>
          </router-link>
          <button
            @click="promoteStatus"
            class="btn btn-secondary"
          >
            <AppIcon name="check" size="xs" class="mr-1.5 text-slate-700" />
            <span>Promote Status</span>
          </button>
        </div>
      </div>

      <!-- Key Metadata Grid -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs font-mono pt-1">
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Hardware Engineer</span>
          <span class="font-bold text-slate-800 text-xs mt-1 block font-sans">{{ revision.engineer || 'Marcus Vance' }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">PCB Designation</span>
          <span class="font-bold text-slate-900 text-xs mt-1 block">{{ revision.pcbRevision || 'PCB-v2.1-B' }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Schematic Revision</span>
          <span class="text-slate-700 text-xs mt-1 block">{{ revision.schematicRevision || 'SCH-v2.1-02' }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Release Date</span>
          <span class="text-brand-600 text-xs mt-1 block">{{ revision.releaseDate || 'In Progress' }}</span>
        </div>
      </div>
    </div>

    <!-- TABS NAVIGATION & CONTENT -->
    <div class="app-card overflow-hidden">
      <div class="px-5 border-b border-slate-200 bg-slate-50/70 flex items-center gap-6 text-xs">
        <button
          @click="activeTab = 'bom'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'bom'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="layers" size="xs" />
          <span>Engineering BOM & Cost</span>
        </button>

        <button
          @click="activeTab = 'changes'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'changes'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="file-text" size="xs" />
          <span>Change Record</span>
        </button>

        <button
          @click="activeTab = 'documents'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'documents'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="file" size="xs" />
          <span>Gerber & Schematic Files</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-slate-200 text-slate-700">
            {{ revision.documents ? revision.documents.length : 0 }}
          </span>
        </button>

        <button
          @click="activeTab = 'history'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'history'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="activity" size="xs" />
          <span>Revision Audit</span>
        </button>

        <button
          @click="activeTab = 'prototypes'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'prototypes'
              ? 'border-indigo-600 text-indigo-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="cpu" size="xs" />
          <span>Prototype Builds (Phase 4)</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200 font-bold">
            {{ revisionPrototypes.length }}
          </span>
        </button>
      </div>

      <div class="p-6">
        <!-- TAB 0: ENGINEERING BOM & COST PREVIEW -->
        <div v-if="activeTab === 'bom'" class="space-y-6">
          <div class="p-5 rounded-2xl bg-gradient-to-br from-slate-900 to-slate-800 text-white flex flex-col md:flex-row items-start md:items-center justify-between gap-4 shadow-lg">
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="px-2 py-0.5 rounded bg-brand-500/20 text-brand-300 font-mono text-3xs font-bold border border-brand-500/30">
                  HARDWARE REVISION BOM
                </span>
                <span class="font-mono text-xs text-slate-300">{{ revision.code }}</span>
              </div>
              <h3 class="text-lg font-bold text-white tracking-tight">
                Multi-Level Engineering Bill of Materials & Cost Rollup
              </h3>
              <p class="text-xs text-slate-300 max-w-2xl leading-relaxed">
                Manages electronic subsystems, PCB components, reference designators (e.g. U1, C1-C4), and recursive line cost rollups for this hardware revision baseline.
              </p>
            </div>

            <router-link
              :to="`/revisions/${revision.id}/bom`"
              class="btn btn-primary px-5 py-2.5 bg-brand-500 hover:bg-brand-400 text-white font-bold text-xs shrink-0 shadow-md flex items-center gap-2"
            >
              <AppIcon name="external-link" size="xs" />
              <span>Launch BOM Workspace</span>
            </router-link>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="p-4 rounded-xl bg-slate-50 border border-slate-200/80">
              <span class="text-3xs uppercase font-bold text-slate-400">Ownership Hierarchy</span>
              <div class="text-xs font-mono font-semibold text-slate-800 mt-1">
                {{ revision.productCode }} ➔ v{{ revision.versionNumber }} ➔ {{ revision.code }}
              </div>
              <span class="text-3xs text-slate-400 mt-1 block">Exclusive Revision Owner</span>
            </div>

            <div class="p-4 rounded-xl bg-slate-50 border border-slate-200/80">
              <span class="text-3xs uppercase font-bold text-slate-400">Multi-Level Architecture</span>
              <div class="text-xs font-semibold text-indigo-700 mt-1">
                Sub-Assemblies + Components
              </div>
              <span class="text-3xs text-slate-400 mt-1 block">Recursive Hierarchy Supported</span>
            </div>

            <div class="p-4 rounded-xl bg-emerald-50/60 border border-emerald-200/80">
              <span class="text-3xs uppercase font-bold text-emerald-700">Cost Estimation Engine</span>
              <div class="text-xs font-mono font-bold text-emerald-900 mt-1">
                Quantity × Unit Cost Rollup
              </div>
              <span class="text-3xs text-emerald-600 mt-1 block">Component Master Cost Foundation</span>
            </div>
          </div>
        </div>

        <!-- TAB 1: CHANGE RECORD -->
        <div v-else-if="activeTab === 'changes'" class="space-y-6">
          <div>
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 mb-2">
              Engineering Change Summary
            </h4>
            <p class="text-xs text-slate-600 leading-relaxed font-medium">
              {{ revision.changeSummary }}
            </p>
          </div>

          <div class="space-y-2">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
              Hardware Modification Checklist
            </h4>
            <div class="space-y-2">
              <div
                v-for="(item, idx) in revision.changesList || []"
                :key="idx"
                class="p-3 rounded-lg border border-slate-200/80 bg-slate-50/50 flex items-start gap-2.5 text-xs text-slate-700 font-medium"
              >
                <span class="w-5 h-5 rounded bg-brand-50 text-brand-600 font-mono font-bold text-3xs flex items-center justify-center shrink-0 mt-0.5 border border-brand-200">
                  {{ idx + 1 }}
                </span>
                <span class="leading-relaxed">{{ item }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 2: DOCUMENTS -->
        <div v-else-if="activeTab === 'documents'" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
            <div
              v-for="(doc, idx) in revision.documents || []"
              :key="idx"
              class="p-4 rounded-xl border border-slate-200/90 bg-slate-50/50 hover:bg-white hover:border-brand-300 hover:shadow-2xs transition-all flex items-start gap-3 group"
            >
              <div class="w-10 h-10 rounded-lg bg-brand-50 text-brand-600 border border-brand-200/80 flex items-center justify-center shrink-0 shadow-2xs">
                <AppIcon name="file-text" size="sm" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="font-mono font-bold text-xs text-slate-900 truncate group-hover:text-brand-600 transition-colors">
                  {{ doc.name || doc }}
                </div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5 flex items-center gap-2">
                  <span>{{ doc.type || 'Artifact' }}</span>
                  <span>•</span>
                  <span>{{ doc.size || '2.4 MB' }}</span>
                </div>
                <button
                  @click="downloadDoc(doc.name || doc)"
                  class="mt-2 text-3xs font-semibold text-brand-600 hover:text-brand-800 flex items-center gap-1"
                >
                  <AppIcon name="download" size="xs" />
                  <span>Download File</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 3: AUDIT HISTORY -->
        <div v-else-if="activeTab === 'history'" class="space-y-4">
          <ChangeHistoryTimeline :items="revisionHistory" />
        </div>

        <!-- TAB 4: PROTOTYPE BUILDS (PHASE 4) -->
        <div v-else-if="activeTab === 'prototypes'" class="space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div>
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800">Physical Prototype Runs for {{ revision.code }}</h3>
              <p class="text-3xs text-slate-400">Physical builds executed against this hardware baseline with frozen BOM snapshots.</p>
            </div>
            <router-link
              :to="`/prototypes/create?revision=${revision.id}`"
              class="btn btn-primary btn-sm text-xs flex items-center gap-1.5 bg-indigo-600 hover:bg-indigo-700"
            >
              <AppIcon name="plus" size="xs" />
              <span>New Prototype Build</span>
            </router-link>
          </div>

          <div v-if="revisionPrototypes.length > 0" class="space-y-3">
            <div
              v-for="p in revisionPrototypes"
              :key="p.id"
              class="p-4 rounded-xl border border-slate-200 bg-white hover:border-indigo-300 hover:shadow-card transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-3"
            >
              <div class="min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="font-mono text-xs font-bold text-indigo-700 bg-indigo-50 px-2 py-0.5 rounded border border-indigo-200">
                    {{ p.code }}
                  </span>
                  <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 border border-slate-200">
                    Build #{{ p.buildNumber }} ({{ p.quantity }} Units)
                  </span>
                  <span class="text-3xs font-medium px-2 py-0.5 rounded-full bg-purple-50 text-purple-700 border border-purple-200">
                    {{ p.purpose }}
                  </span>
                </div>
                <div class="font-bold text-sm text-slate-900">{{ p.name }}</div>
                <div class="text-3xs text-slate-500 mt-0.5">
                  Owner: <strong class="text-slate-700">{{ p.engineeringOwner }}</strong> | Target Date: <strong class="font-mono text-slate-700">{{ p.targetBuildDate }}</strong>
                </div>
              </div>

              <div class="flex items-center gap-3 shrink-0">
                <div class="text-right text-3xs font-mono">
                  <span class="text-slate-400 block">Status</span>
                  <span class="font-bold text-indigo-700">{{ p.status }}</span>
                </div>
                <router-link
                  :to="`/prototypes/${p.id}`"
                  class="btn btn-secondary py-1 px-3 text-xs font-bold flex items-center gap-1"
                >
                  <span>Open Build</span>
                  <AppIcon name="chevron-right" size="2xs" />
                </router-link>
              </div>
            </div>
          </div>

          <div v-else class="text-center py-8 bg-slate-50 rounded-xl border border-slate-200 space-y-2">
            <AppIcon name="cpu" size="lg" class="mx-auto text-slate-300" />
            <div class="text-xs font-bold text-slate-700">No prototype builds created yet for {{ revision.code }}</div>
            <p class="text-3xs text-slate-400">Click below to freeze a BOM snapshot and instantiate the first prototype run.</p>
            <router-link
              :to="`/prototypes/create?revision=${revision.id}`"
              class="btn btn-primary btn-sm text-xs inline-flex items-center gap-1 mt-2 bg-indigo-600 hover:bg-indigo-700"
            >
              <AppIcon name="plus" size="xs" />
              <span>Create First Prototype Build</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import { usePrototypeStore } from '@/stores/prototypeStore';
import RevisionBadge from '@/components/versions/RevisionBadge.vue';
import ChangeHistoryTimeline from '@/components/versions/ChangeHistoryTimeline.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'HardwareRevisionDetailView',
  components: {
    RevisionBadge,
    ChangeHistoryTimeline,
    AppIcon
  },
  data() {
    return {
      activeTab: 'bom'
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    protoStore() {
      return usePrototypeStore();
    },
    revisionId() {
      return this.$route.params.id;
    },
    revision() {
      return this.versionStore.getRevisionById(this.revisionId) || this.versionStore.revisions[0];
    },
    revisionPrototypes() {
      if (!this.revision) return [];
      return this.protoStore.getPrototypesByRevision(this.revision.id || this.revision.code);
    },
    revisionHistory() {
      if (!this.revision) return [];
      return this.versionStore.versionHistory.filter(h => h.revision === this.revision.code || h.targetId === this.revision.id);
    }
  },
  methods: {
    async promoteStatus() {
      const nextStatus = this.revision.status === 'Draft' ? 'In Review' :
        this.revision.status === 'In Review' ? 'Approved' :
        this.revision.status === 'Approved' ? 'Released' : 'Released';

      try {
        await this.versionStore.saveRevision({
          ...this.revision,
          status: nextStatus
        });
        alert(`Revision status promoted to: ${nextStatus}`);
      } catch (err) {
        alert('Failed to promote revision status: ' + (err.response?.data?.message || err.message));
      }
    },
    downloadDoc(name) {
      alert(`Downloading file: ${name}`);
    }
  },
  async mounted() {
    await Promise.all([
      this.versionStore.fetchVersions(),
      this.versionStore.fetchRevisions()
    ]);
  }
};
</script>
