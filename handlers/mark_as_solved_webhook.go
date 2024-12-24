package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfanzain/custom-agent-allocation/config"
	"github.com/alfanzain/custom-agent-allocation/services"
	"github.com/alfanzain/custom-agent-allocation/types"
	"github.com/labstack/echo/v4"
)

type MarkAsSolvedHandler struct {
	QiscusService *services.QiscusService
	QueueService  *services.QueueService
	AgentService  *services.AgentService
}

func NewMarkAsSolvedHandler(
	qiscusService *services.QiscusService,
	queueService *services.QueueService,
	agentService *services.AgentService,
) *MarkAsSolvedHandler {
	return &MarkAsSolvedHandler{
		QiscusService: qiscusService,
		QueueService:  queueService,
		AgentService:  agentService,
	}
}

func (h *MarkAsSolvedHandler) MarkAsSolvedWebhook(c echo.Context) error {
	var payload types.MarkAsResolvedWebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("[Mark As Solved Webhook] Error marshaling payload to JSON: %v", err)
	} else {
		log.Printf("[Mark As Solved Webhook] Payload received: \n%s\n\n", string(jsonData))
	}

	err = h.AgentService.DecreaseAgentCurrentLoad(uint(payload.ResolvedBy.ID))
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error decreasing agent current load: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decrease agent current load"})
	}

	queueExists, err := h.QueueService.DoesQueueExists(config.REDIS_QUEUE_CUSTOMERS_KEY)
	if err != nil {
		log.Printf("[Mark As Solved Webhook] Error checking queue existence: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check queue existence"})
	}

	if !queueExists {
		return c.JSON(http.StatusOK, map[string]string{"message": "Customer resolved successfully"})
	}

	resp, err := h.QiscusService.AllocateAgent()
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error allocating agent: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to allocate agent"})
	}

	agentExists, err := h.AgentService.DoesAgentExist(uint(resp.Data.Agent.ID))
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error checking agent existence: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check agent existence"})
	}

	if !agentExists {
		err = h.AgentService.AddAgent(uint(resp.Data.Agent.ID), resp.Data.Agent.Name, resp.Data.Agent.Count, config.AGENT_DEFAULT_MAX_LOAD)
		if err != nil {
			log.Printf("[Allocate Agent Webhook] Error adding agent: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add agent"})
		}
	}

	agentCurrentLoad := resp.Data.Agent.Count
	if agentCurrentLoad >= config.AGENT_DEFAULT_MAX_LOAD {
		log.Printf("[Allocate Agent Webhook] Agent full. Queueing...")
		return c.JSON(http.StatusOK, map[string]string{"message": "Customer queued successfully"})
	}

	nextRoomId, err := h.QueueService.DequeueCustomer(config.REDIS_QUEUE_CUSTOMERS_KEY)
	if err != nil {
		log.Printf("[Mark As Solved Webhook] Error dequeuing customer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to dequeue customer"})
	}

	_, err = h.QiscusService.AssignAgent(nextRoomId, uint(resp.Data.Agent.ID))
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error assigning agent to customer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign agent to customer"})
	}

	err = h.AgentService.IncreaseAgentCurrentLoad(uint(resp.Data.Agent.ID))
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error increasing agent current load: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to increase agent current load"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Customer resolved successfully"})
}