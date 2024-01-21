// backend/models/agents.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Agents struct {
	gorm.Model
	ID           uint                  `gorm:"primaryKey" json:"agent_id"`
	FirstName    string                `json:"first_name" binding:"required"`
	LastName     string                `json:"last_name" binding:"required"`
	AgentEmail   string                `json:"agent_email" binding:"required,email"`
	Credentials  AgentLoginCredentials `json:"agent_credentials" gorm:"foreignKey:AgentID"`
	Phone        string                `json:"phoneNumber" binding:"required,e164"`
	RoleID       Role                  `json:"role_id" gorm:"embedded"`
	Unit         Unit                  `json:"unit" gorm:"embedded"`
	SupervisorID int                   `json:"supervisor_id"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Agents) TableName() string {
	return "agents"
}

type Unit struct {
	ID        int       `gorm:"primaryKey" json:"unit_id"`
	UnitName  string    `json:"unit_name"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Unit) TableName() string {
	return "unit"
}

type Role struct {
	ID        int       `gorm:"primaryKey" json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Role) TableName() string {
	return "role"
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

type RoleStorage interface {
	CreateRole(*Role) error
	DeleteRole(int) error
	UpdateRole(*Role) error
	GetRoles() ([]*Role, error)
	GetRoleByID(int) (*Role, error)
	GetRoleByNumber(int) (*Role, error)
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
func (as *AgentDBModel) GetRoleByID(id int) (*Role, error) {
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
