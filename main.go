package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"

	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/gorilla/websocket"
)

const (
	screenWidth     = 1024
	screenHeight    = 600
	framesPerSecond = 30
)

type SimConnectData struct {
	// Add the required data fields from the raw simconnect messages
}

func main() {
	raylib.InitWindow(screenWidth, screenHeight, "Aviation Six-Pack")
	raylib.SetTargetFPS(framesPerSecond)

	// Connect to the WebSocket server
	// dataChannel := connectoToWebSocket()

	// Main loop for drawing the six-pack
	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		// select {
		// case data := <-dataChannel:
		// 	// Update the six-pack display using the received data
		// 	updateSixPackDisplay(data)
		// default:
		// 	// If no new data, continue drawing the last frame
		// }

		drawSixPack()

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}

func connectoToWebSocket() chan SimConnectData {
	u := url.URL{Scheme: "ws", Host: "<your-websocket-server-addr>", Path: "/simconnect"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()
	// Set up a channel to receive data from the WebSocket
	dataChannel := make(chan SimConnectData)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading WebSocket message:", err)
				continue
			}

			var data SimConnectData
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println("Error unmarshaling WebSocket message:", err)
				continue
			}

			dataChannel <- data
		}
	}()

	return dataChannel
}

func updateSixPackDisplay(data SimConnectData) {
	// Update the internal state of the six-pack display using the received data
}

func drawSixPack() {
	// Define variables for each of the instruments' values
	// For demo purposes, we use example values
	airSpeed := 47.0
	altitude := 10000.0
	verticalSpeed := 500.0
	heading := 45.0
	pitch := 2.0
	roll := 10.0

	// Airspeed Indicator
	drawAirspeedIndicator(airSpeed, 50, 50, 200)

	// Artificial Horizon
	drawArtificialHorizon(pitch, roll, screenWidth/2-100, 50, 200)

	// Altimeter
	drawAltimeter(altitude, screenWidth-250, 50, 200)

	// Vertical Speed Indicator
	drawVerticalSpeedIndicator(verticalSpeed, 50, screenHeight/2+50, 200)

	// Heading Indicator
	drawHeadingIndicator(heading, screenWidth/2-100, screenHeight/2+50, 200)

	// Turn Coordinator (not part of the six-pack, but often found in Cessna 152)
	drawTurnCoordinator(roll, screenWidth-250, screenHeight/2+50, 200)
}

func drawAirspeedIndicator(airSpeed float64, posX, posY, size int) {
	radius := size / 2
	centerX := posX + radius
	centerY := posY + radius

	// Draw the outer circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius), raylib.DarkGray)

	// Draw the inner circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius-10), raylib.RayWhite)

	// Draw the airspeed marks and labels
	for i := 0; i <= 16; i++ {
		angle := float64(i)*22.5 - 135.0
		startX := centerX + int(float64(radius-15)*math.Cos(angle*math.Pi/180))
		startY := centerY - int(float64(radius-15)*math.Sin(angle*math.Pi/180))
		endX := centerX + int(float64(radius-25)*math.Cos(angle*math.Pi/180))
		endY := centerY - int(float64(radius-25)*math.Sin(angle*math.Pi/180))

		raylib.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), raylib.Black)

		// Draw the speed labels
		if i%2 == 0 {
			labelX := centerX + int(float64(radius-45)*math.Cos(angle*math.Pi/180))
			labelY := centerY - int(float64(radius-45)*math.Sin(angle*math.Pi/180))
			raylib.DrawText(fmt.Sprintf("%d", i*10), int32(labelX-5), int32(labelY-5), 10, raylib.Black)
		}
	}

	// Draw the airspeed needle
	needleAngle := airSpeed*2.25 - 135
	needleX := centerX + int(float64(radius-40)*math.Cos(needleAngle*math.Pi/180))
	needleY := centerY - int(float64(radius-40)*math.Sin(needleAngle*math.Pi/180))
	raylib.DrawLine(int32(centerX), int32(centerY), int32(needleX), int32(needleY), raylib.Red)

	// Draw the center circle
	raylib.DrawCircle(int32(centerX), int32(centerY), 10, raylib.RayWhite)
	raylib.DrawCircleLines(int32(centerX), int32(centerY), 10, raylib.Black)
}

