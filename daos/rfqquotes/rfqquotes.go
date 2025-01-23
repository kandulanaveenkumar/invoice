package rfqquotes

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type IRfqQuotes interface {
	Upsert(ctx *context.Context, m ...*models.RfqQuotes) error
	Get(ctx *context.Context, id string) ([]*models.RfqQuotes, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.RfqQuotes, error)
	Delete(ctx *context.Context, id string) error
	GetByQuote(ctx *context.Context, quoteId string) (*models.RfqQuotes, error)
	GetQuoteDetailsByRfqIds(ctx *context.Context, rfqIDs []string) ([]*models.RfqQuotes, error)
	GetRfqListingDetails(ctx *context.Context, rfqIds []string) ([]*models.RfqListingDetails, error)
	GetFirstQuoteForRfq(ctx *context.Context, rfqId string) (*models.RfqQuotes, error)
}

type RfqQuotes struct {
}

func NewRfqQuotes() IRfqQuotes {
	return &RfqQuotes{}
}

func (t *RfqQuotes) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfq_quotes"
}

func (t *RfqQuotes) Upsert(ctx *context.Context, m ...*models.RfqQuotes) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rfq_id"}, {Name: "quote_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rfq_id", "quote_id", "created_by", "updated_by", "updated_at"}),
	}).Create(m).Error

}
func (t *RfqQuotes) Get(ctx *context.Context, rfqId string) ([]*models.RfqQuotes, error) {
	var result []*models.RfqQuotes
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "rfq_id = ?", rfqId).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqquotes.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqQuotes) Delete(ctx *context.Context, id string) error {
	var result models.RfqQuotes
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfqquotes.", zap.Error(err))
		return err
	}

	return err
}

func (t *RfqQuotes) GetAll(ctx *context.Context, ids []string) ([]*models.RfqQuotes, error) {
	var result []*models.RfqQuotes
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqquotess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqquotess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqQuotes) GetByQuote(ctx *context.Context, quoteId string) (*models.RfqQuotes, error) {
	var result *models.RfqQuotes

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("quote_id = ?", quoteId).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqquotess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqQuotes) GetQuoteDetailsByRfqIds(ctx *context.Context, rfqIDs []string) ([]*models.RfqQuotes, error) {
	var result []*models.RfqQuotes

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" rq").Select("rq.quote_id").
		Where("rfq_id IN (?)", rfqIDs).Order("created_at desc").Scan(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqquotess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqQuotes) GetRfqListingDetails(ctx *context.Context, rfqIds []string) ([]*models.RfqListingDetails, error) {
	var result []*models.RfqListingDetails

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" rq").Select("rq.rfq_id, rq.quote_id, li.partner_id, li.sell,li.region_id").
		Joins("JOIN line_items li ON li.quote_id = rq.quote_id").
		Where("rq.rfq_id IN (?)", rfqIds).
		Where("rq.created_at = (SELECT MIN(rq2.created_at) FROM rfq_quotes rq2 WHERE rq2.rfq_id = rq.rfq_id)").Scan(&result).Error
	if err != nil {
		ctx.Log.Error("failed to get listing details", zap.Error(err), zap.Any("rfq_ids", rfqIds))
		return nil, err
	}
	return result, err
}

func (t *RfqQuotes) GetFirstQuoteForRfq(ctx *context.Context, rfqId string) (*models.RfqQuotes, error) {
	var result *models.RfqQuotes
	subquery := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("MIN(created_at)").Where("rfq_id = ?", rfqId)

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Select("quote_id", "rfq_id").Where("rfq_id = ? AND created_at = (?)", rfqId, subquery).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
