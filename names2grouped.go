package names2grouped

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"iter"
	"os"
	"slices"
)

type Names iter.Seq[string]

type IsSameGroup func(prev string, next string) bool

type GroupedNames iter.Seq[[]string]

func (g GroupedNames) Collect() [][]string {
	return slices.Collect(iter.Seq[[]string](g))
}

func (g GroupedNames) ToJsonToWriter(wtr io.Writer) error {
	var bw *bufio.Writer = bufio.NewWriter(wtr)
	defer bw.Flush()

	var je *json.Encoder = json.NewEncoder(bw)

	for names := range g {
		e := je.Encode(names)
		if nil != e {
			return e
		}
	}

	return nil
}

func (g IsSameGroup) NamesToGrouped(names Names) GroupedNames {
	return func(yield func([]string) bool) {
		var prev string
		var buf []string
		for name := range names {
			var isSameGroup bool = g(prev, name)
			prev = name
			var groupChanged bool = !isSameGroup
			if groupChanged && 0 < len(buf) {
				var cln []string = make([]string, len(buf))
				copy(cln, buf)
				if !yield(cln) {
					return
				}
				buf = buf[:0]
			}

			buf = append(buf, name)
		}

		if 0 < len(buf) {
			yield(buf)
		}
	}
}

type Str string

func (s Str) First() string {
	for _, r := range s {
		return string(r)
	}
	return ""
}

var FirstStringCheck IsSameGroup = func(prev string, next string) bool {
	var p1 string = Str(prev).First()
	var n1 string = Str(next).First()
	return p1 == n1
}

func FirstBytesCheckNew(blen uint16) IsSameGroup {
	var bp bytes.Buffer
	var bn bytes.Buffer

	var ilen int = int(blen)

	return func(prev string, next string) bool {
		bp.Reset()
		bn.Reset()

		bp.WriteString(prev) // no error, panic on OOM
		bn.WriteString(next) // no error  panic on OOM

		var sp []byte = bp.Bytes()
		var sn []byte = bn.Bytes()

		var lpn int = min(len(sp), len(sn))
		var mlen int = min(lpn, ilen)

		var pp []byte = sp[:mlen]
		var pn []byte = sn[:mlen]

		return 0 == bytes.Compare(pp, pn)
	}
}

func ReaderToNames(rdr io.Reader) Names {
	return func(yield func(string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(rdr)
		for s.Scan() {
			var line string = s.Text()
			if !yield(line) {
				return
			}
		}
	}
}

func StdinToNames() Names { return ReaderToNames(os.Stdin) }
