# torshell

lightweight container runtime for Tor-enforced anonymous shells on Linux.

## what it does

run `torshell` and get an isolated shell where **all network traffic is forced through Tor**.

## how it works

- **5 Linux namespaces** (NET, MNT, PID, IPC, UTS) for container-grade isolation
- **C tor daemon** with TransPort + DNSPort for transparent proxying via iptables DNAT
- **veth pair** connecting isolated namespace to tor on host
- **Zero-leak firewall** UDP/ICMP rejected, DNS forced through Tor, no IPv6

## requirements

- Linux kernel ≥ 5.4
- `tor` binary installed
- Root privileges (orchestrator only — shell runs as your user)
- Go 1.26+

## Project Structure

```
cmd/torshell/       -> entry point
internal/
  ├── firewall/     -> iptables/nftables management
  ├── mount/        -> mount namespace (resolv.conf, hosts, tmpfs)
  ├── netns/        -> network namespace (veth, routing)
  ├── session/      -> lifecycle orchestration
  ├── shell/        -> PTY shell management
  ├── state/        -> subnet registry, locks, sweeper
  └── tor/          -> tor daemon management
pkg/
  ├── logger/       -> structured logging
  └── sysutil/      -> system utilities
```

## License

TBD (gnu mereun)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
