# Lit — A Mini Kernel in TinyGo

A bare-metal mini kernel written in [TinyGo](https://tinygo.org/), built from scratch to understand how operating systems work at the lowest level. This project covers everything from the moment the CPU powers on to running a first process.

> This is a learning project.

---

##  The Idea

Most people use an OS every day without knowing what happens in the first milliseconds after pressing the power button. This project explores exactly that.

When a machine powers on:
1. The CPU jumps to a hardcoded address on the motherboard ROM — the **BIOS/UEFI firmware**
2. The firmware runs **POST** (Power-On Self Test) and initializes hardware
3. The firmware finds a **bootloader** on disk (we use GRUB)
4. The bootloader loads **our kernel binary** into RAM and jumps to it
5. **Our code runs** — with no OS underneath, no stdlib, nothing

TinyGo makes this possible in Go by replacing the standard Go runtime with a minimal one that can run without an OS.

---

## What This Kernel Can Do (Minimum Viable Kernel)

| Component | Description |
|---|---|
| **Boot & Init** | Accept control from GRUB, set up the stack, initialize BSS |
| **UART Output** | Write to serial port so we can see what's happening |
| **Interrupt Handling** | Set up IDT, handle CPU exceptions gracefully |
| **Memory Management** | Detect RAM, basic allocator (alloc/free) |
| **First Process** | Run a single function as PID 1 |

If the kernel can do all of the above, it will produce output like this on boot:

```
[BOOT]  Kernel loaded at 0x100000
[MEM]   Detected 128MB RAM
[IDT]   Interrupts initialized
[UART]  Serial output ready
[INIT]  Jumping to first process...
[PROC]  Hello from PID 1!
```

---


##  First Steps

### 1. Install TinyGo

```bash
# macOS
brew install tinygo

# Linux — download from the releases page
wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.0/tinygo0.31.0.linux-amd64.tar.gz
tar -xf tinygo0.31.0.linux-amd64.tar.gz
export PATH=$PATH:$(pwd)/tinygo/bin
```

Verify it works:
```bash
tinygo version
```

### 2. Install QEMU (for testing without real hardware)

```bash
# macOS
brew install qemu

# Ubuntu/Debian
sudo apt install qemu-system-x86
```

## Testing Strategy

One of the core ideas of this project is that as much kernel logic as possible should be testable **without booting anything**.

| Layer | How it's tested |
|---|---|
| Pure logic (memory allocator, IDT encoding) | `go test` directly |
| Hardware output (UART) | Build tag swaps real hardware for a fake buffer |
| Full boot sequence | QEMU scripted with `test_boot.sh` |

The rule: **if it can be a unit test, it should be a unit test.** Only use QEMU for things that genuinely need real hardware behavior.

---

##  Build Order

Build the kernel component by component in this order — each step proves the previous one works:

```
Step 1 — UART output          "Hello from kernel!"
       ↓
Step 2 — IDT + exceptions     survive a crash gracefully
       ↓
Step 3 — Memory detection     know your RAM
       ↓
Step 4 — Memory allocator     alloc/free
       ↓
Step 5 — First process        jump to a function, call it PID 1
       ↓
Step 6 — Timer interrupt      proves hardware interrupts work
```

---

## 🔗 Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [OSDev Wiki](https://wiki.osdev.org/) — the bible for kernel development
- [Writing an OS in Rust](https://os.phil-opp.com/) — great reference even for Go
- [GRUB Multiboot Spec](https://www.gnu.org/software/grub/manual/multiboot/multiboot.html)
- [x86 Interrupt Reference](https://wiki.osdev.org/Interrupts)
