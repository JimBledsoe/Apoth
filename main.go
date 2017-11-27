// A simple utility to look for multiple KUBECONFIG files, look for multiple constexts in thsoe files,
// then present a simple menu listing to the user to select one.
// Upon selection, the ~/.kube/config is overwritten and the program exits.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	// TODO "version"
	// TODO move kubectl applicable funcitons to a separate module
)

const (
	program    = "Apoth"
	kubeconfig = "KUBECONFIG"
	dotkube    = ".kube"
	configfile = "config"
	configdir  = "config.d"
	version    = "NotSetYet" // TODO: move to version file
)

// k8context structure holds the context label and the path where the context is defined
type k8context struct {
	label    string
	filepath string
}

// userHomeDir uses the user library to find the users home directory on
// either Linux or Windows and returns the path
func userHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, err
}

// BuildContexts searches for kubeconfig files from a given directory
// and returns a list of contexts found inside those files
func buildContexts(path string) ([]k8context, error) {
	var k8list []k8context
	var k8c k8context
	// TODO these are hard coded for now
	k8c.label, k8c.filepath = "first-context", "/tmp/foo"
	k8list = append(k8list, k8c)
	k8c.label, k8c.filepath = "second-context", "/tmp/bar"
	k8list = append(k8list, k8c)
	k8c.label, k8c.filepath = "third-context", "/tmp/baz"
	k8list = append(k8list, k8c)
	return k8list, nil
}

// showContexts takes an array of context structures and prints a pleasantly
// formatted list of the index, label, and path of the contexts
func showContexts(k8l []k8context) {
	// Find out how long certain things are just to make a prettier output
	var longestLabel int
	var longestPath int
	for _, k8c := range k8l {
		if len(k8c.label) > longestLabel {
			longestLabel = len(k8c.label)
		}
		if len(k8c.filepath) > longestPath {
			longestPath = len(k8c.filepath)
		}
	}
	listLen := len(k8l)
	lenWidth := 1 + int(math.Log10(float64(listLen)))
	fullWidth := lenWidth + 1 + 2 + longestLabel + 2 + longestPath

	// Display the discovered contexts
	fmt.Printf("%s\n", strings.Repeat("-", fullWidth))
	for i := 0; i < listLen; i++ {
		fmt.Printf("%*d)  %-*s  %s\n",
			lenWidth, i+1,
			longestLabel, k8l[i].label,
			k8l[i].filepath)
	}
	fmt.Printf("%s\n", strings.Repeat("-", fullWidth))
}

// selectContext promts the user to select one of the displayed contexts to use and
// returns a pointer to the selected context structure
func selectContext(k8l []k8context) *k8context {
	if len(k8l) == 0 {
		return nil
	}
	// Keep asking for input until they enter something valid and return
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("Enter 1-%d to set a context, 0 or <enter> to abort: ", len(k8l))
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Print(err)
		} else {
			text = strings.Replace(text, "\r\n", "", -1) // Windows EOL
			text = strings.Replace(text, "\n", "", -1)   // Unix EOL
			// Empty input is the fast ticket out
			if text == "" || text == "0" {
				return nil
			}
			// Validate the input
			num, err := strconv.Atoi(text)
			if err != nil {
				// Silently fail garbage input and stay in loop
			} else if num >= 1 && num <= len(k8l) {
				return &k8l[num-1] // Valid input
			}
		}
	}
	return nil
}

// setContext will overwrite the users .kube/config with the given context
func setContext(k8p *k8context) {
	if k8p == nil {
		return
	}
	fmt.Println("TODO: Set kubeconfig to", k8p.filepath)
	// TODO: write the context info to config file
	return
}

// main loop
func main() {
	fmt.Printf("\n%s %s - a Kubernetes context selector\n", program, version)

	// If the ENV var is being used, it overrides any config file, so we can't use this utility
	envKubeConfig := os.Getenv(kubeconfig)
	if envKubeConfig != "" {
		log.Fatalf("%s: Setting $%s overrides looking for configs in ~/%s", program, kubeconfig, dotkube)
	}

	// Ensure the directory tree down to the config directory exists and bail if it doesn't
	home, err := userHomeDir()
	if err != nil {
		log.Fatalf("%s: locating user's home directory: %s", program, err)
	}
	configDir := path.Join(home, dotkube, configdir)
	if _, err = os.Stat(configDir); err != nil {
		log.Fatalf("%s: config directory does not exist: %s", program, err)
	}

	// If there is an existing config file, see if its in the config directory and put there if not
	// TODO

	// If there are not 2 or more files, then this utility is meaningless
	// TODO

	// Build contexts from config files
	k8contexts, err := buildContexts(configDir)

	// Present menu to user
	showContexts(k8contexts)

	// Accept input from user
	k8p := selectContext(k8contexts)

	// Set context if one was selected
	if k8p != nil {
		setContext(k8p)
	}

	os.Exit(0)
}
