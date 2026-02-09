package eight

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func eight() {
	var (
		server  string
		timeout time.Duration
	)

	flag.StringVar(&server, "server", "pool.ntp.org", "NTP server to query (default: pool.ntp.org)")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "NTP query timeout")
	flag.Parse()

	resp, err := ntp.QueryWithOptions(server, ntp.QueryOptions{Timeout: timeout})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := resp.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	t := time.Now().Add(resp.ClockOffset).UTC()
	fmt.Println(t.Format(time.RFC3339Nano))
}
