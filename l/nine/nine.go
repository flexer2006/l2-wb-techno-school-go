package nine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func UnpackString(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var b strings.Builder
	for i, n := 0, len(s); i < n; {
		r, sz := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && sz == 1 {
			return "", fmt.Errorf("invalid UTF-8 encoding at byte %d", i)
		}

		if r == '\\' {
			if i+sz >= n {
				return "", fmt.Errorf("unfinished escape sequence at byte %d", i)
			}
			nextR, nextSz := utf8.DecodeRuneInString(s[i+sz:])
			if nextR == utf8.RuneError && nextSz == 1 {
				return "", fmt.Errorf("invalid UTF-8 encoding after escape at byte %d", i+sz)
			}
			j := i + sz + nextSz
			start := j
			for j < n {
				if rr, rrSz := utf8.DecodeRuneInString(s[j:]); rr >= '0' && rr <= '9' {
					j += rrSz
					continue
				}
				break
			}
			if start == j {
				b.WriteRune(nextR)
				i = i + sz + nextSz
				continue
			}
			cnt, err := strconv.Atoi(s[start:j])
			if err != nil {
				return "", fmt.Errorf("invalid repeat count at bytes %d..%d: %w", start, j, err)
			}
			if cnt > 0 {
				b.WriteString(strings.Repeat(string(nextR), cnt))
			}
			i = j
			continue
		}

		if r >= '0' && r <= '9' {
			return "", fmt.Errorf("digit %q without preceding symbol at byte %d", r, i)
		}

		j := i + sz
		start := j
		for j < len(s) {
			if rr, rrSz := utf8.DecodeRuneInString(s[j:]); rr >= '0' && rr <= '9' {
				j += rrSz
				continue
			}
			break
		}
		if start == j {
			b.WriteRune(r)
			i += sz
			continue
		}
		cnt, err := strconv.Atoi(s[start:j])
		if err != nil {
			return "", fmt.Errorf("invalid repeat count at bytes %d..%d: %w", start, j, err)
		}
		if cnt > 0 {
			b.WriteString(strings.Repeat(string(r), cnt))
		}
		i = j
	}

	return b.String(), nil
}
