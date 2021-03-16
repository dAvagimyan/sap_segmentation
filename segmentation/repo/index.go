package repo

import "sap_segmentation/model"

type SegmentationRepository interface {
	InsertPackage(items []model.Segmentation) error
}
