# task sections

> **Prerequisites:** Read the [shared rules](../../lark-shared/SKILL.md) for authentication, identity, and operation safety.

## Scope and operation guide

This reference organizes Task `sections` workflows by Catalog method. It currently covers locating sections and listing their tasks. See the [main Skill](../SKILL.md) for group terminology and ownership routing. The current CLI Catalog/schema is authoritative for parameters, permissions, and response fields.

| User goal | Read |
|----------|----------|
| List sections in a resource or locate a section by name | [sections.list](#sectionslist) |
| List tasks in a selected section | [sections.tasks](#sectionstasks) |

## Shared rules

1. Follow the main Skill's discovery flow: run `lark-cli task --help` and `lark-cli task sections --help`, confirm the method, and read its schema before executing a native command.
2. Use `--as user` when querying the current user's personal "My Tasks". For other operations, confirm identity and permissions through the shared rules and the method's help/schema.
3. For queries that support pagination, use `--page-all --page-limit 0` when a complete listing is required. This avoids the default page limit. If the query is interrupted or partially fails, report that results are incomplete; do not conclude that a section or task does not exist.

## sections.list

### When to use

List sections in "My Tasks" or a specific tasklist, or resolve a section name to its GUID.

### Prerequisites

- "My Tasks": use `my_tasks`; no tasklist GUID or `resource_id` is needed.
- A specific tasklist: use `tasklist` and provide its GUID. Use a known GUID directly. If only the tasklist name is available, follow the main Skill's tasklist lookup flow first.
- Resolve unclear ownership from context; ask for the owning resource if it remains unclear. Match a section name against sections, rather than using it as a tasklist search term.

### Examples

**List sections in "My Tasks"**

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type my_tasks --as user --page-all --page-limit 0
```

**List sections in a specific tasklist**

```bash
lark-cli schema task.sections.list
lark-cli task sections list --resource-type tasklist --resource-id "<tasklist_guid>" --as user --page-all --page-limit 0
```

### Result handling

Match `name` in the returned `items`, and use the selected item's `guid` as `section_guid`. If names are duplicated, show candidate GUIDs and their known ownership so the user can select the target. If a complete listing has no match, report that the section was not found in that resource; do not automatically scan other tasklists.

## sections.tasks

### When to use

The user wants to list tasks in a section.

### Prerequisites

The target `section_guid` is required. Use a known GUID directly. If only the section name is available, locate it with [sections.list](#sectionslist) first. A known section GUID does not require another section or tasklist lookup.

### Examples

```bash
lark-cli schema task.sections.tasks
lark-cli task sections tasks --section-guid "<section_guid>" --as user --page-all --page-limit 0
```

### Result handling

Present the returned tasks as requested. Report a complete result only after pagination finishes; otherwise, state that the results are incomplete.

## Adding operation guidance

For each added operation, create a `## sections.<method>` chapter using its Catalog method identifier and update the operation guide. Use the same subsection order: **When to use → Prerequisites → Examples → Result handling**. Keep shared rules at the start of this reference and operation-specific differences in the relevant chapter. Examples must show schema discovery before execution; leave complete parameter, permission, and response definitions to the schema.
