package review

import (
	"encoding/json"
	"reflect"
	"testing"
)

var exampleDashboardResponse = DashboardResponse{
	Code:   200,
	Status: "OK",

	Data: []Data{
		{
			Nid:    "96591",
			D7Nid:  interface{}(nil),
			Action: "Endangered Species Act Consultation (DOI-FWS)",
			Path:   "/proj/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site-cameron-1",
			Agency: "Department of the Interior",
			Bureau: "Fish and Wildlife Service",
			Shpo:   interface{}(nil), ActionIdentifier: interface{}(nil), ActionOutcome: interface{}(nil), ResponsibleAgency: "Fish and Wildlife Service",
			StatusTid:            "3036",
			Status:               "In Progress",
			ActualCompletionDate: 0,
			MilestoneData: []MilestoneData{
				{Tid: "1579001",
					Name:               "Request for ESA Consultation Received",
					OriginalTargetDate: "2021-06-21",
					CurrentTargetDate:  "2021-06-21",
					MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1579006",
					Name:               "Consultation\u00a0Package Deemed Complete – Formal",
					OriginalTargetDate: "2021-10-06",
					CurrentTargetDate:  "2021-10-06",
					MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1579011",
					Name:                      "Conclusion of ESA Consultation",
					OriginalTargetDate:        "2021-12-31",
					CurrentTargetDate:         "2021-12-31",
					MilestoneComplete:         false,
					NumberOfTargetDateChanges: 0,
					Fast41MissedDates:         interface{}(nil),
				},
			},
			HasMissedDate: interface{}(nil), MissedDateExempted: interface{}(nil), StartDate: 1624248000000, EndDate: 1640926800000, NoEndDate: false, ExpectedEndDate: 1640926800000, ActionComplete: false, Pauses: []interface{}{}, PlannedStart: "06/21/2021",
			PlannedEnd: "10/06/2021",
			InProgressData: InProgressData{DaysInPause: 0, StartDate: "10/06/2021",
				EndDate:      "12/04/2021",
				DurationDays: 59}, TotalDuration: TotalDuration{StartDate: "06/21/2021",
				EndDate:      "12/31/2021",
				DurationDays: 193},
		}, {Nid: "96606",
			D7Nid: interface{}(nil), Action: "Section 106 Review",
			Path:             "/proj/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site-cameron-3",
			Agency:           "Department of Transportation",
			Bureau:           "Federal Aviation Administration",
			Shpo:             "Texas",
			ActionIdentifier: interface{}(nil), ActionOutcome: interface{}(nil), ResponsibleAgency: "Federal Aviation Administration",
			StatusTid:            "3036",
			Status:               "In Progress",
			ActualCompletionDate: 0, MilestoneData: []MilestoneData{{Tid: "1579121",
				Name:               "Consultation initiated with SHPO/THPO",
				OriginalTargetDate: "2021-06-25",
				CurrentTargetDate:  "2021-06-25",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1579131",
				Name:               "Section 106 consultation concluded",
				OriginalTargetDate: "2021-12-20",
				CurrentTargetDate:  "2021-12-20",
				MilestoneComplete:  false, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)},
			}, HasMissedDate: interface{}(nil), MissedDateExempted: interface{}(nil), StartDate: 1624593600000, EndDate: 1639976400000, NoEndDate: false, ExpectedEndDate: 1639976400000, ActionComplete: false, Pauses: []interface{}{}, PlannedStart: "06/25/2021",
			PlannedEnd: "06/25/2021",
			InProgressData: InProgressData{DaysInPause: 0, StartDate: "06/25/2021",
				EndDate:      "12/04/2021",
				DurationDays: 162}, TotalDuration: TotalDuration{StartDate: "06/25/2021",
				EndDate:      "12/20/2021",
				DurationDays: 178},
		}, {Nid: "96581",
			D7Nid: interface{}(nil), Action: "Environmental Assessment (EA)",
			Path:   "/proj/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site-cameron-county",
			Agency: "Department of Transportation",
			Bureau: "Federal Aviation Administration",
			Shpo:   interface{}(nil), ActionIdentifier: interface{}(nil), ActionOutcome: interface{}(nil), ResponsibleAgency: "Federal Aviation Administration",
			StatusTid:            "3036",
			Status:               "In Progress",
			ActualCompletionDate: 0, MilestoneData: []MilestoneData{{Tid: "1578906",
				Name:               "Determination to prepare an Environmental Assessment (EA)",
				OriginalTargetDate: "2021-07-01",
				CurrentTargetDate:  "2021-07-01",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1578911",
				Name:               "Issuance of a Draft EA / Release for Public Review",
				OriginalTargetDate: "2021-09-17",
				CurrentTargetDate:  "2021-09-17",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1578916",
				Name:               "Issuance of a Final EA",
				OriginalTargetDate: "2021-12-31",
				CurrentTargetDate:  "2021-12-31",
				MilestoneComplete:  false, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1578926",
				Name:               "EA Process Concluded",
				OriginalTargetDate: "2021-12-31",
				CurrentTargetDate:  "2021-12-31",
				MilestoneComplete:  false, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)},
			}, HasMissedDate: interface{}(nil), MissedDateExempted: interface{}(nil), StartDate: 1625112000000, EndDate: 1640926800000, NoEndDate: false, ExpectedEndDate: 1640926800000, ActionComplete: false, Pauses: []interface{}{}, PlannedStart: "07/01/2021",
			PlannedEnd: "07/01/2021",
			InProgressData: InProgressData{DaysInPause: 0, StartDate: "07/01/2021",
				EndDate:      "12/04/2021",
				DurationDays: 156}, TotalDuration: TotalDuration{StartDate: "07/01/2021",
				EndDate:      "12/31/2021",
				DurationDays: 183},
		}, {Nid: "96601",
			D7Nid: interface{}(nil), Action: "Endangered Species Act Consultation (NOAA-NMFS)",
			Path:   "/proj/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site-cameron-2",
			Agency: "Department of Commerce",
			Bureau: "National Oceanic and Atmospheric Administration",
			Shpo:   interface{}(nil), ActionIdentifier: interface{}(nil), ActionOutcome: interface{}(nil), ResponsibleAgency: "National Oceanic and Atmospheric Administration",
			StatusTid:            "3036",
			Status:               "In Progress",
			ActualCompletionDate: 0, MilestoneData: []MilestoneData{{Tid: "1579081",
				Name:               "Request for ESA Consultation Received",
				OriginalTargetDate: "2021-08-11",
				CurrentTargetDate:  "2021-08-11",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1579096",
				Name:               "Consultation Package Deemed Complete – Informal",
				OriginalTargetDate: "2021-11-04",
				CurrentTargetDate:  "2021-11-04",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1579091",
				Name:               "Conclusion of ESA Consultation",
				OriginalTargetDate: "2021-12-31",
				CurrentTargetDate:  "2021-12-31",
				MilestoneComplete:  false, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)},
			}, HasMissedDate: interface{}(nil), MissedDateExempted: interface{}(nil), StartDate: 1628654400000, EndDate: 1640926800000, NoEndDate: false, ExpectedEndDate: 1640926800000, ActionComplete: false, Pauses: []interface{}{}, PlannedStart: "08/11/2021",
			PlannedEnd: "11/04/2021",
			InProgressData: InProgressData{DaysInPause: 0, StartDate: "11/04/2021",
				EndDate:      "12/04/2021",
				DurationDays: 30}, TotalDuration: TotalDuration{StartDate: "08/11/2021",
				EndDate:      "12/31/2021",
				DurationDays: 142},
		}, {Nid: "96586",
			D7Nid: interface{}(nil), Action: "Magnuson-Stevens Fishery Conservation and Management Act, Section 305 Essential Fish Habitat (EFH) Consultation",
			Path:   "/proj/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site-cameron-0",
			Agency: "Department of Commerce",
			Bureau: "National Oceanic and Atmospheric Administration",
			Shpo:   interface{}(nil), ActionIdentifier: interface{}(nil), ActionOutcome: interface{}(nil), ResponsibleAgency: "National Oceanic and Atmospheric Administration",
			StatusTid:            "3026",
			Status:               "Complete",
			ActualCompletionDate: 0, MilestoneData: []MilestoneData{{Tid: "1578951",
				Name:               "Lead Agency Requests EFH Consultation by submitting an EFH Assessment",
				OriginalTargetDate: "2021-09-17",
				CurrentTargetDate:  "2021-09-17",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1578956",
				Name:               "NOAA Determines the EFH Assessment is complete and Initiates consultation",
				OriginalTargetDate: "2021-10-18",
				CurrentTargetDate:  "2021-10-18",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)}, {Tid: "1578961",
				Name:               "NOAA Issues any EFH conservation recommendations",
				OriginalTargetDate: "2021-10-29",
				CurrentTargetDate:  "2021-10-29",
				MilestoneComplete:  true, NumberOfTargetDateChanges: 0, Fast41MissedDates: interface{}(nil)},
			}, HasMissedDate: interface{}(nil), MissedDateExempted: interface{}(nil), StartDate: 1631851200000, EndDate: 1635480000000, NoEndDate: false, ExpectedEndDate: 1635480000000, ActionComplete: false, Pauses: []interface{}{}, PlannedStart: "09/17/2021",
			PlannedEnd: "10/18/2021",
			InProgressData: InProgressData{DaysInPause: 0, StartDate: "10/18/2021",
				EndDate:      "10/29/2021",
				DurationDays: 11}, TotalDuration: TotalDuration{StartDate: "09/17/2021",
				EndDate:      "10/29/2021",
				DurationDays: 42},
		},
	}, Nid: "96576",
	Title:              "SpaceX Starship/Super Heavy Launch Vehicle Program at the SpaceX Boca Chica Launch Site in Cameron County, Texas",
	IsProjectCompleted: false, ProjectStatus: "In Progress",
	StartDate: 1624248000000, EndDate: 1640926800000, ActualEndDate: 1640926800000, RangeStart: 1611288000000, RangeEnd: 1653886800000, IsFast41: false, HasPausedAction: false, NumOfTicks: 15, Fast41InitiationDate: interface{}(nil), Fast41InitiationDateString: interface{}(nil), Pauses: []interface{}{}, PlannedStart: "06/21/2021",
	PlannedEnd: "06/25/2021",
	InProgressData: InProgressData{DaysInPause: 0, StartDate: "06/25/2021",
		EndDate:      "12/04/2021",
		DurationDays: 162}, TotalDuration: TotalDuration{StartDate: "06/21/2021",
		EndDate:      "12/31/2021",
		DurationDays: 193},
}

