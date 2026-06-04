package phase

import "time"

// getTimerStatus returns the current timer status, calculating remaining time based on whether it's paused or not.
// This should only be called from within the RunEngine loop to ensure thread safety.
func (pm *PhaseManager) getTimerStatus() TimerStatus {
	currentRemaining := pm.remaining
	if !pm.isPaused {
		currentRemaining -= time.Since(pm.startTime)
	}

	return TimerStatus{
		CurrentPhase: pm.currentPhase,
		Remaining:    currentRemaining,
		EndTime:      time.Now().UTC().Add(currentRemaining),
		IsPaused:     pm.isPaused,
	}
}
