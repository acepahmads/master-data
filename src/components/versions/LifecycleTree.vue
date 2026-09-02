<template>
  <div class="space-y-4">
    <div
      v-for="group in productGroups"
      :key="group.productCode"
      class="app-card overflow-hidden"
    >
      <!-- Root Product Header (Level 1) -->
      <div
        @click="toggleGroup(group.productCode)"
        class="px-4 py-3 bg-slate-50/80 border-b border-slate-200/80 flex items-center justify-between cursor-pointer hover:bg-slate-100/80 transition-colors select-none"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <button class="p-1 rounded text-slate-400 hover:text-slate-700">
            <AppIcon
              :name="expandedGroups[group.productCode] ? 'chevron-down' : 'chevron-right'"
              size="xs"
            />
          </button>
          <div class="w-7 h-7 rounded-lg bg-brand-50 border border-brand-200 text-brand-600 flex items-center justify-center shrink-0">
            <AppIcon name="product" size="xs" />
          </div>
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span class="font-bold text-slate-900 text-xs truncate">{{ group.productName }}</span>
              <span class="font-mono font-bold text-3xs text-brand-700 bg-brand-50 px-1.5 py-0.2 rounded border border-brand-200">
                {{ group.productCode }}
              </span>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span class="text-3xs font-mono font-semibold text-slate-500 bg-white px-2 py-0.5 rounded-full border border-slate-200">
            {{ group.versions.length }} Versions
          </span>
          <router-link
            :to="`/versions/create?product=${group.productCode}`"
            @click.native.stop
            class="text-3xs font-semibold text-brand-600 hover:text-brand-800 px-2 py-1 rounded hover:bg-brand-50"
            title="Add Version"
          >
            + Add Version
          </router-link>
        </div>
      </div>

      <!-- Versions & Revisions Hierarchy (Level 2 & Level 3) -->
      <div
        v-show="expandedGroups[group.productCode]"
        class="p-3.5 space-y-3 bg-white"
      >
        <div
          v-for="ver in group.versions"
          :key="ver.id"
          class="rounded-lg border border-slate-200/80 bg-slate-50/40 p-3 hover:border-slate-300 transition-all space-y-2.5"
        >
          <!-- Version Node Header -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div class="flex items-center gap-2 min-w-0">
              <div class="w-6 h-6 rounded-md bg-white border border-slate-200 text-slate-700 flex items-center justify-center font-mono font-bold text-3xs shrink-0 shadow-2xs">
                v{{ ver.versionNumber.split('.')[0] }}
              </div>
              <router-link
                :to="`/versions/${ver.id}`"
                class="font-mono font-bold text-xs text-slate-900 hover:text-brand-600 transition-colors flex items-center gap-1.5"
              >
                <span>v{{ ver.versionNumber }}</span>
                <span class="font-sans font-medium text-slate-600 text-3xs truncate max-w-[200px]">({{ ver.versionName }})</span>
              </router-link>
              <VersionBadge :status="ver.status" />
            </div>

            <div class="flex items-center gap-2 text-3xs font-mono text-slate-400 shrink-0">
              <span class="text-slate-600 font-semibold">{{ ver.owner }}</span>
              <span>•</span>
              <span>{{ ver.releaseDate || ver.createdAt }}</span>
              <router-link
                :to="`/versions/${ver.id}`"
                class="btn-icon p-1"
                title="View Version Details"
              >
                <AppIcon name="eye" size="xs" />
              </router-link>
            </div>
          </div>

          <!-- Hardware Revisions Connected Branch (Level 3) -->
          <div class="pl-5 space-y-1.5 border-l-2 border-slate-200 ml-3">
            <div
              v-for="rev in getRevisionsForVersion(ver.id)"
              :key="rev.id"
              class="flex items-center justify-between p-2 rounded bg-white border border-slate-200/60 hover:border-brand-200 text-xs transition-colors group"
            >
              <div class="flex items-center gap-2 min-w-0">
                <span class="w-1.5 h-1.5 rounded-full bg-slate-300 group-hover:bg-brand-500 shrink-0"></span>
                <router-link
                  :to="`/revisions/${rev.id}`"
                  class="font-mono font-bold text-slate-900 hover:text-brand-600 transition-colors"
                >
                  {{ rev.code }}
                </router-link>
                <span class="text-3xs text-slate-500 truncate max-w-xs font-medium">{{ rev.name }}</span>
              </div>

              <div class="flex items-center gap-2 shrink-0">
                <RevisionBadge :status="rev.status" />
                <span class="text-3xs font-mono text-slate-400">{{ rev.releaseDate || rev.createdAt }}</span>
                <router-link
                  :to="`/revisions/${rev.id}`"
                  class="p-1 text-slate-400 hover:text-brand-600 rounded"
                  title="View Revision"
                >
                  <AppIcon name="eye" size="xs" />
                </router-link>
              </div>
            </div>

            <!-- Empty Revisions Placeholder -->
            <div
              v-if="getRevisionsForVersion(ver.id).length === 0"
              class="text-3xs text-slate-400 italic py-1 pl-2"
            >
              No hardware revisions recorded for this version.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useVersionStore } from '@/stores/versionStore';
import VersionBadge from './VersionBadge.vue';
import RevisionBadge from './RevisionBadge.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'LifecycleTree',
  components: {
    VersionBadge,
    RevisionBadge,
    AppIcon
  },
  props: {
    filterProductCode: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      expandedGroups: {}
    };
  },
  computed: {
    versionStore() {
      return useVersionStore();
    },
    productGroups() {
      const all = this.versionStore.versionsByProduct;
      if (!this.filterProductCode) return all;
      return all.filter(g => g.productCode === this.filterProductCode);
    }
  },
  watch: {
    productGroups: {
      immediate: true,
      handler(groups) {
        if (groups && groups.length > 0) {
          groups.forEach(g => {
            if (this.expandedGroups[g.productCode] === undefined) {
              this.$set(this.expandedGroups, g.productCode, true);
            }
          });
        }
      }
    }
  },
  async mounted() {
    if (this.versionStore.versions.length === 0) {
      await Promise.all([
        this.versionStore.fetchVersions(),
        this.versionStore.fetchRevisions()
      ]);
    }
  },
  methods: {
    toggleGroup(productCode) {
      this.$set(this.expandedGroups, productCode, !this.expandedGroups[productCode]);
    },
    getRevisionsForVersion(versionId) {
      return this.versionStore.getRevisionsForVersion(versionId);
    }
  }
};
</script>
