# Architekturhinweise

## generell
Es handelt sich um ein Projekt in der Sprache golang.

## Steuerung
Es handelt sich um ein Projekt, welches nur über eine Rest-API verfügen soll. 

## Berechtigung
Das Berechtigungssystem ist sehr wichtig. Jeder request muss über eine Session-Id autorisiert werden. Dabei gibt es für jeden endpunkt eine eigene HTTP_* Berechtigung.

## Storage
Die Speicherung findet über einzelene JSON dateien statt. Dabei gilt immer: Jede datei hat eine generische ID und liegt im Verzeichnis zu seiner entität.
 zum Beispiel: lists/alsdkjsfhakdsnjva.json

## Struktur
Es werden module zu sub-domains erstellt. Zu jedem Sub-Modul gelten die folgenden Regeln:
* Rest-Contoller werden in der Datei rest.go implementiert.
* Repositories werden in der dateu repository.go als Interface definiert. Die implementierung der Interfaces erfolt in der datei <type>_repository, wbei typ den Typ des Repos beschreibt. 
Also beispielsweise erfolgt eine In-memory Implementierung in memory_repository.go oder json files in json_repository.go.
* Komplextere Logik erfolt in der datei service.go. Dazu wird immer ein Strukt angelegt.




