<template>
  <div v-if="version" class="space-y-5">
    <!-- 1. VERSION HERO -->
    <div class="app-card p-5 md:p-6 bg-white space-y-4">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-4 border-b border-slate-100">
        <div class="flex items-start gap-4 min-w-0">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-tr from-brand-700 to-brand-500 text-white flex items-center justify-center font-mono font-bold text-lg shadow-md shadow-brand-500/20 shrink-0">
            v{{ version.versionNumber.split('.')[0] }}
          </div>
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="font-mono font-bold text-sm text-slate-900 bg-slate-100 px-2.5 py-0.5 rounded border border-slate-200">
                v{{ version.versionNumber }}
              </span>
              <VersionBadge :status="version.status" />
              <span class="text-3xs font-mono font-bold px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200">
                {{ version.currentRevision }}
              </span>
              <router-link
                :to="`/products/${version.productCode}`"
                class="text-3xs font-mono text-slate-500 hover:text-brand-600 font-semibold"
              >
                {{ version.productCode }}
              </router-link>
            </div>
            <h1 class="text-xl md:text-2xl font-bold text-slate-900 tracking-tight truncate">
              {{ version.productName }}
            </h1>
            <p class="text-xs text-slate-500 mt-1 font-medium max-w-3xl">
              {{ version.versionName }} • {{ version.description }}
            </p>
          </div>
        </div>

        <!-- Action CTAs -->
        <div class="flex items-center gap-2 shrink-0">
          <router-link
            :to="`/versions/compare?vA=${version.id}`"
            class="btn btn-secondary"
          >
            <AppIcon name="git-branch" size="xs" class="mr-1.5 text-slate-500" />
            <span>Compare</span>
          </router-link>
          <button
            @click="openAddRevisionModal"
            class="btn btn-secondary"
          >
            <AppIcon name="plus" size="xs" class="mr-1.5 text-slate-500" />
            <span>Add Revision</span>
          </button>
          <router-link
            :to="`/versions/${version.id}/edit`"
            class="btn btn-primary"
          >
            <AppIcon name="edit" size="xs" class="mr-1.5 text-white" />
            <span>Edit Version</span>
          </router-link>
        </div>
      </div>

      <!-- Hero Key Metadata Strip -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs font-mono pt-1">
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Engineering Owner</span>
          <div class="flex items-center gap-1.5 mt-1">
            <img :src="version.ownerAvatar" class="w-5 h-5 rounded-full object-cover border border-slate-200" />
            <span class="font-bold text-slate-800 font-sans">{{ version.owner }}</span>
          </div>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Current HW Revision</span>
          <span class="font-bold text-slate-900 text-xs mt-1 block">{{ version.currentRevision }} ({{ revisions.length }} total)</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Release Date</span>
          <span class="text-slate-700 text-xs mt-1 block">{{ version.releaseDate || 'In Development' }}</span>
        </div>
        <div>
          <span class="text-3xs uppercase tracking-wider text-slate-400 font-sans font-semibold block">Base Predecessor</span>
          <span class="text-brand-600 text-xs mt-1 block">{{ version.baseVersionId ? `v${version.baseVersionId}` : 'Root Baseline' }}</span>
        </div>
      </div>
    </div>

    <!-- 2. VERSION LIFECYCLE TIMELINE -->
    <div class="app-card p-4 space-y-2">
      <div class="flex items-center justify-between px-2 pb-1 border-b border-slate-100">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-800 flex items-center gap-1.5">
          <AppIcon name="roadmap" size="xs" class="text-brand-600" />
          <span>Product Version Evolution Track</span>
        </h3>
        <span class="text-3xs font-mono text-slate-400">Click node to switch version view</span>
      </div>

      <VersionTimeline
        :versions="productVersions"
        :activeVersionId="version.id"
        @select-version="onSelectTimelineVersion"
      />
    </div>

    <!-- 3. DETAIL TABS & WORKSPACE -->
    <div class="app-card overflow-hidden">
      <!-- Tabs Navigation -->
      <div class="px-5 border-b border-slate-200 bg-slate-50/70 flex items-center gap-6 text-xs">
        <button
          @click="activeTab = 'overview'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'overview'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="file-text" size="xs" />
          <span>Overview</span>
        </button>

        <button
          @click="activeTab = 'revisions'"
          :class="[
            'py-3 font-semibold border-b-2 transition-all flex items-center gap-1.5',
            activeTab === 'revisions'
              ? 'border-brand-600 text-brand-600 font-bold'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          ]"
        >
          <AppIcon name="cpu" size="xs" />
          <span>Hardware Revisions</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-slate-200 text-slate-700">
            {{ revisions.length }}
          </span>
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
          <span>Engineering Documents</span>
          <span class="text-3xs font-mono px-1.5 py-0.2 rounded-full bg-slate-200 text-slate-700">
            {{ version.documents ? version.documents.length : 0 }}
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
          <span>Audit History</span>
        </button>
      </div>

      <!-- Tab Content Area -->
      <div class="p-6">
        <!-- TAB 1: OVERVIEW -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <!-- Release Notes Box -->
          <div class="space-y-2">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
              Engineering Release Notes
            </h4>
            <div class="p-4 rounded-xl bg-slate-50 border border-slate-200 text-xs text-slate-700 whitespace-pre-line leading-relaxed font-mono">
              {{ version.releaseNotes || 'No specific release notes documented for this version baseline.' }}
            </div>
          </div>

          <!-- Description & Scope -->
          <div class="space-y-2">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
              Architecture Description
            </h4>
            <p class="text-xs text-slate-600 leading-relaxed font-medium whitespace-pre-line">
              {{ version.description }}
            </p>
          </div>

          <!-- Quick Hardware Revisions Summary Table -->
          <div class="space-y-2 pt-2 border-t border-slate-100">
            <div class="flex items-center justify-between">
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
                Hardware Revisions Associated
              </h4>
              <button
                @click="activeTab = 'revisions'"
                class="text-3xs font-semibold text-brand-600 hover:text-brand-800"
              >
                View Full Timeline →
              </button>
            </div>

            <div class="overflow-x-auto">
              <table class="dense-table">
                <thead>
                  <tr>
                    <th>Revision</th>
                    <th>Revision Name</th>
                    <th>Status</th>
                    <th>PCB Rev</th>
                    <th>Engineer</th>
                    <th>Release Date</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="rev in revisions"
                    :key="rev.id"
                    @click="$router.push(`/revisions/${rev.id}`)"
                    class="cursor-pointer hover:bg-slate-50"
                  >
                    <td>
                      <span class="font-mono font-bold text-xs text-slate-900 bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">
                        {{ rev.code }}
                      </span>
                    </td>
                    <td class="font-semibold text-slate-800 text-xs">{{ rev.name }}</td>
                    <td>
                      <RevisionBadge :status="rev.status" />
                    </td>
                    <td class="font-mono text-3xs text-slate-600">{{ rev.pcbRevision }}</td>
                    <td class="text-xs text-slate-700">{{ rev.engineer }}</td>
                    <td class="font-mono text-3xs text-slate-500">{{ rev.releaseDate || rev.createdAt }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- TAB 2: HARDWARE REVISIONS TIMELINE -->
        <div v-else-if="activeTab === 'revisions'" class="space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
                Hardware Revision Evolution for v{{ version.versionNumber }}
              </h4>
              <p class="text-3xs text-slate-400">Chronological PCB layout iterations and change records</p>
            </div>
            <button @click="openAddRevisionModal" class="btn btn-primary">
              <AppIcon name="plus" size="xs" class="mr-1.5" />
              <span>Add Revision</span>
            </button>
          </div>

          <RevisionTimeline :revisions="revisions" />
        </div>

        <!-- TAB 3: DOCUMENTS -->
        <div v-else-if="activeTab === 'documents'" class="space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
                Version Engineering Artifacts
              </h4>
              <p class="text-3xs text-slate-400">Schematic PDF, Gerbers, 3D STEP models, and test certificates</p>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
            <div
              v-for="(doc, idx) in version.documents || []"
              :key="idx"
              class="p-4 rounded-xl border border-slate-200/90 bg-slate-50/50 hover:bg-white hover:border-brand-300 hover:shadow-2xs transition-all flex items-start gap-3 group"
            >
              <div class="w-10 h-10 rounded-lg bg-brand-50 text-brand-600 border border-brand-200/80 flex items-center justify-center shrink-0 shadow-2xs">
                <AppIcon name="file-text" size="sm" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="font-mono font-bold text-xs text-slate-900 truncate group-hover:text-brand-600 transition-colors">
                  {{ doc.name }}
                </div>
                <div class="text-3xs text-slate-400 font-mono mt-0.5 flex items-center gap-2">
                  <span>{{ doc.type || 'Document' }}</span>
                  <span>•</span>
                  <span>{{ doc.size || '1.2 MB' }}</span>
                </div>
                <button
                  @click="downloadDoc(doc.name)"
                  class="mt-2 text-3xs font-semibold text-brand-600 hover:text-brand-800 flex items-center gap-1"
                >
                  <AppIcon name="download" size="xs" />
                  <span>Download Artifact</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 4: AUDIT HISTORY -->
        <div v-else-if="activeTab === 'history'" class="space-y-4">
          <ChangeHistoryTimeline :items="versionHistory" />
        </div>
      </div>
    </div>

    <!-- Quick Add Revision Modal -->
    <BaseModal
      :show="addRevisionModalOpen"
      title="Create Hardware Revision"
      maxWidth="md"
      @close="addRevisionModalOpen = false"
    >
      <div class="space-y-3.5 text-xs">
        <div>
          <label class="form-label">Revision Code (e.g. REV-C)</label>
          <input
            v-model="newRevisionForm.code"
            type="text"
            placeholder="REV-C"
            class="form-input font-mono uppercase"
          />
        </div>
        <div>
          <label class="form-label">Revision Name</label>
          <input
            v-model="newRevisionForm.name"
            type="text"
            placeholder="e.g. Sensor Connector & Power Trace Improvement"
            class="form-input"
          />
        </div>
        <div>
          <label class="form-label">Status</label>
          <select v-model="newRevisionForm.status" class="form-select font-mono">
            <option value="Draft">Draft</option>
            <option value="In Review">In Review</option>
            <option value="Approved">Approved</option>
            <option value="Released">Released</option>
          </select>
        </div>
        <div>
          <label class="form-label">Change Summary</label>
          <textarea
            v-model="newRevisionForm.changeSummary"
            rows="3"
            placeholder="Describe hardware changes, trace modifications, or capacitor footprint adjustments..."
            class="form-textarea"
          ></textarea>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="addRevisionModalOpen = false">Cancel</button>
        <button class="btn btn-primary" @click="saveNewRevision">Create Revision</button>
      </template>
    </BaseModal>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import VersionBadge from '@/components/versions/VersionBadge.vue';
import RevisionBadge from '@/components/versions/RevisionBadge.vue';
import VersionTimeline from '@/components/versions/VersionTimeline.vue';
import RevisionTimeline from '@/components/versions/RevisionTimeline.vue';
import ChangeHistoryTimeline from '@/components/versions/ChangeHistoryTimeline.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ProductVersionDetailView',
  components: {
    VersionBadge,
    RevisionBadge,
    VersionTimeline,
    RevisionTimeline,
    ChangeHistoryTimeline,
    BaseModal,
    AppIcon
  },
  data() {
    return {
      activeTab: 'overview',
      addRevisionModalOpen: false,
      newRevisionForm: {
        code: 'REV-C',
        name: '',
        status: 'Draft',
        changeSummary: '',
        engineer: 'Alex Chen'
      }
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    versionId() {
      return this.$route.params.id;
    },
    version() {
      return this.versionStore.getVersionById(this.versionId) || this.versionStore.versions[0];
    },
    productVersions() {
      if (!this.version) return [];
      return this.versionStore.getVersionsForProduct(this.version.productCode);
    },
    revisions() {
      if (!this.version) return [];
      return this.versionStore.getRevisionsForVersion(this.version.id);
    },
    versionHistory() {
      if (!this.version) return [];
      return this.versionStore.versionHistory.filter(h => h.version === `v${this.version.versionNumber}` || h.targetId === this.version.id);
    }
  },
  methods: {
    onSelectTimelineVersion(ver) {
      if (ver && ver.id !== this.version.id) {
        this.$router.push(`/versions/${ver.id}`);
      }
    },
    openAddRevisionModal() {
      this.newRevisionForm = {
        code: `REV-${String.fromCharCode(65 + this.revisions.length)}`,
        name: '',
        status: 'Draft',
        changeSummary: '',
        engineer: 'Alex Chen'
      };
      this.addRevisionModalOpen = true;
    },
    saveNewRevision() {
      if (!this.newRevisionForm.code || !this.newRevisionForm.name) {
        alert('Please specify revision code and name.');
        return;
      }
      this.versionStore.saveRevision({
        ...this.newRevisionForm,
        versionId: this.version.id,
        versionNumber: this.version.versionNumber,
        productCode: this.version.productCode,
        productName: this.version.productName,
        pcbRevision: `PCB-v${this.version.versionNumber}-${this.newRevisionForm.code.replace('REV-', '')}`,
        changesList: [this.newRevisionForm.changeSummary]
      });
      this.addRevisionModalOpen = false;
    },
    downloadDoc(name) {
      alert(`Downloading engineering file: ${name}`);
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
