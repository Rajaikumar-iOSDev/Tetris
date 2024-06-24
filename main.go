package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Some Defines
const (
    SquareSize            = 20
    GridHorizontalSize    = 12
    GridVerticalSize      = 20
    LateralSpeed          = 10
    TurningSpeed          = 12
    FastFallAwaitCounter  = 30
    FadingTime            = 33
    ScreenWidth           = 800
    ScreenHeight          = 450
    GravitySpeedInitial   = 30
)

// GridSquare represents the state of a square in the grid
type GridSquare int

// Enumeration for GridSquare
const (
    Empty GridSquare = iota
    Moving
    Full
    Block
    Fading
)

// Global Variables
var (
    gameOver                 bool
    pause                    bool
    grid                     [GridHorizontalSize][GridVerticalSize]GridSquare
    piece                    [4][4]GridSquare
    incomingPiece            [4][4]GridSquare
    piecePositionX           int
    piecePositionY           int
    fadingColor              rl.Color
    beginPlay                = true
    pieceActive              bool
    detection                bool
    lineToDelete             bool
    lines                    = 0
    gravityMovementCounter   int
    lateralMovementCounter   int
    turnMovementCounter      int
    fastFallMovementCounter  int
    fadeLineCounter          int
    gravitySpeed             = GravitySpeedInitial
)

//------------------------------------------------------------------------------------
// Program main entry point
//------------------------------------------------------------------------------------

func main() {
    rl.InitWindow(ScreenWidth, ScreenHeight, "Tetris in Go")
  

    InitGame()
	rl.SetTargetFPS(60);

    for !rl.WindowShouldClose() {
        UpdateDrawFrame()
    }
	rl.CloseWindow();
}


// InitGame initializes the game
func InitGame() {
    // Initialize game statistics
    lines = 0

    fadingColor = rl.Gray

    piecePositionX = 0
    piecePositionY = 0

    pause = false

    beginPlay = true
    pieceActive = false
    detection = false
    lineToDelete = false

    // Counters
    gravityMovementCounter = 0
    lateralMovementCounter = 0
    turnMovementCounter = 0
    fastFallMovementCounter = 0

    fadeLineCounter = 0
    gravitySpeed = GravitySpeedInitial

    // Initialize grid matrices
    for i := 0; i < GridHorizontalSize; i++ {
        for j := 0; j < GridVerticalSize; j++ {
            if j == GridVerticalSize-1 || i == 0 || i == GridHorizontalSize-1 {
                grid[i][j] = Block
            } else {
                grid[i][j] = Empty
            }
        }
    }

    // Initialize incoming piece matrices
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            incomingPiece[i][j] = Empty
        }
    }
}

// UpdateGame updates the game logic for one frame
func UpdateGame() {
    if !gameOver {
        if rl.IsKeyPressed(rl.KeyP) {
            pause = !pause
        }

        if !pause {
            if !lineToDelete {
                if !pieceActive {
                    // Get another piece
                    pieceActive = CreatePiece()

                    // We leave a little time before starting the fast falling down
                    fastFallMovementCounter = 0
                } else { // Piece falling
                    // Counters update
                    fastFallMovementCounter++
                    gravityMovementCounter++
                    lateralMovementCounter++
                    turnMovementCounter++

                    // We make sure to move if we've pressed the key this frame
                    if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
                        lateralMovementCounter = LateralSpeed
                    }
                    if rl.IsKeyPressed(rl.KeyUp) {
                        turnMovementCounter = TurningSpeed
                    }

                    // Fall down
                    if rl.IsKeyDown(rl.KeyDown) && (fastFallMovementCounter >= FastFallAwaitCounter) {
                        // We make sure the piece is going to fall this frame
                        gravityMovementCounter += gravitySpeed
                    }

                    if gravityMovementCounter >= gravitySpeed {
                        // Basic falling movement
                        CheckDetection(&detection)

                        // Check if the piece has collided with another piece or with the boundings
                        ResolveFallingMovement(&detection, &pieceActive)

                        // Check if we fulfilled a line and if so, erase the line and pull down the lines above
                        CheckCompletion(&lineToDelete)

                        gravityMovementCounter = 0
                    }

                    // Move laterally at player's will
                    if lateralMovementCounter >= LateralSpeed {
                        // Update the lateral movement and if success, reset the lateral counter
                        if ResolveLateralMovement() {
                            lateralMovementCounter = 0
                        }
                    }

                    // Turn the piece at player's will
                    if turnMovementCounter >= TurningSpeed {
                        // Update the turning movement and reset the turning counter
                        if ResolveTurnMovement() {
                            turnMovementCounter = 0
                        }
                    }
                }

                // Game over logic
                for j := 0; j < 2; j++ {
                    for i := 1; i < GridHorizontalSize-1; i++ {
                        if grid[i][j] == Full {
                            gameOver = true
                        }
                    }
                }
            } else {
                // Animation when deleting lines
                fadeLineCounter++

                if fadeLineCounter%8 < 4 {
                    fadingColor = rl.Maroon
                } else {
                    fadingColor = rl.Gray
                }

                if fadeLineCounter >= FadingTime {
                    deletedLines := DeleteCompleteLines()
                    fadeLineCounter = 0
                    lineToDelete = false

                    lines += deletedLines
                }
            }
        }
    } else {
        if rl.IsKeyPressed(rl.KeyEnter) {
            InitGame()
            gameOver = false
        }
    }
}

