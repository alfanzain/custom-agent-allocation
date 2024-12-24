package services

import (
	"errors"
	"fmt"

	"github.com/alfanzain/custom-agent-allocation/models"
	"gorm.io/gorm"
)

type AgentService struct {
	DB *gorm.DB
}

func NewAgentService(db *gorm.DB) *AgentService {
	return &AgentService{
		DB: db,
	}
}

func (s *AgentService) DoesAgentExist(id uint) (bool, error) {
	var agent models.Agent

	if err := s.DB.Where("id = ?", id).First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check agent existence: %w", err)
	}

	return true, nil
}

func (s *AgentService) AddAgent(id uint, name string, currentLoad int, maxLoad int) error {
	agent := models.Agent{
		ID:          id,
		Name:        name,
		CurrentLoad: currentLoad,
		MaxLoad:     maxLoad,
	}

	if err := s.DB.Create(&agent).Error; err != nil {
		return fmt.Errorf("failed to add agent: %w", err)
	}
	return nil
}

func (s *AgentService) IncreaseAgentCurrentLoad(agentID uint) error {
	var agent models.Agent

	if err := s.DB.First(&agent, agentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("agent with ID %d not found", agentID)
		}
		return fmt.Errorf("failed to fetch agent: %w", err)
	}

	if agent.CurrentLoad+1 > agent.MaxLoad {
		return fmt.Errorf("current load exceeds max load for agent %d", agentID)
	}

	agent.CurrentLoad++
	if err := s.DB.Save(&agent).Error; err != nil {
		return fmt.Errorf("failed to update agent load: %w", err)
	}

	return nil
}

func (s *AgentService) DecreaseAgentCurrentLoad(agentID uint) error {
	var agent models.Agent

	if err := s.DB.First(&agent, agentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("agent with ID %d not found", agentID)
		}
		return fmt.Errorf("failed to fetch agent: %w", err)
	}

	if agent.CurrentLoad-1 < 0 {
		return fmt.Errorf("current load exceeds zero for agent %d", agentID)
	}

	agent.CurrentLoad--
	if err := s.DB.Save(&agent).Error; err != nil {
		return fmt.Errorf("failed to update agent load: %w", err)
	}

	return nil
}

func (s *AgentService) GetAgentMaxLoad(agentID uint) (int, error) {
	var agent models.Agent

	if err := s.DB.First(&agent, agentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("agent with ID %d not found", agentID)
		}
		return 0, fmt.Errorf("failed to fetch agent: %w", err)
	}

	return agent.MaxLoad, nil
}
