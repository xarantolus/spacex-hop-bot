package review

import (
	"fmt"
)

func Diff(oldResponse, newResponse *DashboardResponse) (changeDescriptions []string) {
	if oldResponse.Nid != newResponse.Nid || oldResponse.Title != newResponse.Title {
		// Cannot diff different projects
		return nil
	}

	// Has the project status changed?
	if oldResponse.ProjectStatus != newResponse.ProjectStatus && newResponse.ProjectStatus != "" {
		changeDescriptions = append(changeDescriptions, fmt.Sprintf("The project status for the review of %q has changed from %q to %q", newResponse.Title, oldResponse.ProjectStatus, newResponse.ProjectStatus))
	}

	// Check for end date changes
	if oldResponse.TotalDuration.EndDate != newResponse.TotalDuration.EndDate && newResponse.TotalDuration.EndDate != "" {
		changeDescriptions = append(changeDescriptions, fmt.Sprintf("The estimated completion date of the environmental review has changed from %s to %s", oldResponse.TotalDuration.EndDate, newResponse.TotalDuration.EndDate))
	}

	var oldProjectsByID = map[string]Data{}
	for i := range oldResponse.Data {
		v := oldResponse.Data[i]
		oldProjectsByID[v.Nid] = v
	}

	for _, newProject := range newResponse.Data {
		oldProject, ok := oldProjectsByID[newProject.Nid]
		if !ok {
			if newProject.Action != "" && newProject.TotalDuration.EndDate != "" {
				newDesc := fmt.Sprintf("A new project %q with end date %s has been added", newProject.Action, newProject.TotalDuration.EndDate)

				changeDescriptions = append(changeDescriptions, newDesc)
			}

			continue
		}

		// Also check for status changes
		if oldProject.Status != newProject.Status && newProject.Status != "" {
			changeDescriptions = append(changeDescriptions, fmt.Sprintf("The status of the %q has changed from %q to %q", newProject.Action, oldProject.Status, newProject.Status))
		}

		var milestoneDataByID = map[string]MilestoneData{}
		for i := range oldProject.MilestoneData {
			v := oldProject.MilestoneData[i]
			milestoneDataByID[v.Tid] = v
		}

		for _, newMilestone := range newProject.MilestoneData {
			oldMilestone, ok := milestoneDataByID[newMilestone.Tid]
			if !ok {
				if newMilestone.Name != "" && newMilestone.CurrentTargetDate != "" {
					newDesc := fmt.Sprintf("A new milestone %q with target date %q has been added to the %q", newMilestone.Name, newMilestone.CurrentTargetDate, newProject.Action)

					changeDescriptions = append(changeDescriptions, newDesc)
				}
				continue
			}

			if !oldMilestone.MilestoneComplete && newMilestone.MilestoneComplete {
				changeDescriptions = append(changeDescriptions, fmt.Sprintf("The milestone %q of the %q has been completed", newMilestone.Name, newProject.Action))
			}

			if oldMilestone.CurrentTargetDate != newMilestone.CurrentTargetDate && newMilestone.CurrentTargetDate != "" {
				changeDescriptions = append(changeDescriptions, fmt.Sprintf("The target date of milestone %q of the %q has changed from %q to %q", newMilestone.Name, newProject.Action, oldMilestone.CurrentTargetDate, newMilestone.CurrentTargetDate))
			}
		}
	}

	return
}
