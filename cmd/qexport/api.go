package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func apipath(base string) string {
	return filepath.Join(os.Getenv("GOROOT"), "api", base)
}

var sym = regexp.MustCompile(`^pkg (\S+)\s?(.*)?, (?:(var|func|type|const)) ([A-Z]\w*)`)

type GoApi struct {
	Keys map[string]bool
	Ver  string
}

func LoadApi(ver string) (*GoApi, error) {
	f, err := os.Open(apipath(ver + ".txt"))
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)
	keys := make(map[string]bool)
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
			keys[key] = true
		}
	}
	return &GoApi{Ver: ver, Keys: keys}, nil
}

type ApiCheck struct {
	Base map[string]bool
	Apis []*GoApi
}

func NewApiCheck() *ApiCheck {
	ac := &ApiCheck{}
	ac.Base = make(map[string]bool)
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
			if ac.Base[k] {
				delete(api.Keys, k)
			}
		}
		ac.Apis = append(ac.Apis, api)
	}
	return nil
}

func (ac *ApiCheck) FincApis(name string) (vers []string) {
	for _, api := range ac.Apis {
		if api.Keys[name] {
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
