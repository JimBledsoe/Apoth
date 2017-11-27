# Apoth - a kubeconfig management utility

(Apoth is Greek for storage bin - aka container)

## Background

When using Kubernetes where you have access to multiple clusters and need to quickly move between them, each cluster has a unique context and method of accessing the cluster.  The official documentation seems to steer users toward two possible solutions:

* Combining multiple contexts and clusters in the same configuration file instead of keeping each cluster configuration separate.  This seems counter to containerization.
* Using the KUBECONFIG environment variable to chain together multiple files (much as PATH does)

For now, the east way to keep the config files separate is to use the environment variable KUBECONFIG to point to the desired config file, or a chain of files.  This is not friendly.  Aside from editing users configuration scripts, the only easy way to modify session environmental variables is to assign to expressions, which is a little clunky.

My proposal for handling multiple config files is to look for a directory ~/.kube/config.d and if found, treat all of the files in that directory as independent config files.  This would operate identically to the chaining of files in the KUBECONFIG environmental variable together, but allow an easy mechanism for "installation" to copy or link files into the users ~/.kube/config.d directory for semi-permanent persistence.

## How it works

This project is mainly just an exercise at writing a useful Go program that will:

* Make sure the KUBECONFIG env var is not being used to select a specific config file chain
  (Needs further research to determine if simultaneous env var path chain and config file can coexist.  The documentation suggests the ~/.kube/config file is only used if the env var is not set)
  <https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/>
* Locate multiple config files (~/.kube/config.d)
  * Split out a multi-context file into individual files (maybe - still TBD)
  * Look for multiple contexts in each file and present each context as a separate entry in the menu
* Present the list of contexts to the user in a simple menu
* Let the user select one of the configurations (contexts)
* Rewrite the ~/.kube/config file from the selected context

## Future Intentions

If it works out well, I may consider trying to submit this to the Golang group for inclusion into kubectl.  Keeping that in mind, I will attempt to keep everything sympatico with how kubectl currently handles configuration files so that it will be a trivial task to rip out the pertinent code later.  At first glance, it would seem to be simple enough just to include the config.d discovered files as if they were appended to the KUBECONFIG environment variable.
