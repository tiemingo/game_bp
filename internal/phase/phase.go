package phase

import "time"

type Phase string

type CommandType int

const (
	CmdPauseIf CommandType = iota
	CmdResumeIf
	CmdSkipIf
	CmdRestartIf
	CmdGetTimerStatus
	CmdStop
)

type TimerStatus struct {
	CurrentPhase Phase
	EndTime      time.Time
	Remaining    time.Duration
	IsPaused     bool
}

type Command struct {
	Type      CommandType
	SkipIf    func(TimerStatus) bool // Used for CmdSkipIf
	PauseIf   func(TimerStatus) bool // Used for CmdPauseIf
	ResumeIf  func(TimerStatus) bool // Used for CmdResumeIf
	RestartIf func(TimerStatus) bool // Used for CmdRestartIf
	DoneChan  chan struct{}          // Used for CmdSkipIf, CmdPauseIf, CmdResumeIf, CmdRestartIf, closes when command is processed
	ReplyChan chan TimerStatus       // Used for CmdGetTimerStatus
}
