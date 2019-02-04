package tictactoe

import "fmt"

// Maybe a bit much, if zero value of primitives in language spec change
var zeroString string
var zeroRune rune

type game interface {
	// Done so you could implement a board that stores values in db
	// but PoC
	Register(token rune) (game, error)
	StartGame(token rune) (game, error)
	Move(token rune, x int, y int) (game, error)
}

type gameError struct {
	// Done so you could extend errors to be more meaningful
	// but PoC
	msg string
}

func (err *gameError) Error() string {
	return err.msg
}

type gameState int

const (
	registration gameState = iota + 1
	play
	won
	draw
)

type goBoard struct {
	// Struct for storing game state on the heap
	name         string
	tokens       []rune
	whosTurn     int
	board        [][]rune
	state        gameState
	tokensPlaced int
}

func newGoBoard(x, y int, boardName string, admin rune) (*goBoard, error) {
	b := new(goBoard)
	if boardName == zeroString {
		return b, &gameError{"ZeroName"}
	}
	b.name = boardName
	if admin == zeroRune {
		return b, &gameError{"ZeroToken"}
	}
	b.tokens = make([]rune, 1)
	b.tokens[0] = admin
	if x == 0 || y == 0 {
		return b, &gameError{"ZeroBoard"}
	}
	b.board = make([][]rune, y)
	for row := range b.board {
		b.board[row] = make([]rune, x)
	}
	// Enter registration period - doubles as check of successful board
	b.state = registration
	return b, nil
}

func (gb goBoard) Register(newToken rune) (game, error) {
	// Could be made more efficient by handling newToken slice
	// but writing PoC
	for _, takenToken := range gb.tokens {
		if newToken == takenToken {
			return gb, &gameError{"TokenTaken"}
		} else if newToken == zeroRune {
			return gb, &gameError{"ZeroToken"}
		}
	}
	gb.tokens = append(gb.tokens, newToken)
	return gb, nil
}

func (gb goBoard) StartGame(token rune) (game, error) {
	if token == gb.tokens[0] {
		gb.state = play
		return gb, nil
	}
	return gb, &gameError{"NotAdmin"}
}

func (gb goBoard) Move(token rune, x int, y int) (game, error) {
	if gb.state != play {
		return gb, &gameError{fmt.Sprintf("%vState", gb.state)}
	}
	if token == gb.tokens[gb.whosTurn] {
		if !onBoard(gb.board, x, y) {
			return gb, &gameError{"OffBoard"}
		}
		if square := gb.board[y][x]; square == zeroRune {
			gb.board[y][x] = token
			gb.tokensPlaced++
			// Should move into seperate function but PoC
			winningMove := false
			// Search clockwise round placed token, top-left is first visited
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dy == 0 && dx == 0 {
						continue
					}
					// Could be further optimised , but PoC
					if onBoard(gb.board, x+dx, y+dy) {
						// Can we match in this dir
						if token == gb.board[y+dy][x+dx] {
							// Check again in same dir
							if onBoard(gb.board, x+2*dx, y+2*dy) {
								if token == gb.board[y+2*dy][x+2*dx] {
									winningMove = true
								}
							}
							// Check in opposite dir
							if onBoard(gb.board, x-dx, y-dy) {
								if token == gb.board[y-dy][x-dx] {
									winningMove = true
								}
							}
						}
					}
				}
			}
			if winningMove {
				gb.state = won
				return gb, nil
			} else if gb.tokensPlaced >= len(gb.board)*len(gb.board[0]) {
				gb.state = draw
				return gb, nil
			}
			gb.whosTurn++
			if gb.whosTurn >= len(gb.tokens) {
				gb.whosTurn = 0
			}
		} else {
			return gb, &gameError{"SquareTaken"}
		}
		return gb, nil
	}
	// Could be due to player not existing or not his turn
	return gb, &gameError{"InvalidPlayer"}
}

func onBoard(board [][]rune, x, y int) bool {
	dy := len(board)
	dx := len(board[0])
	if x >= dx-1 || y >= dy-1 || x < 0 || y < 0 {
		return false
	}
	return true
}
