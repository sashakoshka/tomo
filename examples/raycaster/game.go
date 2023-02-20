package main

import "time"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

type Game struct {
	*Raycaster
	running bool
	tickChan <- chan time.Time
	stopChan chan bool

	controlState ControlState
}

func NewGame (world World) (game *Game) {
	game = &Game {
		Raycaster: NewRaycaster(world),
		stopChan: make(chan bool),
	}
	game.Raycaster.OnControlStateChange (func (state ControlState) {
		game.controlState = state
	})
	return
}

func (game *Game) DrawTo (canvas canvas.Canvas) {
	if canvas == nil {
		game.stopChan <- true
	} else if !game.running {
		game.running = true
		go game.run()
	}
	game.Raycaster.DrawTo(canvas)
}

func (game *Game) tick () {
	if game.controlState.WalkForward {
		game.Walk(0.1)
	}
	if game.controlState.WalkBackward {
		game.Walk(-0.1)
	}
	if game.controlState.StrafeLeft {
		game.Strafe(-0.1)
	}
	if game.controlState.StrafeRight {
		game.Strafe(0.1)
	}
	if game.controlState.LookLeft {
		game.Rotate(-0.1)
	}
	if game.controlState.LookRight {
		game.Rotate(0.1)
	}

	tomo.Do(game.Draw)
}

func (game *Game) run () {
	ticker := time.NewTicker(time.Second / 30)
	game.tickChan = ticker.C
	for game.running {
		select {
		case <- game.tickChan:
			game.tick()
		case <- game.stopChan:
			ticker.Stop()
		}
	}
}
