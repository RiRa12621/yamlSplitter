# Yaml Splitter
This utility is used to split concatenated kubernetes yamls into separate files.

The use case is if you want to pull a yaml from upstream and integrate it into
your own code, but maintainability and readability is a lot better with separate
files.

An example is https://docs.k0smotron.io/v1.2.0/install.yaml