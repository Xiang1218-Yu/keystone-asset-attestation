package core

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type Module interface {
	Key() string
	Description() string
	Priority() int
	Family() string
	Enabled() bool
	Validate(string) error
	Rewrite(string) string
	Transition(string) bool
	Evidence(string, time.Time) Evidence
	Score(string) int
}

type Evidence struct {
	Module string    `json:"module"`
	Digest string    `json:"digest"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type Record struct {
	ID        string       `json:"id"`
	Stage     string       `json:"stage"`
	Payload   string       `json:"payload"`
	Score     int          `json:"score"`
	Version   int64        `json:"version"`
	Evidence  []Evidence   `json:"evidence"`
	History   []Transition `json:"history"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Transition struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	By   string    `json:"by"`
	At   time.Time `json:"at"`
}

type Engine struct {
	mu       sync.RWMutex
	modules  []Module
	records  map[string]Record
	sequence uint64
}

func NewEngine() *Engine {
	modules := []Module{
		NewAssetModule(),
		NewSupplierModule(),
		NewLineageModule(),
		NewAttestationModule(),
		NewInspectionModule(),
		NewCampaignModule(),
		NewFindingModule(),
		NewHazardModule(),
		NewMaintenanceModule(),
		NewVerifierModule(),
		NewCertificateModule(),
		NewComponentModule(),
		NewLotModule(),
		NewSerialModule(),
		NewLocationModule(),
		NewCustodyModule(),
		NewCalibrationModule(),
		NewWarrantyModule(),
		NewRiskModule(),
		NewControlModule(),
		NewMitigationModule(),
		NewEvidenceModule(),
		NewReferenceModule(),
		NewReleaseModule(),
		NewGateModule(),
		NewDecisionModule(),
		NewApproverModule(),
		NewExceptionModule(),
		NewWaiverModule(),
		NewExpiryModule(),
		NewRenewalModule(),
		NewServiceModule(),
		NewReplacementModule(),
		NewImpactModule(),
		NewDependencyModule(),
		NewCriticalityModule(),
		NewConditionModule(),
		NewFailureModule(),
		NewActionModule(),
		NewClosureModule(),
		NewSignatureModule(),
		NewSealModule(),
		NewDigestModule(),
		NewManifestModule(),
		NewProvenanceModule(),
		NewContractModule(),
		NewObligationModule(),
		NewMilestoneModule(),
		NewVendorModule(),
		NewFacilityModule(),
		NewZoneModule(),
		NewRouteModule(),
		NewInterfaceModule(),
		NewOwnerModule(),
		NewOperatorModule(),
		NewNotificationModule(),
		NewEscalationModule(),
		NewPolicyModule(),
		NewExposureModule(),
		NewResilienceModule(),
		NewForecastModule(),
		NewRecoveryModule(),
	}
	sort.Slice(modules, func(i, j int) bool { return modules[i].Priority() < modules[j].Priority() })
	return &Engine{modules: modules, records: make(map[string]Record)}
}

func (e *Engine) Modules() []Module {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]Module(nil), e.modules...)
}

func (e *Engine) Create(ctx context.Context, id, payload, actor string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	id = strings.TrimSpace(id)
	if id == "" || strings.TrimSpace(payload) == "" {
		return Record{}, errors.New("id and payload are required")
	}
	// The existence check and the write must happen under the same write
	// lock, otherwise two concurrent collectors that both observe the id is
	// absent can each proceed to create the record, producing duplicate
	// successful creates while only the last write survives in the map.
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.records[id]; exists {
		return Record{}, errors.New("record already exists")
	}
	now := time.Now().UTC()
	result := Record{ID: id, Stage: "captured", Payload: strings.TrimSpace(payload), Version: 1, CreatedAt: now, UpdatedAt: now}
	for _, module := range e.modules {
		if !module.Enabled() {
			continue
		}
		if err := module.Validate(result.Payload); err != nil {
			return Record{}, err
		}
		result.Payload = module.Rewrite(result.Payload)
		result.Score += module.Score(result.Payload)
		result.Evidence = append(result.Evidence, module.Evidence(result.Payload, now))
	}
	result.History = append(result.History, Transition{From: "", To: result.Stage, By: actor, At: now})
	e.records[id] = cloneRecord(result)
	return cloneRecord(result), nil
}

func (e *Engine) Get(ctx context.Context, id string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	record, ok := e.records[strings.TrimSpace(id)]
	if !ok {
		return Record{}, errors.New("record not found")
	}
	return cloneRecord(record), nil
}

func (e *Engine) Advance(ctx context.Context, id, stage, actor string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	record, ok := e.records[strings.TrimSpace(id)]
	if !ok {
		return Record{}, errors.New("record not found")
	}
	stage = strings.ToLower(strings.TrimSpace(stage))
	if stage == "" || stage == record.Stage {
		return Record{}, errors.New("a different target stage is required")
	}
	for _, module := range e.modules {
		if !module.Transition(stage) {
			return Record{}, errors.New("stage is not accepted by module " + module.Key())
		}
	}
	now := time.Now().UTC()
	record.History = append(record.History, Transition{From: record.Stage, To: stage, By: actor, At: now})
	record.Stage = stage
	record.Version++
	record.UpdatedAt = now
	record.Payload = strings.TrimSpace(record.Payload)
	e.records[record.ID] = cloneRecord(record)
	return cloneRecord(record), nil
}

func (e *Engine) List(ctx context.Context, stage string, limit int) ([]Record, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]Record, 0, len(e.records))
	for _, record := range e.records {
		if stage != "" && record.Stage != stage {
			continue
		}
		result = append(result, cloneRecord(record))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (e *Engine) ValidatePayload(ctx context.Context, payload string) (int, []string, error) {
	if err := contextError(ctx); err != nil {
		return 0, nil, err
	}
	score := 0
	warnings := make([]string, 0)
	for _, module := range e.modules {
		if err := module.Validate(payload); err != nil {
			return 0, warnings, err
		}
		score += module.Score(payload)
		if len(payload) < module.Priority() {
			warnings = append(warnings, module.Key()+" expects richer context")
		}
	}
	return score, warnings, nil
}

func (e *Engine) Snapshot() map[string]int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	snapshot := map[string]int{"modules": len(e.modules), "records": len(e.records)}
	for _, record := range e.records {
		snapshot["stage_"+record.Stage]++
	}
	return snapshot
}

func contextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func cloneRecord(record Record) Record {
	record.Evidence = append([]Evidence(nil), record.Evidence...)
	record.History = append([]Transition(nil), record.History...)
	return record
}
