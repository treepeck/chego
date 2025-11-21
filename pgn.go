/*
pgn.go implements conversions between Portable Game Notation (PGN) strings and
the [Game] structure.  Functions in this file expect the passed PGN strings and
Game variables to be valid, and may panic if they aren't.

Exported PGN strings consists of 8 parts:
 1. Event (Example - ranked blitz game | normal rapid game);
 2. Site (URL to the game. Example - https://justchess.org/nf8KbSog);
 3. Date (the starting date of the game. Example - 2025.10.12);
 4. Round (? if players haven't play a match yet);
 5. White (Id of the white player);
 6. Black (Id of the black player);
 7. Result (the result of the game);
 8. Termination (the reason of the game ending).
*/

package chego

/*
SerializePGN serializes the specified [Game] into a PGN string.

[Event "rated bullet game"]
[Site "https://lichess.org/uzsRzKS7"]
[Date "2025.10.13"]
[White "chess-art-us"]
[Black "SavvaVetokhin2009"]
[Result "1-0"]
[GameId "uzsRzKS7"]
[UTCDate "2025.10.13"]
[UTCTime "08:09:08"]
[WhiteElo "3159"]
[BlackElo "3073"]
[WhiteRatingDiff "+4"]
[BlackRatingDiff "-5"]
[WhiteTitle "GM"]
[BlackTitle "GM"]
[Variant "Standard"]
[TimeControl "60+0"]
[ECO "A07"]
[Opening "King's Indian Attack: Keres Variation"]
[Termination "Normal"]
[Annotator "lichess.org"]

1. Nf3 { [%clk 0:01:00] } 1... d5 { [%clk 0:01:00] } 2. g3 { [%clk 0:00:59] }
2... Bg4 { [%clk 0:01:00] } { A07 King's Indian Attack: Keres Variation } 3.
Bg2 { [%clk 0:00:59] } 3... c6 { [%clk 0:00:59] } 4. h3 { [%clk 0:00:58] } 4...
Bxf3 { [%clk 0:00:58] } 5. Bxf3 { [%clk 0:00:58] } 5... e6 { [%clk 0:00:57] }
6. Bxd5 { [%clk 0:00:58] } 6... Nf6 { [%clk 0:00:57] } 7. Bf3 { [%clk 0:00:57]
} 7... Bd6 { [%clk 0:00:56] } 8. d3 { [%clk 0:00:56] } 8... h5 { [%clk 0:00:55]
} 9. e4 { [%clk 0:00:55] } 9... h4 { [%clk 0:00:54] } 10. g4 { [%clk 0:00:54] }
10... c5 { [%clk 0:00:53] } 11. Bg2 { [%clk 0:00:53] } 11... Nc6 { [%clk 0:00:52
} 12. f4 { [%clk 0:00:53] } 12... Nd7 { [%clk 0:00:52] } 13. e5 { [%clk 0:00:52]
} 13... Ndxe5 { [%clk 0:00:50] } 14. fxe5 { [%clk 0:00:51] } 14... Bxe5 { [%clk
0:00:50] } 15. O-O { [%clk 0:00:51] } 15... Bc7 { [%clk 0:00:47] } 16. Nc3 {
[%clk 0:00:50] } 16... Qd4+ { [%clk 0:00:47] } 17. Kh1 { [%clk 0:00:49] }
17... g5 { [%clk 0:00:45] } 18. Nb5 { [%clk 0:00:47] } 18... Qe5 { [%clk
0:00:43] } 19. Nxc7+ { [%clk 0:00:47] } 19... Qxc7 { [%clk 0:00:43] } 20. Bxg5
{ [%clk 0:00:46] } 20... Nd4 { [%clk 0:00:42] } 21. c3 { [%clk 0:00:45] } 21...
Nf3 { [%clk 0:00:41] } { Black resigns. } 1-0
*/
func SerializePGN(g Game) string {

	return ""
}
