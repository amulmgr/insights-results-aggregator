# Copyright © 2020 Pavel Tisnovsky
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

"""Simple checker of all JSONs in the given directory (usually repository)."""

from pathlib import Path
from json import load
from sys import exit

passes = 0
failures = 0

files = list(Path(".").rglob("*.json"))

for file in files:
    try:
        with file.open() as fin:
            obj = load(fin)
            print("{} is valid".format(file))
            passes += 1
    except ValueError as e:
        print("{} is invalid".format(file))
        failures += 1
        print(e)

print("{} passes".format(passes))
print("{} failures".format(failures))

if failures > 0:
    exit(1)
