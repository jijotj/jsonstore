package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"jsonstore/pkg/contract"
	"jsonstore/pkg/lib"
	"jsonstore/pkg/service"
)

func GetAllConfigs(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := mgr.GetAll()
		if err != nil {
			zap.S().Errorf("Get all configs: %v", err)
			http.Error(w, fmt.Sprintf("Get all configs: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		lib.WriteResponseJSON(w, res)
	}
}

func GetConfig(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := mux.Vars(r)["name"]
		if !ok {
			zap.S().Errorf("Config name was not provided")
			http.Error(w, "Missing config name", http.StatusBadRequest)
			return
		}

		res, err := mgr.Get(name)
		if err != nil {
			zap.S().Errorf("Get config %s: %v", name, err)
			http.Error(w, fmt.Sprintf("Get config %s: %v", name, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		lib.WriteResponseJSON(w, res)
	}
}

func SearchConfigs(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		if len(queries) != 1 {
			zap.S().Errorf("Search configs: invalid query expression")
			http.Error(w, "Search configs: invalid query expression", http.StatusBadRequest)
			return
		}
		var path, value string
		for k, v := range queries {
			path = k
			if len(v) == 0 {
				v = []string{""}
			}
			value = v[0]
		}

		res, err := mgr.Search(path, value)
		if err != nil {
			zap.S().Errorf("Search configs %s=%s: %v", path, value, err)
			http.Error(w, fmt.Sprintf("Search configs %s=%s: %v", path, value, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		lib.WriteResponseJSON(w, res)
	}
}

func CreateConfig(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request contract.UpsertConfigRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			zap.S().Errorf("Malformed request body: %v", err)
			http.Error(w, "Malformed request body", http.StatusBadRequest)
			return
		}

		err = mgr.Upsert(request)
		if err != nil {
			zap.S().Errorf("Upsert config: %v", err)
			http.Error(w, fmt.Sprintf("Create config: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateConfig(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request contract.UpsertConfigRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			zap.S().Errorf("Malformed request body: %v", err)
			http.Error(w, "Malformed request body", http.StatusBadRequest)
			return
		}

		err = mgr.Upsert(request)
		if err != nil {
			zap.S().Errorf("Upsert config: %v", err)
			http.Error(w, fmt.Sprintf("Update config: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteConfig(mgr service.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := mux.Vars(r)["name"]
		if !ok {
			zap.S().Errorf("Config name was not provided")
			http.Error(w, "Missing config name", http.StatusBadRequest)
			return
		}

		err := mgr.Delete(name)
		if err != nil {
			zap.S().Errorf("Delete config %s: %v", name, err)
			http.Error(w, fmt.Sprintf("Delete config %s: %v", name, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
