<template>
  <div v-if="confidence > 0" class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-3xs font-semibold bg-purple-50 text-purple-700 border border-purple-200">
    <AppIcon name="sparkles" size="xs" class="text-purple-600 animate-pulse" />
    <span>{{ label || 'AI' }}:</span>
    <span class="font-bold font-mono">{{ confidence }}%</span>
  </div>
  <div v-else class="text-3xs text-slate-400 font-mono italic">
    Heuristic
  </div>
</template>

<script>
import AppIcon from '../common/AppIcon.vue';

export default {
  name: 'AiSuggestionBadge',
  components: { AppIcon },
  props: {
    confidence: {
      type: Number,
      default: 0
    },
    label: {
      type: String,
      default: ''
    }
  }
};
</script>
