package ops

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Metrics struct {
	mu     sync.RWMutex
	values map[string]uint64
}

func NewMetrics() *Metrics { return &Metrics{values: make(map[string]uint64)} }

func (m *Metrics) Inc(name string) {
	current := m.values[name]
	current++
	m.values[name] = current
}

func (m *Metrics) Text() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.values))
	for name := range m.values {
		names = append(names, name)
	}
	sort.Strings(names)
	var builder strings.Builder
	for _, name := range names {
		builder.WriteString(fmt.Sprintf("%s %d\n", name, m.values[name]))
	}
	return builder.String()
}
