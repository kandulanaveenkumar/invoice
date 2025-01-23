package shipment

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
)

func (t *Shipment) GetShipmentsForPartner(ctx *context.Context, partnerID string) ([]*models.Shipment, error) {
	var result []*models.Shipment
	subQuery := ctx.DB.Debug().Table("line_items li").Select("distinct quote_id").Where("partner_id = ? ", partnerID)
	query := ctx.DB.Table(t.getTable(ctx)).Joins("JOIN quotes q ON shipments.quote_id = q.id").Where("q.id IN (?) AND shipments.is_deleted = false", subQuery)
	err := query.Debug().Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Shipment) GetShipmentsForCustomerInfo(ctx *context.Context, cid string) ([]*models.Shipment, error) {
	var res []*models.Shipment

	query := ctx.DB.Table(t.getTable(ctx)).
		Where("is_deleted = false")

	query = query.Where("company_id = ?", cid)

	err := query.Debug().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