func copyDefaultWith(f func(*DashboardResponse)) *DashboardResponse {
	// Using this to deep-copy the struct, so modifications are not visible to it
	origJSON, err := json.Marshal(exampleDashboardResponse)
	if err != nil {
		panic(err)
	}

	clone := new(DashboardResponse)
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		panic(err)
	}

	f(clone)

	return clone
}

func TestDiff(t *testing.T) {
	var old = exampleDashboardResponse

	tests := []struct {
		arg                    *DashboardResponse
		wantChangeDescriptions []string
	}{
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.TotalDuration.EndDate = "NewEndDate"
			}),
			wantChangeDescriptions: []string{
				"The estimated completion date of the environmental review has changed from 12/31/2021 to NewEndDate",
			},
		},
		{
			// A change to nothing is ignored
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.TotalDuration.EndDate = ""
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.Nid = "something else"
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.Title = "something else"
			}),
			wantChangeDescriptions: nil,
		},

		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.ProjectStatus = "Completed"
			}),
			wantChangeDescriptions: []string{
				`The project status for the review of "SpaceX Starship/Super Heavy Launch Vehicle Program at the SpaceX Boca Chica Launch Site in Cameron County, Texas" has changed from "In Progress" to "Completed"`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.ProjectStatus = ""
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.TotalDuration.EndDate = ""
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op2.Status = "Completed"
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: []string{
				`The status of the "Endangered Species Act Consultation (DOI-FWS)" has changed from "In Progress" to "Completed"`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op2.Status = ""
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[3]
				op2.Status = "Completed"
				op[3] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: []string{
				`The status of the "Endangered Species Act Consultation (NOAA-NMFS)" has changed from "In Progress" to "Completed"`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[3]
				op2.Status = ""
				op[3] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: nil,
		},

		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op3 := op2.MilestoneData
				op4 := op3[2]
				op4.MilestoneComplete = true
				op3[2] = op4
				op2.MilestoneData = op3
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: []string{
				`The milestone "Conclusion of ESA Consultation" of the "Endangered Species Act Consultation (DOI-FWS)" has been completed`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op3 := op2.MilestoneData
				op4 := op3[2]
				op4.CurrentTargetDate = "2021-11-06"
				op3[2] = op4
				op2.MilestoneData = op3
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: []string{
				`The target date of milestone "Conclusion of ESA Consultation" of the "Endangered Species Act Consultation (DOI-FWS)" has changed from "2021-12-31" to "2021-11-06"`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op3 := op2.MilestoneData
				op4 := op3[2]
				op4.CurrentTargetDate = ""
				op3[2] = op4
				op2.MilestoneData = op3
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: nil,
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				dr.Data = append(dr.Data, Data{
					Action: "Some new project",
					TotalDuration: TotalDuration{
						EndDate: "2022-06-01",
					},
				})
			}),
			wantChangeDescriptions: []string{
				`A new project "Some new project" with end date 2022-06-01 has been added`,
			},
		},
		{
			arg: copyDefaultWith(func(dr *DashboardResponse) {
				op := dr.Data
				op2 := op[0]
				op3 := op2.MilestoneData

				op3 = append(op3, MilestoneData{
					Tid:               "someid",
					Name:              "New Milestone",
					CurrentTargetDate: "2022-03-04",
				})

				op2.MilestoneData = op3
				op[0] = op2
				dr.Data = op
			}),
			wantChangeDescriptions: []string{
				`A new milestone "New Milestone" with target date "2022-03-04" has been added to the "Endangered Species Act Consultation (DOI-FWS)"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if gotChangeDescriptions := Diff(&old, tt.arg); !reflect.DeepEqual(gotChangeDescriptions, tt.wantChangeDescriptions) {
				t.Errorf("Diff() = %v, want %v",
					gotChangeDescriptions, tt.wantChangeDescriptions)
			}
		})
	}
}
