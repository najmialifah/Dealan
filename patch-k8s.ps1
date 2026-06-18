$mapping = @{
  "auth" = 3001; "user" = 3002; "driver" = 3003; "order" = 3004; "matching" = 3005; "pricing" = 3006; "shipment" = 3007; "location" = 3008; "map-route" = 3009; "payment" = 3010; "notification" = 3011; "rating-review" = 3012; "punishment" = 3013; "promo" = 3014; "chat" = 3015
}

foreach ($key in $mapping.Keys) {
    $port = $mapping[$key]
    $file = (Get-ChildItem -Path "$key-service\*-k8s.yaml").FullName
    if ($file) {
        $content = Get-Content $file -Raw
        
        # Replace image
        $content = $content -replace "kelompokdealan/[a-zA-Z0-9_-]+:latest", "dealanregistry.azurecr.io/dealan-$key-service:latest"
        
        # Replace containerPort
        $content = $content -replace "containerPort:\s*\d+", "containerPort: $port"
        
        # Replace targetPort
        $content = $content -replace "targetPort:\s*\d+", "targetPort: $port"
        
        # Inject PORT env var before envFrom if not already there
        # Since envFrom is indented properly, we can just replace envFrom with env+envFrom
        if ($content -notmatch "name:\s*PORT") {
            # Find the indentation of envFrom:
            $content = [regex]::Replace($content, '(?m)^([ \t]*)envFrom:', "`${1}env:`r`n`${1}- name: PORT`r`n`${1}  value: ""$port""`r`n`${1}envFrom:")
        }
        
        Set-Content -Path $file -Value $content
        Write-Host "Updated $file with port $port and image dealanregistry.azurecr.io/dealan-$key-service:latest"
    }
}
