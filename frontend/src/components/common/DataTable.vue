<script setup lang="ts" generic="T extends Record<string, unknown>">
import EmptyState from './EmptyState.vue'

defineProps<{
  title?: string
  description?: string
  items: T[]
  maxHeight?: string
  columns: Array<{
    key: string
    title: string
    width?: string
    align?: 'left' | 'center' | 'right'
  }>
  rowKey: (item: T) => string | number
  emptyTitle?: string
  emptyDescription?: string
}>()
</script>

<template>
  <div class="table-card">
    <header v-if="title || description || $slots.toolbar" class="table-card-header">
      <div>
        <h3 v-if="title">{{ title }}</h3>
        <p v-if="description">{{ description }}</p>
      </div>
      <div v-if="$slots.toolbar" class="table-card-toolbar">
        <slot name="toolbar" />
      </div>
    </header>

    <div
      v-if="items.length"
      class="table-scroll"
      :style="maxHeight ? { '--table-max-height': maxHeight } : undefined"
    >
      <table class="data-table">
        <thead>
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :style="{ width: column.width }"
              :class="column.align ? `align-${column.align}` : ''"
            >
              {{ column.title }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="rowKey(item)">
            <td
              v-for="column in columns"
              :key="column.key"
              :class="column.align ? `align-${column.align}` : ''"
            >
              <slot :name="column.key" :item="item">
                {{ item[column.key] }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <EmptyState
      v-else
      :title="emptyTitle || '暂无记录'"
      :description="emptyDescription || '当前还没有可以展示的记录。'"
    />
  </div>
</template>
