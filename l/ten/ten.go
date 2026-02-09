package ten

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type options struct {
	K            int  // -k N
	Numeric      bool // -n
	Reverse      bool // -r
	Unique       bool // -u
	Month        bool // -M
	TrimTrailing bool // -b
	CheckOnly    bool // -c
	Human        bool // -h
}

type record struct {
	line     string
	key      string
	lineTrim string
	keyTrim  string
	numVal   float64
	humanVal float64
	monthVal int
	flags    uint8
}

const (
	flagNumOK uint8 = 1 << iota
	flagHumanOK
	flagMonthOK
)

var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

var unitPow = map[byte]int{
	'k': 1, 'm': 2, 'g': 3, 't': 4, 'p': 5, 'e': 6,
}

func ten() {
	opts, files, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		usage()
		os.Exit(2)
	}

	recs, err := readRecords(files, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(2)
	}

	cmp := makeComparator(opts)

	if opts.CheckOnly {
		if ok, idx := isSortedRecords(recs, cmp, opts.Reverse); !ok {
			fmt.Fprintf(os.Stderr, "not sorted: line %d: %q > %q\n", idx+1, recs[idx].line, recs[idx+1].line)
			os.Exit(1)
		}
		os.Exit(0)
	}

	idxs := make([]int, len(recs))
	for i := range idxs {
		idxs[i] = i
	}

	sort.SliceStable(idxs, func(i, j int) bool {
		c := cmp(&recs[idxs[i]], &recs[idxs[j]])
		if c == 0 {
			return false
		}
		if opts.Reverse {
			return c > 0
		}
		return c < 0
	})

	w := bufio.NewWriter(os.Stdout)
	defer func() { _ = w.Flush() }()

	if opts.Unique {
		if len(idxs) == 0 {
			return
		}
		prev := idxs[0]
		if _, err := w.WriteString(recs[prev].line); err != nil {
			fmt.Fprintln(os.Stderr, "write error:", err)
			os.Exit(1)
		}
		_ = w.WriteByte('\n')
		for k := 1; k < len(idxs); k++ {
			cur := idxs[k]
			if equalForUniqueness(&recs[prev], &recs[cur], opts) {
				continue
			}
			if _, err := w.WriteString(recs[cur].line); err != nil {
				fmt.Fprintln(os.Stderr, "write error:", err)
				os.Exit(1)
			}
			_ = w.WriteByte('\n')
			prev = cur
		}
		return
	}

	for _, id := range idxs {
		if _, err := w.WriteString(recs[id].line); err != nil {
			fmt.Fprintln(os.Stderr, "write error:", err)
			os.Exit(1)
		}
		_ = w.WriteByte('\n')
	}
}

func parseArgs(args []string) (options, []string, error) {
	var opts options
	files := []string{}
	for i := range len(args) {
		a := args[i]
		if a == "--" {
			files = append(files, args[i+1:]...)
			break
		}
		if a == "-" {
			files = append(files, a)
			continue
		}
		if strings.HasPrefix(a, "-") && len(a) > 1 {
			for j := 1; j < len(a); j++ {
				switch a[j] {
				case 'n':
					opts.Numeric = true
				case 'r':
					opts.Reverse = true
				case 'u':
					opts.Unique = true
				case 'M':
					opts.Month = true
				case 'b':
					opts.TrimTrailing = true
				case 'c':
					opts.CheckOnly = true
				case 'h':
					opts.Human = true
				case 'k':
					if j+1 < len(a) {
						param := a[j+1:]
						k, err := strconv.Atoi(param)
						if err != nil || k <= 0 {
							return opts, nil, fmt.Errorf("invalid -k value: %q", param)
						}
						opts.K = k
						j = len(a)
						break
					}
					i++
					if i >= len(args) {
						return opts, nil, errors.New("missing argument for -k")
					}
					param := args[i]
					k, err := strconv.Atoi(param)
					if err != nil || k <= 0 {
						return opts, nil, fmt.Errorf("invalid -k value: %q", param)
					}
					opts.K = k
				default:
					return opts, nil, fmt.Errorf("invalid option: -%c", a[j])
				}
			}
			continue
		}
		files = append(files, a)
	}
	return opts, files, nil
}

func readRecords(files []string, opts options) ([]record, error) {
	var recs []record
	if len(files) == 0 {
		if err := scanToRecords(os.Stdin, opts, &recs); err != nil {
			return nil, err
		}
		return recs, nil
	}
	for _, fn := range files {
		if fn == "-" {
			if err := scanToRecords(os.Stdin, opts, &recs); err != nil {
				return nil, err
			}
			continue
		}
		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("cannot open %s: %w", fn, err)
		}
		if err := scanToRecords(f, opts, &recs); err != nil {
			_ = f.Close()
			return nil, fmt.Errorf("reading %s: %w", fn, err)
		}
		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("closing %s: %w", fn, err)
		}
	}
	return recs, nil
}

func scanToRecords(r io.Reader, opts options, out *[]record) error {
	s := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, 10*1024*1024)
	for s.Scan() {
		ln := s.Text()
		*out = append(*out, makeRecord(ln, opts))
	}
	return s.Err()
}

