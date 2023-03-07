docker_prune_settings(num_builds=5)

docker_build(
    'localhost:5000/diamond:dev', '',
    dockerfile='Dockerfile',
    only=[
        './go.mod',
        './go.sum',
        './main.go',
     ],
)

k8s_yaml('manifest.yaml')
k8s_resource('diamond', port_forwards=[8080])
