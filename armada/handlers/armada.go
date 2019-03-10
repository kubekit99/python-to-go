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

// from concurrent.futures import ThreadPoolExecutor, as_completed
// from oslo_config import cfg
// from oslo_log import log as logging

// from armada import const
// from armada.conf import set_current_chart
// from armada.exceptions import armada_exceptions
// from armada.exceptions import override_exceptions
// from armada.exceptions import source_exceptions
// from armada.exceptions import tiller_exceptions
// from armada.exceptions import validate_exceptions
// from armada.handlers.chart_deploy import ChartDeploy
// from armada.handlers.manifest import Manifest
// from armada.handlers.override import Override
// from armada.utils.release import release_prefixer
// from armada.utils import source

// LOG := logging.getLogger(__name__)
// CONF := cfg.CONF

// """
// This is the main Armada class handling the Armada
// workflows
// """

type Armada struct {
	// """
	// Initialize the Armada engine and establish a connection to Tiller.

	// :param List[dict] documents: Armada documents.
	// :param tiller: Tiller instance to use.
	// :param bool disable_update_pre: Disable pre-update Tiller operations.
	// :param bool disable_update_post: Disable post-update Tiller
	//     operations.
	// :param bool enable_chart_cleanup: Clean up unmanaged charts.
	// :param bool dry_run: Run charts without installing them.
	// :param bool force_wait: Force Tiller to wait until all charts are
	//     deployed, rather than using each chart"s specified wait policy.
	// :param int timeout: Specifies overall time in seconds that Tiller
	//     should wait for charts until timing out.
	// :param str target_manifest: The target manifest to run. Useful for
	//     specifying which manifest to run when multiple are available.
	// :param int k8s_wait_attempts: The number of times to attempt waiting
	//     for pods to become ready.
	// :param int k8s_wait_attempt_sleep: The time in seconds to sleep
	//     between attempts.
	// """
	enable_chart_cleanup interface{}
	dry_run              interface{}
	force_wait           interface{}
	tiller               *Tiller
	documents            interface{}
	manifest             *Manifest
	chart_cache          interface{}
	chart_deploy         interface{}
}

func (self *Armada) init() {
	// self.enable_chart_cleanup = enable_chart_cleanup
	// self.dry_run = dry_run
	// self.force_wait = force_wait
	// self.tiller = tiller
	// self.documents = Override{ documents, set_ovr, values}.update_manifests()
	// self.manifest = Manifest( self.documents, target_manifest).get_manifest()
	// self.chart_cache = &foo{}
	// self.chart_deploy = ChartDeploy{ disable_update_pre, disable_update_post, self.dry_run,
	// 	k8s_wait_attempts, k8s_wait_attempt_sleep, timeout, self.tiller}

}

