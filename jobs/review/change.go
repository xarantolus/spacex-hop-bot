package review

import (
	"fmt"
	"log"
)

func Diff(old, new *DashboardResponse) (changeDescriptions []string) {
	if old.Nid != new.Nid || old.Title != new.Title {
		// Cannot diff different projects
		return nil
	}

	// Has the project status changed?
	if old.ProjectStatus != new.ProjectStatus && new.ProjectStatus != "" {
		changeDescriptions = append(changeDescriptions, fmt.Sprintf("The project status for the review of %q has changed from %q to %q", new.Title, old.ProjectStatus, new.ProjectStatus))
	}

	// Check for end date changes
	if old.TotalDuration.EndDate != new.TotalDuration.EndDate && new.TotalDuration.EndDate != "" {
		changeDescriptions = append(changeDescriptions, fmt.Sprintf("The estimated completion date of the environmental review has changed from %s to %s", old.TotalDuration.EndDate, new.TotalDuration.EndDate))
	}

	if len(old.Data) != len(new.Data) {
		log.Println("[FAA] It seems like new data has been added")
		return
	}

	for i := 0; i < len(old.Data); i++ {
		var oldProject, newProject = old.Data[i], new.Data[i]

		if oldProject.Nid != newProject.Nid {
			log.Println("[FAA] It seems like old/new data doesn't have same order")
			continue
		}

		// Also check for status changes
		if oldProject.Status != newProject.Status && newProject.Status != "" {
			changeDescriptions = append(changeDescriptions, fmt.Sprintf("The status of the %q has changed from %q to %q", newProject.Action, oldProject.Status, newProject.Status))
		}

		if len(oldProject.MilestoneData) != len(newProject.MilestoneData) {
			log.Printf("[FAA] It seems like new milestone data has been added to %q", newProject.Action)
			continue
		}

		for i := 0; i < len(oldProject.MilestoneData); i++ {
			var oldMilestone, newMilestone = oldProject.MilestoneData[i], newProject.MilestoneData[i]

			if oldMilestone.Name != newMilestone.Name {
				log.Printf("[FAA] It seems like old/new milestone data for %q doesn't have same order", newProject.Action)
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
