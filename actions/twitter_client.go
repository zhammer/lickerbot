package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dghubble/oauth1"
	"github.com/gobuffalo/envy"
)

type TwitterClient struct {
	httpClient *http.Client
	baseURL    string
}

// https://developer.twitter.com/en/docs/tweets/post-and-engage/api-reference/post-statuses-update
type TweetReplyInput struct {
	Status string
	// From docs: "This parameter will be ignored unless the author of the
	// Tweet this parameter references is mentioned within the status text.
	// Therefore, you must include @username , where username is the author
	// of the referenced Tweet, within the update."
	InReplyToStatusID int64
}

// TweetReply tweets a reply to a given status.
func (t *TwitterClient) TweetReply(input *TweetReplyInput) error {
	body := url.Values{}
	body.Set("status", input.Status)
	body.Set("in_reply_to_status_id", strconv.FormatInt(input.InReplyToStatusID, 10))

	resp, err := t.httpClient.PostForm(t.baseURL+"/status/update.json", body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Received non-200 status code from twitter: %d", resp.StatusCode)
	}

	return nil
}

func NewTwitterClient() *TwitterClient {
	config := oauth1.NewConfig(envy.Get("TWITTER_CONSUMER_KEY", ""), envy.Get("TWITTER_CONSUMER_SECRET", ""))
	token := oauth1.NewToken(envy.Get("TWITTER_ACCESS_TOKEN", ""), envy.Get("TWITTER_ACCESS_SECRET", ""))
	httpClient := config.Client(oauth1.NoContext, token)
	return &TwitterClient{
		httpClient: httpClient,
		// can overwrite this in tests
		baseURL: envy.Get("TWITTER_BASE_URL", "https://api.twitter.com/1.1"),
	}
}
