// backend/services/agent_service.go

package services

import (
	"fmt"

	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AgentsServiceInterface provides methods for managing agents.
type AgentServiceInterface interface {
	// Unit Management
	CreateUnit(unit *models.Unit) error
	UpdateUnit(unit *models.Unit) error
	DeleteUnit(unitID uint) error
	GetUnitByID(unitID uint) (*models.Unit, error)
	GetUnits() ([]*models.Unit, error)
	GetUnitByNumber(unitNumber int) (*models.Unit, error)

	// Team Management
	CreateTeam(team *models.Teams) error
	UpdateTeam(team *models.Teams) error
	DeleteTeam(teamID uint) error
	GetTeamByID(teamID uint) (*models.Teams, error)
	GetTeams() ([]*models.Teams, error)
	GetTeamByNumber(teamNumber int) (*models.Teams, error)

	// Role Management
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(roleID uint) error
	GetRoleByID(roleID uint) (*models.Role, error)
	GetRoles() ([]*models.Role, error)
	GetRoleByNumber(roleNumber int) (*models.Role, error)
	GetRoleByName(roleName string) (*models.Role, error)

	// Agent Management
	CreateAgent(agent *models.Agents) error
	UpdateAgent(agent *models.Agents) error
	DeleteAgent(agentID uint) error
	GetAgentByID(agentID uint) (*models.Agents, error)
	GetAllAgents() ([]*models.Agents, error)
	GetAgentByNumber(agentNumber int) (*models.Agents, error)

	// Permission Management
	CreatePermission(permission *models.Permission) error
	UpdatePermission(permission *models.Permission) error
	DeletePermission(permissionID uint) error
	GetPermissionByID(permissionID uint) (*models.Permission, error)
	GetAllPermissions() ([]*models.Permission, error)
	GetPermissionByName(permissionName string) (*models.Permission, error)

	// Agent-Role Assignments
	AssignRoleToAgent(agentID uint, roleName string) error
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetAgentRoles(agentID uint) ([]*models.Role, error)

	// Role-Permission Management
	AddAgentPermissionToRole(roleName string, permissionName string) error
	RemoveAgentPermissionFromRole(roleName string, permissionName string) error
	GetRolePermissions(roleName string) ([]*models.Permission, error)

	// Agent-Permission Assignments
	AssignPermissionsToAgent(agentID uint, permissionNames []string) error
	RevokePermissionFromAgent(agentID uint, permissionID uint) error
	GetAgentPermissions(agentID uint) ([]*models.Permission, error)

	// Team-Agent Assignments
	AddAgentToTeam(teamID uint, agentID uint) error
	GetAgentsByTeam(teamID uint) ([]*models.Agents, error)
	RevokeAgentFromTeam(agentID uint, teamID uint) error

	// Team-Permission Assignments
	AssignPermissionToTeam(teamID uint, permissionName string) error
	RevokePermissionFromTeam(teamID uint, permissionName string) error
	GetTeamPermissions(teamID uint) (*models.TeamPermission, error)

	// User-Agent Assignments
	AddAgentToUser(userID uint, agentID uint) error
	GetAgentsByUser(userID uint) ([]*models.Agents, error)
	GetAllRoles() ([]*models.Role, error)

	// Agent-Permission methods
	AssignPermissionToAgent(agentID uint, permissionName string) error

	// Agent-Permission methods
	GrantPermissionToAgent(agentID uint, permissionID uint) error
}

// Logger interface for logging
type Logger interface {
	Info(message string)
	Error(message string)
}

// EventPublisher interface for publishing events
type EventPublisher interface {
	Publish(event interface{}) error
}

// AgentServiceInterface - unchanged for brevity

type DefaultAgentService struct {
	DB             *gorm.DB
	AgentDBModel   *models.AgentDBModel
	Logger         Logger
	EventPublisher EventPublisher
}

func NewDefaultAgentService(db *gorm.DB, logger Logger, eventPublisher EventPublisher, agentDBModel *models.AgentDBModel) *DefaultAgentService {
	return &DefaultAgentService{
		DB:             db,
		AgentDBModel:   agentDBModel,
		Logger:         logger,
		EventPublisher: eventPublisher,
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

// CreateAgent creates a new agent with detailed logging and event publishing.
func (s *DefaultAgentService) CreateAgent(agent *models.Agents) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.CreateAgent(agent); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to create agent:", "error", err))
			return err
		}

		// Asynchronously publish a creation event
		go func(a *models.Agents) {
			event := struct {
				AgentID uint
				Action  string
			}{
				AgentID: a.ID,
				Action:  "created",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish agent creation event:", "error", err))
			}
		}(agent)
		return nil
	})
}

