package state

import "fmt"

type Stage string

const (
	StageInit       Stage = "init"
	StageResearch   Stage = "research"
	StageArchitect  Stage = "architect"
	StageWrite      Stage = "write"
	StageReview     Stage = "review"
	StageInbox      Stage = "inbox"
	StageFinalize   Stage = "finalize"
	StageDone       Stage = "done"
)

var stageOrder = []Stage{
	StageInit,
	StageResearch,
	StageArchitect,
	StageWrite,
	StageReview,
	StageInbox,
	StageFinalize,
	StageDone,
}

func (s Stage) String() string {
	return string(s)
}

func ParseStage(s string) (Stage, error) {
	for _, stage := range stageOrder {
		if string(stage) == s {
			return stage, nil
		}
	}
	return "", fmt.Errorf("unknown stage: %q", s)
}

func (s Stage) Next() (Stage, error) {
	for i, stage := range stageOrder {
		if stage == s && i+1 < len(stageOrder) {
			return stageOrder[i+1], nil
		}
	}
	return "", fmt.Errorf("no stage after %q", s)
}

func (s Stage) Before(other Stage) bool {
	si, oi := -1, -1
	for i, stage := range stageOrder {
		if stage == s {
			si = i
		}
		if stage == other {
			oi = i
		}
	}
	return si < oi
}

func ValidateTransition(from, to Stage) error {
	if !from.Before(to) && from != to {
		return fmt.Errorf("cannot transition from %q to %q: invalid order", from, to)
	}
	return nil
}
