<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
}>()
</script>

<template>
  <div class="empty-state">
    <strong>{{ title || '暂无数据' }}</strong>
    <p>{{ description || '当前条件下还没有可展示的内容。' }}</p>
  </div>
</template>
