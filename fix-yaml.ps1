foreach ($dir in Get-ChildItem -Directory) {
    if ($dir.Name -match "-service$") {
        $file = (Get-ChildItem -Path "$($dir.FullName)\*-k8s.yaml").FullName
        if ($file) {
            $content = Get-Content $file -Raw
            $pattern = '(?m)^([ \t]*)ports:\s*^([ \t]*)-\s*containerPort:'
            $replacement = "${1}envFrom:
${1}- configMapRef:
${1}    name: dealan-config
${1}ports:
${2}- containerPort:"
            $content = [regex]::Replace($content, $pattern, $replacement)
            Set-Content -Path $file -Value $content
            Write-Host "Patched $file"
        }
    }
}