func (s *DefaultAgentService) CreateAgent2(agent *models.Agents) (*models.Agents, error) {
	var createdAgent *models.Agents
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}
		createdAgent = agent

		// Asynchronously publish an event
		go func(a *models.Agents) {

			if err := s.EventPublisher.Publish(a); err != nil {
				s.Logger.Error("Failed to publish event")
			}
		}(agent)

		return nil
	})

	if err != nil {
		s.Logger.Error("Failed to create agent")
		return nil, err
	}

	return createdAgent, nil
}

// CreateAgent creates a new agent.
func (ps *DefaultAgentService) CreateAgent3(agent *models.Agents) (*models.Agents, error) {
	err := ps.AgentDBModel.CreateAgent(agent)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// GetAgentByID retrieves an agent by their ID with comprehensive error handling.
func (s *DefaultAgentService) GetAgentByID(agentID uint) (*models.Agents, error) {
	agent, err := s.AgentDBModel.GetAgentByID(agentID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to retrieve agent by ID:", "agentID", agentID, "error", err))
		return nil, err
	}
	return agent, nil
}

// UpdateAgent updates an existing agent with transaction management.
// UpdateAgent updates an existing agent with transaction management and event publishing.
func (s *DefaultAgentService) UpdateAgent(agent *models.Agents) error {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(agent).Error; err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to update agent:", err))
			return err // Returning error will rollback the transaction
		}

		// Construct and publish an update event
		event := struct {
			AgentID uint
			Update  string
		}{
			AgentID: agent.ID,
			Update:  "Agent updated successfully",
		}
		if err := s.EventPublisher.Publish(event); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to publish agent update event:", err))
			// Log the error but do not rollback the transaction for event publishing failure
		}

		return nil // Success
	})

	if err != nil {
		s.Logger.Error(fmt.Sprintf("Transaction failed for updating agent:", err))
		return err
	}

	return nil // Successfully updated the agent and potentially published an event
}

// UpdateAgent updates an existing agent.
func (ps *DefaultAgentService) UpdateAgent2(agent *models.Agents) error {
	err := ps.AgentDBModel.UpdateAgent(agent)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAgent deletes an agent by ID with transaction management.
// DeleteAgent deletes an agent by ID with transaction management and event publishing.
func (s *DefaultAgentService) DeleteAgent(agentID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Agents{}, agentID).Error; err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to delete agent:", err))
			return err
		}

		// Construct and asynchronously publish a delete event
		go func() {
			event := struct {
				AgentID uint
				Message string
			}{
				AgentID: agentID,
				Message: "Agent deleted successfully",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish agent delete event:", err))
			}
		}()

		return nil
	})
}

// DeleteAgent deletes an agent by ID.
func (ps *DefaultAgentService) DeleteAgent2(agentID uint) error {
	err := ps.AgentDBModel.DeleteAgent(agentID)
	if err != nil {
		return err
	}
	return nil
}

// backend/services/agent_service.go

// ...

// CreateUnit creates a new unit.
// CreateUnit creates a new unit with transaction management.
func (ps *DefaultAgentService) CreateUnit(unit *models.Unit) error {
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(unit).Error; err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to create unit:", err))
			return err
		}
		// Asynchronously publish unit creation event
		// Construct and publish an update event
		event := struct {
			UnitID uint
			Update string
		}{
			UnitID: unit.ID,
			Update: "Unit created successfully",
		}
		if err := ps.EventPublisher.Publish(event); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to publish create unit event:", err))
			// Log the error but do not rollback the transaction for event publishing failure
		}

		return nil // Success
	})

	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Transaction failed to create unit:", err))
		return err
	}

	return nil
}

func (ps *DefaultAgentService) CreateUnit2(unit *models.Unit) error {
	err := ps.AgentDBModel.CreateUnit(unit)
	if err != nil {
		return err
	}
	return nil
}

// GetUnitByID retrieves a unit by its ID with added logging for error scenarios.
func (ps *DefaultAgentService) GetUnitByID(unitID uint) (*models.Unit, error) {
	unit, err := ps.AgentDBModel.GetUnitByID(unitID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve unit by ID:", err))
		return nil, err
	}
	return unit, nil
}

func (ps *DefaultAgentService) GetUnitByID2(unitID uint) (*models.Unit, error) {
	unit, err := ps.AgentDBModel.GetUnitByID(unitID)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// UpdateUnit updates an existing unit with transaction management.
func (ps *DefaultAgentService) UpdateUnit(unit *models.Unit) error {
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(unit).Error; err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to update unit:", err))
			return err
		}
		// Asynchronously publish unit update event
		event := struct {
			UnitID uint
			Update string
		}{
			UnitID: unit.ID,
			Update: "Unit updated successfully",
		}
		if err := ps.EventPublisher.Publish(event); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to publish update unit event:", err))
			// Log the error but do not rollback the transaction for event publishing failure
		}
		return nil // Success
	})

	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Transaction failed to update unit:", err))
		return err
	}
	return nil
}

