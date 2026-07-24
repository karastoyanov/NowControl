package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// confirm prompts the user with a yes/no question on stderr and reads the
// answer from stdin. When stdin is not a terminal (piped/scripted), it
// returns an error immediately instead of blocking forever, since there is
// no human present to answer — callers should use --yes in that case.
func confirm(prompt string) (bool, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return false, fmt.Errorf("confirmation required for this action; re-run with --yes to skip the prompt when running non-interactively")
	}

	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}
