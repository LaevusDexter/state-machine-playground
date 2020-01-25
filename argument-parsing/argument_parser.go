package argument_parsing

/*
	TODO:
		stick on first matched block
*/

// state machine
type machine = [][][]byte

// state
type state = [][]byte

/*
	jumps to state according to case in context of current state
		[i]:
			<index_of_state>
*/
type jmp = []byte

/*
	executes stuff according to case in context of current state
		[i]:
			0 - skip
			1 - execute
*/
type job = []byte

/*
	jobs:
		[1] append result
		[2] inc bracket counter
		[3] dec bracket counter
		[4] append current match
*/

var stateMachine = machine{
	/*
		start
	*/
	state{ jmp{0, 3, 2, 2, 0, 0, 1}, job{0, 0, 0, 0, 0, 0, 0}, job{0, 0, 1, 1, 0, 0, 0}, job{0, 0, 0, 0, 0, 0, 0}, job{0, 0, 0, 0, 0, 0, 1} }, // 0

	/*
		regular text
	*/
	state{ jmp{1, 3, 2, 2, 1, 1, 1}, job{1, 1, 1, 1, 0, 0, 0}, job{0, 0, 1, 1, 0, 0, 0}, job{0, 0, 0, 0, 0, 0, 0}, job{0, 0, 0, 0, 1, 1, 1} }, // 1

	/*
		text in open_close brackets
	*/
	state{ jmp{2, 2, 2, 2, 2, 1, 2}, job{1, 0, 0, 0, 0, 1, 0}, job{0, 0, 1, 1, 0, 0, 0}, job{0, 0, 0, 0, 1, 1, 0}, job{1, 1, 1, 1, 1, 0, 1} }, // 2

	/*
		text in single block (quotes)
	*/
	state{ jmp{3, 1, 3, 3, 3, 3, 3}, job{0, 1, 0, 0, 0, 0, 0}, job{0, 0, 0, 0, 0, 0, 0}, job{0, 0, 0, 0, 0, 0, 0}, job{1, 0, 1, 1, 1, 1, 1} }, // 3
}

/*
	cases are referring to index of jmp and job cells
*/
const caseSeparator = 0
const caseSingleBlock = 1
const caseOpeningBracket = 2
const caseOpenBracket = 3
const caseCloseBracket = 4
const caseLastClosure = 5
const caseDefault = 6

func Parse(content string, separators []rune, single []rune, brackets []rune) (result [][]rune) {
	var (
		currentState byte
		cs int

		openBrackets, closeBrackets []rune
		bracketCounter int

		buf []rune
		cp int
 	)

	if len(brackets) % 2 != 0 {
		panic("expected even number of brackets (open and closed ones).")
	}

	for i := 0; i < len(brackets); i++ {
		if i % 2 == 0 {
			openBrackets = append(openBrackets, brackets[i])
		} else {
			closeBrackets = append(closeBrackets, brackets[i])
		}
	}

	for _, r := range content {
		switch {
		case containsRune(separators, r):
			cs = caseSeparator
		case containsRune(single, r):
			cs = caseSingleBlock
		case containsRune(openBrackets, r):
			switch bracketCounter {
			case 0:
				cs = caseOpeningBracket
			default:
				cs = caseOpenBracket
			}
		case containsRune(closeBrackets, r):
			switch bracketCounter {
			case 1:
				cs = caseLastClosure
			default:
				cs = caseCloseBracket
			}
		default:
			cs = caseDefault
		}

		state := stateMachine[currentState]

		if check(state[1][cs]) && cp != len(buf) {
			result = append(result, buf[cp:])

			cp = len(buf)
		}

		if check(state[2][cs]) {
			bracketCounter++
		}

		if check(state[3][cs]) {
			bracketCounter--
		}

		if check(state[4][cs]) {
			buf = append(buf, r)
		}

		/*
			jump to next state
		*/
		currentState = stateMachine[currentState][0][cs]
	}

	if len(buf) > cp {
		result = append(result, buf[cp:])
	}

	return
}

func check(b byte) bool {
	if b > 0 {
		return true
	}

	return false
}

func containsRune(slice []rune, r rune) bool {
	for i:=0; i < len(slice); i++ {
		if slice[i] == r {
			return true
		}
	}

	return false
}