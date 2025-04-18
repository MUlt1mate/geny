package commands

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	CommandTypeProtoc = "protoc"
)

var re = regexp.MustCompile(`--(.*)_(opt|out)=(.*)`)

type (
	ProtocCommand struct {
		Type string
		Body ProtocBody
	}
	ProtocBody struct {
		Imports []string
		Plugins []ProtocPlugin
		Files   []string
	}

	ProtocPlugin struct {
		Name       string
		Path       string
		Parameters []ProtocPluginKV
	}
	ProtocPluginKV struct {
		Name  string
		Value string
	}
)

func ParseProtoc(input string) (command *ProtocCommand, err error) {
	command = &ProtocCommand{
		Type: CommandTypeProtoc,
	}
	var (
		parts       = strings.Split(input, " ")
		pluginIndex = make(map[string]int)
	)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "-I="),
			strings.HasPrefix(part, "-IPATH="),
			strings.HasPrefix(part, "--proto_path="):
			importParts := strings.Split(part, "=")
			if len(importParts) != 2 {
				return nil, fmt.Errorf("geny: invalid import parameter: %s in command %s", part, input)
			}
			command.Body.Imports = append(command.Body.Imports, importParts[1])
		case strings.HasSuffix(part, ".proto"):
			command.Body.Files = append(command.Body.Files, part)
		case strings.HasPrefix(part, "--"):
			var plugin *ProtocPlugin
			if plugin, err = parseProtocPluginArg(part); err != nil {
				return nil, fmt.Errorf("%w in command %s", err, input)
			}
			if index, ok := pluginIndex[plugin.Name]; ok {
				command.Body.Plugins[index].Parameters = append(command.Body.Plugins[index].Parameters, plugin.Parameters...)
				if plugin.Path != "" {
					command.Body.Plugins[index].Path = plugin.Path
				}
			} else {
				command.Body.Plugins = append(command.Body.Plugins, *plugin)
				pluginIndex[plugin.Name] = len(command.Body.Plugins) - 1
			}
		case strings.TrimSpace(part) == "", part == "protoc":
			continue
		default:
			return nil, fmt.Errorf("geny: unknown protoc command part: %s in command %s", part, input)
		}
	}
	return command, nil
}

func parseProtocPluginArg(part string) (*ProtocPlugin, error) {
	matches := re.FindStringSubmatch(part)
	if len(matches) != 4 {
		return nil, fmt.Errorf("geny: unknown protoc parameter format: %s", part)
	}
	plugin := &ProtocPlugin{Name: matches[1]}
	parameters := matches[3]
	parametersParts := strings.Split(parameters, ":")
	if len(parametersParts) > 2 {
		return nil, fmt.Errorf("geny: unknown protoc plugin parameter format: %s", part)
	}
	if matches[2] == "out" {
		if len(parametersParts) == 2 {
			plugin.Path = parametersParts[1]
		} else {
			plugin.Path = parametersParts[0]
		}
	}
	if matches[2] != "opt" && len(parametersParts) == 1 {
		return plugin, nil
	}
	pluginParameters := strings.Split(parametersParts[0], ",")
	for _, pluginParameter := range pluginParameters {
		if pluginParameter == "" {
			continue
		}
		pluginParametersKV := strings.Split(pluginParameter, "=")
		if len(pluginParametersKV) != 2 {
			return nil, fmt.Errorf("geny: unknown protoc plugin parameter format: %s", part)
		}
		plugin.Parameters = append(plugin.Parameters, ProtocPluginKV{Name: pluginParametersKV[0], Value: pluginParametersKV[1]})
	}
	return plugin, nil
}

func (s *ProtocCommand) String() string {
	var parts = []string{"protoc"}
	for _, importArg := range s.Body.Imports {
		parts = append(parts, "-I="+importArg)
	}
	for _, plugin := range s.Body.Plugins {
		var pluginParameters []string
		for _, kv := range plugin.Parameters {
			pluginParameters = append(pluginParameters, kv.Name+"="+kv.Value)
		}
		outValue := plugin.Path
		if len(pluginParameters) > 0 {
			outValue = strings.Join(pluginParameters, ",") + ":" + plugin.Path
		}
		parts = append(parts, fmt.Sprintf("--%s_out=%s", plugin.Name, outValue))
	}
	for _, file := range s.Body.Files {
		parts = append(parts, file)
	}
	return strings.Join(parts, " ")
}
