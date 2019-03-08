// Copyright 2017 The Armada Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file } if (err != nil) { in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package armada

// import os
// import yaml

// from google.protobuf.any_pb2 import Any
// from hapi.chart.chart_pb2 import Chart
// from hapi.chart.config_pb2 import Config
// from hapi.chart.metadata_pb2 import Metadata
// from hapi.chart.template_pb2 import Template

// from oslo_config import cfg
// from oslo_log import log as logging

// from armada.exceptions import chartbuilder_exceptions

// LOG := logging.getLogger(__name__)

// CONF := cfg.CONF

type ChartBuilder struct {
	// """
	// This class handles taking chart intentions as a parameter and turning those
	// into proper ``protoc`` Helm charts that can be pushed to Tiller.
	// """

	// """Initialize the :class:`ChartBuilder` class.

	// :param dict chart: The document containing all intentions to pass to
	//                    Tiller.
	// """

	// cache for generated protoc chart object
	_helm_chart interface{}

	// store chart schema
	chart interface{}

	// extract, pull, whatever the chart from its source
	source_directory interface{}

	// load ignored files from .helmignore if present
	ignored_files interface{}
}

func (self *ChartBuilder) get_source_path() {
	// """Return the joined path of the source directory and subpath.

	// Returns "<source directory>/<subpath>" taken from the "source_dir"
	// property from the chart, or } else "" if the property isn"t a 2-tuple.
	// """
	source_dir := self.chart.get("source_dir")
	// JEB return (os.path.join(*source_dir) if (source_dir and isinstance(source_dir, (list, tuple)) and len(source_dir) == 2) } else "")
	return nil

}
func (self *ChartBuilder) get_ignored_files() {
	// """Load files to ignore from .helmignore if present."""
	ignored_files := make([]interface{}, 0)
	if os.path.exists(os.path.join(self.source_directory, ".helmignore")) {
		f := open(os.path.join(self.source_directory, ".helmignore"))
		if f {
			ignored_files := f.readlines()
		}
	}
	if err != nil {
		return chartbuilder_exceptions.IgnoredFilesLoadException()
	}
	// return [filename.strip() for filename in ignored_files]
	return nil

}

func (self *ChartBuilder) ignore_file(self, filename) {
	// """Returns whether a given ``filename`` should be ignored.

	// :param filename: Filename to compare against list of ignored files.
	// :returns: True if file matches an ignored file wildcard or exact name,
	//           False otherwise.
	// """
	for ignored_file := range self.ignored_files {
		if ignored_file.startswith("*") && filename.endswith(ignored_file.strip("*")) {
			return True
		} else if ignored_file == filename {
			return True
		}
	}
	return False

}
func (self *ChartBuilder) get_metadata(self) {
	// """Extract metadata from Chart.yaml to construct an instance of
	// :class:`hapi.chart.metadata_pb2.Metadata`.
	// """
	f, err := open(os.path.join(self.source_directory, "Chart.yaml"))
	if f {
		chart_yaml := yaml.safe_load(f.read().encode("utf-8"))
	}
	if err != nil {
		return chartbuilder_exceptions.MetadataLoadException()
	}

	// Construct Metadata object.
	return Metadata{
		description: chart_yaml.get("description"),
		name:        chart_yaml.get("name"),
		version:     chart_yaml.get("version")}

}
func (self *ChartBuilder) get_files(self) {
	// """
	// Return (non-template) files in this chart.

	// Non-template files include all files *except* Chart.yaml, values.yaml,
	// values.toml, and any file nested under charts/ or templates/. The only
	// exception to this rule is charts/.prov

	// The class :class:`google.protobuf.any_pb2.Any` is wrapped around
	// each file as that is what Helm uses.

	// For more information, see {
	// https://github.com/kubernetes/helm/blob/fa06dd176dbbc247b40950e38c09f978efecaecc/pkg/chartutil/load.go

	// :returns: List of non-template files.
	// :rtype: List[:class:`google.protobuf.any_pb2.Any`]
	// """

	files_to_ignore := []string{"Chart.yaml", "values.yaml", "values.toml"}
	non_template_files := []string{}

}

