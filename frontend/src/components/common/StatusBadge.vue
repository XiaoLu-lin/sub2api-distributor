<script setup lang="ts">
import { computed } from 'vue'
import { statusText } from '../../utils/format'

const props = defineProps<{
  status: string
}>()

const badgeClass = computed(() => ({
  'status-badge': true,
  'status-badge-paying': props.status === 'paying',
  'status-badge-paid': props.status === 'paid',
  'status-badge-cancelled': props.status === 'cancelled',
  'status-badge-active': props.status === 'active',
  'status-badge-disabled': props.status === 'disabled',
}))
</script>

<template>
  <span :class="badgeClass">{{ statusText(status) }}</span>
</template>
