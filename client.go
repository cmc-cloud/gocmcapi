package gocmcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultAPIURL = "https://api.cloud.cmctelecom.vn/ver2"
)

var (
	// ErrNotFound for resource not found status
	ErrNotFound = errors.New("Resource not found")
	// ErrPermissionDenied for permission denied
	ErrPermissionDenied = errors.New("You are not allowed to do this action")
	// ErrCommon for common error
	ErrCommon = errors.New("Error")
)

// OrderResponse response when create a Server
type OrderResponse struct {
	//Success    bool   `json:"success"`
	// ID     string `json:"id"`
	TaskID string `json:"jobid"`
	Price  int    `json:"price"`
	Paid   bool   `json:"paid"`
}

// Client represents CMC Cloud API client.
type Client struct {
	apiURL         string
	apiKey         string
	Server         ServerService
	Task           TaskService
	Volume         VolumeService
	VPC            VPCService
	Network        NetworkService
	FloatingIP     FloatingIPService
	FirewallDirect FirewallDirectService
	FirewallVPC    FirewallVPCService
	Snapshot       SnapshotService
}

// APIError is return when there are an error when call api
type APIError struct {
	Success   bool   `json:"success"`
	ErrorCode int    `json:"error_code"`
	ErrorText string `json:"error_text"`
}

// Timeout is timeout info for a long task
type Timeout struct {
	Delay      time.Duration `default:"geek"` // Wait this time before starting checks
	Timeout    time.Duration // The amount of time to wait before timeout
	MinTimeout time.Duration // Smallest time to wait before refreshes
}

// NewClient creates new CMC Cloud Api client.
func NewClient(apikey string) (*Client, error) {
	c := &Client{
		apiURL: defaultAPIURL,
		apiKey: apikey,
	}
	c.Server = &server{client: c}
	c.Task = &task{client: c}
	c.Volume = &volume{client: c}
	c.FloatingIP = &floatingIP{client: c}
	c.VPC = &vpc{client: c}
	c.Network = &network{client: c}
	c.FirewallDirect = &firewalldirect{client: c}
	c.FirewallVPC = &firewallvpc{client: c}
	c.Snapshot = &snapshot{client: c}
	return c, nil
}

func (c *Client) createRequest(params map[string]string) *resty.Request {
	client := resty.New()

	if params == nil {
		params = make(map[string]string)
	}

	//c.apiKey = "vTMSG7F9mFKnNRYIz8eA9N9XrHJ4zP"
	params["api_key"] = c.apiKey

	//var obj interface{}
	request := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(c.apiKey).
		SetError(&APIError{}).
		SetQueryParams(params)

	return request
}
func (c *Client) parseResponse(response *resty.Response, err error) (string, error) {
	restext := response.String() // fmt.Sprint(response)
	if err != nil {
		return restext, err
	}
	if response.Error() != nil {
		apiError := response.Error().(*APIError)
		if apiError != nil {
			if apiError.ErrorCode == 0 {
				apiError.ErrorCode = response.StatusCode()
			}
			return restext, fmt.Errorf("Error %d: %s", apiError.ErrorCode, apiError.ErrorText)
		}
	}

	if strings.Contains(restext, "error_code") && strings.Contains(restext, "error_text") {
		var apiError APIError
		json.Unmarshal([]byte(restext), &apiError)
		return restext, fmt.Errorf("Error %d: %s", apiError.ErrorCode, apiError.ErrorText)
	}
	return restext, err
}

// Get Request, return resty Response
func (c *Client) Get(path string, params map[string]string) (string, error) {
	resp, err := c.createRequest(params).Get(c.apiURL + "/" + path + ".json")

	restext := fmt.Sprint(resp)
	return restext, err
}

// Post request
func (c *Client) Post(path string, params map[string]interface{}) (string, error) {
	// fmt.Println(c.apiURL+"/"+path+".json", params)
	resp, err := c.createRequest(nil).SetBody(params).Post(c.apiURL + "/" + path + ".json")
	return c.parseResponse(resp, err)
}

// Put request
func (c *Client) Put(path string, params map[string]interface{}) (string, error) {
	resp, err := c.createRequest(nil).SetBody(params).Put(c.apiURL + "/" + path + ".json")

	restext := fmt.Sprint(resp)
	return restext, err
}

// Delete request
func (c *Client) Delete(path string, params map[string]string) (string, error) {
	resp, err := c.createRequest(params).Delete(c.apiURL + "/" + path + ".json")

	restext := fmt.Sprint(resp)
	return restext, err
}

