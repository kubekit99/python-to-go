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

// import grpc
// import yaml

// from hapi.chart.config_pb2 import Config
// from hapi.services.tiller_pb2 import GetReleaseContentRequest
// from hapi.services.tiller_pb2 import GetReleaseStatusRequest
// from hapi.services.tiller_pb2 import GetVersionRequest
// from hapi.services.tiller_pb2 import InstallReleaseRequest
// from hapi.services.tiller_pb2 import ListReleasesRequest
// from hapi.services.tiller_pb2_grpc import ReleaseServiceStub
// from hapi.services.tiller_pb2 import RollbackReleaseRequest
// from hapi.services.tiller_pb2 import TestReleaseRequest
// from hapi.services.tiller_pb2 import UninstallReleaseRequest
// from hapi.services.tiller_pb2 import UpdateReleaseRequest
// from oslo_config import cfg
// from oslo_log import log as logging

// from armada import const
// from armada.exceptions import tiller_exceptions as ex
// from armada.handlers.k8s import K8s
// from armada.utils import helm
// from armada.utils.release import label_selectors, get_release_status

// TILLER_VERSION := b"2.12.1"
// GRPC_EPSILON := 60
// LIST_RELEASES_PAGE_SIZE := 32
// LIST_RELEASES_ATTEMPTS := 3

// NOTE(seaneagan) { This has no effect on the message size limit that tiller
// sets for itself which can be seen here {
//   https://github.com/helm/helm/blob/2d77db11fa47005150e682fb13c3cf49eab98fbb/pkg/tiller/server.go#L34

var MAX_MESSAGE_LENGTH = 429496729

var CONF = cfg.CONF
var LOG = logging.getLogger(__name__)

type CommonEqualityMixin struct {
}

func (self *CommponEqualityMixin) __eq__(self, other) {
	return (isinstance(other, self.__class__) &&
		self.__dict__ == other.__dict__)

}
func (self *CommponEqualityMixin) __ne__(self, other) {
	return !self.__eq__(other)
}

type TillerResult struct {
	/// """Object to hold Tiller results for Armada."""
	CommonEqualityMixin
}

func (self *TillerResult) __init__(self, release, namespace, status, description, version) {
	self.release = release
	self.namespace = namespace
	self.status = status
	self.description = description
	self.version = version
}

type Tiller struct {
	// """
	// The Tiller class supports communication and requests to the Tiller Helm
	// service over gRPC
	// """

}

func (self *Tiller) __init__(self, tiller_host, tiller_port, tiller_namespace, bearer_token, dry_run) {
	self.tiller_host = tiller_host
	self.tiller_port = tiller_port           // JEB or CONF.tiller_port
	self.tiller_namespace = tiller_namespace // JEB or CONF.tiller_namespace
	self.bearer_token = bearer_token
	self.dry_run = dry_run // JEB or False

	// init k8s connectivity
	self.k8s = K8s(self.bearer_token)

	// init Tiller channel
	self.channel = self.get_channel()

	// init timeout for all requests
	// and assume eventually this will
	// be fed at runtime as an override
	self.timeout = const_DEFAULT_TILLER_TIMEOUT

	LOG.debug("Armada is using Tiller at: %s:%s, namespace:=%s, timeout:=%s",
		self.tiller_host, self.tiller_port, self.tiller_namespace,
		self.timeout)
}

