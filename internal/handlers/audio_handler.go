package handlers

import (
	"bytes"
	"io"
	"log"

	"audio-converter/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}

func HandleAudioConversion(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	// Create buffers for audio processing
	wavBuffer := bytes.NewBuffer(nil)
	converter := services.NewConverter()

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
			}
			break
		}

		// Handle different message types
		switch mt {
		case websocket.BinaryMessage:
			// Append WAV data to buffer
			wavBuffer.Write(msg)

			// Process complete WAV chunks
			if wavBuffer.Len() >= 4096 { // Process in 4KB chunks
				// Convert WAV chunk to FLAC
				flacData, err := converter.ConvertChunk(wavBuffer.Next(4096))
				if err != nil {
					log.Printf("conversion error: %v", err)
					continue
				}

				// Send converted FLAC data back to client
				if err := c.WriteMessage(websocket.BinaryMessage, flacData); err != nil {
					log.Printf("write error: %v", err)
					break
				}
			}

		case websocket.CloseMessage:
			return
		}
	}
}