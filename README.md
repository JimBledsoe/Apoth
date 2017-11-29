# Apoth - a kubeconfig management utility

(Apoth is Greek for storage bin - aka container)

## Background

When using Kubernetes where you have access to multiple clusters and need to quickly move between them, each cluster has a unique context and method of accessing the cluster.  The official documentation steers users toward two possible solutions:

* Combining multiple contexts and clusters in the same configuration file instead of keeping each cluster configuration separate.  This seems counter to containerization.
* Using the KUBECONFIG environment variable to chain together multiple files (much as PATH does)

For now, the east way to keep the config files separate is to use the environment variable KUBECONFIG to point to the desired config file, or a chain of files.  This is not super friendly.  Aside from editing users configuration scripts, the only easy way to modify session environmental variables is to assign to expressions, which is a little clunky.

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

If it works out well, I may consider trying to submit this to the K8s group for inclusion into kubectl.  Keeping that in mind, I will attempt to keep everything sympatico with how kubectl currently handles configuration files so that it will be a trivial task to rip out the pertinent code later.  At first glance, it would seem to be simple enough just to include the config.d discovered files as if they were appended to the KUBECONFIG environment variable.

## Research into Existing kubectl Functionality

As mentioned, the documentation steers users to either creating combo configs, or using
the environment variable KUBECONFIG to chain together multiple files.  Looking at the
source code does not contradict this.  Furthermore, from the code or testing:
* If --kubeconfig is specified on the command line, that single file is used and no other locations are checked.
* If the $KUBECONFIG is set, the ~.kube/config file is not checked.
* If the $KUBECONFIG is set to directories instead of files, it will panic.
* A deprecated ~.kube/.kubeconfig file is checked and migrated if found.
* When use-context is used, the very first config file in the $KUBECONFIG chain is the one altered to hold the current-context attribute.  It could easily be setting a context defined in a different file.
* When attributes are set for an entity already defined in a particular config file, the attributes are correctly placed on the attribute in the original file (not use-context).
* When entities are added (set-cluster) they are always stored in the first file in the chain.
* Manual edits to the config files (like setting the current-context) are properly respected on subsequent calls to kubectl.
