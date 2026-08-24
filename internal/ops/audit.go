package ops

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"keystone-asset-attestation/internal/core"
)

type AuditEntry struct {
	ID       string    `json:"id"`
	Actor    string    `json:"actor"`
	Action   string    `json:"action"`
	Resource string    `json:"resource"`
	At       time.Time `json:"at"`
}

type AuditLog struct {
	mu       sync.RWMutex
	entries  []AuditEntry
	sequence uint64
}

func NewAuditLog() *AuditLog { return &AuditLog{entries: make([]AuditEntry, 0, 256)} }

// Add records an audit entry. An empty or blank actor — the case when an
// automated import task creates a record without operator information — is
// resolved to core.DefaultActor so the audit page never has to render a blank
// actor (which previously triggered the empty-data handling exception).
// Explicit actors are preserved verbatim.
func (l *AuditLog) Add(actor, action, resource string, now time.Time) AuditEntry {
	if strings.TrimSpace(actor) == "" {
		actor = core.DefaultActor
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sequence++
	entry := AuditEntry{ID: fmt.Sprintf("audit-%d-%d", now.UnixNano(), l.sequence), Actor: actor, Action: action, Resource: resource, At: now.UTC()}
	l.entries = append(l.entries, entry)
	if len(l.entries) > 5000 {
		l.entries = append([]AuditEntry(nil), l.entries[len(l.entries)-5000:]...)
	}
	return entry
}

func (l *AuditLog) List(limit int) []AuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := append([]AuditEntry(nil), l.entries...)
	if limit > 0 && len(result) > limit {
		result = result[len(result)-limit:]
	}
	return result
}
