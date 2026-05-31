<script setup lang="ts">
import BaseModal from './BaseModal.vue'

defineProps<{
  open: boolean
  title: string
  description: string
  confirmText?: string
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm'): void
}>()
</script>

<template>
  <BaseModal :open="open" :title="title" @close="emit('close')">
    <p class="confirm-text">{{ description }}</p>

    <template #footer>
      <div class="dialog-actions">
        <button class="ghost-button" type="button" @click="emit('close')">取消</button>
        <button type="button" :disabled="loading" @click="emit('confirm')">
          {{ loading ? '处理中...' : confirmText || '确认' }}
        </button>
      </div>
    </template>
  </BaseModal>
</template>
