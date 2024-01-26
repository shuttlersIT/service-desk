// backend/services/agent_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AgentsServiceInterface provides methods for managing agents.
type AgentServiceInterface interface {
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
	GetRoles() ([]*models.Role, error)
	GetRoleByNumber(roleNumber int) (*models.Role, error)
	// Agent methods
	CreateAgent(agent *models.Agents) error
	UpdateAgent(agent *models.Agents) error
	DeleteAgent(agentID uint) error
	GetAllAgents() ([]*models.Agents, error)
	GetAgentByID(agentID uint) (*models.Agents, error)
	GetAgentByNumber(agentNumber int) (*models.Agents, error)

	// Role methods
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(roleID uint) error
	GetAllRoles() ([]*models.Role, error)
	GetRoleByID(roleID uint) (*models.Role, error)
	GetRoleByName(roleName string) (*models.Role, error)

	// Permission methods
	CreatePermission(permission *models.Permission) error
	UpdatePermission(permission *models.Permission) error
	DeletePermission(permissionID uint) error
	GetAllPermissions() ([]*models.Permission, error)
	GetPermissionByID(permissionID uint) (*models.Permission, error)
	GetPermissionByName(permissionName string) (*models.Permission, error)

	// Agent-Role methods
	AssignRoleToAgent(agentID uint, roleName string) error
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetAgentRoles(agentID uint) ([]*models.Role, error)

	// Agent-Permission methods
	AssignPermissionToAgent(agentID uint, permissionName string) error

	// User-Agent methods
	AddAgentToUser(userID uint, agentID uint) error
	GetAgentsByUser(userID uint) ([]*models.Agents, error)

	// Team-Agent methods
	AddAgentToTeam(teamID uint, agentID uint) error
	GetAgentsByTeam(teamID uint) ([]*models.Agents, error)

	// Agent-Permission methods
	GrantPermissionToAgent(agentID uint, permissionID uint) error
	RevokePermissionFromAgent(agentID uint, permissionID uint) error
	GetAgentPermissions(agentID uint) ([]*models.Permission, error)
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
func (ps *DefaultAgentService) GetUnitByID(unitID uint) (*models.Unit, error) {
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
func (ps *DefaultAgentService) DeleteUnit(unitID uint) error {
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
func (ps *DefaultAgentService) GetTeamByID(teamID uint) (*models.Teams, error) {
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
func (ps *DefaultAgentService) DeleteTeam(teamID uint) error {
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
func (ps *DefaultAgentService) GetRoleByID(roleID uint) (*models.Role, error) {
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
func (ps *DefaultAgentService) DeleteRole(roleID uint) error {
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

// GetAllRoles retrieves all roles from the database.
func (as *DefaultAgentService) GetAllRoles() ([]*models.Role, error) {
	return as.AgentDBModel.GetRoles()
}

// GetRoleByName retrieves a role by its name.
func (as *DefaultAgentService) GetRoleByName(roleName string) (*models.Role, error) {
	return as.AgentDBModel.GetRoleByName(roleName)
}

// CreatePermission creates a new permission.
func (as *DefaultAgentService) CreatePermission(permission *models.Permission) error {
	return as.AgentDBModel.CreatePermission(permission)
}

// UpdatePermission updates the details of an existing permission.
func (as *DefaultAgentService) UpdatePermission(permission *models.Permission) error {
	return as.AgentDBModel.UpdatePermission(permission)
}

// DeletePermission deletes a permission from the database.
func (as *DefaultAgentService) DeletePermission(permissionID uint) error {
	return as.AgentDBModel.DeletePermission(permissionID)
}

// GetAllPermissions retrieves all permissions from the database.
func (as *DefaultAgentService) GetAllPermissions() ([]*models.Permission, error) {
	return as.AgentDBModel.GetAllPermissions()
}

// GetPermissionByID retrieves a permission by its ID.
func (as *DefaultAgentService) GetPermissionByID(permissionID uint) (*models.Permission, error) {
	return as.AgentDBModel.GetPermissionByID(permissionID)
}

// GetPermissionByName retrieves a permission by its name.
func (as *DefaultAgentService) GetPermissionByName(permissionName string) (*models.Permission, error) {
	return as.AgentDBModel.GetPermissionByName(permissionName)
}

// AssignRoleToAgent assigns a role to an agent.
func (as *DefaultAgentService) AssignRoleToAgent(agentID uint, roleName string) error {
	r, err := as.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return err
	}
	return as.AgentDBModel.AssignRoleToAgent(agentID, r.ID)
}

// RevokeRoleFromAgent revokes a role from an agent.
func (as *DefaultAgentService) RevokeRoleFromAgent(agentID uint, roleName string) error {
	return as.AgentDBModel.RevokeRoleFromAgent(agentID, roleName)
}

// GetAgentRoles retrieves all roles assigned to an agent.
func (as *DefaultAgentService) GetAgentRoles(agentID uint) ([]*models.Role, error) {
	return as.AgentDBModel.GetAgentRoles(agentID)
}

// AssignPermissionToAgent assigns a permission to an agent.
func (as *DefaultAgentService) AssignPermissionToAgent(agentID uint, permissionName string) error {
	return as.AgentDBModel.AssignPermissionsToAgent(agentID, []string{permissionName})
}

// RevokePermissionFromAgent revokes a permission from an agent.
func (as *DefaultAgentService) RevokePermissionFromAgent(agentID uint, permissionName string) error {
	return as.AgentDBModel.RevokePermissionFromAgent(agentID, permissionName)
}

// GetAgentPermissions retrieves all permissions associated with an agent's roles.
func (as *DefaultAgentService) GetAgentPermissions(agentID uint) ([]*models.Permission, error) {
	return as.AgentDBModel.GetAgentPermissions(agentID)
}

// AddAgentToUser assigns an agent to a user.
func (as *DefaultAgentService) AddAgentToUser(userID uint, agentID uint) error {
	return as.AgentDBModel.CreateUserAgent(userID, agentID)
}

// GetAgentsByUser retrieves all agents assigned to a user.
func (as *DefaultAgentService) GetAgentsByUser(userID uint) ([]*models.Agents, error) {
	return as.AgentDBModel.GetAgentsByUser(userID)
}

// AddAgentToTeam assigns an agent to a team.
func (as *DefaultAgentService) AddAgentToTeam(teamID uint, agentID uint) error {
	return as.AgentDBModel.CreateTeamAgent(teamID, agentID)
}

// GetAgentsByTeam retrieves all agents assigned to a team.
func (as *DefaultAgentService) GetAgentsByTeam(teamID uint) ([]*models.Agents, error) {
	return as.AgentDBModel.GetAgentsByTeam(teamID)
}

// GrantPermissionToAgent grants a permission to an agent.
func (as *DefaultAgentService) GrantPermissionToAgent(agentID uint, permissionID uint) error {
	return as.AgentDBModel.GrantPermissionToAgent(agentID, permissionID)
}

// GetTeamsByAgent retrieves all teams assigned to an agent.
func (das *DefaultAgentService) GetTeamsByAgent(agentID uint) ([]*models.Teams, error) {
	return das.AgentDBModel.GetTeamsByAgent(agentID)
}

// AssignAgentToTeam assigns an agent to a team.
func (das *DefaultAgentService) AssignAgentToTeam(agentID uint, teamID uint) error {
	return das.AgentDBModel.CreateTeamAgent(teamID, agentID)
}

// RevokeAgentFromTeam revokes an agent from a team.
func (das *DefaultAgentService) RevokeAgentFromTeam(agentID uint, teamID uint) error {
	return das.AgentDBModel.DeleteTeamAgent(teamID, agentID)
}

// GetAgentTeams retrieves all teams assigned to an agent.
func (das *DefaultAgentService) GetAgentTeams(agentID uint) ([]*models.Teams, error) {
	return das.AgentDBModel.GetAgentTeams(agentID)
}

// AssignPermissionToTeam assigns a permission to a team.
func (das *DefaultAgentService) AssignPermissionToTeam(teamID uint, permissionName string) error {
	// First, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Then, grant the permission to the team
	return das.AgentDBModel.GrantPermissionToTeam(permission, teamID)
}

// RevokePermissionFromTeam revokes a permission from a team.
func (das *DefaultAgentService) RevokePermissionFromTeam(teamID uint, permissionName string) error {
	// First, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Then, revoke the permission from the team
	return das.AgentDBModel.RevokePermissionFromTeam(teamID, permission.ID)
}

// GetTeamPermissions retrieves all permissions associated with a team's roles.
func (das *DefaultAgentService) GetTeamPermissions(teamID uint) (*models.TeamPermission, error) {
	return das.AgentDBModel.GetTeamPermission(teamID)
}

// AddAgentPermissionToRole adds a permission to a role.
func (das *DefaultAgentService) AddAgentPermissionToRole(roleName string, permissionName string) error {
	// First, get the role by name
	role, err := das.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Then, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Finally, add the permission to the role
	return das.AgentDBModel.AssociatePermissionWithRole(role.ID, "", permission.ID)
}

// RemoveAgentPermissionFromRole removes a permission from a role.
func (das *DefaultAgentService) RemoveAgentPermissionFromRole(roleName string, permissionName string) error {
	// First, get the role by name
	role, err := das.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Then, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Finally, remove the permission from the role
	return das.AgentDBModel.RevokePermissionFromRole(role.RoleName, permission.Name)
}

// GetRolePermissions retrieves all permissions associated with a role.
func (das *DefaultAgentService) GetRolePermissions(roleName string) ([]*models.Permission, error) {
	// First, get the role by name
	role, err := das.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return nil, err
	}

	// Then, retrieve the permissions associated with the role
	return das.AgentDBModel.GetRolePermissions(role.ID)
}