// UpdateUnit updates an existing unit.
func (ps *DefaultAgentService) UpdateUnit2(unit *models.Unit) error {
	err := ps.AgentDBModel.UpdateUnit(unit)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUnit deletes a unit by ID.
// DeleteUnit deletes a unit by ID with transaction management.
func (ps *DefaultAgentService) DeleteUnit(unitID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Unit{}, unitID).Error; err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to delete unit:", err))
			return err
		}
		// Asynchronously publish unit deletion event
		// Construct and asynchronously publish a delete event
		go func() {
			event := struct {
				UnitID  uint
				Message string
			}{
				UnitID:  unitID,
				Message: "Unit deleted successfully",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish agent delete event:", err))
			}
		}()

		return nil
	})
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
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to create team: %v", err))
			return err
		}

		// Asynchronously publish a team creation event
		go func() {
			event := struct {
				TeamID uint
				Action string
			}{
				TeamID: team.ID,
				Action: "team_created",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish team creation event: %v", err))
			}
		}()

		return nil
	})
}

func (ps *DefaultAgentService) CreateTeam2(team *models.Teams) error {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve team by ID %d: %v", teamID, err))
		return nil, err
	}
	return team, nil
}

func (ps *DefaultAgentService) GetTeamByID2(teamID uint) (*models.Teams, error) {
	team, err := ps.AgentDBModel.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// UpdateTeam updates an existing team.
func (ps *DefaultAgentService) UpdateTeam(team *models.Teams) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.UpdateTeam(team); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to update team %d: %v", team.ID, err))
			return err
		}

		go func() {
			event := struct {
				TeamID uint
				Action string
			}{
				TeamID: team.ID,
				Action: "team_updated",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish team update event: %v", err))
			}
		}()
		return nil
	})
}

// DeleteTeam deletes a team by ID.
func (ps *DefaultAgentService) DeleteTeam(teamID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.DeleteTeam(teamID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to delete team %d: %v", teamID, err))
			return err
		}

		go func() {
			event := struct {
				TeamID uint
				Action string
			}{
				TeamID: teamID,
				Action: "team_deleted",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish team deletion event: %v", err))
			}
		}()
		return nil
	})
}

func (ps *DefaultAgentService) DeleteTeam2(teamID uint) error {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve all teams: %v", err))
		return nil, err
	}
	return teams, nil
}

func (ps *DefaultAgentService) GetTeams2() ([]*models.Teams, error) {
	teams, err := ps.AgentDBModel.GetTeams()
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// CreateRole creates a new role with comprehensive error handling and event publishing.
func (s *DefaultAgentService) CreateRole(role *models.Role) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.CreateRole(role); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to create role", "roleName", role.RoleName, "error", err))
			return err
		}

		// Asynchronously publish a creation event
		go func() {
			event := struct {
				RoleID   uint
				RoleName string
				Action   string
			}{
				RoleID:   role.ID,
				RoleName: role.RoleName,
				Action:   "role_created",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish role creation event", "error", err))
			}
		}()
		return nil
	})
}

// CreateRole creates a new role.
func (ps *DefaultAgentService) CreateRole2(role *models.Role) error {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve role by ID %d: %v", roleID, err))
		return nil, err
	}
	return role, nil
}

func (ps *DefaultAgentService) GetRoleByID2(roleID uint) (*models.Role, error) {
	role, err := ps.AgentDBModel.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRole updates an existing role, encapsulating the operation within a transaction and publishing an event.
func (s *DefaultAgentService) UpdateRole(role *models.Role) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.UpdateRole(role); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to update role", "roleID", role.ID, "error", err))
			return err
		}

		// Asynchronously publish an update event
		go func() {
			event := struct {
				RoleID   uint
				RoleName string
				Action   string
			}{
				RoleID:   role.ID,
				RoleName: role.RoleName,
				Action:   "role_updated",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish role update event", "error", err))
			}
		}()
		return nil
	})
}

// UpdateRole updates an existing role.
func (ps *DefaultAgentService) UpdateRole2(role *models.Role) error {
	err := ps.AgentDBModel.UpdateRole(role)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRole deletes a role by ID.
func (ps *DefaultAgentService) DeleteRole(roleID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.DeleteRole(roleID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to delete role %d: %v", roleID, err))
			return err
		}

		go func() {
			event := struct {
				RoleID uint
				Action string
			}{
				RoleID: roleID,
				Action: "role_deleted",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish role deletion event: %v", err))
			}
		}()
		return nil
	})
}

// GetRoles retrieves all roles.
func (ps *DefaultAgentService) GetRoles() ([]*models.Role, error) {
	roles, err := ps.AgentDBModel.GetRoles()
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve all roles: %v", err))
		return nil, err
	}
	return roles, nil
}

func (ps *DefaultAgentService) GetRoles2() ([]*models.Role, error) {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve agent by number %d: %v", agentNumber, err))
		return nil, err
	}
	return agent, nil
}

