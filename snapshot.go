package gocmcapi

import (
	"encoding/json"
)

// SnapshotService interface
type SnapshotService interface {
	Get(id string) (Snapshot, error)
	Create(volumeID string, name string) (OrderResponse, TaskStatus, error)
	Delete(id string) (TaskStatus, error)
	Rename(id string, newName string) error
}

// Snapshot object
type Snapshot struct {
	ID       string `json:"uuid"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	State    string `json:"state"`
	Created  string `json:"created"`
	VolumeID string `json:"volume_id"`
	ServerID string `json:"server_id"`
}

type snapshot struct {
	client *Client
}

// Get snapshot detail
func (v *snapshot) Get(id string) (Snapshot, error) {
	jsonStr, err := v.client.Get("snapshot/info", map[string]string{"id": id})
	var snapshot Snapshot
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &snapshot)
	}
	return snapshot, err
}

// Delete a snapshot
func (v *snapshot) Delete(id string) (TaskStatus, error) {
	return v.client.LongDeleteTask("snapshot/delete", id, nil, MediumTimeSettings)
}
func (v *snapshot) Rename(id string, newName string) error {
	_, err := v.client.Post("snapshot/rename", map[string]interface{}{"id": id, "name": newName})
	return err
}

// Create a new snapshot
func (v *snapshot) Create(volumeID string, name string) (OrderResponse, TaskStatus, error) {
	return v.client.Order("snapshot/create", "", map[string]interface{}{"volume_id": volumeID, "name": name}, HalfDayTimeSettings)
}
