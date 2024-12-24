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

type AllocateAgentHandler struct {
	QiscusService *services.QiscusService
	QueueService  *services.QueueService
	AgentService  *services.AgentService
}

func NewAllocateAgentHandler(
	qiscusService *services.QiscusService,
	queueService *services.QueueService,
	agentService *services.AgentService,
) *AllocateAgentHandler {
	return &AllocateAgentHandler{
		QiscusService: qiscusService,
		QueueService:  queueService,
		AgentService:  agentService,
	}
}

func (h *AllocateAgentHandler) AllocateAgentWebhook(c echo.Context) error {
	var payload types.AllocateAgentWebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error marshaling payload to JSON: %v", err)
	} else {
		log.Printf("[Allocate Agent Webhook] Payload received: \n%s\n\n", string(jsonData))
	}

	err = h.QueueService.EnqueueCustomer(config.REDIS_QUEUE_CUSTOMERS_KEY, payload.RoomID)
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error queuing customer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to queue customer"})
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
	maxLoad, err := h.AgentService.GetAgentMaxLoad(uint(resp.Data.Agent.ID))
	if err != nil {
		log.Printf("[Allocate Agent Webhook] Error getting agent max load: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get agent max load"})
	}

	if agentCurrentLoad >= maxLoad {
		log.Printf("[Allocate Agent Webhook] Agent load exceeds max (%d). Queueing...", maxLoad)
		return c.JSON(http.StatusOK, map[string]string{"message": "Customer queued successfully"})
	} else {
		log.Printf("[Allocate Agent Webhook] Agent load is under max (%d). Assigning agent...", maxLoad)
	}

	for agentCurrentLoad < maxLoad {
		queueExists, err := h.QueueService.DoesQueueExists(config.REDIS_QUEUE_CUSTOMERS_KEY)
		if err != nil {
			log.Printf("[Allocate Agent Webhook] Error checking queue existence: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check queue existence"})
		}

		if !queueExists {
			log.Println("[Allocate Agent Webhook] Queue no longer exists. Exiting.")
			break
		}

		nextRoomId, err := h.QueueService.DequeueCustomer(config.REDIS_QUEUE_CUSTOMERS_KEY)
		if err != nil {
			log.Printf("[Allocate Agent Webhook] Error dequeuing customer: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to dequeue customer"})
		}

		log.Printf("[Allocate Agent Webhook] Assigning room ID '%s' to agent ID '%d'", nextRoomId, resp.Data.Agent.ID)

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

		agentCurrentLoad++
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Customer queued successfully"})
}
