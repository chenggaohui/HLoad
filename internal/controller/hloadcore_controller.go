/*
Copyright 2024 cgh.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	hloadv1 "hload.com/m/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HLoadCoreReconciler reconciles a HLoadCore object
type HLoadCoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hload.hload.com,resources=hloadcores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hload.hload.com,resources=hloadcores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hload.hload.com,resources=hloadcores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HLoadCore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *HLoadCoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	println("收到k8s的Reconcile请求...........")
	// 获取控制器对象中的 template，例如 Deployment 的 Pod 模板
	hLoadCore := &hloadv1.HLoadCore{}
	err := r.Get(ctx, req.NamespacedName, hLoadCore)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Info("hLoadCore delete...")
			//TODO 删除资源逻辑
			return ctrl.Result{}, nil
		}
		return ctrl.Result{RequeueAfter: time.Minute}, err
	}
	podTemplate := hLoadCore.Spec.Template
	// 创建 Pod 对象
	pod := &v1.Pod{
		ObjectMeta: podTemplate.ObjectMeta,
		Spec:       podTemplate.Spec,
	}

	// 设置 OwnerReference，将 Pod 与控制器对象关联起来
	if err := controllerutil.SetControllerReference(hLoadCore, pod, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	//err = r.Create(ctx, pod)
	//if err != nil {
	//	klog.Error("create pod error:", err.Error())
	//	return ctrl.Result{RequeueAfter: time.Minute}, err
	//}
	//klog.Info("create pod success.")
	//hLoadCore.Status.Message = "reconcile success end time:" + time.Now().Format("2006-01-02 15:04:05")
	//err = r.Update(ctx, hLoadCore)
	//if err != nil {
	//	klog.Error("update hLoadCore status error")
	//	return ctrl.Result{RequeueAfter: time.Minute}, err
	//}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HLoadCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hloadv1.HLoadCore{}).
		Complete(r)
}