func (self *Armada) pre_flight_ops() {
	// """Perform a series of checks and operations to ensure proper
	// deployment.
	// """
	LOG.Info("Performing pre-flight operations.")

	// Ensure Tiller is available and manifest is valid
	if !self.tiller.tiller_status() {
		return tiller_exceptions.TillerServicesUnavailableException()
	}

	// Clone the chart sources
	manifest_data := self.manifest.get(const_KEYWORD_ARMADA, &foo{})
	for group := range manifest_data.get(const_KEYWORD_GROUPS, make([]interface{}, 0)) {
		for ch := range group.get(const_KEYWORD_CHARTS, make([]interface{}, 0)) {
			self.get_chart(ch)
		}
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (self *Armada) get_chart(ch interface{}) {
	chart := ch.get("chart", &foo{})
	chart_source := chart.get("source", &foo{})
	location := chart_source.get("location")
	ct_type := chart_source.get("type")
	subpath := chart_source.get("subpath", ".")

	if ct_type == "local" {
		chart["source_dir"] = []string{location, subpath}
	} else if ct_type == "tar" {
		source_key := []string{ct_type, location}

		if !(stringInSlice(source_key, self.chart_cache)) {
			LOG.Info("Downloading tarball from: %s", location)

			if !CONF.certs {
				LOG.warn("Disabling server validation certs to extract charts")
				tarball_dir := source.get_tarball(location, false)
			} else {
				tarball_dir := source.get_tarball(location, CONF.certs)
				self.chart_cache[source_key] = tarball_dir
			}
		}
		chart["source_dir"] = []string{self.chart_cache.get(source_key), subpath}
	} else if ct_type == "git" {
		reference := chart_source.get("reference", "master")
		source_key := []string{ct_type, location, reference}

		if !(stringInSlice(source_key, self.chart_cache)) {
			auth_method := chart_source.get("auth_method")
			proxy_server := chart_source.get("proxy_server")

			logstr := "Cloning repo: {} from branch: {}".format(
				location, reference)
			if proxy_server {
				logstr = logstr + " proxy: {}".format(proxy_server)
			}
			if auth_method {
				logstr := logstr + " auth method: {}".format(auth_method)
			}
			LOG.Info(logstr)

			repo_dir := source.git_clone(
				location,
				reference,
				proxy_server,
				auth_method)

			self.chart_cache[source_key] = repo_dir
		}
		chart["source_dir"] = []string{self.chart_cache.get(source_key), subpath}
	} else {
		chart_name := chart.get("chart_name")
		return source_exceptions.ChartSourceException(ct_type, chart_name)
	}
	for dep := range ch.get("chart", &foo{}).get("dependencies", make([]interface{}, 0)) {
		self.get_chart(dep)
	}
}
func (self *Armada) sync() {
	// """
	// Synchronize Helm with the Armada Config(s)
	// """
	if self.dry_run {
		LOG.Info("Armada is in DRY RUN mode, no changes being made.")
	}
	msg := &foo{
		"install":   make([]interface{}, 0),
		"upgrade":   make([]interface{}, 0),
		"diff":      make([]interface{}, 0),
		"purge":     make([]interface{}, 0),
		"protected": make([]interface{}, 0),
	}

	// TODO: (gardlt) we need to break up this func into
	// a more cleaner format
	self.pre_flight_ops()

	known_releases := self.tiller.list_releases()

	manifest_data := self.manifest.get(const_KEYWORD_ARMADA, &foo{})
	prefix := manifest_data.get(const_KEYWORD_PREFIX)

	for chartgroup := range manifest_data.get(const_KEYWORD_GROUPS, make([]interface{}, 0)) {
		cg_name := chartgroup.get("name", "<missing name>")
		cg_desc := chartgroup.get("description", "<missing description>")
		cg_sequenced := chartgroup.get("sequenced", false) // JEB or self.force_wait

		LOG.Info("Processing ChartGroup: %s (%s), sequenced:=%s%s", cg_name,
			cg_desc, cg_sequenced, " (forced) if self.force_wait else ")

		// TODO(MarshM) Deprecate the `test_charts` key
		cg_test_all_charts := chartgroup.get("test_charts")

		cg_charts := chartgroup.get(const_KEYWORD_CHARTS, make([]interface{}, 0))
		charts := nil // JEB map(lambda x: x.get("chart", &foo{}), cg_charts)

		// JEB: deploy_chart was inlined here

		results := make([]interface{}, 0)
		failures := make([]interface{}, 0)

		// JEB: handle_result was inlined here

		if cg_sequenced {
			for chart := range charts {
				// JEB
				// if (handle_result(chart, lambda: deploy_chart(chart))) {
				//     break
				// }
			}
		} else {
			// JEB
			// with ThreadPoolExecutor(
			//         max_workers:=len(cg_charts)) as executor {
			//     future_to_chart := {
			//         executor.submit(deploy_chart, chart) { chart
			//         for chart in charts
			//     }

			//     for future in as_completed(future_to_chart) {
			//         chart := future_to_chart[future]
			//         handle_result(chart, future.result)
			//     }
			// }
		}

		if failures {
			LOG.Error("Chart deploy(s) failed: %s", failures)
			return armada_exceptions.ChartDeployException(failures)
		}
		for result := range results {
			for k, v := range result.items() {
				msg[k].append(v)
			}
		}

		// End of Charts in ChartGroup
		LOG.Info("All Charts applied in ChartGroup %s.", cg_name)
	}
	self.post_flight_ops()

	if self.enable_chart_cleanup {
		self._chart_cleanup(
			prefix,
			self.manifest[const_KEYWORD_ARMADA][const_KEYWORD_GROUPS], msg)
	}

	LOG.Info("Done applying manifest.")
	return msg
}

func (self *Armada) deploy_chart(chart interface{}) {
	set_current_chart(chart)
	return self.chart_deploy.execute(chart, cg_test_all_charts, prefix, known_releases)
	set_current_chart(None)
}

func (self *Armada) handle_result(chart interface{}, get_result interface{}) {
	// Returns whether or not there was a failure
	name := chart["chart_name"]
	result := get_result()
	if err != nil {
		LOG.exception("Chart deploy [{}] failed".format(name))
		failures.append(name)
		return true
	} else {
		results.append(result)
		return false
	}
}

func (self *Armada) post_flight_ops() {
	// """
	// Operations to run after deployment process has terminated
	// """
	LOG.Info("Performing post-flight operations.")

	// Delete temp dirs used for deployment
	for chart_dir := range self.chart_cache.values() {
		LOG.debug("Removing temp chart directory: %s", chart_dir)
		source.source_cleanup(chart_dir)
	}

}
func (self *Armada) _chart_cleanup(prefix interface{}, charts interface{}, msg interface{}) {
	LOG.Info("Processing chart cleanup to remove unspecified releases.")

	valid_releases := make([]interface{}, 0)
	for gchart := range charts {
		for chart := range gchart.get(const_KEYWORD_CHARTS, make([]interface{}, 0)) {
			valid_releases.append(
				release_prefixer(prefix,
					chart.get("chart", &foo{}).get("release")))

			// JEB actual_releases := [x.name for x in self.tiller.list_releases()]
			release_diff := list(set(actual_releases) - set(valid_releases))

			for release := range release_diff {
				if release.startswith(prefix) {
					LOG.Info("Purging release %s as part of chart cleanup.",
						release)
					self.tiller.uninstall_release(release)
					msg["purge"].append(release)
				}
			}
		}
	}
}
