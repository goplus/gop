/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package gocmd

// -----------------------------------------------------------------------------

type InstallConfig = Config

func Install(arg string, conf *InstallConfig) (err error) {
	return doWithArgs("", "install", conf, arg)
}

func InstallFiles(files []string, conf *InstallConfig) (err error) {
	return doWithArgs("", "install", conf, files...)
}

// -----------------------------------------------------------------------------

type BuildConfig = Config

func Build(arg string, conf *BuildConfig) (err error) {
	return doWithArgs("", "build", conf, arg)
}

func BuildFiles(files []string, conf *BuildConfig) (err error) {
	return doWithArgs("", "build", conf, files...)
}

// -----------------------------------------------------------------------------

type TestConfig = Config

func Test(arg string, conf *TestConfig) (err error) {
	return doWithArgs("", "test", conf, arg)
}

func TestFiles(files []string, conf *TestConfig) (err error) {
	return doWithArgs("", "test", conf, files...)
}

// -----------------------------------------------------------------------------
