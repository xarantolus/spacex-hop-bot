package review

type DashboardResponse struct {
	Code                       int            `json:"code"`
	Status                     string         `json:"status"`
	Data                       []Data         `json:"data"`
	Nid                        string         `json:"nid"`
	Title                      string         `json:"title"`
	IsProjectCompleted         bool           `json:"isProjectCompleted"`
	ProjectStatus              string         `json:"projectStatus"`
	StartDate                  int64          `json:"startDate"`
	EndDate                    int64          `json:"endDate"`
	ActualEndDate              int64          `json:"actualEndDate"`
	RangeStart                 int64          `json:"rangeStart"`
	RangeEnd                   int64          `json:"rangeEnd"`
	IsFast41                   bool           `json:"isFast41"`
	HasPausedAction            bool           `json:"hasPausedAction"`
	NumOfTicks                 int            `json:"numOfTicks"`
	Fast41InitiationDate       interface{}    `json:"fast41InitiationDate"`
	Fast41InitiationDateString interface{}    `json:"fast41InitiationDateString"`
	Pauses                     []interface{}  `json:"pauses"`
	PlannedStart               string         `json:"plannedStart"`
	PlannedEnd                 string         `json:"plannedEnd"`
	InProgressData             InProgressData `json:"inProgressData"`
	TotalDuration              TotalDuration  `json:"totalDuration"`
}

type MilestoneData struct {
	Tid                       string      `json:"tid"`
	Name                      string      `json:"name"`
	OriginalTargetDate        string      `json:"originalTargetDate"`
	CurrentTargetDate         string      `json:"currentTargetDate"`
	MilestoneComplete         bool        `json:"milestoneComplete"`
	NumberOfTargetDateChanges int         `json:"numberOfTargetDateChanges"`
	Fast41MissedDates         interface{} `json:"fast41MissedDates"`
}

type InProgressData struct {
	DaysInPause  int    `json:"daysInPause"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	DurationDays int    `json:"durationDays"`
}

type TotalDuration struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	DurationDays int    `json:"durationDays"`
}

type Data struct {
	Nid                  string          `json:"nid"`
	D7Nid                interface{}     `json:"d7_nid"`
	Action               string          `json:"action"`
	Path                 string          `json:"path"`
	Agency               string          `json:"agency"`
	Bureau               string          `json:"bureau"`
	Shpo                 interface{}     `json:"shpo"`
	ActionIdentifier     interface{}     `json:"actionIdentifier"`
	ActionOutcome        interface{}     `json:"actionOutcome"`
	ResponsibleAgency    string          `json:"responsibleAgency"`
	StatusTid            string          `json:"status_tid"`
	Status               string          `json:"status"`
	ActualCompletionDate int             `json:"actualCompletionDate"`
	MilestoneData        []MilestoneData `json:"milestoneData"`
	HasMissedDate        interface{}     `json:"hasMissedDate"`
	MissedDateExempted   interface{}     `json:"missedDateExempted"`
	StartDate            int64           `json:"startDate"`
	EndDate              int64           `json:"endDate"`
	NoEndDate            bool            `json:"noEndDate"`
	ExpectedEndDate      int64           `json:"expectedEndDate"`
	ActionComplete       bool            `json:"actionComplete"`
	Pauses               []interface{}   `json:"pauses"`
	PlannedStart         string          `json:"plannedStart"`
	PlannedEnd           string          `json:"plannedEnd"`
	InProgressData       InProgressData  `json:"inProgressData"`
	TotalDuration        TotalDuration   `json:"totalDuration"`
}
