package main

import (
	"github.com/dhowden/tag"

	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Sample struct {
	Path  string
	Title string
	Group string

	Metadata tag.Metadata
}

func (self *Sample) _makeTitle() {
	if self.Metadata != nil && self.Metadata.Title() != "" {
		self.Title = fmt.Sprintf("%v (%v)", self.Metadata.Title(), self.Metadata.Artist())
	} else {
		self.Title = _stripSuffix(path.Base(self.Path))
	}
}

func (self *Sample) _makeGroup(basePath string) {
	rel := strings.Replace(self.Path, basePath, "", 1)
	t := strings.Split(rel, "/")
	if len(t) > 1 {
		self.Group = t[1]
	}
}

type Repository struct {
	Path string

	samplesCache []*Sample
}

func NewRepository(path string) *Repository {
	repo := &Repository{
		Path: path,
	}

	return repo
}

func _stripSuffix(path string) string {
	t := strings.Split(path, ".")
	return strings.Join(t[:len(t)-1], ".")
}

func _hasAudioSuffix(path string) bool {
	suffices := []string{
		".wav", ".mp3", ".ogg", ".flac",
	}
	for _, suffix := range suffices {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return false
}

func (self *Repository) listFiles() ([]string, error) {
	log.Println("Searching files in:", self.Path)

	files := []string{}

	err := filepath.Walk(self.Path, func(path string, f os.FileInfo, err error) error {
		// Filter at by suffix
		if _hasAudioSuffix(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return files, nil
}

func (self *Repository) Update() error {
	files, err := self.listFiles()
	if err != nil {
		return err
	}

	self.samplesCache = []*Sample{}

	for _, filename := range files {
		file, err := os.Open(filename)
		defer file.Close()

		if err != nil {
			log.Println("Could not open file:", err)
			continue
		}

		// Read metadata
		meta, _ := tag.ReadFrom(file)

		sample := &Sample{
			Path:     filename,
			Metadata: meta,
		}
		sample._makeTitle()
		sample._makeGroup(self.Path)

		self.samplesCache = append(self.samplesCache, sample)
	}

	// Sort sample cache by group and title

	return nil
}

func (self *Repository) Samples() ([]*Sample, error) {
	err := self.Update()
	if err != nil {
		return []*Sample{}, err
	}

	return self.samplesCache, nil
}
