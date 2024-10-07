package main

import (
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/version"
	"os/exec"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

type RedirectConf struct {
	types.NetConf
	OriginalPort string `json:"originalPort"`
	RedirectPort string `json:"redirectPort"`
}

func cmdAdd(args *skel.CmdArgs) error {
	var conf RedirectConf
	if err := json.Unmarshal(args.StdinData, &conf); err != nil {
		return fmt.Errorf("failed to parse CNI config: %v", err)
	}

	var result types.Result

	// Set up iptables rule to redirect traffic
	fmt.Println("In cmdAdd 0")
	err := setupIptables(conf.OriginalPort, conf.RedirectPort)
	fmt.Println("In cmdAdd 1")
	if err != nil {
		fmt.Println("In cmdAdd 2")
		return fmt.Errorf("failed to set up iptables: %v", err)
	}
	fmt.Println("In cmdAdd 3")

	return types.PrintResult(result, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	var conf RedirectConf
	if err := json.Unmarshal(args.StdinData, &conf); err != nil {
		return fmt.Errorf("failed to parse CNI config: %v", err)
	}

	fmt.Println("In cmdDel 1")
	// Clean up iptables rule
	err := teardownIptables(conf.OriginalPort, conf.RedirectPort)
	fmt.Println("In cmdDel 2")
	if err != nil {
		return fmt.Errorf("failed to clean up iptables: %v", err)
	}
	fmt.Println("In cmdDel 3")

	return nil
}

func setupIptables(originalPort, redirectPort string) error {
	fmt.Println("In setupIptables 1")
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING",
		"-p", "tcp", "--dport", originalPort,
		"-j", "REDIRECT", "--to-port", redirectPort)
	fmt.Println("In setupIptables 2")
	return cmd.Run()
}

func teardownIptables(originalPort, redirectPort string) error {
	fmt.Println("In teardownIptables 1")
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING",
		"-p", "tcp", "--dport", originalPort,
		"-j", "REDIRECT", "--to-port", redirectPort)
	fmt.Println("In teardownIptables 2")
	return cmd.Run()
}

func cmdCheck(args *skel.CmdArgs) error {
	fmt.Println("In cmdCheck 1")
	return nil // No-op for simplicity
}

func main() {
	fmt.Println("In main 1")
	skel.PluginMainFuncs(skel.CNIFuncs{
		Add:   cmdAdd,
		Check: cmdCheck,
		Del:   cmdDel,
		/* FIXME GC */
		/* FIXME Status */
	}, version.All, bv.BuildString("redirect-cni"))
	fmt.Println("In main 2")
}
