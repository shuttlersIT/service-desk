// backend/controllers/agent_controllers.go

package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

var d = godotenv.Load()

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

	err := ctrl.AgentService.CreateAgent(&agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Agent created successfully"})
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

	if err := ctrl.AgentService.DeleteUnit(uint(unitID)); err != nil {
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

	unit, err := ctrl.AgentService.GetUnitByID(uint(unitID))
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

	if err := ctrl.AgentService.DeleteTeam(uint(teamID)); err != nil {
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

	if err := ctrl.AgentService.DeleteRole(uint(roleID)); err != nil {
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

	role, err := ctrl.AgentService.GetRoleByID(uint(roleID))
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

	team, err := ctrl.AgentService.GetTeamByID(uint(teamID))
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

func (ac *AgentController) CreateAgent2(w http.ResponseWriter, r *http.Request) {
	var agent models.Agents
	err := json.NewDecoder(r.Body).Decode(&agent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	erro := ac.AgentService.CreateAgent(&agent)
	if err != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&agent)
}

func (ac *AgentController) GetAgentByID2(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agentID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	agent, err := ac.AgentService.GetAgentByID(uint(agentID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(agent)
}

// Implement other CRUD methods for agents here...

func (ac *AgentController) AssignAgentToTeam2(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agentID, err := strconv.ParseUint(params["agentID"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamID, err := strconv.ParseUint(params["teamID"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ac.AgentService.AssignAgentToTeam(uint(agentID), uint(teamID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// //////////////////////////////////////////////////////////////////////////////////////
func (ac *AgentController) CreateAgent(c *gin.Context) {
	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AgentService.CreateAgent(&agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Agent created successfully"})
}

func (ac *AgentController) GetAgentByID(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := ac.AgentService.GetAgentByID(uint(agentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// Implement other CRUD methods for agents here...

func (ac *AgentController) AssignAgentToTeam(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var team models.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.AssignAgentToTeam(uint(agentID), team.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent assigned to team successfully"})
}

func (ac *AgentController) GetAgentsByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agents, err := ac.AgentService.GetAgentsByUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

func (ac *AgentController) RevokeAgentFromTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.RevokeAgentFromTeam(uint(teamID), agent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent revoked from team successfully"})
}

func (ac *AgentController) GetAgentTeams(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teams, err := ac.AgentService.GetAgentTeams(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// Implement other agent-related controller methods here...

// User-Agent methods

func (ac *AgentController) AddAgentToUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var agent models.Agents
	if err := c.ShouldBindJSON(&agent.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.AddAgentToUser(uint(userID), agent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent added to user successfully"})
}

// Implement other User-Agent methods here...

// Team-Agent methods

func (ac *AgentController) AddAgentToTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.AddAgentToTeam(uint(teamID), agent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent assigned to team successfully"})
}

// Implement other Team-Agent methods here...
// Agent-Permission methods

func (ac *AgentController) GrantPermissionToAgent(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permissionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.GrantPermissionToAgent(uint(agentID), uint(permissionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission granted to agent successfully"})
}

func (ac *AgentController) RevokePermissionFromAgent(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionID, err := strconv.Atoi(c.Param("permissionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, er := ac.AgentService.GetPermissionByID(uint(permissionID))
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.RevokePermissionFromAgent(uint(agentID), permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked from agent successfully"})
}

func (ac *AgentController) GetAgentPermissions(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissions, err := ac.AgentService.GetAgentPermissions(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// Implement other Agent-Permission methods here...

// Team-Permission methods

func (ac *AgentController) AssignPermissionToTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.AssignPermissionToTeam(uint(teamID), permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to team successfully"})
}

func (ac *AgentController) RevokePermissionFromTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.RevokePermissionFromTeam(uint(teamID), permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked from team successfully"})
}

func (ac *AgentController) GetTeamPermissions(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamPermissions, err := ac.AgentService.GetTeamPermissions(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teamPermissions)
}

// Implement other Team-Permission methods here...
// User-Agent methods

// Team-Agent methods

func (ac *AgentController) GetAgentsByTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agents, err := ac.AgentService.GetAgentsByTeam(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// Implement other User-Agent and Team-Agent methods here...

// Agent-Permission methods

func (ac *AgentController) AddAgentPermissionToRole(c *gin.Context) {
	roleName := c.Param("roleName")
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AgentService.AddAgentPermissionToRole(roleName, permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission added to role successfully"})
}

func (ac *AgentController) RemoveAgentPermissionFromRole(c *gin.Context) {
	roleName := c.Param("roleName")
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AgentService.RemoveAgentPermissionFromRole(roleName, permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission removed from role successfully"})
}

func (ac *AgentController) GetRolePermissions(c *gin.Context) {
	roleName := c.Param("roleName")
	role, e := ac.AgentService.GetRoleByName(roleName)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		return
	}

	permissions, err := ac.AgentService.GetRolePermissions(role.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// Implement other Agent-Permission methods here...

// More methods can be added as needed...
// Agent-Permission methods (Continued)

func (ac *AgentController) AssignPermissionToAgent(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	p := c.PostFormArray("perms")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newPermission []*models.Permission
	for _, ps := range p {
		np := &models.Permission{Name: ps}
		if err := c.ShouldBindJSON(&newPermission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		e := ac.AgentService.CreatePermission(np)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			return
		}
	}
	err = ac.AgentService.AssignPermissionsToAgent(uint(agentID), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to agent successfully"})
}

func (ac *AgentController) RevokePermissionFromAgent2(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.RevokePermissionFromAgent(uint(agentID), permission.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked from agent successfully"})
}

func (ac *AgentController) GetAgentPermissions2(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissions, err := ac.AgentService.GetAgentPermissions(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// Agent-Role methods (Continued)

// Agent-Role methods (Continued)

func (ac *AgentController) AssignRoleToAgent(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	y := c.Param("roleName")
	var yArray []string
	yArray = append(yArray, y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role *models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.AssignRoleToAgent(uint(agentID), yArray)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to agent successfully"})
}

func (ac *AgentController) RevokeRoleFromAgent(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AgentService.RevokeRoleFromAgent(uint(agentID), role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role revoked from agent successfully"})
}

func (ac *AgentController) GetAgentRoles(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roles, err := ac.AgentService.GetAgentRoles(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// User-Agent methods (Continued)

func (ac *AgentController) GetTeamsByAgent(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Param("agentID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teams, err := ac.AgentService.GetTeamsByAgent(uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// User-Agent methods (Continued)//
