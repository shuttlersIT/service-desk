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
	gorm.Model
	FirstName    string     `gorm:"size:255;not null" json:"first_name" binding:"required"`
	LastName     string     `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email        string     `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"` // Excluded from JSON responses
	Phone        *string    `gorm:"size:20" json:"phone,omitempty" binding:"omitempty,e164"`
	PositionID   uint       `gorm:"index;type:int unsigned" json:"position_id,omitempty"`
	DepartmentID uint       `gorm:"index;type:int unsigned" json:"department_id,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	ProfilePic   *string    `gorm:"size:255" json:"profile_pic,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	TeamID       *uint      `gorm:"type:int unsigned" json:"team_id,omitempty"`
	SupervisorID *uint      `gorm:"type:int unsigned" json:"supervisor_id,omitempty"`
	Roles        []Role     `gorm:"many2many:agent_roles;" json:"roles"`
	Biography    string     `json:"biography,omitempty"`
	UserID       uint       `gorm:"primaryKey" json:"user_id"`
	AgentDetails Users      `gorm:"foreignKey:UserID" json:"-"`
}

func (Agents) TableName() string {
	return "agents"
}

type AgentProfile struct {
	gorm.Model
	AgentID         uint   `gorm:"primaryKey;autoIncrement:false" json:"agent_id"`
	Bio             string `gorm:"type:text" json:"bio,omitempty"`
	AvatarURL       string `gorm:"type:text" json:"avatar_url,omitempty"`
	Preferences     string `gorm:"type:text" json:"preferences,omitempty"`      // Assuming JSON format
	PrivacySettings string `gorm:"type:text" json:"privacy_settings,omitempty"` // Assuming JSON format
}

func (AgentProfile) TableName() string {
	return "agent_profiles"
}

// Unit represents the schema of the unit table
type Unit struct {
	gorm.Model
	UnitName string  `gorm:"size:255;not null" json:"unit_name"`
	Emoji    *string `gorm:"size:255" json:"emoji,omitempty"`
}

func (Unit) TableName() string {
	return "unit"
}

// Permission represents the schema of the permission table
type Permission struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Teams represents the schema of the teams table
type Teams struct {
	gorm.Model
	TeamName         string  `gorm:"size:255;not null" json:"team_name"`
	Emoji            *string `gorm:"size:255" json:"emoji,omitempty"`
	TeamPermissionID *uint   `gorm:"type:int unsigned" json:"team_permission_id,omitempty"`
}

func (Teams) TableName() string {
	return "team"
}

// TeamPermission links 'teams' with their 'permissions'.
type TeamPermission struct {
	gorm.Model
	TeamID      uint          `gorm:"not null;index:idx_team_id,unique" json:"team_id"`
	Permissions []*Permission `gorm:"many2many:team_permissions_permissions;" json:"permissions,omitempty"`
}

// Role represents the schema of the role table
type Role struct {
	gorm.Model
	RoleName    string  `gorm:"size:255;not null" json:"role_name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Users       []Users `gorm:"many2many:user_roles;" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

// RoleBase represents a foundational role structure that may be used for additional role metadata
type RoleBase struct {
	gorm.Model
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (RoleBase) TableName() string {
	return "role_base"
}

// RolePermission links roles with permissions in a many-to-many relationship
type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"not null;index:idx_role_permission,unique" json:"role_id"`
	PermissionID uint `gorm:"not null;index:idx_role_permission,unique" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// AgentRole links agents with roles in a many-to-many relationship
type AgentRole struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `gorm:"not null;index:idx_agent_role,unique" json:"agent_id"`
	RoleID    uint            `gorm:"not null;index:idx_agent_role,unique" json:"role_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentRole) TableName() string {
	return "agent_roles"
}

type AgentTrainingSession struct {
	gorm.Model
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location" gorm:"type:varchar(255)"`
	TrainerID   uint      `json:"trainer_id" gorm:"index;not null"`
	Attendees   []Agents  `gorm:"many2many:agent_training_attendees;" json:"attendees"`
}

func (AgentTrainingSession) TableName() string {
	return "agent_training_session"
}

// UserAgent represents the relationship between a user and an agent
type UserAgent struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index:idx_user_agent,user_id" json:"user_id"`
	AgentID   uint            `gorm:"not null;index:idx_user_agent,agent_id" json:"agent_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserAgent) TableName() string {
	return "user_agents"
}

// TeamAgent represents the relationship between a team and an agent
type TeamAgent struct {
	gorm.Model
	TeamID  uint `gorm:"not null;index:idx_team_agent,team_id" json:"team_id"`
	AgentID uint `gorm:"not null;index:idx_team_agent,agent_id" json:"agent_id"`
}

func (TeamAgent) TableName() string {
	return "team_agents"
}

// AgentPermission represents the relationship between an agent and their granted permissions
type AgentPermission struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	AgentID      uint            `json:"agent_id" gorm:"index;not null"`
	PermissionID uint            `json:"permission_id" gorm:"index;not null"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentPermission) TableName() string {
	return "agent_permissions"
}

