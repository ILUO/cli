# task sections

> 开始前读取 [共享规则](../../lark-shared/SKILL.md)，确认认证、身份和操作安全规则。

## 适用范围与操作导航

本文按 method 组织 Task `sections` 资源的操作流程，当前提供分组定位和组内任务查询。分组的概念与归属判断见 [主 Skill](../SKILL.md)；接口参数、权限和返回结构以当前 CLI 的 Catalog/schema 为准。

| 用户目标 | 阅读章节 |
|----------|----------|
| 列取指定归属下的分组，或按名称定位分组 | [sections.list](#sectionslist) |
| 查询已选定分组中的任务 | [sections.tasks](#sectionstasks) |

## 共用规则

1. 按主 Skill 的发现流程运行 `lark-cli task --help` 和 `lark-cli task sections --help`，确认 method 后读取对应 schema，再执行原生命令。
2. 代表当前用户查询个人“我的任务”时，使用 `--as user`；其他操作的身份和权限要求按共享规则及对应 method 的 help/schema 确认。
3. 对支持分页的查询，完整列取时使用 `--page-all --page-limit 0`，避免默认页数上限造成遗漏。查询中断或返回部分失败时，说明结果尚不完整，不能断言分组或任务不存在。

## sections.list

### 适用场景

需要列取“我的任务”或指定清单下的分组，或按分组名称取得目标 GUID。

### 前置信息

- “我的任务”：使用 `my_tasks` 归属，无需清单 GUID 或 `resource_id`。
- 指定清单：使用 `tasklist` 归属，需要清单 GUID。已知 GUID 时直接使用；只有清单名称时，先按主 Skill 的清单定位流程取得 GUID。
- 归属不明时，先结合上下文判断，仍不明确再询问归属。分组名称用于匹配分组，不用作清单搜索词。

### 执行示例

**查询“我的任务”中的分组**

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type my_tasks --as user --page-all --page-limit 0
```

**查询指定清单中的分组**

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type tasklist --resource-id "<tasklist_guid>" --as user --page-all --page-limit 0
```

### 结果处理

从返回的 `items` 中按 `name` 定位分组，取其 `guid` 作为 `section_guid`。同名候选展示 GUID 和已知归属，由用户确定目标；完整列取后仍找不到时，说明该归属下未找到目标分组，不自动扫描其他清单。

## sections.tasks

### 适用场景

用户要求查询某个分组中的任务。

### 前置信息

需要目标 `section_guid`。已知分组 GUID 时直接使用；仅有名称时，先按 [sections.list](#sectionslist) 定位分组，无需为已知 GUID 再次查询分组或清单。

### 执行示例

```bash
lark-cli schema task.sections.tasks
lark-cli task sections tasks --section-guid "<section_guid>" --as user --page-all --page-limit 0
```

### 结果处理

按用户要求展示返回的组内任务；完成分页后再报告完整结果，查询未完成时明确说明结果不完整。

## 后续补充约定

新增 section 操作时，以 Catalog 的 method 标识增加 `## sections.<method>` 章节，并同步更新操作导航。每个操作沿用“适用场景 → 前置信息 → 执行示例 → 结果处理”的顺序；共用规则保留在本文前部，场景差异写在对应操作内。执行示例先展示 schema 发现，参数、权限和返回字段的完整定义继续由 schema 提供。
