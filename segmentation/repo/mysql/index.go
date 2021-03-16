package mysql

import (
	"github.com/jmoiron/sqlx"
	"sap_segmentation/model"
)

func NewSegmentationRepository(conn *sqlx.DB) *SegmentationRepository {
	return &SegmentationRepository{
		conn: conn,
	}
}

type SegmentationRepository struct {
	conn *sqlx.DB
}

//	В r.conn.NamedExec - можно передать массив, по документации но я наткнулся на косяк драйвера и так и не смог за это время его переюороть.
// Еще Можно было воспользоваться одним запросом через placeholder  или собирать запрос самому в ручную.
func (r *SegmentationRepository) InsertPackage(items []model.Segmentation) error {
	tx := r.conn.MustBegin()
	for _, item := range items {
		tx.NamedExec(`INSERT INTO segmentations (address_sap_id, adr_segment, segment_id)
       VALUES (:address_sap_id, :adr_segment, :segment_id ) ON DUPLICATE KEY UPDATE adr_segment = values(adr_segment), segment_id = values(segment_id)`, item)
	}
	return tx.Commit()
}
