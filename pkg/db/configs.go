package db

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"

	"jsonstore/pkg/model"
)

type Config interface {
	Get(string) (*model.Config, error)
	GetAll() ([]model.Config, error)
	Search(string, string) ([]model.Config, error)

	Upsert(model.Config) error
	Delete(string) error
}

type configRepo struct {
	data map[string]model.Config
}

func NewConfigRepo() Config {
	return configRepo{map[string]model.Config{}}
}

func (c configRepo) Get(name string) (*model.Config, error) {
	config, ok := c.data[name]
	if !ok {
		return nil, fmt.Errorf("config not found")
	}

	return &config, nil
}

func (c configRepo) GetAll() ([]model.Config, error) {
	if len(c.data) < 1 {
		return nil, fmt.Errorf("no configs found")
	}

	values := make([]model.Config, 0, len(c.data))
	for _, v := range c.data {
		values = append(values, v)
	}

	return values, nil
}

func (c configRepo) Search(path, value string) ([]model.Config, error) {
	var result []model.Config
	for _, v := range c.data {
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("parse stored data")
		}
		got := gjson.ParseBytes(bytes).Get(path).Str
		if got == value {
			result = append(result, v)
		}
	}

	return result, nil
}

func (c configRepo) Upsert(config model.Config) error {
	c.data[config.Name] = config
	return nil
}

func (c configRepo) Delete(name string) error {
	_, ok := c.data[name]
	if !ok {
		return fmt.Errorf("config not found")
	}

	delete(c.data, name)

	return nil
}
