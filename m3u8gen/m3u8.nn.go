package m3u8gen

import "os"
import "strconv"
import "time"
import (
	"bufio"
	"io"
	"strings"
)

type frame struct {
	i            int
	s            string
	line, column int
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch      chan frame
	ch_stop chan bool
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []frame
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	// Since then, I introduced the built-in Line() and Column() functions.
	l, c int

	parseResult interface{}

	// The following line makes it easy for scripts to insert fields in the
	// generated code.
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	yylex := new(Lexer)
	if initFun != nil {
		initFun(yylex)
	}
	yylex.ch = make(chan frame)
	yylex.ch_stop = make(chan bool, 1)
	var scan func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int)
	scan = func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		stopped := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				lcUpdate := func(r rune) {
					if r == '\n' {
						line++
						column = 0
					} else {
						column++
					}
				}
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					lcUpdate(buf[0])
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					select {
					case ch <- frame{matchi, text, line, column}:
						{
						}
					case stopped = <-ch_stop:
						{
						}
					}
					if stopped {
						break
					}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, ch_stop, family[matchi].nest, line, column)
					}
					if atEOF {
						break
					}
					for _, r := range text {
						lcUpdate(r)
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- frame{-1, "", line, column}
	}
	go scan(bufio.NewReader(in), yylex.ch, yylex.ch_stop, dfas, 0, 0)
	return yylex
}

type dfa struct {
	acc          []bool           // Accepting states.
	f            []func(rune) int // Transitions.
	startf, endf []int            // Transitions at start and end of input.
	nest         []dfa
}

