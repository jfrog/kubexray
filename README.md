# kubexray

An open source software project that monitors pods in a Kubernetes cluster to help you detect security & license violations in containers 
running inside the pod. 

KubeXray listens to events from k8s api server, and leverages the metadata from JFrog Xray (commercial product) to ensure that only the pods that comply with your current policy can run on k8s. As an example, KubeXray listens to these event streams:
* Deployment of a new service
* Upgrade of an existing service
* A new license policy, such as a new license type disallowed for runtime.
* A new security issue

And when an issue is detected, KubeXray responds according to the current policy that you have set. 

You can select one of the following possible actions:
* Scaledown to 0. The desired state of a service is updated to 0, making the services inactive but still traceable.
* Delete the corresponding Kubernetes resource that’s pointing to a vulnerable container image(s)
* Ignore and leave the pod running

KubeXray also allows you to enforce policy for running applications that have not been scanned by JFrog Xray and whose risks are unknow. 


## Build Instructions

## Install Instructions
