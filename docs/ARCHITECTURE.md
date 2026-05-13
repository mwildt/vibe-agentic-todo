# Architekturhinweise

## generell
Es handelt sich um ein Projekt in der Sprache golang.

## Steuerung
Es handelt sich um ein Projekt, welches nur über eine Rest-API verfügen soll. 

## Berechtigung ⚠️ WICHTIG
Jeder API-Request muss über eine Session-ID autorisiert werden. Die Autorisierung muss folgende Anforderungen erfüllen:

1. **Session-ID Pflicht**: Jeder Request muss eine gültige Session-ID im Header enthalten
2. **Endpunkt-spezifische Berechtigungen**: Jeder Endpunkt erfordert eine eigene HTTP_* Berechtigung
3. **Middleware-basiert**: Die Autorisierung muss als Middleware implementiert werden, die vor den eigentlichen Handlern ausgeführt wird
4. **Fehlerbehandlung**: Nicht autorisierte Requests müssen mit 401 Unauthorized antworten
5. **Session-Speicherung**: Sessions müssen in einem Session-Store gespeichert und validiert werden

### Implementierungsvorgaben:
- Middleware-Datei: `middleware/auth.go`
- Session-Verwaltung: `auth/session.go`
- Session-Store Interface: `auth.SessionStore`
- In-Memory Implementierung: `auth.InMemorySessionStore`
- Berechtigungsprüfung muss vor jedem Handler aufgerufen werden
- Session-Store muss vor der Handler-Registrierung konfiguriert werden

### Best Practices:
- ✅ Session-Store als Dependency injizieren
- ✅ Thread-sichere Implementierung verwenden (Mutex)
- ✅ Session-Timeout implementieren (z.B. 24 Stunden)
- ✅ Session-IDs mit UUID oder ähnlichem generieren
- ✅ Session-Validierung gegen den Store durchführen

### Anti-Patterns:
- ❌ Sessions manuell verwalten
- ❌ Session-IDs ohne Store validieren
- ❌ Session-Store global ohne Dependency Injection verwenden
- ❌ Thread-Sicherheit ignorieren

### Beispiel:
```go
// Richtige Implementierung
sessionStore := auth.NewInMemorySessionStore()
auth.RegisterHandlers(sessionStore)
middleware.SetSessionStore(sessionStore)

// Falsche Implementierung (Anti-Pattern)
// auth.RegisterHandlers() // Ohne Session-Store
```

## Test Implementation Guidelines

### JSON Handling in Tests

✅ **DO:**
- Use `json.NewDecoder()` to parse responses into structs
- Use `json.Marshal()` to create request bodies
- Define proper Go structs with JSON tags for request/response formats
- Compare struct fields directly (e.g., `if response.Text != expectedText`)

❌ **DON'T:**
- Use `strings.Contains()` or other string operations to check JSON content
- Manually parse JSON strings with `strings.Index()` or regex
- Compare raw JSON strings

### Example of Correct Pattern:

```go
// Define response struct
type NoteResponse struct {
    ID   string `json:"id"`
    Text string `json:"text"`
}

// Parse response properly
var note NoteResponse
if err := json.NewDecoder(response.Body).Decode(&note); err != nil {
    t.Fatalf("Failed to parse response: %v", err)
}

// Compare fields directly
if note.Text != expectedText {
    t.Errorf("Expected text %q, got %q", expectedText, note.Text)
}
```

### Legacy Code Note:
Some older tests may use string operations for convenience. These should be migrated to proper JSON parsing when modified.

## Domain Layer
Die Domain-Schicht enthält die Business-Objekte und ihre Kernlogik. Jedes Domain-Objekt sollte seine eigene Business-Logik kapseln.

### User Domain Objekt
Das User-Objekt in `auth/user/domain.go` enthält:
- Passwort-Hashing (`HashPassword()`)
- Passwort-Verifikation (`VerifyPassword()`)
- Passwort-Änderung (`SetPassword()`)

### Best Practices für Domain-Objekte:
- ✅ Domain-Objekte sollten ihre eigene Validierungslogik enthalten
- ✅ Business-Regeln sollten in Domain-Methoden gekapselt sein
- ✅ Domain-Objekte sollten unabhänig von Storage- oder UI-Logik sein
- ✅ Nutze Konstruktoren für komplexe Objekt-Erstellung

### Anti-Patterns für Domain-Objekte:
- ❌ Domain-Logik in Repositories oder Services auslagern
- ❌ Domain-Objekte mit Storage-Abhängigkeiten
- ❌ Business-Regeln in UI- oder CLI-Schicht implementieren

## Storage
Die Speicherung findet über einzelne JSON-Dateien statt. Dabei gilt immer: Jede Datei hat eine generische ID und liegt im Verzeichnis zu seiner Entität.
Beispiel: lists/alsdkjsfhakdsnjva.json

## Struktur
Es werden Module zu Sub-Domains erstellt. Zu jedem Sub-Modul gelten die folgenden Regeln:
* Rest-Controller werden in der Datei rest.go implementiert.
* Domain-Objekte mit Business-Logik werden in domain.go definiert.
* Repositories werden in der Datei repository.go als Interface definiert. Die Implementierung der Interfaces erfolgt in der Datei <type>_repository, wobei typ den Typ des Repos beschreibt.
* Komplexere Logik erfolgt in der Datei service.go. Dazu wird immer ein Strukt angelegt.
* CLI-Commands werden im Verzeichnis cmd/cli implementiert und folgen dem Cobra-Framework-Pattern.

## CLI-Struktur
Administrative Funktionen werden als CLI-Commands implementiert:
* Command-Dateien folgen dem Muster: cmd/cli/cmd/<command>.go
* Jeder Command hat eine klare Hilfe-Dokumentation
* Commands geben strukturierte Ausgaben (JSON oder Text) zurück
* Fehler werden mit angemessenen Exit-Codes und Fehlermeldungen behandelt

### Best Practices für CLI-Commands:
- ✅ Verwende Cobra oder ein ähnliches Framework für CLI-Struktur
- ✅ Implementiere Hilfe-Flags (--help)
- ✅ Gib klare Erfolgs- und Fehlermeldungen aus
- ✅ Verwende angemessene Exit-Codes (0 für Erfolg, !=0 für Fehler)
- ✅ Unterstütze sowohl JSON- als auch Text-Ausgabe
- ✅ Nutze Domain-Objekte für Business-Logik (z.B. User.HashPassword())

### Anti-Patterns für CLI:
- ❌ Direkte Business-Logik in CLI-Dateien
- ❌ Keine Fehlerbehandlung in Commands
- ❌ Unklare oder fehlende Dokumentation
- ❌ Inkonsistente Ausgabeformate
- ❌ Manuelle Implementierung von Domain-Logik (z.B. Passwort-Hashing in CLI)

## Checkliste für neue Endpunkte
1. [ ] Autorisierungs-Middleware implementieren/einbinden
2. [ ] Session-ID Validierung implementieren
3. [ ] Endpunkt-spezifische Berechtigung prüfen
4. [ ] Fehlerfälle (401, 403) korrekt behandeln




