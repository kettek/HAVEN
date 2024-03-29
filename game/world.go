package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
	"github.com/tinne26/etxt"
)

type World struct {
	PlayerActor      Actor
	Rooms            []*Room
	LastRoom         *Room
	Room             *Room
	Camera           *Camera
	RoutineChan      chan func() bool
	RoutineChans     []func() bool
	Messages         []Message
	Prompts          []*Prompt
	Combat           *Combat
	roomBuilder      func(string) *Room
	Color            color.NRGBA
	colorTicker      int
	postProcessImage *ebiten.Image
	SkipMessages     bool
}

func NewWorld(roomBuilder func(string) *Room) *World {
	return &World{
		Camera:      NewCamera(),
		RoutineChan: make(chan func() bool),
		roomBuilder: roomBuilder,
	}
}

func (w *World) Update() {
	// Process routines.
	for {
		select {
		case fnc := <-w.RoutineChan:
			w.RoutineChans = append(w.RoutineChans, fnc)
		default:
			goto done
		}
	}
done:
	routines := w.RoutineChans[:0]
	for _, r := range w.RoutineChans {
		if !r() {
			routines = append(routines, r)
		}
	}
	w.RoutineChans = routines

	if len(w.Prompts) != 0 {
		w.Prompts[len(w.Prompts)-1].Update()
	}

	if w.Combat != nil {
		if c := w.Combat.Update(w, w.Room); c != nil {
			if cmd, ok := c.(commands.CombatResult); ok {
				px, py, _ := cmd.Winner.(Actor).Position()
				if cmd.Fled {
					// ... nada
				} else if cmd.Winner == w.PlayerActor {
					exp := cmd.ExpGained
					if !cmd.Destroyed {
						exp /= 2 // Half exp for capturing.
						// TODO: Capture glitch.
						w.PlayerActor.(CombatActor).AddGlitch(cmd.Loser.(GlitchActor))
						w.Room.TileMessage(Message{
							X:        px,
							Y:        py,
							Text:     fmt.Sprintf("caught %s", cmd.Loser.(GlitchActor).Name()),
							Color:    color.NRGBA{255, 64, 255, 255},
							Duration: 3 * time.Second,
						})
					}
					// Give EXP
					lvl := cmd.Winner.(CombatActor).AddExp(exp)
					w.Room.TileMessage(Message{
						X:        px,
						Y:        py,
						Text:     fmt.Sprintf("+%d EXP", exp),
						Color:    color.NRGBA{255, 255, 0, 255},
						Duration: 3 * time.Second,
					})
					if lvl > 0 {
						w.Room.TileMessage(Message{
							X:        px,
							Y:        py,
							Text:     fmt.Sprintf("+%d LVL(s)", lvl),
							Color:    color.NRGBA{0, 255, 0, 255},
							Duration: 3 * time.Second,
						})
					}
					w.Room.RemoveActor(cmd.Loser.(Actor))
					w.Room.UpdateGlitchion()
					if w.Room.Glitches == 0 {
						w.Room.Darkness = 0
						res.PlaySound("cleansed")
					}
				} else if cmd.Loser == w.PlayerActor {
					// Penalize the player in each stat by the level of the winner.
					lvl := cmd.Winner.(CombatActor).Level()
					cmd.Loser.(CombatActor).Penalize(lvl, lvl, lvl)
					w.Room.TileMessage(Message{
						X:        px,
						Y:        py,
						Text:     "*bzzt*",
						Color:    color.NRGBA{255, 0, 255, 150},
						Duration: 2 * time.Second,
					})
					px, py, _ = cmd.Loser.(Actor).Position()
					w.Room.TileMessage(Message{
						X:        px,
						Y:        py,
						Text:     fmt.Sprintf("-%d STATS", lvl),
						Color:    color.NRGBA{255, 0, 0, 255},
						Duration: 2 * time.Second,
					})
					w.Room.RemoveActor(cmd.Winner.(Actor))
				}
				w.Combat = nil
				res.Jukebox.Play(w.Room.Song)
			} else if cmd, ok := c.(commands.Prompt); ok {
				w.AddPrompt(cmd.Items, cmd.Message, cmd.Handler, cmd.ShowVersions)
			}
		}
	}

	if w.PlayerActor != nil {
		x, y, _ := w.PlayerActor.Position()
		geom, _ := w.Room.GetTilePositionGeoM(x, y)
		w.Camera.MoveTo(geom.Element(0, 2), geom.Element(1, 2))
	}

	// Process camera.
	w.Camera.Update()

	w.colorTicker++
	if w.LastRoom != nil && w.Room != nil && w.colorTicker < 30 {
		w.Color = color.NRGBA{
			uint8(float64(w.LastRoom.Color.R)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.R)*(float64(w.colorTicker)/30)),
			uint8(float64(w.LastRoom.Color.G)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.G)*(float64(w.colorTicker)/30)),
			uint8(float64(w.LastRoom.Color.B)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.B)*(float64(w.colorTicker)/30)),
			255,
		}
	} else if w.Room != nil {
		w.Color = w.Room.Color
	} else {
		w.Color = color.NRGBA{0, 0, 0, 255}
	}

	// Process room.
	if w.Room != nil {
		cmds := w.Room.Update(w)
		for _, cmd := range cmds {
			switch cmd := cmd.(type) {
			case commands.Prompt:
				w.AddPrompt(cmd.Items, "", cmd.Handler, cmd.ShowVersions)
			case commands.Travel:
				room := w.roomBuilder(cmd.Room)
				var targetActor Actor
				if cmd.Target != nil {
					targetActor = cmd.Target.(Actor)
				} else if w.PlayerActor != nil {
					targetActor = w.PlayerActor
				}
				var x, y int
				if actor := room.GetActorByTag(cmd.Tag); actor != nil {
					x, y, _ = actor.Position()
					x += cmd.OffsetX
					y += cmd.OffsetY
				}
				if targetActor != nil {
					w.Room.RemoveActor(targetActor)
					room.PrependActor(targetActor)
					targetActor.SetPosition(x, y, 0)
				}
				w.EnterRoom(room)
			case commands.Combat:
				attacker := cmd.Attacker.(CombatActor)
				defender := cmd.Defender.(CombatActor)
				//w.Combat = NewCombat(384, 288, attacker, defender)
				w.Combat = NewCombat(500, 388, attacker, defender)
				res.Jukebox.Play("bad-health")
			default:
				fmt.Println("unhandled room->world command", cmd)
			}
		}
	}
}

