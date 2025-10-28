package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type InteractiveHeatmapSpeedGenerator struct {
	InteractiveSpeed int
	Funscript        Script
	Width            int
	Height           int
	NumSegments      int

	DrawRange bool
}

type Script struct {
	// Version of Launchscript
	// #5600 - ignore version, don't validate type
	Version json.RawMessage `json:"version"`
	// Inverted causes up and down movement to be flipped.
	Inverted bool `json:"inverted,omitempty"`
	// Range is the percentage of a full stroke to use.
	Range int `json:"range,omitempty"`
	// Actions are the timed moves.
	Actions []Action `json:"actions"`
}

// Action is a move at a specific time.
type Action struct {
	// At time in milliseconds the action should fire.
	At float64 `json:"at"`
	// Pos is the place in percent to move to.
	Pos int `json:"pos"`

	Speed float64
}

type GradientTable []struct {
	Col    colorful.Color
	Pos    float64
	YRange [2]float64
}

func NewInteractiveHeatmapSpeedGenerator(drawRange bool) *InteractiveHeatmapSpeedGenerator {
	return &InteractiveHeatmapSpeedGenerator{
		Width:       1280,
		Height:      60,
		NumSegments: 600,
		DrawRange:   drawRange,
	}
}

func (g *InteractiveHeatmapSpeedGenerator) Generate(funscriptPath string, heatmapPath string, sceneDuration float64) error {
	funscript, err := g.LoadFunscriptData(funscriptPath, sceneDuration)

	if err != nil {
		return err
	}

	if len(funscript.Actions) == 0 {
		return fmt.Errorf("no valid actions in funscript")
	}

	sceneDurationMilli := int64(sceneDuration * 1000)
	g.Funscript = funscript
	g.Funscript.UpdateIntensityAndSpeed()

	err = g.RenderHeatmap(heatmapPath, sceneDurationMilli)

	if err != nil {
		return err
	}

	g.InteractiveSpeed = g.Funscript.CalculateMedian()

	return nil
}

func (g *InteractiveHeatmapSpeedGenerator) LoadFunscriptData(path string, sceneDuration float64) (Script, error) {
	data, err := os.ReadFile(path)
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
	// #3181 - also trim actions that occur after the scene duration
	loggedBadTimestamp := false
	sceneDurationMilli := sceneDuration * 1000
	isValid := func(x float64) bool {
		return x >= 0 && x < sceneDurationMilli
	}

	i := 0
	for _, x := range funscript.Actions {
		if isValid(x.At) {
			funscript.Actions[i] = x
			i++
		} else if !loggedBadTimestamp {
			loggedBadTimestamp = true
			logger.Warnf("Invalid timestamp %d in %s: subsequent invalid timestamps will not be logged", x.At, path)
		}
	}

	funscript.Actions = funscript.Actions[:i]

	return funscript, nil
}

func (funscript *Script) UpdateIntensityAndSpeed() {

	var t1, t2 float64
	var p1, p2 int
	var intensity float64
	for i := range funscript.Actions {
		if i == 0 {
			continue
		}
		t1 = funscript.Actions[i].At
		t2 = funscript.Actions[i-1].At
		p1 = funscript.Actions[i].Pos
		p2 = funscript.Actions[i-1].Pos

		speed := math.Abs(float64(p1 - p2))
		intensity = float64(speed/float64(t1-t2)) * 1000

		funscript.Actions[i].Speed = intensity
	}
}

// funscript needs to have intensity updated first
func (g *InteractiveHeatmapSpeedGenerator) RenderHeatmap(heatmapPath string, sceneDurationMilli int64) error {
	gradient := g.Funscript.getGradientTable(g.NumSegments, sceneDurationMilli)

	img := image.NewRGBA(image.Rect(0, 0, g.Width, g.Height))
	for x := 0; x < g.Width; x++ {
		xPos := float64(x) / float64(g.Width)
		c := gradient.GetInterpolatedColorFor(xPos)

		y0 := 0
		y1 := g.Height

		if g.DrawRange {
			yRange := gradient.GetYRange(xPos)
			top := int(yRange[0] / 100.0 * float64(g.Height))
			bottom := int(yRange[1] / 100.0 * float64(g.Height))

			y0 = g.Height - top
			y1 = g.Height - bottom
		}

		draw.Draw(img, image.Rect(x, y0, x+1, y1), &image.Uniform{c}, image.Point{}, draw.Src)
	}

	// add 10 minute marks
	maxts := sceneDurationMilli
	const tick = 600000
	var ts int64 = tick
	c, _ := colorful.Hex("#000000")
	for ts < maxts {
		x := int(float64(ts) / float64(maxts) * float64(g.Width))
		draw.Draw(img, image.Rect(x-1, g.Height/2, x+1, g.Height), &image.Uniform{c}, image.Point{}, draw.Src)
		ts += tick
	}

	outpng, err := os.Create(heatmapPath)
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

func (gt GradientTable) GetYRange(t float64) [2]float64 {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// TODO: We are in between c1 and c2. Go blend them!
			return c1.YRange
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].YRange
}

