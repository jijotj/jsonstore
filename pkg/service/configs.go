package service

import (
	"fmt"

	"jsonstore/pkg/contract"
	"jsonstore/pkg/db"
	"jsonstore/pkg/model"
)

type Manager interface {
	Get(string) (*contract.GetConfigResponse, error)
	GetAll() ([]contract.GetConfigResponse, error)
	Search(string, string) ([]contract.GetConfigResponse, error)

	Upsert(contract.UpsertConfigRequest) error
	Delete(string) error
}

type configManager struct {
	configRepo db.Config
}

func NewConfigManager(configRepo db.Config) Manager {
	return configManager{configRepo}
}

func (c configManager) Get(name string) (*contract.GetConfigResponse, error) {
	item, err := c.configRepo.Get(name)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	config := contract.Config{
		Name:     item.Name,
		Metadata: item.Metadata,
	}
	return &contract.GetConfigResponse{Config: config}, nil
}

func (c configManager) GetAll() ([]contract.GetConfigResponse, error) {
	all, err := c.configRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("select all: %w", err)
	}

	resp := make([]contract.GetConfigResponse, 0, len(all))
	for _, item := range all {
		config := contract.Config{
			Name:     item.Name,
			Metadata: item.Metadata,
		}
		resp = append(resp, contract.GetConfigResponse{Config: config})
	}
	return resp, nil
}

func (c configManager) Search(path, value string) ([]contract.GetConfigResponse, error) {
	all, err := c.configRepo.Search(path, value)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	resp := make([]contract.GetConfigResponse, 0, len(all))
	for _, item := range all {
		config := contract.Config{
			Name:     item.Name,
			Metadata: item.Metadata,
		}
		resp = append(resp, contract.GetConfigResponse{Config: config})
	}
	return resp, err
}

func (c configManager) Upsert(req contract.UpsertConfigRequest) error {
	item := model.Config{
		Name:     req.Config.Name,
		Metadata: req.Config.Metadata,
	}
	err := c.configRepo.Upsert(item)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (c configManager) Delete(name string) error {
	err := c.configRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