// DrawGame draws the game for one frame
func DrawGame() {
    rl.BeginDrawing()

    rl.ClearBackground(rl.RayWhite)

    if !gameOver {
        // Draw gameplay area
        offset := rl.Vector2{
            X: float32(ScreenWidth)/2 - (GridHorizontalSize*SquareSize/2) - 50,
            Y: float32(ScreenHeight)/2 - ((GridVerticalSize-1)*SquareSize/2) + SquareSize*2,
        }

        offset.Y -= 50 // NOTE: Hardcoded position!

        controller := offset.X

        for j := 0; j < GridVerticalSize; j++ {
            for i := 0; i < GridHorizontalSize; i++ {
                // Draw each square of the grid
                switch grid[i][j] {
                case Empty:
                    rl.DrawLine(int32(offset.X), int32(offset.Y), int32(offset.X+SquareSize), int32(offset.Y), rl.LightGray)
                    rl.DrawLine(int32(offset.X), int32(offset.Y), int32(offset.X), int32(offset.Y+SquareSize), rl.LightGray)
                    rl.DrawLine(int32(offset.X+SquareSize), int32(offset.Y), int32(offset.X+SquareSize), int32(offset.Y+SquareSize), rl.LightGray)
                    rl.DrawLine(int32(offset.X), int32(offset.Y+SquareSize), int32(offset.X+SquareSize), int32(offset.Y+SquareSize), rl.LightGray)
                case Full:
                    rl.DrawRectangle(int32(offset.X), int32(offset.Y), SquareSize, SquareSize, rl.Gray)
                case Moving:
                    rl.DrawRectangle(int32(offset.X), int32(offset.Y), SquareSize, SquareSize, rl.DarkGray)
                case Block:
                    rl.DrawRectangle(int32(offset.X), int32(offset.Y), SquareSize, SquareSize, rl.LightGray)
                case Fading:
                    rl.DrawRectangle(int32(offset.X), int32(offset.Y), SquareSize, SquareSize, fadingColor)
                }

                offset.X += SquareSize
            }

            offset.X = controller
            offset.Y += SquareSize
        }

        // Draw incoming piece (hardcoded)
        offset.X = 500
        offset.Y = 45

        controler := offset.X

        for j := 0; j < 4; j++ {
            for i := 0; i < 4; i++ {
                if incomingPiece[i][j] == Empty {
                    rl.DrawLine(int32(offset.X), int32(offset.Y), int32(offset.X+SquareSize), int32(offset.Y), rl.LightGray)
                    rl.DrawLine(int32(offset.X), int32(offset.Y), int32(offset.X), int32(offset.Y+SquareSize), rl.LightGray)
                    rl.DrawLine(int32(offset.X+SquareSize), int32(offset.Y), int32(offset.X+SquareSize), int32(offset.Y+SquareSize), rl.LightGray)
                    rl.DrawLine(int32(offset.X), int32(offset.Y+SquareSize), int32(offset.X+SquareSize), int32(offset.Y+SquareSize), rl.LightGray)
                } else if incomingPiece[i][j] == Moving {
                    rl.DrawRectangle(int32(offset.X), int32(offset.Y), SquareSize, SquareSize, rl.Gray)
                }

                offset.X += SquareSize
            }

            offset.X = controler
            offset.Y += SquareSize
        }

        rl.DrawText("INCOMING:", int32(offset.X), int32(offset.Y-100), 10, rl.Gray)
        rl.DrawText(fmt.Sprintf("LINES:      %04d", lines) , int32(offset.X), int32(offset.Y+20), 10, rl.Gray)

        if pause {
            rl.DrawText("GAME PAUSED", int32(ScreenWidth)/2-rl.MeasureText("GAME PAUSED", 40)/2, int32(ScreenWidth)/2-40, 40, rl.Gray)
        }
    } else {
        rl.DrawText("PRESS [ENTER] TO PLAY AGAIN", int32(rl.GetScreenWidth())/2-rl.MeasureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, int32(rl.GetScreenHeight())/2-50, 20, rl.Gray)
    }

    rl.EndDrawing()
}


