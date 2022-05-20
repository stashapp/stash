package astikit

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// PCMLevel computes the PCM level of samples
// https://dsp.stackexchange.com/questions/2951/loudness-of-pcm-stream
// https://dsp.stackexchange.com/questions/290/getting-loudness-of-a-track-with-rms?noredirect=1&lq=1
func PCMLevel(samples []int) float64 {
	// Compute sum of square values
	var sum float64
	for _, s := range samples {
		sum += math.Pow(float64(s), 2)
	}

	// Square root
	return math.Sqrt(sum / float64(len(samples)))
}

func maxPCMSample(bitDepth int) int {
	return int(math.Pow(2, float64(bitDepth))/2.0) - 1
}

// PCMNormalize normalizes the PCM samples
func PCMNormalize(samples []int, bitDepth int) (o []int) {
	// Get max sample
	var m int
	for _, s := range samples {
		if v := int(math.Abs(float64(s))); v > m {
			m = v
		}
	}

	// Get max for bit depth
	max := maxPCMSample(bitDepth)

	// Loop through samples
	for _, s := range samples {
		o = append(o, s*max/m)
	}
	return
}

// ConvertPCMBitDepth converts the PCM bit depth
func ConvertPCMBitDepth(srcSample int, srcBitDepth, dstBitDepth int) (dstSample int, err error) {
	// Nothing to do
	if srcBitDepth == dstBitDepth {
		dstSample = srcSample
		return
	}

	// Convert
	if srcBitDepth < dstBitDepth {
		dstSample = srcSample << uint(dstBitDepth-srcBitDepth)
	} else {
		dstSample = srcSample >> uint(srcBitDepth-dstBitDepth)
	}
	return
}

// PCMSampleFunc is a func that can process a sample
type PCMSampleFunc func(s int) error

// PCMSampleRateConverter is an object capable of converting a PCM's sample rate
type PCMSampleRateConverter struct {
	b                    [][]int
	dstSampleRate        int
	fn                   PCMSampleFunc
	numChannels          int
	numChannelsProcessed int
	numSamplesOutputed   int
	numSamplesProcessed  int
	srcSampleRate        int
}

// NewPCMSampleRateConverter creates a new PCMSampleRateConverter
func NewPCMSampleRateConverter(srcSampleRate, dstSampleRate, numChannels int, fn PCMSampleFunc) *PCMSampleRateConverter {
	return &PCMSampleRateConverter{
		b:             make([][]int, numChannels),
		dstSampleRate: dstSampleRate,
		fn:            fn,
		numChannels:   numChannels,
		srcSampleRate: srcSampleRate,
	}
}

// Reset resets the converter
func (c *PCMSampleRateConverter) Reset() {
	c.b = make([][]int, c.numChannels)
	c.numChannelsProcessed = 0
	c.numSamplesOutputed = 0
	c.numSamplesProcessed = 0
}

