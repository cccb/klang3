package main

import (
	"github.com/dhowden/tag"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Sample struct {
	Id int `json:"id"`

	Path  string `json:"path"`
	Title string `json:"title"`
	Group string `json:"group"`

	Metadata tag.Metadata `json:"-"`

	isPlaying bool
	cmd       *exec.Cmd
	err       error
}

func (self *Sample) Start(playCmd string) (chan bool, error) {
	// Make stopped callback channel
	done := make(chan bool, 1)

	log.Println("Playing file:", self.Path)

	// Start playback
	self.cmd = exec.Command(playCmd, self.Path)
	err := self.cmd.Start()
	if err != nil {
		done <- false
		return done, err
	}

	self.isPlaying = true

	go func() {
		// Wait until finished, reset state
		self.err = self.cmd.Wait()
		if self.err != nil {
			log.Println("Could not play sample:", self.err)
		}

		self.isPlaying = false

		done <- true
	}()

	return done, nil
}

func (self *Sample) Stop() error {
	if self.isPlaying == false {
		return fmt.Errorf("Sample not playing")
	}

	err := self.cmd.Process.Kill()
	if err != nil {
		log.Println("Could not kill player:", err)
		return err
	}

	return nil
}

func (self *Sample) _makeTitle() {
	if self.Metadata != nil && self.Metadata.Title() != "" && self.Metadata.Artist() != "" {
		self.Title = fmt.Sprintf("%v (%v)", self.Metadata.Title(), self.Metadata.Artist())
	} else if self.Metadata != nil && self.Metadata.Title() != "" {
		self.Title = self.Metadata.Title()
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

	for id, filename := range files {
		file, err := os.Open(filename)
		defer file.Close()

		if err != nil {
			log.Println("Could not open file:", err)
			continue
		}

		// Read metadata
		meta, _ := tag.ReadFrom(file)

		sample := &Sample{
			Id:       id,
			Path:     filename,
			Metadata: meta,
		}
		sample._makeTitle()
		sample._makeGroup(self.Path)

		self.samplesCache = append(self.samplesCache, sample)
	}

	return nil
}

func (self *Repository) _enumerateDuplicateTitles() {

}

func (self *Repository) AllSamples() []*Sample {
	return self.samplesCache
}

func (self *Repository) Samples(group string) []*Sample {
	samples := []*Sample{}

	for _, s := range self.samplesCache {
		if group == "" || group == "*" {
			samples = append(samples, s)
			continue
		}

		if s.Group == group {
			samples = append(samples, s)
		}
	}

	return samples
}

func (self *Repository) Groups() []string {
	groups := []string{}
	for _, s := range self.samplesCache {
		hasGroup := false
		for _, g := range groups {
			if s.Group == g {
				hasGroup = true
				break
			}
		}
		if !hasGroup {
			groups = append(groups, s.Group)
		}
	}

	return groups
}

func (self *Repository) GetSampleById(id int) *Sample {
	if id >= len(self.samplesCache) {
		return nil
	}
	return self.samplesCache[id]
}
