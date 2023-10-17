# Changelog

## [1.12.0](https://github.com/soerenschneider/acmevault/compare/v1.11.0...v1.12.0) (2023-10-17)


### Features

* automatically renew vault token before expiry ([e317ba1](https://github.com/soerenschneider/acmevault/commit/e317ba1148403464d2a651421c371fca73953ad1))
* yaml config file support ([#340](https://github.com/soerenschneider/acmevault/issues/340)) ([ab7577b](https://github.com/soerenschneider/acmevault/commit/ab7577bc3eb646fbc23595094b33240d5837c006))


### Bug Fixes

* Bump github.com/go-acme/lego/v4 from 4.14.0 to 4.14.2 ([1393e09](https://github.com/soerenschneider/acmevault/commit/1393e091926fe7bb6b2d0d4a14c1ad57517eec04))
* **deps:** bump github.com/aws/aws-sdk-go-v2/config ([79e1fd3](https://github.com/soerenschneider/acmevault/commit/79e1fd32c31463434ea509c786374435a8174078))
* **deps:** bump github.com/aws/aws-sdk-go-v2/config ([#339](https://github.com/soerenschneider/acmevault/issues/339)) ([3a561ab](https://github.com/soerenschneider/acmevault/commit/3a561abe455c6559ae22ceae1c9537643ad3ddcc))
* **deps:** bump github.com/aws/aws-sdk-go-v2/config ([#344](https://github.com/soerenschneider/acmevault/issues/344)) ([dab082f](https://github.com/soerenschneider/acmevault/commit/dab082feeb7586e3c384fb1c1b68e18d8abdf712))
* **deps:** bump github.com/aws/aws-sdk-go-v2/service/route53 ([586af28](https://github.com/soerenschneider/acmevault/commit/586af282cc7689ff2f4f4c835d9204b202b97b1b))
* **deps:** bump github.com/aws/aws-sdk-go-v2/service/route53 ([#343](https://github.com/soerenschneider/acmevault/issues/343)) ([2151f48](https://github.com/soerenschneider/acmevault/commit/2151f489bbefe9895e0e72d5d309aee0707e1d0d))
* **deps:** bump github.com/go-playground/validator/v10 ([#338](https://github.com/soerenschneider/acmevault/issues/338)) ([699d9bb](https://github.com/soerenschneider/acmevault/commit/699d9bb3269c8700e920aecbb31262b35909a787))
* **deps:** bump github.com/prometheus/client_golang ([4d2c991](https://github.com/soerenschneider/acmevault/commit/4d2c991cf650ee9a895f1d8b2fd8c7cf19d14dbe))
* **deps:** bump github.com/rs/zerolog from 1.30.0 to 1.31.0 ([eb81c37](https://github.com/soerenschneider/acmevault/commit/eb81c37b28ee2ad79f747a73e13cf059cafc0e86))
* **deps:** bump golang from 1.21.1 to 1.21.2 ([#334](https://github.com/soerenschneider/acmevault/issues/334)) ([1022bef](https://github.com/soerenschneider/acmevault/commit/1022bef84c254e5e4f6e8b9e8d24b2a44be61b02))
* **deps:** bump golang to 1.21.3 ([68c02fc](https://github.com/soerenschneider/acmevault/commit/68c02fcd202bdc95059e2205c1a0c769b268fd90))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([#341](https://github.com/soerenschneider/acmevault/issues/341)) ([d8ec969](https://github.com/soerenschneider/acmevault/commit/d8ec969b1c1f6df60129e70dffa6142e8c768baa))

## [1.11.0](https://github.com/soerenschneider/acmevault/compare/v1.10.0...v1.11.0) (2023-09-11)


### Features

* add kubernetes vault auth ([0abdc46](https://github.com/soerenschneider/acmevault/commit/0abdc46f717db8c621e3e8675c07ac250c263b6e))


### Bug Fixes

* add 'implicit' as oneof option ([d3852bf](https://github.com/soerenschneider/acmevault/commit/d3852bf531a70b347add47874977a673c1bf0fcd))
* fix implicit vault auth ([917bfdf](https://github.com/soerenschneider/acmevault/commit/917bfdfaa34eda5a680013d9f2f84aa0528f9b26))
* fix log statement ([43a2b4a](https://github.com/soerenschneider/acmevault/commit/43a2b4a6d30045683f8bf7a0671da3a82cc4b5ac))
* fix logline for good ([75c3744](https://github.com/soerenschneider/acmevault/commit/75c3744c9f021fd28814972e4314088e6cf79b00))

## [1.10.0](https://github.com/soerenschneider/acmevault/compare/v1.9.0...v1.10.0) (2023-07-13)


### Features

* add support for more dynamic renewal ([995437f](https://github.com/soerenschneider/acmevault/commit/995437f690a438740263709e172939e53ac88303))
* validate using validator package ([a3cada0](https://github.com/soerenschneider/acmevault/commit/a3cada06fcbbc660403ec6fd3a3a84e31ae71a83))


### Bug Fixes

* check for error ([a075efc](https://github.com/soerenschneider/acmevault/commit/a075efc746096c07ef879797dc43642e1e494082))
* fix logic ([ba0ffc9](https://github.com/soerenschneider/acmevault/commit/ba0ffc999c680e561b1205cb20e0e47db2f2b883))

## [1.9.0](https://github.com/soerenschneider/acmevault/compare/v1.8.0...v1.9.0) (2022-08-17)


### Features

* Implement custom dns servers ([1759626](https://github.com/soerenschneider/acmevault/commit/17596261e2ddd28c5723c923abf54b1202d8dab5))

## [1.8.0](https://github.com/soerenschneider/acmevault/compare/v1.7.3...v1.8.0) (2022-06-17)


### Features

* Support SANs ([#157](https://github.com/soerenschneider/acmevault/issues/157)) ([c2953ee](https://github.com/soerenschneider/acmevault/commit/c2953eeaee0e94aa3c3334cf88caf6774756c475))

## [1.7.3](https://github.com/soerenschneider/acmevault/compare/v1.7.2...v1.7.3) (2022-06-06)


### Bug Fixes

* Fix metric dimension ([ad633ad](https://github.com/soerenschneider/acmevault/commit/ad633ad40cee2cdb815de5c87536cbfe14bc0d68))

### [1.7.2](https://www.github.com/soerenschneider/acmevault/compare/v1.7.1...v1.7.2) (2022-06-04)


### Miscellaneous Chores

* release 1.7.2 ([ee42f8f](https://www.github.com/soerenschneider/acmevault/commit/ee42f8f21091251f6bcce38a90399f23892688d3))

### [1.7.1](https://www.github.com/soerenschneider/acmevault/compare/v1.7.0...v1.7.1) (2022-02-07)


### Miscellaneous Chores

* release 1.7.1 ([615d63a](https://www.github.com/soerenschneider/acmevault/commit/615d63a159f7160ffdc88d70b81e0c23020a4a38))

## [1.7.0](https://www.github.com/soerenschneider/acmevault/compare/v1.6.2...v1.7.0) (2022-01-29)


### Features

* domain format ([#108](https://www.github.com/soerenschneider/acmevault/issues/108)) ([25f0779](https://www.github.com/soerenschneider/acmevault/commit/25f0779bf921900280ad353f08f59182ba39d9ca))

### [1.6.2](https://www.github.com/soerenschneider/acmevault/compare/v1.6.1...v1.6.2) (2021-12-31)


### Bug Fixes

* Also write metrics in case of no error ([c6b28fb](https://www.github.com/soerenschneider/acmevault/commit/c6b28fb055f980abc35cc1a8728273e0e563b7f4))
* authenticate before reading acme account from storage ([d6d0324](https://www.github.com/soerenschneider/acmevault/commit/d6d0324cf837de3ba9f087ddf06de0c87309bc76))

### [1.6.1](https://www.github.com/soerenschneider/acmevault/compare/v1.6.0...v1.6.1) (2021-10-20)


### Bug Fixes

* Always write prom metrics before exit ([c254b0c](https://www.github.com/soerenschneider/acmevault/commit/c254b0c802692376dc206dd1ffcd073ae9554851))

## [1.6.0](https://www.github.com/soerenschneider/acmevault/compare/v1.5.0...v1.6.0) (2021-10-19)


### Features

* enable loading wrapped token from disk ([d07494c](https://www.github.com/soerenschneider/acmevault/commit/d07494c41295cfedaee4481aea00084cbbdbf792))

## [1.5.0](https://www.github.com/soerenschneider/acmevault/compare/v1.4.0...v1.5.0) (2021-10-18)


### Features

* Support for response wrapped secret_ids ([#79](https://www.github.com/soerenschneider/acmevault/issues/79)) ([8c1c29d](https://www.github.com/soerenschneider/acmevault/commit/8c1c29d2825d8726ee58153fd7340c08111bc5ae))


### Bug Fixes

* always write cert expiry metrics ([33535b7](https://www.github.com/soerenschneider/acmevault/commit/33535b799944a56df16a8313dcd3260b88857ae5))
* eliminate possible newline after the secret_id ([e7452b9](https://www.github.com/soerenschneider/acmevault/commit/e7452b9168f0b52e07cf401f4de5a795f3a74693))

## [1.4.0](https://www.github.com/soerenschneider/acmevault/compare/v1.3.1...v1.4.0) (2021-10-14)


### Features

* Allow writing data as PEM file ([5678f14](https://www.github.com/soerenschneider/acmevault/commit/5678f14cc0380d9621dd06dfec4f957c8a8b9af1))


### Bug Fixes

* check if files are configured before trying to write ([496b116](https://www.github.com/soerenschneider/acmevault/commit/496b116a75f480b27e09374200badb0f4c15ac61))
* do not include issuer cert in pem ([ef0b324](https://www.github.com/soerenschneider/acmevault/commit/ef0b3246edacf0d846777d877c5cf24e865bc86d))
* fix incorrect check using the wrong predicate ([1b9f667](https://www.github.com/soerenschneider/acmevault/commit/1b9f6672a10c5ef632a3a9ff0173423578915c35))
* fix json key name ([1c234e1](https://www.github.com/soerenschneider/acmevault/commit/1c234e1fd17139550bca47e39518459eabd3769a))

### [1.3.1](https://www.github.com/soerenschneider/acmevault/compare/v1.3.0...v1.3.1) (2021-10-13)


### Bug Fixes

* Actually validate config for server ([6d9b969](https://www.github.com/soerenschneider/acmevault/commit/6d9b96992e477fa103952daedb2b3104d962717b))
* Check if pathPrefix is empty ([28cc23f](https://www.github.com/soerenschneider/acmevault/commit/28cc23fd016081486d053c1af1f9a491a901e026))

## [1.3.0](https://www.github.com/soerenschneider/acmevault/compare/v1.2.1...v1.3.0) (2021-10-12)


### Features

* Better validation ([06fec46](https://www.github.com/soerenschneider/acmevault/commit/06fec462ccc9bd48db8c69e99970a4b0057659c3))
* Enable separation for reading private keys and pub cert data ([413e281](https://www.github.com/soerenschneider/acmevault/commit/413e281ff6c3208ea0574b9ee8b0a5c05817a5ab))
* print version information in log ([19c9788](https://www.github.com/soerenschneider/acmevault/commit/19c97887bf5d91f7d5daa849918e3059ae8953db))


### Bug Fixes

* fix log output when printing config ([89751cb](https://www.github.com/soerenschneider/acmevault/commit/89751cbf3e08abe4c040cf7b8a069f4b9a6f9b72))
* increase privacy of not recycling old private keys ([e2cbad8](https://www.github.com/soerenschneider/acmevault/commit/e2cbad89b4fb08c260d3cc3b16d1c0a4e33d633d))

### [1.2.1](https://www.github.com/soerenschneider/acmevault/compare/v1.2.0...v1.2.1) (2021-10-05)


### Bug Fixes

* Remove role_id prefix from domain ([74006b7](https://www.github.com/soerenschneider/acmevault/commit/74006b77c5df5f7176ef3cac7b1743252edfb69f))

## [1.2.0](https://www.github.com/soerenschneider/acmevault/compare/v1.1.0...v1.2.0) (2021-10-05)


### Features

* Add metric for errors retrieving aws credentials ([a17db0a](https://www.github.com/soerenschneider/acmevault/commit/a17db0a0dadc2ed2cb6597876f7b45fa65d1ff68))
* support for KV2 secret engine ([d3d9337](https://www.github.com/soerenschneider/acmevault/commit/d3d93379782a6eabb745504ef8ba60921805475c))


### Bug Fixes

* Fix logic of dynamic credential providers ([46e83f0](https://www.github.com/soerenschneider/acmevault/commit/46e83f0123c22fda6ad28d1a281a0be2439146e3))
* Fix using wrong function ([31cfef6](https://www.github.com/soerenschneider/acmevault/commit/31cfef678760ee51bf2e00f5a10aaa11090650a8))

## [1.1.0](https://www.github.com/soerenschneider/acmevault/compare/v1.0.2...v1.1.0) (2021-09-28)


### Features

* Print version info ([c716245](https://www.github.com/soerenschneider/acmevault/commit/c716245980ef7710023be79a85c697f07c59eb73))

### [1.0.2](https://www.github.com/soerenschneider/acmevault/compare/v1.0.1...v1.0.2) (2021-09-13)


### Bug Fixes

* fix insufficient error handling ([de707ab](https://www.github.com/soerenschneider/acmevault/commit/de707ab49f5435a08310c8eff85d249d400df5d5))
* fix insufficient error handling ([86bcde2](https://www.github.com/soerenschneider/acmevault/commit/86bcde2728fdf38004db070fd2344fd164c61117))

### [1.0.1](https://www.github.com/soerenschneider/acmevault/compare/v1.0.0...v1.0.1) (2021-09-13)


### Miscellaneous Chores

* release 1.0.1 ([6b28d83](https://www.github.com/soerenschneider/acmevault/commit/6b28d837945f77b1c039005fcd91f6aaeecdf4a4))

## 1.0.0 (2021-06-16)


### Features

* Init ([66ef063](https://www.github.com/soerenschneider/acmevault/commit/66ef0639375236678cffccf59a3eb47f553bcb7f))