// LongTask execute a action that return a task
func (c *Client) LongTask(action string, id string, params map[string]interface{}, timeSettings TimeSettings) (TaskStatus, error) {
	if params == nil {
		params = make(map[string]interface{})
	}
	if id != "" {
		params["id"] = id
	}

	jsonStr, err := c.Post(action, params)
	var task Task
	var taskResponse TaskStatus
	json.Unmarshal([]byte(jsonStr), &task)

	if err != nil {
		return taskResponse, err
	}
	taskResponse, err = c.waitForTaskFinished(task.TaskID, timeSettings)
	if err != nil {
		return taskResponse, fmt.Errorf("Error perform action %s: %s, params: %+v", action, err, params)
	}
	return taskResponse, err
}

// LongDeleteTask execute a action that return a task
func (c *Client) LongDeleteTask(action string, id string, params map[string]string, timeSettings TimeSettings) (TaskStatus, error) {
	if params == nil {
		params = make(map[string]string)
	}
	if id != "" {
		params["id"] = id
	}

	jsonStr, err := c.Delete(action, params)
	var task Task
	var taskResponse TaskStatus
	json.Unmarshal([]byte(jsonStr), &task)

	if err != nil {
		return taskResponse, err
	}
	taskResponse, err = c.waitForTaskFinished(task.TaskID, timeSettings)
	if err != nil {
		return taskResponse, fmt.Errorf("Error perform action %s: %s, params: %+v", action, err, params)
	}
	return taskResponse, err
}

// Order create an resource order
func (c *Client) Order(action string, id string, params map[string]interface{}, timeSettings TimeSettings) (OrderResponse, TaskStatus, error) {
	if params == nil {
		params = make(map[string]interface{})
	}
	if id != "" {
		params["id"] = id
	}

	jsonStr, err := c.Post(action, params)
	var order OrderResponse
	var taskStatus TaskStatus

	if err != nil {
		return order, taskStatus, fmt.Errorf("Error perform action %s: %s, params: %+v", action, err, params)
	}

	json.Unmarshal([]byte(jsonStr), &order)
	if !order.Paid {
		return order, taskStatus, fmt.Errorf("Error perform action %s cause order is not paid, input = %+v, response = %s", action, params, jsonStr)
		//errors.New("Can not perform this action cause of payment failed, connect to CMC administrator for your advice")
	}

	taskStatus, err = c.waitForTaskFinished(order.TaskID, timeSettings)
	if err != nil {
		return order, taskStatus, fmt.Errorf("Error perform action %s with task id (%s): %s", action, order.TaskID, err)
	}

	return order, taskStatus, err
}

// TimeSettings object
type TimeSettings struct {
	Delay    int
	Interval int
	Timeout  int
}

// ShortTimeSettings predefined TimeSettings for short task
var ShortTimeSettings = TimeSettings{Delay: 1, Interval: 1, Timeout: 60}

// MediumTimeSettings predefined TimeSettings for medium task
var MediumTimeSettings = TimeSettings{Delay: 3, Interval: 3, Timeout: 5 * 60}

// LongTimeSettings predefined TimeSettings for long task
var LongTimeSettings = TimeSettings{Delay: 10, Interval: 20, Timeout: 20 * 60}

// SuperLongTimeSettings predefined TimeSettings for long task
var SuperLongTimeSettings = TimeSettings{Delay: 20, Interval: 20, Timeout: 5 * 60 * 60}

// HalfDayTimeSettings for long task like take snapshot
var HalfDayTimeSettings = TimeSettings{Delay: 60, Interval: 60, Timeout: 12 * 60 * 60}

// OneDayTimeSettings for long task like take snapshot
var OneDayTimeSettings = TimeSettings{Delay: 60, Interval: 60, Timeout: 24 * 60 * 60}

func (c *Client) waitForTaskFinished(taskID string, timeSettings TimeSettings) (TaskStatus, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", taskID)
	stateConf := &StateChangeConf{
		Pending:    []string{"WAIT", "PROCESSING"},
		Target:     []string{"DONE"},
		Refresh:    c.taskStateRefreshfunc(taskID),
		Timeout:    time.Duration(timeSettings.Timeout) * time.Second,
		Delay:      time.Duration(timeSettings.Delay) * time.Second,
		MinTimeout: time.Duration(timeSettings.Interval) * time.Second,
	}
	res, err := stateConf.WaitForState()
	if err != nil {
		return TaskStatus{}, err
	}
	return res.(TaskStatus), err
}

func (c *Client) taskStateRefreshfunc(taskID string) StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := c.Task.Get(taskID)
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if resp.Status == "ERROR" {
			log.Println("[DEBUG] Task is failed")
			return nil, "", errors.New(fmt.Sprint(resp))
		}

		if resp.Status == "DONE" {
			return resp, "DONE", nil
		}

		log.Println("[DEBUG] Task is not done")
		return nil, "", nil
	}
}
