package cronjobs

import (
	"bytes"
	"html/template"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"bitbucket.org/radarventures/forwarder-adapters/apis/id"
	"bitbucket.org/radarventures/forwarder-adapters/apis/misc"
	"bitbucket.org/radarventures/forwarder-adapters/apis/notifications"
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/misc"
	sdtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	cardSrv "bitbucket.org/radarventures/forwarder-shipments/services/card"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
)

type ExpiryCards struct {
	id      id.ID
	cardSrv cardSrv.ICardService
	misc    misc.Misc
	not     notifications.Notifications
}

func NewExpiryCards() IExpiryCards {
	return &ExpiryCards{
		id:      *id.New(config.Get().IdURL),
		cardSrv: cardSrv.NewCardService(),
		misc:    *misc.New(config.Get().MiscURL),
		not:     *notifications.New(config.Get().MiscURL),
	}
}

type IExpiryCards interface {
	SendExpiryCardsNotifications(ctx *context.Context)
	SendExpiryCardsNotificationsManager(ctx *context.Context) error
	SendExpiryCardsNotificationsExecutive(ctx *context.Context) error
}

func (t *ExpiryCards) SendExpiryCardsNotifications(ctx *context.Context) {
	t.SendExpiryCardsNotificationsManager(ctx)
	t.SendExpiryCardsNotificationsExecutive(ctx)
}

func (t *ExpiryCards) SendExpiryCardsNotificationsManager(ctx *context.Context) error {

	// Parsing the template for manager notifications
	template, err := template.New("ExpiryCardsMailManager_Shared").Parse(config.ExpiryCards_ManagerNotification)
	if err != nil {
		ctx.Log.Error("unable to parse expiry cards manager notification template", zap.Error(err))
		return err
	}

	// Retrieve manager IDs
	managerIds, err := t.id.GetManagerIds(ctx)
	if err != nil {
		ctx.Log.Error("unable to get manager ids", zap.Error(err))
		return err
	}

	// Load configuration values
	blockedUsers := config.Get().ExpiredSummaryMailBlockedUsers
	allowedUsers := config.Get().ExpiredSummaryMailAllowedUsers
	summaryAllowedUsers := config.Get().SummaryMailAllowedUsers
	blockedUsersManagers := config.Get().Blockusers

	// Filter out blocked managers
	filteredManagers := make([]string, 0)
	for _, managerID := range managerIds {
		if !utils.ContainsString(blockedUsersManagers, managerID) {
			filteredManagers = append(filteredManagers, managerID)
		}
	}

	// Iterate over filtered managers
	for _, managerId := range filteredManagers {
		orgTreeWithCount, err := t.cardSrv.GetOrgTreeCount(ctx, managerId, false, "", "")
		if err != nil {
			ctx.Log.Error("unable to get org tree count for manager", zap.Error(err))
			return err
		}

		if orgTreeWithCount == nil {
			ctx.Log.Warn("org tree count empty", zap.Any("manager_id", managerId))
			continue
		}

		emailTemplate := &dtos.Notification{
			ID:              uuid.New().String(),
			Type:            constants.NotTypeEmail,
			Title:           "IMPORTANT! Cards Pending on Your Reportees Workspace",
			Sender:          config.Get().EmailSenderBot,
			IsTransactional: true,
		}

		for adminIdx := range orgTreeWithCount.List {

			managerDetail := &sdtos.ManagerMailDetail{
				Name:          orgTreeWithCount.List[adminIdx].Name,
				Email:         orgTreeWithCount.List[adminIdx].Email,
				TeamViewLink:  config.Get().BaseURL + "/dashboard/manager-workspace?show=managerView",
				ReporteesData: make([]sdtos.ReporteeMailDetail, 0),
			}

			reporteeEmails := make([]string, 0)

			managerID := orgTreeWithCount.List[adminIdx].ID
			blockedUser := utils.ContainsString(blockedUsers, managerID)
			allowedUser := utils.ContainsString(allowedUsers, managerID)

			if (!blockedUser && allowedUser) || (!blockedUser && len(summaryAllowedUsers) == 0) {

				// Iterate over reportees
				for reporteeIdx := range orgTreeWithCount.List[adminIdx].Reportees {
					reportee := &orgTreeWithCount.List[adminIdx].Reportees[reporteeIdx]

					if reportee.BreachedTasks > 0 || reportee.ExpiringTasks > 0 {
						reporteeDetail := sdtos.ReporteeMailDetail{
							Name:          reportee.Name,
							ElapsedCards:  int64(reportee.BreachedTasks),
							ElapsingCards: int64(reportee.ExpiringTasks),
							Email:         reportee.Email,
						}
						reporteeEmails = append(reporteeEmails, reporteeDetail.Email)
						managerDetail.ReporteesData = append(managerDetail.ReporteesData, reporteeDetail)
					}
				}

				if len(managerDetail.Email) > 0 && len(managerDetail.ReporteesData) > 0 {
					buf := new(bytes.Buffer)
					if err := template.Execute(buf, managerDetail); err != nil {
						ctx.Log.Error("unable to execute expiry cards manager notification template", zap.Error(err))
						return err
					}

					emailTemplate.Content = buf.String()
					emailTemplate.Receivers = []string{managerDetail.Email}
					emailTemplate.CC = reporteeEmails

					t.not.SendNotification(ctx, emailTemplate)
				}
			}
		}
	}

	return nil
}

