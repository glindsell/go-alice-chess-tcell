package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/notnil/chess"
)

var logger *log.Logger

type Piece int
type Column int
type Row int

const (
	King Piece = iota
	Queen
	Rook
	Bishop
	Knight
	Pawn
	Empty
)

const (
	CA Column = iota
	CB
	CC
	CD
	CE
	CF
	CG
	CH
)

const (
	R8 Row = iota
	R7
	R6
	R5
	R4
	R3
	R2
	R1
)

func (p Piece) WhiteRune() rune {
	return [...]rune{'\u2654', '\u2655', '\u2656', '\u2657', '\u2658', '\u2659', ' '}[p]
}
func (p Piece) BlackRune() rune {
	return [...]rune{'\u265a', '\u265b', '\u265c', '\u265d', '\u265e', '\u265f', ' '}[p]
}
func (c Column) UpperRune() rune {
	return [...]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'}[c]
}
func (c Column) LowerRune() rune {
	return [...]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}[c]
}
func (r Row) Rune() rune {
	return [...]rune{'8', '7', '6', '5', '4', '3', '2', '1'}[r]
}
func (c Column) UpperString() string {
	return [...]string{"A", "B", "C", "D", "E", "F", "G", "H"}[c]
}
func (c Column) LowerString() string {
	return [...]string{"a", "b", "c", "d", "e", "f", "g", "h"}[c]
}
func (r Row) String() string {
	return [...]string{"8", "7", "6", "5", "4", "3", "2", "1"}[r]
}

var StartingRow = []Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

type Board struct {
	XPos         int
	YPos         int
	XScale       int
	YScale       int
	Screen       tcell.Screen
	ColourScheme *ColourScheme
	StyleScheme  *StyleScheme
	Squares      map[string]Position
}

type Position struct {
	XPos int
	YPos int
}

type Session struct {
	Board1       *Board
	Board2       *Board
	Game1        *chess.Game
	Game2        *chess.Game
	Screen       tcell.Screen
	ColourScheme *ColourScheme
	StyleScheme  *StyleScheme
}

type ColourScheme struct {
	ScreenText            tcell.Color
	ScreenBackground      tcell.Color
	LightSquare           tcell.Color
	DarkSquare            tcell.Color
	WhitePiece            tcell.Color
	BlackPiece            tcell.Color
	MenuText              tcell.Color
	MenuBackground        tcell.Color
	InstructionText       tcell.Color
	InstructionBackground tcell.Color
}

type StyleScheme struct {
	Default            tcell.Style
	WhiteOnLight       tcell.Style
	WhiteOnDark        tcell.Style
	BlackOnLight       tcell.Style
	BlackOnDark        tcell.Style
	WritingMenu        tcell.Style
	WritingInstruction tcell.Style
}

func NewStyleScheme(cs *ColourScheme) *StyleScheme {
	return &StyleScheme{
		Default:            tcell.StyleDefault.Foreground(cs.ScreenText).Background(cs.ScreenBackground),
		WhiteOnLight:       tcell.StyleDefault.Foreground(cs.WhitePiece).Background(cs.LightSquare),
		WhiteOnDark:        tcell.StyleDefault.Foreground(cs.WhitePiece).Background(cs.DarkSquare),
		BlackOnLight:       tcell.StyleDefault.Foreground(cs.BlackPiece).Background(cs.LightSquare),
		BlackOnDark:        tcell.StyleDefault.Foreground(cs.BlackPiece).Background(cs.DarkSquare),
		WritingMenu:        tcell.StyleDefault.Foreground(cs.MenuText).Background(cs.MenuBackground).Bold(true),
		WritingInstruction: tcell.StyleDefault.Foreground(cs.InstructionText).Background(cs.InstructionBackground),
	}
}

func NewSession(b1, b2 *Board, g1, g2 *chess.Game, s tcell.Screen, cScheme *ColourScheme, stScheme *StyleScheme) *Session {
	return &Session{
		Board1:       b1,
		Board2:       b2,
		Game1:        g1,
		Game2:        g2,
		Screen:       s,
		ColourScheme: cScheme,
		StyleScheme:  stScheme,
	}
}

