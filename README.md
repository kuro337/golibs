# Share Library among Apps

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
