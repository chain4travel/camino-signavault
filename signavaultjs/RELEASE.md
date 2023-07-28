# Changelog

All notable changes to this project will be documented in this file.


## [2.0.0]

### Changed
- Multisig tx id generation logic. Now it is based on a hash of the enclosed transaction followed by the current timestamp.

### Fixed
- Bug fix: failure to create a new multisig tx if an identical tx has already been created and expired.
See also [Bugfix/change multisigtx id generation](https://github.com/chain4travel/camino-signavault/pull/56)

## [1.0.11]

Initial release.
