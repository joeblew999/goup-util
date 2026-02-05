---
-- attach_iso.applescript
-- This script attaches an removable drive to a specified virtual machine with given source file.
-- Usage: osascript attach_iso.applescript <VM_ID> --interface <INT> --source <ISO_PATH>
-- Example: osascript attach_iso.applescript test --interface "QdIu" --source "full/path/to/my.iso"
-- add a drive with given interface and source file "full/path/to/my.iso"
on run argv
  set vmId to item 1 of argv # ID of the VM
  -- Parse the --interface argument
  set isoInterface to item 3 of argv
  -- Parse the --source argument
  set isoPath to item 5 of argv as string

  -- Default removable to true
  set removableVal to true

  -- Parse the --removable argument if provided
  repeat with i from 6 to (count argv)
    set currentArg to item i of argv
    if currentArg is "--removable" then
      set removableArg to item (i + 1) of argv
      if removableArg is "false" then
        set removableVal to false
      end if
    end if
  end repeat

  -- gain access to the file, so you can pass it to UTM (which is sandboxed)
  -- isoPath is already a POSIX path string, so we can use it directly
  set isoFile to POSIX file isoPath

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing drives
    set vmDrives to drives of config
    --- create a new drive
    set newDrive to {removable: removableVal, interface: isoInterface, source:isoFile}
    -- Add the new drive to the beginning of the list
    -- set vmDrives to {newDrive} & vmDrives

    -- Add the new drive to the end of the list
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
