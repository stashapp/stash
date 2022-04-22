package manager

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

type InteractiveHeatmapSpeedGenerator struct {
	InteractiveSpeed int
	Funscript        Script
	FunscriptPath    string
	HeatmapPath      string
	Width            int
	Height           int
	NumSegments      int
}

type Script struct {
	// Version of Launchscript
	Version string `json:"version"`
	// Inverted causes up and down movement to be flipped.
	Inverted bool `json:"inverted,omitempty"`
	// Range is the percentage of a full stroke to use.
	Range int `json:"range,omitempty"`
	// Actions are the timed moves.
	Actions      []Action `json:"actions"`
	AvarageSpeed int64
}

// Action is a move at a specific time.
type Action struct {
	// At time in milliseconds the action should fire.
	At int64 `json:"at"`
	// Pos is the place in percent to move to.
	Pos int `json:"pos"`

	Slope     float64
	Intensity int64
	Speed     float64
}

type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

func NewInteractiveHeatmapSpeedGenerator(funscriptPath string, heatmapPath string) *InteractiveHeatmapSpeedGenerator {
	return &InteractiveHeatmapSpeedGenerator{
		FunscriptPath: funscriptPath,
		HeatmapPath:   heatmapPath,
		Width:         320,
		Height:        15,
		NumSegments:   150,
	}
}

func (g *InteractiveHeatmapSpeedGenerator) Generate() error {
	funscript, err := g.LoadFunscriptData(g.FunscriptPath)

	if err != nil {
		return err
	}

	g.Funscript = funscript
	g.Funscript.UpdateIntensityAndSpeed()

	err = g.RenderHeatmap()

	if err != nil {
		return err
	}

	g.InteractiveSpeed = g.Funscript.CalculateMedian()

	return nil
}

func (g *InteractiveHeatmapSpeedGenerator) LoadFunscriptData(path string) (Script, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Script{}, err
	}

	var funscript Script
	err = json.Unmarshal(data, &funscript)
	if err != nil {
		return Script{}, err
	}

	if funscript.Actions == nil {
		return Script{}, fmt.Errorf("actions list missing in %s", path)
	}

	sort.SliceStable(funscript.Actions, func(i, j int) bool { return funscript.Actions[i].At < funscript.Actions[j].At })

	// trim actions with negative timestamps to avoid index range errors when generating heatmap

	isValid := func(x int64) bool { return x >= 0 }

	i := 0
	for _, x := range funscript.Actions {
		if isValid(x.At) {
			funscript.Actions[i] = x
			i++
		}
	}

	funscript.Actions = funscript.Actions[:i]

	return funscript, nil
}

func (funscript *Script) UpdateIntensityAndSpeed() {

	var t1, t2 int64
	var p1, p2 int
	var slope float64
	var intensity int64
	for i := range funscript.Actions {
		if i == 0 {
			continue
		}
		t1 = funscript.Actions[i].At
		t2 = funscript.Actions[i-1].At
		p1 = funscript.Actions[i].Pos
		p2 = funscript.Actions[i-1].Pos

		slope = math.Min(math.Max(1/(2*float64(t1-t2)/1000), 0), 20)
		intensity = int64(slope * math.Abs((float64)(p1-p2)))
		speed := math.Abs(float64(p1-p2)) / float64(t1-t2) * 1000

		funscript.Actions[i].Slope = slope
		funscript.Actions[i].Intensity = intensity
		funscript.Actions[i].Speed = speed
	}
}

// funscript needs to have intensity updated first
func (g *InteractiveHeatmapSpeedGenerator) RenderHeatmap() error {

	gradient := g.Funscript.getGradientTable(g.NumSegments)

	img := image.NewRGBA(image.Rect(0, 0, g.Width, g.Height))
	for x := 0; x < g.Width; x++ {
		c := gradient.GetInterpolatedColorFor(float64(x) / float64(g.Width))
		draw.Draw(img, image.Rect(x, 0, x+1, g.Height), &image.Uniform{c}, image.Point{}, draw.Src)
	}

	// add 10 minute marks
	maxts := g.Funscript.Actions[len(g.Funscript.Actions)-1].At
	const tick = 600000
	var ts int64 = tick
	c, _ := colorful.Hex("#000000")
	for ts < maxts {
		x := int(float64(ts) / float64(maxts) * float64(g.Width))
		draw.Draw(img, image.Rect(x-1, g.Height/2, x+1, g.Height), &image.Uniform{c}, image.Point{}, draw.Src)
		ts += tick
	}

	outpng, err := os.Create(g.HeatmapPath)
	if err != nil {
		return err
	}
	defer outpng.Close()

	err = png.Encode(outpng, img)
	return err
}

func (funscript *Script) CalculateMedian() int {
	sort.Slice(funscript.Actions, func(i, j int) bool {
		return funscript.Actions[i].Speed < funscript.Actions[j].Speed
	})

	mNumber := len(funscript.Actions) / 2

	if len(funscript.Actions)%2 != 0 {
		return int(funscript.Actions[mNumber].Speed)
	}

	return int((funscript.Actions[mNumber-1].Speed + funscript.Actions[mNumber].Speed) / 2)
}

func (gt GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

func (funscript Script) getGradientTable(numSegments int) GradientTable {
	segments := make([]struct {
		count     int
		intensity int
	}, numSegments)
	gradient := make(GradientTable, numSegments)

	maxts := funscript.Actions[len(funscript.Actions)-1].At

	for _, a := range funscript.Actions {
		segment := int(float64(a.At) / float64(maxts+1) * float64(numSegments))
		segments[segment].count++
		segments[segment].intensity += int(a.Intensity)
	}

	for i := 0; i < numSegments; i++ {
		gradient[i].Pos = float64(i) / float64(numSegments-1)
		if segments[i].count > 0 {
			gradient[i].Col = getSegmentColor(float64(segments[i].intensity) / float64(segments[i].count))
		} else {
			gradient[i].Col = getSegmentColor(0.0)
		}
	}

	return gradient
}

func getSegmentColor(intensity float64) colorful.Color {
	colorBlue, _ := colorful.Hex("#1e90ff")   // DodgerBlue
	colorGreen, _ := colorful.Hex("#228b22")  // ForestGreen
	colorYellow, _ := colorful.Hex("#ffd700") // Gold
	colorRed, _ := colorful.Hex("#dc143c")    // Crimson
	colorPurple, _ := colorful.Hex("#800080") // Purple
	colorBlack, _ := colorful.Hex("#0f001e")
	colorBackground, _ := colorful.Hex("#30404d") // Same as GridCard bg

	var stepSize = 60.0
	var f float64
	var c colorful.Color

	switch {
	case intensity <= 0.001:
		c = colorBackground
	case intensity <= 1*stepSize:
		f = (intensity - 0*stepSize) / stepSize
		c = colorBlue.BlendLab(colorGreen, f)
	case intensity <= 2*stepSize:
		f = (intensity - 1*stepSize) / stepSize
		c = colorGreen.BlendLab(colorYellow, f)
	case intensity <= 3*stepSize:
		f = (intensity - 2*stepSize) / stepSize
		c = colorYellow.BlendLab(colorRed, f)
	case intensity <= 4*stepSize:
		f = (intensity - 3*stepSize) / stepSize
		c = colorRed.BlendRgb(colorPurple, f)
	default:
		f = (intensity - 4*stepSize) / (5 * stepSize)
		f = math.Min(f, 1.0)
		c = colorPurple.BlendLab(colorBlack, f)
	}

	return c
}
