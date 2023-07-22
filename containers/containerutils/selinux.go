package containerutils

import (
	"runtime"

	"github.com/opencontainers/selinux/go-selinux"
)

// Ensure SELinux allows us to perform privileged actions across cgroups if it is
// enabled. This is needed on my system to ensure we can start ryuk in testcontainers
// properly.
func VerifySelinuxEnforceMode() {
	if runtime.GOOS == "linux" && selinux.EnforceMode() == selinux.Enforcing {
		panic(
			"SELinux is enabled and set to 'Enforcing'. testcontainers-ryuk cannot start under this setting. " +
				"Please run 'sudo setenforce permissive' to lower the SELinux enforcement level to allow " +
				"Ryuk to start.",
		)
	}
}
