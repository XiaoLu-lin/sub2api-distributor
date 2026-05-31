<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getOpsDistributors, lookupOpsUsers, updateOpsDistributorProfile } from '../../api'
import BaseModal from '../../components/common/BaseModal.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import DataTable from '../../components/common/DataTable.vue'
import PageSection from '../../components/common/PageSection.vue'
import StatusBadge from '../../components/common/StatusBadge.vue'
import ToolbarActions from '../../components/common/ToolbarActions.vue'
import { useSession } from '../../session/useSession'
import type { DistributorProfile, UserLookupItem } from '../../types'

const session = useSession()

const distributors = ref<DistributorProfile[]>([])
const loadingError = ref('')
const successMessage = ref('')

const searchKeyword = ref('')
const searchLoading = ref(false)
const searchResults = ref<UserLookupItem[]>([])
const selectedUser = ref<UserLookupItem | null>(null)

const createOpen = ref(false)
const saving = ref(false)
const toggleLoading = ref(false)
const toggleOpen = ref(false)
const currentDistributor = ref<DistributorProfile | null>(null)

const form = ref<DistributorProfile>({
  user_id: 0,
  status: 'active',
  display_name: '',
  settlement_channel: 'alipay',
  settlement_account_name: '',
  settlement_account_no: '',
  settlement_account_extra: '{}',
  notes: '',
})

const toggleTargetStatus = computed(() => (currentDistributor.value?.status === 'active' ? 'disabled' : 'active'))

function resetForm(): void {
  form.value = {
    user_id: 0,
    status: 'active',
    display_name: '',
    settlement_channel: 'alipay',
    settlement_account_name: '',
    settlement_account_no: '',
    settlement_account_extra: '{}',
    notes: '',
  }
}

async function loadDistributors(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    distributors.value = await getOpsDistributors(session.state.token)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载分销商失败'
  }
}

function openCreateModal(): void {
  createOpen.value = true
  searchKeyword.value = ''
  searchResults.value = []
  selectedUser.value = null
  successMessage.value = ''
  resetForm()
}

function selectUser(user: UserLookupItem): void {
  selectedUser.value = user
  form.value.user_id = user.id
  form.value.display_name = user.username || user.email
  form.value.notes = `由运营开通分销商：${user.email}`
}

async function handleLookup(): Promise<void> {
  if (!session.state.token) {
    return
  }

  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    searchResults.value = []
    return
  }

  searchLoading.value = true
  loadingError.value = ''
  try {
    searchResults.value = await lookupOpsUsers(session.state.token, keyword)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '搜索用户失败'
  } finally {
    searchLoading.value = false
  }
}

async function submitCreate(): Promise<void> {
  if (!session.state.token || !selectedUser.value) {
    loadingError.value = '请先选择一个用户'
    return
  }

  saving.value = true
  loadingError.value = ''
  successMessage.value = ''

  try {
    await updateOpsDistributorProfile(session.state.token, selectedUser.value.id, {
      ...form.value,
      user_id: selectedUser.value.id,
      status: 'active',
    })
    createOpen.value = false
    successMessage.value = `已开通分销商：${selectedUser.value.email}`
    await loadDistributors()
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '开通分销商失败'
  } finally {
    saving.value = false
  }
}

function openToggleDialog(item: DistributorProfile): void {
  currentDistributor.value = item
  toggleOpen.value = true
}

async function submitToggle(): Promise<void> {
  if (!session.state.token || !currentDistributor.value) {
    return
  }

  toggleLoading.value = true
  loadingError.value = ''
  successMessage.value = ''
  try {
    await updateOpsDistributorProfile(session.state.token, currentDistributor.value.user_id, {
      ...currentDistributor.value,
      status: toggleTargetStatus.value,
    })
    toggleOpen.value = false
    successMessage.value =
      toggleTargetStatus.value === 'active'
        ? `已启用分销商：${currentDistributor.value.display_name || currentDistributor.value.user_id}`
        : `已停用分销商：${currentDistributor.value.display_name || currentDistributor.value.user_id}`
    await loadDistributors()
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '更新分销商状态失败'
  } finally {
    toggleLoading.value = false
  }
}

onMounted(loadDistributors)
</script>

