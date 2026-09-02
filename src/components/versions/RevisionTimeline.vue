<template>
  <div class="space-y-4">
    <div class="relative pl-6 space-y-6 before:absolute before:left-2.5 before:top-3 before:bottom-3 before:w-0.5 before:bg-slate-200">
      <div
        v-for="rev in revisions"
        :key="rev.id"
        class="relative group"
      >
        <!-- Timeline node dot -->
        <span
          :class="[
            'absolute -left-6 top-1.5 w-5 h-5 rounded-full flex items-center justify-center font-mono font-bold text-3xs ring-4 ring-white shadow-2xs transition-transform group-hover:scale-110',
            rev.status === 'Released' ? 'bg-emerald-600 text-white' :
            rev.status === 'Approved' ? 'bg-brand-600 text-white' :
            rev.status === 'In Review' ? 'bg-amber-500 text-white' :
            'bg-slate-200 text-slate-600'
          ]"
        >
          {{ rev.code.replace('REV-', '') }}
        </span>

        <!-- Revision Card -->
        <div class="app-card p-4 hover:border-brand-300 transition-all">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-2.5 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <span class="font-mono font-bold text-sm text-slate-900 bg-slate-100 px-2 py-0.5 rounded border border-slate-200">
                {{ rev.code }}
              </span>
              <h4 class="font-bold text-slate-800 text-xs truncate">
                {{ rev.name }}
              </h4>
              <RevisionBadge :status="rev.status" />
            </div>

            <div class="flex items-center gap-2 text-3xs font-mono text-slate-400 shrink-0">
              <span v-if="rev.pcbRevision" class="px-1.5 py-0.5 rounded bg-slate-50 border border-slate-200 text-slate-600">
                {{ rev.pcbRevision }}
              </span>
              <span>•</span>
              <span>{{ rev.releaseDate || rev.createdAt }}</span>
            </div>
          </div>

          <!-- Change Summary & Points -->
          <p class="text-xs text-slate-600 mt-2.5 leading-relaxed font-medium">
            {{ rev.changeSummary }}
          </p>

          <div v-if="rev.changesList && rev.changesList.length > 0" class="mt-2.5 space-y-1 pl-1">
            <div
              v-for="(pt, idx) in rev.changesList"
              :key="idx"
              class="text-3xs text-slate-600 flex items-start gap-1.5 font-medium"
            >
              <span class="text-brand-500 font-bold">•</span>
              <span>{{ pt }}</span>
            </div>
          </div>

          <!-- Bottom Footer: Engineer & Attached Documents -->
          <div class="mt-3 pt-2.5 border-t border-slate-100 flex flex-wrap items-center justify-between gap-2 text-3xs text-slate-400">
            <div class="flex items-center gap-2">
              <span class="font-medium text-slate-700">Engineer:</span>
              <span class="font-semibold text-slate-900">{{ rev.engineer || 'Alex Chen' }}</span>
              <span v-if="rev.approvedBy" class="text-slate-400">(Approved by: {{ rev.approvedBy }})</span>
            </div>

            <div v-if="rev.documents && rev.documents.length > 0" class="flex items-center gap-1.5">
              <span
                v-for="(doc, dIdx) in rev.documents"
                :key="dIdx"
                class="px-2 py-0.5 rounded bg-brand-50 text-brand-700 border border-brand-200 font-mono font-medium flex items-center gap-1"
              >
                <AppIcon name="file" size="xs" />
                <span>{{ doc.name || doc }}</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import RevisionBadge from './RevisionBadge.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'RevisionTimeline',
  components: { RevisionBadge, AppIcon },
  props: {
    revisions: {
      type: Array,
      required: true
    }
  }
};
</script>