type SearchCriteria struct {
	gorm.Model
	Name       string `json:"name,omitempty"`
	Role       string `json:"role,omitempty"`
	Department string `json:"department,omitempty"`
}

func (SearchCriteria) TableName() string {
	return "search_criteria"
}

type AgentStorage interface {
	LogAgentActivity(log AgentActivityLog) error
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

func (db *AgentDBModel) LogAgentActivity(log AgentActivityLog) error {
	return db.DB.Create(&log).Error
}

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

// CreateAgent adds a new agent to the database.
func (db *AgentDBModel) CreateAgent(agent *Agents) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(agent).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAgent deletes an agent by their ID.
func (db *AgentDBModel) DeleteAgent(agentID uint) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Delete(&Agents{}, agentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateAgent updates an existing agent's information.
func (db *AgentDBModel) UpdateAgent(agent *Agents) error {
	if agent == nil {
		return errors.New("provided agent is nil")
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(agent).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
func (db *AgentDBModel) AssignRoleToAgent(agentID uint, roleIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.First(&agent, agentID).Error; err != nil {
			return err
		}

		if err := tx.Model(&agent).Association("Roles").Replace(roleIDs); err != nil {
			return err
		}
		return nil
	})
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

// AssignAgentToTeam assigns an agent to a team, ensuring the assignment is unique.
func (db *AgentDBModel) AssignAgentToTeam(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var teamAgent TeamAgent
		// Check if the agent is already assigned to the team to prevent duplicate entries
		if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&teamAgent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// The assignment does not exist; proceed to create a new one
				teamAgent = TeamAgent{AgentID: agentID, TeamID: teamID}
				if err := tx.Create(&teamAgent).Error; err != nil {
					// Error encountered while creating the assignment; rollback the transaction
					return err
				}
			} else {
				// An unexpected error occurred; rollback the transaction
				return err
			}
		}
		// The assignment already exists or has been successfully created
		return nil
	})
}

// RemoveAgentFromTeam removes an agent from a specific team.
func (db *AgentDBModel) RemoveAgentFromTeam(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error; err != nil {
			// Error encountered while removing the assignment; rollback the transaction
			return err
		}
		// Successfully removed the assignment
		return nil
	})
}

// RemoveAgentFromTeams safely removes an agent from multiple teams.
func (db *AgentDBModel) RemoveAgentFromTeams(agentID uint, teamIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error; err != nil {
				return err // Rollback transaction on error
			}
		}
		return nil // Commit transaction if all deletions succeed
	})
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

// UpdateAgentPermissions updates an agent's permissions comprehensively, ensuring transactional integrity.
func (db *AgentDBModel) UpdateAgentPermissions(agentID uint, newPermissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// First, remove any existing permissions that are not in the newPermissionIDs list
		if err := tx.Where("agent_id = ? AND permission_id NOT IN ?", agentID, newPermissionIDs).Delete(&AgentPermission{}).Error; err != nil {
			return err // Rollback on error
		}

		// Next, add new permissions from newPermissionIDs that the agent does not already have
		for _, permissionID := range newPermissionIDs {
			var existing AgentPermission
			result := tx.Where("agent_id = ? AND permission_id = ?", agentID, permissionID).First(&existing)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Permission does not exist for this agent; proceed to add it
				if err := tx.Create(&AgentPermission{AgentID: agentID, PermissionID: permissionID}).Error; err != nil {
					return err // Rollback on error
				}
			} // Ignore any found records, as we do not need to add existing permissions
		}

		return nil // Commit the transaction
	})
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