func (ps *DefaultAgentService) GetAgentByNumber2(agentNumber int) (*models.Agents, error) {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve unit by number %d: %v", unitNumber, err))
		return nil, err
	}
	return unit, nil
}

func (ps *DefaultAgentService) GetUnitByNumber2(unitNumber int) (*models.Unit, error) {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve team by number %d: %v", teamNumber, err))
		return nil, err
	}
	return team, nil
}

func (ps *DefaultAgentService) GetTeamByNumber2(teamNumber int) (*models.Teams, error) {
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
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve role by number %d: %v", roleNumber, err))
		return nil, err
	}
	return role, nil
}

func (ps *DefaultAgentService) GetRoleByNumber2(roleNumber int) (*models.Role, error) {
	role, err := ps.AgentDBModel.GetRoleByNumber(roleNumber)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetAllRoles retrieves all roles from the database.
func (ps *DefaultAgentService) GetAllRoles() ([]*models.Role, error) {
	roles, err := ps.AgentDBModel.GetRoles()
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve all roles:", err))
		return nil, err
	}
	return roles, nil
}

func (as *DefaultAgentService) GetAllRoles2() ([]*models.Role, error) {
	return as.AgentDBModel.GetRoles()
}

// GetRoleByName retrieves a role by its name.
func (ps *DefaultAgentService) GetRoleByName(roleName string) (*models.Role, error) {
	role, err := ps.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve role by name '%s': %v", roleName, err))
		return nil, err
	}
	return role, nil
}

func (as *DefaultAgentService) GetRoleByName2(roleName string) (*models.Role, error) {
	return as.AgentDBModel.GetRoleByName(roleName)
}

// CreatePermission creates a new permission.
// CreatePermission creates a new permission with transaction management and logs the action.
func (s *DefaultAgentService) CreatePermission(permission *models.Permission) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(permission).Error; err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to create permission", "permissionName", permission.Name, "error", err))
			return err
		}

		// Publish creation event
		go func() {
			event := struct {
				PermissionID uint
				Action       string
			}{
				PermissionID: permission.ID,
				Action:       "permission_created",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish permission creation event", "error", err))
			}
		}()
		return nil
	})
}

func (as *DefaultAgentService) CreatePermission2(permission *models.Permission) error {
	return as.AgentDBModel.CreatePermission(permission)
}

// UpdatePermission updates an existing permission with transaction management, error handling, and event publishing.
func (s *DefaultAgentService) UpdatePermission(permission *models.Permission) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(permission).Error; err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to update permission", "permissionID", permission.ID, "error", err))
			return err
		}

		// Publish update event
		go func() {
			event := struct {
				PermissionID uint
				Action       string
			}{
				PermissionID: permission.ID,
				Action:       "permission_updated",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish permission update event", "error", err))
			}
		}()
		return nil
	})
}

func (ps *DefaultAgentService) UpdatePermission2(permission *models.Permission) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.UpdatePermission(permission); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to update permission '%s': %v", permission.Name, err))
			return err
		}
		return nil
	})
}

// UpdatePermission updates the details of an existing permission.
func (as *DefaultAgentService) UpdatePermission3(permission *models.Permission) error {
	return as.AgentDBModel.UpdatePermission(permission)
}

// DeletePermission deletes a permission from the database.
func (ps *DefaultAgentService) DeletePermission(permissionID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.DeletePermission(permissionID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to delete permission ID %d: %v", permissionID, err))
			return err
		}

		go func() {
			event := struct {
				PermissionID uint
				Action       string
			}{
				PermissionID: permissionID,
				Action:       "permission_deleted",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission deletion event:", err))
			}
		}()
		return nil
	})
}

func (as *DefaultAgentService) DeletePermission2(permissionID uint) error {
	return as.AgentDBModel.DeletePermission(permissionID)
}

// GetAllPermissions retrieves all permissions from the database.
func (ps *DefaultAgentService) GetAllPermissions() ([]*models.Permission, error) {
	permissions, err := ps.AgentDBModel.GetAllPermissions()
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve all permissions:", err))
		return nil, err
	}
	return permissions, nil
}

func (as *DefaultAgentService) GetAllPermissions2() ([]*models.Permission, error) {
	return as.AgentDBModel.GetAllPermissions()
}

// GetPermissionByID retrieves a permission by its ID.
func (ps *DefaultAgentService) GetPermissionByID(permissionID uint) (*models.Permission, error) {
	permission, err := ps.AgentDBModel.GetPermissionByID(permissionID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve permission by ID %d: %v", permissionID, err))
		return nil, err
	}
	return permission, nil
}

func (as *DefaultAgentService) GetPermissionByID2(permissionID uint) (*models.Permission, error) {
	return as.AgentDBModel.GetPermissionByID(permissionID)
}

