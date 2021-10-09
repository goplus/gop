cd cmd
go install -v ./...
cd ..
goproot=`pwd`
if [ ! -e "$HOME/gop" ]; then
    if [ "$goproot" != "$HOME/gop" ]; then
        ln -s $goproot $HOME/gop
    fi
fi
