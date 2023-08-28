# Share Library among Apps

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
