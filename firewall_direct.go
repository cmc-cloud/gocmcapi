package gocmcapi

import (
	"encoding/json"
)

// FirewallDirectService interface
type FirewallDirectService interface {
	Get(serverID string, ipAddress string) (FirewallDirect, error)
	Delete(serverID string, ipAddress string) (TaskStatus, error)
	SaveRules(erverID string, ipAddress string, inboundRules string, outboundRules string) (TaskStatus, error)
}

// FirewallDirectRule object
type FirewallDirectRule struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
	Src       string `json:"src"`
	Dst       string `json:"dst"`
	PortRange string `json:"port_range"`
	Action    string `json:"action"`
	PortType  string `json:"port_type"`
}

// FirewallDirect object
type FirewallDirect struct {
	ServerID      string               `json:"server_id"`
	IPAddress     string               `json:"ip_address"`
	InboundRules  []FirewallDirectRule `json:"inbound_rules"`
	OutboundRules []FirewallDirectRule `json:"outbound_rules"`
}

type firewalldirect struct {
	client *Client
}

// Get FirewallDirect detail
func (v *firewalldirect) Get(serverID string, ipAddress string) (FirewallDirect, error) {
	jsonStr, err := v.client.Get("firewall_direct/info", map[string]string{"server_id": serverID, "ip_address": ipAddress})
	var firewall FirewallDirect
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &firewall)
	}
	return firewall, err
}

// Delete a FirewallDirect
func (v *firewalldirect) Delete(serverID string, ipAddress string) (TaskStatus, error) {
	return v.client.LongDeleteTask("firewall_direct/delete", serverID, map[string]string{"ip_address": ipAddress}, MediumTimeSettings)
}

func (v *firewalldirect) SaveRules(serverID string, ipAddress string, inboundRules string, outboundRules string) (TaskStatus, error) {
	return v.client.LongTask("firewall_direct/save_rules", "", map[string]interface{}{"server_id": serverID, "ip_address": ipAddress, "inbound_rules": inboundRules, "outbound_rules": outboundRules}, MediumTimeSettings)
}
