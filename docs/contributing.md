## Build

Requirements

* Go 1.13+

Update modules

```shell
go get -u ./...
go mod tidy
```

# Snap

Requirements

* snapcraft 3.10+

Build:

```shell
snapcraft
```

Clean:

```shell
snapcraft clean
```

Install / Remove:

```shell
snap install manala_[version]_amd64.snap --dangerous
snap remove manala
```

## Documentation

```shell
docker run --rm -v $(pwd):/data -p 8000:8000 nicksantamaria/mkdocs serve -a 0.0.0.0:8000
```
