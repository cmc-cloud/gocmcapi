package gocmcapi

import (
	"encoding/json"
)

// VPCService interface
type VPCService interface {
	Get(id string) (VPC, error)
	Create(name string, description string, region string, cidr string) (OrderResponse, TaskStatus, error)
	Delete(id string) (TaskStatus, error)
	Update(id string, name string, description string) error
}

// VPC object
type VPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	RegionName  string `json:"zonename"`
	Cidr        string `json:"cidr"`
	Description string `json:"description"`
}
type vpc struct {
	client *Client
}

// Get vpc detail
func (v *vpc) Get(id string) (VPC, error) {
	jsonStr, err := v.client.Get("vpc/info", map[string]string{"id": id})
	var vpc VPC
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vpc)
	}
	return vpc, err
}

// Delete a vpc
func (v *vpc) Delete(id string) (TaskStatus, error) {
	return v.client.LongDeleteTask("vpc/delete", id, nil, MediumTimeSettings)
}
func (v *vpc) Update(id string, name string, description string) error {
	_, err := v.client.LongTask("vpc", id, map[string]interface{}{"name": name, "description": description}, ShortTimeSettings)
	return err
}
func (v *vpc) Create(name string, description string, region string, cidr string) (OrderResponse, TaskStatus, error) {
	return v.client.Order("vpc/create", "", map[string]interface{}{"name": name, "description": description, "region": region, "cidr": cidr}, ShortTimeSettings)
}
