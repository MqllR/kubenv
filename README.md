# kubenv

[![Actions Status](https://github.com/MqllR/kubenv/workflows/Build%20and%20release/badge.svg)](https://github.com/MqllR/kubenv/actions)

**kubenv** is a tool to handle different authentication method on AWS and to deal with multiple kubeconfig files. When working on tens of kubernetes clusters and AWS accounts, it becomes useful to have a tool making the glue between each of them.

Features:
- Wrap tools for AWS authentication when using a federation with an indentity provider. Currently support: [aws-google-auth](https://github.com/cevoaustralia/aws-google-auth), [aws-azure-login](https://github.com/sportradar/aws-azure-login)
- Do an AsumeRole using AWS STS with another AWS profile
- Check the external tools local and remote version
- Install the external tools. Currently for aws-google-auth and aws-iam-authenticator
- Merge different kubeconfig file in different path into a single big kubeconfig
- Bind a kubeconfig to an AWS profile
- Switch between kubernetes context

## Install

Grap the latest release in github or build by yourself:

```
go build
mv kubenv /usr/local/bin/
```

Define your configuration file as [kubenv-example.yaml](https://github.com/MqllR/kubenv/blob/master/kubenv_example.yaml) and export the environment variable:

```
export KUBENV_CONFIG=/path/to/my/config.yaml
```

## Configuration

The configuration is devided in 3 main keys: authProviders, authAccounts, k8sConfigs.

### authProviders

The authProviders give the general authentication provider parameters. 3 available provider:

- aws-google-auth

```yaml
  aws-google-auth:
    IDP: "Bzbj37S"
    SP: "91010906235"
    UserName: "my@e.mail"
```

- aws-azure-login

```yaml
  aws-azure-login:
    TenantID: "Bzbj37S"
    AppIDUri: "91010906235"
    UserName: "my@e.mail"
 ```

- aws-sts

```yaml
  aws-sts:
    UserName: "my@e.mail"
```

### authAccounts

The authAccounts declare the differents AWS profile and the authentication method

- with aws-google-auth:
```yaml
  dev:
    AuthProvider: aws-google-auth
    AWSProfile: dev 
    AWSRole: arn:aws:iam::121638826155:role/AdminDev
    AWSRegion: eu-central-1
```

- with aws-azure-login:

```yaml
  test:
    AuthProvider: aws-azure-login
    AWSProfile: test
    AWSRole: arn:aws:iam::125262463473:role/AdminTest
    Duration: 36000    # In seconds
```

- with aws-sts:

```yaml
  doublejump:
    AuthProvider: aws-sts
    AWSProfile: doublejump
    AWSRole: arn:aws:iam::125262463543:role/DoubleJumpTest
    DependsOn: test
 ```
 
 ### k8sConfigs
 
The k8sConfigs parameter define the path de kubeconfig to sync and optionaly a reference to an authAccount to inject the environment variable AWS_PROFILE into the kubeconfig exec section.

```yaml
  dev:
    sync:
      mode: local
      path: /tmp/k8senv/dev/config
    authAccount: dev 
```

## Getting started

After writing your kubenv config, you'll be able to build your kubeconfig file:

```
# kubenv k8s sync
Sync kubeconfig dev ✔
```

Ensure the external tools are well installed:

```
# kubenv dep check
✔ aws-google-auth: local: 0.0.36        remote: 0.0.36
```

If not, use the dep install command to install the tool

Finally, authenticate yourself:

```
# kubenv auth
✔ doublejump
▸ Authentication using aws-sts...
⚠ Depends on test
▸ Authentication using aws-google-auth...
✔ Token already active. Skipping.
```
