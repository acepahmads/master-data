<template>
  <div class="mb-6 flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-4 border-b border-slate-200/80">
    <div class="min-w-0">
      <div class="flex items-center gap-2.5">
        <h1 class="text-xl md:text-2xl font-bold text-slate-900 tracking-tight">{{ title }}</h1>
        <slot name="badge" />
      </div>
      <p v-if="subtitle" class="text-xs text-slate-500 mt-1 leading-normal font-medium max-w-3xl">
        {{ subtitle }}
      </p>
    </div>
    <div class="flex items-center gap-2.5 shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'PageHeader',
  props: {
    title: {
      type: String,
      required: true
    },
    subtitle: {
      type: String,
      default: ''
    }
  }
};
</script>
