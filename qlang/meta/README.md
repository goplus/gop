qlang package meta
==================

### Usage

```
qlang.Import("", meta.Exports)
```

### Function
		
* pkgs `func() []string`	

```		
>>>pkgs()
[bufio bytes md5 io ioutil hex json errors math os path http reflect runtime strconv strings sync]
```

* dir `func(interface {}) []string`

```
>>>dir(md5)
[_name new sum sumstr hash BlockSize Size]
```

* doc `func(interface {}) string`

```
>>>doc(md5)
package crypto/md5
hash	func(interface {}, ...interface {}) string 
BlockSize	int 
Size	int 
_name	string 
new	func() hash.Hash 
sum	func([]uint8) [16]uint8 
sumstr	func([]uint8) string 
```