// Add adds a new sample to the converter
func (c *PCMSampleRateConverter) Add(i int) (err error) {
	// Forward sample
	if c.srcSampleRate == c.dstSampleRate {
		if err = c.fn(i); err != nil {
			err = fmt.Errorf("astikit: handling sample failed: %w", err)
			return
		}
		return
	}

	// Increment num channels processed
	c.numChannelsProcessed++

	// Reset num channels processed
	if c.numChannelsProcessed > c.numChannels {
		c.numChannelsProcessed = 1
	}

	// Only increment num samples processed if all channels have been processed
	if c.numChannelsProcessed == c.numChannels {
		c.numSamplesProcessed++
	}

	// Append sample to buffer
	c.b[c.numChannelsProcessed-1] = append(c.b[c.numChannelsProcessed-1], i)

	// Throw away data
	if c.srcSampleRate > c.dstSampleRate {
		// Make sure to always keep the first sample but do nothing until we have all channels or target sample has been
		// reached
		if (c.numSamplesOutputed > 0 && float64(c.numSamplesProcessed) < 1.0+float64(c.numSamplesOutputed)*float64(c.srcSampleRate)/float64(c.dstSampleRate)) || c.numChannelsProcessed < c.numChannels {
			return
		}

		// Loop through channels
		for idx, b := range c.b {
			// Merge samples
			var s int
			for _, v := range b {
				s += v
			}
			s /= len(b)

			// Reset buffer
			c.b[idx] = []int{}

			// Custom
			if err = c.fn(s); err != nil {
				err = fmt.Errorf("astikit: handling sample failed: %w", err)
				return
			}
		}

		// Increment num samples outputted
		c.numSamplesOutputed++
		return
	}

	// Do nothing until we have all channels
	if c.numChannelsProcessed < c.numChannels {
		return
	}

	// Repeat data
	for c.numSamplesOutputed == 0 || float64(c.numSamplesProcessed)+1.0 > 1.0+float64(c.numSamplesOutputed)*float64(c.srcSampleRate)/float64(c.dstSampleRate) {
		// Loop through channels
		for _, b := range c.b {
			// Invalid length
			if len(b) != 1 {
				err = fmt.Errorf("astikit: invalid buffer item length %d", len(b))
				return
			}

			// Custom
			if err = c.fn(b[0]); err != nil {
				err = fmt.Errorf("astikit: handling sample failed: %w", err)
				return
			}
		}

		// Increment num samples outputted
		c.numSamplesOutputed++
	}

	// Reset buffer
	c.b = make([][]int, c.numChannels)
	return
}

// PCMChannelsConverter is an object of converting PCM's channels
type PCMChannelsConverter struct {
	dstNumChannels int
	fn             PCMSampleFunc
	srcNumChannels int
	srcSamples     int
}

// NewPCMChannelsConverter creates a new PCMChannelsConverter
func NewPCMChannelsConverter(srcNumChannels, dstNumChannels int, fn PCMSampleFunc) *PCMChannelsConverter {
	return &PCMChannelsConverter{
		dstNumChannels: dstNumChannels,
		fn:             fn,
		srcNumChannels: srcNumChannels,
	}
}

// Reset resets the converter
func (c *PCMChannelsConverter) Reset() {
	c.srcSamples = 0
}

// Add adds a new sample to the converter
func (c *PCMChannelsConverter) Add(i int) (err error) {
	// Forward sample
	if c.srcNumChannels == c.dstNumChannels {
		if err = c.fn(i); err != nil {
			err = fmt.Errorf("astikit: handling sample failed: %w", err)
			return
		}
		return
	}

	// Reset
	if c.srcSamples == c.srcNumChannels {
		c.srcSamples = 0
	}

	// Increment src samples
	c.srcSamples++

	// Throw away data
	if c.srcNumChannels > c.dstNumChannels {
		// Throw away sample
		if c.srcSamples > c.dstNumChannels {
			return
		}

		// Custom
		if err = c.fn(i); err != nil {
			err = fmt.Errorf("astikit: handling sample failed: %w", err)
			return
		}
		return
	}

	// Store
	var ss []int
	if c.srcSamples < c.srcNumChannels {
		ss = []int{i}
	} else {
		// Repeat data
		for idx := c.srcNumChannels; idx <= c.dstNumChannels; idx++ {
			ss = append(ss, i)
		}
	}

	// Loop through samples
	for _, s := range ss {
		// Custom
		if err = c.fn(s); err != nil {
			err = fmt.Errorf("astikit: handling sample failed: %w", err)
			return
		}
	}
	return
}

// PCMSilenceDetector represents a PCM silence detector
type PCMSilenceDetector struct {
	analyses              []pcmSilenceDetectorAnalysis
	buf                   []int
	m                     *sync.Mutex // Locks buf
	minAnalysesPerSilence int
	o                     PCMSilenceDetectorOptions
	samplesPerAnalysis    int
}

type pcmSilenceDetectorAnalysis struct {
	level   float64
	samples []int
}

