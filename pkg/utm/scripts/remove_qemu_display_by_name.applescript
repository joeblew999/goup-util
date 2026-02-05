-- remove_qemu_display.applescript
-- This script removes a display from a specified UTM virtual machine based on the given hardware type.
-- Usage: osascript remove_qemu_display.applescript <VM_UUID> --hardware <HARDWARE>
-- Example: osascript remove_qemu_display.applescript A1B2C3 --hardware "pci"

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

    -- Find and remove the display with the given hardware type
    set updatedDisplays to {}
    repeat with display in vmDisplays
      if hardware of display is not hardwareType then
        set end of updatedDisplays to display
      end if
    end repeat

    -- Set the updated displays list
    set displays of config to updatedDisplays

    -- Save the configuration (VM must be stopped)
    update configuration of vm with config
  end tell
end run