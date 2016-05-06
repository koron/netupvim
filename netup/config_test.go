package netup

import "testing"

func TestLoadConfigEmpty(t *testing.T) {
	c, err := loadConfig("test_data/not_exist.ini")
	if err != nil {
		t.Fatalf("loadConfig(not_exist) should be succeeded: %s", err)
	}
	if c.Source != "" {
		t.Errorf("c.Source should be empty: %q", c.Source)
	}
	if c.TargetDir != "" {
		t.Errorf("c.TargetDir should be empty: %q", c.TargetDir)
	}
	if c.CPU != "" {
		t.Errorf("c.CPU should be empty: %q", c.CPU)
	}
}

func TestLoadConfigSource(t *testing.T) {
	c, err := loadConfig("test_data/source_only.ini")
	if err != nil {
		t.Fatalf("loadConfig(source_only) should be succeeded: %s", err)
	}
	if c.Source != "foo" {
		t.Errorf("c.Source should be \"foo\": %q", c.Source)
	}
	if c.TargetDir != "" {
		t.Errorf("c.TargetDir should be empty: %q", c.TargetDir)
	}
	if c.CPU != "" {
		t.Errorf("c.CPU should be empty: %q", c.CPU)
	}
}

func TestLoadConfigTargetDir(t *testing.T) {
	c, err := loadConfig("test_data/target_dir_only.ini")
	if err != nil {
		t.Fatalf("loadConfig(target_dir_only) should be succeeded: %s", err)
	}
	if c.Source != "" {
		t.Errorf("c.Source should be empty: %q", c.Source)
	}
	if c.TargetDir != "bar" {
		t.Errorf("c.TargetDir should be \"bar\": %q", c.TargetDir)
	}
	if c.CPU != "" {
		t.Errorf("c.CPU should be empty: %q", c.CPU)
	}
}

func TestLoadConfigCPU(t *testing.T) {
	c, err := loadConfig("test_data/cpu_only.ini")
	if err != nil {
		t.Fatalf("loadConfig(target_dir_only) should be succeeded: %s", err)
	}
	if c.Source != "" {
		t.Errorf("c.Source should be empty: %q", c.Source)
	}
	if c.TargetDir != "" {
		t.Errorf("c.TargetDir should be empty: %q", c.TargetDir)
	}
	if c.CPU != "baz" {
		t.Errorf("c.CPU should be \"baz\": %q", c.CPU)
	}
}

func TestLoadConfigAll(t *testing.T) {
	c, err := loadConfig("test_data/all.ini")
	if err != nil {
		t.Fatalf("loadConfig(target_dir_only) should be succeeded: %s", err)
	}
	if c.Source != "foo" {
		t.Errorf("c.Source should be \"foo\": %q", c.Source)
	}
	if c.TargetDir != "bar" {
		t.Errorf("c.TargetDir should be \"bar\": %q", c.TargetDir)
	}
	if c.CPU != "baz" {
		t.Errorf("c.CPU should be \"baz\": %q", c.CPU)
	}
}
