# Lit

**Lit** runs your applications in isolated environments using native Linux kernel features like namespaces, cgroups, and union filesystems — just like Docker, but from scratch and in an minimal way.

---

## Key Features

- ** Docker-like containerization without Docker**  
  Uses pure Linux kernel primitives (no Docker daemon, no containerd, no OCI dependency).

- ** Single config file (lit.yml)**  
  Configure everything (filesystem, resources, env vars, networking) in one declarative file.

- ** Custom networking**  
  Built-in support for isolated networks, port forwarding, and service linking — no `docker-compose` needed.

- ** Rootless & secure by default**  
  Uses Linux user namespaces for rootless containers. Sandboxing options with seccomp + AppArmor/SELinux profiles.

---

## How It Works

Lit creates containers by stitching together core Linux features:

| Component | Linux Feature Used |
|----------|----------------------|
| Process isolation | PID namespace |
| Filesystem isolation | Mount namespace + chroot |
| Network isolation | Network namespace + veth pairs |
| User isolation | User namespaces |
| Resource limits | cgroups v2 |
| Filesystem layering | OverlayFS |
| Runtime optimizations | Custom analyzers and hooks |

---

## Example: `lit.yml`

```yaml
# ===============================
# LIT CONTAINER CONFIGURATION
# ===============================

# Namespace types to isolate for this container.
# Supported types: "net", "pid", "mount"
# These control process/network/filesystem isolation via Linux namespaces.
# Example: isolate network and process IDs
namespace_type: "pid,mount"

# -------------------------------
# RESOURCE LIMITS
# -------------------------------

# Maximum memory the container is allowed to use (in bytes).
# Example: 104857600 = 100 MB
memory_limit: 104857600

# CPU limit in microseconds per 100,000us (100ms) period.
# Example: 50000 = 50% CPU, 20000 = 20% CPU
# Set to -1 for unlimited CPU
cpu_limit: 50000

# -------------------------------
# CONTAINERIZED APPLICATION
# -------------------------------

# The name of the executable inside the container to run.
# Your container root filesystem should contain this binary under: /bin/<image>
#
# Example:
#   If this value is "myserver", then:
#   → Place your binary at: base/bin/myserver
#
# Can be named anything (including "app", "nginx", "hello.exe", etc.)
# Note: this must be a Linux-compatible statically compiled binary.
image: "testApp"
```

## Why Go is Ideal for a Container Runtime

Go is purpose-built for systems like Lit — here’s why it’s a perfect match:

| Reason               | Why It Matters                                                                 |
|----------------------|--------------------------------------------------------------------------------|
| Low-level access   | Interact with Linux syscalls (`unshare`, `clone`, `mount`) via `syscall` or `golang.org/x/sys/unix`. |
| Concurrency model | Goroutines make process lifecycles, I/O, and async optimization smooth.        |
| Static binaries    | Compile to a single, dependency-free binary — portable and easy to distribute.|
| Fast build & deploy| Rapid iteration with tiny binaries. Great for fast dev loops and CI/CD.       |
| Mature ecosystem   | Rich libraries for YAML parsing, CLI, cgroups, netlink, and more.              |
| Community support  | Most container tooling is in Go — easier to learn from and contribute to.     |

---
