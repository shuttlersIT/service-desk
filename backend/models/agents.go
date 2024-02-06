// backend/models/agents.go

package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Agents represents the schema of the agents table
type Agents struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	FirstName    string          `gorm:"size:255;not null" json:"first_name" binding:"required"`
	LastName     string          `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email        string          `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	PasswordHash string          `gorm:"size:60;not null" json:"-"` // Excluded from JSON responses
	Phone        *string         `gorm:"size:20" json:"phone,omitempty" binding:"omitempty,e164"`
	PositionID   *uint           `gorm:"type:int unsigned" json:"position_id,omitempty"`
	DepartmentID *uint           `gorm:"type:int unsigned" json:"department_id,omitempty"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	ProfilePic   *string         `gorm:"size:255" json:"profile_pic,omitempty"`
	LastLoginAt  *time.Time      `json:"last_login_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	TeamID       *uint           `gorm:"type:int unsigned" json:"team_id,omitempty"`
	SupervisorID *uint           `gorm:"type:int unsigned" json:"supervisor_id,omitempty"`
	Roles        []Role          `gorm:"many2many:agent_roles;" json:"roles"`
}

func (Agents) TableName() string {
	return "agents"
}

// Unit represents the schema of the unit table
type Unit struct {
	ID        uint            `gorm:"primaryKey" json:"unit_id"`
	UnitName  string          `gorm:"size:255;not null" json:"unit_name"`
	Emoji     *string         `gorm:"size:255" json:"emoji,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Unit) TableName() string {
	return "unit"
}

// Permission represents the schema of the permission table
type Permission struct {
	ID          uint            `gorm:"primaryKey" json:"permission_id"`
	Name        string          `gorm:"size:255;not null" json:"permission_name"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Teams represents the schema of the teams table
type Teams struct {
	ID               uint            `gorm:"primaryKey" json:"team_id"`
	TeamName         string          `gorm:"size:255;not null" json:"team_name"`
	Emoji            *string         `gorm:"size:255" json:"emoji,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	TeamPermissionID *uint           `json:"team_permission_id,omitempty" gorm:"type:int unsigned"`
}

func (Teams) TableName() string {
	return "team"
}

// TeamPermission links 'teams' with their 'permissions'.
type TeamPermission struct {
	ID          uint           `gorm:"primaryKey" json:"team_permission_id"`
	TeamID      uint           `gorm:"not null;index:,unique" json:"team_id"`
	Permissions []*Permission  `gorm:"many2many:team_permissions_permissions;" json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Role represents the schema of the role table
type Role struct {
	ID          uint            `gorm:"primaryKey" json:"role_id"`
	RoleName    string          `gorm:"size:255;not null" json:"role_name"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Users       []Users         `gorm:"many2many:user_roles;" json:"-"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

// RoleBase represents a foundational role structure that may be used for additional role metadata
type RoleBase struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:255;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (RoleBase) TableName() string {
	return "role_base"
}

// RolePermission links roles with permissions in a many-to-many relationship
type RolePermission struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	RoleID       uint            `gorm:"not null" json:"role_id"`
	PermissionID uint            `gorm:"not null" json:"permission_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// AgentRole links agents with roles in a many-to-many relationship
type AgentRole struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `gorm:"not null" json:"agent_id"`
	RoleID    uint            `gorm:"not null" json:"role_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentRole) TableName() string {
	return "agent_roles"
}

// UserAgent represents the relationship between a user and an agent
type UserAgent struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null" json:"user_id"`
	AgentID   uint            `gorm:"not null" json:"agent_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserAgent) TableName() string {
	return "user_agents"
}

// TeamAgent represents the relationship between a team and an agent
type TeamAgent struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	TeamID    uint            `gorm:"not null" json:"team_id"`
	AgentID   uint            `gorm:"not null" json:"agent_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TeamAgent) TableName() string {
	return "team_agents"
}

// AgentPermission represents the relationship between an agent and their granted permissions
type AgentPermission struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	AgentID      uint            `gorm:"not null" json:"agent_id"`
	PermissionID uint            `gorm:"not null" json:"permission_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentPermission) TableName() string {
	return "agent_permissions"
}

type AgentStorage interface {
	CreateAgent(agent *Agents) error
	DeleteAgent(agentID uint) error
	UpdateAgent(agentID *Agents) error
	GetAllAgents() ([]*Agents, error)
	GetAgentByID(agentID uint) (*Agents, error)
	GetAgentByNumber(agentNumber int) (*Agents, error)
	AssignRolesToAgent(agentID uint, roleNames []string) error
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	GetRolesByAgent(agentID uint) ([]*Role, error)
	Create(agent *Agents) error
	Update(agent *Agents) error
	Delete(agentID uint) error
	GetAll() ([]*Agents, error)
	GetByID(agentID uint) (*Agents, error)
	GetByNumber(agentNumber int) (*Agents, error)
	AssignRoles(agentID uint, roleNames []string) error
	RevokeRole(agentID uint, roleName string) error
	GetPermissions(agentID uint) ([]*Permission, error)
	AssignPermissions(agentID uint, permissionNames []string) error
	RevokePermission(agentID uint, permissionName string) error
	GetAgents() ([]*Agents, error)
	GetAssignedRoles(agentID uint) ([]*Role, error)
	GetAssignedPermissions(agentID uint) ([]*Permission, error)
	AssignPermissionsToAgent(agentID uint, permissionNames []string) error
	RevokePermissionFromAgent(agentID uint, permissionName string) error
	GetAgentPermissions(agentID uint) ([]*Permission, error)
	AssignAgentToTeam(agentID, teamID uint) error
	RemoveAgentFromTeam(agentID, teamID uint) error
	GetAgentTeams(agentID uint) ([]*Teams, error)
	UpdateAgentWithRolesAndPermissions(agent *Agents, roleIDs, permissionIDs []uint) error
	AssignAgentToMultipleTeams(agentID uint, teamIDs []uint) error
	UpdateAgentPermissions(agentID uint, newPermissionIDs []uint) error
}

type UnitStorage interface {
	CreateUnit(unit *Unit) error
	DeleteUnit(unitID uint) error
	UpdateUnit(unitID *Unit) error
	GetUnits() ([]*Unit, error)
	GetUnitByID(unitID uint) (*Unit, error)
	GetUnitByNumber(unitNumber int) (*Unit, error)
}

type TeamStorage interface {
	CreateTeam(team *Teams) error
	UpdateTeam(team *Teams) error
	DeleteTeam(id uint) error
	GetTeamByID(id uint) (*Teams, error)
	GetTeamByNumber(teamNumber int) (*Teams, error)
	GetTeams() ([]*Teams, error)
}

type RoleStorage interface {
	CreateRole(role *Role) error
	DeleteRole(roleID uint) error
	UpdateRole(roleID *Role) error
	GetRoles() ([]*Role, error)
	GetRoleByID(roleID uint) (*Role, error)
	GetRoleByNumber(roleNumber int) (*Role, error)
	AssignRolesToAgent(agentID uint, roleNames []string) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetRoleByName(name string) (*Role, error)
	AssignPermissionsToRole(roleName string, permissionNames []string) error
	RevokePermissionFromRole(roleName string, permissionName string) error
	GetRolePermissions(roleID uint) ([]*Permission, error)
	GetPermissionsByRole(roleName string) ([]*Permission, error)
}

type PermissionStorage interface {
	GetAllPermissions() ([]*Permission, error)
	CreatePermission(permission *Permission) error
	UpdatePermission(permission *Permission) error
	DeletePermission(id uint) error
	GetPermissionByID(id uint) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	GetPermissions() ([]*Permission, error)
}

type AgentRoleStorage interface {
	AssignRoleToAgent(agentID, roleID uint) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	RevokeRoleFromAgent(agentID uint, roleName string) error
}

type RolePermissionStorage interface {
	AssociatePermissionWithRole(roleID uint, permissionName string, permissionID uint) error
	GetRolePermissionPairs(roleID uint) ([]*RolePermission, error)
	GetRolePermissions(roleID uint) ([]*Permission, error)
	AssignPermissionsToRole(roleName string, permissionNames []string) error
	RevokePermissionFromRole(roleName string, permissionName string) error
}

type UserAgentStorage interface {
	CreateUserAgent(agentID uint, userID uint) error
	GetAgentsByUser(userID uint) ([]*Agents, error)
	GetUsersByAgent(agentID uint) ([]*Users, error)
	RemoveUserAgentRelationship(agentID uint, userID uint) error
}

type RoleAgentStorage interface {
	AssignRoleToAgent(agentID uint, roleID uint) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	RevokeRoleFromAgent(agentID uint, roleName string) error
}

type RoleUserStorage interface {
	AssignRoleToUser(userID uint, roleID uint) error
	GetUserRoles(userID uint) ([]*Role, error)
	RevokeRoleFromUser(userID uint, roleName string) error
}

type UserPermissionStorage interface {
	AssignPermissionToUser(userID uint, permissionID uint) error
	GetUserPermissions(userID uint) ([]*Permission, error)
	RevokePermissionFromUser(userID uint, permissionName string) error
}

type TeamAgentStorage interface {
	CreateTeamAgent(teamAgent *TeamAgent) error
	AddAgentToTeam(agentID uint, teamID uint) error
	GetAgentsByTeam(teamID uint) ([]*Agents, error)
	GetTeamsByAgent(agentID uint) ([]*Teams, error)
	RemoveAgentFromTeam(agentID uint, teamID uint) error
}

// Define additional repository interfaces for user-permission relationships, team-agent relationships, etc.

// AgentDBModel handles database operations for Agent
type AgentDBModel struct {
	DB *gorm.DB
}

// NewAgentDBModel creates a new instance of AgentDBModel
func NewAgentDBModel(db *gorm.DB) *AgentDBModel {
	return &AgentDBModel{
		DB: db,
	}
}

// CreateAgent creates a new Agent.
func (as *AgentDBModel) CreateAgent(agent *Agents) (*Agents, error) {
	result := as.DB.Create(agent).Error
	return agent, result
	//return as.GetAgentByID(uint(result.RowsAffected))
}

// DeleteAgent removes an agent from the database.
func (model *AgentDBModel) DeleteAgent(id uint) error {
	return model.DB.Delete(&Agents{}, id).Error
}

// GetAgentByID retrieves an agent by their ID.
func (model *AgentDBModel) GetAgentByID(id uint) (*Agents, error) {
	var agent Agents
	result := model.DB.Preload("Roles").Where("id = ?", id).First(&agent)

	if result.Error != nil {
		return nil, result.Error
	}
	return &agent, nil
}

// UpdateAgent updates an existing agent's details.
func (model *AgentDBModel) UpdateAgent(agent *Agents) error {
	if agent == nil {
		return errors.New("agent is nil")
	}

	return model.DB.Save(agent).Error
}

// GetAgentByNumber retrieves an Agent by their agent number.
func (as *AgentDBModel) GetAgentByNumber(agentNumber int) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("agent_id = ?", agentNumber).First(&agent).Error
	return &agent, err
}

// GetAllAgents retrieves all agents from the database.
func (as *AgentDBModel) GetAllAgents() ([]*Agents, error) {
	var agents []*Agents
	result := as.DB.Preload("Roles").Find(&agents)

	if result.Error != nil {
		return nil, result.Error
	}
	return agents, nil
}

// AssignRoleToAgent assigns a set of roles to an agent.
func (model *AgentDBModel) AssignRoleToAgent(agentID uint, roleIDs []uint) error {
	var agent Agents
	if err := model.DB.First(&agent, agentID).Error; err != nil {
		return err
	}

	return model.DB.Model(&agent).Association("Roles").Replace(roleIDs)
}

// RevokeRoleFromAgent removes a role from an agent.
func (model *AgentDBModel) RevokeRoleFromAgent(agentID, roleID uint) error {
	var agent Agents
	if err := model.DB.First(&agent, agentID).Error; err != nil {
		return err
	}

	return model.DB.Model(&agent).Association("Roles").Delete(roleID)
}

// GetAgentRoles retrieves all roles associated with an agent.
func (model *AgentDBModel) GetAgentRoles(agentID uint) ([]Role, error) {
	var roles []Role
	err := model.DB.Table("roles").
		Joins("join agent_roles on roles.id = agent_roles.role_id").
		Where("agent_roles.agent_id = ?", agentID).
		Scan(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRolesByAgent retrieves all roles associated with an Agent.
func (as *AgentDBModel) GetRolesByAgent(agentID uint) ([]Role, error) {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}
	return agent.Roles, nil
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) RevokePermissionFromRole(roleName string, permissionName string) error {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err = as.DB.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Delete(&RolePermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

// CreateUnit creates a new unit.
func (as *AgentDBModel) CreateUnit(unit *Unit) error {
	return as.DB.Create(unit).Error
}

// DeleteUnit deletes a unit from the database.
func (as *AgentDBModel) DeleteUnit(id uint) error {
	return as.DB.Delete(&Unit{}, id).Error
}

// UpdateUnit updates the details of an existing unit.
func (as *AgentDBModel) UpdateUnit(unit *Unit) error {
	return as.DB.Save(unit).Error
}

// GetUnits retrieves all units from the database.
func (as *AgentDBModel) GetUnits() ([]*Unit, error) {
	var units []*Unit
	err := as.DB.Find(&units).Error
	return units, err
}

// GetUnitByID retrieves a unit by its ID.
func (as *AgentDBModel) GetUnitByID(id uint) (*Unit, error) {
	var unit Unit
	err := as.DB.Where("id = ?", id).First(&unit).Error
	return &unit, err
}

// GetUnitByNumber retrieves a unit by its unit number.
func (as *AgentDBModel) GetUnitByNumber(unitNumber int) (*Unit, error) {
	var unit Unit
	err := as.DB.Where("unit_id = ?", unitNumber).First(&unit).Error
	return &unit, err
}

// CreateTeam creates a new team.
func (as *AgentDBModel) CreateTeam(team *Teams) error {
	return as.DB.Create(team).Error
}

// DeleteTeam deletes a team from the database.
func (as *AgentDBModel) DeleteTeam(id uint) error {
	return as.DB.Delete(&Teams{}, id).Error
}

// UpdateTeam updates the details of an existing team.
func (as *AgentDBModel) UpdateTeam(team *Teams) error {
	return as.DB.Save(team).Error
}

// GetTeams retrieves all teams from the database.
func (as *AgentDBModel) GetTeams() ([]*Teams, error) {
	var teams []*Teams
	err := as.DB.Find(&teams).Error
	return teams, err
}

// GetTeamByID retrieves a team by its ID.
func (as *AgentDBModel) GetTeamByID(id uint) (*Teams, error) {
	var team Teams
	err := as.DB.Where("id = ?", id).First(&team).Error
	return &team, err
}

// AssignAgentToTeams assigns an agent to a list of teams, ensuring each assignment is unique.
func (db *AgentDBModel) AssignAgentToTeams(agentID uint, teamIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			// Check if the agent is already assigned to the team
			var exists int64
			tx.Model(&TeamAgent{}).Where("agent_id = ? AND team_id = ?", agentID, teamID).Count(&exists)
			if exists == 0 {
				// Create the assignment if it doesn't exist
				if err := tx.Create(&TeamAgent{AgentID: agentID, TeamID: teamID}).Error; err != nil {
					// Return error to rollback transaction
					return err
				}
			}
		}
		// Commit transaction
		return nil
	})
}

// RemoveAgentFromTeam removes an agent from a team.
func (model *AgentDBModel) RemoveAgentFromTeam(agentID, teamID uint) error {
	return model.DB.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error
}

// GetTeamsByAgent retrieves all teams associated with an agent.
func (model *AgentDBModel) GetTeamsByAgent(agentID uint) ([]Teams, error) {
	var teams []Teams
	err := model.DB.Joins("JOIN team_agents ON teams.id = team_agents.team_id").
		Where("team_agents.agent_id = ?", agentID).Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// AssignRolesToAgent assigns roles to an agent.
func (as *AgentDBModel) AssignRolesToAgent2(agentID uint, roleNames []string) error {
	if len(roleNames) == 0 {
		return fmt.Errorf("no role names provided")
	}

	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	roles := make([]Role, len(roleNames))
	for i, roleName := range roleNames {
		role, err := as.GetRoleByName(roleName)
		if err != nil {
			return err
		}
		roles[i] = *role
	}

	agent.Roles = roles

	return as.UpdateAgent(agent)
}

// RevokeRoleFromAgent revokes a role from an agent.
func (as *AgentDBModel) RevokeRoleFromAgent2(agentID uint, roleName string) error {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	var updatedRoles []Role
	for _, role := range agent.Roles {
		if role.RoleName != roleName {
			updatedRoles = append(updatedRoles, role)
		}
	}
	agent.Roles = updatedRoles

	return as.UpdateAgent(agent)
}

// GetAgentPermissions retrieves all permissions associated with an agent's roles.
func (as *AgentDBModel) GetAgentPermissions(agentID uint) ([]*Permission, error) {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}

	var permissions []*Permission
	for _, role := range agent.Roles {
		rolePermissions, err := as.GetPermissionsByRole(role.RoleName)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}

// GetPermissionsByAgent retrieves all permissions associated with an agent.
func (model *AgentDBModel) GetPermissionsByAgent(agentID uint) ([]Permission, error) {
	var permissions []Permission
	err := model.DB.Table("permissions").
		Joins("JOIN agent_permissions ON permissions.id = agent_permissions.permission_id").
		Where("agent_permissions.agent_id = ?", agentID).Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateRolePermission(roleID uint, permissionID uint) error {
	rolePermission := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return as.DB.Create(&rolePermission).Error
}

// GrantPermissionToAgent grants a specific permission to an agent.
func (model *AgentDBModel) GrantPermissionToAgent(agentID, permissionID uint) error {
	agentPermission := AgentPermission{AgentID: agentID, PermissionID: permissionID}
	result := model.DB.Where(&AgentPermission{AgentID: agentID, PermissionID: permissionID}).FirstOrCreate(&agentPermission)
	return result.Error
}

// RevokePermissionFromAgent revokes a specific permission from an agent.
func (model *AgentDBModel) RevokePermissionFromAgent(agentID, permissionID uint) error {
	return model.DB.Where("agent_id = ? AND permission_id = ?", agentID, permissionID).Delete(&AgentPermission{}).Error
}

// AssignPermissionsToAgent assigns permissions to an agent's roles.
func (as *AgentDBModel) AssignPermissionsToAgent2(agentID uint, permissionNames []string) error {
	if len(permissionNames) == 0 {
		return fmt.Errorf("no permission names provided")
	}

	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {
		permission, err := as.GetPermissionByName(permissionName)
		if err != nil {
			// If permission doesn't exist, create it
			newPermission := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(newPermission)
			if err != nil {
				return err
			}
			permission, err = as.GetPermissionByName(permissionName)
			if err != nil {
				return err
			}
		}

		// Check if the permission is already assigned
		assigned := false
		for _, role := range agent.Roles {
			rolePermissions, err := as.GetRolePermissionPairs(role.ID)
			if err != nil {
				return err
			}
			for _, rp := range rolePermissions {
				if rp.ID == permission.ID {
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}

		// If not assigned, assign it to all roles of the agent
		if !assigned {
			for i := range agent.Roles {
				rolePermissions := &RolePermission{
					RoleID:       agent.Roles[i].ID,
					PermissionID: permission.ID,
				}
				err := as.CreateRolePermission(rolePermissions.RoleID, rolePermissions.PermissionID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RevokePermissionFromAgent revokes a permission from an agent.
func (as *AgentDBModel) RevokePermissionFromAgent2(agentID uint, permissionName string) error {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	for _, role := range agent.Roles {
		err := as.RevokePermissionFromRole(role.RoleName, permission.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateAgentWithRolesAndPermissions updates agent details, roles, and permissions atomically.
func (model *AgentDBModel) UpdateAgentWithRolesAndPermissions(agent *Agents, roleIDs, permissionIDs []uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// Update agent details
		if err := tx.Save(agent).Error; err != nil {
			return err
		}

		// Update agent roles
		if err := tx.Model(&agent).Association("Roles").Replace(roleIDs); err != nil {
			return err
		}

		// Update agent permissions
		for _, permissionID := range permissionIDs {
			if err := tx.Model(&AgentPermission{}).Where("agent_id = ?", agent.ID).Update("permission_id", permissionID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// AssignAgentToMultipleTeams assigns an agent to multiple teams, ensuring no duplicates.
func (model *AgentDBModel) AssignAgentToMultipleTeams(agentID uint, teamIDs []uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			// Check if the assignment already exists
			exists := tx.Model(&TeamAgent{}).Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&TeamAgent{}).Error
			if errors.Is(exists, gorm.ErrRecordNotFound) {
				// Create new assignment if it does not exist
				if err := tx.Create(&TeamAgent{AgentID: agentID, TeamID: teamID}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// UpdateAgentPermissions updates the list of permissions for an agent.
func (db *AgentDBModel) UpdateAgentPermissions(agentID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Remove permissions not in the new list
		if err := tx.Where("agent_id = ? AND permission_id NOT IN (?)", agentID, permissionIDs).Delete(&AgentPermission{}).Error; err != nil {
			return err
		}

		// Find existing permissions to avoid duplicates
		var existingPermissions []AgentPermission
		tx.Where("agent_id = ?", agentID).Find(&existingPermissions)

		// Map for quick lookup
		existingMap := make(map[uint]bool)
		for _, p := range existingPermissions {
			existingMap[p.PermissionID] = true
		}

		// Add new permissions
		for _, pid := range permissionIDs {
			if !existingMap[pid] {
				if err := tx.Create(&AgentPermission{AgentID: agentID, PermissionID: pid}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// contains checks if a slice contains a specific uint value.
func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// containsPermission checks if a slice of AgentPermission contains a specific PermissionID.
func containsPermission(s []AgentPermission, e uint) bool {
	for _, a := range s {
		if a.PermissionID == e {
			return true
		}
	}
	return false
}

// ... Rest of your code ...

func (as *AgentDBModel) Create(agent *Agents) error {
	return as.DB.Create(agent).Error
}

func (as *AgentDBModel) Update(agent *Agents) error {
	return as.DB.Save(agent).Error
}

func (as *AgentDBModel) Delete(agentID uint) error {
	return as.DB.Delete(&Agents{}, agentID).Error
}

func (as *AgentDBModel) GetAll() ([]*Agents, error) {
	var agents []*Agents
	err := as.DB.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func (as *AgentDBModel) GetByID(agentID uint) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("id = ?", agentID).First(&agent).Error
	return &agent, err
}

func (as *AgentDBModel) GetByNumber(agentNumber int) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("agent_id = ?", agentNumber).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// AssignRoles assigns roles to an agent.
func (as *AgentDBModel) AssignRoles(agentID uint, roleNames []string) error {
	if len(roleNames) == 0 {
		return fmt.Errorf("no role names provided")
	}

	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	var roles []Role
	for _, roleName := range roleNames {
		role, err := as.GetRoleByName(roleName)
		if err != nil {
			return err
		}
		roles = append(roles, *role)
	}

	agent.Roles = roles

	return as.Update(agent)
}

// RevokeRole revokes a role from an agent.
func (as *AgentDBModel) RevokeRole(agentID uint, roleName string) error {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	var updatedRoles []Role
	for _, role := range agent.Roles {
		if role.RoleName != roleName {
			updatedRoles = append(updatedRoles, role)
		}
	}
	agent.Roles = updatedRoles

	return as.Update(agent)
}

// GetPermissions retrieves all permissions associated with an agent's roles.
func (as *AgentDBModel) GetPermissions(agentID uint) ([]*Permission, error) {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return nil, err
	}

	var permissions []*Permission
	for _, role := range agent.Roles {
		rolePermissions, err := as.GetRolePermissions(role.ID)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}

// GetRolePermissions retrieves all permissions associated with a role.
func (as *AgentDBModel) GetRolePermissions(roleID uint) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (as *AgentDBModel) GetPermissionsByRole(roleName string) ([]*Permission, error) {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return nil, err
	}

	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err = as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", role.ID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// AssignPermissions assigns permissions to an agent's roles.
func (as *AgentDBModel) AssignPermissions(agentID uint, permissionNames []string) error {
	if len(permissionNames) == 0 {
		return fmt.Errorf("no permission names provided")
	}

	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {
		permission, err := as.GetPermissionByName(permissionName)
		if err != nil {
			// If permission doesn't exist, create it
			newPermission := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(newPermission)
			if err != nil {
				return err
			}
			permission, err = as.GetPermissionByName(permissionName)
			if err != nil {
				return err
			}
		}

		// Check if the permission is already assigned
		assigned := false
		for _, role := range agent.Roles {
			rolePermissions, err := as.GetRolePermissionPairs(role.ID)
			if err != nil {
				return err
			}
			for _, rp := range rolePermissions {
				if rp.ID == permission.ID {
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}

		// If not assigned, assign it to all roles of the agent
		if !assigned {
			for i := range agent.Roles {
				rolePermissions := &RolePermission{
					RoleID:       agent.Roles[i].ID,
					PermissionID: permission.ID,
				}
				err := as.CreateRolePermission(rolePermissions.PermissionID, rolePermissions.PermissionID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RevokePermission revokes a permission from an agent.
func (as *AgentDBModel) RevokePermission(agentID uint, permissionName string) error {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	for _, role := range agent.Roles {
		err := as.RevokePermissionFromRole(role.RoleName, permission.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAssignedRoles retrieves all roles assigned to an agent.
func (as *AgentDBModel) GetAssignedRoles(agentID uint) ([]*Role, error) {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return nil, err
	}

	var roles []*Role
	for _, role := range agent.Roles {
		roles = append(roles, &role)
	}

	return roles, nil
}

// GetAssignedPermissions retrieves all permissions assigned to an agent.
func (as *AgentDBModel) GetAssignedPermissions(agentID uint) ([]*Permission, error) {
	permissions, err := as.GetPermissions(agentID)
	if err != nil {
		return nil, err
	}

	var assignedPermissions []*Permission
	for i, permission := range permissions {
		assignedPermissions = append(assignedPermissions, permission)
		fmt.Printf("%v permission added", assignedPermissions[i].Name)
	}

	return assignedPermissions, nil
}

// CreateRole creates a new role.
func (as *AgentDBModel) CreateRole(role *Role) error {
	return as.DB.Create(role).Error
}

// GetRoleByID retrieves a role by its ID.
func (as *AgentDBModel) GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := as.DB.Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetRoleByName retrieves a role by its name.
func (as *AgentDBModel) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := as.DB.Where("name = ?", name).First(&role).Error
	return &role, err
}

// UpdateRole updates the details of an existing role.
func (as *AgentDBModel) UpdateRole(role *Role) error {
	return as.DB.Save(role).Error
}

// DeleteRole deletes a role from the database.
func (as *AgentDBModel) DeleteRole(id uint) error {
	return as.DB.Delete(&Role{}, id).Error
}

// CreatePermission creates a new permission.
func (as *AgentDBModel) CreatePermission(permission *Permission) error {
	return as.DB.Create(permission).Error
}

// GetPermissionByID retrieves a permission by its ID.
func (as *AgentDBModel) GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("id = ?", id).First(&permission).Error
	return &permission, err
}

// GetPermissionByName retrieves a permission by its name.
func (as *AgentDBModel) GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

// UpdatePermission updates the details of an existing permission.
func (as *AgentDBModel) UpdatePermission(permission *Permission) error {
	return as.DB.Save(permission).Error
}

// DeletePermission deletes a permission from the database.
func (as *AgentDBModel) DeletePermission(id uint) error {
	return as.DB.Delete(&Permission{}, id).Error
}

// GetRolePermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetRolePermissionPairs(roleID uint) ([]*RolePermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var rolePermissions []*RolePermission
	err := as.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// CreateUserRole creates a new user-role association.
func (as *AgentDBModel) CreateUserRoleBase(agentID uint, roleID uint) error {
	userRole := AgentRole{AgentID: agentID, RoleID: roleID}
	return as.DB.Create(&userRole).Error
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateRoleBasePermission(roleID uint, permissionID uint) error {
	rolePermission := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return as.DB.Create(&rolePermission).Error
}

// GetUserRoles returns a list of roles assigned to a user.
func (as *AgentDBModel) GetAgentRoles2(agentID uint) ([]*Role, error) {
	var roles []*Role
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Where("userRoles.user_id = ?", agentID).
		Find(&roles).Error
	return roles, err
}

// GetRolesByUser retrieves all roles associated with a user.
func (as *AgentDBModel) GetRoleByAgent(agentID uint) ([]*Role, error) {
	// Implement logic to retrieve roles associated with the user.
	// You might need to join userRoles and roles tables.

	var roles []*Role
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Where("userRoles.user_id = ?", agentID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoles retrieves all roles from the database.
func (as *AgentDBModel) GetRoles() ([]*Role, error) {
	var roles []*Role
	err := as.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByNumber retrieves a role by its role number.
func (as *AgentDBModel) GetRoleByNumber(roleNumber int) (*Role, error) {
	var role Role
	err := as.DB.Where("role_id = ?", roleNumber).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles retrieves all roles.
func (as *AgentDBModel) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := as.DB.Find(&roles).Error
	return roles, err
}

// GetUnitByNumber retrieves a unit by its unit number.
func (as *AgentDBModel) GetTeamByNumber(teamNumber int) (*Teams, error) {
	var team Teams
	err := as.DB.Where("team_id = ?", teamNumber).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// AssignPermissionsToRole assigns multiple permissions to a role.
func (as *AgentDBModel) AssignPermissionsToRoleBase(roleName string, permissionNames []string) error {
	var permission Permission
	var errors []error

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {

		p, er := as.GetPermissionByName(permissionName)
		if er != nil {
			permis := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(permis)
			if err != nil {
				return err
			}
			p, er := as.GetPermissionByName(permissionName)
			if er != nil {
				return er
			}
			permission = *p
		} else {
			permission = *p
		}

		erro := as.AssociatePermissionWithRole(role.ID, permission.Name, permission.ID)
		if err != nil {
			errors = append(errors, erro)
			continue
		}
	}
	fmt.Println(errors)
	return nil
}

// AssociatePermissionWithRole associates a permission with a role.
func (as *AgentDBModel) AssociatePermissionWithRole(roleID uint, permissionName string, permissionID uint) error {
	role, err := as.GetRoleByID(roleID)
	var permission *Permission
	if err != nil {
		return err
	}

	if permissionID <= 0 {
		p, err := as.GetPermissionByName(permissionName)
		if err != nil {
			return err
		}
		permission = p

	} else if permissionName == "" {
		p, err := as.GetPermissionByID(permissionID)
		if err != nil {
			return err
		}
		permission = p

	} else if permissionName == "" || permissionID <= 0 {
		return fmt.Errorf("the permission you selected does not exist")

	} else if permissionName != "" || permissionID > 0 {
		p, err := as.GetPermissionByID(permissionID)
		if err != nil {
			return err
		}
		permission = p
	}

	permittedRole := as.CreateRoleBasePermission(role.ID, permission.ID)

	// Implement logic to associate the permission with the role.
	// You might have a separate table or method to manage role-permission relationships.
	// Example: rolePermissions table with columns (roleID, permissionID).

	if permittedRole != nil {
		return permittedRole
	}

	return nil
}

// AssociatePermissionWithRole associates a permission with a team.
func (as *AgentDBModel) GrantPermissionToTeam(permission *Permission, teamID uint) error {
	tp, err := as.GetTeamPermission(teamID)
	if err != nil {
		er := as.CreateTeamPermission(teamID, permission)
		if er != nil {
			return fmt.Errorf("unable to assign permission to team")
		}
	}
	tp.Permissions = append(tp.Permissions, permission)

	return nil
}

// GetTeamPermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetTeamPermission(teamID uint) (*TeamPermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var teamPermissions *TeamPermission
	err := as.DB.Where("team_id = ?", teamID).Find(&teamPermissions).Error
	if err != nil {
		return nil, err
	}

	return teamPermissions, nil
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateTeamPermission(teamID uint, permission *Permission) error {
	var p []*Permission
	p = append(p, permission)

	teamPermission := TeamPermission{TeamID: teamID, Permissions: p}
	return as.DB.Create(&teamPermission).Error
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) RevokePermissionFromTeam(teamID uint, permissionID uint) error {
	role, err := as.GetTeamByID(teamID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err = as.DB.Where("team_id = ? AND permission = ?", role.ID, permission).
		Delete(&TeamPermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetRolePermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetRolePermissionPairs2(roleID uint) ([]*RolePermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var rolePermissions []*RolePermission
	err := as.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// GetRolePermissions retrieves all permissions associated with a role.
func (as *AgentDBModel) GetRolePermissions2(roleID uint) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetAllPermissions retrieves all permissions.
func (as *AgentDBModel) GetAllPermissions() ([]*Permission, error) {
	var permissions []*Permission
	err := as.DB.Find(&permissions).Error
	return permissions, err
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) RevokePermissionFromRole2(roleName string, permissionName string) error {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err = as.DB.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Delete(&RolePermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (as *AgentDBModel) GetPermissionsByRole2(roleName string) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Joins("JOIN roles ON roles.id = rolePermissions.role_id").
		Where("roles.name = ?", roleName).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// RevokeRoleFromAgent revokes a role from an agent.
func (as *AgentDBModel) RevokeRoleBaseFromAgent(agentID uint, roleName string) error {
	_, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the role from the agent.
	// Example: Delete the corresponding agent-role relationship record.
	err = as.DB.Where("agent_id = ? AND role_id = ?", agentID, role.ID).Delete(&AgentRole{}).Error
	if err != nil {
		return err
	}

	return nil
}

// RevokeRoleFromUser revokes a role from a user.
func (as *AgentDBModel) RevokeRoleBaseFromAgent2(agentID uint, roleName string) error {
	_, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the role from the user.
	// Example: Delete the corresponding user-role relationship record.
	err = as.DB.Where("user_id = ? AND role_id = ?", agentID, role.ID).Delete(&AgentRole{}).Error
	if err != nil {
		return err
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////
// UserAgentRepository implementations
func (as *AgentDBModel) CreateUserAgent(agentID uint, userID uint) error {
	// Implement logic to create a user-agent relationship in the database.
	// Example: Create a record in the userAgents table.
	userAgent := UserAgent{AgentID: agentID, UserID: userID}
	return as.DB.Create(&userAgent).Error
}

func (as *AgentDBModel) GetAgentsByUser(userID uint) ([]*Agents, error) {
	// Implement logic to retrieve agents associated with a user.
	// Example: Join userAgents and agents tables to get the agents by user.
	var agents []*Agents
	err := as.DB.Joins("JOIN userAgents ON agents.id = userAgents.agent_id").
		Where("userAgents.user_id = ?", userID).
		Find(&agents).Error
	return agents, err
}

func (as *AgentDBModel) GetUsersByAgent(agentID uint) ([]*Users, error) {
	// Implement logic to retrieve users associated with an agent.
	// Example: Join userAgents and users tables to get the users by agent.
	var users []*Users
	err := as.DB.Joins("JOIN userAgents ON users.id = userAgents.user_id").
		Where("userAgents.agent_id = ?", agentID).
		Find(&users).Error
	return users, err
}

func (as *AgentDBModel) RemoveUserAgentRelationship(agentID uint, userID uint) error {
	// Implement logic to remove a user-agent relationship from the database.
	// Example: Delete the record from the userAgents table.
	return as.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&UserAgent{}).Error
}

// UserPermissionRepository implementations
func (as *AgentDBModel) AssignPermissionToAgent(agentID uint, permissionID uint) error {
	// Implement logic to assign a permission to a user.
	// Example: Create a record in the userPermissions table.
	agentPermission := AgentPermission{AgentID: agentID, PermissionID: permissionID}
	return as.DB.Create(&agentPermission).Error
}

func (as *AgentDBModel) GetAgentsByTeam(teamID uint) ([]*Agents, error) {
	// Implement logic to retrieve agents associated with a team.
	// Example: Join teamAgents and agents tables to get agents by team.
	var agents []*Agents
	err := as.DB.Joins("JOIN teamAgents ON agents.id = teamAgents.agent_id").
		Where("teamAgents.team_id = ?", teamID).
		Find(&agents).Error
	return agents, err
}

func (as *AgentDBModel) GetTeamByName(name string) (*Teams, error) {
	// Implement logic to retrieve a team by its name.
	var team Teams
	err := as.DB.Where("name = ?", name).First(&team).Error
	return &team, err
}

// AgentTeamRepository implementations
func (as *AgentDBModel) GetAgentTeams(agentID uint) ([]*Teams, error) {
	// Implement logic to retrieve teams associated with an agent.
	// Example: Join teamAgents and teams tables to get teams by agent.
	var teams []*Teams
	err := as.DB.Joins("JOIN teamAgents ON teams.id = teamAgents.team_id").
		Where("teamAgents.agent_id = ?", agentID).
		Find(&teams).Error
	return teams, err
}

// Define more repository interfaces and their implementations as needed.

// Example of a custom query method to retrieve agents with specific criteria.
func (as *AgentDBModel) GetAgentsWithCriteria(criteria string) ([]*Agents, error) {
	// Implement a custom query to retrieve agents based on specific criteria.
	var agents []*Agents
	err := as.DB.Where(criteria).Find(&agents).Error
	return agents, err
}

// UserRepository implementations
func (as *AgentDBModel) AddAgentToUser(userID uint, agentID uint) error {
	// Implement logic to associate an agent with a user.
	userAgent := UserAgent{UserID: userID, AgentID: agentID}
	return as.DB.Create(&userAgent).Error
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) DeleteTeamAgent(teamID uint, agentID uint) error {

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err := as.DB.Where("team_id = ? AND agent_id = ?", teamID, agentID).Delete(&TeamAgent{}).Error
	if err != nil {
		return err
	}

	return nil
}
