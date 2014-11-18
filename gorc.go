/*
 * Copyright (c) 2014 MongoDB, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the license is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/winlabs/gowin32"

	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gorc file.json file.exe\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	jsonFile, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open JSON file: %s (%s)\n", args[0], err)
		os.Exit(2)
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)

	var jsonData map[string]interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse JSON file: %s (%s)\n", args[0], err)
		os.Exit(2)
	}
	language, resources, err := ParseResources(jsonData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid resources in JSON file: %s (%s)\n", args[0], err)
		os.Exit(2)
	}

	update, err := gowin32.NewResourceUpdate(args[1], false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open executable file for resource update: %s (%s)\n", args[1], err)
		os.Exit(2)
	}
	defer update.Close()

	for _, res := range resources {
		if err := update.Update(res.Type, gowin32.IntResourceId(res.Id), language, res.Data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to update resource (%s)\n", err)
			os.Exit(2)
		}
	}

	if err := update.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save updated resources to executable file: %s (%s)\n", args[1], err)
		os.Exit(2)
	}
}
