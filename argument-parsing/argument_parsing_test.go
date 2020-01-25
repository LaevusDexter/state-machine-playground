package argument_parsing

import (
	"strings"
	"testing"
)

const text = `,  ,string1"string2"[string3]string4, string5{}
			 string6[123[[string7]]][[string8]][{[string9]}]`

var single = []rune("'\"")
var separators = []rune(", \n\t")
var brackets = []rune("[]{}")

func TestParse(t *testing.T) {
	r := Parse(text, separators, single, brackets)
	for i, str := range r {
		if s := string(str); strings.Count(s, "string") != 1 || !strings.Contains(s, "string" + string(i + 1 + '0'))  {
			t.Fatal(s, i + 1)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for i:=0; i < b.N; i++ {
		_ =  Parse(text, separators, single, brackets)
	}
}