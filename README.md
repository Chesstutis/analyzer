![Go Version](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)
[![Tests](https://github.com/chesstutis/analyzer/actions/workflows/test.yml/badge.svg)](
https://github.com/chesstutis/analyzer/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/chesstutis/analyzer/graph/badge.svg)](
https://codecov.io/github/chesstutis/analyzer)
# Chesstutis Game Analyzer 
- given a chess "Game" from github.com/corentings/chess/v2 output a game analysis and puzzles based on the mistakes each player made
- for now analysis will only output puzzles


# TODO
- Move Classification: 
  - Brilliant, Best, Good, Inaccuracy, Mistake, Blunder, Forced
- Filter Per Player:
  - add a config value for chess.Color so game will only generate puzzles for one color (will save lots of time)
- Eval Per Move:
  - return cpl, classification for each move to create game analysis like chess.com
- Puzzle Lines:
  - add complete lines to puzzles rather than just being one move
- Skip Opening:
  - integrate corentings/opening package so the analysis will skip the "Book Moves"
- GitHub Workflows:
  - setup some github workflows for releases, test coverage...
### test
