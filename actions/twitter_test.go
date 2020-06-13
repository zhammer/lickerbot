package actions

import (
	"bytes"
	"encoding/json"
	"lickerbot/models"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/httptest"
)

const (
	webhookURL = "/webhook/twitter"
)

func (as *ActionSuite) Test_Twitter_Webhook_SecurityCheck() {
	res := as.JSON(webhookURL + "?crc_token=abcde").Get()

	as.Equal(http.StatusOK, res.Result().StatusCode)
	var responseBody map[string]string
	json.Unmarshal(res.Body.Bytes(), &responseBody)

	as.Regexp(regexp.MustCompile(`sha256=.{43}=`), responseBody["response_token"])

}

// we return 200s for any webhook request that we don't care about.
// lazy about this so even invalid messages (like an empty message)
// return a 200.
func (as *ActionSuite) Test_Twitter_Webhook_EmptyMessage() {
	res := as.JSON(webhookURL).Post(nil)

	as.Equal(http.StatusOK, res.Result().StatusCode)
}

// We only care about tweet create events. This is a follow event.
// https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/account-activity-data-objects#favorite_events
func (as *ActionSuite) Test_Twitter_Webhook_NotTweetCreateEvent() {
	followEventBody := `{
		"for_user_id": "1269307350520868866",
		"follow_events": []
	}`

	res := rawRequest(as, "POST", followEventBody)

	as.Equal(http.StatusOK, res.Result().StatusCode)
}

func (as *ActionSuite) Test_Twitter_Webhook_ForDifferentUser() {
	eventBody := `{
		"for_user_id": "123",
		"tweet_create_events": []
	}`

	res := rawRequest(as, "POST", eventBody)

	as.Equal(http.StatusOK, res.Result().StatusCode)
}

