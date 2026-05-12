# Titel: Als Administrator möchte ich über CLI-Commands Benutzer in einer YAML-Datei pflegen können, wobei Passwörter als Hash sicher gespeichert werden

## Inhalt
Als Administrator benötige ich die Möglichkeit, Benutzerdaten über CLI-Commands in einer YAML-Datei zu verwalten. Dies umfasst das Anlegen, Bearbeiten und Löschen von Benutzern. Die Passwörter müssen dabei als kryptografische Hashes gespeichert werden, um die Sicherheit zu gewährleisten. Die YAML-Datei sollte die Benutzerdaten in einem strukturierten Format speichern, das einfach zu lesen und zu bearbeiten ist. Bei der Authentifizierung müssen die eingegebenen Passwörter mit den gespeicherten Hashes verglichen werden.

## Akzeptanzkriterien
- Der Administrator kann Benutzer über CLI-Command `user add` in einer YAML-Datei anlegen
- Der Administrator kann Benutzer über CLI-Command `user update` in einer YAML-Datei bearbeiten
- Der Administrator kann Benutzer über CLI-Command `user delete` in einer YAML-Datei löschen
- Passwörter werden als kryptografische Hashes (z.B. bcrypt) gespeichert
- Die YAML-Datei hat ein definiertes Format für Benutzerdaten
- Die Authentifizierung vergleicht eingegebene Passwörter mit den gespeicherten Hashes
- Die YAML-Datei kann unter einem konfigurierbaren Pfad gespeichert werden
- Bei der Erstellung eines neuen Benutzers wird das Passwort automatisch gehasht
- Die YAML-Datei ist lesbar und editierbar für Administratoren
- Die Implementierung verwendet eine YAML-Bibliothek für das Parsen und Schreiben
- Fehler bei der YAML-Verarbeitung werden angemessen behandelt
- Die YAML-Datei kann Benutzer mit Benutzername, Passwort-Hash und optional Metadaten speichern
- Die Implementierung stellt sicher, dass die YAML-Datei atomar geschrieben wird, um Korruption zu vermeiden
- Die CLI-Commands geben angemessene Erfolgs- oder Fehlermeldungen aus
- Die CLI-Commands unterstützen Hilfe-Optionen (--help)
