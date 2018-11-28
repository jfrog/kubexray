# kubexray

This integration enables run-time monitoring of pods on k8s via JFrog Xray.  

kubexray listens to event streams from both the k8s api server and JFrog Xray to ensure that only the pods that comply with certain policies can run on k8s. It understands the variation between different k8s resources (StatefulSets & Deployments) and accordingly different actions are applied (scale down to 0, delete, ignore). 

What this allows is following -
1. For every pod (running or scheduled to run) on k8s, kubexray checks if there are vulnerabilities or license issues. If yes, then actions (scaledown, delete, ignore) can be easily enforced. 
2. For every pod (running or scheduled to run) on k8s, kubexray checks if the corresponding docker images are scanned by Xray (i.e. coming from Artifactory or not). If not, the same actions (scaledown, delete, ignore) can be applied. 
3. Any time a new policy gets added or updated on Xray, or a new vulnerability is reported, kubexray detects this change and checks if there are issues with existing pods. Same actions (scaledown, delete, ignore) can be applied.


## Install instructions

## How to enforce policies?
