package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jsonstore/pkg/contract"
	"jsonstore/pkg/model"
	"jsonstore/pkg/testlib/mocks"
)

var (
	dc1 = `{
    "name": "datacenter-1",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "false",
          "value": "300m"
        }
      }
    }
  }`
	dc2 = `{
    "name": "datacenter-2",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "true",
          "value": "250m"
        }
      }
    }
  }`
	dc1Item      = model.Config{Name: "datacenter-1", Metadata: []byte(dc1)}
	dc2Item      = model.Config{Name: "datacenter-2", Metadata: []byte(dc2)}
	dc1Data      = contract.Config{Name: "datacenter-1", Metadata: []byte(dc1)}
	dc2Data      = contract.Config{Name: "datacenter-2", Metadata: []byte(dc2)}
	dc1GetResp   = contract.GetConfigResponse{Config: dc1Data}
	dc2GetResp   = contract.GetConfigResponse{Config: dc2Data}
	dc1CreateReq = contract.UpsertConfigRequest{Config: dc1Data}
)

func TestGet(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Get", "datacenter-1").Return(&dc1Item, nil)

	config, err := manager.Get("datacenter-1")

	assert.NoError(t, err, "Unexpected get config error")
	assert.Equal(t, &dc1GetResp, config, "Incorrect config value")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestGetForRepoError(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Get", "datacenter-1").Return(nil, errors.New("some error"))

	_, err := manager.Get("datacenter-1")

	assert.Errorf(t, err, "Missing get config error")
	assert.Contains(t, err.Error(), "select:", "Incorrect get config error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestGetAll(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	all := []model.Config{dc1Item, dc2Item}
	configRepo.On("GetAll").Return(all, nil)

	configs, err := manager.GetAll()

	assert.NoError(t, err, "Unexpected get all configs error")
	assert.Equal(t, []contract.GetConfigResponse{dc1GetResp, dc2GetResp}, configs, "Incorrect configs value")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestGetAllForRepoError(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("GetAll").Return(nil, errors.New("some error"))

	_, err := manager.GetAll()

	assert.Errorf(t, err, "Missing get all configs error")
	assert.Contains(t, err.Error(), "select all:", "Incorrect get all configs error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestSearch(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	all := []model.Config{dc1Item, dc2Item}
	configRepo.On("Search", "metadata.monitoring.enabled", "true").Return(all, nil)

	configs, err := manager.Search("metadata.monitoring.enabled", "true")

	assert.NoError(t, err, "Unexpected get all configs error")
	assert.Equal(t, []contract.GetConfigResponse{dc1GetResp, dc2GetResp}, configs, "Incorrect configs value")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestSearchForRepoError(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Search", "metadata.monitoring.enabled", "true").Return(nil, errors.New("some error"))

	_, err := manager.Search("metadata.monitoring.enabled", "true")

	assert.Errorf(t, err, "Missing get all configs error")
	assert.Contains(t, err.Error(), "search:", "Incorrect get all configs error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestUpsert(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Upsert", dc1Item).Return(nil)

	err := manager.Upsert(dc1CreateReq)

	assert.NoError(t, err, "Unexpected upsert config error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestUpsertForRepoError(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Upsert", dc1Item).Return(errors.New("some error"))

	err := manager.Upsert(dc1CreateReq)

	assert.Errorf(t, err, "Missing upsert config error")
	assert.Contains(t, err.Error(), "insert:", "Incorrect upsert config error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestDelete(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Delete", "datacenter-1").Return(nil)

	err := manager.Delete("datacenter-1")

	assert.NoError(t, err, "Unexpected delete config error")
	mock.AssertExpectationsForObjects(t, configRepo)
}

func TestDeleteForRepoError(t *testing.T) {
	configRepo := new(mocks.Config)
	manager := NewConfigManager(configRepo)

	configRepo.On("Delete", "datacenter-1").Return(errors.New("some error"))

	err := manager.Delete("datacenter-1")

	assert.Errorf(t, err, "Missing delete config error")
	assert.Contains(t, err.Error(), "delete:", "Incorrect delete config error")
	mock.AssertExpectationsForObjects(t, configRepo)
}
