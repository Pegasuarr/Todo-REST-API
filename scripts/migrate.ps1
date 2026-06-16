Get-Content .env | ForEach-Object {
    if ($_ -match '(^[^#][^=]+)=(.+)$')
    {
        Set-Item -Path "env:$($matches[1])" -Value $matches[2]
    }
}

$command = $args[0]
$name = $args[1]

switch($command){
    "up" { migrate -path migrations -database $env:DB_URL up }
    "down" {
        $count = if ($name) { $name } else { "1" }

        Write-Host "Rolling back $count migration(s). Continue? [y/n]"
        $confirm = Read-Host
        if ($confirm -eq 'y')
        {
            migrate -path migrations -database $env:DB_URL down $count
        }
    }
    "create" { migrate create -ext sql -dir migrations -seq $name }
}