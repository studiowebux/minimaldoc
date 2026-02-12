package config

import (
	"reflect"
)

// mergeString sets dst to src if src is non-empty and the CLI flag was not set
func mergeString(dst *string, src string, cliFlags map[string]bool, flagName string) {
	if flagName != "" && cliFlags[flagName] {
		return
	}
	if src != "" {
		*dst = src
	}
}

// mergeStringNoFlag sets dst to src if src is non-empty (no CLI flag check)
func mergeStringNoFlag(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}

// mergeBool sets dst to src if src is true and the CLI flag was not set
// Note: Cannot detect explicit false due to Go zero-value semantics
func mergeBool(dst *bool, src bool, cliFlags map[string]bool, flagName string) {
	if flagName != "" && cliFlags[flagName] {
		return
	}
	if src {
		*dst = src
	}
}

// mergeBoolNoFlag sets dst to src if src is true (no CLI flag check)
func mergeBoolNoFlag(dst *bool, src bool) {
	if src {
		*dst = src
	}
}

// mergeInt sets dst to src if src > 0 and the CLI flag was not set
func mergeInt(dst *int, src int, cliFlags map[string]bool, flagName string) {
	if flagName != "" && cliFlags[flagName] {
		return
	}
	if src > 0 {
		*dst = src
	}
}

// mergeIntNoFlag sets dst to src if src > 0 (no CLI flag check)
func mergeIntNoFlag(dst *int, src int) {
	if src > 0 {
		*dst = src
	}
}

// mergeStringSlice sets dst to src if src has elements and the CLI flag was not set
func mergeStringSlice(dst *[]string, src []string, cliFlags map[string]bool, flagName string) {
	if flagName != "" && cliFlags[flagName] {
		return
	}
	if len(src) > 0 {
		*dst = src
	}
}

// mergeStringSliceNoFlag sets dst to src if src has elements (no CLI flag check)
func mergeStringSliceNoFlag(dst *[]string, src []string) {
	if len(src) > 0 {
		*dst = src
	}
}

// mergeStringFields uses reflection to merge all non-empty string fields from src to dst.
// Both dst and src must be pointers to structs with matching string fields.
// This eliminates repetitive if statements for structs with many string fields.
func mergeStringFields(dst, src any) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)

	// Handle pointer to struct
	if srcVal.Kind() == reflect.Pointer {
		srcVal = srcVal.Elem()
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Field(i)

		// Only handle string fields
		if srcField.Kind() != reflect.String {
			continue
		}

		// Only merge if source is non-empty
		if srcField.String() != "" && dstField.CanSet() {
			dstField.SetString(srcField.String())
		}
	}
}
