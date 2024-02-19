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

// GetAssetByID handles GET /assets/:id
func (ac *AssetController) GetAssetByID(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	asset, err := ac.AssetService.GetAssetByID(uint(assetID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

// GetAllAssets handles GET /assets
func (ac *AssetController) GetAllAssets(ctx *gin.Context) {
	assets, err := ac.AssetService.GetAllAssets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assets"})
		return
	}
	ctx.JSON(http.StatusOK, assets)
}

// GetAssetByNumber handles GET /assets/byNumber/:number
func (ac *AssetController) GetAssetByNumber(ctx *gin.Context) {
	assetNumber, _ := strconv.Atoi(ctx.Param("number"))
	asset, err := ac.AssetService.GetAssetByNumber(assetNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

// CreateAsset handles POST /assets
func (ac *AssetController) CreateAsset(ctx *gin.Context) {
	var asset models.Assets
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := ac.AssetService.CreateAsset(&asset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
		return
	}

	ctx.JSON(http.StatusCreated, asset)
}

// UpdateAsset handles PUT /assets/:id
func (ac *AssetController) UpdateAsset(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	var updatedAsset models.Assets
	if err := ctx.ShouldBindJSON(&updatedAsset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	currentAsset, err := ac.AssetService.GetAssetByID(uint(assetID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// Update the asset fields as needed
	currentAsset.AssetName = updatedAsset.AssetName
	currentAsset.Description = updatedAsset.Description
	// Update other fields accordingly

	if _, err := ac.AssetService.UpdateAsset(uint(assetID), currentAsset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset"})
		return
	}

	ctx.JSON(http.StatusOK, currentAsset)
}

// DeleteAsset handles DELETE /assets/:id
func (ac *AssetController) DeleteAsset(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))

	if err := ac.AssetService.DeleteAsset(uint(assetID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}

// Implement controller methods like GetAssets, CreateAssett, GetAsset, UpdateAsset, DeleteAsset

func (ac *AssetController) CreateAsset2(ctx *gin.Context) {
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

func (ac *AssetController) GetAssetByID2(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	asset, err := ac.AssetService.GetAssetByID(uint(assetID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

// UpdateAsset handles PUT /Asset/:id route.
func (ac *AssetController) UpdateAsset2(ctx *gin.Context) {
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

	updatedAsset, err := ac.AssetService.UpdateAsset(asset.ID, &asset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAsset)
}

// DeleteAsset handles DELETE /assets/:id route.
func (pc *AssetController) DeleteAsset2(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	erro := pc.AssetService.DeleteAsset(uint(id))
	if erro != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "asset deleted successfully"})
}

func (ac *AssetController) GetAllAssets2(ctx *gin.Context) {
	assets, err := ac.AssetService.GetAllAssets()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Assets not found"})
		return
	}
	ctx.JSON(http.StatusOK, assets)
}

// AssignAsset handles POST /assets/assign/:id/:userId/:agentId
func (ac *AssetController) AssignAsset(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	userID, _ := strconv.Atoi(ctx.Param("userId"))
	agentID, _ := strconv.Atoi(ctx.Param("agentId"))
	//session := sessions.Default(ctx)
	//agentID, ok := session.Get("agentID").(uint)
	//if !ok {
	//    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
	//    return
	//}

	if err := ac.AssetService.AssignAssetToUser(uint(assetID), uint(userID), uint(agentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign asset"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset assigned successfully"})
}

// UnassignAsset handles POST /assets/unassign/:id/:agentId
func (ac *AssetController) UnassignAsset(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	agentID, _ := strconv.Atoi(ctx.Param("agentId"))
	//session := sessions.Default(ctx)
	//agentID, ok := session.Get("agentID").(uint)
	//if !ok {
	//    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
	//    return
	//}

	if err := ac.AssetService.UnassignAsset(uint(assetID), uint(agentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign asset"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset unassigned successfully"})
}

// AssignAssetToUser handles POST /assets/assignToUser/:id/:userId/:agentId
func (ac *AssetController) AssignAssetToUser(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	userID, _ := strconv.Atoi(ctx.Param("userId"))
	agentID, _ := strconv.Atoi(ctx.Param("agentId"))
	//session := sessions.Default(ctx)
	//agentID, ok := session.Get("agentID").(uint)
	//if !ok {
	//    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
	//    return
	//}

	err := ac.AssetService.AssignAssetToUser(uint(assetID), uint(userID), uint(agentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign asset to user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset assigned successfully"})
}

// UnassignAssetFromUser handles POST /assets/unassignFromUser/:id/:agentId
func (ac *AssetController) UnassignAssetFromUser(ctx *gin.Context) {
	assetID, _ := strconv.Atoi(ctx.Param("id"))
	agentID, _ := strconv.Atoi(ctx.Param("agentId"))
	userID, _ := strconv.Atoi(ctx.Param("userId"))
	//session := sessions.Default(ctx)
	//agentID, ok := session.Get("agentID").(uint)
	//if !ok {
	//    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
	//    return
	//}

	if err := ac.AssetService.UnassignAssetFromUser(uint(assetID), uint(userID), uint(agentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign asset from user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset unassigned from user successfully"})
}

/*
// AssignAsset handles POST /assets/assign/:id/:userId
func (ac *AssetController) AssignAsset(ctx *gin.Context) {
    assetID, _ := strconv.Atoi(ctx.Param("id"))
    userID, _ := strconv.Atoi(ctx.Param("userId"))

    // Assuming you have a session middleware that stores the agent ID in the session
    session := sessions.Default(ctx)
    agentID, ok := session.Get("agentID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
        return
    }

    if err := ac.AssetService.AssignAsset(uint(assetID), uint(userID), agentID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign asset"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Asset assigned successfully"})
}

// UnassignAsset handles POST /assets/unassign/:id
func (ac *AssetController) UnassignAsset(ctx *gin.Context) {
    assetID, _ := strconv.Atoi(ctx.Param("id"))

    // Assuming you have a session middleware that stores the agent ID in the session
    session := sessions.Default(ctx)
    agentID, ok := session.Get("agentID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
        return
    }

    if err := ac.AssetService.UnassignAsset(uint(assetID), agentID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign asset"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Asset unassigned successfully"})
}

// AssignAssetToUser handles POST /assets/assignToUser/:id/:userId
func (ac *AssetController) AssignAssetToUser(ctx *gin.Context) {
    assetID, _ := strconv.Atoi(ctx.Param("id"))
    userID, _ := strconv.Atoi(ctx.Param("userId"))

    // Assuming you have a session middleware that stores the agent ID in the session
    session := sessions.Default(ctx)
    agentID, ok := session.Get("agentID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
        return
    }

    assetAssignment, err := ac.AssetService.AssignAssetToUser(uint(assetID), uint(userID), agentID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign asset to user"})
        return
    }

    ctx.JSON(http.StatusOK, assetAssignment)
}

// UnassignAssetFromUser handles POST /assets/unassignFromUser/:id
func (ac *AssetController) UnassignAssetFromUser(ctx *gin.Context) {
    assetID, _ := strconv.Atoi(ctx.Param("id"))

    // Assuming you have a session middleware that stores the agent ID in the session
    session := sessions.Default(ctx)
    agentID, ok := session.Get("agentID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
        return
    }

    if err := ac.AssetService.UnassignAssetFromUser(uint(assetID), agentID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign asset from user"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Asset unassigned from user successfully"})
}


// AssignAsset handles POST /assets/:id/assign
func (ac *AssetController) AssignAsset(ctx *gin.Context) {
	// Assuming you have a session middleware that stores the agent ID in the session
    session := sessions.Default(ctx)
    userID, ok := session.Get("userID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    uID, _ := strconv.Atoi(ctx.Param("id"))
	userID := uint(uID)
    //userID := // Get the user ID from the current session (you may use your session management library)

    // Check if the user has permission to assign assets (optional, based on your requirements)

    err := ac.AssetService.AssignAsset(uint(assetID), userID, userID) // Assuming the agent is the user assigning
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign asset"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Asset assigned successfully"})
}

// UnassignAsset handles POST /assets/:id/unassign
func (ac *AssetController) UnassignAsset(ctx *gin.Context) {
    assetID, _ := strconv.Atoi(ctx.Param("id"))
	session := sessions.Default(ctx)
    userID, ok := session.Get("userID").(uint)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    uID, _ := strconv.Atoi(ctx.Param("id"))
	userID := uint(uID)
    //userID := // Get the user ID from the current session (you may use your session management library)

    // Check if the user has permission to unassign assets (optional, based on your requirements)

    err := ac.AssetService.UnassignAsset(uint(assetID), userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign asset"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Asset unassigned successfully"})
}

*/
