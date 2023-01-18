package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/ronenlib/kube-event-alert/config"
	"github.com/ronenlib/kube-event-alert/pkg/controller"
	"github.com/ronenlib/kube-event-alert/pkg/notifier"
	"github.com/ronenlib/kube-event-alert/pkg/util"
	"k8s.io/klog"
)

//go:embed resources/*
var tplFS embed.FS

func main() {
	klog.InitFlags(nil)

	config := config.Load()

	clientset := util.GetKubeClient(config.MasterURL, config.Kubeconfig)
	notifier := notifier.NewWebhookNotifier(config.WebhookURL, tplFS)

	stopCh := make(chan struct{})
	setInterrupt(stopCh)

	controller.Launch(clientset, notifier, stopCh)
}

func setInterrupt(stopCh chan struct{}) {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		close(stopCh)

		<-interrupt
		os.Exit(1)
	}()
}
