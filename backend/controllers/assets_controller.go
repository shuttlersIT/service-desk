// backend/controllers/assets_controllers.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type AssetController struct {
	AssetService *services.DefaultAssetService
}

func NewAssetController(asset *services.DefaultAssetService) *AssetController {
	return &AssetController{
		AssetService: asset,
	}
}

// Implement controller methods like GetAssets, CreateAssett, GetAsset, UpdateAsset, DeleteAsset

func (ac *AssetController) CreateAsset(ctx *gin.Context) {
	var newAsset models.Assets
	if err := ctx.ShouldBindJSON(&newAsset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ac.AssetService.CreateAsset(&newAsset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Asset"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Asset created successfully"})
}

func (ac *AssetController) GetAssetByID(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	asset, err := ac.AssetService.GetAssetByID(uint(assetID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

// UpdateAsset handles PUT /Asset/:id route.
func (ac *AssetController) UpdateAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var asset models.Assets
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = uint(id)

	updatedAsset, err := ac.AssetService.UpdateAsset(&asset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAsset)
}

// DeleteAsset handles DELETE /assets/:id route.
func (pc *AssetController) DeleteAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	status, err := pc.AssetService.DeleteAsset(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, status)
}

func (ac *AssetController) GetAllAssets(ctx *gin.Context) {
	assets, err := ac.AssetService.GetAllAssets()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assets not found"})
		return
	}
	ctx.JSON(http.StatusOK, assets)
}
