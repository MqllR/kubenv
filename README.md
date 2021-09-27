# kubenv

[![Actions Status](https://github.com/MqllR/kubenv/workflows/Publish/badge.svg)](https://github.com/MqllR/kubenv/actions)

**kubenv** is a tool to ease the management of multiple kubernetes cluster.

Features:
- Merge different kubeconfig file in different path into a single big kubeconfig
- Switch between kubernetes context
- Execute a command using either a single or mutliple contexts

## Install

Grap the latest release in github or build by yourself:

```
go get https://github.com/mqllr/kubenv
```

Define your configuration file as [kubenv-example.yaml](https://github.com/MqllR/kubenv/blob/master/example/kubenv_example.yaml) and export the environment variable:

```
export KUBENV_CONFIG=/path/to/my/config.yaml
```

## Configuration

 ### k8sConfigs

The k8sConfigs section define the sync mode. 2 modes available: `local` for local files and `exec` for command execution. The exec mode will capture the command output.

```yaml
  dev:
    sync:
      mode: local
      path: /tmp/k8senv/dev/config
  kind:
    sync:
      mode: exec
      command:
        - bash
        - -c
        - |
          kind export -q kubeconfig --kubeconfig /tmp/test && cat /tmp/test
```

## Usage

Get started by picking up your kubeconfig files:

```
kubenv sync
▸ Start the synchronization of kubeconfig file into /home/mql/.kube/config ...
Sync kubeconfig foo ✔
Sync kubeconfig bar ✔
```

Then you can change your default context by using the `use-context`

```
kubenv use-context
```

or execute a command against a single or multiple contexts:

```
kubenv with-context
```

The tool uses https://github.com/AlecAivazis/survey to navigate between contexts. If you prefer to use j/k to go down or up, press `esc`:

> The user can also press esc to toggle the ability cycle through the options with the j and k keys to do down and up respectively.
