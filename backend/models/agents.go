// backend/models/agents.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Agents struct {
	gorm.Model
	ID                     uint                  `gorm:"primaryKey" json:"agent_id"`
	FirstName              string                `json:"first_name" binding:"required"`
	LastName               string                `json:"last_name" binding:"required"`
	AgentEmail             string                `json:"agent_email" binding:"required,email"`
	Credentials            AgentLoginCredentials `json:"agent_credentials" gorm:"foreignKey:AgentID"`
	Phone                  string                `json:"phoneNumber" binding:"required,e164"`
	RoleID                 Role                  `json:"role_id" gorm:"embedded"`
	Team                   Teams                 `json:"team_id" gorm:"embedded"`
	Unit                   Unit                  `json:"unit" gorm:"embedded"`
	SupervisorID           uint                  `json:"supervisor_id"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	DeletedAt              time.Time             `json:"deleted_at"`
	RoleBase               RoleBase              `json:"role_base" gorm:"embedded"`
	ResetPasswordRequestID uint                  `json:"reset_password_reset" gorm:"foreignKey:AgentID"`
	Roles                  []Role                `json:"roles" gorm:"-"`
}

// TableName sets the table name for the Agent model.
func (Agents) TableName() string {
	return "agents"
}

type Unit struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"unit_id"`
	UnitName  string    `json:"unit_name"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Unit model.
func (Unit) TableName() string {
	return "unit"
}

type Permission struct {
	gorm.Model
	ID          uint   `json:"permission_id"`
	Name        string `json:"permission_name"`
	Description string `json:"description"`
}

// TableName sets the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}

type Teams struct {
	gorm.Model
	ID               uint      `gorm:"primaryKey" json:"team_id"`
	TeamName         string    `json:"team_name"`
	Emoji            string    `json:"emoji"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TeamPermissionID uint      `json:"team_permission" gorm:"foreignKey:TicketID"`
}

// TableName sets the table name for the Teams model.
func (Teams) TableName() string {
	return "team"
}

