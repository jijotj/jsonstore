package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"jsonstore/pkg/contract"
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
	all = `[
  {
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
  },
  {
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
  }
]`
)

func TestGetConfigs(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")
	req = mux.SetURLVars(req, map[string]string{"name": "datacenter-1"})
	var dc1Data contract.GetConfigResponse
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Get", "datacenter-1").Return(&dc1Data, nil)

	GetConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	assert.JSONEq(t, dc1, string(response), "Incorrect config values")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestGetConfigsForMissingConfigNameError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")

	GetConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Missing config name", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestGetConfigsForConfigManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")
	req = mux.SetURLVars(req, map[string]string{"name": "datacenter-1"})

	manager.On("Get", "datacenter-1").Return(nil, errors.New("some error"))

	GetConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Get config datacenter-1:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestGetAllConfigs(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")
	var dc1Data contract.GetConfigResponse
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")
	var dc2Data contract.GetConfigResponse
	err = json.NewDecoder(strings.NewReader(dc2)).Decode(&dc2Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("GetAll").Return([]contract.GetConfigResponse{dc1Data, dc2Data}, nil)

	GetAllConfigs(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	assert.JSONEq(t, all, string(response), "Incorrect config values")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestGetAllConfigsForServiceManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")

	manager.On("GetAll").Return(nil, errors.New("some error"))

	GetAllConfigs(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Get all configs:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestSearchConfigs(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs/search?metadata.monitoring.enabled=true", nil)
	require.NoError(t, err, "Unexpected create request error")
	var data []contract.GetConfigResponse
	err = json.NewDecoder(strings.NewReader(all)).Decode(&data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Search", "metadata.monitoring.enabled", "true").Return(data, nil)

	SearchConfigs(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	assert.JSONEq(t, all, string(response), "Incorrect config values")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestSearchConfigsForMissingQueryParamError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/configs/search", nil)
	require.NoError(t, err, "Unexpected create request error")

	SearchConfigs(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Search configs: invalid query expression", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestSearchConfigsForConfigManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/configs/search?metadata.monitoring.enabled=true", nil)
	require.NoError(t, err, "Unexpected create request error")
	req = mux.SetURLVars(req, map[string]string{"name": "datacenter-1"})

	manager.On("Search", "metadata.monitoring.enabled", "true").Return(nil, errors.New("some error"))

	SearchConfigs(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Search configs metadata.monitoring.enabled=true:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestCreateConfig(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/configs", strings.NewReader(dc1))
	require.NoError(t, err, "Unexpected create request error")
	var dc1Data contract.UpsertConfigRequest
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Upsert", dc1Data).Return(nil)

	CreateConfig(manager).ServeHTTP(rr, req)

	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestCreateConfigForServiceManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/configs", strings.NewReader(dc1))
	require.NoError(t, err, "Unexpected create request error")
	var dc1Data contract.UpsertConfigRequest
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Upsert", dc1Data).Return(errors.New("some error"))

	CreateConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Create config:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestUpdateConfig(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/configs", strings.NewReader(dc1))
	require.NoError(t, err, "Unexpected create request error")
	var dc1Data contract.UpsertConfigRequest
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Upsert", dc1Data).Return(nil)

	UpdateConfig(manager).ServeHTTP(rr, req)

	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestUpdateConfigForServiceManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/configs", strings.NewReader(dc1))
	require.NoError(t, err, "Unexpected create request error")
	var dc1Data contract.UpsertConfigRequest
	err = json.NewDecoder(strings.NewReader(dc1)).Decode(&dc1Data)
	require.NoError(t, err, "Unexpected json decode error")

	manager.On("Upsert", dc1Data).Return(errors.New("some error"))

	UpdateConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Update config:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestDeleteConfig(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")
	req = mux.SetURLVars(req, map[string]string{"name": "datacenter-1"})

	manager.On("Delete", "datacenter-1").Return(nil)

	DeleteConfig(manager).ServeHTTP(rr, req)

	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusOK, rr.Code, "Incorrect http status code")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestDeleteConfigForMissingConfigNameError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")

	DeleteConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Missing config name", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}

func TestDeleteConfigForServiceManagerError(t *testing.T) {
	manager := new(mocks.Manager)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/configs", nil)
	require.NoError(t, err, "Unexpected create request error")
	req = mux.SetURLVars(req, map[string]string{"name": "datacenter-1"})

	manager.On("Delete", "datacenter-1").Return(errors.New("some error"))

	DeleteConfig(manager).ServeHTTP(rr, req)

	response, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err, "Unexpected error while reading response body")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Incorrect http status code")
	assert.Contains(t, string(response), "Delete config datacenter-1:", "Incorrect response")
	mock.AssertExpectationsForObjects(t, manager)
}
