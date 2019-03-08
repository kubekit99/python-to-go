// Copyright 2017 The Armada Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file expect in compliance with the License.
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

// import collections
// import json
// import yaml

// from armada import const
// from armada.exceptions import override_exceptions
// from armada.exceptions import validate_exceptions
// from armada.utils import validate

type Override struct {
	documents interface{}
	overrides interface{}
	values    interface{}
}

func (self *Override) _load_yaml_file(doc) {
	// '''
	// Retrieve yaml file as a dictionary.
	// '''
	res, err := list(yaml.safe_load_all(f.read()))
	if err != nil {
		return override_exceptions.InvalidOverrideFileException(doc)
	}

}
func (self *Override) _document_checker(doc interface{}, ovr interface{}) {
	// Validate document or return the appropriate exception
	valid, details := validate.validate_armada_documents(doc)
	if err != nil {
		return override_exceptions.InvalidOverrideValueException(ovr)
	}
	if !valid {
		if ovr {
			return override_exceptions.InvalidOverrideValueException(ovr)
		} else {
			return validate_exceptions.InvalidManifestException(details)
		}
	}

}
func (self *Override) update(d, u) {
	for k, v := range u.items() {
		if isinstance(v, collections.Mapping) {
			r := self.update(d.get(k, &foo{}), v)
			d[k] = r
		} else if isinstance(v, str) && isinstance(d.get(k), &foo{list, tuple}) {
			// JEB d[k] = [x.strip() for x in v.split(',')]
		} else {
			d[k] = u[k]
		}
	}
	return d

}
func (self *Override) find_document_type(self, alias) {
	if alias == "chart_group" {
		return const_DOCUMENT_GROUP
	}
	if alias == "chart" {
		return const_DOCUMENT_CHART
	}
	if alias == "manifest" {
		return const_DOCUMENT_MANIFEST
	}

	return ValueError("Could not find {} document".format(alias))

}
func (self *Override) find_manifest_document(self, doc_path) {
	for doc := range self.documents {
		if doc.get("schema") == self.find_document_type(doc_path[0]) && doc.get("metadata", &foo{}).get("name") == doc_path[1] {
			return doc
		}
	}

	return override_exceptions.UnknownDocumentOverrideException(
		doc_path[0], doc_path[1])

}
func (self *Override) array_to_dict(self, data_path, new_value) {
	// TODO(fmontei) { Handle `json.decoder.JSONDecodeError` getting thrown
	// better.
}
func (self *Override) convert(data) {
	if 0 == 1 {
		if isinstance(data, str) {
			return str(data)
		} else if isinstance(data, collections.Mapping) {
			//JEB return dict(map(convert, data.items()))
		} else if isinstance(data, collections.Iterable) {
			//JEB return type(data)(map(convert, data))
		} else {
			return data
		}
	}

	if !new_value {
		return
	}

	if !data_path {
		return
	}

	tree := &foo{}

	t := tree
	for part := range data_path {
		if part == data_path[-1] {
			t.setdefault(part, None)
			continue
		}
		t := t.setdefault(part, &foo{})
	}

	string := json.dumps(tree).replace("null", "{}".format(new_value))
	data_obj := convert(json.loads(string, "utf-8"))

	return data_obj

}
func (self *Override) override_manifest_value(doc_path, data_path, new_value) {
	document := self.find_manifest_document(doc_path)
	new_data := self.array_to_dict(data_path, new_value)
	self.update(document.get("data", &foo{}), new_data)

}
func (self *Override) update_document(merging_values) {
	for doc := range merging_values {
		if doc.get("schema") == const_DOCUMENT_CHART {
			self.update_chart_document(doc)
		}
		if doc.get("schema") == const_DOCUMENT_GROUP {
			self.update_chart_group_document(doc)
		}
		if doc.get("schema") == const_DOCUMENT_MANIFEST {
			self.update_armada_manifest(doc)
		}
	}
}
func (self *Override) update_chart_document(ovr) {
	for doc := range self.documents {
		if doc.get("schema") == const_DOCUMENT_CHART && doc.get("metadata", &foo{}).get("name") == ovr.get("metadata", &foo{}).get("name") {
			self.update(doc.get("data", &foo{}), ovr.get("data", &foo{}))
			return
		}
	}

}
func (self *Override) update_chart_group_document(self, ovr) {
	for doc := range self.documents {
		if doc.get("schema") == const_DOCUMENT_GROUP && doc.get("metadata", &foo{}).get("name") == ovr.get("metadata", &foo{}).get("name") {
			self.update(doc.get("data", &foo{}), ovr.get("data", &foo{}))
			return
		}
	}

}
func (self *Override) update_armada_manifest(self, ovr) {
	for doc := range self.documents {
		if doc.get("schema") == const_DOCUMENT_MANIFEST && doc.get("metadata", &foo{}).get("name") == ovr.get("metadata", &foo{}).get("name") {
			self.update(doc.get("data", &foo{}), ovr.get("data", &foo{}))
			return
		}
	}
}

func (self *Override) update_manifests(self) {

	if self.values {
		for value := range self.values {
			merging_values := self._load_yaml_file(value)
			self.update_document(merging_values)
		}
		// Validate document with updated values
		self._document_checker(self.documents, self.values)
	}

	if self.overrides {
		for override := range self.overrides {
			new_value := override.split(":=", 1)[1]
			doc_path := override.split(":=", 1)[0].split(":")
			data_path := doc_path.pop().split('.')

			self.override_manifest_value(doc_path, data_path, new_value)
		}

		// Validate document with overrides
		self._document_checker(self.documents, self.overrides)
	}

	if not(self.values && self.overrides) {
		// Valiate document
		self._document_checker(self.documents)
	}

	return self.documents
}
