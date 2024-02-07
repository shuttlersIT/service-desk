// backend/models/auth.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	SessionID string    `json:"session_id"`
	Expiry    time.Time `json:"expiry"`
}

// TableName sets the table name for the Session model.
func (Session) TableName() string {
	return "sessions"
}

// CreateSession creates a new user session.
func (as *AuthDBModel) CreateSession(session *Session) error {
	return as.DB.Create(session).Error
}

type UserSession struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      Users     `gorm:"foreignKey:UserID" json:"-"`
	SessionID string    `gorm:"size:255;not null;unique" json:"session_id"` // Unique session identifier
	ExpiresAt time.Time `json:"expires_at"`                                 // Session expiration time
	IP        string    `gorm:"size:45" json:"ip"`                          // IP address of the user at session start
	UserAgent string    `gorm:"type:text" json:"user_agent"`                // User agent of the user's browser/device
}

// TableName sets the table name for the Session model.
func (UserSession) TableName() string {
	return "user_sessions"
}

// GetSessionBySessionID retrieves a user session by its session ID.
func (as *AuthDBModel) GetSessionBySessionID(sessionID string) (*Session, error) {
	var userSession Session
	err := as.DB.Where("session_id = ?", sessionID).First(&userSession).Error
	return &userSession, err
}

// DeleteSessionBySessionID deletes a user session by its session ID.
func (as *AuthDBModel) DeleteSessionBySessionID(sessionID string) error {
	return as.DB.Where("session_id = ?", sessionID).Delete(&Session{}).Error
}

// CreateUserRoles creates user roles for a user.
func (as *AuthDBModel) CreateUserRoles(roles []*UserRole) error {
	return as.DB.Create(roles).Error
}

// GetUserRolesByUserID retrieves user roles for a user by their ID.
func (as *AuthDBModel) GetUserRolesByUserID(userID uint) ([]*UserRole, error) {
	var userRoles []*UserRole
	err := as.DB.Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

// DeleteUserRolesByUserID deletes user roles for a user by their ID.
func (as *AuthDBModel) DeleteUserRolesByUserID(userID uint) error {
	return as.DB.Where("user_id = ?", userID).Delete(&UserRole{}).Error
}

// CreateAgentRoles creates agent roles for an agent.
func (as *AuthDBModel) CreateAgentRoles(roles []*AgentRole) error {
	return as.DB.Create(roles).Error
}

// GetAgentRolesByAgentID retrieves agent roles for an agent by their ID.
func (as *AuthDBModel) GetAgentRolesByAgentID(agentID uint) ([]*AgentRole, error) {
	var agentRoles []*AgentRole
	err := as.DB.Where("agent_id = ?", agentID).Find(&agentRoles).Error
	return agentRoles, err
}

// DeleteAgentRolesByAgentID deletes agent roles for an agent by their ID.
func (as *AuthDBModel) DeleteAgentRolesByAgentID(agentID uint) error {
	return as.DB.Where("agent_id = ?", agentID).Delete(&AgentRole{}).Error
}

// HasPermission checks if a user has the required permission.
func (as *AgentDBModel) HasPermission(userID uint, requiredPermission string) (bool, error) {
	// Implement logic to check if the user has the required permission.
	// You can use the user's ID to retrieve their roles and then check if any of those roles have the required permission.

	var hasPermission bool
	status := false
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Joins("JOIN rolePermissions ON roles.id = rolePermissions.role_id").
		Joins("JOIN permissions ON permissions.id = rolePermissions.permission_id").
		Where("userRoles.user_id = ? AND permissions.name = ?", userID, requiredPermission).
		First(&hasPermission).Error

	if err != nil {
		return status, err
	}

	return hasPermission, nil
}

type UserAgentMapping struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	AgentID uint `json:"agent_id"`
}

// TableName sets the table name for the UserAgentMapping model.
func (UserAgentMapping) TableName() string {
	return "userAgentMappings"
}

// CreateUserAgentMapping creates a mapping between a user and an agent.
func (as *AuthDBModel) CreateUserAgentMapping(mapping *UserAgentMapping) error {
	return as.DB.Create(mapping).Error
}

