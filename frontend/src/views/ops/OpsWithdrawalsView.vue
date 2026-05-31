<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  cancelOpsWithdrawal,
  getOpsWithdrawalDetail,
  getOpsWithdrawals,
  markPaid,
} from '../../api'
import BaseModal from '../../components/common/BaseModal.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import DataTable from '../../components/common/DataTable.vue'
import PageSection from '../../components/common/PageSection.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import ToolbarActions from '../../components/common/ToolbarActions.vue'
import { useSession } from '../../session/useSession'
import type { WithdrawalDetailItem, WithdrawalItem } from '../../types'
import { eventActionText, formatCurrency, formatDateTime, formatEventDetail } from '../../utils/format'

const session = useSession()

const withdrawals = ref<WithdrawalItem[]>([])
const loadingError = ref('')
const statusFilter = ref<'all' | 'paying' | 'paid' | 'cancelled'>('all')

const detailOpen = ref(false)
const detailLoading = ref(false)
const detail = ref<WithdrawalDetailItem | null>(null)

const paidOpen = ref(false)
const paidLoading = ref(false)
const paidChannel = ref('bank')
const paidReferenceNo = ref('')
const paidRemark = ref('')

const cancelOpen = ref(false)
const cancelLoading = ref(false)
const currentWithdrawal = ref<WithdrawalItem | null>(null)

const filteredWithdrawals = computed(() =>
  statusFilter.value === 'all'
    ? withdrawals.value
    : withdrawals.value.filter((item) => item.status === statusFilter.value),
)

async function loadWithdrawals(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    withdrawals.value = await getOpsWithdrawals(session.state.token)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载提现管理失败'
  }
}

async function openDetail(item: WithdrawalItem): Promise<void> {
  if (!session.state.token) {
    return
  }

  currentWithdrawal.value = item
  detailOpen.value = true
  detailLoading.value = true

  try {
    detail.value = await getOpsWithdrawalDetail(session.state.token, item.id)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载提现详情失败'
  } finally {
    detailLoading.value = false
  }
}

function openPaidModal(item: WithdrawalItem): void {
  currentWithdrawal.value = item
  paidChannel.value = 'bank'
  paidReferenceNo.value = ''
  paidRemark.value = ''
  paidOpen.value = true
}

async function submitPaid(): Promise<void> {
  if (!session.state.token || !currentWithdrawal.value) {
    return
  }

  paidLoading.value = true
  try {
    await markPaid(
      session.state.token,
      currentWithdrawal.value.id,
      paidChannel.value,
      paidReferenceNo.value,
      paidRemark.value,
    )
    paidOpen.value = false
    await loadWithdrawals()
    if (detailOpen.value) {
      await openDetail(currentWithdrawal.value)
    }
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '标记已打款失败'
  } finally {
    paidLoading.value = false
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
    await cancelOpsWithdrawal(session.state.token, currentWithdrawal.value.id)
    cancelOpen.value = false
    await loadWithdrawals()
    if (detailOpen.value) {
      detailOpen.value = false
      detail.value = null
    }
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '取消申请失败'
  } finally {
    cancelLoading.value = false
  }
}

onMounted(loadWithdrawals)
</script>

