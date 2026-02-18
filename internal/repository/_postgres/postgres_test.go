package _postgres

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"tasks_assignment/pkg/modules"
)

func TestNewPGXDialectRunsMigrations(t *testing.T) {
	repoRoot := projectRoot(t)
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}

	// Migrate looks up file://database/migrations relative to process cwd.
	if err := chdir(repoRoot); err != nil {
		t.Fatalf("change working dir: %v", err)
	}
	t.Cleanup(func() {
		_ = chdir(oldWD)
	})

	cfg := &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5434",
		Username:    "postgres",
		Password:    "postgres",
		DBName:      "go_kbtu",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	var d *Dialect
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("NewPGXDialect panicked: %v", r)
			}
		}()
		d = NewPGXDialect(context.Background(), cfg)
	}()

	if d == nil || d.DB == nil {
		t.Fatal("expected non-nil dialect and db")
	}
	defer d.DB.Close()

	var exists bool
	err = d.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'tasks'
		)
	`)
	if err != nil {
		t.Fatalf("check tasks table existence: %v", err)
	}

	if !exists {
		t.Fatal("expected tasks table to exist after migrations")
	}
}

func projectRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

func chdir(dir string) error {
	return os.Chdir(dir)
}
