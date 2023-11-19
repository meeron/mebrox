package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/meeron/mebrox/server"
)

var (
	commitRouteRegex *regexp.Regexp
)

func Messages(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	commitParams := commitRouteRegex.FindStringSubmatch(r.URL.String())
	if len(commitParams) > 1 && r.Method == http.MethodPost {
		return commitMessage(s, w, commitParams[1])
	}

	return responseNotFound(w, "not found")
}

func commitMessage(s *server.Server, w http.ResponseWriter, id string) error {
	ok, err := s.Broker().CommitMessage(id)
	if err != nil {
		return err
	}

	if !ok {
		return responseNotFound(w, "message not found")
	}

	fmt.Fprintf(w, "ok")

	return nil
}

func init() {
	commitRouteRegex = regexp.MustCompile("/messages/(\\w+)/commit($|/)")
}
