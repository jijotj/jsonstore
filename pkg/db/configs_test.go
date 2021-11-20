package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"jsonstore/pkg/model"
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
	dc1Item = model.Config{
		Name:     "datacenter-1",
		Metadata: dc1,
	}
	dc2Item = model.Config{
		Name:     "datacenter-2",
		Metadata: dc2,
	}
)

func TestGet(t *testing.T) {
	repo := NewConfigRepo()
	err := repo.Upsert(dc1Item)
	require.NoError(t, err, "Unexpected upsert config error")

	config, err := repo.Get("datacenter-1")

	assert.NoError(t, err, "Unexpected get config error")
	assert.Equal(t, &dc1Item, config, "Incorrect config")
}

func TestGetForMissingConfig(t *testing.T) {
	repo := NewConfigRepo()

	_, err := repo.Get("datacenter-1")

	assert.Error(t, err, "Missing get config error")
	assert.Contains(t, err.Error(), "config not found", "Incorrect get config error")
}

func TestGetAll(t *testing.T) {
	repo := NewConfigRepo()
	err := repo.Upsert(dc1Item)
	require.NoError(t, err, "Unexpected upsert config error")
	err = repo.Upsert(dc2Item)
	require.NoError(t, err, "Unexpected upsert config error")

	configs, err := repo.GetAll()

	assert.NoError(t, err, "Unexpected get all configs error")
	assert.Equal(t, []model.Config{dc1Item, dc2Item}, configs, "Incorrect config")
}

func TestGetAllForNoConfigs(t *testing.T) {
	repo := NewConfigRepo()

	_, err := repo.GetAll()

	assert.Error(t, err, "Missing get all configs error")
	assert.Contains(t, err.Error(), "no configs found", "Incorrect get all configs error")
}

func TestSearch(t *testing.T) {
	repo := NewConfigRepo()
	var dc1Item model.Config
	err := json.Unmarshal([]byte(dc1), &dc1Item)
	require.NoError(t, err, "Unexpected unmarshall error")
	err = repo.Upsert(dc1Item)
	require.NoError(t, err, "Unexpected upsert config error")

	configs, err := repo.Search("metadata.monitoring.enabled", "true")

	assert.NoError(t, err, "Unexpected search config error")
	assert.Equal(t, []model.Config{dc1Item}, configs, "Incorrect config")
}

func TestCreate(t *testing.T) {
	repo := NewConfigRepo()

	err := repo.Upsert(dc1Item)

	assert.NoError(t, err, "Unexpected upsert config error")

	config, err := repo.Get("datacenter-1")
	require.NoError(t, err, "Unexpected get config error")
	assert.Equal(t, &dc1Item, config, "Incorrect config")
}

func TestDelete(t *testing.T) {
	repo := NewConfigRepo()
	err := repo.Upsert(dc1Item)
	require.NoError(t, err, "Unexpected upsert config error")
	config, err := repo.Get("datacenter-1")
	require.NoError(t, err, "Unexpected get config error")
	assert.Equal(t, &dc1Item, config, "Incorrect config")

	err = repo.Delete("datacenter-1")

	assert.NoError(t, err, "Unexpected delete config error")

	_, err = repo.Get("datacenter-1")
	assert.Error(t, err, "Missing get config error")
	assert.Contains(t, err.Error(), "config not found", "Incorrect config")
}

func TestDeleteForMissingConfig(t *testing.T) {
	repo := NewConfigRepo()

	err := repo.Delete("datacenter-1")

	assert.Error(t, err, "Missing delete config error")
	assert.Contains(t, err.Error(), "config not found", "Incorrect config")
}
