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
	RoleBase               RoleBase              `json:"role_base" gorm:"embedded"`
	ResetPasswordRequestID uint                  `json:"reset_password_reset" gorm:"embedded"`
	Roles                  []Role                `json:"roles"`
}

// TableName sets the table name for the Agent model.
func (Agents) TableName() string {
	return "agents"
}

type Unit struct {
	ID        uint      `gorm:"primaryKey" json:"unit_id"`
	UnitName  string    `json:"unit_name"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Unit) TableName() string {
	return "unit"
}

type Permission struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName sets the table name for the Agent model.
func (Permission) TableName() string {
	return "permissions"
}

type Teams struct {
	ID        uint      `gorm:"primaryKey" json:"team_id"`
	TeamName  string    `json:"team_name"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Teams) TableName() string {
	return "team"
}

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Role) TableName() string {
	return "role"
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

type AgentStorage interface {
	CreateAgent(*Agents) error
	DeleteAgent(int) error
	UpdateAgent(*Agents) error
	GetAgents() ([]*Agents, error)
	GetAgentByID(int) (*Agents, error)
	GetAgentByNumber(int) (*Agents, error)
}

type UnitStorage interface {
	CreateUnit(*Unit) error
	DeleteUnit(int) error
	UpdateUnit(*Unit) error
	GetUnits() ([]*Unit, error)
	GetUnitByID(int) (*Unit, error)
	GetUnitByNumber(int) (*Unit, error)
}

type TeamStorage interface {
	CreateTeam(*Teams) error
	DeleteTeam(int) error
	UpdateTeam(*Teams) error
	GetTeams() ([]*Teams, error)
	GetTeamByID(int) (*Teams, error)
	GetTeamByNumber(int) (*Teams, error)
}

type RoleStorage interface {
	CreateRole(*Role) error
	DeleteRole(int) error
	UpdateRole(*Role) error
	GetRoles() ([]*Role, error)
	GetRoleByID(uint) (*Role, error)
	GetRoleByNumber(int) (*Role, error)
	AssignRolesToUser(userID uint, roleNames []string) error
}

// AgentModel handles database operations for Agent
type AgentDBModel struct {
	DB *gorm.DB
}

// NewAgentModel creates a new instance of Agent
func NewAgentDBModel(db *gorm.DB) *AgentDBModel {
	return &AgentDBModel{
		DB: db,
	}
}

// //////////////////////////////////////////////////
// //////////////////////////////////////?
type RoleBase struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName sets the table name for the Role model.
func (RoleBase) TableName() string {
	return "role_base"
}

////////////////////////////////////////////////////

// CreateAgent creates a new Agent.
func (as *AgentDBModel) CreateAgent(agent *Agents) error {
	return as.DB.Create(agent).Error
}

// GetAgentByID retrieves a user by its ID.
func (as *AgentDBModel) GetAgentByID(id uint) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("id = ?", id).First(&agent).Error
	return &agent, err
}

