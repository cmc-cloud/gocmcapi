package gocmcapi

import (
	"encoding/json"
)

// FloatingIPService interface
type FloatingIPService interface {
	Get(id string) (FloatingIP, error)
	Create(vpcID string) (OrderResponse, TaskStatus, error)
	Delete(id string) (TaskStatus, error)
}

// FloatingIP object
type FloatingIP struct {
	ID                    string `json:"id"`
	IPAddress             string `json:"ipaddress"`
	RegionName            string `json:"zonename"`
	IsSourceNat           bool   `json:"issourcenat"`
	IsStaticNat           bool   `json:"isstaticnat"`
	AssociatedNetworkID   string `json:"associatednetworkid"`
	AssociatedNetworkName string `json:"associatednetworkname"`
	NetworkID             string `json:"networkid"`
	State                 string `json:"state"`
	VPCID                 string `json:"vpcid"`
}

type floatingIP struct {
	client *Client
}

// Get FloatingIP detail
func (v *floatingIP) Get(id string) (FloatingIP, error) {
	jsonStr, err := v.client.Get("floatingip/info", map[string]string{"id": id})
	var floatingIP FloatingIP
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &floatingIP)
	}
	return floatingIP, err
}

// Delete a FloatingIP
func (v *floatingIP) Delete(id string) (TaskStatus, error) {
	return v.client.LongDeleteTask("floatingip/delete", id, nil, ShortTimeSettings)
}
func (v *floatingIP) Create(vpcID string) (OrderResponse, TaskStatus, error) {
	return v.client.Order("floatingip/create", "", map[string]interface{}{"vpc_id": vpcID}, MediumTimeSettings)
}
