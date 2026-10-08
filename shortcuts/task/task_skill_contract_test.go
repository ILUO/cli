// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package task

import (
	"os"
	"strings"
	"testing"
)

func TestTaskSkillRoutesGroupsBeforeTasklistSearch(t *testing.T) {
	raw, err := os.ReadFile("../../skills/lark-task/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	start := strings.Index(content, "**任务分组定位（必读）**")
	end := strings.Index(content, "**任务搜索技巧**")
	if start < 0 || end <= start {
		t.Fatal("task group routing must precede task and tasklist search guidance")
	}
	section := content[start:end]
	for _, required := range []string{
		"分组 / groups", "`sections`", "`my_tasks`", "`tasklist`",
		"`--as user`", "`--resource-id`", "`section_guid`",
		"references/lark-task-sections.md",
	} {
		if !strings.Contains(section, required) {
			t.Errorf("task group routing missing %q", required)
		}
	}
}

func TestTaskSectionsReferenceDiscoversSchemaBeforeExamples(t *testing.T) {
	raw, err := os.ReadFile("../../skills/lark-task/references/lark-task-sections.md")
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	for _, example := range []struct {
		name    string
		schema  string
		command string
	}{
		{
			name:    "my_tasks",
			schema:  "lark-cli schema task.sections.list",
			command: "lark-cli task sections list --resource-type my_tasks --as user --page-all --page-limit 0",
		},
		{
			name:    "tasklist",
			schema:  "lark-cli schema task.sections.list",
			command: `lark-cli task sections list --resource-type tasklist --resource-id "<tasklist_guid>" --as user --page-all --page-limit 0`,
		},
		{
			name:    "section_tasks",
			schema:  "lark-cli schema task.sections.tasks",
			command: `lark-cli task sections tasks --section-guid "<section_guid>" --as user --page-all --page-limit 0`,
		},
	} {
		t.Run(example.name, func(t *testing.T) {
			commandIndex := strings.Index(content, example.command)
			if commandIndex < 0 {
				t.Fatalf("section reference missing example: %s", example.command)
			}
			blockStart := strings.LastIndex(content[:commandIndex], "```bash\n")
			if blockStart < 0 {
				t.Fatal("section example must be in a bash code block")
			}
			beforeCommand := content[blockStart+len("```bash\n") : commandIndex]
			if strings.Contains(beforeCommand, "```") || !strings.Contains(beforeCommand, example.schema) {
				t.Errorf("section example must follow its schema discovery: %s", example.command)
			}
		})
	}
}