func (funscript Script) getGradientTable(numSegments int, sceneDurationMilli int64) GradientTable {
	const windowSize = 15
	const backfillThreshold = float64(500)

	segments := make([]struct {
		count     int
		intensity int
		yRange    [2]float64
		at        float64
	}, numSegments)
	gradient := make(GradientTable, numSegments)
	posList := []int{}

	maxts := sceneDurationMilli

	for _, a := range funscript.Actions {
		posList = append(posList, a.Pos)

		if len(posList) > windowSize {
			posList = posList[1:]
		}

		sortedPos := make([]int, len(posList))
		copy(sortedPos, posList)
		sort.Ints(sortedPos)

		topHalf := sortedPos[len(sortedPos)/2:]
		bottomHalf := sortedPos[0 : len(sortedPos)/2]

		var totalBottom int
		var totalTop int

		for _, value := range bottomHalf {
			totalBottom += value
		}
		for _, value := range topHalf {
			totalTop += value
		}

		averageBottom := float64(totalBottom) / float64(len(bottomHalf))
		averageTop := float64(totalTop) / float64(len(topHalf))

		segment := int(float64(a.At) / float64(maxts+1) * float64(numSegments))
		// #3181 - sanity check. Clamp segment to numSegments-1
		if segment >= numSegments {
			segment = numSegments - 1
		}
		segments[segment].at = a.At
		segments[segment].count++
		segments[segment].intensity += int(a.Speed)
		segments[segment].yRange[0] = averageTop
		segments[segment].yRange[1] = averageBottom
	}

	lastSegment := segments[0]

	// Fill in gaps in segments
	for i := 0; i < numSegments; i++ {
		segmentTS := float64((maxts / int64(numSegments)) * int64(i))

		// Empty segment - fill it with the previous up to backfillThreshold ms
		if segments[i].count == 0 {
			if segmentTS-lastSegment.at < backfillThreshold {
				segments[i].count = lastSegment.count
				segments[i].intensity = lastSegment.intensity
				segments[i].yRange[0] = lastSegment.yRange[0]
				segments[i].yRange[1] = lastSegment.yRange[1]
			}
		} else {
			lastSegment = segments[i]
		}
	}

	for i := 0; i < numSegments; i++ {
		gradient[i].Pos = float64(i) / float64(numSegments-1)
		gradient[i].YRange = segments[i].yRange
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

	var stepSize = 125.0
	var f float64
	var c colorful.Color

	switch {
	case intensity <= 25:
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

func LoadFunscriptData(path string) (Script, error) {
	data, err := os.ReadFile(path)
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

	return funscript, nil
}

func convertRange(value int, fromLow int, fromHigh int, toLow int, toHigh int) int {
	return ((value-fromLow)*(toHigh-toLow))/(fromHigh-fromLow) + toLow
}

func ConvertFunscriptToCSV(funscriptPath string) ([]byte, error) {
	funscript, err := LoadFunscriptData(funscriptPath)

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	for _, action := range funscript.Actions {
		pos := action.Pos

		if funscript.Inverted {
			pos = convertRange(pos, 0, 100, 100, 0)
		}

		if funscript.Range > 0 {
			pos = convertRange(pos, 0, funscript.Range, 0, 100)
		}

		// I don't know whether the csv format requires int or float, so for now we'll use int
		buffer.WriteString(fmt.Sprintf("%d,%d\r\n", int(math.Round(action.At)), pos))
	}
	return buffer.Bytes(), nil
}

func ConvertFunscriptToCSVFile(funscriptPath string, csvPath string) error {
	csvBytes, err := ConvertFunscriptToCSV(funscriptPath)

	if err != nil {
		return err
	}

	return fsutil.WriteFile(csvPath, csvBytes)
}
