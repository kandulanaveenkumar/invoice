package sis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/radarventures/forwarder-shipments/config"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ISIS interface {
	Get(ctx *context.Context, id string) (*dtos.SISInfoModelAir, error)
	GetForGenerateJobNumber(ctx *context.Context, code string) ([]string, error)
	GetForGenerateSISAirJobNumber(ctx *context.Context, code string) ([]string, error)
	UpsertAirSisInfo(ctx *context.Context, obj *dtos.SISInfoModelAir, fname string, isProcessed bool, data []byte) error
	UpdateSISData(ctx *context.Context, model *dtos.SISInfoModelAir, id string) error
	GetSISInfo(ctx *context.Context, id string) (*dtos.SISInfoModel, error)
	UpdateSisInfoStatus(ctx *context.Context, id string, filename string, data []byte) error
	Upsert(ctx *context.Context, obj *dtos.SISInfo, fname string, data []byte) error
	GetSISInfoReq(ctx *context.Context, req *dtos.GetSISInfoReq) (int64, []*dtos.SISInfo, error)
	GetContent(ctx *context.Context, id string) (string, error)
	GetSISInfoRequest(ctx *context.Context, id string) (*dtos.SISInfo, error)
	SaveSISInfo(ctx *context.Context, model *dtos.SISInfoModel, isISF, isSIS bool, data []byte) error
}

type SIS struct {
}

func NewSIS() ISIS {
	return &SIS{}
}

func (t *SIS) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "sis_info_air"
}

func (t *SIS) getSISTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "sis_info"
}

func (t *SIS) Get(ctx *context.Context, id string) (*dtos.SISInfoModelAir, error) {
	var sisData string

	q := fmt.Sprintf(`SELECT sis_data FROM %s WHERE shipment_id = ?`, t.getTable(ctx))
	err := ctx.DB.WithContext(ctx.Request.Context()).Raw(q, id).Scan(&sisData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.Log.Warn("No data found for the given shipment ID", zap.String("shipment_id", id))
			return nil, nil
		}
		ctx.Log.Error("Unable to get sis info air", zap.Error(err))
		return nil, err
	}

	if sisData == "" {
		ctx.Log.Warn("No JSON data retrieved for the given shipment ID", zap.String("shipment_id", id))
		return nil, nil
	}

	var model dtos.SISInfoModelAir
	err = json.Unmarshal([]byte(sisData), &model)
	if err != nil {
		ctx.Log.Error("Unable to unmarshal sis data", zap.Error(err))
		return nil, err
	}

	return &model, nil
}

func (t *SIS) GetForGenerateJobNumber(ctx *context.Context, code string) ([]string, error) {
	var jobNums []string
	param := "%" + code + "%"

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("sis_data->>'job_no'").Where("sis_data->>'job_no' LIKE ?", param).Pluck("sis_data->>'job_no'", &jobNums).Error

	if err != nil {
		ctx.Log.Error("Unable to retrieve job numbers", zap.Error(err))
		return nil, err
	}

	return jobNums, nil
}

func (t *SIS) GetForGenerateSISAirJobNumber(ctx *context.Context, code string) ([]string, error) {
	var jobNums []string
	param := "%" + code + "%"
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("job_number").Where("job_number LIKE ?", param).Pluck("job_number", &jobNums).Error

	if err != nil {
		ctx.Log.Error("Unable to retrieve job numbers", zap.Error(err))
		return nil, err
	}

	return jobNums, nil
}

