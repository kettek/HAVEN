package res

import (
	"bytes"
	"io"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioContext *audio.Context

type SoundPlayer struct {
	*audio.Player
	Looping bool
	Next    *SoundPlayer
}

type SoundEffect struct {
	bytes []byte
}

func (s SoundEffect) Play() *SoundPlayer {
	p := audioContext.NewPlayerFromBytes(s.bytes)
	p.Play()
	sp := &SoundPlayer{
		Player:  p,
		Looping: false,
	}
	PlayingSounds = append(PlayingSounds, sp)
	return sp
}

func (s SoundEffect) PlayLooped() *SoundPlayer {
	sp := s.Play()
	sp.Looping = true
	return sp
}

var Sounds = map[string]SoundEffect{}
var PlayingSounds []*SoundPlayer

func GetSound(name string) *SoundPlayer {
	if _, ok := Sounds[name]; !ok {
		panic("sound not found")
	}
	p := audioContext.NewPlayerFromBytes(Sounds[name].bytes)
	sp := &SoundPlayer{
		Player:  p,
		Looping: false,
	}
	return sp
}

func PlaySound(name string) *SoundPlayer {
	if _, ok := Sounds[name]; !ok {
		panic("sound not found")
	}
	return Sounds[name].Play()
}

func PlayLoopedSound(name string) *SoundPlayer {
	if _, ok := Sounds[name]; !ok {
		panic("sound not found")
	}
	return Sounds[name].PlayLooped()
}

func UpdateSounds() {
	sounds := PlayingSounds[:0]
	for _, s := range PlayingSounds {
		if s.IsPlaying() {
			sounds = append(sounds, s)
			continue
		} else {
			if s.Looping {
				s.SetPosition(0)
				s.Play()
				sounds = append(sounds, s)
			} else if s.Next != nil {
				s.Next.Play()
				sounds = append(sounds, s.Next)
			}
		}
	}
	PlayingSounds = sounds
}

func init() {
	audioContext = audio.NewContext(48000)

	files, err := FS.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".wav") {
			b, err := FS.ReadFile(f.Name())
			if err != nil {
				panic(err)
			}
			s, err := wav.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(b))
			if err != nil {
				panic(err)
			}
			wavBytes, err := io.ReadAll(s)
			if err != nil {
				panic(err)
			}
			name, _ := strings.CutSuffix(f.Name(), ".wav")
			Sounds[name] = SoundEffect{bytes: wavBytes}
		}
	}
}
