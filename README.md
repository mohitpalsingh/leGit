# leGit
a stripped down implementation of git written in go

# Description
a small Git implementation that's capable of initializing a repository, creating objects in /.git (blobs, trees) and creating commits

# Commands
1. init - to initialize a leGit repo
2. cat-file <blob_sha1> - to read a blob object from it's sha1 hash
3. hash-object -w <filename> - to create a blob object of the file
4. ls-tree <tree_sha1> - to read all contents of leGit tree object (type of object, name, sha1 hash)
5. ls-tree --name-only <tree_sha1> - to read all the files' names of a leGit tree object
6. write-tree - to make a tree object from the files that exists in current directory
7. commit-tree <tree_object_sha1> -p <parent_commit_sha1> -m <commit_message> - to write a commit based on the passed tree object with parent commit reference, author's creds and an optional commit message

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