func (t *ExpiryCards) SendExpiryCardsNotificationsExecutive(ctx *context.Context) error {

	// Parsing the template for executive notifications
	template, err := template.New("ExpiryCardsMailExecutive_Shared").Parse(config.ExpiryCards_ReporteesNotification)
	if err != nil {
		ctx.Log.Error("unable to parse expiry cards share notification template for executives", zap.Error(err))
		return err
	}

	// Retrieve manager IDs
	managerIds, err := t.id.GetManagerIds(ctx)
	if err != nil {
		ctx.Log.Error("Unable to get the manager IDs", zap.Error(err))
		return err
	}

	// Load configuration values
	blockedUsers := config.Get().ExpiredSummaryMailBlockedUsers
	allowedUsers := config.Get().ExpiredSummaryMailAllowedUsersforexec
	summaryAllowedUsers := config.Get().SummaryMailAllowedUsersforexec
	blockedUsersManagers := config.Get().Blockusers

	// Filter out blocked managers
	filteredExecutives := make([]string, 0)
	for _, executiveID := range managerIds {
		if !utils.ContainsString(blockedUsersManagers, executiveID) {
			filteredExecutives = append(filteredExecutives, executiveID)
		}
	}

	// Iterate over filtered executives
	for _, executiveID := range filteredExecutives {
		// Retrieve organization tree with count for the executive
		orgTreeWithCount, err := t.cardSrv.GetOrgTreeCount(ctx, executiveID, false, "", "")
		if err != nil {
			ctx.Log.Error("unable to retrieve organization tree data for executive", zap.Error(err))
			return err
		}

		emailTemplate := &dtos.Notification{
			ID:              uuid.New().String(),
			Type:            constants.NotTypeEmail,
			Title:           "IMPORTANT! Cards Pending on Your Workspace",
			Sender:          config.Get().EmailSenderBot,
			IsTransactional: true,
		}

		// Iterate over executives and their reportees
		for adminIdx := range orgTreeWithCount.List {
			for reporteeIdx := range orgTreeWithCount.List[adminIdx].Reportees {
				reportee := &orgTreeWithCount.List[adminIdx].Reportees[reporteeIdx]
				reporteeID := reportee.ID

				reporteeDetail := sdtos.ReporteeMailDetail{
					Name:          reportee.Name,
					ElapsedCards:  int64(reportee.BreachedTasks),
					ElapsingCards: int64(reportee.ExpiringTasks),
					Email:         reportee.Email,
					WorkspaceLink: config.Get().BaseURL + "/dashboard/executive/" + reporteeID,
				}

				blockedUser := utils.ContainsString(blockedUsers, reporteeID)
				allowedUser := utils.ContainsString(allowedUsers, reporteeID)

				// Check if conditions are met for sending notifications
				if !blockedUser && ((allowedUser || len(summaryAllowedUsers) == 0) &&
					(reporteeDetail.ElapsedCards > 0 || reporteeDetail.ElapsingCards > 0) &&
					len(reporteeDetail.Email) > 0) {

					// Execute the template to create the notification content
					reporteesBuf := new(bytes.Buffer)
					err := template.Execute(reporteesBuf, reporteeDetail)
					if err != nil {
						ctx.Log.Error("unable to execute expiry cards notification template for reportees", zap.Error(err))
						return err
					}

					// Update email template content and receivers
					emailTemplate.Content = reporteesBuf.String()
					emailTemplate.Receivers = []string{reporteeDetail.Email}

					// Send the email notification
					t.not.SendNotification(ctx, emailTemplate)
				}
			}
		}
	}

	return nil
}
