// backend/routes/assets.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAssetsRoutes(router *gin.Engine, assetController *controllers.AssetController) {

	// Define Asset routes
	assetRoutes := router.Group("/assets")
	{
		assetRoutes.POST("/", assetController.CreateAsset)
		assetRoutes.PUT("/:id", assetController.UpdateAsset)
		assetRoutes.GET("/:id", assetController.GetAssetByID)
		assetRoutes.DELETE("/:id", assetController.DeleteAsset)
		assetRoutes.GET("/", assetController.GetAllAssets)
	}
	// Define Asset Assignment routes
	assetAssignmentRoutes := router.Group("/asset-assignments")
	{
		assetAssignmentRoutes.POST("/", assetController.CreateAssetAssignAsset)
		assetAssignmentRoutes.PUT("/:id", assetController.UpdateAssetAssignmentHandler)
		assetAssignmentRoutes.GET("/:id", assetController.GetAssetAssignmentByIDHandler)
		assetAssignmentRoutes.DELETE("/:id", assetController.DeleteAssetAssignmentHandler)
		assetAssignmentRoutes.GET("/", assetController.GetAllAssetAssignmentsHandler)
	}
}
