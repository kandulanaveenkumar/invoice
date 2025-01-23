package shipment

import (
	"fmt"
	"strconv"
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IShipment interface {
	Upsert(ctx *context.Context, m ...*models.Shipment) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.Shipment) error
	Update(ctx *context.Context, m *models.Shipment) error
	UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.Shipment) error

	Get(ctx *context.Context, id string) (*models.Shipment, error)
	GetByQuote(ctx *context.Context, qid string) (*models.Shipment, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Shipment, error)
	Delete(ctx *context.Context, id string) error
	CheckCode(ctx *context.Context, code string) (bool, error)

	GetShipmentCounts(ctx *context.Context, adminIds, customerIds []string) ([]*models.ShipmentCount, error)
	GetShipmentList(ctx *context.Context, adminIds, customerIds []string) ([]*models.ShipmentList, error)
	GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error)

	GetShipmentCardFilter(ctx *context.Context, ids []string, req dtos.CardInstanceReq) ([]*models.Shipment, error)

	GetShipmentsForCompany(ctx *context.Context, cid string) ([]*models.Shipment, error)

	GetPaginatedConsol(ctx *context.Context, req *dtos.ConsolGetReq) ([]*dtos.ConsolParameters, error)
	GetCountsByStatusPaginatedConsol(ctx *context.Context, req *dtos.ConsolGetReq) (map[string]int64, error)
	UpsertConsolShipment(ctx *context.Context, m *models.ConsolShipment) error
	GetShipmentCountWithConsolId(ctx *context.Context, id string) (int, error)
	GetConsol(ctx *context.Context, id string) (*models.ConsolShipment, error)
	DeleteConsol(ctx *context.Context, id string, deleteReason string, adminId uuid.UUID) error
	GetConsolMatchedShipments(ctx *context.Context, regionId string, pol string, pod string, isPolMatchingOnly bool) ([]*models.ConsolMatchedShipment, error)
	UpsertConsol(ctx *context.Context, consol *models.ConsolShipment) error
	UpsertConsolWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ConsolShipment) error
	GetConsolByQuoteId(ctx *context.Context, quoteId string) (*models.ConsolShipment, error)
	GetShipmentsByConsolId(ctx *context.Context, consolId string) ([]*models.Shipment, error)

	GetShipmentsSince(ctx *context.Context, cid string, selectFields []string, shipmentTypes []string, createdSince *time.Time, excludedStatus []string, MasterBillNoCheck bool) ([]*models.Shipment, error)
	GetCompanyDasboardBookingsCount(ctx *context.Context, cids []string) (int64, error)
	GetShipmentsForPartner(ctx *context.Context, partnerID string) ([]*models.Shipment, error)
	GetShipmentsForCustomerInfo(ctx *context.Context, cid string) ([]*models.Shipment, error)
	CustomerPaginationAciveBookingsCount(ctx *context.Context, cids []string) ([]*models.ActiveShipmentCount, error)

	GetLiveBookingNotBreached(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)
	GetLiveBookingBreached(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)

	GetLiveBookings(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error)
	GetFunnelViewCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)
	GetforFunnelFilterBookings(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error)
	GetInsightTimeRangeProcurement(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.BookingJobs, error)
	GetTimeRangeForInsight(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.BookingJobs, error)
	GetInsightDetails(ctx *context.Context, filter *dtos.DashboardFilters, ports []string) (*dtos.InsightDetailsResponse, error)
	GetInsightDetailCompanies(ctx *context.Context, filter *dtos.DashboardFilters, ports []string) ([]*dtos.Company, error)
	GetDSRShipments(ctx *context.Context, req *dtos.DSRShipmentParamters) ([]models.ShipmentQuoteDetailsForDSR, error)
	GetAllCompletedShipments(ctx *context.Context, cid string) ([]*dtos.Shipment, error)
	GetAllUnlockedShipments(ctx *context.Context) ([]*models.Shipment, error)
	UpdateShipmentsAndRfqs(ctx *context.Context, req *dtos.CustomerNameChangeReq) error
	UpdateShipmentLock(ctx *context.Context, id string, isLocked bool) error
}

type Shipment struct {
}

func NewShipment() IShipment {
	return &Shipment{}
}

