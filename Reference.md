# Share Library among Apps

- Pushing Properly

```bash
# Created macro - will run above commands
ghpushtags v0.2.3 "Websockets Upgrades , Profile Guided Optimization"

# From App
go get -u github.com/kuro337/golibs@v0.2.1

import (
	"github.com/kuro337/golibs/utils"
	"github.com/kuro337/golibs/profiling"
)

git add .
git commit -m "some changes"
git tag v1.0.4
git push
git push --tags



# View latest tag
git describe --tags --abbrev=0

#  Recursive Update in all SubDirs
go get -u ./...

# In app using
go clean -cache
go clean -modcache
go mod vendor
go mod tidy
go get -u


# Setting a specific Version
# Get the commit with the hash of the latest commit from Github directly
go get github.com/kuro337/golibs@1739bd0

# Then add the replace statement - so it resolves to that version
require (
	github.com/kuro337/golibs/profiling v0.0.0-20230829052146-b948630de748
	github.com/cockroachdb/pebble v0.0.0-20230826001808-0b401ee526b8
)

replace github.com/kuro337/golibs => github.com/kuro337/golibs v1.0.8

go get -u


```

- Create folder golibs - this is the Repo

- golibs hhas samplefunc folder - samplefunc has

```bash
go.mod

module github.com/kuro337/golibs/samplefunc

go 1.21.0


hello.go
package samplefunc

import "fmt"

func Hello() {
	fmt.Println("Hello from samplefunc!")
}


Then the repo on github is
kuro337/golibs - samplefunc package , etc.


```

To use in different apps -

```go
// randomapp.go

import (
	"proper/logging"
	"proper/webserver"

	"github.com/kuro337/golibs/samplefunc"
)

samplefunc.Hello()

// go mod tidy - to get lib

// If we update a lib

//   go get -u github.com/kuro337/golibs/httpinterface@latest
//   go get -u github.com/kuro337/golibs/samplefunc@latest
```

- Implemening Multiple interfaces using Generics to enforce a Contract

```go

package main

import "fmt"

type R interface {
	Read() bool
}
type W interface {
	Write() bool
}
type RW[T any] interface {
	R
	W
}

type File struct{}

func (f File) Read() bool {
	return true
}

func (f File) Write() bool {
	return true
}

func main() {
	var intSat RW[File] = File{}
	intSat.Read()

	ints := RW[File](File{})
	ints.Read()
}


```

- Git CLI and GH Cli

```bash
cd webinterface

gh auth login

gh repo create go-libs --public

git remote add origin git@github.com:kuro337/go-libs.git

git push origin main


# Creating alias file for zsh
touch ~/alias-config.sh

ghcreate go-libs
l
gs

# Add to ~/.zshrc
source ~/alias-config.sh

# Reload
source ~/.zshrc  # For Zsh

# Usage
ghcreate go-libs # Creates repo

# Deleting Repo
gh auth refresh -h github.com -s delete_repo
gh repo delete testcli
```
