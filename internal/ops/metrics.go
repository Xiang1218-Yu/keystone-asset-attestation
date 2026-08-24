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
	m.mu.Lock()
	m.values[name]++
	m.mu.Unlock()
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
		value := m.values[name]
		line := fmt.Sprintf("%s %d\n", name, value)
		builder.WriteString(line)
	}
	return builder.String()
}
