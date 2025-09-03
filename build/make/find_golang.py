#!/usr/bin/env python3

# Find what we have installed for Golang versus what we are pointing at with GOROOT.
# This is just way too tedious to do in Makefile syntax or even with bash.

# if a target version was not specified or if golang is not installed:
#   exit=1, stdout=None, stderr=usage message
# else:
#   if GOROOT exists and points to a golang installation which is the target version:
#     exit=0, stdout=the path to the go binary file, stderr=None
#   elif GOROOT exists but the go binary does not exist or does not match the target version:
#     exit=1, stdout=None, stderr=no match or bad match message
#   elif GOROOT doesn't exist and /opt/homebrew/Cellar has an @ formula version directory for the target version:
#     exit=0, stdout=the path to the go binary file for the latest matching version, stderr=None
#   elif GOROOT doesn't exist and /opt/homebrew/Cellar/go has a directory for the target version:
#     exit=0, stdout=the path to the go binary file for the latest matching version, stderr=None
#   else:
#     exit=1, stdout=None, stderr=no match message

# To do work on this script within Goland, just add the community-edition python plugin
# then specify the local system python3 as your interpreter since we won't be installing any packages.

import os
from pathlib import Path
import re
import subprocess
import sys
from typing import Optional, Tuple
import importlib
# to avoid needing a pip install we grab the extern'd packaging library from setuptools;
# we want it because packaging.version understands semantic versioning
version = importlib.import_module('setuptools._vendor.packaging.version')


def get_golang_version(go_file: Path) -> version.Version:
    # given the path to the go binary, we can extract the semantic version because 'go' has a 'version' command
    version_string: str = subprocess.check_output([go_file, 'version'], text=True)
    matcher: re.Match = re.search('(?<=go version go)(\S+)\s+.*', version_string)
    unparsed_version: str = matcher.group(1)
    # we are expecting proper semantic versions, which golang uses, so parsing the output should always succeed
    return version.parse(unparsed_version)


def major_minor_versions_equal(v1: version.Version, v2: version.Version) -> bool:
    # we only care about the first two parts of a semantic version because we can't use homebrew to pin any tighter
    return v1.major == v2.major and v1.minor == v2.minor


def go_file_matches_version(go_file: Path, constraint_version: version.Version) -> bool:
    file_version: version.Version = get_golang_version(go_file)
    return major_minor_versions_equal(file_version, constraint_version)


def go_file_from_go_root(go_root: Optional[Path]) -> Optional[Path]:
    # any golang root points to a libexec directory, beneath which we find bin/go
    if not go_root or not go_root.exists():
        return None
    path_to_go: Path = Path(go_root).joinpath('bin/go')
    return path_to_go if path_to_go.exists() else None


def go_root_from_environ() -> Tuple[Optional[Path], bool]:
    # return the path if specified via GOROOT, and a boolean marker to indicate  if there was a GOROOT;
    # GOROOT found means we only will attempt to match to that, while a lack of GOROOT means we search the cellar
    go_root_dir: Optional[str] = os.getenv('GOROOT', None)
    return Path(go_root_dir) if go_root_dir else None, go_root_dir is not None


def go_root_from_parent_dir(go_parent: Path, constraint_version: version.Version) -> Optional[Path]:
    # there could be multiple golang installs within a parent, find the latest version that satisfies the
    # constraint version (yes, this sounds fussy, but 'current' for homebrew will be an ever-evolving situation)
    if not go_parent.exists() or not go_parent.is_dir():
        return None
    match_pattern: str = f"{constraint_version.major}.{constraint_version.minor}.*"
    latest_candidate: Optional[version.Version] = None
    for child in go_parent.iterdir():
        if child.match(match_pattern):
            # we only consider children that are valid for our target version
            if not latest_candidate:
                latest_candidate = child
            elif child > latest_candidate:
                # we use the semantic awareness of version.Version to find the latest installed version
                latest_candidate = child
    return latest_candidate


def go_root_from_brew_formula(constraint_version: version.Version) -> Optional[Path]:
    go_parent: Path = Path(f"/opt/homebrew/Cellar/go@{constraint_version.major}.{constraint_version.minor}")
    return go_root_from_parent_dir(go_parent, constraint_version)


def go_root_from_brew_current(constraint_version: version.Version) -> Optional[Path]:
    go_parent: Path = Path('/opt/homebrew/Cellar/go')
    return go_root_from_parent_dir(go_parent, constraint_version)


EXIT_CODE_SUCCESS: int = 0
EXIT_CODE_FAILURE: int = 1


def succeed_and_show_path(go_file: Path) -> None:
    print(go_file.absolute())
    exit(EXIT_CODE_SUCCESS)


def fail_and_reject_match(go_file: Path, constraint_version: version.Version) -> None:
    print(f"{go_file} does not satisfy version {constraint_version.major}.{constraint_version.minor}", file=sys.stderr)
    exit(EXIT_CODE_FAILURE)


def fail_and_show_no_match(constraint_version: version.Version) -> None:
    print(f"no match found satisfying version {constraint_version.major}.{constraint_version.minor}", file=sys.stderr)
    exit(EXIT_CODE_FAILURE)


def fail_and_show_usage() -> None:
    print(f"usage: python3 {sys.argv[0]} VERSION\n\tVERSION:\tthe target version (n.n[.n])", file=sys.stderr)
    exit(EXIT_CODE_FAILURE)


def main() -> None:
    if len(sys.argv) != 2:
        fail_and_show_usage()
    try:
        constraint_version: version.Version = version.parse(sys.argv[1].strip())
    except version.InvalidVersion:
        # the argument passed was not a valid semantic version
        fail_and_show_usage()
    go_root: Optional[Path]
    go_file: Optional[Path]
    root_specified_via_environment: bool
    go_root, root_specified_via_environment = go_root_from_environ()
    if root_specified_via_environment:
        # GOROOT was specified, so the attempt to find the go binary and match the version is all or nothing
        go_file = go_file_from_go_root(go_root)
        if not go_file:
            fail_and_show_no_match(constraint_version)
        if go_file_matches_version(go_file, constraint_version):
            succeed_and_show_path(go_file)
        else:
            fail_and_reject_match(go_file, constraint_version)
    else:
        # first we try the pinned formula
        go_root = go_root_from_brew_formula(constraint_version)
        go_file = go_file_from_go_root(go_root)
        if go_file:
            succeed_and_show_path(go_file)
        # then we try the current version
        go_root = go_root_from_brew_current(constraint_version)
        go_file = go_file_from_go_root(go_root)
        if go_file:
            succeed_and_show_path(go_file)
    # nothing matched
    fail_and_show_no_match(constraint_version)


if __name__ == '__main__':
    main()
