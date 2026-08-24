package ops

import (
	"fmt"
	"sync"
	"time"
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

func (l *AuditLog) Add(actor, action, resource string, now time.Time) AuditEntry {
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
	// Return an isolated copy so callers (e.g. the audit page) cannot mutate
	// the internal backing array. The copy is built from a fresh allocation,
	// and the limit window slices into that same copy, never into l.entries.
	result := append([]AuditEntry(nil), l.entries...)
	if limit > 0 && len(result) > limit {
		result = result[len(result)-limit:]
	}
	return result
}
