package cronjobs

import (
	"fmt"
	"log"
	"strings"
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/services/ams"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphconfig "github.com/microsoftgraph/msgraph-sdk-go/users"
	"go.uber.org/zap"
)

func CheckMailsNew(ctx *context.Context) {
	amsService := ams.New()

	// Use the client secret credential to authenticate
	cred, err := azidentity.NewClientSecretCredential(config.Get().MSTenant, config.Get().MSClient, config.Get().MSSecret, nil)
	if err != nil {
		log.Fatalf("failed to create credential: %v", err)
	}

	scopes := []string{"https://graph.microsoft.com/.default"} // []string{"Mail.ReadBasic", "Mail.Read", "User.Read", "Mail.ReadWrite"}

	// Create a new Graph client
	graphClient, err := msgraph.NewGraphServiceClientWithCredentials(cred, scopes)

	if err != nil {
		fmt.Println(err)
		return
	}
	var filterTerm = "from/emailAddress/address eq 'noreply@tradetech.net' and isRead eq false"

	requestParameters := &graphconfig.ItemMailFoldersItemMessagesRequestBuilderGetQueryParameters{
		//Select: []string{"sender", "subject"},
		//Search: &searchTerm,

		Filter: &filterTerm,
		//Filter: &searchTerm,
	}
	configuration := &graphconfig.ItemMailFoldersItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}
	userUPN := config.Get().AMSEmail
	result, err := graphClient.Users().ByUserId(userUPN).MailFolders().ByMailFolderId("inbox").Messages().Get(ctx, configuration)

	if err != nil {
		fmt.Println("Unable to fetch", zap.Any("err", err))
		return
	}

	vals := result.GetValue()

	resp := map[string][]dtos.AMSError{}
	msg_ids := []string{}

	for _, v := range vals {

		errs := []dtos.AMSError{}
		//sender := v.GetSender().GetEmailAddress().GetAddress()

		subject := v.GetSubject()
		msg_id := v.GetId()

		if !strings.Contains(*subject, "Rejected File") {

			if strings.Contains(*subject, "Transmission Received: Accepted File:") {
				fileName := strings.Replace(*subject, "ftp_wizlogtec--Transmission Received: Accepted File: ", "", 1)

				if err := amsService.UpdateAmsStatus(ctx, fileName, globals.AMSAcceptedCode, nil); err == nil {
					UpdateReadStatus(ctx, graphClient, *msg_id)
				}
				continue
			}
			subjectSplits := strings.Split(*subject, ":")

			id := strings.Replace(subjectSplits[0], "WZLW HBL ", "", -1)

			//msgPrv := v.GetBodyPreview()

			//For success - 1Y HBL-MBL Linked  (Completed Stage)
			//For received stage - Received 3Z Security Filing
			// For Accepted stage - USA Ocean AMS Filing Accepted

			ctx.Log.Info("Id " + id + " Status : " + subjectSplits[len(subjectSplits)-1])

			if err := amsService.UpdateAmsStatus(ctx, id, subjectSplits[len(subjectSplits)-1], nil); err == nil {
				UpdateReadStatus(ctx, graphClient, *msg_id)
			}

		} else {
			var err dtos.AMSError
			body := v.GetBody().GetContent()
			obj := *body
			msg_ids = append(msg_ids, *v.GetId())
			var bl, NoOfErrs string

			for _, t := range strings.Split(obj, "\n") {

				va := strings.Trim(t, " ")

				if strings.Contains(va, "File Name: ") {

					err.File = replaceR(va, "File Name: ", "")
					//continue
				}

				if strings.Contains(va, "Number of Errors: ") {

					NoOfErrs = replaceR(va, "Number of Errors: ", "")
					//continue
				}

				if strings.Contains(va, "House Bill Number: ") {
					bl = replaceR(va, "House Bill Number: ", "")
				}
				if strings.Contains(va, "Document:") {
					err.Document = replaceR(va, "Document: ", "")
					//continue
				}
				if strings.Contains(va, "Error Value:") {
					err.ErrorValue = replaceR(va, "Error Value: ", "")
					//continue
				}
				if strings.Contains(va, "Error Message:") {
					err.ErrorMessage = replaceR(va, "Error Message: ", "")
					//continue
				}
				if strings.Contains(va, "Element Position:") {
					err.ElementPosition = replaceR(va, "Element Position: ", "")
					errs = append(errs, err)
					resp[err.File] = errs
					err = dtos.AMSError{Document: err.Document, File: err.File}
				}

				// l.Debug(va)

			}

			if err := amsService.UpdateAmsStatus(ctx, bl, fmt.Sprintf("Rejected With %s Errors", NoOfErrs), errs); err == nil {
				UpdateReadStatusNew(ctx, graphClient, *msg_id)
			}

		}

	}

	for k, v := range resp {

		if err := amsService.UpdateErrorStatus(ctx, k, v); err != nil {
			ctx.Log.Error("Unable to update error status ", zap.Error(err))
		}
	}

}

func UpdateReadStatusNew(ctx *context.Context, graphClient *msgraph.GraphServiceClient, message string) error {

	isRead := true
	requestBody := graphmodels.NewMessage()
	requestBody.SetIsRead(&isRead)
	userUPN := config.Get().AMSEmail
	m, err := graphClient.Users().ByUserId(userUPN).Messages().ByMessageId(message).Patch(ctx, requestBody, nil)
	if err != nil {
		fmt.Println(err)
	}

	sender := m.GetSender().GetEmailAddress().GetAddress()

	subject := m.GetSubject()

	msgPrv := m.GetBodyPreview()

	msg_id := m.GetId()

	fmt.Println("DATE : " + m.GetReceivedDateTime().Format(time.DateTime) + "Id " + *msg_id + " Sender : " + *sender + " Subject : " + *subject + " Preview : " + *msgPrv)

	return err
}
