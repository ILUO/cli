// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package common

import (
	"encoding/json"
	"testing"

	"github.com/larksuite/cli/internal/apicatalog"
	"github.com/larksuite/cli/internal/meta"
)

func TestMethodFlagSchemaPrinterProjectsConfiguredInputSubtree(t *testing.T) {
	catalog := apicatalog.New(apicatalog.SourceEmbedded, []meta.Service{
		{
			Name: "task",
			Resources: map[string]meta.Resource{
				"tasks": {
					Methods: map[string]meta.Method{
						"patch": {
							RequestBody: map[string]meta.Field{
								"task": {
									Type: "object",
									Properties: map[string]meta.Field{
										"summary": {Type: "string"},
										"due": {
											Type: "object",
											Properties: map[string]meta.Field{
												"timestamp": {Type: "string"},
											},
										},
									},
								},
								"update_fields": {Type: "array", Required: true},
							},
						},
					},
				},
			},
		},
	})
	printer := methodFlagSchemaPrinter(
		func() apicatalog.Catalog { return catalog },
		"task.tasks.patch",
		map[string]string{"data": "data.task"},
	)

	raw, err := printer("data")
	if err != nil {
		t.Fatalf("printer(data) error = %v", err)
	}
	var schema struct {
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("decode schema: %v\n%s", err, raw)
	}
	if schema.Type != "object" || schema.Properties["summary"] == nil || schema.Properties["update_fields"] != nil {
		t.Fatalf("projected schema = %#v, want task object only", schema)
	}

	nested, err := printer("data.due.timestamp")
	if err != nil {
		t.Fatalf("printer(data.due.timestamp) error = %v", err)
	}
	var timestampSchema map[string]interface{}
	if err := json.Unmarshal(nested, &timestampSchema); err != nil {
		t.Fatalf("decode nested schema: %v\n%s", err, nested)
	}
	if timestampSchema["type"] != "string" {
		t.Fatalf("nested schema = %#v, want string", timestampSchema)
	}
}

func TestMethodFlagSchemaPrinterListsAndValidatesFlags(t *testing.T) {
	catalog := apicatalog.New(apicatalog.SourceEmbedded, []meta.Service{
		{Name: "demo", Resources: map[string]meta.Resource{
			"items": {Methods: map[string]meta.Method{
				"patch": {RequestBody: map[string]meta.Field{"item": {Type: "object"}}},
			}},
		}},
	})
	printer := methodFlagSchemaPrinter(
		func() apicatalog.Catalog { return catalog },
		"demo.items.patch",
		map[string]string{"data": "data.item"},
	)

	listed, err := printer("")
	if err != nil {
		t.Fatalf("printer(list) error = %v", err)
	}
	if string(listed) == "" {
		t.Fatal("printer(list) returned empty output")
	}
	if _, err := printer("unknown"); err == nil {
		t.Fatal("printer(unknown) error = nil, want validation error")
	}
}
