/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package env

import "testing"

func TestVersion(t *testing.T) {
	if Version() != "v"+MainVersion+".x" {
		t.Fatal("TestVersion failed:", Version())
	}
	buildVersion = "v1.0.0-beta1"
	if Version() != buildVersion {
		t.Fatal("TestVersion failed:", Version())
	}
	buildVersion = ""
}

func TestBuild(t *testing.T) {
	if BuildRevision() != "" {
		t.Fatal("BuildCommit failed:", BuildRevision())
	}
	if BuildDate() != "" {
		t.Fatal("BuildInfo failed:", BuildDate())
	}
}
