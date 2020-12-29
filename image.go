package gocmcapi

import (
	"context"
	"encoding/json"
)

// ImageService interface
type ImageService interface {
	List(ctx context.Context) ([]Image, error)
}

// Image object
type Image struct {
	Name string `json:"name"`
	ID   string `json:"uuid"`
	Bits int    `json:"bits"`
}

// type Images []Image

type image struct {
	client *Client
}

// List lists all images. []*Image
func (s *image) List(ctx context.Context) ([]Image, error) {
	restext, err := s.client.Get("server/templates", nil)

	/*
		var data struct {
			Images []*Image
		}
	*/

	images := make([]Image, 0)
	err = json.Unmarshal([]byte(restext), &images)
	Logs("restext")
	Logs(restext)
	Logs("images")
	Logo(images)
	return images, err
}
