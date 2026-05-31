<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getDashboard, getRebates, getWithdrawals } from '../../api'
import PageSection from '../../components/common/PageSection.vue'
import StatCard from '../../components/common/StatCard.vue'
import DataTable from '../../components/common/DataTable.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import { useSession } from '../../session/useSession'
import type { DashboardSummary, RebateItem, WithdrawalItem } from '../../types'
import { formatCurrency, formatDateTime } from '../../utils/format'

const session = useSession()

const dashboard = ref<DashboardSummary | null>(null)
const recentRebates = ref<RebateItem[]>([])
const recentWithdrawals = ref<WithdrawalItem[]>([])
const loadingError = ref('')

const withdrawalPreview = computed(() => recentWithdrawals.value.slice(0, 5))
const rebatePreview = computed(() => recentRebates.value.slice(0, 5))

async function loadData(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''

  try {
    const [dashboardResp, rebateResp, withdrawalResp] = await Promise.all([
      getDashboard(session.state.token),
      getRebates(session.state.token),
      getWithdrawals(session.state.token),
    ])
    dashboard.value = dashboardResp
    recentRebates.value = rebateResp
    recentWithdrawals.value = withdrawalResp
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载失败'
  }
}

onMounted(loadData)
</script>

<template>
  <div class="page-stack">
    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>

    <section class="stats-grid">
      <StatCard label="累计返利" :value="formatCurrency(dashboard?.total_earned)" hint="已累计产生的返利总额" />
      <StatCard label="可申请金额" :value="formatCurrency(dashboard?.withdrawable_amount)" hint="可立即发起提现申请" />
      <StatCard label="打款中" :value="formatCurrency(dashboard?.paying_amount)" hint="已经提交、待线下打款" />
      <StatCard label="已打款" :value="formatCurrency(dashboard?.paid_amount)" hint="已完成打款的历史累计" />
    </section>

    <PageSection title="近期提现申请" description="这里展示最近发起的提现记录和当前状态。" compact>
      <DataTable
        :items="withdrawalPreview"
        :columns="[
          { key: 'request_no', title: '申请单号', width: '220px' },
          { key: 'amount', title: '申请金额', width: '120px', align: 'right' },
          { key: 'status', title: '状态', width: '120px' },
          { key: 'created_at', title: '申请时间', width: '180px' },
        ]"
        :row-key="(item) => item.id"
        empty-title="还没有提现申请"
        empty-description="发起第一笔提现后，这里会展示最近申请记录。"
      >
        <template #amount="{ item }">
          {{ formatCurrency(item.amount) }}
        </template>
        <template #status="{ item }">
          <StatusBadge :status="item.status" />
        </template>
        <template #created_at="{ item }">
          {{ formatDateTime(item.created_at) }}
        </template>
      </DataTable>
    </PageSection>

    <PageSection title="近期返利明细" description="帮助快速确认最近有哪些返利入账。" compact>
      <DataTable
        :items="rebatePreview"
        :columns="[
          { key: 'source_email', title: '来源用户', width: '220px' },
          { key: 'amount', title: '返利金额', width: '120px', align: 'right' },
          { key: 'source_order_id', title: '来源订单', width: '160px' },
          { key: 'created_at', title: '入账时间', width: '180px' },
        ]"
        :row-key="(item) => item.ledger_id"
        empty-title="还没有返利记录"
        empty-description="当被邀请用户触发返利后，这里会展示最近流水。"
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
  </div>
</template>
