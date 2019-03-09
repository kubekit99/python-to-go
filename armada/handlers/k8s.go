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

// import re

// from kubernetes import client
// from kubernetes import config
// from kubernetes import watch
// from kubernetes.client import api_client
// from kubernetes.client.rest import ApiException
// from unittest.mock import Mock
// from oslo_config import cfg
// from oslo_log import log as logging

// from armada.const import DEFAULT_K8S_TIMEOUT
// from armada.exceptions import k8s_exceptions as exceptions

// CONF := cfg.CONF
// LOG := logging.getLogger(__name__)

// TODO: Remove after this bug is fixed and we have uplifted to a fixed version {
//       https://github.com/kubernetes-client/python/issues/411
// Avoid creating thread pools in kubernetes api_client.

// var _dummy_pool = Mock()
// var ThreadPool = nil // JEB lambda *args, **kwargs: _dummy_pool

type K8s struct {
	// """
	// Object to obtain the local kube config file
	// """
}

func (self *K8s) init() {
	// """
	// Initialize connection to Kubernetes
	// """
	self.bearer_token = bearer_token
	api_client := None

	config.load_incluster_config()
	if err != nil {
		config.load_kube_config()
	}

	if self.bearer_token {
		// Configure API key authorization: Bearer Token
		configuration := client.Configuration()
		configuration.api_key_prefix["authorization"] = "Bearer"
		configuration.api_key["authorization"] = self.bearer_token
		api_client := client.ApiClient(configuration)
	}

	self.client = client.CoreV1Api(api_client)
	self.batch_api = client.BatchV1Api(api_client)
	self.batch_v1beta1_api = client.BatchV1beta1Api(api_client)
	self.extension_api = client.ExtensionsV1beta1Api(api_client)
	self.apps_v1_api = client.AppsV1Api(api_client)

}
func (self *K8s) delete_job_action(name interface{}, namespace interface{}, propagation_policy interface{}, timeout interface{}) {
	// """
	// Delete a job from a namespace (see _delete_item_action).

	// :param name: name of job
	// :param namespace: namespace
	// :param propagation_policy: The Kubernetes propagation_policy to apply
	//     to the delete.
	// :param timeout: The timeout to wait for the delete to complete
	// """
	self._delete_item_action(self.batch_api.list_namespaced_job,
		self.batch_api.delete_namespaced_job, "job",
		name, namespace, propagation_policy, timeout)

}
func (self *K8s) delete_cron_job_action(name interface{}, namespace interface{}, propagation_policy interface{}, timeout interface{}) {
	// """
	// Delete a cron job from a namespace (see _delete_item_action).

	// :param name: name of cron job
	// :param namespace: namespace
	// :param propagation_policy: The Kubernetes propagation_policy to apply
	//     to the delete.
	// :param timeout: The timeout to wait for the delete to complete
	// """
	self._delete_item_action(
		self.batch_v1beta1_api.list_namespaced_cron_job,
		self.batch_v1beta1_api.delete_namespaced_cron_job, "cron job",
		name, namespace, propagation_policy, timeout)

}
func (self *K8s) delete_pod_action(name interface{}, namespace interface{}, propagation_policy interface{}, timeout interface{}) {
	// """
	// Delete a pod from a namespace (see _delete_item_action).

	// :param name: name of pod
	// :param namespace: namespace
	// :param propagation_policy: The Kubernetes propagation_policy to apply
	//     to the delete.
	// :param timeout: The timeout to wait for the delete to complete
	// """
	self._delete_item_action(self.client.list_namespaced_pod,
		self.client.delete_namespaced_pod, "pod",
		name, namespace, propagation_policy, timeout)

}
func (self *K8s) _delete_item_action(list_func interface{}, delete_func interface{}, object_type_description interface{}, name interface{}, namespace interface{}, propagation_policy interface{}, timeout interface{}) {
	// """
	// This function takes the action to delete an object (job, cronjob, pod)
	// from kubernetes. It will wait for the object to be fully deleted before
	// returning to processing or timing out.

	// :param list_func: The callback function to list the specified object
	//     type
	// :param delete_func: The callback function to delete the specified
	//     object type
	// :param object_type_description: The types of objects to delete,
	//     in `job`, `cronjob`, or `pod`
	// :param name: The name of the object to delete
	// :param namespace: The namespace of the object
	// :param propagation_policy: The Kubernetes propagation_policy to apply
	//     to the delete. Default "Foreground" means that child objects
	//     will be deleted before the given object is marked as deleted.
	//     See: https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/#controlling-how-the-garbage-collector-deletes-dependents  // noqa
	// :param timeout: The timeout to wait for the delete to complete
	// """
	{
		timeout := self._check_timeout(timeout)

		LOG.debug("Watching to delete %s %s, Wait timeout:=%s",
			object_type_description, name, timeout)
		body := client.V1DeleteOptions()
		w := watch.Watch()
		issue_delete := True
		found_events := False
		for event := range w.stream(list_func, namespace, timeout) {
			if issue_delete {
				delete_func(
					name,
					namespace,
					body,
					propagation_policy)
				issue_delete := False
			}

			event_type := event["type"].upper()
			item_name := event["object"].metadata.name
			LOG.debug("Watch event %s on %s", event_type, item_name)

			if item_name == name {
				found_events := True
				if event_type == "DELETED" {
					LOG.info("Successfully deleted %s %s",
						object_type_description, item_name)
					return
				}
			}
		}

		if !found_events {
			LOG.warn("Saw no delete events for %s %s in namespace:=%s",
				object_type_description, name, namespace)
		}

		err_msg := "Reached timeout while waiting to delete %s: name:=%s, namespace:=%s % (object_type_description, name, namespace)"
		LOG.error(err_msg)
		return exceptions.KubernetesWatchTimeoutException(err_msg)
	}

	if err != nil {
		LOG.exception("Exception when deleting %s: name:=%s, namespace:=%s", object_type_description, name, namespace)
		return e
	}
}
func (self *K8s) get_namespace_job(namespace interface{}, kwargs interface{}) {
	// """
	// :param label_selector: labels of the jobs
	// :param namespace: namespace of the jobs
	// """

	res, err := self.batch_api.list_namespaced_job(namespace, **kwargs)
	if err != nil {
		LOG.error("Exception getting jobs: namespace:=%s, label:=%s: %s",
			namespace, kwargs.get("label_selector", ""), e)
	}
	return res

}
func (self *K8s) get_namespace_cron_job(namespace interface{}, kwargs interface{}) {
	// """
	// :param label_selector: labels of the cron jobs
	// :param namespace: namespace of the cron jobs
	// """

	res, err := self.batch_v1beta1_api.list_namespaced_cron_job(
		namespace, **kwargs)
	if err != nil {
		LOG.error(
			"Exception getting cron jobs: namespace:=%s, label:=%s: %s",
			namespace, kwargs.get("label_selector", ""), e)
	}
	return res

}
func (self *K8s) get_namespace_pod(namespace interface{}, kwargs interface{}) {
	// """
	// :param namespace: namespace of the Pod
	// :param label_selector: filters Pods by label

	// This will return a list of objects req namespace
	// """

	return self.client.list_namespaced_pod(namespace, **kwargs)

}
func (self *K8s) get_namespace_deployment(namespace interface{}, kwargs interface{}) {
	// """
	// :param namespace: namespace of target deamonset
	// :param labels: specify targeted deployment
	// """
	return self.apps_v1_api.list_namespaced_deployment(namespace, **kwargs)

}
func (self *K8s) get_namespace_stateful_set(namespace interface{}, kwargs interface{}) {
	// """
	// :param namespace: namespace of target stateful set
	// :param labels: specify targeted stateful set
	// """
	return self.apps_v1_api.list_namespaced_stateful_set(
		namespace, **kwargs)

}
func (self *K8s) get_namespace_daemon_set(namespace interface{}, kwargs interface{}) {
	// """
	// :param namespace: namespace of target deamonset
	// :param labels: specify targeted daemonset
	// """
	return self.apps_v1_api.list_namespaced_daemon_set(namespace, **kwargs)

}
func (self *K8s) create_daemon_action(namespace interface{}, template interface{}) {
	// """
	// :param: namespace - pod namespace
	// :param: template - deploy daemonset via yaml
	// """
	// we might need to load something here

	self.apps_v1_api.create_namespaced_daemon_set(namespace, body)

}
func (self *K8s) delete_daemon_action(name interface{}, namespace interface{}, body interface{}) {
	// """
	// :param: namespace - pod namespace

	// This will delete daemonset
	// """

	if body {
		body := client.V1DeleteOptions()
	}

	return self.apps_v1_api.delete_namespaced_daemon_set(
		name, namespace, body)

}
func (self *K8s) wait_for_pod_redeployment(old_pod_name interface{}, namespace interface{}) {
	// """
	// :param old_pod_name: name of pods
	// :param namespace: kubernetes namespace
	// """

	base_pod_pattern := re.compile("^(.+)-[a-zA-Z0-9]+$")

	if !base_pod_pattern.match(old_pod_name) {
		LOG.error("Could not identify new pod after purging %s",
			old_pod_name)
		return
	}

	pod_base_name := base_pod_pattern.match(old_pod_name).group(1)

	new_pod_name := ""

	w := watch.Watch()
	for event := range w.stream(self.client.list_namespaced_pod, namespace) {
		event_name := event["object"].metadata.name
		event_match := base_pod_pattern.match(event_name)
		if !event_match || !event_match.group(1) == pod_base_name {
			continue
		}

		pod_conditions := event["object"].status.conditions
		// wait for new pod deployment
		if event["type"] == "ADDED" && !pod_conditions {
			new_pod_name := event_name
		} else if new_pod_name {
			for condition := range pod_conditions {
				if condition.typef == "Ready" && condition.status == "True" {
					LOG.info("New pod %s deployed", new_pod_name)
					w.stop()
				}
			}
		}
	}

}
func (self *K8s) wait_get_completed_podphase(release interface{}, timeout interface{}) {
	// """
	// :param release: part of namespace
	// :param timeout: time before disconnecting stream
	// """
	timeout = self._check_timeout(timeout)

	w := watch.Watch()
	found_events := False
	for event := range w.stream(self.client.list_pod_for_all_namespaces, timeout) {
		resource_name := event["object"].metadata.name

		// JEB if release in resource_name {
		if true {
			found_events := True
			pod_state := event["object"].status.phase
			if pod_state == "Succeeded" {
				w.stop()
				break
			}
		}
	}

	if !found_events {
		LOG.warn("Saw no test events for release %s", release)
	}

}
func (self *K8s) _check_timeout(timeout interface{}) {
	if timeout <= 0 {
		LOG.warn(
			"Kubernetes timeout is invalid or unspecified, using default %ss.", DEFAULT_K8S_TIMEOUT)
		timeout = DEFAULT_K8S_TIMEOUT
	}
	return timeout
}
