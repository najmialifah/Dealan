import os
import re

for root, dirs, files in os.walk('.'):
    for dir in dirs:
        if dir.endswith('-service'):
            for file in os.listdir(dir):
                if file.endswith('-k8s.yaml'):
                    filepath = os.path.join(dir, file)
                    with open(filepath, 'r') as f:
                        content = f.read()
                    
                    # Regex to find "ports:" and its leading spaces
                    # and insert "envFrom:" block right before it with the SAME leading spaces
                    new_content = re.sub(
                        r'^([ \t]*)ports:\s*^([ \t]*)-\s*containerPort:', 
                        r'\1envFrom:\n\1- configMapRef:\n\1    name: dealan-config\n\1ports:\n\2- containerPort:', 
                        content, 
                        flags=re.MULTILINE
                    )
                    
                    with open(filepath, 'w') as f:
                        f.write(new_content)
                    print(f"Patched {filepath}")
