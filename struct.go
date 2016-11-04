package dbmdl

// Interface for all dbmodels
type ifc interface {
	Load() error
	Save() error
}

// Struct is a struct that can be inherited for use with dbmdl
type Struct struct{}

// Load will populate the struct from the database
func (s *Struct) Load() error {

	return nil
}

// Save will save the struct from the database
func (s *Struct) Save() error {

	return nil
}
