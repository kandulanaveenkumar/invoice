package handlers

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/services/airwaybillhouse"
	"bitbucket.org/radarventures/forwarder-shipments/services/airwaybilllabels"
	"bitbucket.org/radarventures/forwarder-shipments/services/airwaybillmaster"
	"github.com/google/uuid"
)

func GetBillDetailsCommon(ctx *context.Context, billName string, req *dtos.GetBillsReq) (*dtos.BillDetailsRes, error) {
	res := dtos.BillDetailsRes{}
	switch billName {
	case globals.MasterAirwayBill:
		res, err := airwaybillmaster.NewAirwayBillMasterService().GetBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	case globals.LabelsAirwayBill:
		res, err := airwaybilllabels.NewAirwayBillLabelsService().GetBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	case globals.HouseAirwayBill:
		res, err := airwaybillhouse.NewAirwayBillHouseService().GetHouseAirwayBill(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	default:
		ctx.Log.Error("Invalid Bill Type: " + billName)
		return &res, nil
	}

}

func SaveBillDetailsCommon(ctx *context.Context, billName string, req *dtos.BillDetailsReq) (*dtos.SaveDetailsRes, error) {
	res := dtos.SaveDetailsRes{}
	switch billName {
	case globals.MasterAirwayBill:
		res, err := airwaybillmaster.NewAirwayBillMasterService().SaveBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	case globals.LabelsAirwayBill:
		res, err := airwaybilllabels.NewAirwayBillLabelsService().SaveBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	case globals.HouseAirwayBill:
		res, err := airwaybillhouse.NewAirwayBillHouseService().SaveHouseAirwayBill(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil
	default:
		ctx.Log.Error("Invalid Bill Type: " + billName)
		return &res, nil
	}

}

func GenerateBillDetailsCommon(ctx *context.Context, billName string, req *dtos.GetBillsReq) (*dtos.SaveDetailsRes, error) {
	res := dtos.SaveDetailsRes{}
	switch billName {
	case globals.MasterAirwayBill:
		res, err := airwaybillmaster.NewAirwayBillMasterService().GenerateBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil

	case globals.HouseAirwayBill:
		result, err := airwaybillhouse.NewAirwayBillHouseService().GenerateHouseAirwayBill(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return result, nil
	default:
		ctx.Log.Error("Invalid Bill Type: " + billName)
		return &res, nil
	}

}

func DownloadBillDetailsCommon(ctx *context.Context, billName string, req *dtos.DownloadReq) (*dtos.DownloadRes, error) {
	ctx.SetLoggingContext(req.ShipmentId,"DownloadBillDetailsCommon")
	res := dtos.DownloadRes{}
	switch billName {
	case globals.MasterAirwayBill:
		res, err := airwaybillmaster.NewAirwayBillMasterService().DownloadBillDetails(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return res, nil

	case globals.HouseAirwayBill:
		result, err := airwaybillhouse.NewAirwayBillHouseService().DownloadHouseAirwayBill(ctx, req)
		if err != nil {
			ctx.Log.Error("unable to fetch bills details" + billName)
			return nil, err
		}
		return result, nil
	default:
		ctx.Log.Error("Invalid Bill Type: " + billName)
		return &res, nil
	}

}

func GenerateMawbLinkCommon(ctx *context.Context, shipmentId uuid.UUID) (*dtos.AvailableMAWBRes, error) {

	res, err := airwaybillmaster.NewAirwayBillMasterService().GetAvailableMAWBForLinking(ctx, shipmentId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
