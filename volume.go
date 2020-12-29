package gocmcapi

import (
	"encoding/json"
)

// VolumeService interface
type VolumeService interface {
	Get(id string) (Volume, error)
	Create(params map[string]interface{}) (OrderResponse, TaskStatus, error)
	Delete(id string) (TaskStatus, error)
	Resize(id string, newSize int) (OrderResponse, TaskStatus, error)
	Rename(id string, newName string) error
	Attach(id string, serverID string) (string, error)
	Detach(id string) (string, error)
}

// Volume object
type Volume struct {
	ID       string `json:"uuid"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Created  string `json:"created"`
	ServerID string `json:"server_id"`
}

type volume struct {
	client *Client
}

// Get volume detail
func (v *volume) Get(id string) (Volume, error) {
	jsonStr, err := v.client.Get("volume/info", map[string]string{"id": id})
	var volume Volume
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &volume)
	}
	return volume, err
}

// Delete a volume
func (v *volume) Delete(id string) (TaskStatus, error) {
	return v.client.LongDeleteTask("volume/delete", id, nil, MediumTimeSettings)
}
func (v *volume) Rename(id string, newName string) error {
	_, err := v.client.Post("volume/rename", map[string]interface{}{"id": id, "name": newName})
	return err
}
func (v *volume) Resize(id string, newSize int) (OrderResponse, TaskStatus, error) {
	return v.client.Order("volume/resize", id, map[string]interface{}{"size": newSize}, LongTimeSettings)
}
func (v *volume) Attach(id string, serverID string) (string, error) {
	return v.client.Post("volume/attach", map[string]interface{}{"id": id, "server_id": serverID})
}
func (v *volume) Detach(id string) (string, error) {
	return v.client.Post("volume/detach", map[string]interface{}{"id": id})
}

// Create a new volume
func (v *volume) Create(params map[string]interface{}) (OrderResponse, TaskStatus, error) {
	return v.client.Order("volume/create", "", params, LongTimeSettings)
}
