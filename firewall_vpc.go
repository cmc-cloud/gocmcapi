package gocmcapi

import (
	"encoding/json"
	"fmt"
)

// FirewallVPCService interface
type FirewallVPCService interface {
	Get(id string) (FirewallVPC, error)
	Create(vpcID, name, description string) (TaskStatus, error)
	Delete(id string) (TaskStatus, error)
	Update(id string, name string, description string) error
	DeleteAllRules(id string) error
	CreateRule(id string, number int, cidrs string, action string, protocol string, ruleType string, portRange string) (TaskStatus, error)
	UpdateRule(id string, number int, cidrs string, action string, protocol string, ruleType string, portRange string) (TaskStatus, error)
	GetRules(id string) ([]interface{}, error)
	ValidateRules(inboundRules string, outboundRules string) ([]string, error)
	SaveRules(id string, inboundRules string, outboundRules string) (TaskStatus, error)
}

// FirewallVPCRule object
type FirewallVPCRule struct {
	ID        string   `json:"id"`
	Protocol  string   `json:"protocol"`
	Action    string   `json:"action"`
	Cidrs     []string `json:"cidrs"`
	PortRange string   `json:"port_range"`
}

// FirewallVPC object
type FirewallVPC struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	VPCID         string            `json:"vpc_id"`
	InboundRules  []FirewallVPCRule `json:"inbound_rules"`
	OutboundRules []FirewallVPCRule `json:"outbound_rules"`
}
type firewallvpc struct {
	client *Client
}

// Get FirewallVPC detail
func (v *firewallvpc) Get(id string) (FirewallVPC, error) {
	jsonStr, err := v.client.Get("firewall_vpc/info", map[string]string{"id": id})
	var firewall FirewallVPC
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &firewall)
	}
	return firewall, err
}

// Delete a FirewallVPC
func (v *firewallvpc) Delete(id string) (TaskStatus, error) {
	return v.client.LongDeleteTask("firewall_vpc/delete", id, nil, MediumTimeSettings)
}
func (v *firewallvpc) Update(id string, name string, description string) error {
	_, err := v.client.LongTask("firewall_vpc/update", id, map[string]interface{}{"name": name, "description": description}, ShortTimeSettings)
	return err
}

/*
func (v *firewallvpc) DeleteAllRules(id string) error {
	_, err := v.client.LongTask("firewall_vpc/delete_all_rules", id, nil, MediumTimeSettings)
	return err
}
*/
func (v *firewallvpc) Create(vpcID, name, description string) (TaskStatus, error) {
	return v.client.LongTask("firewall_vpc/create", "", map[string]interface{}{"name": name, "description": description, "vpc_id": vpcID}, ShortTimeSettings)
}
func (v *firewallvpc) CreateRule(id string, number int, cidrs string, action string, protocol string, ruleType string, portRange string) (TaskStatus, error) {
	return v.client.LongTask("firewall_vpc/create_rule", id, map[string]interface{}{
		"number":     number,
		"cidrs":      cidrs,
		"action":     action,
		"protocol":   protocol,
		"type":       ruleType,
		"port_range": portRange,
	}, ShortTimeSettings)
}

func (v *firewallvpc) UpdateRule(id string, number int, cidrs string, action string, protocol string, ruleType string, portRange string) (TaskStatus, error) {
	return v.client.LongTask("firewall_vpc/update_rule", id, map[string]interface{}{
		"rule_id":    id,
		"number":     number,
		"cidrs":      cidrs,
		"action":     action,
		"protocol":   protocol,
		"type":       ruleType,
		"port_range": portRange,
	}, ShortTimeSettings)
}

func (v *firewallvpc) GetRules(id string) ([]interface{}, error) {
	jsonStr, err := v.client.Get("firewall_vpc/get_rules", map[string]string{"id": id})
	var firewalls []interface{}
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &firewalls)
	}
	return firewalls, err
}

func (v *firewallvpc) SaveRules(id string, inboundRules string, outboundRules string) (TaskStatus, error) {
	return v.client.LongTask("firewall_vpc/save_rules", id, map[string]interface{}{"inbound_rules": inboundRules, "outbound_rules": outboundRules}, MediumTimeSettings)
}

func (v *firewallvpc) DeleteAllRules(id string) error {
	rules, err := v.GetRules(id)
	if err != nil {
		return err
	}
	var tasks = make(map[string]bool)
	for _, rawRule := range rules {
		rule := rawRule.(map[string]interface{})
		ruleID := rule["id"].(string)
		jsonStr, err := v.client.Delete("firewall_vpc/delete_rule", map[string]string{"id": ruleID})
		var task Task
		json.Unmarshal([]byte(jsonStr), &task)
		if err != nil {
			return fmt.Errorf("Error when delete rule id %s: %+v", ruleID, err)
		}
		tasks[task.TaskID] = false
	}
	for taskID, finished := range tasks {
		if !finished {
			_, err = v.client.waitForTaskFinished(taskID, ShortTimeSettings)
			if err != nil {
				return fmt.Errorf("Error when delete rule, task id %s, %+v", taskID, err)
			}
		}
	}
	return nil
}

func (v *firewallvpc) ValidateRules(inboundRules string, outboundRules string) ([]string, error) {
	jsonStr, err := v.client.Post("firewall_vpc/validate_rules", map[string]interface{}{"inbound_rules": inboundRules, "outbound_rules": outboundRules})
	var errors []string
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &errors)
	}
	return errors, err
}
