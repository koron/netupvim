package main

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

func TestLoadGithub(t *testing.T) {
	c, err := loadConfig("test_data/github.ini")
	if err != nil {
		t.Fatalf("loadConfig(github) should be succeeded: %s", err)
	}
	if c.GithubUser != "foo" {
		t.Errorf("c.GithubUser should be %q: %q", "foo", c.GithubUser)
	}
	if c.GithubToken != "0123456789abcdef" {
		t.Errorf("c.GithubToken should be %q: %q", "0123456789abcdef", c.GithubToken)
	}
}

func TestLoadTimeout(t *testing.T) {
	c, err := loadConfig("test_data/timeout.ini")
	if err != nil {
		t.Fatalf("loadConfig(timeout) should be succeeded: %s", err)
	}
	if c.DownloadTimeout != "1200s" {
		t.Errorf("c.DownloadTimeout should be %q: %q", "1200s", c.DownloadTimeout)
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
	if c.LogRotateCount != 1234 {
		t.Errorf("c.LogRotateCount is unexpected: %d", c.LogRotateCount)
	}
	if c.ExeRotateCount != 5678 {
		t.Errorf("c.ExeRotateCount is unexpected: %d", c.ExeRotateCount)
	}
}
