package controllers

import (
	"context"
	"fmt"
	"strings"

	metadata "github.com/abdulhaseeb2/node-label-controller/pkg"
	"github.com/go-logr/logr"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cri-api/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeReconciler reconciles a Node object
type NodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;patch
//+kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=nodes/finalizers,verbs=update

func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Node", req.NamespacedName)
	log.Info("Reconciling Node: " + req.Name)

	node := &corev1.Node{}

	err := r.Get(context.TODO(), req.NamespacedName, node)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Node no longer exists")

			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading the object - requeue the request.
		return reconcilerUtil.RequeueWithError(err)
	}

	// Node is marked for deletion
	if node.DeletionTimestamp != nil {
		log.Info("Deletion timestamp found for node: " + req.Name)
		return reconcilerUtil.DoNotRequeue()
	}

	// If status is not populated with OSImage name(Node might be initializing), then requeue after 2 mins
	if node.Status.NodeInfo.OSImage == "" {
		return reconcilerUtil.RequeueAfter(metadata.NodeRequeueTime)
	}

	// Check node OS image is flatcar
	if strings.HasPrefix(node.Status.NodeInfo.OSImage, metadata.FlatCarContainerLinuxOSImageName) {
		// Does label already exist
		if value, found := node.ObjectMeta.Labels[metadata.FlatCarContainerLinuxLabelKey]; !found || value != metadata.FlatCarContainerLinuxLabelValue {
			// Base object for patch, which patches using the merge-patch strategy with the given object as base.
			nodePatchBase := client.MergeFrom(node.DeepCopy())

			// Add/Update label
			node.ObjectMeta.Labels[metadata.FlatCarContainerLinuxLabelKey] = metadata.FlatCarContainerLinuxLabelValue

			// Patch node
			err = r.Patch(context.Background(), node, nodePatchBase)
			if err != nil {
				return reconcilerUtil.RequeueWithError(err)
			}

			log.Info(fmt.Sprintf("Node %s: Labels updated", node.Name))
			return reconcilerUtil.DoNotRequeue()
		}

		log.Info(fmt.Sprintf("Node %s: FlatCar Container Linux label already set", node.Name))
	}

	return reconcilerUtil.DoNotRequeue()
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		Complete(r)
}
