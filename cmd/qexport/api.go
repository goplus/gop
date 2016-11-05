package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func apipath(base string) string {
	return filepath.Join(os.Getenv("GOROOT"), "api", base)
}

//pkg syscall (windows-386), const CERT_E_CN_NO_MATCH = 2148204815
var sym = regexp.MustCompile(`^pkg (\S+)\s?(.*)?, (?:(var|func|type|const)) ([A-Z]\w*)`)
var num = regexp.MustCompile(`^\-?[0-9]+$`)

type KeyType int

const (
	Normal      KeyType = 1
	ConstInt64  KeyType = 2
	ConstUnit64 KeyType = 3
)

type GoApi struct {
	Keys map[string]KeyType
	Ver  string
}

func LoadApi(ver string) (*GoApi, error) {
	f, err := os.Open(apipath(ver + ".txt"))
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)
	keys := make(map[string]KeyType)
	for sc.Scan() {
		l := sc.Text()
		has := func(v string) bool { return strings.Contains(l, v) }
		if has("interface, ") || has(", method (") {
			continue
		}
		if m := sym.FindStringSubmatch(l); m != nil {
			// 1 pkgname
			// 2 os-arch-cgo
			// 3 var|func|type|const
			// 4 name
			key := m[1] + "." + m[4]
			if _, ok := keys[key]; ok {
				continue
			}
			keys[key] = Normal
			if m[3] == "const" {
				if pos := strings.LastIndex(l, "="); pos != -1 {
					value := strings.TrimSpace(l[pos+1:])
					if num.MatchString(value) {
						_, err := strconv.ParseInt(l[pos+2:], 10, 32)
						if err != nil {
							if value[0] == '-' {
								keys[key] = ConstInt64
							} else {
								keys[key] = ConstUnit64
							}
						}
					}
				}
			}
		}
	}
	return &GoApi{Ver: ver, Keys: keys}, nil
}

type ApiCheck struct {
	Base map[string]KeyType
	Apis []*GoApi
}

func NewApiCheck() *ApiCheck {
	ac := &ApiCheck{}
	ac.Base = make(map[string]KeyType)
	return ac
}

func (ac *ApiCheck) LoadBase(vers ...string) error {
	for _, ver := range vers {
		api, err := LoadApi(ver)
		if err != nil {
			return err
		}
		for k, v := range api.Keys {
			ac.Base[k] = v
		}
	}
	return nil
}

func (ac *ApiCheck) LoadApi(vers ...string) error {
	for _, ver := range vers {
		api, err := LoadApi(ver)
		if err != nil {
			return err
		}
		for k, _ := range api.Keys {
			if _, ok := ac.Base[k]; ok {
				delete(api.Keys, k)
			}
		}
		ac.Apis = append(ac.Apis, api)
	}
	return nil
}

func (ac *ApiCheck) FincApis(name string) (vers []string) {
	for _, api := range ac.Apis {
		if _, ok := api.Keys[name]; ok {
			vers = append(vers, api.Ver)
		}
	}
	return
}

func (ac *ApiCheck) ApiVers() (vers []string) {
	for _, api := range ac.Apis {
		vers = append(vers, api.Ver)
	}
	return
}

func (ac *ApiCheck) CheckConstType(name string) KeyType {
	if typ, ok := ac.Base[name]; ok {
		return typ
	}
	for _, api := range ac.Apis {
		if typ, ok := api.Keys[name]; ok {
			return typ
		}
	}
	return Normal
}