<template>
  <PageSection title="提现管理" description="运营在这里处理所有分销商提现申请，完成线下打款后手动标记。">
    <template #actions>
      <ToolbarActions>
        <label class="compact-select">
          <span>状态筛选</span>
          <select v-model="statusFilter">
            <option value="all">全部</option>
            <option value="paying">打款中</option>
            <option value="paid">已打款</option>
            <option value="cancelled">已取消</option>
          </select>
        </label>
      </ToolbarActions>
    </template>

    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>

    <DataTable
      :items="filteredWithdrawals"
      :columns="[
        { key: 'request_no', title: '申请单号', width: '220px' },
        { key: 'user_id', title: '申请人 ID', width: '110px' },
        { key: 'amount', title: '申请金额', width: '120px', align: 'right' },
        { key: 'status', title: '状态', width: '120px' },
        { key: 'paid_channel', title: '打款渠道', width: '120px' },
        { key: 'created_at', title: '申请时间', width: '180px' },
        { key: 'actions', title: '操作', width: '220px', align: 'center' },
      ]"
      :row-key="(item) => item.id"
      empty-title="暂无提现申请"
      empty-description="分销商提交提现申请后，这里会展示所有待处理和已处理记录。"
    >
      <template #amount="{ item }">
        {{ formatCurrency(item.amount) }}
      </template>
      <template #status="{ item }">
        <StatusBadge :status="item.status" />
      </template>
      <template #paid_channel="{ item }">
        {{ item.paid_channel || '-' }}
      </template>
      <template #created_at="{ item }">
        {{ formatDateTime(item.created_at) }}
      </template>
      <template #actions="{ item }">
        <div class="table-actions">
          <button class="link-button" type="button" @click="openDetail(item)">查看详情</button>
          <button
            v-if="item.status === 'paying'"
            class="link-button"
            type="button"
            @click="openPaidModal(item)"
          >
            标记已打款
          </button>
          <button
            v-if="item.status === 'paying'"
            class="link-button danger-button"
            type="button"
            @click="openCancelDialog(item)"
          >
            取消申请
          </button>
        </div>
      </template>
    </DataTable>

    <BaseModal :open="detailOpen" title="提现申请详情" width="extra-wide" @close="detailOpen = false">
      <div v-if="detailLoading" class="info-text">加载详情中...</div>
      <div v-else-if="detail" class="detail-stack">
        <div class="detail-grid">
          <div class="detail-card">
            <span>申请单号</span>
            <strong>{{ detail.request_no }}</strong>
          </div>
          <div class="detail-card">
            <span>申请金额</span>
            <strong>{{ formatCurrency(detail.amount) }}</strong>
          </div>
          <div class="detail-card">
            <span>当前状态</span>
            <StatusBadge :status="detail.status" />
          </div>
          <div class="detail-card">
            <span>申请时间</span>
            <strong>{{ formatDateTime(detail.created_at) }}</strong>
          </div>
        </div>

        <div class="detail-snapshot">
          <h4>申请快照</h4>
          <div class="detail-grid">
            <div class="detail-card">
              <span>申请前可提现</span>
              <strong>{{ formatCurrency(detail.snapshot_withdrawable_before) }}</strong>
            </div>
            <div class="detail-card">
              <span>申请后可提现</span>
              <strong>{{ formatCurrency(detail.snapshot_withdrawable_after) }}</strong>
            </div>
            <div class="detail-card">
              <span>打款渠道</span>
              <strong>{{ detail.paid_channel || '-' }}</strong>
            </div>
            <div class="detail-card">
              <span>打款流水号</span>
              <strong>{{ detail.paid_reference_no || '-' }}</strong>
            </div>
          </div>
        </div>

        <div class="detail-snapshot">
          <h4>操作时间线</h4>
          <ul class="timeline-list">
            <li v-for="event in detail.events" :key="event.id">
              <strong>{{ eventActionText(event.action) }}</strong>
              <span>{{ formatDateTime(event.created_at) }}</span>
              <p>{{ formatEventDetail(event.action, event.detail) }}</p>
            </li>
          </ul>
        </div>
      </div>
    </BaseModal>

    <BaseModal :open="paidOpen" title="标记已打款" @close="paidOpen = false">
      <div class="form-stack">
        <label class="field-block">
          <span>打款渠道</span>
          <select v-model="paidChannel">
            <option value="bank">银行卡</option>
            <option value="alipay">支付宝</option>
            <option value="wechat">微信</option>
            <option value="usdt">USDT</option>
            <option value="manual">人工协商</option>
          </select>
        </label>
        <label class="field-block">
          <span>打款流水号</span>
          <input v-model="paidReferenceNo" type="text" placeholder="请输入线下打款参考号" />
        </label>
        <label class="field-block">
          <span>打款备注</span>
          <textarea v-model="paidRemark" rows="4" placeholder="可填写到账时间、渠道说明等"></textarea>
        </label>
      </div>

      <template #footer>
        <div class="dialog-actions">
          <button class="ghost-button" type="button" @click="paidOpen = false">取消</button>
          <button type="button" :disabled="paidLoading" @click="submitPaid">
            {{ paidLoading ? '提交中...' : '确认已打款' }}
          </button>
        </div>
      </template>
    </BaseModal>

    <ConfirmDialog
      :open="cancelOpen"
      title="取消提现申请"
      description="取消后该笔申请会失效，并从打款流程中移除。确认继续吗？"
      confirm-text="确认取消"
      :loading="cancelLoading"
      @close="cancelOpen = false"
      @confirm="submitCancel"
    />
  </PageSection>
</template>
