// backend/controllers/agent_controllers.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type AgentController struct {
	AgentService *services.DefaultAgentService
}

func NewAgentController(agentService *services.DefaultAgentService) *AgentController {
	return &AgentController{
		AgentService: agentService,
	}
}

// Implement controller methods like GetAgents, CreateAgents, GetAgent, UpdateAgent, DeleteAgent

// CreateAgentHandler creates a new agent.
func (ctrl *AgentController) CreateAgentHandler(c *gin.Context) {
	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.CreateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// UpdateAgentHandler updates an existing agent.
func (ctrl *AgentController) UpdateAgentHandler(c *gin.Context) {
	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.UpdateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// GetAgentByIDHandler retrieves an agent by their ID.
func (ctrl *AgentController) GetAgentByIDHandler(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Param("id"))

	agent, err := ctrl.AgentService.GetAgentByID(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// DeleteAgentHandler deletes an agent by their ID.
func (ctrl *AgentController) DeleteAgentHandler(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.AgentService.DeleteAgent(uint(agentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// GetAllAgentsHandler retrieves all agents.
func (ctrl *AgentController) GetAllAgentsHandler(c *gin.Context) {
	agents, err := ctrl.AgentService.GetAllAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// CreateUnitHandler creates a new unit.
func (ctrl *AgentController) CreateUnitHandler(c *gin.Context) {
	var unit models.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.CreateUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, unit)
}

// DeleteUnitHandler deletes a unit by its ID.
func (ctrl *AgentController) DeleteUnitHandler(c *gin.Context) {
	unitID, _ := strconv.Atoi(c.Param("unit_id"))

	if err := ctrl.AgentService.DeleteUnit(unitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit deleted successfully"})
}

// UpdateUnitHandler updates an existing unit.
func (ctrl *AgentController) UpdateUnitHandler(c *gin.Context) {
	var unit models.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.UpdateUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, unit)
}

// GetUnitsHandler retrieves all units.
func (ctrl *AgentController) GetUnitsHandler(c *gin.Context) {
	units, err := ctrl.AgentService.GetUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}

// GetUnitByIDHandler retrieves a unit by its ID.
func (ctrl *AgentController) GetUnitByIDHandler(c *gin.Context) {
	unitID, _ := strconv.Atoi(c.Param("unit_id"))

	unit, err := ctrl.AgentService.GetUnitByID(unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, unit)
}

// GetUnitByNumberHandler retrieves a unit by its number.
func (ctrl *AgentController) GetUnitByNumberHandler(c *gin.Context) {
	unitNumber, _ := strconv.Atoi(c.Param("unit_number"))

	unit, err := ctrl.AgentService.GetUnitByNumber(unitNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, unit)
}

// CreateTeamHandler creates a new team.
func (ctrl *AgentController) CreateTeamHandler(c *gin.Context) {
	var team models.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.CreateTeam(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// DeleteTeamHandler deletes a team by its ID.
func (ctrl *AgentController) DeleteTeamHandler(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("team_id"))

	if err := ctrl.AgentService.DeleteTeam(teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// CreateRoleHandler creates a new role.
func (ctrl *AgentController) CreateRoleHandler(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.CreateRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// DeleteRoleHandler deletes a role by its ID.
func (ctrl *AgentController) DeleteRoleHandler(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("role_id"))

	if err := ctrl.AgentService.DeleteRole(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// UpdateRoleHandler updates an existing role.
func (ctrl *AgentController) UpdateRoleHandler(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.UpdateRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetRolesHandler retrieves all roles.
func (ctrl *AgentController) GetRolesHandler(c *gin.Context) {
	roles, err := ctrl.AgentService.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRoleByIDHandler retrieves a role by its ID.
func (ctrl *AgentController) GetRoleByIDHandler(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("role_id"))

	role, err := ctrl.AgentService.GetRoleByID(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetRoleByNumberHandler retrieves a role by its number.
func (ctrl *AgentController) GetRoleByNumberHandler(c *gin.Context) {
	roleNumber, _ := strconv.Atoi(c.Param("role_number"))

	role, err := ctrl.AgentService.GetRoleByNumber(roleNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateTeamHandler updates an existing team.
func (ctrl *AgentController) UpdateTeamHandler(c *gin.Context) {
	var team models.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.AgentService.UpdateTeam(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetTeamsHandler retrieves all teams.
func (ctrl *AgentController) GetTeamsHandler(c *gin.Context) {
	teams, err := ctrl.AgentService.GetTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeamByIDHandler retrieves a team by its ID.
func (ctrl *AgentController) GetTeamByIDHandler(c *gin.Context) {
	teamID, _ := strconv.Atoi((c.Param("team_id")))

	team, err := ctrl.AgentService.GetTeamByID(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetTeamByNumberHandler retrieves a team by its number.
func (ctrl *AgentController) GetTeamByNumberHandler(c *gin.Context) {
	teamNumber, _ := strconv.Atoi(c.Param("team_number"))

	team, err := ctrl.AgentService.GetTeamByNumber(teamNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}
