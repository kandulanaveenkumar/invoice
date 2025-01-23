package cronjobs

import (
	"time"

	"errors"

	"bitbucket.org/radarventures/forwarder-adapters/apis/id"
	"bitbucket.org/radarventures/forwarder-adapters/apis/misc"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/daos/card"
	"bitbucket.org/radarventures/forwarder-shipments/daos/rfq"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipment"
	rfqSer "bitbucket.org/radarventures/forwarder-shipments/services/rfq"
	"go.uber.org/zap"
)

type CardsPostExpiry struct {
	id         id.ID
	cardDb     card.ICard
	shipmentDb shipment.IShipment
	rfqDb      rfq.IRfq
	misc       misc.Misc
	rfqSer     rfqSer.IRfqService
}

func NewCardsPostExpiry() ICardsPostExpiry {
	return &CardsPostExpiry{
		id:         *id.New(config.Get().IdURL),
		cardDb:     card.NewCard(),
		shipmentDb: shipment.NewShipment(),
		rfqDb:      rfq.NewRfq(),
		misc:       *misc.New(config.Get().MiscURL),
		rfqSer:     rfqSer.NewRfqService(),
	}
}

type ICardsPostExpiry interface {
	CardsPostExpiry(ctx *context.Context) error
}

func (c *CardsPostExpiry) CardsPostExpiry(ctx *context.Context) error {

	RfqIdsExpiryList, err := c.rfqSer.GetRfqIdsWithExpiryStatus(ctx)
	if err != nil {
		ctx.Log.Error("Unable to get the bookingRequestIds", zap.Error(err))
		return err
	}
	var isBookedCount []bool
	if RfqIdsExpiryList == nil {
		ctx.Log.Error("Received nil RfqIdsExpiryList")
		return errors.New("received nil booking requests that are expired")
	}

	for _, id := range RfqIdsExpiryList {
		isBooking, err := c.cardDb.DeleteBookingRequestByID(ctx, id, time.Now().UTC())
		if err != nil {
			ctx.Log.Error("Unable to delete booking request", zap.String("BookingRequestID", id), zap.Error(err))
			return err
		}
		if isBooking {
			isBookedCount = append(isBookedCount, isBooking)
		}
	}
	ctx.Log.Info("Deleted expired booking request cards ", zap.Any("BookingRequestID", len(RfqIdsExpiryList)))
	ctx.Log.Info("Total expired booking request ", zap.Any("BookingRequestID", len(isBookedCount)))
	return nil

}