// PCMSilenceDetectorOptions represents a PCM silence detector options
type PCMSilenceDetectorOptions struct {
	MaxSilenceLevel    float64       `toml:"max_silence_level"`
	MinSilenceDuration time.Duration `toml:"min_silence_duration"`
	SampleRate         int           `toml:"sample_rate"`
	StepDuration       time.Duration `toml:"step_duration"`
}

// NewPCMSilenceDetector creates a new silence detector
func NewPCMSilenceDetector(o PCMSilenceDetectorOptions) (d *PCMSilenceDetector) {
	// Create
	d = &PCMSilenceDetector{
		m: &sync.Mutex{},
		o: o,
	}

	// Reset
	d.Reset()

	// Default option values
	if d.o.MinSilenceDuration == 0 {
		d.o.MinSilenceDuration = time.Second
	}
	if d.o.StepDuration == 0 {
		d.o.StepDuration = 30 * time.Millisecond
	}

	// Compute attributes depending on options
	d.samplesPerAnalysis = int(math.Floor(float64(d.o.SampleRate) * d.o.StepDuration.Seconds()))
	d.minAnalysesPerSilence = int(math.Floor(d.o.MinSilenceDuration.Seconds() / d.o.StepDuration.Seconds()))
	return
}

// Reset resets the silence detector
func (d *PCMSilenceDetector) Reset() {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Reset
	d.analyses = []pcmSilenceDetectorAnalysis{}
	d.buf = []int{}
}

// Add adds samples to the buffer and checks whether there are valid samples between silences
func (d *PCMSilenceDetector) Add(samples []int) (validSamples [][]int) {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Append samples to buffer
	d.buf = append(d.buf, samples...)

	// Analyze samples by step
	for len(d.buf) >= d.samplesPerAnalysis {
		// Append analysis
		d.analyses = append(d.analyses, pcmSilenceDetectorAnalysis{
			level:   PCMLevel(d.buf[:d.samplesPerAnalysis]),
			samples: append([]int(nil), d.buf[:d.samplesPerAnalysis]...),
		})

		// Remove samples from buffer
		d.buf = d.buf[d.samplesPerAnalysis:]
	}

	// Loop through analyses
	var leadingSilence, inBetween, trailingSilence int
	for i := 0; i < len(d.analyses); i++ {
		if d.analyses[i].level < d.o.MaxSilenceLevel {
			// This is a silence

			// This is a leading silence
			if inBetween == 0 {
				leadingSilence++

				// The leading silence is valid
				// We can trim its useless part
				if leadingSilence > d.minAnalysesPerSilence {
					d.analyses = d.analyses[leadingSilence-d.minAnalysesPerSilence:]
					i -= leadingSilence - d.minAnalysesPerSilence
					leadingSilence = d.minAnalysesPerSilence
				}
				continue
			}

			// This is a trailing silence
			trailingSilence++

			// Trailing silence is invalid
			if trailingSilence < d.minAnalysesPerSilence {
				continue
			}

			// Trailing silence is valid
			// Loop through analyses
			var ss []int
			for _, a := range d.analyses[:i+1] {
				ss = append(ss, a.samples...)
			}

			// Append valid samples
			validSamples = append(validSamples, ss)

			// Remove leading silence and non silence
			d.analyses = d.analyses[leadingSilence+inBetween:]
			i -= leadingSilence + inBetween

			// Reset counts
			leadingSilence, inBetween, trailingSilence = trailingSilence, 0, 0
		} else {
			// This is not a silence

			// This is a leading non silence
			// We need to remove it
			if i == 0 {
				d.analyses = d.analyses[1:]
				i = -1
				continue
			}

			// This is the first in-between
			if inBetween == 0 {
				// The leading silence is invalid
				// We need to remove it as well as this first non silence
				if leadingSilence < d.minAnalysesPerSilence {
					d.analyses = d.analyses[i+1:]
					i = -1
					continue
				}
			}

			// This non-silence was preceded by a silence not big enough to be a valid trailing silence
			// We incorporate it in the in-between
			if trailingSilence > 0 {
				inBetween += trailingSilence
				trailingSilence = 0
			}

			// This is an in-between
			inBetween++
			continue
		}
	}
	return
}
