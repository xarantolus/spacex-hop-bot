package review

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/xarantolus/spacex-hop-bot/util"
)

type FAAClient struct {
	client *http.Client

	lastReportsLock sync.Mutex
	lastReports     map[int]*DashboardResponse

	fn string
}

func NewFAAClient() *FAAClient {
	var client = &FAAClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},

		fn: "faa-dashboard.json",
	}

	client.lastReports = make(map[int]*DashboardResponse)

	_ = util.LoadJSON(client.fn, &client.lastReports)

	return client
}

func (f *FAAClient) save() error {
	return util.SaveJSON(f.fn, &f.lastReports)
}

func (f *FAAClient) projectReport(projectId int) (r *DashboardResponse, err error) {
	var url = fmt.Sprintf("https://www.permits.performance.gov/api/v1/project?nid=%d", projectId)

	resp, err := f.client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("unexpected status code %d while requesting project with id %d", resp.StatusCode, projectId)
		return
	}

	r = new(DashboardResponse)

	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return
	}

	if r.Code != 200 {
		err = fmt.Errorf("unexpected dashboard status code %d while requesting project with id %d", r.Code, projectId)

		return
	}

	return r, nil
}

func (f *FAAClient) ReportProjectDiff(projectId int) (diff []string, err error) {
	f.lastReportsLock.Lock()
	last, haveLast := f.lastReports[projectId]
	f.lastReportsLock.Unlock()

	project, err := f.projectReport(projectId)
	if err != nil {
		return
	}

	// Save the last known project
	f.lastReportsLock.Lock()
	f.lastReports[projectId] = project
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
