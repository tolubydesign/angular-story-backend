# Information Related to the Maintenance of the Project


Shows the mounted filesystems and their usage statistics

`sudo df -h`

Increase tmp/ size

[Temporarily increase size of tmp folder on Arch linux](https://gist.github.com/ertseyhan/618ab7998bdb66fd6c58)

```bash
#!/bin/bash

sudo mount -o remount,size=10G,noatime /tmp
echo "Done. Please use 'df -h' to make sure folder size is increased."
```

Clear tmp/ folder

Important. Don't use sudo. There are some files in this folder that are being used by background operations.
```bash
$ cd
$ rm -rf ~/tmp/*
```

Find and Remove files recently modified

```bash
$ find -atime 5 -mtime 5 | xargs rm
$ find {.} -type d -name 'go-*' -delete
$ find {.} -type f -name '*go-*' -delete

```

