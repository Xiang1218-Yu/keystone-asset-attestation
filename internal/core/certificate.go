package core

import (
	"fmt"
	"strings"
	"time"
)

type CertificateModule struct {
	key      string
	priority int
	family   string
	focus    string
}

func NewCertificateModule() Module {
	return CertificateModule{key: "certificate", priority: 11, family: "asset", focus: "certificate validity and renewal"}
}

func (m CertificateModule) Key() string {
	return m.key
}

func (m CertificateModule) Description() string {
	return m.focus
}

func (m CertificateModule) Priority() int {
	return m.priority
}

func (m CertificateModule) Family() string {
	return m.family
}

func (m CertificateModule) Enabled() bool {
	return m.priority > 0 && strings.TrimSpace(m.key) != ""
}

func (m CertificateModule) Validate(input string) error {
	value := strings.TrimSpace(input)
	if value == "" {
		return fmt.Errorf("%s requires a non-empty payload", m.key)
	}
	if len(value) > 4096 {
		return fmt.Errorf("%s payload exceeds the domain limit", m.key)
	}
	if strings.Contains(value, "\x00") {
		return fmt.Errorf("%s payload contains an invalid byte", m.key)
	}
	return nil
}

func (m CertificateModule) Rewrite(input string) string {
	value := strings.TrimSpace(input)
	value = strings.Join(strings.Fields(value), " ")
	marker := fmt.Sprintf("[%s:%d]", m.key, m.priority)
	if strings.Contains(value, marker) {
		return value
	}
	return value + " " + marker
}

func (m CertificateModule) Transition(stage string) bool {
	stage = strings.ToLower(strings.TrimSpace(stage))
	switch m.family {
	case "asset":
		return stage == "captured" || stage == "verified" || stage == "sealed" || stage == "released"
	case "lexicon":
		return stage == "draft" || stage == "review" || stage == "published" || stage == "retired"
	default:
		return false
	}
}

func (m CertificateModule) Evidence(input string, now time.Time) Evidence {
	normalized := m.Rewrite(input)
	return Evidence{
		Module: m.key,
		Digest: fmt.Sprintf("%s-%d-%d", m.key, m.priority, len(normalized)),
		Detail: normalized,
		At:     now.UTC(),
	}
}

func (m CertificateModule) Score(input string) int {
	value := strings.TrimSpace(input)
	score := len(value) + m.priority
	if strings.Contains(strings.ToLower(value), m.key) {
		score += m.priority
	}
	if score > 100 {
		return 100
	}
	return score
}
