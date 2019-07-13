package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// "k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type APIResponse struct {
	Namespace   string
	Deployments []Deployment
}

type Deployment struct {
	Name     string
	Replicas int32
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// List Deployments
		list, err := deploymentsClient.List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		var APIResponse APIResponse
		APIResponse.Namespace = apiv1.NamespaceDefault

		var deployments []Deployment

		for _, d := range list.Items {
			deployments = append(APIResponse.Deployments, Deployment{Name: d.Name, Replicas: *d.Spec.Replicas})
		}

		APIResponse.Deployments = deployments

		res, err := json.Marshal(APIResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	fmt.Printf("Registering web server on port 8090...")

	server_err := http.ListenAndServe(":8090", nil)
	if server_err != nil {
		panic(server_err.Error())
	}

	fmt.Printf("Listening on port 8090...")
}
