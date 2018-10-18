package main

import (
    "flag"
    "os"

    "github.com/golang/glog"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/apiserver/pkg/util/logs"

    basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
    "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"

    // make this the path to the provider that you just wrote
    yourprov "github.com/mbobran/test1/pkg/provider"
)

type YourAdapter struct {
    basecmd.AdapterBase

    // the message printed on startup
    Message string
}

func main() {
    logs.InitLogs()
    defer logs.FlushLogs()

    // initialize the flags, with one custom flag for the message
    cmd := &YourAdapter{}
    cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
    cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure you get the glog flags
    cmd.Flags().Parse(os.Args)

    provider := cmd.makeProviderOrDie()
    cmd.WithCustomMetrics(provider)
    // you could also set up external metrics support,
    // if your provider supported it:
    // cmd.WithExternalMetrics(provider)

    glog.Infof(cmd.Message)
    if err := cmd.Run(wait.NeverStop); err != nil {
        glog.Fatalf("unable to run custom metrics adapter: %v", err)
    }
}

func (a *SampleAdapter) makeProviderOrDie() provider.CustomMetricsProvider {
    client, err := a.DynamicClient()
    if err != nil {
        glog.Fatalf("unable to construct dynamic client: %v", err)
    }

    mapper, err := a.RESTMapper()
    if err != nil {
        glog.Fatalf("unable to construct discovery REST mapper: %v", err)
    }

    return yourprov.NewProvider(client, mapper)
}
