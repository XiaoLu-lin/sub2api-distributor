## Context

这次变更来自完整 code review 之后的正式修复，而不是新增功能。目标不是改变“分销商申请提现 -> 运营线下打款 -> 运营标记已打款”这条业务主路径，而是修复当前实现里几个会导致错误业务认知或错误验收结果的缺口。

当前系统存在 4 类核心问题：

1. 邀请入口错误
   - 分销商端展示的邀请码存在，但“邀请链接”拼接成了分销系统自己的域名和 `/register` 路径。
   - 新系统本身没有注册页，因此链接外观正确但行为错误。

2. 分销身份开通不完整
   - 登录分销商端只要求 `distributor_profiles.status = active`。
   - 但邀请码、邀请用户和返利链路依赖主系统 `user_affiliates`。
   - 这导致部分用户会出现“能登录，但 invite-meta 404”的半开通状态。

3. 结算金额模型精度风险
   - Go 服务层当前把 `numeric(20,8)` 转成 `float64`。
   - 返利累计、打款中、已打款、转余额和可提现金额都在服务层继续用浮点做减法。
   - 这会在高频累加或边界金额场景下埋下精度误差。

4. 演示与验收脚本不一致
   - `seed_demo_data.sh` 把被邀请人也开成了分销商。
   - `api_acceptance.sh` 默认账号依然指向旧的演示用户。
   - 这会让团队误以为系统逻辑有问题，实际上是脚本和数据基线错了。

## Goals / Non-Goals

**Goals:**

- 让分销商端复制出去的邀请码链接可以真实进入主系统注册流程。
- 保证运营开通分销商后，该账号在新系统里立即具备完整的邀请码能力，而不是只具备登录能力。
- 让分销系统内部金额计算不再依赖浮点表达。
- 让 seed/acceptance 脚本与当前推荐演示数据一致。
- 为以上修复补充自动化验证和文档说明。

**Non-Goals:**

- 不重写现有提现状态机，不新增审核流。
- 不改主系统原有 affiliate 业务规则，不替换主系统的返利产生逻辑。
- 不新增自动支付或自动打款能力。
- 不扩展为新的注册页或营销页系统。

## Decisions

### 1. 邀请链接改为显式配置主系统注册入口

前端不再使用 `window.location.origin` 拼接注册链接，而是通过明确配置的主系统基准地址生成：

- `VITE_MAIN_APP_BASE_URL`
- 组合结果：`${VITE_MAIN_APP_BASE_URL}/register?aff=<code>`

这样可以避免：

- 分销系统域名被误当成主系统域名
- 本地 5177 / 5176 这类开发端口被复制给业务方
- 未来部署时因为不同域名结构再次踩坑

如果未配置该变量，则前端要降级成明确提示，而不是继续拼接一个假的 `/register`。

### 2. 开通分销商时自动确保主系统 affiliate 档案存在

在后端 `UpsertProfile` 这条链路里，新增“确保 affiliate 档案存在”的步骤：

- 在 `status = active` 的开通或启用路径上执行
- 若主系统 `user_affiliates` 已存在，直接复用
- 若不存在，则创建新的 `aff_code`

这里不直接依赖主系统 service 层，而是在当前独立系统内部新增一个小范围 repository/helper，复用主库表结构完成 `EnsureUserAffiliate` 语义，保持项目独立性。

这样做的好处是：

- 不需要把整个主系统 affiliate service 直接搬进来
- 新系统仍然保持独立部署
- 运营在分销后台“开通分销商”后，账号立即完整可用

### 3. 金额模型在后端切换为 decimal 字符串承载，再在边界层转展示值

这轮不建议简单“继续 float64 并四舍五入”，而是直接把服务层金额模型改成 decimal。

推荐实现：

- 后端内部使用 `github.com/shopspring/decimal`
- SQL 查询不再 cast 到 `double precision`
- 优先 scan 为字符串/数值后转 decimal
- `Summary`、`WithdrawalRequest`、`RebateRecord` 等对前端的 JSON 字段继续输出标准十进制数值字符串可解析形态

考虑到前端当前已经按 number 消费，本轮可以采取折中方式：

- 后端内部全部 decimal
- 对 JSON 输出时格式化成保留 2 位或按原精度可解析的 JSON number / string
- 如果切到 string 会扩大前端改动面，则先保持 JSON number，但只能在最终输出边界做一次受控转换

推荐落地为：

- 内部 decimal
- 对前端响应仍保持 number，减少协议改动
- 所有减法/累加/比较都在 decimal 层完成

### 4. 脚本基线改成“一个分销商，一个被邀请人”

统一演示数据策略：

- 分销商：有 `distributor_profiles`
- 被邀请人：没有 `distributor_profiles`
- 两者都有 `user_affiliates`
- 返利流水只记在分销商账户上

`api_acceptance.sh` 默认账号也切到这组数据，避免默认脚本和当前推荐演示路径分裂。

## Risks / Trade-offs

- [引入 decimal 依赖]：会增加少量实现复杂度，但这是结算系统应承担的复杂度，收益大于成本。
- [前端新增主系统注册链接配置]：部署环境需要补一个变量，但这是显式正确依赖，比继续猜域名更安全。
- [独立系统内部实现 EnsureUserAffiliate 语义]：会复制一小段主系统逻辑，但范围可控，且能避免强耦合主系统 service。
- [脚本切换演示账号]：历史验收文档里的旧账号会失效，因此需要同步文档。

## Migration Plan

1. 先补 OpenSpec 工件，锁定修复范围。
2. 后端先补 affiliate 档案确保逻辑和 decimal 金额逻辑。
3. 前端修邀请链接配置和相关展示/错误提示。
4. 修正 seed 与 acceptance 脚本。
5. 补测试，覆盖 invite meta、profile enablement 和金额口径。
6. 运行后端测试、前端测试、前端构建和验收脚本。

## Open Questions

- 主系统最终线上注册路径是否固定为 `/register`。

当前默认按已有主系统注册入口约定处理；若后续发现线上使用的是不同路径，可通过前端配置继续调整，不阻塞本次修复。
