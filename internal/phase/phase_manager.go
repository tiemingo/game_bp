package phase

import (
	"game_bp/util"
	"time"
)

type PhaseManager struct {
	controlChan chan Command  // Channel for receiving commands to control the phase timer
	doneChan    chan struct{} // Channel to close when the engine is done

	// Phase information
	currentPhase Phase
	phases       map[Phase]func() (Phase, time.Duration, bool) // Map of phase to its logic function and next phase with duration

	// Timer state tracking
	duration  time.Duration // Full duration of the current phase
	remaining time.Duration // Time left in the current phase (used for pause/resume)
	startTime time.Time     // When the current timer started
	isPaused  bool
}

type Config struct {
	InitialPhase      Phase
	InitialDuration   time.Duration
	Phases            map[Phase]func() (Phase, time.Duration, bool)
	CommandBufferSize *int
}

// NewPhaseManager initializes a new PhaseManager with the provided configuration.
func NewPhaseManager(config Config) (*PhaseManager, chan struct{}, error) {

	// Validation of the configuration
	if len(config.Phases) == 0 {
		return nil, nil, util.IErrUnsetPhases
	}

	if config.InitialDuration <= 0 {
		return nil, nil, util.IErrInvalidInitialDuration
	}

	if _, exists := config.Phases[config.InitialPhase]; !exists {
		return nil, nil, util.IErrMissingInitialPhase
	}

	if _, exists := config.Phases[""]; exists {
		return nil, nil, util.IErrEmptyPhase
	}

	// Set command buffer size with a default value
	bufferSize := 64
	if config.CommandBufferSize != nil {
		bufferSize = *config.CommandBufferSize
	}

	doneChan := make(chan struct{})

	return &PhaseManager{
		controlChan:  make(chan Command, bufferSize),
		currentPhase: config.InitialPhase,
		duration:     config.InitialDuration,
		doneChan:     doneChan,
		phases:       config.Phases,
	}, doneChan, nil
}

// advancePhase runs the phase end func and sets the next phase and duration.
// If initial is true, it only sets the first phase without running phase end func.
func (pm *PhaseManager) advancePhase(initial bool) bool {
	shouldContinue := true
	if !initial {
		pm.currentPhase, pm.duration, shouldContinue = pm.phases[pm.currentPhase]()
	}

	// Reset durations for the brand new phase
	pm.remaining = pm.duration
	pm.isPaused = false
	return shouldContinue
}

// stopAndDrain ensures that the timer is stopped and any pending events are cleared.
func (pm *PhaseManager) stopAndDrain(timer *time.Timer) {
	if timer == nil {
		return
	}
	if !timer.Stop() {
		// If timer.Stop() is false, the timer already expired and a value is sitting in timer.C.
		// We MUST drain it to prevent a ghost timeout event.
		select {
		case <-timer.C:
		default:
		}
	}
}

// pauseTimer stops the timer and stores the remaining time.
func (pm *PhaseManager) pauseTimer(timer *time.Timer) {
	if pm.isPaused {
		return
	}
	pm.stopAndDrain(timer)

	// Calculate exactly how much time was consumed before pausing.
	elapsed := time.Since(pm.startTime)
	pm.remaining -= elapsed
	if pm.remaining < 0 {
		pm.remaining = 0
	}
	pm.isPaused = true
}

// resumeTimer restarts the timer with the remaining time.
func (pm *PhaseManager) resumeTimer() *time.Timer {
	if !pm.isPaused {
		return nil
	}
	pm.isPaused = false
	pm.startTime = time.Now().UTC()
	return time.NewTimer(pm.remaining)
}

// restartCurrentPhase resets the current phase's timer back to its full duration, effectively restarting the phase.
func (pm *PhaseManager) restartCurrentPhase(timer *time.Timer) *time.Timer {
	pm.stopAndDrain(timer)
	pm.remaining = pm.duration
	pm.isPaused = false
	pm.startTime = time.Now().UTC()
	return time.NewTimer(pm.duration)
}
