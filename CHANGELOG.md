# Changelog

## [1.2.1](https://github.com/sergiocarracedo/on-a-meet/compare/v1.2.0...v1.2.1) (2026-05-30)


### Bug Fixes

* add bottle :unneeded to Homebrew formula ([f94031a](https://github.com/sergiocarracedo/on-a-meet/commit/f94031a2340813c426c18a9c160bfd58e398bfd8))

## [1.2.0](https://github.com/sergiocarracedo/on-a-meet/compare/v1.1.0...v1.2.0) (2026-05-29)


### Features

* **01-01:** create main.go entry point ([d60a765](https://github.com/sergiocarracedo/on-a-meet/commit/d60a765efb87bf6753e49f64b4fca859f5316b15))
* **01-01:** create root command with Viper config ([ebd33b7](https://github.com/sergiocarracedo/on-a-meet/commit/ebd33b7daaf92f4f5b34b4019169b0e6fd3dec40))
* **01-02:** add --silent/--verbose persistent flags and output wiring ([2446268](https://github.com/sergiocarracedo/on-a-meet/commit/244626856c569d167033df9cef6e516410715b4f))
* **01-02:** create config example and all subcommand stubs ([b1bedc9](https://github.com/sergiocarracedo/on-a-meet/commit/b1bedc9010b79bc4a90e813808f76a0f54a8b88b))
* **01-02:** create internal config package ([738d8d4](https://github.com/sergiocarracedo/on-a-meet/commit/738d8d4cce98d1e793f840eea6ee1d42fb5a28a3))
* **01-02:** create internal output package with pterm wrappers ([f306e53](https://github.com/sergiocarracedo/on-a-meet/commit/f306e532c7b840f2c375f3c2d884584430d6edfd))
* **02-01:** add non-Linux stub for V4L2Detector ([cba7738](https://github.com/sergiocarracedo/on-a-meet/commit/cba7738b8d1b71bc8ded7fafff9986a4225c664d))
* **02-01:** create Detector interface and types ([3e0f334](https://github.com/sergiocarracedo/on-a-meet/commit/3e0f33492ce68196dfacd06f8d374b51a75548b6))
* **02-01:** implement V4L2Detector with ListDevices and Detect ([4f453e6](https://github.com/sergiocarracedo/on-a-meet/commit/4f453e6ccb9d5bafd07c88d2257c19575f1e96c6))
* **02-02:** create polling engine with debounce and hotplug ([a259807](https://github.com/sergiocarracedo/on-a-meet/commit/a25980726323d0173f65d937a54b30a975114d52))
* **02-02:** wire detect command to polling engine ([1cd9630](https://github.com/sergiocarracedo/on-a-meet/commit/1cd9630cd52870f08125ed911cc1cbbd057ab1e6))
* **03-01:** add executor unit tests — timeout, same-state skip, cross-state ([cbf00bd](https://github.com/sergiocarracedo/on-a-meet/commit/cbf00bd7109ad4399315655b496af9f6fe4f5182))
* **03-01:** add Timeout field to Config struct ([9e06233](https://github.com/sergiocarracedo/on-a-meet/commit/9e0623309475c4bd3a8961925f6f8febe602ab60))
* **03-01:** create executor package with TemplateData and Executor struct ([314ea66](https://github.com/sergiocarracedo/on-a-meet/commit/314ea66504dccaebafb8410720bc8940962f0e2c))
* **03-02:** update config.yaml.example with timeout and template docs ([a8da70e](https://github.com/sergiocarracedo/on-a-meet/commit/a8da70e8f4aecbdcd9b2efbe9741a282fa1aae76))
* **03-02:** wire executor into detect command with --timeout flag ([397627e](https://github.com/sergiocarracedo/on-a-meet/commit/397627e3f79de21d3fa82886b61c7a2cfdee1ca8))
* **04-01:** add kardianos/service dependency ([68ef528](https://github.com/sergiocarracedo/on-a-meet/commit/68ef528213b5fa2e1b85a631b78deded63f75811))
* **04-01:** implement install command with kardianos/service ([0cba478](https://github.com/sergiocarracedo/on-a-meet/commit/0cba47876408cea83543e8b928cc9926765eeba0))
* **04-02:** implement uninstall command with kardianos/service ([61b08fc](https://github.com/sergiocarracedo/on-a-meet/commit/61b08fcf0a26987a8db3cf5b5687947e4d3d9a06))
* **05-01:** add lsof backend + factory + tests ([392753c](https://github.com/sergiocarracedo/on-a-meet/commit/392753cffd1ebd29922a7c429f70116e20a9ff49))
* **05-02:** wire factory, list command, goreleaser, README ([4be39d9](https://github.com/sergiocarracedo/on-a-meet/commit/4be39d912bd41d33f74c634575dfc0f0d7d0c1f0))
* add release-please automation with auto-merge (ci, release PR, GoReleaser) ([d21b745](https://github.com/sergiocarracedo/on-a-meet/commit/d21b74547efbbd5ed231467aa90b37076db73d22))
* add release-please automation with auto-merge, CI, and Homebrew tap ([40d5960](https://github.com/sergiocarracedo/on-a-meet/commit/40d5960f778298261be78dbf72895050eb04b647))
* add restart command to restart the service unit ([9af4f39](https://github.com/sergiocarracedo/on-a-meet/commit/9af4f39ba4c2fc7694753fcc3d10a5f9859fdb5e))
* add service subcommand group with start/stop/restart + verbose config ([e81f640](https://github.com/sergiocarracedo/on-a-meet/commit/e81f640b3e5d9f7f314d62f4f88607819202cdf5))
* add version subcommand ([b80a367](https://github.com/sergiocarracedo/on-a-meet/commit/b80a3674315ad882584fb1faea67b53ed069ef15))
* implement macOS detection backend (MacOSDetector) with log stream + system_profiler ([c3b3087](https://github.com/sergiocarracedo/on-a-meet/commit/c3b30870f8a8af6c61e88238d57b14371bd2760d))
* **phase-6:** interactive onboard wizard with huh ([4566166](https://github.com/sergiocarracedo/on-a-meet/commit/456616600852d04a909b21b9a54533dc95882ffa))
* **phase-6:** sudo apply path and auto re-exec ([a106d4e](https://github.com/sergiocarracedo/on-a-meet/commit/a106d4e962ebc10a1c1afa11bff193e1b12b32e1))
* print config on detect startup ([196ad05](https://github.com/sergiocarracedo/on-a-meet/commit/196ad051787760748f294641b4d24d974dcac429))
* **quick-005:** emit initial camera status on detect startup ([d33a9df](https://github.com/sergiocarracedo/on-a-meet/commit/d33a9df1fd5f8ca0d5daebde638c4d777e439735))
* **quick-006:** log on/off command exit code and verbose output ([8ca2245](https://github.com/sergiocarracedo/on-a-meet/commit/8ca224540bf830228c5cbf398ba612e74fc8a02f))
* **quick-007:** print all config values on startup ([45c7fa6](https://github.com/sergiocarracedo/on-a-meet/commit/45c7fa6e2475e2b346e89b4010df6e18e4260b06))
* **quick-008:** show verbose/silent in config line and add executing message ([a30790c](https://github.com/sergiocarracedo/on-a-meet/commit/a30790cb5dad3d70fa08a54fde51742fb2a017d5))
* **quick-009:** redact JWT tokens from CLI output ([cdf32b8](https://github.com/sergiocarracedo/on-a-meet/commit/cdf32b87c7844ed13e8a3afd52e9b3c4f0173262))
* **quick-010:** expand env vars in commands and fix viper defaults ([f7fdbe9](https://github.com/sergiocarracedo/on-a-meet/commit/f7fdbe9c20c7d44ed91a6df0a12d6436075e2cb5))
* **quick-011:** all cameras selected by default, select-all toggle UX ([57ed001](https://github.com/sergiocarracedo/on-a-meet/commit/57ed0012b928130e80c73d9022e4a194fa6f7d35))
* **quick-012:** fix config loading — add /etc/on-a-meet path and viper defaults ([60dbd9d](https://github.com/sergiocarracedo/on-a-meet/commit/60dbd9d764298866f2cd96a2f16c696807a1003b))
* **quick-013:** improve onboard detection test UX and clarify variable naming ([c747722](https://github.com/sergiocarracedo/on-a-meet/commit/c747722166459301844a2540fce7a3981d3ba15d))
* **quick-014:** make detection test optional, show method change only after failure ([895949c](https://github.com/sergiocarracedo/on-a-meet/commit/895949cdf26a02fbe8d170fb71310d5ba3ad3492))
* **quick-015:** stop service before reinstall, prompt before config overwrite ([ee1987b](https://github.com/sergiocarracedo/on-a-meet/commit/ee1987b8c401a43edce2221b004565238b792cf9))
* **quick-023:** add environment-file config for service unit ([5d58d03](https://github.com/sergiocarracedo/on-a-meet/commit/5d58d03e1315f155e791057e9e984a8ebdbed99a))
* **quick-024:** fire on/off commands on engine startup ([14a10f0](https://github.com/sergiocarracedo/on-a-meet/commit/14a10f063919803403eb60e85b1f649d9b4c06e6))
* **quick-025:** support export KEY=VALUE in env files ([9f7c928](https://github.com/sergiocarracedo/on-a-meet/commit/9f7c928bf13958f96047b248509059f01230ff97))
* **quick-026:** add install.sh script with curl-pipe-install in README ([32b284d](https://github.com/sergiocarracedo/on-a-meet/commit/32b284d48d5d21a76343df342ea3751c4da063ed))
* restore Homebrew tap publishing via GoReleaser ([7b100b1](https://github.com/sergiocarracedo/on-a-meet/commit/7b100b167d73c6d83a16397a6bda696cecf74fd3))


### Bug Fixes

* **quick-001:** set --detect flag default to v4l2 ([4b83107](https://github.com/sergiocarracedo/on-a-meet/commit/4b83107f6d432db1e2c7551f8d2fba68fca4c8b2))

## [1.1.0](https://github.com/sergiocarracedo/on-a-meet/compare/v1.0.0...v1.1.0) (2026-05-29)


### Features

* add release-please automation with auto-merge (ci, release PR, GoReleaser) ([d21b745](https://github.com/sergiocarracedo/on-a-meet/commit/d21b74547efbbd5ed231467aa90b37076db73d22))
* add release-please automation with auto-merge, CI, and Homebrew tap ([40d5960](https://github.com/sergiocarracedo/on-a-meet/commit/40d5960f778298261be78dbf72895050eb04b647))
* add version subcommand ([b80a367](https://github.com/sergiocarracedo/on-a-meet/commit/b80a3674315ad882584fb1faea67b53ed069ef15))
* implement macOS detection backend (MacOSDetector) with log stream + system_profiler ([c3b3087](https://github.com/sergiocarracedo/on-a-meet/commit/c3b30870f8a8af6c61e88238d57b14371bd2760d))
* **quick-026:** add install.sh script with curl-pipe-install in README ([32b284d](https://github.com/sergiocarracedo/on-a-meet/commit/32b284d48d5d21a76343df342ea3751c4da063ed))
* restore Homebrew tap publishing via GoReleaser ([7b100b1](https://github.com/sergiocarracedo/on-a-meet/commit/7b100b167d73c6d83a16397a6bda696cecf74fd3))

## 1.0.0 (2026-05-29)


### Features

* **01-01:** create main.go entry point ([d60a765](https://github.com/sergiocarracedo/on-a-meet/commit/d60a765efb87bf6753e49f64b4fca859f5316b15))
* **01-01:** create root command with Viper config ([ebd33b7](https://github.com/sergiocarracedo/on-a-meet/commit/ebd33b7daaf92f4f5b34b4019169b0e6fd3dec40))
* **01-02:** add --silent/--verbose persistent flags and output wiring ([2446268](https://github.com/sergiocarracedo/on-a-meet/commit/244626856c569d167033df9cef6e516410715b4f))
* **01-02:** create config example and all subcommand stubs ([b1bedc9](https://github.com/sergiocarracedo/on-a-meet/commit/b1bedc9010b79bc4a90e813808f76a0f54a8b88b))
* **01-02:** create internal config package ([738d8d4](https://github.com/sergiocarracedo/on-a-meet/commit/738d8d4cce98d1e793f840eea6ee1d42fb5a28a3))
* **01-02:** create internal output package with pterm wrappers ([f306e53](https://github.com/sergiocarracedo/on-a-meet/commit/f306e532c7b840f2c375f3c2d884584430d6edfd))
* **02-01:** add non-Linux stub for V4L2Detector ([cba7738](https://github.com/sergiocarracedo/on-a-meet/commit/cba7738b8d1b71bc8ded7fafff9986a4225c664d))
* **02-01:** create Detector interface and types ([3e0f334](https://github.com/sergiocarracedo/on-a-meet/commit/3e0f33492ce68196dfacd06f8d374b51a75548b6))
* **02-01:** implement V4L2Detector with ListDevices and Detect ([4f453e6](https://github.com/sergiocarracedo/on-a-meet/commit/4f453e6ccb9d5bafd07c88d2257c19575f1e96c6))
* **02-02:** create polling engine with debounce and hotplug ([a259807](https://github.com/sergiocarracedo/on-a-meet/commit/a25980726323d0173f65d937a54b30a975114d52))
* **02-02:** wire detect command to polling engine ([1cd9630](https://github.com/sergiocarracedo/on-a-meet/commit/1cd9630cd52870f08125ed911cc1cbbd057ab1e6))
* **03-01:** add executor unit tests — timeout, same-state skip, cross-state ([cbf00bd](https://github.com/sergiocarracedo/on-a-meet/commit/cbf00bd7109ad4399315655b496af9f6fe4f5182))
* **03-01:** add Timeout field to Config struct ([9e06233](https://github.com/sergiocarracedo/on-a-meet/commit/9e0623309475c4bd3a8961925f6f8febe602ab60))
* **03-01:** create executor package with TemplateData and Executor struct ([314ea66](https://github.com/sergiocarracedo/on-a-meet/commit/314ea66504dccaebafb8410720bc8940962f0e2c))
* **03-02:** update config.yaml.example with timeout and template docs ([a8da70e](https://github.com/sergiocarracedo/on-a-meet/commit/a8da70e8f4aecbdcd9b2efbe9741a282fa1aae76))
* **03-02:** wire executor into detect command with --timeout flag ([397627e](https://github.com/sergiocarracedo/on-a-meet/commit/397627e3f79de21d3fa82886b61c7a2cfdee1ca8))
* **04-01:** add kardianos/service dependency ([68ef528](https://github.com/sergiocarracedo/on-a-meet/commit/68ef528213b5fa2e1b85a631b78deded63f75811))
* **04-01:** implement install command with kardianos/service ([0cba478](https://github.com/sergiocarracedo/on-a-meet/commit/0cba47876408cea83543e8b928cc9926765eeba0))
* **04-02:** implement uninstall command with kardianos/service ([61b08fc](https://github.com/sergiocarracedo/on-a-meet/commit/61b08fcf0a26987a8db3cf5b5687947e4d3d9a06))
* **05-01:** add lsof backend + factory + tests ([392753c](https://github.com/sergiocarracedo/on-a-meet/commit/392753cffd1ebd29922a7c429f70116e20a9ff49))
* **05-02:** wire factory, list command, goreleaser, README ([4be39d9](https://github.com/sergiocarracedo/on-a-meet/commit/4be39d912bd41d33f74c634575dfc0f0d7d0c1f0))
* add release-please automation with auto-merge (ci, release PR, GoReleaser) ([d21b745](https://github.com/sergiocarracedo/on-a-meet/commit/d21b74547efbbd5ed231467aa90b37076db73d22))
* add release-please automation with auto-merge, CI, and Homebrew tap ([40d5960](https://github.com/sergiocarracedo/on-a-meet/commit/40d5960f778298261be78dbf72895050eb04b647))
* add restart command to restart the service unit ([9af4f39](https://github.com/sergiocarracedo/on-a-meet/commit/9af4f39ba4c2fc7694753fcc3d10a5f9859fdb5e))
* add service subcommand group with start/stop/restart + verbose config ([e81f640](https://github.com/sergiocarracedo/on-a-meet/commit/e81f640b3e5d9f7f314d62f4f88607819202cdf5))
* implement macOS detection backend (MacOSDetector) with log stream + system_profiler ([c3b3087](https://github.com/sergiocarracedo/on-a-meet/commit/c3b30870f8a8af6c61e88238d57b14371bd2760d))
* **phase-6:** interactive onboard wizard with huh ([4566166](https://github.com/sergiocarracedo/on-a-meet/commit/456616600852d04a909b21b9a54533dc95882ffa))
* **phase-6:** sudo apply path and auto re-exec ([a106d4e](https://github.com/sergiocarracedo/on-a-meet/commit/a106d4e962ebc10a1c1afa11bff193e1b12b32e1))
* print config on detect startup ([196ad05](https://github.com/sergiocarracedo/on-a-meet/commit/196ad051787760748f294641b4d24d974dcac429))
* **quick-005:** emit initial camera status on detect startup ([d33a9df](https://github.com/sergiocarracedo/on-a-meet/commit/d33a9df1fd5f8ca0d5daebde638c4d777e439735))
* **quick-006:** log on/off command exit code and verbose output ([8ca2245](https://github.com/sergiocarracedo/on-a-meet/commit/8ca224540bf830228c5cbf398ba612e74fc8a02f))
* **quick-007:** print all config values on startup ([45c7fa6](https://github.com/sergiocarracedo/on-a-meet/commit/45c7fa6e2475e2b346e89b4010df6e18e4260b06))
* **quick-008:** show verbose/silent in config line and add executing message ([a30790c](https://github.com/sergiocarracedo/on-a-meet/commit/a30790cb5dad3d70fa08a54fde51742fb2a017d5))
* **quick-009:** redact JWT tokens from CLI output ([cdf32b8](https://github.com/sergiocarracedo/on-a-meet/commit/cdf32b87c7844ed13e8a3afd52e9b3c4f0173262))
* **quick-010:** expand env vars in commands and fix viper defaults ([f7fdbe9](https://github.com/sergiocarracedo/on-a-meet/commit/f7fdbe9c20c7d44ed91a6df0a12d6436075e2cb5))
* **quick-011:** all cameras selected by default, select-all toggle UX ([57ed001](https://github.com/sergiocarracedo/on-a-meet/commit/57ed0012b928130e80c73d9022e4a194fa6f7d35))
* **quick-012:** fix config loading — add /etc/on-a-meet path and viper defaults ([60dbd9d](https://github.com/sergiocarracedo/on-a-meet/commit/60dbd9d764298866f2cd96a2f16c696807a1003b))
* **quick-013:** improve onboard detection test UX and clarify variable naming ([c747722](https://github.com/sergiocarracedo/on-a-meet/commit/c747722166459301844a2540fce7a3981d3ba15d))
* **quick-014:** make detection test optional, show method change only after failure ([895949c](https://github.com/sergiocarracedo/on-a-meet/commit/895949cdf26a02fbe8d170fb71310d5ba3ad3492))
* **quick-015:** stop service before reinstall, prompt before config overwrite ([ee1987b](https://github.com/sergiocarracedo/on-a-meet/commit/ee1987b8c401a43edce2221b004565238b792cf9))
* **quick-023:** add environment-file config for service unit ([5d58d03](https://github.com/sergiocarracedo/on-a-meet/commit/5d58d03e1315f155e791057e9e984a8ebdbed99a))
* **quick-024:** fire on/off commands on engine startup ([14a10f0](https://github.com/sergiocarracedo/on-a-meet/commit/14a10f063919803403eb60e85b1f649d9b4c06e6))
* **quick-025:** support export KEY=VALUE in env files ([9f7c928](https://github.com/sergiocarracedo/on-a-meet/commit/9f7c928bf13958f96047b248509059f01230ff97))
* **quick-026:** add install.sh script with curl-pipe-install in README ([32b284d](https://github.com/sergiocarracedo/on-a-meet/commit/32b284d48d5d21a76343df342ea3751c4da063ed))
* restore Homebrew tap publishing via GoReleaser ([7b100b1](https://github.com/sergiocarracedo/on-a-meet/commit/7b100b167d73c6d83a16397a6bda696cecf74fd3))


### Bug Fixes

* **quick-001:** set --detect flag default to v4l2 ([4b83107](https://github.com/sergiocarracedo/on-a-meet/commit/4b83107f6d432db1e2c7551f8d2fba68fca4c8b2))
