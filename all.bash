cd cmd
go install -v ./... # build all Go+ tools
cd ..
goproot=`pwd`
if [ ! -e "$HOME/gop" ]; then
    if [ "$goproot" != "$HOME/gop" ]; then
        ln -s $goproot $HOME/gop
    fi
fi
gop install ./...   # build all Go+ tutorials
