# task +subscribe-event

> **Prerequisites:** Please read `../lark-shared/SKILL.md` to understand authentication, global parameters, and security rules.
>
> **⚠️ Note:** This API supports both `user` and `bot` identities. Use `user` to subscribe the current user's accessible tasks; use `bot` to subscribe tasks the **application is responsible for**.

Manually create the task update event subscription for the current identity.

For listening to task events, prefer the new unified event consumer:

```bash
lark-cli event consume task.task.update_user_access_v2 --as user
lark-cli event consume task.task.update_user_access_v2 --as bot
```

`event consume task.task.update_user_access_v2` calls this same subscription API during startup, then streams events as NDJSON through the shared event bus.

This shortcut is still useful when you only want to register the subscription and do not want to start a long-running consumer. It is different from `event consume`:
- `task +subscribe-event` registers task-event access for the **current identity**
- `event consume task.task.update_user_access_v2` registers task-event access, starts the local bus daemon if needed, and streams matching events to stdout
- with `--as user`, it subscribes the **current user** to task events for tasks they created, are responsible for, or follow
- with `--as bot`, it subscribes using the **application identity** for tasks the application is responsible for

The task event type is:

```text
task.task.update_user_access_v2
```

Within this event, task changes are represented by commit types (string values). Deduped list:

```text
task_assignees_update
task_completed_update
task_create
task_deleted
task_desc_update
task_followers_update
task_reminders_update
task_start_due_update
task_summary_update
```

`event consume` output uses the standard Lark V2 envelope. Event payload shape (example):

```json
{
  "schema": "2.0",
  "header": {
    "event_id": "evt_xxx",
    "event_type": "task.task.update_user_access_v2",
    "create_time": "1775793266152"
  },
  "event": {
    "event_types": ["task_summary_update"],
    "task_guid": "task_guid_xxx"
  }
}
```

- `.header.event_type`: event type, should be `task.task.update_user_access_v2`
- `.header.event_id`: unique event id (useful for dedup)
- `.header.create_time`: event timestamp (ms)
- `.event.event_types`: list of commit types (see the deduped list above)
- `.event.task_guid`: the task GUID that changed

In practice, this means:
- with `--as user`, the subscribed user can receive updates for tasks visible to them through authorship, assignment, or following
- with `--as bot`, the subscription covers tasks the application is responsible for

To actually receive events, use the standard event consumer:

```bash
lark-cli event consume task.task.update_user_access_v2 --as user
```

Useful projections:

```bash
lark-cli event consume task.task.update_user_access_v2 --as user \
  --jq '{event_id: .header.event_id, task_guid: .event.task_guid, event_types: .event.event_types, timestamp: .header.create_time}'
```

## Recommended Commands

```bash
# Preferred: subscribe and listen in one command
lark-cli event consume task.task.update_user_access_v2 --as user

# Manual subscription only
lark-cli task +subscribe-event --as user

# Manual subscription with app identity
lark-cli task +subscribe-event --as bot
```


## Parameters

This shortcut has no additional parameters.

## Workflow

1. Confirm whether the user wants to subscribe with `user` identity or `bot` identity.
2. If the user wants to listen now, execute `lark-cli event consume task.task.update_user_access_v2 --as <identity>`.
3. If the user only wants to register the subscription, execute `lark-cli task +subscribe-event --as <identity>`.
4. Report whether the subscription or consumer startup succeeded, and clarify which identity it applies to.

> [!CAUTION]
> This is a **Write Operation** -- You must confirm the user's intent before executing.