func (w *World) Input(in inputs.Input) {
	if len(w.Prompts) > 0 {
		w.Prompts[len(w.Prompts)-1].Input(in)
	} else if w.Combat != nil {
		w.Combat.Input(in)
	} else {
		if !w.Room.Input(w, in) {
			switch in := in.(type) {
			case inputs.Key:
				// glitch 1-9 key selection.
				for i := int(ebiten.KeyDigit1); i <= int(ebiten.KeyDigit9); i++ {
					if in.Key != ebiten.Key(i) {
						continue
					}
					if w.PlayerActor != nil {
						glitches := w.PlayerActor.(CombatActor).Glitches()
						j := i - int(ebiten.KeyDigit1)
						if j < len(glitches) {
							w.PlayerActor.(CombatActor).SetGlitch(glitches[j])
						}
					}
				}
			case inputs.Click:
				x, y := int(in.X), int(in.Y)
				// Check for glitch select, absorb, etc.
				if x >= glitchesUIX && x <= glitchesUIX+glitchesUIWidth && y >= glitchesUIY && y <= glitchesUIY+glitchesUIHeight {
					gx := (x - glitchesUIX) / 16
					if w.PlayerActor != nil {
						glitches := w.PlayerActor.(CombatActor).Glitches()
						if gx < len(glitches) {
							w.PlayerActor.(CombatActor).SetGlitch(glitches[gx])
						}
					}
				} else if x >= glitchesUIAbsorbX && x <= glitchesUIAbsorbX+glitchesUIAbsorbWidth && y >= glitchesUIY && y <= glitchesUIY+glitchesUIHeight {
					if w.PlayerActor != nil {
						if glitch := w.PlayerActor.(CombatActor).CurrentGlitch(); glitch != nil {
							w.AddPrompt([]string{"CANCEL", "OK"}, "Absorbing will remove this glitch. It will also also apply a boost to you using its stats.\nAbsorb?", func(i int, s string) bool {
								if i == 0 {
									return true
								} else if i == 1 {
									// Restore stats so the glitch gives a minimal boost. :)
									//glitch.(CombatActor).RestoreStats() // actually, no, having an injured glitch could be a cool tactic -- bring something real low, capture it, then use it for a self-heal.
									p, f, i := glitch.RollBoost()
									w.PlayerActor.(CombatActor).ApplyBoost(p, f, i)
									w.PlayerActor.(CombatActor).RemoveGlitch(glitch)
									res.PlaySound("slurp")
									px, py, _ := w.PlayerActor.Position()
									w.Room.TileMessage(Message{
										X:        px,
										Y:        py,
										Text:     fmt.Sprintf("+%dPEN +%dFIRE +%dINT", p, f, i),
										Color:    color.NRGBA{0, 255, 255, 255},
										Duration: 3 * time.Second,
									})
									return true
								}
								return true
							}, false)
						}
					}
				} else if x >= glitchesUIInfoX && x <= glitchesUIInfoX+glitchesUIInfoWidth && y >= glitchesUIY && y <= glitchesUIY+glitchesUIHeight {
					if w.PlayerActor != nil {
						glitch := w.PlayerActor.(CombatActor).CurrentGlitch()
						if glitch != nil {
							info := fmt.Sprintf("%s (LVL %d)\n\n", glitch.Name(), glitch.Level())
							p, f, i := glitch.(CombatActor).CurrentStats()
							mp, mf, mi := glitch.(CombatActor).MaxStats()
							ability := "-"
							abilityStuff := ""
							abilityDesc := ""
							if glitch.Ability() != nil {
								ability = glitch.Ability().Name + " " + fmt.Sprintf("(%d)", glitch.Ability().Tier)
								abilityStuff = fmt.Sprintf("         TURNS: %d | COOLDOWN: %d", glitch.Ability().Turns, glitch.Ability().Cooldown)
								abilityDesc = string(AbilityDescriptions[AbilityType(glitch.Ability().Name)])
							}
							info += fmt.Sprintf("INTEGRITY %d/%d\n", i, mi)
							info += fmt.Sprintf("FIREWALL %d/%d\n", f, mf)
							info += fmt.Sprintf("PENETRATION %d/%d\n", p, mp)
							info += fmt.Sprintf("ABILITY: %s\n", ability)
							info += fmt.Sprintf("%s\n", abilityStuff)
							info += abilityDesc
							w.AddPrompt([]string{"OK"}, info, func(i int, s string) bool {
								return true
							}, false)
						}
					}
				}
			}
		}
	}
}