// UpdateAgentRolesAndPermissions updates both roles and permissions for an agent.
func (db *AgentDBModel) UpdateAgentRolesAndPermissions(agentID uint, newRoleIDs, newPermissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Update roles
		if err := updateAgentRolesWithinTransaction(tx, agentID, newRoleIDs); err != nil {
			return err
		}
		// Update permissions
		if err := updateAgentPermissionsWithinTransaction(tx, agentID, newPermissionIDs); err != nil {
			return err
		}
		return nil // Commit transaction if both updates succeed
	})
}

func updateAgentRolesWithinTransaction(tx *gorm.DB, agentID uint, newRoleIDs []uint) error {
	// Logic to update agent's roles based on newRoleIDs
	// Similar to the UpdateAgentPermissions logic, with appropriate adjustments for roles
	return nil
}

func updateAgentPermissionsWithinTransaction(tx *gorm.DB, agentID uint, newPermissionIDs []uint) error {
	// Similar logic to UpdateAgentPermissions example provided earlier
	return nil
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
func (db *AgentDBModel) UpdateAgentPermissions2(agentID uint, permissionIDs []uint) error {
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

// AssignMultipleAgentsToTeam facilitates the assignment of multiple agents to a specific team.
func (db *AgentDBModel) AssignMultipleAgentsToTeam(agentIDs []uint, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, agentID := range agentIDs {
			var teamAgent TeamAgent
			if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&teamAgent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Agent is not yet assigned to the team, proceed with the assignment
					teamAgent = TeamAgent{AgentID: agentID, TeamID: teamID}
					if err := tx.Create(&teamAgent).Error; err != nil {
						return err // Rollback on error
					}
				} else {
					return err // Rollback on unexpected error
				}
			}
			// If the agent is already assigned, skip to the next agent without doing anything
		}
		return nil // Commit the transaction
	})
}

// UnassignAgentsFromTeam handles the removal of multiple agents from a team.
func (db *AgentDBModel) UnassignAgentsFromTeam(agentIDs []uint, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Attempt to remove each specified agent from the team
		if err := tx.Where("agent_id IN ? AND team_id = ?", agentIDs, teamID).Delete(&TeamAgent{}).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit the transaction
	})
}

// UpdateTeamPermissions updates a team's permissions.
func (db *AgentDBModel) UpdateTeamPermissions(teamID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Find existing TeamPermission relationship
		var teamPermission TeamPermission
		err := tx.Where("team_id = ?", teamID).First(&teamPermission).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // Rollback on unexpected error
		}

		// Update the permissions associated with the team
		if err := tx.Model(&teamPermission).Association("Permissions").Replace(permissionIDs); err != nil {
			return err // Rollback on error
		}

		return nil // Commit the transaction
	})
}

// PerformComplexOperationWithTransaction demonstrates a complex operation involving multiple steps.
func (db *AgentDBModel) PerformComplexOperationWithTransaction(agentID uint, teamID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Example step 1: Assign an agent to a team
		if err := db.AssignAgentToTeam(agentID, teamID); err != nil {
			return err // Rollback on error
		}

		// Example step 2: Update agent's permissions
		if err := db.UpdateAgentPermissions(agentID, permissionIDs); err != nil {
			return err // Rollback on error
		}

		// Additional steps can be added here, each protected by the transaction
		// ...

		return nil // Commit the transaction if all steps are successful
	})
}

type AgentEvent struct {
	gorm.Model
	Title       *string    `gorm:"size:255;not null" json:"title"` // nullable string
	Description *string    `gorm:"type:text" json:"description"`   // nullable string
	ActionType  *string    `json:"action_type"`
	StartTime   *time.Time `json:"start_time"`               // nullable time.Time
	Details     *string    `gorm:"type:text" json:"details"` // nullable string
	Timestamp   time.Time  `json:"time_stamp"`               // nullable time.Time
	AllDay      bool       `json:"all_day"`
	Location    *string    `gorm:"size:255" json:"location"` // nullable string
	AgentID     uint       `gorm:"not null;index" json:"user_id"`
	Agents      Agents     `gorm:"foreignKey:UserID" json:"-"`
}

func (AgentEvent) TableName() string {
	return "agent_event"
}

type AgentActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `json:"agent_id"`
	Activity  string    `json:"activity"`
	Timestamp time.Time `json:"timestamp"`
}

