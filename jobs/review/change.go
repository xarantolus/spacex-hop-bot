package review

import (
	"fmt"
	"log"
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

	if len(oldResponse.Data) != len(newResponse.Data) {
		log.Println("[Review] It seems like new data has been added")
		return
	}

	for i := 0; i < len(oldResponse.Data); i++ {
		var oldProject, newProject = oldResponse.Data[i], newResponse.Data[i]

		if oldProject.Nid != newProject.Nid {
			log.Println("[Review] It seems like old/new data doesn't have same order")
			continue
		}

		// Also check for status changes
		if oldProject.Status != newProject.Status && newProject.Status != "" {
			changeDescriptions = append(changeDescriptions, fmt.Sprintf("The status of the %q has changed from %q to %q", newProject.Action, oldProject.Status, newProject.Status))
		}

		if len(oldProject.MilestoneData) != len(newProject.MilestoneData) {
			log.Printf("[Review] It seems like new milestone data has been added to %q", newProject.Action)
			continue
		}

		for i := 0; i < len(oldProject.MilestoneData); i++ {
			var oldMilestone, newMilestone = oldProject.MilestoneData[i], newProject.MilestoneData[i]

			if oldMilestone.Name != newMilestone.Name {
				log.Printf("[Review] It seems like old/new milestone data for %q doesn't have same order", newProject.Action)
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
