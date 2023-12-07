package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "github.com/openshift/client-go/operator/clientset/versioned"
    "sigs.k8s.io/yaml"
)

func main() {
    var (
        tokenFlag    = flag.String("token", "", "API token for authentication")
        serverFlag   = flag.String("server", "", "OpenShift API server URL")
    )
    flag.Parse()

    var clientset *versioned.Clientset
    var err error

    if *tokenFlag != "" {
        if *serverFlag == "" {
            fmt.Println("Error: --server flag must be set when using --token")
            os.Exit(1)
        }
        clientset, err = createClientsetWithToken(*serverFlag, *tokenFlag)
    } else {
        clientset, err = createClientsetWithKubeconfig()
    }

    if err != nil {
        fmt.Printf("Error creating OpenShift client: %v\n", err)
        os.Exit(1)
    }

    // Fetch existing policies
    policies, err := getExistingPolicies(clientset)
    if err != nil {
        fmt.Printf("Error retrieving policies: %v\n", err)
        os.Exit(1)
    }

    // Write each policy to a separate file
    for _, policy := range policies.Items {
        fileName := filepath.Join("./src", fmt.Sprintf("%s.yaml", policy.Name))
        policyYAML, err := yaml.Marshal(policy)
        if err != nil {
            fmt.Printf("Error marshalling policy %s: %v\n", policy.Name, err)
            continue
        }
        err = ioutil.WriteFile(fileName, policyYAML, 0644)
        if err != nil {
            fmt.Printf("Error writing policy %s to file: %v\n", policy.Name, err)
            continue
        }
        fmt.Printf("Policy %s successfully written to %s\n", policy.Name, fileName)
    }
}
