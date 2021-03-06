// Copyright 2018 The Armada Authors.
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

// from oslo_log import log as logging

// from armada import const

// from armada.handlers.wait import get_wait_labels
// from armada.utils.release import label_selectors
// from armada.utils.helm import get_test_suite_run_success, is_test_pod

// LOG := logging.getLogger(__name__)

type Test struct {
	// """Initialize a test handler to run Helm tests corresponding to a
	// release.

	// :param chart: The armada chart document
	// :param release_name: Name of a Helm release
	// :param tiller: Tiller object
	// :param cg_test_charts: Chart group `test_charts` key
	// :param cleanup: Triggers cleanup; overrides `test.options.cleanup`
	// :param enable_all: Run tests regardless of the value of `test.enabled`

	// :type chart: dict
	// :type release_name: str
	// :type tiller: Tiller object
	// :type cg_test_charts: bool
	// :type cleanup: bool
	// :type enable_all: bool
	// """

	chart        interface{}
	release_name interface{}
	tiller       *Tiller
	cleanup      interface{}
	k8s_timeout  int
	timeout      int
}

func (self *Test) init() {
	test_values := self.chart.get("test", None)
	// NOTE(drewwalters96) { Support the chart_group `test_charts` key until
	// its deprecation period ends. The `test.enabled`, `enable_all` flag,
	// and deprecated, boolean `test` key override this value if provided.
	if cg_test_charts {
		LOG.warn("Chart group key `test_charts` is deprecated and will be removed. Use `test.enabled` instead.")
		self.test_enabled = cg_test_charts
	} else {
		self.test_enabled = true
	}

	// NOTE: Support old, boolean `test` key until deprecation period ends.
	if typef(test_values) == bool {
		LOG.warn("Boolean value for chart `test` key is deprecated and will be removed. Use `test.enabled` instead.")

		self.test_enabled = test_values

		// NOTE: Use old, default cleanup value (i.e. true) if none is
		// provided.
		if self.cleanup {
			self.cleanup = true
		}
	} else if test_values {
		test_enabled_opt := test_values.get("enabled")
		if !test_enabled_opt {
			self.test_enabled = test_enabled_opt
		}
		// NOTE(drewwalters96) { `self.cleanup`, the cleanup value provided
		// by the API/CLI, takes precedence over the chart value
		// `test.cleanup`.
		if self.cleanup {
			test_options := test_values.get("options", &foo{})
			self.cleanup = test_options.get("cleanup", false)
		}
		self.timeout = test_values.get("timeout", self.timeout)
	} else {
		// Default cleanup value
		if self.cleanup {
			self.cleanup = false
		}
	}

	if enable_all {
		self.test_enabled = true
	}
}

func (self *Test) test_release_for_success() {
	// """Run the Helm tests corresponding to a release for success (i.e. exit
	// code 0).

	// :return: Helm test suite run result
	// """
	LOG.Info("RUNNING: %s tests with timeout:=%ds", self.release_name,
		self.timeout)

	self.delete_test_pods()
	if err != nil {
		LOG.exception("Exception when deleting test pods for release: %s",
			self.release_name)
	}

	test_suite_run := self.tiller.test_release(
		self.release_name, self.timeout, self.cleanup)

	success := get_test_suite_run_success(test_suite_run)
	if success {
		LOG.Info("PASSED: %s", self.release_name)
	} else {
		LOG.Info("FAILED: %s", self.release_name)
	}

	return success
}

func (self *Test) delete_test_pods() {
	// """Deletes any existing test pods for the release, as identified by the
	// wait labels for the chart, to avoid test pod name conflicts when
	// creating the new test pod as well as just for general cleanup since
	// the new test pod should supercede it.
	// """
	labels := get_wait_labels(self.chart)

	// Guard against labels being left empty, so we don"t delete other
	// chart"s test pods.
	if labels {
		label_selector := label_selectors(labels)

		namespace := self.chart["namespace"]

		list_args := &foo{
			"namespace":       namespace,
			"label_selector":  label_selector,
			"timeout_seconds": self.k8s_timeout,
		}

		pod_list := self.tiller.k8s.client.list_namespaced_pod(**list_args)
		// JEB test_pods := [pod for pod in pod_list.items if is_test_pod(pod)]
		test_pods := make([]interface{}, 0)

		if test_pods {
			LOG.Info(
				"Found existing test pods for release with namespace:=%s, labels:=(%s)", namespace, label_selector)
		}

		for test_pod := range test_pods {
			pod_name := test_pod.metadata.name
			LOG.Info("Deleting existing test pod: %s", pod_name)
			self.tiller.k8s.delete_pod_action(
				pod_name, namespace, self.k8s_timeout)
		}
	}
}
