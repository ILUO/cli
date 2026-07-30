// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package task

import (
	"reflect"
	"testing"
)

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
	if out["guid"] != "task-1" {
		t.Fatalf("existing guid changed: %#v", out)
	}
}

func TestProjectTaskFieldsOmitsAbsentFields(t *testing.T) {
	out := map[string]interface{}{}
	projectTaskFields(out, map[string]interface{}{"summary": "Only title"}, standardTaskOutputFields...)

	if out["summary"] != "Only title" {
		t.Fatalf("summary = %#v, want Only title", out["summary"])
	}
	for _, key := range []string{"members", "start", "due", "status"} {
		if _, ok := out[key]; ok {
			t.Fatalf("%s unexpectedly projected: %#v", key, out)
		}
	}
}
