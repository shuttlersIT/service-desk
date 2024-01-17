package controllers

import (
	//"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
)

type AssetController struct {
	AssetDBModel *models.AssetDBModel
}

func NewAssetController() *AssetController {
	return &AssetController{}
}

// Implement controller methods like GetAssets, CreateAssett, GetAsset, UpdateAsset, DeleteAsset