// GetPermissionByName retrieves a permission by its name.
func (ps *DefaultAgentService) GetPermissionByName(permissionName string) (*models.Permission, error) {
	permission, err := ps.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve permission by name '%s': %v", permissionName, err))
		return nil, err
	}
	return permission, nil
}

func (as *DefaultAgentService) GetPermissionByName2(permissionName string) (*models.Permission, error) {
	return as.AgentDBModel.GetPermissionByName(permissionName)
}

func (ps *DefaultAgentService) AssignRoleToAgent(agentID uint, roleNames []string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AssignRolesToAgent(agentID, roleNames); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to assign roles to agent %d: %v", agentID, err))
			return err
		}

		// Asynchronously publish an event for each role assignment
		for _, roleName := range roleNames {
			go func(rName string) {
				event := struct {
					AgentID  uint
					RoleName string
					Action   string
				}{
					AgentID:  agentID,
					RoleName: rName,
					Action:   "role_assigned",
				}
				if err := ps.EventPublisher.Publish(event); err != nil {
					ps.Logger.Error(fmt.Sprintf("Failed to publish role assignment event: %v", err))
				}
			}(roleName)
		}

		return nil
	})
}

// AssignRolesToAgent assigns multiple roles to an agent, ensuring transactional integrity and event publishing.
func (s *DefaultAgentService) AssignRolesToAgent(agentID uint, roles []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.AssignRolesToAgentByRoleID(agentID, roles); err != nil {
			s.Logger.Error(fmt.Sprintf("Error assigning roles to agent", "agentID", agentID, "error", err))
			return err
		}

		// Asynchronously publish event for each role assignment
		for _, roleID := range roles {
			go func(rID uint) {
				event := struct {
					AgentID uint
					RoleID  uint
				}{AgentID: agentID, RoleID: rID}
				if err := s.EventPublisher.Publish(event); err != nil {
					s.Logger.Error(fmt.Sprintf("Failed to publish role assignment event", "error", err))
				}
			}(roleID)
		}
		return nil
	})
}

// AssignRolesToAgent accepts either role names or role IDs to assign roles to an agent.
func (s *DefaultAgentService) AssignRolesToAgent2(agentID uint, roleIdentifiers []interface{}) error {
	// Call the adapted AgentDBModel method with the role identifiers
	err := s.AgentDBModel.AssignRolesToAgentByIDOrName(agentID, roleIdentifiers)
	if err != nil {
		aw := fmt.Sprintf("failed to assign roles to agent", err)
		s.Logger.Error(aw)
		return err
	}

	// Asynchronously publish an event for each role assignment
	for _, identifier := range roleIdentifiers {
		go func(id interface{}) {
			var event interface{}
			switch idType := id.(type) {
			case uint:
				event = struct {
					AgentID uint
					RoleID  uint
				}{
					AgentID: agentID,
					RoleID:  idType,
				}
			case string:
				event = struct {
					AgentID  uint
					RoleName string
				}{
					AgentID:  agentID,
					RoleName: idType,
				}
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				aw := fmt.Sprintf("failed to publish role assignment event", err)
				s.Logger.Error(aw)
			}
		}(identifier)
	}

	return nil
}

// RevokeRolesFromAgent revokes multiple roles from an agent with transaction management and logging.
func (s *DefaultAgentService) RevokeRolesFromAgent(agentID uint, roles []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.RevokeRolesFromAgentByRoleIDs(agentID, roles); err != nil {
			s.Logger.Error(fmt.Sprintf("Error revoking roles from agent", "agentID", agentID, "error", err))
			return err
		}

		// Log event for role revocation
		go func() {
			for _, roleID := range roles {
				event := struct {
					AgentID uint
					RoleID  uint
					Action  string
				}{
					AgentID: agentID,
					RoleID:  roleID,
					Action:  "role_revoked",
				}
				if err := s.EventPublisher.Publish(event); err != nil {
					s.Logger.Error(fmt.Sprintf("Failed to publish role revocation event", "error", err))
				}
			}
		}()
		return nil
	})
}

// AssignRoleToAgent assigns a role to an agent.
func (as *DefaultAgentService) AssignRoleToAgent2(agentID uint, roleName string) error {
	r, err := as.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return err
	}
	var ra []uint
	ra = append(ra, r.ID)
	return as.AgentDBModel.AssignRoleToAgent(agentID, ra)
}

// RevokeRoleFromAgent revokes a role from an agent.
func (ps *DefaultAgentService) RevokeRoleFromAgent(agentID uint, roleName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.RevokeRoleFromAgentByRoleName(agentID, roleName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to revoke role '%s' from agent %d: %v", roleName, agentID, err))
			return err
		}

		go func() {
			event := struct {
				AgentID  uint
				RoleName string
				Action   string
			}{
				AgentID:  agentID,
				RoleName: roleName,
				Action:   "role_revoked",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish role revocation event: %v", err))
			}
		}()

		return nil
	})
}