func (t *Shipment) getTable(ctx *context.Context) string {

	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}

	return ctx.TenantID + "." + "shipments"
}
func (t *Shipment) getConsolShipmentsTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "consol_shipments"
}
func (t *Shipment) getQuotesTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "quotes"
}
func (t *Shipment) getRfqsTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfqs"
}
func (t *Shipment) Upsert(ctx *context.Context, m ...*models.Shipment) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Shipment) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.Shipment) error {
	return tx.Table(t.getTable(ctx)).Save(m).Error
}

func (t *Shipment) Get(ctx *context.Context, id string) (*models.Shipment, error) {
	var result *models.Shipment

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id = ?", id).First(&result).Error
	if err != nil {
		ctx.Log.Error("unable to get shipment", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Shipment) GetByQuote(ctx *context.Context, qid string) (*models.Shipment, error) {
	var result *models.Shipment
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("quote_id = ?", qid).First(&result).Error
	if err != nil {
		ctx.Log.Error("unable to get shipment", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Shipment) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id = ?", id).UpdateColumn("is_deleted", true).UpdateColumn("updated_at", time.Now().UTC()).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipment.", zap.Error(err))
		return err
	}

	return err
}

func (t *Shipment) GetAll(ctx *context.Context, ids []string) ([]*models.Shipment, error) {
	var result []*models.Shipment
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipments.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipments.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Shipment) CheckCode(ctx *context.Context, code string) (bool, error) {
	var result *models.Shipment
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("code = ?", code).First(&result).Error
	return (err == nil && result.Id != uuid.Nil), err
}

func (t *Shipment) GetShipmentCardFilter(ctx *context.Context, ids []string, req dtos.CardInstanceReq) ([]*models.Shipment, error) {
	var results []*models.Shipment
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id in (?) and (consol_id is null or consol_id = ? )", ids, uuid.Nil).Where(&req)
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while getting bulk  --DB", zap.Any("bid", ids), zap.Error(err))
	}
	return results, nil
}

func (t *Shipment) Update(ctx *context.Context, m *models.Shipment) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}

func (t *Shipment) UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.Shipment) error {
	return tx.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}

