-- This script removes the first drive from a specified UTM virtual machine.
-- Usage: osascript remove_first_drive.applescript <VM_UUID>
-- Example: osascript remove_first_drive.applescript A123

on run argv
  set vmId to item 1 of argv -- UUID of the VM

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Get the current drives
    set currentDrives to drives of config

    -- Initialize a new list for the updated drives
    set updatedDrives to {}

    -- Iterate through the current drives and add all except the first one
    repeat with i from 2 to (count of currentDrives)
      set end of updatedDrives to item i of currentDrives
    end repeat

    -- Update the configuration with the new drives list
    set drives of config to updatedDrives
    update configuration of vm with config
  end tell
end run