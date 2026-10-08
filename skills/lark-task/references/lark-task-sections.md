# Task 分组查询

> 开始前读取 [共享规则](../../lark-shared/SKILL.md)，确认认证和身份。按主 Skill 的发现流程运行 `lark-cli task --help` 和 `lark-cli task sections --help`，再读取对应 method 的 schema。

## “我的任务”中的分组

代表当前用户查询其“我的任务”分组，使用 `my_tasks`，无需清单 GUID 或 `resource_id`：

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type my_tasks --as user --page-all --page-limit 0
```

## 某个清单中的分组

已知清单 GUID 时直接使用；只有清单名称时，先按主 Skill 的清单定位流程取得 GUID。分组名称用于匹配分组，不用作清单搜索词。

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type tasklist --resource-id "<tasklist_guid>" --as user --page-all --page-limit 0
```

两种查询均从返回的 `items` 中按 `name` 定位分组，取其 `guid` 作为 `section_guid`。同名候选展示 GUID 和已知归属，由用户确定目标；完整列取后仍找不到时，说明该归属下未找到目标分组，不自动扫描其他清单。

## 分组中的任务

用户要求查询组内任务时，将选定分组的 GUID 传入 `--section-guid`。已知分组 GUID 时，可直接从这一步开始，无需再次查询分组或清单：

```bash
lark-cli schema task.sections.tasks
lark-cli task sections tasks --section-guid "<section_guid>" --as user --page-all --page-limit 0
```

上述示例使用 `--page-all --page-limit 0` 完整列取，避免默认页数上限造成遗漏；若查询中断或返回部分失败，说明尚未完成，不能断言分组或任务不存在。
