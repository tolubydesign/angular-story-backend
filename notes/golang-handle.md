# Notes on Handling, Setting Up, Modifying and Maintaining a GoLang Project

#### Installing GoLang

[Arch Linux Install: Go](https://wiki.archlinux.org/title/Go)

###### Installing Dependencies

List compiler: `go list -e -json -compiled .`

Initialize GoLang Project: `go mod init projectname`

Tidy up Mod file: `go mod tidy`

Install Package: `go install example.com/cmd@latest` OR `go install github.com/lib/pq@latest`


###### Change folder permissions

`sudo chmod 744 ./postgres_database` OR `sudo chmod -R 755 ./postgres_database`

#### Resources

[How to Use Go Modules](https://www.digitalocean.com/community/tutorials/how-to-use-go-modules)

[How to Use Go Modules](https://jogendra.dev/import-cycles-in-golang-and-how-to-deal-with-them)

[Modify File Permissions with chmod](https://www.linode.com/docs/guides/modify-file-permissions-with-chmod/)

