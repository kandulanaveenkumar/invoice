package stock

import (
	"strings"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IStock interface {
	Upsert(ctx *context.Context, m ...*models.StockDetails) error
	Get(ctx *context.Context, id string) (*models.Stock, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Stock, error)
	Delete(ctx *context.Context, id string) error
	GetStockDetailByAirlineAndPort(ctx *context.Context, airline_id string, port_id string, region_id string) (*models.StockDetails, error)
	GetStockDetailsByNumber(ctx *context.Context, stock_number string) (*models.StockDetails, error)
	GetOldestAvailabelStockDetail(ctx *context.Context, airline_id string, port_id string, status string, region_id string, Q string) ([]*models.StockDetails, error)
	GetAllStockDetailsByAirlineAndPort(ctx *context.Context, airline_id string, port_id string, statuses []string, region_id string) ([]*models.StockDetails, error)
	GetAllStockDetailsByStatus(ctx *context.Context, statuses []string, regionId string, q string, Pg int64) ([]*models.StockDetails, error)
	GetAllStockCountsByStatus(ctx *context.Context, statuses []string, regionId string, q string) (int, error)
	GetStockCountByStatusWithAirlineAndPort(ctx *context.Context, statuses []string, q string, regionId string, Pg int64, ids []string) (map[string]map[string]int32, error)
	GetStockCountsPaginationByStatus(ctx *context.Context, statuses []string, q string, regionId string, ids []string) (int, error)
	GetStockDetailsByNumberId(ctx *context.Context, stock_number_id uuid.UUID) (*models.StockDetails, error)
}

type Stock struct {
}

func NewStock() IStock {
	return &Stock{}
}

func (t *Stock) getTable(ctx *context.Context) string {
	ctx.TenantID = "public"
	return ctx.TenantID + "." + "air_stocks"
}

func (t *Stock) Upsert(ctx *context.Context, m ...*models.StockDetails) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Stock) Get(ctx *context.Context, id string) (*models.Stock, error) {
	var result models.Stock
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get stock.", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (t *Stock) Delete(ctx *context.Context, id string) error {
	var result models.Stock
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete stock.", zap.Error(err))
		return err
	}

	return err
}

func (t *Stock) GetAll(ctx *context.Context, ids []string) ([]*models.Stock, error) {
	var result []*models.Stock
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get stocks.", zap.Error(err))
			return nil, err
		}
		return result, nil
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get stocks.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *Stock) GetStockDetailByAirlineAndPort(ctx *context.Context, airline_id string, port_id string, region_id string) (*models.StockDetails, error) {
	var stockDetail *models.StockDetails
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("liner = ? AND port =? AND region_id =?", airline_id, port_id, region_id).Order("created_at desc").Limit(1).Find(&stockDetail).Error
	if err != nil {
		ctx.Log.Error("Unable to get stocks.", zap.Error(err))
		return nil, err
	}
	return stockDetail, nil
}

func (t *Stock) GetStockDetailsByNumber(ctx *context.Context, stock_number string) (*models.StockDetails, error) {
	var stockDetail *models.StockDetails
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("stock_no = ? ", stock_number).Order("created_at desc").First(&stockDetail).Error
	if err != nil {
		ctx.Log.Error("Unable to get stocks.", zap.Error(err))
		return nil, err
	}

	return stockDetail, nil
}

func (t *Stock) GetOldestAvailabelStockDetail(ctx *context.Context, airline_id string, port_id string, status string, region_id string, Q string) ([]*models.StockDetails, error) {
	var stockDetail []*models.StockDetails
	var err error
	if Q == "" {
		err = ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("liner = ? AND port =? AND status=? AND region_id =?", airline_id, port_id, status, region_id).Order("created_at desc").Find(&stockDetail).Error
		if err != nil {
			ctx.Log.Error("Unable to get stocks.", zap.Error(err))
			return nil, err
		}
	} else {
		err = ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("liner = ? AND port =? AND status=? AND region_id =? AND stock_no ILIKE ?", airline_id, port_id, status, region_id, "%"+Q+"%").Order("created_at desc").Find(&stockDetail).Error
		if err != nil {
			ctx.Log.Error("Unable to get stocks.", zap.Error(err))
			return nil, err
		}
	}

	return stockDetail, nil
}

func (t *Stock) GetAllStockDetailsByAirlineAndPort(ctx *context.Context, airline_id string, port_id string, statuses []string, region_id string) ([]*models.StockDetails, error) {
	var stockDetails []*models.StockDetails
	var err error
	for _, status := range statuses {
		var stockDetail []*models.StockDetails
		err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("liner = ? AND port = ? AND status= ? AND region_id = ? ", airline_id, port_id, status, region_id).Order("created_at desc").Find(&stockDetail).Error
		if err != nil {
			ctx.Log.Error("Unable to get stocks.", zap.Error(err))
			return nil, err
		}

		stockDetails = append(stockDetails, stockDetail...)
	}

	return stockDetails, nil
}

