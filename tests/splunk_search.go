package tests

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	MAX_SEARCH_RETRIES = 20 // maximum number of retries for checking search status
)

func checkEventsFromSplunk(t *testing.T, searchQuery string, startTime string, endTimeOptional ...string) []any {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("-->> Splunk Search: checking events in Splunk --")
	user := GetConfigVariable("USER")
	password := GetConfigVariable("PASSWORD")
	baseURL := "https://" + GetConfigVariable("HOST") + ":" + GetConfigVariable("MANAGEMENT_PORT")
	endTime := "now"
	if len(endTimeOptional) > 0 {
		endTime = endTimeOptional[0]
	}
	// post search
	jobID := postSearchRequest(t, user, password, baseURL, searchQuery, startTime, endTime)
	// wait for search status done == true
	for i := 0; i < MAX_SEARCH_RETRIES; i++ { // limit loop - not allowing infinite looping
		logger.Println("Checking Search Status ...")
		isDone := checkSearchJobStatusCode(t, user, password, baseURL, jobID)
		if isDone {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// get events
	reqUrl := fmt.Sprintf("%s/services/search/v2/jobs/%s/events?output_mode=json", baseURL, jobID)
	results := getSplunkSearchResults[any](t, user, password, reqUrl)
	return results
}

func checkStatisticsFromSplunk[V any](t *testing.T, searchQuery string, startTime string, endTimeOptional ...string) []V {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("-->> Splunk Search: checking events in Splunk --")
	user := GetConfigVariable("USER")
	password := GetConfigVariable("PASSWORD")
	baseURL := "https://" + GetConfigVariable("HOST") + ":" + GetConfigVariable("MANAGEMENT_PORT")
	endTime := "now"
	if len(endTimeOptional) > 0 {
		endTime = endTimeOptional[0]
	}
	// post search
	jobID := postSearchRequest(t, user, password, baseURL, searchQuery, startTime, endTime)
	// wait for search status done == true
	for i := 0; i < MAX_SEARCH_RETRIES; i++ { // limit loop - not allowing infinite looping
		logger.Println("Checking Search Status ...")
		isDone := checkSearchJobStatusCode(t, user, password, baseURL, jobID)
		if isDone {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// get events
	reqUrl := fmt.Sprintf("%s/services/search/v2/jobs/%s/results?output_mode=json", baseURL, jobID)
	results := getSplunkSearchResults[V](t, user, password, reqUrl)
	return results
}

func getSplunkSearchResults[V any](t *testing.T, user string, password string, reqUrl string) []V {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("URL: " + reqUrl)
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
	logger.Println("Send Request: Get query status code: " + strconv.Itoa(respEvents.StatusCode))

	bodyEvents, err := io.ReadAll(respEvents.Body)
	require.NoError(t, err, "Error while reading response body")

	var jsonResponseEvents map[string]any
	err = json.Unmarshal(bodyEvents, &jsonResponseEvents)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")

	rawResults := jsonResponseEvents["results"].([]any)

	var results []V
	for _, raw := range rawResults {
		b, err := json.Marshal(raw)
		require.NoError(t, err, "Error while marshalling result item")
		var v V
		err = json.Unmarshal(b, &v)
		require.NoError(t, err, "Error while unmarshalling result item to target type")
		results = append(results, v)
	}

	return results
}

func checkSearchJobStatusCode(t *testing.T, user string, password string, baseURL string, jobID string) bool {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	checkEventURL := baseURL + "/services/search/v2/jobs/" + jobID + "?output_mode=json"
	logger.Println("URL: " + checkEventURL)
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
	logger.Println("Send Request: Check query status code: " + strconv.Itoa(checkResp.StatusCode))
	checkBody, err := io.ReadAll(checkResp.Body)
	require.NoError(t, err, "Error while reading response body")
	var checkJSONResponse map[string]any
	err = json.Unmarshal(checkBody, &checkJSONResponse)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")
	// Print isDone field from response
	isDone := checkJSONResponse["entry"].([]any)[0].(map[string]any)["content"].(map[string]any)["isDone"].(bool)
	logger.Printf("Is Splunk Search completed [isDone flag]: %v\n", isDone)
	return isDone
}
func postSearchRequest(t *testing.T, user string, password string, baseURL string, searchQuery string, startTime string, endTime string) string {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	searchURL := fmt.Sprintf("%s/services/search/v2/jobs?output_mode=json", baseURL)
	query := searchQuery
	logger.Println("Search query: " + query)
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
	logger.Println("Send Request: Post query status code: " + strconv.Itoa(resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Error while reading response body")
	var jsonResponse map[string]any
	err = json.Unmarshal(body, &jsonResponse)
	require.NoError(t, err, "Error while unmarshalling response body to JSON")
	logger.Println(jsonResponse) // debug
	return jsonResponse["sid"].(string)
}
