package model

type Segmentation struct {
	ID           int64  `db:"id" json:"-"`
	AddressSapID string `db:"address_sap_id" json:"ADDRESS_SAP_ID"`
	AdrSegment   string `db:"adr_segment" json:"ADR_SEGMENT"`
	SegmentID    int64  `db:"segment_id" json:"SEGMENT_ID"`
}
