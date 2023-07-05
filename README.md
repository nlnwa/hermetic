# go_container

Template for creating a Go application in a Docker container

## Motivation

This repo exists to more easily create a container based service with `golang`

## TODO

When using this repo as a template, remember to do the following

### Set up Secret scanning

Access `https://github.com/<ORG>/<REPO>/settings/security_analysis` and enable
`Secret scanning` and its sub-option `Push protection`.

### Set up branch protection rules

The default (`main`) branch should be protected, so please ensure that the
branch is sufficiently protected. The following should be enabled:

* Require minimum number of PR reviewers
* Require `Test`, `Lint`, `CodeQl build` and `CodeQL`(from Github Code
  Scanning) workflow jobs to pass before merging
* Require branch is up to date before merging
* Do not allow bypassing the above settings

### Update this readme

Please rewrite this readme to reflect your project.
The actual contents depends on the project you have.

#### Add licensing documents for 3rdParty

If you have 3rd party dependencies, their licenses should be placed in
`/licenses-dependencies/<name-of-dependency>`

### Set up dependabot

Create a new Github label `dependencies` by following
`https://github.com/<ORG>/<REPO>/issues/labels`, as it is required by
`dependabot`.

Set up dependabot by accessing
`https://github.com/<ORG>/<REPO>/settings/security_analysis` and enabling the
following: `Dependabot alerts` and `Dependabot security updates`.
