# Startercode für VSS Blatt 1

Erzeugen Sie als erstes einen `develop`-Branch und arbeiten Sie in diesem bzw.
branchen Sie bei Bedarf von diesem weiter. Die Abgabe erfolgt dann als Merge-Request
vom `develop`- zum `main`-Branch.

## Passen Sie den Code als erstes an Ihre Umgebung an!

Die Änderungen muss nur einer aus dem Team committen und pushen.

1.  Löschen Sie die Dateien `go.mod` und `go.sum` und erzeugen Sie diese mit

    ```
    go mod init <path>
    ```

    neu, wobei Sie `<path>` durch Ihr Repository ersetzen müssen.

2.  Passen Sie die Importe für `customer` und `messages` in den Go-Dateien an
    den neuen Modulpfad an.

3.  Probieren Sie aus, ob der API-Test für den Customer-Service läuft:

    ```
    go run apitest/customer/main.go
    ```

    Wenn das Programm ohne Fehlermeldung durchläuft, passt alles. Als Ausgabe sollten
    Sie etwas in folgender Form sehen:

    ```
    ... [ANSWER] got customer="id:1  name:\"Max Mustermann\""
    ```

4.  Probieren Sie aus, ob der Customer-Service startet:

    ```
    go run customer/service/main.go
    ```

    Sie sollten dann sowas in der Art sehen:

    ```
    ... Starting remote with address address="127.0.0.1:9010"
    ... Started Activator
    ... Started EndpointManager
    ... Starting Proto.Actor server address="127.0.0.1:9010"
    ```

5.  Probieren Sie aus, ob der Client läuft. Dazu muss der Customer-Service
    bereits laufen.

    ```
    go run client/main.go
    ```

## Committen und Pushen

Wenn alles passt, dann committen und pushen Sie jetzt. In Ihrem GitLab-Projekt
wird Ihr Projekt jetzt von [Golangci-lint](https://github.com/golangci/golangci-lint)
überprüft, die Tests werden ausgeführt.

## Ausführen der Services

Öffne ein neues Terminal-Fenster und navigiere in das Verzeichnis des jeweiligen Service.
Anschließend kannst du den Service mit `go run service/main.go` starten.

Dies führt man zunächst für den Book-Service und anschließend für den Customer-Service aus.
Erst danach für den Library-Service.

Als letzes startet man den miniClient mit `go run client/miniClient.go`.

Das Buch wird erfolgreich ausgeliehen, wenn die Ausgabe des miniClients wie folgt aussieht:
`Got book:  Worm`
