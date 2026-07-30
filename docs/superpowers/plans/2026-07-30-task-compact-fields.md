# Task Compact Fields Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Preserve existing task shortcut output contracts while projecting task title, members, start, due, and completion state from Task entities the commands already receive.

**Architecture:** Add one typed-field projection helper at the task output boundary, then call it from existing compact-output builders. Each command selects only missing semantic fields, so existing aliases such as `due_at`, `completed`, and `completed_at` remain unchanged and no redundant fields or new API calls are introduced.

**Tech Stack:** Go, Cobra shortcut runtime, `internal/httpmock`, standard `testing`, `make unit-test`

---

## File Map

- Create `shortcuts/task/task_output.go`: typed task field identifiers and faithful map-to-map projection helper.
- Create `shortcuts/task/task_output_test.go`: helper contract plus root write-output command contracts.
- Modify `shortcuts/task/task_query_helpers.go`: add missing fields to search and related-task projections.
- Modify `shortcuts/task/task_query_helpers_test.go`: pin search/related projection fields and absent-field behavior.
- Modify `shortcuts/task/task_get_my_tasks.go`: add `members` and `start` to my-task items.
- Modify `shortcuts/task/task_get_my_tasks_test.go`: pin final JSON paths.
- Modify `shortcuts/task/task_get_related_tasks_test.go`: pin `start` and `due` in command output.
- Modify `shortcuts/task/shortcuts.go`, `task_reopen.go`, `task_assign.go`, `task_followers.go`, `task_reminder.go`, and `task_complete.go`: project fields already present in each Task response.
- Modify the corresponding command tests to pin root output fields.
- Modify `shortcuts/task/task_update.go`, `tasklist_create.go`, and `tasklist_add_task.go`: project fields into successful nested Task items.
- Modify their tests to pin nested output paths while preserving failure schemas.

### Task 1: Shared faithful projection helper

**Files:**
- Create: `shortcuts/task/task_output.go`
- Create: `shortcuts/task/task_output_test.go`

- [ ] **Step 1: Write the failing helper tests**

Add tests that call the not-yet-defined helper with a full source entity and assert exact preservation, then call it with absent fields and assert those keys are omitted:

```go
func TestProjectTaskFields(t *testing.T) {
	task := map[string]interface{}{
		"summary": "Ship compact fields",
		"members": []interface{}{map[string]interface{}{
			"id": "ou_owner", "name": "Owner", "type": "user", "role": "assignee",
		}},
		"start":  map[string]interface{}{"timestamp": "1000", "is_all_day": false},
		"due":    map[string]interface{}{"timestamp": "2000", "is_all_day": false},
		"status": "todo",
	}
	out := map[string]interface{}{"guid": "task-1"}

	projectTaskFields(out, task, standardTaskOutputFields...)

	for _, field := range standardTaskOutputFields {
		key := string(field)
		if !reflect.DeepEqual(out[key], task[key]) {
			t.Fatalf("%s = %#v, want %#v", key, out[key], task[key])
		}
	}
}

func TestProjectTaskFieldsOmitsAbsentFields(t *testing.T) {
	out := map[string]interface{}{}
	projectTaskFields(out, map[string]interface{}{"summary": "Only title"}, standardTaskOutputFields...)
	if _, ok := out["members"]; ok {
		t.Fatalf("members unexpectedly projected: %#v", out)
	}
}
```

- [ ] **Step 2: Run the tests and verify RED**

Run: `go test ./shortcuts/task -run '^TestProjectTaskFields'`

Expected: build failure because `projectTaskFields` and `standardTaskOutputFields` do not exist.

- [ ] **Step 3: Add the minimal typed-field helper**

```go
package task

type taskOutputField string

const (
	taskOutputSummary taskOutputField = "summary"
	taskOutputMembers taskOutputField = "members"
	taskOutputStart   taskOutputField = "start"
	taskOutputDue     taskOutputField = "due"
	taskOutputStatus  taskOutputField = "status"
)

var standardTaskOutputFields = []taskOutputField{
	taskOutputSummary,
	taskOutputMembers,
	taskOutputStart,
	taskOutputDue,
	taskOutputStatus,
}

func projectTaskFields(dst, task map[string]interface{}, fields ...taskOutputField) {
	for _, field := range fields {
		key := string(field)
		if value, ok := task[key]; ok {
			dst[key] = value
		}
	}
}
```

- [ ] **Step 4: Run helper tests and verify GREEN**

