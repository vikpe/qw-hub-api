package qtvscraper

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

type Scraper struct {
	servers []Server
}

func NewScraper(servers []Server) *Scraper {
	return &Scraper{
		servers: servers,
	}
}

func (s *Scraper) Demos() []Demo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]Demo, 0)
		errs     = make([]error, 0)
	)

	for _, qtvServer := range s.servers {
		wg.Add(1)

		go func(server Server) {
			defer wg.Done()

			demoFilenames, err := server.DemoFilenames()

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, server.Address, err)))
				return
			}

			mutex.Lock()
			for _, filename := range demoFilenames {
				demo := Demo{
					Filename: qdemo.Filename(filename),
					Server:   server,
				}
				allDemos = append(allDemos, demo)
			}
			mutex.Unlock()
		}(qtvServer)
	}

	wg.Wait()

	return allDemos
}
