package utils

import (
	
	"encoding/binary"
	
	"math"

	"audio-converter/internal/models"
)

// AudioBuffer represents a circular buffer for audio processing
type AudioBuffer struct {
	data       []float64
	size       int
	writeIndex int
	readIndex  int
}

func NewAudioBuffer(size int) *AudioBuffer {
	return &AudioBuffer{
		data:       make([]float64, size),
		size:       size,
		writeIndex: 0,
		readIndex:  0,
	}
}

// ValidateWAVHeader checks if the provided data has a valid WAV header
func ValidateWAVHeader(data []byte) (*models.AudioFormat, error) {
	if len(data) < 44 { // Minimum WAV header size
		return nil, &models.ConversionError{
			Code:    models.ErrInvalidFormat,
			Message: "Invalid WAV header size",
		}
	}

	// Check RIFF header
	if string(data[0:4]) != "RIFF" {
		return nil, &models.ConversionError{
			Code:    models.ErrInvalidFormat,
			Message: "Invalid RIFF header",
		}
	}

	// Parse format
	format := &models.AudioFormat{
		NumChannels:   int(binary.LittleEndian.Uint16(data[22:24])),
		SampleRate:    int(binary.LittleEndian.Uint32(data[24:28])),
		BitsPerSample: int(binary.LittleEndian.Uint16(data[34:36])),
	}

	return format, nil
}

// ConvertSampleRate converts audio samples from one sample rate to another
func ConvertSampleRate(input []float64, inputRate, outputRate int) []float64 {
	ratio := float64(outputRate) / float64(inputRate)
	outputLength := int(float64(len(input)) * ratio)
	output := make([]float64, outputLength)

	for i := range output {
		position := float64(i) / ratio
		index := int(position)
		fraction := position - float64(index)

		if index+1 < len(input) {
			output[i] = input[index]*(1-fraction) + input[index+1]*fraction
		} else {
			output[i] = input[index]
		}
	}

	return output
}

// NormalizeSamples normalizes audio samples to prevent clipping
func NormalizeSamples(samples []float64) []float64 {
	// Find the maximum absolute value
	maxAbs := 0.0
	for _, sample := range samples {
		abs := math.Abs(sample)
		if abs > maxAbs {
			maxAbs = abs
		}
	}

	// Normalize if necessary
	if maxAbs > 1.0 {
		normalized := make([]float64, len(samples))
		for i, sample := range samples {
			normalized[i] = sample / maxAbs
		}
		return normalized
	}

	return samples
}

// CalculateRMS calculates the Root Mean Square of audio samples
func CalculateRMS(samples []float64) float64 {
	var sum float64
	for _, sample := range samples {
		sum += sample * sample
	}
	return math.Sqrt(sum / float64(len(samples)))
}

// ApplyGain applies gain to audio samples
func ApplyGain(samples []float64, gainDB float64) []float64 {
	gain := math.Pow(10, gainDB/20.0)
	output := make([]float64, len(samples))
	for i, sample := range samples {
		output[i] = sample * gain
	}
	return output
}