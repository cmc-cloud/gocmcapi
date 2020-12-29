package gocmcapi

import (
	"encoding/json"
)

// ServerService interface
type ServerService interface {
	Get(id string) (Server, error)
	Create(params map[string]interface{}) (OrderResponse, TaskStatus, error)
	Delete(id string) (TaskStatus, error)
	AddSecondaryIP(id string) (OrderResponse, TaskStatus, error)
	RemoveSecondaryIP(id, ip4Address string) (TaskStatus, error)
	AddNic(id, networkID string) (TaskStatus, error)
	RemoveNic(id, nicID string) (TaskStatus, error)
	DisableBackup(id string) (TaskStatus, error)
	EnableBackup(id, intervalType, scheduleTime string) (OrderResponse, TaskStatus, error)
	DisablePrivateNetwork(id string) (TaskStatus, error)
	EnablePrivateNetwork(id string) (TaskStatus, error)
	ResetPassword(id string) (TaskStatus, error)
	Restart(id string) (TaskStatus, error)
	Stop(id string) (TaskStatus, error)
	Start(id string) (TaskStatus, error)
	RestoreSnapshot(id string, snapshotID string) (TaskStatus, error)
	TakeSnapshot(id string, name string) (OrderResponse, TaskStatus, error)
	Resize(id string, cpu, ramGb, rootGb, gpu int) (OrderResponse, TaskStatus, error)
	GetConsoleURL(id string) (string, error)
	Rename(id, newName string) (string, error)
	UpdateScheduleTime(id, intervalType, scheduleTime string) (string, error)
}

// Nic object
type Nic struct {
	ID         string      `json:"uuid"`
	IP4Address string      `json:"ip4_address"`
	IP6Address interface{} `json:"ip6_address"`
	Netmask    string      `json:"netmask"`
	Gateway    string      `json:"gateway"`
	MacAddress string      `json:"mac_address"`
	DefaultNic bool        `json:"default_nic"`
	IsVPC      bool        `json:"is_vpc"`
	IsPrivate  bool        `json:"is_private"`
	IPType     string      `json:"ip_type"`
	NetworkID  string      `json:"network_id"`
}

// Server object
type Server struct {
	ID                string        `json:"uuid"`
	Name              string        `json:"name"`
	DisplayName       string        `json:"display_name"`
	Created           string        `json:"created"`
	Bits              int           `json:"bits"`
	RegionName        string        `json:"zonename"`
	RegionID          string        `json:"zoneid"`
	State             string        `json:"state"`
	MainIPAddress     string        `json:"main_ip_address"`
	ImageName         string        `json:"image_name"`
	ImageID           string        `json:"image_uuid"`
	ImageType         string        `json:"image_type"`
	CPU               int           `json:"cpu"`
	RAM               int           `json:"ram_size"`
	Root              int           `json:"root_size"`
	GPU               int           `json:"gpu"`
	AutoBackup        bool          `json:"autobackup"`
	BackupSchedule    string        `json:"backup_schedule"`
	Nics              []Nic         `json:"nics"`
	Datadisks         []interface{} `json:"datadisks"`
	TotalDatadiskSize int           `json:"total_datadisk_size"`
	Jobs              []interface{} `json:"jobs"`
	Demo              bool          `json:"demo"`
}

type server struct {
	client *Client
}

// Get server detail
func (s *server) Get(id string) (Server, error) {
	jsonStr, err := s.client.Get("server/info", map[string]string{"id": id})
	var server Server
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &server)
	}
	return server, err
}

// Delete a server
func (s *server) Delete(id string) (TaskStatus, error) {
	return s.client.LongDeleteTask("server_action/delete", id, nil, LongTimeSettings)
}
func (s *server) Rename(id, newName string) (string, error) {
	return s.client.Post("server_action/rename", map[string]interface{}{"id": id, "name": newName})
}
func (s *server) UpdateScheduleTime(id, intervalType, scheduleTime string) (string, error) {
	return s.client.Post("server_action/update_schedule_time", map[string]interface{}{"id": id, "interval_type": intervalType, "schedule_time": scheduleTime})
}
func (s *server) AddSecondaryIP(id string) (OrderResponse, TaskStatus, error) {
	return s.client.Order("server_action/add_secondary_ip", id, nil, MediumTimeSettings)
}
func (s *server) RemoveSecondaryIP(id, ip4Address string) (TaskStatus, error) {
	return s.client.LongTask("server_action/remove_secondary_ip", id, map[string]interface{}{"ip4_address": ip4Address}, MediumTimeSettings)
}
func (s *server) AddNic(id, networkID string) (TaskStatus, error) {
	return s.client.LongTask("server_action/add_nic", id, map[string]interface{}{"network_id": networkID}, MediumTimeSettings)
}
func (s *server) RemoveNic(id, nicID string) (TaskStatus, error) {
	return s.client.LongTask("server_action/remove_nic", id, map[string]interface{}{"nic_id": nicID}, MediumTimeSettings)
}
func (s *server) DisableBackup(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/disable_backup", id, nil, ShortTimeSettings)
}
func (s *server) EnableBackup(id, intervalType, scheduleTime string) (OrderResponse, TaskStatus, error) {
	return s.client.Order("server_action/enable_backup", id, map[string]interface{}{"interval_type": intervalType, "schedule_time": scheduleTime}, ShortTimeSettings)
}
func (s *server) DisablePrivateNetwork(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/disable_private_network", id, nil, MediumTimeSettings)
}
func (s *server) EnablePrivateNetwork(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/enable_private_network", id, nil, MediumTimeSettings)
}
func (s *server) ResetPassword(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/reset_pass", id, nil, MediumTimeSettings)
}
func (s *server) Restart(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/restart", id, nil, LongTimeSettings)
}
func (s *server) Stop(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/stop", id, nil, LongTimeSettings)
}
func (s *server) Start(id string) (TaskStatus, error) {
	return s.client.LongTask("server_action/start", id, nil, LongTimeSettings)
}
func (s *server) RestoreSnapshot(id, snapshotID string) (TaskStatus, error) {
	return s.client.LongTask("server_action/restore_snapshot", id, map[string]interface{}{"snapshot_id": snapshotID}, SuperLongTimeSettings)
}
func (s *server) TakeSnapshot(id string, name string) (OrderResponse, TaskStatus, error) {
	return s.client.Order("server_action/take_snapshot", id, map[string]interface{}{"name": name}, HalfDayTimeSettings)
}
func (s *server) Resize(id string, cpu, ramGb, rootGb, gpu int) (OrderResponse, TaskStatus, error) {
	return s.client.Order("server_action/resize", id, map[string]interface{}{"cpu": cpu, "ram": ramGb, "disk": rootGb, "gpu": gpu}, LongTimeSettings)
}
func (s *server) GetConsoleURL(id string) (string, error) {
	jsonStr, err := s.client.Get("server_action/console", map[string]string{"id": id})
	type Console struct {
		URL string `json:"url"`
	}
	var console Console
	json.Unmarshal([]byte(jsonStr), &console)
	if err != nil {
		return "", err
	}
	return console.URL, err
}

// Create a new server
func (s *server) Create(params map[string]interface{}) (OrderResponse, TaskStatus, error) {
	return s.client.Order("server/create", "", params, LongTimeSettings)
}