func (as *DefaultAgentService) RevokeRoleFromAgent2(agentID uint, roleName string) error {
	return as.AgentDBModel.RevokeRoleFromAgent(agentID, roleName)
}

// GetAgentRoles retrieves all roles assigned to an agent.
func (ps *DefaultAgentService) GetAgentRoles(agentID uint) ([]*models.Role, error) {
	roles, err := ps.AgentDBModel.GetAgentRoles(agentID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve roles for agent %d: %v", agentID, err))
		return nil, err
	}
	return roles, nil
}

func (as *DefaultAgentService) GetAgentRoles2(agentID uint) ([]*models.Role, error) {
	// var roles []*models.Role
	r, err := as.AgentDBModel.GetAgentRoles(agentID)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// AssignPermissionsToAgent from the service layer
func (s *DefaultAgentService) AssignPermissionsToAgentWrapper(agentID uint, permissionNames []string) error {
	publishFunc := func(event interface{}) error {
		return s.EventPublisher.Publish(event)
	}

	err := s.AgentDBModel.AssignPermissionsToAgent(agentID, permissionNames, publishFunc)
	if err != nil {
		ae := fmt.Sprintf("Error assigning permissions to agent:", err)
		s.Logger.Error(ae)
	}
	return err
}

// AssignPermissionsToAgent assigns a set of permissions to an agent, with error handling and event publishing.
func (s *DefaultAgentService) AssignPermissionsToAgent(agentID uint, permissionNames []string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Fetch permissions by names
		var permissions []models.Permission
		if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to fetch permissions: %v", err))
			return err
		}

		// Assign permissions to agent
		for _, perm := range permissions {
			if err := tx.Model(&models.Agents{ID: agentID}).Association("Permissions").Append(&perm); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to assign permission to agent: %v", err))
				return err
			}
		}

		// Asynchronously publish a permissions assignment event
		go func() {
			event := struct {
				AgentID       uint
				AssignedPerms []string
			}{
				AgentID:       agentID,
				AssignedPerms: permissionNames,
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish permissions assignment event: %v", err))
			}
		}()

		return nil
	})
}

func (ps *DefaultAgentService) AssignPermissionsToAgent2(agentID uint, permissionNames []string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AssignPermissionsToAgentByPermissionNames(agentID, permissionNames); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to assign permissions to agent %d: %v", agentID, err))
			return err
		}

		// Publish an event for each permission assignment
		for _, permissionName := range permissionNames {
			go func(pn string) {
				event := struct {
					AgentID        uint
					PermissionName string
					Action         string
				}{
					AgentID:        agentID,
					PermissionName: pn,
					Action:         "permission_assigned",
				}
				if err := ps.EventPublisher.Publish(event); err != nil {
					ps.Logger.Error(fmt.Sprintf("Failed to publish permission assignment event: %v", err))
				}
			}(permissionName)
		}

		return nil
	})
}

// AssignPermissionToAgent assigns a permission to an agent.
//func (as *DefaultAgentService) AssignPermissionToAgent2(agentID uint, permissionName string) error {
//	return as.AgentDBModel.AssignPermissionsToAgent(agentID, []string{permissionName})
//}

// RevokePermissionFromAgent revokes a permission from an agent.
func (ps *DefaultAgentService) RevokePermissionFromAgent(agentID uint, permissionName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.RevokePermissionFromAgentByPermissionName(agentID, permissionName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to revoke permission '%s' from agent %d: %v", permissionName, agentID, err))
			return err
		}

		go func() {
			event := struct {
				AgentID        uint
				PermissionName string
				Action         string
			}{
				AgentID:        agentID,
				PermissionName: permissionName,
				Action:         "permission_revoked",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission revocation event: %v", err))
			}
		}()

		return nil
	})
}

func (as *DefaultAgentService) RevokePermissionFromAgent2(agentID uint, permissionID uint) error {
	return as.AgentDBModel.RevokePermissionFromAgent(agentID, permissionID)
}

// GetAgentPermissions retrieves all permissions associated with an agent's roles.
func (as *DefaultAgentService) GetAgentPermissions(agentID uint) ([]*models.Permission, error) {
	return as.AgentDBModel.GetAgentPermissions(agentID)
}

// AddAgentToUser assigns an agent to a user.
// AddAgentToUser assigns an agent to a user and publishes an event for this action.
func (ps *DefaultAgentService) AddAgentToUser(userID uint, agentID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AddAgentToUser(userID, agentID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to add agent %d to user %d: %v", agentID, userID, err))
			return err
		}

		go func() {
			event := struct {
				UserID  uint
				AgentID uint
				Action  string
			}{
				UserID:  userID,
				AgentID: agentID,
				Action:  "agent_assigned_to_user",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish agent assignment to user event: %v", err))
			}
		}()

		return nil
	})
}

