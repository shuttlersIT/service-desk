// backend/routes/agents.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAgentRoutes(r *gin.Engine, agent *controllers.AgentController) {

	a := r.Group("/agents")
	a.GET("/", agent.GetAllAgentsHandler)
	a.GET("/:id", agent.GetAgentByIDHandler)
	a.POST("/", agent.CreateAgentHandler)
	a.PUT("/:id", agent.UpdateAgentHandler)
	a.DELETE("/:id", agent.DeleteAgentHandler)

}