func NewBoard(x, y, xScale, yScale int, s tcell.Screen, cScheme *ColourScheme, stScheme *StyleScheme) *Board {
	squares := make(map[string]Position)
	for row := 0; row <= 7; row++ {
		for col := 0; col <= 7; col++ {
			name := Column(col).LowerString() + Row(row).String()
			squares[name] = Position{
				XPos: x + ((col + 1) * xScale),
				YPos: y + ((row + 1) * yScale),
			}
		}
	}
	logger.Println(squares)
	return &Board{
		XPos:         x,
		YPos:         y,
		XScale:       xScale,
		YScale:       yScale,
		Screen:       s,
		ColourScheme: cScheme,
		StyleScheme:  stScheme,
		Squares:      squares,
	}
}

func main() {
	f, err := os.Create("text.log")
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(f, "go-alice-chess: ", log.LstdFlags)
	defer f.Close()
	logger.Println("start log")
	s, err := tcell.NewScreen()
	if err != nil {
		logger.Fatal(err)
	}
	err = s.Init()
	if err != nil {
		logger.Fatal(err)
	}
	colorScheme := &ColourScheme{
		ScreenText:            tcell.ColorWhite,
		ScreenBackground:      tcell.ColorDarkSlateGrey,
		LightSquare:           tcell.ColorSandyBrown,
		DarkSquare:            tcell.ColorSaddleBrown,
		WhitePiece:            tcell.ColorWhite,
		BlackPiece:            tcell.ColorBlack,
		MenuText:              tcell.ColorSteelBlue,
		MenuBackground:        tcell.ColorSilver,
		InstructionText:       tcell.ColorSeaGreen,
		InstructionBackground: tcell.ColorSilver,
	}
	styleScheme := NewStyleScheme(colorScheme)
	s.SetStyle(styleScheme.Default)
	s.EnableMouse()
	s.Clear()
	b1 := NewBoard(2, 1, 2, 1, s, colorScheme, styleScheme)
	b2 := NewBoard(b1.XPos+11*b1.XScale, b1.YPos, b1.XScale, b1.YScale, s, colorScheme, styleScheme)
	b1.Clear()
	b2.Clear()
	EmitStr(s, b2.XPos+11*b2.XScale, b2.YPos+b2.YScale, styleScheme.WritingMenu, "Press Esc to exit.")
	EmitStr(s, b2.XPos+11*b2.XScale, b2.YPos+2*b2.YScale, styleScheme.WritingMenu, "Press S to start new game.")
	s.Show()
	g1 := chess.NewGame()
	fenStr := "8/8/8/8/8/8/8/8 w - - 0 1"
	fen, err := chess.FEN(fenStr)
	if err != nil {
		log.Fatal(err)
	}
	g2 := chess.NewGame(fen)
	session := NewSession(b1, b2, g1, g2, s, colorScheme, styleScheme)
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			} else if ev.Rune() == 'S' || ev.Rune() == 's' {
				session.Start()
				s.Sync()
				session.PlayRandomTurns(2)
				s.Show()
			}
		case *tcell.EventMouse:
			mx, my := ev.Position()
			session.PrintMouse(mx, my)
			s.Show()
		}
	}
}

func (g *Session) PrintMouse(mx, my int) {
	posfmt := "Mouse: x: %d, y: %d"
	EmitStr(g.Screen, g.Board2.XPos+g.Board2.XScale, g.Board2.YPos+11*g.Board2.YScale, g.StyleScheme.Default, fmt.Sprintf(posfmt, mx, my))
}

func (g *Session) Start() {
	g.Board1.Clear()
	g.Board2.Clear()
	g.Board1.Fill()
	EmitStr(g.Screen, g.Board1.XPos+g.Board1.XScale, g.Board1.YPos+11*g.Board1.YScale, g.StyleScheme.WritingInstruction, "White to move.")
}

func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (b *Board) Clear() {
	for col := 1; col <= 8*b.XScale; col++ {
		for row := 1; row <= 8*b.YScale; row++ {
			if col%b.XScale == 0 && row%b.YScale == 0 {
				if b.LightSquare(col, row) {
					//Light square
					for scCol := 0; scCol < b.XScale; scCol++ {
						for scRow := 0; scRow < b.YScale; scRow++ {
							b.Screen.SetContent(b.XPos+col+scCol, b.YPos+row+scRow, Empty.WhiteRune(), nil, b.StyleScheme.WhiteOnLight)
						}

					}
				} else {
					//Dark square
					for scCol := 0; scCol < b.XScale; scCol++ {
						for scRow := 0; scRow < b.YScale; scRow++ {
							b.Screen.SetContent(b.XPos+col+scCol, b.YPos+row+scRow, Empty.BlackRune(), nil, b.StyleScheme.WhiteOnDark)
						}
					}
				}
			}
		}
	}
	for col := 1; col <= 8*b.XScale; col++ {
		if col%b.XScale == 0 {
			b.Screen.SetContent(b.XPos+col, b.YPos+9*b.YScale, Column((col-1)/b.XScale).UpperRune(), nil, b.StyleScheme.Default)
		}
	}
	for row := 1; row <= 8*b.YScale; row++ {
		if row%b.YScale == 0 {
			b.Screen.SetContent(b.XPos, b.YPos+row, Row((row-1)/b.YScale).Rune(), nil, b.StyleScheme.Default)
		}
	}
}

