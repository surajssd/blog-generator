# blog-generator

This is a blog generator for [Kubernetes Bangalore](https://twitter.com/k8sBLR) community event reports.

**Note**: This is a heavily WIP project and needs lot of work to even have a good starting point.

## Install

```bash
go install
```

## Usage

Generate keys from your twitter account and use as follows:

```bash
blog-generator --consumer-key [REDACTED] \
               --consumer-secret [REDACTED] \
               --access-token [REDACTED] \
               --access-secret [REDACTED] \
               --date "Sat Mar 16"
```

## Roadmap

- Make the tool easier to use via command line
- Add a config file option to generate the blog via one single config file than having many different flags.

