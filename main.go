/*
Copyright 2021.

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

package main

import (
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	appsv1 "github.com/anilkusc/ty-case/api/v1"
	"github.com/anilkusc/ty-case/controllers"
	"log"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(appsv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	log.Print("Creating the Manager.")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatal(err, "Unable to start manager")
		os.Exit(1)
	}

	log.Print("Starting the Controller.")
	if err = (&controllers.TestAppReconciler{
		Client: mgr.GetClient(),
		//Log:    ctrl.Log.WithName("controllers").WithName("TestApp"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create controller", "controller", "TestApp")
		os.Exit(1)
	}

	log.Print("Starting the Manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal(err, "problem running manager")
		os.Exit(1)
	}
}