func (as *ActionSuite) Test_Twitter_Webhook_SomeActionableTweets() {
	// for now i don't know a better way to do this. recreate our app after
	// setting up a mock twitter server and setting that url in envy (which
	// acts as a global key/value store)
	twitterRequests := make(chan *http.Request)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// for now, we parse the form while the request reader is live. if needed we can clone
		// the reader, maybe? idk
		r.ParseForm()
		twitterRequests <- r
	}))
	defer server.Close()
	envy.Set("TWITTER_BASE_URL", server.URL)
	app = nil
	as.App = App()

	// user @sputnik (6666) has licked the boot 1 time already and people have donated $50.
	bootlicker := &models.Bootlicker{
		TwitterUserID: 6666,
		TwitterHandle: "sputnik",
	}
	as.DB.Create(bootlicker)
	lick := &models.Lick{
		TweetID:      4444,
		BootlickerID: bootlicker.ID,
	}
	as.DB.Create(lick)
	pledgedDonation := &models.PledgedDonation{
		Amount:       50,
		BootlickerID: bootlicker.ID,
	}
	as.DB.Create(pledgedDonation)

	// webhook contains several tweets:
	//   1: mentions @lickerbot but is not a reply
	//   2: is a reply but does not mention anyone
	//   3: is a reply but mentions someone other than @lickerbot
	//   4: is a reply that mentions @lickerbot but is a reply *to* lickerbot
	//   5: is an actionable tweet
	//   6: is an actionable tweet, the bootlicker already exists
	eventBody := `{
		"for_user_id": "1269307350520868866",
		"tweet_create_events": [
			{
				"id": 1,
				"in_reply_to_status_id": null,
				"in_reply_to_user_id": null,
				"in_reply_to_screen_name": null,
				"text": "Check out @lickerbot",
				"user": {
					"screen_name": "poster"
				},
				"entities": {
					"user_mentions": [
						{
							"id": 1269307350520868866
						}
					]
				}
			},
			{
				"id": 2,
				"in_reply_to_status_id": 1234,
				"in_reply_to_user_id": 1234,
				"in_reply_to_screen_name": "weird_guy",
				"text": "I'm replying to weird_guy",
				"user": {
					"screen_name": "poster"
				},
				"entities": {
					"user_mentions": []
				}
			},
			{
				"id": 3,
				"in_reply_to_status_id": 1234,
				"in_reply_to_user_id": 1234,
				"in_reply_to_screen_name": "weird_guy",
				"text": "Check out @threadreaderapp",
				"user": {
					"screen_name": "poster"
				},
				"entities": {
					"user_mentions": [
						{
							"id": 123
						}
					]
				}
			},
			{
				"id": 4,
				"in_reply_to_status_id": 1234,
				"in_reply_to_user_id": 1269307350520868866,
				"in_reply_to_screen_name": "lickerbot",
				"text": "That's cool @lickerbot",
				"user": {
					"screen_name": "poster"
				},
				"entities": {
					"user_mentions": [
						{
							"id": 1269307350520868866
						}
					]
				}
			},
			{
				"id": 5,
				"in_reply_to_status_id": 1234,
				"in_reply_to_user_id": 1234,
				"in_reply_to_screen_name": "weird_guy",
				"text": "@lickerbot",
				"user": {
					"screen_name": "poster"
				},
				"entities": {
					"user_mentions": [
						{
							"id": 1269307350520868866
						}
					]
				}
			},
			{
				"id": 6,
				"in_reply_to_status_id": 8888,
				"in_reply_to_user_id": 6666,
				"in_reply_to_screen_name": "sputnik",
				"text": "@lickerbot check this out",
				"user": {
					"screen_name": "someone"
				},
				"entities": {
					"user_mentions": [
						{
							"id": 1269307350520868866
						}
					]
				}
			}
		]
	}`

	res := rawRequest(as, "POST", eventBody)

	as.Equal(http.StatusAccepted, res.Result().StatusCode)

	var responseBody map[string][]int64
	json.Unmarshal(res.Body.Bytes(), &responseBody)
	as.Equal(map[string][]int64{"ingested_tweet_ids": {5, 6}}, responseBody)

	// verify correct tweets sent to twitter
	tweetRequests := receiveTwitterRequests(twitterRequests, 2, 5*time.Second)
	as.Len(tweetRequests, 2)
	for _, request := range tweetRequests {
		as.Equal("/status/update.json", request.URL.Path)
		as.Equal("post", strings.ToLower(request.Method))
	}

	var tweetRequestBodies []url.Values
	for _, request := range tweetRequests {
		tweetRequestBodies = append(tweetRequestBodies, request.Form)
	}
	// sort the requests by "in_reply_to_status_id"
	sort.SliceStable(tweetRequestBodies, func(i int, j int) bool {
		return tweetRequestBodies[i]["in_reply_to_status_id"][0] < tweetRequestBodies[j]["in_reply_to_status_id"][0]
	})

	expected := []url.Values{
		{
			"in_reply_to_status_id": []string{"5"},
			"status":                []string{"@poster Must taste good. weird_guy has licked the boot 1 time. Will you pledge to donate $5 per lick to organizations fighting police brutality? https://lickerbot.com/@weird_guy"},
		},
		{
			"in_reply_to_status_id": []string{"6"},
			"status":                []string{"@someone That's one hell of a lick! sputnik has licked the boot 2 times. In response, people have donated $50 to organizations fighting police brutality. That's $25 per lick! https://lickerbot.com/@sputnik"},
		},
	}
	as.Equal(expected, tweetRequestBodies)

	// for good measure, let's make sure the length of these tweets is well under the 280 limit
	for _, requestBody := range tweetRequestBodies {
		as.Less(len(requestBody["status"][0]), 240)
	}

}

// was having a ton of trouble sending json strings as requests
func rawRequest(as *ActionSuite, method string, body string) *httptest.JSONResponse {
	jsonHelper := as.JSON(webhookURL)
	req, _ := http.NewRequest("POST", jsonHelper.URL, bytes.NewReader([]byte(body)))
	return jsonHelper.Perform(req)
}

// recieve requests to the `requestChannel` until we received `expected` requests or `timeout` elapses.
// returns whatever requests have been received.
func receiveTwitterRequests(requestChannel chan *http.Request, expected int, timeout time.Duration) []*http.Request {
	var received []*http.Request

	// listen for new requests, add them to received
	receivedExpectedRequests := make(chan bool)
	go (func() {
		for i := 0; i < expected; i++ {
			request := <-requestChannel
			received = append(received, request)
		}
		receivedExpectedRequests <- true
	})()

	// listen for requests until either:
	select {
	// we've received the expected number of requests
	case <-receivedExpectedRequests:
	// our timeout has expired
	case <-time.After(timeout):
	}

	// return whatever we were able to receive
	return received
}
