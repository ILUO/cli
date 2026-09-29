# task
> skill: lark-task

## +get-my-tasks
Use when the user requests tasks explicitly assigned to them.

### Avoid when
- The request does not specify an assigned-only relationship, or includes tasks created or followed by the current user → use [[+get-related-tasks]]
- The request gives only a task name without an assigned-only scope → use [[+search]]

### Skills
- lark-task/references/lark-task-get-my-tasks.md

## +get-related-tasks
Use for current-user task lists when the request does not specify a task relationship, or includes tasks the user created or follows.

### Skills
- lark-task/references/lark-task-get-related-tasks.md

## +search
Use when the user gives a task name or keyword without restricting the search to assigned tasks.

### Skills
- lark-task/references/lark-task-search.md
