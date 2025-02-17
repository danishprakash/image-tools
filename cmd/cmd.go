package command

import (
	"fmt"
	"sync"

	"cnrancher.io/image-tools/mirror"
	"cnrancher.io/image-tools/registry"
	u "cnrancher.io/image-tools/utils"
	"github.com/sirupsen/logrus"
)

// PrepareDockerLogin executes docker login command if
// SRC_USERNAME/SRC_PASSWORD, DEST_USERNAME/DEST_PASSWORD
// SRC_REGISTRY, DEST_REGISTRY
// environment variables are set
func ProcessDockerLoginEnv() error {
	if u.EnvSourcePassword != "" && u.EnvSourceUsername != "" {
		logrus.Infof("Running docker login to source registry")
		err := registry.DockerLogin(
			u.EnvSourceRegistry, u.EnvSourceUsername, u.EnvSourcePassword)
		if err != nil {
			return fmt.Errorf("PrepareDockerLogin: failed to login to %s: %w",
				u.EnvSourceRegistry, err)
		}
	}

	if u.EnvDestPassword != "" && u.EnvDestUsername != "" {
		logrus.Infof("Running docker login to destination registry")
		err := registry.DockerLogin(
			u.EnvDestRegistry, u.EnvDestUsername, u.EnvDestPassword)
		if err != nil {
			return fmt.Errorf("PrepareDockerLogin: failed to login to %s: %w",
				u.EnvDestRegistry, err)
		}
	}

	return nil
}

func DockerLoginRegistry(reg string) error {
	logrus.Infof("Logging into %q", reg)
	username, passwd, err := registry.GetDockerPassword(reg)
	if err != nil {
		logrus.Warnf("Please input password of registry %q manually", reg)
		username, passwd, _ = u.ReadUsernamePasswd()
	}
	return registry.DockerLogin(reg, username, passwd)
}

func Worker(num int, failList string, cb func(m *mirror.Mirror)) (
	chan *mirror.Mirror, *sync.WaitGroup) {

	var writeFileMutex = new(sync.Mutex)
	var wg = new(sync.WaitGroup)
	worker := func(ch chan *mirror.Mirror) {
		defer wg.Done()
		for m := range ch {
			err := m.Start()
			if err != nil {
				logrus.WithField("M_ID", m.MID).Errorf("FAILED [%s:%s]",
					m.Source, m.Tag)
				logrus.WithField("M_ID", m.MID).Error(err)
				writeFileMutex.Lock()
				if err := u.AppendFileLine(failList, m.Line); err != nil {
					logrus.Error(err)
				}
				writeFileMutex.Unlock()
			}
			if cb != nil {
				cb(m)
			}
		}
	}
	ch := make(chan *mirror.Mirror)
	for i := 0; i < num; i++ {
		wg.Add(1)
		go worker(ch)
	}
	return ch, wg
}
