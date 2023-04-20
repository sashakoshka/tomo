package main

// import "fmt"
import "math"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

type ControlState struct {
	WalkForward  bool
	WalkBackward bool
	StrafeLeft   bool
	StrafeRight  bool
	LookLeft     bool
	LookRight    bool
	Sprint       bool
}

type Raycaster struct {
	entity tomo.FocusableEntity

	config config.Wrapped

	Camera
	controlState ControlState
	world World
	textures Textures
	onControlStateChange func (ControlState)
	renderDistance int
}

func NewRaycaster (world World, textures Textures) (element *Raycaster) {
	element = &Raycaster {
		Camera: Camera {
			Vector: Vector {
				X: 1,
				Y: 1,
			},
			Angle: math.Pi / 3,
			Fov:   1,
		},
		world: world,
		textures: textures,
		renderDistance: 8,
	}
	element.entity = tomo.NewEntity(element).(tomo.FocusableEntity)
	element.entity.SetMinimumSize(64, 64)
	return
}

func (element *Raycaster) Entity () tomo.Entity {
	return element.entity
}

func (element *Raycaster) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	// artist.FillRectangle(element.core, artist.Uhex(0x000000FF), bounds)
	width   := bounds.Dx()
	height  := bounds.Dy()
	halfway := bounds.Max.Y - height / 2

	ray := Ray { Angle: element.Camera.Angle - element.Camera.Fov / 2 }
	
	for x := 0; x < width; x ++ {
		ray.X = element.Camera.X
		ray.Y = element.Camera.Y
		
		distance, hitPoint, wall, horizontal := ray.Cast (
			element.world, element.renderDistance)
		distance *= math.Cos(ray.Angle - element.Camera.Angle)
		textureX := math.Mod(hitPoint.X + hitPoint.Y, 1)
		if textureX < 0 { textureX += 1 }
		
		wallHeight := height
		if distance > 0 {
			wallHeight = int((float64(height) / 2.0) / float64(distance))
		}

		shade := 1.0
		if horizontal {
			shade *= 0.8
		}
		shade *= 1 - distance / float64(element.renderDistance)
		if shade < 0 { shade = 0 }

		ceilingColor := color.RGBA { 0x00, 0x00, 0x00, 0xFF }
		floorColor   := color.RGBA { 0x39, 0x49, 0x25, 0xFF }

		// draw
		data, stride := destination.Buffer()
		wallStart := halfway - wallHeight
		wallEnd   := halfway + wallHeight

		for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
			switch {
			case y < wallStart:
				data[y * stride + x + bounds.Min.X] = ceilingColor

			case y < wallEnd:
				textureY :=
					float64(y - halfway) / 
					float64(wallEnd - wallStart) + 0.5
				// fmt.Println(textureY)
				
				wallColor := element.textures.At (wall, Vector {
					textureX,
					textureY,
				})
				wallColor = shadeColor(wallColor, shade)
				data[y * stride + x + bounds.Min.X] = wallColor
				
			default:
				data[y * stride + x + bounds.Min.X] = floorColor
			}
		}

		// increment angle
		ray.Angle += element.Camera.Fov / float64(width)
	}

	// element.drawMinimap()
}

func (element *Raycaster) Invalidate () {
	element.entity.Invalidate()
}

func (element *Raycaster) OnControlStateChange (callback func (ControlState)) {
	element.onControlStateChange = callback
}

func (element *Raycaster) Focus () {
	element.entity.Focus()
}

func (element *Raycaster) SetEnabled (bool) { }

func (element *Raycaster) Enabled () bool { return true }

func (element *Raycaster) HandleFocusChange () { }

func (element *Raycaster) HandleMouseDown (x, y int, button input.Button) {
	element.entity.Focus()
}

func (element *Raycaster) HandleMouseUp (x, y int, button input.Button) { }

func (element *Raycaster) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	switch key {
	case input.KeyLeft:  element.controlState.LookLeft  = true
	case input.KeyRight: element.controlState.LookRight = true
	case 'a', 'A': element.controlState.StrafeLeft   = true
	case 'd', 'D': element.controlState.StrafeRight  = true
	case 'w', 'W': element.controlState.WalkForward  = true
	case 's', 'S': element.controlState.WalkBackward = true
	case input.KeyLeftControl: element.controlState.Sprint = true
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
	case input.KeyLeftControl: element.controlState.Sprint = false
	default: return
	}

	if element.onControlStateChange != nil {
		element.onControlStateChange(element.controlState)
	}
}

func shadeColor (c color.RGBA, brightness float64) color.RGBA {
	return color.RGBA {
		uint8(float64(c.R) * brightness),
		uint8(float64(c.G) * brightness),
		uint8(float64(c.B) * brightness),
		c.A,
	}
}

func (element *Raycaster) drawMinimap (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	scale  := 8
	for y := 0; y < len(element.world.Data) / element.world.Stride; y ++ {
	for x := 0; x < element.world.Stride; x ++ {
		cellPt := image.Pt(x, y)
		cell   := element.world.At(cellPt)
		cellBounds :=
			image.Rectangle {
				cellPt.Mul(scale),
				cellPt.Add(image.Pt(1, 1)).Mul(scale),
			}.Add(bounds.Min)
		cellColor  := color.RGBA { 0x22, 0x22, 0x22, 0xFF }
		if cell > 0 {
			cellColor = color.RGBA { 0xFF, 0xFF, 0xFF, 0xFF }
		}
		shapes.FillColorRectangle (
			destination,
			cellColor,
			cellBounds.Inset(1))
	}}

	playerPt := element.Camera.Mul(float64(scale)).Point().Add(bounds.Min)
	playerAnglePt :=
		element.Camera.Add(element.Camera.Delta()).
		Mul(float64(scale)).Point().Add(bounds.Min)
	ray := Ray { Vector: element.Camera.Vector, Angle: element.Camera.Angle }
	_, hit, _, _ := ray.Cast(element.world, 8)
	hitPt := hit.Mul(float64(scale)).Point().Add(bounds.Min)
	
	playerBounds := image.Rectangle { playerPt, playerPt }.Inset(scale / -8)
	shapes.FillColorEllipse (
		destination,
		artist.Hex(0xFFFFFFFF),
		playerBounds)
	shapes.ColorLine (
		destination,
		artist.Hex(0xFFFFFFFF), 1,
		playerPt,
		playerAnglePt)
	shapes.ColorLine (
		destination,
		artist.Hex(0x00FF00FF), 1,
		playerPt,
		hitPt)
}
