package suites

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type message struct {
	ID int `json:"id"`
}

var (
	reMailLink        = regexp.MustCompile(`<a href="(.+)" class="button">.*<\/a>`)
	reMailOneTimeCode = regexp.MustCompile(`<span id="one-time-code" class="otc">([a-zA-Z0-9]+)</span>`)
)

func doGetLinkFromLastMail(t *testing.T) string {
	res := doGetLastEmail(t)

	matches := reMailLink.FindStringSubmatch(string(res))

	require.Len(t, matches, 2, "Number of match for link in email is not equal to one")

	return matches[1]
}

func doGetOneTimeCodeFromLastMail(t *testing.T) string {
	res := doGetLastEmail(t)

	matches := reMailOneTimeCode.FindStringSubmatch(string(res))

	require.Len(t, matches, 2, "Number of match for one-time code spans in email is not equal to one")

	return matches[1]
}

func doGetLastEmail(t *testing.T) []byte {
	res := doHTTPGetQuery(t, fmt.Sprintf("%s/messages", MailBaseURL))
	messages := make([]message, 0)
	err := json.Unmarshal(res, &messages)
	assert.NoError(t, err)
	assert.Greater(t, len(messages), 0)

	messageID := messages[len(messages)-1].ID

	return doHTTPGetQuery(t, fmt.Sprintf("%s/messages/%d.html", MailBaseURL, messageID))
}
