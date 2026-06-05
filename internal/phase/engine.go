package phase

import "time"

func (pm *PhaseManager) RunEngine() {
	defer close(pm.doneChan)

	pm.advancePhase(true) // Initialize first phase
	pm.startTime = time.Now().UTC()
	phaseTimer := time.NewTimer(pm.remaining)

	for {
		// Check for control commands or timer expiration
		select {
		case cmd := <-pm.controlChan:
			switch cmd.Type {

			case CmdSkipIf:
				if cmd.SkipIf(pm.getTimerStatus()) {

					pm.stopAndDrain(phaseTimer)
					if !pm.advancePhase(false) {
						return
					}
					pm.startTime = time.Now().UTC()
					phaseTimer = time.NewTimer(pm.remaining)
				}
				close(cmd.DoneChan)
			case CmdPauseIf:
				if cmd.PauseIf(pm.getTimerStatus()) {
					pm.pauseTimer(phaseTimer)
				}
				close(cmd.DoneChan)

			case CmdResumeIf:
				if cmd.ResumeIf(pm.getTimerStatus()) {
					phaseTimer = pm.resumeTimer()
				}
				close(cmd.DoneChan)

			case CmdRestartIf:
				if cmd.RestartIf(pm.getTimerStatus()) {
					phaseTimer = pm.restartCurrentPhase(phaseTimer)
				}
				close(cmd.DoneChan)

			case CmdResetIf:
				if cmd.ResetIf(pm.getTimerStatus()) {
					pm.resetTimer(phaseTimer)
				}
				close(cmd.DoneChan)

			case CmdGetTimerStatus:
				cmd.ReplyChan <- pm.getTimerStatus()
				close(cmd.ReplyChan)

			case CmdStop:
				return
			}

		case <-phaseTimer.C:
			// Timer naturally expired, go to next phase
			if !pm.advancePhase(false) {
				return
			}
			pm.startTime = time.Now().UTC()
			phaseTimer = time.NewTimer(pm.remaining)
		}
	}
}
