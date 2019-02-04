package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/llamrei/tictactoe/engine"
)

func main() {
	fmt.Println("Starting tictactoe server!")
	// Initially just a 3x3
	classicGame, e := engine.NewGoBoard(3, 3, "classic", 'X')
	if e != nil {
		fmt.Println(e)
	}
	_, e = classicGame.Register('O')
	if e != nil {
		fmt.Println(e)
	}
	_, e = classicGame.StartGame('X')
	if e != nil {
		fmt.Println(e)
	}
	scanner := bufio.NewScanner(os.Stdin)
	var str string
	for {
		fmt.Println("Enter Token:")
		scanner.Scan()
		token := rune(scanner.Bytes()[0])

		fmt.Println("Pick X:")
		scanner.Scan()
		str = scanner.Text()
		x, _ := strconv.Atoi(str)

		fmt.Println("Pick Y:")
		scanner.Scan()
		str = scanner.Text()
		y, _ := strconv.Atoi(str)
		_, err := classicGame.Move(token, x, y)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(classicGame)
		if classicGame.GetState() == engine.Won {
			fmt.Println("You win!")
			break
		} else if classicGame.GetState() == engine.Draw {
			fmt.Println("You drew!")
			break
		}
	}
}
