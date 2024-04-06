# leGit
a stripped down implementation of git written in go

# Description
a small Git implementation that's capable of initializing a repository, creating objects in /.git (blobs, trees) and creating commits

# Testing locally

The `le_git.sh` script is expected to operate on the `.git` folder inside the
current working directory. If you're running this inside the root of this
repository, you might end up accidentally damaging your repository's `.git`
folder.


I suggest executing `le_git.sh` in a different folder when testing locally.
For example:

```sh
mkdir -p /tmp/testing && cd /tmp/testing
/path/to/your/repo/le_git.sh init
```

```sh
alias leGit=/path/to/your/repo/le_git.sh

mkdir -p /tmp/testing && cd /tmp/testing
leGit init
```
