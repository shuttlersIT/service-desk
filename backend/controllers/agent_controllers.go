// backend/controllers/agent_controllers.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type AgentController struct {
	AgentService *services.DefaultAgentService
}

func NewAgentController(agentService *services.DefaultAgentService) *AgentController {
	return &AgentController{
		AgentService: agentService,
	}
}

// Implement controller methods like GetAgents, CreateAgents, GetAgent, UpdateAgent, DeleteAgent

// CreateAgent handles the HTTP request to create a new Agent.
func (pc *AgentController) CreateAgent(ctx *gin.Context) {
	var newAgent models.Agents
	if err := ctx.ShouldBindJSON(&newAgent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := pc.AgentService.CreateAgent(&newAgent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Agent"})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Agents created successfully"})
}

// GetAgentByID handles the HTTP request to retrieve a agents by ID.
func (pc *AgentController) GetAgentByID(ctx *gin.Context) {
	agentID, _ := strconv.Atoi(ctx.Param("id"))
	agent, err := pc.AgentService.GetAgentByID(uint(agentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, agent)
}

// UpdateAgents handles PUT /Agents/:id route.
func (pc *AgentController) UpdateAgent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ad models.Agents
	if err := ctx.ShouldBindJSON(&ad); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad.ID = uint(id)

	updatedAd, err := pc.AgentService.UpdateAgent(&ad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAd)
}

// DeleteAgent handles DELETE /users/:id route.
func (pc *AgentController) DeleteAgent(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	status, err := pc.AgentService.DeleteAgent(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, status)
}

// GetAgentByID handles the HTTP request to retrieve a agents by ID.
func (pc *AgentController) GetAllAgents(ctx *gin.Context) {
	agents, err := pc.AgentService.GetAllAgents()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "agents not found"})
		return
	}
	ctx.JSON(http.StatusOK, agents)
}