// Draw the room, combat, overlays, etc. Don't code like this. :)
func (w *World) Draw(screen *ebiten.Image) {
	if w.postProcessImage == nil {
		w.postProcessImage = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}
	screen.Fill(w.Color)

	geom := ebiten.GeoM{}

	geom.Translate(-w.Camera.W/2, -w.Camera.H/2)
	geom.Rotate(w.Camera.Rotation)
	geom.Translate(w.Camera.W/2, w.Camera.H/2)
	geom.Translate(-w.Camera.X+w.Camera.W/2, -w.Camera.Y+w.Camera.H/2)
	geom.Scale(w.Camera.Zoom, w.Camera.Zoom)
	if w.Room != nil {
		w.Room.Draw(screen, geom)
	}

	w.postProcessImage.Clear()
	if w.Room != nil {
		w.Room.DrawPost(screen, w.postProcessImage, geom)
	}
	screen.DrawImage(w.postProcessImage, nil)

	res.Text.Utils().StoreState()
	if len(w.Messages) > 0 {
		m := w.Messages[0]
		mh := float32(screen.Bounds().Dy()/4) * float32(m.H)
		my := float32(screen.Bounds().Dy())/3.5 - mh/2
		res.Text.SetColor(m.Color)
		if m.Background.A != 0 {
			vector.DrawFilledRect(screen, 0, my, float32(screen.Bounds().Dx()), mh, m.Background, false)
		}
		res.Text.SetAlign(etxt.Center)
		res.Text.DrawWithWrap(screen, m.Text, screen.Bounds().Dx()/2, int(my+mh/2), screen.Bounds().Dx()-32)
	}
	res.Text.Utils().RestoreState()

	if w.Combat != nil {
		geom := ebiten.GeoM{}
		w.Combat.x = float64(screen.Bounds().Dx()/2) - float64(w.Combat.image.Bounds().Dx()/2)
		w.Combat.y = float64(screen.Bounds().Dy()/2) - float64(w.Combat.image.Bounds().Dy()/2)
		geom.Translate(w.Combat.x, w.Combat.y)
		w.Combat.Draw(screen, geom)
	} else if w.PlayerActor != nil {
		p, f, i := w.PlayerActor.(CombatActor).CurrentStats()
		mp, mf, mi := w.PlayerActor.(CombatActor).MaxStats()
		exp := w.PlayerActor.(CombatActor).Exp()
		lvl := w.PlayerActor.(CombatActor).Level()
		// Draw player UI
		x := 6
		y := screen.Bounds().Dy() - 90

		vector.DrawFilledRect(screen, float32(x), float32(y), playerUIWidth, playerUIHeight, color.NRGBA{19, 19, 97, 200}, false)
		vector.StrokeRect(screen, float32(x), float32(y), playerUIWidth, playerUIHeight, 3, color.NRGBA{194, 193, 174, 255}, true)
		x += 4
		y += 3

		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		res.Text.SetColor(color.NRGBA{255, 255, 255, 255})
		res.Text.SetAlign(etxt.Right | etxt.Top)
		res.Text.Draw(screen, w.PlayerActor.Name(), 183, y)
		res.Text.SetAlign(etxt.Left | etxt.Top)
		res.Text.SetColor(color.NRGBA{128, 128, 0, 255})
		res.Text.Draw(screen, fmt.Sprintf("LVL %d", lvl), x, y)
		y += 16
		res.Text.SetColor(color.NRGBA{0, 128, 128, 255})
		res.Text.Draw(screen, fmt.Sprintf("EXP %.3d/100", exp), x, y)
		y += 16
		res.Text.SetColor(color.NRGBA{50, 255, 50, 200})
		res.Text.Draw(screen, fmt.Sprintf("INTEGRITY %d/%d", i, mi), x, y)
		y += 16
		res.Text.SetColor(color.NRGBA{255, 50, 50, 200})
		res.Text.Draw(screen, fmt.Sprintf("FIREWALL %d/%d", f, mf), x, y)
		y += 16
		res.Text.SetColor(color.NRGBA{255, 255, 50, 200})
		res.Text.Draw(screen, fmt.Sprintf("PENETRATION %d/%d", p, mp), x, y)

		// Draw our lil glitchy bois if we got 'em
		glitches := w.PlayerActor.(CombatActor).Glitches()
		if len(glitches) > 0 {
			y := y - 1
			x := x + 180 + 4
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(glitchesUIWidth), glitchesUIHeight, color.NRGBA{19, 19, 97, 200}, false)
			vector.StrokeRect(screen, float32(x), float32(y), float32(glitchesUIWidth), glitchesUIHeight, 3, color.NRGBA{194, 193, 174, 255}, true)
			x += 4
			y += 3
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			glitchesUIX = x
			glitchesUIY = y
			for _, g := range glitches {
				if g == w.PlayerActor.(CombatActor).CurrentGlitch() {
					g.SpriteStack().Highlight = true
				}
				g.SpriteStack().DrawFlat(screen, op.GeoM)
				g.SpriteStack().Highlight = false
				op.GeoM.Translate(16, 0)
				x += 16
			}
			y -= 3
			x = glitchesUIX + glitchesUIWidth + paddingUI
			// Draw an ABSORB button.
			res.Text.SetAlign(etxt.Center)
			glitchesUIAbsorbWidth = res.Text.Measure("ABSORB GLITCH").IntWidth() + paddingUI*2
			glitchesUIAbsorbX = x
			res.Text.SetColor(color.NRGBA{19, 19, 94, 200})
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(glitchesUIAbsorbWidth), glitchesUIHeight, color.NRGBA{194, 193, 174, 200}, false)
			vector.StrokeRect(screen, float32(x), float32(y), float32(glitchesUIAbsorbWidth), glitchesUIHeight, 3, color.NRGBA{19, 19, 94, 255}, true)
			res.Text.Draw(screen, "ABSORB GLITCH", x+glitchesUIAbsorbWidth/2, y+glitchesUIHeight/2+1)
			// Draw an INFO button.
			x += glitchesUIAbsorbWidth
			x += paddingUI * 2
			glitchesUIInfoX = x
			glitchesUIInfoWidth = res.Text.Measure("GLITCH INFO").IntWidth() + paddingUI*2
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(glitchesUIInfoWidth), glitchesUIHeight, color.NRGBA{194, 193, 174, 200}, false)
			vector.StrokeRect(screen, float32(x), float32(y), float32(glitchesUIInfoWidth), glitchesUIHeight, 3, color.NRGBA{19, 19, 94, 255}, true)
			res.Text.Draw(screen, "GLITCH INFO", x+glitchesUIInfoWidth/2, y+glitchesUIHeight/2+1)
		}

		// Draw map UI
		x = screen.Bounds().Dx() - int(mapUIWidth) - 6
		y = 6

		fg := w.Room.Color
		// bg is inverse of fg with wrap around.
		bg := color.NRGBA{255 - fg.R, 255 - fg.G, 255 - fg.B, 255}

		vector.DrawFilledRect(screen, float32(x), float32(y), mapUIWidth, mapUIHeight, bg, false)
		vector.StrokeRect(screen, float32(x), float32(y), mapUIWidth, mapUIHeight, 3, fg, true)
		y += 3

		res.Text.SetColor(fg)
		res.Text.SetAlign(etxt.Right | etxt.Top)
		res.Text.Draw(screen, w.Room.Name, x+int(mapUIWidth)-4, y)
		y += 16
		res.Text.SetFont(res.SmallFont.Font)
		res.Text.SetSize(float64(res.SmallFont.Size))
		res.Text.Draw(screen, w.Room.Song, x+int(mapUIWidth)-4, y)
		res.Text.SetAlign(etxt.Left | etxt.Top)
		y += 12
		if w.Room.MaxGlitches == 0 {
			fg.R /= 2
			fg.G *= 2
			res.Text.SetColor(fg)
			res.Text.Draw(screen, "safe", x+4, y)
		} else if w.Room.Glitches <= 0 {
			fg.R /= 2
			fg.G *= 2
			res.Text.SetColor(fg)
			res.Text.Draw(screen, "it is cleansed", x+4, y)
		} else {
			fg.R *= 2
			fg.G /= 2
			res.Text.SetColor(fg)
			res.Text.Draw(screen, fmt.Sprintf("%d remain", w.Room.Glitches), x+4, y)
		}

		res.Text.Utils().RestoreState()
	}

	if len(w.Prompts) != 0 {
		geom := ebiten.GeoM{}
		prompt := w.Prompts[len(w.Prompts)-1]
		prompt.x = float64(screen.Bounds().Dx()/2) - float64(prompt.image.Bounds().Dx()/2)
		prompt.y = float64(screen.Bounds().Dy()/2) - float64(prompt.image.Bounds().Dy()/2)
		geom.Translate(prompt.x, prompt.y)
		prompt.Draw(screen, geom)
	}
}

