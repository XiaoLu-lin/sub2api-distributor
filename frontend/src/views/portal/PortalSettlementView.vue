<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getProfile, updateProfile } from '../../api'
import PageSection from '../../components/common/PageSection.vue'
import { useSession } from '../../session/useSession'
import type { DistributorProfile } from '../../types'

const session = useSession()

const profile = ref<DistributorProfile | null>(null)
const loadingError = ref('')
const saving = ref(false)
const successMessage = ref('')

async function loadProfile(): Promise<void> {
  if (!session.state.token) {
    return
  }

  loadingError.value = ''
  try {
    profile.value = await getProfile(session.state.token)
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '加载收款信息失败'
  }
}

async function submitProfile(): Promise<void> {
  if (!session.state.token || !profile.value) {
    return
  }

  saving.value = true
  successMessage.value = ''

  try {
    await updateProfile(session.state.token, profile.value)
    successMessage.value = '收款信息已保存'
  } catch (error) {
    loadingError.value = error instanceof Error ? error.message : '保存收款信息失败'
  } finally {
    saving.value = false
  }
}

onMounted(loadProfile)
</script>

<template>
  <PageSection title="收款信息" description="这里保存运营线下打款所需的收款渠道和账户信息。" compact>
    <p v-if="loadingError" class="error-text">{{ loadingError }}</p>
    <p v-if="successMessage" class="success-text">{{ successMessage }}</p>

    <div v-if="profile" class="form-card">
      <div class="form-grid">
        <label class="field-block">
          <span>展示名称</span>
          <input v-model="profile.display_name" type="text" placeholder="例如：华东区域分销商" />
        </label>

        <label class="field-block">
          <span>收款渠道</span>
          <select v-model="profile.settlement_channel">
            <option value="bank">银行卡</option>
            <option value="alipay">支付宝</option>
            <option value="wechat">微信</option>
            <option value="usdt">USDT</option>
            <option value="manual">人工协商</option>
          </select>
        </label>

        <label class="field-block">
          <span>收款户名</span>
          <input v-model="profile.settlement_account_name" type="text" placeholder="请输入收款户名" />
        </label>

        <label class="field-block">
          <span>收款账号</span>
          <input v-model="profile.settlement_account_no" type="text" placeholder="请输入账号或收款地址" />
        </label>
      </div>

      <label class="field-block">
        <span>补充信息</span>
        <textarea
          v-model="profile.settlement_account_extra"
          rows="4"
          placeholder="可填写开户地址、链类型、开户地址备注等 JSON 或文字说明"
        ></textarea>
      </label>

      <label class="field-block">
        <span>备注</span>
        <textarea v-model="profile.notes" rows="3" placeholder="给运营留的额外说明"></textarea>
      </label>

      <div class="form-actions">
        <button type="button" :disabled="saving" @click="submitProfile">
          {{ saving ? '保存中...' : '保存收款信息' }}
        </button>
      </div>
    </div>
  </PageSection>
</template>
