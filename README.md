# matchmaker

This project will spin up a game and match making service. The executable expects 1 command line argument being the address and port where the server listens (for example, `localhost:8888` or `0.0.0.0:443`).

## Compile and run

To compile and run the server:

```shell
go build && ./matchmaker localhost:8888
```

## Packages

Current packages are:

1. `server.go`: implements the main server instance that contains the matchmaker and game services.
2. `lobby`: contains the matchmaking logic and HTTP handler for the `/matchmaker` endpoint.
3. `game`: contains the game logic and HTTP handler for the `/game` endpoint.

`lobby` and `game` services communicate with the clients using websockets.

## Static files

Static files are served from the `static` directory. For example, `http://localhost:8888/index.html` will serve `./static/index.html`.

# Sites

## Matchmaker

To enter the matchmaker start the server and connect through the browser.

## Game

To enter a game go to `/game.html?room_id=<ID>`, where `<ID>` can be any string. If a room with `<ID>` does not exist, one will be created. If it exists, then the game will join that room (unless it's already full).