// GetUserAgentMappingsByUserID retrieves all agent mappings for a user by their ID.
func (as *AuthDBModel) GetUserAgentMappingsByUserID(userID uint) ([]*UserAgentMapping, error) {
	var mappings []*UserAgentMapping
	err := as.DB.Where("user_id = ?", userID).Find(&mappings).Error
	return mappings, err
}

// GetAgentMappingsByAgentID retrieves all user mappings for an agent by their ID.
func (as *AuthDBModel) GetAgentMappingsByAgentID(agentID uint) ([]*UserAgentMapping, error) {
	var mappings []*UserAgentMapping
	err := as.DB.Where("agent_id = ?", agentID).Find(&mappings).Error
	return mappings, err
}

// DeleteUserAgentMapping deletes a mapping between a user and an agent.
func (as *AuthDBModel) DeleteUserAgentMapping(mapping *UserAgentMapping) error {
	return as.DB.Delete(mapping).Error
}

type UserAgentAccess struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	AgentID uint `json:"agent_id"`
	Access  bool `json:"access"`
}

// TableName sets the table name for the UserAgentAccess model.
func (UserAgentAccess) TableName() string {
	return "userAgentAccess"
}

// CreateUserAgentAccess creates an access record between a user and an agent.
func (as *AuthDBModel) CreateUserAgentAccess(access *UserAgentAccess) error {
	return as.DB.Create(access).Error
}

// GetUserAgentAccessByUserID retrieves access records for a user by their ID.
func (as *AuthDBModel) GetUserAgentAccessByUserID(userID uint) ([]*UserAgentAccess, error) {
	var accesses []*UserAgentAccess
	err := as.DB.Where("user_id = ?", userID).Find(&accesses).Error
	return accesses, err
}

// GetAgentAccessByAgentID retrieves access records for an agent by their ID.
func (as *AuthDBModel) GetAgentAccessByAgentID(agentID uint) ([]*UserAgentAccess, error) {
	var accesses []*UserAgentAccess
	err := as.DB.Where("agent_id = ?", agentID).Find(&accesses).Error
	return accesses, err
}

// UpdateUserAgentAccess updates an access record between a user and an agent.
func (as *AuthDBModel) UpdateUserAgentAccess(access *UserAgentAccess) error {
	return as.DB.Save(access).Error
}

// DeleteUserAgentAccess deletes an access record between a user and an agent.
func (as *AuthDBModel) DeleteUserAgentAccess(access *UserAgentAccess) error {
	return as.DB.Delete(access).Error
}

type UserAgentGroup struct {
	gorm.Model
	Name         string         `json:"name"`
	Type         string         `json:"type"` // Example: "user" or "agent"
	GroupMembers []*GroupMember `json:"group_members,omitempty"`
}

type GroupMember struct {
	gorm.Model
	GroupID uint `json:"group_id"`
	UserID  uint `json:"user_id,omitempty"`
	AgentID uint `json:"agent_id,omitempty"`
}

// TableName sets the table name for the UserAgentGroup model.
func (UserAgentGroup) TableName() string {
	return "userAgentGroups"
}

// CreateGroup creates a new user-agent group.
func (as *AuthDBModel) CreateGroup(group *UserAgentGroup) error {
	return as.DB.Create(group).Error
}

// GetGroupByID retrieves a user-agent group by its ID.
func (as *AuthDBModel) GetGroupByID(groupID uint) (*UserAgentGroup, error) {
	var group UserAgentGroup
	err := as.DB.Where("id = ?", groupID).Preload("GroupMembers").First(&group).Error
	return &group, err
}

// UpdateGroup updates the details of a user-agent group.
func (as *AuthDBModel) UpdateGroup(group *UserAgentGroup) error {
	return as.DB.Save(group).Error
}

// DeleteGroup deletes a user-agent group.
func (as *AuthDBModel) DeleteGroup(groupID uint) error {
	return as.DB.Delete(&UserAgentGroup{}, groupID).Error
}

// AddUserToGroup adds a user to a user-agent group.
func (as *AuthDBModel) AddUserToGroup(groupID, userID uint) error {
	member := &GroupMember{GroupID: groupID, UserID: userID}
	return as.DB.Create(member).Error
}

