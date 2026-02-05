on run argv
  set vmId to item 1 of argv # Id of the VM
  tell application "UTM"
    set vm to virtual machine id vmId
    set config to configuration of vm

    -- Initialize the network interfaces list to empty
    set updatedNetworkInterfaces to {}

    -- Update the config with the empty network interfaces
    set network interfaces of config to updatedNetworkInterfaces

    -- Update the VM configuration with the new network interface
    update configuration of vm with config
  end tell
end run