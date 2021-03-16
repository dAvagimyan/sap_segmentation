package segmentation

import (
	"log"
	"sap_segmentation/segmentation/client"
	"sap_segmentation/segmentation/repo"
	"time"
)

func NewLoader(cl client.Client, rep repo.SegmentationRepository, logger *log.Logger) *Updater {
	return &Updater{
		client: cl,
		rep:    rep,
		logger: logger,
	}
}

type Updater struct {
	client client.Client
	rep    repo.SegmentationRepository
	logger *log.Logger
}

func (l *Updater) Load(offset, batchSize, interval int) {
	defer func() {
		if err := recover(); err != nil {
			l.logger.Println(`recover: `, err)
		}
	}()

	for {
		resp, err := l.client.GetItems(offset)
		if err != nil {
			l.logger.Println(err)
			break
		}
		if resp.IsEmpty() {
			l.logger.Println(`finish`)
			break
		}
		offset += batchSize

		if err := l.rep.InsertPackage(resp.Segmentation); err != nil {
			l.logger.Println(err)
		}

		// можно было сделать через time.ticker.
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}
