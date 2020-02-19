# Deprecation Notice
*Note: KubeXray is no longer maintained or supported by JFrog.  Feel free to review this code for your own POC concepts, but we are not continuing to update it or add features.  For people looking for great tools to help for enforcement in Kubernetes, we do continue to have [KubeNab](https://github.com/jfrog/kubenab) which allows enforcement of what repositories a kubernetes cluster pulls from (which then can leverage enforcement of Xray policies in Artifactory).*

# JFrog KubeXray scanner on Kubernetes

An open source software project that monitors pods in a [Kubernetes](https://kubernetes.io/) cluster to help you detect security & license violations in containers 
running inside the pod. 

KubeXray listens to events from Kubernetes API server, and leverages the metadata from [JFrog Xray](https://jfrog.com/xray/) (commercial product) to ensure that only the pods that comply with your current policy can run on Kubernetes. As an example, KubeXray listens to these event streams:
* Deployment of a new service
* Upgrade of an existing service
* A new license policy, such as a new license type disallowed for runtime.
* A new security issue

And when an issue is detected, KubeXray responds according to the current policy that you have set. 

You can select one of the following possible actions:
* Scaledown to 0. The desired state of a service's replica count is updated to 0, making the services inactive but still traceable.
* Delete the corresponding Kubernetes resource that’s pointing to a vulnerable container image(s)
* Ignore and leave the pod running

KubeXray also allows you to enforce policy for running applications that have not been scanned by JFrog Xray and whose risks are unknown. 

## Install Instructions

The easiest way to install KubeXray is using the Helm [chart](https://github.com/jfrog/charts/tree/master/stable/kubexray)

Please follow install instruction from chart's [readme](https://github.com/jfrog/charts/blob/master/stable/kubexray/README.md)

## Local development and testing

### Building binary

To build `kubexray` locally 

  ```console
  make build
  ```

### Docker

To build `kubexray` docker image locally (testing docker image build)

  ```console
  make image
  ```

## Contributing Code
We welcome community contribution through pull requests.

<a name="License"/>

## License
This tool is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).


(c) All rights reserved JFrog
