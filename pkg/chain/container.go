package chain

import (
	"fmt"

	"github.com/InfraBlockchain/infrablockspace-operator/pkg/render"
	corev1 "k8s.io/api/core/v1"
)

func CreateInitContainer(name, image string, commands []string, volumeMounts []corev1.VolumeMount) corev1.Container {
	return corev1.Container{
		Name:         name,
		Image:        image,
		Command:      commands,
		VolumeMounts: volumeMounts,
	}
}

func CreateChainContainer(name, imageVersion string, commands, args []string, volumeMounts []corev1.VolumeMount) corev1.Container {
	return corev1.Container{
		Name:         name,
		Image:        fmt.Sprintf("public.ecr.aws/v8x3j0k5/infrablockspace:%s", imageVersion),
		Args:         args,
		Command:      commands,
		VolumeMounts: volumeMounts,
	}
}

func GetDownloadSpecCommand(chainSpecUrl, fileName string) []string {
	return []string{
		"curl",
		"-L",
		chainSpecUrl,
		"-o",
		fmt.Sprintf("/tmp/%s", fileName),
	}
}

func GetInjectKeyCommandAndArgs(keys []Key) ([]string, []string) {
	commands := []string{"/bin/sh"}
	injectKeyStoreScript := render.RenderingInTemplate(InjectKeyScript, keys)
	args := []string{
		"-c",
		injectKeyStoreScript,
	}
	return commands, args
}

func GetRelayChainArgs(port Port, isBoot bool, bootNodesUrl []string) []string {
	validatorArgs := getChainArgs(RelayChain)
	rpc := fmt.Sprintf("--rpc-port=%d", port.RPCPort)
	validatorArgs = appendRelayChainArgs(validatorArgs, rpc)
	if !(isBoot) {
		validatorArgs = appendBootNods(validatorArgs, bootNodesUrl...)
	}
	return validatorArgs

}

func appendRelayChainArgs(args []string, ports ...string) []string {
	relayArgs := []string{
		"--base-path",
		"/data/relay",
		"--chain",
		"/tmp/relay-chain-spec.json",
		"--prometheus-external",
		"--prometheus-port=9615",
		"--unsafe-rpc-external",
		"--rpc-cors",
		"all",
		"--keystore-path=/keystore",
		"--rpc-max-connections=16000",
	}
	args = append(args, relayArgs...)
	args = append(args, ports...)
	return args
}

func appendBootNods(args []string, bootNodesUrl ...string) []string {
	args = append(args, "--bootnodes")
	args = append(args, bootNodesUrl...)
	return args
}

func getChainArgs(chainPosition ChainType) []string {
	if chainPosition == "relay" {
		return []string{
			"--validator",
		}
	} else {
		return []string{
			"--collator",
		}
	}
}