type TeamPermission struct {
	gorm.Model
	ID          uint         `json:"team_permission_id" gorm:"primaryKey"`
	TeamID      uint         `json:"team_id" gorm:"foreignKey:TeamID"`
	Permissions []Permission `json:"permission_id" gorm:"embedded"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Role struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Role model.
func (Role) TableName() string {
	return "role"
}

type RoleBase struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName sets the table name for the RoleBase model.
func (RoleBase) TableName() string {
	return "role_base"
}

// Define a model for rolePermissions to associate roles with permissions.
type RolePermission struct {
	gorm.Model
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
}

// TableName sets the table name for the RoleBase model.
func (RolePermission) TableName() string {
	return "role_permission"
}

type AgentRole struct {
	gorm.Model
	AgentID uint `json:"agent_id"`
	RoleID  uint `json:"role_id"`
}

// TableName sets the table name for the AgentRole model.
func (AgentRole) TableName() string {
	return "agentRoles"
}

// UserAgent represents the relationship between a user and an agent.
type UserAgent struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	UserID  uint `json:"user_id"`
	AgentID uint `json:"agent_id"`
}

// TableName sets the table name for the AgentRole model.
func (UserAgent) TableName() string {
	return "userAgent"
}

// TeamAgent represents the relationship between a team and an agent.
type TeamAgent struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	TeamID  uint `json:"team_id"`
	AgentID uint `json:"agent_id"`
}

// TableName sets the table name for the AgentRole model.
func (TeamAgent) TableName() string {
	return "teamAgent"
}

// AgentPermissions represents the relationship between an agent and their granted permissions.
type AgentPermission struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	AgentID      uint `json:"agent_id" gorm:"foreignKey:AgentID"`
	PermissionID uint `json:"permission_id" gorm:"foreignKey:PermissionID"`
}

// TableName sets the table name for the AgentRole model.
func (AgentPermission) TableName() string {
	return "agentPermissions"
}

type AgentStorage interface {
	CreateAgent(agent *Agents) error
	DeleteAgent(agentID uint) error
	UpdateAgent(agentID *Agents) error
	GetAllAgents() ([]*Agents, error)
	GetAgentByID(agentID uint) (*Agents, error)
	GetAgentByNumber(agentNumber int) (*Agents, error)
	AssignRolesToAgent(agentID uint, roleNames []string) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	RevokeRoleFromAgent(agentID uint, roleName string) error
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
	GetAgentPermissions(agentID uint) ([]*Permission, error)
	AssignPermissionsToAgent(agentID uint, permissionNames []string) error
	RevokePermissionFromAgent(agentID uint, permissionName string) error
	GetAssignedRoles(agentID uint) ([]*Role, error)
	GetAssignedPermissions(agentID uint) ([]*Permission, error)
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
	result := as.DB.Create(agent)

	return as.GetAgentByID(uint(result.RowsAffected))
}

// DeleteAgent deletes an Agent from the database.
func (as *AgentDBModel) DeleteAgent(id uint) error {
	return as.DB.Delete(&Agents{}, id).Error
}

// UpdateAgent updates the details of an existing Agent.
func (as *AgentDBModel) UpdateAgent(agent *Agents) error {
	return as.DB.Save(agent).Error
}

// GetAgentByID retrieves an Agent by its ID.
func (as *AgentDBModel) GetAgentByID(id uint) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("id = ?", id).First(&agent).Error
	return &agent, err
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
	err := as.DB.Find(&agents).Error
	return agents, err
}

// AssignRolesToAgent assigns roles to an Agent.
func (as *AgentDBModel) AssignRolesToAgent(agentID uint, roleNames []string) error {
	if len(roleNames) == 0 {
		return fmt.Errorf("no role names provided")
	}

	// Retrieve the Agent
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	// Fetch the roles and assign them to the Agent
	var roles []Role
	for _, roleName := range roleNames {
		role, err := as.GetRoleByName(roleName)
		if err != nil {
			return err
		}
		roles = append(roles, *role)
	}

	agent.Roles = roles

	// Update the Agent with assigned roles
	if err := as.UpdateAgent(agent); err != nil {
		return err
	}

	return nil
}

// GetAgentRoles retrieves all roles associated with an Agent.
func (as *AgentDBModel) GetAgentRoles(agentID uint) ([]*Role, error) {
	var roles []*Role
	err := as.DB.Model(&Agents{}).Where("id = ?", agentID).Association("Roles").Find(&roles)
	return roles, err
}

// RevokeRoleFromAgent revokes a role from an Agent.
func (as *AgentDBModel) RevokeRoleFromAgent(agentID uint, roleName string) error {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	var updatedRoles []Role
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	for _, r := range agent.Roles {
		if r.ID != role.ID {
			updatedRoles = append(updatedRoles, r)
		}
	}

	agent.Roles = updatedRoles

	if err := as.UpdateAgent(agent); err != nil {
		return err
	}

	return nil
}

// GetRolesByAgent retrieves all roles associated with an Agent.
func (as *AgentDBModel) GetRolesByAgent(agentID uint) ([]Role, error) {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}
	return agent.Roles, nil
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

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateRolePermission(roleID uint, permissionID uint) error {
	rolePermission := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return as.DB.Create(&rolePermission).Error
}

// AssignPermissionsToAgent assigns permissions to an agent's roles.
func (as *AgentDBModel) AssignPermissionsToAgent(agentID uint, permissionNames []string) error {
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
				Description: "",
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
func (as *AgentDBModel) RevokePermissionFromAgent(agentID uint, permissionName string) error {
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
				Description: "",
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

// UserRole struct with UserID and RoleID fields (already defined).

// AssignRoleToUser assigns a role to a user.
func (as *AgentDBModel) AssignRoleToAgent(agentID, roleID uint) error {
	userRole := AgentRole{AgentID: agentID, RoleID: roleID}
	return as.DB.Create(&userRole).Error
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
				Description: "",
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
	tp.Permissions = append(tp.Permissions, *permission)

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
	var p []Permission
	p = append(p, *permission)

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

// TeamAgent Storage Implementations
func (as *AgentDBModel) CreateTeamAgent(agentID uint, teamID uint) error {

	return as.AddAgentToTeam(agentID, teamID)
}

// TeamAgentRepository implementations
func (as *AgentDBModel) AddAgentToTeam(agentID uint, teamID uint) error {
	// Implement logic to add an agent to a team.
	// Example: Create a record in the teamAgents table.
	teamAgent := TeamAgent{AgentID: agentID, TeamID: teamID}
	return as.DB.Create(&teamAgent).Error
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

func (as *AgentDBModel) GetTeamsByAgent(agentID uint) ([]*Teams, error) {
	// Implement logic to retrieve teams associated with an agent.
	// Example: Join teamAgents and teams tables to get teams by agent.
	var teams []*Teams
	err := as.DB.Joins("JOIN teamAgents ON teams.id = teamAgents.team_id").
		Where("teamAgents.agent_id = ?", agentID).
		Find(&teams).Error
	return teams, err
}

func (as *AgentDBModel) RemoveAgentFromTeam(agentID uint, teamID uint) error {
	// Implement logic to remove an agent from a team.
	// Example: Delete the record from the teamAgents table.
	return as.DB.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error
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

func (as *AgentDBModel) GrantPermissionToAgent(agentID uint, permissionID uint) error {
	// Implement logic to grant a permission to an agent.
	agentPermission := AgentPermission{AgentID: agentID, PermissionID: permissionID}
	return as.DB.Create(&agentPermission).Error
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