func (t *SIS) UpsertAirSisInfo(ctx *context.Context, obj *dtos.SISInfoModelAir, fname string, isProcessed bool, data []byte) error {
	if obj == nil {
		return errors.New("nil sis info cannot be saved")
	}

	sisData, err := json.Marshal(obj)
	if err != nil {
		ctx.Log.Error("Error serializing sis data", zap.Error(err))
		return err
	}

	values := map[string]interface{}{
		"shipment_id":           obj.ShipmentId,
		"shipment_code":         obj.ShipmentCode,
		"job_number":            obj.JobNo,
		"file_name":             fname,
		"service_type":          obj.ShipmentType,
		"sis_data":              sisData,
		"sis_processed":         isProcessed,
		"sis_data_out_contents": string(data),
		"created_at":            time.Now().UTC(),
		"updated_at":            time.Now().UTC(),
	}

	err = ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shipment_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"job_number", "file_name", "service_type", "sis_data", "sis_processed", "sis_data_out_contents", "updated_at"}),
		}).
		Create(values).Error

	if err != nil {
		ctx.Log.Error("Error saving the request", zap.Error(err), zap.Any("sis", obj))
		return err
	}

	return nil
}

func (t *SIS) UpdateSISData(ctx *context.Context, model *dtos.SISInfoModelAir, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("booking_id = ?", id).
		Update("sis_data", model).Error

	if err != nil {
		ctx.Log.Error("[sis] UpdateSISData: Failed to update", zap.Error(err))
		return err
	}

	return nil
}

func (s *SIS) GetSISInfo(ctx *context.Context, id string) (*dtos.SISInfoModel, error) {
	var model dtos.SISInfoModel
	q := fmt.Sprintf(`SELECT sis_data FROM %s WHERE shipment_id = ?`, s.getSISTable(ctx))
	err := ctx.DB.WithContext(ctx.Request.Context()).Raw(q, id).Scan(&model).Error
	if err != nil {
		ctx.Log.Error("unable to get sis info", zap.Error(err))
		return nil, err
	}

	return &model, nil
}

func (s *SIS) UpdateSisInfoStatus(ctx *context.Context, id string, filename string, data []byte) error {

	query := fmt.Sprintf(`update %s set status_file_name=$1,sis_status_data_out_contents=$2 where id = $3`, s.getSISTable(ctx))
	err := ctx.DB.WithContext(ctx.Request.Context()).Raw(query, filename, string(data), id).Error
	return err
}

func (s *SIS) Upsert(ctx *context.Context, obj *dtos.SISInfo, fname string, data []byte) error {

	if obj == nil {
		return errors.New("nil sis info cannot be saved")
	}

	id, _ := uuid.Parse(obj.Id)
	var q string

	if id == uuid.Nil {
		obj.Id = uuid.New().String()
		q = fmt.Sprintf(`INSERT INTO %s (id, sis_booking, shipment_code, quote_code, cargo_ready_date, service_type, is_dangerous, hs_code, confirmed_date,
			origin_port_id, origin_port_name, dest_port_id, dest_port_name, imo_code, file_name, contents, incoterms, handling_agent_code, released_agent_code) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) 
			ON CONFLICT (file_name) DO UPDATE SET 
				cargo_ready_date = EXCLUDED.cargo_ready_date,
				service_type = EXCLUDED.service_type,
				is_dangerous = EXCLUDED.is_dangerous,
				hs_code = EXCLUDED.hs_code,
				confirmed_date = EXCLUDED.confirmed_date,
				origin_port_id = EXCLUDED.origin_port_id,
				origin_port_name = EXCLUDED.origin_port_name,
				dest_port_id = EXCLUDED.dest_port_id,
				dest_port_name = EXCLUDED.dest_port_name,
				imo_code = EXCLUDED.imo_code,
				contents = EXCLUDED.contents,
				incoterms = EXCLUDED.incoterms,
				handling_agent_code = EXCLUDED.handling_agent_code,
				released_agent_code = EXCLUDED.released_agent_code,
				updated_at = now()`, s.getSISTable(ctx))
	} else {
		q = fmt.Sprintf(`UPDATE %s SET sis_booking = $2,
			shipment_code = $3,
			quote_code = $4,
			cargo_ready_date = $5,
			service_type = $6,
			is_dangerous = $7,
			hs_code = $8,
			confirmed_date = $9,
			origin_port_id = $10,
			origin_port_name = $11,
			dest_port_id = $12,
			dest_port_name = $13,
			imo_code = $14,
			file_name = $15,
			contents = $16,
			incoterms = $17,
			handling_agent_code = $18,
			released_agent_code = $19,
			updated_at = now()
			WHERE id = $1`, s.getSISTable(ctx))
	}

	err := ctx.DB.WithContext(ctx.Request.Context()).Debug().Exec(q,
		obj.Id,
		obj.SisBooking,
		obj.ShipmentCode,
		obj.ShipmentRequestCode,
		obj.CargoReadyDate,
		obj.ServiceType,
		obj.IsDangerous,
		obj.HsCode,
		obj.RequestedOn,
		obj.OriginPortId,
		obj.OriginPortName,
		obj.DestPortId,
		obj.DestPortName,
		obj.ImoCode,
		fname,
		string(data),
		obj.Incoterms,
		obj.HandlingAgentCode,
		obj.ReleasedAgentCode,
	).Error

	if err != nil {
		ctx.Log.Error("Error saving the request", zap.Error(err), zap.Any("sis", obj.SisBooking))
		return err
	}

	return nil
}