func (w *World) EnterRoom(room *Room) {
	go func() {
		<-w.FuncR(func() {
			if w.Room != nil {
				if w.Room.OnLeave != nil {
					w.Room.OnLeave(w, w.Room)
				}
				w.Room.active = false
			}
		})
		<-w.FuncR(func() {
			w.LastRoom = w.Room
			w.Room = room
			if w.LastRoom != nil {
				w.Room.DrawMode = w.LastRoom.DrawMode
			}
			if w.Room.Song != "" {
				res.Jukebox.Play(w.Room.Song)
			}
			w.Room.UpdateGlitchion()
			w.colorTicker = 0
			if w.PlayerActor != nil {
				x, y, _ := w.PlayerActor.Position()
				geom, _ := w.Room.GetTilePositionGeoM(x, y)
				w.Camera.SetPosition(geom.Element(0, 2), geom.Element(1, 2))
			}

		})
		if w.Room.OnEnter != nil {
			w.Room.OnEnter(w, w.Room)
		}
		<-w.FuncR(func() {
			w.Room.Activate()
			if w.Room.Glitches > 0 {
				res.PlaySound("glitched")
			} else {
				res.PlaySound("cleansed")
			}
		})
	}()
}

func (w *World) AddPrompt(items []string, msg string, cb func(int, string) bool, showExtra bool) {
	w.Prompts = append(w.Prompts, NewPrompt(320, 200, items, msg, func(i int, s string) bool {
		done := false
		if i == -1 {
			done = true
			cb(i, s)
		} else {
			done = cb(i, s)
		}
		if done {
			if len(w.Prompts) >= 1 {
				w.Prompts = w.Prompts[:len(w.Prompts)-1]
			}
		}
		return done
	}, showExtra))
}

