foreach ($dir in Get-ChildItem -Directory) {
    if ($dir.Name -match "-service$") {
        $file = (Get-ChildItem -Path "$($dir.FullName)\*-k8s.yaml").FullName
        if ($file) {
            $content = Get-Content $file -Raw
            $content = $content -replace '(?m)^([ \t]*)ports:\s*\r?\n([ \t]*)-\s*containerPort:', '$1envFrom:
$1- configMapRef:
$1    name: dealan-config
$1ports:
$2- containerPort:'
            Set-Content -Path $file -Value $content
        }
    }
}
