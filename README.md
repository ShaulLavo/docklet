# Ducklet

A lightweight Docker implementation written in Go.

## Getting Started

First, clone the repository:

```sh
git clone https://github.com/ShaulLavo/docklet.git
```

You will get:

```
docklet/
├─ .gitignore
├─ .vscode/
│  └─ settings.json
├─ build.sh
├─ cli.go
├─ go.mod
├─ go.sum
└─ main.go
```

Make the build script executable and run it:

```sh
chmod +x build.sh
./build.sh
```

## Usage Example

```sh
./docklet build --tag my_image_name --path DockletFile
./docklet run /bin/bash
```
