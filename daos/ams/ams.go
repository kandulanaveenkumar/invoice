package ams

import (
	"encoding/json"
	"fmt"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/thirdparty/ams"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AMSDBI interface {
	AddAmsInfo(*context.Context, *models.AmsInfo) error
	UpdateAmsStatus(*context.Context, string, string) (string, error)
	GetAmsInfo(*context.Context, string) (*models.AmsInfo, error)
	UpdateFileGenerated(*context.Context, *ams.CargoSecurity, uuid.UUID, string) error
	GetBaseInfoFor(*context.Context, string) (*models.AmsInfo, error)
	UpdateErrorStatus(*context.Context, string, []dtos.AMSError) error
	GetAmsInfoWithId(*context.Context, string) (*models.AmsInfo, error)
	GetAllAddressDetails(*context.Context, string) ([]models.AmsDaoContactResp, error)
	AddAllAddressDetails(*context.Context, []models.AmsDaoContactResp) error
	UpdateAmsStatusManual(*context.Context, string, string, string) error
}

func NewAMSInfo() AMSDBI {

	return &AmsDB{}
}

type AmsDB struct {
}

// AddAmsInfo implements AMSDBI
func (a *AmsDB) AddAmsInfo(ctx *context.Context, info *models.AmsInfo) error {

	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Omit(clause.Associations).Table("ams_info").Create(&info).Error; err != nil {
			// return any error will rollback
			return err
		}

		vessels := info.AMSVesselInfo

		if err := tx.Omit(clause.Associations).Table("ams_vessel_info").Create(vessels).Error; err != nil {
			return err
		}

		routes := info.AMSRouteInfo
		if err := tx.Omit(clause.Associations).Table("ams_route_info").Create(routes).Error; err != nil {
			return err
		}
		containers := info.AMSContainerInfo
		if err := tx.Omit(clause.Associations).Table("ams_container_info").Create(containers).Error; err != nil {
			return err
		}
		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {

		return err
	}

	return err
}

// GetAmsInfo implements AMSDBI
func (a *AmsDB) GetAmsInfo(ctx *context.Context, hbl_no string) (*models.AmsInfo, error) {

	var resp models.AmsInfo

	var vessel []*models.AMSVesselInfo
	var route []*models.AMSRouteInfo
	var container []*models.AMSContainerInfo
	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit(clause.Associations).Table("ams_info").Where("hbl_no = ?", hbl_no).Order("created_at desc").Limit(1).Scan(&resp).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_vessel_info").Where("ams_info_id = ?", resp.ID).Scan(&vessel).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_route_info").Where("ams_info_id = ?", resp.ID).Scan(&route).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_container_info").Where("ams_info_id = ?", resp.ID).Scan(&container).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil

	})

	if err != nil {
		ctx.Log.Error(fmt.Sprintf(`unable to fetch amsinfo for hbl_no = %s`, hbl_no), zap.Error(err))
		return nil, err

	}

	resp.AMSContainerInfo = container
	resp.AMSRouteInfo = route
	resp.AMSVesselInfo = vessel

	return &resp, err

}

func (a *AmsDB) GetAmsInfoWithId(ctx *context.Context, id string) (*models.AmsInfo, error) {

	var resp models.AmsInfo

	var vessel []*models.AMSVesselInfo
	var route []*models.AMSRouteInfo
	var container []*models.AMSContainerInfo
	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit(clause.Associations).Table("ams_info").Where("id = ?", id).Order("created_at desc").Limit(1).Scan(&resp).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_vessel_info").Where("ams_info_id = ?", resp.ID).Scan(&vessel).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_route_info").Where("ams_info_id = ?", resp.ID).Scan(&route).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_container_info").Where("ams_info_id = ?", resp.ID).Scan(&container).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil

	})

	if err != nil {
		ctx.Log.Error(fmt.Sprintf(`unable to fetch amsinfo for id = %s`, id), zap.Error(err))
		return nil, err

	}

	resp.AMSContainerInfo = container
	resp.AMSRouteInfo = route
	resp.AMSVesselInfo = vessel

	return &resp, err

}

