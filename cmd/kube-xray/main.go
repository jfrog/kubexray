package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_v1 "k8s.io/client-go/informers/core/v1"
	api_core_v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/rest"

	// Import to initialize client auth plugins.
	// Only GCP GKE auth is supported, Azure auth crashes with go-client v9.0.0
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// retrieve the k8s cluster client from within/outside of the cluster
func getKubernetesClient() kubernetes.Interface {

	var client kubernetes.Interface

	clusterConfig, err := rest.InClusterConfig()
	if err == nil {
		client, err := kubernetes.NewForConfig(clusterConfig)
		if err == nil {
			return client
		}
	}

	if err != nil {
		log.Warnf("Get in-cluster config: %v", err)
	}

	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
		// generate the client based off of the config
	}

	// generate the client based off of the config
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Debug("Successfully constructed k8s client")
	return client
}

// main code path
func main() {
	client := getKubernetesClient()

	namespace := os.Getenv("JFROG_K8S_NS")
	if namespace == "" {
		namespace = meta_v1.NamespaceDefault
	}
	namespace = ""

	//Create the filtered informer
	//See: cache.NewFilteredListWatchFromClient
	informer := core_v1.NewFilteredPodInformer(client, namespace,
		//how often to call the update function on the informer with all the objects
		//in the cache in order to retransmit events - no resync (0)
		//TODO: to report to artifactory, we want to use resync
		0,
		cache.Indexers{},
		//TODO: customize a returned ListOptions to watch on, such as pod annotations (using envs).
		func(options *meta_v1.ListOptions) {},
	)

	// create a new queue for the informer to put watched resources as keys for the handler to take
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	//Set event handlers for the 3 event types by adding the resource key to the queue
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// convert the resource object into a key (in the format of 'namespace/name') - just for debugging
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Debugf("Add pod: %s", key)
			if err == nil {
				enqueuePod(obj, queue, true)
			}
		},
		//Only called on re-sync
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			log.Debugf("Update pod: %s", key)
			if err == nil {
				store := informer.GetStore()
				store = store
				enqueuePod(newObj, queue, true)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// DeletionHandlingMetaNamsespaceKeyFunc is a helper function that allows
			// us to check the DeletedFinalStateUnknown existence, in the event that
			// a resource was deleted but it is still contained in the index
			//
			// this then in turn calls MetaNamespaceKeyFunc
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Debugf("Delete pod: %s", key)
			if err == nil {
				enqueuePod(obj, queue, false)
			}
		},
	})

	handler := &TestHandler{}
	if handler.Init(client) != nil {
		os.Exit(1)
	}

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := Controller{
		logger:    log.NewEntry(log.New()),
		clientset: client,
		informer:  informer,
		queue:     queue,
		handler:   handler,
	}

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}

func enqueuePod(obj interface{}, queue workqueue.RateLimitingInterface, includeOnlyRunning bool) bool {
	//We copy the object to the cache rather than store it by a string key and getting
	//it back from the index by this key, since when we get it by key it is already deleted
	//and does not contain the object in the index (from which we need to extract the containers)
	pod := obj.(*api_core_v1.Pod)
	//Filter pending pods
	if includeOnlyRunning && pod.Status.Phase != api_core_v1.PodRunning {
		return false
	}
	copy := pod.DeepCopy()
	queue.Add(copy)
	return true
}
