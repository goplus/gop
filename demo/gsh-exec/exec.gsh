gop "run", "./foo"
exec "gop run ./foo"
exec "FOO=100 gop run ./foo"
exec {"FOO": "101"}, "gop", "run", "./foo"
exec "gop", "run", "./foo"
exec "ls $HOME"
ls ${HOME}
