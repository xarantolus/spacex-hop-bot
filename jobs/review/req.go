package review

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/xarantolus/spacex-hop-bot/util"
)

type Client struct {
	httpClient *http.Client

	lastReportsLock sync.Mutex
	lastReports     map[int]*DashboardResponse

	fn string
}

func NewReviewClient() *Client {
	var client = &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},

		lastReports: make(map[int]*DashboardResponse),
		fn:          "faa-dashboard.json",
	}

	_ = util.LoadJSON(client.fn, &client.lastReports)

	return client
}

func (f *Client) save() error {
	return util.SaveJSON(f.fn, &f.lastReports)
}

func (f *Client) projectReport(projectID int) (r *DashboardResponse, err error) {
	var url = fmt.Sprintf("https://www.permits.performance.gov/api/v1/project?nid=%d", projectID)

	resp, err := f.httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("unexpected status code %d while requesting project with id %d", resp.StatusCode, projectID)
		return
	}

	r = new(DashboardResponse)

	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return
	}

	if r.Code != 200 {
		err = fmt.Errorf("unexpected dashboard status code %d while requesting project with id %d", r.Code, projectID)

		return
	}

	return r, nil
}

func (f *Client) ReportProjectDiff(projectID int) (diff []string, err error) {
	f.lastReportsLock.Lock()
	last, haveLast := f.lastReports[projectID]
	f.lastReportsLock.Unlock()

	project, err := f.projectReport(projectID)
	if err != nil {
		return
	}

	// Save the last known project
	f.lastReportsLock.Lock()
	f.lastReports[projectID] = project
	f.lastReportsLock.Unlock()

	err = f.save()
	if err != nil {
		return
	}
	if !haveLast {
		return nil, nil
	}

	// OK, we have a last project. Let's do a diff
	return Diff(last, project), nil
}
