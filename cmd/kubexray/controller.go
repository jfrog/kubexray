package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	logger    *log.Entry
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	handler   Handler
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// handle a panic with logging and exiting
	defer utilruntime.HandleCrash()
	// ignore new items in the queue but when all goroutines
	// have completed existing items then shutdown
	defer c.queue.ShutDown()

	c.logger.Debug("Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.informer.Run(stopCh)

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.logger.Debug("Controller.Run: cache sync complete")

	// run the runWorker method every second with a stop channel
	wait.Until(c.runWorker, time.Second, stopCh)
}

// execute the loop to process new items added to the queue
func (c *Controller) runWorker() {
	c.logger.Debug("Controller.runWorker: starting")

	// invoke processNextQueueItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextQueueItem() {
		c.logger.Debug("Controller.runWorker: processing next item")
	}

	c.logger.Debug("Controller.runWorker: completed")
}

// processNextQueueItem retrieves each queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextQueueItem() bool {
	c.logger.Debug("Controller.processNextQueueItem: start")

	// fetch the next item (blocking) from the queue to process or
	// if a shutdown is requested then return out of this to stop
	// processing
	item, quit := c.queue.Get()
	c.logger.Debug("Controller.processNextQueueItem: item fetched")

	// stop the worker loop from running as this indicates we
	// have sent a shutdown message that the queue has indicated
	// from the Get method
	if quit {
		return false
	}

	defer c.queue.Done(item)

	// assert the string out of the item (format `namespace/name`)
	indexKey, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		c.queue.Forget(item)
		utilruntime.HandleError(err)
		return true
	}


	// take the string item and get the object out of the indexer
	//
	// item will contain the complex object for the resource and
	// exists is a bool that'll indicate whether or not the
	// resource was created (true) or deleted (false)
	//
	// if there is an error in getting the item from the index
	// then we want to retry this particular queue item a certain
	// number of times (5 here) before we forget the queue item
	// and throw an error
	_, exists, err := c.informer.GetIndexer().GetByKey(indexKey)
	if err != nil {
		if c.queue.NumRequeues(item) < 5 {
			c.logger.Errorf("Controller.processNextQueueItem: Failed processing item with item %s with error %v, retrying", item, err)
			c.queue.AddRateLimited(item)
		} else {
			c.logger.Errorf("Controller.processNextQueueItem: Failed processing item with item %s with error %v, no more retries", item, err)
			c.queue.Forget(item)
			utilruntime.HandleError(err)
		}
	}

	// if the item doesn't exist then it was deleted and we need to fire off the handler's
	// ObjectDeleted method. but if the object does exist that indicates that the object
	// was created (or updated) so run the ObjectCreated method
	//
	// after both instances, we want to forget the item from the queue, as this indicates
	// a code path of successful queue item processing
	if !exists {
		c.logger.Debugf("Controller.processNextQueueItem: object deleted detected: %s", indexKey)
		c.handler.ObjectDeleted(c.clientset, item)
		c.queue.Forget(item)
	} else {
		c.logger.Debugf("Controller.processNextQueueItem: object created detected: %s", indexKey)
		c.handler.ObjectCreated(c.clientset, item)
		c.queue.Forget(item)
	}

	// keep the worker loop running by returning true
	return true
}
