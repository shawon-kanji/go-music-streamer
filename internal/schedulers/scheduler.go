package schedulers

import "go-music-streamer/internal/domain/entity"

type Queue interface {
	Start()
	Stop()
	AddTask(task *entity.Song)
}

type Scheduler struct {
	TagGenerator Queue
}

func InitiateSchedulers(s *Scheduler) {
	if s == nil {
		return
	}

	s.TagGenerator = NewTagGenerator()
	s.TagGenerator.Start()

}
