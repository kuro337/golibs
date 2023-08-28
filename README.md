# Share Library among Apps

- Pushing Properly

```bash
git commit -a - m "say_things - some changes"
git tag v1.0.1
git push
git push --tags

# In the project-
# In github releases should have the tags on Right Side
github.com/Chinmay337/golibs/httpinterface v1.0.1

# Now delete go.mod and go.sum and run go mod tidy for latest deps





```

- Create folder golibs - this is the Repo

- golibs hhas samplefunc folder - samplefunc has

```bash
go.mod

module github.com/Chinmay337/golibs/samplefunc

go 1.21.0


hello.go
package samplefunc

import "fmt"

func Hello() {
	fmt.Println("Hello from samplefunc!")
}


Then the repo on github is
Chinmay337/golibs - samplefunc package , etc.


```

To use in different apps -

```go
// randomapp.go

import (
	"proper/logging"
	"proper/webserver"

	"github.com/Chinmay337/golibs/samplefunc"
)

samplefunc.Hello()

// go mod tidy - to get lib

// If we update a lib

//   go get -u github.com/Chinmay337/golibs/httpinterface@latest
//   go get -u github.com/Chinmay337/golibs/samplefunc@latest
```

- Implemening Multiple interfaces

```go

type Writable interface {
	Write() string
}

type Readable interface {
	Read() string
}

type ReadWrite interface {
	Writeable
	Readable
}

type SatisfyIface struct {
	content string
}

func (s *SatisfyIface) Read() string {
	return s.content
}

func (s *SatisfyIface) Write() string {
	return s.content
}

var satisfiesIface ReadWrite = &SatisfyIface{content:"abcdef"}




```

- Git CLI and GH Cli

```bash
cd webinterface

gh auth login

gh repo create go-libs --public

git remote add origin git@github.com:Chinmay337/go-libs.git

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
