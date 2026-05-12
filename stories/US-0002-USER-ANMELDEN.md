# Titel: Als Nutzer möchte ich mich mit meinem Usernamen und Passwort anmelden können um eine gültige Session zu erzeugen

## Inhalt
Als Nutzer möchte ich die Möglichkeit haben, mich mit meinen Anmeldedaten (Benutzername und Passwort) zu authentifizieren, um eine gültige Session zu erhalten. Diese Session wird für alle weiteren API-Aufrufe benötigt, wie in der Architektur-Dokumentation festgelegt. Das System sollte die Anmeldedaten validieren und bei erfolgreicher Authentifizierung eine Session-ID zurückgeben, die in subsequenten Requests im X-Session-ID Header verwendet werden kann.

## Akzeptanzkriterien
- Der Nutzer kann sich mit Benutzername und Passwort anmelden
- Bei erfolgreicher Anmeldung wird eine gültige Session-ID zurückgegeben
- Die Session-ID muss mindestens 40 Zeichen lang sein und zufäälig erzeugt werden(gemäß Architektur-Dokumentation)
- Bei falschen Anmeldedaten wird ein entsprechender Fehler zurückgegeben
- Die Session-ID kann in subsequenten Requests im X-Session-ID Header verwendet werden
- Die Anmeldung erfolgt über einen REST-Endpunkt /login mit POST-Methode
- Erfolgreiche Anmeldung gibt HTTP-Status 200 zurück
- Fehlgeschlagene Anmeldung gibt HTTP-Status 401 zurück
- Die Session-ID muss den Anforderungen der Auth-Middleware entsprechen
