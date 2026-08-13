// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package common

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/apicatalog"
	"github.com/larksuite/cli/internal/registry"
	metaschema "github.com/larksuite/cli/internal/schema"
)

// MethodFlagSchemaPrinter projects shortcut composite-flag schemas from the
// canonical API method schema. Each binding maps a shortcut flag name to a
// dotted path inside the method input schema, for example data -> data.task.
func MethodFlagSchemaPrinter(methodPath string, bindings map[string]string) func(string) ([]byte, error) {
	return methodFlagSchemaPrinter(registry.SchemaCatalog, methodPath, bindings)
}

func methodFlagSchemaPrinter(
	catalog func() apicatalog.Catalog,
	methodPath string,
	bindings map[string]string,
) func(string) ([]byte, error) {
	ownedBindings := make(map[string]string, len(bindings))
	for flag, path := range bindings {
		ownedBindings[flag] = path
	}

	return func(flagName string) ([]byte, error) {
		available := make([]string, 0, len(ownedBindings))
		for flag := range ownedBindings {
			available = append(available, flag)
		}
		sort.Strings(available)

		flagName = strings.TrimSpace(flagName)
		if flagName == "" {
			return json.MarshalIndent(map[string]interface{}{
				"method":               methodPath,
				"introspectable_flags": available,
				"hint":                 "run again with --flag-name <name> to dump that flag's JSON Schema; append a dotted property path to inspect one nested field",
			}, "", "  ")
		}

		requested := strings.Split(flagName, ".")
		flag := strings.ReplaceAll(requested[0], "_", "-")
		boundPath, ok := ownedBindings[flag]
		if !ok {
			return nil, errs.NewValidationError(
				errs.SubtypeInvalidArgument,
				"no JSON Schema registered for --%s; available: %v",
				requested[0],
				available,
			).WithParam("--flag-name")
		}

		target, err := catalog().Resolve(apicatalog.ParsePath([]string{methodPath}))
		if err != nil || target.Kind != apicatalog.TargetMethod || target.Method == nil {
			validationErr := errs.NewValidationError(
				errs.SubtypeFailedPrecondition,
				"API schema source %s is unavailable",
				methodPath,
			).WithHint("run lark-cli schema " + methodPath + " to verify the installed API metadata")
			if err != nil {
				validationErr.WithCause(err)
			}
			return nil, validationErr
		}

		envelope := metaschema.EnvelopeOf(*target.Method)
		path := append(strings.Split(boundPath, "."), requested[1:]...)
		property, ok := inputSchemaProperty(envelope.InputSchema, path)
		if !ok {
			return nil, errs.NewValidationError(
				errs.SubtypeInvalidArgument,
				"JSON Schema path %s is unavailable in %s",
				strings.Join(path, "."),
				methodPath,
			).WithParam("--flag-name")
		}
		return json.MarshalIndent(property, "", "  ")
	}
}

func inputSchemaProperty(input *metaschema.InputSchema, path []string) (metaschema.Property, bool) {
	if input == nil || input.Properties == nil || len(path) == 0 {
		return metaschema.Property{}, false
	}
	property, ok := input.Properties.Map[path[0]]
	if !ok {
		return metaschema.Property{}, false
	}
	for _, segment := range path[1:] {
		property, ok = nestedSchemaProperty(property, segment)
		if !ok {
			return metaschema.Property{}, false
		}
	}
	return property, true
}

func nestedSchemaProperty(property metaschema.Property, segment string) (metaschema.Property, bool) {
	for depth := 0; depth < 8; depth++ {
		if property.Properties != nil {
			child, ok := property.Properties.Map[segment]
			return child, ok
		}
		if property.Items == nil {
			return metaschema.Property{}, false
		}
		property = *property.Items
	}
	return metaschema.Property{}, false
}
