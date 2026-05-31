<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    width?: 'normal' | 'wide' | 'extra-wide'
    closeOnMask?: boolean
  }>(),
  {
    width: 'normal',
    closeOnMask: true,
  },
)

const emit = defineEmits<{
  (e: 'close'): void
}>()

const widthClass = computed(() => ({
  'base-modal-panel': true,
  'base-modal-wide': props.width === 'wide',
  'base-modal-extra-wide': props.width === 'extra-wide',
}))

function onMaskClick(): void {
  if (props.closeOnMask) {
    emit('close')
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="base-modal-root">
      <div class="base-modal-backdrop" @click="onMaskClick"></div>
      <div :class="widthClass" role="dialog" aria-modal="true">
        <header class="base-modal-header">
          <div>
            <h3>{{ title }}</h3>
          </div>
          <button class="icon-button" type="button" @click="$emit('close')">×</button>
        </header>

        <div class="base-modal-body">
          <slot />
        </div>

        <footer v-if="$slots.footer" class="base-modal-footer">
          <slot name="footer" />
        </footer>
      </div>
    </div>
  </Teleport>
</template>
