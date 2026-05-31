<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getRebates } from '../../api'
import PageSection from '../../components/common/PageSection.vue'
import DataTable from '../../components/common/DataTable.vue'
import { useSession } from '../../session/useSession'
import type { RebateItem } from '../../types'
import { formatCurrency, formatDateTime } from '../../utils/format'

const session = useSession()
const rebates = ref<RebateItem[]>([])
const loadingError = ref('')

async function loadRebates(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    rebates.value = await getRebates(session.state.token)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载返利明细失败'
  }
}

onMounted(loadRebates)
</script>

<template>
  <PageSection title="返利明细" description="每一笔返利流水都会在这里累计展示。">
    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>
    <DataTable
      :items="rebates"
      :columns="[
        { key: 'ledger_id', title: '流水 ID', width: '120px' },
        { key: 'amount', title: '返利金额', width: '140px', align: 'right' },
        { key: 'source_email', title: '来源用户', width: '260px' },
        { key: 'source_order_id', title: '来源订单', width: '160px' },
        { key: 'created_at', title: '入账时间', width: '180px' },
      ]"
      :row-key="(item) => item.ledger_id"
      empty-title="暂无返利流水"
      empty-description="当被邀请用户完成触发返利的行为后，这里会出现对应明细。"
    >
      <template #amount="{ item }">
        {{ formatCurrency(item.amount) }}
      </template>
      <template #source_order_id="{ item }">
        {{ item.source_order_id || '-' }}
      </template>
      <template #created_at="{ item }">
        {{ formatDateTime(item.created_at) }}
      </template>
    </DataTable>
  </PageSection>
</template>
