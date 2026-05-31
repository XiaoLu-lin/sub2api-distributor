<template>
  <main class="auth-layout">
    <section class="auth-showcase">
      <p class="brand-chip">Sub2API Distributor</p>
      <h1>分销商结算后台</h1>
      <p class="auth-copy">
        统一查看返利、提现申请与线下打款进度。登录后会按你的身份自动进入分销商端或运营端。
      </p>
      <div class="auth-note-list">
        <article>
          <strong>分销商端</strong>
          <span>查看返利概览、邀请明细、提现状态与收款信息。</span>
        </article>
        <article>
          <strong>运营端</strong>
          <span>集中处理打款中申请，查看详情并手动标记已打款。</span>
        </article>
      </div>
    </section>

    <section class="auth-panel">
      <slot />
    </section>
  </main>
</template>
