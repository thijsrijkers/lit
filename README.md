# Lit

**Lit** is a container runtime that brings performance, simplicity, and transparency to modern app deployment — without relying on Docker or any external container engines.

Lit runs your applications in isolated environments using native Linux kernel features like namespaces, cgroups, and union filesystems — just like Docker, but from scratch and in a fully customizable, minimal way.

---

## 🚀 Key Features

- **🔧 Docker-like containerization without Docker**  
  Uses pure Linux kernel primitives (no Docker daemon, no containerd, no OCI dependency).

- **📁 Single config file (lit.yml)**  
  Configure everything (filesystem, resources, env vars, networking) in one declarative file.

- **♻️ Auto-Optimizing Containers (Optional)**  
  Enable runtime optimizations:
  - **Memory tuning** based on usage.
  - **Dead code stripping** for interpreted languages (like Python/Node).
  - **Layer pruning & image slimming** in the background.

- **🌐 Custom networking**  
  Built-in support for isolated networks, port forwarding, and service linking — no `docker-compose` needed.

- **🔒 Rootless & secure by default**  
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

## 📄 Example: `lit.yml`

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

## ✅ Why Go is Ideal for a Container Runtime

Go is purpose-built for systems like Slimbox — here’s why it’s a perfect match:

| Reason               | Why It Matters                                                                 |
|----------------------|--------------------------------------------------------------------------------|
| 🧠 Low-level access   | Interact with Linux syscalls (`unshare`, `clone`, `mount`) via `syscall` or `golang.org/x/sys/unix`. |
| ⚙️ Concurrency model | Goroutines make process lifecycles, I/O, and async optimization smooth.        |
| 💼 Static binaries    | Compile to a single, dependency-free binary — portable and easy to distribute.|
| 🌱 Fast build & deploy| Rapid iteration with tiny binaries. Great for fast dev loops and CI/CD.       |
| 🛠 Mature ecosystem   | Rich libraries for YAML parsing, CLI, cgroups, netlink, and more.              |
| 👥 Community support  | Most container tooling is in Go — easier to learn from and contribute to.     |

---

## 🧠 Optimizer: Behind the Scenes

Lit goes beyond traditional containers by offering **on-the-fly optimization** to improve performance, reduce size, and clean up unnecessary overhead — optionally enabled via the config.

Here’s how each optimization works:

### 🔍 Memory Tuning (`auto_memory_tune`)
At runtime, Slimbox monitors container memory usage and automatically adjusts `cgroup` limits. It starts with a safe allocation and can scale limits up/down based on behavior, preventing out-of-memory crashes or unused reservations.

- Detects idle vs active memory patterns
- Uses `memory.stat` and `memory.current` for profiling
- Graceful resizing without needing container restarts


### 📦 Image Slimming (`shrink_image`)
Once the container is running, Lit starts a background slimming task:

- Merges intermediate layers
- Deletes cache files, package managers’ leftovers (e.g., `apt`, `npm`, `pip`)
- Compresses the final filesystem layout

All this happens **non-blockingly** — app starts immediately while slimming happens in parallel.

---

## 🧪 Use Cases

Lit is ideal for:

- 🔹 Developers building **minimalist, high-performance containers**
- 🔹 Teams who want **Docker-like workflows without Docker**
- 🔹 Systems with limited resources (IoT, embedded, edge computing)
- 🔹 Security-focused deployments with rootless containers
- 🔹 CI/CD pipelines that auto-tune containers for speed and size


## 🛣 Roadmap

Here’s what’s planned:

- [ ] Namespace + cgroup-based runtime
- [ ] Config-driven container launcher (`lit.yml`)
- [ ] Layered filesystem with OverlayFS
- [ ] Language-specific optimization modules
- [ ] Dynamic resource tuning engine
- [ ] Built-in network isolation and port mapping
- [ ] Cross-platform support (Linux-first, BSD later)
- [ ] Runtime plugin system for community features

---
