package report

import (
	"fmt"

	"github.com/jpb/barrelman/internal/cli"
)

// GitHubActionAnnotation prints a GitHub Actions annotation to stdout
func GitHubActionAnnotation(w cli.Warning) {
	title := "barrelman"

	if w.StartLine > 0 && w.EndLine > 0 {
		fmt.Printf(
			"::warning file=%s,line=%d,endLine=%d,title=%s::%s\n",
			w.Filepath,
			w.StartLine,
			w.EndLine,
			title,
			w.Message,
		)
		return
	}

	fmt.Printf(
		"::warning file=%s,title=%s::%s\n",
		w.Filepath,
		title,
		w.Message,
	)
}
