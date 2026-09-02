<template>
  <BaseModal :show="show" :title="title" maxWidth="sm" @close="cancel">
    <div class="flex items-start gap-3">
      <div class="w-9 h-9 rounded-full bg-rose-50 border border-rose-200 flex items-center justify-center text-rose-600 shrink-0">
        <AppIcon name="alert-circle" size="md" />
      </div>
      <div>
        <p class="text-xs text-slate-600 leading-relaxed">{{ message }}</p>
      </div>
    </div>

    <template #footer>
      <button type="button" class="btn btn-secondary" @click="cancel">
        {{ cancelText }}
      </button>
      <button type="button" class="btn btn-danger" @click="confirm">
        {{ confirmText }}
      </button>
    </template>
  </BaseModal>
</template>

<script>
import BaseModal from './BaseModal.vue';
import AppIcon from './AppIcon.vue';

export default {
  name: 'ConfirmDialog',
  components: { BaseModal, AppIcon },
  props: {
    show: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: 'Confirm Deletion'
    },
    message: {
      type: String,
      default: 'Are you sure you want to delete this record? This action cannot be undone.'
    },
    confirmText: {
      type: String,
      default: 'Delete'
    },
    cancelText: {
      type: String,
      default: 'Cancel'
    }
  },
  methods: {
    confirm() {
      this.$emit('confirm');
    },
    cancel() {
      this.$emit('cancel');
    }
  }
};
</script>
