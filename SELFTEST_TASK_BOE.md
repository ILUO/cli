# Task Shortcuts — BOE 自测报告（boe_task_tasklist_oapi_support）

- 执行时间：2026-04-10 19:59:37 +0800
- 环境切换：`. ./lark-env.sh boe --boe-env-name boe_task_tasklist_oapi_support`
- 代理提示：`HTTPS_PROXY=http://127.0.0.1:8899`

## 覆盖命令

- 读：
  - `task +get-related-tasks`
  - `task +search`
  - `task +tasklist-search`
  - `task +get-my-tasks`（cursor 自测）
- 写：
  - `task +subscribe-event`（user / bot）
  - `task +set-ancestor`
  - `task +create`（用于 set-ancestor 验证）

## 结果摘要

- user 路径：全部返回 `need_user_authorization (user: ou_...)`，当前 BOE 用户未完成授权，网络链路正常但鉴权阻塞。
- bot 路径：命令链路正常，见下述详情。

## 详细记录

### Subscribe Event

- 命令：`task +subscribe-event --as bot --format json`
- 输出：`{"ok": true, "identity": "bot", "data": {"ok": true}}`
- 结论：成功注册应用身份的任务事件订阅。

### Create Tasks（用于 set-ancestor）

- 命令：`task +create --as bot --summary "CLI_SELFTEST_PARENT_20260410_2000" --format json`
- 输出：返回 `guid` 与 `url`，示例 `guid=1af0da5e-a8ae-40fa-846e-b897f3b4ac02`
- 命令：`task +create --as bot --summary "CLI_SELFTEST_CHILD_20260410_2000" --format json`
- 输出：返回 `guid` 与 `url`，示例 `guid=c89d2d7d-9958-43c3-b583-06521022e5cd`

### Set Ancestor（设置与清空）

- 命令：`task +set-ancestor --as bot --task-id c89d2d7d-... --ancestor-id 1af0da5e-... --format json`
- 输出：`{"ok": true, "data": {"guid": "c89d2d7d-..."}}`
- 命令：`task +set-ancestor --as bot --task-id c89d2d7d-... --format json`（清空）
- 输出：`{"ok": true, "data": {"guid": "c89d2d7d-..."}}`
- 结论：设置/清空祖先链路正常，返回值与预期一致。

### Dry-Run（user 路径形状验证）

- `task +get-related-tasks --as user --page-token 1752730590582902 --created-by-me --dry-run`
  - 请求：`GET /open-apis/task/v2/task_v2/list_related_task`，带 `page_size=100`、`page_token=1752730590582902`、`user_id_type=open_id`
- `task +search --as user --query cli-selftest-nohit-20260410 --page-limit 1 --dry-run`
  - 请求：`POST /open-apis/task/v2/tasks/search`，`{"query":"cli-selftest-nohit-20260410"}`
- `task +tasklist-search --as user --query cli-selftest-nohit-20260410 --page-limit 1 --dry-run`
  - 请求：`POST /open-apis/task/v2/tasklists/search`，`{"query":"cli-selftest-nohit-20260410"}`
- `task +get-my-tasks --as user --page-token 1752730590582902 --page-limit 1 --dry-run`
  - 请求：`GET /open-apis/task/v2/tasks`，带 `type=my_tasks`、`page_token=1752730590582902` 等

### user 路径实跑（均被鉴权拦截）

- `task +get-related-tasks --as user --page-limit 1 --format json`
- `task +search --as user --query cli-selftest-nohit-20260410 --page-limit 1 --format json`
- `task +tasklist-search --as user --query cli-selftest-nohit-20260410 --page-limit 1 --format json`
- `task +get-my-tasks --as user --page-token 1752730590582902 --page-limit 1 --format json`
- `task +subscribe-event --as user --format json`
- 统一失败：`need_user_authorization (user: ou_...)`

## 结论与后续

- Bot 路径（subscribe-event、create、set-ancestor）功能验证通过。
- User 路径当前 BOE 需要先完成授权；建议用 `lark-cli auth login --as user --scope "task:task:read task:task:write task:tasklist:read task:tasklist:write"` 完成授权后再跑实测。
- 文档已澄清：
  - `+subscribe-event` 为身份级订阅注册，不需要 task GUID
  - `+get-related-tasks` 的 `has_more/page_token` 是上游游标；`--created-by-me` / `--followed-by-me` 为本地过滤

