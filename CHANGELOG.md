<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- [#16] Add [authz module](https://github.com/cosmos/cosmos-sdk/tree/master/x/authz/spec)

## [v1.1.0]((https://github.com/ChihuahuaChain/chihuahua/releases/tag/v1.1.0) - 2022-01-02
- [#1](https://github.com/ChihuahuaChain/chihuahua/pull/2) Version bumps, add mainnet files, many improvements and fixes
- [#2](https://github.com/pomifer/chihuahua/pull/1) Add a minimum validator commission of 5% based on proposal [#1](https://omniflix.chihuahua.wtf/proposals)

## [v1.0.0](https://github.com/ChihuahuaChain/chihuahua/releases/tag/v1.0.0) - 2021-12-18

Release mainnet