func (t *Shipment) GetPaginatedConsol(ctx *context.Context, req *dtos.ConsolGetReq) ([]*dtos.ConsolParameters, error) {
	var results []*dtos.ConsolParameters
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" consol").
		Select(`
	consol.id AS consol_id, 
	consol.quote_id AS consol_quote_id,
	consol.code AS consol_code,
	sc.type AS container_info,
	consol.occupied_cbm AS occupied_volume,
	q.etd AS etd,
	consol.pol_name AS pol,
	consol.pod_name AS pod,
	consol.created_at AS created_on,
	consol.delete_reason as drop_job_reason,
	CASE 
		WHEN consol.is_deleted = true THEN 'dropped' 
		WHEN consol.status ILIKE '%Booking completed%' THEN 'completed' 
		ELSE 'active' 
	END AS status,
	sc.cbm AS total_volume`).
		Joins("JOIN quotes AS q ON consol.quote_id = q.id").
		Joins("JOIN shipment_containers AS sc ON consol.id = sc.shipment_id").
		Where("consol.region_id = ? and consol.type = ? ", req.RegionId, globals.BookingTypeCONSOL)

	if req.Q == "" {
		switch req.Status {
		case globals.ConsolDropped:
			query.Where("consol.is_deleted = true")
		case globals.ConsolActive:
			query.Where("consol.status NOT ILIKE '%Booking completed%' AND consol.is_deleted = false")
		case globals.ConsolCompleted:
			query.Where("consol.status ILIKE '%Booking completed%' AND consol.is_deleted = false")
		case globals.ConsolShift:
			query.Where("lower(consol.status) != lower(?) AND consol.is_deleted = false AND consol.can_link_to_consol = true", constants.ShipmentCompleted)
		}

		if req.Id != "" {
			query.Where("consol.id = ?", req.Id)
		}

		if req.Pol != "" {
			query.Where("consol.pol = ?", req.Pol)
		}

	} else {
		query.Where("consol.code ilike (?)", req.Q+"%")

	}
	query.Order("consol.created_at DESC")

	size := config.Get().PageSize
	pg, err := strconv.Atoi(ctx.Query("pg"))
	if err != nil {
		ctx.Log.Error("unable to parse pageNo")
		return nil, err
	}

	offset := size * (pg - 1)
	if req.Q == "" && req.Status != globals.ConsolShift {
		if offset > 0 {
			query.Offset(offset)
		}

		query.Limit(size)
	}

	err = query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting consols", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (t *Shipment) GetCountsByStatusPaginatedConsol(ctx *context.Context, req *dtos.ConsolGetReq) (map[string]int64, error) {

	counts := make(map[string]int64)

	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" consol").
		Select(`
			COUNT(CASE WHEN consol.status NOT ILIKE '%Booking completed%' AND consol.is_deleted = false THEN 1 END) AS active_count,
			COUNT(CASE WHEN consol.is_deleted = true THEN 1 END) AS dropped_count,
			COUNT(CASE WHEN lower(status) != lower('`+constants.ShipmentCompleted+`') AND consol.is_deleted = false AND consol.can_link_to_consol = true THEN 1 END) AS can_shift,
			COUNT(CASE WHEN consol.status ILIKE '%Booking completed%' AND consol.is_deleted = false THEN 1 END) AS completed_count`).
		Where("consol.type = ? AND region_id = ?", globals.BookingTypeCONSOL, req.RegionId)

	if req.Id != "" {
		query = query.Where("consol.id = ?", req.Id)
	}
	if req.Q != "" {
		query = query.Where("consol.code ILIKE ?", req.Q)
	}
	if req.Pol != "" {
		query = query.Where("consol.pol = ?", req.Pol)
	}

	var result struct {
		ActiveCount    int64
		DroppedCount   int64
		CanShift       int64
		CompletedCount int64
	}
	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	counts[globals.ConsolActive] = result.ActiveCount
	counts[globals.ConsolDropped] = result.DroppedCount
	counts[globals.ConsolShift] = result.CanShift
	counts[globals.ConsolCompleted] = result.CompletedCount

	// pendingQuery := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getRfqsTable(ctx)+" rfq").
	// 	Select("count(rfq.id)").
	// 	Where("rfq.type = ? AND rfq.status = ? AND rfq.is_deleted = false", globals.BookingTypeCONSOL, globals.ConsolStatusBuyTBA).
	// 	Where("rfq.region_id = ?", req.RegionId)

	// if req.Id != "" {
	// 	pendingQuery = pendingQuery.Where("rfq.id = ?", req.Id)
	// }
	// if req.Q != "" {
	// 	pendingQuery = pendingQuery.Where("rfq.code = ?", req.Q)
	// }
	// if req.Pol != "" {
	// 	pendingQuery = pendingQuery.Where("rfq.pol = ?", req.Pol)
	// }

	// var pendingCount int64
	// err = pendingQuery.Scan(&pendingCount).Error
	// if err != nil {
	// 	return nil, err
	// }
	// counts[globals.ConsolPending] = int(pendingCount)

	return counts, nil

}

func (t *Shipment) UpsertConsolShipment(ctx *context.Context, m *models.ConsolShipment) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getConsolShipmentsTable(ctx)).Save(m).Error
}

func (t *Shipment) GetShipmentCountWithConsolId(ctx *context.Context, id string) (int, error) {
	var shipments_count int64

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" shipment").Where("shipment.consol_id =?", id).
		Select("shipment.id").Where("shipment.is_deleted::BOOLEAN = false")

	q.Count(&shipments_count)
	err := q.Begin().Error
	if err != nil {
		return int(shipments_count), err
	}

	return int(shipments_count), nil
}
func (t *Shipment) GetConsol(ctx *context.Context, id string) (*models.ConsolShipment, error) {
	var result *models.ConsolShipment
	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getConsolShipmentsTable(ctx)).Where("id =?", id)
	err := q.Scan(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Log.Error("error while getting Consol")
	}
	return result, nil
}

func (t *Shipment) DeleteConsol(ctx *context.Context, id string, deleteReason string, adminId uuid.UUID) error {

	consol, err := t.GetConsol(ctx, id)
	if err != nil {
		ctx.Log.Error("error while getting Consol")
	}
	consol.DeleteReason = deleteReason
	consol.UpdatedAt = time.Now().UTC()
	consol.UpdatedBy = adminId.String()
	consol.IsDeleted = true

	return t.UpsertConsol(ctx, consol)
}

