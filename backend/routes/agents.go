// backend/routes/agents.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAgentRoutes(r *gin.Engine, agent *controllers.AgentController) {

	a := r.Group("/agents")
	a.GET("/", agent.GetAllAgents)
	a.GET("/:id", agent.GetAgentByID)
	a.POST("/", agent.CreateAgent)
	a.PUT("/:id", agent.UpdateAgent)
	a.DELETE("/:id", agent.DeleteAgent)

}