Run: `go test ./shortcuts/task -run '^TestProjectTaskFields'`

Expected: PASS.

- [ ] **Step 5: Commit the helper**

```bash
git add shortcuts/task/task_output.go shortcuts/task/task_output_test.go
git commit -m "feat: add task output field projector"
```

### Task 2: Read-summary outputs

**Files:**
- Modify: `shortcuts/task/task_query_helpers.go:136-201`
- Modify: `shortcuts/task/task_query_helpers_test.go:68-225`
- Modify: `shortcuts/task/task_get_my_tasks.go:199-225`
- Modify: `shortcuts/task/task_get_my_tasks_test.go:13-101`
- Modify: `shortcuts/task/task_get_related_tasks_test.go:95-185`

- [ ] **Step 1: Extend read-output tests first**

Add `members`, `start`, `due`, and `status` fixtures to the direct helper tests. Assert:

```go
projected := outputTaskSummary(task)
if !reflect.DeepEqual(projected["members"], task["members"]) ||
	!reflect.DeepEqual(projected["start"], task["start"]) ||
	projected["status"] != "todo" {
	t.Fatalf("search projection lost compact fields: %#v", projected)
}
if _, duplicated := projected["due"]; duplicated {
	t.Fatalf("search projection duplicated due alongside due_at: %#v", projected)
}
```

For related tasks, assert `start` and `due` are present. In the JSON execution fixtures for `+get-my-tasks`, add `members` and `start` and assert their serialized paths while retaining `due_at` and `completed`. In `+get-related-tasks`, add `start` and `due` to the fixture and expected JSON fragments.

- [ ] **Step 2: Run focused read tests and verify RED**

Run:

```bash
go test ./shortcuts/task -run 'Test(OutputTaskSummary|OutputRelatedTaskAndTimeRangeFilter|GetMyTasks_LocalTimeFormatting|GetRelatedTasks_Execute)'
```

Expected: FAIL because the new fields are absent.

- [ ] **Step 3: Add only the missing projections**

After each existing output map is built:

```go
// outputTaskSummary: due_at/completed_at already carry those meanings.
projectTaskFields(out, task, taskOutputMembers, taskOutputStart, taskOutputStatus)

// outputRelatedTask: members/status already exist.
projectTaskFields(out, task, taskOutputStart, taskOutputDue)

// GetMyTasks outputItem: due_at/completed already exist.
projectTaskFields(outputItem, item, taskOutputMembers, taskOutputStart)
```

- [ ] **Step 4: Run focused tests and verify GREEN**

Run the command from Step 2.

Expected: PASS.

- [ ] **Step 5: Commit read-summary behavior**

```bash
git add shortcuts/task/task_query_helpers.go shortcuts/task/task_query_helpers_test.go shortcuts/task/task_get_my_tasks.go shortcuts/task/task_get_my_tasks_test.go shortcuts/task/task_get_related_tasks_test.go
git commit -m "feat: expose missing fields in task summaries"
```

### Task 3: Root write-result outputs

**Files:**
- Modify: `shortcuts/task/shortcuts.go:234-250`
- Modify: `shortcuts/task/task_reopen.go:42-59`
- Modify: `shortcuts/task/task_assign.go:84-95`
- Modify: `shortcuts/task/task_followers.go:85-96`
- Modify: `shortcuts/task/task_reminder.go:75-173`
- Modify: `shortcuts/task/task_complete.go:86-103`
- Test: `shortcuts/task/task_output_test.go`
- Test: `shortcuts/task/task_complete_test.go`
- Test: `shortcuts/task/task_reminder_test.go`

- [ ] **Step 1: Write root output contract tests**

Use `taskShortcutTestFactory`, `warmTenantToken`, and `httpmock.Stub` to run each shortcut with a Task response containing the five upstream fields. Decode the JSON envelope and assert existing `guid`/`url` plus the new fields. Use a table for create, reopen, assign, and followers; keep complete and reminder in their existing specialized test files.

The common assertion must be exact:

```go
func assertStandardTaskFields(t *testing.T, data map[string]interface{}) {
	t.Helper()
	for _, key := range []string{"summary", "members", "start", "due", "status"} {
		if _, ok := data[key]; !ok {
			t.Errorf("data missing %q: %#v", key, data)
		}
	}
}
```

For `+complete`, assert `status` remains the existing derived value and that `summary`, `members`, `start`, and `due` are added. For `+reminder --remove` with no reminders, assert the early-success output also uses the Task data already fetched.

