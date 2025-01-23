package cronjobs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"bitbucket.org/radarventures/forwarder-adapters/apis/id"
	"bitbucket.org/radarventures/forwarder-adapters/apis/misc"
	miscdtos "bitbucket.org/radarventures/forwarder-adapters/dtos/misc"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/upload"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/daos/document"
	"bitbucket.org/radarventures/forwarder-shipments/daos/invoice"
	"bitbucket.org/radarventures/forwarder-shipments/daos/invoicerequest"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	inv "bitbucket.org/radarventures/forwarder-shipments/services/invoice"
	"bitbucket.org/radarventures/forwarder-shipments/services/invoice/helper"
	"go.uber.org/zap"
)

type PendingInvoice struct {
	No               string
	InvoiceType      string
	ShipmentId       uuid.UUID
	CompanyId        uuid.UUID
	Id               uuid.UUID
	InvoiceRequestId uuid.UUID
	RegionId         uuid.UUID
	DocId            uuid.UUID
	CreatedBy        uuid.UUID
}

type IcaInvoice struct {
	no string
	Id uuid.UUID
}

func removeSpecialChars(str string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	cleanedStr := regex.ReplaceAllString(str, "")
	cleanedStr = strings.Join(strings.Fields(cleanedStr), " ")
	return cleanedStr
}

