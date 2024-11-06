package models

import "fmt"

type AudioFormat struct {
	SampleRate    int
	NumChannels   int
	BitsPerSample int
}

type ConversionJob struct {
	ID           string
	InputFormat  AudioFormat
	OutputFormat AudioFormat
	Status       string
	ErrorMessage string
	CreatedAt    int64
	CompletedAt  int64
}

type ConversionStats struct {
	TotalBytesProcessed int64
	ConversionTime      int64
	InputFormat         AudioFormat
	OutputFormat        AudioFormat
}

type ConversionError struct {
	Code    string
	Message string
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error constants
const (
	ErrInvalidFormat    = "INVALID_FORMAT"
	ErrConversionFailed = "CONVERSION_FAILED"
	ErrInvalidChunkSize = "INVALID_CHUNK_SIZE"
	ErrStreamCorrupted  = "STREAM_CORRUPTED"
)

// Status constants
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)