package config

import "testing"

func boolPtr(v bool) *bool { return &v }
func intPtr(v int) *int    { return &v }

func TestMergeBool_ExplicitFalse(t *testing.T) {
	dst := true
	mergeBool(&dst, boolPtr(false), nil, "")
	if dst != false {
		t.Error("mergeBool should override with explicit false")
	}
}

func TestMergeBool_Nil(t *testing.T) {
	dst := true
	mergeBool(&dst, nil, nil, "")
	if dst != true {
		t.Error("mergeBool should not change dst when src is nil")
	}
}

func TestMergeBool_CLIFlagPrecedence(t *testing.T) {
	dst := true
	flags := map[string]bool{"test": true}
	mergeBool(&dst, boolPtr(false), flags, "test")
	if dst != true {
		t.Error("mergeBool should not override when CLI flag is set")
	}
}

func TestMergeBoolNoFlag_ExplicitFalse(t *testing.T) {
	dst := true
	mergeBoolNoFlag(&dst, boolPtr(false))
	if dst != false {
		t.Error("mergeBoolNoFlag should override with explicit false")
	}
}

func TestMergeBoolNoFlag_Nil(t *testing.T) {
	dst := true
	mergeBoolNoFlag(&dst, nil)
	if dst != true {
		t.Error("mergeBoolNoFlag should not change dst when src is nil")
	}
}

func TestMergeInt_ExplicitZero(t *testing.T) {
	dst := 365
	mergeInt(&dst, intPtr(0), nil, "")
	if dst != 0 {
		t.Error("mergeInt should override with explicit 0")
	}
}

func TestMergeInt_Nil(t *testing.T) {
	dst := 365
	mergeInt(&dst, nil, nil, "")
	if dst != 365 {
		t.Error("mergeInt should not change dst when src is nil")
	}
}

func TestMergeInt_CLIFlagPrecedence(t *testing.T) {
	dst := 365
	flags := map[string]bool{"test": true}
	mergeInt(&dst, intPtr(0), flags, "test")
	if dst != 365 {
		t.Error("mergeInt should not override when CLI flag is set")
	}
}

func TestMergeIntNoFlag_ExplicitZero(t *testing.T) {
	dst := 365
	mergeIntNoFlag(&dst, intPtr(0))
	if dst != 0 {
		t.Error("mergeIntNoFlag should override with explicit 0")
	}
}

func TestMergeIntNoFlag_Nil(t *testing.T) {
	dst := 365
	mergeIntNoFlag(&dst, nil)
	if dst != 365 {
		t.Error("mergeIntNoFlag should not change dst when src is nil")
	}
}
