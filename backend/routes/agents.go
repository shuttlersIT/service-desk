// backend/routes/agents.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAgentRoutes(router *gin.Engine, agentController *controllers.AgentController) {

	// Define Agent routes
	agentRoutes := router.Group("/agents")
	{
		agentRoutes.POST("/", agentController.CreateAgentHandler)
		agentRoutes.PUT("/:id", agentController.UpdateAgentHandler)
		agentRoutes.GET("/:id", agentController.GetAgentByIDHandler)
		agentRoutes.DELETE("/:id", agentController.DeleteAgentHandler)
		agentRoutes.GET("/", agentController.GetAllAgentsHandler)
		agentRoutes.POST("/units", agentController.CreateUnitHandler)
		agentRoutes.PUT("/units/:id", agentController.UpdateUnitHandler)
		agentRoutes.DELETE("/units/:id", agentController.DeleteUnitHandler)
		agentRoutes.GET("/units", agentController.GetUnitsHandler)
		agentRoutes.GET("/units/:number", agentController.GetUnitByNumberHandler)
		agentRoutes.POST("/teams", agentController.CreateTeamHandler)
		agentRoutes.PUT("/teams/:id", agentController.UpdateTeamHandler)
		agentRoutes.DELETE("/teams/:id", agentController.DeleteTeamHandler)
		agentRoutes.GET("/teams", agentController.GetTeamsHandler)
		agentRoutes.GET("/teams/:number", agentController.GetTeamByNumberHandler)
		agentRoutes.POST("/roles", agentController.CreateRoleHandler)
		agentRoutes.PUT("/roles/:id", agentController.UpdateRoleHandler)
		agentRoutes.DELETE("/roles/:id", agentController.DeleteRoleHandler)
		agentRoutes.GET("/roles", agentController.GetRolesHandler)
		agentRoutes.GET("/roles/:number", agentController.GetRoleByNumberHandler)
	}

	// Define Team routes
	teamRoutes := router.Group("/teams")
	{
		teamRoutes.POST("/", agentController.CreateTeamHandler)
		teamRoutes.PUT("/:id", agentController.UpdateTeamHandler)
		teamRoutes.GET("/:id", agentController.GetTeamByIDHandler)
		teamRoutes.DELETE("/:id", agentController.DeleteTeamHandler)
		teamRoutes.GET("/", agentController.GetTeamsHandler)
	}

}
