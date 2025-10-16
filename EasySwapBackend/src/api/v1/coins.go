package v1

import (
	"github.com/ProjectsTask/EasySwapBackend/src/service/svc"
	"github.com/ProjectsTask/EasySwapBackend/src/service/v1"
	"github.com/ProjectsTask/EasySwapBackend/src/types/v1"
	"github.com/ProjectsTask/EasySwapBase/errcode"
	"github.com/ProjectsTask/EasySwapBase/xhttp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchCoinsHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get query param "keyword"
		keyword := c.Query("keyword")

		res, err := service.SearchCoins(c, svcCtx, keyword)
		if err != nil {
			xhttp.Error(c, err)
			return
		}

		err = svcCtx.Dao.AddCoinsMetadata(res)
		if err != nil {
			xhttp.Error(c, err)
		}
		xhttp.OkJson(c, res)
	}
}

func ListUserCoinFavoritesHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddr := c.Params.ByName("address")
		if userAddr == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}

		userFavoriteCoins, err := svcCtx.Dao.ListUserFavoriteCoinMetadata(userAddr)
		if err != nil {
			xhttp.Error(c, errcode.NewCustomErr("list user favorite coins err:"+err.Error()))
			return
		}

		var coinIds []string
		for _, coin := range userFavoriteCoins {
			coinIds = append(coinIds, coin.Id)
		}

		resultInMap, err := service.GetCoinsPrice(c, svcCtx, coinIds, "usd")
		if err != nil {
			xhttp.Error(c, errcode.ErrUnexpected)
			return
		}

		var result []types.CoinWithLatestPrice
		for _, coin := range userFavoriteCoins {
			result = append(result, types.CoinWithLatestPrice{
				Id:       coin.Id,
				Symbol:   coin.Symbol,
				Name:     coin.Name,
				Price:    resultInMap[coin.Id].Usd,
				Currency: "usd",
			})
		}

		xhttp.OkJson(c, result)

	}
}

func AddUserCoinFavoriteHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddr := c.Params.ByName("address")
		if userAddr == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}

		req := types.AddCoinFavoriteRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			xhttp.Error(c, errcode.NewCustomErr("invalid request:"+err.Error()))
			return
		}
		if req.Id == "" {
			xhttp.Error(c, errcode.NewCustomErr("coin id is null"))
			return
		}

		err := svcCtx.Dao.AddFavoriteCoin(userAddr, req.Id)
		if err != nil {
			xhttp.Error(c, errcode.NewCustomErr("get coin metadata err:"+err.Error()))
			return
		}

		xhttp.OkJson(c, gin.H{
			"message": "success",
		})
	}
}

func RemoveUserCoinFavoriteHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddr := c.Params.ByName("address")
		if userAddr == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}

		coinId := c.Params.ByName("coin_id")
		if coinId == "" {
			xhttp.Error(c, errcode.NewCustomErr("coin id is null"))
			return
		}

		err := svcCtx.Dao.RemoveFavoriteCoin(userAddr, coinId)
		if err != nil {
			xhttp.Error(c, errcode.NewCustomErr("remove coin favorite err:"+err.Error()))
			return
		}

		xhttp.OkJson(c, gin.H{
			"message": "success",
		})
	}
}

func ListUserCoinAlertsHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")

		alerts, err := svcCtx.Dao.ListUserCoinAlerts(address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch alerts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"alerts": alerts})
	}
}

type AddAlertRequest struct {
	CoinID      string  `json:"coin_id" binding:"required"`
	AlertType   string  `json:"alert_type" binding:"required,oneof=above below"`
	TargetPrice float64 `json:"target_price" binding:"required"`
}

func AddUserCoinAlertHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		var req AddAlertRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// todo check input params

		err := svcCtx.Dao.AddCoinPriceAlert(address, req.CoinID, req.AlertType, req.TargetPrice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save alert"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "alert saved successfully"})
	}
}

func DeleteUserCoinAlertHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		coinID := c.Param("coin_id")
		alertType := c.Param("alert_type")

		// todo check input params
		err := svcCtx.Dao.RemoveCoinPriceAlert(address, coinID, alertType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete alert"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "alert deleted"})
	}
}

func GetEtherBalanceHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		if address == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}

		balance, err := service.GetAccountBalance(c.Request.Context(), svcCtx, address, 1)
		if err != nil {
			xhttp.Error(c, errcode.NewCustomErr("get ether balance err:"+err.Error()))
			return
		}

		xhttp.OkJson(c, gin.H{
			"address": address,
			"balance": balance.Result,
		})
	}
}