// UpdateAmsStatus implements AMSDBI
func (a *AmsDB) UpdateAmsStatus(ctx *context.Context, id, msg string) (string, error) {
	ctx.Log.Debug("Update ams status " + id + " " + msg)

	var ams_string string
	err := ctx.DB.Debug().Raw(`select ams_info_id from ams_generated where ams_file like $1 order by created_at desc limit 1`, id+"%").Scan(&ams_string).Error

	if err != nil {
		// return any error will rollback
		ctx.Log.Error("AMS ID ERR :  ", zap.Error(err))
		return "", err
	}

	ams_info_id, err := uuid.Parse(ams_string)

	if ams_info_id == uuid.Nil || err != nil {
		return "", nil
	}

	var bl string
	err = ctx.DB.Debug().Raw(`Update ams_info set ams_status = $1 where id = $2 returning hbl_no`, msg, ams_info_id).Scan(&bl).Error
	if err != nil {
		// return any error will rollback
		ctx.Log.Error("AMS Update ERR  :  ", zap.Error(err))

		return "", err
	}

	return ams_string, err

}

func (a *AmsDB) UpdateFileGenerated(ctx *context.Context, cargo *ams.CargoSecurity, amsID uuid.UUID, file string) error {
	data, err := cargo.Marshal()

	if err != nil {
		return err
	}

	val := models.AMSGenerated{
		ID:        uuid.New(),
		AmsInfoId: amsID.String(),
		AmsFile:   file,
		Cargo:     string(data),
	}
	err = ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Table("ams_info").Where("id = ?", amsID).Update("is_generated", true).Error; err != nil {
			// return any error will rollback
			return err
		}

		if err := tx.Omit(clause.Associations).Table("ams_generated").Create(&val).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil
	})

	return err
}

func (a *AmsDB) GetBaseInfoFor(ctx *context.Context, booking_id string) (*models.AmsInfo, error) {

	var resp models.AmsInfo
	var vessel []*models.AMSVesselInfo
	var route []*models.AMSRouteInfo

	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit(clause.Associations).Table("ams_info").Where("booking_id = ?", booking_id).Order("created_at desc").Limit(1).Scan(&resp).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_vessel_info").Where("ams_info_id = ?", resp.ID).Scan(&vessel).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Omit(clause.Associations).Table("ams_route_info").Where("ams_info_id = ?", resp.ID).Scan(&route).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil

	})

	if err != nil {
		ctx.Log.Error(fmt.Sprintf(`unable to fetch Base AMSInfo for booking = %s`, booking_id), zap.Error(err))
		return nil, err

	}

	resp.AMSRouteInfo = route
	resp.AMSVesselInfo = vessel

	return &resp, err

}

func (a *AmsDB) UpdateErrorStatus(ctx *context.Context, file string, errorObj []dtos.AMSError) error {
	val, err := json.Marshal(errorObj)
	if err != nil {
		ctx.Log.Error("Unable to update error", zap.Error(err))
		return err
	}

	var ams_info_id string
	query := `update ams_generated set error_response = $1::jsonb,updated_at = now() where ams_file = $2 returning ams_info_id`
	err = ctx.DB.Debug().Raw(query, val, file).Scan(&ams_info_id).Error

	if err != nil {
		ctx.Log.Error("Unable to update error", zap.Error(err))
		return err
	}

	return nil
}

func (a *AmsDB) GetAllAddressDetails(ctx *context.Context, bl_no string) ([]models.AmsDaoContactResp, error) {

	contacts := []models.AmsDaoContactResp{}
	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit(clause.Associations).Table("ams_contact_info").Where("bl_no = ?", bl_no).Scan(&contacts).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil

	})
	return contacts, err
}

func (a *AmsDB) AddAllAddressDetails(ctx *context.Context, contacts []models.AmsDaoContactResp) error {

	err := ctx.DB.Debug().Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Omit(clause.Associations).Table("ams_contact_info").Save(&contacts).Error; err != nil {
			// return any error will rollback
			return err
		}
		return nil
	})
	return err
}
func (a *AmsDB) UpdateAmsStatusManual(ctx *context.Context, bookingId string, hblNo string, msg string) error {
	ctx.Log.Debug("Started Update ams status bookingId : " + bookingId + " hblNo :" + hblNo + " msg: " + msg)

	err := ctx.DB.Debug().Exec(`UPDATE ams_info SET ams_status = $1 WHERE id = (SELECT id FROM ams_info WHERE shipment_id = $2 AND hbl_no = $3 ORDER BY created_at DESC LIMIT 1)`, msg, bookingId, hblNo).Error
	if err != nil {
		ctx.Log.Error("AMS Update ERR  :  ", zap.Error(err))
		return err
	}

	ctx.Log.Debug("Completed Update ams status bookingId : " + bookingId + " hblNo :" + hblNo + " msg: " + msg)

	return nil
}
