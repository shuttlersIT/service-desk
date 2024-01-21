// backend/services/agent_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AgentsServiceInterface provides methods for managing agents.
type AgentServiceInterface interface {
	CreateAgent(agent *models.Agents) error
	GetAgentByID(id uint) (*models.Agents, error)
	UpdateAgent(agent *models.Agents) (*models.Agents, error)
	DeleteAgent(agentID uint) (bool, error)
	GetAllAgents() []*models.Agents
}

// DefaultAgentService is the default implementation of AgentService
type DefaultAgentService struct {
	DB           *gorm.DB
	AgentDBModel *models.AgentDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAgentService creates a new DefaultAdvertisementService.
func NewDefaultAgentService(agentDBModel *models.AgentDBModel) *DefaultAgentService {
	return &DefaultAgentService{
		AgentDBModel: agentDBModel,
	}
}

// GetAllAgents retrieves all agents.
func (ps *DefaultAgentService) GetAllAgents() ([]*models.Agents, error) {
	agents, err := ps.AgentDBModel.GetAllAgents()
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// CreateAgent creates a new agent.
func (ps *DefaultAgentService) CreateAgent(agent *models.Agents) error {
	err := ps.AgentDBModel.CreateAgent(agent)
	if err != nil {
		return err
	}
	return nil
}

// CreateAgent creates a new agent.
func (ps *DefaultAgentService) GetAgentByID(id uint) (*models.Agents, error) {
	agent, err := ps.AgentDBModel.GetAgentByID(id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// UpdateAgent updates an existing agent.
func (ps *DefaultAgentService) UpdateAgent(agent *models.Agents) (*models.Agents, error) {
	err := ps.AgentDBModel.UpdateAgent(agent)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// DeleteAgent deletes an agent by ID.
func (ps *DefaultAgentService) DeleteAgent(agentID uint) (bool, error) {
	status := false
	err := ps.AgentDBModel.DeleteAgent(agentID)
	if err != nil {
		return status, err
	}
	status = true
	return status, nil
}
