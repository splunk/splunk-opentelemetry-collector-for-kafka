package common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Statistic struct {
	EarliestTime string `json:"earliest_time"`
	LatestTime   string `json:"latest_time"`
	TotalEvents  string `json:"count"`
}

const (
	MAX_SEARCH_RETRIES = 20 // maximum number of retries for checking search status
)

func GetEventsFromSplunk(t *testing.T, searchQuery string, startTime string, endTimeOptional ...string) []any {
	t.Logf("-->> Splunk Search: checking events in Splunk --")
	user := GetConfigVariable("USER")
	password := GetConfigVariable("PASSWORD")
	baseURL := "https://" + GetConfigVariable("HOST") + ":" + GetConfigVariable("MANAGEMENT_PORT")
	endTime := "now"
	if len(endTimeOptional) > 0 {
		endTime = endTimeOptional[0]
	}

	// post search
	jobID := postSearchRequest(t, user, password, baseURL, searchQuery, startTime, endTime)
	err := waitForSearchCompletion(t, user, password, baseURL, jobID)
	require.NoError(t, err, "Search job did not complete in expected time")

	// get events
	reqUrl := fmt.Sprintf("%s/services/search/v2/jobs/%s/events?output_mode=json", baseURL, jobID)
	rawResults := getRawSplunkSearchResults(t, user, password, reqUrl)
	return rawResults
}

func GetStatisticsFromSplunk(t *testing.T, searchQuery string, startTime string, endTimeOptional ...string) []Statistic {
	t.Logf("-->> Splunk Search: checking statistics in Splunk --")
	user := GetConfigVariable("USER")
	password := GetConfigVariable("PASSWORD")
	baseURL := "https://" + GetConfigVariable("HOST") + ":" + GetConfigVariable("MANAGEMENT_PORT")
	endTime := "now"
	if len(endTimeOptional) > 0 {
		endTime = endTimeOptional[0]
	}

	// post search
	jobID := postSearchRequest(t, user, password, baseURL, searchQuery, startTime, endTime)
	err := waitForSearchCompletion(t, user, password, baseURL, jobID)
	require.NoError(t, err, "Search job did not complete in expected time")

	// get events
	reqUrl := fmt.Sprintf("%s/services/search/v2/jobs/%s/results?output_mode=json", baseURL, jobID)
	rawResults := getRawSplunkSearchResults(t, user, password, reqUrl)

	var results []Statistic
	for _, raw := range rawResults {
		b, err := json.Marshal(raw)
		require.NoError(t, err, "Error while marshalling result item")
		var v Statistic
		err = json.Unmarshal(b, &v)
		require.NoError(t, err, "Error while unmarshalling result item to target type")
		results = append(results, v)
	}

	return results
}

func getRawSplunkSearchResults(t *testing.T, user string, password string, reqUrl string) []any {
	t.Logf("URL: " + reqUrl)
	reqEvents, err := http.NewRequest("GET", reqUrl, nil)
	require.NoError(t, err, "Error while preparing GET request to get search results")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	reqEvents.SetBasicAuth(user, password)
	respEvents, err := client.Do(reqEvents)
	require.NoError(t, err, "Error while executing Http GET request to get search results")
	defer respEvents.Body.Close()
	t.Logf("Send Request: Got query status code: " + strconv.Itoa(respEvents.StatusCode))

	bodyEvents, err := io.ReadAll(respEvents.Body)
	require.NoError(t, err, "Error while reading response body")

	var jsonResponseEvents map[string]any
	err = json.Unmarshal(bodyEvents, &jsonResponseEvents)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")

	//t.Logf("%+v", jsonResponseEvents) // debug
	rawResults := jsonResponseEvents["results"].([]any)

	return rawResults
}

func checkSearchJobStatusCode(t *testing.T, user string, password string, baseURL string, jobID string) bool {
	checkEventURL := baseURL + "/services/search/v2/jobs/" + jobID + "?output_mode=json"
	t.Logf("URL: " + checkEventURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	checkReqEvents, err := http.NewRequest("GET", checkEventURL, nil)
	require.NoError(t, err, "Error while preparing GET request to check search job status")
	checkReqEvents.SetBasicAuth(user, password)
	checkResp, err := client.Do(checkReqEvents)
	require.NoError(t, err, "Error while executing Http GET request to check search job status")
	defer checkResp.Body.Close()
	t.Logf("Send Request: Check query status code: " + strconv.Itoa(checkResp.StatusCode))
	checkBody, err := io.ReadAll(checkResp.Body)
	require.NoError(t, err, "Error while reading response body")
	var checkJSONResponse map[string]any
	err = json.Unmarshal(checkBody, &checkJSONResponse)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")
	// Print isDone field from response
	isDone := checkJSONResponse["entry"].([]any)[0].(map[string]any)["content"].(map[string]any)["isDone"].(bool)
	t.Logf("Is Splunk Search completed [isDone flag]: %v\n", isDone)
	return isDone
}
func postSearchRequest(t *testing.T, user string, password string, baseURL string, searchQuery string, startTime string, endTime string) string {
	searchURL := fmt.Sprintf("%s/services/search/v2/jobs?output_mode=json", baseURL)
	query := searchQuery
	t.Logf("Search query: " + query)
	data := url.Values{}
	data.Set("search", query)
	data.Set("earliest_time", startTime)
	data.Set("latest_time", endTime)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", searchURL, strings.NewReader(data.Encode()))
	require.NoError(t, err, "Error while preparing POST request")
	req.SetBasicAuth(user, password)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	require.NoError(t, err, "Error while executing Http POST request")
	defer resp.Body.Close()
	t.Logf("Send Request: Post query status code: " + strconv.Itoa(resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Error while reading response body")
	var jsonResponse map[string]any
	err = json.Unmarshal(body, &jsonResponse)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")
	t.Logf("%+v", jsonResponse) // debug
	return jsonResponse["sid"].(string)
}

func waitForSearchCompletion(t *testing.T, user string, password string, baseURL string, jobID string) error {
	for i := 0; i < MAX_SEARCH_RETRIES; i++ { // limit loop - not allowing infinite looping
		t.Logf("Checking Search Status ...")
		isDone := checkSearchJobStatusCode(t, user, password, baseURL, jobID)
		if isDone {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("search job %s did not complete within the expected time", jobID)
}
