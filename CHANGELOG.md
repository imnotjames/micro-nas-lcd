# Changelog

## [0.3.2](https://github.com/imnotjames/micro-nas-lcd/compare/v0.3.1...v0.3.2) (2025-08-25)


### Bug Fixes

* drop memory percent that doesn't show anyway ([234390b](https://github.com/imnotjames/micro-nas-lcd/commit/234390b53e8260cd7935c5039b9a6b8ab8fb069d))
* handle spaces in columns more gracefully ([0a5b1da](https://github.com/imnotjames/micro-nas-lcd/commit/0a5b1da0f9f6816783ea000709d7b3f07918095f))

## [0.3.1](https://github.com/imnotjames/micro-nas-lcd/compare/v0.3.0...v0.3.1) (2025-08-25)


### Bug Fixes

* build with CGO_ENABLED off ([91d6867](https://github.com/imnotjames/micro-nas-lcd/commit/91d68672ba060145a55bbcc8df5190df9ace4391))

## [0.3.0](https://github.com/imnotjames/micro-nas-lcd/compare/v0.2.1...v0.3.0) (2025-08-25)


### Features

* format memory as gigabytes only ([87930ba](https://github.com/imnotjames/micro-nas-lcd/commit/87930ba7c9dfd42492c033262619bcfa6b0f0382))


### Bug Fixes

* properly pull local IP ([c8849d5](https://github.com/imnotjames/micro-nas-lcd/commit/c8849d5a3e7fea2134539d70ca3fb9436d8eb752))

## [0.2.0](https://github.com/imnotjames/micro-nas-lcd/compare/v0.1.0...v0.2.0) (2025-08-25)


### Features

* allow specifying duration between "pages" ([67ec6f4](https://github.com/imnotjames/micro-nas-lcd/commit/67ec6f47767552bd94965c1ec8e39533d4dcef10))
* show disk information ([51e5701](https://github.com/imnotjames/micro-nas-lcd/commit/51e570199539a544a6dfa59787b3d302e4612406))


### Bug Fixes

* use columns for the key/val length ([69e4d13](https://github.com/imnotjames/micro-nas-lcd/commit/69e4d13b8cb4400b1a00526c263538b4697aad66))

## [0.1.0](https://github.com/imnotjames/micro-nas-lcd/compare/v0.0.1...v0.1.0) (2025-08-24)


### Features

* allow specifying the address / column size / row size ([3e92e69](https://github.com/imnotjames/micro-nas-lcd/commit/3e92e6974f3c68eabe7e51cc65b1b04d85d72fc6))


### Bug Fixes

* remove extra new line when crafting text ([f66abde](https://github.com/imnotjames/micro-nas-lcd/commit/f66abde366e4d9abcec6c09fb78b6782156f11c7))
* track bus and close both bus and mcp ([1f91162](https://github.com/imnotjames/micro-nas-lcd/commit/1f9116243dfabe92f0617ce560ada3444a76415b))
* trim before trying to truncate or match to col width ([08105b0](https://github.com/imnotjames/micro-nas-lcd/commit/08105b01edbd04208d43a35a84c7c4cbcb404f40))
* truncate text instead of error ([78f2783](https://github.com/imnotjames/micro-nas-lcd/commit/78f278341c0adbce397c7541a823190cab6112fb))
