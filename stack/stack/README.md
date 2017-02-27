# Docker/docker/cli/command/stack

This code is directly copied from the docker repo, except that internal methods
have been exported so that we can use them in this tool.

Elements that were removed were done only to reduce the dependency base:

1. github.com/spf13/pflag was added manually
2. Any CobraCommand references and methods were removed
3. The various Context references were replaced with the core context package.