var dfas = []dfa{
	// #EXTM3U
	{[]bool{false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 35:
				return 1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return 2
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return 4
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return 5
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return 6
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return 7
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			case 51:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-VERSION:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return 3
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return 9
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return 10
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 11
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return 12
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return 13
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return 14
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 15
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 16
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-INDEPENDENT-SEGMENTS
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return 9
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 10
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return 11
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 12
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return 13
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 14
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 15
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return 16
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 17
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 18
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 19
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 20
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 21
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 22
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return 23
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return 24
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return 25
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 26
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 27
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 28
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-MEDIA:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return 9
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 10
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return 11
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return 12
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 13
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 14
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-STREAM-INF:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return 3
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return 9
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 10
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return 11
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return 12
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 13
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return 14
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 15
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return 16
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 17
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return 18
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 19
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-TARGETDURATION:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return 9
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 10
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 11
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return 12
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 13
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return 14
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return 15
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return 16
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 17
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 18
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return 19
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return 20
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return 21
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return 22
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 23
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-SERVER-CONTROL:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return 3
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return 9
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return 10
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 11
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return 12
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return 13
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 14
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 15
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return 16
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 17
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return 18
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 19
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 20
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 21
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return 22
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 23
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-PART-INF:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return 3
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return 9
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 10
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return 11
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 12
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 13
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return 14
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return 15
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return 16
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 17
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-MEDIA-SEQUENCE:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return 9
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 10
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return 11
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return 12
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 13
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 14
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return 15
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 16
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return 17
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return 18
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 19
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return 20
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return 21
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 22
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 23
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 81:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-SKIP:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return 3
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return 9
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return 10
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return 11
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return 12
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 13
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXTINF:
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return 3
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return 6
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return 7
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return 8
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return 9
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 58:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-PROGRAM-DATE-TIME:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return 9
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 10
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return 11
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return 12
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 13
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 14
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return 15
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 16
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return 17
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 18
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 19
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 20
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 21
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 22
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return 23
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return 24
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 25
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 26
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 77:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-PART:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return 3
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return 9
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 10
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return 11
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 12
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 13
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-PRELOAD-HINT:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return 9
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 10
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 11
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return 12
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 13
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 14
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return 15
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 16
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return 17
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return 18
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return 19
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 20
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 21
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-RENDITION-REPORT:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 9
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return 10
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 11
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return 12
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return 13
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 14
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return 15
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return 16
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 17
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 18
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 19
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return 20
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return 21
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return 22
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 23
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 24
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 25
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n#EXT-X-MAP:
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return 3
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return 5
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 6
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return 8
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return 9
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return 10
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return 11
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return 12
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return -1
			case 45:
				return -1
			case 58:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n[A-Za-z][^\"\n, #=]+
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return -1
			case 97 <= r && r <= 122:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 3
			case 97 <= r && r <= 122:
				return 3
			}
			return 3
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 3
			case 97 <= r && r <= 122:
				return 3
			}
			return 3
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// \n[ \t]*
	{[]bool{false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 9:
				return -1
			case 10:
				return 1
			case 32:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return 2
			case 10:
				return -1
			case 32:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return 2
			case 10:
				return -1
			case 32:
				return 2
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// \n#[^(EXT)].*
	{[]bool{false, false, false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			case 35:
				return -1
			case 40:
				return -1
			case 41:
				return -1
			case 69:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 35:
				return 2
			case 40:
				return -1
			case 41:
				return -1
			case 69:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return 3
			case 35:
				return 3
			case 40:
				return -1
			case 41:
				return -1
			case 69:
				return -1
			case 84:
				return -1
			case 88:
				return -1
			}
			return 3
		},
		func(r rune) int {
			switch r {
			case 10:
				return 4
			case 35:
				return 4
			case 40:
				return 4
			case 41:
				return 4
			case 69:
				return 4
			case 84:
				return 4
			case 88:
				return 4
			}
			return 4
		},
		func(r rune) int {
			switch r {
			case 10:
				return 4
			case 35:
				return 4
			case 40:
				return 4
			case 41:
				return 4
			case 69:
				return 4
			case 84:
				return 4
			case 88:
				return 4
			}
			return 4
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// BANDWIDTH=
	{[]bool{false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return 1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 2
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return 3
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return 4
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return 6
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return 7
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return 8
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return 9
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 10
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 84:
				return -1
			case 87:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// AVERAGE-BANDWIDTH=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return 2
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return 3
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return 4
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 5
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return 6
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return 7
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 8
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return 9
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 10
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return 11
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return 12
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return 13
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return 14
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return 15
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return 16
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return 17
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 18
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 72:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 86:
				return -1
			case 87:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// RESOLUTION=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return 2
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return 3
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 4
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return 5
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return 6
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 7
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return 8
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 9
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return 10
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 11
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// FRAME-RATE=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return 1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return 2
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 3
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return 4
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return 5
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 6
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return 7
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 8
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return 10
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 11
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 77:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// CODECS=
	{[]bool{false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return 1
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return 2
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return -1
			case 68:
				return 3
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 4
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return 5
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return 6
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 7
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// AUDIO=
	{[]bool{false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 1
			case 68:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 85:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return 3
			case 73:
				return -1
			case 79:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return 4
			case 79:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 79:
				return 5
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 6
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, nil},

	// TYPE=
	{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 84:
				return 1
			case 89:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 89:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 80:
				return 3
			case 84:
				return -1
			case 89:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return 4
			case 80:
				return -1
			case 84:
				return -1
			case 89:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 5
			case 69:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 89:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 69:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			case 89:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

	// GROUP-ID=
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return 1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 2
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return 3
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return 5
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 6
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return 7
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return 8
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 9
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// NAME=
	{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 78:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 2
			case 69:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return 3
			case 78:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return 4
			case 77:
				return -1
			case 78:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 5
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

	// DEFAULT=
	{[]bool{false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return 1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return 2
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return 3
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 4
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return 6
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return 7
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 8
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 70:
				return -1
			case 76:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// AUTOSELECT=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return 3
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return 4
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return 5
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return 6
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return 7
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return 8
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return 9
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return 10
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 11
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// LANGUAGE=
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return 1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 2
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return 3
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return 4
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 6
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return 7
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return 8
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 9
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// CHANNELS=
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return 1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return 2
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 3
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return 4
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return 5
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return 6
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return 7
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return 8
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 9
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 69:
				return -1
			case 72:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// URI=
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 73:
				return -1
			case 82:
				return -1
			case 85:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 73:
				return -1
			case 82:
				return 2
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 73:
				return 3
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 4
			case 73:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 73:
				return -1
			case 82:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// CAN-BLOCK-RELOAD=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return 1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return 3
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 4
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return 5
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return 6
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 7
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return 8
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return 9
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 10
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return 12
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return 13
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return 14
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 15
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return 16
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 17
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// CAN-SKIP-UNTIL=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return 1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return 3
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 4
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 5
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return 6
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return 7
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return 8
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 9
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return 10
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return 11
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 12
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return 13
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return 14
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 15
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 67:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// PART-HOLD-BACK=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return 1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return 3
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 5
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return 6
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return 7
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return 8
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return 9
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 10
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return 11
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 12
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return 13
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return 14
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 15
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 66:
				return -1
			case 67:
				return -1
			case 68:
				return -1
			case 72:
				return -1
			case 75:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// PART-TARGET=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return 1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return 3
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 5
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 6
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 7
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return 8
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return 9
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return 10
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 12
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// SKIPPED-SEGMENTS=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return 2
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return 3
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return 4
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return 5
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 6
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return 7
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 8
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 9
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 10
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return 11
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return 12
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 13
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return 14
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return 15
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return 16
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 17
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 71:
				return -1
			case 73:
				return -1
			case 75:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// DURATION=
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return 1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return 3
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return 4
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return 5
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return 6
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return 7
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return 8
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 9
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 65:
				return -1
			case 68:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 79:
				return -1
			case 82:
				return -1
			case 84:
				return -1
			case 85:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// INDEPENDENT=
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return 1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 2
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return 3
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 4
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return 5
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 6
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 7
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return 8
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return 9
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return 10
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 12
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 68:
				return -1
			case 69:
				return -1
			case 73:
				return -1
			case 78:
				return -1
			case 80:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// LAST-MSN=
	{[]bool{false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return 1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return 3
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 5
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return 6
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return 7
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return 8
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 9
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 77:
				return -1
			case 78:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// LAST-PART=
	{[]bool{false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return 1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 2
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return 3
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 5
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return 6
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return 7
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return 8
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return 10
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			case 65:
				return -1
			case 76:
				return -1
			case 80:
				return -1
			case 82:
				return -1
			case 83:
				return -1
			case 84:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// [A-Za-z\-]+=
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return 1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 1
			case 97 <= r && r <= 122:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 1
			case 61:
				return 2
			}
			switch {
			case 65 <= r && r <= 90:
				return 1
			case 97 <= r && r <= 122:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return -1
			case 97 <= r && r <= 122:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// [0-9]+\-[0-9]+\-[0-9]+T[0-9]+:[0-9]+:[0-9]+\.[0-9]+Z[+-][0-9]+:[0-9]+
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return 2
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return 4
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return 6
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return 8
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return 10
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return 12
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 13
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return 14
			}
			switch {
			case 48 <= r && r <= 57:
				return 13
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return 15
			case 45:
				return 15
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 16
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return 17
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 16
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 18
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 18
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// [0-9]+\-[0-9]+\-[0-9]+T[0-9]+:[0-9]+:[0-9]+\.[0-9]+Z
	{[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 2
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 4
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return 6
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return 8
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 7
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return 10
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return 12
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 13
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return 14
			}
			switch {
			case 48 <= r && r <= 57:
				return 13
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			case 58:
				return -1
			case 84:
				return -1
			case 90:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// [0-9]+\.[0-9]*
	{[]bool{false, false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 46:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// -[0-9]+\.[0-9]*
	{[]bool{false, false, false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return 1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return 3
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 4
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// [0-9]+x[0-9]+
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 120:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// [0-9]+
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// 0[xX][0-9A-Fa-f]+
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 48:
				return 1
			case 88:
				return -1
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			case 65 <= r && r <= 70:
				return -1
			case 97 <= r && r <= 102:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 48:
				return -1
			case 88:
				return 2
			case 120:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			case 65 <= r && r <= 70:
				return -1
			case 97 <= r && r <= 102:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 48:
				return 3
			case 88:
				return -1
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			case 65 <= r && r <= 70:
				return 3
			case 97 <= r && r <= 102:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 48:
				return 3
			case 88:
				return -1
			case 120:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			case 65 <= r && r <= 70:
				return 3
			case 97 <= r && r <= 102:
				return 3
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// \"[^\"\n\r]+\"
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 34:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 34:
				return -1
			}
			return 2
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 34:
				return 3
			}
			return 2
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 34:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// [A-Za-z][^\"\n, #=]+
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 1
			case 97 <= r && r <= 122:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return 2
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 32:
				return -1
			case 34:
				return -1
			case 35:
				return -1
			case 44:
				return -1
			case 61:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return 2
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// ,
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 44:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 44:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},
}

func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}

func (yyLex *Lexer) Stop() {
	yyLex.ch_stop <- true
}

// Text returns the matched text.
func (yylex *Lexer) Text() string {
	return yylex.stack[len(yylex.stack)-1].s
}

// Line returns the current line number.
// The first line is 0.
func (yylex *Lexer) Line() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].line
}

// Column returns the current column number.
// The first column is 0.
func (yylex *Lexer) Column() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].column
}

func (yylex *Lexer) next(lvl int) int {
	if lvl == len(yylex.stack) {
		l, c := 0, 0
		if lvl > 0 {
			l, c = yylex.stack[lvl-1].line, yylex.stack[lvl-1].column
		}
		yylex.stack = append(yylex.stack, frame{0, "", l, c})
	}
	if lvl == len(yylex.stack)-1 {
		p := &yylex.stack[lvl]
		*p = <-yylex.ch
		yylex.stale = false
	} else {
		yylex.stale = true
	}
	return yylex.stack[lvl].i
}
func (yylex *Lexer) pop() {
	yylex.stack = yylex.stack[:len(yylex.stack)-1]
}
func (yylex Lexer) Error(e string) {
	panic(e)
}

// Lex runs the lexer. Always returns 0.
// When the -s option is given, this function is not generated;
// instead, the NN_FUN macro runs the lexer.
func (yylex *Lexer) Lex(lval *yySymType) int {
OUTER0:
	for {
		switch yylex.next(0) {
		case 0:
			{
				lval.i = TAG_EXTM3U
				return lval.i
			}
		case 1:
			{
				lval.i = TAG_EXT_X_VERSION
				return lval.i
			}
		case 2:
			{
				lval.i = TAG_EXT_X_INDEPENDENT_SEGMENTS
				return lval.i
			}
		case 3:
			{
				lval.i = TAG_EXT_X_MEDIA
				return lval.i
			}
		case 4:
			{
				lval.i = TAG_EXT_X_STREAM_INF
				return lval.i
			}
		case 5:
			{
				lval.i = TAG_EXT_X_TARGETDURATION
				return lval.i
			}
		case 6:
			{
				lval.i = TAG_EXT_X_SERVER_CONTROL
				return lval.i
			}
		case 7:
			{
				lval.i = TAG_EXT_X_PART_INF
				return lval.i
			}
		case 8:
			{
				lval.i = TAG_EXT_X_MEDIA_SEQUENCE
				return lval.i
			}
		case 9:
			{
				lval.i = TAG_EXT_X_SKIP
				return lval.i
			}
		case 10:
			{
				lval.i = TAG_EXTINF
				return lval.i
			}
		case 11:
			{
				lval.i = TAG_EXT_X_PROGRAM_DATE_TIME
				return lval.i
			}
		case 12:
			{
				lval.i = TAG_EXT_X_PART
				return lval.i
			}
		case 13:
			{
				lval.i = TAG_EXT_X_PRELOAD_HINT
				return lval.i
			}
		case 14:
			{
				lval.i = TAG_EXT_X_RENDITION_REPORT
				return lval.i
			}
		case 15:
			{
				lval.i = TAG_EXT_X_MAP
				return lval.i
			}
		case 16:
			{
				t := yylex.Text()
				lval.s = t[1:]
				return SECONDLINEVALUE
			}
		case 17:
			{ /* ignore empty line */
			}
		case 18:
			{ /* ignore #comment lines */
			}
		case 19:
			{
				lval.i = ATTR_BANDWIDTH
				return lval.i
			}
		case 20:
			{
				lval.i = ATTR_AVERAGE_BANDWIDTH
				return lval.i
			}
		case 21:
			{
				lval.i = ATTR_RESOLUTION
				return lval.i
			}
		case 22:
			{
				lval.i = ATTR_FRAME_RATE
				return lval.i
			}
		case 23:
			{
				lval.i = ATTR_CODECS
				return lval.i
			}
		case 24:
			{
				lval.i = ATTR_AUDIO
				return lval.i
			}
		case 25:
			{
				lval.i = ATTR_TYPE
				return lval.i
			}
		case 26:
			{
				lval.i = ATTR_GROUP_ID
				return lval.i
			}
		case 27:
			{
				lval.i = ATTR_NAME
				return lval.i
			}
		case 28:
			{
				lval.i = ATTR_DEFAULT
				return lval.i
			}
		case 29:
			{
				lval.i = ATTR_AUTOSELECT
				return lval.i
			}
		case 30:
			{
				lval.i = ATTR_LANGUAGE
				return lval.i
			}
		case 31:
			{
				lval.i = ATTR_CHANNELS
				return lval.i
			}
		case 32:
			{
				lval.i = ATTR_URI
				return lval.i
			}
		case 33:
			{
				lval.i = ATTR_CAN_BLOCK_RELOAD
				return lval.i
			}
		case 34:
			{
				lval.i = ATTR_CAN_SKIP_UNTIL
				return lval.i
			}
		case 35:
			{
				lval.i = ATTR_PART_HOLD_BACK
				return lval.i
			}
		case 36:
			{
				lval.i = ATTR_PART_TARGET
				return lval.i
			}
		case 37:
			{
				lval.i = ATTR_SKIPPED_SEGMENTS
				return lval.i
			}
		case 38:
			{
				lval.i = ATTR_DURATION
				return lval.i
			}
		case 39:
			{
				lval.i = ATTR_INDEPENDENT
				return lval.i
			}
		case 40:
			{
				lval.i = ATTR_LAST_MSN
				return lval.i
			}
		case 41:
			{
				lval.i = ATTR_LAST_PART
				return lval.i
			}
		case 42:
			{
				t := yylex.Text()
				lval.s = t[0 : len(t)-1]
				return ATTRKEY
			}
		case 43:
			{
				lval.t, _ = time.Parse(time.RFC3339, yylex.Text())
				return TIMEVAL
			}
		case 44:
			{
				lval.t, _ = time.Parse(time.RFC3339, yylex.Text())
				return TIMEVAL
			}
		case 45:
			{
				lval.f, _ = strconv.ParseFloat(yylex.Text(), 64)
				return FLOATVAL
			}
		case 46:
			{
				lval.f, _ = strconv.ParseFloat(yylex.Text(), 64)
				return FLOATVAL
			}
		case 47:
			{
				lval.r = yylex.Text()
				return RESOLUTIONVAL
			}
		case 48:
			{
				lval.i64, _ = strconv.ParseInt(yylex.Text(), 10, 64)
				return INTEGERVAL
			}
		case 49:
			{
				lval.i64, _ = strconv.ParseInt(yylex.Text(), 16, 64)
				return INTEGERVAL
			}
		case 50:
			{
				t := yylex.Text()
				lval.s = t[1 : len(t)-2]
				return STRINGVAL
			}
		case 51:
			{
				t := yylex.Text()
				lval.s = t
				return STRINGVAL
			}
		case 52:
			{
				lval.i = COMMA
				return lval.i
			}
		default:
			break OUTER0
		}
		continue
	}
	yylex.pop()

	return 0
}
func main() {
	yyParse(NewLexer(os.Stdin))
}
