package arbeit

// ArbeitstagStatus ...
type ArbeitstagStatus string

// ArbeitstagStatus ...
const (
	StatusFrei             ArbeitstagStatus = "-"
	StatusArbeitstag       ArbeitstagStatus = "A"
	StatusArbeitsvormittag ArbeitstagStatus = "V"
	StatusFeiertag         ArbeitstagStatus = "S"
	StatusBetriebsfrei     ArbeitstagStatus = "F"
)

// ArbeitstagKategorie ...
type ArbeitstagKategorie string

// Buero ...
const (
	Buero             ArbeitstagKategorie = "B"
	Homeoffice        ArbeitstagKategorie = "H"
	Krank             ArbeitstagKategorie = "K"
	Urlaub            ArbeitstagKategorie = "U"
	Sonderurlaub      ArbeitstagKategorie = "S"
	Freizeitausgleich ArbeitstagKategorie = "F"
)

// ZeitspanneKategorie ...
type ZeitspanneKategorie string

// ZeitspanneKategorie ...
const (
	StatusAZ              ZeitspanneKategorie = "A"
	StatusPause           ZeitspanneKategorie = "P"
	StatusExtra           ZeitspanneKategorie = "E"
	StatusWeg             ZeitspanneKategorie = "W"
	StatusRestpausenabzug ZeitspanneKategorie = "R"
)
