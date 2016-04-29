package main

import "testing"

func TestLoadConfigEmpty(t *testing.T) {
	c, err := loadConfig("test_data/not_exist.ini")
	if err != nil {
		t.Fatalf("loadConfig(not_exist) should be succeeded: %s", err)
	}
	if c.Source != "" {
		t.Error("c.Source should be empty: %q", c.Source)
	}
	if c.TargetDir != "" {
		t.Error("c.TargetDir should be empty: %q", c.TargetDir)
	}
}

func TestLoadConfigSource(t *testing.T) {
	c, err := loadConfig("test_data/source_only.ini")
	if err != nil {
		t.Fatalf("loadConfig(source_only) should be succeeded: %s", err)
	}
	if c.Source != "foo" {
		t.Error("c.Source should be \"foo\": %q", c.Source)
	}
	if c.TargetDir != "" {
		t.Error("c.TargetDir should be empty: %q", c.TargetDir)
	}
}

func TestLoadConfigTargetDir(t *testing.T) {
	c, err := loadConfig("test_data/target_dir_only.ini")
	if err != nil {
		t.Fatalf("loadConfig(target_dir_only) should be succeeded: %s", err)
	}
	if c.Source != "" {
		t.Error("c.Source should be empty: %q", c.Source)
	}
	if c.TargetDir != "bar" {
		t.Error("c.TargetDir should be \"bar\": %q", c.TargetDir)
	}
}

func TestLoadConfigAll(t *testing.T) {
	c, err := loadConfig("test_data/all.ini")
	if err != nil {
		t.Fatalf("loadConfig(target_dir_only) should be succeeded: %s", err)
	}
	if c.Source != "foo" {
		t.Error("c.Source should be \"foo\": %q", c.Source)
	}
	if c.TargetDir != "bar" {
		t.Error("c.TargetDir should be \"bar\": %q", c.TargetDir)
	}
}