<template>
  <PageSection title="分销商管理" description="集中查看、开通和停用分销身份，并维护其收款资料。">
    <template #actions>
      <ToolbarActions>
        <button type="button" @click="openCreateModal">开通分销商</button>
      </ToolbarActions>
    </template>

    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>
    <p v-if="successMessage" class="success-text">{{ successMessage }}</p>

    <DataTable
      :items="distributors"
      :columns="[
        { key: 'user_id', title: '用户 ID', width: '100px' },
        { key: 'display_name', title: '显示名称', width: '180px' },
        { key: 'status', title: '状态', width: '100px' },
        { key: 'settlement_channel', title: '收款方式', width: '120px' },
        { key: 'settlement_account_name', title: '收款户名', width: '160px' },
        { key: 'settlement_account_no', title: '收款账号', width: '240px' },
        { key: 'notes', title: '备注', width: '220px' },
        { key: 'actions', title: '操作', width: '140px', align: 'center' },
      ]"
      :row-key="(item) => item.user_id"
      empty-title="暂无分销商"
      empty-description="开通分销商后，这里会展示其资料与当前状态。"
    >
      <template #display_name="{ item }">
        {{ item.display_name || '-' }}
      </template>
      <template #status="{ item }">
        <StatusBadge :status="item.status" />
      </template>
      <template #settlement_account_name="{ item }">
        {{ item.settlement_account_name || '-' }}
      </template>
      <template #settlement_account_no="{ item }">
        {{ item.settlement_account_no || '-' }}
      </template>
      <template #notes="{ item }">
        {{ item.notes || '-' }}
      </template>
      <template #actions="{ item }">
        <div class="table-actions">
          <button class="link-button" type="button" @click="openToggleDialog(item)">
            {{ item.status === 'active' ? '停用' : '启用' }}
          </button>
        </div>
      </template>
    </DataTable>

    <BaseModal :open="createOpen" title="开通分销商" width="wide" @close="createOpen = false">
      <div class="detail-stack">
        <div class="lookup-panel">
          <div class="lookup-toolbar">
            <label class="field-block lookup-field">
              <span>搜索主系统用户</span>
              <input
                v-model="searchKeyword"
                type="text"
                placeholder="输入用户 ID、邮箱或用户名"
                @keydown.enter.prevent="handleLookup"
              />
            </label>
            <button type="button" :disabled="searchLoading" @click="handleLookup">
              {{ searchLoading ? '搜索中...' : '搜索用户' }}
            </button>
          </div>

          <div v-if="searchResults.length" class="lookup-results">
            <button
              v-for="user in searchResults"
              :key="user.id"
              type="button"
              class="lookup-result-item"
              :class="{ 'lookup-result-item-active': selectedUser?.id === user.id }"
              @click="selectUser(user)"
            >
              <div>
                <strong>{{ user.email }}</strong>
                <p>ID {{ user.id }} · {{ user.username || '未设置用户名' }}</p>
              </div>
              <div class="lookup-result-meta">
                <StatusBadge :status="user.status" />
                <span v-if="user.is_distributor" class="table-muted">已是分销商</span>
              </div>
            </button>
          </div>
          <p v-else class="info-text">先搜索用户，再选择要开通分销身份的账号。</p>
        </div>

        <div v-if="selectedUser" class="form-card form-card-embedded">
          <div class="detail-grid">
            <div class="detail-card">
              <span>已选用户</span>
              <strong>{{ selectedUser.email }}</strong>
            </div>
            <div class="detail-card">
              <span>主系统状态</span>
              <StatusBadge :status="selectedUser.status" />
            </div>
          </div>

          <div class="form-grid">
            <label class="field-block">
              <span>显示名称</span>
              <input v-model="form.display_name" type="text" placeholder="例如：华东区域分销商" />
            </label>

            <label class="field-block">
              <span>收款渠道</span>
              <select v-model="form.settlement_channel">
                <option value="bank">银行卡</option>
                <option value="alipay">支付宝</option>
                <option value="wechat">微信</option>
                <option value="usdt">USDT</option>
                <option value="manual">人工协商</option>
              </select>
            </label>

            <label class="field-block">
              <span>收款户名</span>
              <input v-model="form.settlement_account_name" type="text" placeholder="请输入收款户名" />
            </label>

            <label class="field-block">
              <span>收款账号</span>
              <input v-model="form.settlement_account_no" type="text" placeholder="请输入账号或收款地址" />
            </label>
          </div>

          <label class="field-block">
            <span>补充信息</span>
            <textarea
              v-model="form.settlement_account_extra"
              rows="4"
              placeholder="可填写 JSON 或文字说明，例如链类型、开户地址、备注信息"
            ></textarea>
          </label>

          <label class="field-block">
            <span>备注</span>
            <textarea v-model="form.notes" rows="3" placeholder="给运营留的额外说明"></textarea>
          </label>
        </div>
      </div>

      <template #footer>
        <div class="dialog-actions">
          <button class="ghost-button" type="button" @click="createOpen = false">取消</button>
          <button type="button" :disabled="saving || !selectedUser" @click="submitCreate">
            {{ saving ? '开通中...' : '确认开通' }}
          </button>
        </div>
      </template>
    </BaseModal>

    <ConfirmDialog
      :open="toggleOpen"
      :title="toggleTargetStatus === 'active' ? '启用分销商' : '停用分销商'"
      :description="
        toggleTargetStatus === 'active'
          ? '启用后，该用户可以登录分销商后台查看返利和提现。'
          : '停用后，该用户将不能再登录分销商后台。'
      "
      :confirm-text="toggleTargetStatus === 'active' ? '确认启用' : '确认停用'"
      :loading="toggleLoading"
      @close="toggleOpen = false"
      @confirm="submitToggle"
    />
  </PageSection>
</template>
