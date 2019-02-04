package engine

import (
	"fmt"
	"strings"
)

// Maybe a bit much, enables control of 'zero' values and in case language spec changes
var zeroString string
var zeroRune rune = ' '

type Game interface {
	// Done so you could implement a board that stores values in db
	// but PoC
	Register(token rune) (Game, error)
	StartGame(token rune) (Game, error)
	Move(token rune, x int, y int) (Game, error)
	GetState() GameState
	String() string
}

type gameError struct {
	// Done so you could extend errors to be more meaningful
	// but PoC
	msg string
}

func (err *gameError) Error() string {
	return err.msg
}

type GameState int

const (
	Registration GameState = iota + 1
	Play
	Won
	Draw
)

type goBoard struct {
	// Struct for storing Game state on the heap
	name         string
	tokens       []rune
	whosTurn     int
	board        [][]rune
	state        GameState
	tokensPlaced int
}

// NewGoBoard Creates a new tictactoe board with in-memory (not database) storage.
// Admin token is needed to start the Game, will be first in turn order.
func NewGoBoard(x, y int, boardName string, admin rune) (*goBoard, error) {
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
		for column := range b.board[row] {
			b.board[row][column] = zeroRune
		}
	}
	// Enter Registration period - doubles as check of successful board
	b.state = Registration
	return b, nil
}

// Register Adds a token to the list of valid tokens - can only happen during registration
func (gb *goBoard) Register(newToken rune) (Game, error) {
	if gb.state == Registration {
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
	return gb, &gameError{"NotRegistering"}
}

// StartGame Given the Game's admin token, can finish registration period
func (gb *goBoard) StartGame(token rune) (Game, error) {
	if token == gb.tokens[0] {
		gb.state = Play
		return gb, nil
	}
	return gb, &gameError{"NotAdmin"}
}

// Move Try to place a token on x,y - can only happen during play
func (gb *goBoard) Move(token rune, x int, y int) (Game, error) {
	if gb.state != Play {
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
				gb.state = Won
				return gb, nil
			} else if gb.tokensPlaced >= len(gb.board)*len(gb.board[0]) {
				gb.state = Draw
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
	// Could be due to Player not existing or not his turn
	return gb, &gameError{"InvalidPlayer"}
}

func (gb *goBoard) GetState() GameState {
	return gb.state
}

func (gb *goBoard) String() string {
	var sb strings.Builder
	for row := range gb.board {
		sb.WriteRune('|')
		for column := range gb.board[row] {
			sb.WriteRune(gb.board[row][column])
			sb.WriteRune('|')
		}
		sb.WriteString("\n")
	}
	if gb.state == Won {
		sb.WriteString(fmt.Sprintf("%v won!\n", string(gb.tokens[gb.whosTurn])))
	} else if gb.state == Draw {
		sb.WriteString("Draw!\n")
	} else {
		sb.WriteString(fmt.Sprintf("%v to move\n", string(gb.tokens[gb.whosTurn])))
	}
	return sb.String()
}

func onBoard(board [][]rune, x, y int) bool {
	dy := len(board)
	dx := len(board[0])
	if x > dx-1 || y > dy-1 || x < 0 || y < 0 {
		return false
	}
	return true
}