func makeRecord(line string, opts options) record {
	var r record
	r.line = line
	if opts.K > 0 {
		r.key = extractColumn(line, opts.K)
	} else {
		r.key = line
	}
	if opts.TrimTrailing {
		r.keyTrim = strings.TrimRight(r.key, " \t")
		r.lineTrim = strings.TrimRight(line, " \t")
	} else {
		r.keyTrim = r.key
		r.lineTrim = line
	}
	if opts.Numeric {
		if v, err := strconv.ParseFloat(strings.TrimSpace(r.keyTrim), 64); err == nil {
			r.numVal = v
			r.flags |= flagNumOK
		}
	}
	if opts.Human {
		if v, ok := parseHuman(strings.TrimSpace(r.keyTrim)); ok {
			r.humanVal = v
			r.flags |= flagHumanOK
		}
	}
	if opts.Month {
		if v, ok := parseMonth(strings.TrimSpace(r.keyTrim)); ok {
			r.monthVal = v
			r.flags |= flagMonthOK
		}
	}
	return r
}

func extractColumn(s string, k int) string {
	if k <= 0 {
		return s
	}
	start := 0
	for col := 1; col < k; col++ {
		i := strings.IndexByte(s[start:], '\t')
		if i < 0 {
			return ""
		}
		start += i + 1
		if start >= len(s) {
			return ""
		}
	}
	end := strings.IndexByte(s[start:], '\t')
	if end < 0 {
		return s[start:]
	}
	return s[start : start+end]
}

func parseHuman(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	last := s[len(s)-1]
	power := 0
	if ('A' <= last && last <= 'Z') || ('a' <= last && last <= 'z') {
		last |= 0x20
		p, ok := unitPow[last]
		if !ok {
			return 0, false
		}
		power = p
		s = strings.TrimSpace(s[:len(s)-1])
		if s == "" {
			return 0, false
		}
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	if power > 0 {
		v *= math.Pow(1024, float64(power))
	}
	return v, true
}

func parseMonth(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 3 {
		return 0, false
	}
	m := strings.ToLower(s[:3])
	if v, ok := monthMap[m]; ok {
		return v, true
	}
	return 0, false
}

func makeComparator(opts options) func(a, b *record) int {
	if opts.Month {
		return func(a, b *record) int {
			if a.flags&flagMonthOK != 0 && b.flags&flagMonthOK != 0 {
				if a.monthVal < b.monthVal {
					return -1
				} else if a.monthVal > b.monthVal {
					return 1
				}
				return 0
			}
			if c := strings.Compare(a.keyTrim, b.keyTrim); c != 0 {
				return c
			}
			return strings.Compare(a.lineTrim, b.lineTrim)
		}
	}
	if opts.Human {
		return func(a, b *record) int {
			if a.flags&flagHumanOK != 0 && b.flags&flagHumanOK != 0 {
				if a.humanVal < b.humanVal {
					return -1
				} else if a.humanVal > b.humanVal {
					return 1
				}
				return 0
			}
			if c := strings.Compare(a.keyTrim, b.keyTrim); c != 0 {
				return c
			}
			return strings.Compare(a.lineTrim, b.lineTrim)
		}
	}
	if opts.Numeric {
		return func(a, b *record) int {
			if a.flags&flagNumOK != 0 && b.flags&flagNumOK != 0 {
				if a.numVal < b.numVal {
					return -1
				} else if a.numVal > b.numVal {
					return 1
				}
				return 0
			}
			if c := strings.Compare(a.keyTrim, b.keyTrim); c != 0 {
				return c
			}
			return strings.Compare(a.lineTrim, b.lineTrim)
		}
	}
	return func(a, b *record) int {
		if c := strings.Compare(a.keyTrim, b.keyTrim); c != 0 {
			return c
		}
		return strings.Compare(a.lineTrim, b.lineTrim)
	}
}

func isSortedRecords(recs []record, cmp func(a, b *record) int, reverse bool) (bool, int) {
	for i := 0; i+1 < len(recs); i++ {
		c := cmp(&recs[i], &recs[i+1])
		if reverse {
			if c < 0 {
				return false, i
			}
		} else {
			if c > 0 {
				return false, i
			}
		}
	}
	return true, -1
}

func equalForUniqueness(a, b *record, opts options) bool {
	if opts.TrimTrailing {
		return a.lineTrim == b.lineTrim
	}
	return a.line == b.line
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: sort [-n] [-r] [-u] [-M] [-b] [-c] [-h] [-k N] [file ...]")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -k N  sort by N-th tab-separated column (1-based)")
	fmt.Fprintln(os.Stderr, "  -n    numeric sort")
	fmt.Fprintln(os.Stderr, "  -r    reverse order")
	fmt.Fprintln(os.Stderr, "  -u    unique lines")
	fmt.Fprintln(os.Stderr, "  -M    month sort (Jan,Feb,...)")
	fmt.Fprintln(os.Stderr, "  -b    ignore trailing blanks when comparing/unique")
	fmt.Fprintln(os.Stderr, "  -c    check if input is sorted (no output on success)")
	fmt.Fprintln(os.Stderr, "  -h    human numeric sort (K/M/G/T base 1024)")
}
