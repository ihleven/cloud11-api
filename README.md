# cloud11-api




# package httperror

erfuellt error interface und ermoeglicht in der applikation zeitnah fehelrcodes zu definieren
erkennt diverse standard fehlertypen
httperror.GetStatus(err) liefert http code und detailiert message entweder direkt aus einem httperror oder aus normaelem fehler (status 500)


The action takes HTTP requests (URLs and their methods) and uses that input to interact with the domain, after which it passes the domain's output to one and only one responder.
The domain can modify state, interacting with storage and/or manipulating data as needed. It contains the business logic.
The responder builds the entire HTTP response from the domain's output which is given to it by the action.


Interfaces:
 * Responder -> Respond(http.Responsewriter, http.Reuquest, interface{}) kann anhand des requests unterschiede machen
 * Renderer -> Render(http.Responsewriter, interface{}) eryeugt ausgabe auf w fuer interface{}

 * noch zu definierendes interface fuer domain logic

Actioneer ist struct, das http.HAndler erfuellt 

Drive ist interface, Hidrive und FS sind implemtierungen davon 


DriveHandler und DriveServeHAndler sind structs mit Drive und REsponder als komponenten
{
    drive Drive
    responder Responder
}