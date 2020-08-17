package kacker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"reflect"
	"text/template"
)

// Variable represents a template variable that will be provided to the kickstart
// template.
type Variable struct {
	Name      string   `yaml:"name"`
	Value     string   `yaml:"value"`
	Values    []string `yaml:"values"`
	URL       string   `yaml:"url"`
	URLs      []string `yaml:"urls"`
	Fragment  string   `yaml:"fragment"`
	Fragments []string `yaml:"fragments"`
}

// Kickstart collects the kickstart template file and any variable
// substitutions.
type Kickstart struct {
	From      []string   `yaml:"from"`
	Variables []Variable `yaml:"variables"`
}

type resolvedVariables struct {
	Variables map[string]string
	Arrays    map[string][]string
}

// Require panics if the given variable has not been resolved as a string. It can
// be used inline.
func (t *resolvedVariables) Require(name string) string {
	ret, ok := t.Variables[name]
	if !ok {
		panic(fmt.Sprintf("Unresolved variable %s", name))
	}
	return ret
}

// Optional grabs the given variable and returns it, otherwise returns fallback.
// This is equivalent to:
//    `{{ with index .Variables "name" }}{{ . }}{{ else }} "fallback" {{ end }}`
func (t *resolvedVariables) Optional(name string, fallback string) string {
	ret, ok := t.Variables[name]
	if !ok {
		return fallback
	}
	return ret
}

// RequireArray panics if the given variable has not been resolved as an array.
// It can be used inline.
func (t *resolvedVariables) RequireArray(name string) []string {
	ret, ok := t.Arrays[name]
	if !ok {
		panic(fmt.Sprintf("Unresolved array %s", name))
	}
	return ret
}

func (v Variable) validate() error {
	if len(v.Name) == 0 {
		return fmt.Errorf("Variable name cannot be empty")
	}
	v.Name = ""

	ref := reflect.ValueOf(v)
	filled := false
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Len() > 0 {
			if filled {
				return fmt.Errorf("Cannot specify more than one value for variable")
			}
			filled = true
		}
	}

	if !filled {
		return fmt.Errorf("Must specify at least one value")
	}
	return nil
}

func resolveSlice(sl []string, fn func(string) (string, error)) ([]string, error) {
	ret := []string{}
	for _, elem := range sl {
		app, err := fn(elem)
		if err != nil {
			return nil, err
		}
		ret = append(ret, app)
	}
	return ret, nil
}

const resolveErrFmt = "Error resolving variable %s: %s"

func resolveVariables(vars []Variable) (*resolvedVariables, error) {
	var ret resolvedVariables
	ret.Arrays = make(map[string][]string)
	ret.Variables = make(map[string]string)

	for _, elem := range vars {
		err := elem.validate()
		if err != nil {
			return nil, fmt.Errorf(resolveErrFmt, elem.Name, err.Error())
		}

		_, okv := ret.Variables[elem.Name]
		_, oka := ret.Arrays[elem.Name]
		if okv || oka {
			return nil, fmt.Errorf("Variable %s is redefined", elem.Name)
		}

		if len(elem.Value) > 0 || len(elem.URL) > 0 || len(elem.Fragment) > 0 {
			var content string
			switch {
			case len(elem.Value) > 0:
				content = elem.Value
			case len(elem.URL) > 0:
				content, err = resolveURL(elem.URL)
			case len(elem.Fragment) > 0:
				content, err = resolveFragment(elem.Fragment)
			}
			ret.Variables[elem.Name] = content
		} else {
			var content []string
			switch {
			case len(elem.Values) > 0:
				content = elem.Values
			case len(elem.URLs) > 0:
				content, err = resolveSlice(elem.URLs, resolveURL)
			case len(elem.Fragments) > 0:
				content, err = resolveSlice(elem.Fragments, resolveFragment)
			}
			ret.Arrays[elem.Name] = content
		}
		if err != nil {
			return nil, fmt.Errorf(resolveErrFmt, elem.Name, err.Error())
		}
	}
	return &ret, nil
}

func (kc *Kickstart) ResolveTempFile() (string, error) {
	return resolveToTempFile("./", "kacker-ks-*.cfg", kc.Resolve)
}

func (kc *Kickstart) Resolve(writer io.Writer) error {
	funcMap := template.FuncMap{
		"escape": escape,
	}
	basename := path.Base(kc.From[0])
	temp, err := template.New(basename).Funcs(funcMap).ParseFiles(kc.From...)
	if err != nil {
		return err
	}

	vars, err := resolveVariables(kc.Variables)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := temp.Execute(&buf, vars); err != nil {
		return err
	}
	if globalConf.VerboseLogging {
		log.Printf("Resulting kickstart file:\n%s\n", buf.String())
	}
	_, err = writer.Write(buf.Bytes())
	return err
}