func (as *DefaultAgentService) AddAgentToUser2(userID uint, agentID uint) error {
	return as.AgentDBModel.CreateUserAgent(userID, agentID)
}

// GetAgentsByUser retrieves all agents assigned to a user with comprehensive error handling.
func (ps *DefaultAgentService) GetAgentsByUser(userID uint) ([]*models.Agents, error) {
	agents, err := ps.AgentDBModel.GetAgentsByUser(userID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve agents for user %d: %v", userID, err))
		return nil, err
	}
	return agents, nil
}

func (as *DefaultAgentService) GetAgentsByUser2(userID uint) ([]*models.Agents, error) {
	return as.AgentDBModel.GetAgentsByUser(userID)
}

// AddAgentToTeam assigns an agent to a team, encapsulating the operation in a transaction and publishing an event.
func (ps *DefaultAgentService) AddAgentToTeam(teamID uint, agentID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AddAgentToTeam(teamID, agentID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to add agent %d to team %d: %v", agentID, teamID, err))
			return err
		}

		go func() {
			event := struct {
				TeamID  uint
				AgentID uint
				Action  string
			}{
				TeamID:  teamID,
				AgentID: agentID,
				Action:  "agent_added_to_team",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish agent add to team event: %v", err))
			}
		}()

		return nil
	})
}

// AddAgentToTeam assigns an agent to a team.
func (as *DefaultAgentService) AddAgentToTeam2(teamID uint, agentID uint) error {
	return as.AgentDBModel.AddAgentToTeam(teamID, agentID)
}

// GetAgentsByTeam retrieves all agents assigned to a team.
func (as *DefaultAgentService) GetAgentsByTeam(teamID uint) ([]*models.Agents, error) {
	return as.AgentDBModel.GetAgentsByTeam(teamID)
}

// GrantPermissionToAgent grants a permission to an agent.
func (ps *DefaultAgentService) GrantPermissionToAgent(agentID, permissionID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.GrantPermissionToAgent(agentID, permissionID); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to grant permission %d to agent %d: %v", permissionID, agentID, err))
			return err
		}

		// Asynchronously publish a permission grant event
		go func() {
			event := struct {
				AgentID      uint
				PermissionID uint
				Action       string
			}{
				AgentID:      agentID,
				PermissionID: permissionID,
				Action:       "permission_granted",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission grant event: %v", err))
			}
		}()

		return nil
	})
}

func (as *DefaultAgentService) GrantPermissionToAgent2(agentID uint, permissionID uint) error {
	return as.AgentDBModel.GrantPermissionToAgent(agentID, permissionID)
}

// GetTeamsByAgent retrieves all teams assigned to an agent.
func (ps *DefaultAgentService) GetTeamsByAgent(agentID uint) ([]models.Teams, error) {
	teams, err := ps.AgentDBModel.GetTeamsByAgent(agentID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve teams for agent %d: %v", agentID, err))
		return nil, err
	}
	return teams, nil
}

func (das *DefaultAgentService) GetTeamsByAgent2(agentID uint) ([]models.Teams, error) {
	teams, err := das.AgentDBModel.GetTeamsByAgent(agentID)
	if err != nil {
		return nil, err
	}
	return teams, err
}

// AssignAgentToTeam assigns an agent to a team.
func (das *DefaultAgentService) AssignAgentToTeam(agentID uint, teamID uint) error {
	return das.AgentDBModel.AssignAgentToTeam(teamID, agentID)
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
func (ps *DefaultAgentService) AssignPermissionToTeam(teamID uint, permissionName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AssignPermissionsToTeamByPermissionName(teamID, permissionName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to assign permission '%s' to team %d: %v", permissionName, teamID, err))
			return err
		}

		go func() {
			event := struct {
				TeamID         uint
				PermissionName string
				Action         string
			}{
				TeamID:         teamID,
				PermissionName: permissionName,
				Action:         "permission_assigned_to_team",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission assignment to team event: %v", err))
			}
		}()

		return nil
	})
}

func (das *DefaultAgentService) AssignPermissionToTeam2(teamID uint, permissionName string) error {
	// First, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Then, grant the permission to the team
	return das.AgentDBModel.GrantPermissionToTeam(permission, teamID)
}

// RevokePermissionFromTeam revokes a permission from a team.
func (ps *DefaultAgentService) RevokePermissionFromTeam(teamID uint, permissionName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.RevokePermissionFromTeamByPermissionName(teamID, permissionName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to revoke permission '%s' from team %d: %v", permissionName, teamID, err))
			return err
		}

		go func() {
			event := struct {
				TeamID         uint
				PermissionName string
				Action         string
			}{
				TeamID:         teamID,
				PermissionName: permissionName,
				Action:         "permission_revoked_from_team",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission revocation from team event: %v", err))
			}
		}()

		return nil
	})
}

