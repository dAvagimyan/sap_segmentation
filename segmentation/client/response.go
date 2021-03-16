package client

import "sap_segmentation/model"

type Response struct {
	Segmentation []model.Segmentation
	EndPoint     string `json:"-"`
}

func (r *Response) IsEmpty() bool {
	return r.Segmentation == nil || len(r.Segmentation) == 0
}
