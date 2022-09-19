load('ext://restart_process', 'docker_build_with_restart')

compile_cmd = 'make compile'

local_resource(
  'build-api',
  compile_cmd,
  deps=['.'], ignore=['**', '!*.go']
)

docker_build_with_restart(
  'tschwaa/api',
  '.',
  dockerfile='Dockerfile.dev',
  entrypoint=['/src/build/main', '--'],
  only=[
    './build',
  ],
  live_update=[
    sync('./build', '/src/build')
  ]
)
