package res

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioContext *audio.Context

type SoundEffect struct {
	bytes []byte
}

func (s SoundEffect) Play() *audio.Player {
	p := audioContext.NewPlayerFromBytes(s.bytes)
	p.Play()
	return p
}

var Sounds = map[string]SoundEffect{}

func PlaySound(name string) *audio.Player {
	if _, ok := Sounds[name]; !ok {
		panic("sound not found")
	}
	return Sounds[name].Play()
}

func init() {
	audioContext = audio.NewContext(48000)

	files, err := FS.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println("hmm", f.Name())
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
