package handlers
import (
	"net/http"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)
func GetPartnerInfoDashboard(c *context.Context) {
	Id := c.Params.ByName("cid")
	companyId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().GetPartnerDashboard(c, companyId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetCustomerInfoDashboard(c *context.Context) {
	Id := c.Params.ByName("cid")
	companyId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().GetCustomerProfileInfoDashboard(c, companyId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}