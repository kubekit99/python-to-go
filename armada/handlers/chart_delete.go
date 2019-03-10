// Copyright 2019 The Armada Authors.
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

// from armada import const

type ChartDelete struct {
	// """Initialize a chart delete handler.

	// :param chart: The armada chart document
	// :param release_name: Name of a Helm release
	// :param tiller: Tiller object
	// :param purge: Whether to purge the release

	// :type chart: object
	// :type release_name: str
	// :type tiller: Tiller object
	// :type purge: bool
	// """

	chart         interface{}
	release_name  interface{}
	tiller        *Tiller
	purge         interface{}
	delete_config interface{}
	// TODO(seaneagan) { Consider allowing this to be a percentage of the
	// chart's `wait.timeout` so that the timeouts can scale together, and
	// likely default to some reasonable value, e.g. "50%".
	timeout interface{}
}

// :wait
func (self *ChartDelete) get_timeout() {
	return self.timeout

}
func (self *ChartDelete) delete() {
	// """Delete the release associated with the chart"
	// """
	self.tiller.uninstall_release(release_name, self.get_timeout(), self.purge)
}
