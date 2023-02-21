package main

// import "fmt"
import "math"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type ControlState struct {
	WalkForward  bool
	WalkBackward bool
	StrafeLeft   bool
	StrafeRight  bool
	LookLeft     bool
	LookRight    bool
}

type Raycaster struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	config config.Wrapped

	Camera
	controlState ControlState
	world World
	onControlStateChange func (ControlState)
}

func NewRaycaster (world World) (element *Raycaster) {
	element = &Raycaster {
		Camera: Camera {
			Vector: Vector {
				X: 1.5,
				Y: 1.5,
			},
			Angle: math.Pi / 3,
			Fov:   1,
		},
		world: world,
	}
	element.Core, element.core = core.NewCore(element.drawAll)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.Draw)
	element.core.SetMinimumSize(64, 64)
	return
}

func (element *Raycaster) OnControlStateChange (callback func (ControlState)) {
	element.onControlStateChange = callback
}

func (element *Raycaster) Draw () {
	if element.core.HasImage() {
		element.drawAll()
		element.core.DamageAll()
	}
}

func (element *Raycaster) HandleMouseDown (x, y int, button input.Button) {
	if !element.Focused() { element.Focus() }
}

func (element *Raycaster) HandleMouseUp (x, y int, button input.Button) { }
func (element *Raycaster) HandleMouseMove (x, y int) { }
func (element *Raycaster) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Raycaster) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	switch key {
	case input.KeyLeft:  element.controlState.LookLeft  = true
	case input.KeyRight: element.controlState.LookRight = true
	case 'a', 'A': element.controlState.StrafeLeft   = true
	case 'd', 'D': element.controlState.StrafeRight  = true
	case 'w', 'W': element.controlState.WalkForward  = true
	case 's', 'S': element.controlState.WalkBackward = true
	default: return
	}

	if element.onControlStateChange != nil {
		element.onControlStateChange(element.controlState)
	}
}

func (element *Raycaster) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	switch key {
	case input.KeyLeft:  element.controlState.LookLeft  = false
	case input.KeyRight: element.controlState.LookRight = false
	case 'a', 'A': element.controlState.StrafeLeft   = false
	case 'd', 'D': element.controlState.StrafeRight  = false
	case 'w', 'W': element.controlState.WalkForward  = false
	case 's', 'S': element.controlState.WalkBackward = false
	default: return
	}

	if element.onControlStateChange != nil {
		element.onControlStateChange(element.controlState)
	}
}

func (element *Raycaster) drawAll () {
	bounds := element.Bounds()
	// artist.FillRectangle(element.core, artist.Uhex(0x000000FF), bounds)
	width  := bounds.Dx()
	height := bounds.Dy()

	ray := Ray { Angle: element.Camera.Angle - element.Camera.Fov / 2 }
	
	for x := 0; x < width; x ++ {
		ray.X = element.Camera.X
		ray.Y = element.Camera.Y
		
		distance, _ := ray.Cast(element.world, 8)
		distanceFac := float64(distance) / 8
		distance    *= math.Cos(ray.Angle - element.Camera.Angle)
		
		wallHeight := height
		if distance > 0 {
			wallHeight = int((float64(height) / 2.0) / float64(distance))
		}

		ceilingColor := color.RGBA { 0x00, 0x00, 0x00, 0xFF }
		wallColor    := color.RGBA { 0xCC, 0x33, 0x22, 0xFF }
		floorColor   := color.RGBA { 0x11, 0x50, 0x22, 0xFF }

		// fmt.Println(float64(distance) / 32)

		wallColor  = artist.LerpRGBA(wallColor, ceilingColor, distanceFac)

		// draw
		data, stride := element.core.Buffer()
		wallStart := height / 2 - wallHeight + bounds.Min.Y
		wallEnd   := height / 2 + wallHeight + bounds.Min.Y
		if wallStart < 0            { wallStart = 0 }
		if wallEnd   > bounds.Max.Y { wallEnd   = bounds.Max.Y }
		for y := bounds.Min.Y; y < wallStart; y ++ {
			data[y * stride + x + bounds.Min.X] = ceilingColor
		}
		for y := wallStart; y < wallEnd; y ++ {
			data[y * stride + x + bounds.Min.X] = wallColor
		}
		for y := wallEnd; y < bounds.Max.Y; y ++ {
			floorFac := float64(y - (height / 2)) / float64(height / 2)
			data[y * stride + x + bounds.Min.X] =
				artist.LerpRGBA(ceilingColor, floorColor, floorFac)
		}

		// increment angle
		ray.Angle += element.Camera.Fov / float64(width)
	}

	// element.drawMinimap()
}

func (element *Raycaster) drawMinimap () {
	bounds := element.Bounds()
	scale  := 16
	for y := 0; y < 10; y ++ {
	for x := 0; x < 10; x ++ {
		cellPt := image.Pt(x, y)
		cell   := element.world.At(cellPt)
		cellBounds :=
			image.Rectangle {
				cellPt.Mul(scale),
				cellPt.Add(image.Pt(1, 1)).Mul(scale),
			}.Add(bounds.Min)
		cellColor  := color.RGBA { 0x22, 0x22, 0x22, 0xFF }
		if cell == 1 {
			cellColor = color.RGBA { 0xFF, 0xFF, 0xFF, 0xFF }
		}
		artist.FillRectangle (
			element.core,
			artist.NewUniform(cellColor),
			cellBounds.Inset(1))
	}}

	playerPt := element.Camera.Mul(float64(scale)).Point().Add(bounds.Min)
	playerAnglePt :=
		element.Camera.Add(element.Camera.Delta()).
		Mul(float64(scale)).Point().Add(bounds.Min)
	ray := Ray { Vector: element.Camera.Vector, Angle: element.Camera.Angle }
	_, hit := ray.Cast(element.world, 8)
	hitPt := hit.Mul(float64(scale)).Point().Add(bounds.Min)
	// fmt.Println(rayDistance)
	
	playerBounds := image.Rectangle { playerPt, playerPt }.Inset(scale / -8)
	artist.FillEllipse (
		element.core,
		artist.Uhex(0xFFFFFFFF),
		playerBounds)
	artist.Line (
		element.core,
		artist.Uhex(0xFFFFFFFF), 1,
		playerPt,
		playerAnglePt)
	artist.Line (
		element.core,
		artist.Uhex(0x00FF00FF), 1,
		playerPt,
		hitPt)
}
