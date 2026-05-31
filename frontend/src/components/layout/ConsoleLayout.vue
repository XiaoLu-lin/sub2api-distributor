<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { distributorNavItems, operatorNavItems } from '../../config/navigation'
import { useSession } from '../../session/useSession'
import type { NavItem } from '../../types'

const route = useRoute()
const router = useRouter()
const session = useSession()
const mobileNavOpen = ref(false)

const navItems = computed<NavItem[]>(() =>
  session.state.portalRole === 'operator' ? operatorNavItems : distributorNavItems,
)

const currentTitle = computed(() => route.meta.title || '控制台')
const currentPortalLabel = computed(() => (session.state.portalRole === 'operator' ? '运营后台' : '分销商端'))
const pageDescription = computed(() => {
  if (route.path === '/portal/dashboard') return '总览返利、提现进度与近期动态'
  if (route.path === '/portal/invitees') return '查看邀请用户，并快速复制邀请码和注册链接'
  if (route.path === '/portal/rebates') return '按流水查看返利入账明细'
  if (route.path === '/portal/withdrawals') return '发起提现申请，跟踪打款状态'
  if (route.path === '/portal/settlement') return '维护运营线下打款所需的收款信息'
  if (route.path === '/ops/distributors') return '查看已开通分销身份的用户及其资料'
  if (route.path === '/ops/withdrawals') return '处理打款中申请，并记录打款结果'
  return currentPortalLabel.value
})

function isActive(path: string): boolean {
  return route.path === path
}

async function handleLogout(): Promise<void> {
  await session.logout()
  router.push('/login')
}

function closeMobileNav(): void {
  mobileNavOpen.value = false
}
</script>

<template>
  <div class="console-shell">
    <aside class="console-sidebar" :class="{ 'console-sidebar-open': mobileNavOpen }">
      <div class="console-sidebar-head">
        <div>
          <p class="sidebar-kicker">Sub2API</p>
          <h2>Distributor</h2>
        </div>
        <button class="icon-button mobile-only" type="button" @click="closeMobileNav">×</button>
      </div>

      <p class="sidebar-role">{{ currentPortalLabel }}</p>

      <nav class="console-nav">
        <RouterLink
          v-for="item in navItems"
          :key="item.key"
          :to="item.path"
          class="console-nav-link"
          :class="{ 'console-nav-link-active': isActive(item.path) }"
          @click="closeMobileNav"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="console-sidebar-foot">
        <p>{{ session.state.email || '未登录用户' }}</p>
        <span>{{ currentPortalLabel }}</span>
      </div>
    </aside>

    <div v-if="mobileNavOpen" class="console-overlay" @click="closeMobileNav"></div>

    <section class="console-main">
      <header class="console-header">
        <div class="console-header-left">
          <button class="icon-button mobile-only" type="button" @click="mobileNavOpen = true">☰</button>
          <div>
            <p class="page-kicker">{{ currentPortalLabel }}</p>
            <h1>{{ currentTitle }}</h1>
            <p class="page-subcopy">{{ pageDescription }}</p>
          </div>
        </div>

        <div class="console-header-right">
          <div class="console-user-chip">
            <strong>{{ session.state.email || '未登录用户' }}</strong>
            <span>{{ currentPortalLabel }}</span>
          </div>
          <button class="secondary-button" type="button" @click="handleLogout">退出登录</button>
        </div>
      </header>

      <div class="console-page">
        <RouterView />
      </div>
    </section>
  </div>
</template>