- [ ] **Step 2: Run root-output tests and verify RED**

Run:

```bash
go test ./shortcuts/task -run 'Test(TaskRootOutputFields|CompleteTask|ReminderTask)'
```

Expected: FAIL on missing projected fields.

- [ ] **Step 3: Project fields without changing identifiers or requests**

For create/reopen/assign/followers/reminder, keep the existing output map and append:

```go
projectTaskFields(outData, task, standardTaskOutputFields...)
```

For reminder, use the already fetched `taskObj`, including the no-existing-reminder early return. For complete, do not overwrite its derived `status`:

```go
projectTaskFields(outData, task, taskOutputSummary, taskOutputMembers, taskOutputStart, taskOutputDue)
```

- [ ] **Step 4: Run root-output tests and verify GREEN**

Run the command from Step 2.

Expected: PASS with no additional registered HTTP stubs.

- [ ] **Step 5: Commit root write outputs**

```bash
git add shortcuts/task/shortcuts.go shortcuts/task/task_reopen.go shortcuts/task/task_assign.go shortcuts/task/task_followers.go shortcuts/task/task_reminder.go shortcuts/task/task_complete.go shortcuts/task/task_output_test.go shortcuts/task/task_complete_test.go shortcuts/task/task_reminder_test.go
git commit -m "feat: expand task write result fields"
```

### Task 4: Nested successful Task items

**Files:**
- Modify: `shortcuts/task/task_update.go:88-110`
- Modify: `shortcuts/task/task_update_test.go:85-165`
- Modify: `shortcuts/task/tasklist_create.go:122-142`
- Modify: `shortcuts/task/tasklist_create_test.go:15-90`
- Modify: `shortcuts/task/tasklist_add_task.go:89-106`
- Modify: `shortcuts/task/tasklist_add_task_test.go:13-48`

- [ ] **Step 1: Extend nested-item tests first**

Populate successful Task fixtures with all five upstream fields. After decoding output, call `assertStandardTaskFields` for `tasks[]`, `created_tasks[]`, and `successful_tasks[]`. Retain assertions for `confirmed`, partial-failure exit codes, and failure item shapes.

- [ ] **Step 2: Run nested-output tests and verify RED**

Run:

```bash
go test ./shortcuts/task -run 'Test(TaskUpdateNormalizesAllIDsAndReturnsConfirmedFields|CreateTasklist_PartialFailure|AddTaskToTasklist_Success)'
```

Expected: FAIL because successful nested items contain only identifiers and command-specific fields.

- [ ] **Step 3: Add projection to each successful item**

Build each existing item first, then project standard fields:

```go
item := map[string]interface{}{
	"guid": guid,
	"url":  urlVal,
}
projectTaskFields(item, task, standardTaskOutputFields...)
successful = append(successful, item)
```

For update, retain `confirmed` in the item before projecting. Do not touch failed item construction.

- [ ] **Step 4: Run nested-output tests and verify GREEN**

Run the command from Step 2.

Expected: PASS; partial-failure tests still report unchanged failed item schemas.

- [ ] **Step 5: Commit nested result behavior**

```bash
git add shortcuts/task/task_update.go shortcuts/task/task_update_test.go shortcuts/task/tasklist_create.go shortcuts/task/tasklist_create_test.go shortcuts/task/tasklist_add_task.go shortcuts/task/tasklist_add_task_test.go
git commit -m "feat: expose fields in nested task results"
```

### Task 5: Formatting, regression, and repository gates

**Files:**
- Verify all modified Go files and module metadata.

- [ ] **Step 1: Format modified Go files**

Run: `gofmt -w shortcuts/task`

Expected: command exits 0.

- [ ] **Step 2: Run task package tests**

Run: `go test ./shortcuts/task`

Expected: PASS.

- [ ] **Step 3: Run required unit-test gate**

Run: `make unit-test`

Expected: PASS with zero failed packages.

- [ ] **Step 4: Run static and repository cleanliness checks**

Run:

```bash
go vet ./...
gofmt -l .
git diff --check
git status --short
```

Expected: `go vet` exits 0, `gofmt -l .` prints nothing, `git diff --check` exits 0, and status contains only intentional task changes.

- [ ] **Step 5: Commit any formatting-only changes if present**

```bash
git add shortcuts/task
git commit -m "chore: format task compact field changes"
```

Skip this commit when formatting produced no additional diff.
