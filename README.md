<div align="center">

# MCP Server Template

[![Backstage][repo_backstage_img]][repo_backstage_url]

A repo to help bootstrap the creation of [MCP Servers][mcp_information_url] for use within N-able.

</div>

## 📖 Project Documentation

The best way to better explore all the features of the **MCP Server Template**
is to read the AI [Backstage documentation][repo_backstage_url].

## ⚡️ Quick start

1. Install [brew][brew_install_url]
2. Install [make][make_brew_url]
3. [Clone][repo_clone_url] this repo
4. Run the following command to install the pre-requisites tools.
    ```bash
    make run-housekeeping
    ```
5. Setup your Artifactory credentials following this [guide][artifactory_credentials_url].
6. Run:
    ```shell
    make go-env
    ```
7. Setup local cluster
    ```shell
    make clean-setup-local-kind-cluster
    ```
8. Build!
    ```shell
    make dev
    ```
9. You can use [postman][postman_docs_url] to interact with the MCP server.

## Contribution

We welcome contributions to this project, but we do ask that the hooks and commit conventions are followed to ensure PRs
are easy to review and a common style is applied across the codebase.

### Pre-Commit

[Precommit][precommit_repo_url] hooks are a set of scripts that run before a commit or push is made. They help ensure that the code adheres
to the project style, that security checks are run, and that the code is formatted correctly.

To create the various hooks in your environment, run the following command:

```shell
   make pre-commit-install
```

#### Goland

Setup code completion by [following this guide][goland_linter_url].

#### VS Code

Setup code completion by [following this guide][vscode_linter_url].

### Conventional Commits

[Conventional Commits][conventional_commits_url] are used in this repo to provide a consistent approach to commit
messages which aids with the release process, eg. [NextStage pipelines][nextstage_pipeline_url] will read the commit messages for
determining the next release version.

[Commitizen][commitizen_utility_url] is used to help with the commit message creation. To use it, run:

```shell
  cz commit
```

and follow the prompts to create a commit message that adheres to the conventional commits standard.

#### Goland

Install the conventional commit plugin via the [Goland marketplace][goland_conventional_commit_plugin_url]

#### VS Code

Setup code completion by [following this guide][vscode_conventional_commit_plugin_url].

### ☁️ AWS

N-able uses AWS to host the applications across the various environments and regions. To access the AWS resources, you
need
to have the AWS CLI (see [Quick Start](#run-cluster-locally)) installed and configured on your local machine.

1. Run the following [anvil][anvil_repo_url] command to add all the relevant AWS credentials to your local environment.
   ```bash
   anvil aws config \
    --idp-url https://access.n-able.dev/auth/realms/NABL/protocol/saml/clients/amazon-aws \
    --idp-realm nabl
   ```
2. Now run the following command to add the AWS [EKS][eks_cluster_url] kubecontexts to your local environment.
   ```bash
   make setup-k8s-contexts
   ```

### Gotcha's

#### Kustomize Version

The template is configured to work with a specific version of kustomize - v5.0.3.

To change the version, and then pin the proper version, run the following command.

- WSL2
    ```shell
    brew uninstall kustomize || curl -L -H "Authorization: Bearer QQ==" -o kustomize-5.0.3.x86_64_linux.bottle.tar.gz https://ghcr.io/v2/homebrew/core/kustomize/blobs/sha256:64fafd82844704588435ef680441bc24d32e2f69c81fc981b07d91cba11e6903 -v && brew install kustomize-5.0.3.x86_64_linux.bottle.tar.gz && brew pin kustomize
    ```
- 🍎
   ```shell
   brew uninstall kustomize || curl -L -H "Authorization: Bearer QQ==" -o kustomize-5.0.3.sonoma.bottle.tar.gz https://ghcr.io/v2/homebrew/core/kustomize/blobs/sha256:dbdc6f593427a2b984ded62f9630f855b89a4b51f4a2073d2cb2472b1fe9508e -v && brew install kustomize-5.0.3.sonoma.bottle.tar.gz
   ```
---

[//]: # (Links should be alphabetised for readability)
<!-- Anvil -->

[anvil_repo_url]: https://github.com/logicnow/anvil

<!-- Artifactory -->

[artifactory_credentials_url]: https://portal.n-able.dev/docs/rm/component/fusion-metrics-validator/howto/artifactory/

<!-- AWS -->

[eks_cluster_url]: https://console.aws.amazon.com/eks/clusters

<!-- Brew -->

[brew_install_url]: https://brew.sh/

<!-- Conventional Commits -->

[conventional_commits_url]: https://www.conventionalcommits.org/en/v1.0.0/

[commitizen_utility_url]: https://commitizen-tools.github.io/commitizen/

[goland_conventional_commit_plugin_url]: https://plugins.jetbrains.com/plugin/13389-conventional-commit

[vscode_conventional_commit_plugin_url]: https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits

<!-- Linter -->

[goland_linter_url]: https://github.com/mvdan/gofumpt?tab=readme-ov-file#goland

[vscode_linter_url]: https://github.com/mvdan/gofumpt?tab=readme-ov-file#visual-studio-code

<!-- Make -->

[make_brew_url]: https://formulae.brew.sh/formula/make

<!-- MCP -->

[mcp_information_url]: https://modelcontextprotocol.io/quickstart/server

<!-- NextStage -->

[nextstage_pipeline_url]: https://portal.n-able.dev/docs/nextstage/system/nextstage/NSPipelines/NextStage-Pipelines/

<!-- Postman -->
[postman_docs_url]: ../reference/postman-mcp.md

<!-- Precommit -->
[precommit_repo_url]: https://pre-commit.com/


<!-- Repository -->

[howto_local_env_url]: docs/howto/set-up-local-environment.md

[repo_backstage_url]: https://portal.n-able.dev/catalog/default/system/mcp/docs

[repo_backstage_img]: https://img.shields.io/badge/Docs-Backstage?logo=backstage&logoColor=white&labelColor=%239BF0E1&color=white

[repo_clone_url]: https://github.com/nable-fusion/mcp-server-template

[repo_logo_img]: https://raw.githubusercontent.com/nable-fusion/mcp-server-template/main/docs/mcp-server-logo.png
