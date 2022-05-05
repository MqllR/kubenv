# kubenv

[![Actions Status](https://github.com/MqllR/kubenv/workflows/Publish/badge.svg)](https://github.com/MqllR/kubenv/actions)

**kubenv** is a tool to ease the management of multiple kubernetes cluster.

Features:
- Merge different kubeconfig files from different path into a single big kubeconfig
- Switch between kubernetes context
- Execute a command using either a single or mutliple contexts

## Install

Grabe the latest release in github release or build it by yourself:

```
go get https://github.com/mqllr/kubenv
```

## Usage

Get started by picking up your kubeconfig files:

```
kubenv sync -a=false -m=exec --command="cat /path/to/kubeconfig"
▸ Start the synchronization of kubeconfig file into /home/john/.kube/config ...
```

Which is equivalent to:

```
kubenv sync -a=false -m=local --path="/path/to/kubeconfig"
▸ Start the synchronization of kubeconfig file into /home/john/.kube/config ...
```

And also:

```
kubenv sync -a=false -m=glob --glob="path/to/kubeconfig"
▸ Start the synchronization of kubeconfig file into /home/john/.kube/config ...
```

Then you can change your default context by using the `use-context`

```
kubenv use-context
```

or execute a command against a single or multiple contexts:

```
kubenv with-context
```

The CLI uses https://github.com/AlecAivazis/survey to navigate between contexts. If you prefer to use j/k to go down or up, press `esc`:

> The user can also press esc to toggle the ability cycle through the options with the j and k keys to do down and up respectively.

## Bonus

Want the CLI to be part of kubectl?

```
ln -s /home/john/go/bin/kubenv /usr/local/bin/kubectl-env
```
