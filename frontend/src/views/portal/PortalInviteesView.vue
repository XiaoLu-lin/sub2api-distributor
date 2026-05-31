<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getInvitees, getInviteMeta } from '../../api'
import PageSection from '../../components/common/PageSection.vue'
import DataTable from '../../components/common/DataTable.vue'
import { useSession } from '../../session/useSession'
import type { InviteMeta, InviteeItem } from '../../types'
import { buildInviteLink, formatCurrency, formatDateTime } from '../../utils/format'

const session = useSession()
const invitees = ref<InviteeItem[]>([])
const inviteMeta = ref<InviteMeta | null>(null)
const loadingError = ref('')
const copyMessage = ref('')
const mainAppBaseUrl = ((import.meta.env.VITE_MAIN_APP_BASE_URL as string | undefined) ?? '').trim()

const inviteLink = computed(() => {
  return buildInviteLink(inviteMeta.value?.aff_code, mainAppBaseUrl)
})

const inviteLinkMissingConfig = computed(() => Boolean(inviteMeta.value?.aff_code) && !inviteLink.value)

async function loadInvitees(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    const [metaResp, inviteeResp] = await Promise.all([
      getInviteMeta(session.state.token),
      getInvitees(session.state.token),
    ])
    inviteMeta.value = metaResp
    invitees.value = inviteeResp
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载邀请用户失败'
  }
}

async function copyText(value: string, successText: string): Promise<void> {
  if (!value) {
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    copyMessage.value = successText
    window.setTimeout(() => {
      copyMessage.value = ''
    }, 1800)
  } catch {
    copyMessage.value = '复制失败，请手动复制'
  }
}

onMounted(loadInvitees)
</script>

<template>
  <div class="page-stack">
    <PageSection
      title="邀请入口"
      description="分销商可以在这里快速拿到自己的邀请码和邀请注册链接。"
      compact
    >
      <div class="invite-entry-card">
        <div class="invite-entry-main">
          <div class="invite-entry-block">
            <span class="invite-entry-label">邀请码</span>
            <strong class="invite-entry-code">{{ inviteMeta?.aff_code || '-' }}</strong>
          </div>
          <div class="invite-entry-block">
            <span class="invite-entry-label">邀请链接</span>
            <p class="invite-entry-link">
              {{ inviteLink || '未配置主系统注册地址' }}
            </p>
          </div>
        </div>

        <div class="invite-entry-actions">
          <button class="secondary-button" type="button" @click="copyText(inviteMeta?.aff_code || '', '邀请码已复制')">
            复制邀请码
          </button>
          <button type="button" @click="copyText(inviteLink, '邀请链接已复制')">复制邀请链接</button>
        </div>
      </div>

      <p v-if="copyMessage" class="success-text">{{ copyMessage }}</p>
      <p v-if="inviteLinkMissingConfig" class="error-text">
        当前还没有配置主系统注册地址，暂时不能生成可用邀请链接，请先配置 `VITE_MAIN_APP_BASE_URL`。
      </p>
      <p class="info-text">新用户注册时带上 `aff` 参数或填写邀请码，即可绑定到该分销商名下。</p>
    </PageSection>

    <PageSection title="邀请用户" description="查看通过你的邀请进入系统的用户及其累计返利贡献。" compact>
    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>
    <DataTable
      :items="invitees"
      :columns="[
        { key: 'user_id', title: '用户 ID', width: '120px' },
        { key: 'email', title: '邮箱', width: '260px' },
        { key: 'username', title: '用户名', width: '180px' },
        { key: 'total_rebate', title: '累计返利', width: '140px', align: 'right' },
        { key: 'created_at', title: '邀请时间', width: '180px' },
      ]"
      :row-key="(item) => item.user_id"
      empty-title="暂无邀请用户"
      empty-description="当有用户通过你的邀请进入系统后，这里会出现明细。"
    >
      <template #username="{ item }">
        {{ item.username || '-' }}
      </template>
      <template #total_rebate="{ item }">
        {{ formatCurrency(item.total_rebate) }}
      </template>
        <template #created_at="{ item }">
          {{ formatDateTime(item.created_at) }}
        </template>
      </DataTable>
    </PageSection>
  </div>
</template>
