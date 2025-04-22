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
name: myapp
entrypoint: python app.py

filesystem:
  base: ubuntu:22.04
  copy:
    - ./:/app
  workdir: /app

env:
  - DEBUG=true
  - PORT=8080

resources:
  cpu: "1"
  memory: "512M"

network:
  expose: [8080]

optimize:
  auto_memory_tune: true
  shrink_image: true
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


## Use Cases

Lit is ideal for:

- 🔹 Developers building **minimalist, high-performance containers**
- 🔹 Teams who want **Docker-like workflows without Docker**
- 🔹 Systems with limited resources (IoT, embedded, edge computing)
- 🔹 Security-focused deployments with rootless containers
- [ ] Cross-platform support (Linux-first, BSD later)
- [ ] Runtime plugin system for community features

---
