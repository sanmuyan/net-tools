$commands = @{
    "nettools" = "nts"
}
$build_data = @{
    "linux"   = @{
        "arch_list" = @("amd64", "arm64")
        "suffix"    = ""
    }
    "darwin"  = @{
        "arch_list" = @("amd64", "arm64")
        "suffix"    = ""
    }
    "windows" = @{
        "arch_list" = @("amd64", "arm64")
        "suffix"    = ".exe"
    }
}

Remove-Item .\nts\* -Recurse
Remove-Item .\pkg\* -Recurse

foreach ($os in $build_data.Keys) {
    foreach ($arch in $build_data[$os].arch_list) {
        $env:GOOS=$os
        $env:GOARCH=$arch
        foreach ($command in $commands.Keys) {
            $suffix = $build_data[$os].suffix
            $path = ".\nts\$os\$arch"
            $command_name = $commands[$command]
            $bin = "$path\$command_name$suffix"
            $upx_bin = "$path\$command_name$suffix"

            go build -ldflags "-w -s" -o  $bin ../cmd/$command
            #upx -1 -o $upx_bin $bin
            #Remove-Item $bin
        }
    }
}

Compress-Archive -Path ./nts  -DestinationPath pkg/nts.zip -Force