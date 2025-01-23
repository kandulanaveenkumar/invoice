package handlers

import (
	"net/http"
	"strconv"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/stock"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"go.uber.org/zap"
)

func CreateStockDetails(c *context.Context) {

	req := dtos.StocksReq{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("Unable to bing json", zap.Error(err))
	}
	req.RegionId = c.Account.RegionID
	res, err := stock.NewStockService().CreateNewStockDetails(c, &req)
	if err != nil {
		c.Log.Error("Error creating stock details", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.Log.Info("checking res", zap.Any("new stock details", res))
	c.JSON(http.StatusOK, res)

}

func CreateStockNumbers(c *context.Context) {

	req := dtos.StockNumberGenerateReq{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("Error binding json", zap.Error(err))
	}
	queryType := c.Query("type")
	req.Type = queryType
	airlineId := c.Params.ByName("airline_id")
	req.Airline = airlineId
	res, err := stock.NewStockService().CreateStocksNumbers(c, &req)
	if err != nil {
		c.Log.Error("Error creating new stock numbers", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetAirlinePortDetails(c *context.Context) {

	queryAirlineId := c.Query("airline_id")
	queryPortId := c.Query("port_id")
	req := dtos.AirLineDetailsReq{
		AirLine: queryAirlineId,
		PortId:  queryPortId,
	}
	res, err := stock.NewStockService().GetAirlinePortDetails(c, &req)
	if err != nil {
		c.Log.Error("Error fetching airline port details", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.Log.Info("checking res", zap.Any("arilines res", res))
	c.JSON(http.StatusOK, res)

}

func GetAvailableStock(c *context.Context) {

	queryAirlineId := c.Query("airline_id")
	queryPortId := c.Query("port_id")
	queryQ := c.Query("q")
	regionId := c.Query("region_id")
	req := dtos.AirLineDetailsReq{
		AirLine:  queryAirlineId,
		PortId:   queryPortId,
		Q:        queryQ,
		RegionId: regionId,
	}
	req.RegionId = c.Account.RegionID
	res, err := stock.NewStockService().GetAvailableStockNumber(c, &req)
	if err != nil {
		c.Log.Error("Error creating new stock numbers", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)

}

func GetStockNumberDetails(c *context.Context) {

	queryAirlineId := c.Query("airline_id")
	queryPortId := c.Query("port_id")
	queryStatus := c.Query("status")
	queryRegion := c.Query("region_id")

	req := dtos.AirLineDetailsReq{
		AirLine:  queryAirlineId,
		PortId:   queryPortId,
		Status:   queryStatus,
		RegionId: queryRegion,
	}
	req.RegionId = c.Account.RegionID
	res, err := stock.NewStockService().GetStockNumberDetails(c, &req)
	if err != nil {
		c.Log.Error("Error creating new stock numbers", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetAllAirLineStocks(c *context.Context) {

	page, _ := strconv.Atoi(c.Query("pg"))

	req := &dtos.GetReq{
		Q:        c.Query("q"),
		Pg:       int64(page),
		RegionId: c.Query("region_id"),
	}

	if req.RegionId == "" {
		req.RegionId = c.Account.RegionID
	}

	if c.Query("status") == "available" {

		res, err := stock.NewStockService().GetAllAirLineStocksAvailable(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.GetResponse(http.StatusInternalServerError, "", "Failed"),
			)
			return
		}
		c.JSON(http.StatusOK, res)

	} else if c.Query("status") == "completed" {

		res, err := stock.NewStockService().GetAllAirLineStocksCompleted(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.GetResponse(http.StatusInternalServerError, "", "Failed"),
			)
			return
		}
		c.JSON(http.StatusOK, res)

	} else {

		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", "proper parameter for status not status"),
		)
		return
	}
}

func UpdateStockStatus(c *context.Context) {

	req := dtos.AirLineDetailsReq{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("unable to bind json", zap.Error(err))
		return
	}
	req.RegionId = c.Account.RegionID
	res, err := stock.NewStockService().UpdateStockStatus(c, &req)
	if err != nil {
		c.Log.Error("Error creating new stock numbers", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", "Failed"),
		)
		return
	}
	c.JSON(http.StatusOK, res)

}