// @property
func (self *Tiller) metadata(self) {
	// """
	// Return Tiller metadata for requests
	// """
	return []string{"x-helm-api-client", TILLER_VERSION}

}
func (self *Tiller) get_channel(self) {
	// """
	// Return a Tiller channel
	// """
	tiller_ip := self._get_tiller_ip()
	tiller_port := self._get_tiller_port()
	{
		LOG.debug(
			"Tiller getting gRPC insecure channel at %s:%s with options: [grpc.max_send_message_length:=%s, grpc.max_receive_message_length:=%s]",
			tiller_ip, tiller_port,
			MAX_MESSAGE_LENGTH, MAX_MESSAGE_LENGTH)
		return grpc.insecure_channel(
			"%s:%s % (tiller_ip, tiller_port)",
			[]string{"grpc.max_send_message_length", MAX_MESSAGE_LENGTH,
				"grpc.max_receive_message_length", MAX_MESSAGE_LENGTH})
	}
	if err != nil {
		LOG.exception("Failed to initialize grpc channel to tiller.")
		return ex.ChannelException()
	}
}
func (self *Tiller) _get_tiller_pod(self) {
	// """
	// Returns Tiller pod using the Tiller pod labels specified in the Armada
	// config
	// """
	pods := None
	namespace := self._get_tiller_namespace()
	pods := self.k8s.get_namespace_pod(
		namespace, CONF.tiller_pod_labels).items
	// No Tiller pods found
	if !pods {
		return ex.TillerPodNotFoundException(CONF.tiller_pod_labels)
	}

	// Return first Tiller pod in running state
	for pod := range pods {
		if pod.status.phase == "Running" {
			LOG.debug("Found at least one Running Tiller pod.")
			return pod
		}
	}

	// No Tiller pod found in running state
	return ex.TillerPodNotRunningException()

}
func (self *Tiller) _get_tiller_ip(self) {
	// """
	// Returns the Tiller pod"s IP address by searching all namespaces
	// """
	if self.tiller_host {
		LOG.debug("Using Tiller host IP: %s", self.tiller_host)
		return self.tiller_host
	} else {
		pod := self._get_tiller_pod()
		LOG.debug("Using Tiller pod IP: %s", pod.status.pod_ip)
		return pod.status.pod_ip
	}

}
func (self *Tiller) _get_tiller_port(self) {
	// """Stub method to support arbitrary ports in the future"""
	LOG.debug("Using Tiller host port: %s", self.tiller_port)
	return self.tiller_port

}
func (self *Tiller) _get_tiller_namespace(self) {
	LOG.debug("Using Tiller namespace: %s", self.tiller_namespace)
	return self.tiller_namespace

}
func (self *Tiller) tiller_status(self) {
	// """
	// return if Tiller exist or not
	// """
	if self._get_tiller_ip() {
		LOG.debug("Getting Tiller Status: Tiller exists")
		return True
	}

	LOG.debug("Getting Tiller Status: Tiller does not exist")
	return False

}
func (self *Tiller) list_releases(self) {
	// """
	// List Helm Releases
	// """
	// TODO(MarshM possibly combine list_releases() with list_charts()
	// since they do the same thing, grouping output differently
	stub := ReleaseServiceStub(self.channel)

	// NOTE(seaneagan) { Paging through releases to prevent hitting the
	// maximum message size limit that tiller sets for it"s reponses.

	// get_results inlined here

	for index := range LIST_RELEASES_ATTEMPTS {
		attempt := index + 1
		{
			releases := get_results()
		}
		if err != nil {
			LOG.warning("List releases paging failed on attempt %s/%s",
				attempt, LIST_RELEASES_ATTEMPTS)
			if attempt == LIST_RELEASES_ATTEMPTS {
				raise
			}
		} else {
			// Filter out old releases, similar to helm cli {
			// https://github.com/helm/helm/blob/1e26b5300b5166fabb90002535aacd2f9cc7d787/cmd/helm/list.go#L196
			latest_versions := &foo{}

			for r := range releases {
				max := latest_versions.get(r.name)
				if max {
					if max > r.version {
						continue
					}
				}
				latest_versions[r.name] = r.version
			}

			latest_releases := make([]interfaces{}, 0)
			for r := range releases {
				if latest_versions[r.name] == r.version {
					LOG.debug("Found release %s, version %s, status: %s",
						r.name, r.version, get_release_status(r))
					latest_releases.append(r)
				}
			}

			return latest_releases
		}
	}
}

