package containers

import (
	"runtime"

	"github.com/opencontainers/selinux/go-selinux"
)

func enforceSeLinuxPermissive() {
	if runtime.GOOS == "linux" && selinux.EnforceMode() == selinux.Enforcing {
		panic(
			"SELinux is enabled and set to 'Enforcing'. testcontainers-ryuk cannot start under this setting. " +
				"Please run 'sudo setenforce permissive' to lower the SELinux enforcement level to allow " +
				"Ryuk to start.",
		)
	}
}
