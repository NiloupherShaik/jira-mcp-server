## Will initialise Gitguardian
run-gitguardian-login:
	@ggshield auth login --lifetime 180

## Will run pre-commit on staged changes
pre-commit-update:
	@pre-commit autoupdate --config .pre-commit-config.yaml

## Will run pre-commit on staged changes
pre-commit-run:
	@pre-commit run --config .pre-commit-config.yaml

## Will run pre-commit on full project
pre-commit-run-all:
	@pre-commit run --all-files --config .pre-commit-config.yaml

## Will initialise pre-commit
pre-commit-install:
	@pre-commit install --config .pre-commit-config.yaml --install-hooks

## Will remove pre-commit
pre-commit-uninstall:
	@pre-commit uninstall --config .pre-commit-config.yaml

## Fully set up your pre-commit environment
## PSA: This adds a directory in your root to add git-templates
pre-commit-setup:
	@git config --global init.templateDir ~/.git-template
	@pre-commit init-templatedir ~/.git-template
	@make pre-commit-install

# Runs commitizen to check the validity of the commit message
commitizen-run-check:
	@cz check --commit-msg-file .git/COMMIT_EDITMSG
