<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { cancelWithdrawal, createWithdrawal, getDashboard, getWithdrawals } from '../../api'
import BaseModal from '../../components/common/BaseModal.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import DataTable from '../../components/common/DataTable.vue'
import PageSection from '../../components/common/PageSection.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import ToolbarActions from '../../components/common/ToolbarActions.vue'
import { useSession } from '../../session/useSession'
import type { DashboardSummary, WithdrawalItem } from '../../types'
import { formatCurrency, formatDateTime } from '../../utils/format'

const session = useSession()

const withdrawals = ref<WithdrawalItem[]>([])
const dashboard = ref<DashboardSummary | null>(null)
const loadingError = ref('')

const createOpen = ref(false)
const createLoading = ref(false)
const createError = ref('')
const withdrawalAmount = ref<number | null>(null)
const withdrawalRemark = ref('')

const cancelOpen = ref(false)
const cancelLoading = ref(false)
const currentWithdrawal = ref<WithdrawalItem | null>(null)

const withdrawableText = computed(() => formatCurrency(dashboard.value?.withdrawable_amount))

async function loadData(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    const [withdrawalResp, dashboardResp] = await Promise.all([
      getWithdrawals(session.state.token),
      getDashboard(session.state.token),
    ])
    withdrawals.value = withdrawalResp
    dashboard.value = dashboardResp
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载提现申请失败'
  }
}

function openCreateModal(): void {
  withdrawalAmount.value = null
  withdrawalRemark.value = ''
  createError.value = ''
  createOpen.value = true
}

async function submitCreate(): Promise<void> {
  if (!session.state.token) {
    return
  }
  if (!withdrawalAmount.value || withdrawalAmount.value <= 0) {
    createError.value = '请输入大于 0 的提现金额'
    return
  }

  createLoading.value = true
  createError.value = ''

  try {
    await createWithdrawal(session.state.token, withdrawalAmount.value, withdrawalRemark.value)
    createOpen.value = false
    await loadData()
  } catch (error) {
    createError.value = error instanceof Error ? error.message : '提交提现申请失败'
  } finally {
    createLoading.value = false
  }
}

function openCancelDialog(item: WithdrawalItem): void {
  currentWithdrawal.value = item
  cancelOpen.value = true
}

async function submitCancel(): Promise<void> {
  if (!session.state.token || !currentWithdrawal.value) {
    return
  }

  cancelLoading.value = true
  try {
    await cancelWithdrawal(session.state.token, currentWithdrawal.value.id)
    cancelOpen.value = false
    currentWithdrawal.value = null
    await loadData()
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '取消提现申请失败'
  } finally {
    cancelLoading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <PageSection title="提现申请" description="分销商提交申请后立即显示打款中，运营线下处理后再标记已打款。">
    <template #actions>
      <ToolbarActions>
        <span class="toolbar-tip">当前可申请金额：{{ withdrawableText }}</span>
        <button type="button" @click="openCreateModal">发起提现</button>
      </ToolbarActions>
    </template>

    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>

    <DataTable
      :items="withdrawals"
      :columns="[
        { key: 'request_no', title: '申请单号', width: '220px' },
        { key: 'amount', title: '申请金额', width: '120px', align: 'right' },
        { key: 'status', title: '状态', width: '120px' },
        { key: 'snapshot_withdrawable_before', title: '申请前可提', width: '140px', align: 'right' },
        { key: 'created_at', title: '申请时间', width: '180px' },
        { key: 'actions', title: '操作', width: '140px', align: 'center' },
      ]"
      :row-key="(item) => item.id"
      empty-title="暂无提现申请"
      empty-description="发起提现申请后，列表会按状态展示每一笔申请。"
    >
      <template #amount="{ item }">
        {{ formatCurrency(item.amount) }}
      </template>
      <template #status="{ item }">
        <StatusBadge :status="item.status" />
      </template>
      <template #snapshot_withdrawable_before="{ item }">
        {{ formatCurrency(item.snapshot_withdrawable_before) }}
      </template>
      <template #created_at="{ item }">
        {{ formatDateTime(item.created_at) }}
      </template>
      <template #actions="{ item }">
        <button
          v-if="item.status === 'paying'"
          class="link-button"
          type="button"
          @click="openCancelDialog(item)"
        >
          取消申请
        </button>
        <span v-else class="table-muted">-</span>
      </template>
    </DataTable>

    <BaseModal :open="createOpen" title="发起提现申请" @close="createOpen = false">
      <div class="form-stack">
        <p class="modal-helper">提交成功后，该笔申请会立即进入“打款中”状态。</p>
        <label class="field-block">
          <span>申请金额</span>
          <input v-model.number="withdrawalAmount" type="number" min="0" step="0.01" placeholder="请输入提现金额" />
        </label>
        <label class="field-block">
          <span>备注</span>
          <textarea v-model="withdrawalRemark" rows="4" placeholder="可选，补充给运营的说明"></textarea>
        </label>
        <p class="info-text">当前可申请金额：{{ withdrawableText }}</p>
        <p v-if="createError" class="error-text">{{ createError }}</p>
      </div>

      <template #footer>
        <div class="dialog-actions">
          <button class="ghost-button" type="button" @click="createOpen = false">取消</button>
          <button type="button" :disabled="createLoading" @click="submitCreate">
            {{ createLoading ? '提交中...' : '确认提交' }}
          </button>
        </div>
      </template>
    </BaseModal>

    <ConfirmDialog
      :open="cancelOpen"
      title="取消提现申请"
      description="取消后，这笔金额会释放回可申请金额。确认继续吗？"
      confirm-text="确认取消"
      :loading="cancelLoading"
      @close="cancelOpen = false"
      @confirm="submitCancel"
    />
  </PageSection>
</template>