// UpdateDrawFrame updates the game state and draws one frame
func UpdateDrawFrame() {
    UpdateGame()
    DrawGame()
}

// CreatePiece initializes a new piece and places it at the top of the grid
func CreatePiece() bool {
    piecePositionX = (GridHorizontalSize - 4) / 2
    piecePositionY = 0

    // If the game is starting and you are going to create the first piece, we create an extra one
    if beginPlay {
        GetRandomPiece()
        beginPlay = false
    }

    // We assign the incoming piece to the actual piece
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            piece[i][j] = incomingPiece[i][j]
        }
    }

    // We assign a random piece to the incoming one
    GetRandomPiece()

    // Assign the piece to the grid
    for i := piecePositionX; i < piecePositionX+4; i++ {
        for j := 0; j < 4; j++ {
            if piece[i-piecePositionX][j] == Moving {
                grid[i][j] = Moving
            }
        }
    }

    return true
}


// GetRandomPiece generates a random piece and assigns it to the incomingPiece variable
func GetRandomPiece() {
    random := rl.GetRandomValue(0, 6) // Assuming rl is a package that provides GetRandomValue

    // Reset incomingPiece to Empty
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            incomingPiece[i][j] = Empty
        }
    }

    // Assign a new shape to incomingPiece based on the random value
    switch random {
    case 0:
        // Cube
        incomingPiece[1][1] = Moving
        incomingPiece[2][1] = Moving
        incomingPiece[1][2] = Moving
        incomingPiece[2][2] = Moving
    case 1:
        // L
        incomingPiece[1][0] = Moving
        incomingPiece[1][1] = Moving
        incomingPiece[1][2] = Moving
        incomingPiece[2][2] = Moving
    case 2:
        // L inversa
        incomingPiece[1][2] = Moving
        incomingPiece[2][0] = Moving
        incomingPiece[2][1] = Moving
        incomingPiece[2][2] = Moving
    case 3:
        // Recta
        incomingPiece[0][1] = Moving
        incomingPiece[1][1] = Moving
        incomingPiece[2][1] = Moving
        incomingPiece[3][1] = Moving
    case 4:
        // Creu tallada
        incomingPiece[1][0] = Moving
        incomingPiece[1][1] = Moving
        incomingPiece[1][2] = Moving
        incomingPiece[2][1] = Moving
    case 5:
        // S
        incomingPiece[1][1] = Moving
        incomingPiece[2][1] = Moving
        incomingPiece[2][2] = Moving
        incomingPiece[3][2] = Moving
    case 6:
        // S inversa
        incomingPiece[1][2] = Moving
        incomingPiece[2][2] = Moving
        incomingPiece[2][1] = Moving
        incomingPiece[3][1] = Moving
    }
}

// ResolveFallingMovement checks if the current piece should stop Moving (if it has landed) or continue falling.
func ResolveFallingMovement(detection *bool, pieceActive *bool) {
    if *detection {
        // If we finished Moving this piece, we stop it
        for j := GridVerticalSize - 2; j >= 0; j-- {
            for i := 1; i < GridHorizontalSize-1; i++ {
                if grid[i][j] == Moving {
                    grid[i][j] = Full
                    *detection = false
                    *pieceActive = false
                }
            }
        }
    } else {
        // We move down the piece
        for j := GridVerticalSize - 2; j > 0; j-- { // Adjusted loop to prevent index out of range
            for i := 1; i < GridHorizontalSize-1; i++ {
                if grid[i][j] == Moving {
                    grid[i][j+1] = Moving
                    grid[i][j] = Empty
                }
            }
        }

        piecePositionY++
    }
}

// ResolveLateralMovement checks and performs lateral movement of the current piece, returning true if a collision occurs.
func ResolveLateralMovement() bool {
    collision := false

    // Piece movement
    if rl.IsKeyDown(rl.KeyLeft) { // Move left
        // Check if it is possible to move to the left
        for j := GridVerticalSize - 2; j >= 0; j-- {
            for i := 1; i < GridHorizontalSize-1; i++ {
                if grid[i][j] == Moving {
                    // Check if we are touching the left wall or we have a full square at the left
                    if i-1 == 0 || grid[i-1][j] == Full {
                        collision = true
                    }
                }
            }
        }

        // If able, move left
        if !collision {
            for j := GridVerticalSize - 2; j >= 0; j-- {
                for i := 1; i < GridHorizontalSize-1; i++ {
                    if grid[i][j] == Moving {
                        grid[i-1][j] = Moving
                        grid[i][j] = Empty
                    }
                }
            }
            piecePositionX--
        }
    } else if rl.IsKeyDown(rl.KeyRight) { // Move right
        // Check if it is possible to move to the right
        for j := GridVerticalSize - 2; j >= 0; j-- {
            for i := 1; i < GridHorizontalSize-1; i++ {
                if grid[i][j] == Moving {
                    // Check if we are touching the right wall or we have a full square at the right
                    if i+1 == GridHorizontalSize-1 || grid[i+1][j] == Full {
                        collision = true
                    }
                }
            }
        }

        // If able, move right
        if !collision {
            for j := GridVerticalSize - 2; j >= 0; j-- {
                for i := GridHorizontalSize - 1; i >= 1; i-- {
                    if grid[i][j] == Moving {
                        grid[i+1][j] = Moving
                        grid[i][j] = Empty
                    }
                }
            }
            piecePositionX++
        }
    }

    return collision
}

