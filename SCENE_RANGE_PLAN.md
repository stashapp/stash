# stashapp #3530 — 单文件多场景（Scene Range）实现方案

> 目标：$450 赏金（OpenCollective，支持 PayPal 提现）
> 认领：2026-08-06 10:00 UTC 已在 issue #3530 评论认领
> 仓库：stashapp/stash (Go 后端 + React 前端, 12.8k star)
> Fork：benkwok1983-cmd/stash → 本地 E:\Hermes\bounties\stash

## 需求
一个视频文件可以包含多个场景（按时间范围切分），无需重新编码。
维护者倾向"类似 Markers 的机制"（issue #260/#779 相关）。
prodxdack 2026-08-04 提出完整保守方案，我认可并采纳。

## 方案（prodxdack 设计 + 我的确认）

### 1. 数据模型
`scenes_files` 表（关联场景-文件）加两个可空列：
- `start_time` REAL NULL — 场景在文件中的开始时间（秒）
- `end_time` REAL NULL — 结束时间（秒）

范围属于"场景-文件关系"而非文件本身 → 多个场景可引用同一文件的不同范围。
- 双 NULL：完整文件（现有行为，向后兼容）
- 仅 start_time：从该点播到文件尾
- 双值：有界范围
- 范围仅对场景的 primary file 有效

### 2. 校验
`0 <= start_time < end_time <= 文件时长`
无范围的场景保持现有行为。

### 3. GraphQL/API
- Scene 暴露 `start_time`/`end_time`/派生 `duration`
- duration 从文件时长 + 范围计算，不单独存储
- 新增显式操作：`sceneCreateFromFileRange(file_id, start_time, end_time, title)`
- UI 动作："Create scene from range"

### 4. 播放
- 有范围场景：bounded MP4/WebM 转码，FFmpeg 传 start + duration
- 播放时间相对场景范围
- HLS/DASH 初始不启用（除非可 range-aware）
- 原文件不复制、不重编码

### 5. UI
- 场景编辑面板加 start/end 字段
- 文件信息视图加 "Create scene from range" 动作
- Markers 保持场景相对，播放器加载媒体时换算为源文件偏移

### 6. 测试
迁移、非法/边界范围、旧场景兼容、多场景共享一文件、duration 及过滤、范围播放、markers、循环、范围创建

### 初始范围（不做）
自动场景检测、高级时间线编辑、range-aware HLS/DASH

## 代码改动点（已探明）

### A. 数据库迁移
- 新建 `pkg/sqlite/migrations/86_scene_file_range.up.sql`（当前最高 85）
  ```sql
  ALTER TABLE `scenes_files` ADD COLUMN `start_time` REAL NULL;
  ALTER TABLE `scenes_files` ADD COLUMN `end_time` REAL NULL;
  ```
- 更新 `pkg/sqlite/database.go` 的 `appSchemaVersion` 到 86

### B. SQLite 层 — pkg/sqlite/table.go（relatedFilesTable，行 779-874）
- `insertJoin`（812）：加 start/end 参数（可选）
- `insertJoins`（824）：透传
- 新增 `getRange(ctx, sceneID)` 读 start/end
- 新增 `setRange(ctx, sceneID, start, end)` 更新
- `replaceJoins`（834）：范围由后续 setRange 设置

### C. 模型层 — pkg/models/
- `model_scene.go` Scene struct：加 `StartTime *float64` / `EndTime *float64`（transient，不持久化到 scenes 表）
- 场景的 file 关联读取时填充 range
- `SceneFileType`（行 271）：Duration 改为派生（文件时长 - 范围）

### D. API 层 — internal/api/
- `resolver_model_scene.go`：Scene 加 `start_time`/`end_time`/`duration` resolver
- `resolver_mutation_scene.go`：新增 `SceneCreateFromFileRange` mutation
- GraphQL schema：`graphql/schema/types/scene.graphql`
  - SceneFileType 或 Scene 加 start_time/end_time
  - Mutation 加 sceneCreateFromFileRange
  - gqlgen 重新生成（go generate ./...）

### E. 前端 — ui/
- Scene 编辑面板：start/end 输入
- 文件信息视图："Create scene from range"
- Player：范围播放支持（video 元素 time range）

## 里程碑
1. [ ] Go 环境装好，项目编译通过
2. [ ] 迁移文件 + appSchemaVersion
3. [ ] SQLite 层 range 读写 + 单元测试
4. [ ] 模型 + GraphQL schema + gqlgen 生成
5. [ ] API resolver + mutation
6. [ ] 后端测试通过
7. [ ] 前端 UI
8. [ ] 前端测试/构建
9. [ ] PR 提交 + 维护者 review
