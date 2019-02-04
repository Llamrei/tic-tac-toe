# A simple implementation of a tic-tac-toe (noughts and crosses) http server.

Can have one server with arbitrary number of games of nxn grids with arbitrary amount of players.

To set-up, on server with latest Go installation:
1. Go into /server/main.go and configure your game how you would like - the default is 3x3 grid with two players: 'X' and 'O' - X starts first. For guidance on setting up more abitrary game look inside /examples/.
2. > go run server/main.go

To play, on client (if using browser, disable caching for best results):
1. Send GET requests to [server_endpoint/?token=[token]&x=[x]&x=[y]]() to make a move.
* Given an invalid move you will be given a hint as to what you did wrong.
2. Send GET requests to [server_endpoint/]() to check state of board, who's turn it is, and whether the game is over

#   Future work:

* All interaction is through GET request as front-end was not part of assignment and manually sending POSTs isn't very user friendly. Ideally would want a nice interactive front end.

* Ability for # of games, nxn grids and # of player games to be manageable through requests to server.

* Implement some tests