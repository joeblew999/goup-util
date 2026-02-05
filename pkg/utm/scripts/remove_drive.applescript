-- remove_drive.applescript
-- This script removes a drive from a specified UTM virtual machine with the given drive ID.
-- Usage: osascript remove_drive.applescript <VM_UUID> <DRIVE_ID>
-- Example: osascript remove_drive.applescript A1B2C3 7FB247A3-DC9F-4A61-A123-0AEE1BEEC636

on run argv
  set vmId to item 1 of argv # UUID of the VM
  set driveId to item 2 of argv # ID of the drive to remove

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing drives
    set vmDrives to drives of config

    -- Find and remove the drive with the given ID
    set updatedDrives to {}
    repeat with drive in vmDrives
      if id of drive is not driveId then
        set end of updatedDrives to drive
      end if
    end repeat

    -- Set the updated drives list
    set drives of config to updatedDrives

    -- Save the configuration (VM must be stopped)
    update configuration of vm with config
  end tell
end run