// backend/routes/assets.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupAssetsRoutes(router *gin.Engine, assetController *controllers.AssetController) {

	// Define Asset routes
	assetRoutes := router.Group("/assets")
	{
		assetRoutes.POST("/", assetController.CreateAsset)
		assetRoutes.PUT("/:id", assetController.UpdateAsset)
		assetRoutes.GET("/:id", assetController.GetAssetByID)
		assetRoutes.DELETE("/:id", assetController.DeleteAsset)
		assetRoutes.GET("/", assetController.GetAllAssets)
		assetRoutes.POST("/asset-assignments/assign/:id", assetController.AssignAssetToUser)
		assetRoutes.PUT("/asset-assignments/assign/unassign/:id", assetController.UnassignAssetFromUser)
		//assetRoutes.GET("/:id", assetController.GetAssetAssignmentByIDHandler)
		//assetRoutes.DELETE("/:id", assetController.DeleteAssetAssignmentHandler)
		//assetRoutes.GET("/", assetController.GetAllAssetAssignmentsHandler)
	}
}
