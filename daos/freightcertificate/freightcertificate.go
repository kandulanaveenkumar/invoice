package freightcertificate
 
import (
    "bitbucket.org/radarventures/forwarder-adapters/utils/context"
    "bitbucket.org/radarventures/forwarder-shipments/database/models"
    "go.uber.org/zap"
    "github.com/google/uuid"
)
 
type IFreightCertificate interface {
    Upsert(ctx *context.Context, m ...*models.FreightCertificate) error
    Get(ctx *context.Context, id uuid.UUID) (*models.FreightCertificate, error)
}
 
type FreightCertificate struct {
}
 
func NewFreightCertificate() IFreightCertificate {
    return &FreightCertificate{}
}
 
func (t *FreightCertificate) getTable(ctx *context.Context) string {
    return ctx.TenantID + "." + "freight_certificates"
}
func (t *FreightCertificate) Get(ctx *context.Context, shipmentId uuid.UUID) (*models.FreightCertificate,error) {
    
    result := &models.FreightCertificate{}
    err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "shipment_id = ?", shipmentId).Error
    if err != nil {
        ctx.Log.Error("Unable to get fc.", zap.Error(err))
        return nil, err
    }
 
    return result, err
}
 
func (fc *FreightCertificate) Upsert(ctx *context.Context, m ...*models.FreightCertificate) error {
    return ctx.DB.WithContext(ctx.Request.Context()).Table(fc.getTable(ctx)).Save(m).Error
}