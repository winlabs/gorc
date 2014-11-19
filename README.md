gorc
====

This utility is a Win32 resource compiler written in Go.  It currently only supports version stamp resources,
although support for other resource types may be added in the future.  It may be used to add resources to a Win32
executable compiled from Go.  It is run on the executable after it is built and modifies it to add the resources.

### Usage

The utility should be invoked with two paths on the command line, the first specifying a JSON file containing the
resources to be compiled and the second specifying the executable file to which they should be added.

	gorc hello_resources.json hello.exe

There is also a flag `--discard` that can be used to instruct the utility to delete any resources that already exist in
the executable.

	gorc --discard hello_resources.json hello.exe

### JSON File Format

The following is an example JSON file showing the format used to specify the resources.  All version information fields
are optional and may be omitted if not needed.

	{
		"language": "en-us",
		"version": {
			"fileVersion": "1.0.0.0",
			"productVersion": "1.0.0.0",
			"fileFlags": ["VS_FF_DEBUG", "VS_FF_PRERELEASE"],
			"fileOS": "VOS_NT_WINDOWS32",
			"fileType": "VFT_APP",
			"fileSubtype": "VFT2_UNKNOWN",
			"stringFileInfo": {
				"comments": "Example executable to demonstrate version stamping with gorc.",
				"companyName": "MongoDB",
				"fileDescription": "Hello",
				"fileVersion": "1.0",
				"internalName": "hello",
				"legalCopyright": "Copyright (C) 2014 MongoDB, Inc.",
				"legalTrademarks": "MongoDB is a registered trademark of MongoDB, Inc.",
				"originalFilename": "hello.exe",
				"productName": "Hello",
				"productVersion": "1.0"
			}
		}
	}
