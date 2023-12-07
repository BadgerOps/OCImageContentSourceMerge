package main

import (
    "flag"
    "fmt"
    "os"
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
        // Token is provided; use token-based authentication
        if *serverFlag == "" {
            fmt.Println("Error: --server flag must be set when using --token")
            os.Exit(1)
        }
        clientset, err = createClientsetWithToken(*serverFlag, *tokenFlag)
    } else {
        // No token provided; use kubeconfig for authentication
        clientset, err = createClientsetWithKubeconfig()
    }

    if err != nil {
        fmt.Printf("Error creating OpenShift client: %v\n", err)
        os.Exit(1)
    }

    // Rest of your code to handle ImageContentSourcePolicies...
}
