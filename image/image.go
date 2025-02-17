package image

import (
	"fmt"
	"strings"
)

type Image struct {
	Source      string
	Destination string
	Tag         string
	Arch        string
	Variant     string
	OS          string
	OsVersion   string

	// Digest is the source image manifest sha256sum
	Digest string

	// Directory is the absolute path to save the image
	Directory string

	// SavedFolder is the folder name of the saved image
	SavedFolder string

	Copied bool
	Saved  bool
	Loaded bool

	SourceMediaType string

	// IID is the ID of the Image
	IID int
	// MID is the ID of the Mirror
	MID int
}

type ImageOptions struct {
	Source      string
	Destination string
	Tag         string
	Arch        string
	Variant     string
	OS          string
	OsVersion   string
	Digest      string
	Directory   string
	SavedFolder string

	SourceMediaType string

	MID int
}

func NewImage(opts *ImageOptions) *Image {
	return &Image{
		Source:          opts.Source,
		Destination:     opts.Destination,
		Tag:             opts.Tag,
		Arch:            opts.Arch,
		Variant:         opts.Variant,
		OS:              opts.OS,
		OsVersion:       opts.OsVersion,
		Digest:          opts.Digest,
		Directory:       opts.Directory,
		SavedFolder:     opts.SavedFolder,
		SourceMediaType: opts.SourceMediaType,
		MID:             opts.MID,
	}
}

// CopiedTag gets the tag of the copied image,
// the format should be: ${VERSION}-${ARCH}${VARIANT}${EXTRA}
//
// If the OS is not linux, such as windows, darwin, etc
// the format should be: ${VERSION}-${OS}-${ARCH}${VARIANT}
//
// If the extra is not empty, the tag format should be:
// ${VERSION}-${OS}-${ARCH}${VARIANT}-${EXTRA}
func CopiedTag(tag, OS, arch, variant string, extra ...string) string {
	var (
		prefix string // ${VERSION}-${OS} or // ${VERSION} only if linux
		suffix string // ${ARCH}${VARIANT} (variant can be empty)
	)
	switch OS {
	case "":
		prefix = tag
	case "linux":
		prefix = tag
	default:
		prefix = fmt.Sprintf("%s-%s", tag, OS)
	}

	switch arch {
	case "amd64":
		suffix = arch
	case "arm64":
		// there is only one variant of arm64 is v8, so discard it
		suffix = arch
	case "arm":
		// arm has variant v5, v7, etc...
		suffix = fmt.Sprintf("%s%s", arch, variant)
	default:
		// other arch: s390x, ppc64...
		suffix = fmt.Sprintf("%s%s", arch, variant)
	}
	if len(extra) > 0 {
		suffix = fmt.Sprintf("%s-%s", suffix, strings.Join(extra, "-"))
	}

	return fmt.Sprintf("%s-%s", prefix, suffix)
}