func (t *Shipment) GetConsolMatchedShipments(ctx *context.Context, regionId string, pol string, pod string, isPolMatchingOnly bool) ([]*models.ConsolMatchedShipment, error) {
	var result []*models.ConsolMatchedShipment
	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" s").Select("s.id, s.code ,s.pol, s.pod , s.cargo_ready_date,s.created_at,s.created_by,s.region_id, s.origin_region_id, s.dest_region_id,s.occupied_cbm,s.company_id, q.etd").Joins("join quotes q on q.id =s.quote_id ").Where(" s.type = 'LCL' and s.is_deleted::BOOLEAN =false").Where("s.pol =? ", pol)
	if regionId != "" {
		q.Where("s.region_id = ? ", regionId)
	}
	if isPolMatchingOnly {
		q.Where("s.pod != ?", pod)
	} else {
		q.Where("s.pod = ?", pod)
	}
	q.Order("s.created_at ASC")

	err := q.Find(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Log.Error("error while getting Consol")
	}
	return result, nil
}

func (t *Shipment) UpsertConsol(ctx *context.Context, consol *models.ConsolShipment) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getConsolShipmentsTable(ctx)).Save(consol).Error
}

func (t *Shipment) GetShipmentsForCompany(ctx *context.Context, cid string) ([]*models.Shipment, error) {
	var result []*models.Shipment
	err := ctx.DB.Table(t.getTable(ctx)).Where("is_deleted = false AND company_id = ?", cid).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipments.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *Shipment) UpsertConsolWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ConsolShipment) error {
	return tx.Table(t.getConsolShipmentsTable(ctx)).Save(m).Error
}

func (t *Shipment) GetConsolByQuoteId(ctx *context.Context, quoteId string) (*models.ConsolShipment, error) {
	var result *models.ConsolShipment
	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getConsolShipmentsTable(ctx)).Where("quote_id =?", quoteId)
	err := q.Scan(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Log.Error("error while getting Consol")
	}
	return result, nil
}

