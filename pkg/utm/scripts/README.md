# UTM AppleScript Files

These AppleScript files provide automation for UTM virtual machine management.

## Source

Copied from: https://github.com/naveenrajm7/packer-plugin-utm

Original location: `builder/utm/common/scripts/`

## Scripts

| Script | Purpose |
|--------|---------|
| `create_vm.applescript` | Create a new VM with specified backend and architecture |
| `customize_vm.applescript` | Set CPU, RAM, UEFI boot, and other VM settings |
| `add_drive.applescript` | Add a disk drive to a VM |
| `attach_iso.applescript` | Attach an ISO file to a VM |
| `add_network_interface.applescript` | Add a network interface (shared, emulated, bridged) |
| `add_port_forwards.applescript` | Add port forwarding rules |
| `clear_port_forwards.applescript` | Remove all port forwarding rules |
| `clear_network_interfaces.applescript` | Remove all network interfaces |
| `add_qemu_display.applescript` | Add a QEMU display (VNC, SPICE, etc.) |
| `add_qemu_additional_args.applescript` | Add custom QEMU arguments |
| `remove_drive.applescript` | Remove a drive from a VM |
| `remove_first_drive.applescript` | Remove the first drive from a VM |
| `remove_qemu_display_by_name.applescript` | Remove a display by name |
| `remove_qemu_additional_args.applescript` | Remove custom QEMU arguments |

## UTM AppleScript Enum Codes

UTM uses 4-character enum codes in AppleScript:

### Backend
- `QeMu` - QEMU emulation
- `ApLe` - Apple Virtualization framework

### Drive Interface
- `QdIv` - virtio
- `QdIu` - USB
- `QdIi` - IDE
- `QdIs` - SCSI
- `QdIn` - NVMe

### Network Mode
- `ShRd` - Shared Network (NAT)
- `EmUd` - Emulated VLAN (for port forwarding)
- `BrDg` - Bridged
- `HsOn` - Host Only

### Protocol
- `TcPp` - TCP
- `UdPp` - UDP

## License

Original code is licensed under MPL-2.0 (HashiCorp).
