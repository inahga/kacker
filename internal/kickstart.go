package kacker

import (
	"fmt"
	"reflect"
)

// Variable represents a template variable that will be provided to the kickstart
// template.
type Variable struct {
	Name      string   `yaml:"name"`
	Value     string   `yaml:"value"`
	Values    []string `yaml:"values"`
	URL       string   `yaml:"url"`
	Fragment  string   `yaml:"fragment"`
	Fragments []string `yaml:"fragments"`
}

// KickstartCustomization collects the kickstart template file and any variable substitutions.
type KickstartCustomization struct {
	From      []string   `yaml:"from"`
	Variables []Variable `yaml:"variables"`
}

type resolvedVariables struct {
	Variables map[string]string
	Arrays    map[string][]string
}

func (t *resolvedVariables) GetOptionalVariable(name string, fallback string) string {
	ret, ok := t.Variables[name]
	if !ok {
		return fallback
	}
	return ret
}

func (t *resolvedVariables) GetVariable(name string) string {
	ret, ok := t.Variables[name]
	if !ok {
		panic(fmt.Sprintf("Unresolved variable %s", name))
	}
	return ret
}

func (t *resolvedVariables) GetOptionalArray(name string) []string {
	ret, ok := t.Arrays[name]
	if !ok {
		return []string{}
	}
	return ret
}

func (t *resolvedVariables) GetArray(name string) []string {
	ret, ok := t.Arrays[name]
	if !ok {
		panic(fmt.Sprintf("Unresolved array %s", name))
	}
	return ret
}

func resolveFragment(fragment string) string {
	return ""
}

func resolveURL(url string) string {
	return ""
}

// IsValid validates whether a variable has a name and only one of the other
// fields present
func (v Variable) IsValid() (bool, error) {
	if len(v.Name) == 0 {
		return false, fmt.Errorf("Variable name cannot be empty")
	}
	v.Name = ""

	ref := reflect.ValueOf(v)
	filled := false
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Len() > 0 {
			if filled {
				return false, fmt.Errorf("Cannot specify more than one value for variable")
			}
			filled = true
		}
	}

	if !filled {
		return false, fmt.Errorf("Must specify at least one value")
	}
	return true, nil
}

func resolveVariables(vars []Variable) (*resolvedVariables, error) {
	var ret resolvedVariables
	ret.Arrays = make(map[string][]string)
	ret.Variables = make(map[string]string)

	for _, elem := range vars {
		_, err := elem.IsValid()
		if err != nil {
			return nil, fmt.Errorf("Error parsing variable %s: %s", elem.Name, err.Error())
		}

		_, okv := ret.Variables[elem.Name]
		_, oka := ret.Arrays[elem.Name]
		if okv || oka {
			return nil, fmt.Errorf("Variable %s is redefined", elem.Name)
		}

		switch {
		case len(elem.Value) > 0:
			ret.Variables[elem.Name] = elem.Value
		case len(elem.Values) > 0:
			ret.Arrays[elem.Name] = elem.Values
		case len(elem.URL) > 0:
			ret.Variables[elem.Name] = resolveURL(elem.URL)
		case len(elem.Fragment) > 0:
			ret.Variables[elem.Name] = resolveFragment(elem.Fragment)
		case len(elem.Fragments) > 0:
			for _, frag := range elem.Fragments {
				ret.Arrays[elem.Name] = append(ret.Arrays[elem.Name], resolveFragment(frag))
			}
		}
	}
	return &ret, nil
}

// ResolveKickstart takes a customization and applies the given variables to the
// given templates. The result is a complete kickstart file.
func ResolveKickstart(kc *KickstartCustomization) string {
	// get templates

	// get resolved variables

	// apply variables to template and return result

	return ""
}