func drawArtificialHorizon(pitch, roll float64, posX, posY, size int) {
	radius := size / 2
	centerX := posX + radius
	centerY := posY + radius

	// Rotate the pitch lines according to the roll angle
	rotationCenter := raylib.Vector2{float32(centerX), float32(centerY)}
	raylib.BeginMode2D(raylib.NewCamera2D(rotationCenter, rotationCenter, float32(roll), 1))

	// Draw the sky and ground
	raylib.DrawRectangle(int32(posX), int32(posY)-int32(pitch*10), int32(size), int32(size)/2, raylib.SkyBlue)
	raylib.DrawRectangle(int32(posX), int32(posY+size/2)-int32(pitch*10), int32(size), int32(size)/2, raylib.DarkGreen)

	// Draw pitch lines
	pitchLineLength := 80
	for i := -9; i <= 9; i++ {
		if i == 0 {
			continue
		}
		pitchOffset := int(pitch*10) + i*10
		startX := centerX - pitchLineLength/2
		startY := centerY - pitchOffset
		endX := centerX + pitchLineLength/2
		endY := centerY - pitchOffset

		raylib.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), raylib.Black)
	}

	raylib.EndMode2D()

	// Draw the outer circle
	//	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius), raylib.DarkGray)

	// Draw the roll indicator
	raylib.DrawRectangle(int32(centerX-2), int32(centerY)-int32(radius), 4, 20, raylib.Black)
	raylib.DrawRectangle(int32(centerX-2), int32(centerY)+int32(radius)-20, 4, 20, raylib.Black)

	// Draw a triangle at the top
	raylib.DrawTriangle(raylib.Vector2{float32(centerX), float32(centerY) - float32(radius) + 20},
		raylib.Vector2{float32(centerX) - 5, float32(centerY) - float32(radius) + 35},
		raylib.Vector2{float32(centerX) + 5, float32(centerY) - float32(radius) + 35},
		raylib.Red)

	// Draw the center circle
	raylib.DrawCircle(int32(centerX), int32(centerY), 10, raylib.RayWhite)
	raylib.DrawCircleLines(int32(centerX), int32(centerY), 10, raylib.Black)
}

func drawAltimeter(altitude float64, posX, posY, size int) {
	radius := size / 2
	centerX := posX + radius
	centerY := posY + radius

	// Draw the outer circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius), raylib.Gray)

	// Draw the inner circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius-10), raylib.RayWhite)

	// Draw altitude marks and labels
	for i := 0; i < 36; i++ {
		angle := float64(i) * 10.0
		startX := centerX + int(float64(radius-15)*math.Cos(angle*math.Pi/180))
		startY := centerY - int(float64(radius-15)*math.Sin(angle*math.Pi/180))
		endX := centerX + int(float64(radius-25)*math.Cos(angle*math.Pi/180))
		endY := centerY - int(float64(radius-25)*math.Sin(angle*math.Pi/180))

		raylib.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), raylib.Black)

		// Draw the altitude labels
		if i%6 == 0 {
			labelX := centerX + int(float64(radius-40)*math.Cos(angle*math.Pi/180))
			labelY := centerY - int(float64(radius-40)*math.Sin(angle*math.Pi/180))
			raylib.DrawText(fmt.Sprintf("%d0", i), int32(labelX-5), int32(labelY-5), 10, raylib.Black)
		}
	}

	// Draw the altitude needle
	needleAngle := (altitude / 1000.0) * 360.0
	needleX := centerX + int(float64(radius-40)*math.Cos(needleAngle*math.Pi/180))
	needleY := centerY - int(float64(radius-40)*math.Sin(needleAngle*math.Pi/180))
	raylib.DrawLine(int32(centerX), int32(centerY), int32(needleX), int32(needleY), raylib.Red)

	// Draw the center circle
	raylib.DrawCircle(int32(centerX), int32(centerY), 10, raylib.RayWhite)
}

