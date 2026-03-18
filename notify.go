// notify.go — Team notification support via ntfy.
//
// See TODO/J016-team-notifications-ntfy.md for the integration plan
// describing how to wire these types into the existing workflow functions.
package main

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// EventType identifies the kind of workflow event that occurred.
type EventType string

const (
	// EventMerged is sent after a successful merge + push.
	EventMerged EventType = "merged"
	// EventStarted is sent after a twig is started and pushed.
	EventStarted EventType = "started"
	// EventJoined is sent after a user joins a twig and pushes.
	EventJoined EventType = "joined"
	// EventAutoCommitted is sent after dirty changes are auto-committed and pushed.
	EventAutoCommitted EventType = "auto-committed"
)

// Event describes a workflow event that should be communicated to the team.
type Event struct {
	Type    EventType
	User    string // who triggered the event
	Twig    string // shared twig name (e.g. "feature-x")
	Branch  string // full branch name (e.g. "alice/feature-x")
	Message string // human-readable summary
}

// Notifier sends workflow event notifications to the team.
type Notifier interface {
	// Notify sends an event notification asynchronously. It must never
	// block the caller or return an error — failures are silently ignored.
	Notify(ctx context.Context, event Event)
}

// ntfyNotifier sends notifications via an ntfy server using HTTP POST.
type ntfyNotifier struct {
	topic  string // ntfy topic name
	server string // ntfy server URL (e.g. "https://ntfy.sh")
}

// noopNotifier discards all events. Returned when no topic is configured.
type noopNotifier struct{}

func (n *noopNotifier) Notify(_ context.Context, _ Event) {}

// newNotifier reads ntfy configuration from git config and returns the
// appropriate Notifier implementation. Returns a noopNotifier if no topic
// is configured or notifications are explicitly disabled.
func newNotifier(ctx context.Context) Notifier {
	topic, _ := gitOutputTrimmed(ctx, "config", "--get", "mob-consensus.ntfyTopic")
	if topic == "" {
		return &noopNotifier{}
	}

	enabled, _ := gitOutputTrimmed(ctx, "config", "--get", "mob-consensus.ntfyEnabled")
	if strings.EqualFold(enabled, "false") {
		return &noopNotifier{}
	}

	server, _ := gitOutputTrimmed(ctx, "config", "--get", "mob-consensus.ntfyServer")
	if server == "" {
		server = "https://ntfy.sh"
	}

	return &ntfyNotifier{
		topic:  topic,
		server: server,
	}
}

// Notify sends the event to the configured ntfy topic in a detached
// goroutine with a short timeout. It never blocks and silently ignores
// all errors.
func (n *ntfyNotifier) Notify(_ context.Context, event Event) {
	url := n.server + "/" + n.topic

	// Build a title from the event type.
	title := "mob-consensus: " + string(event.Type)

	// Use the human-readable message as the body.
	body := event.Message
	if body == "" {
		body = string(event.Type)
	}

	// Tags for ntfy filtering.
	tags := string(event.Type)

	// Fire and forget — detached context so the goroutine survives if the
	// caller's context is cancelled.
	sendCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	go func() {
		defer cancel()
		req, err := http.NewRequestWithContext(sendCtx, "POST", url, strings.NewReader(body))
		if err != nil {
			return
		}
		req.Header.Set("Title", title)
		req.Header.Set("Tags", tags)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		_ = resp.Body.Close()
	}()
}
