set -e

cd cmd
go install -v ./... # build all Go+ tools
cd ..
goproot=`pwd`
if [ ! -e "$HOME/gop" ]; then
    if [ "$goproot" != "$HOME/gop" ]; then
        ln -s $goproot $HOME/gop
    fi
fi

if  ! command -v gop &> /dev/null ; then
    echo "gop has been installed. To use it conveniently, you need put $GOPATH/bin or $HOME/go/bin directory into your PATH environment variable, see \`go help install\`"
else
    gop install ./...   # build all Go+ tutorials
fi
