---
-- add_qemu_display.applescript
-- This script adds a display to a specified UTM virtual machine with the given hardware type.
-- Usage: osascript add_qemu_display.applescript <VM_UUID> --hardware <HARDWARE>
-- Example: osascript add_qemu_display.applescript A1B2C3 --hardware "virtio-ramfb"
-- Adds a display with the specified hardware type.
-- Returns the hardware type for tracking (compatible with all UTM versions).

on run argv
  set vmId to item 1 of argv # UUID of the VM
  -- Parse the --hardware argument
  set hardwareType to item 3 of argv

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing displays
    set vmDisplays to displays of config
    --- create a new display
    set newDisplay to {hardware: hardwareType}
    --- add the display to the end of the list
    copy newDisplay to end of vmDisplays
    --- set displays with new display list
    set displays of config to vmDisplays

    --- save the configuration (VM must be stopped)
    update configuration of vm with config

    -- Return the hardware type for tracking (compatible with all UTM versions)
    -- We don't access id or index properties to maintain compatibility
    return hardwareType
  end tell
end run