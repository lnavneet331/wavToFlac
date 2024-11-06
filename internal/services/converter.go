package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// Converter holds conversion settings for sample rate, channels, etc.
type Converter struct {
	sampleRate    int
	numChannels   int
	bitsPerSample int
}

// NewConverter initializes a new Converter instance with default values.
func NewConverter() *Converter {
	return &Converter{
		sampleRate:    44100,
		numChannels:   2,
		bitsPerSample: 16,
	}
}

// ConvertChunk converts WAV data to FLAC using external ffmpeg tool for encoding.
func (c *Converter) ConvertChunk(wavData []byte) ([]byte, error) {
	// Create a WAV decoder to read the audio data.
	wavReader := bytes.NewReader(wavData)
	decoder := wav.NewDecoder(wavReader)

	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV data")
	}

	// Decode the WAV data into an audio buffer
	buf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: c.numChannels,
			SampleRate:  c.sampleRate,
		},
	}
	if _, err := decoder.PCMBuffer(buf); err != nil {
		return nil, err
	}

	// Write WAV data to a temporary file for ffmpeg to read
	tmpWavFile, err := ioutil.TempFile("", "input*.wav")
	if err != nil {
		return nil, err
	}
	defer tmpWavFile.Close()
	defer os.Remove(tmpWavFile.Name())

	if _, err := tmpWavFile.Write(wavData); err != nil {
		return nil, err
	}

	// Set up the ffmpeg command to read the WAV file and output FLAC data
	cmd := exec.Command("ffmpeg", "-i", tmpWavFile.Name(), "-f", "flac", "pipe:1")
	var flacBuffer bytes.Buffer
	cmd.Stdout = &flacBuffer

	// Run the ffmpeg command to convert the WAV to FLAC
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error converting WAV to FLAC: %v", err)
	}

	return flacBuffer.Bytes(), nil
}
