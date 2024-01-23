// backend/services/agent_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AgentsServiceInterface provides methods for managing agents.
type AgentServiceInterface interface {
	CreateAgent(agent *models.Agents) error
	UpdateAgent(agent *models.Agents) error
	GetAgentByID(id uint) (*models.Agents, error)
	DeleteAgent(agentID uint) error
	GetAllAgents() ([]*models.Agents, error)
	CreateUnit(unit *models.Unit) error
	DeleteUnit(unitID int) error
	UpdateUnit(unit *models.Unit) error
	GetUnits() ([]*models.Unit, error)
	GetUnitByID(unitID int) (*models.Unit, error)
	GetUnitByNumber(unitNumber int) (*models.Unit, error)
	CreateTeam(team *models.Teams) error
	DeleteTeam(teamID int) error
	UpdateTeam(team *models.Teams) error
	GetTeams() ([]*models.Teams, error)
	GetTeamByID(teamID int) (*models.Teams, error)
	GetTeamByNumber(teamNumber int) (*models.Teams, error)
	CreateRole(role *models.Role) error
	DeleteRole(roleID int) error
	UpdateRole(role *models.Role) error
	GetRoles() ([]*models.Role, error)
	GetRoleByID(roleID int) (*models.Role, error)
	GetRoleByNumber(roleNumber int) (*models.Role, error)
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

// GetAgentByID retrieves an agent by its ID.
func (ps *DefaultAgentService) GetAgentByID(id uint) (*models.Agents, error) {
	agent, err := ps.AgentDBModel.GetAgentByID(id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// UpdateAgent updates an existing agent.
func (ps *DefaultAgentService) UpdateAgent(agent *models.Agents) error {
	err := ps.AgentDBModel.UpdateAgent(agent)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAgent deletes an agent by ID.
func (ps *DefaultAgentService) DeleteAgent(agentID uint) error {
	err := ps.AgentDBModel.DeleteAgent(agentID)
	if err != nil {
		return err
	}
	return nil
}

// backend/services/agent_service.go

// ...

// CreateUnit creates a new unit.
func (ps *DefaultAgentService) CreateUnit(unit *models.Unit) error {
	err := ps.AgentDBModel.CreateUnit(unit)
	if err != nil {
		return err
	}
	return nil
}

// GetUnitByID retrieves a unit by its ID.
func (ps *DefaultAgentService) GetUnitByID(unitID int) (*models.Unit, error) {
	unit, err := ps.AgentDBModel.GetUnitByID(unitID)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// UpdateUnit updates an existing unit.
func (ps *DefaultAgentService) UpdateUnit(unit *models.Unit) error {
	err := ps.AgentDBModel.UpdateUnit(unit)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUnit deletes a unit by ID.
func (ps *DefaultAgentService) DeleteUnit(unitID int) error {
	err := ps.AgentDBModel.DeleteUnit(unitID)
	if err != nil {
		return err
	}
	return nil
}

// GetUnits retrieves all units.
func (ps *DefaultAgentService) GetUnits() ([]*models.Unit, error) {
	units, err := ps.AgentDBModel.GetUnits()
	if err != nil {
		return nil, err
	}
	return units, nil
}

// CreateTeam creates a new team.
func (ps *DefaultAgentService) CreateTeam(team *models.Teams) error {
	err := ps.AgentDBModel.CreateTeam(team)
	if err != nil {
		return err
	}
	return nil
}

// GetTeamByID retrieves a team by its ID.
func (ps *DefaultAgentService) GetTeamByID(teamID int) (*models.Teams, error) {
	team, err := ps.AgentDBModel.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// UpdateTeam updates an existing team.
func (ps *DefaultAgentService) UpdateTeam(team *models.Teams) error {
	err := ps.AgentDBModel.UpdateTeam(team)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTeam deletes a team by ID.
func (ps *DefaultAgentService) DeleteTeam(teamID int) error {
	err := ps.AgentDBModel.DeleteTeam(teamID)
	if err != nil {
		return err
	}
	return nil
}

// GetTeams retrieves all teams.
func (ps *DefaultAgentService) GetTeams() ([]*models.Teams, error) {
	teams, err := ps.AgentDBModel.GetTeams()
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// CreateRole creates a new role.
func (ps *DefaultAgentService) CreateRole(role *models.Role) error {
	err := ps.AgentDBModel.CreateRole(role)
	if err != nil {
		return err
	}
	return nil
}

// GetRoleByID retrieves a role by its ID.
func (ps *DefaultAgentService) GetRoleByID(roleID int) (*models.Role, error) {
	role, err := ps.AgentDBModel.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRole updates an existing role.
func (ps *DefaultAgentService) UpdateRole(role *models.Role) error {
	err := ps.AgentDBModel.UpdateRole(role)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRole deletes a role by ID.
func (ps *DefaultAgentService) DeleteRole(roleID int) error {
	err := ps.AgentDBModel.DeleteRole(roleID)
	if err != nil {
		return err
	}
	return nil
}

// GetRoles retrieves all roles.
func (ps *DefaultAgentService) GetRoles() ([]*models.Role, error) {
	roles, err := ps.AgentDBModel.GetRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// backend/services/agent_service.go

// ...

// GetAgentByNumber retrieves an agent by their agent number.
func (ps *DefaultAgentService) GetAgentByNumber(agentNumber int) (*models.Agents, error) {
	agent, err := ps.AgentDBModel.GetAgentByNumber(agentNumber)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// GetUnitByNumber retrieves a unit by its unit number.
func (ps *DefaultAgentService) GetUnitByNumber(unitNumber int) (*models.Unit, error) {
	unit, err := ps.AgentDBModel.GetUnitByNumber(unitNumber)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// GetTeamByNumber retrieves a team by its team number.
func (ps *DefaultAgentService) GetTeamByNumber(teamNumber int) (*models.Teams, error) {
	team, err := ps.AgentDBModel.GetTeamByNumber(teamNumber)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// GetRoleByNumber retrieves a role by its role number.
func (ps *DefaultAgentService) GetRoleByNumber(roleNumber int) (*models.Role, error) {
	role, err := ps.AgentDBModel.GetRoleByNumber(roleNumber)
	if err != nil {
		return nil, err
	}
	return role, nil
}