func (AgentActivityLog) TableName() string {
	return "agent_activity_log"
}

// LogAction records an action performed by an agent for auditing purposes.
func (db *AgentDBModel) LogAgentAction(agentID uint, actionType string, details string) error {
	actionLog := AgentEvent{
		AgentID:    agentID,
		ActionType: &actionType,
		Details:    &details,
		Timestamp:  time.Now(),
	}

	return db.DB.Create(&actionLog).Error
}

// GetAgentActions retrieves the recent actions performed by a specific agent.
func (db *AgentDBModel) GetAgentActions(agentID uint, limit int) ([]AgentEvent, error) {
	var agentActions []AgentEvent
	if err := db.DB.Where("agent_id = ?", agentID).Order("timestamp desc").Limit(limit).Find(&agentActions).Error; err != nil {
		return nil, err
	}
	return agentActions, nil
}

// GetAllActions retrieves recent actions performed by all agents.
func (db *AgentDBModel) GetAllActions(limit int) ([]AgentEvent, error) {
	var allActions []AgentEvent
	if err := db.DB.Order("timestamp desc").Limit(limit).Find(&allActions).Error; err != nil {
		return nil, err
	}
	return allActions, nil
}

// SearchAgents searches for agents based on the provided criteria with pagination and sorting options.
func (db *AgentDBModel) SearchAgents(criteria SearchCriteria, page, pageSize int, sortBy string, sortOrder string) ([]Agents, error) {
	var agents []Agents
	query := db.DB

	// Apply search criteria
	if criteria.Name != "" {
		query = query.Where("name LIKE ?", "%"+criteria.Name+"%")
	}
	if criteria.Role != "" {
		query = query.Where("role = ?", criteria.Role)
	}
	if criteria.Department != "" {
		query = query.Where("department = ?", criteria.Department)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Apply sorting
	if sortBy != "" {
		orderBy := sortBy + " " + sortOrder
		query = query.Order(orderBy)
	}

	// Execute the query
	if err := query.Find(&agents).Error; err != nil {
		return nil, err
	}

	return agents, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////
type AgentSchedule struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ShiftType string    `json:"shift_type" gorm:"type:varchar(100);not null"` // E.g., "morning", "night"
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}
type AgentShift struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	ShiftDate time.Time `json:"shift_date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	ShiftType string    `json:"shift_type" gorm:"type:varchar(100);not null"` // E.g., "Morning", "Evening", "Night"
}

type AgentSkill struct {
	gorm.Model
	AgentID   uint   `json:"agent_id" gorm:"index;not null"`
	SkillName string `json:"skill_name" gorm:"type:varchar(255);not null"`
	Level     int    `json:"level" gorm:"not null"` // E.g., 1 to 5, where 5 is expert level
}

type AgentFeedback struct {
	gorm.Model
	AgentID      uint      `json:"agent_id" gorm:"index;not null"`
	FeedbackType string    `json:"feedback_type" gorm:"type:varchar(100);not null"` // E.g., "Customer", "Supervisor"
	Score        int       `json:"score" gorm:"type:int;not null"`                  // Typically a numerical score, e.g., 1-10
	Comments     string    `json:"comments" gorm:"type:text"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

type AgentKPI struct {
	gorm.Model
	AgentID     uint    `json:"agent_id" gorm:"index;not null"`
	KPIName     string  `json:"kpi_name" gorm:"type:varchar(255);not null"`
	Value       float64 `json:"value" gorm:"type:decimal(10,2);not null"` // Example: Average resolution time, Customer satisfaction score
	TargetValue float64 `json:"target_value" gorm:"type:decimal(10,2)"`
	Period      string  `json:"period" gorm:"type:varchar(100);not null"` // E.g., "Monthly", "Quarterly"
}

type AgentOnboarding struct {
	gorm.Model
	AgentID        uint       `json:"agent_id" gorm:"index;not null"`
	OnboardingStep string     `json:"onboarding_step" gorm:"type:varchar(255);not null"` // E.g., "Documentation", "Training", "Mentoring"
	Status         string     `json:"status" gorm:"type:varchar(100);not null"`          // E.g., "Pending", "Completed"
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type AgentTeam struct {
	gorm.Model
	TeamName    string   `json:"team_name" gorm:"type:varchar(255);not null;unique"`
	Description string   `json:"description" gorm:"type:text"`
	LeaderID    uint     `json:"leader_id" gorm:"index"` // Optional: ID of the team leader
	Members     []Agents `gorm:"many2many:agent_teams_members;"`
}

type AgentTrainingRecord struct {
	gorm.Model
	AgentID          uint      `json:"agent_id" gorm:"index;not null"`
	TrainingModuleID uint      `json:"training_module_id" gorm:"index;not null"`
	Score            int       `json:"score"`
	Feedback         string    `json:"feedback" gorm:"type:text"`
	CompletedAt      time.Time `json:"completed_at"`
}

type AgentAvailability struct {
	gorm.Model
	AgentID       uint       `json:"agent_id" gorm:"index;not null"`
	Availability  string     `json:"availability" gorm:"type:varchar(100);not null"` // E.g., "Available", "Busy", "Offline"
	LastUpdated   time.Time  `json:"last_updated"`
	NextAvailable *time.Time `json:"next_available,omitempty"`
}

type AgentContactInfo struct {
	gorm.Model
	AgentID      uint   `json:"agent_id" gorm:"index;not null"`
	ContactType  string `json:"contact_type" gorm:"type:varchar(100);not null"` // E.g., "Phone", "Email", "Skype"
	ContactValue string `json:"contact_value" gorm:"type:varchar(255);not null"`
}

type AgentLoginActivity struct {
	gorm.Model
	AgentID    uint       `json:"agent_id" gorm:"index;not null"`
	LoginTime  time.Time  `json:"login_time"`
	LogoutTime *time.Time `json:"logout_time,omitempty"`
	IP         string     `json:"ip" gorm:"type:varchar(45)"`
}

type AgentVacation struct {
	gorm.Model
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Reason     string    `json:"reason" gorm:"type:text"`
	ApprovedBy uint      `json:"approved_by"` // Optional: Manager or supervisor who approved the vacation
}
type CustomerInteraction struct {
	gorm.Model
	CustomerID      uint      `json:"customer_id" gorm:"index;not null"`
	AgentID         uint      `json:"agent_id" gorm:"index;not null"`
	Channel         string    `json:"channel" gorm:"type:varchar(100);not null"` // E.g., "Email", "Phone", "Chat"
	Content         string    `json:"content" gorm:"type:text;not null"`
	InteractionTime time.Time `json:"interaction_time"`
}

type FeedbackReview struct {
	gorm.Model
	FeedbackID uint      `json:"feedback_id" gorm:"index;not null"`
	ReviewerID uint      `json:"reviewer_id" gorm:"index;not null"` // Manager or QA specialist
	Review     string    `json:"review" gorm:"type:text"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

type AgentTicketAssignment struct {
	gorm.Model
	TicketID   uint      `json:"ticket_id" gorm:"index;not null"`
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	AssignedAt time.Time `json:"assigned_at"`
}

type AgentTrainingModule struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	ModuleType  string `json:"module_type" gorm:"type:varchar(100);not null"` // E.g., "Online", "In-Person"
	Duration    int    `json:"duration"`                                      // Duration in minutes
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}

type AgentCertification struct {
	gorm.Model
	AgentID       uint       `json:"agent_id" gorm:"index;not null"`
	Certification string     `json:"certification" gorm:"type:varchar(255);not null"`
	IssuedBy      string     `json:"issued_by" gorm:"type:varchar(255)"`
	IssuedDate    time.Time  `json:"issued_date"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
}

type AgentPerformanceReview struct {
	gorm.Model
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	ReviewDate time.Time `json:"review_date"`
	Score      float64   `json:"score"`
	Feedback   string    `json:"feedback" gorm:"type:text"`
}

type AgentLeaveRequest struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	LeaveType string    `json:"leave_type" gorm:"type:varchar(100);not null"` // E.g., "Annual", "Sick", "Personal"
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status" gorm:"type:varchar(100);not null"` // E.g., "Pending", "Approved", "Denied"
}

type AgentScheduleOverride struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Reason    string    `json:"reason" gorm:"type:text"`
}

type AgentSkillSet struct {
	gorm.Model
	AgentID uint   `json:"agent_id" gorm:"index;not null"`
	Skill   string `json:"skill" gorm:"type:varchar(255);not null"`
	Level   string `json:"level" gorm:"type:varchar(100);not null"` // E.g., "Beginner", "Intermediate", "Expert"
}
