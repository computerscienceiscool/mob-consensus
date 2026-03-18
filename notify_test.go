package main

// Unit tests for notify.go — Notifier interface and implementations.
//
// These tests exercise the notification types in isolation without requiring
// any wired-up workflow. The ntfyNotifier tests use httptest to verify HTTP
// behavior; the newNotifier tests use temporary git repos to exercise
// git-config-based construction.
//
// Note: newNotifier tests use os.Chdir (process-wide state) and therefore
// cannot run in parallel with each other. They are grouped under a single
// parent test to run sequentially.

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNoopNotifierDoesNotPanic ensures the noop path is safe to call.
func TestNoopNotifierDoesNotPanic(t *testing.T) {
	t.Parallel()
	n := &noopNotifier{}
	n.Notify(context.Background(), Event{
		Type:    EventMerged,
		User:    "alice",
		Twig:    "feat",
		Branch:  "alice/feat",
		Message: "merged feat",
	})
}

// TestNtfyNotifierSendsRequest verifies that Notify POSTs to the correct
// URL with the expected headers and body.
func TestNtfyNotifierSendsRequest(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var gotReq struct {
		method string
		path   string
		title  string
		tags   string
		body   string
	}

	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotReq.method = r.Method
		gotReq.path = r.URL.Path
		gotReq.title = r.Header.Get("Title")
		gotReq.tags = r.Header.Get("Tags")
		b, _ := io.ReadAll(r.Body)
		gotReq.body = string(b)
		w.WriteHeader(http.StatusOK)
		close(done)
	}))
	defer srv.Close()

	n := &ntfyNotifier{
		topic:  "test-topic",
		server: srv.URL,
	}

	n.Notify(context.Background(), Event{
		Type:    EventStarted,
		User:    "bob",
		Twig:    "feat",
		Branch:  "bob/feat",
		Message: "bob started feat",
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for HTTP request")
	}

	mu.Lock()
	defer mu.Unlock()

	if gotReq.method != "POST" {
		t.Fatalf("method=%q, want POST", gotReq.method)
	}
	if gotReq.path != "/test-topic" {
		t.Fatalf("path=%q, want /test-topic", gotReq.path)
	}
	if gotReq.title != "mob-consensus: started" {
		t.Fatalf("title=%q, want %q", gotReq.title, "mob-consensus: started")
	}
	if gotReq.tags != "started" {
		t.Fatalf("tags=%q, want %q", gotReq.tags, "started")
	}
	if gotReq.body != "bob started feat" {
		t.Fatalf("body=%q, want %q", gotReq.body, "bob started feat")
	}
}

// TestNtfyNotifierEmptyMessageFallback verifies that an empty Message field
// falls back to the event type string.
func TestNtfyNotifierEmptyMessageFallback(t *testing.T) {
	t.Parallel()

	done := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		done <- string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := &ntfyNotifier{topic: "t", server: srv.URL}
	n.Notify(context.Background(), Event{Type: EventJoined})

	select {
	case body := <-done:
		if body != "joined" {
			t.Fatalf("body=%q, want %q", body, "joined")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for HTTP request")
	}
}

// TestNewNotifier groups all newNotifier unit tests. These use os.Chdir
// (process-wide) so they must run sequentially, not in parallel.
func TestNewNotifier(t *testing.T) {
	t.Run("NoTopic", func(t *testing.T) {
		dir := initTestGitDir(t)
		ctx := chdirCtx(t, dir)

		n := newNotifier(ctx)
		if _, ok := n.(*noopNotifier); !ok {
			t.Fatalf("newNotifier() returned %T, want *noopNotifier", n)
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		dir := initTestGitDir(t)
		gitCfg(t, dir, "mob-consensus.ntfyTopic", "my-topic")
		gitCfg(t, dir, "mob-consensus.ntfyEnabled", "false")
		ctx := chdirCtx(t, dir)

		n := newNotifier(ctx)
		if _, ok := n.(*noopNotifier); !ok {
			t.Fatalf("newNotifier() returned %T, want *noopNotifier", n)
		}
	})

	t.Run("WithTopicAndServer", func(t *testing.T) {
		dir := initTestGitDir(t)
		gitCfg(t, dir, "mob-consensus.ntfyTopic", "my-topic")
		gitCfg(t, dir, "mob-consensus.ntfyServer", "https://custom.example.com")
		ctx := chdirCtx(t, dir)

		n := newNotifier(ctx)
		nn, ok := n.(*ntfyNotifier)
		if !ok {
			t.Fatalf("newNotifier() returned %T, want *ntfyNotifier", n)
		}
		if nn.topic != "my-topic" {
			t.Fatalf("topic=%q, want %q", nn.topic, "my-topic")
		}
		if nn.server != "https://custom.example.com" {
			t.Fatalf("server=%q, want %q", nn.server, "https://custom.example.com")
		}
	})

	t.Run("DefaultServer", func(t *testing.T) {
		dir := initTestGitDir(t)
		gitCfg(t, dir, "mob-consensus.ntfyTopic", "my-topic")
		ctx := chdirCtx(t, dir)

		n := newNotifier(ctx)
		nn, ok := n.(*ntfyNotifier)
		if !ok {
			t.Fatalf("newNotifier() returned %T, want *ntfyNotifier", n)
		}
		if nn.server != "https://ntfy.sh" {
			t.Fatalf("server=%q, want %q", nn.server, "https://ntfy.sh")
		}
	})
}

// --- helpers ---

// initTestGitDir creates a minimal git repo in a temp directory and returns
// the path. The repo has an initial commit so git-config works normally.
func initTestGitDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "git", "init")
	runGit(t, dir, "git", "config", "user.email", "test@test.com")
	runGit(t, dir, "git", "config", "user.name", "Test")
	runGit(t, dir, "git", "commit", "--allow-empty", "-m", "init")
	return dir
}

// gitCfg sets a git config value in the given repo directory.
func gitCfg(t *testing.T, dir, key, value string) {
	t.Helper()
	runGit(t, dir, "git", "config", key, value)
}

// runGit executes a command in the given directory and fails the test on error.
func runGit(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %s: %v", name, strings.Join(args, " "), err)
	}
}

// chdirCtx changes the working directory for the duration of the test and
// returns a background context. This is needed because gitOutputTrimmed
// runs git in the current working directory.
func chdirCtx(t *testing.T, dir string) context.Context {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
	return context.Background()
}
