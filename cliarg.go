package cliarg

import (
	"os"
	"strconv"
	"strings"
)

type aliasOption struct {
	name         string
	defaultValue interface{}
}

type Arguments struct {
	Commands []string
	Options  map[string]interface{}
	Aliases  map[string]aliasOption
}

func NewArguments() *Arguments {
	return &Arguments{
		Commands: []string{},
		Options:  map[string]interface{}{},
		Aliases:  make(map[string]aliasOption),
	}
}

func (a *Arguments) Alias(flag string, value string, defaultValue interface{}) {
	a.Aliases[flag] = aliasOption{
		name:         value,
		defaultValue: defaultValue,
	}

	// set defaultValue
	if defaultValue != nil {
		a.Options[value] = defaultValue
	}
}

func (a *Arguments) Parse() {
	args := os.Args[1:]
	size := len(args)

	for i := 0; i < size; i++ {
		r := []byte(args[i])

		if string(r[0]) != "-" {
			a.Commands = append(a.Commands, args[i])
			continue
		}

		if string(r[1]) == "-" {
			splitted := strings.Split(args[i], "=")

			if len(splitted) > 1 {
				b := []byte(splitted[0])
				a.Options[string(b[2:])] = splitted[1]
			} else {
				b := []byte(args[i])
				a.Options[string(b[2:])] = ""
			}

		} else if alias, ok := a.Aliases[string(r[1])]; ok {
			// spaced value
			if len(args[i]) == 2 {
				if alias.defaultValue == nil {
					a.Options[alias.name] = true
				} else if _, ok := alias.defaultValue.(bool); ok {
					a.Options[alias.name] = true
				} else if i+1 < size {
					a.Options[alias.name] = args[i+1]
					i++
				}
				// alias with value
			} else {
				b := []byte(args[i])
				a.Options[alias.name] = string(b[2:])
			}
		}
	}
}

func (a *Arguments) GetOptionAsString(sign string) (string, bool) {
	if value, ok := a.Options[sign]; ok {
		return value.(string), true
	}
	return "", false
}

func (a *Arguments) GetOptionAsInt(sign string) (int, bool) {
	if value, ok := a.Options[sign]; ok {
		if v, err := strconv.Atoi(value.(string)); err != nil {
			return 0, false
		} else {
			return v, true
		}
	}
	return 0, false
}

func (a *Arguments) GetOptionAsBool(sign string) (bool, bool) {
	if value, ok := a.Options[sign]; ok {
		return value.(bool), true
	}
	return false, false
}

func (a *Arguments) GetOptionAsFloat(sign string) (float64, bool) {
	if value, ok := a.Options[sign]; ok {
		if v, err := strconv.ParseFloat(value.(string), 32); err != nil {
			return 0.0, false
		} else {
			return v, true
		}
	}
	return 0.0, false
}

func (a *Arguments) GetOption(sign string) (v interface{}, ok bool) {
	v, ok = a.Options[sign]
	return
}

func (a *Arguments) GetCommands() []string {
	return a.Commands
}

func (a *Arguments) GetCommandSize() int {
	return len(a.Commands)
}

func (a *Arguments) GetCommandAt(index int) (v string, ok bool) {
	if len(a.Commands) > index-1 {
		v = a.Commands[index-1]
		ok = true
	}

	return
}
