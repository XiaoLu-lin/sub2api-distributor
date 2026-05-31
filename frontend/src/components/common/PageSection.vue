<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  compact?: boolean
}>()
</script>

<template>
  <section class="page-section" :class="{ 'page-section-compact': compact }">
    <header class="page-section-header">
      <div>
        <h2>{{ title }}</h2>
        <p v-if="description">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="page-section-actions">
        <slot name="actions" />
      </div>
    </header>

    <div class="page-section-body">
      <slot />
    </div>
  </section>
</template>
