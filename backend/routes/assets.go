// backend/routes/assets.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAssetsRoutes(r *gin.Engine, assets *controllers.AssetController) {

	a := r.Group("/assets")
	a.GET("/", assets.GetAllAssets)
	a.GET("/:id", assets.GetAssetByID)
	a.POST("/", assets.CreateAsset)
	a.PUT("/:id", assets.UpdateAsset)
	a.DELETE("/:id", assets.DeleteAsset)

}
