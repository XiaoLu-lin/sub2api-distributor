<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login, me } from '../../api'
import AuthLayout from '../../components/layout/AuthLayout.vue'
import { useSession } from '../../session/useSession'

const router = useRouter()
const session = useSession()

const email = ref(session.state.email)
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

async function handleLogin(): Promise<void> {
  loading.value = true
  errorMessage.value = ''

  try {
    const result = await login(email.value, password.value)
    session.setSession(result.token, result.portal_role, email.value)
    session.setMeProfile(await me(result.token))
    await router.push(result.portal_role === 'operator' ? '/ops/withdrawals' : '/portal/dashboard')
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout>
    <div class="login-card">
      <div class="login-card-header">
        <p class="panel-kicker">统一登录</p>
        <h2>进入分销结算后台</h2>
        <p>复用主系统账号登录，系统会按角色自动切换到对应工作台。</p>
      </div>

      <form class="form-stack" @submit.prevent="handleLogin">
        <label class="field-block">
          <span>邮箱</span>
          <input v-model="email" type="email" autocomplete="username" placeholder="请输入账号邮箱" />
        </label>

        <label class="field-block">
          <span>密码</span>
          <input v-model="password" type="password" autocomplete="current-password" placeholder="请输入登录密码" />
        </label>

        <button type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>

        <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
      </form>
    </div>
  </AuthLayout>
</template>
