# object-cache

This project is a object cache which is used to cache dependencies of CI projects. It provides a store API which can be used to store and retrieve objects.

# Building

This project uses [Bazel](https://bazel.build/) as the build system and Open Container Initiative (OCI) for packaging.

To build the project, run the following command:

```bash
bazel build //...
```

To  run the tests, run the following command:

```bash
bazel test //...
```

To export the container image, run the following command:

```bash
bazel run //services/store:load
```

## Disclaimer

This project is in the early stages of development and is not yet ready for use.

## License

This project is released under Apache 2.0 license.
