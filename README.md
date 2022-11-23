<p align="center">
  <img style="width:44%" src="images/kubefirst.svg" alt="Kubefirst Logo"/>
</p>

<p align="center">
  GitOps Infrastructure & Application Delivery Platform
</p>

<p align="center">
  <a href="https://docs.kubefirst.io/">Install</a>&nbsp;|&nbsp;
  <a href="https://twitter.com/kubefirst">Twitter</a>&nbsp;|&nbsp;
  <a href="https://www.linkedin.com/company/kubefirst">LinkedIn</a>&nbsp;|&nbsp;
  <a href="https://join.slack.com/t/kubefirst/shared_invite/zt-r0r9cfts-OVnH0ooELDLm9n9p2aU7fw">Slack</a>&nbsp;|&nbsp;
  <a href="https://kubeshop.io/blog-projects/kubefirst">Blog</a>
</p>

<p align="center">
  <a href="https://github.com/kubefirst/kubefirst/releases"><img title="Release" src="https://img.shields.io/github/v/release/kubefirst/kubefirst"/></a>
  <!-- <a href=""><img title="Docker builds" src="https://img.shields.io/docker/automated/kubeshop/tracetest"/></a> -->
  <a href="https://github.com/kubefirst/kubefirst/releases"><img title="Release date" src="https://img.shields.io/github/release-date/kubefirst/kubefirst"/></a>
</p>


---

# Kubefirst CLI

The Kubefirst CLI is a cloud provisioning tool. With simple setup and two CLI commands, we create a kubernetes cluster 
managed with automated Infrastructure as Code, GitOps asset management and application delivery, secrets management, a 
sample application delivered to development, staging, and production, and so much more. It's an open source platform 
ready to be customized to suit your company's needs.

- [Kubefirst CLI](#kubefirst-cli)
  - [DNS Setup](#dns-setup)
  - [Installing the CLI](#installing-the-cli)
  - [Other installation techniques:](#other-installation-techniques)
  - [Initialization](#initialization)
    - [localhost](#localhost)
    - [GitLab](#gitlab)
    - [GitHub](#github)
  - [Creation](#creation)
  - [Destroy](#destroy)
- [What to do next](#what-to-do-next)
- [If you want learn more](#if-you-want-learn-more)
  - [Access ArgoCD](#access-argocd)
  - [Available Commands](#available-commands)
  - [The Provisioning Process](#the-provisioning-process)
  - [Experimental Mode (we discourage it)](#experimental-mode-we-discourage-it)
  - [Feed K-Ray](#feed-k-ray)

![kubefirst architecture diagram](/images/kubefirst-arch.png)

## DNS Setup

In order to install Kubefirst it's required to have a public domain. For root domains, setting the `--hosted-zone-name`
is enough, in case you want to use subdomains, and the domain is hosted on AWS, please follow the
[AWS documentation](https://aws.amazon.com/premiumsupport/knowledge-center/create-subdomain-route-53/).

Provisioned services on root domain will be hosted as:

```bash
argocd.example.com
gitlab.example.com
...
```

Provisioned services on subdomains will be hosted as:

```bash
argocd.subdomain.example.com
gitlab.subdomain.example.com
...
```

## Installing the CLI

```bash 
brew install kubefirst/tools/kubefirst
```
## Other installation techniques:

[Details Here](./build/README.md)

## Initialization

Since Kubefirst 1.9 version, GitHub is also available as a Git platform provider alongside GitLab.

### localhost

localhost enables Kubefirst to be installed in your local machine, not requiring an AWS account, using localhost you can:

```bash
kubefirst local
```

### GitLab

To prepare the installation using GitLab you can:

```bash
kubefirst init \
--cloud aws \
--profile default \
--region eu-central-1 \
--admin-email user@example.com \
--cluster-name your_cluster_name \
--hosted-zone-name domain.example \
--s3-suffix you-s3-bucket-name \
--gitops-branch main \
--metaphor-branch main \
--git-provider gitlab \
--aws-nodes-spot
```

### GitHub

To prepare the installation using GitHub you can:

```bash
export KUBEFIRST_GITHUB_AUTH_TOKEN=your_github_auth_token

kubefirst init \
--admin-email yourname@example.com \
--cloud aws \
--hosted-zone-name example.com \
--region eu-central-1 \
--cluster-name example_com \
--profile default \
--github-user yourgithubhandle \
--github-owner yourgithuborganization \
--gitops-branch main \
--metaphor-branch main
```

## Creation

At this point, everything is ready to start provisioning the cloud services, and for that run:

```bash
kubefirst cluster create
```

## Destroy

It will destroy the kubefirst management cluster, and clean up every change made in the cloud.

```bash
kubefirst cluster destroy
```

# What to do next

[Learn More - Getting Started](https://docs.kubefirst.io/kubefirst/getting-started.html)


# If you want learn more 

## Access ArgoCD

```bash
aws eks update-kubeconfig --name your_cluster_name
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
kubectl -n argocd port-forward svc/argocd-server 8080:80
```

## Available Commands

Kubefirst provides extra tooling for handling the provisioning work.

| Command        | Description                                               |
|:---------------|:----------------------------------------------------------|
| clean          | removes all kubefirst resources locally for new execution |
| cluster create | create a kubefirst management cluster                     |
| cluster destroy| destroy the kubefirst management cluster                  |
| info           | provides general Kubefirst setup data                     |
| init           | initialize your local machine to execute `create`         |
| version        | print the version number for kubefirst-cli"               |

---

## The Provisioning Process

![kubefirst provisioning diagram](/images/provisioning.png)


## Experimental Mode (we discourage it)

If you want to run the main branch (`go run . init ...`, for example), you need to inform the version you are trying to use. To achieve this, you need to pass the `ldflags` in your command: `go run -ldflags="-X github.com/kubefirst/kubefirst/configs.K1Version=1.9.3" . version`. 

We use external repositories during the provisioning process, and the version/tag in these repositories matches with the version of kubefirst cli; due to this approach, you need to perform the steps described above to run the code from the source.

## Feed K-Ray

Did you know our superhero mascot K-Ray gets its frictionless superpowers from a healthy diet of GitHub stars? K-Ray gets soooo hungry too - you wouldn't believe it. Feed K-Ray a GitHub star ⭐ above to bookmark our project and keep K-Ray happy!!

[![Star History Chart](https://api.star-history.com/svg?repos=kubefirst/kubefirst&type=Date)](https://star-history.com/#kubefirst/kubefirst&Date)


todo: 
- standardize URL / Url
- move all config to config or to pkg but not both
- remove duplicate constants
- todo githubUser or githubUsername
- prereq aws-iam-authenticator
- docs updates
- specify ~/.k1 directory for ci
- github organization vs owner final
- discuss log levels
- KUBEFIRST_GITHUB_AUTH_TOKEN --> GITHUB_TOKEN