# Apoth - a kubeconfig management utility
(Apoth is Greek for storage bin - aka container) 

## Background
When using Kubernetes where you have access to multiple clusters and need to quickly move between them, each cluster has a unique context and method of accessing the cluster.  The official documentation seems to steer users toward combining multiple contexts and clusters in the same configuration file instead of keeping each cluster configuration separate.  This seems counter to containerization.  

For now, the east way to keep the config files separate is to use the environment variable KUBECONFIG to point to the desired config file.  This is not friendly.

## Proposal
This project is mainly just an exercise at writing a useful Go program that will:
* Locate multiple config files (~/.kube/config.d)
  * Look for multiple contexts in each file and present each context as a separate entry in the menu
  * Split out a multi-context file into individual files (TBD).
* Present the list of contexts to the user in a simple menu
* Let the user select one of the configurations (contexts) 
* Rewrite the ~/.kube/config file from the selected context

## Future Intentions
If it works out well, I may consider trying to submit this to the Golang group for inclusion into kubectl.  Keeping that in mind, I will attempt to keep everything sympatico with how kubectl currently handles configuration files so that it will be a trivial task to rip out the pertinent code later.

My proposal for handling multiple config files is to look for a directory ~/.kube/config.d and if found, treat all of the files in that directory as independent config files.