package main

import (
	"io"
	"reflect"
	"strings"
)

// InOutPlugins struct for holding references to plugins
type InOutPlugins struct {
	Inputs  []io.Reader
	Outputs []io.Writer
}

// Plugins holds all the plugin objects
var Plugins *InOutPlugins = new(InOutPlugins)

// extractLimitOptions detects if plugin get called with limiter support
// Returns address and limit
func extractLimitOptions(options string) (string, string) {
	split := strings.Split(options, "|")

	if len(split) > 1 {
		return split[0], split[1]
	}

	return split[0], ""
}

// Automatically detects type of plugin and initialize it
//
// See this article if curious about relfect stuff below: http://blog.burntsushi.net/type-parametric-functions-golang
func registerPlugin(constructor interface{}, options ...interface{}) {
	vc := reflect.ValueOf(constructor)

	// Pre-processing options to make it work with reflect
	vo := []reflect.Value{}
	for _, oi := range options {
		vo = append(vo, reflect.ValueOf(oi))
	}

	// Removing limit options from path
	path, limit := extractLimitOptions(vo[0].String())

	// Writing value back without limiter "|" options
	vo[0] = reflect.ValueOf(path)

	// Calling our constructor with list of given options
	plugin := vc.Call(vo)[0].Interface()
	pluginWrapper := plugin

	if limit != "" {
		pluginWrapper = NewLimiter(plugin, limit)
	} else {
		pluginWrapper = plugin
	}

	if _, ok := plugin.(io.Reader); ok {
		Plugins.Inputs = append(Plugins.Inputs, pluginWrapper.(io.Reader))
	}

	if _, ok := plugin.(io.Writer); ok {
		Plugins.Outputs = append(Plugins.Outputs, pluginWrapper.(io.Writer))
	}
}

// InitPlugins specify and initialize all available plugins
func InitPlugins() {
	for _, options := range Settings.inputDummy {
		registerPlugin(NewDummyInput, options)
	}

	for _, options := range Settings.outputDummy {
		registerPlugin(NewDummyOutput, options)
	}

	for _, options := range Settings.inputRAW {
		registerPlugin(NewRAWInput, options)
	}

	for _, options := range Settings.inputTCP {
		registerPlugin(NewTCPInput, options)
	}

	for _, options := range Settings.outputTCP {
		registerPlugin(NewTCPOutput, options)
	}

	for _, options := range Settings.inputFile {
		registerPlugin(NewFileInput, options)
	}

	for _, options := range Settings.outputFile {
		registerPlugin(NewFileOutput, options)
	}

	for _, options := range Settings.inputHTTP {
		registerPlugin(NewHTTPInput, options)
	}

	for _, options := range Settings.outputHTTP {
		registerPlugin(NewHTTPOutput, options, &Settings.outputHTTPConfig)
	}
}
