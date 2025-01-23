package cronjobs

import (
	"fmt"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipment"
	shipmenttracking "bitbucket.org/radarventures/forwarder-shipments/services/shipment-tracking"
	"go.uber.org/zap"
)

func RunContainerTrackingAutomation(ctx *context.Context) {

	ctx.Log.Info("RunContainerTrackingAutomation started")

	createdSince := time.Date(2024, time.October, 01, 0, 0, 0, 0, time.UTC)

	shipments, err := shipment.NewShipment().GetShipmentsSince(ctx, "", []string{"id"}, []string{constants.ShipmentTypeFCL}, &createdSince, []string{constants.ShipmentCompleted, constants.ShipmentCreated}, true)
	if err != nil {
		ctx.Log.Error("error while fetching shipments", zap.Error(err))
	}

	for i, shipment := range shipments {

		ctx.Log.Info(fmt.Sprintf("Processing Shipment %d/%d", i+1, len(shipments)),
			zap.String("id", shipment.Id.String()),
		)

		err := shipmenttracking.NewShipmentTrackingService().AutomateOceanShipment(ctx, shipment.Id.String())
		if err != nil {
			ctx.Log.Error("error automating ocean shpment")
			continue
		}

		ctx.Log.Info("Automation completed for shipment")
	}

	ctx.Log.Info("RunContainerTrackingAutomation completed")

}
