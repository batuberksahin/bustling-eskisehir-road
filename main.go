package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	gameName = "Bustling Eskisehir Road"
	sWidth   = 720
	sHeight  = 1280
)

type Human struct {
	FirstCoord rl.Vector2
	Coord      rl.Vector2
	Correct    bool
}

type Customer struct {
	Coord  rl.Vector2
	To     float32
	isTake bool
	isGet  bool
}

type Building struct {
	Texture rl.Texture2D
	Coord   rl.Vector2
}

func main() {
	// START WINDOW
	rl.InitWindow(sWidth, sHeight, gameName)
	rl.SetTargetFPS(60)

	// VARIABLES
	var Score int

	var Humans []Human
	var Customers []Customer
	var Buildings []Building

	var buildingTextures []rl.Texture2D

	Passenger := 50
	PassengerText := false

	isStarted := false
	isMoving := false
	isPassedStation := false
	isPoliceComing := false

	doingCarCombo := false
	carComboFlag := true
	comboText := false
	comboCont := false
	comboMultiplier := 1

	gearSituation := false
	driverSituation := 0 // 0: Normal, 1: Speeding Up, 2: Speeding Down
	firstSpeed := float32(8)
	currentSpeed := firstSpeed

	direction := 0 // 0: Direct, -1: Left, 1: Right
	directionFlag := false

	playerPosition := rl.Vector2{440, 750}
	playerAngle := float32(0)
	stationPosition := rl.Vector2{550, 0}
	camera := rl.NewCamera2D(rl.Vector2{0, 0}, playerPosition, 0, 1)
	policePosition := rl.Vector2{0, -200}
	randomPoliceDistance := float32(0)

	// SOUNDS
	rl.InitAudioDevice()

	engineSound := rl.LoadMusicStream("assets/engineloop.mp3")
	rl.PlayMusicStream(engineSound)

	successJump := rl.LoadSound("assets/successjump.mp3")
	failJump := rl.LoadSound("assets/failjump.mp3")
	akbil := rl.LoadSound("assets/akbil.mp3")

	// TEXTURES
	background := rl.LoadTexture("assets/background.png")

	roadTexture := rl.LoadTexture("assets/road.png")
	busTexture := rl.LoadTexture("assets/ego.png")
	stationTexture := rl.LoadTexture("assets/station.png")
	policeTexture := rl.LoadTexture("assets/police.png")
	manCorrectTexture := rl.LoadTexture("assets/man-true.png")
	manInvalidTexture := rl.LoadTexture("assets/man-false.png")

	guvenparkTexture := rl.LoadTexture("assets/guvenpark.png")
	buildingTextures = append(buildingTextures, guvenparkTexture)

	milliKutuphaneTexture := rl.LoadTexture("assets/milliKutuphane.png")
	buildingTextures = append(buildingTextures, milliKutuphaneTexture)

	astiTexture := rl.LoadTexture("assets/asti.png")
	buildingTextures = append(buildingTextures, astiTexture)

	hacettepeTexture := rl.LoadTexture("assets/hacettepe.png")
	buildingTextures = append(buildingTextures, hacettepeTexture)

	havalimaniTexture := rl.LoadTexture("assets/havalimani.png")
	buildingTextures = append(buildingTextures, havalimaniTexture)

	odtuTexture := rl.LoadTexture("assets/odtu.png")
	buildingTextures = append(buildingTextures, odtuTexture)

	firstBuilding := Building{
		Texture: buildingTextures[0],
		Coord: rl.Vector2{
			X: 20,
			Y: 1000,
		},
	}

	Buildings = append(Buildings, firstBuilding)

	// SHAPES
	backgroundRectangle := rl.NewRectangle(0, 0, 720, 1280)

	roadRectangle := rl.NewRectangle(0, 0, 256, 50000)

	sourceBusRectangle := rl.NewRectangle(0, 0, 43, 150)
	busRectangle := rl.NewRectangle(playerPosition.X, playerPosition.Y, 43*1.6, 150*1.6)

	sourcePoliceRectangle := rl.NewRectangle(0, 0, 33, 57)
	policeRectangle := rl.NewRectangle(policePosition.X, policePosition.Y, 33*1.6, 57*1.6)

	// APPLICATION
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		if isStarted {
			rl.ClearBackground(rl.DarkGreen)

			if isMoving {
				rand.Seed(time.Now().UnixNano())

				// ENGINE SOUND
				rl.UpdateMusicStream(engineSound)

				// COLLISON CHECK
				if rl.CheckCollisionRecs(busRectangle, policeRectangle) {
					collisionRectangle := rl.GetCollisionRec(busRectangle, policeRectangle)
					y1 := busRectangle.Y
					y2 := collisionRectangle.Y

					radian := playerAngle * 3.14 / 180

					forbiddenWith := math.Tan(float64(radian)) * math.Abs(float64(y2-y1))

					if collisionRectangle.Width > float32(math.Abs(forbiddenWith)) {
						isMoving = false
					}
				}

				// GEAR CONTROL
				if !gearSituation {
					go func() {
						gearSituation = true

						currentSpeed = firstSpeed

						driverSituation = rand.Intn(3)

						time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)

						gearSituation = false
					}()
				}

				acceleration := rand.Float32()

				if driverSituation == 1 {
					currentSpeed = currentSpeed + acceleration/1000
					playerPosition.Y -= currentSpeed
					camera.Offset.Y += currentSpeed
				} else if driverSituation == 2 {
					currentSpeed = currentSpeed - acceleration/1000
					playerPosition.Y -= currentSpeed
					camera.Offset.Y += currentSpeed
				} else {
					playerPosition.Y -= currentSpeed
					camera.Offset.Y += currentSpeed
				}

				// STATION HANDLER
				if !isPassedStation {
					randomPosition := rand.Intn(300)

					stationPosition.Y = playerPosition.Y - 1250 + float32(randomPosition)

					isPassedStation = true

					// CREATE CUSTOMERS
					customerOnStation := rand.Intn(6)

					for i := 0; i < customerOnStation; i++ {
						customer := Customer{
							Coord: rl.Vector2{
								X: stationPosition.X + 20,
								Y: stationPosition.Y - 60 + float32(rand.Intn(220)),
							},
							To:     playerPosition.X,
							isTake: false,
							isGet:  false,
						}

						Customers = append(Customers, customer)
					}
				}

				if isPassedStation && playerPosition.Y+400 < stationPosition.Y {
					isPassedStation = false

					Customers = []Customer{}
				}

				// POLICE HANDLER
				if !isPoliceComing {
					LoR := rand.Intn(2)

					if LoR == 0 {
						policePosition.X = 360
						policePosition.Y = playerPosition.Y - 1000
					} else {
						policePosition.X = 470
						policePosition.Y = playerPosition.Y - 1000
					}

					rand.Seed(time.Now().UnixNano())
					randomPoliceDistance = float32(rand.Intn(1200))

					isPoliceComing = true
				}

				if isPoliceComing {
					policePosition.Y += 2
				}

				if isPoliceComing && playerPosition.Y+500+randomPoliceDistance < policePosition.Y {
					isPoliceComing = false

					comboCont = false

					carComboFlag = true
					doingCarCombo = false
				}

				// EVENT HANDLER
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) || rl.IsKeyPressed(rl.KeySpace) {
					touchPosition := rl.GetMousePosition()

					// RELEASE PASSENGER
					if touchPosition.Y < sHeight/2 && touchPosition.X > sWidth/2 {
						if Passenger > 0 {
							if playerPosition.Y+10 > stationPosition.Y-60 && playerPosition.Y+10 < stationPosition.Y+220 && playerPosition.X+100 > stationPosition.X {
								human := Human{
									FirstCoord: rl.Vector2{
										X: playerPosition.X + 40,
										Y: playerPosition.Y + 10,
									},
									Coord: rl.Vector2{
										X: playerPosition.X + 40,
										Y: playerPosition.Y + 10,
									},
									Correct: true,
								}

								Score += 2

								Humans = append(Humans, human)

								if rl.IsSoundPlaying(successJump) {
									rl.StopSound(successJump)
									rl.PlaySound(successJump)
								} else {
									rl.PlaySound(successJump)
								}
							} else {
								human := Human{
									FirstCoord: rl.Vector2{
										X: playerPosition.X + 40,
										Y: playerPosition.Y + 10,
									},
									Coord: rl.Vector2{
										X: playerPosition.X + 40,
										Y: playerPosition.Y + 10,
									},
									Correct: false,
								}

								Score--

								Humans = append(Humans, human)

								rl.PlaySound(failJump)
							}

							Passenger--
						} else {
							go func() {
								PassengerText = true

								time.Sleep(time.Duration(3) * time.Second)

								PassengerText = false
							}()
						}
					} else if touchPosition.Y < sHeight/2 && touchPosition.X < sWidth/2 {
						if playerPosition.Y+10 > stationPosition.Y-60 && playerPosition.Y+10 < stationPosition.Y+220 && playerPosition.X+100 > stationPosition.X {
							if len(Customers) > 0 {
								whichCustomer := rand.Intn(len(Customers))

								Customers[whichCustomer].isTake = true

								Passenger++

								if !rl.IsSoundPlaying(akbil) {
									rl.SetSoundVolume(akbil, 0.5)
									rl.PlaySound(akbil)
								}
							}
						}
					}

				}

				if rl.IsMouseButtonDown(rl.MouseLeftButton) || rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyRight) {
					touchPosition := rl.GetMousePosition()

					// TURN LEFT
					if (touchPosition.Y > sHeight/2 && touchPosition.X < sWidth/2) || rl.IsKeyDown(rl.KeyLeft) {
						if playerPosition.X > 350 {
							playerPosition.X -= 5

							if playerAngle >= -15 {
								playerAngle -= 3
							} else {
								playerAngle = -15
							}
							direction = -1
						}
					}

					// TURN RIGHT
					if (touchPosition.Y > sHeight/2 && touchPosition.X > sWidth/2) || rl.IsKeyDown(rl.KeyRight) {
						if playerPosition.X < 482 {
							playerPosition.X += 5

							if playerAngle <= 15 {
								playerAngle += 3
							} else {
								playerAngle = 15
							}
							direction = 1
						}
					}
				}

				if rl.IsMouseButtonUp(rl.MouseLeftButton) {
					directionFlag = true
				}

				// DIRECTION CONTROLLER
				if directionFlag {
					if playerAngle >= 5 || playerAngle <= -5 {
						if direction == -1 {
							playerAngle += 5
						} else if direction == 1 {
							playerAngle -= 5
						} else {
							playerAngle = 0
							directionFlag = false
						}
					} else {
						playerAngle = 0
						directionFlag = false
					}
				}

				// GARBAGE COLLECTOR XD
				if len(Humans) > 1 {
					if playerPosition.Y+400 < Humans[0].Coord.Y {
						Humans = Humans[1:]
					}
				}

				if len(Buildings) > 6 {
					if playerPosition.Y+400 < Buildings[0].Coord.Y {
						Buildings = Buildings[1:]
					}
				}

				// BUILDINGS
				if len(Buildings) < 10 {
					buildingRandom := rand.Intn(len(buildingTextures))

					building := Building{
						Texture: buildingTextures[buildingRandom],
						Coord: rl.Vector2{
							X: Buildings[len(Buildings)-1].Coord.X,
							Y: Buildings[len(Buildings)-1].Coord.Y - 270,
						},
					}

					Buildings = append(Buildings, building)
				}

				// COMBO

				if !doingCarCombo {

					if playerPosition.X-10 < policePosition.X+40+30 && playerPosition.Y-60 > policePosition.Y-10-30 && playerPosition.Y-60 < policePosition.Y-10+150 {
						doingCarCombo = true
						comboCont = true
						// POLIS SAĞDAN GELIYO
					} else if playerPosition.X-10+43 > policePosition.X+40-33-30 && playerPosition.Y-60 > policePosition.Y-10-30 && playerPosition.Y-60 < policePosition.Y-10+150 {
						doingCarCombo = true
						comboCont = true
						// POLIS SOLDAN GELIYO
					}
				}

				if isPoliceComing && playerPosition.Y+43*1.6 < policePosition.Y {
					if !comboCont {
						comboMultiplier = 1
					}
				}

				if carComboFlag {
					if doingCarCombo {
						go func() {
							comboText = true

							time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)

							comboText = false
						}()

						comboMultiplier *= 2

						carComboFlag = false
					}
				}

				// DRAWING THINGS
				rl.BeginMode2D(camera)

				rl.DrawTextureRec(roadTexture, roadRectangle, rl.Vector2{300, -49000}, rl.White)

				for _, building := range Buildings {
					rl.DrawTextureEx(building.Texture, building.Coord, 0, 1, rl.White)
				}

				busRectangle.X = playerPosition.X
				busRectangle.Y = playerPosition.Y
				rl.DrawTexturePro(busTexture, sourceBusRectangle, busRectangle, rl.Vector2{21, 75}, playerAngle, rl.White)

				rl.DrawTextureEx(stationTexture, stationPosition, 0, 2, rl.White)

				policeRectangle.X = policePosition.X
				policeRectangle.Y = policePosition.Y
				rl.DrawTexturePro(policeTexture, sourcePoliceRectangle, policeRectangle, rl.Vector2{16, 28}, 0, rl.White)

				for id, human := range Humans {
					if human.Correct {
						Humans[id].Coord.X += 1
						rl.DrawTextureEx(manCorrectTexture, human.Coord, 90, 1.5, rl.White)
					} else {
						Humans[id].Coord.Y += 3
						rl.DrawTextureEx(manInvalidTexture, human.Coord, 0, 1.5, rl.White)
					}
				}

				for id, customer := range Customers {
					if !customer.isGet {
						if customer.isTake {
							if customer.Coord.X > customer.To+120 {
								Customers[id].Coord.X -= 3
								rl.DrawTextureEx(manCorrectTexture, customer.Coord, 90, 1.5, rl.White)
							} else {
								customer.isGet = true
							}
						} else {
							rl.DrawTextureEx(manCorrectTexture, customer.Coord, 0, 1.5, rl.White)
						}
					}
				}

				// rl.DrawCircleV(rl.Vector2{stationPosition.X, stationPosition.Y - 60}, 3, rl.Red)
				// rl.DrawCircleV(rl.Vector2{stationPosition.X, stationPosition.Y + 220}, 3, rl.Red)
				// rl.DrawCircleV(rl.Vector2{playerPosition.X + 40, playerPosition.Y + 10}, 3, rl.Red)

				// rl.DrawCircleV(rl.Vector2{policePosition.X - 10 + 50 + 30, policePosition.Y - 10}, 3, rl.Red)
				// rl.DrawCircleV(rl.Vector2{playerPosition.X - 10, playerPosition.Y - 60}, 3, rl.Red)
				// rl.DrawCircleV(rl.Vector2{policeOrigin.X, policeOrigin.Y}, 3, rl.Red)

				rl.EndMode2D()

				rl.DrawText("Passenger: "+strconv.Itoa(Passenger), 20, 20, 30, rl.Black)
				rl.DrawText("Score: "+strconv.Itoa(Score), 20, 60, 30, rl.Black)

				if comboText {
					// rl.DrawText("WOW! "+strconv.Itoa(comboMultiplier)+"X COMBO", (sWidth-rl.MeasureText("WOW! "+strconv.Itoa(comboMultiplier)+"X COMBO", 60))/2, 290, 60, rl.Yellow)
					// COMBO KAPALI
				}

				if PassengerText {
					rl.DrawText("TAKE PASSENGERS", (sWidth-rl.MeasureText("TAKE PASSENGERS", 60))/2, 340, 60, rl.Red)
				}

			} else {
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					Score = 0

					Humans = []Human{}
					Buildings = []Building{}
					Customers = []Customer{}

					Passenger = 50

					// doingStationCombo = false
					doingCarCombo = false
					carComboFlag = true
					comboText = false
					comboMultiplier = 1

					isMoving = true
					isPassedStation = false
					isPoliceComing = false

					gearSituation = false
					driverSituation = 0 // 0: Normal, 1: Speeding Up, 2: Speeding Down
					firstSpeed = float32(8)
					currentSpeed = firstSpeed

					direction = 0 // 0: Direct, -1: Left, 1: Right
					directionFlag = false

					playerPosition = rl.Vector2{440, 750}
					playerAngle = float32(0)
					stationPosition = rl.Vector2{550, 0}
					camera = rl.NewCamera2D(rl.Vector2{0, 0}, playerPosition, 0, 1)
					policePosition = rl.Vector2{0, -200}
					randomPoliceDistance = float32(0)

					firstBuilding = Building{
						Texture: buildingTextures[0],
						Coord: rl.Vector2{
							X: 20,
							Y: 1000,
						},
					}

					Buildings = append(Buildings, firstBuilding)

					// BUILDINGS
					if len(Buildings) < 10 {
						buildingRandom := rand.Intn(len(buildingTextures))

						building := Building{
							Texture: buildingTextures[buildingRandom],
							Coord: rl.Vector2{
								X: Buildings[len(Buildings)-1].Coord.X,
								Y: Buildings[len(Buildings)-1].Coord.Y - 270,
							},
						}

						Buildings = append(Buildings, building)
					}
				}

				rl.BeginMode2D(camera)

				rl.DrawTextureRec(roadTexture, roadRectangle, rl.Vector2{300, -49000}, rl.White)

				for _, building := range Buildings {
					rl.DrawTextureEx(building.Texture, building.Coord, 0, 1, rl.White)
				}

				busRectangle.X = playerPosition.X
				busRectangle.Y = playerPosition.Y
				rl.DrawTexturePro(busTexture, sourceBusRectangle, busRectangle, rl.Vector2{21, 75}, playerAngle, rl.White)

				rl.DrawTextureEx(stationTexture, stationPosition, 0, 3, rl.White)

				policeRectangle.X = policePosition.X
				policeRectangle.Y = policePosition.Y
				rl.DrawTexturePro(policeTexture, sourcePoliceRectangle, policeRectangle, rl.Vector2{16, 28}, 0, rl.White)

				for id, human := range Humans {
					if human.Correct {
						Humans[id].Coord.X += 1
						rl.DrawTextureEx(manCorrectTexture, human.Coord, 90, 1.5, rl.White)
					} else {
						Humans[id].Coord.Y += 3
						rl.DrawTextureEx(manInvalidTexture, human.Coord, 0, 1.5, rl.White)
					}
				}

				rl.EndMode2D()

				rl.DrawText("Game Over", (sWidth-rl.MeasureText("Game Over", 60))/2, sHeight/2-40, 60, rl.DarkPurple)
				rl.DrawText("Tap to Play Again", (sWidth-rl.MeasureText("Tap to Play Again", 50))/2, sHeight/2+40, 50, rl.DarkPurple)

			}

		} else {
			// WELCOME MENU
			rl.ClearBackground(rl.LightGray)

			// rl.DrawText("Welcome to Ring Simulator", 30, sHeight/2-100, 50, rl.DarkPurple)

			// rl.DrawText("Tap to Start", 180, sHeight/2, 50, rl.DarkPurple)

			rl.DrawTextureRec(background, backgroundRectangle, rl.Vector2{0, 0}, rl.White)

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				isStarted = true
				isMoving = true
			}
		}

		rl.EndDrawing()
	}
}
