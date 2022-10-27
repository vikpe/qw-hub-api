package sources

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vikpe/qw-hub-api/internal/scrape/qtv"
	"github.com/vikpe/qw-hub-api/types"
)

type QtvServerConfig struct {
	Address        string `json:"address"`
	DemoDateFormat string `json:"date_format"`
}

type QtvDemoScraper struct {
	qtvServers []types.QtvServer
}

func NewQtvDemoScraper(qtvServerConfigs []QtvServerConfig) *QtvDemoScraper {
	qtvs := make([]types.QtvServer, 0)

	for _, config := range qtvServerConfigs {
		qtvs = append(qtvs, types.QtvServer{
			Address:        config.Address,
			DemoDateFormat: config.DemoDateFormat,
		})
	}

	return &QtvDemoScraper{
		qtvServers: qtvs,
	}
}

func (s *QtvDemoScraper) Demos() []types.QtvHostedDemo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]types.QtvHostedDemo, 0)
		errs     = make([]error, 0)
	)

	for _, qtvServer := range s.qtvServers {
		wg.Add(1)

		go func(qtvServer types.QtvServer) {
			defer wg.Done()

			demoFilenames, err := qtv.GetDemoFilenames(qtvServer.Address)

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, qtvServer.Address, err)))
				return
			}

			mutex.Lock()
			for _, filename := range demoFilenames {
				demo := types.QtvHostedDemo{
					Filename: types.DemoFilename(filename),
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