// AddAgentToGroup adds an agent to a user-agent group.
func (as *AuthDBModel) AddAgentToGroup(groupID, agentID uint) error {
	member := &GroupMember{GroupID: groupID, AgentID: agentID}
	return as.DB.Create(member).Error
}

// RemoveUserFromGroup removes a user from a user-agent group.
func (as *AuthDBModel) RemoveUserFromGroup(groupID, userID uint) error {
	return as.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&GroupMember{}).Error
}

// RemoveAgentFromGroup removes an agent from a user-agent group.
func (as *AuthDBModel) RemoveAgentFromGroup(groupID, agentID uint) error {
	return as.DB.Where("group_id = ? AND agent_id = ?", groupID, agentID).Delete(&GroupMember{}).Error
}

// GetSessionByUserID retrieves an active session for a user by their ID.
func (as *AuthDBModel) GetSessionByUserID(userID uint) (*Session, error) {
	var session Session
	err := as.DB.Where("user_id = ? AND expiry > ?", userID, time.Now()).First(&session).Error
	return &session, err
}

// DeleteSession deletes a user's session.
func (as *AuthDBModel) DeleteSession(session *Session) error {
	return as.DB.Delete(session).Error
}

func (sm *Session) ExpireSessionsForUser(userID uint) error {
	// Implementation details...
}

func (am *AuthDBModel) OAuthLogin(provider string, token string) (*Agents, error) {
	// Implementation to handle OAuth flow with different providers
}

func (am *AuthDBModel) CheckAccess(agentID uint, resource string, action string) (bool, error) {
	// Check if the agent has permission to perform the action on the resource
}