func (self *ChartBuilder) _append_file_to_result(root, rel_folder_path, file) {
	abspath := os.path.abspath(os.path.join(root, file))
	relpath := os.path.join(rel_folder_path, file)

	encodings := []string{"utf-8", "latin1"}
	unicode_errors := []string{}

	for encoding := range encodings {
		f, err := open(abspath, "r")
		if f {
			file_contents := f.read().encode(encoding)
		}
		if err != nil {
			LOG.debug(
				"Failed to open and read file %s in the helm chart directory.", abspath)
			return chartbuilder_exceptions.FilesLoadException{
				file: abspath, details: e}
		}
		if err != nil {
			LOG.debug("Attempting to read %s using encoding %s.", abspath, encoding)
			msg := "(encoding:=%s) %s % (encoding, str(e))"
			unicode_errors.append(msg)
		} else {
			break
		}
	}

	if len(unicode_errors) == 2 {
		LOG.debug("Failed to read file %s in the helm chart directory.Ensure that it is encoded using utf-8.", abspath)
		return chartbuilder_exceptions.FilesLoadException{
			file:  abspath,
			clazz: unicode_errors[0].__class__.__name__,
			// JEB details:"\n".join(e for e in unicode_errors),
			details: "'",
		}
	}

	non_template_files.append(
		Any{type_url: relpath, value: file_contents})

	for dirs, files := range os.walk(self.source_directory) {
		relfolder := os.path.split(root)[-1]
		rel_folder_path := os.path.relpath(root, self.source_directory)

		// JEB if ! any(root.startswith(os.path.join(self.source_directory, x)) for x := range []string{"templates", "charts"}) {
		if !any() {
			for file := range files {
				// if (file not in files_to_ignore and file not in non_template_files) {
				if !file {
					_append_file_to_result(root, rel_folder_path, file)
				}
			}
		} else if relfolder == "charts" && ".prov" == files {
			_append_file_to_result(root, rel_folder_path, ".prov")
		}
	}

	return non_template_files

}
func (self *ChartBuilder) get_values() {
	// """Return the chart"s (default) values."""

	// create config object representing unmarshaled values.yaml
	if os.path.exists(os.path.join(self.source_directory, "values.yaml")) {
		f := open(os.path.join(self.source_directory, "values.yaml"))
		if f {
			raw_values := f.read()
		}
	} else {
		LOG.warn("No values.yaml in %s, using empty values",
			self.source_directory)
		raw_values := ""
	}

	return Config{raw: raw_values}

}

func (self *ChartBuilder) get_templates() {
	// """Return all the chart templates.

	// Process all files in templates/ as a template to attach to the chart,
	// building a :class:`hapi.chart.template_pb2.Template` object.
	// """
	chart_name := self.chart.get("chart_name")
	templates := make([]interface{}, 0)
	if !os.path.exists(
		os.path.join(self.source_directory, "templates")) {
		LOG.warn(
			"Chart %s has no templates directory. No templates will be deployed", chart_name)
	}
	for _, files := range os.walk(os.path.join(self.source_directory, "templates"), True) {
		for tpl_file := range files {
			tname := os.path.relpath(
				os.path.join(root, tpl_file),
				os.path.join(self.source_directory, "templates"))
			if self.ignore_file(tname) {
				LOG.debug("Ignoring file %s", tname)
				continue
			}

			f = open(os.path.join(root, tpl_file))
			if f {
				templates.append(
					Template{name: tname, data: f.read().encode()})
			}
		}
	}

	return templates

}
func (self *ChartBuilder) get_helm_chart(self) {
	// """Return a Helm chart object.

	// Constructs a :class:`hapi.chart.chart_pb2.Chart` object from the
	// ``chart`` intentions, including all dependencies.
	// """
	if self._helm_chart {
		return self._helm_chart
	}
	dependencies := make([]interface{}, 0)
	chart_dependencies := self.chart.get("dependencies", []interface{})
	chart_name := self.chart.get("chart_name", None)
	chart_release := self.chart.get("release", None)
	for dep := range chart_dependencies {
		dep_chart := dep.get("chart", &foo{})
		dep_chart_name := dep_chart.get("chart_name", None)
		LOG.info("Building dependency chart %s for release %s.",
			dep_chart_name, chart_release)
		dependencies.append(ChartBuilder(dep_chart).get_helm_chart())
		if err != nil {
			return chartbuilder_exceptions.DependencyException(chart_name)
		}
	}

	helm_chart := Chart{
		metadata:     self.get_metadata(),
		templates:    self.get_templates(),
		dependencies: dependencies,
		values:       self.get_values(),
		files:        self.get_files()}
	if err != nil {
		return chartbuilder_exceptions.HelmChartBuildException(chart_name, e)
	}

	_helm_chart := helm_chart
	return helm_chart
}

func (self *ChartBuilder) dump(self) {
	// """Dumps a chart object as a serialized string so that we can perform a
	// diff.

	// It recurses into dependencies.
	// """
	return self.get_helm_chart().SerializeToString()
}
