package gravity

import (
	"encoding/gob"
	"os"
	"strings"
)

const filename = "gravity.gob"

type StorageService struct {
	g *Gravity
}

func newStorageService(g *Gravity) *StorageService {
	return &StorageService{
		g: g,
	}
}

func (s *StorageService) readStorage(fname string) (c *Credentials, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return
	}
	defer file.Close()

	c = &Credentials{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return nil, err
	}

	return
}

func (s *StorageService) writeStorage(fname string) (err error) {
	file, err := os.Create(fname)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(s.g.State.cred)
	if err != nil {
		return
	}

	return
}

// load() transports local storage data into Gravity.State.
// If Gravity.State doesn't match the local straage, returns an error.
func (s *StorageService) load() (err error) {
	c, err := s.readStorage(filename)
	if err != nil {
		return
	}

	if !(c.identifier == s.g.State.cred.identifier && c.pwd == s.g.State.cred.pwd) {
		return ErrStorageDoesNotMatch
	}

	s.g.State.cred = c

	return
}

// save() exports current Gravity.State as local storage data
func (s *StorageService) save() error {
	// Check idtype just in case.
	idtype := getIDType(s.g.State.cred.identifier)
	if idtype == -1 {
		return ErrInvalidIdentifier
	}

	return s.writeStorage(filename)
}

func (s *StorageService) createOneAndSave() error {
	gaid, _ := generateUUID()
	uuid, _ := generateUUID()

	s.g.State.cred.gaid = gaid
	s.g.State.cred.uuid = strings.ToUpper(uuid)

	return s.save()
}