/*
// CreateRole creates a new role.
func (as *AuthDBModel) CreateRoleBase(role *Role) error {
	return as.DB.Create(role).Error
}

// GetRoleByName retrieves a role by its name.
func (as *AuthDBModel) GetRoleBase(name string) (*RoleBase, error) {
	var roleBase RoleBase
	err := as.DB.Where("name = ?", name).First(&roleBase).Error
	return &roleBase, err
}

// AssignRoleToUser assigns a role to a user.
func (as *AuthDBModel) AssignRoleToUser(userID uint, roleBase string) error {
	role, err := as.GetRoleBase(roleBase)
	if err != nil {
		return err
	}

	user, err := as.UserDBModel.GetUserByID(userID)

	user.RoleBase = *role

	// Implement logic to associate the user with the role.
	// You might have a separate table or method to manage user-role relationships.
	// Example: userRoles table with columns (userID, roleID).

	return nil
}

// ///////////////////////////////////////////////////////////////////////////
type Permission struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName sets the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}

// CreatePermission creates a new permission.
func (as *AuthDBModel) CreatePermission(permission *Permission) error {
	return as.DB.Create(permission).Error
}

// GetPermissionByName retrieves a permission by its name.
func (as *AuthDBModel) GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

// AssociatePermissionWithRole associates a permission with a role.
func (as *AuthDBModel) AssociatePermissionWithRoleBase(roleBase string, permissionName string) error {
	role, err := as.GetRoleBase(roleBase)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Implement logic to associate the permission with the role.
	// You might have a separate table or method to manage role-permission relationships.
	// Example: rolePermissions table with columns (roleID, permissionID).

	return nil
}

// AssignRolesToUser assigns multiple roles to a user.
func (as *AuthDBModel) AssignRoleBaseToUser(userID uint, roleBases []string) error {
	for _, roleBase := range roleBases {
		if err := as.AssignRoleToUser(userID, roleBase); err != nil {
			return err
		}
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
func (as *AuthDBModel) AssignPermissionsToRoleBase(roleBase string, permissionNames []string) error {
	role, err := as.GetRoleBase(roleBase)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {
		if err := as.AssociatePermissionWithRoleBase(roleBase, permissionName); err != nil {
			return err
		}
	}
	return nil
}

// GetRolePermissions retrieves all permissions associated with a role.
func (as *AuthDBModel) GetRoleBasePermissions(roleBase string) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	// Example query:
	// err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
	//     Where("rolePermissions.role_id = ?", roleID).
	//     Find(&permissions).Error

	return permissions, nil
}

// AssignRoleToUser assigns a role to a user.
func (as *AuthDBModel) AssignRoleBaseToUser(userID uint, roleBase string) error {
	user, err := as.UserDBModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleBase(roleBase)
	if err != nil {
		return err
	}

	// Implement logic to associate the role with the user.
	// You might have a separate table or method to manage user-role relationships.
	// Example: userRoles table with columns (user_id, role_id).

	return nil
}

// RevokeRoleFromUser revokes a role from a user.
func (as *AuthDBModel) RevokeRoleBaseFromUser(userID uint, roleName string) error {
	user, err := as.UserDBModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleBase(roleName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the role from the user.
	// Example: Delete the corresponding user-role relationship record.

	return nil
}

// GetRolesByUser retrieves all roles associated with a user.
func (as *AuthDBModel) GetRoleBaseByUser(userID uint) ([]*Role, error) {
	// Implement logic to retrieve roles associated with the user.
	// You might need to join userRoles and roles tables.

	var roles []*Role
	// Example query:
	// err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
	//     Where("userRoles.user_id = ?", userID).
	//     Find(&roles).Error

	return roles, nil
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AuthDBModel) RevokePermissionFromRoleBase(roleName string, permissionName string) error {
	role, err := as.GetRoleBase(roleName)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.

	return nil
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (as *AuthDBModel) GetPermissionsByRole(roleName string) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	// Example query:
	// err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
	//     Where("rolePermissions.role_id = ?", roleID).
	//     Find(&permissions).Error

	return permissions, nil
}

// Define a model for userRoles to associate users with roles.
type UserRoleBase struct {
	gorm.Model
	UserID     uint
	RoleBaseID uint
}

// Define a model for rolePermissions to associate roles with permissions.
type RoleBasePermission struct {
	gorm.Model
	RoleBaseID   uint
	PermissionID uint
}

// CreateUserRole creates a new user-role association.
func (as *AuthDBModel) CreateUserRoleBase(userID uint, roleBaseID uint) error {
	userRole := UserRoleBase{UserID: userID, RoleBaseID: roleBaseID}
	return as.DB.Create(&userRole).Error
}

// CreateRolePermission creates a new role-permission association.
func (as *AuthDBModel) CreateRoleBasePermission(roleID, permissionID uint) error {
	rolePermission := RoleBasePermission{RoleBaseID: roleID, PermissionID: permissionID}
	return as.DB.Create(&rolePermission).Error
}

// HasPermission checks if a user has the required permission.
func (as *AuthDBModel) HasPermission(userID uint, requiredPermission string) (bool, error) {
	// Implement logic to check if the user has the required permission.
	// You can use the user's ID to retrieve their roles and then check if any of those roles have the required permission.

	var hasPermission bool
	// Example query:
	// err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
	//     Joins("JOIN rolePermissions ON roles.id = rolePermissions.role_id").
	//     Joins("JOIN permissions ON permissions.id = rolePermissions.permission_id").
	//     Where("userRoles.user_id = ? AND permissions.name = ?", userID, requiredPermission).
	//     First(&hasPermission).Error

	return hasPermission, nil
}

// CreateRole creates a new role.
func (as *AuthDBModel) CreateRoleBase(roleBase *Role) error {
	return as.DB.Create(roleBase).Error
}

// UpdateRole updates an existing role.
func (as *AuthDBModel) UpdateRoleBase(roleBase *RoleBase) error {
	return as.DB.Save(roleBase).Error
}

// DeleteRole deletes a role by ID.
func (as *AuthDBModel) DeleteRoleBase(roleID uint) error {
	return as.DB.Delete(&Role{}, roleID).Error
}

// GetRoleByID retrieves a role by its ID.
func (as *AuthDBModel) GetRoleBaseByID(roleBaseID uint) (*RoleBase, error) {
	var roleBase RoleBase
	err := as.DB.Where("id = ?", roleBaseID).First(&roleBase).Error
	return &roleBase, err
}

// GetAllRoles retrieves all roles.
func (as *AuthDBModel) GetAllRoleBase() ([]*RoleBase, error) {
	var roleBases []*RoleBase
	err := as.DB.Find(&roleBases).Error
	return roleBases, err
}

*/