func (s *SIS) GetSISInfoReq(ctx *context.Context, req *dtos.GetSISInfoReq) (int64, []*dtos.SISInfo, error) {
	var sisInfos []*dtos.SISInfo
	var totalCount int64

	pageSize := int(config.Get().PageSize)
	page := int(req.Pg)
	offset := (page - 1) * pageSize

	query := ctx.DB.WithContext(ctx.Request.Context()).Debug().
		Table(s.getSISTable(ctx)).
		Where("is_booked = false").
		Where(s.sisFilters(req))

	err := query.Count(&totalCount).Error
	if err != nil {
		ctx.Log.Error("Unable to count SIS objects", zap.Error(err))
		return 0, nil, err
	}

	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sisInfos).Error
	if err != nil {
		ctx.Log.Error("Unable to query the SIS objects", zap.Error(err))
		return totalCount, nil, err
	}

	return totalCount, sisInfos, nil
}

func (s *SIS) sisFilters(req *dtos.GetSISInfoReq) string {
	var filters []string

	if len(req.Id) > 0 {
		filters = append(filters, "id = '"+req.Id+"'")
	}

	if len(req.Code) > 0 {
		codeFilter := `(sis_booking ILIKE '` + req.Code + `%' OR shipment_code ILIKE '` + req.Code + `%' OR quote_code ILIKE '` + req.Code + `%')`
		filters = append(filters, codeFilter)
	}

	if len(filters) > 0 {
		return " AND " + strings.Join(filters, " AND ")
	}

	return ""
}

func (s *SIS) GetContent(ctx *context.Context, id string) (string, error) {

	var result string

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(s.getTable(ctx)).Select("contents").Where("id = ?", id).Scan(&result).Error
	if err != nil {
		ctx.Log.Error("unable to get content of sis", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (s *SIS) GetSISInfoRequest(ctx *context.Context, id string) (*dtos.SISInfo, error) {

	var model dtos.SISInfo

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(s.getSISTable(ctx)).Where("id =?", id).Scan(&model).Error
	if err != nil {
		ctx.Log.Error("unable to get sis info", zap.Error(err))
		return nil, err
	}

	return &model, nil
}

func (s *SIS) SaveSISInfo(ctx *context.Context, model *dtos.SISInfoModel, isISF, isSIS bool, data []byte) error {

	var SOIds []string
	for _, v := range model.Info {
		SOIds = append(SOIds, fmt.Sprintf("'%s'", v.ShipmentOrder))
	}

	subQuery := fmt.Sprintf("where sis_booking in (%s)", strings.Join(SOIds, ","))

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(s.getSISTable(ctx)).Exec(`update sis_info set sis_data = $1,isf_processed=$2,sis_processed=$3, sis_data_out_contents = $4 `+subQuery, model, isISF, isSIS, string(data)).Error
	if err != nil {
		ctx.Log.Error("Unable to save sis info", zap.Error(err))
		return err
	}

	return nil
}
