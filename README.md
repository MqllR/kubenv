# kubenv

[![Actions Status](https://github.com/MqllR/kubenv/workflows/Publish/badge.svg)](https://github.com/MqllR/kubenv/actions)

**kubenv** makes working with multiple Kubernetes clusters painless. Merge all your kubeconfig files into one, switch contexts interactively, and run commands across several clusters at once.

## Features

- **Merge** kubeconfig files from different locations into a single kubeconfig
- **Switch** between Kubernetes contexts with an interactive prompt
- **Run** a command against one or several contexts in one go
- **Inspect, edit and delete** contexts without opening the kubeconfig by hand
- **Plug into kubectl** as a plugin (see [Bonus](#bonus))

## Installation

Download the latest binary from the [GitHub releases](https://github.com/MqllR/kubenv/releases), or install it with Go:

```
go install github.com/mqllr/kubenv@latest
```

## Usage

### 1. Import your kubeconfig files

Start by collecting your kubeconfig files. kubenv can read them from a local file, a glob pattern, or the output of a command:

```
# From a command
kubenv sync -a=false -m=exec --command="cat /path/to/kubeconfig"
▸ Start the synchronization of kubeconfig file into /home/john/.kube/config ...

# From a local file
kubenv sync -a=false -m=local --path="/path/to/kubeconfig"

# From a glob pattern
kubenv sync -a=false -m=glob --glob="path/to/kubeconfig"
```

### 2. Switch context

Pick your default context from an interactive list:

```
kubenv use-context
```

### 3. Run a command on one or several contexts

```
kubenv with-context
```

### Other commands

| Command | Description |
|---|---|
| `kubenv show` | Show information about your contexts |
| `kubenv edit` | Edit a context manually |
| `kubenv delete` | Delete a context |
| `kubenv version` | Print the version |

Run `kubenv --help` to see every command and flag.

### Navigation tips

kubenv uses [survey](https://github.com/AlecAivazis/survey) for its prompts. If you prefer `j`/`k` to move down and up, press `esc`:

> The user can also press esc to toggle the ability to cycle through the options with the j and k keys to go down and up respectively.

## Bonus

Want kubenv to feel like part of kubectl? Create a symlink named `kubectl-env`, and kubectl will pick it up as a plugin:

```
ln -s /home/john/go/bin/kubenv /usr/local/bin/kubectl-env
```

You can then run `kubectl env use-context`.
