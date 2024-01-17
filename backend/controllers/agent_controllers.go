package controllers

import (

	// "github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
)

type AgentController struct {
	AgentDBModel *models.AgentDBModel
}

func NewAgentController() *AgentController {
	return &AgentController{}
}

// Implement controller methods like GetAgents, CreateAgents, GetAgent, UpdateAgent, DeleteAgent