func FetchPendingInvoices(ctx *context.Context) {

	var pendingInvoices []PendingInvoice

	err := ctx.DB.Table("invoice_requests JOIN invoices ON invoices.invoice_request_id::TEXT = invoice_requests.id::TEXT").
		Where("is_completed", "false").
		Find(&pendingInvoices).
		Error
	if err != nil {
		ctx.Log.Error("unable to get pending invoices", zap.Error(err))
		return
	}

	shipmentDocs := []*models.Document{}

	miscService := *misc.New(config.Get().MiscURL)
	idService := id.New(config.Get().IdURL)

	for _, pendingInvoice := range pendingInvoices {

		ctx.Log.Info("prending invoices", zap.Any("shipment_id", pendingInvoice.ShipmentId), zap.Any("invoice_request_id", pendingInvoice.InvoiceRequestId))

		if pendingInvoice.DocId != uuid.Nil {
			continue
		}

		res, err := inv.NewInvoiceService().GetInvoiceRetrieval(ctx, pendingInvoice.Id)
		if err != nil {
			ctx.Log.Error("unable to get the response", zap.Error(err))
			return
		}

		shipment, err := shipment.NewShipment().Get(ctx, pendingInvoice.ShipmentId.String())
		if err != nil {
			ctx.Log.Error("unable to shipment", zap.Error(err))
			return
		}

		owner := globals.Internal
		var partnerName string
		if shipment.CompanyId == pendingInvoice.CompanyId && (shipment.RegionId == shipment.OriginRegionId && shipment.RegionId == shipment.DestRegionId) {
			if pendingInvoice.InvoiceType == constants.CustomerInvoice || pendingInvoice.InvoiceType == constants.CreditNote {
				owner = "Customer"
			}
		}

		if shipment.Type == constants.ShipmentTypeMisc {
			owner = globals.Internal
		}

		if pendingInvoice.InvoiceType == constants.DebitNote || pendingInvoice.InvoiceType == constants.VendorInvoice {

			partnerAccount, err := idService.GetPartner(ctx, pendingInvoice.CompanyId.String())
			if err != nil {
				return
			}

			if partnerAccount != nil {
				partnerName = removeSpecialChars(partnerAccount.Company.Name)
			}

			ctx.Log.Info("partner details", zap.Any("partnerName", partnerName))

			owner = "Partner"
		}

		if res != nil {

			fileName := pendingInvoice.InvoiceType + "-" + pendingInvoice.No + "." + upload.FileFormatPDF

			docRes, err := upload.New(config.Get().MiscURL).UploadToS3(ctx, &upload.UploadReq{
				File:        res.([]byte),
				Folder:      fmt.Sprintf("/companies/%v/shipments/%v", pendingInvoice.CompanyId, pendingInvoice.ShipmentId),
				FileName:    fileName,
				FileFormat:  upload.FileFormatPDF,
				ContentType: upload.ContentTypeApplication,
			})
			if err != nil {
				ctx.Log.Error("unable to upload to s3", zap.Error(err))
				return
			}

			shipmentDocs = append(shipmentDocs, &models.Document{
				Id:           uuid.New(),
				DocumentId:   docRes.DocumentId,
				Name:         fileName,
				Type:         pendingInvoice.InvoiceType,
				Owner:        owner,
				InstanceId:   pendingInvoice.ShipmentId,
				InstanceType: constants.WorkflowTypeShipment,
				RegionId:     pendingInvoice.RegionId,
				CreatedBy:    uuid.MustParse(config.Get().WizBotID),
				UpdatedBy:    uuid.MustParse(config.Get().WizBotID),
			})

			err = invoice.NewInvoice().Update(ctx, &models.Invoice{
				DocId: docRes.DocumentId,
				ID:    pendingInvoice.Id,
			})
			if err != nil {
				ctx.Log.Error("unable to update the invoice", zap.Error(err))
				return
			}

			err = invoicerequest.NewInvoiceRequest().Update(ctx, &models.InvoiceRequest{
				IsCompleted:  true,
				IsSuccessful: true,
				ID:           pendingInvoice.InvoiceRequestId,
			})
			if err != nil {
				ctx.Log.Error("unable to update the invoice", zap.Error(err))
				return
			}

			adminName := ""
			updatedBy := pendingInvoice.CreatedBy.String()
			if ctx.Account != nil {
				adminDetails, err := idService.GetAccountInternal(ctx, pendingInvoice.CreatedBy.String())
				if err != nil {
					ctx.Log.Error("error while getting account", zap.Error(err))
					return
				}
				if adminDetails != nil {
					adminName = adminDetails.Name
				}
			}

			taggedMembers := make(map[string]interface{})
			taggedMembers[updatedBy] = adminName

			ctx.Log.Info("admin details", zap.Any("adminName", adminName), zap.Any("updatedBy", updatedBy))

			_, err = miscService.SendCollab(ctx, &miscdtos.CollabMsg{
				RefID:         pendingInvoice.ShipmentId.String(),
				RefType:       "shipment",
				Msg:           fmt.Sprintf("Hi %s invoice copy has been retrieved for the invoice number %s. Please refer to the document tab for the uploaded invoice copy", adminName, pendingInvoice.No),
				TaggedMembers: taggedMembers,
				TaskRegionID:  pendingInvoice.RegionId.String(),
				ChatType:      "internal_chat",
				CreatedBy:     config.Get().WizBotID,
				AccountId:     updatedBy,
			})
			if err != nil {
				ctx.Log.Error("failed to send collab message", zap.Error(err))
			}

		} else {

			err = helper.NewHelper().DeleteInvoice(ctx, pendingInvoice.Id)
			if err != nil {
				ctx.Log.Error("unable to delete invoice", zap.Error(err))
				return
			}

			err = invoicerequest.NewInvoiceRequest().Update(ctx, &models.InvoiceRequest{
				IsCompleted:  true,
				IsSuccessful: false,
				ID:           pendingInvoice.InvoiceRequestId,
			})
			if err != nil {
				ctx.Log.Error("unable to update the invoice", zap.Error(err))
				return
			}
		}

	}

	if len(shipmentDocs) > 0 {
		err := document.NewDocument().Upsert(ctx, shipmentDocs...)
		if err != nil {
			ctx.Log.Error("unable to upsert invoice docs", zap.Error(err))
			return
		}
	}

}

func IcaInvoices(ctx *context.Context) {
	var icaInvoices []IcaInvoice

	err := ctx.DB.Table("invoices").
		Where("is_ica_invoice = ? AND cost_booking_generated = ?", true, false).
		Find(&icaInvoices).
		Error

	if err != nil {
		ctx.Log.Error("unable to get ica invoices", zap.Error(err))
		return
	}

	shipmentDocs := []*models.Document{}

	miscService := *misc.New(config.Get().MiscURL)
	idService := id.New(config.Get().IdURL)

	for _, icaInvoice := range icaInvoices {

		res, err := inv.NewInvoiceService().GenerateIcaInvoice(ctx, icaInvoice.Id)
		if err != nil {
			ctx.Log.Error("unable to get the response", zap.Error(err))
			return
		}

		
	}

}