func (b *Board) Fill() {
	for col := 1; col <= 8*b.XScale; col++ {
		for row := 1; row <= 8*b.YScale; row++ {
			if col%b.XScale == 0 && row%b.YScale == 0 {
				switch row {
				case 1:
					if b.LightSquare(col, row) {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, StartingRow[(col-1)/b.XScale].BlackRune(), nil, b.StyleScheme.BlackOnLight)
					} else {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, StartingRow[(col-1)/b.XScale].BlackRune(), nil, b.StyleScheme.BlackOnDark)
					}
				case 2 * b.YScale:
					if b.LightSquare(col, row) {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, Pawn.BlackRune(), nil, b.StyleScheme.BlackOnLight)
					} else {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, Pawn.BlackRune(), nil, b.StyleScheme.BlackOnDark)
					}
				case 7 * b.YScale:
					if b.LightSquare(col, row) {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, Pawn.WhiteRune(), nil, b.StyleScheme.WhiteOnLight)
					} else {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, Pawn.WhiteRune(), nil, b.StyleScheme.WhiteOnDark)
					}
				case 8 * b.YScale:
					if b.LightSquare(col, row) {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, StartingRow[(col-1)/b.XScale].WhiteRune(), nil, b.StyleScheme.WhiteOnLight)
					} else {
						b.Screen.SetContent(b.XPos+col, b.YPos+row, StartingRow[(col-1)/b.XScale].WhiteRune(), nil, b.StyleScheme.WhiteOnDark)
					}
				}
			}
		}
	}
}

func (b *Board) LightSquare(col, row int) bool {
	if Abs(col/b.XScale-row/b.YScale)%2 == 0 {
		return true
	}
	return false
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *Session) PlayRandomTurns(turns int) {
	for turn := 1; turn <= turns; turn++ {
		time.Sleep(1 * time.Second)
		// select a random move
		moves := s.Game1.ValidMoves()
		move := moves[rand.Intn(len(moves))]
		logger.Println(move.String())
		err := s.Game1.Move(move)
		if err != nil {
			logger.Fatal(err)
		}
		/*err = s.Game2.Move(move)
		if err != nil {
			s.Screen.Fini()
			logger.Fatal(err)
		}*/
		s.Move1to2(move, turn)
		s.Screen.Show()
	}
	return
}

func (g *Session) Move1to2(move *chess.Move, turn int) {
	colFrom := move.S1().File()
	rowFrom := move.S1().Rank()
	colTo := move.S2().File()
	rowTo := move.S2().Rank()
	pos1 := g.Board1.Squares[colFrom.String()+rowFrom.String()]
	pos2 := g.Board2.Squares[colTo.String()+rowTo.String()]
	piece, _, st1, _ := g.Screen.GetContent(pos1.XPos, pos1.YPos)
	g.Board1.Screen.SetContent(pos1.XPos, pos1.YPos, Empty.WhiteRune(), nil, st1)
	_, _, st2, _ := g.Screen.GetContent(pos2.XPos, pos2.YPos)
	if turn%2 == 0 {
		st2 = st2.Foreground(g.ColourScheme.BlackPiece)
	} else {
		st2 = st2.Foreground(g.ColourScheme.WhitePiece)
	}
	g.Board2.Screen.SetContent(pos2.XPos, pos2.YPos, piece, nil, st2)
}

func (b *Board) Add(move *chess.Move, piece rune) {
	col := move.S1().File()
	row := move.S1().Rank()
	pos := b.Squares[col.String()+row.String()]
	_, _, st, _ := b.Screen.GetContent(pos.XPos, pos.YPos)
	b.Screen.SetContent(pos.XPos, pos.YPos, piece, nil, st)
}
