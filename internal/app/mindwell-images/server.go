package images

import (
	"database/sql"
	"go.uber.org/zap"
	"log"

	"github.com/sevings/mindwell-server/utils"

	"github.com/zpatrick/go-config"
)

type MindwellImages struct {
	cfg     *config.Config
	db      *sql.DB
	log     *zap.Logger
	acts    chan ImageProcessor
	stop    chan bool
	folder  string
	baseURL string
}

func NewMindwellImages(cfg *config.Config) *MindwellImages {
	logger, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Println(err)
	}

	mi := &MindwellImages{
		cfg:  cfg,
		log:  logger,
		acts: make(chan ImageProcessor, 50),
		stop: make(chan bool),
	}

	_, err = zap.RedirectStdLogAt(mi.LogSystem(), zap.ErrorLevel)
	if err != nil {
		mi.LogSystem().Error(err.Error())
	}

	mi.db = utils.OpenDatabase(cfg)
	mi.baseURL = mi.ConfigString("images.base_url")
	mi.folder = mi.ConfigString("images.folder")

	go func() {
		for act := range mi.acts {
			act.Work()
		}
		mi.stop <- true
	}()

	return mi
}

func (mi *MindwellImages) ConfigString(key string) string {
	value, err := mi.cfg.String(key)
	if err != nil {
		mi.LogSystem().Warn(err.Error())
	}

	return value
}

func (mi *MindwellImages) ConfigBytes(key string) []byte {
	return []byte(mi.ConfigString(key))
}

func (mi *MindwellImages) TokenHash() utils.TokenHash {
	return utils.NewTokenHash(mi)
}

func (mi *MindwellImages) Folder() string {
	return mi.folder
}

func (mi *MindwellImages) BaseURL() string {
	return mi.baseURL
}

func (mi *MindwellImages) DB() *sql.DB {
	return mi.db
}

func (mi *MindwellImages) LogApi() *zap.Logger {
	return mi.log.With(zap.String("type", "api"))
}

func (mi *MindwellImages) LogImages() *zap.Logger {
	return mi.log.With(zap.String("type", "images"))
}

func (mi *MindwellImages) LogSystem() *zap.Logger {
	return mi.log.With(zap.String("type", "system"))
}

func (mi *MindwellImages) QueueAction(is *imageStore, ID int64, action string) {
	mi.LogImages().Info("queue",
		zap.Int64("id", ID),
		zap.String("action", action),
		zap.String("path", is.FileName()),
	)

	mi.acts <- ImageProcessor{is: is, ID: ID, act: action, mi: mi}
}

func (mi *MindwellImages) Shutdown() {
	close(mi.acts)
	<-mi.stop
}