func (das *DefaultAgentService) RevokePermissionFromTeam2(teamID uint, permissionName string) error {
	// First, get the permission by name
	permission, err := das.AgentDBModel.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Then, revoke the permission from the team
	return das.AgentDBModel.RevokePermissionFromTeam(teamID, permission.ID)
}

// GetTeamPermissions retrieves all permissions associated with a team's roles with error handling.
func (ps *DefaultAgentService) GetTeamPermissions(teamID uint) ([]*models.TeamPermission, error) {
	permissions, err := ps.AgentDBModel.GetTeamPermissions(teamID)
	if err != nil {
		ps.Logger.Error(fmt.Sprintf("Failed to retrieve permissions for team %d: %v", teamID, err))
		return nil, err
	}
	return permissions, nil
}

// GetTeamPermissions retrieves all permissions associated with a team's roles.
func (das *DefaultAgentService) GetTeamPermissions2(teamID uint) (*models.TeamPermission, error) {
	return das.AgentDBModel.GetTeamPermission(teamID)
}

// AddAgentPermissionToRole adds a permission to a role.
// AddAgentPermissionToRole adds a permission to a role with comprehensive transaction management and event publishing.
func (ps *DefaultAgentService) AddAgentPermissionToRole(roleName string, permissionName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.AddPermissionToRoleByPermissionName(roleName, permissionName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to add permission '%s' to role '%s': %v", permissionName, roleName, err))
			return err
		}

		go func() {
			event := struct {
				RoleName       string
				PermissionName string
				Action         string
			}{
				RoleName:       roleName,
				PermissionName: permissionName,
				Action:         "permission_added_to_role",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission addition to role event: %v", err))
			}
		}()

		return nil
	})
}

// RemoveAgentPermissionFromRole removes a permission from a role.
// RemoveAgentPermissionFromRole removes a permission from a role with detailed logging and asynchronous event publishing.
func (ps *DefaultAgentService) RemoveAgentPermissionFromRole(roleName string, permissionName string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := ps.AgentDBModel.RemovePermissionFromRoleByPermissionName(roleName, permissionName); err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to remove permission '%s' from role '%s': %v", permissionName, roleName, err))
			return err
		}

		go func() {
			event := struct {
				RoleName       string
				PermissionName string
				Action         string
			}{
				RoleName:       roleName,
				PermissionName: permissionName,
				Action:         "permission_removed_from_role",
			}
			if err := ps.EventPublisher.Publish(event); err != nil {
				ps.Logger.Error(fmt.Sprintf("Failed to publish permission removal from role event: %v", err))
			}
		}()

		return nil
	})
}

func (das *DefaultAgentService) RemoveAgentPermissionFromRole2(roleName string, permissionName string) error {
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

// GetRolePermissions retrieves permissions associated with a role with enhanced error handling and logging.
func (s *DefaultAgentService) GetRolePermissions(roleID uint) ([]*models.Permission, error) {
	permissions, err := s.AgentDBModel.GetRolePermissions(roleID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to get permissions for role", "roleID", roleID, "error", err))
		return nil, err
	}
	return permissions, nil
}

func (das *DefaultAgentService) GetRolePermissions2(roleName string) ([]*models.Permission, error) {
	// First, get the role by name
	role, err := das.AgentDBModel.GetRoleByName(roleName)
	if err != nil {
		return nil, err
	}

	// Then, retrieve the permissions associated with the role
	return das.AgentDBModel.GetRolePermissions(role.ID)
}

// GrantPermissionToAgent grants a permission to an agent with transactional integrity and detailed logging.
func (s *DefaultAgentService) GrantPermissionToAgent3(agentID, permissionID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.GrantPermissionToAgent(agentID, permissionID); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to grant permission to agent", "agentID", agentID, "permissionID", permissionID, "error", err))
			return err
		}

		// Publish event for permission granted
		go func() {
			event := struct {
				AgentID      uint
				PermissionID uint
				Action       string
			}{
				AgentID:      agentID,
				PermissionID: permissionID,
				Action:       "permission_granted",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish permission granted event", "error", err))
			}
		}()
		return nil
	})
}

// RemoveAgentFromTeam removes an agent from a team, ensuring transactional integrity and event publishing.
func (s *DefaultAgentService) RemoveAgentFromTeam(agentID, teamID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.AgentDBModel.RemoveAgentFromTeam(agentID, teamID); err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to remove agent from team", "agentID", agentID, "teamID", teamID, "error", err))
			return err
		}

		// Asynchronously publish team removal event
		go func() {
			event := struct {
				AgentID uint
				TeamID  uint
				Action  string
			}{
				AgentID: agentID,
				TeamID:  teamID,
				Action:  "agent_removed_from_team",
			}
			if err := s.EventPublisher.Publish(event); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to publish agent removal from team event", "error", err))
			}
		}()
		return nil
	})
}
