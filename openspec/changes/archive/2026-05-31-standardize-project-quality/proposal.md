## Why

`sub2api-distributor` 已经具备可运行的 MVP 能力，但项目层面的质量基线还不完整：文档存在端口和账号信息漂移，前端缺少自动化测试，后端测试覆盖偏薄，后端代码注释不足，导致后续交接、扩展和排错成本偏高。现在补齐这些基础设施，可以把项目从“能演示”推进到“可维护、可验证、可交付”。

## What Changes

- 统一根文档、测试文档和系统说明文档，修正端口、账号、启动方式和验收说明的漂移信息。
- 新增并完善 OpenSpec 变更工件，明确项目质量基线、验证范围和实现任务。
- 为前端补充轻量自动化测试基线，覆盖关键工具函数、会话存储和基础请求行为。
- 为后端补充更多单元测试，覆盖认证、配置和已有核心逻辑的关键行为。
- 为后端导出类型、导出函数和关键业务逻辑补齐说明性注释，降低阅读和维护门槛。
- 更新 API 验收脚本，使其覆盖当前实际已实现的接口能力。

## Capabilities

### New Capabilities
- `project-quality-baseline`: 定义该项目在文档一致性、自动化验证、验收覆盖和后端代码可读性方面的基线要求。

### Modified Capabilities
- None.

## Impact

- Affected code:
  - `backend/internal/**`
  - `frontend/src/**`
  - `frontend/package.json`
  - `scripts/api_acceptance.sh`
- Affected docs:
  - `README.md`
  - `TESTING.md`
  - `docs/distributor-system-overview.md`
  - `frontend/README.md`
- Affected systems:
  - OpenSpec change artifacts under `openspec/changes/standardize-project-quality/`
