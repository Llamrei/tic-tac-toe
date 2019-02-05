package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/llamrei/tictactoe/engine"
)

func main() {

	fmt.Println("Creating tictactoe engine!")
	// Initially just a 3x3
	classicGame, e := engine.NewGoBoard(3, 3, "classic", 'X')
	if e != nil {
		log.Panic(e)
	}
	_, e = classicGame.Register('O')
	if e != nil {
		log.Panic(e)
	}
	_, e = classicGame.StartGame('X')
	if e != nil {
		log.Panic(e)
	}

	fmt.Println("Starting server")
	http.HandleFunc("/", gameHandler(classicGame))
	log.Fatal(http.ListenAndServe(":80", nil))
}

func gameHandler(g engine.Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			values := r.URL.Query()
			if len(values) == 3 {
				tokenStrs, ok := values["token"]
				if !ok {
					fmt.Fprintf(w, "Please send queries of form in README")
				}
				xStrs, ok := values["x"]
				if !ok {
					fmt.Fprintf(w, "Please send queries of form in README")
				}
				yStrs, ok := values["y"]
				if !ok {
					fmt.Fprintf(w, "Please send queries of form in README")
				}
				if len(tokenStrs) != 1 || len(xStrs) != 1 || len(yStrs) != 1 {
					fmt.Fprintf(w, "Please send queries of form in README")
				} else {
					tokenSlice := []rune(tokenStrs[0])
					x, err := strconv.Atoi(xStrs[0])
					if err != nil {
						fmt.Fprintf(w, "Please send queries of form in README")
					}
					y, err := strconv.Atoi(yStrs[0])
					if err != nil {
						fmt.Fprintf(w, "Please send queries of form in README")
					}
					_, e := g.Move(tokenSlice[0], x, y)
					if e != nil {
						fmt.Fprintf(w, "%s", e)
					} else {
						http.Redirect(w, r, "/", 301)
					}
				}
			} else if len(values) > 0 {
				fmt.Fprintf(w, "Please send queries of form in README")
			} else {
				fmt.Fprintf(w, g.String())
			}
		} else {
			fmt.Fprintf(w, "Please only send GET")
		}
	}
}
