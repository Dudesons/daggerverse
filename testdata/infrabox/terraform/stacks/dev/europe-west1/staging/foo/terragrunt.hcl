terraform {
  source = "/terraform/modules/random_id"
}

dependencies {
  paths = [
    "/terraform/stacks/dev/europe-west1/staging/foo",
  ]
}

inputs = {
  length        = 10
}
