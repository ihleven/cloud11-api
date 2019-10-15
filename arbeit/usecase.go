package arbeit

// Facade ...
var Facade AFacade

// AFacade ...
type AFacade struct {
	repo Repository
}

// Arbeitstag ..
func (a AFacade) Arbeitstag(id int) (*Arbeitstag, error) {

	tag, err := a.repo.ReadArbeitstag(id)

	return tag, err
}
