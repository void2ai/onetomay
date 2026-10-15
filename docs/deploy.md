# Deploy onetomany System-Wide

## Install

```bash
# Build
go build -o onetomany .

# Copy to system PATH (requires sudo)
sudo cp ./onetomany /usr/local/bin/
```

## Verify

```bash
which onetomany
onetomany --help
```

## Update

```bash
go build -o onetomany .
sudo cp ./onetomany /usr/local/bin/
```

## Uninstall

```bash
sudo rm /usr/local/bin/onetomany
```
