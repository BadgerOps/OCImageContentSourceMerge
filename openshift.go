package main

import (
    "flag"
    "fmt"
    "os"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/util/retry"
    operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
    "github.com/openshift/client-go/operator/clientset/versioned"
)

/*
* So, playing around with this for now, might change it. Create twomethods of interacting, either .kube/config
* or `--token= and `--server=`
*
*/

func createClientsetWithKubeconfig() (*versioned.Clientset, error) {
    var kubeconfig *string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // Use the current context in kubeconfig TODO:do we allow ourselves to set context for multi context kubeconfig?
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        return nil, fmt.Errorf("error in building kubeconfig: %v", err)
    }

    // Create the clientset
    clientset, err := versioned.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("error in creating clientset: %v", err)
    }

    return clientset, nil
}

// Function to create a clientset to interact with OpenShift API using token
func createClientsetWithToken(serverURL, token string) (*versioned.Clientset, error) {
    config := &rest.Config{
        Host: "serverURL",
        BearerToken: token,
        TLSClientConfig: rest.TLSClientConfig{Insecure: true}, // or configure proper TLS
    }

    clientset, err := versioned.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("error in creating clientset with token: %v", err)
    }

    return clientset, nil
}

// Function to get existing ImageContentSourcePolicies
func getExistingPolicies(clientset *versioned.Clientset) (*operatorv1alpha1.ImageContentSourcePolicyList, error) {
    policies, err := clientset.OperatorV1alpha1().ImageContentSourcePolicies().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("error retrieving policies: %v", err)
    }

    return policies, nil
}