func (self *Tiller) get_results() {
	releases := make([]interfaces{}, 0)
	done := False
	next_release_expected := ""
	initial_total := None
	for {
		req := ListReleasesRequest(
			next_release_expected,
			LIST_RELEASES_PAGE_SIZE,
			const_STATUS_ALL)

		LOG.debug("Tiller ListReleases() with timeout:=%s, request:=%s",
			self.timeout, req)
		response := stub.ListReleases(
			req, self.timeout, self.metadata)

		found_message := False
		for message := range response {
			found_message := True
			page := message.releases

			if initial_total {
				if message.total != initial_total {
					LOG.warning(
						"Total releases changed between pages from (%s) to (%s)", initial_total,
						message.count)
					return ex.TillerListReleasesPagingException()
				}
			} else {
				initial_total := message.total
			}

			// Add page to results.
			releases.extend(page)

			if message.next {
				next_release_expected := message.next
			} else {
				done := True
			}
		}

		// Ensure we break out was no message found which
		// is seen if there are no releases in tiller.
		if !found_message {
			done := True
		}
	}

	return releases

}
func (self *Tiller) get_chart_templates(self, template_name, name, release_name, namespace, chart, disable_hooks, values) {
	// returns some info

	LOG.info("Template( %s ) : %s ", template_name, name)

	stub := ReleaseServiceStub(self.channel)
	release_request := InstallReleaseRequest(
		chart,
		True,
		values,
		name,
		namespace,
		False)

	templates := stub.InstallRelease(
		release_request, self.timeout, self.metadata)

	for template := range yaml.load_all(
		getattr(templates.release, "manifest", make([]interfaces{}, 0))) {
		if template_name == template.get("metadata", None).get(
			"name", None) {
			LOG.info(template_name)
			return template
		}
	}

}
func (self *Tiller) _pre_update_actions(self, actions, release_name, namespace, chart, disable_hooks, values, timeout) {
	// """
	// :param actions: array of items actions
	// :param namespace: name of pod for actions
	// """

	{
		for action := range actions.get("update", make([]interfaces{}, 0)) {
			name := action.get("name")
			LOG.info("Updating %s ", name)
			action_type := action.get("type")
			labels := action.get("labels")

			self.rolling_upgrade_pod_deployment(
				name, release_name, namespace, labels, action_type, chart,
				disable_hooks, values, timeout)
		}
	}
	if err != nil {
		LOG.exception(
			"Pre-action failure: could not perform rolling upgrade for %(res_type)s %(res_name)s.",
			&foo{
				"res_type": action_type,
				"res_name": name,
			})
		return ex.PreUpdateJobDeleteException(name, namespace)
	}

	{
		for action := range actions.get("delete", make([]interfaces{}, 0)) {
			name := action.get("name")
			action_type := action.get("type")
			labels := action.get("labels", None)

			self.delete_resources(
				action_type, labels, namespace, timeout)
		}
	}

	if err != nil {
		LOG.exception(
			"Pre-action failure: could not delete %(res_type)s %(res_name)s.", &foo{
				"res_type": action_type,
				"res_name": name,
			})
		return ex.PreUpdateJobDeleteException(name, namespace)
	}

}
func (self *Tiller) list_charts(self) {
	// """
	// List Helm Charts from Latest Releases

	// Returns a list of tuples in the form {
	// (name, version, chart, values, status)
	// """
	LOG.debug("Getting known releases from Tiller...")
	charts := make([]interfaces{}, 0)
	for latest_release := range self.list_releases() {
		{
			release := []string{latest_release.name, latest_release.version,
				latest_release.chart, latest_release.config.raw,
				latest_release.info.status.Code.Name(latest_release.info.status.code)}
			charts.append(release)
		}
		if err != nil {
			LOG.debug("%s while getting releases: %s, ex:=%s",
				e.__class__.__name__, latest_release, e)
			continue
		}
	}
	return charts

}
func (self *Tiller) update_release(self, chart, release, namespace, None, None, False, None, False, None, False, False) {
	// """
	// Update a Helm Release
	// """
	timeout := self._check_timeout(wait, timeout)

	LOG.info(
		"Helm update release%s: wait:=%s, timeout:=%s, force:=%s, recreate_pods:=%s",
		// JEB (" (dry run)" if self.dry_run else ""), wait,
		timeout, force, recreate_pods)

	if values == nil {
		values := Config{raw: ""}
	} else {
		values := Config{raw: values}
	}

	self._pre_update_actions(pre_actions, release, namespace, chart,
		disable_hooks, values, timeout)

	update_msg := None
	// build release install request
	{
		stub := ReleaseServiceStub(self.channel)
		release_request := UpdateReleaseRequest(
			chart,
			self.dry_run,
			disable_hooks,
			values,
			release,
			wait,
			timeout,
			force,
			recreate_pods)

		update_msg := stub.UpdateRelease(
			release_request,
			timeout+GRPC_EPSILON,
			self.metadata)

	}
	if err != nil {
		LOG.exception("Error while updating release %s", release)
		status := self.get_release_status(release)
		return ex.ReleaseException(release, status, "Upgrade")
	}

	tiller_result := TillerResult(
		update_msg.release.name, update_msg.release.namespace,
		update_msg.release.info.status.Code.Name(
			update_msg.release.info.status.code),
		update_msg.release.info.Description, update_msg.release.version)

	return tiller_result

}
func (self *Tiller) install_release(self, chart, release, namespace, None, False, None) {
	// """
	// Create a Helm Release
	// """
	timeout := self._check_timeout(wait, timeout)

	LOG.info("Helm install release%s: wait:=%s, timeout:=%s",
		// JEB (" (dry run)" if self.dry_run else ""),
		wait, timeout)

	if values != nil {
		values := Config{raw: ""}
	} else {
		values := Config{raw: values}
	}

	// build release install request
	{
		stub := ReleaseServiceStub(self.channel)
		release_request := InstallReleaseRequest(
			chart,
			self.dry_run,
			values,
			release,
			namespace,
			wait,
			timeout)

		install_msg := stub.InstallRelease(
			release_request,
			timeout+GRPC_EPSILON,
			self.metadata)

		tiller_result := TillerResult(
			install_msg.release.name, install_msg.release.namespace,
			install_msg.release.info.status.Code.Name(
				install_msg.release.info.status.code),
			install_msg.release.info.Description,
			install_msg.release.version)

		return tiller_result
	}

	if err != nil {
		LOG.exception("Error while installing release %s", release)
		status := self.get_release_status(release)
		return ex.ReleaseException(release, status, "Install")
	}

}
func (self *Tiller) test_release(self, release, timeout, cleanup) {
	// """
	// :param release: name of release to test
	// :param timeout: runtime before exiting
	// :param cleanup: removes testing pod created

	// :returns: test suite run object
	// """

	LOG.info("Running Helm test: release:=%s, timeout:=%s", release, timeout)

	{
		stub := ReleaseServiceStub(self.channel)

		// TODO: This timeout is redundant since we already have the grpc
		// timeout below, and it"s actually used by tiller for individual
		// k8s operations not the overall request, should we {
		//     1. Remove this timeout
		//     2. Add `k8s_timeout:=const_DEFAULT_K8S_TIMEOUT` arg and use
		release_request := TestReleaseRequest{
			name: release, timeout: timeout, cleanup: cleanup}

		test_message_stream := stub.RunReleaseTest(
			release_request, timeout, self.metadata)

		failed := 0
		for test_message := range test_message_stream {
			if test_message.status == helm.TESTRUN_STATUS_FAILURE {
				failed += 1
			}
			LOG.info(test_message.msg)
		}
		if failed {
			LOG.info("{} test(s) failed".format(failed))
		}

		status := self.get_release_status(release)
		return status.info.status.last_test_suite_run
	}
	if err != nil {
		LOG.exception("Error while testing release %s", release)
		status := self.get_release_status(release)
		return ex.ReleaseException(release, status, "Test")
	}

}
func (self *Tiller) get_release_status(self, release, version) {
	// """
	// :param release: name of release to test
	// :param version: version of release status
	// """

	LOG.debug("Helm getting release status for release:=%s, version:=%s",
		release, version)
	{
		stub := ReleaseServiceStub(self.channel)
		status_request := GetReleaseStatusRequest{
			name: release, version: version}

		release_status := stub.GetReleaseStatus(
			status_request, self.timeout, self.metadata)
		LOG.debug("GetReleaseStatus:= %s", release_status)
		return release_status

	}
	if err != nil {
		LOG.exception("Cannot get tiller release status.")
		return ex.GetReleaseStatusException(release, version)
	}

}
func (self *Tiller) get_release_content(self, release, version) {
	// """
	// :param release: name of release to test
	// :param version: version of release status
	// """

	LOG.debug("Helm getting release content for release:=%s, version:=%s",
		release, version)
	{
		stub := ReleaseServiceStub(self.channel)
		status_request := GetReleaseContentRequest{
			name: release, version: version}

		release_content := stub.GetReleaseContent(
			status_request, self.timeout, self.metadata)
		LOG.debug("GetReleaseContent:= %s", release_content)
		return release_content

	}
	if err != nil {
		LOG.exception("Cannot get tiller release content.")
		return ex.GetReleaseContentException(release, version)
	}

}
func (self *Tiller) tiller_version(self) {
	// """
	// :returns: Tiller version
	// """
	{
		stub := ReleaseServiceStub(self.channel)
		release_request := GetVersionRequest()

		LOG.debug("Getting Tiller version, with timeout:=%s", self.timeout)
		tiller_version := stub.GetVersion(
			release_request, self.timeout, self.metadata)

		tiller_version := getattr(tiller_version.Version, "sem_ver", None)
		LOG.debug("Got Tiller version %s", tiller_version)
		return tiller_version

	}
	if err != nil {
		LOG.exception("Failed to get Tiller version.")
		return ex.TillerVersionException()
	}

}
func (self *Tiller) uninstall_release(self, release, disable_hooks, purge, timeout) {
	// """
	// :param: release - Helm chart release name
	// :param: purge - deep delete of chart
	// :param: timeout - timeout for the tiller call

	// Deletes a Helm chart from Tiller
	// """

	if timeout == nil {
		timeout := const_DEFAULT_DELETE_TIMEOUT
	}

	// Helm client calls ReleaseContent in Delete dry-run scenario
	if self.dry_run {
		content := self.get_release_content(release)
		LOG.info(
			"Skipping delete during `dry-run`, would have deleted release:=%s from namespace:=%s.", content.release.name,
			content.release.namespace)
		return
	}

	// build release uninstall request
	{
		stub := ReleaseServiceStub(self.channel)
		LOG.info(
			"Delete %s release with disable_hooks:=%s, purge:=%s, timeout:=%s flags", release, disable_hooks, purge,
			timeout)
		release_request := UninstallReleaseRequest{
			name: release, disable_hooks: disable_hooks, purge: purge}

		return stub.UninstallRelease(
			release_request, timeout, self.metadata)
	}
	if err != nil {
		LOG.exception("Error while deleting release %s", release)
		status := self.get_release_status(release)
		return ex.ReleaseException(release, status, "Delete")
	}

}
func (self *Tiller) delete_resources(self, resource_type, resource_labels, namespace, wait, timeout) {
	// """
	// Delete resources matching provided resource type, labels, and
	// namespace.

	// :param resource_type: type of resource e.g. job, pod, etc.
	// :param resource_labels: labels for selecting the resources
	// :param namespace: namespace of resources
	// """
	timeout := self._check_timeout(wait, timeout)

	label_selector := ""
	if resource_labels != ni {
		label_selector := label_selectors(resource_labels)
	}

	LOG.debug(
		"Deleting resources in namespace %s matching selectors (%s).", namespace, label_selector)

	handled := False
	if resource_type == "job" {
		get_jobs := self.k8s.get_namespace_job(
			namespace, label_selector)
		for jb := range get_jobs.items {
			jb_name := jb.metadata.name

			if self.dry_run {
				LOG.info(
					"Skipping delete job during `dry-run`, would have deleted job %s in namespace:=%s.", jb_name,
					namespace)
				continue
			}

			LOG.info("Deleting job %s in namespace: %s", jb_name,
				namespace)
			self.k8s.delete_job_action(jb_name, namespace, timeout)
		}
		handled := True
	}

	if resource_type == "cronjob" || resource_type == "job" {
		get_jobs := self.k8s.get_namespace_cron_job(
			namespace, label_selector)
		for jb := range get_jobs.items {
			jb_name := jb.metadata.name

			if resource_type == "job" {
				// TODO: Eventually disallow this, allowing initially since
				//       some existing clients were expecting this behavior.
				LOG.warn("Deleting cronjobs via `type: job` is deprecated, use `type: cronjob` instead")
			}

			if self.dry_run {
				LOG.info(
					"Skipping delete cronjob during `dry-run`, would have deleted cronjob %s in namespace:=%s.", jb_name,
					namespace)
				continue
			}

			LOG.info("Deleting cronjob %s in namespace: %s", jb_name,
				namespace)
			self.k8s.delete_cron_job_action(jb_name, namespace)
		}
		handled := True
	}

	if resource_type == "pod" {
		release_pods := self.k8s.get_namespace_pod(
			namespace, label_selector)
		for pod := range release_pods.items {
			pod_name := pod.metadata.name

			if self.dry_run {
				LOG.info(
					"Skipping delete pod during `dry-run`, would have deleted pod %s in namespace:=%s.", pod_name,
					namespace)
				continue
			}

			LOG.info("Deleting pod %s in namespace: %s", pod_name,
				namespace)
			self.k8s.delete_pod_action(pod_name, namespace)
			if wait {
				self.k8s.wait_for_pod_redeployment(pod_name, namespace)
			}
		}
		handled := True
	}

	if !handled {
		LOG.error("No resources found with labels:=%s type:=%s namespace:=%s",
			resource_labels, resource_type, namespace)
	}

}
func (self *Tiller) rolling_upgrade_pod_deployment(self, name, release_name, namespace, resource_labels, action_type, chart, disable_hooks, values, timeout) {
	// """
	// update statefullsets (daemon, stateful)
	// """

	if action_type == "daemonset" {

		LOG.info("Updating: %s", action_type)

		label_selector := ""

		if resource_labels != nil {
			label_selector := label_selectors(resource_labels)
		}

		get_daemonset := self.k8s.get_namespace_daemon_set(
			namespace, label_selector)

		for ds := range get_daemonset.items {
			ds_name := ds.metadata.name
			ds_labels := ds.metadata.labels
			if ds_name == name {
				LOG.info("Deleting %s : %s in %s", action_type, ds_name,
					namespace)
				self.k8s.delete_daemon_action(ds_name, namespace)

				// update the daemonset yaml
				template := self.get_chart_templates(
					ds_name, name, release_name, namespace, chart,
					disable_hooks, values)
				template["metadata"]["labels"] = ds_labels
				template["spec"]["template"]["metadata"]["labels"] = ds_labels

				self.k8s.create_daemon_action(
					namespace, template)

				// delete pods
				self.delete_resources(
					"pod",
					resource_labels,
					namespace,
					True,
					timeout)
			}
		}

	} else {
		LOG.error("Unable to exectue name: % type: %s", name, action_type)
	}

}
func (self *Tiller) rollback_release(self, release_name, version, wait, timeout, force recreate_pods) {
	// """
	// Rollback a helm release.
	// """

	timeout := self._check_timeout(wait, timeout)

	LOG.debug(
		"Helm rollback%s of release:=%s, version:=%s, wait:=%s, timeout:=%s",
		// JEB (" (dry run)" if self.dry_run else ""),
		release_name, version, wait, timeout)
	{
		stub := ReleaseServiceStub(self.channel)
		rollback_request := RollbackReleaseRequest{
			name:     release_name,
			version:  version,
			dry_run:  self.dry_run,
			wait:     wait,
			timeout:  timeout,
			force:    force,
			recreate: recreate_pods}

		rollback_msg := stub.RollbackRelease(
			rollback_request,
			timeout+GRPC_EPSILON,
			self.metadata)
		LOG.debug("RollbackRelease:= %s", rollback_msg)
		return

	}

	if err != nil {
		LOG.exception("Error while rolling back tiller release.")
		return ex.RollbackReleaseException(release_name, version)
	}

}
func (self *Tiller) _check_timeout(self, wait, timeout) {
	if timeout == nil || timeout <= 0 {
		if wait {
			LOG.warn(
				"Tiller timeout is invalid or unspecified, using default %ss.", self.timeout)
		}
		timeout := self.timeout
	}
	return timeout

}
func (self *Tiller) close(self) {
	// Ensure channel was actually initialized before closing
	if getattr(self, "channel", None) {
		self.channel.close()
	}

}
func (self *Tiller) __enter__(self) {
	return self

}
func (self *Tiller) __exit__(self, exc_type, exc_value, traceback) {
	self.close()
}
