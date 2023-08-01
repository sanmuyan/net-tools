$tools = @("portscan", "tcpping", "speedtests" , "speedtestc")
$build_data = @{
    "linux"   = @{
        "arch_list" = @("amd64", "arm")
        "suffix"    = ""
    }
    "darwin"  = @{
        "arch_list" = @("amd64", "arm64")
        "suffix"    = ""
    }
    "windows" = @{
        "arch_list" = @("amd64")
        "suffix"    = ".exe"
    }
}

foreach ($os in $build_data.Keys) {
    foreach ($arch in $build_data[$os].arch_list) {
        $env:GOOS=$os
        $env:GOARCH=$arch
        foreach ($tool in $tools) {
            $suffix = $build_data[$os].suffix
            go build -o ./net-tools/$os/$arch/$tool$suffix ../cmd/$tool
        }
    }
}

Compress-Archive -Path ./net-tools  -DestinationPath pkg/net-tools.zip -Force