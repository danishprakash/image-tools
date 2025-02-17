package mirror

import (
	"fmt"

	u "cnrancher.io/image-tools/utils"
	"github.com/sirupsen/logrus"
)

func (m *Mirror) StartLoad() error {
	if m == nil {
		return fmt.Errorf("StartLoad: %w", u.ErrNilPointer)
	}
	if m.Directory == "" {
		return fmt.Errorf("StartLoad: directory is empty string")
	}
	if m.Mode != MODE_LOAD {
		return fmt.Errorf("StartLoad: mirrorer is not in LOAD mode")
	}

	logrus.WithField("M_ID", m.MID).
		Infof("DEST: [%v] TAG: [%v]", m.Destination, m.Tag)

	for _, img := range m.images {
		img.MID = m.MID
		if err := img.Load(); err != nil {
			return fmt.Errorf("StartLoad: %w", err)
		}
	}

	if !m.compareSourceDestManifest() {
		logrus.WithField("M_ID", m.MID).
			Info("Creating dest manifest list...")
		// Create a new dest manifest list
		if err := m.updateDestManifest(); err != nil {
			return fmt.Errorf("Mirror: %w", err)
		}
	} else {
		logrus.WithField("M_ID", m.MID).
			Info("Dest manifest list already exists, no need to recreate")
	}

	// logrus.WithField("M_ID", m.MID).
	// 	Info("Creating dest manifest list...")
	// if err := m.updateDestManifest(); err != nil {
	// 	return fmt.Errorf("StartLoad: %w", err)
	// }

	logrus.WithField("M_ID", m.MID).
		Infof("Loaded \"%s:%s\"", m.Destination, m.Tag)

	return nil
}

func (m *Mirror) Loaded() int {
	var num int = 0
	for _, img := range m.images {
		if img.Loaded {
			num++
		}
	}
	return num
}
