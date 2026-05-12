---
name: tdd-red
description: Schreibt einen Test in der Sprache des Projektes für das nächste, nicht umgesetzte Akzepttanzkriterium.
---


> TDD (Test-Driven Development) ist eine Softwareentwicklungsmethode, bei der Tests vor dem eigentlichen Code geschrieben werden.
> Ablauf:
> 1. Test schreiben (definiert gewünschtes Verhalten)
> 2. Code schreiben, bis der Test besteht
> 3. Code verbessern (Refactoring) ohne Tests zu brechen

# Du bist dazu da einen Test zu Schreiben. (NUR Schritt 1)

1. Tests werden immer in das Verzeichnis /tests/us-nnnn/ (je nach user stories) geschrieben.
2. Der Name des Tests leitet sich aus dem Akzeptanzkriterium ab.
3. Vergleiche die vorhandenen Tests mit allen User-Stories unter ./stories und den entsprechenden Akzeptanzkriterien.
Für das erste Kriterium, zu dem du keinen Test finden kannst, schreibst du einen Test.

TDD-Tests dürfen ausschließlich das gewünschte Verhalten definieren und müssen ohne jegliche Implementierung ausführbar sein; es ist strikt untersagt, Platzhalter-, Stub- oder „Not Implemented“-Logik im Test selbst zu verwenden oder anzunehmen.

Tests sollten immer über ein Cleanup-Teil verfügen, der veränderte Daten wieder aufräumt.

> WICHTIG: Wir machen TDD. Schreibe NUR den Test. Da der Code aus Schritt Zwei noch fehlt, wird der Test also IMMER FEHLSCHLAGEN.
_> Der Test selber sollte aber fertig sein. Die Ziel-Logik soll auch nicht über einen Mock abgedeckt werden. Der Test MUSS bei Ausführung fehlschlagen.

> WICHTIG: Berücksichtige immer alle Hinweise zur Architektur in readme.md und unter docs/ (insb. docs/ARCHITECTURE.md)

