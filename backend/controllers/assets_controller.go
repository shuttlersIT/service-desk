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

// CreateAsset handles the HTTP request to create a new Asset.
func (pc *AssetController) CreateAsset(ctx *gin.Context) {
	var newAsset models.Assets
	if err := ctx.ShouldBindJSON(&newAsset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := pc.AssetService.CreateAsset(&newAsset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Assett"})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Assets created successfully"})
}

// GetAssetByID handles the HTTP request to retrieve a assets by ID.
func (pc *AssetController) GetAssetByID(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	asset, err := pc.AssetService.GetAssetByID(uint(assetID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

// UpdateAsset handles PUT /Asset/:id route.
func (pc *AssetController) UpdateAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ad models.Assets
	if err := ctx.ShouldBindJSON(&ad); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad.ID = uint(id)

	updatedAd, err := pc.AssetService.UpdateAsset(&ad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAd)
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

// GetAssetByID handles the HTTP request to retrieve a assets by ID.
func (pc *AssetController) GetAllAssets(ctx *gin.Context) {
	assets, err := pc.AssetService.GetAllAssets()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "assets not found"})
		return
	}
	ctx.JSON(http.StatusOK, assets)
}
