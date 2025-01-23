package blno

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"go.uber.org/zap"
)

type IBlNo interface {
	Get(ctx *context.Context) (int64, error)
}

type BlNo struct {
}

func NewBlNo() IBlNo {
	return &BlNo{}
}

func (t *BlNo) Get(ctx *context.Context) (int64, error) {

	bl_no := int64(0)

	err := ctx.DB.WithContext(ctx.Request.Context()).Raw("SELECT nextval('bl_no_seq')").Scan(&bl_no).Error
	if err != nil {
		ctx.Log.Error("Unable to get blno.", zap.Error(err))
	}

	return bl_no, err
}
