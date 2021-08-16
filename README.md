# Git-Cache (WIP)

An attempt to re-write https://gitlab.com/grouperenault/git_cdn in Golang

## Problem Statements

When a CI system is sufficiently scaled, CI workers running jobs/tasks in parallel could create
a significant load onto the VCS system(Gitlab, Github, Bitbucket etc...).  For monorepo which has
hundreds to thousands of engineers actively contributing to it, this became a challenge to keep the
VCS system running against the load of CI workers trying to fetch the latest commits constantly.

Additionally, CI workers could be deployed in geo-location that is isolated and/or remotely compare
to the VCS servers.  I.e. There could be a need to have separate workers pool in different timezones
that is closer to your data or to your engineering team.  These remote worker pool can suffer from
high latency when trying to reach the centralized VCS server deployment.  Futhermore, distant networks
bandwidth could be expensive to serve the CI traffic load.

To solve these problems, a local mirror of a VCS server coupling with each CI worker pool is needed.
These mirror shall act as a delayed-cache to the VCS server and deduplicate/serve the requests of CI
workers while still transparently forward the AuthN/AuthZ to centralized VCS server to validate.

## Design

(WIP)

## References

- Git CDN: https://gitlab.com/grouperenault/git_cdn
- Git Protocol V2 Parser: https://github.com/google/gitprotocolio
- Google Unofficial Git Caching Proxy: https://github.com/google/goblet
- Gitlab's Gitaly Pack-Objects Cache: https://gitlab.com/gitlab-org/gitaly/-/blob/master/doc/design_pack_objects_cache.md