func (t *Shipment) GetShipmentsByConsolId(ctx *context.Context, consolId string) ([]*models.Shipment, error) {
	var result []*models.Shipment
	err := ctx.DB.Table(t.getTable(ctx)).Where("consol_id = ?", consolId).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipments.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *Shipment) GetShipmentsSince(ctx *context.Context, cid string, selectFields []string, shipmentTypes []string, createdSince *time.Time, excludedStatus []string, MasterBillNoCheck bool) ([]*models.Shipment, error) {
	var result []*models.Shipment
	tx := ctx.DB.Table(t.getTable(ctx))

	if selectFields != nil {
		tx.Select(selectFields)
	}

	tx.Where("is_deleted = false")

	if MasterBillNoCheck {
		tx.Where("master_bill_no != ''")
	}

	if cid != "" {
		tx.Where("company_id = ?", cid)
	}

	if len(excludedStatus) > 0 {
		tx.Where("status not IN (?)", excludedStatus)
	}

	if shipmentTypes != nil {
		tx.Where("type IN (?)", shipmentTypes)
	}

	if createdSince != nil {
		tx.Where("created_at >= ?", createdSince)
	}

	err := tx.Debug().Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipments.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

const DashboardListCount = 5

func (t *Shipment) GetCompanyDasboardBookingsCount(ctx *context.Context, cids []string) (int64, error) {
	var bookingCount int64

	query := ctx.DB.Table(t.getTable(ctx)).
		Where("is_deleted = false").
		Not("status IN (?)", "Booking completed")

	query = query.Where("company_id IN (?)", cids)

	err := query.Debug().Count(&bookingCount).Error
	if err != nil {
		return 0, err
	}

	return bookingCount, nil
}

func (t *Shipment) CustomerPaginationAciveBookingsCount(ctx *context.Context, cids []string) ([]*models.ActiveShipmentCount, error) {
	var bookingCount []*models.ActiveShipmentCount

	query := ctx.DB.Debug().Table(t.getTable(ctx)).
		Where("is_deleted = ?", false).
		Not("status = ?", "Booking completed").
		Where("company_id IN (?)", cids).
		Select("count(distinct(id)) as count, company_id").
		Group("company_id")

	err := query.Debug().Find(&bookingCount).Error
	if err != nil {
		return nil, err
	}

	return bookingCount, nil
}

func (t *Shipment) GetDSRShipments(ctx *context.Context, req *dtos.DSRShipmentParamters) ([]models.ShipmentQuoteDetailsForDSR, error) {
	var shipments []models.ShipmentQuoteDetailsForDSR
	query := ctx.DB.Debug().Table(t.getTable(ctx)+" s").Select(`s.id, s.code, s.house_bill_nos, s.cargo_ready_date, 
	s.pol_name, s.pod_name, s.status, s.incoterm, s.occupied_cbm, s.occupied_weight, s.occupied_volume_weight, 
	s.created_at, s.teus, q.etd, q.eta, q.free_days, q.transit_days,q.liner,q.vessel_name,q.voyage_no,s.type,COALESCE(SUM(CASE WHEN s.type IN ('LCL', 'AIR') THEN sp.count ELSE 0 END), 0) AS count,
	MAX(CASE WHEN s.is_door_pickup = TRUE AND sl.type = 'origin' THEN sl.address ELSE NULL END) AS door_pickup, 
    MAX(CASE WHEN s.is_door_delivery = TRUE AND sl.type = 'destination' THEN sl.address ELSE NULL END) AS door_delivery`).
		Joins("JOIN quotes as q ON s.quote_id = q.id").
		Joins("LEFT JOIN shipment_products AS sp ON s.id = sp.shipment_id AND s.type IN ('LCL', 'AIR')").
		Joins("LEFT JOIN shipment_locations AS sl ON s.id = sl.shipment_id").
		Where("s.company_id = ? and s.type in ? and s.is_deleted = false ", req.CompanyId, req.Types).
		Group("s.id, s.code, s.house_bill_nos, s.cargo_ready_date, s.pol_name, s.pod_name, s.status, s.incoterm, s.occupied_cbm, s.occupied_weight, s.occupied_volume_weight, s.created_at, s.teus, q.etd, q.eta, q.free_days, q.transit_days,q.liner,q.vessel_name,q.voyage_no, s.type").
		Order("s.created_at DESC")
	err := query.Scan(&shipments).Error
	if err != nil {
		ctx.Log.Error("Error fetching shipment details for DSR", zap.Error(err))
		return nil, err
	}

	return shipments, nil
}

func (t *Shipment) GetAllCompletedShipments(ctx *context.Context, cid string) ([]*dtos.Shipment, error) {
	var shipments []*dtos.Shipment
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Select("id,type").
		Where("is_deleted = false").
		Where("status = ?", constants.ShipmentCompleted).
		Where("company_id = ?", cid).
		Order("created_at DESC").
		Scan(&shipments).Error
	if err != nil {
		ctx.Log.Error("Error fetching shipment details", zap.Error(err))
		return nil, err
	}
	return shipments, err
}

func (t *Shipment) GetAllUnlockedShipments(ctx *context.Context) ([]*models.Shipment, error) {
	var result []*models.Shipment
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Where("is_deleted::BOOLEAN = false").
		Where("is_shipment_locked::BOOLEAN = false").
		Order("created_at DESC").
		Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipments.", zap.Error(err))
		return nil, err
	}

	return result, err

}

func (t *Shipment) UpdateShipmentsAndRfqs(ctx *context.Context, req *dtos.CustomerNameChangeReq) error {

	err := ctx.DB.Debug().Exec(fmt.Sprintf(`UPDATE %s set company_name = ? WHERE company_id = ? `, t.getTable(ctx)), req.CompanyName, req.CompanyId).Error
	if err != nil {
		ctx.Log.Error("Unable to update shipments", zap.Error(err))
		return err
	}

	err = ctx.DB.Debug().Exec(fmt.Sprintf(`UPDATE %s set company_name = ? WHERE company_id = ? `, t.getRfqsTable(ctx)), req.CompanyName, req.CompanyId).Error
	if err != nil {
		ctx.Log.Error("Unable to update rfqs", zap.Error(err))
		return err
	}

	return nil
}

func (t *Shipment) UpdateShipmentLock(ctx *context.Context, id string, isLocked bool) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", id).Update("is_shipment_locked", isLocked).Error
}
