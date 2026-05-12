---
name: refactor
description: Dieser Skill führt Code-Refactorings durch und dokumentiert die Änderungen in der Architektur-Dokumentation.
---


# Refactoring Skill

## Zweck
Dieser Skill führt Code-Refactorings durch und dokumentiert die Änderungen in der Architektur-Dokumentation, um zukünftige Fehler zu vermeiden.

## Arbeitsweise

### 1. Refactoring-Anweisungen entgegennehmen
- Analysiere die gegebenen Refactoring-Anforderungen
- Identifiziere betroffene Codebereiche
- Prüfe bestehende Tests

### 2. Refactoring durchführen
- Führe schrittweise Änderungen durch
- Stelle sicher, dass alle Tests weiterhin bestehen
- Halte die Architekturprinzipien ein

### 3. Dokumentation aktualisieren
- Aktualisiere `docs/ARCHITECTURE.md` mit den Änderungen
- Dokumentiere Best Practices
- Füge Beispiele hinzu
- Ergänze Anti-Patterns, die vermieden werden sollen

## Refactoring-Prozess

### Schritt 1: Analyse
```
1. Lese die Refactoring-Anforderungen
2. Identifiziere betroffene Dateien und Module
3. Prüfe bestehende Tests
4. Analysiere Abhängigkeiten
```

### Schritt 2: Planung
```
1. Erstelle eine Checkliste der Änderungen
2. Identifiziere Breaking Changes
3. Plane Test-Anpassungen
4. Dokumentiere erwartete Verbesserungen
```

### Schritt 3: Umsetzung
```
1. Führe Änderungen schrittweise durch
2. Teste nach jedem Schritt
3. Behebe aufgetretene Fehler
4. Optimiere Performance
```

### Schritt 4: Dokumentation
```
1. Aktualisiere Architektur-Dokumentation
2. Dokumentiere neue Patterns
3. Füge Code-Beispiele hinzu
4. Ergänze Warnungen vor Anti-Patterns
```

## Dokumentationsstruktur

### Für jedes Refactoring:
1. **Problem**: Was war das ursprüngliche Problem?
2. **Lösung**: Wie wurde es gelöst?
3. **Best Practices**: Was sollte in Zukunft beachtet werden?
4. **Code-Beispiele**: Vorher/Nachher-Vergleiche
5. **Tests**: Welche Tests wurden angepasst?

## Beispiel-Refactoring

### Problem:
```
"Die Authentifizierungs-Logik ist über mehrere Dateien verteilt und schwer wartbar"
```

### Lösung:
```
1. Konsolidiere Auth-Logik in middleware/auth.go
2. Erstelle klare Schnittstellen
3. Trenne Session-Verwaltung von Request-Handling
```

### Dokumentation:
```markdown
## Authentifizierung (refaktorisiert)

### Best Practices:
- Auth-Logik sollte in einer zentralen Middleware liegen
- Session-Verwaltung sollte interface-basiert sein
- Fehlerbehandlung sollte konsistent sein (401 für Auth-Fehler)

### Anti-Patterns:
- ❌ Auth-Logik in Handlern verteilen
- ❌ Session-IDs manuell generieren
- ❌ Fehlercodes inkonsistent verwenden

### Beispiel:
```go
// Gut: Zentrale Middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !validateSession(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```
```

## Ausführung

Um ein Refactoring durchzuführen:
1. Gib klare Refactoring-Anweisungen
2. Spezifiziere betroffene Module
3. Definiere erwartete Ergebnisse
4. Ich führe das Refactoring durch und dokumentiere es

## Qualitätskriterien

- Alle Tests müssen bestehen
- Code muss lesbarer werden
- Performance darf nicht leiden
- Architekturprinzipien müssen eingehalten werden
- Dokumentation muss aktuell sein