func drawVerticalSpeedIndicator(verticalSpeed float64, posX, posY, size int) {
	// Draw the vertical speed indicator rectangle
	raylib.DrawRectangle(int32(posX), int32(posY), int32(size), int32(size), raylib.Gray)
	raylib.DrawRectangle(int32(posX+5), int32(posY+5), int32(size-10), int32(size-10), raylib.RayWhite)

	// Draw the vertical speed scale
	scaleHeight := size - 40
	scaleStep := scaleHeight / 12
	for i := 0; i <= 12; i++ {
		y := posY + 20 + i*scaleStep
		raylib.DrawLine(int32(posX+20), int32(y), int32(posX+size-20), int32(y), raylib.Black)
	}

	// Draw the vertical speed labels
	labels := []string{"6", "4", "2", "0", "2", "4", "6"}
	for i, label := range labels {
		y := posY + 20 + i*scaleStep
		raylib.DrawText(label, int32(posX+size/2-10), int32(y-5), 10, raylib.Black)
	}

	// Draw the vertical speed needle
	needleY := posY + 20 + int(scaleHeight/2) - int(verticalSpeed/2000.0*float64(scaleHeight/2))
	raylib.DrawLine(int32(posX+10), int32(needleY), int32(posX+size-10), int32(needleY), raylib.Red)
}

func drawHeadingIndicator(heading float64, posX, posY, size int) {
	radius := size / 2
	centerX := posX + radius
	centerY := posY + radius

	// Draw the outer circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius), raylib.Gray)

	// Draw the inner circle
	raylib.DrawCircle(int32(centerX), int32(centerY), float32(radius-10), raylib.RayWhite)

	// Draw heading marks and labels
	for i := 0; i < 36; i++ {
		angle := float64(i) * 10.0
		startX := centerX + int(float64(radius-15)*math.Cos(angle*math.Pi/180))
		startY := centerY - int(float64(radius-15)*math.Sin(angle*math.Pi/180))
		endX := centerX + int(float64(radius-25)*math.Cos(angle*math.Pi/180))
		endY := centerY - int(float64(radius-25)*math.Sin(angle*math.Pi/180))

		raylib.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), raylib.Black)

		// Draw the heading labels
		if i%9 == 0 {
			labelX := centerX + int(float64(radius-40)*math.Cos(angle*math.Pi/180))
			labelY := centerY - int(float64(radius-40)*math.Sin(angle*math.Pi/180))
			raylib.DrawText(fmt.Sprintf("%d", i*10), int32(labelX-5), int32(labelY-5), 10, raylib.Black)
		}
	}

	// Draw the heading needle
	needleX := centerX + int(float64(radius-20)*math.Sin(heading*math.Pi/180))
	needleY := centerY - int(float64(radius-20)*math.Cos(heading*math.Pi/180))
	raylib.DrawLine(int32(centerX), int32(centerY), int32(needleX), int32(needleY), raylib.Red)

	// Draw the center circle
	raylib.DrawCircle(int32(centerX), int32(centerY), 10, raylib.RayWhite)
}

func drawTurnCoordinator(roll float64, posX, posY, size int) {
	// Draw the turn coordinator rectangle
	raylib.DrawRectangle(int32(posX), int32(posY), int32(size), int32(size/2), raylib.Gray)
	raylib.DrawRectangle(int32(posX+5), int32(posY+5), int32(size-10), int32(size/2-10), raylib.RayWhite)

	// Draw the roll marks and labels
	marks := []float64{-30, -20, -10, 0, 10, 20, 30}
	for _, mark := range marks {
		x := posX + size/2 + int(float64(size/2-20)*math.Sin(mark*math.Pi/180))
		y := posY + size/4 - int(float64(size/2-20)*math.Cos(mark*math.Pi/180))
		raylib.DrawCircle(int32(x), int32(y), 3, raylib.Black)
	}

	// Draw the roll needle
	needleX := posX + size/2 + int(float64(size/2-20)*math.Sin(roll*math.Pi/180))
	needleY := posY + size/4 - int(float64(size/2-20)*math.Cos(roll*math.Pi/180))
	raylib.DrawLine(int32(posX+size/2), int32(posY+size/4), int32(needleX), int32(needleY), raylib.Red)
}
