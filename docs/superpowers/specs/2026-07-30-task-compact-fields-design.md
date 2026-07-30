# Task Compact Fields Design

## Goal

Expose the task title, members, start time, due time, and completion state when
task shortcuts already receive a Task entity but currently discard those
fields while building compact output.

The change must remain backward compatible: existing JSON fields, paths, and
types stay unchanged, and the implementation must not add API requests.

## Compatibility Rules

- Do not remove, rename, or change the type of an existing output field.
- Add a field only when the command does not already expose the same meaning.
- Reuse Task API names and values for newly projected entity fields:
  `summary`, `members`, `start`, `due`, and `status`.
- Preserve existing compact aliases such as `due_at`, `completed`, and
  `completed_at`; do not emit a second field for the same meaning in that
  command.
- Preserve current fallback behavior. In particular, a `+search` item whose
  detail request fails continues to contain only `guid` and `url`.
- Empty API values are transcribed faithfully. No synthetic member, time, or
  status value is invented.

## Affected Outputs

### Resource-reference write results

When their existing API response contains `data.task`, extend the current
`guid`/`url` projection with `summary`, `members`, `start`, `due`, and `status`:

- `task +create`
- `task +reopen`
- `task +assign`
- `task +followers`
- `task +reminder`

`task +complete` already returns `status`; add the missing `summary`, `members`,
`start`, and `due` fields without changing `completed_at` or
`already_completed`.

### Nested write results

Extend each successful Task item using the same projection:

- `task +update`: `tasks[]`
- `task +tasklist-create`: `created_tasks[]`
- `task +tasklist-task-add`: `successful_tasks[]`

Batch failure item schemas remain unchanged.

### Read summaries

Add only missing meanings to avoid redundant fields:

| Command | Existing equivalents kept | Newly projected fields |
| --- | --- | --- |
| `task +search` | `summary`, `due_at`, `completed_at` | `members`, `start`, `status` |
| `task +get-my-tasks` | `summary`, `due_at`, `completed` | `members`, `start` |
| `task +get-related-tasks` | `summary`, `members`, `status` | `start`, `due` |

The existing command-specific formatted time fields remain unchanged.

## Excluded Outputs

- Generated raw API commands already return the full upstream entity.
- `task +set-ancestor` is excluded because its endpoint does not provide a
  usable Task entity; fetching one would introduce an extra request and new
  post-mutation failure semantics.
- Comment, attachment, and tasklist outputs represent different entity types.
- Tasklist-level fields in `task +tasklist-create` remain unchanged; only its
  `created_tasks[]` items are Task entities.

## Data Flow

At the API boundary, existing shortcut code receives loose Task maps. A shared
projection helper copies the requested upstream fields into the existing
compact result map. Member data is preserved as returned, including `role`, so
callers can identify responsible users with `role == "assignee"` without losing
followers.

Every affected command uses Task data it already has:

- create/update/complete/reopen/member/tasklist mutations use their existing
  Task response;
- reminder uses its existing pre-mutation Task GET because reminder mutations
  do not change the projected fields;
- list-related and my-tasks use the returned list items;
- search keeps its existing per-hit detail query.

No new GET or other API call is introduced.

## Error and Fallback Behavior

Projection is best-effort and cannot turn a successful API operation into an
error. Missing upstream fields are omitted rather than guessed. Existing typed
API errors, partial-failure envelopes, exit codes, and per-item failure objects
remain unchanged.

## Testing

- Add focused helper tests that pin faithful projection and omission of absent
  fields.
- Extend command tests to assert newly exposed fields at the final JSON paths.
- Cover root write results, nested batch results, and all three read-summary
  variants.
- Keep existing fallback tests green, especially search detail failure and
  partial batch failures.
- Run task package tests, `gofmt`, and the repository-required unit test gate
  before completion.
