on run argv
    tell application "UTM"
        set vmId to item 1 of argv -- VM id is given as the first argument
        set vmName to ""
        set cpuCount to 0
        set memorySize to 0
        set vmNotes to ""
        set useHypervisor to null
        set uefiBoot to null
        set directoryShareMode to null


        -- Parse arguments
        repeat with i from 2 to (count argv)
            set currentArg to item i of argv
            if currentArg is "--name" then
                set vmName to item (i + 1) of argv
            else if currentArg is "--cpus" then
                set cpuCount to item (i + 1) of argv
            else if currentArg is "--memory" then
                set memorySize to item (i + 1) of argv
            else if currentArg is "--notes" then
                set vmNotes to item (i + 1) of argv
            else if currentArg is "--use-hypervisor" then
                set hypervisorArg to item (i + 1) of argv
                if hypervisorArg is "true" then
                    set useHypervisor to true
                else if hypervisorArg is "false" then
                    set useHypervisor to false
                end if
            else if currentArg is "--uefi-boot" then
                set uefiBootArg to item (i + 1) of argv
                if uefiBootArg is "true" then
                    set uefiBoot to true
                else if uefiBootArg is "false" then
                    set uefiBoot to false
                end if
            else if currentArg is "--directory-share-mode" then
                set directoryShareMode to item (i + 1) of argv
            end if
        end repeat
        
        -- Get the VM and its configuration
        set vm to virtual machine id vmId -- ID is assumed to be valid
        set config to configuration of vm
        
        -- Set VM name if provided
        if vmName is not "" then
            set name of config to vmName
        end if
        
        -- Set CPU count if provided
        if cpuCount is not 0 then
            set cpu cores of config to cpuCount
        end if
        
        -- Set memory size if provided
        if memorySize is not 0 then
            set memory of config to memorySize
        end if
        
        -- Set the notes if --notes is provided (existing notes will be overwritten)
        if vmNotes is not "" then
            set notes of config to vmNotes
        end if

        -- Set Use Hypervisor if provided
        if useHypervisor is not null then
            set hypervisor of config to useHypervisor
        end if

        -- Set UEFI boot if provided
        if uefiBoot is not null then
            set uefi of config to uefiBoot
        end if

        -- Set Directory Sharing mode if provided
        if directoryShareMode is not null then
            set directory share mode of config to directoryShareMode -- mode is assumed to be enum value
        end if

        -- Save the configuration
        update configuration of vm with config
 
    end tell
end run