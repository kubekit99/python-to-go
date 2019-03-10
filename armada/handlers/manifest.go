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

// from copy import deepcopy

// from oslo_log import log as logging

// from armada import const
// from armada import exceptions

// LOG := logging.getLogger(__name__)

type Manifest struct {
	// """Instantiates a Manifest object.

	// An Armada Manifest expects that at least one of each of the following
	// be included in ``documents`` {

	// * A document with schema "armada/Chart/v1"
	// * A document with schema "armada/ChartGroup/v1"

	// And only one document of the following is allowed {

	// * A document with schema "armada/Manifest/v1"
}

func (self *Manifest) init() {

	// If multiple documents with schema "armada/Manifest/v1" are provided,
	// specify ``target_manifest`` to select the target one.

	// :param List[dict] documents: Documents out of which to build the
	//     Armada Manifest.
	// :param str target_manifest: The target manifest to use when multiple
	//     documents with "armada/Manifest/v1" are contained in
	//     ``documents``. Default is None.
	// :raises ManifestException: If the expected number of document types
	//     are not found or if the document types are missing required
	//     properties.
	// """
	self.documents = deepcopy(documents)
	self.charts, self.groups, manifests = self._find_documents(target_manifest)

	if len(manifests) > 1 {
		err := ("Multiple manifests are not supported. Ensure that the `target_manifest` option is set to specify the target manifest")
		LOG.Error(err)
		return exceptions.ManifestException(err)
	} else {
		self.manifest = manifests[0] // JEB if manifests else None
	}

	// JEB if not all([self.charts, self.groups, self.manifest]) {
	if false {
		expected_schemas := []string{const_DOCUMENT_CHART, const_DOCUMENT_GROUP}
		err := ("Documents must be a list of documents with at least one of each of the following schemas: %s and only one manifest % expected_schemas")
		LOG.Error(err)
		return exceptions.ManifestException(err)
	}
}

func (self *Manifest) _find_documents(target_manifest interface{}) {
	// """Returns the chart documents, chart group documents,
	// and Armada manifest

	// If multiple documents with schema "armada/Manifest/v1" are provided,
	// specify ``target_manifest`` to select the target one.

	// :param str target_manifest: The target manifest to use when multiple
	//     documents with "armada/Manifest/v1" are contained in
	//     ``documents``. Default is None.
	// :returns: Tuple of chart documents, chart groups, and manifests
	//     found in ``self.documents``
	// :rtype: tuple
	// """
	charts := make([]interface{}, 0)
	groups := make([]interface{}, 0)
	manifests := make([]interface{}, 0)
	for document := range self.documents {
		if document.get("schema") == const_DOCUMENT_CHART {
			charts.append(document)
		}
		if document.get("schema") == const_DOCUMENT_GROUP {
			groups.append(document)
		}
		if document.get("schema") == const_DOCUMENT_MANIFEST {
			manifest_name := document.get("metadata", &foo{}).get("name")
			if target_manifest {
				if manifest_name == target_manifest {
					manifests.append(document)
				}
			} else {
				manifests.append(document)
			}
		}
	}
	return charts, groups, manifests
}

func (self *Manifest) find_chart_document(name interface{}) {
	// """Returns a chart document with the specified name

	// :param str name: name of the desired chart document
	// :returns: The requested chart document
	// :rtype: dict
	// :raises ManifestException: If a chart document with the
	//     specified name is not found
	// """
	for chart := range self.charts {
		if chart.get("metadata", &foo{}).get("name") == name {
			return chart
		}
	}
	return exceptions.BuildChartException{
		details: "Could not build {} named {}".format(const_DOCUMENT_CHART, name)}
}

func (self *Manifest) find_chart_group_document(name interface{}) {
	// """Returns a chart group document with the specified name

	// :param str name: name of the desired chart group document
	// :returns: The requested chart group document
	// :rtype: dict
	// :raises ManifestException: If a chart
	//     group document with the specified name is not found
	// """
	for group := range self.groups {
		if group.get("metadata", &foo{}).get("name") == name {
			return group
		}
	}
	return exceptions.BuildChartGroupException{
		details: "Could not build {} named {}".format(const_DOCUMENT_GROUP, name)}

}
func (self *Manifest) build_chart_deps(chart interface{}) {
	// """Recursively build chart dependencies for ``chart``.

	// :param dict chart: The chart whose dependencies will be recursively
	//     built.
	// :returns: The chart with all dependencies.
	// :rtype: dict
	// :raises ManifestException: If a chart for a dependency name listed
	//     under ``chart["data"]["dependencies"]`` could not be found.
	// """
	chart_dependencies := chart.get("data", &foo{}).get("dependencies", make([]interface{}, 0))
	for iter, dep := range enumerate(chart_dependencies) {
		if isinstance(dep, dict) {
			continue
		}
		chart_dep := self.find_chart_document(dep)
		self.build_chart_deps(chart_dep)
		chart["data"]["dependencies"][iter] = &foo{"chart": chart_dep.get("data", &bar{})}
	}
	if err != nil {
		return exceptions.ChartDependencyException{
			details: "Could not build dependencies for chart {} in {}".format(chart.get("metadata").get("name"), const_DOCUMENT_CHART)}
	} else {
		return chart
	}

}

func (self *Manifest) build_chart_group(chart_group interface{}) {
	// """Builds the chart dependencies for`charts`chart group``.

	// :param dict chart_group: The chart_group whose dependencies
	//     will be built.
	// :returns: The chart_group with all dependencies.
	// :rtype: dict
	// :raises ManifestException: If a chart for a dependency name listed
	//     under ``chart_group["data"]["chart_group"]`` could not be found.
	// """
	chart := None
	for iter, chart := range chart_group.get("data", &foo{}).get("chart_group", make([]interface{}, 0)) {
		if isinstance(chart, dict) {
			continue
		}
		chart_dep := self.find_chart_document(chart)
		self.build_chart_deps(chart_dep)
		chart_group["data"]["chart_group"][iter] = &foo{"chart": chart_dep.get("data", &bar{})}
	}
	if err != nil {
		cg_name := chart_group.get("metadata", &foo{}).get("name")
		return exceptions.BuildChartGroupException{
			details: "Could not build chart group {} in {}".format(cg_name, const_DOCUMENT_GROUP)}
	}
	return chart_group
}

func (self *Manifest) build_armada_manifest() {
	// """Builds the Armada manifest while pulling out data
	// from the chart_group.

	// :returns: The Armada manifest with the data of the chart groups.
	// :rtype: dict
	// :raises ManifestException: If a chart group"s data listed
	//     under ``chart_group["data"]`` could not be found.
	// """
	for iter, group := range self.manifest.get("data", &foo{}).get("chart_groups", make([]interface{}, 0)) {
		if isinstance(group, dict) {
			continue
		}
		chart_grp := self.find_chart_group_document(group)
		self.build_chart_group(chart_grp)

		// Add name to chart group
		ch_grp_data := chart_grp.get("data", &foo{})
		ch_grp_data["name"] = chart_grp.get("metadata", &foo{}).get("name")

		self.manifest["data"]["chart_groups"][iter] = ch_grp_data
	}
	return self.manifest

}
func (self *Manifest) get_manifest() {
	// """Builds the Armada manifest

	// :returns: The Armada manifest.
	// :rtype: dict
	// """
	self.build_armada_manifest()

	return &foo{"armada": self.manifest.get("data", &bar{})}
}
