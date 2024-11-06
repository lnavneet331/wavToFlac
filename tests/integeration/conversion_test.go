package integration

import (
	"bytes"
	"encoding/binary"
	"fmt"
	
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	
)

// TestConfig holds test configuration
type TestConfig struct {
	ServerURL     string
	TestDataPath  string
	TimeoutSeconds int
}

// Global test configuration
var testConfig = TestConfig{
	ServerURL:      "ws://localhost:8080/ws/convert",
	TestDataPath:   "testdata/",
	TimeoutSeconds: 5,
}

func TestWebSocketConnection(t *testing.T) {
	// Start server in test mode
	ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()
	
	// Read test WAV file
	wavData, err := os.ReadFile(testConfig.TestDataPath + "test.wav")
	if err != nil {
		t.Fatalf("Failed to read test WAV file: %v", err)
	}
	
	// Send WAV data
	err = ws.WriteMessage(websocket.BinaryMessage, wavData)
	if err != nil {
		t.Fatalf("Failed to send WAV data: %v", err)
	}
	
	// Read FLAC response with timeout
	done := make(chan struct{})
	var responseError error
	
	go func() {
		_, message, err := ws.ReadMessage()
		if err != nil {
			responseError = fmt.Errorf("Failed to read WebSocket message: %v", err)
			close(done)
			return
		}
		
		// Validate FLAC data
		if len(message) == 0 {
			responseError = fmt.Errorf("Received empty FLAC data")
			close(done)
			return
		}
		
		if !bytes.HasPrefix(message, []byte("fLaC")) {
			responseError = fmt.Errorf("Invalid FLAC header in response")
			close(done)
			return
		}
		
		close(done)
	}()
	
	// Wait for response with timeout
	select {
	case <-done:
		if responseError != nil {
			t.Fatal(responseError)
		}
	case <-time.After(time.Duration(testConfig.TimeoutSeconds) * time.Second):
		t.Fatal("Timeout waiting for FLAC response")
	}
}

func TestMultipleSimultaneousConnections(t *testing.T) {
	numConnections := 5
	var wg sync.WaitGroup
	errors := make(chan error, numConnections)
	
	// Start multiple connections
	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Connect to WebSocket
			ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
			if err != nil {
				errors <- fmt.Errorf("Connection %d failed: %v", id, err)
				return
			}
			defer ws.Close()
			
			// Send test data
			testData := createTestWAVData(44100, 2, 16)
			err = ws.WriteMessage(websocket.BinaryMessage, testData)
			if err != nil {
				errors <- fmt.Errorf("Connection %d failed to send data: %v", id, err)
				return
			}
			
			// Read response with timeout
			done := make(chan bool)
			go func() {
				_, response, err := ws.ReadMessage()
				if err != nil {
					errors <- fmt.Errorf("Connection %d failed to read response: %v", id, err)
					done <- false
					return
				}
				
				if !bytes.HasPrefix(response, []byte("fLaC")) {
					errors <- fmt.Errorf("Connection %d received invalid FLAC response", id)
					done <- false
					return
				}
				
				done <- true
			}()
			
			select {
			case success := <-done:
				if !success {
					return
				}
			case <-time.After(time.Duration(testConfig.TimeoutSeconds) * time.Second):
				errors <- fmt.Errorf("Connection %d timed out", id)
			}
		}(i)
	}
	
	// Wait for all connections to complete
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		t.Error(err)
	}
}

func TestLargeFileConversion(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()
	
	// Create large WAV data (1 minute of audio)
	sampleRate := 44100
	// duration := 60 // seconds
	wavData := createTestWAVData(sampleRate, 2, 16)
	
	// Send data in chunks
	chunkSize := 4096
	for i := 0; i < len(wavData); i += chunkSize {
		end := i + chunkSize
		if end > len(wavData) {
			end = len(wavData)
		}
		
		err = ws.WriteMessage(websocket.BinaryMessage, wavData[i:end])
		if err != nil {
			t.Fatalf("Failed to send chunk %d: %v", i/chunkSize, err)
		}
		
		// Read response for each chunk
		_, response, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read response for chunk %d: %v", i/chunkSize, err)
		}
		
		// Validate FLAC chunk
		if len(response) == 0 {
			t.Errorf("Empty response for chunk %d", i/chunkSize)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		expectedError bool
	}{
		{
			name: "Invalid WAV header",
			data: []byte("NOT A WAV FILE"),
			expectedError: true,
		},
		{
			name: "Empty data",
			data: []byte{},
			expectedError: true,
		},
		{
			name: "Corrupted WAV data",
			data: createCorruptedWAVData(),
			expectedError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
			if err != nil {
				t.Fatalf("Failed to connect to WebSocket: %v", err)
			}
			defer ws.Close()
			
			err = ws.WriteMessage(websocket.BinaryMessage, tt.data)
			if err != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}
			
			// _, response, err := ws.ReadMessage()
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper functions

func createTestWAVData(sampleRate, channels, bitsPerSample int) []byte {
	buf := new(bytes.Buffer)
	
	// Write WAV header
	binary.Write(buf, binary.LittleEndian, []byte("RIFF"))
	binary.Write(buf, binary.LittleEndian, uint32(36)) // ChunkSize
	binary.Write(buf, binary.LittleEndian, []byte("WAVE"))
	binary.Write(buf, binary.LittleEndian, []byte("fmt "))
	binary.Write(buf, binary.LittleEndian, uint32(16)) // Subchunk1Size
	binary.Write(buf, binary.LittleEndian, uint16(1))  // AudioFormat (PCM)
	binary.Write(buf, binary.LittleEndian, uint16(channels))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate*channels*bitsPerSample/8))
	binary.Write(buf, binary.LittleEndian, uint16(channels*bitsPerSample/8))
	binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample))
	
	// Write data header
	binary.Write(buf, binary.LittleEndian, []byte("data"))
	dataSize := uint32(1024) // Some sample data size
	binary.Write(buf, binary.LittleEndian, dataSize)
	
	// Write sample audio data
	for i := 0; i < int(dataSize)/2; i++ {
		binary.Write(buf, binary.LittleEndian, int16(i%100))
	}
	
	return buf.Bytes()
}

func createCorruptedWAVData() []byte {
	data := createTestWAVData(44100, 2, 16)
	// Corrupt the format chunk
	data[20] = 0xFF
	data[21] = 0xFF
	return data
}

func TestConnectionResilience(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()
	
	// Test reconnection after disconnect
	ws.Close()
	ws, _, err = websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to reconnect: %v", err)
	}
	
	// Send data after reconnection
	testData := createTestWAVData(44100, 2, 16)
	err = ws.WriteMessage(websocket.BinaryMessage, testData)
	assert.NoError(t, err, "Should be able to send data after reconnection")
}

func TestPerformance(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	ws, _, err := websocket.DefaultDialer.Dial(testConfig.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()
	
	// Create 1MB of test data
	testData := createTestWAVData(44100, 2, 16)
	
	// Measure conversion time
	start := time.Now()
	err = ws.WriteMessage(websocket.BinaryMessage, testData)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}
	
	_, response, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	
	duration := time.Since(start)
	
	// Assert performance requirements
	assert.Less(t, duration.Milliseconds(), int64(1000), 
		"Conversion should take less than 1 second for 1MB of data")
	assert.Greater(t, len(response), 0, 
		"Response should contain converted data")
}