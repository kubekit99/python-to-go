#!/usr/bin/python3

# Copyright 2019 Kubekit team
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import yaml
import os
import argparse
import re
import shutil
import tempfile


def changesignature(signature):
    newsignature = signature
    if "()" in newsignature:
        return newsignature
    newsignature = newsignature.replace(","," interface{},")
    newsignature = newsignature.replace(") {"," interface{}) {")
    return newsignature

def changefunc(filename):
    pattern = re.compile("^func \(")

    with tempfile.NamedTemporaryFile(mode='w', delete=False) as tmp_file:
        with open(filename) as src_file:
            for line in src_file:
                if pattern.search(line):
                    # tmp_file.write(pattern_compiled.sub(repl, line))
                    tmp_file.write(changesignature(line))
                else:
                    tmp_file.write(line)

    shutil.copystat(filename, tmp_file.name)
    shutil.move(tmp_file.name, filename + ".new")

if __name__ == "__main__":
    if os.path.isfile(".settings.yaml"):
        with open('.settings.yaml') as stream2:
            localsettings = yaml.safe_load(stream2)
        defaultuser = localsettings["user"]
        defaultpassword = localsettings["password"]
    else:
        defaultuser = "username"
        defaultpassword = "xxxpasswordxxx"

    parser = argparse.ArgumentParser()
    parser.add_argument('-c', '--command',
                        help="Command to run",
                        type=str, choices=["changefunc"],
                        default="changefunc")
    parser.add_argument('-f', '--filename',
                        help="Go file",
                        type=str,
                        default="armada.go")

    args = parser.parse_args()

    if (args.command == "changefunc"):
        changefunc(args.filename)
    else:
        pass
