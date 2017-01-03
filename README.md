# Schedule-Buddy
A site to simplify the scheduling process, specifically at Miami University.

# Run
Install Go

    brew install go

Add this to your `.bashrc`
```shell
export PATH=$PATH:/usr/local/opt/go/libexec/bin
export GOPATH=$HOME/golang
export GOROOT=/usr/local/opt/go/libexec
```
then from the root directory

    go get -u github.com/labstack/echo
    go run *.go
