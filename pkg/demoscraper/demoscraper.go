package demoscraper

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vikpe/qw-hub-api/pkg/demoscraper/qtv"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

type DemoScraper struct {
	qtvServers []qtv.Server
}

func New(servers []qtv.Server) *DemoScraper {
	return &DemoScraper{
		qtvServers: servers,
	}
}

func (s *DemoScraper) Demos() []qtv.Demo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]qtv.Demo, 0)
		errs     = make([]error, 0)
	)

	for _, qtvServer := range s.qtvServers {
		wg.Add(1)

		go func(qtvServer qtv.Server) {
			defer wg.Done()

			demoFilenames, err := qtv.GetDemoFilenames(qtvServer.Address)

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, qtvServer.Address, err)))
				return
			}

			mutex.Lock()
			for _, filename := range demoFilenames {
				demo := qtv.Demo{
					Filename: qdemo.Filename(filename),
					Server:   qtvServer,
				}
				allDemos = append(allDemos, demo)
			}
			mutex.Unlock()
		}(qtvServer)
	}

	wg.Wait()

	return allDemos
}
