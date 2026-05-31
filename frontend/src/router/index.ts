import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import type { PortalRole, RouteMetaConfig } from '../types'
import { useSession } from '../session/useSession'

declare module 'vue-router' {
  interface RouteMeta extends RouteMetaConfig {}
}

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/auth/LoginView.vue'),
    meta: {
      title: '登录',
      requiresAuth: false,
    },
  },
  {
    path: '/portal',
    component: () => import('../components/layout/ConsoleLayout.vue'),
    meta: {
      title: '分销商后台',
      requiresAuth: true,
      role: 'distributor',
    },
    children: [
      {
        path: '',
        redirect: '/portal/dashboard',
      },
      {
        path: 'dashboard',
        name: 'portal-dashboard',
        component: () => import('../views/portal/PortalDashboardView.vue'),
        meta: {
          title: '概览',
          requiresAuth: true,
          role: 'distributor',
        },
      },
      {
        path: 'invitees',
        name: 'portal-invitees',
        component: () => import('../views/portal/PortalInviteesView.vue'),
        meta: {
          title: '邀请用户',
          requiresAuth: true,
          role: 'distributor',
        },
      },
      {
        path: 'rebates',
        name: 'portal-rebates',
        component: () => import('../views/portal/PortalRebatesView.vue'),
        meta: {
          title: '返利明细',
          requiresAuth: true,
          role: 'distributor',
        },
      },
      {
        path: 'withdrawals',
        name: 'portal-withdrawals',
        component: () => import('../views/portal/PortalWithdrawalsView.vue'),
        meta: {
          title: '提现申请',
          requiresAuth: true,
          role: 'distributor',
        },
      },
      {
        path: 'settlement',
        name: 'portal-settlement',
        component: () => import('../views/portal/PortalSettlementView.vue'),
        meta: {
          title: '收款信息',
          requiresAuth: true,
          role: 'distributor',
        },
      },
    ],
  },
  {
    path: '/ops',
    component: () => import('../components/layout/ConsoleLayout.vue'),
    meta: {
      title: '运营后台',
      requiresAuth: true,
      role: 'operator',
    },
    children: [
      {
        path: '',
        redirect: '/ops/withdrawals',
      },
      {
        path: 'distributors',
        name: 'ops-distributors',
        component: () => import('../views/ops/OpsDistributorsView.vue'),
        meta: {
          title: '分销商管理',
          requiresAuth: true,
          role: 'operator',
        },
      },
      {
        path: 'withdrawals',
        name: 'ops-withdrawals',
        component: () => import('../views/ops/OpsWithdrawalsView.vue'),
        meta: {
          title: '提现管理',
          requiresAuth: true,
          role: 'operator',
        },
      },
    ],
  },
]

function resolveDefaultPath(role: PortalRole | ''): string {
  if (role === 'operator') {
    return '/ops/withdrawals'
  }
  return '/portal/dashboard'
}

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const session = useSession()
  await session.hydrate()

  if (to.path === '/') {
    return resolveDefaultPath(session.state.portalRole)
  }

  if (!to.meta.requiresAuth && session.isAuthenticated.value && to.path === '/login') {
    return resolveDefaultPath(session.state.portalRole)
  }

  if (to.meta.requiresAuth && !session.isAuthenticated.value) {
    return '/login'
  }

  if (to.meta.role && to.meta.role !== session.state.portalRole) {
    return resolveDefaultPath(session.state.portalRole)
  }

  return true
})

router.afterEach((to) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} | Sub2API Distributor`
  }
})
