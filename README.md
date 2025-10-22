# ASCII Grid Runner

A little terminal game in Go where you navigate a random 2D grid of walls (`#`) and open spaces (`.`) and try to escape by moving off the map. Great practice for slices, structs, enums, and simple game loops.

## Features

- Prompt for grid size (n × n) and wall density (0.0–1.0)
- Randomly placed walls (`#`) and one player (`^`, `>`, `v`, or `<`)
- WASD controls to turn and move
- Step counter and “you win” when you step off the grid
- Clears and re-renders the screen each turn

## Requirements

- Go 1.20 or newer

## Installation

1. Clone this repo  
2. `cd` into the project directory  
3. Run:
