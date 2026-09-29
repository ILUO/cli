// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package affordance

import (
	"os"
	"strings"
	"testing"

	"github.com/larksuite/cli/internal/apicatalog"
	"github.com/larksuite/cli/internal/meta"
)

func TestTaskQueryAffordanceExplainsUnspecifiedScope(t *testing.T) {
	previous := mdSource
	t.Cleanup(func() { SetSource(previous) })
	SetSource(os.DirFS("../../affordance"))

	parse := func(command string) meta.Affordance {
		t.Helper()
		raw, ok := For(apicatalog.Catalog{}, "task", command)
		if !ok {
			t.Fatalf("For(task, %s) returned no command guidance", command)
		}
		guidance, ok := (meta.Method{Affordance: raw}).ParsedAffordance()
		if !ok {
			t.Fatalf("Task %s guidance did not parse", command)
		}
		return guidance
	}

	assigned := parse("+get-my-tasks")
	if !taskGuidanceContains(assigned.UseWhen, "explicitly assigned") ||
		!taskGuidanceContains(assigned.AvoidWhen, "+get-related-tasks") {
		t.Fatalf("assigned-only guidance does not route unspecified scope to related tasks: %#v", assigned)
	}
	related := parse("+get-related-tasks")
	if !taskGuidanceContains(related.UseWhen, "does not specify a task relationship") {
		t.Fatalf("related-task guidance does not cover unspecified scope: %#v", related)
	}
	search := parse("+search")
	if !taskGuidanceContains(search.UseWhen, "task name or keyword") {
		t.Fatalf("search guidance does not cover name-only lookup: %#v", search)
	}
}

func taskGuidanceContains(items []string, text string) bool {
	for _, item := range items {
		if strings.Contains(item, text) {
			return true
		}
	}
	return false
}
