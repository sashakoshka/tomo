package main

import "time"
import "git.tebibyte.media/sashakoshka/tomo"

type Game struct {
	*Raycaster
	running bool
	tickChan <- chan time.Time
	stopChan chan bool

	stamina float64
	health  float64

	controlState ControlState
	onStatUpdate func ()
}

func NewGame (world World, textures Textures) (game *Game) {
	game = &Game {
		Raycaster: NewRaycaster(world, textures),
		stopChan: make(chan bool),
	}
	game.Raycaster.OnControlStateChange (func (state ControlState) {
		game.controlState = state
	})
	game.stamina = 0.5
	game.health  = 1
	return
}

func (game *Game) Start () {
	if game.running == true { return }
	game.running = true
	go game.run()
}

func (game *Game) Stop () {
	select {
	case game.stopChan <- true:
	default:
	}	
}

func (game *Game) Stamina () float64 {
	return game.stamina
}

func (game *Game) Health () float64 {
	return game.health
}

func (game *Game) OnStatUpdate (callback func ()) {
	game.onStatUpdate = callback
}

func (game *Game) tick () {
	moved := false
	statUpdate := false
	
	speed := 0.07
	if game.controlState.Sprint {
		speed = 0.16
	}
	if game.stamina <= 0 {
		speed = 0
	}
	
	if game.controlState.WalkForward {
		game.Walk(speed)
		moved = true
	}
	if game.controlState.WalkBackward {
		game.Walk(-speed)
		moved = true
	}
	if game.controlState.StrafeLeft {
		game.Strafe(-speed)
		moved = true
	}
	if game.controlState.StrafeRight {
		game.Strafe(speed)
		moved = true
	}
	if game.controlState.LookLeft {
		game.Rotate(-0.1)
	}
	if game.controlState.LookRight {
		game.Rotate(0.1)
	}

	if moved {
		game.stamina -= speed / 50
		statUpdate = true
	} else if game.stamina < 1 {
		game.stamina += 0.005
		statUpdate = true
	} 
	
	if game.stamina > 1 {
		game.stamina = 1
	}
	if game.stamina < 0 {
		game.stamina = 0
	}

	tomo.Do(game.Invalidate)
	if statUpdate && game.onStatUpdate != nil {
		tomo.Do(game.onStatUpdate)
	}
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
