## Why

`sub2api-distributor` 当前已经能展示分销商返利、提现和运营打款流转，但最近一轮 code review 证明这套系统在“邀请链路正确性”和“结算数据完整性”上还有几处会直接影响业务验收的问题：

- 分销商页面展示的邀请链接指向的是分销系统自己的 `/register`，而不是主系统真实注册入口，导致复制出去的链接无法真正完成主系统注册绑定。
- 运营开通分销商时只写入 `distributor_profiles`，没有确保主系统 `user_affiliates` 档案存在，历史用户或外部导入用户可能出现“可以登录分销后台，但拿不到邀请码”的半开通状态。
- 分销系统金额链路当前使用 `float64` 承载货币字段，后续在累计返利、已打款、打款中和可提现金额相互扣减时存在精度风险。
- 演示 seed 脚本会错误地把被邀请人也开成分销商，API 验收脚本默认账号也已经和当前最小演示数据脱节，容易让联调和验收得到错误结论。

如果不修复这些问题，系统虽然“可以跑起来”，但邀请入口、分销商开通和结算口径都无法稳定支撑真实业务验收。

## What Changes

- 修正分销商端邀请入口展示逻辑，让邀请码链接明确指向主系统的真实注册地址。
- 补齐“开通分销商”链路，在启用分销身份时自动确保主系统 `user_affiliates` 档案存在。
- 将分销系统内部结算金额模型从 `float64` 调整为更安全的 decimal 表达，并保持前后端展示稳定。
- 修正演示 seed 脚本和 API 验收脚本，确保默认演示数据、默认账号和真实业务角色一致。
- 为上述修复补充测试与文档，确保邀请链路、分销商开通和金额口径可重复验证。

## Capabilities

### New Capabilities
- `distributor-integrity-baseline`: 定义分销商邀请入口、分销身份开通和结算金额口径必须满足的完整性要求。

### Modified Capabilities
- None.

## Impact

- Affected backend:
  - `backend/internal/distributor/**`
  - `backend/internal/httpapi/**`
  - `backend/internal/app/**`
- Affected frontend:
  - `frontend/src/views/portal/**`
  - `frontend/src/views/ops/**`
  - `frontend/src/api/**`
  - `frontend/src/types.ts`
- Affected scripts and docs:
  - `scripts/seed_demo_data.sh`
  - `scripts/api_acceptance.sh`
  - `README.md`
  - `TESTING.md`
  - `docs/distributor-system-overview.md`