// UpdateAgent updates the details of an existing agent.
func (as *AgentDBModel) UpdateAgent(agent *Agents) error {
	if err := as.DB.Save(agent).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a Agent from the database.
func (as *AgentDBModel) DeleteAgent(id uint) error {
	if err := as.DB.Delete(&Agents{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllAgents retrieves all agents from the database.
func (as *AgentDBModel) GetAllAgents() ([]*Agents, error) {
	var agents []*Agents
	err := as.DB.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// GetAgentByNumber retrieves an agent by their agent number.
func (as *AgentDBModel) GetAgentByNumber(agentNumber int) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("agent_id = ?", agentNumber).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// CreateUnit creates a new unit.
func (as *AgentDBModel) CreateUnit(unit *Unit) error {
	return as.DB.Create(unit).Error
}

// DeleteUnit deletes a unit from the database.
func (as *AgentDBModel) DeleteUnit(id int) error {
	if err := as.DB.Delete(&Unit{}, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUnit updates the details of an existing unit.
func (as *AgentDBModel) UpdateUnit(unit *Unit) error {
	return as.DB.Save(unit).Error
}

// GetUnits retrieves all units from the database.
func (as *AgentDBModel) GetUnits() ([]*Unit, error) {
	var units []*Unit
	err := as.DB.Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// GetUnitByID retrieves a unit by its ID.
func (as *AgentDBModel) GetUnitByID(id int) (*Unit, error) {
	var unit Unit
	err := as.DB.Where("unit_id = ?", id).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// GetUnitByNumber retrieves a unit by its unit number.
func (as *AgentDBModel) GetUnitByNumber(unitNumber int) (*Unit, error) {
	var unit Unit
	err := as.DB.Where("unit_id = ?", unitNumber).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// CreateRole creates a new role.
func (as *AgentDBModel) CreateRole(role *Role) error {
	return as.DB.Create(role).Error
}

// DeleteRole deletes a role from the database.
func (as *AgentDBModel) DeleteRole(id int) error {
	if err := as.DB.Delete(&Role{}, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRole updates the details of an existing role.
func (as *AgentDBModel) UpdateRole(role *Role) error {
	return as.DB.Save(role).Error
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

// GetRoleByID retrieves a role by its ID.
func (as *AgentDBModel) GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := as.DB.Where("role_id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
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

// CreateUnit creates a new unit.
func (as *AgentDBModel) CreateTeam(team *Teams) error {
	return as.DB.Create(team).Error
}

// DeleteUnit deletes a unit from the database.
func (as *AgentDBModel) DeleteTeam(id int) error {
	if err := as.DB.Delete(&Teams{}, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUnit updates the details of an existing unit.
func (as *AgentDBModel) UpdateTeam(team *Teams) error {
	return as.DB.Save(team).Error
}

// GetUnits retrieves all units from the database.
func (as *AgentDBModel) GetTeams() ([]*Teams, error) {
	var teams []*Teams
	err := as.DB.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// GetUnitByID retrieves a unit by its ID.
func (as *AgentDBModel) GetTeamByID(id int) (*Teams, error) {
	var team Teams
	err := as.DB.Where("team_id = ?", id).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
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

// GetRoleByName retrieves a role by its name.
func (as *AgentDBModel) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := as.DB.Where("name = ?", name).First(&role).Error
	return &role, err
}

// AssignRoleToUser assigns a role to a agent.
func (as *AgentDBModel) AssignRoleToUser(agentID uint, roleName string) error {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	agent.RoleID = *role

	erro := as.UpdateAgent(agent)
	if erro != nil {
		return erro
	}
	// Implement logic to associate the user with the role.
	// You might have a separate table or method to manage user-role relationships.
	// Example: userRoles table with columns (userID, roleID).

	return nil
}

// AssignRoleToUser assigns a role to a agent.
func (as *AgentDBModel) AssignRolesToAgent(agentID uint, roleNames []string) error {
	if len(roleNames) <= 0 {
		return fmt.Errorf("no role names provided")
	}
	var error []error
	var roles []Role
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	for _, r := range roleNames {

		role, err := as.GetRoleByName(r)
		if err != nil {
			error = append(error, err)
		}
		roles = append(agent.Roles, *role)
	}

	agent.Roles = roles

	erro := as.UpdateAgent(agent)
	if erro != nil {
		return erro
	}
	// Implement logic to associate the user with the role.
	// You might have a separate table or method to manage user-role relationships.
	// Example: userRoles table with columns (userID, roleID).
	fmt.Println(error)
	return nil
}

// GetAllRoles retrieves all roles.
func (as *AgentDBModel) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := as.DB.Find(&roles).Error
	return roles, err
}

// CreatePermission creates a new permission.
func (as *AgentDBModel) CreatePermission(permission *Permission) error {
	return as.DB.Create(permission).Error
}

// UpdatePermission updates an existing permission.
func (as *AgentDBModel) UpdatePermission(permission *Permission) error {
	return as.DB.Save(permission).Error
}

// DeletePermission deletes a permission by ID.
func (as *AgentDBModel) DeletePermission(permissionID uint) error {
	return as.DB.Delete(&Permission{}, permissionID).Error
}

// GetPermissionByID retrieves a permission by its ID.
func (as *AgentDBModel) GetPermissionByID(permissionID uint) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("id = ?", permissionID).First(&permission).Error
	return &permission, err
}

// GetAllPermissions retrieves all permissions.
func (as *AgentDBModel) GetAllPermissions() ([]*Permission, error) {
	var permissions []*Permission
	err := as.DB.Find(&permissions).Error
	return permissions, err
}

// GetPermissionByName retrieves a permission by its name.
func (as *AgentDBModel) GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("name = ?", name).First(&permission).Error
	return &permission, err
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

// GetUserRoles retrieves all roles associated with a user.
func (as *AuthDBModel) GetUserRoleBase(userID uint) ([]*RoleBase, error) {
	// Implement logic to retrieve roles associated with the user.
	// You might need to join userRoles and roles tables.

	var roles []*RoleBase
	// Example query:
	// err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
	//     Where("userRoles.user_id = ?", userID).
	//     Find(&roles).Error

	return roles, nil
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

// GetRolesByAgent retrieves all roles associated with an agent.
func (as *AgentDBModel) GetRolesByAgent(agentID uint) ([]*Role, error) {
	// Implement logic to retrieve roles associated with the agent.
	// You might need to join agentRoles and roles tables.

	var roles []*Role
	err := as.DB.Joins("JOIN agentRoles ON roles.id = agentRoles.role_id").
		Where("agentRoles.agent_id = ?", agentID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
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

// Define a model for rolePermissions to associate roles with permissions.
type RolePermission struct {
	gorm.Model
	RoleID       uint
	PermissionID uint
}

// CreateUserRole creates a new user-role association.
func (as *AgentDBModel) CreateUserRoleBase(agentID uint, roleID uint) error {
	userRole := AgentRole{AgentID: agentID, RoleID: roleID}
	return as.DB.Create(&userRole).Error
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateRoleBasePermission(roleID, permissionID uint) error {
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
func (as *AgentDBModel) GetAgentRoles(agentID uint) ([]*Role, error) {
	var roles []*Role
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Where("userRoles.user_id = ?", agentID).
		Find(&roles).Error
	return roles, err
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