func (w *World) MessageR(msg Message) chan bool {
	done := make(chan bool)

	msg.start = time.Now()
	msg.id = messageID
	if msg.Color.A == 0 {
		msg.Color = color.NRGBA{0, 0, 0, 255}
	}
	if msg.Font == nil {
		msg.Font = &res.DefFont
	}
	messageID++

	first := true
	fnc := func() bool {
		if first {
			w.Messages = append(w.Messages, msg)
			if w.SkipMessages {
				msg.Duration = 0
			}
			first = false
		}

		delta := time.Since(msg.start)
		if delta < 200*time.Millisecond {
			w.Messages[0].Color.A = uint8(float64(delta) / float64(200*time.Millisecond) * 200)
			w.Messages[0].Background.A = uint8(float64(delta) / float64(200*time.Millisecond) * 200)
			w.Messages[0].H = float64(delta) / float64(200*time.Millisecond)
		} else if delta > msg.Duration-200*time.Millisecond {
			w.Messages[0].Color.A = uint8(float64(msg.Duration-delta) / float64(200*time.Millisecond) * 200)
			w.Messages[0].Background.A = uint8(float64(msg.Duration-delta) / float64(200*time.Millisecond) * 200)
			w.Messages[0].H = float64(msg.Duration-delta) / float64(200*time.Millisecond)
		}

		if time.Since(msg.start) >= msg.Duration {
			w.Messages = w.Messages[1:]
			done <- true
			return true
		}
		return false
	}

	w.RoutineChan <- fnc
	return done
}

func (w *World) FuncR(fnc func()) chan bool {
	done := make(chan bool)
	f := func() bool {
		fnc()
		done <- true
		return true
	}
	w.RoutineChan <- f
	return done
}

const playerUIHeight = 84
const playerUIWidth = 180
const playerUIPadding = 4 // maybe?
const paddingUI = 4

var glitchesUIX = 0
var glitchesUIY = 0

const glitchesUIHeight = 18
const glitchesUIWidth = 16*9 + 2

var glitchesUIAbsorbX = 0
var glitchesUIAbsorbWidth = 0
var glitchesUIInfoX = 0
var glitchesUIInfoWidth = 0

const mapUIHeight = 43
const mapUIWidth = 150
