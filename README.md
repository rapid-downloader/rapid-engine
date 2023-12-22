# Rapid Download Manager
Rapid is a download manager that capable to download in chunks

## Run the server
Install the dependency
```bash
go mod tidy
```

And then run the server
```bash
go run .
```

## Run the client
There is 2 available clients, CLI and GUI. 

### CLI
Build the CLI
```bash
cd cli
go build -o build .
```

Example usage
```bash
./build/cli download https://link.testfile.org/PDF50MB
```

### GUI
The GUI client developed with Wails. Currently stil in WIP. To open it, use the following command
```bash
cd gui
wails dev
```