// ResolveTurnMovement checks if the UP key is pressed and rotates the piece if possible.
func ResolveTurnMovement() bool {
    // Input for turning the piece
    if rl.IsKeyDown(rl.KeyUp) {
        var aux GridSquare
        checker := false

        // Check all turning possibilities
        conditions := []bool{
            grid[piecePositionX+3][piecePositionY] == Moving && grid[piecePositionX][piecePositionY] != Empty && grid[piecePositionX][piecePositionY] != Moving,
            grid[piecePositionX+3][piecePositionY+3] == Moving && grid[piecePositionX+3][piecePositionY] != Empty && grid[piecePositionX+3][piecePositionY] != Moving,
            grid[piecePositionX][piecePositionY+3] == Moving && grid[piecePositionX+3][piecePositionY+3] != Empty && grid[piecePositionX+3][piecePositionY+3] != Moving,
            // Add other conditions following the same pattern
        }

        for _, condition := range conditions {
            if condition {
                checker = true
                break
            }
        }

        if !checker {
            // Rotate the piece
            aux = piece[0][0]
            piece[0][0], piece[3][0], piece[3][3], piece[0][3] = piece[3][0], piece[3][3], piece[0][3], aux
            aux = piece[1][0]
            piece[1][0], piece[3][1], piece[2][3], piece[0][2] = piece[3][1], piece[2][3], piece[0][2], aux
            aux = piece[2][0]
            piece[2][0], piece[3][2], piece[1][3], piece[0][1] = piece[3][2], piece[1][3], piece[0][1], aux
            aux = piece[1][1]
            piece[1][1], piece[2][1], piece[2][2], piece[1][2] = piece[2][1], piece[2][2], piece[1][2], aux
        }

        // Clear the Moving piece from the grid
        for j := GridVerticalSize - 2; j >= 0; j-- {
            for i := 1; i < GridHorizontalSize - 1; i++ {
                if grid[i][j] == Moving {
                    grid[i][j] = Empty
                }
            }
        }

        // Place the piece in the new position
        for i := piecePositionX; i < piecePositionX+4; i++ {
            for j := piecePositionY; j < piecePositionY+4; j++ {
                if piece[i-piecePositionX][j-piecePositionY] == Moving {
                    grid[i][j] = Moving
                }
            }
        }

        return true
    }

    return false
}


// CheckDetection checks for detection in a grid.
func CheckDetection(detection *bool) {
    for j := GridVerticalSize - 2; j >= 0; j-- {
        for i := 1; i < GridHorizontalSize-1; i++ {
            if (grid[i][j] == Moving) && (grid[i][j+1] == Full || grid[i][j+1] == Block) {
                *detection = true
            }
        }
    }
}

// CheckCompletion checks each line of the grid to see if it's completely filled.
func CheckCompletion(lineToDelete *bool) {
    for j := GridVerticalSize - 2; j >= 0; j-- {
        calculator := 0
        for i := 1; i < GridHorizontalSize-1; i++ {
            if grid[i][j] == Full {
                calculator++
            }

            if calculator == GridHorizontalSize-2 {
                *lineToDelete = true
                // Reset calculator for the next line
                calculator = 0

                // Mark the completed line for deletion
                for z := 1; z < GridHorizontalSize-1; z++ {
                    grid[z][j] = Fading
                }
            }
        }
    }
}

// DeleteCompleteLines goes through the grid and deletes any lines marked as complete.
func DeleteCompleteLines() int {
    deletedLines := 0

    for j := GridVerticalSize - 2; j >= 0; j-- {
        for grid[1][j] == Fading {
            // Clear the line
            for i := 1; i < GridHorizontalSize-1; i++ {
                grid[i][j] = Empty
            }

            // Move all lines above down
            for j2 := j - 1; j2 >= 0; j2-- {
                for i2 := 1; i2 < GridHorizontalSize-1; i2++ {
                    if grid[i2][j2] == Full || grid[i2][j2] == Fading {
                        grid[i2][j2+1] = grid[i2][j2]
                        grid[i2][j2] = Empty
                    }
                }
            }

            deletedLines++
        }
    }

    return deletedLines
}