func (t *Stock) GetAllStockDetailsByStatus(ctx *context.Context, statuses []string, regionId string, q string, Pg int64) ([]*models.StockDetails, error) {
	var limit, offset int64
	PageSize := 10
	offset = int64(PageSize) * (Pg - 1)
	limit = int64(PageSize)
	var stocks []*models.StockDetails
	if q != "" {
		err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("status in ? AND region_id =? AND stock_no ILIKE ?", statuses, regionId, q+"%").Limit(10).Scan(&stocks).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("status in ? AND region_id =? ", statuses, regionId).Limit(int(limit)).Offset(int(offset)).Scan(&stocks).Error
		if err != nil {
			return nil, err
		}
	}

	return stocks, nil
}

func (t *Stock) GetAllStockCountsByStatus(ctx *context.Context, statuses []string, regionId string, q string) (int, error) {
	stockCount := 0
	ctx.TenantID = "public"
	tablename := t.getTable(ctx)
	for _, status := range statuses {
		var count int

		query := `SELECT count(id) from ` + tablename + ` sd WHERE sd.status = $1 AND sd.region_id = $2 `

		if q != "" {
			query += `AND sd.stock_no ilike $3 `
			err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Raw(query, status, regionId, q+"%").Scan(&count).Error
			if err != nil {
				ctx.Log.Error("unable to fetch counts 1", zap.Error(err))
			}
		} else {
			err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Raw(query, status, regionId).Scan(&count).Error
			if err != nil {
				ctx.Log.Error("unable to fetch counts 2", zap.Error(err))
			}
		}
		stockCount += count
	}
	return stockCount, nil
}

func (t *Stock) GetStockCountByStatusWithAirlineAndPort(ctx *context.Context, statuses []string, q string, regionId string, Pg int64, ids []string) (map[string]map[string]int32, error) {
	stockCountMap := make(map[string]map[string]int32)
	var err error
	var limit, offset int64
	tablename := t.getTable(ctx)
	offset = int64(globals.PageSize) * (Pg - 1)
	limit = int64(globals.PageSize)
	for _, status := range statuses {
		var StocksCounts []*models.StockCount
		query := `SELECT sc.liner, sc.port, count(sc.id) AS stock_count FROM ` + tablename + ` sc `
		query += `INNER JOIN (SELECT liner , count(id) AS airline_stocks_count from ` + tablename + ` sd `
		if q != "" {
			query += ` WHERE sd.liner IN ('` + strings.Join(ids, "','") + `') AND sd.status=$2 AND sd.region_id = $1`
			query += ` GROUP BY liner ORDER BY airline_stocks_count DESC LIMIT 10) AS airline_stocks ON airline_stocks.liner = sc.liner WHERE sc.region_id = $1`
			query += ` GROUP BY sc.liner,sc.port`
			err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Raw(query, regionId, status).Scan(&StocksCounts).Error
			if err != nil {
				return nil, err
			}
		} else {
			query += ` WHERE sd.status = $1 AND sd.region_id = $2  `
			query += ` GROUP BY liner ORDER BY airline_stocks_count DESC LIMIT $3 OFFSET $4) AS airline_stocks ON airline_stocks.liner = sc.liner WHERE sc.region_id = $2 `
			query += ` GROUP BY sc.liner,sc.port `
			err = ctx.DB.WithContext(ctx.Request.Context()).Raw(query, status, regionId, limit, offset).Scan(&StocksCounts).Error
			if err != nil {
				return nil, err
			}
		}

		for _, StocksCount := range StocksCounts {
			if _, ok := stockCountMap[StocksCount.Liner]; !ok {
				stockCountMap[StocksCount.Liner] = make(map[string]int32, 0)
			}
			if _, ok := stockCountMap[StocksCount.Liner][StocksCount.Port]; !ok {
				stockCountMap[StocksCount.Liner][StocksCount.Port] = 0
			}
			stockCountMap[StocksCount.Liner][StocksCount.Port] += StocksCount.StockCount
		}
	}
	return stockCountMap, nil
}

func (t *Stock) GetStockCountsPaginationByStatus(ctx *context.Context, statuses []string, q string, regionId string, ids []string) (int, error) {
	stockCount := 0
	tablename := t.getTable(ctx)
	for _, status := range statuses {
		var count int
		if q != "" {
			query := `SELECT count(distinct  liner)  from ` + tablename + ` sd WHERE sd.status = $1 AND sd.region_id = $2 `
			query += `AND sd.liner in ('` + strings.Join(ids, "','") + `')`
			err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Raw(query, status, regionId, q).Scan(&count).Error
			if err != nil {
				ctx.Log.Error("unable to fetch counts 1", zap.Error(err))
			}

		} else {
			query := `SELECT count(distinct  liner)  from ` + tablename + ` sd WHERE sd.status = $1 AND sd.region_id = $2 `
			err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Raw(query, status, regionId).Scan(&count).Error
			if err != nil {
				ctx.Log.Error("unable to fetch counts 2", zap.Error(err))
			}
		}
		stockCount += count
	}
	return stockCount, nil
}

func (t *Stock) GetStockDetailsByNumberId(ctx *context.Context, stock_number_id uuid.UUID) (*models.StockDetails, error) {
	var stockDetail *models.StockDetails
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id = ? ", stock_number_id).Order("created_at desc").First(&stockDetail).Error
	if err != nil {
		ctx.Log.Error("Unable to get stocks.", zap.Error(err))
		return nil, err
	}

	return stockDetail, nil
}
