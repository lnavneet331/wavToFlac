package unit

import (
	"bytes"
	"encoding/binary"
	"testing"

	"audio-converter/internal/models"
	"audio-converter/internal/services"
	"audio-converter/pkg/utils"
)

func TestConverter_ConvertChunk(t *testing.T) {
	// Create a test WAV chunk
	wavChunk := createTestWAVChunk(t)
	
	// Initialize converter
	converter := services.NewConverter()
	
	// Test conversion
	flacData, err := converter.ConvertChunk(wavChunk)
	if err != nil {
		t.Fatalf("Failed to convert chunk: %v", err)
	}
	
	// Validate FLAC output
	if len(flacData) == 0 {
		t.Error("FLAC data is empty")
	}
	
	// Validate FLAC header
	if !bytes.HasPrefix(flacData, []byte("fLaC")) {
		t.Error("Invalid FLAC header")
	}
}

func TestAudioFormat_Validation(t *testing.T) {
	tests := []struct {
		name        string
		sampleRate  int
		channels    int
		bitDepth    int
		shouldError bool
	}{
		{"valid format", 44100, 2, 16, false},
		{"invalid sample rate", 0, 2, 16, true},
		{"invalid channels", 44100, 0, 16, true},
		{"invalid bit depth", 44100, 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := &models.AudioFormat{
				SampleRate:    tt.sampleRate,
				NumChannels:   tt.channels,
				BitsPerSample: tt.bitDepth,
			}
			
			_, err := utils.ValidateWAVHeader(createWAVHeader(format))
			if (err != nil) != tt.shouldError {
				t.Errorf("ValidateWAVHeader() error = %v, shouldError = %v", err, tt.shouldError)
			}
		})
	}
}

func createTestWAVChunk(t *testing.T) []byte {
	// Create a buffer for WAV data
	buf := new(bytes.Buffer)
	
	// Write WAV header
	header := []byte{
		'R', 'I', 'F', 'F', // ChunkID
		0x24, 0x00, 0x00, 0x00, // ChunkSize
		'W', 'A', 'V', 'E', // Format
		'f', 'm', 't', ' ', // Subchunk1ID
		0x10, 0x00, 0x00, 0x00, // Subchunk1Size
		0x01, 0x00, // AudioFormat (PCM)
		0x02, 0x00, // NumChannels
		0x44, 0xAC, 0x00, 0x00, // SampleRate (44100)
		0x10, 0xB1, 0x02, 0x00, // ByteRate
		0x04, 0x00, // BlockAlign
		0x10, 0x00, // BitsPerSample
	}
	
	buf.Write(header)
	
	// Write some sample audio data
	samples := []int16{0, 100, -100, 200, -200, 300, -300, 400}
	for _, sample := range samples {
		binary.Write(buf, binary.LittleEndian, sample)
	}
	
	return buf.Bytes()
}

func createWAVHeader(format *models.AudioFormat) []byte {
	buf := new(bytes.Buffer)
	// Write minimal WAV header for testing
	binary.Write(buf, binary.LittleEndian, []byte("RIFF"))
	binary.Write(buf, binary.LittleEndian, uint32(36)) // ChunkSize
	binary.Write(buf, binary.LittleEndian, []byte("WAVE"))
	binary.Write(buf, binary.LittleEndian, []byte("fmt "))
	binary.Write(buf, binary.LittleEndian, uint32(16)) // Subchunk1Size
	binary.Write(buf, binary.LittleEndian, uint16(1))  // AudioFormat (PCM)
	binary.Write(buf, binary.LittleEndian, uint16(format.NumChannels))
	binary.Write(buf, binary.LittleEndian, uint32(format.SampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(format.SampleRate*format.NumChannels*format.BitsPerSample/8))
	binary.Write(buf, binary.LittleEndian, uint16(format.NumChannels*format.BitsPerSample/8))
	binary.Write(buf, binary.LittleEndian, uint16(format.BitsPerSample))
	return buf.Bytes()
}