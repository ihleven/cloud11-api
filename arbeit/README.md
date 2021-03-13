# Arbeit

## TODO

arbeit.handler wird eigentlich nicht gebraucht und kann in arbeit integriert werden.
In Main würde dann 
	uc := arbeit.NewUsecase(repo)
durch
	uc := arbeit.NewHandler(repo)
ersetzt werden.

Rückgabe sollte für 

{
    "account": 1,
    "datum": "2018-08-07",
    "job": "IC",
    "kalendertag": {
        "date": "2018-08-07T00:00:00Z",
        "year": 2018,
        "month": 8,
        "day": 7,
        "jahrtag": 219,
        "kw_jahr": 2018,
        "kw_nr": 32,
        "kw_tag": 2
    },
    "arbeitstag": {
        "IC": {
            "status": "S",
            "kategorie": "H",
            "soll": 8,
            "beginn": "2018-08-07T07:06:00+02:00",
            "ende": "2018-08-07T16:57:00+02:00",
            "brutto": 9.85,
            "pausen": 0,
            "extra": 0,
            "netto": 9.85,
            "diff": -1.8499999999999996,
            "saldo": 0,
            "kommentar": ""
        }
    }
}