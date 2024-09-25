# pong-multiplayer-go

The classic pong game written in Go with multiplayer support.

<img src="assets/doc/pongo.gif" width="320" height="240" alt="PONGO" />

## Available game modes

- Single player
- Two players
- Multiplayer

### Single player

- Use `Up` and `Down` to move the left paddle up and down.

### Two players

- Player 1: Use `Q` and `A` to move the left paddle up and down.
- Player 2: Use `Up` and `Down` to move the right paddle up and down.

### Multiplayer

To play in multiplayer mode, you need to run the [server](https://github.com/reneepc/pongo-server/) and the game.

#### Server

```bash
make run-server
```

- Player 1: Use `Up` and `Down` to move the left paddle up and down.
- Player 2: Use `Up` and `Down` to move the right paddle up and down.

## How to run the game

```bash
make run
```

Made with :heart: by Gandarez Labs.
