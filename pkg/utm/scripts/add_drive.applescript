---
-- add_drive.applescript
-- This script adds a drive to a specified UTM virtual machine with given interface and size.
-- Usage: osascript add_drive.applescript <VM_UUID> --interface <INTERFACE> --size <SIZE>
-- Example: osascript add_drive.applescript A1B2C3  --interface "QdIu" --size 65536
-- creates driver with USB interface and size 65536

on run argv
  set vmId to item 1 of argv # UUID of the VM
  -- Parse the --interface argument
  set driveInterface to item 3 of argv
  -- Parse the --size argument
  set driveSize to item 5 of argv 

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing drives
    set vmDrives to drives of config
    --- create a new drive
    set newDrive to {interface: driveInterface, guest size: driveSize}
    --- add the drive to the end of the list
    copy newDrive to end of vmDrives
    --- set drives with new drive list
    set drives of config to vmDrives

    --- save the configuration (VM must be stopped)
    update configuration of vm with config
    
    -- Get the updated drive id
    set updatedConfig to configuration of vm
    set updatedDrives to drives of updatedConfig
    set updatedDrive to item -1 of updatedDrives
    set updatedDriveId to id of updatedDrive

    -- return the new drive id
    return updatedDriveId
  end tell
end run