package githooks

import "bytes"

func indentLines(b []byte, level int) (out []byte) {
	lineSep := []byte("\n")
	indent := bytes.Repeat([]byte("\t"), level)

	buf := bytes.NewBuffer(out)

	for _, line := range bytes.Split(b, lineSep) {
		if len(line) == 0 {
			continue
		}
		buf.Write(indent)
		buf.Write(line)
		buf.Write(lineSep)
	}

	return buf.Bytes()
}
